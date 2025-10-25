<template>
  <div class="moments-page">
    <div class="moments-container">
      <n-spin :show="loading">
        <div v-if="moments.length > 0" class="moments-list">
          <div v-for="moment in moments" :key="moment.id" class="moment-card">
            <div class="moment-header">
              <n-space :size="12">
                <n-avatar :src="moment.user.avatar" :size="48" round />
                <div class="user-info">
                  <div class="username">{{ moment.user.nickname || moment.user.username }}</div>
                  <div class="time">{{ formatDate(moment.created_at, 'YYYY-MM-DD HH:mm') }}</div>
                </div>
              </n-space>
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
                  {{ moment.like_count }}
                </span>
              </n-space>
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import { HeartOutline, Heart } from '@vicons/ionicons5'
import { getMoments, likeMoment } from '@/api/moment'
import type { Moment } from '@/api/moment'
import { formatDate } from '@/utils/format'

const message = useMessage()

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
    const res = await getMoments({
      page: currentPage.value,
      page_size: pageSize.value,
      status: 1
    })

    if (res.data) {
      moments.value = res.data.list
      total.value = res.data.total
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
  try {
    const wasLiked = moment.liked
    
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

.moments-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 24px;
}

.moments-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.moment-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  transition: all 0.3s;
}

.moment-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

html.dark .moment-card {
  background: #1f1f1f;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

html.dark .moment-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.4);
}

.moment-header {
  margin-bottom: 16px;
}

.user-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.username {
  font-size: 16px;
  font-weight: 600;
  color: #1a202c;
}

html.dark .username {
  color: #e5e5e5;
}

.time {
  font-size: 13px;
  color: #94a3b8;
}

.moment-content {
  margin-bottom: 16px;
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

.moment-images {
  margin-bottom: 16px;
}

.moment-image {
  width: 200px;
  height: 200px;
  border-radius: 8px;
  cursor: pointer;
}

.moment-footer {
  padding-top: 12px;
  border-top: 1px solid #e5e7eb;
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
}

.like-button:hover {
  color: #ef4444;
  transform: scale(1.1);
}

.like-button:active {
  transform: scale(0.95);
}

.like-button .liked {
  color: #ef4444;
  animation: likeAnimation 0.3s ease;
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
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 40px;
  padding: 24px 0;
}
</style>

