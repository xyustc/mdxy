package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xyu/mdxy/internal/pkg/response"
	"github.com/xyu/mdxy/internal/service"
)

type CheatSheetHandler struct {
	service *service.CheatSheetService
}

func NewCheatSheetHandler(service *service.CheatSheetService) *CheatSheetHandler {
	return &CheatSheetHandler{service: service}
}

func (h *CheatSheetHandler) GetPublished(c *gin.Context) {
	data, err := h.service.GetPublished(c.Param("slug"))
	if err != nil {
		response.NotFound(c, "未找到已发布的速查表")
		return
	}
	response.Success(c, data)
}

func (h *CheatSheetHandler) AdminListSnapshots(c *gin.Context) {
	limit := 10
	if raw := c.Query("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}

	data, err := h.service.ListSnapshots(c.Param("slug"), limit)
	if err != nil {
		response.InternalServerError(c, "获取速查表快照失败")
		return
	}
	response.Success(c, data)
}

func (h *CheatSheetHandler) AdminSync(c *gin.Context) {
	data, err := h.service.Sync(c.Param("slug"))
	if err != nil {
		response.InternalServerError(c, "同步速查表失败")
		return
	}
	response.Success(c, data)
}

func (h *CheatSheetHandler) AdminPublish(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的快照 ID")
		return
	}

	data, err := h.service.Publish(c.Param("slug"), uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}
