<template>
  <div class="archive-page">
    <n-card>
      <h1 style="text-align: center; margin-bottom: 32px">文章归档</h1>

      <n-spin :show="loading">
        <n-timeline>
          <n-timeline-item
            v-for="(group, index) in groupedPosts"
            :key="index"
            :title="group.date"
            :time="`${group.posts.length} 篇`"
          >
            <n-list>
              <n-list-item
                v-for="post in group.posts"
                :key="post.id"
                class="archive-item"
                @click="router.push(`/post/${post.id}`)"
              >
                <n-thing :title="post.title" :description="post.summary">
                  <template #footer>
                    <n-space size="small">
                      <n-tag size="small" :bordered="false" type="info">
                        {{ post.category.name }}
                      </n-tag>
                      <span style="color: #999; font-size: 12px">
                        {{ formatDate(post.created_at, 'MM-DD') }}
                      </span>
                    </n-space>
                  </template>
                </n-thing>
              </n-list-item>
            </n-list>
          </n-timeline-item>
        </n-timeline>

        <n-empty v-if="!loading && groupedPosts.length === 0" description="暂无归档" />
      </n-spin>
    </n-card>
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
  max-width: 900px;
  margin: 0 auto;
}

.archive-item {
  cursor: pointer;
  transition: all 0.2s;
}

.archive-item:hover {
  background: #f5f5f5;
}
</style>

