# 个人网站 v2.0

基于 Go + Vue 3 + TypeScript 的现代化个人网站系统。

## 技术栈

### 后端
- Go 1.21+
- Gin Web Framework
- GORM (SQLite)
- JWT 认证
- Viper 配置管理

### 前端
- Vue 3 + TypeScript
- Vite 5
- Vue Router 4
- Pinia 状态管理
- Naive UI (前台)
- Element Plus (后台)
- Markdown-it + Highlight.js

## 功能特性

### Phase 1 (当前版本)

- ✅ 个人信息展示与管理
- ✅ 八股笔记浏览系统
- ✅ 后台管理面板
- ✅ JWT 认证
- ✅ 暗色模式支持
- ✅ 响应式设计
- ✅ Docker 部署

### Phase 2 (计划中)

- 📝 技术博客系统（支持 Front-matter）
- 🔍 文章分类与标签
- 📊 阅读量统计
- 🛠️ 工具箱模块
- 💬 评论系统集成（Giscus）

## 项目结构

```
mdxy/
├── backend-go/          # Go 后端
│   ├── cmd/server/      # 程序入口
│   ├── internal/        # 私有代码
│   │   ├── config/      # 配置管理
│   │   ├── model/       # 数据模型
│   │   ├── handler/     # HTTP 处理器
│   │   ├── service/     # 业务逻辑
│   │   ├── repository/  # 数据访问
│   │   ├── middleware/  # 中间件
│   │   └── router/      # 路由
│   └── Dockerfile
│
├── frontend/            # Vue 3 前端
│   ├── src/
│   │   ├── api/         # API 接口
│   │   ├── layouts/     # 布局组件
│   │   ├── views/       # 页面组件
│   │   ├── stores/      # Pinia stores
│   │   └── router/      # 路由配置
│   └── Dockerfile
│
├── content/             # 内容目录
│   ├── notes/           # 八股笔记
│   └── articles/        # 技术文章
│
└── docker-compose.yml   # Docker 编排
```

## 快速开始

### 前置要求

- Docker & Docker Compose
- 或者：Go 1.21+ 和 Node.js 20+

### 使用 Docker 部署（推荐）

1. 克隆仓库
```bash
git clone <your-repo>
cd mdxy
```

2. 配置环境变量（可选）
```bash
cp .env.example .env
# 编辑 .env 文件，修改 JWT_SECRET 等配置
```

3. 启动服务
```bash
docker-compose up -d
```

4. 访问应用
- 前台：http://localhost
- 后台：http://localhost/admin/login
- 默认账号：admin / admin123

### 本地开发

#### 后端开发

```bash
cd backend-go

# 安装依赖
go mod download

# 运行
go run cmd/server/main.go
```

后端将运行在 http://localhost:8080

#### 前端开发

```bash
cd frontend

# 安装依赖
npm install

# 运行开发服务器
npm run dev
```

前端将运行在 http://localhost:5173

## API 文档

### 公开接口

- `GET /api/v1/profile` - 获取个人信息
- `GET /api/v1/notes/tree` - 获取笔记目录树
- `GET /api/v1/notes/search?q=keyword` - 搜索笔记
- `GET /api/v1/notes/*path` - 获取笔记内容

### 管理员接口

- `POST /api/v1/admin/login` - 登录
- `PUT /api/v1/admin/profile` - 更新个人信息（需认证）

## 配置说明

### 后端配置

通过 `backend-go/config.yaml` 或环境变量配置：

```yaml
server:
  port: "8080"
  mode: "release"  # debug / release

database:
  path: "./data/mdxy.db"

jwt:
  secret: "change-this-secret-in-production"
  expire_hours: 720

content:
  notes_dir: "../content/notes"
  articles_dir: "../content/articles"
```

### 前端配置

通过 `frontend/vite.config.ts` 配置代理和构建选项。

## 部署指南

### 生产环境部署

1. 修改 `.env` 文件中的敏感配置
2. 修改默认管理员密码
3. 配置 HTTPS（推荐使用 Nginx + Let's Encrypt）
4. 设置防火墙规则

### 更新部署

```bash
git pull
docker-compose down
docker-compose build
docker-compose up -d
```

## 开发指南

### 添加新的 API 接口

1. 在 `internal/model/` 定义数据模型
2. 在 `internal/repository/` 实现数据访问
3. 在 `internal/service/` 实现业务逻辑
4. 在 `internal/handler/` 实现 HTTP 处理器
5. 在 `internal/router/router.go` 注册路由

### 添加新的前端页面

1. 在 `src/views/` 创建页面组件
2. 在 `src/router/index.ts` 添加路由
3. 在 `src/api/` 添加 API 接口调用

## 常见问题

### 1. 如何修改管理员密码？

默认密码是 `admin123`。首次部署后，请立即通过数据库修改密码哈希，或者在代码中修改 `internal/database/database.go` 中的初始化逻辑。

### 2. 如何添加笔记？

将 Markdown 文件放入 `content/notes/` 目录，系统会自动扫描并展示。

### 3. 如何自定义主题？

修改 `frontend/src/styles/variables.css` 中的 CSS 变量。

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！
