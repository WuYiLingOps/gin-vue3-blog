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
        >
          <div class="category-item-content">
            <div class="category-name-wrapper">
              <span class="category-dot" :style="{ backgroundColor: category.color || '#2196F3' }"></span>
              <span class="category-name">{{ category.name }}</span>
            </div>
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
  gap: 8px;
}

.category-item {
  cursor: pointer;
  padding: 10px 12px;
  border-radius: 8px;
  transition: all 0.2s ease;
  background: rgba(0, 0, 0, 0.02);
}

.category-item:hover {
  background: rgba(8, 145, 178, 0.08);
  transform: translateX(4px);
}

html.dark .category-item {
  background: rgba(255, 255, 255, 0.05);
}

html.dark .category-item:hover {
  background: rgba(56, 189, 248, 0.15);
}

.category-item-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.category-name-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  min-width: 0;
}

.category-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.category-name {
  font-size: 14px;
  color: #1a202c;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

html.dark .category-name {
  color: #e5e5e5;
}

.category-count {
  font-size: 12px;
  color: #64748b;
  font-weight: 500;
  flex-shrink: 0;
}

html.dark .category-count {
  color: #94a3b8;
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

