package database

import (
	"time"

	"gorm.io/gorm"
)

// AccessLog 访问日志表
type AccessLog struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	IPAddress    string    `gorm:"index;size:45" json:"ip_address"` // 支持IPv6
	VisitorID    string    `gorm:"index;size:36" json:"visitor_id"` // 用户唯一标识
	UserAgent    string    `gorm:"size:500" json:"user_agent"`
	Path         string    `gorm:"index;size:500" json:"path"`
	Method       string    `gorm:"size:10" json:"method"`
	StatusCode   int       `json:"status_code"`
	ResponseTime float64   `json:"response_time"` // 毫秒
	Referer      string    `gorm:"size:500" json:"referer"`
	DeviceType   string    `gorm:"size:50" json:"device_type"` // PC/Mobile/Tablet
	OS           string    `gorm:"size:50" json:"os"`
	Browser      string    `gorm:"size:50" json:"browser"`
	CreatedAt    time.Time `gorm:"index" json:"created_at"`
}

// TableName 设置表名
func (AccessLog) TableName() string {
	return "access_logs"
}

// ToMap 转换为map
func (a *AccessLog) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":            a.ID,
		"ip_address":    a.IPAddress,
		"visitor_id":    a.VisitorID,
		"user_agent":    a.UserAgent,
		"path":          a.Path,
		"method":        a.Method,
		"status_code":   a.StatusCode,
		"response_time": a.ResponseTime,
		"referer":       a.Referer,
		"device_type":   a.DeviceType,
		"os":            a.OS,
		"browser":       a.Browser,
		"created_at":    a.CreatedAt.Format(time.RFC3339),
	}
}

// InitDatabase 初始化数据库表
func InitDatabase(db *gorm.DB) {
	// 自动迁移数据库表
	db.AutoMigrate(&AccessLog{})
}
