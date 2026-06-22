<!--
 * 项目名称：blog-frontend
 * 文件名称：MyRevisions.vue
 * 创建时间：2026-04-28 18:30:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：编辑记录页面，普通管理员查看和管理自己提交的文章修订
-->
<template>
  <div class="my-revisions-container">
    <n-card title="编辑记录">
      <!-- 状态筛选 -->
      <div class="filter-bar">
        <n-radio-group v-model:value="statusFilter" @update:value="handleStatusChange">
          <n-radio-button value="">全部</n-radio-button>
          <n-radio-button value="pending">待审批</n-radio-button>
          <n-radio-button value="approved">已通过</n-radio-button>
          <n-radio-button value="rejected">已拒绝</n-radio-button>
          <n-radio-button value="withdrawn">已撤回</n-radio-button>
        </n-radio-group>
      </div>

      <!-- 修订列表 -->
      <n-data-table
        :columns="columns"
        :data="revisions"
        :loading="loading"
        :pagination="paginationReactive"
        :bordered="false"
      />
    </n-card>

    <!-- 对比弹窗 -->
    <n-modal
      v-model:show="showDiffModal"
      preset="card"
      title="修改对比"
      style="width: 80%; max-width: 1200px"
      :bordered="false"
    >
      <div v-if="currentDiff" class="diff-container">
        <div class="diff-header">
          <p><strong>文章：</strong>{{ currentDiff.post_title }}</p>
          <p><strong>修改说明：</strong>{{ currentRevision?.editor_comment || '无' }}</p>
          <p><strong>提交时间：</strong>{{ formatDate(currentRevision?.created_at || '') }}</p>
        </div>

        <n-divider />

        <div v-if="currentDiff.changes && currentDiff.changes.length > 0">
          <div v-for="(change, index) in currentDiff.changes" :key="index" class="change-item">
            <h4>{{ getFieldLabel(change.field) }}</h4>
            <DiffViewer
              :old-value="formatValue(change.old, change.field)"
              :new-value="formatValue(change.new, change.field)"
              :field="change.field"
            />
          </div>
        </div>
        <n-empty v-else description="无修改内容" />
      </div>

      <template #footer>
        <n-space justify="space-between">
          <n-button
            v-if="currentRevision?.status === 'pending'"
            type="warning"
            @click="handleWithdraw(currentRevision.id)"
          >
            撤回
          </n-button>
          <div v-else></div>
          <n-button @click="showDiffModal = false">关闭</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, h, computed } from 'vue'
import { NButton, NTag, NSpace, useMessage, useDialog } from 'naive-ui'
import type { DataTableColumns, PaginationProps } from 'naive-ui'
import { getMyRevisions, getRevisionDiff, withdrawRevision } from '@/api/postRevision'
import type { PostRevision } from '@/types/blog'
import DiffViewer from '@/components/DiffViewer.vue'

const message = useMessage()
const dialog = useDialog()

// 数据
const loading = ref(false)
const revisions = ref<PostRevision[]>([])
const statusFilter = ref('')
const pagination = ref({
  page: 1,
  pageSize: 15,
  total: 0
})

// 对比弹窗
const showDiffModal = ref(false)
const currentDiff = ref<any>(null)
const currentRevision = ref<PostRevision | null>(null)

// 表格列定义
const columns: DataTableColumns<PostRevision> = [
  {
    title: 'ID',
    key: 'id',
    width: 80
  },
  {
    title: '文章标题',
    key: 'post.title',
    render: row => row.post?.title || '未知文章'
  },
  {
    title: '修改说明',
    key: 'editor_comment',
    render: row => row.editor_comment || '无'
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: row => {
      return h(
        NTag,
        { type: getStatusType(row.status) },
        { default: () => getStatusText(row.status) }
      )
    }
  },
  {
    title: '提交时间',
    key: 'created_at',
    width: 180,
    render: row => formatDate(row.created_at)
  },
  {
    title: '审批人',
    key: 'reviewer',
    width: 120,
    render: row => row.reviewer?.username || '-'
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    fixed: 'right',
    render: row => {
      return h(
        NSpace,
        {},
        {
          default: () => [
            h(
              NButton,
              {
                type: 'primary',
                size: 'small',
                onClick: () => viewDiff(row)
              },
              { default: () => '查看对比' }
            ),
            row.status === 'pending'
              ? h(
                  NButton,
                  {
                    type: 'error',
                    size: 'small',
                    onClick: () => handleWithdraw(row.id)
                  },
                  { default: () => '撤回' }
                )
              : null
          ]
        }
      )
    }
  }
]

// 分页配置
const paginationReactive = computed<PaginationProps>(() => ({
  page: pagination.value.page,
  pageSize: pagination.value.pageSize,
  pageCount: Math.ceil(pagination.value.total / pagination.value.pageSize),
  showSizePicker: true,
  pageSizes: [10, 15, 20, 50],
  onChange: (page: number) => {
    pagination.value.page = page
    fetchRevisions()
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.value.pageSize = pageSize
    pagination.value.page = 1
    fetchRevisions()
  }
}))

// 获取修订列表
const fetchRevisions = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.value.page,
      page_size: pagination.value.pageSize
    }
    if (statusFilter.value) {
      params.status = statusFilter.value
    }

    const res = await getMyRevisions(params)
    revisions.value = res.data?.list || []
    pagination.value.total = res.data?.total || 0
  } catch (error: any) {
    message.error(error.message || '获取列表失败')
  } finally {
    loading.value = false
  }
}

// 状态筛选变化
const handleStatusChange = () => {
  pagination.value.page = 1
  fetchRevisions()
}

// 查看对比
const viewDiff = async (revision: PostRevision) => {
  try {
    const res = await getRevisionDiff(revision.id)
    currentDiff.value = res.data
    currentRevision.value = revision
    showDiffModal.value = true
  } catch (error: any) {
    message.error(error.message || '获取对比失败')
  }
}

// 撤回修订
const handleWithdraw = async (id: number) => {
  dialog.warning({
    title: '确认撤回',
    content: '确定要撤回此修订吗？撤回后将无法恢复。',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await withdrawRevision(id)
        message.success('撤回成功')
        showDiffModal.value = false
        fetchRevisions()
      } catch (error: any) {
        message.error(error.message || '撤回失败')
      }
    }
  })
}

// 获取状态类型
const getStatusType = (status: string): any => {
  const typeMap: Record<string, any> = {
    pending: 'warning',
    approved: 'success',
    rejected: 'error',
    withdrawn: 'default'
  }
  return typeMap[status] || 'default'
}

// 获取状态文本
const getStatusText = (status: string) => {
  const textMap: Record<string, string> = {
    pending: '待审批',
    approved: '已通过',
    rejected: '已拒绝',
    withdrawn: '已撤回'
  }
  return textMap[status] || status
}

// 获取字段标签
const getFieldLabel = (field: string) => {
  const labelMap: Record<string, string> = {
    title: '标题',
    content: '内容',
    summary: '摘要',
    cover: '封面',
    category_id: '分类',
    visibility: '可见性',
    is_top: '置顶',
    tag_ids: '标签'
  }
  return labelMap[field] || field
}

// 格式化值
const formatValue = (value: any, field: string) => {
  if (value === null || value === undefined) return '无'
  if (field === 'is_top') return value ? '是' : '否'
  if (field === 'visibility') {
    const visMap: Record<string, string> = {
      public: '公开',
      private: '私密',
      password: '密码保护'
    }
    return visMap[value] || value
  }
  if (Array.isArray(value)) return value.join(', ')
  return value
}

// 格式化日期
const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

onMounted(() => {
  fetchRevisions()
})
</script>

<style scoped>
.my-revisions-container {
  padding: 20px;
}

.filter-bar {
  margin-bottom: 20px;
}

.diff-container {
  max-height: 600px;
  overflow-y: auto;
}

.diff-header p {
  margin: 8px 0;
}

.change-item {
  margin-bottom: 24px;

  h4 {
    margin: 0 0 12px 0;
    color: #303133;
  }
}
</style>
