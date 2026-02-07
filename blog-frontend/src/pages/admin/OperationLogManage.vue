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
          :label-placement="isMobile ? 'top' : 'left'"
          :label-width="isMobile ? undefined : 80"
          :show-feedback="false"
          class="filter-form"
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
        <div v-if="isMobile" class="card-list">
          <n-card v-for="log in logs" :key="log.id" class="list-card" size="small">
            <template #header>
              <div class="card-header-content">
                <div class="header-left">
                  <n-checkbox 
                    :checked="selectedRowKeys.includes(log.id)"
                    @update:checked="(checked) => handleCardSelect(log.id, checked)"
                  />
                  <span class="user-name">{{ log.username }}</span>
                </div>
                <n-space size="small">
                  <n-tag :type="getModuleType(log.module)" size="tiny">
                    {{ getModuleLabel(log.module) }}
                  </n-tag>
                  <n-tag :type="getActionType(log.action)" size="tiny">
                    {{ getActionLabel(log.action) }}
                  </n-tag>
                </n-space>
              </div>
            </template>
            <div class="card-content">
              <div class="info-item">
                <span class="label">目标：</span>
                <span class="value">{{ log.target_name || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="label">描述：</span>
                <span class="value">{{ log.description || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="label">IP：</span>
                <span class="value">{{ log.ip }}</span>
              </div>
              <div class="info-item">
                <span class="label">时间：</span>
                <span class="value">{{ formatDate(log.created_at, 'YYYY-MM-DD HH:mm:ss') }}</span>
              </div>
            </div>
            <template #footer>
              <n-space justify="end">
                <n-button size="tiny" type="error" @click="handleDelete(log)">
                  删除
                </n-button>
              </n-space>
            </template>
          </n-card>
        </div>

        <n-data-table
          v-else
          :columns="columns"
          :data="logs"
          :loading="loading"
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
            :page-slot="isMobile ? 3 : 7"
            :simple="isMobile"
            @update:page="handlePageChange"
          />
        </div>
      </n-spin>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, h } from 'vue'
import { useMessage, useDialog, NButton, NTag, NSpace, NCard, NForm, NFormItem, NSelect, NInput, NPagination, NSpin, NText, NCheckbox } from 'naive-ui'
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

// 辅助函数：获取模块显示文本
function getModuleLabel(module: string) {
  const moduleMap: Record<string, string> = {
    post: '文章',
    category: '分类',
    tag: '标签',
    user: '用户',
    comment: '评论',
    moment: '说说',
    chat: '聊天室'
  }
  return moduleMap[module] || module
}

// 辅助函数：获取模块标签类型
function getModuleType(_module: string): 'info' | 'success' | 'warning' | 'error' | 'default' {
  return 'info'
}

// 辅助函数：获取操作类型显示文本
function getActionLabel(action: string) {
  const actionMap: Record<string, string> = {
    create: '创建',
    update: '更新',
    delete: '删除'
  }
  return actionMap[action] || action
}

// 辅助函数：获取操作类型标签类型
function getActionType(action: string): 'success' | 'warning' | 'error' | 'info' {
  const actionMap: Record<string, 'success' | 'warning' | 'error'> = {
    create: 'success',
    update: 'warning',
    delete: 'error'
  }
  return actionMap[action] || 'info'
}

// 移动端卡片选择处理
function handleCardSelect(id: number, checked: boolean) {
  if (checked) {
    if (!selectedRowKeys.value.includes(id)) {
      selectedRowKeys.value.push(id)
    }
  } else {
    selectedRowKeys.value = selectedRowKeys.value.filter(k => k !== id)
  }
}

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
      return h(
        NTag,
        { type: getModuleType(row.module), size: 'small' },
        { default: () => getModuleLabel(row.module) }
      )
    }
  },
  {
    title: '操作类型',
    key: 'action',
    width: 100,
    render: row => {
      return h(
        NTag,
        { type: getActionType(row.action), size: 'small' },
        { default: () => getActionLabel(row.action) }
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
  isMobile.value = window.innerWidth <= 1100
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

/* 移动端样式 (断点调整为 1100px) */
@media (max-width: 1100px) {
  .operation-log-manage-page {
    padding: 12px;
  }
  
  .pagination-wrapper {
    justify-content: center;
  }

  .filter-form {
    margin-bottom: 16px;
  }

  /* 让筛选表单在移动端更好的排列 */
  :deep(.n-form.n-form--inline) {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }
  
  :deep(.n-form-item) {
    margin-right: 0 !important;
    width: 100%;
  }
  
  :deep(.n-form-item-blank) {
    width: 100% !important;
  }
  
  :deep(.n-select), :deep(.n-input) {
    width: 100% !important;
  }

  :deep(.n-form-item .n-form-item-label) {
    padding-bottom: 4px;
    font-weight: 500;
  }
}

/* 卡片列表样式 */
.card-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 8px 0;
}

.list-card {
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.card-header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.user-name {
  font-weight: 500;
  font-size: 14px;
}

.card-content {
  padding: 4px 0;
}

.info-item {
  display: flex;
  align-items: flex-start;
  margin-bottom: 6px;
  font-size: 12px;
  line-height: 1.4;
}

.info-item .label {
  color: #888;
  width: 45px;
  flex-shrink: 0;
}

.info-item .value {
  color: #555;
  flex: 1;
  word-break: break-all;
}
</style>
