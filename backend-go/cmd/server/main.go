package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xyu/mdxy/internal/config"
	"github.com/xyu/mdxy/internal/database"
	"github.com/xyu/mdxy/internal/router"
)

func main() {
	// 加载配置
	if err := config.Load(); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库
	if err := database.Init(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 设置 Gin 模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 创建路由
	r := gin.Default()

	// 设置路由
	router.Setup(r)

	// 启动服务器
	addr := ":" + config.AppConfig.Server.Port
	log.Printf("服务器启动在 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
