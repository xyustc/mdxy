package service

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/xyu/mdxy/internal/config"
)

type NoteService struct {
	metaMu sync.Mutex
}

func NewNoteService() *NoteService {
	return &NoteService{}
}

// metaFilePath 返回 meta.json 的绝对路径（NotesDir 同级）
func (s *NoteService) metaFilePath() string {
	return filepath.Join(filepath.Dir(config.AppConfig.Content.NotesDir), "notes-meta.json")
}

type notesMeta struct {
	Featured []string `json:"featured"`
}

// loadMeta 读取 meta.json，容错：不存在或解析失败均返回空结构
func (s *NoteService) loadMeta() notesMeta {
	data, err := os.ReadFile(s.metaFilePath())
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("notes-meta.json 读取失败: %v", err)
		}
		return notesMeta{}
	}
	var m notesMeta
	if err := json.Unmarshal(data, &m); err != nil {
		log.Printf("notes-meta.json 解析失败，已忽略: %v", err)
		return notesMeta{}
	}
	if m.Featured == nil {
		m.Featured = []string{}
	}
	return m
}

// saveMeta 原子写入 meta.json（写临时文件再 rename）
func (s *NoteService) saveMeta(m notesMeta) error {
	if m.Featured == nil {
		m.Featured = []string{}
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	target := s.metaFilePath()
	tmp := target + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	if err := os.Rename(tmp, target); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}

// SetFeatured 设置笔记的 featured 状态
func (s *NoteService) SetFeatured(notePath string, featured bool) error {
	s.metaMu.Lock()
	defer s.metaMu.Unlock()

	m := s.loadMeta()
	set := make(map[string]bool, len(m.Featured))
	for _, p := range m.Featured {
		set[p] = true
	}

	if featured {
		set[notePath] = true
	} else {
		delete(set, notePath)
	}

	m.Featured = make([]string, 0, len(set))
	for p := range set {
		m.Featured = append(m.Featured, p)
	}
	return s.saveMeta(m)
}

type NoteNode struct {
	Name       string     `json:"name"`
	Type       string     `json:"type"` // file/directory
	Path       string     `json:"path"`
	Featured   bool       `json:"featured,omitempty"`
	Children   []NoteNode `json:"children,omitempty"`
}

// GetTree 获取笔记目录树
func (s *NoteService) GetTree() ([]NoteNode, error) {
	notesDir := config.AppConfig.Content.NotesDir
	s.metaMu.Lock()
	m := s.loadMeta()
	s.metaMu.Unlock()
	featured := make(map[string]bool, len(m.Featured))
	for _, p := range m.Featured {
		featured[p] = true
	}
	return s.scanDirectory(notesDir, notesDir, featured)
}

func (s *NoteService) scanDirectory(dir, baseDir string, featured map[string]bool) ([]NoteNode, error) {
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
			children, err := s.scanDirectory(fullPath, baseDir, featured)
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
				Name:     name,
				Type:     "file",
				Path:     relPath,
				Featured: featured[relPath],
			})
		}
	}

	return nodes, nil
}

// validatePath 校验路径安全性，返回绝对路径
func (s *NoteService) validatePath(notePath string) (string, error) {
	notePath = filepath.Clean(notePath)
	if strings.Contains(notePath, "..") || filepath.IsAbs(notePath) {
		return "", os.ErrPermission
	}

	fullPath := filepath.Join(config.AppConfig.Content.NotesDir, notePath)

	absNotesDir, _ := filepath.Abs(config.AppConfig.Content.NotesDir)
	absFullPath, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absFullPath, absNotesDir+string(filepath.Separator)) {
		return "", os.ErrPermission
	}

	return fullPath, nil
}

// validateNotePath 在 validatePath 基础上额外要求 .md 后缀
func (s *NoteService) validateNotePath(notePath string) (string, error) {
	if !strings.HasSuffix(notePath, ".md") {
		return "", os.ErrPermission
	}
	return s.validatePath(notePath)
}

// GetContent 获取笔记内容
func (s *NoteService) GetContent(notePath, clientIP string) (string, error) {
	fullPath, err := s.validateNotePath(notePath)
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// GetRawContent 获取笔记原始内容（不注入水印）
func (s *NoteService) GetRawContent(notePath string) (string, error) {
	fullPath, err := s.validateNotePath(notePath)
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// SaveNote 写入笔记内容
func (s *NoteService) SaveNote(notePath, content string) error {
	fullPath, err := s.validateNotePath(notePath)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(fullPath, []byte(content), 0644)
}

// DeleteNote 删除笔记文件
func (s *NoteService) DeleteNote(notePath string) error {
	fullPath, err := s.validateNotePath(notePath)
	if err != nil {
		return err
	}

	return os.Remove(fullPath)
}

// CreateDirectory 创建目录
func (s *NoteService) CreateDirectory(dirPath string) error {
	fullPath, err := s.validatePath(dirPath)
	if err != nil {
		return err
	}

	return os.MkdirAll(fullPath, 0755)
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

// parseFeaturedFrontmatter 已废弃，featured 状态改由 notes-meta.json 管理

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
