<template>
  <n-card class="hot-posts-card" title="热门文章推荐" :bordered="false">
    <n-spin :show="loading">
      <div v-if="posts.length" class="list">
        <div
          v-for="item in posts"
          :key="item.id"
          class="list-item"
          @click="goPost(item.id)"
        >
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
      <n-empty v-else description="暂无热门文章" size="small" />
    </n-spin>
  </n-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { TimeOutline, EyeOutline, HeartOutline } from '@vicons/ionicons5'
import { getHotPosts } from '@/api/post'
import { formatDate } from '@/utils/format'
import type { Post } from '@/types/blog'

const router = useRouter()
const posts = ref<Post[]>([])
const loading = ref(false)
const LIMIT = 5

async function fetchHot() {
  loading.value = true
  try {
    const res = await getHotPosts(LIMIT)
    posts.value = res.data || []
  } catch (error) {
    console.error('获取热门文章失败', error)
  } finally {
    loading.value = false
  }
}

function goPost(id: number) {
  router.push(`/post/${id}`)
}

onMounted(fetchHot)
</script>

<style scoped>
.hot-posts-card {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 16px;
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
}

.list-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.08);
}

.item-title {
  font-size: 14px;
  font-weight: 600;
  color: #111827;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
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

html.dark .hot-posts-card {
  background: rgba(30, 41, 59, 0.7);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

html.dark .list-item {
  background: rgba(255, 255, 255, 0.06);
  border-color: rgba(255, 255, 255, 0.08);
}

html.dark .item-title {
  color: #e5e7eb;
}
</style>

