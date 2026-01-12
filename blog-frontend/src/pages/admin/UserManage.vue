<template>
  <div class="user-manage-page">
    <h1 class="page-title">用户管理</h1>

    <!-- 注册开关 -->
    <n-card class="register-settings-card" style="margin-bottom: 16px">
      <n-space align="center" justify="space-between">
        <div>
          <n-text strong>限制用户注册</n-text>
          <n-text depth="3" style="margin-left: 8px; font-size: 13px">
            开启后将禁止新用户注册
          </n-text>
        </div>
        <n-switch
          v-model:value="registerDisabled"
          :loading="settingLoading"
          :disabled="settingLoading"
          @update:value="handleToggleRegister"
        />
      </n-space>
    </n-card>

    <!-- 内容区域 -->
    <div class="content-area">
      <n-data-table
        :columns="columns"
        :data="users"
        :loading="loading"
        :scroll-x="isMobile ? 900 : undefined"
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
import { useMessage, useDialog, NButton, NTag, NSpace, NAvatar, NCard, NText, NSwitch } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { getUsers, updateUserStatus, deleteUser, getRegisterSettings, updateRegisterSettings } from '@/api/user'
import { formatDate } from '@/utils/format'
import type { User } from '@/types/auth'

const message = useMessage()
const dialog = useDialog()

const loading = ref(false)
const users = ref<User[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = 15 // 固定每页显示15个用户
const isMobile = ref(false)
const registerDisabled = ref(false)
const settingLoading = ref(false)

// 检测移动设备
function checkMobile() {
  isMobile.value = window.innerWidth <= 768
}

// 计算总页数
const totalPages = computed(() => Math.ceil(total.value / pageSize))

const columns: DataTableColumns<User> = [
  { 
    title: 'ID', 
    key: 'id', 
    width: 60,
    render: (_row, index) => {
      return (currentPage.value - 1) * pageSize + index + 1
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
  // 异步获取注册配置，不阻塞页面渲染
  fetchRegisterSettings().catch(() => {
    // 静默处理错误
  })
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

async function fetchUsers() {
  try {
    loading.value = true
    const res = await getUsers({
      page: currentPage.value,
      page_size: pageSize
    })

    if (res.data) {
      users.value = res.data.list
      total.value = res.data.total
      // 确保页码不超过最大页数
      const maxPage = Math.ceil(total.value / pageSize) || 1
      if (currentPage.value > maxPage && maxPage > 0) {
        currentPage.value = maxPage
      }
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

async function fetchRegisterSettings() {
  try {
    settingLoading.value = true
    const res = await getRegisterSettings()
    // 安全地获取配置值
    if (res?.data?.disable_register !== undefined) {
      const value = String(res.data.disable_register)
      registerDisabled.value = value === '1' || value === 'true'
    } else {
      // 如果数据格式不对，使用默认值
      registerDisabled.value = false
    }
  } catch (error: any) {
    // 静默处理错误，不影响页面正常显示
    // 如果配置不存在或获取失败，使用默认值（允许注册）
    console.warn('获取注册配置失败，使用默认值（允许注册）:', error?.message || error)
    registerDisabled.value = false
  } finally {
    settingLoading.value = false
  }
}

async function handleToggleRegister(value: boolean) {
  const oldValue = !value // 保存旧值
  try {
    settingLoading.value = true
    await updateRegisterSettings({ disable_register: value ? '1' : '0' })
    message.success(value ? '已关闭用户注册' : '已开启用户注册')
  } catch (error: any) {
    // 恢复原值
    registerDisabled.value = oldValue
    message.error(error.message || '更新配置失败')
  } finally {
    settingLoading.value = false
  }
}
</script>

<style scoped>
.user-manage-page {
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
  
  .user-manage-page :deep(.n-data-table) {
    font-size: 13px;
  }
  
  .pagination-wrapper {
    bottom: 8px;
    right: 16px;
  }
}
</style>

