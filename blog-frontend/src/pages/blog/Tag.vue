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
      <div v-else class="tag-detail">
        <!-- 标签头部卡片 -->
        <div v-if="tag" class="tag-detail-header">
          <div class="tag-info-card">
            <div class="tag-badge" :style="{ backgroundColor: tag.color || '#2196F3' }">
              <span class="tag-badge-text">{{ tag.name }}</span>
            </div>
            <div class="tag-stats-row">
              <div class="stat-box">
                <span class="stat-value">{{ tag.post_count }}</span>
                <span class="stat-desc">篇文章</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 文章列表 -->
        <div class="posts-container">
          <div
            v-for="post in posts"
            :key="post.id"
            class="post-item"
            @click="router.push(`/post/${post.slug}`)"
          >
            <div class="post-item-header">
              <h2 class="post-item-title">{{ post.title }}</h2>
              <n-tag 
                :color="{ color: post.category.color || '#2196F3', textColor: '#fff' }" 
                size="small"
                class="post-category-tag"
              >
                {{ post.category.name }}
              </n-tag>
            </div>
            <p class="post-item-summary">{{ post.summary || '暂无摘要' }}</p>
            <div class="post-item-footer">
              <span class="post-date">{{ formatDate(post.created_at, 'YYYY年MM月DD日') }}</span>
              <div class="post-actions">
                <span class="post-action-item">
                  <i class="icon">👁️</i>
                  {{ post.view_count || 0 }}
                </span>
                <span class="post-action-item">
                  <i class="icon">❤️</i>
                  {{ post.like_count || 0 }}
                </span>
              </div>
            </div>
          </div>

          <n-empty 
            v-if="!loading && posts.length === 0" 
            description="该标签下暂无文章" 
            style="margin: 80px 0"
          />

          <div v-if="total > 0" class="pagination-wrapper">
            <n-pagination
              v-model:page="currentPage"
              :page-count="totalPages"
              :page-slot="7"
              @update:page="handlePageChange"
            />
          </div>
        </div>
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

/* 标签详情页样式 */
.tag-detail {
  max-width: 900px;
  margin: 0 auto;
}

.tag-detail-header {
  margin-bottom: 48px;
}

.tag-info-card {
  background: linear-gradient(135deg, #f8fafc 0%, #ffffff 100%);
  border-radius: 20px;
  padding: 40px;
  text-align: center;
  box-shadow: 
    0 10px 40px rgba(0, 0, 0, 0.08),
    0 2px 8px rgba(0, 0, 0, 0.04);
  position: relative;
  overflow: hidden;
  border: 1px solid rgba(0, 0, 0, 0.06);
}

html.dark .tag-info-card {
  background: linear-gradient(135deg, #1e293b 0%, #0f172a 100%);
  box-shadow: 
    0 10px 40px rgba(0, 0, 0, 0.4),
    0 2px 8px rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

/* 装饰性背景 */
.tag-info-card::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -20%;
  width: 300px;
  height: 300px;
  background: radial-gradient(circle, rgba(8, 145, 178, 0.08) 0%, transparent 70%);
  border-radius: 50%;
  pointer-events: none;
}

.tag-info-card::after {
  content: '';
  position: absolute;
  bottom: -30%;
  left: -15%;
  width: 250px;
  height: 250px;
  background: radial-gradient(circle, rgba(5, 150, 105, 0.08) 0%, transparent 70%);
  border-radius: 50%;
  pointer-events: none;
}

html.dark .tag-info-card::before {
  background: radial-gradient(circle, rgba(56, 189, 248, 0.15) 0%, transparent 70%);
}

html.dark .tag-info-card::after {
  background: radial-gradient(circle, rgba(74, 222, 128, 0.15) 0%, transparent 70%);
}

.tag-badge {
  display: inline-flex;
  align-items: center;
  padding: 14px 32px;
  border-radius: 50px;
  box-shadow: 
    0 8px 24px rgba(0, 0, 0, 0.15),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
  margin-bottom: 32px;
  position: relative;
  z-index: 1;
  animation: badgePulse 3s ease-in-out infinite;
}

@keyframes badgePulse {
  0%, 100% {
    transform: scale(1);
    box-shadow: 
      0 8px 24px rgba(0, 0, 0, 0.15),
      inset 0 1px 0 rgba(255, 255, 255, 0.2);
  }
  50% {
    transform: scale(1.05);
    box-shadow: 
      0 12px 32px rgba(0, 0, 0, 0.2),
      inset 0 1px 0 rgba(255, 255, 255, 0.3);
  }
}

.tag-badge-text {
  font-size: 28px;
  font-weight: 800;
  color: white;
  letter-spacing: 0.5px;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.tag-stats-row {
  display: flex;
  justify-content: center;
  gap: 24px;
  position: relative;
  z-index: 1;
}

.stat-box {
  background: rgba(255, 255, 255, 0.5);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 20px 40px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  border: 1px solid rgba(255, 255, 255, 0.6);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.06);
  transition: all 0.3s ease;
}

html.dark .stat-box {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.stat-box:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
}

html.dark .stat-box:hover {
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.4);
}

.stat-value {
  font-size: 36px;
  font-weight: 900;
  background: linear-gradient(135deg, #0891b2 0%, #059669 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  line-height: 1;
}

html.dark .stat-value {
  background: linear-gradient(135deg, #38bdf8 0%, #4ade80 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.stat-desc {
  font-size: 14px;
  color: #64748b;
  font-weight: 600;
  letter-spacing: 0.5px;
}

html.dark .stat-desc {
  color: #94a3b8;
}

/* 文章列表容器 */
.posts-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.post-item {
  background: white;
  border-radius: 12px;
  padding: 24px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 
    0 2px 8px rgba(0, 0, 0, 0.05),
    0 1px 4px rgba(0, 0, 0, 0.03);
  border: 1px solid rgba(0, 0, 0, 0.05);
}

html.dark .post-item {
  background: #1a1a1a;
  box-shadow: 
    0 2px 8px rgba(0, 0, 0, 0.3),
    0 1px 4px rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.post-item:hover {
  transform: translateY(-4px);
  box-shadow: 
    0 12px 32px rgba(0, 0, 0, 0.12),
    0 4px 12px rgba(0, 0, 0, 0.08);
  border-color: rgba(8, 145, 178, 0.3);
}

html.dark .post-item:hover {
  box-shadow: 
    0 12px 32px rgba(0, 0, 0, 0.5),
    0 4px 12px rgba(0, 0, 0, 0.3);
  border-color: rgba(56, 189, 248, 0.3);
}

.post-item-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 12px;
}

.post-item-title {
  font-size: 22px;
  font-weight: 700;
  margin: 0;
  color: #1a202c;
  line-height: 1.4;
  flex: 1;
  transition: color 0.2s;
}

html.dark .post-item-title {
  color: #e5e5e5;
}

.post-item:hover .post-item-title {
  color: #0891b2;
}

html.dark .post-item:hover .post-item-title {
  color: #38bdf8;
}

.post-category-tag {
  flex-shrink: 0;
  font-weight: 600;
}

.post-item-summary {
  color: #64748b;
  line-height: 1.7;
  margin: 0 0 16px 0;
  font-size: 15px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

html.dark .post-item-summary {
  color: #94a3b8;
}

.post-item-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 16px;
  border-top: 1px solid rgba(0, 0, 0, 0.06);
}

html.dark .post-item-footer {
  border-top-color: rgba(255, 255, 255, 0.06);
}

.post-date {
  font-size: 14px;
  color: #94a3b8;
  font-weight: 500;
}

html.dark .post-date {
  color: #64748b;
}

.post-actions {
  display: flex;
  gap: 20px;
}

.post-action-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: #94a3b8;
  font-weight: 500;
}

html.dark .post-action-item {
  color: #64748b;
}

.post-action-item .icon {
  font-style: normal;
  font-size: 16px;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 48px;
  padding: 24px 0;
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
  
  /* 标签详情页响应式 */
  .tag-info-card {
    padding: 32px 20px;
  }
  
  .tag-badge {
    padding: 12px 24px;
    margin-bottom: 24px;
  }
  
  .tag-badge-text {
    font-size: 22px;
  }
  
  .stat-box {
    padding: 16px 28px;
  }
  
  .stat-value {
    font-size: 28px;
  }
  
  .stat-desc {
    font-size: 13px;
  }
  
  .post-item {
    padding: 20px;
  }
  
  .post-item-title {
    font-size: 18px;
  }
  
  .post-item-summary {
    font-size: 14px;
  }
  
  .post-item-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .post-category-tag {
    align-self: flex-start;
  }
}
</style>

