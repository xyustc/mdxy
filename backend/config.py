"""应用配置"""
import os
from pathlib import Path

# 项目根目录
BASE_DIR = Path(__file__).resolve().parent.parent

# 笔记存放目录
NOTES_DIR = BASE_DIR / "notes"

# 允许的文件扩展名
ALLOWED_EXTENSIONS = {".md", ".markdown"}

# CORS 配置
CORS_ORIGINS = [
    "http://localhost:5173",  # Vite 默认端口
    "http://localhost:3000",
    "http://127.0.0.1:5173",
    "http://127.0.0.1:3000",
]
