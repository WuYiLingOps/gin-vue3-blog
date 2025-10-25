# 个人博客系统 - 前端需求文档

## 技术栈
- **框架**: Vue 3 + Vite + TypeScript
- **状态管理**: Pinia
- **HTTP 客户端**: Axios
- **UI 组件库**: Naive UI
- **Markdown 编辑器**: @kangc/v-md-editor
- **路由管理**: Vue Router 4
- **代码规范**: ESLint + Prettier

## 项目结构
```
blog-frontend/
├── public/
│   ├── favicon.ico
│   └── index.html
├── src/
│   ├── api/                    # API 接口
│   │   ├── auth.ts            # 认证相关接口
│   │   ├── post.ts            # 文章相关接口
│   │   ├── category.ts        # 分类相关接口
│   │   ├── tag.ts             # 标签相关接口
│   │   ├── comment.ts         # 评论相关接口
│   │   └── index.ts           # 接口统一导出
│   ├── assets/                 # 静态资源
│   │   ├── images/
│   │   ├── icons/
│   │   └── styles/
│   │       ├── global.css     # 全局样式
│   │       ├── variables.css  # CSS 变量
│   │       └── themes/        # 主题样式
│   ├── components/             # 公共组件
│   │   ├── common/            # 通用组件
│   │   │   ├── AppHeader.vue  # 头部导航
│   │   │   ├── AppFooter.vue  # 底部信息
│   │   │   ├── BackTop.vue    # 返回顶部
│   │   │   ├── Loading.vue    # 加载组件
│   │   │   └── Pagination.vue # 分页组件
│   │   ├── blog/              # 博客相关组件
│   │   │   ├── PostCard.vue   # 文章卡片
│   │   │   ├── PostList.vue   # 文章列表
│   │   │   ├── CategoryList.vue # 分类列表
│   │   │   ├── TagCloud.vue   # 标签云
│   │   │   ├── CommentList.vue # 评论列表
│   │   │   └── CommentForm.vue # 评论表单
│   │   └── admin/             # 管理后台组件
│   │       ├── Sidebar.vue    # 侧边栏
│   │       ├── PostEditor.vue # 文章编辑器
│   │       └── DataTable.vue  # 数据表格
│   ├── composables/            # 组合式函数
│   │   ├── useAuth.ts         # 认证相关
│   │   ├── usePagination.ts   # 分页逻辑
│   │   ├── useTheme.ts        # 主题切换
│   │   └── useSearch.ts       # 搜索功能
│   ├── layouts/                # 布局组件
│   │   ├── DefaultLayout.vue  # 默认布局
│   │   ├── AdminLayout.vue    # 管理后台布局
│   │   └── AuthLayout.vue     # 认证页面布局
│   ├── pages/                  # 页面组件
│   │   ├── blog/              # 博客页面
│   │   │   ├── Home.vue       # 首页
│   │   │   ├── PostDetail.vue # 文章详情
│   │   │   ├── Category.vue   # 分类页面
│   │   │   ├── Tag.vue        # 标签页面
│   │   │   ├── Archive.vue    # 归档页面
│   │   │   └── About.vue      # 关于页面
│   │   ├── auth/              # 认证页面
│   │   │   ├── Login.vue      # 登录页面
│   │   │   ├── Register.vue   # 注册页面
│   │   │   └── Profile.vue    # 个人资料
│   │   └── admin/             # 管理后台页面
│   │       ├── Dashboard.vue  # 仪表盘
│   │       ├── PostManage.vue # 文章管理
│   │       ├── CategoryManage.vue # 分类管理
│   │       ├── TagManage.vue  # 标签管理
│   │       ├── CommentManage.vue # 评论管理
│   │       └── UserManage.vue # 用户管理
│   ├── router/                 # 路由配置
│   │   ├── index.ts           # 路由主文件
│   │   ├── guards.ts          # 路由守卫
│   │   └── routes/            # 路由模块
│   │       ├── blog.ts        # 博客路由
│   │       ├── auth.ts        # 认证路由
│   │       └── admin.ts       # 管理后台路由
│   ├── stores/                 # Pinia 状态管理
│   │   ├── auth.ts            # 认证状态
│   │   ├── blog.ts            # 博客状态
│   │   ├── theme.ts           # 主题状态
│   │   └── app.ts             # 应用全局状态
│   ├── types/                  # TypeScript 类型定义
│   │   ├── api.ts             # API 响应类型
│   │   ├── blog.ts            # 博客相关类型
│   │   ├── auth.ts            # 认证相关类型
│   │   └── common.ts          # 通用类型
│   ├── utils/                  # 工具函数
│   │   ├── request.ts         # Axios 封装
│   │   ├── auth.ts            # 认证工具
│   │   ├── format.ts          # 格式化工具
│   │   ├── storage.ts         # 本地存储工具
│   │   └── constants.ts       # 常量定义
│   ├── App.vue                 # 根组件
│   └── main.ts                 # 应用入口
├── .env                        # 环境变量
├── .env.development           # 开发环境变量
├── .env.production            # 生产环境变量
├── .eslintrc.js               # ESLint 配置
├── .prettierrc                # Prettier 配置
├── tsconfig.json              # TypeScript 配置
├── vite.config.ts             # Vite 配置
└── package.json               # 项目依赖
```

## 页面设计

### 1. 博客前台页面

#### 首页 (/)
- **功能**: 展示最新文章、热门文章、分类导航
- **组件**: 
  - 头部导航栏（搜索框、主题切换、用户菜单）
  - 轮播图/置顶文章
  - 文章列表（分页）
  - 侧边栏（分类、标签云、最新评论）
  - 底部信息

#### 文章详情页 (/post/:id)
- **功能**: 展示文章内容、评论系统
- **组件**:
  - 文章标题、作者、发布时间、分类、标签
  - Markdown 渲染的文章内容
  - 文章操作（点赞、分享）
  - 评论列表和评论表单
  - 相关文章推荐
  - 文章目录导航

#### 分类页面 (/category/:id)
- **功能**: 展示特定分类下的文章
- **组件**:
  - 分类信息展示
  - 该分类下的文章列表
  - 分页组件

#### 标签页面 (/tag/:id)
- **功能**: 展示特定标签下的文章
- **组件**:
  - 标签信息展示
  - 该标签下的文章列表
  - 分页组件

#### 归档页面 (/archive)
- **功能**: 按时间归档展示文章
- **组件**:
  - 年份/月份分组的文章列表
  - 时间轴样式展示

#### 关于页面 (/about)
- **功能**: 展示博主个人信息
- **组件**:
  - 个人介绍
  - 联系方式
  - 技能展示

### 2. 用户认证页面

#### 登录页面 (/login)
- **功能**: 用户登录
- **表单字段**:
  - 用户名/邮箱
  - 密码
  - 记住我选项
- **验证**: 表单验证、错误提示

#### 注册页面 (/register)
- **功能**: 用户注册
- **表单字段**:
  - 用户名
  - 邮箱
  - 密码
  - 确认密码
- **验证**: 表单验证、重复性检查

#### 个人资料页面 (/profile)
- **功能**: 查看和编辑个人信息
- **表单字段**:
  - 昵称
  - 头像
  - 个人简介
  - 邮箱
  - 密码修改

### 3. 管理后台页面

#### 仪表盘 (/admin/dashboard)
- **功能**: 展示网站统计数据
- **组件**:
  - 数据统计卡片（文章数、用户数、评论数、访问量）
  - 图表展示（访问趋势、文章发布趋势）
  - 最新评论列表
  - 系统信息

#### 文章管理 (/admin/posts)
- **功能**: 文章 CRUD 操作
- **组件**:
  - 文章列表表格（标题、状态、分类、发布时间、操作）
  - 搜索和筛选功能
  - 批量操作
  - 文章编辑器（Markdown）

#### 分类管理 (/admin/categories)
- **功能**: 分类 CRUD 操作
- **组件**:
  - 分类列表表格
  - 添加/编辑分类表单

#### 标签管理 (/admin/tags)
- **功能**: 标签 CRUD 操作
- **组件**:
  - 标签列表表格
  - 添加/编辑标签表单

#### 评论管理 (/admin/comments)
- **功能**: 评论审核和管理
- **组件**:
  - 评论列表表格
  - 评论状态切换
  - 批量操作

#### 用户管理 (/admin/users)
- **功能**: 用户管理（仅超级管理员）
- **组件**:
  - 用户列表表格
  - 用户状态管理

## 功能需求

### 1. 用户认证系统
- 用户注册、登录、登出
- JWT Token 自动管理
- 登录状态持久化
- 权限控制（普通用户/管理员）
- 个人资料管理

### 2. 文章展示系统
- 文章列表展示（支持分页）
- 文章详情展示
- Markdown 渲染
- 代码高亮
- 图片预览
- 文章搜索
- 文章分类筛选
- 文章标签筛选
- 文章归档

### 3. 评论系统
- 评论列表展示
- 添加评论
- 评论回复（多级回复）
- 评论点赞
- 评论管理

### 4. 管理后台
- 仪表盘数据展示
- 文章 CRUD 操作
- Markdown 编辑器
- 图片上传
- 分类和标签管理
- 评论管理
- 用户管理

### 5. 用户体验功能
- 响应式设计
- 主题切换（明暗模式）
- 搜索功能
- 返回顶部
- 加载状态
- 错误处理
- 无限滚动/分页
- 面包屑导航

## 状态管理设计

### Auth Store (认证状态)
```typescript
interface AuthState {
  user: User | null
  token: string | null
  isLoggedIn: boolean
  permissions: string[]
}

interface AuthActions {
  login(credentials: LoginForm): Promise<void>
  logout(): void
  register(userData: RegisterForm): Promise<void>
  updateProfile(profile: ProfileForm): Promise<void>
  refreshToken(): Promise<void>
}
```

### Blog Store (博客状态)
```typescript
interface BlogState {
  posts: Post[]
  currentPost: Post | null
  categories: Category[]
  tags: Tag[]
  comments: Comment[]
  loading: boolean
  pagination: PaginationInfo
}

interface BlogActions {
  fetchPosts(params?: QueryParams): Promise<void>
  fetchPostById(id: number): Promise<void>
  fetchCategories(): Promise<void>
  fetchTags(): Promise<void>
  fetchComments(postId: number): Promise<void>
  addComment(comment: CommentForm): Promise<void>
}
```

### Theme Store (主题状态)
```typescript
interface ThemeState {
  isDark: boolean
  primaryColor: string
  fontSize: 'small' | 'medium' | 'large'
}

interface ThemeActions {
  toggleTheme(): void
  setPrimaryColor(color: string): void
  setFontSize(size: string): void
}
```

## 类型定义

### 用户相关类型
```typescript
interface User {
  id: number
  username: string
  email: string
  nickname: string
  avatar: string
  bio: string
  role: 'admin' | 'user'
  status: number
  created_at: string
  updated_at: string
}

interface LoginForm {
  username: string
  password: string
  remember?: boolean
}

interface RegisterForm {
  username: string
  email: string
  password: string
  confirmPassword: string
}
```

### 文章相关类型
```typescript
interface Post {
  id: number
  title: string
  content: string
  summary: string
  cover: string
  status: number
  is_top: boolean
  view_count: number
  like_count: number
  user_id: number
  category_id: number
  published_at: string
  created_at: string
  updated_at: string
  user: User
  category: Category
  tags: Tag[]
  comments?: Comment[]
}

interface Category {
  id: number
  name: string
  description: string
  color: string
  sort: number
  post_count: number
  created_at: string
  updated_at: string
}

interface Tag {
  id: number
  name: string
  color: string
  post_count: number
  created_at: string
  updated_at: string
}
```

### 评论相关类型
```typescript
interface Comment {
  id: number
  content: string
  post_id: number
  user_id: number
  parent_id?: number
  status: number
  created_at: string
  updated_at: string
  user: User
  parent?: Comment
  children?: Comment[]
}
```

## 路由设计

### 博客前台路由
```typescript
const blogRoutes = [
  { path: '/', component: Home, name: 'Home' },
  { path: '/post/:id', component: PostDetail, name: 'PostDetail' },
  { path: '/category/:id', component: Category, name: 'Category' },
  { path: '/tag/:id', component: Tag, name: 'Tag' },
  { path: '/archive', component: Archive, name: 'Archive' },
  { path: '/about', component: About, name: 'About' },
]
```

### 认证路由
```typescript
const authRoutes = [
  { path: '/login', component: Login, name: 'Login' },
  { path: '/register', component: Register, name: 'Register' },
  { path: '/profile', component: Profile, name: 'Profile', meta: { requiresAuth: true } },
]
```

### 管理后台路由
```typescript
const adminRoutes = [
  {
    path: '/admin',
    component: AdminLayout,
    meta: { requiresAuth: true, requiresAdmin: true },
    children: [
      { path: 'dashboard', component: Dashboard, name: 'AdminDashboard' },
      { path: 'posts', component: PostManage, name: 'PostManage' },
      { path: 'categories', component: CategoryManage, name: 'CategoryManage' },
      { path: 'tags', component: TagManage, name: 'TagManage' },
      { path: 'comments', component: CommentManage, name: 'CommentManage' },
      { path: 'users', component: UserManage, name: 'UserManage' },
    ]
  }
]
```

## UI/UX 设计要求

### 1. 设计风格
- 现代简洁的设计风格
- 响应式布局（支持桌面端、平板、移动端）
- 统一的色彩方案
- 优雅的动画效果
- 良好的可访问性

### 2. 主题系统
- 支持明暗两种主题
- 主题切换动画
- 主题色彩自定义
- 字体大小调节

### 3. 交互体验
- 流畅的页面切换
- 友好的加载状态
- 清晰的错误提示
- 直观的操作反馈
- 键盘快捷键支持

### 4. 性能优化
- 路由懒加载
- 图片懒加载
- 组件按需加载
- 虚拟滚动（长列表）
- 缓存策略

## 开发配置

### Vite 配置
```typescript
// vite.config.ts
export default defineConfig({
  plugins: [vue(), vueJsx()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
    },
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['vue', 'vue-router', 'pinia'],
          ui: ['naive-ui'],
          editor: ['@kangc/v-md-editor'],
        },
      },
    },
  },
})
```

### ESLint 配置
```javascript
// .eslintrc.js
module.exports = {
  extends: [
    '@vue/eslint-config-typescript',
    '@vue/eslint-config-prettier',
  ],
  rules: {
    'vue/multi-word-component-names': 'off',
    '@typescript-eslint/no-unused-vars': 'error',
    'prettier/prettier': 'error',
  },
}
```

### Prettier 配置
```json
// .prettierrc
{
  "semi": false,
  "singleQuote": true,
  "tabWidth": 2,
  "trailingComma": "es5",
  "printWidth": 100
}
```

## 第三方依赖

### 核心依赖
```json
{
  "vue": "^3.3.0",
  "vue-router": "^4.2.0",
  "pinia": "^2.1.0",
  "axios": "^1.4.0",
  "naive-ui": "^2.34.0",
  "@kangc/v-md-editor": "^2.3.0",
  "typescript": "^5.0.0",
  "vite": "^4.3.0"
}
```

### 开发依赖
```json
{
  "@vitejs/plugin-vue": "^4.2.0",
  "@vue/eslint-config-prettier": "^7.1.0",
  "@vue/eslint-config-typescript": "^11.0.0",
  "eslint": "^8.40.0",
  "prettier": "^2.8.0",
  "sass": "^1.62.0",
  "unplugin-auto-import": "^0.16.0",
  "unplugin-vue-components": "^0.25.0"
}
```

### 可选依赖
- `@vueuse/core` - Vue 组合式工具库
- `dayjs` - 日期处理库
- `lodash-es` - 工具函数库
- `nprogress` - 页面加载进度条
- `vue-virtual-scroller` - 虚拟滚动
- `@headlessui/vue` - 无样式 UI 组件

## 部署要求

### 1. 构建配置
- 生产环境构建优化
- 静态资源压缩
- 代码分割
- 浏览器兼容性

### 2. 部署方式
- 静态文件部署（Nginx）
- CDN 加速
- Docker 容器化
- CI/CD 自动部署

### 3. 环境变量
```bash
# .env.production
VITE_API_BASE_URL=https://api.yourblog.com
VITE_APP_TITLE=我的博客
VITE_UPLOAD_URL=https://cdn.yourblog.com
```

## 测试要求

### 1. 单元测试
- 组件测试
- 工具函数测试
- Store 测试

### 2. 集成测试
- 页面功能测试
- API 集成测试
- 路由测试

### 3. E2E 测试
- 用户流程测试
- 跨浏览器测试
- 性能测试
