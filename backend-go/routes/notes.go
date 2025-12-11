package routes

import (
	"net/http"

	"mdxy-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// registerNotesRoutes 注册笔记相关路由
func registerNotesRoutes(router *gin.Engine, db *gorm.DB) {
	notesGroup := router.Group("/api/notes")

	// 获取笔记目录树
	notesGroup.GET("", func(c *gin.Context) {
		tree := services.GetNoteTree()
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    tree,
		})
	})

	// 搜索笔记
	notesGroup.GET("/search", func(c *gin.Context) {
		q := c.Query("q")
		if q == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "搜索关键词不能为空",
			})
			return
		}

		results := services.SearchNotes(q)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    results,
			"total":   len(results),
		})
	})

	// 获取指定笔记的内容 - 使用参数化路由而不是通配符
	notesGroup.GET("/:notePath", func(c *gin.Context) {
		notePath := c.Param("notePath")

		// 处理嵌套路径
		if c.Query("path") != "" {
			notePath = c.Query("path")
		}

		content := services.GetNoteContent(notePath)
		if content == "" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "笔记不存在",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"path":    notePath,
				"content": content,
			},
		})
	})
}
