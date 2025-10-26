<template>
  <div class="archive-page">
    <n-spin :show="loading">
      <div class="timeline-container">
        <div v-for="(group, index) in groupedPosts" :key="index" class="timeline-group">
          <!-- 年份标题 -->
          <div class="year-header">
            <h2 class="year-title">{{ group.date }}</h2>
          </div>
          
          <!-- 文章列表 -->
          <div class="posts-list">
            <div
              v-for="post in group.posts"
              :key="post.id"
              class="post-item"
              @click="router.push(`/post/${post.id}`)"
            >
              <div class="post-date-box">
                <span class="date-month">{{ formatDate(post.created_at, 'M月') }}</span>
                <span class="date-day">{{ formatDate(post.created_at, 'DD') }}</span>
              </div>
              <div class="post-content">
                <h3 class="post-title">{{ post.title }}</h3>
                <p v-if="post.summary" class="post-summary">{{ post.summary }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <n-empty v-if="!loading && groupedPosts.length === 0" description="暂无归档" />
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { getPosts } from '@/api/post'
import { formatDate } from '@/utils/format'
import type { Post } from '@/types/blog'

const router = useRouter()
const message = useMessage()

const loading = ref(false)
const posts = ref<Post[]>([])

onMounted(() => {
  fetchPosts()
})

async function fetchPosts() {
  try {
    loading.value = true
    const res = await getPosts({
      page: 1,
      page_size: 1000,
      status: 1
    })

    if (res.data) {
      posts.value = res.data.list
    }
  } catch (error: any) {
    message.error('获取文章列表失败')
  } finally {
    loading.value = false
  }
}

// 按年份分组
const groupedPosts = computed(() => {
  const groups: { [key: string]: Post[] } = {}

  posts.value.forEach(post => {
    const year = new Date(post.created_at).getFullYear().toString()
    if (!groups[year]) {
      groups[year] = []
    }
    groups[year].push(post)
  })

  return Object.keys(groups)
    .sort((a, b) => parseInt(b) - parseInt(a))
    .map(year => ({
      date: year,
      posts: groups[year].sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
    }))
})
</script>

<style scoped>
.archive-page {
  max-width: 1000px;
  margin: 0 auto;
  padding: 40px 20px;
}

.timeline-container {
  position: relative;
}

.timeline-group {
  position: relative;
  margin-bottom: 48px;
  padding-left: 60px;
  border-left: 3px solid #0891b2;
}

html.dark .timeline-group {
  border-left-color: #38bdf8;
}

/* 年份标题 */
.year-header {
  position: relative;
  margin-bottom: 32px;
  margin-left: -60px;
  padding-left: 20px;
}

.year-title {
  font-size: 32px;
  font-weight: 700;
  margin: 0;
  padding-left: 20px;
  color: #1a202c;
  display: inline-block;
}

html.dark .year-title {
  color: #e5e5e5;
}

/* 年份标题左侧的横线装饰 */
.year-header::before {
  content: '';
  position: absolute;
  left: -3px;
  top: 50%;
  transform: translateY(-50%);
  width: 30px;
  height: 3px;
  background: #0891b2;
}

html.dark .year-header::before {
  background: #38bdf8;
}

/* 文章列表 */
.posts-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.post-item {
  display: flex;
  gap: 24px;
  padding: 20px 24px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid rgba(8, 145, 178, 0.1);
}

html.dark .post-item {
  background: rgba(30, 41, 59, 0.8);
  border-color: rgba(56, 189, 248, 0.1);
}

.post-item:hover {
  transform: translateX(8px);
  box-shadow: 0 4px 16px rgba(8, 145, 178, 0.15);
  border-color: #0891b2;
}

html.dark .post-item:hover {
  box-shadow: 0 4px 16px rgba(56, 189, 248, 0.2);
  border-color: #38bdf8;
}

/* 日期框 */
.post-date-box {
  flex-shrink: 0;
  width: 80px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 8px;
  padding: 12px;
  border: 1px solid rgba(8, 145, 178, 0.2);
}

html.dark .post-date-box {
  background: rgba(15, 23, 42, 0.8);
  border-color: rgba(56, 189, 248, 0.2);
}

.date-month {
  font-size: 13px;
  color: #0891b2;
  font-weight: 600;
}

html.dark .date-month {
  color: #38bdf8;
}

.date-day {
  font-size: 28px;
  font-weight: 700;
  color: #0e7490;
  line-height: 1;
  margin-top: 4px;
}

html.dark .date-day {
  color: #0ea5e9;
}

/* 文章内容 */
.post-content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.post-title {
  font-size: 18px;
  font-weight: 600;
  margin: 0 0 8px 0;
  color: #1a202c;
  line-height: 1.4;
}

html.dark .post-title {
  color: #e5e5e5;
}

.post-summary {
  color: #64748b;
  line-height: 1.6;
  margin: 0;
  font-size: 14px;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
}

html.dark .post-summary {
  color: #94a3b8;
}

/* 响应式 */
@media (max-width: 768px) {
  .archive-page {
    padding: 24px 16px;
  }

  .timeline-group {
    padding-left: 40px;
  }

  .year-header {
    margin-left: -40px;
  }

  .year-title {
    font-size: 24px;
  }

  .post-item {
    flex-direction: column;
    gap: 16px;
  }

  .post-date-box {
    width: 100%;
    flex-direction: row;
    justify-content: center;
    gap: 8px;
  }
}
</style>
