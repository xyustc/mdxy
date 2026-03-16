# 个人网站 v2.0（Go + Vue）

当前仓库为 Go 后端 + Vue 3 前端的单后端架构。

## 技术栈

- 后端：Go 1.21+、Gin、GORM、SQLite、JWT
- 前端：Vue 3、TypeScript、Vite 5

## 目录结构

```text
mdxy/
├── backend-go/        # Go 后端
├── frontend/          # Vue 前端
├── content/           # 内容目录（notes/articles）
├── deploy/            # 部署脚本
└── docker-compose.yml
```

## 快速开始

### 本地开发

```bash
./dev.sh start
```

默认会输出本机与局域网访问地址；如需显式指定绑定地址可使用：

```bash
FRONTEND_HOST=0.0.0.0 ./dev.sh start
```

或分别启动：

```bash
cd backend-go && go run cmd/server/main.go
cd frontend && npm install && npm run dev
```

### Docker 部署

```bash
docker-compose up -d
```

### 非 Docker 远端部署（HTTPS）

```bash
cp .env.prod.example .env.prod
bash deploy/non-docker-https-deploy.sh up
```

管理员改密（一键）：

```bash
bash deploy/reset-admin-password.sh
```

## 常用文档

- 详细项目说明与部署：`README-v2.md`
- Go 后端说明：`backend-go/README.md`
