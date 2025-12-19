<template>
  <div class="category-manage-page">
    <div class="header">
      <h1>友链分类管理</h1>
      <n-button type="primary" :size="isMobile ? 'small' : 'medium'" @click="handleCreate">
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
      :scroll-x="isMobile ? 800 : undefined"
      :single-line="false"
    />

    <!-- 创建/编辑对话框 -->
    <n-modal 
      v-model:show="showModal" 
      preset="card" 
      :title="editingId ? '编辑分类' : '新建分类'" 
      :style="{ width: isMobile ? '95%' : '600px', maxWidth: isMobile ? '95vw' : '600px' }"
    >
      <n-form ref="formRef" :model="formData" :rules="rules">
        <n-form-item label="分类名称" path="name">
          <n-input v-model:value="formData.name" placeholder="例如：推荐" />
        </n-form-item>

        <n-form-item label="分类描述">
          <n-input
            v-model:value="formData.description"
            type="textarea"
            :rows="2"
            placeholder="分类描述，例如：都是大佬,推荐关注"
          />
        </n-form-item>

        <n-form-item label="排序">
          <n-input-number v-model:value="formData.sort_order" :min="0" style="width: 100%" />
          <template #feedback>
            <span style="color: #909399; font-size: 12px;">数字越大越靠前</span>
          </template>
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
import { useMessage, useDialog, NButton, NSpace } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { AddOutline } from '@vicons/ionicons5'
import { 
  getFriendLinkCategoriesAdmin, 
  createFriendLinkCategory, 
  updateFriendLinkCategory, 
  deleteFriendLinkCategory 
} from '@/api/friendlink'
import type { FriendLinkCategory, FriendLinkCategoryForm } from '@/api/friendlink'

const message = useMessage()
const dialog = useDialog()

const loading = ref(false)
const submitting = ref(false)
const showModal = ref(false)
const categories = ref<FriendLinkCategory[]>([])
const editingId = ref<number | null>(null)
const isMobile = ref(false)

// 检测移动设备
function checkMobile() {
  isMobile.value = window.innerWidth <= 768
}

const formData = reactive<FriendLinkCategoryForm>({
  name: '',
  description: '',
  sort_order: 0
})

const rules = {
  name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }]
}

const columns: DataTableColumns<FriendLinkCategory> = [
  { 
    title: 'ID', 
    key: 'id', 
    width: 60
  },
  {
    title: '分类名称',
    key: 'name',
    width: 150
  },
  {
    title: '描述',
    key: 'description',
    width: 250,
    ellipsis: { tooltip: true }
  },
  {
    title: '排序',
    key: 'sort_order',
    width: 100
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    fixed: 'right',
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
    const res = await getFriendLinkCategoriesAdmin()
    if (res && res.data) {
      categories.value = Array.isArray(res.data) ? res.data : []
    }
  } catch (error: any) {
    message.error(error.message || '获取分类列表失败')
  } finally {
    loading.value = false
  }
}

function resetForm() {
  formData.name = ''
  formData.description = ''
  formData.sort_order = 0
}

function handleCreate() {
  editingId.value = null
  resetForm()
  showModal.value = true
}

function handleEdit(category: FriendLinkCategory) {
  editingId.value = category.id
  formData.name = category.name
  formData.description = category.description || ''
  formData.sort_order = category.sort_order
  showModal.value = true
}

async function handleSubmit() {
  try {
    await formRef.value?.validate()
    submitting.value = true

    if (editingId.value) {
      await updateFriendLinkCategory(editingId.value, formData)
      message.success('分类更新成功')
    } else {
      await createFriendLinkCategory(formData)
      message.success('分类创建成功')
    }

    showModal.value = false
    fetchCategories()
  } catch (error: any) {
    if (error.message && !error.message.includes('验证')) {
      message.error(error.message || '操作失败')
    }
  } finally {
    submitting.value = false
  }
}

function handleDelete(id: number) {
  const category = categories.value.find(c => c.id === id)
  const categoryName = category?.name || '该分类'
  
  dialog.warning({
    title: '确认删除',
    content: `确定要删除分类"${categoryName}"吗？删除后该分类下的友链需要重新分配分类。`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteFriendLinkCategory(id)
        message.success('删除成功')
        fetchCategories()
      } catch (error: any) {
        message.error(error.message || '删除失败')
      }
    }
  })
}

const formRef = ref()
</script>

<style scoped>
.category-manage-page {
  padding: 24px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.header h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
}

@media (max-width: 768px) {
  .category-manage-page {
    padding: 16px;
  }

  .header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .header h1 {
    font-size: 20px;
  }
}
</style>

