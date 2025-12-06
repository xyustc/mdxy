# Markdown 笔记系统 (mdxy)

一个基于 Vue 3 + FastAPI 的在线 Markdown 笔记查看系统，支持目录浏览、全文搜索和代码高亮等功能。

## 功能特性

- 📁 目录结构浏览：清晰展示笔记文件夹结构
- 🔍 全文搜索：快速查找笔记内容
- ✨ 代码高亮：支持多种编程语言的语法高亮
- 📱 响应式设计：适配不同屏幕尺寸
- 📋 自动生成目录：自动提取 Markdown 标题生成文章目录

## 技术栈

### 前端
- Vue 3 (Composition API)
- Vue Router
- Axios
- Markdown-it
- Highlight.js
- Vite 构建工具

### 后端
- FastAPI (Python)
- Uvicorn (ASGI 服务器)
- Python-multipart (文件处理)

## 项目结构

```
.
├── backend/                 # 后端代码
│   ├── routers/             # API 路由
│   ├── services/            # 业务逻辑
│   ├── config.py            # 配置文件
│   ├── main.py              # 应用入口
│   └── requirements.txt     # Python 依赖
├── frontend/                # 前端代码
│   ├── src/
│   │   ├── api/             # API 接口封装
│   │   ├── components/      # Vue 组件
│   │   ├── router/          # 路由配置
│   │   ├── views/           # 页面视图
│   │   ├── App.vue          # 根组件
│   │   └── main.js          # 应用入口
│   ├── index.html           # HTML 模板
│   ├── package.json         # npm 依赖
│   └── vite.config.js       # Vite 配置
└── notes/                   # 笔记文件目录 (运行时创建)
```

## 快速开始

### 环境要求

- Node.js >= 16
- Python >= 3.8
- npm 或 yarn

### 后端启动

1. 进入后端目录：
   ```bash
   cd backend
   ```

2. 安装 Python 依赖：
   ```bash
   pip install -r requirements.txt
   ```

3. 启动后端服务：
   ```bash
   python main.py
   ```
   
   默认运行在 `http://localhost:8000`

### 前端启动

1. 进入前端目录：
   ```bash
   cd frontend
   ```

2. 安装依赖：
   ```bash
   npm install
   ```

3. 启动开发服务器：
   ```bash
   npm run dev
   ```
   
   默认运行在 `http://localhost:5173`

## 使用说明

1. 将您的 Markdown 笔记文件放入 `notes` 目录中
2. 在浏览器中访问 `http://localhost:5173`（开发环境）或 `http://localhost`（Docker 部署）
3. 通过左侧目录浏览笔记
4. 使用顶部搜索框进行全文搜索

## API 接口

- `GET /api/notes` - 获取笔记目录树
- `GET /api/notes/{path}` - 获取指定笔记内容
- `GET /api/notes/search?q={keyword}` - 搜索笔记

## 开发指南

### 添加新功能

1. 前端组件位于 `frontend/src/components/`
2. 页面视图位于 `frontend/src/views/`
3. API 接口封装在 `frontend/src/api/`
4. 后端路由位于 `backend/routers/`
5. 后端业务逻辑位于 `backend/services/`

### 构建部署

#### 开发环境

前端构建：
```bash
cd frontend
npm run build
```

构建后的静态文件位于 `frontend/dist/` 目录。

#### 生产环境（Docker）

推荐使用 Docker 部署到生产环境（如阿里云服务器）。

## Docker 部署

### 前置要求

- Docker >= 20.10

### 使用 Docker 命令

```bash
# 构建镜像
docker build -t mdxy .

# 运行容器
docker run -d --name mdxy-app -p 80:80 -v /root/code/mdxy/notes:/app/notes -e PYTHONUNBUFFERED=1 -e NOTES_DIR=/app/notes --restart unless-stopped mdxy:latest

# 查看日志
docker logs -f mdxy-app

# 停止并删除容器
docker stop mdxy-app && docker rm mdxy-app
```

### 阿里云服务器部署步骤

1. 登录阿里云服务器

2. 安装 Docker：
   ```bash
   # 安装 Docker
   curl -fsSL https://get.docker.com | bash -s docker --mirror Aliyun
   
   # 启动 Docker
   sudo systemctl start docker
   sudo systemctl enable docker
   ```

3. 上传项目到服务器：
   ```bash
   # 在本地打包
   tar -czf mdxy.tar.gz --exclude=node_modules --exclude=dist --exclude=.git .
   
   # 上传到服务器
   scp mdxy.tar.gz user@your-server-ip:/home/user/
   
   # 在服务器上解压
   tar -xzf mdxy.tar.gz -C /home/user/mdxy
   ```

4. 创建笔记目录并启动：
   ```bash
   cd /home/user/mdxy
   mkdir -p notes
   docker build -t mdxy .
   docker run -d --name mdxy-app -p 80:80 -v /home/user/mdxy/notes:/app/notes -e PYTHONUNBUFFERED=1 -e NOTES_DIR=/app/notes --restart unless-stopped mdxy:latest
   ```

5. 配置防火墙开放 80 端口：
   ```bash
   # 阿里云需要在控制台安全组规则中开放 80 端口
   # 本地防火墙也需要开放
   sudo firewall-cmd --permanent --add-port=80/tcp
   sudo firewall-cmd --reload
   ```

6. 访问应用：
   - 浏览器打开：`http://your-server-ip`

### 更新部署

```bash
# 拉取最新代码（如果使用 git）
git pull

# 重新构建并启动
docker build -t mdxy .
docker stop mdxy-app && docker rm mdxy-app
docker run -d --name mdxy-app -p 80:80 -v /root/code/mdxy/notes:/app/notes -e PYTHONUNBUFFERED=1 -e NOTES_DIR=/app/notes --restart unless-stopped mdxy:latest

# 清理旧镜像
docker image prune -f
```

### 注意事项

- 笔记目录 `notes/` 通过卷挂载，数据会持久化保存
- 默认监听 80 端口，如需更改可修改 Dockerfile 或运行命令
- 生产环境建议配置 HTTPS（使用 Nginx 反向代理 + Let's Encrypt）
- 日志会输出到 Docker 容器日志中，使用 `docker logs` 查看

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目。