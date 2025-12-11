package config

import (
	"os"
	"path/filepath"
)

var (
	// 笔记存放目录（支持环境变量配置）
	NotesDir = getEnvOrDefault("NOTES_DIR", filepath.Join(getBaseDir(), "..", "notes"))

	// CORS配置
	CORSOrigins = []string{
		"http://localhost:5173", // Vite默认端口
		"http://localhost:3000",
		"http://127.0.0.1:5173",
		"http://127.0.0.1:3000",
		"http://localhost", // Docker部署
		"http://127.0.0.1",
		"*", // 允许所有来源（生产环境建议限制）
	}

	// 管理员配置
	AdminPassword  = getEnvOrDefault("ADMIN_PASSWORD", "panzai")
	JWTSecretKey   = getEnvOrDefault("JWT_SECRET_KEY", "your-secret-key-change-in-production-2024")
	JWTAlgorithm   = "HS256"
	JWTExpireHours = 24 * 30 // token有效期30天

	// 数据库配置
	DataDir     = getEnvOrDefault("DATA_DIR", filepath.Join(getBaseDir(), "data"))
	DatabaseURL = "file:" + filepath.Join(DataDir, "analytics.db")

	// 文件上传配置
	MaxUploadSize = int64(10 * 1024 * 1024) // 10MB

	// 允许的文件扩展名
	AllowedExtensions = map[string]bool{
		".md":       true,
		".markdown": true,
	}
)

// getEnvOrDefault 获取环境变量，如果不存在则使用默认值
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getBaseDir 获取项目根目录
func getBaseDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "."
	}
	return dir
}
