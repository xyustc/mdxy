package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"mdxy-backend/config"
	"mdxy-backend/database"
	"mdxy-backend/middleware"
	"mdxy-backend/routes"
)

func main() {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 初始化数据库
	db, err := gorm.Open(sqlite.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移数据库表
	database.InitDatabase(db)

	// 创建Gin路由器
	router := gin.Default()

	// 配置CORS跨域
	router.Use(middleware.CORSMiddleware())

	// 添加访问日志中间件
	router.Use(middleware.AccessLogMiddleware(db))

	// 注册路由
	routes.RegisterRoutes(router, db)

	// 根路径
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":   "Markdown 笔记系统 API",
			"docs":      "/docs",
			"notes_dir": config.NotesDir,
			"timestamp": middleware.GetCurrentTimestamp(),
		})
	})

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		// 检查笔记目录是否存在且可访问
		if _, err := os.Stat(config.NotesDir); os.IsNotExist(err) {
			c.JSON(503, gin.H{
				"status":               "error",
				"message":              "笔记目录不存在",
				"timestamp":            middleware.GetCurrentTimestamp(),
				"notes_dir_accessible": false,
				"service":              "mdxy-backend",
			})
			return
		}

		c.JSON(200, gin.H{
			"status":               "ok",
			"timestamp":            middleware.GetCurrentTimestamp(),
			"notes_dir_accessible": true,
			"service":              "mdxy-backend",
		})
	})

	// 启动服务器
	log.Printf("Server starting on port 8000")
	err = router.Run(":8000")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
