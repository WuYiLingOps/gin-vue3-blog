<!--
 * @ProjectName: go-vue3-blog
 * @FileName: CategoryManage.vue
 * @CreateTime: 2026-02-02 11:43:17
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 分类管理页面组件，提供文章分类的增删改查功能
 -->
<template>
  <div class="category-manage-page">
    <div class="header">
      <h1>分类管理</h1>
      <n-button type="primary" :size="isMobile ? 'small' : 'medium'" @click="showModal = true">
        <template #icon>
          <n-icon :component="AddOutline" />
        </template>
        <span v-if="!isMobile">新建分类</span>
        <span v-else>新建</span>
      </n-button>
    </div>

    <div v-if="isMobile" class="card-list">
      <n-card v-for="category in categories" :key="category.id" class="list-card" size="small">
        <template #header>
          <n-space align="center">
            <n-tag :color="{ color: category.color, textColor: '#fff' }" size="small">
              {{ category.name }}
            </n-tag>
          </n-space>
        </template>
        <div class="card-content">
          <div class="info-item">
            <span class="label">描述：</span>
            <span class="value">{{ category.description || '-' }}</span>
          </div>
          <div class="info-item">
            <span class="label">文章数：</span>
            <span class="value">{{ category.post_count }} 篇</span>
          </div>
          <div class="info-item">
            <span class="label">排序：</span>
            <span class="value">{{ category.sort }}</span>
          </div>
        </div>
        <template #footer>
          <n-space justify="end" size="small">
            <n-button size="small" @click="handleEdit(category)">编辑</n-button>
            <n-button size="small" type="error" @click="handleDelete(category.id)">删除</n-button>
          </n-space>
        </template>
      </n-card>
    </div>

    <n-data-table 
      v-else
      :columns="columns" 
      :data="categories" 
      :loading="loading"
      :single-line="false"
    />

    <!-- 创建/编辑对话框 -->
    <n-modal 
      v-model:show="showModal" 
      preset="card" 
      :title="editingId ? '编辑分类' : '新建分类'" 
      :style="{ width: isMobile ? '95%' : '500px', maxWidth: isMobile ? '95vw' : '500px' }"
    >
      <n-form ref="formRef" :model="formData" :rules="rules">
        <n-form-item label="名称" path="name">
          <n-input v-model:value="formData.name" placeholder="请输入分类名称" />
        </n-form-item>

        <n-form-item label="描述">
          <n-input
            v-model:value="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入分类描述"
          />
        </n-form-item>

        <n-form-item label="颜色">
          <n-color-picker v-model:value="formData.color" :swatches="colorSwatches" />
        </n-form-item>

        <n-form-item label="排序">
          <n-input-number v-model:value="formData.sort" :min="0" style="width: 100%" />
        </n-form-item>
      </n-form>

      <template #footer>
        <n-space justify="end">
          <n-button @click="showModal = false">取消</n-button>
          <n-button type="primary" :loading="submitting" @click="handleSubmit">
            保存
          </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted, h } from 'vue'
import { useMessage, useDialog, NButton, NTag, NSpace } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { AddOutline } from '@vicons/ionicons5'
import { getCategories, createCategory, updateCategory, deleteCategory } from '@/api/category'
import type { Category, CategoryForm } from '@/types/blog'
import { DEFAULT_COLORS } from '@/utils/constants'

const message = useMessage()
const dialog = useDialog()

const loading = ref(false)
const submitting = ref(false)
const showModal = ref(false)
const categories = ref<Category[]>([])
const editingId = ref<number | null>(null)
const isMobile = ref(false)

// 检测移动设备
function checkMobile() {
  isMobile.value = window.innerWidth <= 1100
}

const formData = reactive<CategoryForm>({
  name: '',
  description: '',
  color: '#2196F3',
  sort: 0
})

const colorSwatches = DEFAULT_COLORS

const rules = {
  name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }]
}

const columns: DataTableColumns<Category> = [
  { 
    title: 'ID', 
    key: 'id', 
    width: 60,
    render: (_row, index) => index + 1
  },
  { title: '名称', key: 'name' },
  {
    title: '颜色',
    key: 'color',
    width: 100,
    render: row =>
      h(NTag, { color: { color: row.color, textColor: '#fff' } }, { default: () => row.color })
  },
  { title: '文章数', key: 'post_count', width: 100 },
  { title: '排序', key: 'sort', width: 80 },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render: row =>
      h(NSpace, null, {
        default: () => [
          h(
            NButton,
            { size: 'small', onClick: () => handleEdit(row) },
            { default: () => '编辑' }
          ),
          h(
            NButton,
            { size: 'small', type: 'error', onClick: () => handleDelete(row.id) },
            { default: () => '删除' }
          )
        ]
      })
  }
]

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  fetchCategories()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

async function fetchCategories() {
  try {
    loading.value = true
    const res = await getCategories()
    if (res.data) {
      categories.value = res.data
    }
  } catch (error: any) {
    message.error(error.message || '获取分类列表失败')
  } finally {
    loading.value = false
  }
}

function handleEdit(category: Category) {
  editingId.value = category.id
  formData.name = category.name
  formData.description = category.description
  formData.color = category.color
  formData.sort = category.sort
  showModal.value = true
}

async function handleSubmit() {
  try {
    submitting.value = true
    if (editingId.value) {
      await updateCategory(editingId.value, formData)
      message.success('更新成功')
    } else {
      await createCategory(formData)
      message.success('创建成功')
    }
    showModal.value = false
    resetForm()
    fetchCategories()
  } catch (error: any) {
    message.error(error.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

function handleDelete(id: number) {
  const category = categories.value.find(c => c.id === id)
  const categoryName = category?.name || '该分类'
  const postCount = category?.post_count || 0
  
  dialog.warning({
    title: '确认删除',
    content: postCount > 0 
      ? `确定要删除分类"${categoryName}"吗？该分类下还有 ${postCount} 篇文章！` 
      : `确定要删除分类"${categoryName}"吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteCategory(id)
        message.success('删除成功')
        fetchCategories()
      } catch (error: any) {
        message.error(error.message || '删除失败')
      }
    }
  })
}

function resetForm() {
  editingId.value = null
  formData.name = ''
  formData.description = ''
  formData.color = '#2196F3'
  formData.sort = 0
}
</script>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  gap: 12px;
}

.header h1 {
  margin: 0;
  font-size: 24px;
}

/* 移动端样式 (断点调整为 1100px) */
@media (max-width: 1100px) {
  .header h1 {
    font-size: 20px;
  }
  
  .category-manage-page :deep(.n-data-table) {
    font-size: 13px;
  }
}

/* 卡片列表样式 */
.card-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 8px 0;
}

.list-card {
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.card-content {
  padding: 8px 0;
}

.info-item {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
  font-size: 13px;
  line-height: 1.5;
}

.info-item .label {
  color: #666;
  width: 60px;
  flex-shrink: 0;
}

.info-item .value {
  color: #333;
  word-break: break-all;
}
</style>

