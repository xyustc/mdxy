package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xyu/mdxy/internal/pkg/response"
	"github.com/xyu/mdxy/internal/service"
)

type HomeHandler struct {
	noteService *service.NoteService
	toolService *service.ToolService
}

func NewHomeHandler(noteService *service.NoteService, toolService *service.ToolService) *HomeHandler {
	return &HomeHandler{noteService: noteService, toolService: toolService}
}

type HomeResponse struct {
	FeaturedNotes []service.NoteNode `json:"featured_notes"`
	NoteCount     int                `json:"note_count"`
	FeaturedTools []interface{}      `json:"featured_tools"`
	ToolCount     int                `json:"tool_count"`
}

// Get 返回首页所需数据
func (h *HomeHandler) Get(c *gin.Context) {
	tree, err := h.noteService.GetTree()
	if err != nil {
		response.InternalServerError(c, "获取笔记失败")
		return
	}

	allNotes := flattenNoteNodes(tree)
	featuredNotes := make([]service.NoteNode, 0)
	for _, n := range allNotes {
		if n.Featured {
			featuredNotes = append(featuredNotes, n)
		}
	}

	allTools, err := h.toolService.List("", true)
	if err != nil {
		response.InternalServerError(c, "获取工具失败")
		return
	}

	featuredTools, err := h.toolService.ListFeatured()
	if err != nil {
		response.InternalServerError(c, "获取精选工具失败")
		return
	}

	featuredToolsIface := make([]interface{}, len(featuredTools))
	for i, t := range featuredTools {
		featuredToolsIface[i] = t
	}

	response.Success(c, gin.H{
		"featured_notes": featuredNotes,
		"note_count":     len(allNotes),
		"featured_tools": featuredToolsIface,
		"tool_count":     len(allTools),
	})
}

func flattenNoteNodes(nodes []service.NoteNode) []service.NoteNode {
	var result []service.NoteNode
	for _, n := range nodes {
		if n.Type == "file" {
			result = append(result, n)
		} else {
			result = append(result, flattenNoteNodes(n.Children)...)
		}
	}
	return result
}
