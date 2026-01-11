<template>
  <div class="category-list-card">
    <n-card title="📂 分类列表" size="small" :bordered="false" class="category-list">
      <div v-if="loading" class="category-loading">加载中...</div>
      <div v-else-if="error" class="category-error">{{ error }}</div>
      <div v-else class="categories-wrap">
        <div
          v-for="category in categories"
          :key="category.id"
          class="category-item"
          @click="handleClick(category.id)"
          :style="{ '--category-color': category.color || '#0891b2' }"
        >
          <div class="category-item-indicator" :style="{ backgroundColor: category.color || '#0891b2' }"></div>
          <div class="category-item-content">
            <span class="category-name">{{ category.name }}</span>
            <span class="category-count">{{ category.post_count }}</span>
          </div>
        </div>
        <n-empty v-if="!categories.length" size="small" description="暂无分类" />
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { getCategories } from '@/api/category'
import type { Category } from '@/types/blog'

const router = useRouter()
const message = useMessage()

const loading = ref(false)
const error = ref('')
const categories = ref<Category[]>([])

async function fetchCategories() {
  loading.value = true
  error.value = ''
  try {
    const res = await getCategories()
    categories.value = res.data || []
    // 按文章数量降序排序
    categories.value.sort((a, b) => b.post_count - a.post_count)
  } catch (e: any) {
    error.value = e.message || '获取分类列表失败'
    message.error(error.value)
  } finally {
    loading.value = false
  }
}

function handleClick(id: number) {
  router.push(`/category/${id}`)
}

onMounted(() => {
  fetchCategories()
})
</script>

<style scoped>
.category-list-card {
  width: 100%;
}

.category-list {
  width: 100%;
}

.category-loading,
.category-error {
  padding: 12px;
  font-size: 12px;
  color: #64748b;
}

.category-error {
  color: #d14343;
}

.categories-wrap {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.category-item {
  position: relative;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 8px;
  transition: all 0.25s ease;
  background: rgba(0, 0, 0, 0.02);
  border: 1px solid transparent;
  display: flex;
  align-items: center;
  gap: 10px;
  overflow: hidden;
}

html.dark .category-item {
  background: rgba(255, 255, 255, 0.03);
}

.category-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 2px;
  background: var(--category-color);
  transform: scaleY(0);
  transition: transform 0.25s ease;
}

.category-item:hover {
  transform: translateX(4px);
  background: linear-gradient(90deg, color-mix(in srgb, var(--category-color) 8%, transparent) 0%, transparent 100%);
  border-color: color-mix(in srgb, var(--category-color) 20%, transparent);
}

html.dark .category-item:hover {
  background: linear-gradient(90deg, color-mix(in srgb, var(--category-color) 12%, transparent) 0%, transparent 100%);
}

.category-item:hover::before {
  transform: scaleY(1);
}

.category-item-indicator {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
  box-shadow: 0 0 0 2px rgba(255, 255, 255, 0.8), 0 0 0 3px color-mix(in srgb, var(--category-color) 20%, transparent);
  transition: all 0.25s ease;
}

html.dark .category-item-indicator {
  box-shadow: 0 0 0 2px rgba(30, 30, 30, 0.8), 0 0 0 3px color-mix(in srgb, var(--category-color) 20%, transparent);
}

.category-item:hover .category-item-indicator {
  transform: scale(1.3);
  box-shadow: 0 0 0 2px rgba(255, 255, 255, 0.8), 0 0 0 4px color-mix(in srgb, var(--category-color) 30%, transparent), 0 0 8px color-mix(in srgb, var(--category-color) 40%, transparent);
}

html.dark .category-item:hover .category-item-indicator {
  box-shadow: 0 0 0 2px rgba(30, 30, 30, 0.8), 0 0 0 4px color-mix(in srgb, var(--category-color) 30%, transparent), 0 0 8px color-mix(in srgb, var(--category-color) 40%, transparent);
}

.category-item-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex: 1;
  min-width: 0;
  gap: 8px;
}

.category-name {
  font-size: 14px;
  color: #1a202c;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  transition: color 0.25s ease;
}

html.dark .category-name {
  color: #e5e5e5;
}

.category-item:hover .category-name {
  color: var(--category-color);
}

.category-count {
  font-size: 12px;
  color: #64748b;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 10px;
  background: rgba(0, 0, 0, 0.04);
  transition: all 0.25s ease;
  flex-shrink: 0;
}

html.dark .category-count {
  color: #94a3b8;
  background: rgba(255, 255, 255, 0.08);
}

.category-item:hover .category-count {
  background: color-mix(in srgb, var(--category-color) 15%, transparent);
  color: var(--category-color);
}

@media (max-width: 1024px) {
  .category-list :deep(.n-card__content) {
    padding: 16px !important;
  }
}

@media (max-width: 768px) {
  .category-list-card {
    min-width: 0;
    overflow: visible;
  }

  .category-list {
    min-width: 0;
    overflow: visible;
  }

  .categories-wrap {
    min-width: 0;
    overflow: visible;
  }
}
</style>

