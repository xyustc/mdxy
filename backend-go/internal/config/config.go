package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Content  ContentConfig
}

type ServerConfig struct {
	Port string
	Mode string // debug, release
}

type DatabaseConfig struct {
	Path string
}

type JWTConfig struct {
	Secret     string
	ExpireHours int
}

type ContentConfig struct {
	NotesDir    string
	ArticlesDir string
}

var AppConfig *Config

// Load 加载配置
func Load() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// 设置默认值
	setDefaults()

	// 环境变量优先
	viper.AutomaticEnv()

	// 读取配置文件（可选）
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
		log.Println("配置文件未找到，使用默认配置和环境变量")
	}

	AppConfig = &Config{
		Server: ServerConfig{
			Port: viper.GetString("server.port"),
			Mode: viper.GetString("server.mode"),
		},
		Database: DatabaseConfig{
			Path: viper.GetString("database.path"),
		},
		JWT: JWTConfig{
			Secret:      viper.GetString("jwt.secret"),
			ExpireHours: viper.GetInt("jwt.expire_hours"),
		},
		Content: ContentConfig{
			NotesDir:    viper.GetString("content.notes_dir"),
			ArticlesDir: viper.GetString("content.articles_dir"),
		},
	}

	// 确保目录存在
	ensureDir(AppConfig.Content.NotesDir)
	ensureDir(AppConfig.Content.ArticlesDir)
	ensureDir(filepath.Dir(AppConfig.Database.Path))

	return nil
}

func setDefaults() {
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.path", "./data/mdxy.db")
	viper.SetDefault("jwt.secret", "change-this-secret-in-production")
	viper.SetDefault("jwt.expire_hours", 720) // 30天
	viper.SetDefault("content.notes_dir", "../content/notes")
	viper.SetDefault("content.articles_dir", "../content/articles")
}

func ensureDir(path string) {
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Printf("创建目录失败 %s: %v", path, err)
	}
}
