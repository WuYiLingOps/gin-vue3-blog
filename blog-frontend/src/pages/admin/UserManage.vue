<template>
  <div class="user-manage-page">
    <h1 style="margin-bottom: 16px">用户管理</h1>

    <n-data-table
      :columns="columns"
      :data="users"
      :loading="loading"
      :pagination="pagination"
      @update:page="handlePageChange"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, h } from 'vue'
import { useMessage, NButton, NTag, NSpace, NAvatar } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { getUsers, updateUserStatus } from '@/api/user'
import { formatDate } from '@/utils/format'
import { normalizeImageUrl } from '@/utils/url'
import type { User } from '@/types/auth'

const message = useMessage()

const loading = ref(false)
const users = ref<User[]>([])
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
    width: 120,
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
          )
        ]
      })
  }
]

onMounted(() => {
  fetchUsers()
})

async function fetchUsers() {
  try {
    loading.value = true
    const res = await getUsers({
      page: currentPage.value,
      page_size: pageSize.value
    })

    if (res.data) {
      // 标准化用户头像 URL
      users.value = res.data.list.map(user => ({
        ...user,
        avatar: normalizeImageUrl(user.avatar)
      }))
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

async function handleToggleStatus(user: User) {
  try {
    const newStatus = user.status === 1 ? 0 : 1
    await updateUserStatus(user.id, newStatus)
    message.success('状态更新成功')
    fetchUsers()
  } catch (error: any) {
    message.error(error.message || '操作失败')
  }
}
</script>

