<template>
  <div class="archive-page">
    <div class="archive-header">
      <h1>文章归档</h1>
      <p class="archive-subtitle">共 {{ posts.length }} 篇文章</p>
    </div>

    <n-spin :show="loading">
      <div class="timeline-container">
        <div v-for="(group, index) in groupedPosts" :key="index" class="timeline-group">
          <div class="timeline-date">
            <div class="date-circle"></div>
            <h2>{{ group.date }}</h2>
            <span class="post-count">{{ group.posts.length }} 篇</span>
          </div>
          
          <div class="posts-list">
            <div
              v-for="post in group.posts"
              :key="post.id"
              class="post-card"
              @click="router.push(`/post/${post.id}`)"
            >
              <div class="post-card-content">
                <h3 class="post-title">{{ post.title }}</h3>
                <p class="post-summary">{{ post.summary }}</p>
                <div class="post-meta">
                  <n-tag size="small" :bordered="false" :color="{ color: post.category.color, textColor: '#fff' }">
                    {{ post.category.name }}
                  </n-tag>
                  <span class="post-date">
                    {{ formatDate(post.created_at, 'MM-DD') }}
                  </span>
                </div>
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
import dayjs from 'dayjs'

const router = useRouter()
const message = useMessage()

const loading = ref(false)
const posts = ref<Post[]>([])

const groupedPosts = computed(() => {
  const groups: { date: string; posts: Post[] }[] = []
  const postsByMonth: Record<string, Post[]> = {}

  posts.value.forEach(post => {
    const month = dayjs(post.created_at).format('YYYY年MM月')
    if (!postsByMonth[month]) {
      postsByMonth[month] = []
    }
    postsByMonth[month].push(post)
  })

  Object.keys(postsByMonth)
    .sort((a, b) => b.localeCompare(a))
    .forEach(month => {
      groups.push({
        date: month,
        posts: postsByMonth[month]
      })
    })

  return groups
})

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
</script>

<style scoped>
.archive-page {
  max-width: 1000px;
  margin: 0 auto;
  padding: 40px 20px;
}

.archive-header {
  text-align: center;
  margin-bottom: 60px;
}

.archive-header h1 {
  font-size: 48px;
  font-weight: 800;
  margin: 0 0 16px 0;
  background: linear-gradient(135deg, #10b981 0%, #06b6d4 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.archive-subtitle {
  font-size: 16px;
  color: #94a3b8;
  margin: 0;
}

.timeline-container {
  position: relative;
}

.timeline-group {
  position: relative;
  margin-bottom: 60px;
}

.timeline-group::before {
  content: '';
  position: absolute;
  left: 19px;
  top: 40px;
  bottom: -40px;
  width: 2px;
  background: linear-gradient(180deg, #10b981 0%, #06b6d4 100%);
  opacity: 0.3;
}

.timeline-group:last-child::before {
  display: none;
}

.timeline-date {
  position: relative;
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
  padding-left: 60px;
}

.date-circle {
  position: absolute;
  left: 0;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #10b981 0%, #06b6d4 100%);
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
}

.date-circle::after {
  content: '';
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: white;
}

.timeline-date h2 {
  font-size: 28px;
  font-weight: 700;
  margin: 0;
  color: #1a202c;
}

html.dark .timeline-date h2 {
  color: #e5e5e5;
}

.post-count {
  font-size: 14px;
  color: #10b981;
  padding: 4px 12px;
  background: rgba(16, 185, 129, 0.1);
  border-radius: 12px;
}

.posts-list {
  display: grid;
  gap: 16px;
  padding-left: 60px;
}

.post-card {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.5);
  border-radius: 16px;
  padding: 24px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

html.dark .post-card {
  background: rgba(30, 41, 59, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.post-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(16, 185, 129, 0.2);
  border-color: rgba(16, 185, 129, 0.3);
}

html.dark .post-card:hover {
  box-shadow: 0 12px 24px rgba(16, 185, 129, 0.15);
}

.post-card-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.post-title {
  font-size: 20px;
  font-weight: 600;
  margin: 0;
  color: #1a202c;
  line-height: 1.4;
}

html.dark .post-title {
  color: #e5e5e5;
}

.post-summary {
  font-size: 14px;
  color: #64748b;
  margin: 0;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

html.dark .post-summary {
  color: #94a3b8;
}

.post-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.post-date {
  font-size: 13px;
  color: #94a3b8;
}

@media (max-width: 768px) {
  .archive-page {
    padding: 24px 16px;
  }

  .archive-header h1 {
    font-size: 36px;
  }

  .timeline-date {
    padding-left: 50px;
  }

  .timeline-date h2 {
    font-size: 22px;
  }

  .posts-list {
    padding-left: 50px;
  }

  .post-card {
    padding: 20px;
  }

  .post-title {
    font-size: 18px;
  }
}
</style>

