package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	// 注册笔记相关路由
	registerNotesRoutes(router, db)

	// 注册管理相关路由
	registerAdminRoutes(router, db)

	// 注册统计相关路由
	registerAnalyticsRoutes(router, db)
}
