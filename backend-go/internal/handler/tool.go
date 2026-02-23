package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xyu/mdxy/internal/model"
	"github.com/xyu/mdxy/internal/pkg/response"
	"github.com/xyu/mdxy/internal/service"
)

type ToolHandler struct {
	service *service.ToolService
}

func NewToolHandler(service *service.ToolService) *ToolHandler {
	return &ToolHandler{service: service}
}

// List 公开接口 - 获取工具列表
func (h *ToolHandler) List(c *gin.Context) {
	category := c.Query("category")
	tools, err := h.service.List(category, true)
	if err != nil {
		response.InternalServerError(c, "获取工具列表失败")
		return
	}
	response.Success(c, tools)
}

// GetCategories 公开接口 - 获取分类列表
func (h *ToolHandler) GetCategories(c *gin.Context) {
	categories, err := h.service.GetCategories()
	if err != nil {
		response.InternalServerError(c, "获取分类列表失败")
		return
	}
	response.Success(c, categories)
}

// AdminGetByID 管理接口 - 获取单个工具
func (h *ToolHandler) AdminGetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的工具ID")
		return
	}
	tool, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "工具不存在")
		return
	}
	response.Success(c, tool)
}

// AdminCreate 管理接口 - 创建工具
func (h *ToolHandler) AdminCreate(c *gin.Context) {
	var tool model.Tool
	if err := c.ShouldBindJSON(&tool); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.service.Create(&tool); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, tool)
}

// AdminUpdate 管理接口 - 更新工具
func (h *ToolHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的工具ID")
		return
	}
	existing, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "工具不存在")
		return
	}
	if err := c.ShouldBindJSON(existing); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	existing.ID = uint(id)
	if err := h.service.Update(existing); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, existing)
}

// AdminDelete 管理接口 - 删除工具
func (h *ToolHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的工具ID")
		return
	}
	if err := h.service.Delete(uint(id)); err != nil {
		response.InternalServerError(c, "删除失败")
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}
