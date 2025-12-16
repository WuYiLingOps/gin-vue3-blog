<template>
  <n-card class="author-card" :bordered="false">
    <n-spin :show="loading">
      <div class="author-content">
      <!-- 头像 -->
      <div class="avatar-wrapper">
        <n-avatar
          :src="authorProfile?.author.avatar || ''"
          :size="80"
          round
          :fallback-src="defaultAvatar"
        >
          <template v-if="!authorProfile?.author.avatar">
            {{ (authorProfile?.author.nickname || authorProfile?.author.username || '博主').charAt(0).toUpperCase() }}
          </template>
        </n-avatar>
      </div>

      <!-- 用户名 -->
      <h3 class="author-name">{{ authorProfile?.author.nickname || authorProfile?.author.username }}</h3>

      <!-- 座右铭 -->
      <p v-if="authorProfile?.author.bio" class="author-bio">
        {{ authorProfile.author.bio }}
      </p>

      <!-- 统计数据 -->
      <div class="stats-section">
        <div class="stat-item" @click="goToPosts" title="查看所有文章">
          <div class="stat-label">文章</div>
          <div class="stat-value">{{ authorProfile?.stats.posts || 0 }}</div>
        </div>
        <div class="stat-item" @click="goToTags" title="查看所有标签">
          <div class="stat-label">标签</div>
          <div class="stat-value">{{ authorProfile?.stats.tags || 0 }}</div>
        </div>
        <div class="stat-item" @click="goToCategories" title="查看所有分类">
          <div class="stat-label">分类</div>
          <div class="stat-value">{{ authorProfile?.stats.categories || 0 }}</div>
        </div>
      </div>

      <!-- 社交链接（最多展示 5 个） -->
      <div v-if="visibleSocialLinks.length" class="social-links">
        <template v-for="link in visibleSocialLinks" :key="link.type">
          <a
            v-if="link.type !== 'wechat'"
            :href="link.href"
            target="_blank"
            rel="noopener noreferrer"
            class="social-icon"
            :class="link.type"
            :title="link.title"
          >
            <SocialIcons :type="link.type" />
          </a>
          <a
            v-else
            href="javascript:void(0)"
            class="social-icon wechat"
            :title="link.title"
            @click="showWechatQR = true"
          >
            <SocialIcons type="wechat" />
          </a>
        </template>
      </div>
      </div>
    </n-spin>

    <!-- 微信二维码弹窗 -->
    <n-modal v-model:show="showWechatQR">
      <n-card
        title="微信二维码"
        :bordered="false"
        style="max-width: 300px"
        closable
        @close="showWechatQR = false"
      >
        <div style="text-align: center">
          <n-image
            v-if="socialLinks.wechat"
            :src="socialLinks.wechat"
            width="200"
            height="200"
            object-fit="contain"
          />
          <p v-else style="color: #999">未配置微信二维码</p>
        </div>
      </n-card>
    </n-modal>
  </n-card>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getAuthorProfile, type AuthorProfile } from '@/api/blog'
import { getPublicSettings } from '@/api/setting'
import type { SiteSettings } from '@/api/setting'
import SocialIcons from './SocialIcons.vue'

const MAX_SOCIAL_LINKS = 5
type SocialLinkType = 'github' | 'gitee' | 'email' | 'rss' | 'qq' | 'wechat'
interface SocialLink {
  type: SocialLinkType
  href?: string
  title: string
}

const authorProfile = ref<AuthorProfile | null>(null)
const siteSettings = ref<SiteSettings>({})
const socialLinks = ref({
  github: '',
  gitee: '',
  email: '',
  rss: '',
  qq: '',
  wechat: ''
})
const showWechatQR = ref(false)
const loading = ref(false)
const defaultAvatar = '/default-avatar.png'
const router = useRouter()

// 计算需展示的社交链接（最多 5 个，按优先顺序）
const visibleSocialLinks = computed<SocialLink[]>(() => {
  const links: SocialLink[] = []
  const data = socialLinks.value

  if (data.github?.trim()) {
    links.push({ type: 'github', href: data.github.trim(), title: 'GitHub' })
  }
  if (data.gitee?.trim()) {
    links.push({ type: 'gitee', href: data.gitee.trim(), title: 'Gitee' })
  }
  if (data.email?.trim()) {
    links.push({ type: 'email', href: `mailto:${data.email.trim()}`, title: 'Email' })
  }
  if (data.rss?.trim()) {
    links.push({ type: 'rss', href: data.rss.trim(), title: 'RSS' })
  }
  if (data.qq?.trim()) {
    links.push({ type: 'qq', href: `tencent://message/?uin=${data.qq.trim()}`, title: 'QQ' })
  }
  if (data.wechat?.trim()) {
    links.push({ type: 'wechat', title: '微信' })
  }

  return links.slice(0, MAX_SOCIAL_LINKS)
})

// 跳转到文章列表
function goToPosts() {
  router.push('/')
}

// 跳转到标签列表
function goToTags() {
  router.push('/tag')
}

// 跳转到分类列表
function goToCategories() {
  router.push('/category')
}

// 获取博主信息
async function fetchAuthorProfile() {
  try {
    loading.value = true
    const res = await getAuthorProfile()
    if (res.data) {
      authorProfile.value = res.data
      console.log('博主信息:', res.data)
    }
  } catch (error: any) {
    console.error('获取博主信息失败:', error)
    // 如果获取失败，设置默认值避免显示错误
    authorProfile.value = {
      author: {
        id: 0,
        username: '博主',
        nickname: '博主',
        avatar: '',
        bio: ''
      },
      stats: {
        posts: 0,
        tags: 0,
        categories: 0
      }
    }
  } finally {
    loading.value = false
  }
}

// 获取网站设置（包含社交链接）
async function fetchSettings() {
  try {
    const res = await getPublicSettings()
    console.log('网站设置响应:', res)
    if (res.data) {
      siteSettings.value = res.data
      // 从设置中提取社交链接（如果存在）
      socialLinks.value = {
        github: (res.data.social_github || '').trim(),
        gitee: (res.data.social_gitee || '').trim(),
        email: (res.data.social_email || '').trim(),
        rss: (res.data.social_rss || '').trim(),
        qq: (res.data.social_qq || '').trim(),
        wechat: (res.data.social_wechat || '').trim()
      }
      console.log('社交链接:', socialLinks.value)
    }
  } catch (error: any) {
    console.error('获取网站设置失败:', error)
    console.error('错误详情:', error.response?.data || error.message)
  }
}

onMounted(() => {
  fetchAuthorProfile()
  fetchSettings()
})
</script>

<style scoped>
.author-card {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.5);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  border-radius: 16px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  position: sticky;
  top: 100px;
  z-index: 20; /* 确保个人名片在文章卡片之上 */
}

.author-card:hover {
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.12);
  border-color: rgba(8, 145, 178, 0.3);
}

html.dark .author-card {
  background: rgba(30, 41, 59, 0.7);
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

html.dark .author-card:hover {
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.5);
  border-color: rgba(56, 189, 248, 0.3);
}

.author-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px;
}

.avatar-wrapper {
  margin-bottom: 16px;
}

.author-name {
  margin: 0 0 8px 0;
  font-size: 20px;
  font-weight: 600;
  color: #1a202c;
  text-align: center;
}

html.dark .author-name {
  color: #e5e5e5;
}

.author-bio {
  margin: 0 0 20px 0;
  font-size: 14px;
  color: #64748b;
  text-align: center;
  line-height: 1.6;
  max-width: 100%;
  word-break: break-word;
}

html.dark .author-bio {
  color: #94a3b8;
}

.stats-section {
  display: flex;
  justify-content: space-around;
  width: 100%;
  margin-bottom: 20px;
  padding: 16px 0;
  border-top: 1px solid rgba(0, 0, 0, 0.06);
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
}

html.dark .stats-section {
  border-top-color: rgba(255, 255, 255, 0.1);
  border-bottom-color: rgba(255, 255, 255, 0.1);
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  transition: all 0.3s;
  padding: 8px;
  border-radius: 8px;
  flex: 1;
}

.stat-item:hover {
  background: rgba(8, 145, 178, 0.1);
  transform: translateY(-2px);
}

html.dark .stat-item:hover {
  background: rgba(56, 189, 248, 0.15);
}

.stat-label {
  font-size: 12px;
  color: #64748b;
}

html.dark .stat-label {
  color: #94a3b8;
}

.stat-value {
  font-size: 20px;
  font-weight: 700;
  color: #1a202c;
}

html.dark .stat-value {
  color: #e5e5e5;
}

.social-links {
  display: flex;
  justify-content: center;
  gap: 12px;
  flex-wrap: wrap;
}

.social-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  color: #fff;
  text-decoration: none;
  transition: all 0.3s;
  cursor: pointer;
  padding: 6px;
  box-sizing: border-box;
  overflow: hidden;
}

.social-icon svg {
  width: 24px;
  height: 24px;
  flex-shrink: 0;
}

.social-icon:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 禁用状态的样式 */
.social-icon.disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background: #d1d5db !important;
}

html.dark .social-icon.disabled {
  background: #4b5563 !important;
}

.social-icon.disabled:hover {
  transform: none;
  box-shadow: none;
}

.social-icon.github {
  background: #24292e;
}

.social-icon.email {
  background: #4285f4;
}

.social-icon.rss {
  background: #ffa500;
}

.social-icon.qq {
  background: #12b7f5;
}

.social-icon.wechat {
  background: #07c160;
}

.social-icon.gitee {
  background: transparent;
  /* Gitee 图标 SVG 内部已包含红色背景，无需额外设置 */
}

/* 移动端适配 */
@media (max-width: 768px) {
  .author-card {
    position: static;
    margin-bottom: 24px;
  }
}
</style>

