package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/xyu/mdxy/internal/model"
)

type AnalyticsRepository struct {
	db *gorm.DB
}

var excludedPopularPaths = []string{
	"/api/v1/profile",
	"/api/v1/tools",
	"/api/v1/tools/categories",
	"/api/v1/notes/tree",
}

func NewAnalyticsRepository(db *gorm.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

type DailyStat struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type PageStat struct {
	Path  string `json:"path"`
	Count int64  `json:"count"`
}

type DeviceStat struct {
	DeviceType string `json:"device_type"`
	Count      int64  `json:"count"`
}

type BrowserStat struct {
	Browser string `json:"browser"`
	Count   int64  `json:"count"`
}

type GeoStat struct {
	Country string `json:"country"`
	Region  string `json:"region"`
	Count   int64  `json:"count"`
}

type Overview struct {
	TotalPV int64 `json:"total_pv"`
	TotalUV int64 `json:"total_uv"`
	TodayPV int64 `json:"today_pv"`
	TodayUV int64 `json:"today_uv"`
}

func (r *AnalyticsRepository) GetOverview() (Overview, error) {
	var o Overview
	today := time.Now().In(time.Local).Format("2006-01-02")
	if err := r.db.Model(&model.AccessLog{}).Count(&o.TotalPV).Error; err != nil {
		return o, err
	}
	if err := r.db.Model(&model.AccessLog{}).
		Distinct("COALESCE(NULLIF(visitor_id, ''), ip_address)").
		Count(&o.TotalUV).Error; err != nil {
		return o, err
	}
	if err := r.db.Model(&model.AccessLog{}).
		Where("DATE(created_at, 'localtime') = ?", today).
		Count(&o.TodayPV).Error; err != nil {
		return o, err
	}
	if err := r.db.Model(&model.AccessLog{}).
		Where("DATE(created_at, 'localtime') = ?", today).
		Distinct("COALESCE(NULLIF(visitor_id, ''), ip_address)").
		Count(&o.TodayUV).Error; err != nil {
		return o, err
	}
	return o, nil
}

func (r *AnalyticsRepository) GetDailyPV(days int) ([]DailyStat, error) {
	var stats []DailyStat
	since := time.Now().In(time.Local).AddDate(0, 0, -days).Format("2006-01-02")
	err := r.db.Model(&model.AccessLog{}).
		Select("DATE(created_at, 'localtime') as date, COUNT(*) as count").
		Where("DATE(created_at, 'localtime') >= ?", since).
		Group("DATE(created_at, 'localtime')").Order("date ASC").
		Find(&stats).Error
	return stats, err
}

func (r *AnalyticsRepository) GetDailyUV(days int) ([]DailyStat, error) {
	var stats []DailyStat
	since := time.Now().In(time.Local).AddDate(0, 0, -days).Format("2006-01-02")
	err := r.db.Model(&model.AccessLog{}).
		Select("DATE(created_at, 'localtime') as date, COUNT(DISTINCT COALESCE(NULLIF(visitor_id, ''), ip_address)) as count").
		Where("DATE(created_at, 'localtime') >= ?", since).
		Group("DATE(created_at, 'localtime')").Order("date ASC").
		Find(&stats).Error
	return stats, err
}

func (r *AnalyticsRepository) GetPopularPages(limit int) ([]PageStat, error) {
	var stats []PageStat
	err := r.db.Model(&model.AccessLog{}).
		Select("path, COUNT(*) as count").
		Where("method = 'GET'").
		Where("path NOT IN ?", excludedPopularPaths).
		Group("path").Order("count DESC").Limit(limit).
		Find(&stats).Error
	return stats, err
}

func (r *AnalyticsRepository) GetDeviceStats() ([]DeviceStat, error) {
	var stats []DeviceStat
	err := r.db.Model(&model.AccessLog{}).
		Select("device_type, COUNT(*) as count").
		Where("device_type != ''").
		Group("device_type").Order("count DESC").
		Find(&stats).Error
	return stats, err
}

func (r *AnalyticsRepository) GetBrowserStats() ([]BrowserStat, error) {
	var stats []BrowserStat
	err := r.db.Model(&model.AccessLog{}).
		Select("browser, COUNT(*) as count").
		Where("browser != ''").
		Group("browser").Order("count DESC").
		Find(&stats).Error
	return stats, err
}

func (r *AnalyticsRepository) GetGeoStats() ([]GeoStat, error) {
	var stats []GeoStat
	err := r.db.Model(&model.AccessLog{}).
		Select("country, region, COUNT(*) as count").
		Where("country != ''").
		Group("country, region").Order("count DESC").
		Find(&stats).Error
	return stats, err
}

func (r *AnalyticsRepository) GetToolCount() (int64, error) {
	var count int64
	err := r.db.Model(&model.Tool{}).Count(&count).Error
	return count, err
}
