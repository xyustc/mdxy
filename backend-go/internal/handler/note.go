package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xyu/mdxy/internal/pkg/response"
	"github.com/xyu/mdxy/internal/service"
)

type NoteHandler struct {
	service *service.NoteService
}

func NewNoteHandler(service *service.NoteService) *NoteHandler {
	return &NoteHandler{service: service}
}

// GetTree 获取笔记目录树
func (h *NoteHandler) GetTree(c *gin.Context) {
	tree, err := h.service.GetTree()
	if err != nil {
		response.InternalServerError(c, "获取笔记目录失败")
		return
	}

	response.Success(c, tree)
}

// GetContent 获取笔记内容
func (h *NoteHandler) GetContent(c *gin.Context) {
	notePath := c.Param("path")

	content, err := h.service.GetContent(notePath)
	if err != nil {
		response.NotFound(c, "笔记不存在")
		return
	}

	response.Success(c, gin.H{
		"path":    notePath,
		"content": content,
	})
}

// Search 搜索笔记
func (h *NoteHandler) Search(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		response.BadRequest(c, "搜索关键词不能为空")
		return
	}

	results, err := h.service.Search(keyword)
	if err != nil {
		response.InternalServerError(c, "搜索失败")
		return
	}

	response.Success(c, gin.H{
		"results": results,
		"total":   len(results),
	})
}
