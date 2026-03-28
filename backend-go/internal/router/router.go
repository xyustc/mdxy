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

	toolRepo := repository.NewToolRepository(database.DB)
	toolService := service.NewToolService(toolRepo)
	toolHandler := handler.NewToolHandler(toolService)

	homeHandler := handler.NewHomeHandler(noteService, toolService)

	analyticsRepo := repository.NewAnalyticsRepository(database.DB)
	analyticsService := service.NewAnalyticsService(analyticsRepo, noteService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	searchService := service.NewSearchService(noteService, toolRepo, database.DB)
	searchHandler := handler.NewSearchHandler(searchService)

	gameScoreHandler := handler.NewGameScoreHandler(database.DB)

	cheatSheetRepo := repository.NewCheatSheetRepository(database.DB)
	cheatSheetService := service.NewCheatSheetService(cheatSheetRepo)
	cheatSheetHandler := handler.NewCheatSheetHandler(cheatSheetService)
	cheatSheetService.StartScheduler()

	// API v1
	v1 := r.Group("/api/v1")
	{
		// 公开接口
		v1.GET("/profile", middleware.RateLimitPublicAPI(), profileHandler.Get)
		v1.GET("/home", middleware.RateLimitPublicAPI(), homeHandler.Get)

		// 笔记接口
		notes := v1.Group("/notes")
		{
			notes.GET("/tree", middleware.RateLimitPublicAPI(), noteHandler.GetTree)
			notes.GET("/search", middleware.RateLimitPublicAPI(), noteHandler.Search)
			notes.GET("/content/*path", middleware.RateLimitPublicAPI(), middleware.RateLimitNoteContent(), noteHandler.GetContent)
		}

		// 工具公开接口
		tools := v1.Group("/tools")
		{
			tools.GET("", middleware.RateLimitPublicAPI(), toolHandler.List)
			tools.GET("/categories", middleware.RateLimitPublicAPI(), toolHandler.GetCategories)
		}

		// 搜索公开接口
		search := v1.Group("/search")
		{
			search.GET("", middleware.RateLimitPublicAPI(), searchHandler.Search)
			search.GET("/popular", middleware.RateLimitPublicAPI(), searchHandler.Popular)
		}

		cheatSheets := v1.Group("/cheat-sheets")
		{
			cheatSheets.GET("/:slug", middleware.RateLimitPublicAPI(), cheatSheetHandler.GetPublished)
		}

		// 游戏分数接口
		games := v1.Group("/games")
		{
			games.POST("/score", middleware.RateLimitPublicAPI(), gameScoreHandler.UpdateScore)
			games.GET("/score", middleware.RateLimitPublicAPI(), gameScoreHandler.GetBestScore)
			games.GET("/leaderboard", middleware.RateLimitPublicAPI(), gameScoreHandler.GetLeaderboard)
		}

		// 管理员接口
		admin := v1.Group("/admin")
		{
			admin.POST("/login", adminHandler.Login)

			// 需要认证的接口
			authorized := admin.Use(middleware.AuthMiddleware())
			{
				authorized.PUT("/profile", profileHandler.Update)

				// 工具管理
				authorized.GET("/tools", toolHandler.AdminList)
				authorized.GET("/tools/:id", toolHandler.AdminGetByID)
				authorized.POST("/tools", toolHandler.AdminCreate)
				authorized.PUT("/tools/:id", toolHandler.AdminUpdate)
				authorized.DELETE("/tools/:id", toolHandler.AdminDelete)

				// 数据统计
				authorized.GET("/analytics/overview", analyticsHandler.Overview)
				authorized.GET("/analytics/trends", analyticsHandler.Trends)
				authorized.GET("/analytics/popular", analyticsHandler.PopularPages)
				authorized.GET("/analytics/devices", analyticsHandler.Devices)
				authorized.GET("/analytics/browsers", analyticsHandler.Browsers)
				authorized.GET("/analytics/geo", analyticsHandler.Geo)

				// 笔记管理
			authorized.GET("/notes/tree", noteHandler.GetTree)
			authorized.GET("/notes/content/*path", noteHandler.AdminGetContent)
			authorized.POST("/notes/save", noteHandler.AdminSave)
			authorized.POST("/notes/featured", noteHandler.AdminSetFeatured)
			authorized.POST("/notes/directory", noteHandler.AdminCreateDir)
			authorized.DELETE("/notes/*path", noteHandler.AdminDelete)

			authorized.GET("/cheat-sheets/:slug/snapshots", cheatSheetHandler.AdminListSnapshots)
				authorized.POST("/cheat-sheets/:slug/sync", cheatSheetHandler.AdminSync)
				authorized.POST("/cheat-sheets/:slug/publish/:id", cheatSheetHandler.AdminPublish)
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
