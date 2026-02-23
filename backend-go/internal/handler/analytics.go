package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xyu/mdxy/internal/pkg/response"
	"github.com/xyu/mdxy/internal/service"
)

type AnalyticsHandler struct {
	service *service.AnalyticsService
}

func NewAnalyticsHandler(service *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: service}
}

func (h *AnalyticsHandler) Overview(c *gin.Context) {
	data, err := h.service.GetOverview()
	if err != nil {
		response.InternalServerError(c, "获取概览数据失败")
		return
	}
	response.Success(c, data)
}

func (h *AnalyticsHandler) Trends(c *gin.Context) {
	days := 30
	if d := c.Query("days"); d != "" {
		if v, err := strconv.Atoi(d); err == nil && v > 0 && v <= 365 {
			days = v
		}
	}
	pv, uv, err := h.service.GetTrends(days)
	if err != nil {
		response.InternalServerError(c, "获取趋势数据失败")
		return
	}
	response.Success(c, gin.H{"pv": pv, "uv": uv})
}

func (h *AnalyticsHandler) PopularPages(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}
	data, err := h.service.GetPopularPages(limit)
	if err != nil {
		response.InternalServerError(c, "获取热门页面失败")
		return
	}
	response.Success(c, data)
}

func (h *AnalyticsHandler) Devices(c *gin.Context) {
	data, err := h.service.GetDeviceStats()
	if err != nil {
		response.InternalServerError(c, "获取设备统计失败")
		return
	}
	response.Success(c, data)
}

func (h *AnalyticsHandler) Browsers(c *gin.Context) {
	data, err := h.service.GetBrowserStats()
	if err != nil {
		response.InternalServerError(c, "获取浏览器统计失败")
		return
	}
	response.Success(c, data)
}

func (h *AnalyticsHandler) Geo(c *gin.Context) {
	data, err := h.service.GetGeoStats()
	if err != nil {
		response.InternalServerError(c, "获取地域统计失败")
		return
	}
	response.Success(c, data)
}
