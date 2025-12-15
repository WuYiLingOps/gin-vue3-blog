<template>
  <div class="ip-whitelist-manage-page">
    <h1 class="page-title">IP 白名单管理</h1>

    <!-- 操作栏 -->
    <n-space style="margin-bottom: 16px" justify="space-between" align="center">
      <n-space>
        <n-button type="primary" @click="showAddModal = true">
          <template #icon>
            <n-icon>
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                <path fill="currentColor" d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
              </svg>
            </n-icon>
          </template>
          添加 IP
        </n-button>
        <n-button @click="showCheckModal = true">
          <template #icon>
            <n-icon>
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                <path fill="currentColor" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
              </svg>
            </n-icon>
          </template>
          检查 IP
        </n-button>
        <n-button type="warning" @click="handleCleanExpired">
          <template #icon>
            <n-icon>
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                <path fill="currentColor" d="M19 4h-3.5l-1-1h-5l-1 1H5v2h14M6 19a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2V7H6v12z"/>
              </svg>
            </n-icon>
          </template>
          清理过期记录
        </n-button>
      </n-space>
      <n-button @click="fetchWhitelist">
        <template #icon>
          <n-icon>
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
              <path fill="currentColor" d="M17.65 6.35A7.958 7.958 0 0 0 12 4c-4.42 0-7.99 3.58-7.99 8s3.57 8 7.99 8c3.73 0 6.84-2.55 7.73-6h-2.08A5.99 5.99 0 0 1 12 18c-3.31 0-6-2.69-6-6s2.69-6 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4l-2.35 2.35z"/>
            </svg>
          </n-icon>
        </template>
        刷新
      </n-button>
    </n-space>

    <!-- 数据表格 -->
    <n-data-table
      :columns="columns"
      :data="whitelist"
      :loading="loading"
      :pagination="pagination"
      :scroll-x="isMobile ? 1000 : undefined"
      :single-line="false"
      @update:page="handlePageChange"
      @update:page-size="handlePageSizeChange"
    />

    <!-- 添加 IP 对话框 -->
    <n-modal v-model:show="showAddModal">
      <n-card
        title="添加 IP 到白名单"
        :bordered="false"
        size="large"
        style="max-width: 600px"
        closable
        @close="handleAddModalClose"
      >
        <n-form ref="addFormRef" :model="addForm" :rules="addFormRules" label-width="100px">
          <n-form-item label="IP 地址" path="ip">
            <n-input
              v-model:value="addForm.ip"
              placeholder="请输入 IP 地址或 CIDR（如：192.168.1.100 或 192.168.1.0/24）"
              @keyup.enter="handleAdd"
            />
          </n-form-item>
          <n-form-item label="添加原因" path="reason">
            <n-input
              v-model:value="addForm.reason"
              type="textarea"
              placeholder="请输入添加原因（可选）"
              :autosize="{ minRows: 3, maxRows: 6 }"
            />
          </n-form-item>
          <n-form-item label="有效期" path="duration">
            <n-input-number
              v-model:value="addForm.duration"
              :min="0"
              :max="8760"
              placeholder="小时"
              style="width: 100%"
            >
              <template #suffix>小时</template>
            </n-input-number>
            <template #feedback>
              <span style="font-size: 12px; color: #999">
                0 表示永久有效，最大 8760 小时（1年）
              </span>
            </template>
          </n-form-item>
        </n-form>
        <template #footer>
          <n-space justify="end">
            <n-button @click="handleAddModalClose">取消</n-button>
            <n-button type="primary" :loading="adding" @click="handleAdd">
              确定添加
            </n-button>
          </n-space>
        </template>
      </n-card>
    </n-modal>

    <!-- 检查 IP 对话框 -->
    <n-modal v-model:show="showCheckModal">
      <n-card
        title="检查 IP 状态"
        :bordered="false"
        size="large"
        style="max-width: 600px"
        closable
        @close="handleCheckModalClose"
      >
        <n-form ref="checkFormRef" :model="checkForm" :rules="checkFormRules" label-width="100px">
          <n-form-item label="IP 地址" path="ip">
            <n-input
              v-model:value="checkForm.ip"
              placeholder="请输入要检查的 IP 地址"
              @keyup.enter="handleCheck"
            />
          </n-form-item>
        </n-form>
        <template #footer>
          <n-space justify="end">
            <n-button @click="handleCheckModalClose">关闭</n-button>
            <n-button type="primary" :loading="checking" @click="handleCheck">
              检查
            </n-button>
          </n-space>
        </template>
      </n-card>
    </n-modal>

    <!-- 检查结果对话框 -->
    <n-modal v-model:show="showCheckResultModal">
      <n-card
        title="IP 检查结果"
        :bordered="false"
        size="large"
        style="max-width: 600px"
        closable
        @close="showCheckResultModal = false"
      >
        <div v-if="checkResult">
          <n-alert
            :type="checkResult.whitelisted ? 'success' : 'warning'"
            style="margin-bottom: 16px"
          >
            <template #header>
              {{ checkResult.whitelisted ? '该 IP 已在白名单中' : '该 IP 不在白名单中' }}
            </template>
          </n-alert>
          <div v-if="checkResult.whitelisted && checkResult.info" class="check-result-info">
            <n-descriptions :column="1" bordered>
              <n-descriptions-item label="IP 地址">
                {{ checkResult.info.ip }}
              </n-descriptions-item>
              <n-descriptions-item label="添加原因">
                {{ checkResult.info.reason || '无' }}
              </n-descriptions-item>
              <n-descriptions-item label="过期时间">
                {{ checkResult.info.expire_at ? formatDate(checkResult.info.expire_at) : '永久有效' }}
              </n-descriptions-item>
              <n-descriptions-item label="添加时间">
                {{ formatDate(checkResult.info.created_at) }}
              </n-descriptions-item>
            </n-descriptions>
          </div>
        </div>
        <template #footer>
          <n-space justify="end">
            <n-button @click="showCheckResultModal = false">关闭</n-button>
            <n-button
              v-if="checkResult?.whitelisted && checkResult.info"
              type="error"
              @click="handleDeleteFromCheck(checkResult.info.id)"
            >
              移除白名单
            </n-button>
          </n-space>
        </template>
      </n-card>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, h } from 'vue'
import {
  useMessage,
  useDialog,
  NButton,
  NSpace,
  NIcon,
  NModal,
  NCard,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NAlert,
  NDescriptions,
  NDescriptionsItem,
  type DataTableColumns,
  type FormInst,
  type FormRules
} from 'naive-ui'
import {
  getIPWhitelist,
  addIPWhitelist,
  deleteIPWhitelist,
  checkIPWhitelistStatus,
  cleanExpiredIPWhitelist,
  type IPWhitelist,
  type IPWhitelistCheckResult
} from '@/api/ipWhitelist'
import { formatDate } from '@/utils/format'

const message = useMessage()
const dialog = useDialog()

const loading = ref(false)
const whitelist = ref<IPWhitelist[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const isMobile = ref(false)

// 添加表单
const showAddModal = ref(false)
const adding = ref(false)
const addFormRef = ref<FormInst | null>(null)
const addForm = ref({
  ip: '',
  reason: '',
  duration: 0 // 默认永久
})

// 检查表单
const showCheckModal = ref(false)
const checking = ref(false)
const checkFormRef = ref<FormInst | null>(null)
const checkForm = ref({
  ip: ''
})

// 检查结果
const showCheckResultModal = ref(false)
const checkResult = ref<IPWhitelistCheckResult | null>(null)

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

// IP 地址或 CIDR 验证
const validateIPOrCIDR = (_rule: any, value: string): boolean | Error => {
  if (!value) {
    return new Error('请输入 IP 地址或 CIDR')
  }
  // IPv4 格式
  const ipv4Regex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/
  // IPv6 格式
  const ipv6Regex = /^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$|^::1$|^::$/
  // CIDR 格式
  const cidrRegex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\/([0-9]|[1-2][0-9]|3[0-2])$/
  
  if (ipv4Regex.test(value) || ipv6Regex.test(value) || cidrRegex.test(value)) {
    return true
  }
  return new Error('IP 地址格式不正确，支持 IPv4、IPv6 和 CIDR 格式（如：192.168.1.0/24）')
}

const addFormRules: FormRules = {
  ip: [
    { required: true, validator: validateIPOrCIDR, trigger: ['blur', 'input'] }
  ],
  duration: [
    { required: true, type: 'number', min: 0, max: 8760, message: '有效期必须在 0-8760 小时之间', trigger: 'blur' }
  ]
}

const checkFormRules: FormRules = {
  ip: [
    { required: true, validator: validateIPOrCIDR, trigger: ['blur', 'input'] }
  ]
}

// 表格列
const columns: DataTableColumns<IPWhitelist> = [
  {
    title: 'ID',
    key: 'id',
    width: 80
  },
  {
    title: 'IP 地址',
    key: 'ip',
    width: 180
  },
  {
    title: '添加原因',
    key: 'reason',
    ellipsis: { tooltip: true }
  },
  {
    title: '过期时间',
    key: 'expire_at',
    width: 180,
    render: (row) => row.expire_at ? formatDate(row.expire_at) : '永久有效'
  },
  {
    title: '添加时间',
    key: 'created_at',
    width: 180,
    render: (row) => formatDate(row.created_at)
  },
  {
    title: '操作',
    key: 'actions',
    width: 100,
    render: (row) =>
      h(
        NButton,
        {
          size: 'small',
          type: 'error',
          onClick: () => handleDelete(row.id, row.ip)
        },
        { default: () => '删除' }
      )
  }
]

// 获取白名单列表
async function fetchWhitelist() {
  try {
    loading.value = true
    const res = await getIPWhitelist({
      page: currentPage.value,
      page_size: pageSize.value
    })
    if (res.data) {
      whitelist.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch (error: any) {
    message.error(error.message || '获取白名单列表失败')
  } finally {
    loading.value = false
  }
}

// 添加 IP
async function handleAdd() {
  try {
    await addFormRef.value?.validate()
    adding.value = true
    await addIPWhitelist({
      ip: addForm.value.ip,
      reason: addForm.value.reason || undefined,
      duration: addForm.value.duration
    })
    message.success('添加成功')
    handleAddModalClose()
    fetchWhitelist()
  } catch (error: any) {
    if (!error?.errors) {
      message.error(error.message || '添加失败')
    }
  } finally {
    adding.value = false
  }
}

// 关闭添加对话框
function handleAddModalClose() {
  showAddModal.value = false
  addForm.value = {
    ip: '',
    reason: '',
    duration: 0
  }
  addFormRef.value?.restoreValidation()
}

// 删除 IP
function handleDelete(id: number, ip: string) {
  dialog.warning({
    title: '确认删除',
    content: `确定要从白名单中删除 IP "${ip}" 吗？删除后该 IP 将不再享受白名单豁免。`,
    positiveText: '确定删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteIPWhitelist(id)
        message.success('删除成功')
        fetchWhitelist()
      } catch (error: any) {
        message.error(error.message || '删除失败')
      }
    }
  })
}

// 检查 IP
async function handleCheck() {
  try {
    await checkFormRef.value?.validate()
    checking.value = true
    const res = await checkIPWhitelistStatus(checkForm.value.ip)
    if (res.data) {
      checkResult.value = res.data
      showCheckModal.value = false
      showCheckResultModal.value = true
    }
  } catch (error: any) {
    if (!error?.errors) {
      message.error(error.message || '检查失败')
    }
  } finally {
    checking.value = false
  }
}

// 关闭检查对话框
function handleCheckModalClose() {
  showCheckModal.value = false
  checkForm.value.ip = ''
  checkFormRef.value?.restoreValidation()
}

// 从检查结果中删除
function handleDeleteFromCheck(id: number) {
  handleDelete(id, checkResult.value?.info?.ip || '')
  showCheckResultModal.value = false
}

// 清理过期记录
function handleCleanExpired() {
  dialog.info({
    title: '清理过期记录',
    content: '确定要清理所有过期的白名单记录吗？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        const res = await cleanExpiredIPWhitelist()
        const count = res.data?.deleted_count || 0
        message.success(`清理成功，共删除 ${count} 条过期记录`)
        fetchWhitelist()
      } catch (error: any) {
        message.error(error.message || '清理失败')
      }
    }
  })
}

// 分页处理
function handlePageChange(page: number) {
  currentPage.value = page
  fetchWhitelist()
}

function handlePageSizeChange(size: number) {
  pageSize.value = size
  currentPage.value = 1
  fetchWhitelist()
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  fetchWhitelist()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped>
.ip-whitelist-manage-page {
  padding: 20px;
}

.page-title {
  margin-bottom: 16px;
  font-size: 24px;
  font-weight: 600;
}

.check-result-info {
  margin-top: 16px;
}

/* 移动端样式 */
@media (max-width: 768px) {
  .page-title {
    font-size: 20px;
  }
  
  .ip-whitelist-manage-page :deep(.n-data-table) {
    font-size: 13px;
  }
}
</style>

