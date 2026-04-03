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

## 配置

通过 `config.yaml` 或环境变量配置：

- `SERVER_PORT`: 服务端口（默认 8080）
- `SERVER_MODE`: 运行模式 debug/release
- `DATABASE_PATH`: 数据库文件路径
- `JWT_SECRET`: JWT 密钥
- `JWT_EXPIRE_HOURS`: Token 过期时间（小时，默认 720）
- `CONTENT_NOTES_DIR`: 笔记目录
- `CONTENT_ARTICLES_DIR`: 文章目录

环境变量可以直接使用下划线形式，例如：

- `DATABASE_PATH=/app/data/mdxy.db`
- `JWT_SECRET=my-secret`
- `CONTENT_NOTES_DIR=/app/notes`

## API 接口

### 公开接口

- `GET /api/v1/profile` - 获取个人信息
- `GET /api/v1/home` - 获取首页数据（笔记+工具聚合）
- `GET /api/v1/notes/tree` - 获取笔记目录树
- `GET /api/v1/notes/search?q=keyword` - 搜索笔记
- `GET /api/v1/notes/content/*path` - 获取笔记内容
- `GET /api/v1/tools` - 获取工具列表
- `GET /api/v1/tools/categories` - 获取工具分类列表
- `GET /api/v1/search` - 全局搜索
- `GET /api/v1/search/popular` - 热门搜索
- `GET /api/v1/cheat-sheets/:slug` - 获取已发布速查表
- `POST /api/v1/games/score` - 提交游戏分数
- `GET /api/v1/games/score` - 获取最佳分数
- `GET /api/v1/games/leaderboard` - 游戏排行榜

### 管理员接口（需 JWT 认证）

- `POST /api/v1/admin/login` - 登录
- `PUT /api/v1/admin/profile` - 更新个人信息
- `GET/POST/PUT/DELETE /api/v1/admin/tools[/:id]` - 工具 CRUD
- `GET /api/v1/admin/notes/tree` - 笔记目录树
- `GET /api/v1/admin/notes/content/*path` - 笔记内容读取
- `POST /api/v1/admin/notes/save` - 保存笔记
- `POST /api/v1/admin/notes/featured` - 置顶设置
- `POST /api/v1/admin/notes/directory` - 创建目录
- `DELETE /api/v1/admin/notes/*path` - 删除笔记
- `GET /api/v1/admin/analytics/*` - 数据统计（概览、趋势、热门页面、设备、浏览器、地域）
- `GET /api/v1/admin/cheat-sheets/:slug/snapshots` - 速查表快照列表
- `POST /api/v1/admin/cheat-sheets/:slug/sync` - 速查表同步
- `POST /api/v1/admin/cheat-sheets/:slug/publish/:id` - 速查表发布

## 默认账号

- 用户名: `admin`
- 密码: `admin123`

**生产环境请立即修改！**
