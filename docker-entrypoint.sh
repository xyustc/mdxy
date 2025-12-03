#!/bin/bash

# 启动后端服务（后台运行）
cd /app/backend
python -m uvicorn main:app --host 0.0.0.0 --port 8000 &

# 等待后端启动
sleep 2

# 启动 nginx（前台运行）
nginx -g "daemon off;"
