<template>
  <div class="tag-page">
    <n-spin :show="loading">
      <!-- 标签列表视图 -->
      <div v-if="!tagId" class="tags-list">
        <!-- 外层白色卡片 -->
        <div class="outer-card">
          <!-- 标题和统计信息 -->
          <div class="page-header">
            <n-icon :component="PricetagsOutline" size="48" class="header-icon" />
            <div class="stats-text">
              目前共计 <span class="stats-number">{{ tags.length }}</span> 个标签
            </div>
          </div>

          <!-- 内层浅灰色卡片 - 标签云 -->
          <div class="inner-card">
            <div class="tags-cloud">
              <div
                v-for="(t, index) in tags"
                :key="t.id"
                class="tag-bubble"
                :class="t.font_size ? '' : getTagSizeClass(t.post_count)"
                :style="{ ...getTagStyle(t), ...getTagPosition(index) }"
                @click="router.push(`/tag/${t.id}`)"
              >
                <span class="tag-hash">#</span>
                <span class="tag-name">{{ t.name }}</span>
              </div>
            </div>
            
            <n-empty v-if="!loading && tags.length === 0" description="暂无标签" style="margin-top: 48px" />
          </div>
        </div>
      </div>

      <!-- 标签详情视图 -->
      <div v-else>
        <n-card v-if="tag">
          <div class="tag-header">
            <n-tag :color="{ color: tag.color || '#2196F3', textColor: '#fff' }" size="large">
              # {{ tag.name }}
            </n-tag>
            <p style="margin-top: 16px; color: #666">{{ tag.post_count }} 篇文章</p>
          </div>
        </n-card>

        <!-- 文章列表 -->
        <n-space vertical :size="16" style="margin-top: 24px">
          <n-card
            v-for="post in posts"
            :key="post.id"
            hoverable
            class="post-card"
            @click="router.push(`/post/${post.id}`)"
          >
            <h2 class="post-title">{{ post.title }}</h2>
            <p class="post-summary">{{ post.summary }}</p>
            <div class="post-meta">
              <n-space>
                <n-tag :bordered="false" type="info">{{ post.category.name }}</n-tag>
                <n-divider vertical />
                <span>{{ formatDate(post.created_at, 'YYYY-MM-DD') }}</span>
              </n-space>
            </div>
          </n-card>

          <n-empty v-if="!loading && posts.length === 0" description="该标签下暂无文章" />

          <n-pagination
            v-if="total > 0"
            v-model:page="currentPage"
            :page-count="totalPages"
            @update:page="handlePageChange"
          />
        </n-space>
      </div>
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useMessage } from 'naive-ui'
import { PricetagsOutline } from '@vicons/ionicons5'
import { getTagById, getPostsByTag, getTags } from '@/api/tag'
import { formatDate } from '@/utils/format'
import type { Tag, Post } from '@/types/blog'

const router = useRouter()
const route = useRoute()
const message = useMessage()

const loading = ref(false)
const tags = ref<Tag[]>([])
const tag = ref<Tag | null>(null)
const posts = ref<Post[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = 10

const tagId = computed(() => route.params.id ? Number(route.params.id) : null)
const totalPages = computed(() => Math.ceil(total.value / pageSize))

// 根据文章数量返回标签大小类名（如果标签有自定义大小则不使用）
function getTagSizeClass(postCount: number): string {
  if (postCount >= 10) return 'tag-large'
  if (postCount >= 5) return 'tag-medium'
  return 'tag-small'
}

// 获取标签样式（使用后台设置的颜色、文字颜色和文字大小）
function getTagStyle(tag: Tag) {
  const baseColor = tag.color || '#2196F3' // 使用后台设置的背景颜色，如果没有则使用默认蓝色
  const textColor = tag.text_color || '#fff' // 使用后台设置的文字颜色，如果没有则使用白色
  const fontSize = tag.font_size // 使用后台设置的文字大小
  
  const style: any = {
    backgroundColor: baseColor, // 使用纯色背景，不用渐变
    '--tag-text-color': textColor // 自定义CSS变量
  }
  
  // 如果有自定义文字大小，则覆盖默认大小
  if (fontSize) {
    style.fontSize = `${fontSize}px`
  }
  
  return style
}

// 获取标签位置（完全打乱的布局）
function getTagPosition(index: number) {
  // 使用随机种子算法确保每次刷新位置一致
  const seed2 = (index * 4253 + 17389) % 233280
  const randomY = seed2 / 233280
  
  // 创建随机的垂直偏移，让每行高度不同（减小幅度避免重叠）
  const rowOffset = Math.sin(index * 1.5) * 20 // 减小波浪幅度
  const randomOffset = (randomY - 0.5) * 30 // 减小随机偏移范围
  const totalOffsetY = rowOffset + randomOffset
  
  return {
    // 使用 margin 来定位，不使用 transform，这样不会和动画冲突
    marginTop: `${totalOffsetY}px`,
    marginLeft: '18px', // 增加左边距
    marginRight: '18px', // 增加右边距
    marginBottom: '24px' // 增加底部边距，防止垂直重叠
  }
}

// 监听路由变化
watch(() => route.params.id, (newId) => {
  // 重置状态
  tag.value = null
  posts.value = []
  total.value = 0
  currentPage.value = 1
  
  if (newId) {
    fetchTag()
    fetchPosts()
  } else {
    fetchTags()
  }
}, { immediate: true })

onMounted(() => {
  if (tagId.value) {
    fetchTag()
    fetchPosts()
  } else {
    fetchTags()
  }
})

async function fetchTags() {
  try {
    loading.value = true
    const res = await getTags()
    if (res.data) {
      tags.value = res.data
    }
  } catch (error: any) {
    message.error('获取标签列表失败')
  } finally {
    loading.value = false
  }
}

async function fetchTag() {
  try {
    const res = await getTagById(tagId.value!)
    if (res.data) {
      tag.value = res.data
    }
  } catch (error: any) {
    message.error('获取标签信息失败')
  }
}

async function fetchPosts() {
  try {
    loading.value = true
    const res = await getPostsByTag(tagId.value!, {
      page: currentPage.value,
      page_size: pageSize
    })

    if (res.data) {
      posts.value = res.data.list
      total.value = res.data.total
    }
  } catch (error: any) {
    message.error('获取文章列表失败')
  } finally {
    loading.value = false
  }
}

function handlePageChange() {
  fetchPosts()
}
</script>

<style scoped>
.tag-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
  position: relative;
  min-height: 100vh;
}

/* 背景装饰 */
.tag-page::before {
  content: '';
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: 
    radial-gradient(circle at 20% 30%, rgba(8, 145, 178, 0.08) 0%, transparent 50%),
    radial-gradient(circle at 80% 70%, rgba(5, 150, 105, 0.08) 0%, transparent 50%),
    radial-gradient(circle at 50% 50%, rgba(56, 189, 248, 0.05) 0%, transparent 70%);
  pointer-events: none;
  z-index: 0;
}

/* 标签列表样式 */
.tags-list {
  padding: 24px 0;
  position: relative;
  z-index: 1;
}

/* 外层白色卡片 */
.outer-card {
  background: white;
  border-radius: 16px;
  box-shadow: 
    0 8px 32px rgba(0, 0, 0, 0.08),
    0 2px 8px rgba(0, 0, 0, 0.04);
  padding: 48px 40px;
  position: relative;
  overflow: hidden;
}

html.dark .outer-card {
  background: #1a1a1a;
  box-shadow: 
    0 8px 32px rgba(0, 0, 0, 0.4),
    0 2px 8px rgba(0, 0, 0, 0.2);
}

/* 内层浅灰色卡片 */
.inner-card {
  background: #f9f9f9;
  border-radius: 12px;
  padding: 60px 40px 80px;
  margin-top: 32px;
  box-shadow: 
    inset 0 2px 8px rgba(0, 0, 0, 0.05),
    0 2px 4px rgba(0, 0, 0, 0.02);
}

html.dark .inner-card {
  background: #2a2a2a;
  box-shadow: 
    inset 0 2px 8px rgba(0, 0, 0, 0.3),
    0 2px 4px rgba(0, 0, 0, 0.1);
}

/* 页面头部 */
.page-header {
  text-align: center;
  margin-bottom: 0;
  position: relative;
  z-index: 1;
}

.header-icon {
  color: #0891b2;
  margin-bottom: 20px;
  animation: float 3s ease-in-out infinite;
  filter: drop-shadow(0 4px 12px rgba(8, 145, 178, 0.3));
  position: relative;
  z-index: 1;
}

html.dark .header-icon {
  color: #38bdf8;
  filter: drop-shadow(0 4px 12px rgba(56, 189, 248, 0.4));
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
  }
  50% {
    transform: translateY(-12px) rotate(5deg);
  }
}

.page-title {
  font-size: 48px;
  font-weight: 900;
  text-align: center;
  margin: 0 0 20px 0;
  background: linear-gradient(135deg, #0891b2 0%, #059669 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: -0.03em;
  position: relative;
  z-index: 1;
  text-shadow: 0 2px 20px rgba(8, 145, 178, 0.2);
}

html.dark .page-title {
  background: linear-gradient(135deg, #38bdf8 0%, #4ade80 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-shadow: 0 2px 20px rgba(56, 189, 248, 0.3);
}

/* 统计信息 */
.stats-text {
  font-size: 20px;
  font-weight: 600;
  color: #475569;
  margin-top: 12px;
  position: relative;
  z-index: 1;
  letter-spacing: 0.02em;
}

html.dark .stats-text {
  color: #cbd5e1;
}

.stats-number {
  font-size: 32px;
  font-weight: 800;
  background: linear-gradient(135deg, #0891b2 0%, #059669 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0 10px;
  display: inline-block;
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
}

html.dark .stats-number {
  background: linear-gradient(135deg, #38bdf8 0%, #4ade80 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* 标签云 */
.tags-cloud {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  align-items: baseline;
  align-content: flex-start;
  padding: 40px 60px 60px;
  min-height: 400px;
  position: relative;
  gap: 0;
  line-height: 2.5;
  z-index: 1;
}

/* 气泡标签 */
.tag-bubble {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 12px 24px;
  border-radius: 50px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  position: relative;
  overflow: hidden;
  animation: fadeInUp 0.4s ease-out, floatTag 6s ease-in-out infinite;
  will-change: transform;
}

/* 每个标签不同的动画延迟，让它们不同步 */
.tag-bubble:nth-child(1) { animation-delay: 0s, 0s; }
.tag-bubble:nth-child(2) { animation-delay: 0s, -1s; }
.tag-bubble:nth-child(3) { animation-delay: 0s, -2s; }
.tag-bubble:nth-child(4) { animation-delay: 0s, -3s; }
.tag-bubble:nth-child(5) { animation-delay: 0s, -4s; }
.tag-bubble:nth-child(6) { animation-delay: 0s, -5s; }
.tag-bubble:nth-child(7) { animation-delay: 0s, -0.5s; }
.tag-bubble:nth-child(8) { animation-delay: 0s, -1.5s; }
.tag-bubble:nth-child(9) { animation-delay: 0s, -2.5s; }
.tag-bubble:nth-child(10) { animation-delay: 0s, -3.5s; }
.tag-bubble:nth-child(11) { animation-delay: 0s, -4.5s; }
.tag-bubble:nth-child(12) { animation-delay: 0s, -5.5s; }
.tag-bubble:nth-child(13) { animation-delay: 0s, -1.2s; }
.tag-bubble:nth-child(14) { animation-delay: 0s, -2.4s; }
.tag-bubble:nth-child(15) { animation-delay: 0s, -3.6s; }
.tag-bubble:nth-child(16) { animation-delay: 0s, -4.8s; }
.tag-bubble:nth-child(17) { animation-delay: 0s, -0.8s; }
.tag-bubble:nth-child(18) { animation-delay: 0s, -1.8s; }
.tag-bubble:nth-child(19) { animation-delay: 0s, -2.8s; }
.tag-bubble:nth-child(20) { animation-delay: 0s, -3.8s; }
.tag-bubble:nth-child(n+21) { animation-delay: 0s, -2.2s; }

/* 入场动画 - 更快更流畅 */
@keyframes fadeInUp {
  0% {
    opacity: 0;
  }
  100% {
    opacity: 1;
  }
}

/* 浮动动画 - 像跷跷板一样摇摆 */
@keyframes floatTag {
  0% {
    transform: translate(0, 0) rotate(0deg);
  }
  25% {
    transform: translate(4px, -6px) rotate(2deg);
  }
  50% {
    transform: translate(0, -12px) rotate(0deg);
  }
  75% {
    transform: translate(-4px, -6px) rotate(-2deg);
  }
  100% {
    transform: translate(0, 0) rotate(0deg);
  }
}

.tag-bubble::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(255, 255, 255, 0.2);
  opacity: 0;
  transition: opacity 0.3s;
}

.tag-bubble:hover::before {
  opacity: 1;
}

.tag-bubble:hover {
  transform: translate(0, -8px) scale(1.08);
  box-shadow: 0 12px 28px rgba(0, 0, 0, 0.25);
  animation-play-state: paused;
}

.tag-bubble:active {
  transform: translate(0, -4px) scale(1.05);
}

/* 标签大小 */
.tag-small {
  font-size: 16px;
  padding: 12px 22px;
}

.tag-medium {
  font-size: 18px;
  padding: 14px 26px;
}

.tag-large {
  font-size: 20px;
  padding: 16px 30px;
}

.tag-hash {
  font-weight: 700;
  color: #726c69 !important;
  opacity: 1;
  margin-right: 2px;
}

.tag-name {
  letter-spacing: 0.02em;
  color: var(--tag-text-color, white);
  font-weight: 600;
}

/* 标签详情样式 */
.tag-header {
  text-align: center;
  padding: 24px 0;
}

.post-card {
  cursor: pointer;
  transition: all 0.3s;
}

.post-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

html.dark .post-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.5);
}

.post-title {
  font-size: 20px;
  margin: 0 0 12px 0;
  color: #1a202c;
}

html.dark .post-title {
  color: #e5e5e5;
}

.post-summary {
  color: #64748b;
  line-height: 1.6;
  margin: 0 0 12px 0;
}

html.dark .post-summary {
  color: #94a3b8;
}

.post-meta {
  font-size: 14px;
}

/* 响应式 */
@media (max-width: 768px) {
  .tag-page {
    padding: 16px;
  }
  
  .outer-card {
    padding: 32px 24px;
    border-radius: 12px;
  }
  
  .inner-card {
    padding: 40px 24px 60px;
    border-radius: 8px;
    margin-top: 24px;
  }
  
  .page-header {
    margin-bottom: 0;
  }
  
  .header-icon {
    font-size: 36px !important;
  }
  
  .page-title {
    font-size: 32px;
    margin-bottom: 12px;
  }
  
  .stats-text {
    font-size: 16px;
  }
  
  .stats-number {
    font-size: 24px;
  }
  
  .tags-cloud {
    padding: 32px 16px;
    gap: 12px;
    min-height: 200px;
  }
  
  .tag-bubble {
    font-size: 13px !important;
    padding: 8px 16px !important;
  }
  
  .tag-small {
    font-size: 12px !important;
    padding: 8px 14px !important;
  }
  
  .tag-medium {
    font-size: 13px !important;
    padding: 9px 18px !important;
  }
  
  .tag-large {
    font-size: 15px !important;
    padding: 10px 20px !important;
  }
}
</style>
