<template>
  <div class="category-page">
    <n-spin :show="loading">
      <!-- 分类列表视图 -->
      <div v-if="!categoryId" class="categories-list">
        <h1 class="page-title">文章分类</h1>
        <n-grid :x-gap="16" :y-gap="16" :cols="responsive">
          <n-grid-item v-for="cat in categories" :key="cat.id">
            <n-card
              hoverable
              class="category-card"
              @click="router.push(`/category/${cat.id}`)"
              :style="{ borderTop: `4px solid ${cat.color}` }"
            >
              <div class="category-info">
                <h2 class="category-name">{{ cat.name }}</h2>
                <p v-if="cat.description" class="category-desc">{{ cat.description }}</p>
                <div class="category-count">
                  <n-icon :component="DocumentTextOutline" size="18" />
                  <span>{{ cat.post_count }} 篇文章</span>
                </div>
              </div>
            </n-card>
          </n-grid-item>
        </n-grid>
        <n-empty v-if="!loading && categories.length === 0" description="暂无分类" style="margin-top: 48px" />
      </div>

      <!-- 分类详情视图 -->
      <div v-else>
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
      </div>
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useMessage } from 'naive-ui'
import { DocumentTextOutline } from '@vicons/ionicons5'
import { getCategoryById, getCategories } from '@/api/category'
import { getPosts } from '@/api/post'
import { formatDate } from '@/utils/format'
import type { Category, Post } from '@/types/blog'

const router = useRouter()
const route = useRoute()
const message = useMessage()

const loading = ref(false)
const categories = ref<Category[]>([])
const category = ref<Category | null>(null)
const posts = ref<Post[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = 10

const categoryId = computed(() => route.params.id ? Number(route.params.id) : null)
const totalPages = computed(() => Math.ceil(total.value / pageSize))

// 响应式列数
const responsive = computed(() => {
  if (typeof window === 'undefined') return 3
  const width = window.innerWidth
  if (width < 768) return 1
  if (width < 1024) return 2
  return 3
})

// 监听路由变化
watch(() => route.params.id, (newId) => {
  // 重置状态
  category.value = null
  posts.value = []
  total.value = 0
  currentPage.value = 1
  
  if (newId) {
    fetchCategory()
    fetchPosts()
  } else {
    fetchCategories()
  }
}, { immediate: true })

onMounted(() => {
  if (categoryId.value) {
    fetchCategory()
    fetchPosts()
  } else {
    fetchCategories()
  }
})

async function fetchCategories() {
  try {
    loading.value = true
    const res = await getCategories()
    if (res.data) {
      categories.value = res.data
    }
  } catch (error: any) {
    message.error('获取分类列表失败')
  } finally {
    loading.value = false
  }
}

async function fetchCategory() {
  try {
    const res = await getCategoryById(categoryId.value!)
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
      category_id: categoryId.value!,
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
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
}

/* 分类列表样式 */
.categories-list {
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

.category-card {
  cursor: pointer;
  transition: all 0.3s ease;
  height: 100%;
}

.category-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.15);
}

html.dark .category-card:hover {
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.5);
}

.category-info {
  padding: 8px 0;
}

.category-name {
  font-size: 22px;
  font-weight: 600;
  margin: 0 0 12px 0;
  color: #1a202c;
}

html.dark .category-name {
  color: #e5e5e5;
}

.category-desc {
  color: #64748b;
  font-size: 14px;
  line-height: 1.6;
  margin: 0 0 16px 0;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
}

html.dark .category-desc {
  color: #94a3b8;
}

.category-count {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #0891b2;
  font-size: 14px;
  font-weight: 500;
}

html.dark .category-count {
  color: #38bdf8;
}

/* 分类详情样式 */
.category-header {
  text-align: center;
  padding: 24px 0;
}

.category-header h1 {
  font-size: 32px;
  margin: 0 0 16px 0;
  color: #1a202c;
}

html.dark .category-header h1 {
  color: #e5e5e5;
}

.category-header p {
  color: #64748b;
  margin: 0 0 16px 0;
}

html.dark .category-header p {
  color: #94a3b8;
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
  color: #94a3b8;
  font-size: 14px;
}

/* 响应式 */
@media (max-width: 768px) {
  .category-page {
    padding: 16px;
  }
  
  .page-title {
    font-size: 24px;
    margin-bottom: 32px;
  }
}
</style>

