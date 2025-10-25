<template>
  <div class="comment-manage-page">
    <h1 style="margin-bottom: 16px">评论管理</h1>

    <n-data-table
      :columns="columns"
      :data="comments"
      :loading="loading"
      :pagination="pagination"
      @update:page="handlePageChange"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, h } from 'vue'
import { useMessage, NButton, NTag, NSpace, NEllipsis } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { getAllComments, updateCommentStatus, deleteComment } from '@/api/comment'
import { formatDate } from '@/utils/format'
import type { Comment } from '@/types/blog'

const message = useMessage()

const loading = ref(false)
const comments = ref<Comment[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const pagination = computed(() => ({
  page: currentPage.value,
  pageSize: pageSize.value,
  pageCount: Math.ceil(total.value / pageSize.value),
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50]
}))

const columns: DataTableColumns<Comment> = [
  { 
    title: 'ID', 
    key: 'id', 
    width: 60,
    render: (row, index) => {
      return (currentPage.value - 1) * pageSize.value + index + 1
    }
  },
  {
    title: '内容',
    key: 'content',
    ellipsis: {
      tooltip: true
    },
    render: row => h(NEllipsis, { style: 'max-width: 300px' }, { default: () => row.content })
  },
  {
    title: '用户',
    key: 'user',
    width: 120,
    render: row => row.user.nickname
  },
  {
    title: '文章',
    key: 'post',
    width: 150,
    ellipsis: { tooltip: true },
    render: row => row.post?.title || '-'
  },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: row =>
      h(
        NTag,
        { type: row.status === 1 ? 'success' : 'default', size: 'small' },
        { default: () => (row.status === 1 ? '正常' : '隐藏') }
      )
  },
  {
    title: '时间',
    key: 'created_at',
    width: 160,
    render: row => formatDate(row.created_at, 'YYYY-MM-DD HH:mm')
  },
  {
    title: '操作',
    key: 'actions',
    width: 180,
    render: row =>
      h(NSpace, null, {
        default: () => [
          h(
            NButton,
            {
              size: 'small',
              onClick: () => handleToggleStatus(row)
            },
            { default: () => (row.status === 1 ? '隐藏' : '显示') }
          ),
          h(
            NButton,
            {
              size: 'small',
              type: 'error',
              onClick: () => handleDelete(row.id)
            },
            { default: () => '删除' }
          )
        ]
      })
  }
]

onMounted(() => {
  fetchComments()
})

async function fetchComments() {
  try {
    loading.value = true
    const res = await getAllComments({
      page: currentPage.value,
      page_size: pageSize.value
    })

    if (res.data) {
      comments.value = res.data.list
      total.value = res.data.total
    }
  } catch (error: any) {
    message.error(error.message || '获取评论列表失败')
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  currentPage.value = page
  fetchComments()
}

async function handleToggleStatus(comment: Comment) {
  try {
    const newStatus = comment.status === 1 ? 0 : 1
    await updateCommentStatus(comment.id, newStatus)
    message.success('状态更新成功')
    fetchComments()
  } catch (error: any) {
    message.error(error.message || '操作失败')
  }
}

async function handleDelete(id: number) {
  try {
    await deleteComment(id)
    message.success('删除成功')
    fetchComments()
  } catch (error: any) {
    message.error(error.message || '删除失败')
  }
}
</script>

