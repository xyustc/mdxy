package routes

import (
	"net/http"

	"mdxy-backend/services"
	"mdxy-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// registerAdminRoutes 注册管理相关路由
func registerAdminRoutes(router *gin.Engine, db *gorm.DB) {
	adminGroup := router.Group("/api/admin")

	// 管理员登录
	adminGroup.POST("/login", func(c *gin.Context) {
		var req struct {
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "请求参数错误",
			})
			return
		}

		result := services.AuthenticateAdmin(req.Password)
		if !result.Success {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": result.Message,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"token": result.Token,
			},
			"message": result.Message,
		})
	})

	// 验证token是否有效
	adminGroup.GET("/verify", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "缺少认证信息",
			})
			return
		}

		// 提取token
		tokenString := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "认证格式错误",
			})
			return
		}

		// 验证token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "无效的认证凭证",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    claims,
			"message": "Token有效",
		})
	})

	// 上传笔记文件
	adminGroup.POST("/notes/upload", func(c *gin.Context) {
		// TODO: 实现上传笔记文件功能
		c.JSON(http.StatusNotImplemented, gin.H{
			"success": false,
			"message": "功能待实现",
		})
	})

	// 删除笔记文件
	adminGroup.DELETE("/notes/*notePath", func(c *gin.Context) {
		// TODO: 实现删除笔记文件功能
		c.JSON(http.StatusNotImplemented, gin.H{
			"success": false,
			"message": "功能待实现",
		})
	})

	// 创建目录
	adminGroup.POST("/notes/mkdir", func(c *gin.Context) {
		// TODO: 实现创建目录功能
		c.JSON(http.StatusNotImplemented, gin.H{
			"success": false,
			"message": "功能待实现",
		})
	})

	// 移动/重命名文件
	adminGroup.PUT("/notes/move", func(c *gin.Context) {
		// TODO: 实现移动/重命名文件功能
		c.JSON(http.StatusNotImplemented, gin.H{
			"success": false,
			"message": "功能待实现",
		})
	})
}
