<template>
  <div class="tag-manage-page">
    <div class="header">
      <h1>标签管理</h1>
      <n-button type="primary" @click="showModal = true">
        <template #icon>
          <n-icon :component="AddOutline" />
        </template>
        新建标签
      </n-button>
    </div>

    <n-data-table :columns="columns" :data="tags" :loading="loading" />

    <!-- 创建/编辑对话框 -->
    <n-modal v-model:show="showModal" preset="card" :title="editingId ? '编辑标签' : '新建标签'" style="width: 500px">
      <n-form ref="formRef" :model="formData" :rules="rules">
        <n-form-item label="名称" path="name">
          <n-input v-model:value="formData.name" placeholder="请输入标签名称" />
        </n-form-item>

        <n-form-item label="颜色">
          <n-color-picker v-model:value="formData.color" :swatches="colorSwatches" />
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
import { ref, reactive, onMounted, h } from 'vue'
import { useMessage, NButton, NTag, NSpace } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { AddOutline } from '@vicons/ionicons5'
import { getTags, createTag, updateTag, deleteTag } from '@/api/tag'
import type { Tag, TagForm } from '@/types/blog'
import { DEFAULT_COLORS } from '@/utils/constants'

const message = useMessage()

const loading = ref(false)
const submitting = ref(false)
const showModal = ref(false)
const tags = ref<Tag[]>([])
const editingId = ref<number | null>(null)

const formData = reactive<TagForm>({
  name: '',
  color: '#2196F3'
})

const colorSwatches = DEFAULT_COLORS

const rules = {
  name: [{ required: true, message: '请输入标签名称', trigger: 'blur' }]
}

const columns: DataTableColumns<Tag> = [
  { 
    title: 'ID', 
    key: 'id', 
    width: 60,
    render: (row, index) => index + 1
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
  fetchTags()
})

async function fetchTags() {
  try {
    loading.value = true
    const res = await getTags()
    if (res.data) {
      tags.value = res.data
    }
  } catch (error: any) {
    message.error(error.message || '获取标签列表失败')
  } finally {
    loading.value = false
  }
}

function handleEdit(tag: Tag) {
  editingId.value = tag.id
  formData.name = tag.name
  formData.color = tag.color
  showModal.value = true
}

async function handleSubmit() {
  try {
    submitting.value = true
    if (editingId.value) {
      await updateTag(editingId.value, formData)
      message.success('更新成功')
    } else {
      await createTag(formData)
      message.success('创建成功')
    }
    showModal.value = false
    resetForm()
    fetchTags()
  } catch (error: any) {
    message.error(error.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

async function handleDelete(id: number) {
  try {
    await deleteTag(id)
    message.success('删除成功')
    fetchTags()
  } catch (error: any) {
    message.error(error.message || '删除失败')
  }
}

function resetForm() {
  editingId.value = null
  formData.name = ''
  formData.color = '#2196F3'
}
</script>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.header h1 {
  margin: 0;
}
</style>

