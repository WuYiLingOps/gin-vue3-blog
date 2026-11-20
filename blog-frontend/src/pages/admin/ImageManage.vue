<!--
 * @ProjectName: go-vue3-blog
 * @FileName: ImageManage.vue
 * @CreateTime: 2026-11-20 10:00:00
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 图片管理页面组件，提供全站图片盘点和图床间迁移功能
 -->
<template>
  <div class="image-manage-page">
    <div class="header">
      <n-space>
        <template v-if="activeTab === 'assets'">
          <n-button
            type="primary"
            :size="isMobile ? 'small' : 'medium'"
            :loading="scanning"
            @click="handleScan(true)"
          >
            <template #icon>
              <n-icon :component="RefreshOutline" />
            </template>
            <span v-if="!isMobile">扫描全站图片</span>
            <span v-else>扫描</span>
          </n-button>
          <n-button
            type="info"
            :size="isMobile ? 'small' : 'medium'"
            :disabled="migrateStatus.running && selectedUrls.length === 0 && !showProgressModal"
            @click="migrateStatus.running ? (showProgressModal = true) : (showMigrateModal = true)"
          >
            <template #icon>
              <n-icon :component="SwapHorizontalOutline" />
            </template>
            {{
              migrateStatus.running
                ? '查看迁移进度'
                : `迁移选中图片（${selectedUrls.length}）`
            }}
          </n-button>
        </template>
        <template v-else>
          <n-button
            :size="isMobile ? 'small' : 'medium'"
            :loading="logsLoading"
            @click="fetchLogs"
          >
            <template #icon>
              <n-icon :component="RefreshOutline" />
            </template>
            刷新日志
          </n-button>
          <n-button
            type="error"
            :size="isMobile ? 'small' : 'medium'"
            :disabled="logs.length === 0"
            @click="handleClearLogs"
          >
            清空日志
          </n-button>
        </template>
      </n-space>
    </div>

    <!-- 主视图切换：我的相册 / 图片资产 / 迁移日志 -->
    <n-tabs v-model:value="activeTab" type="line" size="small" class="main-tabs">
      <n-tab name="album">我的相册</n-tab>
      <n-tab name="assets">图片资产</n-tab>
      <n-tab name="logs">迁移日志</n-tab>
    </n-tabs>

    <template v-if="activeTab === 'assets'">
    <n-alert type="info" :show-icon="true" class="tips-alert" closable>
      盘点范围覆盖文章封面/正文、用户头像、说说、相册、友链、聊天头像及系统配置中的全部图片。
      点击「使用位置」可查看图片来源明细；外链图片仅可迁入自有图床；EasyImage 图床不支持 API 删除源文件，迁移后需登录图床手动清理。
    </n-alert>

    <!-- 筛选栏（本地过滤，不发请求） -->
    <div class="filter-bar">
      <n-tabs v-model:value="filterStorage" type="segment" size="small">
        <n-tab name="">全部</n-tab>
        <n-tab name="local">本地存储</n-tab>
        <n-tab name="oss">阿里云OSS</n-tab>
        <n-tab name="cos">腾讯云COS</n-tab>
        <n-tab name="easyimage">EasyImage</n-tab>
        <n-tab name="external">外链图片</n-tab>
      </n-tabs>
      <n-input
        v-model:value="keyword"
        placeholder="按 URL 关键字筛选"
        clearable
        size="small"
        class="keyword-input"
      />
    </div>

    <!-- 资产表格（客户端分页，每页 20 条，图片懒加载） -->
    <!-- remote：数据在前端已完成全量+本地过滤，由本组件自行切片，避免 n-data-table 按切片长度钳制页码 -->
    <n-data-table
      remote
      :columns="columns"
      :data="pagedAssets"
      :loading="scanning"
      :single-line="false"
      :row-key="(row: ImageAsset) => row.url"
      :checked-row-keys="checkedKeys"
      :pagination="{
        page: pagination.page,
        pageSize: pagination.pageSize,
        itemCount: filteredAssets.length,
        showSizePicker: true,
        pageSizes: [20, 50, 100],
        onChange: (page: number) => (pagination.page = page),
        onUpdatePageSize: (size: number) => {
          pagination.pageSize = size
          pagination.page = 1
        }
      }"
      max-height="calc(100vh - 380px)"
      @update:checked-row-keys="handleCheck"
    />
    </template>

    <!-- 迁移日志视图 -->
    <template v-else-if="activeTab === 'logs'">
      <!-- remote：数据为后端分页返回的当前页，由本组件控制页码，避免按切片长度钳制页码 -->
      <n-data-table
        remote
        :columns="logColumns"
        :data="logs"
        :loading="logsLoading"
        :single-line="false"
        :row-key="(row: ImageMigrateLog) => row.id"
        :pagination="logPagination"
        max-height="calc(100vh - 320px)"
      />
    </template>

    <!-- 我的相册视图 -->
    <template v-else>
      <AlbumPanel />
    </template>

    <!-- 图片来源明细弹窗 -->
    <n-modal
      v-model:show="showSourceModal"
      preset="card"
      title="图片来源明细"
      :style="{ width: isMobile ? '95%' : '640px', maxWidth: '95vw' }"
    >
      <n-space vertical size="small">
        <n-text depth="3" style="font-size: 12px; word-break: break-all">
          {{ sourceAsset?.url }}
        </n-text>
        <n-text>
          共被引用 <b>{{ sourceAsset?.ref_count }}</b> 次，出现在以下 {{ sourceAsset?.locations.length }} 处：
        </n-text>
        <div
          v-for="(loc, idx) in sourceAsset?.locations || []"
          :key="idx"
          class="source-item"
        >
          <n-tag size="small" type="info" :bordered="false">{{ loc.table_label }}</n-tag>
          <n-text depth="3" style="font-size: 12px">
            {{ loc.table }}.{{ loc.field }} · 记录 #{{ loc.row_id }}
            <template v-if="loc.kind === 'markdown'">（正文内嵌图）</template>
            <template v-else-if="loc.kind === 'json_array'">（JSON 数组）</template>
            <template v-else-if="loc.kind === 'settings'">（系统配置）</template>
          </n-text>
        </div>
      </n-space>
    </n-modal>

    <!-- 迁移日志明细弹窗 -->
    <n-modal
      v-model:show="showLogDetailModal"
      preset="card"
      :title="logDetailType === 'success' ? '迁移成功映射' : '迁移失败明细'"
      :style="{ width: isMobile ? '95%' : '720px', maxWidth: '95vw' }"
    >
      <n-text depth="3" style="font-size: 12px">
        {{ logDetailAsset?.message }}
      </n-text>
      <div class="result-section" style="max-height: 380px; margin-top: 12px">
        <template v-if="logDetailType === 'success'">
          <div
            v-for="item in logDetailAsset?.mappings || []"
            :key="item.old_url"
            class="result-item"
          >
            <span class="result-url">{{ item.old_url }}</span>
            <n-icon :component="ArrowForwardOutline" size="14" />
            <span class="result-url new">{{ item.new_url }}</span>
          </div>
        </template>
        <template v-else>
          <div
            v-for="item in logDetailAsset?.failed_items || []"
            :key="item.url"
            class="result-item error"
          >
            <span class="result-url">{{ item.url }}</span>
            <n-text type="error" depth="3">{{ item.error }}</n-text>
          </div>
        </template>
      </div>
    </n-modal>

    <!-- 迁移确认弹窗 -->
    <n-modal
      v-model:show="showMigrateModal"
      preset="card"
      title="迁移图片"
      :style="{ width: isMobile ? '95%' : '520px', maxWidth: '95vw' }"
    >
      <n-space vertical size="large">
        <n-alert type="warning" :show-icon="true">
          即将迁移 <b>{{ selectedUrls.length }}</b> 张图片到
          <b>{{ targetStorageLabel }}</b>。迁移过程中请勿关闭页面，可在任务结束后查看结果。
        </n-alert>

        <n-form-item label="目标存储" label-placement="left">
          <n-select v-model:value="targetStorage" :options="targetOptions" />
        </n-form-item>

        <n-form-item label-placement="left">
          <template #label>
            <span>删除源文件</span>
          </template>
          <n-switch v-model:value="deleteSource" />
          <span class="delete-source-tip">
            迁移成功后删除源图床上的文件（EasyImage 与外链无法删除，将自动跳过）
          </span>
        </n-form-item>

        <n-alert v-if="hasEasyImageSource" type="warning" :show-icon="true">
          选中图片中包含 EasyImage 来源，该图床不支持 API 删除，迁移后需登录图床手动清理源文件。
        </n-alert>
      </n-space>

      <template #footer>
        <n-space justify="end">
          <n-button @click="showMigrateModal = false">取消</n-button>
          <n-button type="primary" :loading="migrateStatus.running" @click="handleMigrate">
            开始迁移
          </n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- 迁移进度弹窗（运行中也可关闭，任务继续在后台执行，点「查看进度」可重新打开） -->
    <n-modal
      v-model:show="showProgressModal"
      preset="card"
      title="迁移进度"
      :style="{ width: isMobile ? '95%' : '640px', maxWidth: '95vw' }"
    >
      <n-space vertical size="large">
        <n-progress
          type="line"
          :percentage="progressPercent"
          :status="migrateStatus.running ? 'default' : hasFailure ? 'error' : 'success'"
          :show-indicator="true"
        />
        <n-text>
          {{ progressText }}
        </n-text>

        <!-- 实时步骤（运行中展示当前处理到哪张图、哪个环节、已耗时） -->
        <n-alert
          v-if="migrateStatus.running"
          type="info"
          :show-icon="false"
          class="current-step-alert"
        >
          <div>当前步骤：{{ migrateStatus.current_step || '准备中…' }}</div>
          <div v-if="migrateStatus.current_url" class="current-url">
            {{ migrateStatus.current_url }}
          </div>
          <n-text depth="3" style="font-size: 12px">
            已进行 {{ migrateStatus.elapsed_sec || 0 }} 秒 · 单图最长约 2 分钟（下载/上传各 60 秒超时），超时会自动标记失败并继续
          </n-text>
        </n-alert>

        <template v-if="!migrateStatus.running && migrateStatus.finished">
          <n-alert :type="hasFailure ? 'warning' : 'success'" :show-icon="true">
            {{ migrateStatus.message }}
          </n-alert>

          <div v-if="migrateStatus.failed.length > 0" class="result-section">
            <n-h6>失败列表（{{ migrateStatus.failed.length }}）</n-h6>
            <div v-for="item in migrateStatus.failed" :key="item.url" class="result-item error">
              <span class="result-url">{{ item.url }}</span>
              <n-text type="error" depth="3">{{ item.error }}</n-text>
            </div>
          </div>

          <div v-if="migrateStatus.success.length > 0" class="result-section">
            <n-h6>成功映射（{{ migrateStatus.success.length }}）</n-h6>
            <div v-for="item in migrateStatus.success" :key="item.old_url" class="result-item">
              <span class="result-url">{{ item.old_url }}</span>
              <n-icon :component="ArrowForwardOutline" size="14" />
              <span class="result-url new">{{ item.new_url }}</span>
            </div>
          </div>
        </template>
      </n-space>

      <template #footer>
        <n-space justify="end">
          <n-button
            v-if="!migrateStatus.running && hasFailure"
            type="warning"
            @click="handleRetryFailed"
          >
            重试失败项
          </n-button>
          <n-button @click="showProgressModal = false">
            {{ migrateStatus.running ? '后台运行' : '关闭' }}
          </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, h, onMounted, onUnmounted } from 'vue'
import { useMessage, useDialog, NTag, NSpace, NButton } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import {
  RefreshOutline,
  SwapHorizontalOutline,
  ArrowForwardOutline
} from '@vicons/ionicons5'
import {
  scanImageAssets,
  migrateImages,
  getMigrateStatus,
  getMigrateLogs,
  deleteMigrateLogs,
  type ImageAsset,
  type ImageMigrateLog,
  type MigrateStatus
} from '@/api/imageAsset'
import { getFileUrl } from '@/utils/request'
import AlbumPanel from '@/components/AlbumPanel.vue'

const message = useMessage()
const dialog = useDialog()
const isMobile = ref(window.innerWidth <= 1100)

// 主视图 Tab：我的相册 / 图片资产 / 迁移日志（相册为默认视图）
const activeTab = ref<'album' | 'assets' | 'logs'>('album')

// 图片来源明细弹窗
const showSourceModal = ref(false)
const sourceAsset = ref<ImageAsset | null>(null)

// 迁移日志
const logs = ref<ImageMigrateLog[]>([])
const logsLoading = ref(false)
const logPagination = reactive({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page: number) => {
    logPagination.page = page
    fetchLogs()
  },
  onUpdatePageSize: (size: number) => {
    logPagination.pageSize = size
    logPagination.page = 1
    fetchLogs()
  }
})

// 归一化成功映射：新格式为数组 [{old_url,new_url}]，历史格式为对象 {old_url: new_url}
function normalizeMappings(
  raw: { old_url: string; new_url: string }[] | Record<string, string> | null | undefined
): { old_url: string; new_url: string }[] {
  if (!raw) return []
  if (Array.isArray(raw)) return raw
  // 对象格式转换（历史日志）
  return Object.entries(raw).map(([oldUrl, newUrl]) => ({ old_url: oldUrl, new_url: newUrl }))
}

// 查询迁移日志
async function fetchLogs() {
  try {
    logsLoading.value = true
    const res = await getMigrateLogs(logPagination.page, logPagination.pageSize)
    if (res.data) {
      // 后端 JSONB 字段可能返回数组或 null，统一归一化；
      // 历史日志的 mappings 曾存为 {old: new} 对象格式，转换为数组格式以正常展示
      logs.value = (res.data.list || []).map(item => ({
        ...item,
        mappings: normalizeMappings(item.mappings),
        failed_items: item.failed_items || []
      }))
      logPagination.itemCount = res.data.total || 0
    }
  } catch (error: any) {
    message.error(error.message || '获取迁移日志失败')
    console.error(error)
  } finally {
    logsLoading.value = false
  }
}

// 清空全部迁移日志
function handleClearLogs() {
  dialog.warning({
    title: '清空迁移日志',
    content: '确定要清空全部迁移日志吗？此操作不可恢复（不影响已迁移的图片和已替换的地址）。',
    positiveText: '清空',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteMigrateLogs([])
        message.success('迁移日志已清空')
        fetchLogs()
      } catch (error: any) {
        message.error(error.message || '清空迁移日志失败')
      }
    }
  })
}

// 打开图片来源明细弹窗
function openSourceModal(asset: ImageAsset) {
  sourceAsset.value = asset
  showSourceModal.value = true
}

// 全部图片资产（后端一次返回的全集，筛选在前端本地完成，切换图床分类零请求）
const assets = ref<ImageAsset[]>([])
const scanning = ref(false)
const filterStorage = ref('')
const keyword = ref('')
const checkedKeys = ref<string[]>([])

// 客户端分页：避免一次渲染几百行表格和图片
const pagination = reactive({
  page: 1,
  pageSize: 20
})
const pagedAssets = computed(() => {
  const start = (pagination.page - 1) * pagination.pageSize
  return filteredAssets.value.slice(start, start + pagination.pageSize)
})
// 筛选条件变化时回到第一页
watch([filterStorage, keyword], () => {
  pagination.page = 1
})

// 本地过滤：按存储类型 + URL 关键字
const filteredAssets = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  return assets.value.filter(asset => {
    if (filterStorage.value && asset.storage_type !== filterStorage.value) return false
    if (kw && !asset.url.toLowerCase().includes(kw)) return false
    return true
  })
})

// 迁移表单
const showMigrateModal = ref(false)
const showProgressModal = ref(false)
const targetStorage = ref<'local' | 'oss' | 'cos' | 'easyimage'>('easyimage')
const deleteSource = ref(false)
let pollTimer: number | null = null

const migrateStatus = ref<MigrateStatus>({
  running: false,
  total: 0,
  migrated: 0,
  failed: [],
  success: [],
  finished: false,
  message: ''
})

const targetOptions = [
  { label: '本地存储', value: 'local' },
  { label: '阿里云OSS', value: 'oss' },
  { label: '腾讯云COS', value: 'cos' },
  { label: 'EasyImage 图床', value: 'easyimage' }
]

const targetStorageLabel = computed(
  () => targetOptions.find(opt => opt.value === targetStorage.value)?.label || targetStorage.value
)

// 选中的 URL 列表
const selectedUrls = computed(() =>
  assets.value
    .filter(asset => checkedKeys.value.includes(asset.url))
    .map(asset => asset.url)
)

// 选中图片中是否包含 EasyImage 来源
const hasEasyImageSource = computed(() =>
  assets.value.some(
    asset => checkedKeys.value.includes(asset.url) && asset.storage_type === 'easyimage'
  )
)

const hasFailure = computed(() => migrateStatus.value.failed.length > 0)

const progressPercent = computed(() => {
  if (migrateStatus.value.total === 0) return 0
  return Math.round((migrateStatus.value.migrated / migrateStatus.value.total) * 100)
})

const progressText = computed(() => {
  const { running, migrated, total } = migrateStatus.value
  if (running) {
    return `正在迁移：${migrated} / ${total}（5 张并发处理，请耐心等待）`
  }
  if (!migrateStatus.value.finished) return '等待任务开始…'
  return `任务结束：成功 ${migrateStatus.value.success.length} 张，失败 ${migrateStatus.value.failed.length} 张`
})

const storageLabelMap: Record<string, string> = {
  local: '本地存储',
  oss: '阿里云OSS',
  cos: '腾讯云COS',
  easyimage: 'EasyImage',
  external: '外链图片'
}

// 表格列定义
const columns: DataTableColumns<ImageAsset> = [
  { type: 'selection' },
  {
    title: '预览',
    key: 'preview',
    width: 90,
    render(row) {
      return h('img', {
        src: getFileUrl(row.url),
        loading: 'lazy',
        style: {
          width: '56px',
          height: '56px',
          objectFit: 'cover',
          borderRadius: '4px',
          cursor: 'pointer',
          display: 'block',
          margin: '4px auto',
          background: 'rgba(128, 128, 128, 0.1)'
        },
        onClick: () => window.open(getFileUrl(row.url), '_blank')
      })
    }
  },
  {
    title: 'URL',
    key: 'url',
    minWidth: 260,
    ellipsis: { tooltip: true },
    render(row) {
      return h(
        'a',
        {
          href: getFileUrl(row.url),
          target: '_blank',
          style: { wordBreak: 'break-all', fontSize: '12px' },
          onClick: (e: Event) => {
            e.preventDefault()
            window.open(getFileUrl(row.url), '_blank')
          }
        },
        row.url
      )
    }
  },
  {
    title: '来源存储',
    key: 'storage_type',
    width: 110,
    render(row) {
      return h(
        'span',
        {
          class: `storage-tag storage-${row.storage_type}`
        },
        storageLabelMap[row.storage_type] || row.storage_type
      )
    }
  },
  {
    title: '使用位置',
    key: 'locations',
    width: 190,
    render(row) {
      // 展示前两处中文位置，点击打开来源明细弹窗
      const labels = (row.locations || []).map(loc => loc.table_label)
      const display = labels.slice(0, 2).join('、') + (labels.length > 2 ? ` 等${labels.length}处` : '')
      return h(
        'a',
        {
          style: { fontSize: '12px', cursor: 'pointer', color: '#2080f0' },
          onClick: () => openSourceModal(row)
        },
        display || '查看来源'
      )
    }
  },
  {
    title: '引用',
    key: 'ref_count',
    width: 70
  },
  {
    title: '可删源',
    key: 'deletable',
    width: 80,
    render(row) {
      return h(
        'span',
        { style: { fontSize: '12px', color: row.deletable ? '#18a058' : '#d03050' } },
        row.deletable ? '是' : '否'
      )
    }
  }
]

// 迁移日志表格列定义
const targetStorageTextMap: Record<string, string> = {
  local: '本地存储',
  oss: '阿里云OSS',
  cos: '腾讯云COS',
  easyimage: 'EasyImage'
}

const logColumns: DataTableColumns<ImageMigrateLog> = [
  {
    title: '完成时间',
    key: 'created_at',
    width: 160,
    render(row) {
      return new Date(row.created_at).toLocaleString('zh-CN')
    }
  },
  {
    title: '操作人',
    key: 'operator_name',
    width: 100,
    render(row) {
      return row.operator_name || '-'
    }
  },
  {
    title: '目标图床',
    key: 'target_storage',
    width: 110,
    render(row) {
      return h(
        'span',
        { class: `storage-tag storage-${row.target_storage}` },
        targetStorageTextMap[row.target_storage] || row.target_storage
      )
    }
  },
  {
    title: '结果',
    key: 'counts',
    width: 190,
    render(row) {
      return h('div', { style: { fontSize: '12px', lineHeight: '1.7' } }, [
        h('div', `迁移 ${row.total_count} 张：成功 ${row.success_count}，失败 ${row.failed_count}`),
        h('div', `替换数据库 ${row.updated_rows} 处${row.delete_source ? `，删源 ${row.deleted_files} 个` : ''}`),
        h('div', { style: { color: '#888' } }, `耗时 ${row.duration_sec} 秒`)
      ])
    }
  },
  {
    title: '汇总信息',
    key: 'message',
    minWidth: 220,
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', { style: { fontSize: '12px' } }, row.message || '-')
    }
  },
  {
    title: '明细',
    key: 'detail',
    width: 130,
    fixed: 'right',
    render(row) {
      return h(NSpace, { size: 'small' }, {
        default: () => [
          h(
            NButton,
            {
              size: 'tiny',
              onClick: () => openLogDetail(row, 'success')
            },
            { default: () => `映射 ${row.mappings.length}` }
          ),
          row.failed_items.length > 0
            ? h(
                NButton,
                {
                  size: 'tiny',
                  type: 'error',
                  onClick: () => openLogDetail(row, 'failed')
                },
                { default: () => `失败 ${row.failed_items.length}` }
              )
            : null
        ]
      })
    }
  }
]

// 日志明细弹窗
const showLogDetailModal = ref(false)
const logDetailAsset = ref<ImageMigrateLog | null>(null)
const logDetailType = ref<'success' | 'failed'>('success')

function openLogDetail(log: ImageMigrateLog, type: 'success' | 'failed') {
  logDetailAsset.value = log
  logDetailType.value = type
  showLogDetailModal.value = true
}

// 扫描全站图片
// refresh=true 时后端强制刷新缓存重新扫描数据库；日常进入页面/迁移完成后的刷新走 Redis 缓存
async function handleScan(refresh: boolean = false) {
  try {
    scanning.value = true
    const res = await scanImageAssets(refresh)
    assets.value = res.data || []
    // 过滤掉不在当前列表里的勾选项
    const urlSet = new Set(assets.value.map(a => a.url))
    checkedKeys.value = checkedKeys.value.filter(url => urlSet.has(url))
    if (assets.value.length === 0) {
      message.info(refresh ? '未扫描到图片' : '未扫描到图片，可点击「扫描全站图片」刷新')
    }
  } catch (error: any) {
    message.error(error.message || '扫描图片失败')
    console.error(error)
  } finally {
    scanning.value = false
  }
}

function handleCheck(keys: (string | number)[]) {
  checkedKeys.value = keys.map(String)
}

// 发起迁移
async function handleMigrate() {
  try {
    await migrateImages({
      urls: selectedUrls.value,
      target_storage: targetStorage.value,
      delete_source: deleteSource.value
    })
    showMigrateModal.value = false
    showProgressModal.value = true
    startPolling()
  } catch (error: any) {
    message.error(error.message || '启动迁移任务失败')
    console.error(error)
  }
}

// 重试失败项
function handleRetryFailed() {
  const failedUrls = migrateStatus.value.failed.map(item => item.url)
  const failedAssets = assets.value.filter(asset => failedUrls.includes(asset.url))
  checkedKeys.value = failedUrls
  // 直接打开迁移弹窗复用同一流程
  if (failedAssets.length === 0) {
    // 失败 URL 可能不在当前筛选列表里，直接发起
    startMigration(failedUrls)
  } else {
    showMigrateModal.value = true
  }
}

async function startMigration(urls: string[]) {
  try {
    await migrateImages({
      urls,
      target_storage: targetStorage.value,
      delete_source: deleteSource.value
    })
    showMigrateModal.value = false
    showProgressModal.value = true
    startPolling()
  } catch (error: any) {
    message.error(error.message || '启动迁移任务失败')
    console.error(error)
  }
}

// 轮询迁移进度
function startPolling() {
  stopPolling()
  pollMigrateStatus()
  pollTimer = window.setInterval(pollMigrateStatus, 2000)
}

function stopPolling() {
  if (pollTimer !== null) {
    clearInterval(pollTimer)
    pollTimer = null
  }
}

async function pollMigrateStatus() {
  try {
    const res = await getMigrateStatus()
    if (res.data) {
      // 防御性归一化：后端异常或旧版本返回 null 时兜底为空数组
      const status = {
        ...res.data,
        failed: res.data.failed || [],
        success: res.data.success || []
      }
      const prevFailedCount = migrateStatus.value.failed?.length || 0
      migrateStatus.value = status
      // 运行中发现新的失败项，立即提示（避免直到结束才发现全部失败）
      if (status.running && status.failed.length > prevFailedCount) {
        const latest = status.failed[status.failed.length - 1]
        message.warning(`图片迁移失败：${latest.url}（${latest.error}）`)
      }
      if (!status.running && status.finished) {
        stopPolling()
        // 任务结束后刷新列表（后端已在替换数据库后失效缓存，此处走重新扫描）
        handleScan(true)
      }
    }
  } catch (error: any) {
    console.error('查询迁移进度失败:', error)
  }
}

onUnmounted(() => {
  stopPolling()
})

// 主 Tab 切换：首次切到迁移日志时加载
watch(activeTab, tab => {
  if (tab === 'logs' && logs.value.length === 0 && logPagination.itemCount === 0) {
    fetchLogs()
  }
})

onMounted(() => {
  // 初始加载走 Redis 缓存；若服务端有进行中/已结束的迁移任务，恢复进度展示并继续轮询
  handleScan().then(() => {
    getMigrateStatus()
      .then(res => {
        if (res.data && (res.data.running || (res.data.finished && ((res.data.success?.length || 0) > 0 || (res.data.failed?.length || 0) > 0)))) {
          migrateStatus.value = {
            ...res.data,
            failed: res.data.failed || [],
            success: res.data.success || []
          }
          if (res.data.running) {
            startPolling()
          }
        }
      })
      .catch(() => {})
  })
})
</script>

<style scoped>
.image-manage-page {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  margin-bottom: 16px;
}

.tips-alert {
  margin-bottom: 16px;
}

.filter-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.filter-bar .n-tabs {
  flex: 1;
}

.keyword-input {
  width: 220px;
}

.storage-tag {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.storage-local {
  background: rgba(24, 160, 88, 0.12);
  color: #18a058;
}

.storage-oss {
  background: rgba(32, 128, 240, 0.12);
  color: #2080f0;
}

.storage-cos {
  background: rgba(128, 128, 128, 0.12);
  color: #888;
}

.storage-easyimage {
  background: rgba(240, 160, 32, 0.12);
  color: #f0a020;
}

.storage-external {
  background: rgba(208, 48, 80, 0.12);
  color: #d03050;
}

.main-tabs {
  margin-bottom: 12px;
}

.source-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 8px;
  border-radius: 4px;
  background: rgba(128, 128, 128, 0.06);
}

.delete-source-tip {
  font-size: 12px;
  color: #888;
  margin-left: 8px;
}

.result-section {
  max-height: 200px;
  overflow-y: auto;
}

.current-step-alert .current-url {
  font-size: 12px;
  word-break: break-all;
  color: #2080f0;
  margin: 4px 0;
}

.result-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 0;
  font-size: 12px;
  border-bottom: 1px dashed rgba(128, 128, 128, 0.2);
}

.result-item .result-url {
  word-break: break-all;
  flex: 1;
}

.result-item .result-url.new {
  color: #18a058;
}

@media (max-width: 1100px) {
  .image-manage-page {
    padding: 12px;
  }

  .header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .filter-bar {
    flex-direction: column;
    align-items: stretch;
  }

  .keyword-input {
    width: 100%;
  }
}
</style>
