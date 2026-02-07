<!--
 * @ProjectName: go-vue3-blog
 * @FileName: OperationLogManage.vue
 * @CreateTime: 2026-02-06 22:00:00
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 操作日志管理页面组件，提供操作日志的查询和查看功能（仅超级管理员）
 -->
<template>
  <div class="operation-log-manage-page">
    <n-card title="操作日志管理">
      <n-spin :show="loading">
        <!-- 筛选条件 -->
        <n-form
          :model="filterForm"
          inline
          label-placement="left"
          :label-width="80"
          style="margin-bottom: 20px"
        >
          <n-form-item label="操作模块">
            <n-select
              v-model:value="filterForm.module"
              placeholder="请选择模块"
              clearable
              :options="moduleOptions"
              style="width: 150px"
            />
          </n-form-item>
          <n-form-item label="操作类型">
            <n-select
              v-model:value="filterForm.action"
              placeholder="请选择操作类型"
              clearable
              :options="actionOptions"
              style="width: 150px"
            />
          </n-form-item>
          <n-form-item label="用户名">
            <n-input
              v-model:value="filterForm.username"
              placeholder="请输入用户名"
              clearable
              style="width: 150px"
            />
          </n-form-item>
          <n-form-item>
            <n-button type="primary" @click="handleSearch">查询</n-button>
            <n-button style="margin-left: 8px" @click="handleReset">重置</n-button>
          </n-form-item>
        </n-form>

        <!-- 批量操作工具栏 -->
        <n-card v-if="selectedRowKeys.length > 0" style="margin-bottom: 16px" size="small">
          <n-space align="center" justify="space-between">
            <n-text strong>已选择 {{ selectedRowKeys.length }} 条记录</n-text>
            <n-space>
              <n-button type="error" @click="handleBatchDelete">
                批量删除
              </n-button>
              <n-button @click="selectedRowKeys = []">
                取消选择
              </n-button>
            </n-space>
          </n-space>
        </n-card>

        <!-- 数据表格 -->
        <n-data-table
          :columns="columns"
          :data="logs"
          :loading="loading"
          :scroll-x="isMobile ? 1200 : undefined"
          :single-line="false"
          :row-key="(row: OperationLog) => row.id"
          v-model:checked-row-keys="selectedRowKeys"
          @update:checked-row-keys="handleCheckedRowKeysChange"
        />

        <!-- 分页 -->
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
      </n-spin>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, h } from 'vue'
import { useMessage, useDialog, NButton, NTag, NSpace, NCard, NForm, NFormItem, NSelect, NInput, NPagination, NSpin, NText } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { getOperationLogs, deleteOperationLog, batchDeleteOperationLogs } from '@/api/operationLog'
import type { OperationLog, OperationLogParams } from '@/api/operationLog'
import { formatDate } from '@/utils/format'

const message = useMessage()
const dialog = useDialog()

const loading = ref(false)
const logs = ref<OperationLog[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = 15
const isMobile = ref(false)
const selectedRowKeys = ref<number[]>([])

const filterForm = ref<OperationLogParams>({
  module: undefined,
  action: undefined,
  username: undefined
})

// 模块选项
const moduleOptions = [
  { label: '文章', value: 'post' },
  { label: '分类', value: 'category' },
  { label: '标签', value: 'tag' },
  { label: '用户', value: 'user' },
  { label: '评论', value: 'comment' },
  { label: '说说', value: 'moment' },
  { label: '聊天室', value: 'chat' }
]

// 操作类型选项
const actionOptions = [
  { label: '创建', value: 'create' },
  { label: '更新', value: 'update' },
  { label: '删除', value: 'delete' }
]

// 计算总页数
const totalPages = computed(() => Math.ceil(total.value / pageSize))

// 表格列定义
const columns: DataTableColumns<OperationLog> = [
  {
    type: 'selection'
  },
  {
    title: 'ID',
    key: 'id',
    width: 60,
    render: (_row, index) => {
      return (currentPage.value - 1) * pageSize + index + 1
    }
  },
  {
    title: '操作用户',
    key: 'username',
    width: 120
  },
  {
    title: '操作模块',
    key: 'module',
    width: 100,
    render: row => {
      const moduleMap: Record<string, string> = {
        post: '文章',
        category: '分类',
        tag: '标签',
        user: '用户',
        comment: '评论',
        moment: '说说',
        chat: '聊天室'
      }
      return h(
        NTag,
        { type: 'info', size: 'small' },
        { default: () => moduleMap[row.module] || row.module }
      )
    }
  },
  {
    title: '操作类型',
    key: 'action',
    width: 100,
    render: row => {
      const actionMap: Record<string, { text: string; type: 'success' | 'warning' | 'error' }> = {
        create: { text: '创建', type: 'success' },
        update: { text: '更新', type: 'warning' },
        delete: { text: '删除', type: 'error' }
      }
      const actionInfo = actionMap[row.action] || { text: row.action, type: 'info' }
      return h(
        NTag,
        { type: actionInfo.type, size: 'small' },
        { default: () => actionInfo.text }
      )
    }
  },
  {
    title: '目标名称',
    key: 'target_name',
    width: 200,
    ellipsis: { tooltip: true }
  },
  {
    title: '操作描述',
    key: 'description',
    width: 300,
    ellipsis: { tooltip: true }
  },
  {
    title: 'IP地址',
    key: 'ip',
    width: 140
  },
  {
    title: '操作时间',
    key: 'created_at',
    width: 160,
    render: row => formatDate(row.created_at, 'YYYY-MM-DD HH:mm:ss')
  },
  {
    title: '操作',
    key: 'actions',
    width: 100,
    fixed: 'right',
    render: row => {
      return h(
        NButton,
        {
          size: 'small',
          type: 'error',
          onClick: () => handleDelete(row)
        },
        { default: () => '删除' }
      )
    }
  }
]

// 检测移动设备
function checkMobile() {
  isMobile.value = window.innerWidth <= 768
}

// 获取操作日志列表
async function fetchLogs() {
  try {
    loading.value = true
    const params: OperationLogParams = {
      page: currentPage.value,
      page_size: pageSize
    }
    
    if (filterForm.value.module) {
      params.module = filterForm.value.module
    }
    if (filterForm.value.action) {
      params.action = filterForm.value.action
    }
    if (filterForm.value.username) {
      params.username = filterForm.value.username
    }

    const res = await getOperationLogs(params)

    if (res.data) {
      logs.value = res.data.list
      total.value = res.data.total
      // 确保页码不超过最大页数
      const maxPage = Math.ceil(total.value / pageSize) || 1
      if (currentPage.value > maxPage && maxPage > 0) {
        currentPage.value = maxPage
      }
    }
  } catch (error: any) {
    message.error(error.message || '获取操作日志列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
function handleSearch() {
  currentPage.value = 1
  fetchLogs()
}

// 重置
function handleReset() {
  filterForm.value = {
    module: undefined,
    action: undefined,
    username: undefined
  }
  currentPage.value = 1
  fetchLogs()
}

// 分页变化
function handlePageChange(page: number) {
  currentPage.value = page
  selectedRowKeys.value = [] // 切换页面时清空选择
  fetchLogs()
}

// 处理选择变化
function handleCheckedRowKeysChange(keys: Array<string | number>) {
  selectedRowKeys.value = keys as number[]
}

// 删除单个操作日志
function handleDelete(log: OperationLog) {
  dialog.error({
    title: '确认删除',
    content: `确定要删除这条操作日志吗？此操作不可恢复！`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteOperationLog(log.id)
        message.success('删除成功')
        selectedRowKeys.value = []
        fetchLogs()
      } catch (error: any) {
        message.error(error.message || '删除失败')
      }
    }
  })
}

// 批量删除操作日志
function handleBatchDelete() {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择要删除的日志')
    return
  }

  dialog.error({
    title: '确认批量删除',
    content: `确定要删除选中的 ${selectedRowKeys.value.length} 条操作日志吗？此操作不可恢复！`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await batchDeleteOperationLogs(selectedRowKeys.value)
        message.success('批量删除成功')
        selectedRowKeys.value = []
        fetchLogs()
      } catch (error: any) {
        message.error(error.message || '批量删除失败')
      }
    }
  })
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  fetchLogs()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped>
.operation-log-manage-page {
  padding: 20px;
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 768px) {
  .operation-log-manage-page {
    padding: 12px;
  }
}
</style>
