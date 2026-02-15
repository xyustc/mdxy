package database

import (
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
			PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // admin123
		}
		DB.Create(admin)
		log.Println("✓ 创建默认管理员: admin / admin123")
	}

	// 检查是否已有个人信息
	DB.Model(&model.Profile{}).Count(&count)
	if count == 0 {
		profile := &model.Profile{
			Name:  "Your Name",
			Title: "Software Engineer",
			Bio:   "Welcome to my personal website!",
		}
		DB.Create(profile)
		log.Println("✓ 创建默认个人信息")
	}
}
