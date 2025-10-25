<template>
  <div class="category-page">
    <n-spin :show="loading">
      <n-card v-if="category">
        <div class="category-header">
          <h1>{{ category.name }}</h1>
          <p v-if="category.description">{{ category.description }}</p>
          <n-tag :color="{ color: category.color, textColor: '#fff' }">
            {{ category.post_count }} 篇文章
          </n-tag>
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
              <span>{{ formatDate(post.created_at, 'YYYY-MM-DD') }}</span>
              <n-divider vertical />
              <span>{{ post.view_count }} 阅读</span>
            </n-space>
          </div>
        </n-card>

        <n-empty v-if="!loading && posts.length === 0" description="该分类下暂无文章" />

        <n-pagination
          v-if="total > 0"
          v-model:page="currentPage"
          :page-count="totalPages"
          @update:page="handlePageChange"
        />
      </n-space>
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useMessage } from 'naive-ui'
import { getCategoryById } from '@/api/category'
import { getPosts } from '@/api/post'
import { formatDate } from '@/utils/format'
import type { Category, Post } from '@/types/blog'

const router = useRouter()
const route = useRoute()
const message = useMessage()

const loading = ref(false)
const category = ref<Category | null>(null)
const posts = ref<Post[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = 10

const categoryId = computed(() => Number(route.params.id))
const totalPages = computed(() => Math.ceil(total.value / pageSize))

onMounted(() => {
  fetchCategory()
  fetchPosts()
})

async function fetchCategory() {
  try {
    const res = await getCategoryById(categoryId.value)
    if (res.data) {
      category.value = res.data
    }
  } catch (error: any) {
    message.error('获取分类信息失败')
  }
}

async function fetchPosts() {
  try {
    loading.value = true
    const res = await getPosts({
      page: currentPage.value,
      page_size: pageSize,
      category_id: categoryId.value,
      status: 1
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
.category-page {
  max-width: 900px;
  margin: 0 auto;
}

.category-header {
  text-align: center;
  padding: 24px 0;
}

.category-header h1 {
  font-size: 32px;
  margin: 0 0 16px 0;
}

.category-header p {
  color: #666;
  margin: 0 0 16px 0;
}

.post-card {
  cursor: pointer;
  transition: all 0.3s;
}

.post-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.post-title {
  font-size: 20px;
  margin: 0 0 12px 0;
}

.post-summary {
  color: #666;
  line-height: 1.6;
  margin: 0 0 12px 0;
}

.post-meta {
  color: #999;
  font-size: 14px;
}
</style>

