<template>
  <n-card
    class="announcement-card"
    :bordered="false"
  >
    <template #header>
      <div class="card-title">
        <span>📢</span>
        <span>公告栏</span>
      </div>
    </template>
    <n-spin :show="loading">
      <div class="announcement-list" v-if="announcements.length">
        <div
          v-for="item in announcements"
          :key="item.id"
          class="announcement-item"
          @click="openDetail(item)"
        >
          <div class="item-header">
            <div class="item-title">
              <n-tag
                v-if="item.priority === 1"
                type="error"
                size="small"
                round
                class="pin-tag"
              >
                置顶
              </n-tag>
              <span class="title-text">
                {{ getTitle(item) }}
              </span>
            </div>
          </div>
          <div class="item-content" v-html="item.content"></div>
          <div class="item-footer">
            <span class="item-time">
              <n-time :time="new Date(item.created_at)" format="yyyy-MM-dd HH:mm" />
            </span>
          </div>
        </div>
      </div>
      <n-empty v-else description="暂无公告" />
    </n-spin>

    <!-- 欢迎图片 -->
    <div class="welcome-image-wrapper">
      <div class="welcome-image-card">
        <img src="/公告栏.gif" alt="欢迎" class="welcome-image" />
      </div>
    </div>

    <n-modal v-model:show="showDetail" preset="card" style="max-width: 640px" :bordered="false">
      <template #header>
        <div class="detail-header">
          <span class="detail-title">
            {{ selected?.title || '系统公告' }}
          </span>
          <n-tag v-if="selected?.priority === 1" type="error" size="small" round>置顶</n-tag>
        </div>
      </template>
      <div class="detail-meta">
        <span class="meta-time">
          发布时间：
          <n-time :time="selectedTime" format="yyyy-MM-dd HH:mm" />
        </span>
        <span class="meta-author">来源：{{ selected?.username || '系统消息' }}</span>
      </div>
      <n-scrollbar style="max-height: 420px; margin-top: 12px">
        <div class="detail-content" v-html="selected?.content"></div>
      </n-scrollbar>
    </n-modal>
  </n-card>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useMessage } from 'naive-ui'
import { getAnnouncements, type Announcement } from '@/api/blog'

interface Props {
  limit?: number
}

const props = withDefaults(defineProps<Props>(), {
  limit: 3
})

const announcements = ref<Announcement[]>([])
const loading = ref(false)
const showDetail = ref(false)
const selected = ref<Announcement | null>(null)
const message = useMessage()

const selectedTime = computed(() => (selected.value ? new Date(selected.value.created_at) : new Date()))

function stripHTML(html: string) {
  const div = document.createElement('div')
  div.innerHTML = html
  return div.textContent || div.innerText || ''
}

function getTitle(item: Announcement) {
  // 若未来支持独立标题，可直接返回 title 字段；当前用内容前若干字符代替
  const text = stripHTML(item.content)
  // 取第一行作为标题，如果没有换行则取前30个字符
  const firstLine = text.split('\n')[0] || text
  return firstLine ? (firstLine.length > 30 ? `${firstLine.slice(0, 30)}...` : firstLine) : '系统公告'
}

async function fetchAnnouncements() {
  loading.value = true
  try {
    const res = await getAnnouncements(props.limit)
    announcements.value = res.data || []
  } catch (err: any) {
    message.error(err?.response?.data?.message || '获取公告失败')
  } finally {
    loading.value = false
  }
}

function openDetail(item: Announcement) {
  selected.value = item
  showDetail.value = true
}

onMounted(fetchAnnouncements)
</script>

<style scoped>
.announcement-card {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.5);
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
}

.card-title {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-weight: 700;
  color: #1f2937;
}

.announcement-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.announcement-item {
  padding: 12px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.7);
  border: 1px solid rgba(0, 0, 0, 0.04);
  cursor: pointer;
  transition: all 0.2s ease;
}

.announcement-item:hover {
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.08);
  transform: translateY(-2px);
}

.item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-bottom: 4px;
}

.item-title {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-weight: 700;
  color: #1f2937;
}

.title-text {
  max-width: 100%;
  display: inline-block;
  white-space: normal;
  word-wrap: break-word;
  word-break: break-word;
}

.pin-tag {
  flex-shrink: 0;
}

.item-time {
  color: #94a3b8;
  font-size: 12px;
  white-space: nowrap;
}

.item-content {
  color: #4b5563;
  font-size: 13px;
  line-height: 1.8;
  margin-top: 8px;
  min-height: 20px;
  white-space: pre-wrap;
  word-wrap: break-word;
  word-break: break-word;
  overflow-wrap: break-word;
}

.item-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 6px;
}


.detail-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.detail-title {
  font-size: 16px;
  font-weight: 700;
}

.detail-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #94a3b8;
  font-size: 13px;
  gap: 12px;
}

.meta-author {
  white-space: nowrap;
}

.detail-content {
  font-size: 14px;
  line-height: 1.7;
  color: #1f2937;
  word-break: break-word;
}

html.dark .announcement-card {
  background: rgba(30, 41, 59, 0.7);
  border: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
}

html.dark .card-title {
  color: #fff;
}

html.dark .announcement-item {
  background: rgba(255, 255, 255, 0.03);
  border-color: rgba(255, 255, 255, 0.06);
}

html.dark .item-title {
  color: #fff;
}

html.dark .item-content {
  color: #e5e7eb;
}

html.dark .item-time {
  color: #94a3b8;
}

html.dark .detail-content {
  color: #e2e8f0;
}

.welcome-image-wrapper {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid rgba(0, 0, 0, 0.06);
}

html.dark .welcome-image-wrapper {
  border-top-color: rgba(255, 255, 255, 0.08);
}

.welcome-image-card {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
  padding: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.06);
  transition: all 0.3s ease;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.welcome-image-card:hover {
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

html.dark .welcome-image-card {
  background: rgba(30, 41, 59, 0.8);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

html.dark .welcome-image-card:hover {
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.5);
}

.welcome-image {
  width: 100%;
  max-width: 100%;
  height: auto;
  display: block;
  border-radius: 8px;
  object-fit: contain;
  image-rendering: auto;
  -webkit-backface-visibility: hidden;
  backface-visibility: hidden;
  -webkit-transform: translateZ(0);
  transform: translateZ(0);
}
</style>

