package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gorm.io/gorm"

	"github.com/xyu/mdxy/internal/config"
	"github.com/xyu/mdxy/internal/model"
	"github.com/xyu/mdxy/internal/repository"
)

const ClaudeCodeCheatSheetSlug = "claude-code"

type CheatSheetCodeItem struct {
	Code string `json:"code"`
	Text string `json:"text,omitempty"`
}

type CheatSheetFooterRow struct {
	Label string               `json:"label"`
	Items []CheatSheetCodeItem `json:"items"`
}

type CheatSheetRow struct {
	Key        string `json:"key"`
	Desc       string `json:"desc,omitempty"`
	AddedAt    string `json:"added_at,omitempty"`
	KeyVariant string `json:"key_variant,omitempty"`
}

type CheatSheetGroup struct {
	Title string          `json:"title"`
	Rows  []CheatSheetRow `json:"rows"`
}

type CheatSheetSection struct {
	ID     string            `json:"id"`
	Title  string            `json:"title"`
	Theme  string            `json:"theme"`
	Groups []CheatSheetGroup `json:"groups"`
}

type CheatSheetMeta struct {
	SourceURL   string `json:"source_url"`
	SourceTitle string `json:"source_title"`
	Version     string `json:"version"`
	UpdatedAt   string `json:"updated_at"`
}

type CheatSheetContent struct {
	Meta      CheatSheetMeta         `json:"meta"`
	Changelog []CheatSheetCodeItem   `json:"changelog"`
	Footer    []CheatSheetFooterRow  `json:"footer"`
	Columns   [][]CheatSheetSection  `json:"columns"`
}

type CheatSheetPublishedPayload struct {
	Slug        string            `json:"slug"`
	Status      string            `json:"status"`
	SourceURL   string            `json:"source_url"`
	SyncedAt    time.Time         `json:"synced_at"`
	PublishedAt *time.Time        `json:"published_at,omitempty"`
	Content     CheatSheetContent `json:"content"`
}

type CheatSheetSnapshotSummary struct {
	ID                  uint       `json:"id"`
	Slug                string     `json:"slug"`
	SourceVersion       string     `json:"source_version"`
	SourceUpdatedAtText string     `json:"source_updated_at_text"`
	Status              string     `json:"status"`
	SyncError           string     `json:"sync_error,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	PublishedAt         *time.Time `json:"published_at,omitempty"`
}

type CheatSheetSyncResult struct {
	Changed       bool                     `json:"changed"`
	AutoPublished bool                     `json:"auto_published"`
	Snapshot      *CheatSheetSnapshotSummary `json:"snapshot,omitempty"`
}

type CheatSheetService struct {
	repo          *repository.CheatSheetRepository
	httpClient    *http.Client
	schedulerOnce sync.Once
}

func NewCheatSheetService(repo *repository.CheatSheetRepository) *CheatSheetService {
	return &CheatSheetService{
		repo: repo,
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

func (s *CheatSheetService) StartScheduler() {
	if !config.AppConfig.CheatSheetSync.Enabled {
		return
	}

	s.schedulerOnce.Do(func() {
		interval := time.Duration(config.AppConfig.CheatSheetSync.IntervalHours) * time.Hour
		if interval <= 0 {
			interval = 24 * time.Hour
		}

		go func() {
			s.runScheduledSync()

			ticker := time.NewTicker(interval)
			defer ticker.Stop()

			for range ticker.C {
				s.runScheduledSync()
			}
		}()
	})
}

func (s *CheatSheetService) runScheduledSync() {
	if _, err := s.Sync(ClaudeCodeCheatSheetSlug); err != nil {
		log.Printf("自动同步速查表失败: %v", err)
	}
}

func (s *CheatSheetService) GetPublished(slug string) (*CheatSheetPublishedPayload, error) {
	if slug != ClaudeCodeCheatSheetSlug {
		return nil, gorm.ErrRecordNotFound
	}

	snapshot, err := s.repo.GetPublished(slug)
	if err != nil {
		return nil, err
	}

	content, err := decodeCheatSheetContent(snapshot.ContentJSON)
	if err != nil {
		return nil, err
	}

	return &CheatSheetPublishedPayload{
		Slug:        snapshot.Slug,
		Status:      snapshot.Status,
		SourceURL:   snapshot.SourceURL,
		SyncedAt:    snapshot.CreatedAt,
		PublishedAt: snapshot.PublishedAt,
		Content:     content,
	}, nil
}

func (s *CheatSheetService) ListSnapshots(slug string, limit int) ([]CheatSheetSnapshotSummary, error) {
	if slug != ClaudeCodeCheatSheetSlug {
		return nil, gorm.ErrRecordNotFound
	}

	items, err := s.repo.ListSnapshots(slug, limit)
	if err != nil {
		return nil, err
	}

	result := make([]CheatSheetSnapshotSummary, 0, len(items))
	for _, item := range items {
		result = append(result, mapSnapshotSummary(item))
	}
	return result, nil
}

func (s *CheatSheetService) Publish(slug string, id uint) (*CheatSheetSnapshotSummary, error) {
	if slug != ClaudeCodeCheatSheetSlug {
		return nil, gorm.ErrRecordNotFound
	}

	snapshot, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if snapshot.Slug != slug {
		return nil, gorm.ErrRecordNotFound
	}
	if snapshot.Status == "failed" {
		return nil, errors.New("失败快照不能发布")
	}

	if err := s.repo.Publish(slug, id); err != nil {
		return nil, err
	}

	updated, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	summary := mapSnapshotSummary(*updated)
	return &summary, nil
}

func (s *CheatSheetService) Sync(slug string) (*CheatSheetSyncResult, error) {
	if slug != ClaudeCodeCheatSheetSlug {
		return nil, gorm.ErrRecordNotFound
	}

	sourceURL := config.AppConfig.CheatSheetSync.SourceURL
	html, err := s.fetchHTML(sourceURL)
	if err != nil {
		s.recordFailedSnapshot(slug, sourceURL, err)
		return nil, err
	}

	content, err := parseClaudeCodeCheatSheet(html, sourceURL)
	if err != nil {
		s.recordFailedSnapshot(slug, sourceURL, err)
		return nil, err
	}

	contentJSON, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}

	contentHash := sha256.Sum256(contentJSON)
	hashValue := hex.EncodeToString(contentHash[:])

	existing, err := s.repo.FindByHash(slug, hashValue)
	if err == nil {
		summary := mapSnapshotSummary(*existing)
		return &CheatSheetSyncResult{Changed: false, AutoPublished: false, Snapshot: &summary}, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	snapshot := &model.CheatSheetSnapshot{
		Slug:                slug,
		SourceURL:           sourceURL,
		SourceVersion:       content.Meta.Version,
		SourceUpdatedAtText: content.Meta.UpdatedAt,
		ContentJSON:         string(contentJSON),
		ContentHash:         hashValue,
		Status:              "draft",
	}
	if err := s.repo.Create(snapshot); err != nil {
		return nil, err
	}

	autoPublished := false
	if _, err := s.repo.GetPublished(slug); errors.Is(err, gorm.ErrRecordNotFound) {
		if err := s.repo.Publish(slug, snapshot.ID); err == nil {
			autoPublished = true
			snapshot, _ = s.repo.GetByID(snapshot.ID)
		}
	}

	summary := mapSnapshotSummary(*snapshot)
	return &CheatSheetSyncResult{
		Changed:       true,
		AutoPublished: autoPublished,
		Snapshot:      &summary,
	}, nil
}

func (s *CheatSheetService) fetchHTML(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "mdxy-cheatsheet-sync/1.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("抓取源页面失败: HTTP %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	html, err := doc.Html()
	if err != nil {
		return "", err
	}
	return html, nil
}

func (s *CheatSheetService) recordFailedSnapshot(slug, sourceURL string, syncErr error) {
	_ = s.repo.Create(&model.CheatSheetSnapshot{
		Slug:      slug,
		SourceURL: sourceURL,
		Status:    "failed",
		SyncError: syncErr.Error(),
	})
}

func parseClaudeCodeCheatSheet(rawHTML, sourceURL string) (CheatSheetContent, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawHTML))
	if err != nil {
		return CheatSheetContent{}, err
	}

	page := doc.Find(".page")
	if page.Length() == 0 {
		return CheatSheetContent{}, errors.New("未找到速查表页面容器")
	}

	content := CheatSheetContent{
		Meta: CheatSheetMeta{
			SourceURL:   sourceURL,
			SourceTitle: cleanText(page.Find(".header h1").First().Text()),
			Version:     cleanText(page.Find(".version-info").First().Text()),
			UpdatedAt:   strings.TrimPrefix(cleanText(page.Find(".last-updated").First().Text()), "最后更新："),
		},
		Changelog: parseChangelogItems(page.Find(".changelog-list li")),
		Footer:    parseFooterRows(page.Find(".footer-row")),
		Columns:   parseColumns(page.Find(".main-grid").Children()),
	}

	if content.Meta.SourceTitle == "" || len(content.Columns) == 0 {
		return CheatSheetContent{}, errors.New("源页面结构不完整，无法解析")
	}

	return content, nil
}

func parseChangelogItems(items *goquery.Selection) []CheatSheetCodeItem {
	result := make([]CheatSheetCodeItem, 0, items.Length())
	items.Each(func(_ int, item *goquery.Selection) {
		code := cleanText(item.Find("code").First().Text())
		text := cleanText(strings.TrimPrefix(item.Text(), code))
		result = append(result, CheatSheetCodeItem{
			Code: code,
			Text: strings.TrimLeft(text, "—- "),
		})
	})
	return result
}

func parseFooterRows(rows *goquery.Selection) []CheatSheetFooterRow {
	result := make([]CheatSheetFooterRow, 0, rows.Length())
	rows.Each(func(_ int, row *goquery.Selection) {
		label := strings.TrimSuffix(cleanText(row.Find(".footer-label").First().Text()), "：")
		items := make([]CheatSheetCodeItem, 0)
		row.Find(".footer-item").Each(func(_ int, item *goquery.Selection) {
			code := cleanText(item.Find("code").First().Text())
			text := cleanText(strings.TrimSpace(strings.TrimPrefix(item.Text(), code)))
			items = append(items, CheatSheetCodeItem{Code: code, Text: text})
		})
		result = append(result, CheatSheetFooterRow{Label: label, Items: items})
	})
	return result
}

func parseColumns(columns *goquery.Selection) [][]CheatSheetSection {
	result := make([][]CheatSheetSection, 0)
	columns.Each(func(_ int, column *goquery.Selection) {
		sections := make([]CheatSheetSection, 0)
		column.ChildrenFiltered(".section").Each(func(_ int, section *goquery.Selection) {
			title := cleanText(section.Find(".section-header").First().Text())
			sections = append(sections, CheatSheetSection{
				ID:     inferSectionID(title),
				Title:  title,
				Theme:  inferTheme(section),
				Groups: parseSectionGroups(section.Find(".section-content").First()),
			})
		})
		if len(sections) > 0 {
			result = append(result, sections)
		}
	})
	return result
}

func parseSectionGroups(content *goquery.Selection) []CheatSheetGroup {
	groups := make([]CheatSheetGroup, 0)
	current := CheatSheetGroup{}

	content.Children().Each(func(_ int, child *goquery.Selection) {
		switch {
		case child.HasClass("sub-header"):
			if current.Title != "" {
				groups = append(groups, current)
			}
			current = CheatSheetGroup{Title: cleanText(child.Text())}
		case child.HasClass("row"):
			row := parseRow(child)
			if row.Key == "" && row.Desc == "" {
				return
			}
			if current.Title == "" {
				current = CheatSheetGroup{Title: "未分类"}
			}
			if row.Key == "" && row.Desc != "" && len(current.Rows) > 0 {
				current.Rows[len(current.Rows)-1].Desc = row.Desc
				return
			}
			current.Rows = append(current.Rows, row)
		}
	})

	if current.Title != "" {
		groups = append(groups, current)
	}

	return groups
}

func parseRow(row *goquery.Selection) CheatSheetRow {
	keySelection := row.Find(".key").First()
	descSelection := row.Find(".desc").First()

	desc := cleanText(descSelection.Text())
	addedAt, _ := descSelection.Find(".badge-new").Attr("data-added")
	key := extractKeyText(keySelection)

	result := CheatSheetRow{
		Key:     key,
		Desc:    desc,
		AddedAt: strings.TrimSpace(addedAt),
	}
	if desc == "粘贴图片" {
		result.KeyVariant = "paste-image"
	}
	return result
}

func extractKeyText(selection *goquery.Selection) string {
	if selection.Length() == 0 {
		return ""
	}

	parts := make([]string, 0)
	selection.Contents().Each(func(_ int, node *goquery.Selection) {
		text := cleanText(node.Text())
		if node.HasClass("keycap") {
			parts = append(parts, "["+text+"]")
			return
		}
		if text != "" {
			parts = append(parts, text)
		}
	})

	return strings.Join(parts, "")
}

func inferTheme(section *goquery.Selection) string {
	for _, className := range strings.Fields(section.AttrOr("class", "")) {
		if strings.HasPrefix(className, "section-") && className != "section" {
			return strings.TrimPrefix(className, "section-")
		}
	}
	return "cli"
}

func inferSectionID(title string) string {
	switch title {
	case "⌨️ 键盘快捷键":
		return "keyboard"
	case "🔌 MCP 服务器":
		return "mcp"
	case "⚡ 斜杠命令":
		return "slash"
	case "📁 记忆与文件":
		return "memory"
	case "🧠 工作流与技巧":
		return "workflows"
	case "⚙️ 配置与环境":
		return "config"
	case "🔧 技能与代理":
		return "skills"
	case "🖥️ CLI 与标志":
		return "cli"
	default:
		return "section"
	}
}

func cleanText(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

func decodeCheatSheetContent(raw string) (CheatSheetContent, error) {
	var content CheatSheetContent
	if err := json.Unmarshal([]byte(raw), &content); err != nil {
		return CheatSheetContent{}, err
	}
	return content, nil
}

func mapSnapshotSummary(snapshot model.CheatSheetSnapshot) CheatSheetSnapshotSummary {
	return CheatSheetSnapshotSummary{
		ID:                  snapshot.ID,
		Slug:                snapshot.Slug,
		SourceVersion:       snapshot.SourceVersion,
		SourceUpdatedAtText: snapshot.SourceUpdatedAtText,
		Status:              snapshot.Status,
		SyncError:           snapshot.SyncError,
		CreatedAt:           snapshot.CreatedAt,
		PublishedAt:         snapshot.PublishedAt,
	}
}
