package service

import (
	"fmt"
	"html"
	"strings"

	"gorm.io/gorm"

	"github.com/xyu/mdxy/internal/model"
	"github.com/xyu/mdxy/internal/repository"
)

func escapeLike(s string) string {
	r := strings.NewReplacer("%", "\\%", "_", "\\_")
	return r.Replace(s)
}

type SearchService struct {
	noteService *NoteService
	toolRepo    *repository.ToolRepository
	db          *gorm.DB
}

func NewSearchService(noteService *NoteService, toolRepo *repository.ToolRepository, db *gorm.DB) *SearchService {
	return &SearchService{noteService: noteService, toolRepo: toolRepo, db: db}
}

type SearchResultItem struct {
	Type    string `json:"type"` // note / tool
	Title   string `json:"title"`
	Path    string `json:"path"`
	Context string `json:"context"`
	URL     string `json:"url,omitempty"`
}

type SearchGroup struct {
	Type    string             `json:"type"`
	Label   string             `json:"label"`
	Results []SearchResultItem `json:"results"`
}

type UnifiedSearchResult struct {
	Groups []SearchGroup `json:"groups"`
	Total  int           `json:"total"`
}

func (s *SearchService) Search(query, searchType, ip string) (UnifiedSearchResult, error) {
	result := UnifiedSearchResult{}

	if searchType == "all" || searchType == "notes" {
		noteResults, err := s.noteService.Search(query)
		if err == nil && len(noteResults) > 0 {
			var items []SearchResultItem
			for _, r := range noteResults {
				items = append(items, SearchResultItem{
					Type:    "note",
					Title:   r.Name,
					Path:    "/notes/" + r.Path,
					Context: highlightText(r.Context, query),
				})
			}
			result.Groups = append(result.Groups, SearchGroup{
				Type: "notes", Label: "笔记", Results: items,
			})
			result.Total += len(items)
		}
	}

	if searchType == "all" || searchType == "tools" {
		var tools []model.Tool
		escaped := "%" + escapeLike(query) + "%"
		s.db.Where("is_visible = ? AND (name LIKE ? ESCAPE '\\' OR description LIKE ? ESCAPE '\\' OR category LIKE ? ESCAPE '\\')",
			true, escaped, escaped, escaped).
			Find(&tools)
		if len(tools) > 0 {
			var items []SearchResultItem
			for _, t := range tools {
				items = append(items, SearchResultItem{
					Type:    "tool",
					Title:   highlightText(t.Name, query),
					Path:    "/tools",
					Context: highlightText(t.Description, query),
					URL:     t.URL,
				})
			}
			result.Groups = append(result.Groups, SearchGroup{
				Type: "tools", Label: "工具", Results: items,
			})
			result.Total += len(items)
		}
	}

	// 异步记录搜索日志
	go func() {
		s.db.Create(&model.SearchLog{
			Query:       query,
			Source:      searchType,
			ResultCount: result.Total,
			IP:          ip,
		})
	}()

	return result, nil
}

func (s *SearchService) PopularKeywords(limit int) ([]map[string]interface{}, error) {
	var results []struct {
		Query string
		Count int64
	}
	err := s.db.Model(&model.SearchLog{}).
		Select("query, COUNT(*) as count").
		Group("query").Order("count DESC").Limit(limit).
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	var keywords []map[string]interface{}
	for _, r := range results {
		keywords = append(keywords, map[string]interface{}{
			"keyword": r.Query,
			"count":   r.Count,
		})
	}
	return keywords, nil
}

func highlightText(text, keyword string) string {
	if keyword == "" || text == "" {
		return html.EscapeString(text)
	}
	lower := strings.ToLower(text)
	lowerKw := strings.ToLower(keyword)
	idx := strings.Index(lower, lowerKw)
	if idx == -1 {
		return html.EscapeString(text)
	}
	return fmt.Sprintf("%s<mark>%s</mark>%s",
		html.EscapeString(text[:idx]),
		html.EscapeString(text[idx:idx+len(keyword)]),
		html.EscapeString(text[idx+len(keyword):]))
}
