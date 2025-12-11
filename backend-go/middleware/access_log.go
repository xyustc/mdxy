package middleware

import (
	"time"

	"mdxy-backend/database"
	"mdxy-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetCurrentTimestamp 获取当前时间戳
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// AccessLogMiddleware 访问日志记录中间件
func AccessLogMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 排除不需要记录的路径
		excludedPaths := []string{
			"/api/admin/", // 管理接口
			"/docs",       // API文档
			"/openapi.json",
			"/health",
			"/favicon.ico",
		}

		// 检查是否需要记录
		shouldLog := true
		for _, path := range excludedPaths {
			if len(c.Request.URL.Path) >= len(path) && c.Request.URL.Path[:len(path)] == path {
				shouldLog = false
				break
			}
		}

		// 记录开始时间
		startTime := time.Now()

		// 执行请求
		c.Next()

		// 计算响应时间
		responseTime := time.Since(startTime).Seconds() * 1000 // 转换为毫秒

		// 记录访问日志（仅记录笔记访问）
		if shouldLog && len(c.Request.URL.Path) >= 11 && c.Request.URL.Path[:11] == "/api/notes" {
			// 解析User-Agent
			userAgentString := c.GetHeader("User-Agent")
			uaInfo := utils.ParseUserAgent(userAgentString)

			// 获取客户端IP
			clientIP := c.ClientIP()

			// 获取用户标识符
			visitorID := c.GetHeader("X-Visitor-Id")
			if visitorID == "" {
				// 如果请求头中没有，则尝试从Cookie中获取
				visitorID, _ = c.Cookie("visitor_id")
			}

			// 准备日志数据
			logData := database.AccessLog{
				IPAddress:    clientIP,
				VisitorID:    visitorID,
				UserAgent:    userAgentString,
				Path:         c.Request.URL.Path,
				Method:       c.Request.Method,
				StatusCode:   c.Writer.Status(),
				ResponseTime: responseTime,
				Referer:      c.GetHeader("Referer"),
				DeviceType:   uaInfo.DeviceType,
				OS:           uaInfo.OS,
				Browser:      uaInfo.Browser,
				CreatedAt:    time.Now(),
			}

			// 异步写入数据库
			go func() {
				db.Create(&logData)
			}()
		}
	}
}
