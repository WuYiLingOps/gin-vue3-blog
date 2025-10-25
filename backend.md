# 个人博客系统 - 后端需求文档

## 技术栈
- **框架**: Go + Gin
- **数据库**: PostgreSQL + GORM
- **认证**: JWT
- **日志**: Zap
- **配置管理**: Viper
- **跨域处理**: CORS

## 项目结构
```
blog-backend/
├── cmd/
│   └── server/
│       └── main.go          # 应用入口
├── config/
│   ├── config.go           # 配置结构体和加载逻辑
│   └── config.yml          # 配置文件
├── handler/                # HTTP 处理器
│   ├── auth.go            # 认证相关接口
│   ├── post.go            # 文章相关接口
│   ├── category.go        # 分类相关接口
│   ├── tag.go             # 标签相关接口
│   ├── comment.go         # 评论相关接口
│   └── user.go            # 用户相关接口
├── middleware/             # 中间件
│   ├── auth.go            # JWT 认证中间件
│   ├── cors.go            # 跨域中间件
│   ├── logger.go          # 日志中间件
│   └── rate_limit.go      # 限流中间件
├── model/                  # 数据模型
│   └── models.go          # 数据库模型定义
├── repository/             # 数据访问层
│   ├── post.go            # 文章数据访问
│   ├── user.go            # 用户数据访问
│   ├── category.go        # 分类数据访问
│   ├── tag.go             # 标签数据访问
│   └── comment.go         # 评论数据访问
├── service/                # 业务逻辑层
│   ├── auth.go            # 认证业务逻辑
│   ├── post.go            # 文章业务逻辑
│   ├── user.go            # 用户业务逻辑
│   ├── category.go        # 分类业务逻辑
│   ├── tag.go             # 标签业务逻辑
│   └── comment.go         # 评论业务逻辑
├── util/                   # 工具函数
│   ├── jwt.go             # JWT 工具
│   ├── password.go        # 密码加密工具
│   ├── response.go        # 统一响应格式
│   └── validator.go       # 参数验证工具
├── logger/                 # 日志配置
│   └── logger.go          # Zap 日志初始化
├── db/                     # 数据库相关
│   └── database.go        # 数据库连接和初始化
├── sql/                    # SQL 脚本
│   └── init.sql           # 数据库初始化脚本
├── go.mod
└── go.sum
```

## 数据模型设计

### 用户表 (users)
```go
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Username  string    `json:"username" gorm:"uniqueIndex;not null"`
    Email     string    `json:"email" gorm:"uniqueIndex;not null"`
    Password  string    `json:"-" gorm:"not null"`
    Nickname  string    `json:"nickname"`
    Avatar    string    `json:"avatar"`
    Bio       string    `json:"bio"`
    Role      string    `json:"role" gorm:"default:user"` // admin, user
    Status    int       `json:"status" gorm:"default:1"`  // 1:正常 0:禁用
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 文章表 (posts)
```go
type Post struct {
    ID          uint       `json:"id" gorm:"primaryKey"`
    Title       string     `json:"title" gorm:"not null"`
    Content     string     `json:"content" gorm:"type:text"`
    Summary     string     `json:"summary"`
    Cover       string     `json:"cover"`
    Status      int        `json:"status" gorm:"default:1"` // 1:发布 0:草稿 -1:删除
    IsTop       bool       `json:"is_top" gorm:"default:false"`
    ViewCount   int        `json:"view_count" gorm:"default:0"`
    LikeCount   int        `json:"like_count" gorm:"default:0"`
    UserID      uint       `json:"user_id"`
    CategoryID  uint       `json:"category_id"`
    PublishedAt *time.Time `json:"published_at"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
    
    // 关联关系
    User     User     `json:"user" gorm:"foreignKey:UserID"`
    Category Category `json:"category" gorm:"foreignKey:CategoryID"`
    Tags     []Tag    `json:"tags" gorm:"many2many:post_tags;"`
    Comments []Comment `json:"comments,omitempty"`
}
```

### 分类表 (categories)
```go
type Category struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"uniqueIndex;not null"`
    Description string    `json:"description"`
    Color       string    `json:"color"`
    Sort        int       `json:"sort" gorm:"default:0"`
    PostCount   int       `json:"post_count" gorm:"default:0"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### 标签表 (tags)
```go
type Tag struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"uniqueIndex;not null"`
    Color     string    `json:"color"`
    PostCount int       `json:"post_count" gorm:"default:0"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 评论表 (comments)
```go
type Comment struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Content   string    `json:"content" gorm:"not null"`
    PostID    uint      `json:"post_id"`
    UserID    uint      `json:"user_id"`
    ParentID  *uint     `json:"parent_id"` // 父评论ID，用于回复
    Status    int       `json:"status" gorm:"default:1"` // 1:正常 0:隐藏
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    
    // 关联关系
    User     User      `json:"user" gorm:"foreignKey:UserID"`
    Post     Post      `json:"post" gorm:"foreignKey:PostID"`
    Parent   *Comment  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
    Children []Comment `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}
```

## API 接口设计

### 认证相关接口
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/logout` - 用户登出
- `POST /api/auth/refresh` - 刷新 Token
- `GET /api/auth/profile` - 获取用户信息
- `PUT /api/auth/profile` - 更新用户信息
- `PUT /api/auth/password` - 修改密码

### 文章相关接口
- `GET /api/posts` - 获取文章列表（支持分页、搜索、分类筛选）
- `GET /api/posts/:id` - 获取文章详情
- `POST /api/posts` - 创建文章（需要认证）
- `PUT /api/posts/:id` - 更新文章（需要认证）
- `DELETE /api/posts/:id` - 删除文章（需要认证）
- `POST /api/posts/:id/like` - 点赞文章
- `GET /api/posts/archives` - 获取文章归档
- `GET /api/posts/hot` - 获取热门文章
- `GET /api/posts/recent` - 获取最新文章

### 分类相关接口
- `GET /api/categories` - 获取分类列表
- `GET /api/categories/:id` - 获取分类详情
- `POST /api/categories` - 创建分类（需要管理员权限）
- `PUT /api/categories/:id` - 更新分类（需要管理员权限）
- `DELETE /api/categories/:id` - 删除分类（需要管理员权限）

### 标签相关接口
- `GET /api/tags` - 获取标签列表
- `GET /api/tags/:id` - 获取标签详情
- `POST /api/tags` - 创建标签（需要认证）
- `PUT /api/tags/:id` - 更新标签（需要认证）
- `DELETE /api/tags/:id` - 删除标签（需要认证）

### 评论相关接口
- `GET /api/posts/:id/comments` - 获取文章评论列表
- `POST /api/posts/:id/comments` - 添加评论（需要认证）
- `PUT /api/comments/:id` - 更新评论（需要认证）
- `DELETE /api/comments/:id` - 删除评论（需要认证）

### 管理后台接口
- `GET /api/admin/dashboard` - 获取仪表盘数据
- `GET /api/admin/users` - 获取用户列表
- `PUT /api/admin/users/:id/status` - 更新用户状态
- `GET /api/admin/posts` - 获取所有文章（包括草稿）
- `GET /api/admin/comments` - 获取所有评论

## 功能需求

### 1. 用户认证与权限管理
- JWT Token 认证机制
- 用户注册、登录、登出
- 密码加密存储（bcrypt）
- 角色权限控制（管理员/普通用户）
- Token 自动刷新机制

### 2. 文章管理系统
- 文章 CRUD 操作
- Markdown 格式支持
- 文章分类和标签
- 文章状态管理（草稿/发布/删除）
- 文章置顶功能
- 文章浏览量统计
- 文章点赞功能
- 文章搜索功能
- 文章归档功能

### 3. 评论系统
- 评论 CRUD 操作
- 评论回复功能（支持多级回复）
- 评论状态管理
- 评论审核功能

### 4. 分类和标签管理
- 分类 CRUD 操作
- 标签 CRUD 操作
- 分类和标签统计

### 5. 文件上传
- 图片上传功能
- 文件类型和大小限制
- 图片压缩和处理

### 6. 数据统计
- 文章数量统计
- 用户数量统计
- 评论数量统计
- 访问量统计

## 非功能需求

### 1. 性能要求
- 接口响应时间 < 200ms
- 支持并发访问
- 数据库查询优化
- 缓存机制（Redis）

### 2. 安全要求
- SQL 注入防护
- XSS 攻击防护
- CSRF 攻击防护
- 接口限流
- 敏感信息加密

### 3. 可维护性
- 代码规范和注释
- 错误处理和日志记录
- 单元测试覆盖
- API 文档生成

### 4. 可扩展性
- 模块化设计
- 插件化架构
- 微服务架构支持

## 配置管理

### config.yml 示例
```yaml
app:
  port: 8080

server:
  mode: debug # debug, release

db:
  host: localhost
  port: 5432
  user: postgres
  password: password
  dbname: blog
  sslmode: disable

jwt:
  secret: your-secret-key
  expire_hours: 24

log:
  level: info # debug, info, warn, error

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

upload:
  max_size: 10485760 # 10MB
  allowed_types: ["jpg", "jpeg", "png", "gif"]
  upload_path: "./uploads"
```

## 部署要求

### 1. 环境要求
- Go 1.21+
- PostgreSQL 13+
- Redis 6+

### 2. 部署方式
- Docker 容器化部署
- 支持 Docker Compose
- 支持 Kubernetes 部署

### 3. 监控和日志
- 应用性能监控
- 错误日志收集
- 访问日志记录
- 健康检查接口

## 开发规范

### 1. 代码规范
- 遵循 Go 官方代码规范
- 使用 golangci-lint 进行代码检查
- 统一的错误处理机制
- 统一的响应格式

### 2. Git 规范
- 使用 Git Flow 工作流
- 提交信息规范
- 代码审查机制

### 3. 测试规范
- 单元测试覆盖率 > 80%
- 集成测试
- API 测试

## 第三方依赖

### 核心依赖
- `github.com/gin-gonic/gin` - Web 框架
- `gorm.io/gorm` - ORM 框架
- `gorm.io/driver/postgres` - PostgreSQL 驱动
- `github.com/golang-jwt/jwt/v5` - JWT 认证
- `go.uber.org/zap` - 日志库
- `github.com/spf13/viper` - 配置管理
- `golang.org/x/crypto/bcrypt` - 密码加密

### 可选依赖
- `github.com/go-redis/redis/v8` - Redis 客户端
- `github.com/swaggo/gin-swagger` - API 文档生成
- `github.com/stretchr/testify` - 测试框架
- `github.com/gin-contrib/cors` - CORS 中间件
