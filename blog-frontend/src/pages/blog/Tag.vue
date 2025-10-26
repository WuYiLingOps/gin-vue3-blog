<template>
  <div class="tag-page">
    <n-spin :show="loading">
      <!-- 标签列表视图 -->
      <div v-if="!tagId" class="tags-list">
        <h1 class="page-title">文章标签</h1>
        <div class="tags-cloud">
          <n-tag
            v-for="t in tags"
            :key="t.id"
            :color="{ color: t.color || '#2196F3', textColor: '#fff' }"
            :size="getTagSize(t.post_count)"
            class="tag-item"
            @click="router.push(`/tag/${t.id}`)"
          >
            # {{ t.name }} ({{ t.post_count }})
          </n-tag>
        </div>
        <n-empty v-if="!loading && tags.length === 0" description="暂无标签" style="margin-top: 48px" />
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

// 根据文章数量返回标签大小
function getTagSize(postCount: number): 'small' | 'medium' | 'large' {
  if (postCount >= 10) return 'large'
  if (postCount >= 5) return 'medium'
  return 'small'
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
}

/* 标签列表样式 */
.tags-list {
  padding: 24px 0;
}

.page-title {
  font-size: 32px;
  font-weight: 700;
  text-align: center;
  margin: 0 0 48px 0;
  color: #1a202c;
}

html.dark .page-title {
  color: #e5e5e5;
}

.tags-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  justify-content: center;
  padding: 32px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 16px;
  backdrop-filter: blur(10px);
}

html.dark .tags-cloud {
  background: rgba(30, 41, 59, 0.5);
}

.tag-item {
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: 500;
}

.tag-item:hover {
  transform: translateY(-4px) scale(1.05);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
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
  
  .page-title {
    font-size: 24px;
    margin-bottom: 32px;
  }
  
  .tags-cloud {
    padding: 24px 16px;
    gap: 12px;
  }
}
</style>
