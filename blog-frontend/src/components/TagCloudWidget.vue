<template>
  <div class="tag-cloud-card">
    <n-card title="🏷️ 标签列表" size="small" :bordered="false" class="tag-cloud">
      <div v-if="loading" class="tag-loading">加载中...</div>
      <div v-else-if="error" class="tag-error">{{ error }}</div>
      <div v-else class="tags-wrap">
        <div
          v-for="tag in tags"
          :key="tag.id"
          class="tag-chip"
          :class="sizeClass(tag.post_count)"
          :style="chipStyle(tag)"
          @click="handleClick(tag.id)"
        >
          #{{ tag.name }}
        </div>
        <n-empty v-if="!tags.length" size="small" description="暂无标签" />
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { getTags } from '@/api/tag'
import type { Tag } from '@/types/blog'

const router = useRouter()
const message = useMessage()

const loading = ref(false)
const error = ref('')
const tags = ref<Tag[]>([])

function sizeClass(count: number) {
  if (count >= 15) return 'tag-lg'
  if (count >= 8) return 'tag-md'
  return 'tag-sm'
}

function chipStyle(tag: Tag) {
  const bg = tag.color || '#2196F3'
  const text = tag.text_color || '#fff'
  const fontSize = tag.font_size ? `${tag.font_size}px` : undefined
  return {
    backgroundColor: bg,
    color: text,
    fontSize
  }
}

async function fetchTags() {
  loading.value = true
  error.value = ''
  try {
    const res = await getTags()
    tags.value = res.data || []
  } catch (e: any) {
    error.value = e.message || '获取标签失败'
    message.error(error.value)
  } finally {
    loading.value = false
  }
}

function handleClick(id: number) {
  router.push(`/tag/${id}`)
}

onMounted(() => {
  fetchTags()
})
</script>

<style scoped>
.tag-cloud-card {
  width: 100%;
}

.tag-cloud {
  width: 100%;
}

.tag-loading,
.tag-error {
  padding: 12px;
  font-size: 12px;
  color: #64748b;
}

.tag-error {
  color: #d14343;
}

.tags-wrap {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.tag-chip {
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 14px;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.08);
  line-height: 1.1;
}

.tag-chip:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 14px rgba(0, 0, 0, 0.12);
}

.tag-sm {
  font-size: 13px;
}

.tag-md {
  font-size: 14px;
}

.tag-lg {
  font-size: 16px;
  font-weight: 700;
}

@media (max-width: 1024px) {
  .tag-cloud :deep(.n-card__content) {
    padding: 16px !important;
  }
}

@media (max-width: 768px) {
  .tag-cloud-card {
    min-width: 0; /* 防止内容被压缩 */
    overflow: visible; /* 确保内容不被裁剪 */
  }

  .tag-cloud {
    min-width: 0; /* 防止内容被压缩 */
    overflow: visible; /* 确保内容不被裁剪 */
  }

  .tags-wrap {
    min-width: 0; /* 防止内容被压缩 */
    overflow: visible; /* 确保内容不被裁剪 */
  }
}
</style>

