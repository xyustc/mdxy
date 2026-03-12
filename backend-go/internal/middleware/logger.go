package middleware

import (
	"crypto/md5"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xyu/mdxy/internal/database"
	"github.com/xyu/mdxy/internal/model"
	"github.com/xyu/mdxy/internal/pkg/geo"
)

var logChan chan model.AccessLog

func init() {
	logChan = make(chan model.AccessLog, 1000)
	go asyncLogWriter()
}

func asyncLogWriter() {
	batch := make([]model.AccessLog, 0, 50)
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 || database.DB == nil {
			return
		}
		if err := database.DB.Create(&batch).Error; err != nil {
			log.Printf("批量写入访问日志失败: %v", err)
		}
		batch = batch[:0]
	}

	for {
		select {
		case entry, ok := <-logChan:
			if !ok {
				flush()
				return
			}
			batch = append(batch, entry)
			if len(batch) >= 50 {
				flush()
			}
		case <-ticker.C:
			flush()
		}
	}
}

// Logger 访问日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		ua := c.Request.UserAgent()

		log.Printf("[%s] %s %s %d %v", method, path, clientIP, statusCode, latency)

		// 只记录 API 请求和页面请求，跳过静态资源
		if strings.HasPrefix(path, "/api/") || path == "/" || !strings.Contains(path, ".") {
			deviceType, osName, browser := parseUserAgent(ua)
			visitorID := generateVisitorID(clientIP, ua)
			country, region := geo.Lookup(clientIP)

			entry := model.AccessLog{
				IPAddress:    clientIP,
				VisitorID:    visitorID,
				UserAgent:    truncate(ua, 500),
				Path:         path,
				Method:       method,
				StatusCode:   statusCode,
				ResponseTime: float64(latency.Milliseconds()),
				Referer:      truncate(c.Request.Referer(), 500),
				DeviceType:   deviceType,
				OS:           osName,
				Browser:      browser,
				Country:      country,
				Region:       region,
			}

			select {
			case logChan <- entry:
			default:
				// channel 满了就丢弃，不阻塞请求
			}
		}
	}
}

func generateVisitorID(ip, ua string) string {
	hash := md5.Sum([]byte(ip + "|" + ua))
	return fmt.Sprintf("%x", hash)[:16]
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	// 确保不在 UTF-8 字符中间截断
	return strings.ToValidUTF8(s[:maxLen], "")
}

func parseUserAgent(ua string) (deviceType, osName, browser string) {
	lower := strings.ToLower(ua)

	// Device type
	switch {
	case (strings.Contains(lower, "mobile") || strings.Contains(lower, "android")) && !strings.Contains(lower, "tablet"):
		deviceType = "Mobile"
	case strings.Contains(lower, "tablet") || strings.Contains(lower, "ipad"):
		deviceType = "Tablet"
	case strings.Contains(lower, "bot") || strings.Contains(lower, "spider") || strings.Contains(lower, "crawl"):
		deviceType = "Bot"
	default:
		deviceType = "Desktop"
	}

	// OS
	switch {
	case strings.Contains(lower, "windows"):
		osName = "Windows"
	case strings.Contains(lower, "mac os") || strings.Contains(lower, "macintosh"):
		osName = "macOS"
	case strings.Contains(lower, "linux"):
		osName = "Linux"
	case strings.Contains(lower, "android"):
		osName = "Android"
	case strings.Contains(lower, "iphone") || strings.Contains(lower, "ipad") || strings.Contains(lower, "ios"):
		osName = "iOS"
	default:
		osName = "Other"
	}

	// Browser
	switch {
	case strings.Contains(lower, "edg"):
		browser = "Edge"
	case strings.Contains(lower, "chrome") && !strings.Contains(lower, "edg"):
		browser = "Chrome"
	case strings.Contains(lower, "firefox"):
		browser = "Firefox"
	case strings.Contains(lower, "safari") && !strings.Contains(lower, "chrome"):
		browser = "Safari"
	case strings.Contains(lower, "opera") || strings.Contains(lower, "opr"):
		browser = "Opera"
	default:
		browser = "Other"
	}

	return
}
