# ===== 阶段1: 构建前端 =====
FROM node:20-alpine AS frontend-builder

# 设置 npm 镜像
RUN npm config set registry https://registry.npmmirror.com

WORKDIR /app/frontend

# 复制 package 文件
COPY frontend/package*.json ./

# 安装依赖
RUN npm ci

# 复制前端源码
COPY frontend/ ./

# 构建前端
RUN npm run build


# ===== 阶段2: 最终镜像 =====
FROM python:3.11-slim

WORKDIR /app

# 安装 nginx 和 curl
RUN apt-get update && apt-get install -y nginx curl && rm -rf /var/lib/apt/lists/*

# 复制后端代码
COPY backend/ ./backend/

# 安装 Python 依赖
RUN pip install --no-cache-dir -r backend/requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple

# 从构建阶段复制前端构建产物
COPY --from=frontend-builder /app/frontend/dist /usr/share/nginx/html

# 复制 nginx 配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 复制启动脚本
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

# 创建笔记目录
RUN mkdir -p /app/notes

# 暴露端口
EXPOSE 80

# 设置环境变量
ENV PYTHONUNBUFFERED=1
ENV NOTES_DIR=/app/notes

# 启动
ENTRYPOINT ["/docker-entrypoint.sh"]
