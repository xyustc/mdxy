package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server         ServerConfig
	Database       DatabaseConfig
	JWT            JWTConfig
	Content        ContentConfig
	CheatSheetSync CheatSheetSyncConfig
}

type ServerConfig struct {
	Port string
	Mode string // debug, release
}

type DatabaseConfig struct {
	Path string
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

type ContentConfig struct {
	NotesDir    string
	ArticlesDir string
}

type CheatSheetSyncConfig struct {
	Enabled       bool
	IntervalHours int
	SourceURL     string
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
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
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
		CheatSheetSync: CheatSheetSyncConfig{
			Enabled:       viper.GetBool("cheat_sheet_sync.enabled"),
			IntervalHours: viper.GetInt("cheat_sheet_sync.interval_hours"),
			SourceURL:     viper.GetString("cheat_sheet_sync.source_url"),
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
	viper.SetDefault("cheat_sheet_sync.enabled", true)
	viper.SetDefault("cheat_sheet_sync.interval_hours", 24)
	viper.SetDefault("cheat_sheet_sync.source_url", "https://banwagong1.com/claude-code.html")
}

func ensureDir(path string) {
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Printf("创建目录失败 %s: %v", path, err)
	}
}
