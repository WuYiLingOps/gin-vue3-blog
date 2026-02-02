<!--
 * @ProjectName: go-vue3-blog
 * @FileName: CommentManage.vue
 * @CreateTime: 2026-02-02 11:43:17
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 评论管理页面组件，提供评论的审核、删除等管理功能
 -->
<template>
  <div class="comment-manage-page">
    <h1 class="page-title">评论管理</h1>

    <!-- 内容区域 -->
    <div class="content-area">
      <n-data-table
        :columns="columns"
        :data="comments"
        :loading="loading"
        :scroll-x="isMobile ? 1000 : undefined"
        :single-line="false"
      />
      
      <!-- 分页 - 位于表格右下角 -->
      <div class="pagination-wrapper">
        <n-pagination
          v-if="total > 0"
          v-model:page="currentPage"
          :page-count="totalPages"
          :page-size="pageSize"
          :page-slot="7"
          @update:page="handlePageChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, h } from 'vue'
import { useMessage, useDialog, NButton, NTag, NSpace, NEllipsis } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { getAllComments, updateCommentStatus, deleteComment } from '@/api/comment'
import { formatDate } from '@/utils/format'
import type { Comment } from '@/types/blog'

const message = useMessage()
const dialog = useDialog()

const loading = ref(false)
const comments = ref<Comment[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = 15 // 固定每页显示15条评论
const isMobile = ref(false)

// 检测移动设备
function checkMobile() {
  isMobile.value = window.innerWidth <= 768
}

// 计算总页数
const totalPages = computed(() => Math.ceil(total.value / pageSize))

const columns: DataTableColumns<Comment> = [
  { 
    title: 'ID', 
    key: 'id', 
    width: 60,
    render: (_row, index) => {
      return (currentPage.value - 1) * pageSize + index + 1
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
  checkMobile()
  window.addEventListener('resize', checkMobile)
  fetchComments()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

async function fetchComments() {
  try {
    loading.value = true
    const res = await getAllComments({
      page: currentPage.value,
      page_size: pageSize
    })

    if (res.data) {
      comments.value = res.data.list
      total.value = res.data.total
      // 确保页码不超过最大页数
      const maxPage = Math.ceil(total.value / pageSize) || 1
      if (currentPage.value > maxPage && maxPage > 0) {
        currentPage.value = maxPage
      }
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

function handleDelete(id: number) {
  const comment = comments.value.find(c => c.id === id)
  const userName = comment?.user.nickname || '该用户'
  const contentPreview = comment?.content.substring(0, 20) || ''
  
  dialog.warning({
    title: '确认删除',
    content: `确定要删除 ${userName} 的评论吗？\n"${contentPreview}${comment && comment.content.length > 20 ? '...' : ''}"\n删除后无法恢复！`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteComment(id)
        message.success('删除成功')
        fetchComments()
      } catch (error: any) {
        message.error(error.message || '删除失败')
      }
    }
  })
}
</script>

<style scoped>
.comment-manage-page {
  position: relative;
}

.page-title {
  margin-bottom: 16px;
  font-size: 24px;
}

/* 内容区域 */
.content-area {
  position: relative;
  padding-bottom: 50px; /* 为分页器预留空间 */
}

/* 分页样式 - 位于表格右下角 */
.pagination-wrapper {
  position: absolute;
  bottom: 10px; /* 距离表格底部 */
  right: 20px; /* 距离右侧 */
  z-index: 10;
  display: flex;
  justify-content: flex-end;
  align-items: center;
  box-sizing: border-box;
}

/* 移动端样式 */
@media (max-width: 768px) {
  .page-title {
    font-size: 20px;
  }
  
  .comment-manage-page :deep(.n-data-table) {
    font-size: 13px;
  }
  
  .pagination-wrapper {
    bottom: 16px;
    right: 16px;
  }
}
</style>

