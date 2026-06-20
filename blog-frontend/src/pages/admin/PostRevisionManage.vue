<!--
 * @ProjectName: go-vue3-blog
 * @FileName: PostRevisionManage.vue
 * @CreateTime: 2026-04-20
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Description: 文章修订审批管理页面，超级管理员审批管理员提交的文章修改
 -->
<template>
  <div class="revision-manage-page">
    <div class="page-header">
      <h1 class="page-title">待审批修订</h1>
      <n-button-group v-if="!isMobile" size="small" class="view-toggle-group">
        <n-button :type="viewMode === 'table' ? 'primary' : 'default'" @click="viewMode = 'table'">
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
    </div>

    <!-- 内容区域 -->
    <div class="content-area">
      <div v-if="isMobile || viewMode === 'card'" class="card-list">
        <n-card
          v-for="revision in revisions"
          :key="revision.id"
          class="list-card"
          :size="isMobile ? 'small' : 'medium'"
        >
          <template #header>
            <div class="card-header-content">
              <span class="post-title">{{ revision.title }}</span>
              <n-tag type="warning" :size="isMobile ? 'tiny' : 'small'">待审批</n-tag>
            </div>
          </template>
          <div class="card-content">
            <div class="info-item">
              <span class="label">编辑者：</span>
              <span class="value">{{ revision.editor?.nickname || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">提交时间：</span>
              <span class="value">{{ formatDate(revision.created_at, 'YYYY-MM-DD HH:mm') }}</span>
            </div>
            <div class="info-item">
              <span class="label">修改说明：</span>
              <span class="value">{{ revision.editor_comment || '无' }}</span>
            </div>
          </div>
          <template #footer>
            <n-space justify="end" :size="isMobile ? 'small' : 'medium'">
              <n-button :size="isMobile ? 'tiny' : 'small'" @click="handleViewDiff(revision.id)">
                查看对比
              </n-button>
              <n-button
                v-if="canWithdraw(revision)"
                :size="isMobile ? 'tiny' : 'small'"
                type="warning"
                @click="handleWithdraw(revision.id)"
              >
                撤回
              </n-button>
              <n-button
                :size="isMobile ? 'tiny' : 'small'"
                type="success"
                @click="handleApprove(revision.id)"
              >
                通过
              </n-button>
              <n-button
                :size="isMobile ? 'tiny' : 'small'"
                type="error"
                @click="handleReject(revision.id)"
              >
                拒绝
              </n-button>
            </n-space>
          </template>
        </n-card>
      </div>

      <n-data-table
        v-else-if="viewMode === 'table'"
        :columns="columns"
        :data="revisions"
        :loading="loading"
        :single-line="false"
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
    </div>

    <!-- 对比弹窗 -->
    <n-modal
      v-model:show="showDiffModal"
      preset="card"
      title="修改对比"
      style="width: 90%; max-width: 1200px"
    >
      <div v-if="diffData" class="diff-container">
        <n-space vertical :size="16">
          <!-- 基本信息 -->
          <n-card title="基本信息" size="small">
            <n-descriptions :column="2" bordered>
              <n-descriptions-item label="文章标题">{{ diffData.post_title }}</n-descriptions-item>
              <n-descriptions-item label="修改者">{{
                diffData.editor.username
              }}</n-descriptions-item>
              <n-descriptions-item label="修改时间">{{
                formatDate(diffData.created_at)
              }}</n-descriptions-item>
              <n-descriptions-item label="修改说明">{{
                diffData.editor_comment || '无'
              }}</n-descriptions-item>
            </n-descriptions>
          </n-card>

          <!-- 修改内容 -->
          <n-card title="修改内容" size="small">
            <n-alert v-if="diffData.changes_count === 0" type="info"> 无修改内容 </n-alert>
            <n-space v-else vertical :size="16">
              <div v-for="change in diffData.changes" :key="change.field" class="change-item">
                <n-divider title-placement="left">
                  <n-tag type="info">{{ getFieldLabel(change.field) }}</n-tag>
                </n-divider>
                <DiffViewer
                  :old-value="formatValue(change.field, change.old)"
                  :new-value="formatValue(change.field, change.new)"
                  :field="change.field"
                />
              </div>
            </n-space>
          </n-card>
        </n-space>
      </div>
      <template #footer>
        <n-space justify="space-between">
          <n-button
            v-if="diffData && diffData.editor?.id === authStore.user?.id"
            type="warning"
            @click="handleWithdraw(currentRevisionId)"
          >
            撤回
          </n-button>
          <div v-else></div>
          <n-space>
            <n-button @click="showDiffModal = false">关闭</n-button>
            <n-button type="success" @click="handleApproveFromDiff">通过</n-button>
            <n-button type="error" @click="handleRejectFromDiff">拒绝</n-button>
          </n-space>
        </n-space>
      </template>
    </n-modal>

    <!-- 拒绝原因弹窗 -->
    <n-modal v-model:show="showRejectModal" preset="dialog" title="拒绝修订">
      <n-form>
        <n-form-item label="拒绝原因">
          <n-input
            v-model:value="rejectReason"
            type="textarea"
            placeholder="请输入拒绝原因（选填）"
            :rows="4"
          />
        </n-form-item>
      </n-form>
      <template #action>
        <n-space justify="end">
          <n-button @click="showRejectModal = false">取消</n-button>
          <n-button type="error" @click="confirmReject">确认拒绝</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, h, watch } from 'vue'
import { useMessage, useDialog, NButton, NIcon, NTag, NSpace, NEllipsis } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { GridOutline, AppsOutline } from '@vicons/ionicons5'
import {
  getPendingRevisions,
  getRevisionDiff,
  approveRevision,
  rejectRevision,
  withdrawRevision
} from '@/api/postRevision'
import { useAuthStore } from '@/stores'
import { formatDate } from '@/utils/format'
import type { PostRevision, RevisionDiff } from '@/types/blog'
import DiffViewer from '@/components/DiffViewer.vue'

const message = useMessage()
const dialog = useDialog()
const authStore = useAuthStore()

const loading = ref(false)
const revisions = ref<PostRevision[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = 15
const isMobile = ref(false)
const viewMode = ref<'table' | 'card'>('table')

const showDiffModal = ref(false)
const diffData = ref<RevisionDiff | null>(null)
const currentRevisionId = ref<number>(0)

const showRejectModal = ref(false)
const rejectReason = ref('')
const pendingRejectId = ref<number>(0)

function checkMobile() {
  isMobile.value = window.innerWidth <= 1100
}

const totalPages = computed(() => Math.ceil(total.value / pageSize))

const columns: DataTableColumns<PostRevision> = [
  {
    title: 'ID',
    key: 'id',
    width: 60,
    render: (_row, index) => {
      return (currentPage.value - 1) * pageSize + index + 1
    }
  },
  {
    title: '文章标题',
    key: 'title',
    ellipsis: {
      tooltip: true
    },
    render: row =>
      h(NEllipsis, { style: 'max-width: 250px' }, { default: () => row.post?.title || '-' })
  },
  {
    title: '编辑者',
    key: 'editor',
    width: 120,
    render: row => row.editor?.nickname || '-'
  },
  {
    title: '修改说明',
    key: 'editor_comment',
    width: 200,
    ellipsis: { tooltip: true },
    render: row => row.editor_comment || '无'
  },
  {
    title: '提交时间',
    key: 'created_at',
    width: 160,
    render: row => formatDate(row.created_at, 'YYYY-MM-DD HH:mm')
  },
  {
    title: '操作',
    key: 'actions',
    width: 240,
    render: row =>
      h(NSpace, null, {
        default: () =>
          [
            h(
              NButton,
              {
                size: 'small',
                onClick: () => handleViewDiff(row.id)
              },
              { default: () => '查看对比' }
            ),
            canWithdraw(row)
              ? h(
                  NButton,
                  {
                    size: 'small',
                    type: 'warning',
                    onClick: () => handleWithdraw(row.id)
                  },
                  { default: () => '撤回' }
                )
              : null,
            h(
              NButton,
              {
                size: 'small',
                type: 'success',
                onClick: () => handleApprove(row.id)
              },
              { default: () => '通过' }
            ),
            h(
              NButton,
              {
                size: 'small',
                type: 'error',
                onClick: () => handleReject(row.id)
              },
              { default: () => '拒绝' }
            )
          ].filter(Boolean)
      })
  }
]

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)

  const savedViewMode = localStorage.getItem('revision-manage-view-mode')
  if (savedViewMode === 'card' || savedViewMode === 'table') {
    viewMode.value = savedViewMode
  }

  fetchRevisions()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

watch(viewMode, newMode => {
  localStorage.setItem('revision-manage-view-mode', newMode)
})

async function fetchRevisions() {
  try {
    loading.value = true
    const res = await getPendingRevisions({
      page: currentPage.value,
      page_size: pageSize
    })

    if (res.data) {
      revisions.value = res.data.list
      total.value = res.data.total
      const maxPage = Math.ceil(total.value / pageSize) || 1
      if (currentPage.value > maxPage && maxPage > 0) {
        currentPage.value = maxPage
      }
    }
  } catch (error: any) {
    message.error(error.message || '获取待审批列表失败')
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  currentPage.value = page
  fetchRevisions()
}

async function handleViewDiff(id: number) {
  try {
    const res = await getRevisionDiff(id)
    if (res.data) {
      diffData.value = res.data
      currentRevisionId.value = id
      showDiffModal.value = true
    }
  } catch (error: any) {
    message.error(error.message || '获取修改对比失败')
  }
}

async function handleApprove(id: number) {
  dialog.warning({
    title: '确认通过',
    content: '确定要通过这个修订吗？修改将立即应用到文章。',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await approveRevision(id)
        message.success('审批通过')
        fetchRevisions()
        showDiffModal.value = false
      } catch (error: any) {
        message.error(error.message || '审批失败')
      }
    }
  })
}

async function handleReject(id: number) {
  pendingRejectId.value = id
  rejectReason.value = ''
  showRejectModal.value = true
}

async function confirmReject() {
  try {
    await rejectRevision(pendingRejectId.value, { reject_reason: rejectReason.value })
    message.success('已拒绝修订')
    showRejectModal.value = false
    showDiffModal.value = false
    fetchRevisions()
  } catch (error: any) {
    message.error(error.message || '拒绝失败')
  }
}

function handleApproveFromDiff() {
  handleApprove(currentRevisionId.value)
}

function handleRejectFromDiff() {
  handleReject(currentRevisionId.value)
}

// 撤回修订
async function handleWithdraw(id: number) {
  dialog.warning({
    title: '确认撤回',
    content: '确定要撤回这个修订吗？撤回后将无法恢复。',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await withdrawRevision(id)
        message.success('撤回成功')
        fetchRevisions()
      } catch (error: any) {
        message.error(error.message || '撤回失败')
      }
    }
  })
}

// 判断是否可以撤回（只有创建者可以撤回）
function canWithdraw(revision: PostRevision): boolean {
  return revision.editor_id === authStore.user?.id
}

// 获取字段标签
function getFieldLabel(field: string): string {
  const labels: Record<string, string> = {
    title: '标题',
    content: '内容',
    summary: '摘要',
    cover: '封面',
    category_id: '分类',
    tag_ids: '标签',
    visibility: '可见性',
    is_top: '置顶状态'
  }
  return labels[field] || field
}

// 格式化值
function formatValue(field: string, value: any): string {
  if (value === null || value === undefined) {
    return '无'
  }

  if (field === 'visibility') {
    return value === 1 ? '公开' : '私密'
  }

  if (field === 'is_top') {
    return value ? '是' : '否'
  }

  if (field === 'tag_ids' && Array.isArray(value)) {
    return value.join(', ')
  }

  // 不再截断内容，让 DiffViewer 组件处理
  return String(value)
}
</script>

<style scoped lang="scss">
.revision-manage-page {
  padding: 20px;

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    .page-title {
      font-size: 24px;
      font-weight: 600;
      margin: 0;
    }
  }

  .content-area {
    background: var(--card-color);
    border-radius: 8px;
    padding: 20px;

    .card-list {
      display: grid;
      gap: 16px;

      .list-card {
        .card-header-content {
          display: flex;
          justify-content: space-between;
          align-items: center;

          .post-title {
            font-weight: 500;
            flex: 1;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }
        }

        .card-content {
          .info-item {
            margin-bottom: 8px;

            .label {
              color: var(--text-color-3);
              margin-right: 8px;
            }

            .value {
              color: var(--text-color-1);
            }
          }
        }
      }
    }

    .pagination-wrapper {
      display: flex;
      justify-content: flex-end;
      margin-top: 20px;
    }
  }
}

.diff-container {
  .change-item {
    margin-bottom: 16px;
  }
}

@media (max-width: 1100px) {
  .revision-manage-page {
    padding: 12px;

    .page-header {
      .page-title {
        font-size: 20px;
      }
    }

    .content-area {
      padding: 12px;
    }
  }
}
</style>
