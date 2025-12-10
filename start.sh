#!/bin/bash
set -e

# 启动后端服务（后台运行）
cd /app/backend
nohup python -m uvicorn main:app --host 0.0.0.0 --port 8000 > app.log 2>&1

# 启动 nginx（前台运行）
nginx -g "daemon off;"
