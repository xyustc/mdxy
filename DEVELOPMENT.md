# 开发指南

## 项目概述

这是一个基于 Go + Vue 3 + TypeScript 的个人网站系统，采用前后端分离架构。

## 技术架构

### 后端架构

```
cmd/server/main.go (入口)
    ↓
internal/router (路由注册)
    ↓
internal/handler (HTTP 处理器)
    ↓
internal/service (业务逻辑)
    ↓
internal/repository (数据访问)
    ↓
internal/model (数据模型)
```

### 前端架构

```
main.ts (入口)
    ↓
router (路由配置)
    ↓
layouts (布局组件)
    ↓
views (页面组件)
    ↓
api (接口调用)
    ↓
stores (状态管理)
```

## 本地开发环境搭建

### 1. 安装依赖

#### Go 环境
- Go 1.21 或更高版本
- 安装：https://golang.org/dl/

#### Node.js 环境
- Node.js 20 或更高版本
- 安装：https://nodejs.org/

### 2. 启动后端

```bash
cd backend-go

# 首次运行需要下载依赖
go mod download

# 启动开发服务器
go run cmd/server/main.go
```

后端将运行在 http://localhost:8080

#### 后端开发工具推荐

使用 Air 实现热重载：

```bash
# 安装 Air
go install github.com/cosmtrek/air@latest

# 使用 Air 启动
air
```

### 3. 启动前端

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端将运行在 http://localhost:5173

#### 前端开发工具推荐

- VS Code + Volar 插件
- 启用 TypeScript 类型检查：`npm run type-check`

## 开发流程

### 添加新的 API 接口

#### 1. 定义数据模型

在 `internal/model/models.go` 中添加：

```go
type NewModel struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Name      string    `gorm:"size:100;not null" json:"name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

#### 2. 创建 Repository

在 `internal/repository/` 创建新文件：

```go
type NewModelRepository struct {
    db *gorm.DB
}

func NewNewModelRepository(db *gorm.DB) *NewModelRepository {
    return &NewModelRepository{db: db}
}

func (r *NewModelRepository) Create(model *model.NewModel) error {
    return r.db.Create(model).Error
}
```

#### 3. 创建 Service

在 `internal/service/` 创建新文件：

```go
type NewModelService struct {
    repo *repository.NewModelRepository
}

func NewNewModelService(repo *repository.NewModelRepository) *NewModelService {
    return &NewModelService{repo: repo}
}

func (s *NewModelService) Create(model *model.NewModel) error {
    return s.repo.Create(model)
}
```

#### 4. 创建 Handler

在 `internal/handler/` 创建新文件：

```go
type NewModelHandler struct {
    service *service.NewModelService
}

func NewNewModelHandler(service *service.NewModelService) *NewModelHandler {
    return &NewModelHandler{service: service}
}

func (h *NewModelHandler) Create(c *gin.Context) {
    var req model.NewModel
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "参数错误")
        return
    }

    if err := h.service.Create(&req); err != nil {
        response.InternalServerError(c, "创建失败")
        return
    }

    response.Success(c, req)
}
```

#### 5. 注册路由

在 `internal/router/router.go` 中添加：

```go
newModelRepo := repository.NewNewModelRepository(database.DB)
newModelService := service.NewNewModelService(newModelRepo)
newModelHandler := handler.NewNewModelHandler(newModelService)

v1.POST("/newmodel", newModelHandler.Create)
```

### 添加新的前端页面

#### 1. 定义类型

在 `src/api/types.ts` 中添加：

```typescript
export interface NewModel {
  id: number
  name: string
  created_at: string
  updated_at: string
}
```

#### 2. 创建 API 接口

在 `src/api/` 创建新文件：

```typescript
import request from './request'
import type { NewModel, ApiResponse } from './types'

export const newModelApi = {
  create(data: Partial<NewModel>): Promise<ApiResponse<NewModel>> {
    return request.post('/newmodel', data)
  }
}
```

#### 3. 创建页面组件

在 `src/views/` 创建新文件：

```vue
<template>
  <div class="new-page">
    <h1>New Page</h1>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { newModelApi } from '@/api/newmodel'
</script>

<style scoped>
.new-page {
  padding: 20px;
}
</style>
```

#### 4. 添加路由

在 `src/router/index.ts` 中添加：

```typescript
{
  path: '/new',
  name: 'new-page',
  component: () => import('@/views/NewPage.vue')
}
```

## 数据库管理

### 查看数据库

```bash
cd backend-go/data
sqlite3 mdxy.db

# SQLite 命令
.tables          # 查看所有表
.schema profiles # 查看表结构
SELECT * FROM profiles; # 查询数据
```

### 重置数据库

删除 `backend-go/data/mdxy.db` 文件，重启后端会自动重建。

## 调试技巧

### 后端调试

使用 Delve 调试器：

```bash
# 安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 启动调试
dlv debug cmd/server/main.go
```

### 前端调试

使用浏览器开发者工具：
- Chrome DevTools
- Vue DevTools 扩展

### API 测试

推荐使用：
- Postman
- curl
- HTTPie

示例：

```bash
# 登录
curl -X POST http://localhost:8080/api/v1/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 获取个人信息
curl http://localhost:8080/api/v1/profile
```

## 代码规范

### Go 代码规范

- 使用 `gofmt` 格式化代码
- 遵循 Go 官方代码规范
- 错误处理不能忽略
- 导出的函数和类型必须有注释

### TypeScript 代码规范

- 使用 2 空格缩进
- 优先使用 `const`，避免 `var`
- 所有函数参数和返回值必须有类型注解
- 组件使用 `<script setup lang="ts">`

### Vue 组件规范

- 组件名使用 PascalCase
- Props 必须定义类型
- 使用 Composition API
- 样式使用 scoped

## 常见问题

### 1. Go 依赖下载失败

设置 Go 代理：

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

### 2. 前端依赖安装失败

使用淘宝镜像：

```bash
npm config set registry https://registry.npmmirror.com
```

### 3. 跨域问题

开发环境已配置代理，生产环境通过 Nginx 反向代理解决。

### 4. TypeScript 类型错误

运行类型检查：

```bash
npm run type-check
```

## 性能优化建议

### 后端优化

- 使用数据库索引
- 实现缓存机制（Redis）
- 使用连接池
- 启用 Gzip 压缩

### 前端优化

- 路由懒加载
- 图片懒加载
- 使用 CDN
- 代码分割
- Tree Shaking

## 部署检查清单

- [ ] 修改 JWT_SECRET
- [ ] 修改默认管理员密码
- [ ] 设置 SERVER_MODE=release
- [ ] 配置 HTTPS
- [ ] 设置防火墙规则
- [ ] 配置日志收集
- [ ] 设置自动备份
- [ ] 配置监控告警

## 下一步计划

### 待实现

1. 技术博客系统
   - Front-matter 解析
   - 文章同步接口
   - 分类和标签管理
   - 阅读量统计

2. 评论系统
   - 集成 Giscus
   - 评论管理

### 长期优化

1. CI/CD
   - GitHub Actions
   - 自动化测试
   - 自动部署

### 已完成（原 Phase 2/Phase 3）

- 工具箱模块（CRUD + 分类 + 前台展示）✅
- 全局搜索（笔记 + 工具搜索 + 热门）✅
- 数据统计（访问量趋势、设备/浏览器/地域分布）✅
- Claude Code 速查表（自动同步 + 发布）✅
- ID 证件照裁剪 + AI 背景移除 ✅
- EtC 图片混淆器 ✅
- 游戏模块（2048、贪吃蛇、Stark Shapes）✅
- 非 Docker + Alibaba Cloud HTTPS 部署 ✅
- 管理员一键改密 ✅

## 参考资源

- [Go 官方文档](https://golang.org/doc/)
- [Gin 文档](https://gin-gonic.com/docs/)
- [GORM 文档](https://gorm.io/docs/)
- [Vue 3 文档](https://vuejs.org/)
- [TypeScript 文档](https://www.typescriptlang.org/docs/)
- [Naive UI 文档](https://www.naiveui.com/)
- [Element Plus 文档](https://element-plus.org/)
