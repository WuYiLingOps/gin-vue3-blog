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

        <n-form-item label="背景颜色">
          <n-color-picker v-model:value="formData.color" :swatches="colorSwatches" :modes="['hex']" />
        </n-form-item>

        <n-form-item label="文字颜色">
          <n-color-picker v-model:value="formData.text_color" :swatches="colorSwatches" :modes="['hex']" />
          <template #feedback>
            <span style="font-size: 12px; color: #999;">不设置则默认为白色</span>
          </template>
        </n-form-item>

        <n-form-item label="文字大小">
          <n-input-number v-model:value="formData.font_size" :min="12" :max="32" placeholder="16" style="width: 100%">
            <template #suffix>px</template>
          </n-input-number>
          <template #feedback>
            <span style="font-size: 12px; color: #999;">不设置则根据文章数量自动调整（16-20px）</span>
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
import { ref, reactive, onMounted, h } from 'vue'
import { useMessage, useDialog, NButton, NTag, NSpace } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { AddOutline } from '@vicons/ionicons5'
import { getTags, createTag, updateTag, deleteTag } from '@/api/tag'
import type { Tag, TagForm } from '@/types/blog'
import { DEFAULT_COLORS } from '@/utils/constants'

const message = useMessage()
const dialog = useDialog()

const loading = ref(false)
const submitting = ref(false)
const showModal = ref(false)
const tags = ref<Tag[]>([])
const editingId = ref<number | null>(null)

const formData = reactive<TagForm>({
  name: '',
  color: '#2196F3',
  text_color: undefined,
  font_size: undefined
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
  formData.text_color = tag.text_color
  formData.font_size = tag.font_size
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

function handleDelete(id: number) {
  const tag = tags.value.find(t => t.id === id)
  const tagName = tag?.name || '该标签'
  const postCount = tag?.post_count || 0
  
  dialog.warning({
    title: '确认删除',
    content: postCount > 0 
      ? `确定要删除标签"${tagName}"吗？该标签下还有 ${postCount} 篇文章！` 
      : `确定要删除标签"${tagName}"吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteTag(id)
        message.success('删除成功')
        fetchTags()
      } catch (error: any) {
        message.error(error.message || '删除失败')
      }
    }
  })
}

function resetForm() {
  editingId.value = null
  formData.name = ''
  formData.color = '#2196F3'
  formData.text_color = undefined
  formData.font_size = undefined
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

