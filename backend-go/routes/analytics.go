package routes

import (
	"net/http"
	"time"

	"mdxy-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// registerAnalyticsRoutes 注册统计相关路由
func registerAnalyticsRoutes(router *gin.Engine, db *gorm.DB) {
	analyticsGroup := router.Group("/api/admin/analytics")

	// 获取访问日志列表
	analyticsGroup.GET("/logs", func(c *gin.Context) {
		// 解析查询参数
		page := c.DefaultQuery("page", "1")
		limit := c.DefaultQuery("limit", "50")
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")
		ip := c.Query("ip")
		path := c.Query("path")

		// 调用服务获取日志
		result := services.GetAccessLogs(db, page, limit, startDate, endDate, ip, path)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    result.Logs,
			"total":   result.Total,
			"page":    result.Page,
			"limit":   result.Limit,
			"pages":   result.Pages,
		})
	})

	// 获取统计概览
	analyticsGroup.GET("/overview", func(c *gin.Context) {
		// 解析查询参数
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")

		// 如果没有指定日期，默认最近30天
		if startDate == "" {
			startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02T15:04:05Z07:00")
		}
		if endDate == "" {
			endDate = time.Now().Format("2006-01-02T15:04:05Z07:00")
		}

		// 调用服务获取统计数据
		stats := services.GetOverviewStats(db, startDate, endDate)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    stats,
		})
	})
}
