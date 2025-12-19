<template>
  <div class="admin-layout">
    <n-loading-bar-provider>
      <loading-bar-handler />
    </n-loading-bar-provider>
    <n-layout has-sider>
      <!-- 侧边栏 -->
      <n-layout-sider
        collapse-mode="width"
        :collapsed-width="64"
        :width="240"
        :collapsed="collapsed"
        :show-trigger="!isMobile"
        :native-scrollbar="false"
        bordered
        @collapse="collapsed = true"
        @expand="collapsed = false"
      >
        <div class="logo">
          <h3 v-if="!collapsed">管理后台</h3>
          <h3 v-else>后台</h3>
        </div>
        <n-menu
          v-model:value="activeKey"
          :collapsed="collapsed"
          :collapsed-width="64"
          :collapsed-icon-size="22"
          :options="menuOptions"
          @update:value="handleMenuSelect"
        />
      </n-layout-sider>

      <n-layout>
        <!-- 头部 -->
        <n-layout-header class="admin-header">
          <div class="header-content">
            <!-- 移动端菜单按钮 -->
            <n-button
              v-if="isMobile"
              text
              class="mobile-menu-btn"
              @click="mobileMenuVisible = true"
            >
              <template #icon>
                <n-icon :component="MenuOutline" size="24" />
              </template>
            </n-button>

            <n-breadcrumb class="breadcrumb-wrapper">
              <n-breadcrumb-item v-if="!isMobile">管理后台</n-breadcrumb-item>
              <n-breadcrumb-item>{{ currentTitle }}</n-breadcrumb-item>
            </n-breadcrumb>

            <div class="header-actions">
              <n-button v-if="!isMobile" text @click="router.push('/')">
                <template #icon>
                  <n-icon :component="HomeOutline" />
                </template>
                返回首页
              </n-button>
              <n-button v-else circle text @click="router.push('/')">
                <template #icon>
                  <n-icon :component="HomeOutline" />
                </template>
              </n-button>

              <n-dropdown :options="userMenuOptions" @select="handleUserMenu">
                <n-button text>
                  <n-avatar round size="small" :src="authStore.user?.avatar" />
                  <span v-if="!isMobile" class="ml-2">{{ authStore.user?.nickname }}</span>
                </n-button>
              </n-dropdown>
            </div>
          </div>
        </n-layout-header>

        <!-- 内容区域 -->
        <n-layout-content class="admin-content">
          <div class="content-wrapper">
            <router-view />
          </div>
        </n-layout-content>
      </n-layout>
    </n-layout>

    <!-- 移动端抽屉菜单 -->
    <n-drawer v-model:show="mobileMenuVisible" :width="240" placement="left">
      <n-drawer-content title="管理后台" :native-scrollbar="false">
        <div class="mobile-menu">
          <n-menu
            v-model:value="activeKey"
            :options="menuOptions"
            @update:value="handleMobileMenuSelect"
          />
        </div>
      </n-drawer-content>
    </n-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h, defineComponent, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  HomeOutline,
  GridOutline,
  DocumentTextOutline,
  PricetagsOutline,
  FolderOutline,
  ChatbubblesOutline,
  ChatboxEllipsesOutline,
  PeopleOutline,
  PersonOutline,
  LogOutOutline,
  SettingsOutline,
  MenuOutline,
  ChatbubbleEllipsesOutline,
  ShieldCheckmarkOutline,
  LinkOutline
} from '@vicons/ionicons5'
import { useAuthStore } from '@/stores'
import { NIcon, useLoadingBar } from 'naive-ui'

// LoadingBar 处理组件
const LoadingBarHandler = defineComponent({
  name: 'LoadingBarHandler',
  setup() {
    const loadingBar = useLoadingBar()
    const router = useRouter()

    router.beforeEach(() => {
      loadingBar.start()
    })

    router.afterEach(() => {
      loadingBar.finish()
    })

    router.onError(() => {
      loadingBar.error()
    })

    return () => null
  }
})

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const collapsed = ref(false)
const activeKey = ref(route.name as string)
const mobileMenuVisible = ref(false)
const isMobile = ref(false)

const currentTitle = computed(() => {
  return route.meta.title || ''
})

// 检测是否为移动设备
function checkMobile() {
  isMobile.value = window.innerWidth <= 768
  if (isMobile.value) {
    collapsed.value = true
  }
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

// 渲染图标
const renderIcon = (icon: any) => {
  return () => h(NIcon, null, { default: () => h(icon) })
}

// 菜单选项
const menuOptions = [
  {
    label: '仪表盘',
    key: 'Dashboard',
    icon: renderIcon(GridOutline)
  },
  {
    label: '文章管理',
    key: 'PostManage',
    icon: renderIcon(DocumentTextOutline)
  },
  {
    label: '说说管理',
    key: 'MomentManage',
    icon: renderIcon(ChatboxEllipsesOutline)
  },
  {
    label: '聊天室管理',
    key: 'ChatManage',
    icon: renderIcon(ChatbubbleEllipsesOutline)
  },
  {
    label: '分类管理',
    key: 'CategoryManage',
    icon: renderIcon(FolderOutline)
  },
  {
    label: '标签管理',
    key: 'TagManage',
    icon: renderIcon(PricetagsOutline)
  },
  {
    label: '评论管理',
    key: 'CommentManage',
    icon: renderIcon(ChatbubblesOutline)
  },
  {
    label: '用户管理',
    key: 'UserManage',
    icon: renderIcon(PeopleOutline)
  },
  {
    label: 'IP访问控制',
    key: 'IPAccessControl',
    icon: renderIcon(ShieldCheckmarkOutline)
  },
  {
    label: '友链管理',
    key: 'FriendLinkManage',
    icon: renderIcon(LinkOutline)
  },
  {
    label: '友链分类',
    key: 'FriendLinkCategoryManage',
    icon: renderIcon(FolderOutline)
  },
  {
    label: '网站设置',
    key: 'SiteSettings',
    icon: renderIcon(SettingsOutline)
  }
]

// 用户菜单选项
const userMenuOptions = [
  {
    label: '个人资料',
    key: 'profile',
    icon: renderIcon(PersonOutline)
  },
  {
    label: '退出登录',
    key: 'logout',
    icon: renderIcon(LogOutOutline)
  }
]

// 处理菜单选择
function handleMenuSelect(key: string) {
  activeKey.value = key
  router.push({ name: key })
}

// 处理移动端菜单选择
function handleMobileMenuSelect(key: string) {
  activeKey.value = key
  mobileMenuVisible.value = false
  router.push({ name: key })
}

// 处理用户菜单
function handleUserMenu(key: string) {
  switch (key) {
    case 'profile':
      router.push('/profile')
      break
    case 'logout':
      authStore.logout()
      router.push('/')
      break
  }
}
</script>

<style scoped>
.admin-layout {
  height: 100vh;
  position: relative;
}

/* 侧边栏样式 */
.admin-layout :deep(.n-layout-sider) {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(20px) saturate(180%);
  border-right: 1px solid rgba(8, 145, 178, 0.1);
  box-shadow: 4px 0 24px rgba(0, 0, 0, 0.06);
}

html.dark .admin-layout :deep(.n-layout-sider) {
  background: rgba(15, 23, 42, 0.9);
  border-right: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 4px 0 24px rgba(0, 0, 0, 0.3);
}

/* 移动端隐藏侧边栏 */
@media (max-width: 768px) {
  .admin-layout :deep(.n-layout-sider) {
    display: none;
  }
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-bottom: 1px solid rgba(8, 145, 178, 0.1);
}

html.dark .logo {
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.logo h3 {
  margin: 0;
  font-size: 20px;
  font-weight: 800;
  background: linear-gradient(135deg, #0891b2 0%, #059669 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

html.dark .logo h3 {
  background: linear-gradient(135deg, #38bdf8 0%, #4ade80 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.admin-header {
  height: 64px;
  display: flex;
  align-items: center;
  padding: 0 24px;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(20px) saturate(180%);
  border-bottom: 1px solid rgba(8, 145, 178, 0.1);
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.06);
}

html.dark .admin-header {
  background: rgba(15, 23, 42, 0.8);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.3);
}

.header-content {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.mobile-menu-btn {
  margin-right: 8px;
}

.breadcrumb-wrapper {
  flex: 1;
  min-width: 0;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

/* 移动端头部样式 */
@media (max-width: 768px) {
  .admin-header {
    padding: 0 12px;
    height: 56px;
  }

  .header-actions {
    gap: 8px;
  }
  
  .breadcrumb-wrapper {
    font-size: 14px;
  }
}

.admin-content {
  padding: 24px;
  overflow-y: auto;
  height: calc(100vh - 64px);
}

.content-wrapper {
  background: transparent;
  border-radius: 16px;
  padding: 0;
  min-height: calc(100vh - 112px);
}

.ml-2 {
  margin-left: 8px;
}

/* 移动端内容区域样式 */
@media (max-width: 768px) {
  .admin-content {
    padding: 12px;
    height: calc(100vh - 56px);
  }

  .content-wrapper {
    min-height: calc(100vh - 80px);
  }
}

/* 移动端抽屉菜单样式 */
.mobile-menu {
  height: 100%;
}
</style>

