package services

import (
	"math"
	"strconv"
	"time"

	"mdxy-backend/database"

	"gorm.io/gorm"
)

// LogsResult 日志查询结果
type LogsResult struct {
	Logs  []map[string]interface{} `json:"logs"`
	Total int64                    `json:"total"`
	Page  int                      `json:"page"`
	Limit int                      `json:"limit"`
	Pages int                      `json:"pages"`
}

// GetAccessLogs 获取访问日志列表
func GetAccessLogs(
	db *gorm.DB,
	pageStr string,
	limitStr string,
	startDate string,
	endDate string,
	ip string,
	path string,
) LogsResult {
	// 解析分页参数
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 50
	}

	// 构建查询
	query := db.Model(&database.AccessLog{})

	// 日期过滤
	if startDate != "" {
		if start, err := time.Parse(time.RFC3339, startDate); err == nil {
			query = query.Where("created_at >= ?", start)
		}
	}

	if endDate != "" {
		if end, err := time.Parse(time.RFC3339, endDate); err == nil {
			query = query.Where("created_at <= ?", end)
		}
	}

	// IP过滤
	if ip != "" {
		query = query.Where("ip_address LIKE ?", "%"+ip+"%")
	}

	// 路径过滤
	if path != "" {
		query = query.Where("path LIKE ?", "%"+path+"%")
	}

	// 总数
	var total int64
	query.Count(&total)

	// 分页
	offset := (page - 1) * limit
	var logs []database.AccessLog
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs)

	// 转换为map
	logMaps := make([]map[string]interface{}, len(logs))
	for i, log := range logs {
		logMaps[i] = log.ToMap()
	}

	return LogsResult{
		Logs:  logMaps,
		Total: total,
		Page:  page,
		Limit: limit,
		Pages: int(math.Ceil(float64(total) / float64(limit))),
	}
}

// OverviewStats 统计概览
type OverviewStats struct {
	TotalVisits     int64                    `json:"total_visits"`
	UniqueVisitors  int64                    `json:"unique_visitors"`
	AvgResponseTime float64                  `json:"avg_response_time"`
	TopPages        []map[string]interface{} `json:"top_pages"`
	VisitorTrends   []map[string]interface{} `json:"visitor_trends"`
	DeviceStats     map[string]int64         `json:"device_stats"`
	OsStats         []map[string]interface{} `json:"os_stats"`
	BrowserStats    []map[string]interface{} `json:"browser_stats"`
}

// GetOverviewStats 获取统计概览
func GetOverviewStats(
	db *gorm.DB,
	startDate string,
	endDate string,
) OverviewStats {
	// 构建基础查询
	query := db.Model(&database.AccessLog{})

	// 日期过滤
	if startDate != "" {
		if start, err := time.Parse(time.RFC3339, startDate); err == nil {
			query = query.Where("created_at >= ?", start)
		}
	}

	if endDate != "" {
		if end, err := time.Parse(time.RFC3339, endDate); err == nil {
			query = query.Where("created_at <= ?", end)
		}
	}

	// 总访问量
	var totalVisits int64
	query.Count(&totalVisits)

	// 独立访客数（优先使用visitor_id，备选IP地址）
	var uniqueVisitors int64
	// 构建带日期过滤的查询
	filteredQuery := db.Model(&database.AccessLog{})
	if startDate != "" {
		if start, err := time.Parse(time.RFC3339, startDate); err == nil {
			filteredQuery = filteredQuery.Where("created_at >= ?", start)
		}
	}
	if endDate != "" {
		if end, err := time.Parse(time.RFC3339, endDate); err == nil {
			filteredQuery = filteredQuery.Where("created_at <= ?", end)
		}
	}

	// 先尝试使用visitor_id统计
	var uniqueVisitorsByID int64
	filteredQuery.Where("visitor_id IS NOT NULL").
		Distinct("visitor_id").
		Count(&uniqueVisitorsByID)

	// 如果没有visitor_id记录，则使用IP地址统计
	if uniqueVisitorsByID == 0 {
		filteredQuery.Distinct("ip_address").
			Count(&uniqueVisitors)
	} else {
		uniqueVisitors = uniqueVisitorsByID
	}

	// 平均响应时间
	var avgResponseTime float64
	// 构建带日期过滤的查询
	filteredQuery2 := db.Model(&database.AccessLog{})
	if startDate != "" {
		if start, err := time.Parse(time.RFC3339, startDate); err == nil {
			filteredQuery2 = filteredQuery2.Where("created_at >= ?", start)
		}
	}
	if endDate != "" {
		if end, err := time.Parse(time.RFC3339, endDate); err == nil {
			filteredQuery2 = filteredQuery2.Where("created_at <= ?", end)
		}
	}
	filteredQuery2.Select("AVG(response_time)").Scan(&avgResponseTime)

	// Top访问页面
	type TopPage struct {
		Path  string
		Count int64
	}
	var topPagesResult []TopPage
	// 构建带日期过滤的查询
	filteredQuery3 := db.Model(&database.AccessLog{})
	if startDate != "" {
		if start, err := time.Parse(time.RFC3339, startDate); err == nil {
			filteredQuery3 = filteredQuery3.Where("created_at >= ?", start)
		}
	}
	if endDate != "" {
		if end, err := time.Parse(time.RFC3339, endDate); err == nil {
			filteredQuery3 = filteredQuery3.Where("created_at <= ?", end)
		}
	}
	filteredQuery3.Select("path, COUNT(id) as count").
		Group("path").
		Order("count DESC").
		Limit(10).
		Scan(&topPagesResult)

	topPages := make([]map[string]interface{}, len(topPagesResult))
	for i, tp := range topPagesResult {
		topPages[i] = map[string]interface{}{
			"path":  tp.Path,
			"count": tp.Count,
		}
	}

	// 访问趋势（按天）
	type VisitorTrend struct {
		Date  time.Time
		Count int64
	}
	var visitorTrendsResult []VisitorTrend
	// 构建带日期过滤的查询
	filteredQuery4 := db.Model(&database.AccessLog{})
	if startDate != "" {
		if start, err := time.Parse(time.RFC3339, startDate); err == nil {
			filteredQuery4 = filteredQuery4.Where("created_at >= ?", start)
		}
	}
	if endDate != "" {
		if end, err := time.Parse(time.RFC3339, endDate); err == nil {
			filteredQuery4 = filteredQuery4.Where("created_at <= ?", end)
		}
	}
	filteredQuery4.Select("DATE(created_at) as date, COUNT(id) as count").
		Group("DATE(created_at)").
		Order("date").
		Scan(&visitorTrendsResult)

	visitorTrends := make([]map[string]interface{}, len(visitorTrendsResult))
	for i, vt := range visitorTrendsResult {
		visitorTrends[i] = map[string]interface{}{
			"date":  vt.Date.Format("2006-01-02"),
			"count": vt.Count,
		}
	}

	// 设备统计
	type DeviceStat struct {
		DeviceType string
		Count      int64
	}
	var deviceStatsResult []DeviceStat
	// 构建带日期过滤的查询
	filteredQuery5 := db.Model(&database.AccessLog{})
	if startDate != "" {
		if start, err := time.Parse(time.RFC3339, startDate); err == nil {
			filteredQuery5 = filteredQuery5.Where("created_at >= ?", start)
		}
	}
	if endDate != "" {
		if end, err := time.Parse(time.RFC3339, endDate); err == nil {
			filteredQuery5 = filteredQuery5.Where("created_at <= ?", end)
		}
	}
	filteredQuery5.Select("device_type, COUNT(id) as count").
		Group("device_type").
		Scan(&deviceStatsResult)

	deviceStats := make(map[string]int64)
	for _, ds := range deviceStatsResult {
		deviceStats[ds.DeviceType] = ds.Count
	}

	// 操作系统统计
	type OsStat struct {
		OS    string
		Count int64
	}
	var osStatsResult []OsStat
	// 构建带日期过滤的查询
	filteredQuery6 := db.Model(&database.AccessLog{})
	if startDate != "" {
		if start, err := time.Parse(time.RFC3339, startDate); err == nil {
			filteredQuery6 = filteredQuery6.Where("created_at >= ?", start)
		}
	}
	if endDate != "" {
		if end, err := time.Parse(time.RFC3339, endDate); err == nil {
			filteredQuery6 = filteredQuery6.Where("created_at <= ?", end)
		}
	}
	filteredQuery6.Select("os, COUNT(id) as count").
		Group("os").
		Order("count DESC").
		Limit(10).
		Scan(&osStatsResult)

	osStats := make([]map[string]interface{}, len(osStatsResult))
	for i, os := range osStatsResult {
		osStats[i] = map[string]interface{}{
			"os":    os.OS,
			"count": os.Count,
		}
	}

	// 浏览器统计
	type BrowserStat struct {
		Browser string
		Count   int64
	}
	var browserStatsResult []BrowserStat
	// 构建带日期过滤的查询
	filteredQuery7 := db.Model(&database.AccessLog{})
	if startDate != "" {
		if start, err := time.Parse(time.RFC3339, startDate); err == nil {
			filteredQuery7 = filteredQuery7.Where("created_at >= ?", start)
		}
	}
	if endDate != "" {
		if end, err := time.Parse(time.RFC3339, endDate); err == nil {
			filteredQuery7 = filteredQuery7.Where("created_at <= ?", end)
		}
	}
	filteredQuery7.Select("browser, COUNT(id) as count").
		Group("browser").
		Order("count DESC").
		Limit(10).
		Scan(&browserStatsResult)

	browserStats := make([]map[string]interface{}, len(browserStatsResult))
	for i, bs := range browserStatsResult {
		browserStats[i] = map[string]interface{}{
			"browser": bs.Browser,
			"count":   bs.Count,
		}
	}

	return OverviewStats{
		TotalVisits:     totalVisits,
		UniqueVisitors:  uniqueVisitors,
		AvgResponseTime: math.Round(avgResponseTime*100) / 100,
		TopPages:        topPages,
		VisitorTrends:   visitorTrends,
		DeviceStats:     deviceStats,
		OsStats:         osStats,
		BrowserStats:    browserStats,
	}
}
