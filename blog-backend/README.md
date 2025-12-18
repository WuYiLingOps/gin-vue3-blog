# Blog Backend - Go 后端服务

基于 Go + Gin + PostgreSQL + Redis 的博客系统后端服务。

## 📋 目录

- [快速开始](#快速开始)
- [项目结构](#项目结构)
- [配置说明](#配置说明)
- [功能特性](#功能特性)
- [管理员 IP 豁免功能](#管理员-ip-豁免功能)
- [图片上传存储](#图片上传存储)
- [常见问题排查](#常见问题排查)

## 🚀 快速开始

### 环境要求

- Go >= 1.21
- PostgreSQL >= 15
- Redis

### 安装依赖

```bash
go mod download
```

### 配置数据库

1. 创建 PostgreSQL 数据库
2. 导入数据库结构：
```bash
psql -U postgres -d blogdb -f sql/init.sql
```

### 配置文件

编辑 `config/config-dev.yml`，配置数据库、Redis、邮箱等信息。

### 运行服务

```bash
cd blog-backend
go run cmd/server/main.go
```

服务默认运行在 `http://localhost:8080`

## 📁 项目结构

```
blog-backend/
├── cmd/server/        # 程序入口
├── config/            # 配置文件
├── db/                # 数据库连接
├── handler/           # 请求处理器（Controller层）
├── service/           # 业务逻辑层
├── repository/        # 数据访问层（DAO层）
├── model/             # 数据模型
├── middleware/        # 中间件
├── router/            # 路由配置
├── util/              # 工具函数
└── uploads/           # 上传文件目录
```

## ⚙️ 配置说明

### 配置文件结构

- `config/config.yml` - 环境配置（dev/prod）
- `config/config-dev.yml` - 开发环境配置
- `config/config-prod.yml` - 生产环境配置

### 主要配置项

- **数据库配置**：PostgreSQL 连接信息
- **Redis 配置**：缓存服务配置
- **JWT 配置**：Token 密钥和过期时间
- **邮箱配置**：SMTP 服务配置（用于密码重置）
- **OSS 配置**：阿里云 OSS 配置（可选）
- **安全配置**：管理员 IP 白名单

## ✨ 功能特性

- ✅ 用户认证（注册、登录、JWT）
- ✅ 文章管理（CRUD、分类、标签）
- ✅ 评论系统（嵌套回复）
- ✅ 说说动态
- ✅ 实时聊天室（WebSocket）
- ✅ 文件上传（本地存储/OSS）
- ✅ IP 黑名单和频率限制
- ✅ IP 黑名单管理 API（查看、添加、删除、检查、清理过期）
- ✅ 管理员 IP 豁免（角色豁免 + IP 白名单）
- ✅ 数据统计
- ✅ 系统广播可定向投递：支持 `announcement`（仅公告栏）、`chat`（仅聊天室）、`both`（同时），数据库 `chat_messages.target` 默认为 `announcement`

### 系统广播投递说明

- 接口：`POST /admin/chat/broadcast`
- 请求体示例：
```json
{
  "content": "系统维护通知",
  "priority": 1,
  "target": "announcement" // announcement | chat | both，默认 both
}
```
- 公告列表接口仅返回 `target` 为 `announcement` / `both` 的广播；聊天室历史与实时消息仅展示 `chat` / `both`。

---

## 🚫 IP 黑名单管理

### 功能概述

系统提供了完整的 IP 黑名单管理功能，支持自动封禁和手动封禁两种方式，并提供管理后台 API 接口。

### 封禁类型

1. **自动封禁** (`ban_type = 1`)
   - 由系统根据访问频率自动触发
   - 当 IP 访问频率超过限制时自动添加到黑名单

2. **手动封禁** (`ban_type = 2`)
   - 由管理员通过管理后台手动添加
   - 可设置封禁原因和封禁时长

### API 接口

#### 1. 获取 IP 黑名单列表

```
GET /admin/ip-blacklist?page=1&page_size=20
```

**响应示例**：
```json
{
  "code": 200,
  "data": {
    "list": [
      {
        "id": 1,
        "ip": "192.168.1.100",
        "reason": "恶意访问",
        "ban_type": 2,
        "expire_at": "2025-01-15T10:00:00Z",
        "created_at": "2025-01-14T10:00:00Z",
        "updated_at": "2025-01-14T10:00:00Z"
      }
    ],
    "total": 10,
    "page": 1,
    "page_size": 20
  }
}
```

#### 2. 添加 IP 到黑名单

```
POST /admin/ip-blacklist
```

**请求体**：
```json
{
  "ip": "192.168.1.100",
  "reason": "恶意访问",
  "duration": 24  // 封禁时长（小时），0 表示永久封禁
}
```

#### 3. 删除 IP（解除封禁）

```
DELETE /admin/ip-blacklist/:id
```

#### 4. 检查 IP 状态

```
GET /admin/ip-blacklist/check?ip=192.168.1.100
```

**响应示例**：
```json
{
  "code": 200,
  "data": {
    "banned": true,
    "info": {
      "id": 1,
      "ip": "192.168.1.100",
      "reason": "恶意访问",
      "ban_type": 2,
      "expire_at": "2025-01-15T10:00:00Z"
    }
  }
}
```

#### 5. 清理过期记录

```
POST /admin/ip-blacklist/clean-expired
```

**响应示例**：
```json
{
  "code": 200,
  "message": "清理成功",
  "data": {
    "deleted_count": 5
  }
}
```

### 封禁时长设置

- **duration = 0**：永久封禁（`expire_at = null`）
- **duration > 0**：临时封禁，在指定小时后自动过期

### 自动清理机制

系统会定期自动清理过期的黑名单记录，也可以通过 API 手动触发清理。

### 权限要求

所有 IP 黑名单管理 API 都需要：
- 用户已登录
- 用户角色为管理员（`admin`）

---

## 🔐 管理员 IP 豁免功能

### 功能概述

为了解决管理员账号在后台操作时因访问频率过高而被自动封禁的问题，实现了**双重豁免机制**：

1. **方案一：基于用户角色的豁免** - 已登录的管理员自动豁免
2. **方案二：基于 IP 白名单的豁免** - 配置文件中指定的管理员 IP 豁免

### 工作原理

```
请求到达
    ↓
检查是否是管理员用户（解析 Token）
    ↓ 是 → 直接放行 ✅
    ↓ 否
检查 IP 是否在白名单中
    ↓ 是 → 直接放行 ✅
    ↓ 否
执行黑名单检查
    ↓
执行频率限制检查
```

### 配置方法

编辑 `config/config-dev.yml`：

```yaml
security:
  admin_ip_whitelist:
    - "127.0.0.1"           # 本地回环地址（IPv4）
    - "::1"                 # 本地回环地址（IPv6）
    - "192.168.1.100"       # 精确 IP 地址
    - "192.168.1.0/24"      # CIDR 格式（整个网段）
```

### 支持格式

- **精确 IP**：`192.168.1.100`
- **CIDR 格式**：`192.168.1.0/24`（支持整个网段）
- **IPv6**：`::1`、`2001:db8::/32`

### 豁免条件

满足以下**任一条件**即可豁免：

1. **已登录且角色为管理员**
   - 通过 JWT Token 解析获取用户角色
   - 如果 Token 有效且角色为 `admin`，自动豁免

2. **IP 在配置的白名单中**
   - 从配置文件读取
   - 支持精确匹配和 CIDR 格式
   - 不依赖登录状态

### 注意事项

1. **配置文件更新后需要重启服务**
2. **生产环境配置**：需要在 `config-prod.yml` 中同步配置
3. **安全建议**：白名单 IP 应该谨慎添加，只添加可信的管理员 IP

---

## 📸 图片上传存储

### 本地存储位置

当选择**本地存储**时，图片保存在以下位置：

```
blog-backend/
└── uploads/              # 主上传目录
    ├── avatars/          # 头像目录
    │   └── 20251214051715_bdbcaf98.jpg
    └── 20251214055346_5161a8e6.JPG  # 普通图片
```

### 文件命名规则

格式：`YYYYMMDDHHmmss_随机8位字符.扩展名`

例如：`20251214051715_bdbcaf98.jpg`

### 配置说明

#### 后端配置

**路由配置** (`router/router.go`)：
```go
// 静态文件服务（用于访问上传的文件）
// 使用绝对路径，确保无论从哪个目录运行都能找到 uploads 目录
uploadsPath, _ := filepath.Abs("./uploads")
r.Static("/uploads", uploadsPath)
```

**中间件配置** (`middleware/ip_blacklist.go`)：
```go
// 排除静态文件路径，不进行频率限制和黑名单检查
path := c.Request.URL.Path
if strings.HasPrefix(path, "/uploads") {
    c.Next()
    return
}
```

#### 前端配置

**Vite 代理配置** (`vite.config.ts`)：
```typescript
proxy: {
  '/api': {
    target: 'http://localhost:8080',
    changeOrigin: true
  },
  '/uploads': {  // ✅ 已添加，用于访问上传的图片
    target: 'http://localhost:8080',
    changeOrigin: true
  }
}
```

### 访问方式

#### 开发环境

- **直接访问后端**：`http://localhost:8080/uploads/avatars/xxx.jpg`
- **通过前端代理**：`http://localhost:3000/uploads/avatars/xxx.jpg`

#### 生产环境

如果前后端分离部署：
- 前端需要配置反向代理（Nginx）将 `/uploads` 请求转发到后端
- 或者将 `uploads` 目录通过 Nginx 直接提供静态文件服务

### 存储类型

当前支持三种存储方式：

1. **本地存储（local）**：文件保存在服务器本地 `uploads` 目录
2. **阿里云 OSS 存储（oss）**：文件上传到阿里云对象存储
3. **腾讯云 COS 存储（cos）**：文件上传到腾讯云对象存储

可以在管理后台的「网站设置 → 上传存储配置」中切换存储类型。

#### 阿里云 OSS 配置示例

在 `config/config-dev.yml` 或 `config/config-prod.yml` 中配置 `oss` 节点（示例）：

```yaml
oss:
  endpoint: "oss-cn-hangzhou.aliyuncs.com"
  access_key_id: "your-ak"
  access_key_secret: "your-sk"
  bucket_name: "your-bucket"
  # 可选，自定义访问域名（如接入 CDN）
  # domain: "https://static.example.com"
```

#### 腾讯云 COS 配置示例

在 `config/config-dev.yml` 或 `config/config-prod.yml` 中配置 `cos` 节点（示例）：

```yaml
cos:
  # Bucket 访问地址，格式必须是完整 URL：
  # https://<bucket-name>.cos.<region>.myqcloud.com
  # 示例：你的 COS Bucket 叫 xx-blog-images-xxxx，地域是 ap-guangzhou
  bucket_url: "https://xx-blog-images-xxxx.cos.ap-guangzhou.myqcloud.com"

  # 腾讯云访问密钥（建议只在本地真实 config-dev.yml / config-prod.yml 中填写，不要提交到仓库）
  secret_id: "你的 SecretId"
  secret_key: "你的 SecretKey"

  # 可选：自定义访问域名（例如绑定了 CDN 或自定义域名）
  # 如果配置了 domain，返回的文件 URL 会优先使用这个域名
  # 没有自定义域名时可以先留空或注释掉
  # domain: "https://static.example.com"
```

> 说明：
> - 存储类型的开关（local / oss / cos）保存在数据库的系统设置中，通过管理后台页面修改；
> - OSS / COS 的连接参数通过 `config-dev.yml` / `config-prod.yml` 加载；
> - 如需在不同环境中安全管理敏感信息（如密钥、密码），可在后端项目根目录（`blog-backend/`）创建对应的 `.env.config.<env>` 文件（如 `.env.config.dev` / `.env.config.prod`），覆盖配置文件中的默认值。

---

## 🔧 常见问题排查

### 图片无法查看

#### 问题现象

上传图片后，返回的 URL 正确，但在浏览器中无法查看图片。

#### 已修复的问题

1. ✅ **Vite 代理配置**：已添加 `/uploads` 路径的代理配置
2. ✅ **中间件拦截问题**：已在 `IPBlacklistMiddleware` 中排除 `/uploads` 路径
3. ✅ **路径问题**：已使用绝对路径 `filepath.Abs("./uploads")`

#### 排查步骤

**步骤 1：检查文件是否存在**

```bash
cd blog-backend
dir uploads
dir uploads\avatars
```

**步骤 2：测试直接访问后端**

在浏览器中直接访问：
```
http://localhost:8080/uploads/avatars/文件名.jpg
```

**预期结果**：应该能看到图片

**如果 404**：
- 检查后端服务是否运行在 8080 端口
- 检查文件路径是否正确
- 检查后端控制台是否有错误日志

**步骤 3：测试通过前端代理访问**

在浏览器中访问：
```
http://localhost:3000/uploads/avatars/文件名.jpg
```

**预期结果**：应该能看到图片（通过 Vite 代理）

**如果 404**：
- 检查前端服务是否运行在 3000 端口
- 检查 Vite 配置是否正确
- 重启前端开发服务器

#### 常见错误

**问题 1：浏览器显示 404**

**可能原因**：
1. 文件不存在
2. 路径不正确
3. 静态文件服务未正确配置

**排查方法**：
```bash
# 检查文件是否存在
cd blog-backend
dir uploads\avatars

# 测试直接访问
# 在浏览器访问：http://localhost:8080/uploads/avatars/文件名.jpg
```

**问题 2：浏览器显示 403 Forbidden**

**可能原因**：
1. 文件权限问题
2. 中间件拦截（已修复）

**排查方法**：
```bash
# 检查文件权限（Linux/Mac）
ls -la blog-backend/uploads/avatars/

# Windows 检查文件属性，确保文件可读
```

**问题 3：浏览器显示 429 Too Many Requests**

**已修复**：静态文件请求不再受频率限制

**问题 4：CORS 跨域问题**

**检查**：
- 查看浏览器控制台的 CORS 错误
- 检查 `middleware/cors.go` 配置

#### 检查清单

- [ ] 后端服务正常运行（端口 8080）
- [ ] 前端服务正常运行（端口 3000）
- [ ] `blog-backend/uploads` 目录存在
- [ ] 上传的文件确实保存在 `uploads` 目录
- [ ] Vite 配置包含 `/uploads` 代理
- [ ] 中间件已排除 `/uploads` 路径
- [ ] 静态文件服务使用绝对路径
- [ ] 浏览器控制台无错误
- [ ] 网络请求返回 200 状态码

#### 调试技巧

**1. 查看后端日志**

启动后端时，查看控制台输出：
- 是否有错误信息
- 静态文件服务是否正常启动

**2. 查看浏览器开发者工具**

1. 打开开发者工具（F12）
2. 切换到 Network 标签
3. 刷新页面或访问图片 URL
4. 查看请求详情：
   - **Request URL**：请求的完整 URL
   - **Status Code**：状态码（应该是 200）
   - **Response**：响应内容

**3. 测试 curl 命令**

```bash
# 测试后端直接访问
curl -I http://localhost:8080/uploads/avatars/文件名.jpg

# 应该返回 200 OK
```

### 管理员被误封禁

如果管理员账号因访问频率过高被自动封禁：

1. **检查是否已配置 IP 白名单**
2. **使用管理员账号登录**（已登录管理员自动豁免）
3. **手动解除封禁**：在管理后台的 IP 黑名单管理中删除相关记录

### 其他问题

如果遇到其他问题，请检查：

1. 后端日志输出
2. 数据库连接是否正常
3. Redis 连接是否正常
4. 配置文件是否正确

---

## 📚 相关文档

- [项目根目录 README](../README.md) - 完整项目说明
- [学习指南](./学习指南.md) - Go 后端学习路径（如果存在）

## 📝 更新日志

### 最新更新

- ✅ IP 黑名单管理 API（查看、添加、删除、检查、清理过期）
- ✅ 管理员 IP 豁免功能（双重机制：角色豁免 + IP 白名单）
- ✅ 图片上传存储优化（绝对路径、中间件排除）
- ✅ Vite 代理配置完善

---

**最后更新**：2025-01-XX

