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
- ✅ 八股笔记浏览系统（目录树、搜索、内容）
- ✅ 后台管理面板（仪表盘、笔记管理、工具管理、数据统计、个人信息编辑）
- ✅ JWT 认证
- ✅ 暗色模式支持
- ✅ 响应式设计 / 移动端适配
- ✅ Docker 部署
- ✅ Alibaba Cloud 一键 HTTPS 部署（Caddy）
- ✅ 非 Docker 一键 HTTPS 部署（Nginx + Certbot）
- ✅ 访问数据统计（趋势、地域、设备、浏览器分布）
- ✅ 全局搜索（笔记 + 工具）
- ✅ 工具箱模块（分类展示、管理 CRUD）
- ✅ Claude Code 速查表（自动同步 + 发布）
- ✅ ID 证件照裁剪 + AI 背景移除
- ✅ EtC 图片混淆器（防截图保护）
- ✅ 游戏模块（2048、贪吃蛇、Stark Shapes 3D 手势粒子）
- ✅ 管理员一键改密脚本
- ✅ 内容水印保护
- ✅ 移动端手势锁定与布局稳定

### Phase 2 (计划中)

- 📝 技术博客系统（支持 Front-matter）
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
│   │   ├── router/      # 路由
│   │   └── pkg/         # 工具包（JWT/响应/地理/水印）
│   └── Dockerfile
│
├── frontend/            # Vue 3 前端
│   ├── src/
│   │   ├── api/         # API 接口
│   │   ├── layouts/     # 布局组件
│   │   ├── views/       # 页面组件
│   │   │   ├── client/  # 公开页面
│   │   │   └── admin/   # 管理后台
│   │   ├── components/  # 公共组件
│   │   ├── stores/      # Pinia stores
│   │   ├── router/      # 路由配置
│   │   ├── styles/      # CSS tokens / themes
│   │   └── utils/       # 工具函数
│   └── Dockerfile
│
├── content/             # 内容目录
│   ├── notes/           # 八股笔记
│   └── articles/        # 技术文章（预留）
│
├── deploy/              # 部署脚本
│   ├── alicloud-deploy.sh
│   ├── non-docker-https-deploy.sh
│   └── reset-admin-password.sh
│
├── dev.sh               # 本地开发管理脚本
├── docker-compose.yml   # Docker 编排
└── docker-compose.prod.yml  # 生产环境编排
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

如需让局域网设备访问开发前端，可使用：

```bash
FRONTEND_HOST=0.0.0.0 ./dev.sh start
```

## API 文档

### 公开接口

- `GET /api/v1/profile` - 获取个人信息
- `GET /api/v1/home` - 获取首页数据（笔记+工具）
- `GET /api/v1/notes/tree` - 获取笔记目录树
- `GET /api/v1/notes/search?q=keyword` - 搜索笔记
- `GET /api/v1/notes/content/*path` - 获取笔记内容
- `GET /api/v1/tools` - 获取工具列表
- `GET /api/v1/tools/categories` - 获取工具分类
- `GET /api/v1/search` - 全局搜索（笔记+工具）
- `GET /api/v1/search/popular` - 热门搜索
- `GET /api/v1/cheat-sheets/:slug` - 获取已发布速查表
- `POST /api/v1/games/score` - 提交游戏分数
- `GET /api/v1/games/score` - 查询最佳分数
- `GET /api/v1/games/leaderboard` - 排行榜

### 管理员接口（需 JWT 认证）

- `POST /api/v1/admin/login` - 登录
- `PUT /api/v1/admin/profile` - 更新个人信息
- **笔记管理**:
  - `GET /api/v1/admin/notes/tree` - 目录树
  - `GET /api/v1/admin/notes/content/*path` - 内容读取
  - `POST /api/v1/admin/notes/save` - 保存笔记
  - `POST /api/v1/admin/notes/featured` - 置顶/取消置顶
  - `POST /api/v1/admin/notes/directory` - 创建目录
  - `DELETE /api/v1/admin/notes/*path` - 删除笔记
- **工具管理**: CRUD (`GET/POST/PUT/DELETE /api/v1/admin/tools[/:id]`)
- **数据统计**:
  - `GET /api/v1/admin/analytics/overview` - 概览
  - `GET /api/v1/admin/analytics/trends` - 趋势
  - `GET /api/v1/admin/analytics/popular` - 热门页面
  - `GET /api/v1/admin/analytics/devices` - 设备统计
  - `GET /api/v1/admin/analytics/browsers` - 浏览器统计
  - `GET /api/v1/admin/analytics/geo` -地域分布
- **速查表管理**:
  - `GET /api/v1/admin/cheat-sheets/:slug/snapshots` - 快照列表
  - `POST /api/v1/admin/cheat-sheets/:slug/sync` - 同步
  - `POST /api/v1/admin/cheat-sheets/:slug/publish/:id` - 发布

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

### Alibaba Cloud Linux 生产部署（推荐：Caddy 自动 HTTPS）

#### 1) 服务器准备

Alibaba Cloud Linux 3：

```bash
sudo dnf install -y docker docker-compose-plugin git
sudo systemctl enable --now docker
```

Alibaba Cloud Linux 2（如需）：

```bash
sudo yum install -y docker docker-compose-plugin git
sudo systemctl enable --now docker
```

同时请在阿里云安全组放行 `80`、`443`（如需直连后端再放行 `8080`）。

#### 2) 拉取代码并切换分支

```bash
git clone https://github.com/xyustc/mdxy.git
cd mdxy
git checkout Copilot/mobile-ui-guideline-optimizations
```

#### 3) 配置生产环境变量

```bash
cp .env.prod.example .env.prod
# 编辑 .env.prod，至少填写：
# SITE_DOMAIN=你的域名
# ACME_EMAIL=你的邮箱
# JWT_SECRET=高强度随机字符串
```

#### 4) 一键启动（含 HTTPS）

```bash
bash deploy/alicloud-deploy.sh up
```

Caddy 会自动申请并续期证书（Let's Encrypt）。

#### 5) 日常运维命令

```bash
bash deploy/alicloud-deploy.sh status
bash deploy/alicloud-deploy.sh logs
bash deploy/alicloud-deploy.sh restart
bash deploy/alicloud-deploy.sh down
```

#### 6) 更新上线

```bash
git pull
bash deploy/alicloud-deploy.sh up
```

#### 7) 常见问题

- 如果证书签发失败，先检查：
  - 域名 A 记录是否已解析到服务器公网 IP
  - 安全组和系统防火墙是否开放 80/443
  - 80/443 是否被其他进程占用
- 中国大陆网络环境偶发 Let's Encrypt 不稳定时，可在 `deploy/Caddyfile` 切换到 ZeroSSL（文件内已给注释示例）。

### 非 Docker 一键部署（Nginx + systemd + Certbot）

如果你已经在服务器上拉取代码，并希望直接以本机进程部署（不使用 Docker），可以使用：

```bash
cp .env.prod.example .env.prod
# 编辑 .env.prod，至少设置：
# SITE_DOMAIN=你的域名
# ACME_EMAIL=你的邮箱
# JWT_SECRET=高强度随机字符串

bash deploy/non-docker-https-deploy.sh up
```

脚本会自动完成：

- 构建后端（Go）与前端（Vite）
- 发布前端静态文件到 `/var/www/mdxy`
- 写入并启动 `mdxy-backend` systemd 服务
- 写入 Nginx 反代配置（`/api -> 127.0.0.1:8080`）
- 使用 Certbot 申请/配置 HTTPS 证书并开启 80->443 跳转

常用命令：

```bash
bash deploy/non-docker-https-deploy.sh build
bash deploy/non-docker-https-deploy.sh start
bash deploy/non-docker-https-deploy.sh stop
bash deploy/non-docker-https-deploy.sh health
bash deploy/non-docker-https-deploy.sh status
bash deploy/non-docker-https-deploy.sh restart
bash deploy/non-docker-https-deploy.sh logs
bash deploy/non-docker-https-deploy.sh renew
```

可选参数（示例）：

```bash
bash deploy/non-docker-https-deploy.sh up --npm-ci auto --node-options "--max-old-space-size=1024"

# 默认已是降级构建（跳过 vue-tsc，直接 vite build）
# 如需完整构建（vue-tsc && vite build）
bash deploy/non-docker-https-deploy.sh up --full-build

# 如需同时申请 xingyu.top + www.xingyu.top 证书
bash deploy/non-docker-https-deploy.sh up --with-www

# 强制重新构建（忽略缓存）
bash deploy/non-docker-https-deploy.sh up --force-rebuild

# 关闭缓存（每次都重新构建）
bash deploy/non-docker-https-deploy.sh up --no-cache

# stop 时一并停止 nginx
bash deploy/non-docker-https-deploy.sh stop --stop-nginx
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

可直接使用一键脚本（推荐）：

```bash
bash deploy/reset-admin-password.sh
```

脚本会交互式输入新密码，并自动更新 SQLite 中 `admins` 表的密码哈希。

如需指定数据库路径：

```bash
bash deploy/reset-admin-password.sh --db-path /var/lib/mdxy/mdxy.db
```

### 2. 如何添加笔记？

将 Markdown 文件放入 `content/notes/` 目录，系统会自动扫描并展示。

### 3. 如何自定义主题？

修改 `frontend/src/styles/variables.css` 中的 CSS 变量。

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！
