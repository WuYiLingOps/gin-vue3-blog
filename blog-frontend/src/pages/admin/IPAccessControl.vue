<!--
 * @ProjectName: go-vue3-blog
 * @FileName: IPAccessControl.vue
 * @CreateTime: 2026-02-02 11:43:17
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: IP访问控制页面组件，提供IP黑名单和白名单的管理功能
 -->
<template>
  <div class="ip-access-control-page">
    <h1 class="page-title">IP 访问控制</h1>

    <!-- 标签页 -->
    <n-tabs v-model:value="activeTab" type="line" animated @update:value="onTabChange">
      <!-- 黑名单标签页 -->
      <n-tab-pane name="blacklist" tab="IP 黑名单">
        <!-- 操作栏 -->
        <n-space style="margin-bottom: 16px; margin-top: 16px" justify="space-between" align="center">
          <n-space>
            <n-button type="primary" @click="showAddBlacklistModal = true">
              <template #icon>
                <n-icon>
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                    <path fill="currentColor" d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
                  </svg>
                </n-icon>
              </template>
              添加 IP
            </n-button>
            <n-button @click="showCheckBlacklistModal = true">
              <template #icon>
                <n-icon>
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                    <path fill="currentColor" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
                  </svg>
                </n-icon>
              </template>
              检查 IP
            </n-button>
            <n-button type="warning" @click="handleCleanExpiredBlacklist">
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
          <n-button @click="fetchBlacklist">
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
          :columns="blacklistColumns"
          :data="blacklist"
          :loading="blacklistLoading"
          :pagination="blacklistPagination"
          :scroll-x="isMobile ? 1000 : undefined"
          :single-line="false"
          @update:page="handleBlacklistPageChange"
          @update:page-size="handleBlacklistPageSizeChange"
        />
      </n-tab-pane>

      <!-- 白名单标签页 -->
      <n-tab-pane name="whitelist" tab="IP 白名单">
        <!-- 操作栏 -->
        <n-space style="margin-bottom: 16px; margin-top: 16px" justify="space-between" align="center">
          <n-space>
            <n-button type="primary" @click="showAddWhitelistModal = true">
              <template #icon>
                <n-icon>
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                    <path fill="currentColor" d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
                  </svg>
                </n-icon>
              </template>
              添加 IP
            </n-button>
            <n-button @click="showCheckWhitelistModal = true">
              <template #icon>
                <n-icon>
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                    <path fill="currentColor" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
                  </svg>
                </n-icon>
              </template>
              检查 IP
            </n-button>
            <n-button type="warning" @click="handleCleanExpiredWhitelist">
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
          :columns="whitelistColumns"
          :data="whitelist"
          :loading="whitelistLoading"
          :pagination="whitelistPagination"
          :scroll-x="isMobile ? 1000 : undefined"
          :single-line="false"
          @update:page="handleWhitelistPageChange"
          @update:page-size="handleWhitelistPageSizeChange"
        />
      </n-tab-pane>
    </n-tabs>

    <!-- 添加黑名单 IP 对话框 -->
    <n-modal v-model:show="showAddBlacklistModal">
      <n-card
        title="添加 IP 到黑名单"
        :bordered="false"
        size="large"
        style="max-width: 600px"
        closable
        @close="handleAddBlacklistModalClose"
      >
        <n-form ref="addBlacklistFormRef" :model="addBlacklistForm" :rules="addBlacklistFormRules" label-width="100px">
          <n-form-item label="IP 地址" path="ip">
            <n-input
              v-model:value="addBlacklistForm.ip"
              placeholder="请输入 IP 地址（如：192.168.1.100）"
              @keyup.enter="handleAddBlacklist"
            />
          </n-form-item>
          <n-form-item label="封禁原因" path="reason">
            <n-input
              v-model:value="addBlacklistForm.reason"
              type="textarea"
              placeholder="请输入封禁原因（可选）"
              :autosize="{ minRows: 3, maxRows: 6 }"
            />
          </n-form-item>
          <n-form-item label="封禁时长" path="duration">
            <n-input-number
              v-model:value="addBlacklistForm.duration"
              :min="0"
              :max="8760"
              placeholder="小时"
              style="width: 100%"
            >
              <template #suffix>小时</template>
            </n-input-number>
            <template #feedback>
              <span style="font-size: 12px; color: #999">
                0 表示永久封禁，最大 8760 小时（1年）
              </span>
            </template>
          </n-form-item>
        </n-form>
        <template #footer>
          <n-space justify="end">
            <n-button @click="handleAddBlacklistModalClose">取消</n-button>
            <n-button type="primary" :loading="addingBlacklist" @click="handleAddBlacklist">
              确定添加
            </n-button>
          </n-space>
        </template>
      </n-card>
    </n-modal>

    <!-- 检查黑名单 IP 对话框 -->
    <n-modal v-model:show="showCheckBlacklistModal">
      <n-card
        title="检查 IP 状态"
        :bordered="false"
        size="large"
        style="max-width: 600px"
        closable
        @close="handleCheckBlacklistModalClose"
      >
        <n-form ref="checkBlacklistFormRef" :model="checkBlacklistForm" :rules="checkBlacklistFormRules" label-width="100px">
          <n-form-item label="IP 地址" path="ip">
            <n-input
              v-model:value="checkBlacklistForm.ip"
              placeholder="请输入要检查的 IP 地址"
              @keyup.enter="handleCheckBlacklist"
            />
          </n-form-item>
        </n-form>
        <template #footer>
          <n-space justify="end">
            <n-button @click="handleCheckBlacklistModalClose">关闭</n-button>
            <n-button type="primary" :loading="checkingBlacklist" @click="handleCheckBlacklist">
              检查
            </n-button>
          </n-space>
        </template>
      </n-card>
    </n-modal>

    <!-- 黑名单检查结果对话框 -->
    <n-modal v-model:show="showCheckBlacklistResultModal">
      <n-card
        title="IP 检查结果"
        :bordered="false"
        size="large"
        style="max-width: 600px"
        closable
        @close="showCheckBlacklistResultModal = false"
      >
        <div v-if="checkBlacklistResult">
          <n-alert
            :type="checkBlacklistResult.banned ? 'error' : 'success'"
            style="margin-bottom: 16px"
          >
            <template #header>
              {{ checkBlacklistResult.banned ? '该 IP 已被封禁' : '该 IP 未被封禁' }}
            </template>
          </n-alert>
          <div v-if="checkBlacklistResult.banned && checkBlacklistResult.info" class="check-result-info">
            <n-descriptions :column="1" bordered>
              <n-descriptions-item label="IP 地址">
                {{ checkBlacklistResult.info.ip }}
              </n-descriptions-item>
              <n-descriptions-item label="封禁类型">
                <n-tag :type="checkBlacklistResult.info.ban_type === 1 ? 'warning' : 'error'">
                  {{ checkBlacklistResult.info.ban_type === 1 ? '自动封禁' : '手动封禁' }}
                </n-tag>
              </n-descriptions-item>
              <n-descriptions-item label="封禁原因">
                {{ checkBlacklistResult.info.reason || '无' }}
              </n-descriptions-item>
              <n-descriptions-item label="过期时间">
                {{ checkBlacklistResult.info.expire_at ? formatDate(checkBlacklistResult.info.expire_at) : '永久封禁' }}
              </n-descriptions-item>
              <n-descriptions-item label="封禁时间">
                {{ formatDate(checkBlacklistResult.info.created_at) }}
              </n-descriptions-item>
            </n-descriptions>
          </div>
        </div>
        <template #footer>
          <n-space justify="end">
            <n-button @click="showCheckBlacklistResultModal = false">关闭</n-button>
            <n-button
              v-if="checkBlacklistResult?.banned && checkBlacklistResult.info"
              type="error"
              @click="handleDeleteBlacklistFromCheck(checkBlacklistResult.info.id)"
            >
              解除封禁
            </n-button>
          </n-space>
        </template>
      </n-card>
    </n-modal>

    <!-- 添加白名单 IP 对话框 -->
    <n-modal v-model:show="showAddWhitelistModal">
      <n-card
        title="添加 IP 到白名单"
        :bordered="false"
        size="large"
        style="max-width: 600px"
        closable
        @close="handleAddWhitelistModalClose"
      >
        <n-form ref="addWhitelistFormRef" :model="addWhitelistForm" :rules="addWhitelistFormRules" label-width="100px">
          <n-form-item label="IP 地址" path="ip">
            <n-input
              v-model:value="addWhitelistForm.ip"
              placeholder="请输入 IP 地址或 CIDR（如：192.168.1.100 或 192.168.1.0/24）"
              @keyup.enter="handleAddWhitelist"
            />
          </n-form-item>
          <n-form-item label="添加原因" path="reason">
            <n-input
              v-model:value="addWhitelistForm.reason"
              type="textarea"
              placeholder="请输入添加原因（可选）"
              :autosize="{ minRows: 3, maxRows: 6 }"
            />
          </n-form-item>
          <n-form-item label="有效期" path="duration">
            <n-input-number
              v-model:value="addWhitelistForm.duration"
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
            <n-button @click="handleAddWhitelistModalClose">取消</n-button>
            <n-button type="primary" :loading="addingWhitelist" @click="handleAddWhitelist">
              确定添加
            </n-button>
          </n-space>
        </template>
      </n-card>
    </n-modal>

    <!-- 检查白名单 IP 对话框 -->
    <n-modal v-model:show="showCheckWhitelistModal">
      <n-card
        title="检查 IP 状态"
        :bordered="false"
        size="large"
        style="max-width: 600px"
        closable
        @close="handleCheckWhitelistModalClose"
      >
        <n-form ref="checkWhitelistFormRef" :model="checkWhitelistForm" :rules="checkWhitelistFormRules" label-width="100px">
          <n-form-item label="IP 地址" path="ip">
            <n-input
              v-model:value="checkWhitelistForm.ip"
              placeholder="请输入要检查的 IP 地址"
              @keyup.enter="handleCheckWhitelist"
            />
          </n-form-item>
        </n-form>
        <template #footer>
          <n-space justify="end">
            <n-button @click="handleCheckWhitelistModalClose">关闭</n-button>
            <n-button type="primary" :loading="checkingWhitelist" @click="handleCheckWhitelist">
              检查
            </n-button>
          </n-space>
        </template>
      </n-card>
    </n-modal>

    <!-- 白名单检查结果对话框 -->
    <n-modal v-model:show="showCheckWhitelistResultModal">
      <n-card
        title="IP 检查结果"
        :bordered="false"
        size="large"
        style="max-width: 600px"
        closable
        @close="showCheckWhitelistResultModal = false"
      >
        <div v-if="checkWhitelistResult">
          <n-alert
            :type="checkWhitelistResult.whitelisted ? 'success' : 'warning'"
            style="margin-bottom: 16px"
          >
            <template #header>
              {{ checkWhitelistResult.whitelisted ? '该 IP 已在白名单中' : '该 IP 不在白名单中' }}
            </template>
          </n-alert>
          <div v-if="checkWhitelistResult.whitelisted && checkWhitelistResult.info" class="check-result-info">
            <n-descriptions :column="1" bordered>
              <n-descriptions-item label="IP 地址">
                {{ checkWhitelistResult.info.ip }}
              </n-descriptions-item>
              <n-descriptions-item label="添加原因">
                {{ checkWhitelistResult.info.reason || '无' }}
              </n-descriptions-item>
              <n-descriptions-item label="过期时间">
                {{ checkWhitelistResult.info.expire_at ? formatDate(checkWhitelistResult.info.expire_at) : '永久有效' }}
              </n-descriptions-item>
              <n-descriptions-item label="添加时间">
                {{ formatDate(checkWhitelistResult.info.created_at) }}
              </n-descriptions-item>
            </n-descriptions>
          </div>
        </div>
        <template #footer>
          <n-space justify="end">
            <n-button @click="showCheckWhitelistResultModal = false">关闭</n-button>
            <n-button
              v-if="checkWhitelistResult?.whitelisted && checkWhitelistResult.info"
              type="error"
              @click="handleDeleteWhitelistFromCheck(checkWhitelistResult.info.id)"
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
  NTag,
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
  NTabs,
  NTabPane,
  type DataTableColumns,
  type FormInst,
  type FormRules
} from 'naive-ui'
import {
  getIPBlacklist,
  addIPBlacklist,
  deleteIPBlacklist,
  checkIPStatus,
  cleanExpiredIPBlacklist,
  type IPBlacklist,
  type IPCheckResult
} from '@/api/ipBlacklist'
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

// 标签页
const activeTab = ref<'blacklist' | 'whitelist'>('blacklist')

// 黑名单相关
const blacklistLoading = ref(false)
const blacklist = ref<IPBlacklist[]>([])
const blacklistTotal = ref(0)
const blacklistCurrentPage = ref(1)
const blacklistPageSize = ref(20)

// 白名单相关
const whitelistLoading = ref(false)
const whitelist = ref<IPWhitelist[]>([])
const whitelistTotal = ref(0)
const whitelistCurrentPage = ref(1)
const whitelistPageSize = ref(20)

const isMobile = ref(false)

// 黑名单表单
const showAddBlacklistModal = ref(false)
const addingBlacklist = ref(false)
const addBlacklistFormRef = ref<FormInst | null>(null)
const addBlacklistForm = ref({
  ip: '',
  reason: '',
  duration: 24
})

const showCheckBlacklistModal = ref(false)
const checkingBlacklist = ref(false)
const checkBlacklistFormRef = ref<FormInst | null>(null)
const checkBlacklistForm = ref({
  ip: ''
})

const showCheckBlacklistResultModal = ref(false)
const checkBlacklistResult = ref<IPCheckResult | null>(null)

// 白名单表单
const showAddWhitelistModal = ref(false)
const addingWhitelist = ref(false)
const addWhitelistFormRef = ref<FormInst | null>(null)
const addWhitelistForm = ref({
  ip: '',
  reason: '',
  duration: 0
})

const showCheckWhitelistModal = ref(false)
const checkingWhitelist = ref(false)
const checkWhitelistFormRef = ref<FormInst | null>(null)
const checkWhitelistForm = ref({
  ip: ''
})

const showCheckWhitelistResultModal = ref(false)
const checkWhitelistResult = ref<IPWhitelistCheckResult | null>(null)

// 检测移动设备
function checkMobile() {
  isMobile.value = window.innerWidth <= 768
}

// 黑名单分页
const blacklistPagination = computed(() => ({
  page: blacklistCurrentPage.value,
  pageSize: blacklistPageSize.value,
  pageCount: Math.ceil(blacklistTotal.value / blacklistPageSize.value),
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50]
}))

// 白名单分页
const whitelistPagination = computed(() => ({
  page: whitelistCurrentPage.value,
  pageSize: whitelistPageSize.value,
  pageCount: Math.ceil(whitelistTotal.value / whitelistPageSize.value),
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50]
}))

// IP 地址验证
const validateIP = (_rule: any, value: string): boolean | Error => {
  if (!value) {
    return new Error('请输入 IP 地址')
  }
  const ipRegex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/
  const ipv6Regex = /^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$|^::1$|^::$/
  if (!ipRegex.test(value) && !ipv6Regex.test(value)) {
    return new Error('IP 地址格式不正确')
  }
  return true
}

// IP 地址或 CIDR 验证
const validateIPOrCIDR = (_rule: any, value: string): boolean | Error => {
  if (!value) {
    return new Error('请输入 IP 地址或 CIDR')
  }
  const ipv4Regex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/
  const ipv6Regex = /^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$|^::1$|^::$/
  const cidrRegex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\/([0-9]|[1-2][0-9]|3[0-2])$/
  
  if (ipv4Regex.test(value) || ipv6Regex.test(value) || cidrRegex.test(value)) {
    return true
  }
  return new Error('IP 地址格式不正确，支持 IPv4、IPv6 和 CIDR 格式（如：192.168.1.0/24）')
}

const addBlacklistFormRules: FormRules = {
  ip: [
    { required: true, validator: validateIP, trigger: ['blur', 'input'] }
  ],
  duration: [
    { required: true, type: 'number', min: 0, max: 8760, message: '封禁时长必须在 0-8760 小时之间', trigger: 'blur' }
  ]
}

const checkBlacklistFormRules: FormRules = {
  ip: [
    { required: true, validator: validateIP, trigger: ['blur', 'input'] }
  ]
}

const addWhitelistFormRules: FormRules = {
  ip: [
    { required: true, validator: validateIPOrCIDR, trigger: ['blur', 'input'] }
  ],
  duration: [
    { required: true, type: 'number', min: 0, max: 8760, message: '有效期必须在 0-8760 小时之间', trigger: 'blur' }
  ]
}

const checkWhitelistFormRules: FormRules = {
  ip: [
    { required: true, validator: validateIPOrCIDR, trigger: ['blur', 'input'] }
  ]
}

// 黑名单表格列
const blacklistColumns: DataTableColumns<IPBlacklist> = [
  {
    title: 'ID',
    key: 'id',
    width: 80
  },
  {
    title: 'IP 地址',
    key: 'ip',
    width: 160
  },
  {
    title: '封禁类型',
    key: 'ban_type',
    width: 120,
    render: (row) =>
      h(
        NTag,
        { type: row.ban_type === 1 ? 'warning' : 'error', size: 'small' },
        { default: () => (row.ban_type === 1 ? '自动封禁' : '手动封禁') }
      )
  },
  {
    title: '封禁原因',
    key: 'reason',
    ellipsis: { tooltip: true }
  },
  {
    title: '过期时间',
    key: 'expire_at',
    width: 180,
    render: (row) => row.expire_at ? formatDate(row.expire_at) : '永久封禁'
  },
  {
    title: '封禁时间',
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
          onClick: () => handleDeleteBlacklist(row.id, row.ip)
        },
        { default: () => '删除' }
      )
  }
]

// 白名单表格列
const whitelistColumns: DataTableColumns<IPWhitelist> = [
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
          onClick: () => handleDeleteWhitelist(row.id, row.ip)
        },
        { default: () => '删除' }
      )
  }
]

// 获取黑名单列表
async function fetchBlacklist() {
  try {
    blacklistLoading.value = true
    const res = await getIPBlacklist({
      page: blacklistCurrentPage.value,
      page_size: blacklistPageSize.value
    })
    if (res.data) {
      blacklist.value = res.data.list || []
      blacklistTotal.value = res.data.total || 0
    }
  } catch (error: any) {
    message.error(error.message || '获取黑名单列表失败')
  } finally {
    blacklistLoading.value = false
  }
}

// 获取白名单列表
async function fetchWhitelist() {
  try {
    whitelistLoading.value = true
    const res = await getIPWhitelist({
      page: whitelistCurrentPage.value,
      page_size: whitelistPageSize.value
    })
    if (res.data) {
      whitelist.value = res.data.list || []
      whitelistTotal.value = res.data.total || 0
    }
  } catch (error: any) {
    message.error(error.message || '获取白名单列表失败')
  } finally {
    whitelistLoading.value = false
  }
}

// 添加黑名单 IP
async function handleAddBlacklist() {
  try {
    await addBlacklistFormRef.value?.validate()
    addingBlacklist.value = true
    await addIPBlacklist({
      ip: addBlacklistForm.value.ip,
      reason: addBlacklistForm.value.reason || undefined,
      duration: addBlacklistForm.value.duration
    })
    message.success('添加成功')
    handleAddBlacklistModalClose()
    fetchBlacklist()
  } catch (error: any) {
    if (!error?.errors) {
      message.error(error.message || '添加失败')
    }
  } finally {
    addingBlacklist.value = false
  }
}

// 关闭添加黑名单对话框
function handleAddBlacklistModalClose() {
  showAddBlacklistModal.value = false
  addBlacklistForm.value = {
    ip: '',
    reason: '',
    duration: 24
  }
  addBlacklistFormRef.value?.restoreValidation()
}

// 删除黑名单 IP
function handleDeleteBlacklist(id: number, ip: string) {
  dialog.warning({
    title: '确认删除',
    content: `确定要从黑名单中删除 IP "${ip}" 吗？删除后该 IP 将可以正常访问。`,
    positiveText: '确定删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteIPBlacklist(id)
        message.success('删除成功')
        fetchBlacklist()
      } catch (error: any) {
        message.error(error.message || '删除失败')
      }
    }
  })
}

// 检查黑名单 IP
async function handleCheckBlacklist() {
  try {
    await checkBlacklistFormRef.value?.validate()
    checkingBlacklist.value = true
    const res = await checkIPStatus(checkBlacklistForm.value.ip)
    if (res.data) {
      checkBlacklistResult.value = res.data
      showCheckBlacklistModal.value = false
      showCheckBlacklistResultModal.value = true
    }
  } catch (error: any) {
    if (!error?.errors) {
      message.error(error.message || '检查失败')
    }
  } finally {
    checkingBlacklist.value = false
  }
}

// 关闭检查黑名单对话框
function handleCheckBlacklistModalClose() {
  showCheckBlacklistModal.value = false
  checkBlacklistForm.value.ip = ''
  checkBlacklistFormRef.value?.restoreValidation()
}

// 从检查结果中删除黑名单
function handleDeleteBlacklistFromCheck(id: number) {
  handleDeleteBlacklist(id, checkBlacklistResult.value?.info?.ip || '')
  showCheckBlacklistResultModal.value = false
}

// 清理过期黑名单记录
function handleCleanExpiredBlacklist() {
  dialog.info({
    title: '清理过期记录',
    content: '确定要清理所有过期的黑名单记录吗？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        const res = await cleanExpiredIPBlacklist()
        const count = res.data?.deleted_count || 0
        message.success(`清理成功，共删除 ${count} 条过期记录`)
        fetchBlacklist()
      } catch (error: any) {
        message.error(error.message || '清理失败')
      }
    }
  })
}

// 黑名单分页处理
function handleBlacklistPageChange(page: number) {
  blacklistCurrentPage.value = page
  fetchBlacklist()
}

function handleBlacklistPageSizeChange(size: number) {
  blacklistPageSize.value = size
  blacklistCurrentPage.value = 1
  fetchBlacklist()
}

// 添加白名单 IP
async function handleAddWhitelist() {
  try {
    await addWhitelistFormRef.value?.validate()
    addingWhitelist.value = true
    await addIPWhitelist({
      ip: addWhitelistForm.value.ip,
      reason: addWhitelistForm.value.reason || undefined,
      duration: addWhitelistForm.value.duration
    })
    message.success('添加成功')
    handleAddWhitelistModalClose()
    fetchWhitelist()
  } catch (error: any) {
    if (!error?.errors) {
      message.error(error.message || '添加失败')
    }
  } finally {
    addingWhitelist.value = false
  }
}

// 关闭添加白名单对话框
function handleAddWhitelistModalClose() {
  showAddWhitelistModal.value = false
  addWhitelistForm.value = {
    ip: '',
    reason: '',
    duration: 0
  }
  addWhitelistFormRef.value?.restoreValidation()
}

// 删除白名单 IP
function handleDeleteWhitelist(id: number, ip: string) {
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

// 检查白名单 IP
async function handleCheckWhitelist() {
  try {
    await checkWhitelistFormRef.value?.validate()
    checkingWhitelist.value = true
    const res = await checkIPWhitelistStatus(checkWhitelistForm.value.ip)
    if (res.data) {
      checkWhitelistResult.value = res.data
      showCheckWhitelistModal.value = false
      showCheckWhitelistResultModal.value = true
    }
  } catch (error: any) {
    if (!error?.errors) {
      message.error(error.message || '检查失败')
    }
  } finally {
    checkingWhitelist.value = false
  }
}

// 关闭检查白名单对话框
function handleCheckWhitelistModalClose() {
  showCheckWhitelistModal.value = false
  checkWhitelistForm.value.ip = ''
  checkWhitelistFormRef.value?.restoreValidation()
}

// 从检查结果中删除白名单
function handleDeleteWhitelistFromCheck(id: number) {
  handleDeleteWhitelist(id, checkWhitelistResult.value?.info?.ip || '')
  showCheckWhitelistResultModal.value = false
}

// 清理过期白名单记录
function handleCleanExpiredWhitelist() {
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

// 白名单分页处理
function handleWhitelistPageChange(page: number) {
  whitelistCurrentPage.value = page
  fetchWhitelist()
}

function handleWhitelistPageSizeChange(size: number) {
  whitelistPageSize.value = size
  whitelistCurrentPage.value = 1
  fetchWhitelist()
}

// 标签页切换时加载对应数据
function onTabChange(tab: string) {
  if (tab === 'blacklist' && blacklist.value.length === 0) {
    fetchBlacklist()
  } else if (tab === 'whitelist' && whitelist.value.length === 0) {
    fetchWhitelist()
  }
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  fetchBlacklist()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped>
.ip-access-control-page {
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
  
  .ip-access-control-page :deep(.n-data-table) {
    font-size: 13px;
  }
}
</style>

