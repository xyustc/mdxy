  #!/bin/sh

# 检查笔记目录是否存在，不存在则创建
if [ ! -d "$NOTES_DIR" ]; then
    echo "Creating notes directory: $NOTES_DIR"
    mkdir -p "$NOTES_DIR"
fi

# 检查数据目录是否存在，不存在则创建
if [ ! -d "$DATA_DIR" ]; then
    echo "Creating data directory: $DATA_DIR"
    mkdir -p "$DATA_DIR"
fi

echo "Starting services..."

# 启动 Go 后端（后台运行）
echo "Starting Go backend on port 8000..."
./mdxy-backend &

# 等待后端启动
sleep 3

# 测试后端是否正常运行
echo "Testing backend connectivity..."
curl -f http://localhost:8000/health || {
    echo "Backend failed to start properly"
    exit 1
}

echo "Backend is running successfully"

# 启动 Nginx（前台运行）
echo "Starting Nginx..."
exec nginx -g 'daemon off;'