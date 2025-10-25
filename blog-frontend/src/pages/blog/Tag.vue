<template>
  <div class="tag-page">
    <n-spin :show="loading">
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
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useMessage } from 'naive-ui'
import { getTagById, getPostsByTag } from '@/api/tag'
import { formatDate } from '@/utils/format'
import type { Tag, Post } from '@/types/blog'

const router = useRouter()
const route = useRoute()
const message = useMessage()

const loading = ref(false)
const tag = ref<Tag | null>(null)
const posts = ref<Post[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = 10

const tagId = computed(() => Number(route.params.id))
const totalPages = computed(() => Math.ceil(total.value / pageSize))

onMounted(() => {
  fetchTag()
  fetchPosts()
})

async function fetchTag() {
  try {
    const res = await getTagById(tagId.value)
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
    const res = await getPostsByTag(tagId.value, {
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
  max-width: 900px;
  margin: 0 auto;
}

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
  font-size: 14px;
}
</style>

