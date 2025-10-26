<template>
  <div class="default-layout">
    <n-layout position="absolute">
      <!-- 头部 -->
      <n-layout-header class="header" position="absolute">
        <div class="header-content">
          <div class="logo" @click="router.push('/')">
            <h2>{{ siteSettings.site_name || '我的博客' }}</h2>
          </div>

          <!-- 桌面端导航菜单 -->
          <n-menu
            v-model:value="activeKey"
            mode="horizontal"
            :options="menuOptions"
            class="nav-menu desktop-only"
            @update:value="handleMenuSelect"
          />

          <div class="header-actions">
            <!-- 移动端菜单按钮 -->
            <n-button text class="mobile-menu-btn" @click="showMobileMenu = true">
              <template #icon>
                <n-icon :component="MenuOutline" size="24" />
              </template>
            </n-button>
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
            <p class="running-time" v-html="runningTime"></p>
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

    <!-- 移动端侧边抽屉菜单 -->
    <n-drawer v-model:show="showMobileMenu" :width="280" placement="left">
      <n-drawer-content title="菜单" closable>
        <n-menu
          v-model:value="activeKey"
          :options="menuOptions"
          @update:value="handleMobileMenuSelect"
        />
        
        <n-divider style="margin: 24px 0" />
        
        <!-- 用户信息 -->
        <div v-if="authStore.isLoggedIn" class="mobile-user-info">
          <n-space vertical :size="12">
            <div class="user-profile">
              <n-avatar :src="authStore.user?.avatar" :size="48" round />
              <div class="user-details">
                <div class="user-nickname">{{ authStore.user?.nickname }}</div>
                <div class="user-role">{{ authStore.isAdmin ? '管理员' : '普通用户' }}</div>
              </div>
            </div>
            
            <n-button block @click="handleMobileUserAction('profile')">
              <template #icon>
                <n-icon :component="PersonOutline" />
              </template>
              个人资料
            </n-button>
            
            <n-button v-if="authStore.isAdmin" block type="info" @click="handleMobileUserAction('admin')">
              <template #icon>
                <n-icon :component="SettingsOutline" />
              </template>
              管理后台
            </n-button>
            
            <n-button block type="error" @click="handleMobileUserAction('logout')">
              <template #icon>
                <n-icon :component="LogOutOutline" />
              </template>
              退出登录
            </n-button>
          </n-space>
        </div>
        
        <!-- 未登录状态 -->
        <div v-else class="mobile-login">
          <n-button block type="primary" @click="handleMobileLogin">
            登录 / 注册
          </n-button>
        </div>
      </n-drawer-content>
    </n-drawer>

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
              <div class="result-summary" v-html="getHighlightedSummary(post)"></div>
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
          <n-button type="primary" size="small" @click="handleSearch">
            查看全部搜索结果
          </n-button>
        </div>
      </div>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h, onMounted, onBeforeUnmount } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { MoonOutline, SunnyOutline, PersonOutline, LogOutOutline, SettingsOutline, SearchOutline, MenuOutline, HomeOutline, ArchiveOutline, ChatbubblesOutline, InformationCircleOutline } from '@vicons/ionicons5'
import { useAuthStore, useAppStore } from '@/stores'
import { NIcon } from 'naive-ui'
import { getPublicSettings } from '@/api/setting'
import type { SiteSettings } from '@/api/setting'
import { getPosts } from '@/api/post'
import { formatDate } from '@/utils/format'
import { highlightKeyword, extractHighlightSnippet } from '@/utils/highlight'
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
const showMobileMenu = ref(false)
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
    path: '/',
    icon: () => h(NIcon, null, { default: () => h(HomeOutline) })
  },
  {
    label: '文章归档',
    key: 'ArchiveMenu',
    icon: () => h(NIcon, null, { default: () => h(ArchiveOutline) }),
    children: [
      {
        label: '归档',
        key: 'Archive',
        path: '/archive'
      },
      {
        label: '分类',
        key: 'Category',
        path: '/category'
      },
      {
        label: '标签',
        key: 'Tag',
        path: '/tag'
      }
    ]
  },
  {
    label: '说说',
    key: 'Moments',
    path: '/moments',
    icon: () => h(NIcon, null, { default: () => h(ChatbubblesOutline) })
  },
  {
    label: '关于',
    key: 'About',
    path: '/about',
    icon: () => h(NIcon, null, { default: () => h(InformationCircleOutline) })
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
  // 如果是父菜单（ArchiveMenu），不做任何操作
  if (key === 'ArchiveMenu') {
    return
  }
  
  activeKey.value = key
  router.push({ name: key })
}

// 处理移动端菜单选择
function handleMobileMenuSelect(key: string) {
  // 如果是父菜单（ArchiveMenu），不做任何操作
  if (key === 'ArchiveMenu') {
    return
  }
  
  activeKey.value = key
  showMobileMenu.value = false
  router.push({ name: key })
}

// 处理移动端用户操作
function handleMobileUserAction(action: string) {
  showMobileMenu.value = false
  handleUserMenu(action)
}

// 处理移动端登录
function handleMobileLogin() {
  showMobileMenu.value = false
  router.push('/auth/login')
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
  
  runningTime.value = `网站已运行 <span class="time-number">${days}</span> 天 <span class="time-number">${hours}</span> 小时 <span class="time-number">${minutes}</span> 分钟 <span class="time-number">${seconds}</span> 秒`
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

// 高亮摘要（如果摘要中没有关键词，则从内容中提取）
function getHighlightedSummary(post: Post): string {
  const summary = post.summary || ''
  const keyword = searchKeyword.value
  
  if (!keyword) {
    return summary
  }
  
  // 检查摘要中是否包含关键词
  const lowerSummary = summary.toLowerCase()
  const lowerKeyword = keyword.toLowerCase()
  
  if (lowerSummary.includes(lowerKeyword)) {
    // 如果摘要中包含关键词，直接高亮
    return highlightKeyword(summary, keyword)
  } else if (post.content) {
    // 如果摘要中不包含关键词，但内容存在，从内容中提取包含关键词的片段
    const snippet = extractHighlightSnippet(post.content, keyword, 150)
    // 如果提取到了包含关键词的片段，就使用它；否则使用原摘要
    if (snippet && snippet.toLowerCase().includes(lowerKeyword)) {
      return highlightKeyword(snippet, keyword)
    }
  }
  
  // 默认返回原摘要
  return summary
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
  max-width: 500px;
  min-width: 400px;
}

/* 自定义导航菜单样式 */
.nav-menu :deep(.n-menu-item) {
  font-weight: 600;
  transition: all 0.3s;
}

.nav-menu :deep(.n-menu-item:hover) {
  color: #0891b2;
}

/* 确保菜单文字完整显示 */
.nav-menu :deep(.n-menu-item-content) {
  overflow: visible !important;
  white-space: nowrap;
}

.nav-menu :deep(.n-menu-item-content-header) {
  overflow: visible !important;
  white-space: nowrap;
}

.nav-menu :deep(.n-menu-item-content__icon) {
  margin-right: 8px;
  flex-shrink: 0;
}

.nav-menu :deep(.n-submenu-children .n-menu-item-content-header) {
  overflow: visible !important;
  white-space: nowrap;
}

/* 菜单项文本 */
.nav-menu :deep(.n-menu-item-content-header__title) {
  overflow: visible !important;
  text-overflow: clip !important;
  white-space: nowrap;
}

/* 子菜单样式 */
.nav-menu :deep(.n-submenu) {
  overflow: visible !important;
}

.nav-menu :deep(.n-submenu-children) {
  overflow: visible !important;
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

/* 玻璃态底部 - 占满全屏宽度 */
.footer {
  padding: 20px 24px;
  text-align: center;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px) saturate(180%);
  border-top: 1px solid rgba(8, 145, 178, 0.1);
  box-shadow: 0 -4px 24px rgba(0, 0, 0, 0.04);
  margin-top: 40px;
  /* 占满全屏宽度 */
  width: 100vw;
  margin-left: calc(-50vw + 50%);
  position: relative;
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
  margin: 4px 0;
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
}

.header-actions :deep(.n-avatar:hover) {
  transform: scale(1.1);
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

/* 运行时间样式 */
.running-time :deep(.time-number) {
  font-weight: 600;
  background: linear-gradient(135deg, #0891b2 0%, #059669 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  font-size: 1.05em;
  margin: 0 2px;
}

html.dark .running-time :deep(.time-number) {
  background: linear-gradient(135deg, #38bdf8 0%, #4ade80 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* ===== 移动端适配 ===== */

/* 移动端菜单按钮 - 默认隐藏 */
.mobile-menu-btn {
  display: none;
}

/* 平板以下显示移动端菜单按钮，隐藏桌面导航 */
@media (max-width: 768px) {
  .mobile-menu-btn {
    display: inline-flex;
  }
  
  .desktop-only {
    display: none !important;
  }
  
  .header {
    padding: 0 16px;
    height: 60px;
  }
  
  .logo h2 {
    font-size: 20px;
  }
  
  .header-actions {
    gap: 8px;
  }
  
  .main-content {
    padding: 16px 12px 0 12px;
    height: calc(100vh - 60px);
  }
  
  .content-wrapper {
    min-height: calc(100vh - 60px - 160px);
  }
  
  .footer {
    padding: 16px 12px;
    margin-top: 24px;
  }
  
  .footer-content p {
    font-size: 12px;
  }
  
  .running-time {
    font-size: 11px !important;
  }
}

/* 移动端用户信息样式 */
.mobile-user-info {
  padding: 16px 0;
}

.user-profile {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: rgba(8, 145, 178, 0.05);
  border-radius: 12px;
  margin-bottom: 16px;
}

html.dark .user-profile {
  background: rgba(56, 189, 248, 0.1);
}

.user-details {
  flex: 1;
}

.user-nickname {
  font-size: 16px;
  font-weight: 600;
  color: #1a202c;
  margin-bottom: 4px;
}

html.dark .user-nickname {
  color: #e5e5e5;
}

.user-role {
  font-size: 13px;
  color: #64748b;
}

html.dark .user-role {
  color: #94a3b8;
}

.mobile-login {
  padding: 16px 0;
}

/* 小屏幕优化 */
@media (max-width: 480px) {
  .logo h2 {
    font-size: 18px;
  }
  
  .header-actions :deep(.n-button .n-button__content > span) {
    display: none;
  }
}
</style>

