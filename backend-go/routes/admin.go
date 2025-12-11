package routes

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"mdxy-backend/config"
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
		// 获取上传的文件
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "未找到上传文件",
			})
			return
		}

		// 获取目标路径参数
		path := c.PostForm("path")

		// 验证文件扩展名
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !config.AllowedExtensions[ext] {
			allowedExts := make([]string, 0, len(config.AllowedExtensions))
			for k := range config.AllowedExtensions {
				allowedExts = append(allowedExts, k)
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fmt.Sprintf("不支持的文件类型，仅允许: %s", strings.Join(allowedExts, ", ")),
			})
			return
		}

		// 构建目标路径
		var targetDir string
		if path != "" {
			targetDir = services.SafePath(path)
			if targetDir == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "无效的路径",
				})
				return
			}

			// 如果目标是文件路径，获取其父目录
			if filepath.Ext(targetDir) != "" {
				targetDir = filepath.Dir(targetDir)
			}
		} else {
			targetDir = config.NotesDir
		}

		// 确保目录存在
		err = os.MkdirAll(targetDir, 0755)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "创建目录失败",
			})
			return
		}

		// 构建完整文件路径
		targetFilePath := filepath.Join(targetDir, file.Filename)

		// 检查文件是否已存在
		if _, err := os.Stat(targetFilePath); err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "文件已存在",
			})
			return
		}

		// 检查文件大小
		if file.Size > config.MaxUploadSize {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fmt.Sprintf("文件过大，最大允许 %.1fMB", float64(config.MaxUploadSize)/(1024*1024)),
			})
			return
		}

		// 保存文件
		if err := c.SaveUploadedFile(file, targetFilePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": fmt.Sprintf("上传失败: %s", err.Error()),
			})
			return
		}

		// 计算相对路径
		relPath, err := filepath.Rel(config.NotesDir, targetFilePath)
		if err != nil {
			relPath = targetFilePath
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"path": strings.ReplaceAll(relPath, "\\", "/"),
				"name": file.Filename,
				"size": file.Size,
			},
			"message": "上传成功",
		})
	})

	// 删除笔记文件
	adminGroup.DELETE("/notes/*notePath", func(c *gin.Context) {
		notePath := c.Param("notePath")
		// 清理路径，移除前导斜杠
		notePath = strings.TrimPrefix(notePath, "/")

		// 获取安全路径
		safePath := services.SafePath(notePath)
		if safePath == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "无效的路径",
			})
			return
		}

		// 检查文件是否存在
		if _, err := os.Stat(safePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "文件不存在",
			})
			return
		}

		// 删除文件或目录
		var err error
		fileInfo, _ := os.Stat(safePath)
		if fileInfo.IsDir() {
			err = os.RemoveAll(safePath)
		} else {
			err = os.Remove(safePath)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": fmt.Sprintf("删除失败: %s", err.Error()),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "删除成功",
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
