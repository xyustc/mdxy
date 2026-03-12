package service

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/xyu/mdxy/internal/config"
)

type NoteService struct{}

func NewNoteService() *NoteService {
	return &NoteService{}
}

type NoteNode struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"` // file/directory
	Path     string      `json:"path"`
	Children []NoteNode  `json:"children,omitempty"`
}

// GetTree 获取笔记目录树
func (s *NoteService) GetTree() ([]NoteNode, error) {
	notesDir := config.AppConfig.Content.NotesDir
	return s.scanDirectory(notesDir, notesDir)
}

func (s *NoteService) scanDirectory(dir, baseDir string) ([]NoteNode, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var nodes []NoteNode

	for _, entry := range entries {
		// 跳过隐藏文件
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		fullPath := filepath.Join(dir, entry.Name())
		relPath, _ := filepath.Rel(baseDir, fullPath)
		relPath = filepath.ToSlash(relPath)

		if entry.IsDir() {
			children, err := s.scanDirectory(fullPath, baseDir)
			if err != nil || len(children) == 0 {
				continue
			}
			nodes = append(nodes, NoteNode{
				Name:     entry.Name(),
				Type:     "directory",
				Path:     relPath,
				Children: children,
			})
		} else if strings.HasSuffix(entry.Name(), ".md") {
			name := strings.TrimSuffix(entry.Name(), ".md")
			nodes = append(nodes, NoteNode{
				Name: name,
				Type: "file",
				Path: relPath,
			})
		}
	}

	return nodes, nil
}

// GetContent 获取笔记内容
func (s *NoteService) GetContent(notePath, clientIP string) (string, error) {
	// 安全检查：防止路径穿越
	notePath = filepath.Clean(notePath)
	if strings.Contains(notePath, "..") || filepath.IsAbs(notePath) {
		return "", os.ErrPermission
	}

	// 只允许读取 .md 文件
	if !strings.HasSuffix(notePath, ".md") {
		return "", os.ErrPermission
	}

	fullPath := filepath.Join(config.AppConfig.Content.NotesDir, notePath)

	// 二次验证：确保最终路径在 NotesDir 内
	absNotesDir, _ := filepath.Abs(config.AppConfig.Content.NotesDir)
	absFullPath, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absFullPath, absNotesDir+string(filepath.Separator)) {
		return "", os.ErrPermission
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

type SearchResult struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Context string `json:"context"`
}

// Search 搜索笔记
func (s *NoteService) Search(keyword string) ([]SearchResult, error) {
	var results []SearchResult
	notesDir := config.AppConfig.Content.NotesDir

	err := filepath.Walk(notesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		contentStr := string(content)
		lowerContent := strings.ToLower(contentStr)
		lowerKeyword := strings.ToLower(keyword)

		// 检查文件名或内容是否包含关键词
		if strings.Contains(strings.ToLower(info.Name()), lowerKeyword) ||
			strings.Contains(lowerContent, lowerKeyword) {

			relPath, _ := filepath.Rel(notesDir, path)
			relPath = filepath.ToSlash(relPath)

			name := strings.TrimSuffix(info.Name(), ".md")
			context := extractContext(contentStr, keyword, 100)

			results = append(results, SearchResult{
				Name:    name,
				Path:    relPath,
				Context: context,
			})
		}

		return nil
	})

	return results, err
}

func extractContext(content, keyword string, length int) string {
	lowerContent := strings.ToLower(content)
	lowerKeyword := strings.ToLower(keyword)

	pos := strings.Index(lowerContent, lowerKeyword)
	if pos == -1 {
		if len(content) > length {
			return content[:length] + "..."
		}
		return content
	}

	start := pos - length/2
	if start < 0 {
		start = 0
	}

	end := pos + len(keyword) + length/2
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
