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
│   │   └── index.ts
│   ├── assets/            # 资源文件
│   │   └── styles/
│   │       └── global.css
│   ├── components/        # 公共组件
│   ├── layouts/          # 布局组件
│   ├── pages/            # 页面组件
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

创建 `.env.local` 文件配置本地环境变量：

```bash
VITE_APP_TITLE=我的博客
VITE_API_BASE_URL=http://localhost:8080
VITE_UPLOAD_URL=http://localhost:8080/api/upload
```

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

### 开发中

- 🚧 页面组件
- 🚧 路由配置
- 🚧 公共组件
- 🚧 布局组件

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

## 许可证

MIT License

