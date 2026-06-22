<!--
  项目名称：blog-frontend
  文件名称：RecentPostsCard.vue
  创建时间：2026-02-01 20:03:19

  系统用户：Administrator
  作　　者：無以菱
  联系邮箱：huangjing510@126.com
  功能描述：最新发布文章卡片组件，展示最新发布的文章列表，包含文章封面、标题、发布时间、浏览量、点赞数等信息，点击可跳转到文章详情页。
-->
<template>
  <n-card class="recent-posts-card" :bordered="false">
    <template #header>
      <div class="card-title">
        <span>📰</span>
        <span>最新发布文章</span>
      </div>
    </template>
    <n-spin :show="loading">
      <div v-if="posts.length" class="list">
        <div v-for="item in posts" :key="item.id" class="list-item" @click="goPost(item)">
          <div class="item-cover">
            <n-image
              v-if="item.cover"
              :src="item.cover"
              :alt="item.title"
              class="cover-image"
              preview-disabled
            />
            <div v-else class="cover-placeholder">
              <n-icon :component="ImageOutline" size="20" />
            </div>
          </div>
          <div class="item-content">
            <div class="item-title" :title="item.title">{{ item.title }}</div>
            <div class="item-meta">
              <span class="meta">
                <n-icon :component="TimeOutline" size="14" />
                {{ formatDate(item.created_at, 'YYYY-MM-DD') }}
              </span>
              <span class="meta">
                <n-icon :component="EyeOutline" size="14" />
                {{ item.view_count }}
              </span>
              <span class="meta">
                <n-icon :component="HeartOutline" size="14" />
                {{ item.like_count }}
              </span>
            </div>
          </div>
        </div>
      </div>
      <n-empty v-else description="暂无最新文章" size="small" />
    </n-spin>
  </n-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { TimeOutline, EyeOutline, HeartOutline, ImageOutline } from '@vicons/ionicons5'
import { getRecentPosts } from '@/api/post'
import { getPublicDisplaySettings } from '@/api/setting'
import { formatDate } from '@/utils/format'
import type { Post } from '@/types/blog'

const router = useRouter()
const posts = ref<Post[]>([])
const loading = ref(false)
const limit = ref(5)

async function fetchRecent() {
  loading.value = true
  try {
    const res = await getRecentPosts(limit.value)
    posts.value = res.data || []
  } catch (error) {
    console.error('获取最新文章失败', error)
  } finally {
    loading.value = false
  }
}

function goPost(post: { id: number; slug: string }) {
  router.push(`/post/${post.slug}`)
}

onMounted(async () => {
  try {
    const res = await getPublicDisplaySettings()
    if (res.data?.recent_posts_limit) {
      const parsed = parseInt(res.data.recent_posts_limit, 10)
      if (!isNaN(parsed) && parsed > 0) {
        limit.value = parsed
      }
    }
  } catch {
    // 配置加载失败时使用默认值 5
  }
  fetchRecent()
})
</script>

<style scoped>
.recent-posts-card {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 16px;
}

.card-title {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-weight: 700;
  color: #1f2937;
}

.list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.list-item {
  padding: 10px 12px;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.7);
  border: 1px solid rgba(0, 0, 0, 0.04);
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  gap: 12px;
  align-items: center;
}

.list-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.08);
}

.item-cover {
  flex-shrink: 0;
  width: 120px;
  height: 80px;
  border-radius: 10px;
  overflow: hidden;
  background: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
}

.cover-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.cover-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #c0c4cc;
  background: linear-gradient(135deg, #f5f7fa 0%, #f0f2f5 100%);
}

.item-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.item-title {
  font-size: 14px;
  font-weight: 600;
  color: #111827;
  line-height: 1.5;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-clamp: 2;
  word-break: break-word;
}

.item-meta {
  display: flex;
  gap: 10px;
  align-items: center;
  color: #94a3b8;
  font-size: 12px;
  margin-top: 4px;
}

.meta {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

html.dark .recent-posts-card {
  background: rgba(30, 41, 59, 0.7);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

html.dark .card-title {
  color: #fff;
}

html.dark .list-item {
  background: rgba(255, 255, 255, 0.06);
  border-color: rgba(255, 255, 255, 0.08);
}

html.dark .item-title {
  color: #fff;
}

html.dark .item-meta {
  color: #94a3b8;
}
</style>
