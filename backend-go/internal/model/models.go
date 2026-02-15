package model

import (
	"time"

	"gorm.io/gorm"
)

// Profile 个人信息
type Profile struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Title       string    `gorm:"size:200" json:"title"`
	Bio         string    `gorm:"type:text" json:"bio"`
	Avatar      string    `gorm:"size:500" json:"avatar"`
	Email       string    `gorm:"size:100" json:"email"`
	GitHub      string    `gorm:"size:200" json:"github"`
	LinkedIn    string    `gorm:"size:200" json:"linkedin"`
	Twitter     string    `gorm:"size:200" json:"twitter"`
	Website     string    `gorm:"size:200" json:"website"`
	Skills      string    `gorm:"type:text" json:"skills"` // JSON array
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Admin 管理员
type Admin struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:50;uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	LastLoginAt  *time.Time `json:"last_login_at"`
}

// AccessLog 访问日志
type AccessLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	IPAddress    string    `gorm:"size:45;index" json:"ip_address"`
	VisitorID    string    `gorm:"size:36;index" json:"visitor_id"`
	UserAgent    string    `gorm:"size:500" json:"user_agent"`
	Path         string    `gorm:"size:500;index" json:"path"`
	Method       string    `gorm:"size:10" json:"method"`
	StatusCode   int       `json:"status_code"`
	ResponseTime float64   `json:"response_time"` // 毫秒
	Referer      string    `gorm:"size:500" json:"referer"`
	DeviceType   string    `gorm:"size:50" json:"device_type"`
	OS           string    `gorm:"size:50" json:"os"`
	Browser      string    `gorm:"size:50" json:"browser"`
	CreatedAt    time.Time `gorm:"index" json:"created_at"`
}

// Category 文章分类
type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Slug      string    `gorm:"size:50;uniqueIndex;not null" json:"slug"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

// Tag 标签
type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Slug      string    `gorm:"size:50;uniqueIndex;not null" json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

// Article 文章
type Article struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:200;not null" json:"title"`
	Slug        string    `gorm:"size:200;uniqueIndex;not null" json:"slug"`
	FilePath    string    `gorm:"size:500;not null" json:"file_path"`
	Summary     string    `gorm:"type:text" json:"summary"`
	CoverImage  string    `gorm:"size:500" json:"cover_image"`
	CategoryID  uint      `gorm:"index" json:"category_id"`
	Category    *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags        []Tag     `gorm:"many2many:article_tags;" json:"tags,omitempty"`
	Status      int8      `gorm:"default:0;index" json:"status"` // 0=草稿 1=已发布
	ViewCount   uint      `gorm:"default:0" json:"view_count"`
	PublishedAt *time.Time `gorm:"index" json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Tool 工具
type Tool struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Type        string    `gorm:"size:20;index" json:"type"` // video/software/game/link
	URL         string    `gorm:"size:500" json:"url"`
	Icon        string    `gorm:"size:500" json:"icon"`
	Category    string    `gorm:"size:50;index" json:"category"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	IsVisible   bool      `gorm:"default:true" json:"is_visible"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AutoMigrate 自动迁移所有表
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Profile{},
		&Admin{},
		&AccessLog{},
		&Category{},
		&Tag{},
		&Article{},
		&Tool{},
	)
}
