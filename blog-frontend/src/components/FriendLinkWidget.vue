<template>
  <n-card title="🔗 友情链接" size="small" :bordered="false" class="friendlink-widget">
    <div v-if="loading" class="friendlink-loading">加载中...</div>
    <div v-else-if="error" class="friendlink-error">{{ error }}</div>
    <div v-else class="friendlinks-list">
      <div
        v-for="link in friendLinks"
        :key="link.id"
        class="friendlink-item"
        @click="handleClick(link.url)"
      >
        <div class="friendlink-header">
          <div class="friendlink-icon">
            <n-image
              v-if="link.icon"
              :src="link.icon"
              :alt="link.name"
              width="32"
              height="32"
              object-fit="cover"
              :preview-disabled="true"
              fallback-src="/logo.jpg"
            />
            <div v-else class="icon-placeholder">
              <n-icon :component="LinkOutline" size="20" />
            </div>
          </div>
          <div class="friendlink-info">
            <div class="friendlink-name">{{ link.name }}</div>
            <div v-if="link.description" class="friendlink-desc">{{ link.description }}</div>
          </div>
        </div>
        <div v-if="link.screenshot" class="friendlink-screenshot">
          <n-image
            :src="link.screenshot"
            :alt="link.name"
            width="100%"
            height="120"
            object-fit="cover"
            :preview-disabled="true"
            class="screenshot-image"
          />
        </div>
      </div>
      <n-empty v-if="!friendLinks.length" size="small" description="暂无友链" />
    </div>
  </n-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useMessage } from 'naive-ui'
import { LinkOutline } from '@vicons/ionicons5'
import { getFriendLinks } from '@/api/friendlink'
import type { FriendLink } from '@/api/friendlink'

const message = useMessage()

const loading = ref(false)
const error = ref('')
const friendLinks = ref<FriendLink[]>([])

function handleClick(url: string) {
  window.open(url, '_blank', 'noopener,noreferrer')
}

async function fetchFriendLinks() {
  loading.value = true
  error.value = ''
  try {
    const res = await getFriendLinks()
    friendLinks.value = res.data || []
  } catch (e: any) {
    error.value = e.message || '获取友链失败'
    message.error(error.value)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchFriendLinks()
})
</script>

<style scoped>
.friendlink-widget {
  width: 100%;
}

.friendlink-loading,
.friendlink-error {
  padding: 12px;
  font-size: 12px;
  color: #64748b;
  text-align: center;
}

.friendlink-error {
  color: #d14343;
}

.friendlinks-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.friendlink-item {
  padding: 12px;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.7);
  border: 1px solid rgba(0, 0, 0, 0.04);
  cursor: pointer;
  transition: all 0.2s ease;
}

.friendlink-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.08);
  background: rgba(255, 255, 255, 0.9);
}

.friendlink-header {
  display: flex;
  gap: 10px;
  align-items: flex-start;
  margin-bottom: 8px;
}

.friendlink-icon {
  flex-shrink: 0;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  overflow: hidden;
  background: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-placeholder {
  color: #94a3b8;
}

.friendlink-info {
  flex: 1;
  min-width: 0;
}

.friendlink-name {
  font-size: 14px;
  font-weight: 600;
  color: #111827;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.friendlink-desc {
  font-size: 12px;
  color: #64748b;
  line-height: 1.5;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-clamp: 2;
}

.friendlink-screenshot {
  margin-top: 8px;
  border-radius: 8px;
  overflow: hidden;
  background: #f5f5f5;
}

.screenshot-image {
  width: 100%;
  display: block;
}

html.dark .friendlink-item {
  background: rgba(255, 255, 255, 0.06);
  border-color: rgba(255, 255, 255, 0.08);
}

html.dark .friendlink-item:hover {
  background: rgba(255, 255, 255, 0.1);
}

html.dark .friendlink-name {
  color: #e5e7eb;
}

html.dark .friendlink-desc {
  color: #94a3b8;
}

html.dark .friendlink-icon {
  background: rgba(255, 255, 255, 0.1);
}

html.dark .friendlink-screenshot {
  background: rgba(255, 255, 255, 0.05);
}
</style>

