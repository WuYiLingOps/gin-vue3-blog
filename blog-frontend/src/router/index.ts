// 路由配置

import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { setupRouterGuards } from './guards'

// 布局
const DefaultLayout = () => import('@/layouts/DefaultLayout.vue')
const AdminLayout = () => import('@/layouts/AdminLayout.vue')
const AuthLayout = () => import('@/layouts/AuthLayout.vue')

// 博客页面
const Home = () => import('@/pages/blog/Home.vue')
const PostDetail = () => import('@/pages/blog/PostDetail.vue')
const Category = () => import('@/pages/blog/Category.vue')
const Tag = () => import('@/pages/blog/Tag.vue')
const Archive = () => import('@/pages/blog/Archive.vue')
const Moments = () => import('@/pages/blog/Moments.vue')

// 认证页面
const Login = () => import('@/pages/auth/Login.vue')
const Register = () => import('@/pages/auth/Register.vue')
const Profile = () => import('@/pages/auth/Profile.vue')
const ForgotPassword = () => import('@/pages/auth/ForgotPassword.vue')

// 管理后台页面
const Dashboard = () => import('@/pages/admin/Dashboard.vue')
const PostManage = () => import('@/pages/admin/PostManage.vue')
const PostEdit = () => import('@/pages/admin/PostEdit.vue')
const CategoryManage = () => import('@/pages/admin/CategoryManage.vue')
const TagManage = () => import('@/pages/admin/TagManage.vue')
const CommentManage = () => import('@/pages/admin/CommentManage.vue')
const UserManage = () => import('@/pages/admin/UserManage.vue')
const SiteSettings = () => import('@/pages/admin/SiteSettings.vue')
const MomentManage = () => import('@/pages/admin/MomentManage.vue')

const routes: RouteRecordRaw[] = [
  // 博客前台路由
  {
    path: '/',
    component: DefaultLayout,
    children: [
      {
        path: '',
        name: 'Home',
        component: Home,
        meta: { title: '首页' }
      },
      {
        path: 'post/:id',
        name: 'PostDetail',
        component: PostDetail,
        meta: { title: '文章详情' }
      },
      {
        path: 'category',
        name: 'Category',
        component: Category,
        meta: { title: '分类' }
      },
      {
        path: 'category/:id',
        name: 'CategoryDetail',
        component: Category,
        meta: { title: '分类详情' }
      },
      {
        path: 'tag',
        name: 'Tag',
        component: Tag,
        meta: { title: '标签' }
      },
      {
        path: 'tag/:id',
        name: 'TagDetail',
        component: Tag,
        meta: { title: '标签详情' }
      },
      {
        path: 'archive',
        name: 'Archive',
        component: Archive,
        meta: { title: '归档' }
      },
      {
        path: 'moments',
        name: 'Moments',
        component: Moments,
        meta: { title: '说说' }
      }
    ]
  },

  // 认证路由
  {
    path: '/auth',
    component: AuthLayout,
    children: [
      {
        path: 'login',
        name: 'Login',
        component: Login,
        meta: { title: '登录' }
      },
      {
        path: 'register',
        name: 'Register',
        component: Register,
        meta: { title: '注册' }
      },
      {
        path: 'forgot-password',
        name: 'ForgotPassword',
        component: ForgotPassword,
        meta: { title: '找回密码' }
      }
    ]
  },

  // 个人资料（需要认证）
  {
    path: '/profile',
    component: DefaultLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Profile',
        component: Profile,
        meta: { title: '个人资料', requiresAuth: true }
      }
    ]
  },

  // 管理后台路由
  {
    path: '/admin',
    component: AdminLayout,
    meta: { requiresAuth: true, requiresAdmin: true },
    children: [
      {
        path: '',
        redirect: '/admin/dashboard'
      },
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: Dashboard,
        meta: { title: '仪表盘', requiresAuth: true, requiresAdmin: true }
      },
      {
        path: 'posts',
        name: 'PostManage',
        component: PostManage,
        meta: { title: '文章管理', requiresAuth: true, requiresAdmin: true }
      },
      {
        path: 'posts/edit/:id',
        name: 'PostEdit',
        component: PostEdit,
        meta: { title: '编辑文章', requiresAuth: true, requiresAdmin: true }
      },
      {
        path: 'categories',
        name: 'CategoryManage',
        component: CategoryManage,
        meta: { title: '分类管理', requiresAuth: true, requiresAdmin: true }
      },
      {
        path: 'tags',
        name: 'TagManage',
        component: TagManage,
        meta: { title: '标签管理', requiresAuth: true, requiresAdmin: true }
      },
      {
        path: 'comments',
        name: 'CommentManage',
        component: CommentManage,
        meta: { title: '评论管理', requiresAuth: true, requiresAdmin: true }
      },
      {
        path: 'users',
        name: 'UserManage',
        component: UserManage,
        meta: { title: '用户管理', requiresAuth: true, requiresAdmin: true }
      },
      {
        path: 'site',
        name: 'SiteSettings',
        component: SiteSettings,
        meta: { title: '网站设置', requiresAuth: true, requiresAdmin: true }
      },
      {
        path: 'moments',
        name: 'MomentManage',
        component: MomentManage,
        meta: { title: '说说管理', requiresAuth: true, requiresAdmin: true }
      }
    ]
  },

  // 404 页面
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/pages/NotFound.vue'),
    meta: { title: '404' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(_to, _from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// 设置路由守卫
setupRouterGuards(router)

export default router

