"""应用配置"""
import os
from pathlib import Path

# 项目根目录
BASE_DIR = Path(__file__).resolve().parent.parent

# 笔记存放目录（支持环境变量配置）
NOTES_DIR = Path(os.environ.get("NOTES_DIR", BASE_DIR / "notes"))

# 允许的文件扩展名
ALLOWED_EXTENSIONS = {".md", ".markdown"}

# CORS 配置
CORS_ORIGINS = [
    "http://localhost:5173",  # Vite 默认端口
    "http://localhost:3000",
    "http://127.0.0.1:5173",
    "http://127.0.0.1:3000",
    "http://localhost",       # Docker 部署
    "http://127.0.0.1",
    "*",                      # 允许所有来源（生产环境建议限制）
]

# 管理员配置
ADMIN_PASSWORD = os.environ.get("ADMIN_PASSWORD", "panzai")
JWT_SECRET_KEY = os.environ.get("JWT_SECRET_KEY", "your-secret-key-change-in-production-2024")
JWT_ALGORITHM = "HS256"
JWT_EXPIRE_MINUTES = 60 * 24 * 30  # token有效期30天

# 数据库配置
DATA_DIR = Path(os.environ.get("DATA_DIR", BASE_DIR / "backend" / "data"))
DATA_DIR.mkdir(parents=True, exist_ok=True)
DATABASE_URL = f"sqlite:///{DATA_DIR / 'analytics.db'}"

# 文件上传配置
MAX_UPLOAD_SIZE = 10 * 1024 * 1024  # 10MB
