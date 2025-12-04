#!/bin/bash
set -e

# 启动后端服务（后台运行）
cd /app/backend
python -m uvicorn main:app --host 0.0.0.0 --port 8000 &
BACKEND_PID=$!

# 等待后端启动（带健康检查）
echo "等待后端服务启动..."
for i in {1..30}; do
    if curl -s http://127.0.0.1:8000/health > /dev/null 2>&1; then
        echo "后端服务已就绪"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "后端服务启动超时"
        exit 1
    fi
    sleep 1
done

# 启动 nginx（前台运行）
nginx -g "daemon off;"
