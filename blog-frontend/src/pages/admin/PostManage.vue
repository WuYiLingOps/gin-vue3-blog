<!--
 * @ProjectName: go-vue3-blog
 * @FileName: PostManage.vue
 * @CreateTime: 2026-02-02 11:43:17
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 文章管理页面组件，提供文章的增删改查功能
 -->
<template>
  <div class="post-manage-page">
    <div class="header">
      <h1>文章管理</h1>
      <n-space>
        <n-button type="info" :size="isMobile ? 'small' : 'medium'" @click="showUploadModal = true">
          <template #icon>
            <n-icon :component="DocumentOutline" />
          </template>
          <span v-if="!isMobile">上传Markdown</span>
          <span v-else>上传</span>
        </n-button>
        <n-button type="primary" :size="isMobile ? 'small' : 'medium'" @click="showCreateModal = true">
          <template #icon>
            <n-icon :component="AddOutline" />
          </template>
          <span v-if="!isMobile">新建文章</span>
          <span v-else>新建</span>
        </n-button>
      </n-space>
    </div>

    <!-- 筛选 -->
    <n-space :vertical="isMobile" style="margin: 16px 0" :wrap="!isMobile">
      <n-input
        v-model:value="searchKeyword"
        placeholder="搜索文章..."
        clearable
        :style="{ width: isMobile ? '100%' : '250px' }"
        @keyup.enter="handleSearch"
      />
      <n-select
        v-model:value="filterCategory"
        placeholder="分类"
        clearable
        :style="{ width: isMobile ? '100%' : '150px' }"
        :options="categoryOptions"
        @update:value="handleFilterChange"
      />
      <n-select
        v-model:value="filterStatus"
        placeholder="状态"
        clearable
        :style="{ width: isMobile ? '100%' : '120px' }"
        :options="statusOptions"
        @update:value="handleFilterChange"
      />
      <n-button :block="isMobile" @click="handleSearch">搜索</n-button>
    </n-space>

    <!-- 内容区域 -->
    <div class="content-area">
      <!-- 文章列表 -->
      <n-data-table
        :columns="columns"
        :data="posts"
        :loading="loading"
        :scroll-x="isMobile ? undefined : 800"
        :single-line="false"
      />
    </div>

    <!-- 分页 - 固定在右下角 -->
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

    <!-- 上传Markdown文件对话框 -->
    <n-modal 
      v-model:show="showUploadModal" 
      preset="card" 
      title="上传Markdown文档" 
      :style="{ width: isMobile ? '95%' : '600px', maxWidth: isMobile ? '95vw' : '600px' }"
      :mask-closable="false"
      :close-on-esc="false"
    >
      <n-upload
        :file-list="markdownFileList"
        :max="1"
        accept=".md,.markdown"
        :show-file-list="true"
        @before-upload="handleBeforeUploadMarkdown"
        @remove="handleRemoveMarkdown"
      >
        <n-upload-dragger>
          <div style="margin-bottom: 12px">
            <n-icon size="48" :depth="3">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                <polyline points="14 2 14 8 20 8"></polyline>
                <line x1="16" y1="13" x2="8" y2="13"></line>
                <line x1="16" y1="17" x2="8" y2="17"></line>
                <polyline points="10 9 9 9 8 9"></polyline>
              </svg>
            </n-icon>
          </div>
          <n-text style="font-size: 16px">
            点击或拖拽文件到此区域上传
          </n-text>
          <n-p depth="3" style="margin: 8px 0 0 0">
            支持上传 .md 或 .markdown 格式的文件
          </n-p>
        </n-upload-dragger>
      </n-upload>
      <template #footer>
        <n-space justify="end">
          <n-button @click="handleCancelUpload">取消</n-button>
          <n-button type="primary" :disabled="!markdownContent" @click="handleParseMarkdown">
            解析并创建文章
          </n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- 创建/编辑文章对话框 -->
    <n-modal 
      v-model:show="showCreateModal" 
      preset="card" 
      title="创建文章" 
      :style="{ width: isMobile ? '95%' : '800px', maxWidth: isMobile ? '95vw' : '800px' }"
      :mask-closable="false"
      :close-on-esc="false"
      @close="handleModalClose"
    >
      <n-form ref="formRef" :model="formData" :rules="rules">
        <n-form-item label="标题" path="title">
          <n-input v-model:value="formData.title" placeholder="请输入文章标题" />
        </n-form-item>

        <n-form-item label="摘要" path="summary">
          <n-input
            v-model:value="formData.summary"
            type="textarea"
            :rows="3"
            placeholder="请输入文章摘要"
          />
        </n-form-item>

        <n-form-item label="封面图" path="cover">
          <n-space vertical style="width: 100%">
            <image-upload
              v-model="formData.cover"
              :width="isMobile ? 300 : 500"
              :height="isMobile ? 170 : 280"
              :max-size-m-b="5"
              alt="文章封面"
              @success="handleCoverSuccess"
            />
            <n-text depth="3" style="font-size: 12px">
              建议尺寸：1200x680 或 16:9 比例，支持 jpg、png、gif 格式（可选）
            </n-text>
          </n-space>
        </n-form-item>

        <n-form-item label="内容" path="content">
          <markdown-editor v-model="formData.content" height="400px" />
        </n-form-item>

        <n-form-item label="分类" path="category_id">
          <n-select
            v-model:value="formData.category_id"
            :options="categoryOptions"
            placeholder="选择分类"
          />
        </n-form-item>

        <n-form-item label="标签" path="tag_ids">
          <n-select
            v-model:value="formData.tag_ids"
            :options="tagOptions"
            placeholder="选择标签"
            multiple
          />
        </n-form-item>

        <n-form-item label="状态">
          <n-radio-group v-model:value="formData.status">
            <n-radio :value="0">草稿</n-radio>
            <n-radio :value="1">发布</n-radio>
          </n-radio-group>
        </n-form-item>

        <n-form-item label="可见性">
          <n-radio-group v-model:value="formData.visibility">
            <n-radio :value="1">公开</n-radio>
            <n-radio :value="0">私密（仅管理员可见）</n-radio>
          </n-radio-group>
        </n-form-item>

        <n-form-item label="置顶">
          <n-switch v-model:value="formData.is_top" />
        </n-form-item>
      </n-form>

      <template #footer>
        <n-space justify="end">
          <n-button @click="handleCancel">取消</n-button>
          <n-button type="primary" :loading="submitting" @click="handleSubmit">
            保存
          </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, h } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage, useDialog, NButton, NTag, NSpace } from 'naive-ui'
import type { DataTableColumns, FormInst, UploadFileInfo } from 'naive-ui'
import { AddOutline, DocumentOutline } from '@vicons/ionicons5'
import { getPosts, createPost, deletePost, exportPost } from '@/api/post'
import { useBlogStore, useAuthStore } from '@/stores'
import { formatDate } from '@/utils/format'
import type { Post, PostForm } from '@/types/blog'
import MarkdownEditor from '@/components/MarkdownEditor.vue'
import ImageUpload from '@/components/ImageUpload.vue'

const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const blogStore = useBlogStore()
const authStore = useAuthStore()

const loading = ref(false)
const submitting = ref(false)
const showCreateModal = ref(false)
const showUploadModal = ref(false)
const formRef = ref<FormInst | null>(null)
const posts = ref<Post[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = 10 // 固定每页显示10篇文章
const searchKeyword = ref('')
const filterCategory = ref<number | null>(null)
const filterStatus = ref<number | null>(null)
const isMobile = ref(false)

// Markdown上传相关
const markdownFileList = ref<UploadFileInfo[]>([])
const markdownContent = ref('')
const markdownFileName = ref('')

// 计算总页数
const totalPages = computed(() => Math.ceil(total.value / pageSize))

// 检测是否为「非桌面端」（移动端 + 平板）
// 小于等于 1024 视为移动/平板，只展示精简列
function checkMobile() {
  isMobile.value = window.innerWidth <= 1024
}

const formData = reactive<PostForm>({
  title: '',
  content: '',
  summary: '',
  cover: '',
  category_id: null,
  tag_ids: [],
  status: 1,
  visibility: 1,
  is_top: false
})

const statusOptions = [
  { label: '草稿', value: 0 },
  { label: '已发布', value: 1 }
]

const categoryOptions = computed(() =>
  blogStore.categories.map(c => ({ label: c.name, value: c.id }))
)

const tagOptions = computed(() => blogStore.tags.map(t => ({ label: t.name, value: t.id })))

// 桌面端完整列
const desktopColumns: DataTableColumns<Post> = [
  { 
    title: 'ID', 
    key: 'id', 
    width: 60,
    render: (_row, index) => {
      return (currentPage.value - 1) * pageSize + index + 1
    }
  },
  { title: '标题', key: 'title', ellipsis: { tooltip: true } },
  {
    title: '分类',
    key: 'category',
    width: 150,
    ellipsis: { tooltip: true },
    render: row =>
      h(
        NTag,
        {
          type: 'info',
          size: 'small'
        },
        { default: () => row.category?.name || '' }
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
        { default: () => (row.status === 1 ? '已发布' : '草稿') }
      )
  },
  {
    title: '可见性',
    key: 'visibility',
    width: 80,
    render: row =>
      h(
        NTag,
        { type: row.visibility === 1 ? 'success' : 'warning', size: 'small' },
        { default: () => (row.visibility === 1 ? '公开' : '私密') }
      )
  },
  {
    title: '浏览',
    key: 'view_count',
    width: 80
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 160,
    render: row => formatDate(row.created_at, 'YYYY-MM-DD HH:mm')
  },
  {
    title: '操作',
    key: 'actions',
    width: 220,
    render: row => {
      // 检查删除权限：普通管理员不能删除超级管理员创建的文章
      const currentUserRole = authStore.user?.role || 'user'
      const postAuthorRole = row.user?.role || 'user'
      const canDelete = currentUserRole === 'super_admin' || 
                       (currentUserRole === 'admin' && postAuthorRole !== 'super_admin') ||
                       (currentUserRole === 'user' && row.user_id === authStore.user?.id)
      
      return h(NSpace, null, {
        default: () => [
          h(
            NButton,
            { size: 'small', onClick: () => handleEdit(row.id) },
            { default: () => '编辑' }
          ),
          h(
            NButton,
            { 
              size: 'small', 
              type: 'error', 
              disabled: !canDelete,
              onClick: () => handleDelete(row.id) 
            },
            { default: () => '删除' }
          ),
          h(
            NButton,
            { size: 'small', type: 'primary', ghost: true, onClick: () => handleExport(row.id, row.title) },
            { default: () => '导出MD' }
          )
        ]
      })
    }
  }
]

// 移动端精简列：突出标题，保留分类和状态
const mobileColumns: DataTableColumns<Post> = [
  {
    title: '标题',
    key: 'title',
    ellipsis: { tooltip: true }
  },
  {
    title: '分类',
    key: 'category',
    width: 120,
    ellipsis: { tooltip: true },
    render: row =>
      h(
        NTag,
        {
          type: 'info',
          size: 'small'
        },
        { default: () => row.category?.name || '' }
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
        { default: () => (row.status === 1 ? '已发布' : '草稿') }
      )
  }
]

// 根据当前是否为移动端切换列配置
const columns = computed<DataTableColumns<Post>>(() =>
  isMobile.value ? mobileColumns : desktopColumns
)

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }],
  category_id: [
    { 
      required: true, 
      message: '请选择分类', 
      trigger: ['blur', 'change'],
      validator: (_rule: any, value: any) => {
        if (value === null || value === undefined || value === 0) {
          return new Error('请选择分类')
        }
        return true
      }
    }
  ]
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  blogStore.init()
  fetchPosts()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

async function fetchPosts() {
  try {
    loading.value = true
    const res = await getPosts({
      page: currentPage.value,
      page_size: pageSize,
      keyword: searchKeyword.value,
      category_id: filterCategory.value ?? undefined,
      status: filterStatus.value ?? undefined
    })

    if (res.data) {
      posts.value = res.data.list
      total.value = res.data.total
      // 确保页码不超过最大页数
      const maxPage = Math.ceil(total.value / pageSize) || 1
      if (currentPage.value > maxPage && maxPage > 0) {
        currentPage.value = maxPage
      }
    }
  } catch (error: any) {
    message.error(error.message || '获取文章列表失败')
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  currentPage.value = page
  fetchPosts()
}

function handleSearch() {
  currentPage.value = 1
  fetchPosts()
}

function handleFilterChange() {
  currentPage.value = 1
  fetchPosts()
}

function handleEdit(id: number) {
  router.push({ name: 'PostEdit', params: { id } })
}

function handleCoverSuccess(url: string) {
  formData.cover = url
  message.success('封面图上传成功')
}

// 处理Markdown文件上传前的验证
function handleBeforeUploadMarkdown(data: { file: UploadFileInfo }) {
  const file = data.file.file
  if (!file) {
    message.error('文件读取失败')
    return false
  }
  
  // 检查文件类型
  const fileName = file.name.toLowerCase()
  if (!fileName.endsWith('.md') && !fileName.endsWith('.markdown')) {
    message.error('只能上传 .md 或 .markdown 格式的文件')
    return false
  }
  
  // 检查文件大小（10MB）
  const maxSize = 10 * 1024 * 1024
  if (file.size > maxSize) {
    message.error('文件大小不能超过 10MB')
    return false
  }
  
  // 添加到文件列表
  markdownFileList.value = [data.file]
  
  // 读取文件内容
  const reader = new FileReader()
  reader.onload = (e) => {
    const content = e.target?.result as string
    markdownContent.value = content
    markdownFileName.value = file.name
    message.success('文件读取成功，点击"解析并创建文章"继续')
  }
  reader.onerror = () => {
    message.error('文件读取失败')
    markdownFileList.value = []
  }
  reader.readAsText(file, 'UTF-8')
  
  return false // 阻止自动上传，我们手动处理
}

// 移除Markdown文件
function handleRemoveMarkdown() {
  markdownFileList.value = []
  markdownContent.value = ''
  markdownFileName.value = ''
}

// 解析Markdown内容并打开创建文章对话框
function handleParseMarkdown() {
  if (!markdownContent.value) {
    message.warning('请先上传Markdown文件')
    return
  }
  
  try {
    // 解析Markdown内容，尝试提取front matter
    const parsed = parseMarkdownContent(markdownContent.value)
    
    // 填充表单数据
    formData.title = parsed.title || markdownFileName.value.replace(/\.(md|markdown)$/i, '')
    formData.content = parsed.content
    formData.summary = parsed.summary || ''
    formData.cover = parsed.cover || ''
    
    // 关闭上传对话框，打开创建文章对话框
    showUploadModal.value = false
    showCreateModal.value = true
    
    // 清空上传相关数据
    markdownFileList.value = []
    markdownContent.value = ''
    markdownFileName.value = ''
    
    message.success('Markdown内容已解析，请完善文章信息后保存')
  } catch (error: any) {
    message.error('解析Markdown文件失败：' + (error.message || '未知错误'))
  }
}

// 将HTML格式的img标签转换为Markdown格式
function convertHtmlImgToMarkdown(content: string): string {
  // 匹配HTML img标签，支持自闭合和普通格式
  // 匹配 <img src="url" alt="alt" ... /> 或 <img src="url" alt="alt" ...>
  const imgTagRegex = /<img\s+([^>]*?)\/?>/gi
  
  return content.replace(imgTagRegex, (match, attributes) => {
    // 提取src属性
    const srcMatch = attributes.match(/src\s*=\s*["']([^"']+)["']/i)
    if (!srcMatch) {
      return match // 如果没有src属性，保留原标签
    }
    
    const imageUrl = srcMatch[1]
    
    // 提取alt属性，如果没有则从URL中提取文件名作为alt
    const altMatch = attributes.match(/alt\s*=\s*["']([^"']*)["']/i)
    let altText = altMatch ? altMatch[1] : ''
    
    // 如果alt为空，从URL中提取文件名（不含扩展名）
    if (!altText && imageUrl) {
      const urlMatch = imageUrl.match(/\/([^\/]+)\.(jpg|jpeg|png|gif|webp|svg)$/i)
      if (urlMatch) {
        altText = urlMatch[1]
      } else {
        // 如果无法提取文件名，使用默认值
        altText = 'image'
      }
    }
    
    // 转换为Markdown格式
    return `![${altText}](${imageUrl})`
  })
}

// 解析Markdown内容，提取front matter和正文
function parseMarkdownContent(content: string) {
  const result: {
    title?: string
    summary?: string
    cover?: string
    content: string
  } = {
    content: content
  }
  
  // 先将HTML格式的img标签转换为Markdown格式
  content = convertHtmlImgToMarkdown(content)
  
  // 尝试解析YAML front matter
  const frontMatterRegex = /^---\s*\n([\s\S]*?)\n---\s*\n([\s\S]*)$/
  const match = content.match(frontMatterRegex)
  
  if (match) {
    const frontMatter = match[1]
    // 内容部分也需要转换HTML img标签
    result.content = convertHtmlImgToMarkdown(match[2])
    
    // 解析YAML front matter中的字段
    const titleMatch = frontMatter.match(/^title:\s*(.+)$/m)
    const summaryMatch = frontMatter.match(/^summary:\s*(.+)$/m)
    const coverMatch = frontMatter.match(/^cover:\s*(.+)$/m)
    
    if (titleMatch) {
      result.title = titleMatch[1].trim().replace(/^["']|["']$/g, '')
    }
    if (summaryMatch) {
      result.summary = summaryMatch[1].trim().replace(/^["']|["']$/g, '')
    }
    if (coverMatch) {
      result.cover = coverMatch[1].trim().replace(/^["']|["']$/g, '')
    }
  } else {
    // 如果没有front matter，尝试从内容中提取标题（第一个#标题）
    const titleMatch = content.match(/^#\s+(.+)$/m)
    if (titleMatch) {
      result.title = titleMatch[1].trim()
      // 移除第一个标题行，内容已经转换过HTML img标签
      result.content = content.replace(/^#\s+.+$/m, '').trim()
    } else {
      // 如果没有标题，直接使用转换后的内容
      result.content = content
    }
    
    // 尝试提取摘要（前200个字符，去除markdown语法和HTML标签）
    // 先移除代码块（包括多行和单行）
    let plainText = result.content.replace(/```[\s\S]*?```/g, '')
    // 移除行内代码
    plainText = plainText.replace(/`[^`]+`/g, '')
    // 移除所有HTML标签（包括自闭合标签如<img />和普通标签如<div></div>）
    plainText = plainText.replace(/<[^>]*>/g, '')
    // 移除markdown图片语法 ![alt](url)
    plainText = plainText.replace(/!\[([^\]]*)\]\([^\)]+\)/g, '')
    // 移除markdown语法符号（但保留必要的标点）
    plainText = plainText.replace(/[#*_\[\]()!`]/g, '')
    // 移除链接，保留文字
    plainText = plainText.replace(/\[([^\]]+)\]\([^\)]+\)/g, '$1')
    // 移除HTML实体编码（如&nbsp;等）
    plainText = plainText.replace(/&[a-z]+;/gi, ' ')
    // 移除多余的空格和换行
    plainText = plainText.replace(/\n+/g, ' ').replace(/\s+/g, ' ').trim()
    
    if (plainText.length > 0) {
      result.summary = plainText.substring(0, 200)
      if (plainText.length > 200) {
        result.summary += '...'
      }
    }
  }
  
  return result
}

// 取消上传
function handleCancelUpload() {
  markdownFileList.value = []
  markdownContent.value = ''
  markdownFileName.value = ''
  showUploadModal.value = false
}

async function handleSubmit() {
  try {
    // 先进行前端表单验证
    await formRef.value?.validate()
    
    submitting.value = true
    await createPost(formData)
    message.success('文章创建成功')
    clearForm()
    showCreateModal.value = false
    fetchPosts()
  } catch (error: any) {
    // 如果是表单验证错误，不显示错误提示（Naive UI 会自动显示）
    if (error?.errors) {
      return
    }
    
    // 处理后端返回的错误
    let errorMessage = '创建失败'
    if (error.response?.data?.message) {
      errorMessage = error.response.data.message
    } else if (error.message) {
      errorMessage = error.message
    }
    
    message.error(errorMessage)
  } finally {
    submitting.value = false
  }
}

function handleDelete(id: number) {
  // 查找要删除的文章
  const post = posts.value.find(p => p.id === id)
  if (!post) {
    message.error('文章不存在')
    return
  }

  // 前端权限检查：普通管理员不能删除超级管理员创建的文章
  const currentUserRole = authStore.user?.role || 'user'
  const postAuthorRole = post.user?.role || 'user'
  
  if (currentUserRole === 'admin' && postAuthorRole === 'super_admin') {
    message.error('普通管理员无权删除超级管理员创建的文章')
    return
  }

  dialog.warning({
    title: '确认删除',
    content: '确定要删除这篇文章吗？删除后无法恢复！',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deletePost(id)
        message.success('删除成功')
        fetchPosts()
      } catch (error: any) {
        message.error(error.message || '删除失败')
      }
    }
  })
}

// 导出 Markdown
function safeFileName(name: string | null | undefined, fallback: string) {
  // 去掉首尾引号，避免被当成内容字符
  const raw = name?.trim().replace(/^"+|"+$/g, '') || fallback
  const cleaned = raw.replace(/[\\/:*?"<>|]/g, '')
  return cleaned || fallback
}

function parseFilenameFromHeader(disposition?: string): string | null {
  if (!disposition) return null
  // 兼容 filename* 和 filename
  const utf8Match = disposition.match(/filename\\*=(?:UTF-8'')?([^;]+)/i)
  if (utf8Match && utf8Match[1]) {
    try {
      return decodeURIComponent(utf8Match[1])
    } catch {
      return utf8Match[1]
    }
  }
  const asciiMatch = disposition.match(/filename="?([^";]+)"?/i)
  return asciiMatch ? asciiMatch[1] : null
}

// 导出 Markdown
async function handleExport(id: number, title: string) {
  try {
    const res = await exportPost(id, 'md')
    const blob = new Blob([res.data], { type: 'text/markdown;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    const headerName = parseFilenameFromHeader(res.headers['content-disposition'])
    a.download = safeFileName(headerName || `${title || `post-${id}`}.md`, `post-${id}.md`)
    a.click()
    URL.revokeObjectURL(url)
    message.success('导出成功')
  } catch (error: any) {
    message.error(error?.response?.data?.message || error?.message || '导出失败')
  }
}

// 检查是否有未保存的内容
function hasUnsavedContent(): boolean {
  return !!(
    formData.title.trim() ||
    formData.content.trim() ||
    formData.summary.trim() ||
    formData.cover ||
    (formData.category_id !== null && formData.category_id !== 0) ||
    formData.tag_ids.length > 0
  )
}

// 保存为草稿
async function saveAsDraft() {
  try {
    // 验证必填字段
    if (!formData.title.trim()) {
      message.error('标题不能为空')
      return
    }
    if (!formData.content.trim()) {
      message.error('内容不能为空')
      return
    }
    if (!formData.category_id || formData.category_id === 0) {
      message.error('请选择分类')
      return
    }

    submitting.value = true
    
    // 设置为草稿状态
    const draftData = { ...formData, status: 0 }
    await createPost(draftData)
    message.success('已保存为草稿')
    clearForm()
    showCreateModal.value = false
    fetchPosts()
  } catch (error: any) {
    message.error(error.message || '保存草稿失败')
  } finally {
    submitting.value = false
  }
}

// 清空表单
function clearForm() {
  formData.title = ''
  formData.content = ''
  formData.summary = ''
  formData.cover = ''
  formData.category_id = null
  formData.tag_ids = []
  formData.status = 1
  formData.is_top = false
}

// 处理取消操作
function handleCancel() {
  // 如果有未保存的内容，弹出确认框
  if (hasUnsavedContent()) {
    dialog.warning({
      title: '提示',
      content: '检测到您有未保存的内容，是否要保存为草稿？',
      positiveText: '保存草稿',
      negativeText: '直接离开',
      onPositiveClick: async () => {
        await saveAsDraft()
      },
      onNegativeClick: () => {
        clearForm()
        showCreateModal.value = false
      }
    })
  } else {
    showCreateModal.value = false
  }
}

// 处理模态框关闭事件（点击 X 号）
function handleModalClose() {
  // 和取消按钮逻辑一样
  handleCancel()
}
</script>

<style scoped>
.post-manage-page {
  position: relative;
  /* 移除 min-height 和 padding-bottom，让页面高度自适应内容 */
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  gap: 12px;
}

.header h1 {
  margin: 0;
  font-size: 24px;
}

/* 内容区域 */
.content-area {
  position: relative;
}

/* 分页样式 - 自适应布局，跟随内容流动，保持在表格下方居右 */
.pagination-wrapper {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
  align-items: center;
  box-sizing: border-box;
}

/* 移动端样式 */
@media (max-width: 768px) {
  .header h1 {
    font-size: 20px;
  }
  
  .post-manage-page :deep(.n-data-table) {
    font-size: 13px;
  }
  
  .pagination-wrapper {
    justify-content: center; /* 手机上居中显示分页器，避免贴边 */
    margin-top: 12px;
  }
}

</style>

