# 前端部署指南

## 📦 安装依赖

```bash
cd blog-frontend
npm install
```

## 🚀 启动开发服务器

```bash
npm run dev
```

应用将在 http://localhost:3000 启动

## 📝 注意事项

### 1. 需要添加的依赖包

由于项目使用了一些 Naive UI 的图标库，需要安装：

```bash
npm install @vicons/ionicons5
```

### 2. 路由类型声明

在 `src/types/` 目录下添加 `router.d.ts`：

```typescript
import 'vue-router'

declare module 'vue-router' {
  interface RouteMeta {
    title?: string
    requiresAuth?: boolean
    requiresAdmin?: boolean
  }
}
```

### 3. 缺失的页面组件

以下页面组件需要根据实际需求补充实现：

**博客页面**:
- `src/pages/blog/PostDetail.vue` - 文章详情页
- `src/pages/blog/Category.vue` - 分类页面
- `src/pages/blog/Tag.vue` - 标签页面  
- `src/pages/blog/Archive.vue` - 归档页面
- `src/pages/blog/About.vue` - 关于页面

**用户页面**:
- `src/pages/auth/Profile.vue` - 个人资料页

**管理后台页面**:
- `src/pages/admin/Dashboard.vue` - 仪表盘
- `src/pages/admin/PostManage.vue` - 文章管理
- `src/pages/admin/CategoryManage.vue` - 分类管理
- `src/pages/admin/TagManage.vue` - 标签管理
- `src/pages/admin/CommentManage.vue` - 评论管理
- `src/pages/admin/UserManage.vue` - 用户管理

### 4. 后端连接

确保后端服务已启动在 `http://localhost:8080`

## 🔧 快速创建缺失的页面

### 简单占位页面模板

```vue
<template>
  <div class="page">
    <n-card>
      <h2>{{ title }}</h2>
      <p>页面开发中...</p>
    </n-card>
  </div>
</template>

<script setup lang="ts">
const title = '页面标题'
</script>

<style scoped>
.page {
  max-width: 1200px;
  margin: 0 auto;
}
</style>
```

## 📊 已完成功能

- ✅ 完整的项目配置
- ✅ 类型定义系统
- ✅ API 接口封装
- ✅ 状态管理
- ✅ 路由配置
- ✅ 三种布局组件
- ✅ 登录/注册页面
- ✅ 首页（文章列表）
- ✅ 404 页面

## 🎯 最小可运行版本

为了快速运行，可以先创建简单的占位页面，然后逐步完善功能。

### 构建生产版本

```bash
npm run build
```

### 预览生产版本

```bash
npm run preview
```

## 🌐 部署到服务器

### 使用 Nginx

1. 构建生产版本
2. 将 `dist/` 目录内容上传到服务器
3. 配置 Nginx:

```nginx
server {
    listen 80;
    server_name yourdomain.com;
    
    location / {
        root /var/www/blog-frontend;
        try_files $uri $uri/ /index.html;
    }
    
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 📞 故障排查

### 问题1: 依赖安装失败
```bash
# 清除缓存重新安装
rm -rf node_modules package-lock.json
npm install
```

### 问题2: 路由404
检查 Vite 配置中的 base 路径和服务器配置

### 问题3: API 请求失败
检查代理配置和后端服务状态

## 💡 开发建议

1. 先运行后端服务
2. 安装前端依赖
3. 启动前端开发服务器
4. 逐步完善页面功能

