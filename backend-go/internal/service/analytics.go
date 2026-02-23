package service

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/xyu/mdxy/internal/config"
	"github.com/xyu/mdxy/internal/repository"
)

type AnalyticsService struct {
	repo           *repository.AnalyticsRepository
	noteService    *NoteService
	noteCountCache int
	noteCountTime  time.Time
	noteCountMu    sync.Mutex
}

func NewAnalyticsService(repo *repository.AnalyticsRepository, noteService *NoteService) *AnalyticsService {
	return &AnalyticsService{repo: repo, noteService: noteService}
}

type OverviewData struct {
	TotalPV   int64 `json:"total_pv"`
	TotalUV   int64 `json:"total_uv"`
	TodayPV   int64 `json:"today_pv"`
	TodayUV   int64 `json:"today_uv"`
	NoteCount int   `json:"note_count"`
	ToolCount int64 `json:"tool_count"`
}

func (s *AnalyticsService) GetOverview() (OverviewData, error) {
	o, err := s.repo.GetOverview()
	if err != nil {
		return OverviewData{}, err
	}
	toolCount, _ := s.repo.GetToolCount()
	noteCount := s.countNotes()
	return OverviewData{
		TotalPV:   o.TotalPV,
		TotalUV:   o.TotalUV,
		TodayPV:   o.TodayPV,
		TodayUV:   o.TodayUV,
		NoteCount: noteCount,
		ToolCount: toolCount,
	}, nil
}

func (s *AnalyticsService) countNotes() int {
	s.noteCountMu.Lock()
	defer s.noteCountMu.Unlock()
	if time.Since(s.noteCountTime) < 5*time.Minute && s.noteCountCache > 0 {
		return s.noteCountCache
	}
	count := 0
	notesDir := config.AppConfig.Content.NotesDir
	filepath.Walk(notesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			count++
		}
		return nil
	})
	s.noteCountCache = count
	s.noteCountTime = time.Now()
	return count
}

func (s *AnalyticsService) GetTrends(days int) ([]repository.DailyStat, []repository.DailyStat, error) {
	pv, err := s.repo.GetDailyPV(days)
	if err != nil {
		return nil, nil, err
	}
	uv, err := s.repo.GetDailyUV(days)
	if err != nil {
		return nil, nil, err
	}
	return pv, uv, nil
}

func (s *AnalyticsService) GetPopularPages(limit int) ([]repository.PageStat, error) {
	return s.repo.GetPopularPages(limit)
}

func (s *AnalyticsService) GetDeviceStats() ([]repository.DeviceStat, error) {
	return s.repo.GetDeviceStats()
}

func (s *AnalyticsService) GetBrowserStats() ([]repository.BrowserStat, error) {
	return s.repo.GetBrowserStats()
}

func (s *AnalyticsService) GetGeoStats() ([]repository.GeoStat, error) {
	return s.repo.GetGeoStats()
}
