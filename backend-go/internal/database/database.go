package database

import (
	"errors"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/xyu/mdxy/internal/config"
	"github.com/xyu/mdxy/internal/model"
)

var DB *gorm.DB

// Init 初始化数据库
func Init() error {
	cfg := config.AppConfig.Database

	var err error
	DB, err = gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return err
	}

	log.Println("✓ 数据库连接成功")

	// 自动迁移
	if err := model.AutoMigrate(DB); err != nil {
		return err
	}

	log.Println("✓ 数据库表迁移完成")

	// 初始化默认数据
	initDefaultData()

	return nil
}

// initDefaultData 初始化默认数据
func initDefaultData() {
	// 检查是否已有管理员
	var count int64
	DB.Model(&model.Admin{}).Count(&count)
	if count == 0 {
		// 创建默认管理员（密码: admin123）
		// 注意：生产环境应该通过环境变量或命令行设置
		admin := &model.Admin{
			Username:     "admin",
			PasswordHash: "$2a$10$X4qrjUf5JWhyssmqgUdJseNcdCJm9j5PA.EPPnDz9V1bdn/1XzRw2", // admin123
		}
		DB.Create(admin)
		log.Println("✓ 创建默认管理员: admin / admin123")
	}

	// 检查是否已有个人信息
	DB.Model(&model.Profile{}).Count(&count)
	if count == 0 {
		profile := &model.Profile{
			Name:   "Xingyu",
			Title:  "🎓 软件工程硕士在读,💻 AI编程技术爱好者,🚀 创新学习者",
			Bio:    "软件工程硕士在读，热爱编程，专注于学习前沿技术，积极参与项目实践。\n用代码探索世界，用学习拓展未来。\n🎯 专注AI编程研究 · 💡 追求技术创新 · 🌟 热爱开源分享",
			Email:  "xingyujafk@gmail.com",
			GitHub: "https://github.com/xyustc",
			Skills: `{
  "stats": [
    {"label": "硕士", "value": "二年级", "icon": "school", "color": "cyan"},
    {"label": "实习经验", "value": "2段", "icon": "briefcase", "color": "coral"},
    {"label": "编程语言", "value": "4+", "icon": "code", "color": "teal"},
    {"label": "学习热情", "value": "∞", "icon": "heart", "color": "orange"}
  ],
  "categories": [
    {"name": "语言", "icon": "code", "items": ["Java", "Python", "Go"]},
    {"name": "Web服务", "icon": "server", "items": ["SpringBoot", "Flask", "Nginx", "JS"]},
    {"name": "云服务", "icon": "cloud", "items": ["微服务", "Docker", "RPC"]},
    {"name": "数据库", "icon": "database", "items": ["MySQL", "Redis", "Hive"]}
  ],
  "education": [
    {"degree": "软件工程硕士在读", "school": "中国科学技术大学", "period": "2023 - 至今", "detail": "研究方向：后端开发工程 · AI编程提效"},
    {"degree": "软件工程学士", "school": "大连理工大学", "period": "2018 - 2022", "detail": "主修：编译原理 · 数据结构与算法 · 软件工程 · 计算机科学"}
  ],
  "experience": [
    {"company": "腾讯 · AI LAB", "role": "后端开发 - Go、Python", "period": "2024.09 - 2025.03", "detail": "负责大模型应用平台的开发与优化"},
    {"company": "美团 · 数据平台部", "role": "后端开发 - Java、Python、全栈AiCoding探索", "period": "2024.04 - 2024.07", "detail": "参与数据平台后端服务开发，优化数据处理流程和系统性能"}
  ],
  "hobbies": ["AI CODING", "编程工具", "技术社区", "开源贡献"]
}`,
		}
		DB.Create(profile)
		log.Println("✓ 创建默认个人信息")
	}

	ensureDefaultTool(
		"Claude Code 速查表",
		"Claude Code 键盘快捷键、斜杠命令、MCP、CLI 标志与工作流的高密度速查页。",
		"link",
		"/tools/claude-code-cheatsheet",
		"⌨️",
		"AI Coding",
		90,
	)
}

func ensureDefaultTool(name, description, toolType, url, icon, category string, sortOrder int) {
	var tool model.Tool
	err := DB.Where("name = ? OR url = ?", name, url).First(&tool).Error
	if err == nil {
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		DB.Create(&model.Tool{
			Name:        name,
			Description: description,
			Type:        toolType,
			URL:         url,
			Icon:        icon,
			Category:    category,
			SortOrder:   sortOrder,
			IsVisible:   true,
		})
		log.Printf("✓ 创建默认工具: %s", name)
	}
}
