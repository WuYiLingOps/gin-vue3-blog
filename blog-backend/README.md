# 个人博客系统 - 后端

基于 Go + Gin + GORM + PostgreSQL 构建的现代化博客后端系统。

## 技术栈

- **框架**: Gin Web Framework
- **ORM**: GORM
- **数据库**: PostgreSQL
- **认证**: JWT (JSON Web Token)
- **日志**: Zap
- **配置**: Viper
- **容器化**: Docker & Docker Compose

## 功能特性

- ✅ 用户认证与授权（JWT）
- ✅ 文章管理（CRUD、分类、标签）
- ✅ 评论系统（支持多级回复）
- ✅ 分类和标签管理
- ✅ 用户权限控制
- ✅ 文章点赞和浏览统计
- ✅ RESTful API 设计
- ✅ 日志记录
- ✅ CORS 跨域处理
- ✅ 接口限流
- ✅ 数据分页

## 项目结构

```
blog-backend/
├── cmd/
│   └── server/
│       └── main.go          # 应用入口
├── config/
│   ├── config.go           # 配置加载
│   └── config.yml          # 配置文件
├── db/
│   └── database.go         # 数据库连接
├── handler/                # HTTP 处理器
├── middleware/             # 中间件
├── model/                  # 数据模型
├── repository/             # 数据访问层
├── service/                # 业务逻辑层
├── util/                   # 工具函数
├── logger/                 # 日志系统
├── sql/                    # SQL 脚本
├── Dockerfile              # Docker 镜像构建
├── docker-compose.yml      # Docker Compose 配置
├── Makefile               # 构建脚本
├── go.mod                 # Go 模块
└── README.md              # 项目文档
```

## 快速开始

### 前置要求

- Go 1.21+
- PostgreSQL 13+
- Docker & Docker Compose (可选)

### 本地开发

1. **克隆项目**
```bash
git clone <repository-url>
cd blog-backend
```

2. **安装依赖**
```bash
go mod download
```

3. **配置数据库**
```bash
# 创建 PostgreSQL 数据库
createdb blogdb

# 初始化数据库表结构
psql -U postgres -d blogdb -f sql/init.sql
```

4. **配置文件**
```bash
# 复制配置文件
cp config/config.example.yml config/config.yml

# 编辑配置文件，修改数据库连接信息
vim config/config.yml
```

5. **运行应用**
```bash
go run cmd/server/main.go
# 或使用 Makefile
make run
```

服务器将在 http://localhost:8080 启动

### 使用 Docker

1. **启动所有服务**
```bash
docker-compose up -d
```

2. **查看日志**
```bash
docker-compose logs -f backend
```

3. **停止服务**
```bash
docker-compose down
```

## API 文档

### 认证接口

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | `/api/auth/register` | 用户注册 | 否 |
| POST | `/api/auth/login` | 用户登录 | 否 |
| POST | `/api/auth/logout` | 用户登出 | 是 |
| GET | `/api/auth/profile` | 获取用户信息 | 是 |
| PUT | `/api/auth/profile` | 更新用户信息 | 是 |
| PUT | `/api/auth/password` | 修改密码 | 是 |
| POST | `/api/auth/refresh` | 刷新 Token | 是 |

### 文章接口

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/api/posts` | 获取文章列表 | 否 |
| GET | `/api/posts/:id` | 获取文章详情 | 否 |
| POST | `/api/posts` | 创建文章 | 是 |
| PUT | `/api/posts/:id` | 更新文章 | 是 |
| DELETE | `/api/posts/:id` | 删除文章 | 是 |
| POST | `/api/posts/:id/like` | 点赞文章 | 否 |
| GET | `/api/posts/archives` | 获取归档 | 否 |
| GET | `/api/posts/hot` | 获取热门文章 | 否 |
| GET | `/api/posts/recent` | 获取最新文章 | 否 |

### 分类接口

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/api/categories` | 获取分类列表 | 否 |
| GET | `/api/categories/:id` | 获取分类详情 | 否 |
| POST | `/api/categories` | 创建分类 | 管理员 |
| PUT | `/api/categories/:id` | 更新分类 | 管理员 |
| DELETE | `/api/categories/:id` | 删除分类 | 管理员 |

### 标签接口

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/api/tags` | 获取标签列表 | 否 |
| GET | `/api/tags/:id` | 获取标签详情 | 否 |
| GET | `/api/tags/:id/posts` | 获取标签下的文章 | 否 |
| POST | `/api/tags` | 创建标签 | 是 |
| PUT | `/api/tags/:id` | 更新标签 | 是 |
| DELETE | `/api/tags/:id` | 删除标签 | 是 |

### 评论接口

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/api/comments/post/:id` | 获取文章评论 | 否 |
| POST | `/api/comments` | 创建评论 | 是 |
| PUT | `/api/comments/:id` | 更新评论 | 是 |
| DELETE | `/api/comments/:id` | 删除评论 | 是 |

### 管理后台接口

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/api/admin/dashboard` | 仪表盘数据 | 管理员 |
| GET | `/api/admin/users` | 用户列表 | 管理员 |
| PUT | `/api/admin/users/:id/status` | 更新用户状态 | 管理员 |
| GET | `/api/admin/posts` | 所有文章 | 管理员 |
| GET | `/api/admin/comments` | 所有评论 | 管理员 |
| PUT | `/api/admin/comments/:id/status` | 更新评论状态 | 管理员 |

## 请求示例

### 用户注册
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 用户登录
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### 创建文章（需要 Token）
```bash
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "title": "我的第一篇文章",
    "content": "文章内容...",
    "summary": "文章摘要",
    "category_id": 1,
    "tag_ids": [1, 2],
    "status": 1
  }'
```

## 配置说明

配置文件 `config/config.yml`:

```yaml
app:
  port: 8080          # 服务端口

server:
  mode: debug         # 运行模式: debug, release

db:
  host: localhost     # 数据库主机
  port: 5432         # 数据库端口
  user: postgres     # 数据库用户名
  password: password # 数据库密码
  dbname: blogdb     # 数据库名称
  sslmode: disable   # SSL 模式

jwt:
  secret: your-secret-key  # JWT 密钥（生产环境务必修改）
  expire_hours: 72         # Token 过期时间（小时）

log:
  level: info        # 日志级别: debug, info, warn, error
```

## 环境变量

支持通过环境变量覆盖配置：

- `APP_PORT`: 应用端口
- `SERVER_MODE`: 运行模式
- `DB_HOST`: 数据库主机
- `DB_PORT`: 数据库端口
- `DB_USER`: 数据库用户
- `DB_PASSWORD`: 数据库密码
- `DB_NAME`: 数据库名称
- `JWT_SECRET`: JWT 密钥
- `LOG_LEVEL`: 日志级别

## 开发指南

### 添加新功能

1. 在 `model/` 中定义数据模型
2. 在 `repository/` 中实现数据访问层
3. 在 `service/` 中实现业务逻辑
4. 在 `handler/` 中实现 HTTP 处理器
5. 在 `cmd/server/main.go` 中注册路由

### 代码规范

- 遵循 Go 官方代码规范
- 使用 `go fmt` 格式化代码
- 使用 `golangci-lint` 进行代码检查

### 测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
go test -v ./service/...
```

## 部署

### 生产环境配置

1. 修改 `config/config.yml`，设置 `server.mode` 为 `release`
2. 修改 JWT 密钥为强密码
3. 配置生产环境数据库
4. 启用 HTTPS

### Docker 部署

```bash
# 构建镜像
docker build -t blog-backend:latest .

# 运行容器
docker run -d \
  -p 8080:8080 \
  -v ./config:/root/config \
  --name blog-backend \
  blog-backend:latest
```

## 故障排查

### 数据库连接失败

- 检查 PostgreSQL 是否正在运行
- 验证数据库配置是否正确
- 确认数据库用户有足够的权限

### JWT Token 无效

- 确认 Token 未过期
- 检查 JWT 密钥配置
- 验证 Authorization 头格式正确

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

