package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xyu/mdxy/internal/pkg/response"
	"github.com/xyu/mdxy/internal/service"
	"gorm.io/gorm"
)

type GameScoreHandler struct {
	service *service.GameScoreService
}

func NewGameScoreHandler(db *gorm.DB) *GameScoreHandler {
	return &GameScoreHandler{
		service: service.NewGameScoreService(db),
	}
}

// UpdateScore 更新分数
func (h *GameScoreHandler) UpdateScore(c *gin.Context) {
	var req struct {
		ToolID   uint   `json:"tool_id" binding:"required"`
		PlayerID string `json:"player_id" binding:"required"`
		Score    int    `json:"score" binding:"required,min=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	bestScore, err := h.service.UpdateScore(req.ToolID, req.PlayerID, req.Score)
	if err != nil {
		response.InternalServerError(c, "更新分数失败")
		return
	}

	response.Success(c, gin.H{"best_score": bestScore})
}

// GetBestScore 获取最高分
func (h *GameScoreHandler) GetBestScore(c *gin.Context) {
	toolIDStr := c.Query("tool_id")
	playerID := c.Query("player_id")

	if toolIDStr == "" || playerID == "" {
		response.BadRequest(c, "参数错误")
		return
	}

	toolID, err := strconv.ParseUint(toolIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "tool_id 格式错误")
		return
	}

	bestScore, err := h.service.GetBestScore(uint(toolID), playerID)
	if err != nil {
		response.InternalServerError(c, "获取最高分失败")
		return
	}

	response.Success(c, gin.H{"best_score": bestScore})
}

// GetLeaderboard 获取排行榜
func (h *GameScoreHandler) GetLeaderboard(c *gin.Context) {
	toolIDStr := c.Query("tool_id")
	if toolIDStr == "" {
		response.BadRequest(c, "tool_id 参数必填")
		return
	}

	toolID, err := strconv.ParseUint(toolIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "tool_id 格式错误")
		return
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	leaderboard, err := h.service.GetLeaderboard(uint(toolID), limit)
	if err != nil {
		response.InternalServerError(c, "获取排行榜失败")
		return
	}

	response.Success(c, leaderboard)
}
