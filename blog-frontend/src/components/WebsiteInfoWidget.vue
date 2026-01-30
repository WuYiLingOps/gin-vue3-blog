<template>
  <div class="website-info-card">
    <n-card title="📊 网站资讯" size="small" :bordered="false" class="website-info">
      <div v-if="loading" class="info-loading">加载中...</div>
      <div v-else-if="error" class="info-error">{{ error }}</div>
      <div v-else class="info-list">
        <div class="info-item">
          <span class="info-label">本站总字数</span>
          <span class="info-value">{{ formatNumber(info.total_words) }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">本站访客数</span>
          <span class="info-value">{{ formatNumber(info.total_visitors) }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">本站总访问量</span>
          <span class="info-value">{{ formatNumber(info.total_views) }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">最后更新时间</span>
          <span class="info-value">{{ formatLastUpdateTime(info.last_update_time) }}</span>
        </div>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useMessage } from 'naive-ui'
import { getWebsiteInfo, type WebsiteInfo } from '@/api/blog'

const message = useMessage()

const loading = ref(false)
const error = ref('')
const info = ref<WebsiteInfo>({
  total_words: 0,
  total_visitors: 0,
  total_views: 0,
  last_update_time: ''
})

function formatNumber(num: number): string {
  if (num < 1000) return num.toString()
  if (num < 10000) return (num / 1000).toFixed(1) + 'k'
  if (num < 1000000) return (num / 10000).toFixed(1) + 'w'
  return (num / 1000000).toFixed(1) + 'm'
}

function formatLastUpdateTime(timeStr: string): string {
  if (!timeStr) return '暂无'
  try {
    const date = new Date(timeStr)
    const now = new Date()
    const diff = now.getTime() - date.getTime()

    // 未来时间或解析异常导致的负数：视为刚刚
    if (diff <= 0) return '刚刚'

    const days = Math.floor(diff / (1000 * 60 * 60 * 24))
    
    if (days === 0) {
      const hours = Math.floor(diff / (1000 * 60 * 60))
      if (hours === 0) {
        const minutes = Math.floor(diff / (1000 * 60))
        return minutes <= 0 ? '刚刚' : `${minutes} 分钟前`
      }
      return `${hours} 小时前`
    }

    // 超过 1 天：统一显示具体天数，不显示周/月/年前
    return `${days} 天前`
  } catch {
    return timeStr
  }
}

async function fetchWebsiteInfo() {
  loading.value = true
  error.value = ''
  try {
    const res = await getWebsiteInfo()
    info.value = res.data || info.value
  } catch (e: any) {
    error.value = e.message || '获取网站资讯失败'
    message.error(error.value)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchWebsiteInfo()
})
</script>

<style scoped>
.website-info-card {
  width: 100%;
}

.website-info {
  width: 100%;
}

.info-loading,
.info-error {
  padding: 12px;
  font-size: 12px;
  color: #64748b;
}

.info-error {
  color: #d14343;
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.info-item:last-child {
  border-bottom: none;
}

html.dark .info-item {
  border-bottom-color: rgba(255, 255, 255, 0.08);
}

.info-label {
  font-size: 13px;
  color: #64748b;
  font-weight: 400;
}

html.dark .info-label {
  color: #94a3b8;
}

.info-value {
  font-size: 14px;
  color: #1a202c;
  font-weight: 500;
  text-align: right;
}

html.dark .info-value {
  color: #e5e7eb;
}

@media (max-width: 1024px) {
  .website-info :deep(.n-card__content) {
    padding: 16px !important;
  }
}

@media (max-width: 768px) {
  .website-info-card {
    min-width: 0;
    overflow: visible;
  }

  .website-info {
    min-width: 0;
    overflow: visible;
  }

  .info-list {
    min-width: 0;
    overflow: visible;
  }
}
</style>
