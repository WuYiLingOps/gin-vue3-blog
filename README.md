# 🎨 个人博客系统

一个基于 Vue 3 + Go 的现代化全栈博客系统，采用前后端分离架构，具有优雅的 UI 设计和完善的功能。

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Vue](https://img.shields.io/badge/Vue-3.3+-green.svg)
![Go](https://img.shields.io/badge/Go-1.21+-blue.svg)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-blue.svg)

## 🌐 在线演示

**网站地址**: [https://blog.leadcode.top/](https://blog.leadcode.top/)

- 📱 支持 PC、平板、手机访问
- 💬 可以访问聊天室与其他用户实时交流
- 👤 支持匿名访问或注册登录

## 📸 项目截图

### 首页
![首页](./screenshots/home.png)

### 文章详情
![文章详情](./screenshots/post-detail.png)

### 聊天室
![聊天室](./screenshots/chatroom.png)

### 管理后台
![管理后台](./screenshots/admin-dashboard.png)

## ✨ 特性

### 🎯 核心功能
- 📝 **文章管理** - Markdown 编辑器，支持代码高亮、图片上传
- 🏷️ **分类标签** - 灵活的分类和标签系统
- 💬 **评论系统** - 支持嵌套回复的评论功能，支持多种评论类型（文章评论、友链评论等）
- 🔗 **友情链接** - 友链管理功能，支持独立评论系统
- 💭 **说说动态** - 类似朋友圈的动态发布
- 💬 **实时聊天室** - WebSocket 实时通信，支持登录用户和匿名访问
- 👤 **用户系统** - 完整的用户注册、登录、权限管理
- 🔐 **安全认证** - 密码重置、邮箱修改、邮件验证码
- 🛡️ **验证码系统** - 图形验证码，基于 Redis 存储，支持 IP 限流和防暴力破解
- 🔒 **权限控制** - 基于角色的访问控制（RBAC）
- 📊 **数据统计** - 最近 7 天访问量趋势、文章统计、用户统计
- 🎨 **主题切换** - 支持亮色/暗色主题

### 🛠️ 技术特性
- 🚀 **现代化技术栈** - Vue 3 + TypeScript + Go + PostgreSQL
- 🎨 **优雅 UI** - Naive UI 组件库 + 玻璃态设计
- 📱 **响应式设计** - 完美适配各种设备（PC/平板/手机）
- ⚡ **高性能** - Vite 构建 + Pinia 状态管理 + 异步邮件发送
- 💬 **实时通信** - WebSocket 长连接 + 自动重连 + 心跳检测
- 🔐 **安全可靠** - JWT 认证 + 密码加密 + 邮箱验证 + IP 黑名单
- 🧹 **自动清理** - 定时清理过期数据，保持数据库整洁
- 📦 **易于部署** - Docker 支持 + 详细部署文档

## 🏗️ 项目结构

```bash
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
- **Redis** >= 3.0 (建议 >= 6.2 以避免客户端警告)
- **pnpm** (推荐) 或 npm

###  0️⃣docker快速启动环境（可选）

创建`pgsql`指令：

```bash
docker run --name pg-prod \
  -e POSTGRES_PASSWORD="123456ok!" \
  -v /data/PgSqlData:/var/lib/postgresql/data \
  -p 5432:5432 \
  -d postgres:17-alpine
```

创建`redis`指令:

```bash
docker run --name redis-prod \
  -p 6379:6379 \
  -v /data/redisData:/data \
  -e REDIS_PASSWORD=123456 \
  --restart=always \
  -d redis:7-alpine \
  redis-server --requirepass 123456 --appendonly yes
```

查看是否创建成功：

```bash
[root@docker-server ~]# docker ps
CONTAINER ID   IMAGE                COMMAND                  CREATED          STATUS          PORTS                                         NAMES
51e019841d66   redis:7-alpine       "docker-entrypoint.s…"   18 minutes ago   Up 18 minutes   0.0.0.0:6379->6379/tcp, [::]:6379->6379/tcp   redis-prod
22205f8e78c6   postgres:17-alpine   "docker-entrypoint.s…"   34 minutes ago   Up 34 minutes   0.0.0.0:5432->5432/tcp, [::]:5432->5432/tcp   pg-prod
```

完成

### 1️⃣ 克隆项目

```bash
git clone https://gitee.com/qssq9398/go-vue3-blog.git
cd myBlog
```

### 2️⃣ 数据库配置

1. 创建 PostgreSQL 数据库（指定编码/排序，便于跨版本迁移一致）：
```sql
CREATE DATABASE blogdb
  WITH OWNER = postgres
       ENCODING = 'UTF8'
       LC_COLLATE = 'en_US.utf8'
       LC_CTYPE   = 'en_US.utf8'
       TEMPLATE   = template0;
```

2. 导入数据库结构：
```bash
cd blog-backend/sql
psql -U postgres -d blogdb -f init.sql
```

> 如果是docker启动的pgsql则看下面

```bash
第一步：
# 进入容器内的 psql 交互界面
docker exec -it pg-prod psql -U postgres

# 在 psql 中创建 blogdb 库（执行后输入 \q 退出）
CREATE DATABASE blogdb
  WITH OWNER = postgres
       ENCODING = 'UTF8'
       LC_COLLATE = 'en_US.utf8'
       LC_CTYPE   = 'en_US.utf8'
       TEMPLATE   = template0;
# 退出
\q

第二步：
# 将 init.sql 传入容器
docker cp go-vue3-blog/blog-backend/sql/init.sql pg-prod:/tmp/init.sql

# 导入数据
docker exec -it pg-prod psql -U postgres -d blogdb -f /tmp/init.sql
```

### 3️⃣ 后端配置与启动

> 如果没有配置go的镜像代理，可以参考[Go 国内加速：Go 国内加速镜像 | Go 技术论坛](https://learnku.com/go/wikis/38122)

```bash
cd blog-backend

# 1. 安装依赖
go mod download

# 2. 配置数据库连接和邮箱服务（推荐：YAML + .env.config.dev 组合）
# （1）先编辑 config/config-dev.yml，填入非敏感的默认配置（如主机、端口等）
#     可参考 blog-backend/config/env.config.example 中的示例字段说明
vim config/config-dev.yml

# （2）在后端项目根目录（blog-backend）创建 .env.config.dev，仅填写真正的账号/密码等敏感信息，例如：
# DB_HOST=127.0.0.1
# DB_PORT=5432
# DB_USER=postgres
# DB_PASSWORD=安全的数据库密码
# DB_NAME=blogdb
# REDIS_HOST=127.0.0.1
# REDIS_PORT=6379
# REDIS_PASSWORD=安全的redis密码
# JWT_SECRET=更复杂的JWT密钥
# EMAIL_PASSWORD=邮箱授权码



# 配置邮箱服务（用于密码重置）
# email:
#   host: smtp.qq.com
#   port: 587
#   username: your-email@qq.com
#   password: your-auth-code  # QQ邮箱授权码
#   from_name: 情迁阁

# 3. 配置 Gitee 贡献热力图 API（可选）
# 如果需要首页贡献热力图，需要在 config/config-dev.yml 中配置 gitee-calendar-api 地址：
# gitee_calendar:
#   api_url: "http://localhost:8081/api"  # gitee-calendar-api 服务地址
# 
# 然后启动 gitee-calendar-api 服务（默认端口 8081，路径 /api）：
# ./gitee-calendar-api   # 前台运行
# nohup ./gitee-calendar-api > gitee-calendar-api.log 2>&1 &   # 后台运行
# 
# 注意：gitee-calendar-api 现在由后端调用，前端不再直接访问该服务

# 4. 运行后端服务
go run cmd/server/main.go
```

后端服务默认运行在 `http://localhost:8080`

### 4️⃣ 前端配置与启动

```bash
cd blog-frontend

# 1. 安装依赖
# 如果没有安装pnpm，可全局安装：npm install -g pnpm
pnpm install

# 2. 配置 API 地址（可选）
# 创建 .env.development 文件
echo "VITE_API_BASE_URL=http://localhost:8080" > .env.development
# 注意：贡献热力图现在通过后端 API 获取，无需单独配置 VITE_GITEE_CALENDAR_API

# 3. 启动开发服务器
pnpm dev
```

前端服务默认运行在 `http://localhost:3000`

### 5️⃣ 访问系统

- **前台首页**: http://localhost:3000
- **管理后台**: http://localhost:3000/admin
- **默认管理员账号**: 
  - 用户名: `admin`
  - 密码: `password`

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

### 第一步：后端部署

#### 方式一：Docker Compose（推荐）

使用 Docker Compose 一键部署后端服务（包含 PostgreSQL、Redis、后端应用）：

#### 1. 编译后端程序

```bash
cd blog-backend

# Windows PowerShell
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o blog-backend ./cmd/server

# Linux/Mac
GOOS=linux GOARCH=amd64 go build -o blog-backend ./cmd/server
```

#### 2. 配置环境变量（推荐）

> 生产环境同样推荐使用「YAML + .env.config.prod」方案，将敏感信息放到环境变量文件中，而不是写死在 `config-prod.yml` 里。

- **后端应用内部配置覆盖**  
  在后端可执行文件的工作目录下创建 `.env.config.prod`，用于覆盖 `config/config-prod.yml` 中的敏感字段；可参考 `blog-backend/config/env.config.example` 中的字段说明，例如：

  ```env
  # 数据库
  DB_HOST=postgres
  DB_PORT=5432
  DB_USER=postgres
  DB_PASSWORD=your_postgres_password
  DB_NAME=blogdb
  
  # Redis
  REDIS_HOST=redis
  REDIS_PORT=6379
  REDIS_PASSWORD=your_redis_password
  
  # Gitee Calendar API（如使用）
  GITEE_CALENDAR_API_URL=http://127.0.0.1:8081/api
  
  # JWT
  JWT_SECRET=your_jwt_secret
  
  # 阿里云 OSS（如使用）
  OSS_ENDPOINT=oss-cn-hangzhou.aliyuncs.com
  OSS_ACCESS_KEY_ID=your-ak
  OSS_ACCESS_KEY_SECRET=your-sk
  OSS_BUCKET_NAME=your-bucket
  
  # 腾讯云 COS（如使用）
  COS_BUCKET_URL=https://your-bucket.cos.ap-guangzhou.myqcloud.com
  COS_SECRET_ID=your-cos-secret-id
  COS_SECRET_KEY=your-cos-secret-key
  ```

- **环境变量字段说明（生产建议放在 `.env.config.prod`，模板：`blog-backend/config/env.config.example`）**

  - 基础必填：`DB_HOST` `DB_PORT` `DB_USER` `DB_PASSWORD` `DB_NAME`  
  - Redis：`REDIS_HOST` `REDIS_PORT` `REDIS_PASSWORD`  
  - JWT：`JWT_SECRET`（可选 `JWT_EXPIRE_HOURS`）  
  - 邮件（找回密码/验证邮件）：`EMAIL_HOST` `EMAIL_PORT` `EMAIL_USERNAME` `EMAIL_PASSWORD`（`EMAIL_FROM_NAME` 如需覆盖）  
  - Gitee 贡献热力图：`GITEE_CALENDAR_API_URL`（必填，后端调用 gitee-calendar-api）  
  - 对象存储可选：阿里云 OSS（`OSS_ENDPOINT` `OSS_ACCESS_KEY_ID` `OSS_ACCESS_KEY_SECRET` `OSS_BUCKET_NAME` `OSS_DOMAIN`），腾讯云 COS（`COS_BUCKET_URL` `COS_SECRET_ID` `COS_SECRET_KEY` `COS_DOMAIN`）  
  - 其他按需：如有自定义新增字段，统一放入 `.env.config.prod` 并在后端读取。

- **数据库 / Redis 容器自身变量**  
  你仍然可以在 `docker-compose.yml` 同级的 `.env` 文件中，为 PostgreSQL / Redis 设置自身的密码等变量，例如：

```env
  # PostgreSQL 容器内部初始化密码
POSTGRES_PASSWORD=your_postgres_password
POSTGRES_DB=blogdb

  # Redis 容器内部初始化密码
REDIS_PASSWORD=your_redis_password
```

  或者直接在 `docker-compose.yml` 文件中内联这些环境变量。

#### 3. 启动所有服务

```bash
# 构建并启动所有服务
docker compose up -d --build

# 查看服务状态
docker compose ps

# 查看日志
docker compose logs -f backend
```

#### 4. 初始化数据库

```bash
# 进入数据库容器
docker exec -it pg-prod psql -U postgres -d blogdb

# 创建数据库
CREATE DATABASE blogdb
  WITH OWNER = postgres
       ENCODING = 'UTF8'
       LC_COLLATE = 'en_US.utf8'
       LC_CTYPE   = 'en_US.utf8'
       TEMPLATE   = template0;

# 开始外部导入 SQL
docker exec -i pg-prod psql -U postgres -d blogdb < sql/init.sql
# 如果是导入备份数据则如下：
docker cp blog_backup.dump pg-prod:/tmp
docker exec -it pg-prod pg_restore -U postgres -d blogdb --clean --if-exists /tmp/blog_backup.dump
```

#### 5. 服务管理

```bash
# 停止所有服务
docker compose down

# 停止并删除数据卷（谨慎使用！）
docker compose down -v

# 重启服务
docker compose restart backend

# 查看日志
docker compose logs -f backend
```

#### 服务访问地址：
- **后端 API**: `http://localhost:8080`
- **PostgreSQL**: `localhost:5632`
- **Redis**: `localhost:6379`

#### 方式二：本地编译部署✅

##### 1. 启动 Gitee Calendar API 服务（必选，贡献热力图依赖）

1. **部署 gitee-calendar-api 服务**  
   本仓库已自带编译好的 `gitee-calendar-api`（根目录），可直接赋予执行权限使用（默认占用端口为8081）。若需查看/自行编译源码，可访问：`https://gitee.com/wylblog/go-code-calendar-api.git`

   ```bash
   cd /web/go-vue3-blog/gitee-calendar-api
   chmod +x gitee-calendar-api
   ```

2. **启动 gitee-calendar-api 服务**（示例为简单后台运行方式，生产环境可用 systemd 管理）：

   ```bash
   # 前台调试运行
   ./gitee-calendar-api

   # 或后台运行（输出到 gitee-calendar-api.log）
   nohup ./gitee-calendar-api > gitee-calendar-api.log 2>&1 &
   ```

   默认监听端口为 `8081`，路径为 `/api`，即本机访问地址为：`http://127.0.0.1:8081/api?user=你的Gitee用户名`。

   > **说明**：
   > - `gitee-calendar-api` 现在由**后端调用**，前端不再直接访问该服务
   > - 前端通过后端 API `/api/calendar/gitee` 获取热力图数据，后端会调用 `gitee-calendar-api` 并缓存结果（20分钟过期）
   > - `gitee-calendar-api` 通常部署在与后端相同的服务器上，通过 `127.0.0.1:8081` 访问，无需通过 Nginx 暴露给外部

##### 2. 配置后端 Gitee Calendar API 地址（必选，生产环境统一用环境变量）

在后端根目录创建（或编辑）`.env.config.prod`，通过环境变量配置 `gitee-calendar-api` 服务地址（不再修改 `config/config-prod.yml`）。**模板已提供：`blog-backend/config/env.config.example`**，可直接复制为 `.env.config.prod` 后按需修改：

```env
# .env.config.prod
GITEE_CALENDAR_API_URL=http://127.0.0.1:8081/api
```

**通过 Nginx 代理（HTTP 或 HTTPS）**

如果 `gitee-calendar-api` 通过 Nginx 代理访问，可配置为：

```env
# HTTP 方式
GITEE_CALENDAR_API_URL=http://your-domain.com/gitee-calendar-api

# HTTPS 方式（SSL 证书）
GITEE_CALENDAR_API_URL=https://your-domain.com/gitee-calendar-api
# 示例：GITEE_CALENDAR_API_URL=https://huangjingblog.cn/gitee-calendar-api
```

##### 3. 构建并启动后端服务

```bash
cd blog-backend

# 构建后端可执行文件
go build -o blog-backend cmd/server/main.go

# 前台运行（调试用）
./blog-backend

# 后台运行（简单方式，生产环境建议配合 systemd 等守护进程管理）
nohup ./blog-backend > app.log 2>&1 &
```

手动在主机安装并启动 PostgreSQL、Redis，按需配置 `config/config-prod.yml`，再以服务方式管理可执行文件。

### 第二步：前端构建与配置

#### 2.1 前端 `.env.production` 环境变量配置

1. 在 `blog-frontend` 目录下创建或编辑 `.env.production`：

   ```bash
   cd blog-frontend
   vim .env.production
   ```

2. 写入（或补充）如下内容（根据你的实际域名调整）：

   **HTTP 方式：**
   ```env
   # 后端主 API（博客业务接口）
   # 方式一：只配置域名（推荐，更简洁）
   VITE_API_BASE_URL=http://your-domain.com
   
   # 方式二：配置包含 /api 的完整路径（也可以，函数会自动处理）
   # VITE_API_BASE_URL=http://your-domain.com/api
   ```

   **HTTPS 方式（SSL 证书）：**
   ```env
   # 后端主 API（博客业务接口）
   # 方式一：只配置域名（推荐，更简洁）
   VITE_API_BASE_URL=https://your-domain.com
   # 示例：VITE_API_BASE_URL=https://huangjingblog.cn
   
   # 方式二：配置包含 /api 的完整路径（也可以，函数会自动处理）
   # VITE_API_BASE_URL=https://your-domain.com/api
   # 示例：VITE_API_BASE_URL=https://huangjingblog.cn/api
   ```

   - `VITE_API_BASE_URL`：博客后端（Gin 服务）的基础地址，前端所有业务接口都会基于该地址请求，包括贡献热力图数据（通过 `/api/calendar/gitee` 接口获取）。
   - **推荐配置方式**：只配置域名（如 `https://your-domain.com`），贡献热力图组件会自动添加 `/api` 前缀。如果已配置包含 `/api` 的路径，也能正常工作。

#### 2.2 构建前端项目

```bash
cd blog-frontend
pnpm install
pnpm build
```

构建产物在 `dist` 目录，可部署到任何静态服务器（Nginx、Vercel、Netlify 等）。部署新的 `dist` 后，首页热力图会自动通过后端 API 获取数据。

> 说明：前端首页热力图组件位置为 `blog-frontend/src/components/GiteeCalendar.vue`，其数据源现在通过后端 API `/api/calendar/gitee` 获取，后端会自动调用 `gitee-calendar-api` 并缓存结果。

### 第三步：Nginx 部署与反向代理

1. 在服务器上准备前端目录（例如 `/web/go-vue3-blog/blog-frontend/dist`），**将本地 `dist` 目录中的所有文件和子目录整体上传到该目录**，保持结构不变，例如：

   ```text
   /web/go-vue3-blog/blog-frontend/dist/
   ├── index.html
   ├── assets/
   ├── logo.svg
   └── 备案图标.png
   ```

   Nginx 中的 `root` 应指向 **包含 `index.html` 的目录本身**（如 `/web/go-vue3-blog/blog-frontend/dist`，可按实际路径调整），而不是上级目录。
2. 配置 Nginx（按需替换域名/路径/证书），HTTP 示例：

```nginx
server {
    listen 80;
    server_name your-domain.com;   # 修改为你的域名/主机名，例如：huangjingblog.cn

    # 前端静态资源目录（dist 构建产物）
    root /web/go-vue3-blog/blog-frontend/dist;  # 按实际部署路径修改
    index index.html;

    # 前端路由回退到 index.html（适配前端 history 模式）
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Gitee 贡献日历 API（go-code-calendar-api）反向代理
    # 对应后端 .env.config.prod 中的 GITEE_CALENDAR_API_URL=http://your-domain.com/gitee-calendar-api
    location /gitee-calendar-api {
        proxy_pass http://127.0.0.1:8081/api;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 本地存储上传文件访问（通过后端读取 /uploads 下资源）
    # 使用 ^~ 确保优先级高于下面的静态资源正则 location
    location ^~ /uploads/ {
        proxy_pass http://127.0.0.1:8080;  # 与后端 API 相同地址
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # WebSocket（聊天室连接：/api/chat/ws），需保持连接与协议升级
    location /api/chat/ws {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 后端 API 反向代理（其余 /api/* 请求）
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 可选：静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|svg|ico|woff2?)$ {
        try_files $uri =404;
        expires 30d;
        access_log off;
    }
}
```

3. HTTPS 示例（含 80→443 跳转，请替换证书路径）：

```nginx
# 80 强制跳转到 443
server {
    listen 80;
    server_name your-domain.com;   # 修改为你的域名/主机名，例如：huangjingblog.cn
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;   # 修改为你的域名/主机名，例如：huangjingblog.cn

    # 证书路径（替换为实际证书文件）
    ssl_certificate     /usr/local/nginx/ssl/your-domain.com.pem;  # 例如：/usr/local/nginx/ssl/huangjingblog.cn.pem
    ssl_certificate_key /usr/local/nginx/ssl/your-domain.com.key;  # 例如：/usr/local/nginx/ssl/huangjingblog.cn.key
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # 前端静态资源目录（dist 构建产物）
    root /web/go-vue3-blog/blog-frontend/dist;  # 按实际部署路径修改
    index index.html;

    # 前端路由回退
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Gitee 贡献日历 API（go-code-calendar-api）反向代理
    # 对应后端 .env.config.prod 中的 GITEE_CALENDAR_API_URL=https://your-domain.com/gitee-calendar-api
    location /gitee-calendar-api {
        proxy_pass http://127.0.0.1:8081/api;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 本地存储上传文件访问（通过后端读取 /uploads 下资源）
    # 使用 ^~ 确保优先级高于下面的静态资源正则 location
    location ^~ /uploads/ {
        proxy_pass http://127.0.0.1:8080;  # 与后端 API 相同地址
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # WebSocket（聊天室连接：/api/chat/ws）
    location /api/chat/ws {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 后端 API 反代
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 可选：静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|svg|ico|woff2?)$ {
        try_files $uri =404;
        expires 30d;
        access_log off;
    }
}
```

4. 重载 Nginx：`nginx -s reload` 或 `systemctl reload nginx`。

### 第四步：数据迁移（可选，新旧数据库切换时使用）

1. 备份旧库（PostgreSQL 示例）：
```bash
pg_dump -h old-host -U old_user -d old_db -Fc -f backup.dump
```
参数说明：
- `-h old-host`：旧数据库主机
- `-U old_user`：旧库用户名
- `-d old_db`：旧库库名
- `-Fc`：自定义格式（便于 pg_restore 增强选项）
- `-f backup.dump`：输出文件路径
示例：
```bash
pg_dump -h 10.0.0.5 -U bloguser -d blogdb -Fc -f /tmp/blog_backup_2025-12-15.dump
```

2. 恢复到新库：
```bash
pg_restore -h new-host -U new_user -d new_db --clean --if-exists backup.dump
```
参数说明：
- `-h new-host`：新数据库主机
- `-U new_user`：新库用户名
- `-d new_db`：新库库名
- `--clean`：导入前先 DROP 现有对象
- `--if-exists`：仅在对象存在时执行 DROP，减少报错
- `backup.dump`：备份文件路径
示例：
```bash
pg_restore -h 127.0.0.1 -U postgres -d blogdb --clean --if-exists /tmp/blog_backup_2025-12-15.dump
```
3. 如使用 Docker Compose，可在容器内执行：
```bash
docker exec -i pg-prod pg_restore -U postgres -d blogdb --clean --if-exists < backup.dump
```
4. 确认新库的连接信息已写入 `config/config-prod.yml` 或环境变量，并与 Nginx/后端代理地址匹配。
5. **迁移上传文件目录**（重要）：
   - 如果使用**本地存储**，需要将旧服务器的 `blog-backend/uploads/` 目录完整复制到新服务器
   - 包括 `uploads/avatars/`（用户头像）和 `uploads/` 下的其他图片文件
   - 迁移方法：
     ```bash
     # 在旧服务器上打包
     cd blog-backend
     tar -czf uploads-backup.tar.gz uploads/
     
     # 传输到新服务器（使用 scp 或其他方式）
     scp uploads-backup.tar.gz user@new-server:/path/to/blog-backend/
     
     # 在新服务器上解压
     cd blog-backend
     tar -xzf uploads-backup.tar.gz
     ```
   - **注意**：如果使用 OSS/COS 存储，文件在云端，无需迁移本地文件；但普通用户头像强制使用本地存储，仍需迁移 `uploads/avatars/` 目录
6. 如不需要历史会话，可清理 Redis（如果有登录态存储），再重启后端。

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
- **多类型评论支持** - 支持文章评论、友链评论等多种评论类型
- **独立评论系统** - 友链页面拥有独立的评论系统，不依赖文章
- **评论通知功能** - 管理员可配置是否接收评论通知邮件（由于普通用户无权限写文章，文章作者只能是管理员，统一通过管理员通知处理）
- **邮件通知优化** - 智能处理SMTP服务器响应，确保邮件发送稳定性

### 💭 说说动态
- 发布图文动态
- 多图上传（最多9张）
- 点赞功能
- 公开/私密状态

### 🔗 友情链接
- 友链管理（后台添加、编辑、删除）
- 友链展示（前台展示友链列表）
- 友链申请（支持在友链页面申请）
- **独立评论系统** - 友链页面拥有独立的评论功能，支持嵌套回复
- 我的友链信息配置（名称、描述、URL、头像、站点图片、RSS订阅等）
- YAML 格式友链信息导出

### 👤 用户中心
- 用户注册和登录
- 个人资料编辑
- 头像上传（普通用户默认使用本地存储，管理员可使用 OSS/COS）
- 密码修改
- 忘记密码（邮箱验证码）
- 邮箱修改（限制一年2次）

### 🛡️ 验证码系统
- **图形验证码** - 基于 Redis 存储，支持数字验证码
- **验证码存储** - 验证码答案存储在 Redis 中，2分钟自动过期
- **IP 限流** - 每个 IP 每分钟最多获取 10 次验证码，防止频繁请求
- **防暴力破解** - 5 分钟内最多错误 5 次，超过限制需等待
- **安全机制** - 验证成功后自动删除验证码，确保一次性使用

### 💬 聊天室功能
- 实时 WebSocket 通信
- 支持登录用户和匿名访问
- 在线人数统计（按用户去重）
- 消息历史记录（最近50条）
- 表情符号选择器
- 管理员右键菜单功能：
  - 右键消息删除
  - 右键头像踢出用户
- 管理后台功能：
  - 批量删除消息
  - 发送系统广播
  - 踢出在线用户
  - 封禁IP地址
- 移动端响应式适配

### 🔧 管理后台
- 📊 仪表盘数据统计（最近 7 天访问量折线图、文章分类占比）
- 📝 文章管理
- 🏷️ 分类标签管理
- 💬 评论管理
- 🔗 友链管理（添加、编辑、删除友链，配置我的友链信息）
- 💭 说说管理
- 💬 聊天室管理（消息管理、用户管理）
- 👥 用户管理
- ⚙️ 网站设置（包含通知配置）
- 🚫 IP 访问控制（黑名单/白名单统一管理）

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
- **缓存**: Redis
  - 验证码存储与管理（2分钟过期）
  - IP 限流（防止频繁获取验证码）
  - 防暴力破解（错误次数限制）
  - Gitee 贡献热力图数据缓存（20分钟过期）
- **认证**: JWT
- **WebSocket**: Gorilla WebSocket
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
- `GET /api/comments/type?type={type}&target_id={target_id}` - 根据评论类型和目标ID获取评论（用于友链等特殊页面）
  - `type`: 评论类型（`post`-文章评论，`friendlink`-友链评论）
  - `target_id`: 目标ID（友链评论时固定为0）
- `POST /api/comments` - 创建评论（需认证）
  - 文章评论：`{ "content": "...", "comment_type": "post", "post_id": 1 }`
  - 友链评论：`{ "content": "...", "comment_type": "friendlink", "target_id": 0 }`
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
  - 普通用户：强制使用本地存储，不上传到 OSS/COS
  - 管理员：根据后台配置选择存储方式（本地/OSS/COS）
- `POST /api/upload/image` - 上传图片（需认证，根据配置选择存储方式）

### 友链相关
- `GET /api/friend-links` - 获取友链列表（公开）
- `GET /api/admin/friend-links` - 获取友链列表（管理员）
- `GET /api/admin/friend-links/:id` - 获取友链详情（管理员）
- `POST /api/admin/friend-links` - 创建友链（管理员）
- `PUT /api/admin/friend-links/:id` - 更新友链（管理员）
- `DELETE /api/admin/friend-links/:id` - 删除友链（管理员）
- `GET /api/settings/friendlink-info` - 获取我的友链信息（公开）
- `PUT /api/admin/settings/friendlink-info` - 更新我的友链信息（管理员）

### 设置相关
- `GET /api/settings/public` - 获取公开设置
- `GET /api/settings/site` - 获取网站设置（管理员）
- `PUT /api/settings/site` - 更新网站设置（管理员）
- `GET /api/settings/notification` - 获取通知配置（管理员）
- `PUT /api/settings/notification` - 更新通知配置（管理员）
  - 支持配置管理员评论通知开关

### 验证码相关
- `GET /api/captcha` - 获取图形验证码
  - 基于 Redis 存储验证码答案（2分钟过期）
  - IP 限流：每个 IP 每分钟最多获取 10 次
  - 防暴力破解：5 分钟内最多错误 5 次

### 日历相关
- `GET /api/calendar/gitee?user={username}` - 获取 Gitee 贡献热力图数据
  - 后端会调用 `gitee-calendar-api` 并缓存结果（20分钟过期）
  - 参数：`user` - Gitee 用户名

### 聊天室相关
- `WS /api/chat/ws` - WebSocket 连接（支持登录用户和匿名访问）
- `GET /api/chat/messages` - 获取消息列表
- `GET /api/chat/online` - 获取在线信息

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
- `GET /api/admin/chat/messages` - 聊天消息列表（管理员）
- `DELETE /api/admin/chat/messages/:id` - 删除消息（管理员）
- `POST /api/admin/chat/broadcast` - 发送系统广播（管理员）
- `POST /api/admin/chat/kick` - 踢出用户（管理员）
- `POST /api/admin/chat/ban` - 封禁IP（管理员）

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

### v1.4.0 (2025-12-20)
- ✨ 新增评论通知功能
  - 📧 管理员评论通知（可配置开关）
  - 🔔 当有用户评论文章时，可选择性通知所有管理员
  - ⚙️ 通知配置可在管理后台"网站设置"中开启/关闭
  - 🌐 智能URL获取：优先从"我的友链信息"中获取网站地址
  - 📝 邮件模板支持Markdown内容预览
  - 🎨 优化邮件模板样式，采用专业HTML邮件模板设计
  - 🐛 修复SMTP "short response" 错误处理，优化邮件发送稳定性
- 🔧 优化评论通知逻辑
  - 由于普通用户无权限写文章，文章作者只能是管理员，统一通过管理员通知处理
  - 避免重复通知：评论者本人不会收到通知邮件
  - 智能错误处理：忽略SMTP服务器的非标准响应，确保邮件发送成功时不误报错误
- 🎨 优化管理后台通知配置界面
  - 清晰展示管理员通知开关
  - 完善通知说明文档

### v1.3.0 (2025-12-19)
- ✨ 新增友情链接功能
  - 🔗 友链管理（后台添加、编辑、删除友链）
  - 📋 友链展示页面（前台展示友链列表）
  - 💬 友链独立评论系统（支持嵌套回复，不依赖文章）
  - ⚙️ 我的友链信息配置（名称、描述、URL、头像、站点图片、RSS订阅等）
  - 📄 YAML 格式友链信息导出
- 🔧 评论系统重构
  - 🎯 支持多种评论类型（文章评论、友链评论等）
  - 📊 评论表扩展：添加 `comment_type` 和 `target_id` 字段
  - 🔄 向后兼容：保持文章评论功能不变
  - 🗑️ 移除特殊文章依赖（不再需要创建ID=999999的特殊文章）

### v1.2.0 (2025-11-04)
- ✨ 新增实时聊天室功能
  - 💬 WebSocket 实时通信
  - 👥 支持登录用户和匿名访问
  - 📊 在线人数统计（按用户去重）
  - 🔄 自动重连机制
  - ❤️ 心跳检测保持连接
  - 📜 消息历史记录
  - 😀 表情符号选择器
- 🛡️ 管理员权限功能
  - 右键消息删除
  - 右键头像踢出用户
  - 批量删除消息
  - 发送系统广播
  - IP 封禁管理
- 📱 移动端完美适配
- 🐛 修复在线人数统计bug

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

