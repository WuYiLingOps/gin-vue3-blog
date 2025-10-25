<template>
  <div class="default-layout">
    <n-layout position="absolute">
      <!-- 头部 -->
      <n-layout-header class="header" position="absolute">
        <div class="header-content">
          <div class="logo" @click="router.push('/')">
            <h2>{{ siteSettings.site_name || '我的博客' }}</h2>
          </div>

          <n-menu
            v-model:value="activeKey"
            mode="horizontal"
            :options="menuOptions"
            class="nav-menu"
            @update:value="handleMenuSelect"
          />

          <div class="header-actions">
            <!-- 搜索按钮 -->
            <n-button text @click="showSearchModal = true">
              <template #icon>
                <n-icon :component="SearchOutline" size="20" />
              </template>
            </n-button>

            <n-button text @click="toggleTheme">
              <template #icon>
                <n-icon :component="isDark ? SunnyOutline : MoonOutline" />
              </template>
            </n-button>

            <n-dropdown v-if="authStore.isLoggedIn" :options="userMenuOptions" @select="handleUserMenu">
              <n-button text>
                <n-avatar round size="small" :src="authStore.user?.avatar" />
                <span class="ml-2">{{ authStore.user?.nickname }}</span>
              </n-button>
            </n-dropdown>

            <n-button v-else type="primary" @click="router.push('/auth/login')">
              登录
            </n-button>
          </div>
        </div>
      </n-layout-header>

      <!-- 主体内容 -->
      <n-layout-content class="main-content" content-style="padding: 0;" :native-scrollbar="false">
        <div class="content-wrapper">
          <router-view />
        </div>
        
        <!-- 底部 -->
        <div class="footer">
          <div class="footer-content">
            <p>&copy; 2025 {{ siteSettings.site_name || '我的博客' }}. All rights reserved.</p>
            <p class="running-time">{{ runningTime }}</p>
            <div v-if="siteSettings.site_icp || siteSettings.site_police" class="filing-info">
              <a v-if="siteSettings.site_icp" href="https://beian.miit.gov.cn/" target="_blank" rel="noopener noreferrer">
                {{ siteSettings.site_icp }}
              </a>
              <a v-if="siteSettings.site_police" href="https://www.beian.gov.cn/" target="_blank" rel="noopener noreferrer">
                {{ siteSettings.site_police }}
              </a>
            </div>
          </div>
        </div>
      </n-layout-content>
    </n-layout>

    <!-- 搜索对话框 -->
    <n-modal
      v-model:show="showSearchModal"
      preset="card"
      title="搜索文章"
      style="width: 800px; max-width: 90vw; margin-top: 10vh"
      :bordered="false"
      :segmented="false"
      @mask-click="showSearchModal = false"
    >
      <div class="search-modal-content">
        <n-input
          v-model:value="searchKeyword"
          placeholder="输入关键词搜索文章..."
          size="large"
          clearable
          autofocus
          @input="handleSearchInput"
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <n-icon :component="SearchOutline" />
          </template>
        </n-input>

        <!-- 搜索结果列表 -->
        <div v-if="searchResults.length > 0 || (searchKeyword && searchLoading)" class="search-results">
          <n-divider style="margin: 20px 0" />
          <n-spin :show="searchLoading">
            <div class="search-result-item" v-for="post in searchResults" :key="post.id" @click="goToPost(post.id)">
              <div class="result-title" v-html="highlightText(post.title)"></div>
              <div class="result-meta">
                <span>{{ post.category.name }}</span>
                <n-divider vertical />
                <span>{{ formatDate(post.created_at, 'YYYY-MM-DD') }}</span>
                <n-divider vertical />
                <span>{{ post.view_count }} 阅读</span>
              </div>
              <div v-if="post.summary" class="result-summary" v-html="highlightText(post.summary)"></div>
            </div>
            <n-empty v-if="searchKeyword && !searchLoading && searchResults.length === 0" description="未找到相关文章" style="margin: 32px 0" />
          </n-spin>
        </div>

        <!-- 空状态占位 -->
        <div v-else-if="!searchKeyword" class="search-empty-placeholder">
          <n-icon :component="SearchOutline" size="64" style="color: #d1d5db; margin-bottom: 16px" />
          <p style="color: #9ca3af; font-size: 15px; margin: 0">输入关键词开始搜索文章</p>
        </div>

        <!-- 提示信息 -->
        <div v-if="searchResults.length > 0" class="search-footer">
          <span class="search-count">找到 {{ searchResults.length }} 篇文章</span>
          <span class="search-hint">点击文章标题查看详情</span>
        </div>
      </div>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h, onMounted, onBeforeUnmount } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { MoonOutline, SunnyOutline, PersonOutline, LogOutOutline, SettingsOutline, SearchOutline } from '@vicons/ionicons5'
import { useAuthStore, useAppStore } from '@/stores'
import { NIcon } from 'naive-ui'
import { getPublicSettings } from '@/api/setting'
import type { SiteSettings } from '@/api/setting'
import { getPosts } from '@/api/post'
import { formatDate } from '@/utils/format'
import { highlightKeyword } from '@/utils/highlight'
import type { Post } from '@/types/blog'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const appStore = useAppStore()

const activeKey = ref(route.name as string)
const isDark = computed(() => appStore.theme === 'dark')
const siteSettings = ref<SiteSettings>({})
const runningTime = ref('')
const searchKeyword = ref('')
const showSearchModal = ref(false)
const searchResults = ref<Post[]>([])
const searchLoading = ref(false)
let searchTimer: number | null = null

// 网站启动时间（可以在这里设置你的网站上线日期）
const siteStartDate = new Date('2025-10-25 00:00:00')

// 菜单选项
const menuOptions = [
  {
    label: '首页',
    key: 'Home',
    path: '/'
  },
  {
    label: '归档',
    key: 'Archive',
    path: '/archive'
  },
  {
    label: '说说',
    key: 'Moments',
    path: '/moments'
  },
  {
    label: '关于',
    key: 'About',
    path: '/about'
  }
]

// 用户菜单选项
const userMenuOptions = computed(() => {
  const options = [
    {
      label: '个人资料',
      key: 'profile',
      icon: () => h(NIcon, null, { default: () => h(PersonOutline) })
    }
  ]

  if (authStore.isAdmin) {
    options.push({
      label: '管理后台',
      key: 'admin',
      icon: () => h(NIcon, null, { default: () => h(SettingsOutline) })
    })
  }

  options.push({
    label: '退出登录',
    key: 'logout',
    icon: () => h(NIcon, null, { default: () => h(LogOutOutline) })
  })

  return options
})

// 处理菜单选择
function handleMenuSelect(key: string) {
  activeKey.value = key
  router.push({ name: key })
}

// 切换主题
function toggleTheme() {
  appStore.toggleTheme()
}

// 获取网站配置
async function fetchSiteSettings() {
  try {
    const res = await getPublicSettings()
    if (res.data) {
      siteSettings.value = res.data
      // 获取配置后更新页面标题
      updatePageTitle()
    }
  } catch (error) {
    console.error('获取网站配置失败:', error)
  }
}

// 计算网站运行时间
function calculateRunningTime() {
  const now = new Date()
  const diff = now.getTime() - siteStartDate.getTime()
  
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  const seconds = Math.floor((diff % (1000 * 60)) / 1000)
  
  runningTime.value = `网站已运行 ${days} 天 ${hours} 小时 ${minutes} 分钟 ${seconds} 秒`
}

let timer: number | null = null

onMounted(() => {
  fetchSiteSettings()
  calculateRunningTime()
  // 每秒更新一次运行时间
  timer = window.setInterval(calculateRunningTime, 1000)
  
  // 更新页面标题
  updatePageTitle()
})

// 更新页面标题
function updatePageTitle() {
  const siteName = siteSettings.value.site_name || '我的博客'
  const currentTitle = route.meta.title as string || '首页'
  document.title = `${currentTitle} | ${siteName}`
}

onBeforeUnmount(() => {
  if (timer) {
    clearInterval(timer)
  }
})

// 实时搜索
async function handleSearchInput() {
  // 清除之前的定时器
  if (searchTimer) {
    clearTimeout(searchTimer)
  }

  const keyword = searchKeyword.value.trim()
  
  if (!keyword) {
    searchResults.value = []
    return
  }

  // 防抖：延迟500ms执行搜索
  searchTimer = window.setTimeout(async () => {
    try {
      searchLoading.value = true
      const res = await getPosts({
        page: 1,
        page_size: 10,
        keyword: keyword,
        status: 1
      })
      
      if (res.data) {
        searchResults.value = res.data.list
      }
    } catch (error) {
      console.error('搜索失败:', error)
      searchResults.value = []
    } finally {
      searchLoading.value = false
    }
  }, 500)
}

// 跳转到文章详情
function goToPost(postId: number) {
  showSearchModal.value = false
  router.push(`/post/${postId}`)
  // 清空搜索
  setTimeout(() => {
    searchKeyword.value = ''
    searchResults.value = []
  }, 300)
}

// 查看全部搜索结果
function handleSearch() {
  if (searchKeyword.value.trim()) {
    showSearchModal.value = false
    router.push({
      path: '/',
      query: { keyword: searchKeyword.value }
    })
    // 清空搜索框
    setTimeout(() => {
      searchKeyword.value = ''
      searchResults.value = []
    }, 300)
  }
}

// 高亮文本
function highlightText(text: string): string {
  if (!searchKeyword.value || !text) {
    return text || ''
  }
  return highlightKeyword(text, searchKeyword.value)
}

// 处理用户菜单
function handleUserMenu(key: string) {
  switch (key) {
    case 'profile':
      router.push('/profile')
      break
    case 'admin':
      router.push('/admin')
      break
    case 'logout':
      authStore.logout()
      router.push('/')
      break
  }
}
</script>

<style scoped>
.default-layout {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow: hidden;
}

.default-layout :deep(.n-layout) {
  height: 100vh;
}

/* 玻璃态顶部导航栏 */
.header {
  padding: 0 24px;
  height: 72px;
  display: flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(20px) saturate(180%);
  border-bottom: 1px solid rgba(8, 145, 178, 0.1);
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.06);
  position: sticky;
  top: 0;
  z-index: 100;
  transition: all 0.3s;
}

html.dark .header {
  background: rgba(15, 23, 42, 0.8);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.3);
}

.header-content {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  cursor: pointer;
  display: flex;
  align-items: center;
  transition: transform 0.3s;
}

.logo:hover {
  transform: scale(1.05);
}

.logo h2 {
  margin: 0;
  font-size: 26px;
  font-weight: 800;
  background: linear-gradient(135deg, #0891b2 0%, #059669 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: -0.02em;
}

html.dark .logo h2 {
  background: linear-gradient(135deg, #38bdf8 0%, #4ade80 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.nav-menu {
  flex: 1;
  margin: 0 48px;
  max-width: 400px;
}

/* 自定义导航菜单样式 */
.nav-menu :deep(.n-menu-item) {
  font-weight: 600;
  transition: all 0.3s;
}

.nav-menu :deep(.n-menu-item:hover) {
  color: #0891b2;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-actions :deep(.n-button) {
  font-weight: 600;
  transition: all 0.3s;
}

.main-content {
  padding: 32px 24px 0 24px;
  position: relative;
  z-index: 1;
  overflow-y: auto;
  height: calc(100vh - 72px);
}

.main-content :deep(.n-scrollbar-content) {
  min-height: 100%;
}

.content-wrapper {
  max-width: 1400px;
  margin: 0 auto;
  min-height: calc(100vh - 72px - 180px);
  padding-bottom: 20px;
}

/* 玻璃态底部 */
.footer {
  padding: 32px 24px;
  text-align: center;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px) saturate(180%);
  border-top: 1px solid rgba(8, 145, 178, 0.1);
  box-shadow: 0 -4px 24px rgba(0, 0, 0, 0.04);
  margin-top: 40px;
}

html.dark .footer {
  background: rgba(15, 23, 42, 0.7);
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 0 -4px 24px rgba(0, 0, 0, 0.2);
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
}

.footer-content p {
  margin: 6px 0;
  color: #64748b;
  font-size: 14px;
  font-weight: 500;
}

html.dark .footer-content p {
  color: #94a3b8;
}

.running-time {
  font-family: 'Courier New', Consolas, monospace;
  font-size: 13px !important;
  color: #0891b2 !important;
  font-weight: 600;
  opacity: 0.9 !important;
}

html.dark .running-time {
  color: #38bdf8 !important;
}

.filing-info {
  margin-top: 8px;
  font-size: 13px;
  display: flex;
  gap: 16px;
  justify-content: center;
  flex-wrap: wrap;
}

.filing-info a {
  color: rgba(8, 145, 178, 0.8);
  text-decoration: none;
  transition: all 0.3s;
}

.filing-info a:hover {
  color: rgba(8, 145, 178, 1);
}

html.dark .filing-info a {
  color: rgba(56, 189, 248, 0.8);
}

html.dark .filing-info a:hover {
  color: rgba(56, 189, 248, 1);
}

.ml-2 {
  margin-left: 8px;
}

/* 添加头像悬停效果 */
.header-actions :deep(.n-avatar) {
  transition: all 0.3s;
  border: 2px solid rgba(8, 145, 178, 0.2);
}

.header-actions :deep(.n-avatar:hover) {
  transform: scale(1.1);
  border-color: rgba(8, 145, 178, 0.5);
  box-shadow: 0 4px 12px rgba(8, 145, 178, 0.3);
}

/* 搜索模态框内容 */
.search-modal-content {
  min-height: 280px;
}

/* 搜索结果样式 */
.search-results {
  max-height: 500px;
  overflow-y: auto;
  margin-bottom: 16px;
}

/* 空状态占位 */
.search-empty-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  min-height: 240px;
}

html.dark .search-empty-placeholder p {
  color: #6b7280 !important;
}

html.dark .search-empty-placeholder :deep(.n-icon) {
  color: #4b5563 !important;
}

.search-result-item {
  padding: 16px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  margin-bottom: 12px;
  border: 1px solid transparent;
}

.search-result-item:hover {
  background: rgba(8, 145, 178, 0.05);
  border-color: rgba(8, 145, 178, 0.2);
  transform: translateX(4px);
}

html.dark .search-result-item:hover {
  background: rgba(56, 189, 248, 0.1);
  border-color: rgba(56, 189, 248, 0.2);
}

.result-title {
  font-size: 16px;
  font-weight: 600;
  color: #1a202c;
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

html.dark .result-title {
  color: #e5e5e5;
}

.result-meta {
  font-size: 13px;
  color: #64748b;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
}

html.dark .result-meta {
  color: #94a3b8;
}

.result-summary {
  font-size: 14px;
  color: #475569;
  line-height: 1.6;
  overflow: hidden;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  text-overflow: ellipsis;
}

html.dark .result-summary {
  color: #cbd5e1;
}

.search-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0 8px 0;
  border-top: 1px solid rgba(8, 145, 178, 0.1);
}

html.dark .search-footer {
  border-top-color: rgba(56, 189, 248, 0.1);
}

.search-count {
  font-size: 14px;
  font-weight: 600;
  color: #0891b2;
}

html.dark .search-count {
  color: #38bdf8;
}

.search-hint {
  font-size: 13px;
  color: #94a3b8;
}

html.dark .search-hint {
  color: #64748b;
}

/* 搜索高亮样式 */
:deep(.search-highlight) {
  background: linear-gradient(120deg, #fef08a 0%, #fde047 100%);
  color: #854d0e;
  padding: 2px 4px;
  border-radius: 3px;
  font-weight: 600;
  box-shadow: 0 1px 3px rgba(251, 191, 36, 0.3);
}

html.dark :deep(.search-highlight) {
  background: linear-gradient(120deg, #fbbf24 0%, #f59e0b 100%);
  color: #1f2937;
  box-shadow: 0 1px 3px rgba(251, 191, 36, 0.5);
}
</style>

