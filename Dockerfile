# ===== 阶段1: 构建前端 =====
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# 设置淘宝 npm 镜像源
RUN npm config set registry https://registry.npmmirror.com

# 先复制依赖文件，利用 Docker 缓存
COPY frontend/package*.json ./

# 安装依赖（变化少，优先利用缓存）
RUN npm ci --fetch-retries=5 --fetch-retry-mintimeout=20000 --fetch-retry-maxtimeout=120000

# 复制源码并构建（变化多，缓存失效时重新构建）
COPY frontend/ ./ 

# 构建前端
USER root
RUN chmod +x node_modules/.bin/*
RUN npm run build

# ===== 阶段2: 最终镜像 =====
FROM python:3.11-slim

WORKDIR /app

# 替换国内源、安装 nginx 和 curl、清理缓存（合并为单层）
RUN sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list.d/debian.sources && \
    sed -i 's/security.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list.d/debian.sources && \
    apt-get update && \
    apt-get install -y --no-install-recommends nginx curl && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# 先复制 requirements.txt 并安装依赖（利用缓存）
COPY backend/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple

# 复制后端代码
COPY backend/ ./backend/

# 从构建阶段复制前端构建产物
COPY --from=frontend-builder /app/frontend/dist /usr/share/nginx/html

# 复制 nginx 配置并删除默认站点
COPY nginx.conf /etc/nginx/conf.d/default.conf
RUN rm -f /etc/nginx/sites-enabled/default

# 复制启动脚本
COPY start.sh .
RUN chmod +x start.sh

# 创建目录
RUN mkdir -p /app/notes /app/backend/data

# 暴露端口
EXPOSE 80

# 设置环境变量
ENV PYTHONUNBUFFERED=1
ENV NOTES_DIR=/app/notes
ENV DATA_DIR=/app/backend/data

# 启动
ENTRYPOINT ["./start.sh"]

# 使用说明:
# 构建镜像: docker build -t mdxy .
# 运行容器: docker run -d --name mdxy-app -p 80:80 \
#   -v $(pwd)/notes:/app/notes \
#   -v $(pwd)/backend/data:/app/backend/data \
#   -e PYTHONUNBUFFERED=1 \
#   -e NOTES_DIR=/app/notes \
#   -e DATA_DIR=/app/backend/data \
#   --restart unless-stopped mdxy:latest
# 查看日志: docker logs -f mdxy-app
# 停止容器: docker stop mdxy-app && docker rm mdxy-app