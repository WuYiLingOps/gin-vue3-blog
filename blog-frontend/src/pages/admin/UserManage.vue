<template>
  <div class="user-manage-page">
    <h1 class="page-title">用户管理</h1>

    <n-data-table
      :columns="columns"
      :data="users"
      :loading="loading"
      :pagination="pagination"
      :scroll-x="isMobile ? 900 : undefined"
      :single-line="false"
      @update:page="handlePageChange"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, h } from 'vue'
import { useMessage, useDialog, NButton, NTag, NSpace, NAvatar } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { getUsers, updateUserStatus, deleteUser } from '@/api/user'
import { formatDate } from '@/utils/format'
import type { User } from '@/types/auth'

const message = useMessage()
const dialog = useDialog()

const loading = ref(false)
const users = ref<User[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const isMobile = ref(false)

// 检测移动设备
function checkMobile() {
  isMobile.value = window.innerWidth <= 768
}

const pagination = computed(() => ({
  page: currentPage.value,
  pageSize: pageSize.value,
  pageCount: Math.ceil(total.value / pageSize.value),
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50]
}))

const columns: DataTableColumns<User> = [
  { 
    title: 'ID', 
    key: 'id', 
    width: 60,
    render: (_row, index) => {
      return (currentPage.value - 1) * pageSize.value + index + 1
    }
  },
  {
    title: '头像',
    key: 'avatar',
    width: 80,
    render: row => h(NAvatar, { src: row.avatar, round: true })
  },
  { title: '用户名', key: 'username' },
  { title: '昵称', key: 'nickname' },
  { title: '邮箱', key: 'email', ellipsis: { tooltip: true } },
  {
    title: '角色',
    key: 'role',
    width: 100,
    render: row =>
      h(
        NTag,
        { type: row.role === 'admin' ? 'error' : 'info', size: 'small' },
        { default: () => (row.role === 'admin' ? '管理员' : '用户') }
      )
  },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: row =>
      h(
        NTag,
        { type: row.status === 1 ? 'success' : 'default', size: 'small' },
        { default: () => (row.status === 1 ? '正常' : '禁用') }
      )
  },
  {
    title: '注册时间',
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
            { default: () => (row.status === 1 ? '禁用' : '启用') }
          ),
          h(
            NButton,
            {
              size: 'small',
              type: 'error',
              onClick: () => handleDelete(row)
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
  fetchUsers()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

async function fetchUsers() {
  try {
    loading.value = true
    const res = await getUsers({
      page: currentPage.value,
      page_size: pageSize.value
    })

    if (res.data) {
      users.value = res.data.list
      total.value = res.data.total
    }
  } catch (error: any) {
    message.error(error.message || '获取用户列表失败')
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  currentPage.value = page
  fetchUsers()
}

function handleToggleStatus(user: User) {
  const newStatus = user.status === 1 ? 0 : 1
  const action = newStatus === 0 ? '禁用' : '启用'
  const actionType = newStatus === 0 ? 'warning' : 'info'
  
  dialog[actionType]({
    title: `确认${action}`,
    content: newStatus === 0 
      ? `确定要禁用用户"${user.nickname || user.username}"吗？禁用后该用户将无法登录！` 
      : `确定要启用用户"${user.nickname || user.username}"吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await updateUserStatus(user.id, newStatus)
        message.success(`${action}成功`)
        fetchUsers()
      } catch (error: any) {
        message.error(error.message || '操作失败')
      }
    }
  })
}

function handleDelete(user: User) {
  dialog.error({
    title: '确认删除',
    content: `确定要删除用户"${user.nickname || user.username}"吗？此操作不可恢复！`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteUser(user.id)
        message.success('删除成功')
        fetchUsers()
      } catch (error: any) {
        message.error(error.message || '删除失败')
      }
    }
  })
}
</script>

<style scoped>
.page-title {
  margin-bottom: 16px;
  font-size: 24px;
}

/* 移动端样式 */
@media (max-width: 768px) {
  .page-title {
    font-size: 20px;
  }
  
  .user-manage-page :deep(.n-data-table) {
    font-size: 13px;
  }
}
</style>

