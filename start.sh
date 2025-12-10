#!/bin/bash
set -e

# 启动后端服务（后台运行）
cd /app/backend

# 使用更可靠的后台启动方式
python -m uvicorn main:app --host 0.0.0.0 --port 8000 > app.log 2>&1 &
BACKEND_PID=$!

# 等待后端服务启动完成
echo "等待后端服务启动..."
for i in {1..30}; do
    if curl -s http://localhost:8000/health > /dev/null; then
        echo "后端服务已启动"
        break
    fi
    echo "等待中... ($i/30)"
    sleep 2
done

# 检查后端服务是否真正启动
if ! kill -0 $BACKEND_PID 2>/dev/null; then
    echo "后端服务启动失败"
    cat app.log
    exit 1
fi

# 启动 nginx（前台运行）
nginx -g "daemon off;"