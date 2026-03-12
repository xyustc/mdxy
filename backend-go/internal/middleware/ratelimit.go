package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xyu/mdxy/internal/pkg/response"
)

// 滑动窗口限流器
type slidingWindowLimiter struct {
	mu      sync.RWMutex
	windows map[string]*window
	limit   int           // 窗口内最大请求数
	window  time.Duration // 窗口大小
}

type window struct {
	requests []time.Time
	mu       sync.Mutex
}

var noteContentLimiter *slidingWindowLimiter

func init() {
	// 初始化限流器：每分钟30次
	noteContentLimiter = &slidingWindowLimiter{
		windows: make(map[string]*window),
		limit:   30,
		window:  time.Minute,
	}

	// 定期清理过期窗口
	go noteContentLimiter.cleanup()
}

// RateLimitNoteContent 笔记内容接口限流中间件
func RateLimitNoteContent() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if !noteContentLimiter.allow(clientIP) {
			c.Header("Retry-After", "60")
			response.TooManyRequests(c, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

// allow 检查是否允许请求
func (l *slidingWindowLimiter) allow(key string) bool {
	l.mu.Lock()
	w, exists := l.windows[key]
	if !exists {
		w = &window{requests: make([]time.Time, 0)}
		l.windows[key] = w
	}
	l.mu.Unlock()

	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-l.window)

	// 移除过期请求
	validRequests := make([]time.Time, 0)
	for _, t := range w.requests {
		if t.After(cutoff) {
			validRequests = append(validRequests, t)
		}
	}
	w.requests = validRequests

	// 检查是否超限
	if len(w.requests) >= l.limit {
		return false
	}

	// 记录本次请求
	w.requests = append(w.requests, now)
	return true
}

// cleanup 定期清理过期窗口
func (l *slidingWindowLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		l.mu.Lock()
		now := time.Now()
		for key, w := range l.windows {
			w.mu.Lock()
			// 如果窗口内没有有效请求，删除该窗口
			if len(w.requests) == 0 || w.requests[len(w.requests)-1].Add(l.window).Before(now) {
				delete(l.windows, key)
			}
			w.mu.Unlock()
		}
		l.mu.Unlock()
	}
}
