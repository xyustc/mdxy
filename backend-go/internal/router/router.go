package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xyu/mdxy/internal/database"
	"github.com/xyu/mdxy/internal/handler"
	"github.com/xyu/mdxy/internal/middleware"
	"github.com/xyu/mdxy/internal/repository"
	"github.com/xyu/mdxy/internal/service"
)

func Setup(r *gin.Engine) {
	// 全局中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// 初始化依赖
	profileRepo := repository.NewProfileRepository(database.DB)
	profileService := service.NewProfileService(profileRepo)
	profileHandler := handler.NewProfileHandler(profileService)

	adminRepo := repository.NewAdminRepository(database.DB)
	adminService := service.NewAdminService(adminRepo)
	adminHandler := handler.NewAdminHandler(adminService)

	noteService := service.NewNoteService()
	noteHandler := handler.NewNoteHandler(noteService)

	// API v1
	v1 := r.Group("/api/v1")
	{
		// 公开接口
		v1.GET("/profile", profileHandler.Get)

		// 笔记接口
		notes := v1.Group("/notes")
		{
			notes.GET("/tree", noteHandler.GetTree)
			notes.GET("/search", noteHandler.Search)
			notes.GET("/content/*path", noteHandler.GetContent)
		}

		// 管理员接口
		admin := v1.Group("/admin")
		{
			admin.POST("/login", adminHandler.Login)

			// 需要认证的接口
			authorized := admin.Use(middleware.AuthMiddleware())
			{
				authorized.PUT("/profile", profileHandler.Update)
			}
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}
