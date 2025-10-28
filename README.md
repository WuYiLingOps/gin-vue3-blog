# 🎨 个人博客系统

一个基于 Vue 3 + Go 的现代化全栈博客系统，采用前后端分离架构，具有优雅的 UI 设计和完善的功能。

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Vue](https://img.shields.io/badge/Vue-3.3+-green.svg)
![Go](https://img.shields.io/badge/Go-1.21+-blue.svg)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-blue.svg)

## ✨ 特性

### 🎯 核心功能
- 📝 **文章管理** - Markdown 编辑器，支持代码高亮、图片上传
- 🏷️ **分类标签** - 灵活的分类和标签系统
- 💬 **评论系统** - 支持嵌套回复的评论功能
- 💭 **说说动态** - 类似朋友圈的动态发布
- 👤 **用户系统** - 完整的用户注册、登录、权限管理
- 🔐 **安全认证** - 密码重置、邮箱修改、邮件验证码
- 🔒 **权限控制** - 基于角色的访问控制（RBAC）
- 📊 **数据统计** - 访问统计、文章统计、用户统计
- 🎨 **主题切换** - 支持亮色/暗色主题

### 🛠️ 技术特性
- 🚀 **现代化技术栈** - Vue 3 + TypeScript + Go + PostgreSQL
- 🎨 **优雅 UI** - Naive UI 组件库 + 玻璃态设计
- 📱 **响应式设计** - 完美适配各种设备
- ⚡ **高性能** - Vite 构建 + Pinia 状态管理 + 异步邮件发送
- 🔐 **安全可靠** - JWT 认证 + 密码加密 + 邮箱验证 + IP 黑名单
- 🧹 **自动清理** - 定时清理过期数据，保持数据库整洁
- 📦 **易于部署** - Docker 支持 + 详细部署文档

## 🏗️ 项目结构

```
myBlog/
├── blog-frontend/          # 前端项目
│   ├── src/
│   │   ├── api/           # API 接口
│   │   ├── assets/        # 静态资源
│   │   ├── components/    # 公共组件
│   │   ├── layouts/       # 布局组件
│   │   ├── pages/         # 页面组件
│   │   ├── router/        # 路由配置
│   │   ├── stores/        # 状态管理
│   │   ├── types/         # TypeScript 类型
│   │   └── utils/         # 工具函数
│   ├── public/            # 公共资源
│   └── package.json
│
├── blog-backend/          # 后端项目
│   ├── cmd/
│   │   └── server/        # 服务入口
│   ├── config/            # 配置文件
│   ├── db/                # 数据库连接
│   ├── handler/           # 请求处理器
│   ├── middleware/        # 中间件
│   ├── model/             # 数据模型
│   ├── repository/        # 数据访问层
│   ├── service/           # 业务逻辑层
│   ├── util/              # 工具函数
│   ├── uploads/           # 上传文件
│   └── go.mod
│
└── README.md              # 项目说明
```

## 🚀 快速开始

### 环境要求

- **Node.js** >= 18.0.0
- **Go** >= 1.21
- **PostgreSQL** >= 15
- **pnpm** (推荐) 或 npm

### 1️⃣ 克隆项目

```bash
git clone <repository-url>
cd myBlog
```

### 2️⃣ 数据库配置

1. 创建 PostgreSQL 数据库：
```sql
CREATE DATABASE blogdb;
```

2. 导入数据库结构：
```bash
cd blog-backend/sql
psql -U postgres -d blogdb -f init.sql
```

### 3️⃣ 后端配置与启动

```bash
cd blog-backend

# 1. 安装依赖
go mod download

# 2. 配置数据库连接和邮箱服务
# 编辑 config/config-dev.yml
vim config/config-dev.yml

# 配置邮箱服务（用于密码重置）
# email:
#   host: smtp.qq.com
#   port: 587
#   username: your-email@qq.com
#   password: your-auth-code  # QQ邮箱授权码
#   from_name: 情迁阁

# 3. 运行后端服务
go run cmd/server/main.go
```

后端服务默认运行在 `http://localhost:8080`

### 4️⃣ 前端配置与启动

```bash
cd blog-frontend

# 1. 安装依赖
pnpm install

# 2. 配置 API 地址（可选）
# 创建 .env.development 文件
echo "VITE_API_BASE_URL=http://localhost:8080" > .env.development

# 3. 启动开发服务器
pnpm dev
```

前端服务默认运行在 `http://localhost:3000`

### 5️⃣ 访问系统

- **前台首页**: http://localhost:3000
- **管理后台**: http://localhost:3000/admin
- **默认管理员账号**: 
  - 用户名: `admin`
  - 密码: `admin123`

## ⚙️ 邮箱配置说明

### QQ邮箱授权码获取

1. 登录QQ邮箱网页版
2. 进入 **设置** → **账户**
3. 找到 **POP3/IMAP/SMTP/Exchange/CardDAV/CalDAV服务**
4. 开启 **POP3/SMTP服务** 或 **IMAP/SMTP服务**
5. 点击 **生成授权码**，按提示发送短信
6. 获得16位授权码，填入配置文件的 `password` 字段

### 其他邮箱配置

**163邮箱**:
```yaml
email:
  host: smtp.163.com
  port: 465
  username: your-email@163.com
  password: your-auth-code
  from_name: 情迁阁
```

**Gmail**:
```yaml
email:
  host: smtp.gmail.com
  port: 587
  username: your-email@gmail.com
  password: your-app-password
  from_name: 情迁阁
```

## 📦 生产部署

### 前端构建

```bash
cd blog-frontend
pnpm build
```

构建产物在 `dist` 目录，可部署到任何静态服务器（Nginx、Vercel、Netlify 等）。

### 后端编译

```bash
cd blog-backend
go build -o blog-backend cmd/server/main.go
```

### Docker 部署

```bash
# 后端
cd blog-backend
docker build -t blog-backend .
docker run -p 8080:8080 blog-backend

# 前端（需要先构建）
cd blog-frontend
pnpm build
# 使用 Nginx 等服务器部署 dist 目录
```

## 🎨 主要功能模块

### 📝 文章管理
- Markdown 编辑器，支持实时预览
- 代码高亮（支持多种编程语言）
- 图片上传和管理
- 文章分类和标签
- 文章置顶和草稿
- 点赞和浏览统计

### 💬 评论系统
- 支持嵌套回复
- 评论审核
- 评论状态管理
- 用户头像显示

### 💭 说说动态
- 发布图文动态
- 多图上传（最多9张）
- 点赞功能
- 公开/私密状态

### 👤 用户中心
- 用户注册和登录
- 个人资料编辑
- 头像上传
- 密码修改
- 忘记密码（邮箱验证码）
- 邮箱修改（限制一年2次）

### 🔧 管理后台
- 📊 仪表盘数据统计（访问统计、分类统计）
- 📝 文章管理
- 🏷️ 分类标签管理
- 💬 评论管理
- 💭 说说管理
- 👥 用户管理
- ⚙️ 网站设置
- 🚫 IP 黑名单管理

## 🛠️ 技术栈

### 前端技术
- **框架**: Vue 3.3 + TypeScript
- **构建工具**: Vite 5
- **UI 组件**: Naive UI
- **状态管理**: Pinia
- **路由**: Vue Router 4
- **HTTP 客户端**: Axios
- **Markdown**: @kangc/v-md-editor
- **代码高亮**: Prism.js
- **图表**: ECharts
- **工具库**: VueUse、Day.js

### 后端技术
- **语言**: Go 1.21+
- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: PostgreSQL
- **认证**: JWT
- **邮件发送**: SMTP (支持QQ邮箱、163邮箱等)
- **日志**: 自定义日志中间件
- **配置**: Viper (YAML)
- **密码加密**: bcrypt
- **定时任务**: Go 原生 Goroutine + Timer

## 📖 API 文档

### 认证相关
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/logout` - 用户登出
- `POST /api/auth/refresh` - 刷新Token
- `GET /api/auth/profile` - 获取用户信息
- `PUT /api/auth/profile` - 更新用户信息
- `PUT /api/auth/password` - 修改密码
- `POST /api/auth/forgot-password` - 忘记密码（发送验证码）
- `POST /api/auth/reset-password` - 重置密码
- `PUT /api/auth/email` - 修改邮箱
- `GET /api/auth/email-change-info` - 获取邮箱修改信息

### 文章相关
- `GET /api/posts` - 获取文章列表
- `GET /api/posts/:id` - 获取文章详情
- `GET /api/posts/archives` - 获取归档
- `GET /api/posts/hot` - 获取热门文章
- `GET /api/posts/recent` - 获取最新文章
- `POST /api/posts` - 创建文章（需认证）
- `PUT /api/posts/:id` - 更新文章（需认证）
- `DELETE /api/posts/:id` - 删除文章（需认证）
- `POST /api/posts/:id/like` - 点赞文章

### 分类相关
- `GET /api/categories` - 获取分类列表
- `GET /api/categories/:id` - 获取分类详情
- `POST /api/categories` - 创建分类（管理员）
- `PUT /api/categories/:id` - 更新分类（管理员）
- `DELETE /api/categories/:id` - 删除分类（管理员）

### 标签相关
- `GET /api/tags` - 获取标签列表
- `GET /api/tags/:id` - 获取标签详情
- `GET /api/tags/:id/posts` - 获取标签下的文章
- `POST /api/tags` - 创建标签（需认证）
- `PUT /api/tags/:id` - 更新标签（需认证）
- `DELETE /api/tags/:id` - 删除标签（需认证）

### 评论相关
- `GET /api/comments/post/:id` - 获取文章评论
- `POST /api/comments` - 创建评论（需认证）
- `PUT /api/comments/:id` - 更新评论（需认证）
- `DELETE /api/comments/:id` - 删除评论（需认证）

### 说说相关
- `GET /api/moments` - 获取说说列表
- `GET /api/moments/:id` - 获取说说详情
- `GET /api/moments/recent` - 获取最新说说
- `POST /api/moments` - 发布说说（需认证）
- `PUT /api/moments/:id` - 更新说说（需认证）
- `DELETE /api/moments/:id` - 删除说说（需认证）
- `POST /api/moments/:id/like` - 点赞说说

### 上传相关
- `POST /api/upload/avatar` - 上传头像（需认证）
- `POST /api/upload/image` - 上传图片（需认证）

### 设置相关
- `GET /api/settings/public` - 获取公开设置
- `GET /api/settings/site` - 获取网站设置（管理员）
- `PUT /api/settings/site` - 更新网站设置（管理员）

### 验证码相关
- `GET /api/captcha` - 获取图形验证码

### 管理后台相关
- `GET /api/admin/dashboard/stats` - 仪表盘统计
- `GET /api/admin/dashboard/category-stats` - 分类统计
- `GET /api/admin/dashboard/visit-stats` - 访问统计
- `GET /api/admin/users` - 用户列表
- `PUT /api/admin/users/:id/status` - 更新用户状态
- `DELETE /api/admin/users/:id` - 删除用户
- `GET /api/admin/posts` - 所有文章
- `GET /api/admin/comments` - 所有评论
- `PUT /api/admin/comments/:id/status` - 更新评论状态
- `GET /api/admin/moments` - 所有说说
- `GET /api/admin/ip-blacklist` - IP黑名单列表
- `POST /api/admin/ip-blacklist` - 添加IP黑名单
- `DELETE /api/admin/ip-blacklist/:id` - 删除IP黑名单
- `GET /api/admin/ip-blacklist/check` - 检查IP状态
- `POST /api/admin/ip-blacklist/clean-expired` - 清理过期IP

更多详细说明请参考 [后端文档](./blog-backend/README.md)

## 🎯 开发指南

### 代码规范
- 使用 ESLint + Prettier 格式化代码
- 遵循 Vue 3 Composition API 风格
- TypeScript 严格模式
- Git Commit 规范

### 目录说明
详见各子项目的 README：
- [前端文档](./blog-frontend/README.md)
- [后端文档](./blog-backend/README.md)

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

## 📝 更新日志

### v1.1.0 (2025-10-28)
- ✨ 新增忘记密码功能（邮箱验证码）
- ✨ 新增邮箱修改功能（限制一年2次）
- 🔐 邮件验证码安全认证
- ⚡ 异步邮件发送，提升响应速度
- 🧹 自动清理过期验证码（每小时）
- 🎨 优化认证页面UI

### v1.0.0 (2025-10-25)
- ✨ 初始版本发布
- 📝 完整的文章管理系统
- 💬 评论系统
- 💭 说说功能
- 🎨 主题切换
- 📊 数据统计
- 🔐 权限管理

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 👨‍💻 作者

**情随事迁**

## 🙏 致谢

- [Vue.js](https://vuejs.org/)
- [Naive UI](https://www.naiveui.com/)
- [Gin](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- 所有贡献者

## 📮 联系方式

如有问题或建议，欢迎通过以下方式联系：

- 📧 Email: 2665339398@qq.com


---

⭐ 如果这个项目对你有帮助，请给个 Star 支持一下！

