<template>
  <div class="home-page">
    <!-- 顶部：贡献热力图（左） + 个人名片（右） -->
    <div class="top-row">
      <div class="calendar-wrapper">
        <GiteeCalendar />
      </div>
      <div class="top-author-wrapper">
        <AuthorCard />
        <div class="tablet-calendar">
          <GiteeCalendar />
        </div>
        <div class="mobile-announcement">
          <AnnouncementBoard :limit="3" />
        </div>
        <div class="mobile-tag-cloud">
          <TagCloudWidget />
        </div>
      </div>
    </div>

    <div class="home-layout">
      <!-- 左侧：文章列表 -->
      <div class="posts-section">
        <n-space vertical :size="24">
          <!-- 文章列表 -->
          <n-spin :show="loading">
            <n-space vertical :size="16">
              <n-card
                v-for="post in posts"
                :key="post.id"
                hoverable
                class="post-card"
                @click="router.push(`/post/${post.id}`)"
              >
            <div class="post-card-wrapper">
              <div class="post-card-content">
                <!-- 文章信息 -->
                <div class="post-info">
                  <h2 class="post-title">
                    <n-tag v-if="post.is_top" type="error" size="small" class="top-tag">置顶</n-tag>
                    <span v-html="getHighlightedTitle(post.title)"></span>
                  </h2>

                  <p class="post-summary" v-html="getHighlightedSummary(post)"></p>
                </div>

                <!-- 封面图 -->
                <div v-if="post.cover" class="post-cover">
                  <n-image
                    :src="post.cover"
                    :alt="post.title"
                    object-fit="cover"
                    :preview-disabled="true"
                    class="cover-image"
                  />
                </div>
              </div>

              <!-- 底部信息栏 -->
              <div class="post-footer">
                <div class="post-meta">
                  <n-space :size="8">
                    <n-avatar :src="post.user.avatar" :size="24" round />
                    <span class="meta-item">{{ post.user.nickname }}</span>
                    <n-divider vertical />
                    <span class="meta-item">
                      <n-icon :component="TimeOutline" size="14" />
                      {{ formatDate(post.created_at, 'YYYY-MM-DD') }}
                    </span>
                    <n-divider vertical />
                    <span class="meta-item">
                      <n-icon :component="EyeOutline" size="14" />
                      {{ post.view_count }}
                    </span>
                  </n-space>
                </div>

                <div class="post-tags">
                  <n-tag :bordered="false" type="info" size="small">{{ post.category.name }}</n-tag>
                  <n-tag v-for="tag in post.tags" :key="tag.id" :bordered="false" size="small">
                    {{ tag.name }}
                  </n-tag>
                </div>
              </div>
            </div>
          </n-card>

          <!-- 空状态 -->
          <n-empty v-if="!loading && posts.length === 0" description="暂无文章" />

          <!-- 分页 -->
          <div v-if="total > 0" class="pagination-wrapper">
            <n-pagination
              v-model:page="currentPage"
              :page-count="totalPages"
              :page-size="pageSize"
              :page-slot="7"
              @update:page="handlePageChange"
            />
          </div>
        </n-space>
      </n-spin>
    </n-space>
      </div>

    <!-- 右侧：公告栏 + 热门文章（保持原有组件和顺序） -->
    <div class="sidebar-section">
      <div class="sidebar-card-wrapper sidebar-announcement">
        <AnnouncementBoard :limit="3" />
      </div>
      <div class="sidebar-card-wrapper sidebar-hot-posts">
        <HotPostsCard />
      </div>
      <div class="sidebar-card-wrapper sidebar-tag-cloud">
        <TagCloudWidget />
      </div>
    </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useMessage } from 'naive-ui'
import { TimeOutline, EyeOutline } from '@vicons/ionicons5'
import { getPosts } from '@/api'
import { formatDate } from '@/utils/format'
import { highlightKeyword, extractHighlightSnippet } from '@/utils/highlight'
import { useBlogStore } from '@/stores'
import type { Post } from '@/types/blog'
import AuthorCard from '@/components/AuthorCard.vue'
import AnnouncementBoard from '@/components/AnnouncementBoard.vue'
import HotPostsCard from '@/components/HotPostsCard.vue'
import TagCloudWidget from '@/components/TagCloudWidget.vue'
import GiteeCalendar from '@/components/GiteeCalendar.vue'

const router = useRouter()
const route = useRoute()
const message = useMessage()
const blogStore = useBlogStore()

const loading = ref(false)
const posts = ref<Post[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(6)
const searchKeyword = ref('')

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

// 监听路由变化，处理搜索
watch(() => route.query.keyword, (newKeyword) => {
  if (newKeyword) {
    searchKeyword.value = newKeyword as string
    currentPage.value = 1
    fetchPosts()
  }
}, { immediate: true })

onMounted(() => {
  blogStore.init()
  // 从URL获取搜索关键词
  if (route.query.keyword) {
    searchKeyword.value = route.query.keyword as string
  }
  fetchPosts()
})

async function fetchPosts() {
  try {
    loading.value = true
    const res = await getPosts({
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: searchKeyword.value,
      status: 1
    })

    if (res.data) {
      posts.value = res.data.list
      total.value = res.data.total
    }
  } catch (error: any) {
    message.error(error.message || '获取文章列表失败')
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  currentPage.value = page
  fetchPosts()
}

// 高亮标题
function getHighlightedTitle(title: string): string {
  if (!searchKeyword.value) {
    return title
  }
  return highlightKeyword(title, searchKeyword.value)
}

// 高亮摘要
function getHighlightedSummary(post: Post): string {
  const summary = post.summary || ''
  
  if (!searchKeyword.value) {
    return summary
  }
  
  // 检查摘要中是否包含关键词
  const lowerSummary = summary.toLowerCase()
  const lowerKeyword = searchKeyword.value.toLowerCase()
  
  if (lowerSummary.includes(lowerKeyword)) {
    // 如果摘要中包含关键词，直接高亮
    return highlightKeyword(summary, searchKeyword.value)
  } else if (post.content) {
    // 如果摘要中不包含关键词，但内容存在，从内容中提取包含关键词的片段
    const snippet = extractHighlightSnippet(post.content, searchKeyword.value, 150)
    // 如果提取到了包含关键词的片段，就使用它；否则使用原摘要
    if (snippet && snippet.toLowerCase().includes(lowerKeyword)) {
      return highlightKeyword(snippet, searchKeyword.value)
    }
  }
  
  // 默认返回原摘要（可能不包含关键词，但至少显示文章的原始摘要）
  return summary
}
</script>

<style scoped>
.home-page {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 20px;
  position: relative;
  z-index: 1;
}

.top-row {
  display: grid;
  /* 左侧自适应内容，右侧固定为与侧边栏相同的宽度（360px），保持个人名片大小一致 */
  grid-template-columns: minmax(0, 1fr) 360px;
  gap: 32px;
  margin-bottom: 24px;
}

.calendar-wrapper,
.top-author-wrapper {
  display: flex;
  min-width: 0; /* 防止 flex 子元素被压缩 */
  overflow: visible; /* 确保内容不被裁剪 */
}

.calendar-wrapper :deep(.hexo-calendar-card) {
  width: 100%;
}

/* 首页顶部区域的个人名片不需要吸顶，避免在特定分辨率下出现悬停不随滚动的问题 */
.top-author-wrapper :deep(.author-card) {
  width: 100%;
  position: static;
  top: auto;
}

.tablet-calendar {
  display: none;
  width: 100%;
  margin-top: 16px;
}

.mobile-announcement,
.mobile-tag-cloud {
  display: none;
  width: 100%;
  margin-top: 12px;
}

/* 桌面端：让顶部个人名片与侧边栏卡片左边对齐，保持视觉宽度一致 */
@media (min-width: 1025px) {
  .top-author-wrapper {
    margin-left: 12px;
  }
}

.home-layout {
  display: grid;
  grid-template-columns: 1fr 360px; /* 扩大右侧区域，容纳更大的个人名片 */
  gap: 32px; /* 增加间距，避免重叠 */
  align-items: start;
}

.posts-section {
  min-width: 0; /* 防止内容溢出 */
}

.sidebar-section {
  position: relative;
  z-index: 10; /* 确保侧边栏在文章卡片之上 */
  margin-left: 12px; /* 稍微往右移动，留出呼吸空间 */
  display: flex;
  flex-direction: column;
  gap: 16px; /* 卡片之间的间距，避免重叠 */
}

/* 平板端和移动端布局 */
@media (max-width: 1024px) {
  .home-layout {
    grid-template-columns: 1fr;
  }

  /* 桌面端的热力图在平板和移动端都不显示 */
  .calendar-wrapper {
    display: none;
  }

  .sidebar-announcement,
  .sidebar-hot-posts,
  .sidebar-tag-cloud {
    display: none;
  }

  .top-row {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .top-author-wrapper {
    margin-left: 0;
    flex-direction: column;
    min-width: 0; /* 防止内容被压缩 */
    overflow: visible; /* 确保内容不被裁剪 */
  }

  .sidebar-section {
    order: 0;
    margin-left: 0;
  }
}

/* 平板端布局（768px ~ 1024px） */
@media (min-width: 768px) and (max-width: 1024px) {
  .tablet-calendar {
    display: block;
    /* 移除固定缩放，让组件内部的自动缩放逻辑处理自适应 */
  }

  .mobile-announcement,
  .mobile-tag-cloud {
    display: block;
  }
}

/* 移动端布局（< 768px） */
@media (max-width: 767px) {
  /* 移动端也显示热力图，位置在个人名片下方、公告栏上方 */
  .tablet-calendar {
    display: block;
    width: 100%;
    margin-top: 16px;
    margin-bottom: 0;
    min-width: 0; /* 防止内容被压缩 */
    overflow: visible; /* 确保内容不被裁剪 */
  }

  /* 移动端热力图卡片优化 */
  .tablet-calendar :deep(.hexo-calendar-card) {
    width: 100%;
    min-width: 0;
  }

  .top-author-wrapper {
    min-width: 0; /* 防止内容被压缩 */
    overflow: visible; /* 确保内容不被裁剪 */
  }

  .mobile-announcement,
  .mobile-tag-cloud {
    display: block;
    min-width: 0; /* 防止内容被压缩 */
    overflow: visible; /* 确保内容不被裁剪 */
  }
}

/* 玻璃态卡片效果 */
.home-page :deep(.n-card) {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.5);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  border-radius: 16px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.home-page :deep(.n-card):hover {
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.12);
  border-color: rgba(8, 145, 178, 0.3);
}

.home-page :deep(.n-card .n-card__content) {
  padding: 20px !important;
}

.post-card {
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
  position: relative;
}

.post-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(8, 145, 178, 0.05), transparent);
  transition: left 0.6s;
  z-index: 1;
}

.post-card:hover::before {
  left: 100%;
}

.post-card:hover {
  transform: translateY(-8px) scale(1.02); /* 减小放大倍数，避免向右扩展太多 */
  box-shadow: 0 20px 48px rgba(0, 0, 0, 0.15);
  border-color: rgba(8, 145, 178, 0.4);
  z-index: 5; /* 设置z-index，但低于侧边栏 */
}

.post-card-wrapper {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-height: 160px;
}

.post-card-content {
  display: flex;
  gap: 24px;
  position: relative;
  z-index: 2;
  flex: 1;
}

.post-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.post-title {
  font-size: 20px;
  font-weight: 700;
  margin: 0 0 12px 0;
  color: #1a202c;
  letter-spacing: -0.01em;
  line-height: 1.4;
  display: flex;
  align-items: center;
  gap: 8px;
}

.top-tag {
  flex-shrink: 0;
}

html.dark .post-title {
  color: #e5e5e5;
}

.post-summary {
  color: #64748b;
  line-height: 1.7;
  margin: 0;
  font-size: 14px;
  overflow: hidden;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  text-overflow: ellipsis;
}

html.dark .post-summary {
  color: #94a3b8;
}

.post-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  width: 100%;
}

.post-meta {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.meta-item {
  color: #64748b;
  font-size: 13px;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

html.dark .meta-item {
  color: #94a3b8;
}

.post-tags {
  display: flex;
  gap: 8px;
  align-items: center;
  flex-wrap: wrap;
  justify-content: flex-end;
}

.post-cover {
  flex-shrink: 0;
  width: 220px;
  aspect-ratio: 4 / 3; /* 提升可视面积，减少留白 */
  border-radius: 8px;
  overflow: hidden;
  background: #f5f5f5;
}

html.dark .post-cover {
  background: #2a2a2a;
}

.cover-image {
  display: block;
  width: 100%;
  height: 100%;
  transition: transform 0.4s;
}

.cover-image :deep(.n-image) {
  width: 100%;
  height: 100%;
  background: transparent;
}

.cover-image :deep(img) {
  width: 100%;
  height: 100%;
  object-fit: contain; /* 展示全貌，避免裁剪 */
  display: block;
  background: #f5f5f5;
}

.post-card:hover .cover-image {
  transform: scale(1.05);
}

/* 响应式布局 */
@media (max-width: 768px) {
  .post-card-content {
    flex-direction: column-reverse;
  }
  
  .post-cover {
    width: 100%;
    height: 200px;
  }
  
  .post-footer {
    flex-direction: column;
    align-items: flex-start;
  }
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

/* 深色模式卡片 */
html.dark .home-page :deep(.n-card) {
  background: rgba(30, 41, 59, 0.7);
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

html.dark .home-page :deep(.n-card):hover {
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.5);
  border-color: rgba(56, 189, 248, 0.3);
}

html.dark .post-card:hover {
  box-shadow: 0 20px 48px rgba(0, 0, 0, 0.6);
  border-color: rgba(56, 189, 248, 0.4);
}

/* 分页居中 */
.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 8px;
}
</style>

