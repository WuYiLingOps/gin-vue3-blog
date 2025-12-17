<template>
  <div class="moments-page">
    <div class="moments-layout">
      <div class="moments-main">
        <n-spin :show="loading">
          <div v-if="moments.length > 0" class="moments-list">
            <div v-for="moment in moments" :key="moment.id" class="moment-card">
              <!-- 左侧头像 -->
              <div class="moment-avatar">
                <n-avatar :src="moment.user.avatar" :size="56" round />
              </div>

              <!-- 右侧内容区 -->
              <div class="moment-main">
                <div class="moment-header">
                  <div class="user-info">
                    <span class="username">{{ moment.user.nickname || moment.user.username }}</span>
                    <span class="time">{{ formatDate(moment.created_at, 'YYYY-MM-DD HH:mm') }}</span>
                  </div>
                </div>

                <div class="moment-content">
                  <p>{{ moment.content }}</p>
                </div>

                <div v-if="moment.images" class="moment-images">
                  <n-image-group>
                    <n-space :size="12">
                      <n-image
                        v-for="(img, index) in parseImages(moment.images)"
                        :key="index"
                        :src="img"
                        object-fit="cover"
                        class="moment-image"
                      />
                    </n-space>
                  </n-image-group>
                </div>

                <div class="moment-footer">
                  <n-space :size="16">
                    <span class="stat-item like-button" @click="handleLike(moment)">
                      <n-icon :component="moment.liked ? Heart : HeartOutline" :class="{ liked: moment.liked }" />
                      <span v-if="moment.like_count > 0" class="like-text">{{ moment.like_count }}</span>
                    </span>
                  </n-space>
                </div>
              </div>
            </div>
          </div>

          <n-empty v-else description="还没有发布说说" style="margin: 60px 0" />
        </n-spin>

        <!-- 分页 -->
        <div v-if="total > pageSize" class="pagination-wrapper">
          <n-pagination
            v-model:page="currentPage"
            :page-count="Math.ceil(total / pageSize)"
            :page-size="pageSize"
            show-size-picker
            :page-sizes="[10, 20, 30]"
            @update:page="handlePageChange"
            @update:page-size="handlePageSizeChange"
          />
        </div>
      </div>

      <div class="sidebar-section">
        <div class="sidebar-card-wrapper">
          <AuthorCard />
        </div>
        <div class="sidebar-card-wrapper">
          <AnnouncementBoard :limit="3" />
        </div>
        <div class="sidebar-card-wrapper">
          <TagCloudWidget />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import { HeartOutline, Heart } from '@vicons/ionicons5'
import { getMoments, likeMoment } from '@/api/moment'
import type { Moment, MomentParams } from '@/api/moment'
import { formatDate } from '@/utils/format'
import { useAuthStore } from '@/stores/auth'
import AuthorCard from '@/components/AuthorCard.vue'
import AnnouncementBoard from '@/components/AnnouncementBoard.vue'
import TagCloudWidget from '@/components/TagCloudWidget.vue'

const message = useMessage()
const authStore = useAuthStore()

const loading = ref(false)
const moments = ref<Moment[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 解析图片JSON
function parseImages(images: string): string[] {
  try {
    return images ? JSON.parse(images) : []
  } catch {
    return []
  }
}

// 获取说说列表
async function fetchMoments() {
  try {
    loading.value = true
    const isAdmin = authStore.isLoggedIn && authStore.user?.role === 'admin'

    const params: MomentParams = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (!isAdmin) {
      params.status = 1
    }

    const res = await getMoments(params) as any

    if (res && res.data) {
      moments.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch (error) {
    console.error('获取说说列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 页码变化
function handlePageChange(page: number) {
  currentPage.value = page
  fetchMoments()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

// 每页数量变化
function handlePageSizeChange(size: number) {
  pageSize.value = size
  currentPage.value = 1
  fetchMoments()
}

// 点赞/取消点赞说说
async function handleLike(moment: Moment) {
  const wasLiked = moment.liked
  
  try {
    // 乐观更新 UI
    if (wasLiked) {
      // 取消点赞
      moment.like_count--
      moment.liked = false
    } else {
      // 点赞
      moment.like_count++
      moment.liked = true
    }
    
    // 调用后端 API
    await likeMoment(moment.id)
  } catch (error) {
    // 失败时回滚
    if (wasLiked) {
      moment.like_count++
      moment.liked = true
    } else {
      moment.like_count--
      moment.liked = false
    }
    console.error('操作失败:', error)
    message.error('操作失败，请重试')
  }
}

onMounted(() => {
  fetchMoments()
})
</script>

<style scoped>
.moments-page {
  min-height: calc(100vh - 180px);
  padding: 32px 0;
}

.moments-layout {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 20px;
  display: grid;
  grid-template-columns: minmax(0, 1fr) 360px;
  gap: 32px;
  align-items: start;
}

.moments-main {
  min-width: 0;
}

.moments-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* 卡片布局 - 左侧头像，右侧内容 */
.moment-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  transition: all 0.3s;
  display: flex;
  gap: 16px;
}

.moment-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

html.dark .moment-card {
  background: #1f1f1f;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

html.dark .moment-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.4);
}

/* 左侧头像区域 */
.moment-avatar {
  flex-shrink: 0;
}

.moment-avatar :deep(.n-avatar) {
  transition: all 0.3s;
}

.moment-card:hover .moment-avatar :deep(.n-avatar) {
  transform: scale(1.05);
}

/* 右侧主内容区域 */
.moment-main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.moment-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.user-info {
  display: flex;
  align-items: baseline;
  gap: 12px;
  flex-wrap: wrap;
}

.username {
  font-size: 16px;
  font-weight: 700;
  color: #1a202c;
  letter-spacing: -0.01em;
}

html.dark .username {
  color: #e5e5e5;
}

.time {
  font-size: 13px;
  color: #94a3b8;
  font-weight: 400;
}

html.dark .time {
  color: #64748b;
}

/* 说说内容 */
.moment-content {
  line-height: 1.8;
}

.moment-content p {
  font-size: 15px;
  line-height: 1.8;
  color: #334155;
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
}

html.dark .moment-content p {
  color: #cbd5e1;
}

/* 图片展示区域 */
.moment-images {
  margin-top: 8px;
}

.moment-image {
  width: 200px;
  height: 200px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  border: 1px solid rgba(0, 0, 0, 0.06);
}

.moment-image:hover {
  transform: scale(1.02);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

html.dark .moment-image {
  border-color: rgba(255, 255, 255, 0.1);
}

/* 底部互动区域 */
.moment-footer {
  padding-top: 12px;
  border-top: 1px solid #e5e7eb;
  margin-top: 4px;
}

html.dark .moment-footer {
  border-top-color: #374151;
}

.stat-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: #64748b;
  transition: all 0.3s;
}

.like-button {
  cursor: pointer;
  user-select: none;
  padding: 6px 12px;
  border-radius: 6px;
  background: transparent;
  transition: all 0.3s;
}

.like-button:hover {
  color: #ef4444;
  background: rgba(239, 68, 68, 0.08);
  transform: translateY(-1px);
}

.like-button:active {
  transform: translateY(0) scale(0.98);
}

.like-button .liked {
  color: #ef4444;
  animation: likeAnimation 0.3s ease;
}

.like-text {
  font-weight: 500;
}

@keyframes likeAnimation {
  0% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.3);
  }
  100% {
    transform: scale(1);
  }
}

html.dark .stat-item {
  color: #94a3b8;
}

html.dark .like-button:hover {
  color: #f87171;
  background: rgba(248, 113, 113, 0.12);
}

/* 分页 */
.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 40px;
  padding: 24px 0;
}

.sidebar-section {
  position: relative;
  z-index: 10;
  margin-left: 12px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

@media (max-width: 1024px) {
  .moments-layout {
    grid-template-columns: 1fr;
    padding: 0 16px;
  }

  .sidebar-section {
    display: none;
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .moment-card {
    padding: 16px;
    gap: 12px;
  }

  .moment-avatar :deep(.n-avatar) {
    width: 44px !important;
    height: 44px !important;
  }

  .username {
    font-size: 15px;
  }

  .time {
    font-size: 12px;
  }

  .moment-content p {
    font-size: 14px;
  }

  .moment-image {
    width: 150px;
    height: 150px;
  }
}

@media (max-width: 480px) {
  .moments-container {
    padding: 0 16px;
  }

  .moment-card {
    padding: 12px;
  }

  .moment-image {
    width: 120px;
    height: 120px;
  }

  .user-info {
    flex-direction: column;
    gap: 4px;
    align-items: flex-start;
  }
}
</style>

