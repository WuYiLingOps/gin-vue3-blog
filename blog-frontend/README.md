# 个人博客系统 - 前端

基于 Vue 3 + TypeScript + Vite + Naive UI 构建的现代化博客前端系统。

## 技术栈

- **框架**: Vue 3 (Composition API)
- **构建工具**: Vite
- **语言**: TypeScript
- **状态管理**: Pinia
- **HTTP 客户端**: Axios
- **UI 组件库**: Naive UI
- **Markdown 编辑器**: @kangc/v-md-editor
- **路由**: Vue Router 4
- **工具库**: @vueuse/core, dayjs
- **代码规范**: ESLint + Prettier

## 项目结构

```
blog-frontend/
├── public/                 # 静态资源
├── src/
│   ├── api/               # API 接口
│   │   ├── auth.ts
│   │   ├── post.ts
│   │   ├── category.ts
│   │   ├── tag.ts
│   │   ├── comment.ts
│   │   ├── user.ts
│   │   ├── ipBlacklist.ts
│   │   └── index.ts
│   ├── assets/            # 资源文件
│   │   └── styles/
│   │       └── global.css
│   ├── components/        # 公共组件
│   ├── layouts/          # 布局组件
│   ├── pages/            # 页面组件
│   │   ├── admin/       # 管理后台页面
│   │   │   ├── IPAccessControl.vue
│   │   │   └── ...
│   ├── router/           # 路由配置
│   ├── stores/           # Pinia 状态管理
│   │   ├── auth.ts
│   │   ├── blog.ts
│   │   ├── app.ts
│   │   └── index.ts
│   ├── types/            # TypeScript 类型
│   │   ├── common.ts
│   │   ├── auth.ts
│   │   └── blog.ts
│   ├── utils/            # 工具函数
│   │   ├── request.ts
│   │   ├── storage.ts
│   │   ├── format.ts
│   │   ├── validator.ts
│   │   └── constants.ts
│   ├── App.vue
│   └── main.ts
├── index.html
├── vite.config.ts
├── tsconfig.json
├── package.json
└── README.md
```

## 快速开始

### 安装依赖

```bash
npm install
# 或
pnpm install
```

### 开发

```bash
npm run dev
```

应用将在 http://localhost:3000 启动

### 构建

```bash
npm run build
```

### 预览

```bash
npm run preview
```

### 代码检查

```bash
npm run lint
```

### 代码格式化

```bash
npm run format
```

## 环境变量

前端使用 Vite 的环境变量机制，推荐分别在开发 / 生产环境下创建对应文件：

- 开发环境：`blog-frontend/.env.development`
- 生产环境：`blog-frontend/.env.production`

示例（开发环境）：

```bash
VITE_APP_TITLE=我的博客
VITE_API_BASE_URL=http://localhost:8080
VITE_UPLOAD_URL=http://localhost:8080/api/upload
```

> 说明：
> - `VITE_API_BASE_URL` 用于配置后端 API 基础地址，开发时也会被 `vite.config.ts` 读取，用于 Vite 代理 `/api` 与 `/uploads` 请求。
> - 如未配置，开发代理会默认指向 `http://localhost:8080`。

### 更换后端端口（如从 `8080` 改为 `8090`）

如果后端端口修改为 `8090`，开发环境下需同步调整前端配置：

1. 修改后端 `config/config-dev.yml` 中的：
   ```yaml
   app:
     port: 8090
   ```
2. 在前端目录创建或修改 `.env.development`：
   ```bash
   VITE_API_BASE_URL=http://localhost:8090
   ```
3. 重启前端开发服务器：
   ```bash
   pnpm dev
   ```

生产环境同理，在 `.env.production` 中将 `VITE_API_BASE_URL` 指向后端实际的对外地址（域名或 IP + 端口），并确保 Nginx 的 `/api`、`/uploads` 反向代理目标与之保持一致。

## 功能特性

### 已完成

- ✅ 项目基础配置
- ✅ TypeScript 类型定义
- ✅ API 接口封装
- ✅ Pinia 状态管理
- ✅ 工具函数
- ✅ 请求拦截和响应处理
- ✅ 本地存储管理
- ✅ 日期格式化
- ✅ 表单验证
- ✅ IP 黑名单管理页面
- ✅ 管理后台页面（文章、分类、标签、评论、用户、说说、聊天、IP 黑名单等）
- ✅ 路由配置和权限守卫
- ✅ 布局组件（默认布局、管理布局、认证布局）
- ✅ 公共组件（Markdown 编辑器、图片上传、头像上传等）

## 开发规范

### 组件命名

- 页面组件使用 PascalCase
- 公共组件使用 PascalCase
- 组件文件名与组件名保持一致

### 代码风格

- 使用 Composition API
- 使用 `<script setup>` 语法
- 优先使用 TypeScript 类型推断
- 遵循 ESLint 和 Prettier 配置

### Git 提交规范

- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式
- `refactor`: 重构
- `test`: 测试
- `chore`: 构建/工具

## 浏览器支持

- Chrome >= 87
- Firefox >= 78
- Safari >= 14
- Edge >= 88

## 功能特性

### 管理后台功能

- ✅ **文章管理**：创建、编辑、删除文章，支持 Markdown 编辑
- ✅ **分类管理**：文章分类的增删改查
- ✅ **标签管理**：文章标签的增删改查
- ✅ **评论管理**：查看、删除评论
- ✅ **用户管理**：用户列表、状态管理、角色管理
  - 支持三种角色：超级管理员（super_admin）、管理员（admin）、普通用户（user）
  - 超级管理员可修改用户角色和状态
  - 管理员可查看用户列表，但无法修改用户角色
- ✅ **角色权限系统**
  - 基于角色的访问控制（RBAC）
  - 路由守卫实现基于角色的访问控制
  - 菜单根据用户角色动态显示
  - 超级管理员可访问所有功能，管理员部分功能受限
- ✅ **操作日志管理**（仅超级管理员）
  - 查看所有管理员和超级管理员的操作记录
  - 支持按模块、操作类型、用户名筛选
  - 支持单个删除和批量删除
  - 记录详细的操作信息（操作者、操作类型、目标、IP、时间等）
- ✅ **说说管理**：动态发布和管理
- ✅ **聊天管理**：实时聊天室管理，系统广播可选择投递范围（公告栏/聊天室/同时）
- ✅ **IP 黑名单管理**：查看、添加、删除被封禁的 IP，检查 IP 状态，清理过期记录
- ✅ **网站设置**：系统配置管理
- ✅ **数据统计**：访问量、文章数等数据统计

### IP 黑名单管理

IP 黑名单管理页面 (`/admin/ip-blacklist`) 提供以下功能：

1. **查看所有被封禁的 IP 列表**
   - 支持分页显示
   - 显示封禁类型（自动/手动）
   - 显示封禁原因和过期时间

2. **手动添加 IP 到黑名单**
   - 输入 IP 地址（支持 IPv4）
   - 设置封禁原因
   - 设置封禁时长（小时），0 表示永久封禁

3. **删除黑名单 IP（解除封禁）**
   - 一键解除封禁

4. **检查 IP 状态**
   - 快速查询任意 IP 是否被封禁
   - 显示封禁详情（如果被封禁）

5. **清理过期记录**
   - 一键清理所有过期的黑名单记录

### 博客前台功能

- ✅ 文章列表和详情页
- ✅ 分类和标签浏览
- ✅ 评论系统
- ✅ 说说动态
- ✅ 实时聊天室
- ✅ 用户注册和登录
- ✅ 个人资料管理

## 许可证

MIT License

