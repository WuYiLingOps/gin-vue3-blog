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

    <n-data-table 
      :columns="columns" 
      :data="categories" 
      :loading="loading"
      :scroll-x="isMobile ? 600 : undefined"
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
  isMobile.value = window.innerWidth <= 768
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

/* 移动端样式 */
@media (max-width: 768px) {
  .header h1 {
    font-size: 20px;
  }
  
  .category-manage-page :deep(.n-data-table) {
    font-size: 13px;
  }
}
</style>

