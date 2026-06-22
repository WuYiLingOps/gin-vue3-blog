<!--
 * @ProjectName: go-vue3-blog
 * @FileName: PushHistoryTab.vue
 * @CreateTime: 2026-04-28 13:30:00
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 推送历史记录 Tab 组件
 -->
<template>
  <div class="push-history-tab">
    <!-- 统计卡片 -->
    <n-grid :cols="isMobile ? 2 : 4" :x-gap="16" :y-gap="16" style="margin-bottom: 16px">
      <n-gi>
        <n-statistic label="总推送次数" :value="stats.total_count">
          <template #prefix>
            <n-icon :component="SendOutline" />
          </template>
        </n-statistic>
      </n-gi>
      <n-gi>
        <n-statistic label="成功推送" :value="stats.total_success">
          <template #prefix>
            <n-icon :component="CheckmarkCircleOutline" color="#18a058" />
          </template>
        </n-statistic>
      </n-gi>
      <n-gi>
        <n-statistic label="失败推送" :value="stats.total_failed">
          <template #prefix>
            <n-icon :component="CloseCircleOutline" color="#d03050" />
          </template>
        </n-statistic>
      </n-gi>
      <n-gi>
        <n-statistic label="最近推送">
          <template #prefix>
            <n-icon :component="TimeOutline" />
          </template>
          {{ stats.last_push_at ? formatDate(stats.last_push_at, 'MM-DD HH:mm') : '暂无' }}
        </n-statistic>
      </n-gi>
    </n-grid>

    <!-- 推送历史列表 -->
    <n-card title="推送历史记录" :bordered="false">
      <template #header-extra>
        <!-- 视图切换按钮（仅桌面端显示） -->
        <n-button-group v-if="!isMobile" size="small" class="view-toggle-group">
          <n-button
            :type="viewMode === 'table' ? 'primary' : 'default'"
            @click="viewMode = 'table'"
          >
            <template #icon>
              <n-icon :component="GridOutline" />
            </template>
            表格
          </n-button>
          <n-button :type="viewMode === 'card' ? 'primary' : 'default'" @click="viewMode = 'card'">
            <template #icon>
              <n-icon :component="AppsOutline" />
            </template>
            卡片
          </n-button>
        </n-button-group>
      </template>

      <!-- 卡片视图 -->
      <div v-if="isMobile || viewMode === 'card'" class="card-list">
        <n-card
          v-for="history in histories"
          :key="history.id"
          class="list-card"
          :size="isMobile ? 'small' : 'medium'"
        >
          <template #header>
            <div class="card-header-content">
              <span class="post-title">{{ history.post_title }}</span>
              <n-tag :type="getStatusType(history.status)" :size="isMobile ? 'tiny' : 'small'">
                {{ getStatusText(history.status) }}
              </n-tag>
            </div>
          </template>
          <div class="card-content">
            <div class="info-item">
              <span class="label">推送统计：</span>
              <span class="value">
                成功 <n-text type="success">{{ history.success_count }}</n-text> / 失败
                <n-text type="error">{{ history.failed_count }}</n-text> / 总计
                {{ history.total_count }}
              </span>
            </div>
            <div class="info-item">
              <span class="label">开始时间：</span>
              <span class="value">{{ formatDate(history.started_at, 'YYYY-MM-DD HH:mm') }}</span>
            </div>
            <div v-if="history.completed_at" class="info-item">
              <span class="label">完成时间：</span>
              <span class="value">{{ formatDate(history.completed_at, 'YYYY-MM-DD HH:mm') }}</span>
            </div>
          </div>
          <template #footer>
            <n-space justify="end" :size="isMobile ? 'small' : 'medium'">
              <n-button :size="isMobile ? 'tiny' : 'small'" @click="handleViewDetail(history)">
                查看详情
              </n-button>
              <n-button
                :size="isMobile ? 'tiny' : 'small'"
                type="error"
                @click="handleDelete(history.id)"
              >
                删除
              </n-button>
            </n-space>
          </template>
        </n-card>
      </div>

      <!-- 表格视图 -->
      <n-data-table
        v-else
        :columns="columns"
        :data="histories"
        :loading="loading"
        :pagination="pagination"
        :single-line="false"
        @update:page="handlePageChange"
      />
    </n-card>

    <!-- 推送详情抽屉 -->
    <n-drawer
      v-model:show="showDetail"
      :width="isMobile ? '100%' : 600"
      placement="right"
      :closable="true"
    >
      <n-drawer-content :title="`推送详情 - ${currentHistory?.post_title}`" :closable="true">
        <div v-if="currentHistory" class="detail-content">
          <!-- 推送概览 -->
          <n-descriptions :column="1" bordered>
            <n-descriptions-item label="文章标题">
              {{ currentHistory.post_title }}
            </n-descriptions-item>
            <n-descriptions-item label="推送状态">
              <n-tag :type="getStatusType(currentHistory.status)">
                {{ getStatusText(currentHistory.status) }}
              </n-tag>
            </n-descriptions-item>
            <n-descriptions-item label="推送统计">
              成功 {{ currentHistory.success_count }} / 失败 {{ currentHistory.failed_count }} /
              总计 {{ currentHistory.total_count }}
            </n-descriptions-item>
            <n-descriptions-item label="开始时间">
              {{ formatDate(currentHistory.started_at, 'YYYY-MM-DD HH:mm:ss') }}
            </n-descriptions-item>
            <n-descriptions-item label="完成时间">
              {{
                currentHistory.completed_at
                  ? formatDate(currentHistory.completed_at, 'YYYY-MM-DD HH:mm:ss')
                  : '进行中'
              }}
            </n-descriptions-item>
          </n-descriptions>

          <!-- 推送详情列表 -->
          <n-divider />
          <h3>推送详情</h3>
          <n-data-table
            :columns="detailColumns"
            :data="details"
            :loading="detailLoading"
            :pagination="detailPagination"
            size="small"
            @update:page="handleDetailPageChange"
          />
        </div>
      </n-drawer-content>
    </n-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted, h, watch } from 'vue'
import { useMessage, useDialog, NButton, NButtonGroup, NTag, NIcon } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import {
  SendOutline,
  CheckmarkCircleOutline,
  CloseCircleOutline,
  TimeOutline,
  EyeOutline,
  TrashOutline,
  GridOutline,
  AppsOutline
} from '@vicons/ionicons5'
import { formatDate } from '@/utils/format'
import {
  getPushHistories,
  getPushHistoryDetail,
  deletePushHistory,
  getPushStats
} from '@/api/subscribe'
import type { PushHistory, PushDetail, PushStats } from '@/api/subscribe'

const message = useMessage()
const dialog = useDialog()
const isMobile = ref(window.innerWidth < 768)
const viewMode = ref<'table' | 'card'>('table')

// 统计数据
const stats = reactive<PushStats>({
  total_count: 0,
  total_success: 0,
  total_failed: 0,
  last_push_at: ''
})

// 推送历史列表
const loading = ref(false)
const histories = ref<PushHistory[]>([])
const pagination = reactive({
  page: 1,
  pageSize: 10,
  pageCount: 1,
  showSizePicker: true,
  pageSizes: [10, 20, 50]
})

// 推送详情
const showDetail = ref(false)
const detailLoading = ref(false)
const currentHistory = ref<PushHistory | null>(null)
const details = ref<PushDetail[]>([])
const detailPagination = reactive({
  page: 1,
  pageSize: 20,
  pageCount: 1
})

// 表格列定义
const columns: DataTableColumns<PushHistory> = [
  {
    title: '文章标题',
    key: 'post_title',
    ellipsis: { tooltip: true }
  },
  {
    title: '推送状态',
    key: 'status',
    width: 100,
    render: row =>
      h(NTag, { type: getStatusType(row.status) }, { default: () => getStatusText(row.status) })
  },
  {
    title: '推送统计',
    key: 'stats',
    width: 150,
    render: row => `${row.success_count}/${row.failed_count}/${row.total_count}`
  },
  {
    title: '开始时间',
    key: 'started_at',
    width: 160,
    render: row => formatDate(row.started_at, 'YYYY-MM-DD HH:mm')
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render: row =>
      h('div', { style: 'display: flex; gap: 8px;' }, [
        h(
          NButton,
          {
            size: 'small',
            onClick: () => handleViewDetail(row)
          },
          { default: () => '查看详情', icon: () => h(NIcon, { component: EyeOutline }) }
        ),
        h(
          NButton,
          {
            size: 'small',
            type: 'error',
            onClick: () => handleDelete(row.id)
          },
          { icon: () => h(NIcon, { component: TrashOutline }) }
        )
      ])
  }
]

// 详情表格列定义
const detailColumns: DataTableColumns<PushDetail> = [
  {
    title: '邮箱',
    key: 'subscriber_email',
    ellipsis: { tooltip: true }
  },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: row =>
      h(
        NTag,
        {
          type: row.status === 1 ? 'success' : row.status === 2 ? 'error' : 'default',
          size: 'small'
        },
        { default: () => (row.status === 1 ? '成功' : row.status === 2 ? '失败' : '待发送') }
      )
  },
  {
    title: '发送时间',
    key: 'sent_at',
    width: 140,
    render: row => (row.sent_at ? formatDate(row.sent_at, 'MM-DD HH:mm:ss') : '-')
  },
  {
    title: '错误信息',
    key: 'error_message',
    ellipsis: { tooltip: true },
    render: row => row.error_message || '-'
  }
]

// 获取状态类型
const getStatusType = (status: number) => {
  return status === 1 ? 'success' : status === 2 ? 'warning' : 'default'
}

// 获取状态文本
const getStatusText = (status: number) => {
  return status === 1 ? '已完成' : status === 2 ? '部分失败' : '进行中'
}

// 获取统计数据
const fetchStats = async () => {
  try {
    const res = await getPushStats()
    Object.assign(stats, res.data)
  } catch (error: any) {
    message.error(error?.message || '获取统计数据失败')
  }
}

// 获取推送历史列表
const fetchHistories = async () => {
  try {
    loading.value = true
    const res = await getPushHistories(pagination.page, pagination.pageSize)
    histories.value = res.data?.list || []
    pagination.pageCount = Math.ceil((res.data?.total || 0) / pagination.pageSize)
  } catch (error: any) {
    message.error(error?.message || '获取推送历史失败')
  } finally {
    loading.value = false
  }
}

// 查看详情
const handleViewDetail = async (history: PushHistory) => {
  currentHistory.value = history
  showDetail.value = true
  detailPagination.page = 1
  await fetchDetails()
}

// 获取推送详情
const fetchDetails = async () => {
  if (!currentHistory.value) return

  try {
    detailLoading.value = true
    const res = await getPushHistoryDetail(
      currentHistory.value.id,
      detailPagination.page,
      detailPagination.pageSize
    )
    details.value = res.data?.details || []
    detailPagination.pageCount = Math.ceil((res.data?.total || 0) / detailPagination.pageSize)
  } catch (error: any) {
    message.error(error?.message || '获取推送详情失败')
  } finally {
    detailLoading.value = false
  }
}

// 删除推送历史
const handleDelete = (id: number) => {
  dialog.warning({
    title: '确认删除',
    content: '确定要删除这条推送历史记录吗？此操作不可恢复。',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deletePushHistory(id)
        message.success('删除成功')
        await fetchHistories()
        await fetchStats()
      } catch (error: any) {
        message.error(error?.message || '删除失败')
      }
    }
  })
}

// 分页变化
const handlePageChange = (page: number) => {
  pagination.page = page
  fetchHistories()
}

// 详情分页变化
const handleDetailPageChange = (page: number) => {
  detailPagination.page = page
  fetchDetails()
}

// 检测屏幕尺寸
const checkMobile = () => {
  isMobile.value = window.innerWidth < 768
}

// 监听视图模式变化，保存到 localStorage
watch(viewMode, newMode => {
  localStorage.setItem('push-history-view-mode', newMode)
})

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)

  // 从本地存储恢复视图模式
  const savedMode = localStorage.getItem('push-history-view-mode')
  if (savedMode === 'card' || savedMode === 'table') {
    viewMode.value = savedMode
  }

  fetchStats()
  fetchHistories()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped lang="scss">
.push-history-tab {
  .detail-content {
    h3 {
      margin: 16px 0;
      font-size: 16px;
      font-weight: 600;
    }
  }

  .view-toggle-group {
    flex-shrink: 0;
  }

  .card-list {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
    padding: 8px 0;
  }

  @media (max-width: 1100px) {
    .card-list {
      grid-template-columns: 1fr;
    }
  }

  .list-card {
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
    transition: all 0.3s ease;
  }

  .list-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  }

  .card-header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 8px;
  }

  .post-title {
    font-weight: 500;
    font-size: 14px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    flex: 1;
    min-width: 0;
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
    width: 65px;
    flex-shrink: 0;
  }

  .info-item .value {
    color: #555;
    flex: 1;
  }
}

html.dark .push-history-tab {
  .list-card {
    :deep(.n-card__content) {
      background: rgba(30, 41, 59, 0.85);
    }
  }
}

// 移动端抽屉样式优化
@media (max-width: 768px) {
  :deep(.n-drawer) {
    .n-drawer-content {
      .n-drawer-body-content-wrapper {
        padding: 16px;
      }
    }

    .n-drawer-header {
      .n-drawer-header__main {
        flex: 1;
        min-width: 0;
      }

      .n-drawer-header__title {
        font-size: 16px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
  }
}
</style>
