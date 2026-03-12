package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xyu/mdxy/internal/pkg/response"
	"github.com/xyu/mdxy/internal/service"
)

type SearchHandler struct {
	service *service.SearchService
}

func NewSearchHandler(service *service.SearchService) *SearchHandler {
	return &SearchHandler{service: service}
}

func (h *SearchHandler) Search(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		response.BadRequest(c, "搜索关键词不能为空")
		return
	}
	searchType := c.DefaultQuery("type", "all")
	results, err := h.service.Search(q, searchType, c.ClientIP())
	if err != nil {
		response.InternalServerError(c, "搜索失败")
		return
	}
	response.Success(c, results)
}

func (h *SearchHandler) Popular(c *gin.Context) {
	keywords, err := h.service.PopularKeywords(10)
	if err != nil {
		response.InternalServerError(c, "获取热门搜索失败")
		return
	}
	response.Success(c, keywords)
}
