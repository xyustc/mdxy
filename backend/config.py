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
