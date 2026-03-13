# Go 后端

基于 Gin + GORM 的 RESTful API 服务。

## 技术栈

- Gin Web Framework
- GORM (SQLite)
- JWT 认证
- Viper 配置管理

## 目录结构

```
backend-go/
├── cmd/server/          # 程序入口
├── internal/            # 私有代码
│   ├── config/          # 配置管理
│   ├── model/           # 数据模型
│   ├── handler/         # HTTP 处理器
│   ├── service/         # 业务逻辑
│   ├── repository/      # 数据访问
│   ├── middleware/      # 中间件
│   ├── router/          # 路由
│   ├── database/        # 数据库初始化
│   └── pkg/             # 内部工具包
└── config.yaml          # 配置文件
```

## 快速开始

### 安装依赖

```bash
go mod download
```

### 运行

```bash
go run cmd/server/main.go
```

### 构建

```bash
go build -o bin/server cmd/server/main.go
```

## 迁移 Python 访问数据

如果线上之前运行的是 Python 后端，历史访问数据保存在旧的 `analytics.db`，可以先迁移再切换到 Go：

```bash
go run ./cmd/migrate-python-db \
  -source /path/to/backend/data/analytics.db \
  -target /path/to/backend-go/data/mdxy.db
```

这个命令会自动备份源库和目标库，并把 Python 的 `access_logs` 合并进 Go 数据库，补齐 Go 统计需要的 `country` / `region` 字段。旧记录里 `visitor_id` 为空时，Go 侧统计会自动回退到 `ip_address`，所以历史 UV 不会直接丢失。

### 推荐的远端切换脚本

如果远端是按仓库根目录的 `docker-compose.yml` 部署 Go 后端，可以直接使用脚本完成“备份 -> 迁移 -> 重启 -> 健康检查”：

```bash
bash backend-go/scripts/cutover-from-python.sh \
  --source-db /path/to/python/backend/data/analytics.db \
  --target-db /path/to/repo/backend-go/data/mdxy.db
```

脚本会在 `backend-go/data/cutover-backups/cutover-时间戳/` 下保留：

- Python 源库备份
- 切换前的 Go 目标库快照
- 迁移命令自己的备份产物
- 一个 `rollback.env`，供回滚脚本直接使用

如果需要恢复到切换前的 Go 数据库快照，可执行：

```bash
bash backend-go/scripts/rollback-go-cutover.sh \
  --cutover-dir /path/to/repo/backend-go/data/cutover-backups/cutover-YYYYmmdd-HHMMSS
```

说明：

- 这个回滚脚本恢复的是“切换前 Go 数据库快照”，适合 Go 切换失败后快速回到旧的 Go 数据状态。
- 如果你要回退到旧 Python 服务本身，最稳妥的方式仍然是保留旧 Python 容器/进程，不删镜像和数据卷，在流量层切回旧服务。

## 配置

通过 `config.yaml` 或环境变量配置：

- `SERVER_PORT`: 服务端口（默认 8080）
- `SERVER_MODE`: 运行模式 debug/release
- `DATABASE_PATH`: 数据库文件路径
- `JWT_SECRET`: JWT 密钥
- `CONTENT_NOTES_DIR`: 笔记目录
- `CONTENT_ARTICLES_DIR`: 文章目录

环境变量可以直接使用下划线形式，例如：

- `DATABASE_PATH=/app/data/mdxy.db`
- `CONTENT_NOTES_DIR=/app/content/notes`
- `CONTENT_ARTICLES_DIR=/app/content/articles`

## API 接口

### 公开接口

- `GET /api/v1/profile` - 获取个人信息
- `GET /api/v1/notes/tree` - 获取笔记目录树
- `GET /api/v1/notes/search?q=keyword` - 搜索笔记
- `GET /api/v1/notes/*path` - 获取笔记内容

### 管理员接口

- `POST /api/v1/admin/login` - 登录
- `PUT /api/v1/admin/profile` - 更新个人信息（需认证）

## 默认账号

- 用户名: `admin`
- 密码: `admin123`

**生产环境请立即修改！**
