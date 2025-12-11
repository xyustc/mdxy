package services

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"mdxy-backend/config"
)

// TreeNode 目录树节点
type TreeNode struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Path     string      `json:"path"`
	Children []*TreeNode `json:"children,omitempty"`
}

// SearchResult 搜索结果
type SearchResult struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Context string `json:"context"`
}

// GetNoteTree 获取笔记目录树
func GetNoteTree() []*TreeNode {
	return scanDirectory(config.NotesDir, config.NotesDir)
}

// scanDirectory 递归扫描目录
func scanDirectory(directory string, relativeBase string) []*TreeNode {
	items := []*TreeNode{}

	// 检查目录是否存在
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return items
	}

	// 读取目录内容
	entries, err := ioutil.ReadDir(directory)
	if err != nil {
		return items
	}

	// 对条目进行排序（文件夹优先，然后按名称排序）
	// 这里简化处理，实际可以更复杂

	for _, entry := range entries {
		// 跳过隐藏文件和目录
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		fullPath := filepath.Join(directory, entry.Name())
		relativePath, _ := filepath.Rel(relativeBase, fullPath)

		if entry.IsDir() {
			children := scanDirectory(fullPath, relativeBase)
			// 只添加非空文件夹或包含md文件的文件夹
			if len(children) > 0 {
				items = append(items, &TreeNode{
					Name:     entry.Name(),
					Type:     "directory",
					Path:     strings.ReplaceAll(relativePath, "\\", "/"),
					Children: children,
				})
			}
		} else {
			// 检查文件扩展名是否允许
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if config.AllowedExtensions[ext] {
				name := strings.TrimSuffix(entry.Name(), ext)
				items = append(items, &TreeNode{
					Name: name,
					Type: "file",
					Path: strings.ReplaceAll(relativePath, "\\", "/"),
				})
			}
		}
	}

	return items
}

// GetNoteContent 获取笔记内容
func GetNoteContent(notePath string) string {
	// 安全检查：防止目录穿越攻击
	safePath := SafePath(notePath)
	if safePath == "" {
		return ""
	}

	// 检查文件是否存在
	if _, err := os.Stat(safePath); os.IsNotExist(err) {
		return ""
	}

	// 检查文件扩展名是否允许
	ext := strings.ToLower(filepath.Ext(safePath))
	if !config.AllowedExtensions[ext] {
		return ""
	}

	// 读取文件内容
	content, err := ioutil.ReadFile(safePath)
	if err != nil {
		return ""
	}

	return string(content)
}

// SafePath 验证并返回安全的文件路径，防止目录穿越攻击
func SafePath(notePath string) string {
	// 规范化路径
	normalized := filepath.Clean(notePath)
	// 移除开头的斜杠
	normalized = strings.TrimPrefix(normalized, "/")

	// 构建完整路径
	fullPath := filepath.Join(config.NotesDir, normalized)

	// 确保路径在 NOTES_DIR 内
	rel, err := filepath.Rel(config.NotesDir, fullPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return ""
	}

	return fullPath
}

// SearchNotes 搜索笔记内容
func SearchNotes(keyword string) []SearchResult {
	results := []SearchResult{}
	keywordLower := strings.ToLower(keyword)

	// 遍历 NOTES_DIR 目录下的所有 .md 文件
	err := filepath.Walk(config.NotesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// 只处理文件
		if !info.IsDir() {
			// 检查文件扩展名
			ext := strings.ToLower(filepath.Ext(path))
			if config.AllowedExtensions[ext] {
				// 读取文件内容
				content, readErr := ioutil.ReadFile(path)
				if readErr != nil {
					return nil
				}

				contentStr := string(content)
				contentLower := strings.ToLower(contentStr)

				// 检查是否包含关键词
				if strings.Contains(contentLower, keywordLower) || strings.Contains(strings.ToLower(info.Name()), keywordLower) {
					relativePath, _ := filepath.Rel(config.NotesDir, path)
					// 提取匹配的上下文
					context := extractContext(contentStr, keyword)
					results = append(results, SearchResult{
						Name:    strings.TrimSuffix(info.Name(), ext),
						Path:    strings.ReplaceAll(relativePath, "\\", "/"),
						Context: context,
					})
				}
			}
		}

		return nil
	})

	if err != nil {
		return results
	}

	return results
}

// extractContext 提取关键词周围的上下文
func extractContext(content string, keyword string) string {
	contextLength := 100
	keywordLower := strings.ToLower(keyword)
	contentLower := strings.ToLower(content)

	pos := strings.Index(contentLower, keywordLower)
	if pos == -1 {
		// 如果内容中没找到，返回开头部分
		if len(content) > contextLength {
			return content[:contextLength] + "..."
		}
		return content
	}

	start := pos - contextLength/2
	if start < 0 {
		start = 0
	}

	end := pos + len(keyword) + contextLength/2
	if end > len(content) {
		end = len(content)
	}

	context := content[start:end]
	if start > 0 {
		context = "..." + context
	}
	if end < len(content) {
		context = context + "..."
	}

	return context
}
