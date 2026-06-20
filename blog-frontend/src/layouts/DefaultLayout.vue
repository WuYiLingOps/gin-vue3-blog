<!--
 * @ProjectName: go-vue3-blog
 * @FileName: DefaultLayout.vue
 * @CreateTime: 2026-02-02 11:40:46
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 默认布局组件，提供网站的主要布局结构，包括顶部导航栏、主体内容区域、底部信息栏、移动端菜单、搜索功能等
 -->
<template>
  <div class="default-layout">
    <GlobalBackground />
    <CanvasBackground />
    <n-layout position="absolute">
      <!-- 头部 -->
      <n-layout-header
        class="header"
        :class="{ 'header-hidden': headerHidden, 'header-transparent': isOnCover }"
        :style="
          isOnCover
            ? {
                background: 'transparent',
                backdropFilter: 'none',
                borderBottom: 'none',
                boxShadow: 'none'
              }
            : {}
        "
        position="absolute"
      >
        <div class="header-content">
          <div class="logo" @click="router.push('/')">
            <h2>{{ siteSettings.site_name || defaultSiteName }}</h2>
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
            <!-- 打赏按钮 -->
            <n-button text @click="showRewardModal = true">
              <template #icon>
                <n-icon :component="CafeOutline" size="20" />
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

            <n-dropdown
              v-if="authStore.isLoggedIn"
              :options="userMenuOptions"
              @select="handleUserMenu"
            >
              <n-button text>
                <n-avatar round size="small" :src="authStore.user?.avatar" />
                <span class="ml-2">{{ authStore.user?.nickname }}</span>
              </n-button>
            </n-dropdown>

            <n-button v-else type="primary" @click="router.push('/auth/login')"> 登录 </n-button>
          </div>
        </div>
      </n-layout-header>

      <!-- 主体内容 -->
      <n-layout-content
        ref="mainContentRef"
        :class="['main-content', { 'main-content-home': route.name === 'Home' }]"
        content-style="padding: 0;"
        :native-scrollbar="false"
        @scroll="handleNaiveScroll"
      >
        <div class="content-wrapper" :class="{ 'content-wrapper-home': route.name === 'Home' }">
          <router-view />
        </div>

        <!-- 底部 -->
        <div class="footer">
          <div class="footer-content">
            <p>
              &copy; {{ copyrightYear }} {{ siteSettings.site_name || defaultSiteName }}. All rights
              reserved.
            </p>
            <p class="running-time" v-html="runningTime"></p>
            <div v-if="siteSettings.site_icp || siteSettings.site_police" class="filing-info">
              <a
                v-if="siteSettings.site_icp"
                href="https://beian.miit.gov.cn/"
                target="_blank"
                rel="noopener noreferrer"
              >
                {{ siteSettings.site_icp }}
              </a>
              <a
                v-if="siteSettings.site_police"
                href="https://www.beian.gov.cn/"
                target="_blank"
                rel="noopener noreferrer"
                class="police-filing"
              >
                <img src="/备案图标.png" alt="公安备案" class="police-icon" />
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

            <n-button
              v-if="authStore.isAdmin"
              block
              type="info"
              @click="handleMobileUserAction('admin')"
            >
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
          <n-button block type="primary" @click="handleMobileLogin"> 登录 / 注册 </n-button>
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
        <div
          v-if="searchResults.length > 0 || (searchKeyword && searchLoading)"
          class="search-results"
        >
          <n-divider style="margin: 20px 0" />
          <n-spin :show="searchLoading">
            <div
              class="search-result-item"
              v-for="post in searchResults"
              :key="post.id"
              @click="goToPost(post)"
            >
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
            <n-empty
              v-if="searchKeyword && !searchLoading && searchResults.length === 0"
              description="未找到相关文章"
              style="margin: 32px 0"
            />
          </n-spin>
        </div>

        <!-- 空状态占位 -->
        <div v-else-if="!searchKeyword" class="search-empty-placeholder">
          <n-icon
            :component="SearchOutline"
            size="64"
            style="color: #d1d5db; margin-bottom: 16px"
          />
          <p style="color: #9ca3af; font-size: 15px; margin: 0">输入关键词开始搜索文章</p>
        </div>

        <!-- 提示信息 -->
        <div v-if="searchResults.length > 0" class="search-footer">
          <span class="search-count">找到 {{ searchResults.length }} 篇文章</span>
          <n-button type="primary" size="small" @click="handleSearch"> 查看全部搜索结果 </n-button>
        </div>
      </div>
    </n-modal>

    <!-- 打赏对话框 -->
    <n-modal
      v-model:show="showRewardModal"
      preset="card"
      title="请作者喝杯咖啡 ☕"
      style="width: 420px; max-width: 90vw"
      :bordered="false"
      :segmented="false"
      @mask-click="showRewardModal = false"
    >
      <div class="reward-modal-content">
        <template v-if="siteSettings.reward_wechat || siteSettings.reward_alipay">
          <div class="reward-qrcodes">
            <div v-if="siteSettings.reward_wechat" class="reward-item">
              <img :src="siteSettings.reward_wechat" alt="微信收款码" class="reward-qrcode" />
              <p class="reward-label">微信</p>
            </div>
            <div v-if="siteSettings.reward_alipay" class="reward-item">
              <img :src="siteSettings.reward_alipay" alt="支付宝收款码" class="reward-qrcode" />
              <p class="reward-label">支付宝</p>
            </div>
          </div>
          <p class="reward-tip">感谢您的支持与鼓励！</p>
        </template>
        <template v-else>
          <n-empty description="暂未开放打赏功能" style="padding: 32px 0" />
        </template>
      </div>
    </n-modal>

    <!-- 修改密码对话框 -->
    <n-modal
      v-model:show="showPasswordModal"
      preset="card"
      title="修改密码"
      style="width: 500px; max-width: 90vw"
      :bordered="false"
      :segmented="false"
    >
      <n-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules">
        <n-form-item path="old_password" label="当前密码">
          <n-input
            v-model:value="passwordForm.old_password"
            type="password"
            show-password-on="click"
            placeholder="请输入当前密码"
          />
        </n-form-item>

        <n-form-item path="new_password" label="新密码">
          <n-input
            v-model:value="passwordForm.new_password"
            type="password"
            show-password-on="click"
            placeholder="至少6个字符"
          />
        </n-form-item>

        <n-form-item path="confirm_password" label="确认新密码">
          <n-input
            v-model:value="passwordForm.confirm_password"
            type="password"
            show-password-on="click"
            placeholder="再次输入新密码"
          />
        </n-form-item>
      </n-form>

      <template #footer>
        <n-space justify="end">
          <n-button @click="showPasswordModal = false">取消</n-button>
          <n-button type="primary" :loading="passwordSubmitting" @click="handlePasswordSubmit">
            确定
          </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import {
  ref,
  computed,
  h,
  onMounted,
  onBeforeUnmount,
  reactive,
  nextTick,
  provide,
  watch
} from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  MoonOutline,
  SunnyOutline,
  PersonOutline,
  LogOutOutline,
  SettingsOutline,
  SearchOutline,
  MenuOutline,
  HomeOutline,
  ArchiveOutline,
  ChatbubblesOutline,
  ChatboxEllipsesOutline,
  LinkOutline,
  InformationCircleOutline,
  CafeOutline
} from '@vicons/ionicons5'
import { useAuthStore, useAppStore } from '@/stores'
import { NIcon, useMessage, useDialog } from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui'
import { getPublicSettings } from '@/api/setting'
import type { SiteSettings } from '@/api/setting'
import { getPosts } from '@/api/post'
import { updatePassword } from '@/api/auth'
import type { PasswordForm } from '@/types/auth'
import { formatDate } from '@/utils/format'
import { highlightKeyword, extractHighlightSnippet } from '@/utils/highlight'
import type { Post } from '@/types/blog'
import CanvasBackground from '@/components/CanvasBackground.vue'
import GlobalBackground from '@/components/GlobalBackground.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const appStore = useAppStore()
const message = useMessage()
const dialog = useDialog()

const activeKey = ref(route.name as string)
const isDark = computed(() => appStore.theme === 'dark')
const siteSettings = ref<SiteSettings>({})
const defaultSiteName = '菱风叙'
const runningTime = ref('')
const copyrightYear = computed(() => {
  const startYear = siteStartDate.value.getFullYear()
  const currentYear = new Date().getFullYear()
  return currentYear <= startYear ? `${startYear}` : `${startYear}-${currentYear}`
})
const headerHidden = ref(false)
const isOnCover = ref(false)
const mainContentRef = ref<any>(null) // 取到组件实例后通过 $el 获取真实 DOM
const scrollEl = ref<HTMLElement | null>(null)
let lastScrollTop = 0
const searchKeyword = ref('')
const showSearchModal = ref(false)
const showMobileMenu = ref(false)
const showRewardModal = ref(false)
const searchResults = ref<Post[]>([])
const searchLoading = ref(false)
let searchTimer: number | null = null

// 修改密码相关
const showPasswordModal = ref(false)
const passwordFormRef = ref<FormInst | null>(null)
const passwordSubmitting = ref(false)
const passwordForm = reactive<PasswordForm>({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const passwordRules: FormRules = {
  old_password: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码至少6个字符', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    {
      validator: (_rule, value) => value === passwordForm.new_password,
      message: '两次密码不一致',
      trigger: ['blur', 'input']
    }
  ]
}

// 网站启动时间（优先从后台配置读取，未配置时使用默认值）
const defaultStartDate = '2025-10-13 00:00:00'
const siteStartDate = computed(() => {
  const configured = siteSettings.value.site_start_date
  return new Date(configured || defaultStartDate)
})

// 菜单选项（部分菜单根据登录状态动态生成）
const menuOptions = computed(() => {
  const base = [
    {
      label: '首页',
      key: 'Home',
      path: '/',
      icon: () => h(NIcon, null, { default: () => h(HomeOutline) })
    },
    {
      label: '文章归档',
      key: 'Archive',
      path: '/archive',
      icon: () => h(NIcon, null, { default: () => h(ArchiveOutline) })
    },
    {
      label: '友情链接',
      key: 'FriendLinks',
      path: '/friend-links',
      icon: () => h(NIcon, null, { default: () => h(LinkOutline) })
    },
    {
      label: '说说',
      key: 'Moments',
      path: '/moments',
      icon: () => h(NIcon, null, { default: () => h(ChatbubblesOutline) })
    },
    {
      label: '聊天室',
      key: 'Chat',
      path: '/chat',
      icon: () => h(NIcon, null, { default: () => h(ChatboxEllipsesOutline) })
    },
    {
      label: '关于我',
      key: 'About',
      path: '/about',
      icon: () => h(NIcon, null, { default: () => h(InformationCircleOutline) })
    }
  ] as any[]

  return base
})

// 用户菜单选项
const userMenuOptions = computed(() => {
  const options = [
    {
      label: '个人资料',
      key: 'profile',
      icon: () => h(NIcon, null, { default: () => h(PersonOutline) })
    },
    {
      label: '修改密码',
      key: 'change-password',
      icon: () => h(NIcon, null, { default: () => h(SettingsOutline) })
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

  // 查找菜单项，优先使用 path，如果没有 path 则使用 name
  let targetPath: string | undefined = undefined

  for (const item of menuOptions.value) {
    // 检查是否是当前项
    if (item.key === key) {
      targetPath = (item as any).path
      break
    }
    // 检查子项
    if ((item as any).children) {
      const child = (item as any).children.find((child: any) => child.key === key)
      if (child) {
        targetPath = child.path
        break
      }
    }
  }

  if (targetPath) {
    router.push(targetPath).catch(err => {
      // 忽略导航重复的错误
      if (err.name !== 'NavigationDuplicated') {
        console.error('路由跳转失败:', err)
      }
    })
  } else {
    router.push({ name: key }).catch(err => {
      // 忽略导航重复的错误
      if (err.name !== 'NavigationDuplicated') {
        console.error('路由跳转失败:', err)
      }
    })
  }
}

// 处理移动端菜单选择
function handleMobileMenuSelect(key: string) {
  activeKey.value = key
  showMobileMenu.value = false

  // 查找菜单项，优先使用 path，如果没有 path 则使用 name
  let targetPath: string | undefined = undefined

  for (const item of menuOptions.value) {
    // 检查是否是当前项
    if (item.key === key) {
      targetPath = (item as any).path
      break
    }
    // 检查子项
    if ((item as any).children) {
      const child = (item as any).children.find((child: any) => child.key === key)
      if (child) {
        targetPath = child.path
        break
      }
    }
  }

  if (targetPath) {
    router.push(targetPath).catch(err => {
      // 忽略导航重复的错误
      if (err.name !== 'NavigationDuplicated') {
        console.error('路由跳转失败:', err)
      }
    })
  } else {
    router.push({ name: key }).catch(err => {
      // 忽略导航重复的错误
      if (err.name !== 'NavigationDuplicated') {
        console.error('路由跳转失败:', err)
      }
    })
  }
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
      // 将封面背景图数组存入全局 store，供 GlobalBackground 使用
      let bgImages: string[] = []
      if (res.data.cover_bg_images) {
        try {
          bgImages = JSON.parse(res.data.cover_bg_images)
        } catch (e) {
          console.error('解析背景图数组失败:', e)
          bgImages = []
        }
      }
      appStore.setBgImages(bgImages)
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
  const diff = now.getTime() - siteStartDate.value.getTime()

  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  const seconds = Math.floor((diff % (1000 * 60)) / 1000)

  runningTime.value = `网站已运行 <span class="time-number">${days}</span> 天 <span class="time-number">${hours}</span> 小时 <span class="time-number">${minutes}</span> 分钟 <span class="time-number">${seconds}</span> 秒`
}

let timer: number | null = null

// 向子路由提供滚动容器引用和封面状态
provide('layoutScrollEl', scrollEl)
provide('isOnCover', isOnCover)

// 路由切换时重置封面状态
watch(
  () => route.name,
  name => {
    if (name === 'Home') {
      isOnCover.value = true
      headerHidden.value = false
      // 通过 Naive UI 组件方法滚回顶部
      mainContentRef.value?.scrollTo({ top: 0 })
      lastScrollTop = 0
    } else {
      isOnCover.value = false
    }
  }
)

onMounted(async () => {
  await nextTick()

  // 首页初始状态直接设为封面态（未滚动时导航栏透明）
  if (route.name === 'Home') {
    isOnCover.value = true
  }

  // 配置已在 main.ts 中预加载，这里只需获取即可
  fetchSiteSettings()
  calculateRunningTime()
  timer = window.setInterval(calculateRunningTime, 1000)
  updatePageTitle()
})

// 更新页面标题
function updatePageTitle() {
  const siteName = siteSettings.value.site_name || defaultSiteName
  const currentTitle = (route.meta.title as string) || '首页'
  document.title = `${currentTitle} | ${siteName}`
}

onBeforeUnmount(() => {
  if (timer) {
    clearInterval(timer)
  }
})

// Naive UI n-layout-content 的 scroll 回调，直接从事件拿 scrollTop
function handleNaiveScroll(e: Event) {
  const target = e.target as HTMLElement
  if (target) {
    scrollEl.value = target
  }
  handleScroll()
}

function handleScroll() {
  const current = scrollEl.value ? scrollEl.value.scrollTop : 0
  const delta = current - lastScrollTop

  // 判断是否在首页封面区域（用于导航栏透明效果）
  const isHomePage = route.name === 'Home'

  if (isHomePage) {
    isOnCover.value = current < window.innerHeight * 0.8
  } else {
    isOnCover.value = false
  }

  // 抖动过滤
  if (Math.abs(delta) < 5) {
    return
  }

  // 向下滚动超过 40px 就隐藏导航栏（无论是否在封面区域）
  if (delta > 0 && current > 40) {
    headerHidden.value = true
  } else if (delta < 0) {
    headerHidden.value = false
  }

  lastScrollTop = current
}

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
function goToPost(post: { id: number; slug: string }) {
  showSearchModal.value = false
  router.push(`/post/${post.slug}`)
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
    case 'change-password':
      showPasswordModal.value = true
      break
    case 'admin':
      router.push('/admin')
      break
    case 'logout':
      handleLogout()
      break
  }
}

// 处理退出登录
function handleLogout() {
  dialog.warning({
    title: '退出登录',
    content: '确定要退出登录吗？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: () => {
      authStore.logout()
      message.success('已退出登录')
      router.push('/')
    }
  })
}

// 提交修改密码
async function handlePasswordSubmit() {
  try {
    await passwordFormRef.value?.validate()
    passwordSubmitting.value = true

    await updatePassword(passwordForm)
    message.success('密码修改成功，请重新登录')

    // 重置表单
    passwordForm.old_password = ''
    passwordForm.new_password = ''
    passwordForm.confirm_password = ''
    showPasswordModal.value = false

    // 退出登录
    authStore.logout()
    router.push('/auth/login')
  } catch (error: any) {
    if (error.message) {
      message.error(error.message)
    }
  } finally {
    passwordSubmitting.value = false
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
  --header-height: 72px;
}

@media (max-width: 768px) {
  .default-layout {
    --header-height: 60px;
  }
}

.default-layout :deep(.n-layout) {
  height: 100vh;
}

/* 玻璃态顶部导航栏 — 深色半透明底 + 白色文字（适配全局背景图） */
.header {
  padding: 0 24px;
  height: 72px;
  display: flex;
  align-items: center;
  background: rgba(0, 0, 0, 0.35);
  backdrop-filter: blur(20px) saturate(180%);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.15);
  position: sticky;
  top: 0;
  z-index: 100;
  transition: all 0.3s;
}

.header.header-hidden {
  transform: translateY(-100%);
  opacity: 0;
}

html.dark .header {
  background: rgba(0, 0, 0, 0.5);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.4);
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
  gap: 12px;
  transition: transform 0.3s;
  flex-shrink: 0;
}

.logo:hover {
  transform: scale(1.05);
}

.logo-image {
  width: 32px;
  height: 32px;
  object-fit: contain;
  flex-shrink: 0;
  border-radius: 8px;
}

.logo h2 {
  margin: 0;
  font-size: 26px;
  font-weight: 800;
  color: #fff;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  letter-spacing: -0.02em;
}

html.dark .logo h2 {
  color: #fff;
}

.nav-menu {
  flex: 1 1 auto;
  margin: 0 24px;
  max-width: 100%;
  min-width: 0;
}

/* 自定义导航菜单样式 — 白色文字 */
.nav-menu :deep(.n-menu-item) {
  font-weight: 600;
  transition: all 0.3s;
}

.nav-menu :deep(.n-menu-item-content-header) {
  color: rgba(255, 255, 255, 0.9) !important;
  text-shadow: 0 1px 4px rgba(0, 0, 0, 0.3);
}

.nav-menu :deep(.n-menu-item-content__icon),
.nav-menu :deep(.n-icon) {
  color: rgba(255, 255, 255, 0.9) !important;
}

.nav-menu :deep(.n-menu-item:hover .n-menu-item-content-header) {
  color: #fff !important;
}

.nav-menu :deep(.n-menu-item-content::before) {
  background-color: rgba(255, 255, 255, 0.1) !important;
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
  color: rgba(255, 255, 255, 0.9) !important;
  text-shadow: 0 1px 4px rgba(0, 0, 0, 0.3);
}

.header-actions :deep(.n-icon) {
  color: rgba(255, 255, 255, 0.9) !important;
}

.main-content {
  padding: 24px 24px 0 24px;
  position: relative;
  z-index: 2;
  overflow-y: auto;
  height: calc(100vh - 72px);
  transition:
    transform 0.3s ease,
    height 0.3s ease;
  transform: translateY(0);
}

/* 首页去掉左右 padding，让封面可以铺满全宽 */
.main-content-home {
  padding-top: 0 !important;
  padding-left: 0 !important;
  padding-right: 0 !important;
  height: 100vh !important;
  margin-top: -72px;
}

.header + .main-content {
  transition:
    transform 0.3s ease,
    height 0.3s ease;
}

.header.header-hidden + .main-content {
  transform: translateY(-72px);
  height: 100vh;
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

/* 首页时去除顶部内边距，让封面紧贴视口顶部，且不限制宽度 */
.content-wrapper-home {
  padding-top: 0 !important;
  max-width: 100% !important;
  margin-left: 0 !important;
  margin-right: 0 !important;
}

/* 导航栏全透明态：仅在首页封面区域内去掉背景 */
.header.header-transparent {
  background: transparent !important;
  backdrop-filter: none !important;
  border-bottom-color: transparent !important;
  box-shadow: none !important;
}

.header.header-transparent :deep(.n-layout-header) {
  background: transparent !important;
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

.police-filing {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.police-icon {
  width: 16px;
  height: 16px;
  vertical-align: middle;
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
@media (max-width: 1150px) {
  .mobile-menu-btn {
    display: inline-flex;
  }

  .desktop-only {
    display: none !important;
  }

  .header {
    padding: 0 12px;
    height: 60px;
  }

  .logo-image {
    width: 24px;
    height: 24px;
  }

  .logo h2 {
    font-size: 20px;
  }

  .header-actions {
    gap: 8px;
  }

  .main-content {
    padding: 20px 12px 0 12px;
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

@media (max-width: 1024px) {
  .nav-menu {
    margin: 0 12px;
  }

  .nav-menu :deep(.n-menu-item) {
    font-size: 13px;
    padding-left: 8px;
    padding-right: 8px;
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

/* 打赏模态框样式 */
.reward-modal-content {
  text-align: center;
  padding: 8px 0;
}

.reward-qrcodes {
  display: flex;
  justify-content: center;
  gap: 32px;
  flex-wrap: wrap;
}

.reward-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.reward-qrcode {
  width: 160px;
  height: 160px;
  object-fit: contain;
  border-radius: 8px;
  border: 1px solid rgba(8, 145, 178, 0.15);
}

.reward-label {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #0891b2;
}

html.dark .reward-label {
  color: #38bdf8;
}

.reward-tip {
  margin: 20px 0 4px;
  font-size: 14px;
  color: #64748b;
}

html.dark .reward-tip {
  color: #94a3b8;
}

/* 小屏幕优化 */
@media (max-width: 480px) {
  .logo-image {
    width: 20px;
    height: 20px;
  }

  .logo h2 {
    font-size: 18px;
  }

  .header-actions :deep(.n-button .n-button__content > span) {
    display: none;
  }
}
</style>
