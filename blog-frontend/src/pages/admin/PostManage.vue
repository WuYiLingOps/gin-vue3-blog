<template>
  <div class="post-manage-page">
    <div class="header">
      <h1>文章管理</h1>
      <n-button type="primary" @click="showCreateModal = true">
        <template #icon>
          <n-icon :component="AddOutline" />
        </template>
        新建文章
      </n-button>
    </div>

    <!-- 筛选 -->
    <n-space style="margin: 16px 0">
      <n-input
        v-model:value="searchKeyword"
        placeholder="搜索文章..."
        clearable
        style="width: 250px"
        @keyup.enter="fetchPosts"
      />
      <n-select
        v-model:value="filterStatus"
        placeholder="状态"
        clearable
        style="width: 120px"
        :options="statusOptions"
        @update:value="fetchPosts"
      />
      <n-button @click="fetchPosts">搜索</n-button>
    </n-space>

    <!-- 文章列表 -->
    <n-data-table
      :columns="columns"
      :data="posts"
      :loading="loading"
      :pagination="pagination"
      @update:page="handlePageChange"
    />

    <!-- 创建/编辑文章对话框 -->
    <n-modal 
      v-model:show="showCreateModal" 
      preset="card" 
      title="创建文章" 
      style="width: 800px"
      :mask-closable="false"
      :close-on-esc="false"
      @close="handleModalClose"
    >
      <n-form ref="formRef" :model="formData" :rules="rules">
        <n-form-item label="标题" path="title">
          <n-input v-model:value="formData.title" placeholder="请输入文章标题" />
        </n-form-item>

        <n-form-item label="摘要" path="summary">
          <n-input
            v-model:value="formData.summary"
            type="textarea"
            :rows="3"
            placeholder="请输入文章摘要"
          />
        </n-form-item>

        <n-form-item label="内容" path="content">
          <markdown-editor v-model="formData.content" height="400px" />
        </n-form-item>

        <n-form-item label="分类" path="category_id">
          <n-select
            v-model:value="formData.category_id"
            :options="categoryOptions"
            placeholder="选择分类"
          />
        </n-form-item>

        <n-form-item label="标签" path="tag_ids">
          <n-select
            v-model:value="formData.tag_ids"
            :options="tagOptions"
            placeholder="选择标签"
            multiple
          />
        </n-form-item>

        <n-form-item label="状态">
          <n-radio-group v-model:value="formData.status">
            <n-radio :value="0">草稿</n-radio>
            <n-radio :value="1">发布</n-radio>
          </n-radio-group>
        </n-form-item>

        <n-form-item label="置顶">
          <n-switch v-model:value="formData.is_top" />
        </n-form-item>
      </n-form>

      <template #footer>
        <n-space justify="end">
          <n-button @click="handleCancel">取消</n-button>
          <n-button type="primary" :loading="submitting" @click="handleSubmit">
            保存
          </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage, useDialog, NButton, NTag, NSpace } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { AddOutline } from '@vicons/ionicons5'
import { getPosts, createPost, deletePost } from '@/api/post'
import { useBlogStore } from '@/stores'
import { formatDate } from '@/utils/format'
import type { Post, PostForm } from '@/types/blog'
import MarkdownEditor from '@/components/MarkdownEditor.vue'

const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const blogStore = useBlogStore()

const loading = ref(false)
const submitting = ref(false)
const showCreateModal = ref(false)
const posts = ref<Post[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchKeyword = ref('')
const filterStatus = ref<number | null>(null)

const formData = reactive<PostForm>({
  title: '',
  content: '',
  summary: '',
  cover: '',
  category_id: null,
  tag_ids: [],
  status: 1,
  is_top: false
})

const statusOptions = [
  { label: '草稿', value: 0 },
  { label: '已发布', value: 1 }
]

const categoryOptions = computed(() =>
  blogStore.categories.map(c => ({ label: c.name, value: c.id }))
)

const tagOptions = computed(() => blogStore.tags.map(t => ({ label: t.name, value: t.id })))

const pagination = computed(() => ({
  page: currentPage.value,
  pageSize: pageSize.value,
  pageCount: Math.ceil(total.value / pageSize.value),
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50]
}))

const columns: DataTableColumns<Post> = [
  { 
    title: 'ID', 
    key: 'id', 
    width: 60,
    render: (_row, index) => {
      return (currentPage.value - 1) * pageSize.value + index + 1
    }
  },
  { title: '标题', key: 'title', ellipsis: { tooltip: true } },
  {
    title: '分类',
    key: 'category',
    width: 100,
    render: row => h(NTag, { type: 'info', size: 'small' }, { default: () => row.category.name })
  },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: row =>
      h(NTag, { type: row.status === 1 ? 'success' : 'default', size: 'small' }, { default: () => (row.status === 1 ? '已发布' : '草稿') })
  },
  {
    title: '浏览',
    key: 'view_count',
    width: 80
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 160,
    render: row => formatDate(row.created_at, 'YYYY-MM-DD HH:mm')
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render: row =>
      h(NSpace, null, {
        default: () => [
          h(
            NButton,
            { size: 'small', onClick: () => handleEdit(row.id) },
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

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }],
  category_id: [
    { 
      required: true, 
      message: '请选择分类', 
      trigger: ['blur', 'change'],
      validator: (_rule: any, value: any) => {
        if (value === null || value === undefined || value === 0) {
          return new Error('请选择分类')
        }
        return true
      }
    }
  ]
}

onMounted(() => {
  blogStore.init()
  fetchPosts()
})

async function fetchPosts() {
  try {
    loading.value = true
    const res = await getPosts({
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: searchKeyword.value,
      status: filterStatus.value ?? undefined
    })

    if (res.data) {
      posts.value = res.data.list
      total.value = res.data.total
    }
  } catch (error: any) {
    message.error(error.message || '获取文章列表失败')
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  currentPage.value = page
  fetchPosts()
}

function handleEdit(id: number) {
  router.push({ name: 'PostEdit', params: { id } })
}

async function handleSubmit() {
  try {
    submitting.value = true
    await createPost(formData)
    message.success('文章创建成功')
    clearForm()
    showCreateModal.value = false
    fetchPosts()
  } catch (error: any) {
    message.error(error.message || '创建失败')
  } finally {
    submitting.value = false
  }
}

function handleDelete(id: number) {
  dialog.warning({
    title: '确认删除',
    content: '确定要删除这篇文章吗？删除后无法恢复！',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deletePost(id)
        message.success('删除成功')
        fetchPosts()
      } catch (error: any) {
        message.error(error.message || '删除失败')
      }
    }
  })
}

// 检查是否有未保存的内容
function hasUnsavedContent(): boolean {
  return !!(
    formData.title.trim() ||
    formData.content.trim() ||
    formData.summary.trim() ||
    formData.cover ||
    (formData.category_id !== null && formData.category_id !== 0) ||
    formData.tag_ids.length > 0
  )
}

// 保存为草稿
async function saveAsDraft() {
  try {
    // 验证必填字段
    if (!formData.title.trim()) {
      message.error('标题不能为空')
      return
    }
    if (!formData.content.trim()) {
      message.error('内容不能为空')
      return
    }
    if (!formData.category_id || formData.category_id === 0) {
      message.error('请选择分类')
      return
    }

    submitting.value = true
    
    // 设置为草稿状态
    const draftData = { ...formData, status: 0 }
    await createPost(draftData)
    message.success('已保存为草稿')
    clearForm()
    showCreateModal.value = false
    fetchPosts()
  } catch (error: any) {
    message.error(error.message || '保存草稿失败')
  } finally {
    submitting.value = false
  }
}

// 清空表单
function clearForm() {
  formData.title = ''
  formData.content = ''
  formData.summary = ''
  formData.cover = ''
  formData.category_id = null
  formData.tag_ids = []
  formData.status = 1
  formData.is_top = false
}

// 处理取消操作
function handleCancel() {
  // 如果有未保存的内容，弹出确认框
  if (hasUnsavedContent()) {
    dialog.warning({
      title: '提示',
      content: '检测到您有未保存的内容，是否要保存为草稿？',
      positiveText: '保存草稿',
      negativeText: '直接离开',
      onPositiveClick: async () => {
        await saveAsDraft()
      },
      onNegativeClick: () => {
        clearForm()
        showCreateModal.value = false
      }
    })
  } else {
    showCreateModal.value = false
  }
}

// 处理模态框关闭事件（点击 X 号）
function handleModalClose() {
  // 和取消按钮逻辑一样
  handleCancel()
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

