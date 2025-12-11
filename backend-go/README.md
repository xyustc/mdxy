# Markdown 笔记系统 - Go 后端

这是使用 Go 语言重构的 Markdown 笔记系统后端。

## 项目结构

```
backend-go/
├── config/           # 配置文件
├── database/         # 数据库模型和初始化
├── middleware/       # 中间件
├── routes/           # 路由定义
├── services/         # 业务逻辑服务
├── utils/            # 工具函数
├── go.mod           # Go 模块定义
├── go.sum           # Go 依赖校验和
└── main.go          # 程序入口
```

## 功能特性

### 1. 笔记管理
- 获取笔记目录树
- 读取笔记内容
- 搜索笔记

### 2. 管理功能
- 管理员登录认证
- 笔记文件上传
- 笔记文件删除
- 目录创建
- 文件移动/重命名

### 3. 访问统计
- 访问日志记录
- 统计数据查询
- 访问趋势分析

## API 接口

### 笔记相关接口
- `GET /api/notes` - 获取笔记目录树
- `GET /api/notes/search?q={keyword}` - 搜索笔记
- `GET /api/notes/{path}` - 获取指定笔记内容

### 管理接口
- `POST /api/admin/login` - 管理员登录
- `GET /api/admin/verify` - 验证token
- `POST /api/admin/notes/upload` - 上传笔记
- `DELETE /api/admin/notes/{path}` - 删除笔记
- `POST /api/admin/notes/mkdir` - 创建目录
- `PUT /api/admin/notes/move` - 移动/重命名文件

### 统计接口
- `GET /api/admin/analytics/logs` - 获取访问日志
- `GET /api/admin/analytics/overview` - 获取统计概览

## 配置说明

环境变量：
- `NOTES_DIR` - 笔记文件存储目录
- `ADMIN_PASSWORD` - 管理员密码
- `JWT_SECRET_KEY` - JWT密钥
- `DATA_DIR` - 数据库存储目录

## 部署说明

1. 安装 Go 1.21 或更高版本
2. 克隆项目代码
3. 运行 `go mod tidy` 安装依赖
4. 设置必要的环境变量
5. 运行 `go run main.go` 启动服务

服务默认监听在 `:8000` 端口。