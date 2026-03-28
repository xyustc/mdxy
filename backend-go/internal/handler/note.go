package handler

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xyu/mdxy/internal/pkg/response"
	"github.com/xyu/mdxy/internal/pkg/watermark"
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
	// Gin 的 *path 通配符会包含前导斜杠，需要去掉
	notePath = strings.TrimPrefix(notePath, "/")

	if notePath == "" {
		response.BadRequest(c, "笔记路径不能为空")
		return
	}

	clientIP := c.ClientIP()
	content, err := h.service.GetContent(notePath, clientIP)
	if err != nil {
		if os.IsPermission(err) {
			response.Forbidden(c, "无权访问该笔记")
			return
		}
		response.NotFound(c, "笔记不存在")
		return
	}

	// 注入水印
	watermarkedContent := watermark.Inject(content, clientIP)
	c.Header("X-Robots-Tag", "noindex, noarchive, nosnippet")

	response.Success(c, gin.H{
		"path":    notePath,
		"content": watermarkedContent,
	})
}

// AdminSetFeatured 设置笔记 featured 状态
func (h *NoteHandler) AdminSetFeatured(c *gin.Context) {
	var body struct {
		Path     string `json:"path" binding:"required"`
		Featured bool   `json:"featured"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.service.SetFeatured(body.Path, body.Featured); err != nil {
		response.InternalServerError(c, "设置失败")
		return
	}

	response.Success(c, nil)
}

// AdminGetContent 获取笔记原始内容（管理员，不注入水印）
func (h *NoteHandler) AdminGetContent(c *gin.Context) {
	notePath := strings.TrimPrefix(c.Param("path"), "/")
	if notePath == "" {
		response.BadRequest(c, "笔记路径不能为空")
		return
	}

	content, err := h.service.GetRawContent(notePath)
	if err != nil {
		if os.IsPermission(err) {
			response.Forbidden(c, "无权访问该笔记")
			return
		}
		response.NotFound(c, "笔记不存在")
		return
	}

	response.Success(c, gin.H{"path": notePath, "content": content})
}

// AdminSave 保存笔记内容
func (h *NoteHandler) AdminSave(c *gin.Context) {
	var body struct {
		Path    string `json:"path" binding:"required"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.service.SaveNote(body.Path, body.Content); err != nil {
		if os.IsPermission(err) {
			response.Forbidden(c, "无权操作该路径")
			return
		}
		response.InternalServerError(c, "保存失败")
		return
	}

	response.Success(c, nil)
}

// AdminDelete 删除笔记
func (h *NoteHandler) AdminDelete(c *gin.Context) {
	notePath := strings.TrimPrefix(c.Param("path"), "/")
	if notePath == "" {
		response.BadRequest(c, "笔记路径不能为空")
		return
	}

	if err := h.service.DeleteNote(notePath); err != nil {
		if os.IsPermission(err) {
			response.Forbidden(c, "无权删除该笔记")
			return
		}
		response.NotFound(c, "笔记不存在")
		return
	}

	response.Success(c, nil)
}

// AdminCreateDir 创建目录
func (h *NoteHandler) AdminCreateDir(c *gin.Context) {
	var body struct {
		Path string `json:"path" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.service.CreateDirectory(body.Path); err != nil {
		if os.IsPermission(err) {
			response.Forbidden(c, "无权操作该路径")
			return
		}
		response.InternalServerError(c, "创建目录失败")
		return
	}

	response.Success(c, nil)
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
