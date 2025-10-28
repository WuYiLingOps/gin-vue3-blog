<template>
  <div class="post-detail-container">
    <n-spin :show="loading" class="post-detail-page">
      <!-- 文章主体 -->
      <div class="post-main">
        <n-card v-if="post">
        <!-- 文章头部 -->
        <div class="post-header">
          <h1 class="post-title">{{ post.title }}</h1>
          <div class="post-meta">
            <n-space>
              <n-avatar :src="post.user.avatar" round />
              <span>{{ post.user.nickname }}</span>
              <n-divider vertical />
              <span>{{ formatDate(post.created_at, 'YYYY-MM-DD HH:mm') }}</span>
              <n-divider vertical />
              <span>
                <n-icon :component="EyeOutline" />
                {{ post.view_count }}
              </span>
              <n-divider vertical />
              <span>
                <n-icon :component="HeartOutline" />
                {{ post.like_count }}
              </span>
            </n-space>
          </div>
          <div class="post-tags">
            <n-space>
              <n-tag :bordered="false" type="info">{{ post.category.name }}</n-tag>
              <n-tag v-for="tag in post.tags" :key="tag.id" :bordered="false">
                {{ tag.name }}
              </n-tag>
            </n-space>
          </div>
        </div>

        <n-divider />

        <!-- 文章内容 -->
        <div class="post-content">
          <markdown-preview :content="post.content" />
        </div>

        <n-divider />

        <!-- 文章操作 -->
        <div class="post-actions">
          <n-space justify="center">
            <n-button :type="liked ? 'primary' : 'default'" @click="handleLike">
              <template #icon>
                <n-icon :component="liked ? Heart : HeartOutline" />
              </template>
              {{ post.like_count }}
            </n-button>
            <n-button v-if="canEdit" @click="handleEdit">
              <template #icon>
                <n-icon :component="CreateOutline" />
              </template>
              编辑
            </n-button>
          </n-space>
        </div>

        <n-divider />

        <!-- 评论区 -->
        <div class="comments-section">
          <h3>评论 ({{ comments.length }})</h3>

          <!-- 评论表单 -->
          <n-card v-if="authStore.isLoggedIn" class="comment-form">
            <!-- 回复提示 -->
            <n-alert
              v-if="replyToComment"
              type="info"
              closable
              style="margin-bottom: 12px"
              @close="replyToComment = null; replyToUser = null; commentContent = ''"
            >
              正在回复 <strong>@{{ (replyToUser || replyToComment).user.nickname }}</strong> 的评论
            </n-alert>
            
            <n-input
              v-model:value="commentContent"
              type="textarea"
              placeholder="写下你的评论..."
              :rows="3"
            />
            <div style="margin-top: 12px; text-align: right">
              <n-button type="primary" :loading="submitting" @click="handleSubmitComment">
                {{ replyToComment ? '发表回复' : '发表评论' }}
              </n-button>
            </div>
          </n-card>

          <n-alert v-else type="info" style="margin-bottom: 16px">
            请 <n-button text type="primary" @click="router.push('/auth/login')">登录</n-button> 后发表评论
          </n-alert>

          <!-- 评论列表 -->
          <div class="comments-list">
            <div v-for="comment in comments" :key="comment.id" class="comment-item">
              <n-space align="start">
                <n-avatar :src="comment.user.avatar" round />
                <div class="comment-content">
                  <div class="comment-header">
                    <strong>{{ comment.user.nickname }}</strong>
                    <span class="comment-time">{{ formatRelativeTime(comment.created_at) }}</span>
                  </div>
                  <p>{{ comment.content }}</p>
                  <div class="comment-actions">
                    <n-button text size="small" @click="handleReply(comment)">回复</n-button>
                    <n-button 
                      v-if="comment.children && comment.children.length > 0"
                      text 
                      size="small" 
                      @click="toggleExpand(comment.id)"
                    >
                      {{ expandedComments.has(comment.id) ? '收起' : `展开 ${comment.children.length} 条回复` }}
                    </n-button>
                    <n-popconfirm
                      v-if="canDeleteComment(comment)"
                      @positive-click="handleDeleteComment(comment.id)"
                    >
                      <template #trigger>
                        <n-button text size="small" type="error">删除</n-button>
                      </template>
                      确定要删除这条评论吗？
                    </n-popconfirm>
                  </div>

                  <!-- 子评论 -->
                  <div v-if="comment.children && comment.children.length > 0 && expandedComments.has(comment.id)" class="reply-list">
                    <div
                      v-for="reply in comment.children"
                      :key="reply.id"
                      class="reply-item"
                    >
                      <n-space align="start">
                        <n-avatar :src="reply.user.avatar" round size="small" />
                        <div class="reply-content">
                          <div class="reply-header">
                            <strong>{{ reply.user.nickname }}</strong>
                            <span class="reply-to">回复 @{{ getReplyTargetName(reply, comment) }}</span>
                            <span class="comment-time">{{
                              formatRelativeTime(reply.created_at)
                            }}</span>
                          </div>
                          <p>{{ removeAtMention(reply.content) }}</p>
                          <div class="comment-actions">
                            <n-button text size="small" @click="handleReply(comment, reply)">回复</n-button>
                            <n-popconfirm
                              v-if="canDeleteComment(reply)"
                              @positive-click="handleDeleteComment(reply.id)"
                            >
                              <template #trigger>
                                <n-button text size="small" type="error">删除</n-button>
                              </template>
                              确定要删除这条回复吗？
                            </n-popconfirm>
                          </div>
                        </div>
                      </n-space>
                    </div>
                  </div>
                </div>
              </n-space>
            </div>
          </div>

          <n-empty v-if="comments.length === 0" description="暂无评论" />
        </div>
      </n-card>
      </div>
    </n-spin>

    <!-- 右侧目录 - 独立于 spin -->
    <div class="post-toc" v-if="post && tocItems.length > 0">
      <n-card size="small">
        <template #header>
          <div class="toc-header">
            <span class="toc-title">目录</span>
            <span class="toc-progress-text">{{ readingProgress }}%</span>
          </div>
          <!-- 阅读进度条 -->
          <div class="reading-progress-bar">
            <div class="reading-progress-fill" :style="{ width: `${readingProgress}%` }"></div>
          </div>
        </template>
        <div class="toc-list">
          <a
            v-for="item in tocItems"
            :key="item.id"
            :class="['toc-item', `toc-level-${item.level}`, { active: activeHeading === item.id }]"
            :href="`#${item.id}`"
            @click.prevent="scrollToHeading(item.id)"
          >
            {{ item.text }}
          </a>
        </div>
      </n-card>
    </div>

    <!-- 返回顶部按钮 -->
    <n-back-top v-if="scrollContainer" :right="40" :bottom="80" :listen-to="scrollContainer" class="back-to-top">
      <div class="back-top-button">
        <n-icon size="24" :component="ArrowUpOutline" />
      </div>
    </n-back-top>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick, onBeforeUnmount } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useMessage } from 'naive-ui'
import {
  EyeOutline,
  HeartOutline,
  Heart,
  CreateOutline,
  ArrowUpOutline
} from '@vicons/ionicons5'
import { getPostById, likePost } from '@/api/post'
import { getCommentsByPostId, createComment, deleteComment } from '@/api/comment'
import { formatDate, formatRelativeTime } from '@/utils/format'
import { useAuthStore } from '@/stores'
import type { Post, Comment } from '@/types/blog'
import MarkdownPreview from '@/components/MarkdownPreview.vue'

const router = useRouter()
const route = useRoute()
const message = useMessage()
const authStore = useAuthStore()

const loading = ref(false)
const submitting = ref(false)
const post = ref<Post | null>(null)
const comments = ref<Comment[]>([])
const commentContent = ref('')
const liked = ref(false)
const scrollContainer = ref<HTMLElement | null>(null)
const replyToComment = ref<Comment | null>(null) // 记录正在回复的主评论
const replyToUser = ref<Comment | null>(null) // 记录正在回复的具体用户
const expandedComments = ref<Set<number>>(new Set()) // 记录展开的评论ID

// TOC 相关
interface TocItem {
  id: string
  text: string
  level: number
}

const tocItems = ref<TocItem[]>([])
const activeHeading = ref('')
const readingProgress = ref(0)

const postId = computed(() => Number(route.params.id))

const canEdit = computed(() => {
  if (!post.value || !authStore.isLoggedIn) return false
  return authStore.isAdmin || post.value.user_id === authStore.user?.id
})

onMounted(() => {
  fetchPost()
  fetchComments()
  
  // 监听滚动 - 需要监听 n-layout-content 的滚动
  nextTick(() => {
    // 尝试多个可能的滚动容器
    let container = document.querySelector('.main-content .n-scrollbar-container') as HTMLElement
    if (!container) {
      container = document.querySelector('.main-content') as HTMLElement
    }
    if (!container) {
      container = document.querySelector('.n-layout-scroll-container') as HTMLElement
    }
    
    scrollContainer.value = container
    
    if (container) {
      container.addEventListener('scroll', handleScroll)
    }
  })
})

onBeforeUnmount(() => {
  if (scrollContainer.value) {
    scrollContainer.value.removeEventListener('scroll', handleScroll)
  }
})

async function fetchPost() {
  try {
    loading.value = true
    const res = await getPostById(postId.value)
    if (res.data) {
      post.value = res.data
      liked.value = res.data.liked || false
      // 等待 DOM 更新后生成目录
      nextTick(() => {
        generateToc()
      })
    }
  } catch (error: any) {
    message.error(error.message || '获取文章失败')
    router.push('/')
  } finally {
    loading.value = false
  }
}

async function fetchComments() {
  try {
    const res = await getCommentsByPostId(postId.value)
    if (res.data) {
      comments.value = res.data
    }
  } catch (error: any) {
    console.error('获取评论失败:', error)
  }
}

async function handleLike() {
  if (!authStore.isLoggedIn) {
    message.warning('请先登录')
    return
  }

  try {
    const res = await likePost(postId.value)
    if (res.data) {
      const isLiked = res.data.liked
      liked.value = isLiked
      if (post.value) {
        // 根据返回的状态更新点赞数
        if (isLiked) {
          post.value.like_count++
          message.success('点赞成功')
        } else {
          post.value.like_count--
          message.success('取消点赞')
        }
      }
    }
  } catch (error: any) {
    message.error(error.message || '操作失败')
  }
}

async function handleSubmitComment() {
  if (!commentContent.value.trim()) {
    message.warning('请输入评论内容')
    return
  }

  try {
    submitting.value = true
    const commentData: any = {
      content: commentContent.value,
      post_id: postId.value
    }
    
    // 如果是回复评论，添加 parent_id
    if (replyToComment.value) {
      commentData.parent_id = replyToComment.value.id
    }
    
    await createComment(commentData)
    message.success(replyToComment.value ? '回复成功' : '评论成功')
    commentContent.value = ''
    replyToComment.value = null
    fetchComments()
  } catch (error: any) {
    message.error(error.message || '评论失败')
  } finally {
    submitting.value = false
  }
}

function handleReply(parentComment: Comment, targetUser?: Comment) {
  // parentComment 是主评论（顶级评论）
  // targetUser 是要回复的具体用户（可能是主评论作者，也可能是回复者）
  replyToComment.value = parentComment
  replyToUser.value = targetUser || parentComment
  commentContent.value = `@${(targetUser || parentComment).user.nickname} `
  // 滚动到评论框
  nextTick(() => {
    const commentForm = document.querySelector('.comment-form textarea')
    if (commentForm) {
      (commentForm as HTMLElement).focus()
    }
  })
}

// 获取回复目标的名称
function getReplyTargetName(reply: Comment, parentComment: Comment): string {
  // 如果回复的内容以 @xxx 开头，尝试提取用户名
  const match = reply.content.match(/^@(\S+)\s/)
  if (match) {
    return match[1]
  }
  // 否则默认显示主评论作者
  return parentComment.user.nickname
}

// 移除评论内容开头的 @xxx
function removeAtMention(content: string): string {
  return content.replace(/^@\S+\s/, '')
}

// 切换评论展开/收起
function toggleExpand(commentId: number) {
  if (expandedComments.value.has(commentId)) {
    expandedComments.value.delete(commentId)
  } else {
    expandedComments.value.add(commentId)
  }
}

// 判断是否可以删除评论
function canDeleteComment(comment: Comment): boolean {
  if (!authStore.isLoggedIn) return false
  // 管理员可以删除所有评论，普通用户只能删除自己的评论
  return authStore.isAdmin || comment.user_id === authStore.user?.id
}

// 删除评论
async function handleDeleteComment(commentId: number) {
  try {
    await deleteComment(commentId)
    message.success('删除成功')
    fetchComments()
  } catch (error: any) {
    message.error(error.message || '删除失败')
  }
}

function handleEdit() {
  router.push(`/admin/posts/edit/${postId.value}`)
}

// 生成目录
function generateToc() {
  const article = document.querySelector('.post-content')
  if (!article) return

  const headings = article.querySelectorAll('h1, h2, h3, h4, h5, h6')
  const items: TocItem[] = []

  headings.forEach((heading, index) => {
    const level = parseInt(heading.tagName.substring(1))
    const text = heading.textContent || ''
    const id = `heading-${index}`
    
    // 给标题添加 id
    heading.id = id
    
    items.push({ id, text, level })
  })

  tocItems.value = items
}

// 滚动到指定标题
function scrollToHeading(id: string) {
  const element = document.getElementById(id)
  if (!element) return
  
  // 尝试多个可能的滚动容器
  let scrollContainer = document.querySelector('.main-content .n-scrollbar-container') as HTMLElement
  if (!scrollContainer) {
    scrollContainer = document.querySelector('.main-content') as HTMLElement
  }
  if (!scrollContainer) {
    scrollContainer = document.querySelector('.n-layout-scroll-container') as HTMLElement
  }
  
  if (scrollContainer) {
    // 获取元素相对于滚动容器的位置
    const containerRect = scrollContainer.getBoundingClientRect()
    const elementRect = element.getBoundingClientRect()
    const currentScroll = scrollContainer.scrollTop
    
    // 计算目标滚动位置
    const targetScroll = currentScroll + (elementRect.top - containerRect.top) - 100
    
    // 平滑滚动
    scrollContainer.scrollTo({ top: targetScroll, behavior: 'smooth' })
    activeHeading.value = id
  }
}

// 监听滚动，高亮当前标题
function handleScroll() {
  let scrollContainer = document.querySelector('.main-content .n-scrollbar-container') as HTMLElement
  if (!scrollContainer) {
    scrollContainer = document.querySelector('.main-content') as HTMLElement
  }
  if (!scrollContainer) {
    scrollContainer = document.querySelector('.n-layout-scroll-container') as HTMLElement
  }
  
  if (!scrollContainer) return
  
  // 计算阅读进度
  const scrollTop = scrollContainer.scrollTop
  const scrollHeight = scrollContainer.scrollHeight
  const clientHeight = scrollContainer.clientHeight
  const maxScroll = scrollHeight - clientHeight
  
  if (maxScroll > 0) {
    const progress = Math.round((scrollTop / maxScroll) * 100)
    readingProgress.value = Math.min(100, Math.max(0, progress))
  } else {
    readingProgress.value = 100
  }
  
  const containerRect = scrollContainer.getBoundingClientRect()
  
  const headings = tocItems.value.map(item => {
    const element = document.getElementById(item.id)
    if (element) {
      const elementRect = element.getBoundingClientRect()
      // 相对于容器顶部的位置
      const relativeTop = elementRect.top - containerRect.top
      return {
        id: item.id,
        top: relativeTop
      }
    }
    return null
  }).filter(Boolean)

  // 找到最接近容器顶部的标题（在可视区域内）
  const current = headings.find(h => h && h.top > 0 && h.top < 200)
  if (current) {
    activeHeading.value = current.id
  } else if (headings.length > 0 && headings[0] && headings[0].top < 0) {
    // 找到已经滚动过去的最后一个标题
    const passed = headings.filter(h => h && h.top < 0)
    if (passed.length > 0) {
      activeHeading.value = passed[passed.length - 1]!.id
    }
  }
}

</script>

<style scoped>
.post-detail-container {
  max-width: 1400px;
  margin: 0 auto;
  position: relative;
  z-index: 1;
  padding: 0 24px;
}

.post-detail-page {
  max-width: 820px;
  margin: 0 auto;
}

.post-main {
  width: 100%;
}

.post-toc {
  width: 260px;
  position: fixed;
  right: max(24px, calc((100vw - 1400px) / 2 + 24px));
  top: 104px;
  height: 500px;
  max-height: calc(100vh - 124px);
  overflow-y: auto;
  z-index: 10;
  scrollbar-width: thin;
  scrollbar-color: rgba(8, 145, 178, 0.3) transparent;
}

.post-toc::-webkit-scrollbar {
  width: 6px;
}

.post-toc::-webkit-scrollbar-track {
  background: transparent;
}

.post-toc::-webkit-scrollbar-thumb {
  background: rgba(8, 145, 178, 0.3);
  border-radius: 3px;
}

.post-toc::-webkit-scrollbar-thumb:hover {
  background: rgba(8, 145, 178, 0.5);
}

html.dark .post-toc {
  scrollbar-color: rgba(56, 189, 248, 0.3) transparent;
}

html.dark .post-toc::-webkit-scrollbar-thumb {
  background: rgba(56, 189, 248, 0.3);
}

html.dark .post-toc::-webkit-scrollbar-thumb:hover {
  background: rgba(56, 189, 248, 0.5);
}

@media (max-width: 1400px) {
  .post-toc {
    right: 24px;
  }
}

@media (max-width: 1200px) {
  .post-toc {
    display: none;
  }
  
  .post-detail-container {
    padding: 0 16px;
  }
  
  .post-detail-page {
    max-width: 100%;
  }
}

/* 移动端优化 */
@media (max-width: 768px) {
  .post-detail-container {
    padding: 0 8px;
  }
  
  .post-title {
    font-size: 24px;
  }
  
  .post-meta {
    font-size: 13px;
  }
  
  .post-meta :deep(.n-space) {
    flex-wrap: wrap;
    row-gap: 8px;
  }
  
  .post-content {
    font-size: 15px;
  }
  
  .comments-section h3 {
    font-size: 18px;
  }
  
  .comment-item {
    padding: 12px 0;
  }
  
  .reply-list {
    margin-left: 24px;
    padding-left: 12px;
  }
  
  /* 减小卡片内边距，增加内容区域 */
  .post-detail-page :deep(.n-card .n-card__content) {
    padding: 16px 12px !important;
  }
}

/* 极小屏幕优化 */
@media (max-width: 480px) {
  .post-detail-container {
    padding: 0 6px;
  }
  
  .post-detail-page :deep(.n-card .n-card__content) {
    padding: 14px 10px !important;
  }
}

/* 玻璃态卡片效果 */
.post-detail-page :deep(.n-card) {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.5);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  border-radius: 16px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

html.dark .post-detail-page :deep(.n-card) {
  background: rgba(30, 41, 59, 0.7);
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

.post-header {
  margin-bottom: 24px;
}

.post-title {
  font-size: 36px;
  font-weight: 800;
  margin: 0 0 16px 0;
  color: #1a202c;
  letter-spacing: -0.02em;
  line-height: 1.2;
}

html.dark .post-title {
  color: #e5e5e5;
}

.post-meta {
  color: #64748b;
  margin-bottom: 16px;
  font-weight: 500;
}

html.dark .post-meta {
  color: #94a3b8;
}

.post-meta span {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.post-tags {
  margin-top: 16px;
}

.post-content {
  min-height: 200px;
  line-height: 1.8;
  color: #374151;
}

html.dark .post-content {
  color: #d1d5db;
}

.post-actions {
  padding: 24px 0;
}

.comments-section {
  margin-top: 32px;
}

.comments-section h3 {
  margin-bottom: 16px;
  font-weight: 700;
  color: #1a202c;
}

html.dark .comments-section h3 {
  color: #e5e5e5;
}

.comment-form {
  margin-bottom: 24px;
}

.comments-list {
  margin-top: 24px;
}

.comment-item {
  padding: 16px 0;
  border-bottom: 1px solid rgba(148, 163, 184, 0.15);
}

html.dark .comment-item {
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-content {
  flex: 1;
}

.comment-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.comment-time {
  color: #999;
  font-size: 12px;
}

.comment-content p {
  margin: 8px 0;
  line-height: 1.6;
}

.comment-actions {
  margin-top: 8px;
  display: flex;
  gap: 16px;
}

.reply-list {
  margin-top: 12px;
  padding-left: 16px;
  margin-left: 48px;
  border-left: 3px solid rgba(16, 185, 129, 0.15);
  background: rgba(16, 185, 129, 0.02);
  border-radius: 0 8px 8px 0;
  padding-top: 8px;
  padding-bottom: 4px;
}

html.dark .reply-list {
  border-left-color: rgba(52, 211, 153, 0.2);
  background: rgba(52, 211, 153, 0.03);
}

.reply-item {
  padding: 8px 0;
  position: relative;
}

.reply-item:not(:last-child) {
  border-bottom: 1px solid rgba(148, 163, 184, 0.08);
}

html.dark .reply-item:not(:last-child) {
  border-bottom-color: rgba(255, 255, 255, 0.05);
}

.reply-content {
  flex: 1;
}

.reply-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.reply-to {
  color: #10b981;
  font-size: 12px;
  font-weight: 400;
  padding: 2px 6px;
  background: rgba(16, 185, 129, 0.08);
  border-radius: 4px;
}

html.dark .reply-to {
  color: #34d399;
  background: rgba(52, 211, 153, 0.12);
}

.reply-content p {
  margin: 4px 0;
  line-height: 1.6;
}

/* 目录样式 */
.post-toc :deep(.n-card) {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.5);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  border-radius: 16px;
}

html.dark .post-toc :deep(.n-card) {
  background: rgba(30, 41, 59, 0.7);
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

/* 目录头部 */
.toc-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.toc-title {
  font-weight: 600;
  font-size: 16px;
  color: #1a202c;
}

html.dark .toc-title {
  color: #e5e5e5;
}

.toc-progress-text {
  font-size: 14px;
  font-weight: 600;
  color: #0891b2;
  font-family: 'Courier New', Consolas, monospace;
}

html.dark .toc-progress-text {
  color: #38bdf8;
}

/* 阅读进度条 */
.reading-progress-bar {
  width: 100%;
  height: 4px;
  background: rgba(8, 145, 178, 0.1);
  border-radius: 2px;
  overflow: hidden;
  margin-top: 8px;
}

html.dark .reading-progress-bar {
  background: rgba(56, 189, 248, 0.15);
}

.reading-progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #0891b2 0%, #059669 100%);
  border-radius: 2px;
  transition: width 0.3s ease;
  box-shadow: 0 0 8px rgba(8, 145, 178, 0.5);
}

html.dark .reading-progress-fill {
  background: linear-gradient(90deg, #38bdf8 0%, #4ade80 100%);
  box-shadow: 0 0 8px rgba(56, 189, 248, 0.5);
}

.toc-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.toc-item {
  display: block;
  padding: 8px 12px;
  color: #64748b;
  text-decoration: none;
  border-radius: 8px;
  transition: all 0.3s;
  cursor: pointer;
  font-size: 14px;
  line-height: 1.5;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

html.dark .toc-item {
  color: #94a3b8;
}

.toc-item:hover {
  background: rgba(8, 145, 178, 0.1);
  color: #0891b2;
  transform: translateX(4px);
}

html.dark .toc-item:hover {
  background: rgba(56, 189, 248, 0.1);
  color: #38bdf8;
}

.toc-item.active {
  background: rgba(8, 145, 178, 0.15);
  color: #0891b2;
  font-weight: 600;
  border-left: 3px solid #0891b2;
}

html.dark .toc-item.active {
  background: rgba(56, 189, 248, 0.15);
  color: #38bdf8;
  border-left-color: #38bdf8;
}

/* 标题层级缩进 */
.toc-level-1 {
  padding-left: 12px;
}

.toc-level-2 {
  padding-left: 24px;
  font-size: 13px;
}

.toc-level-3 {
  padding-left: 36px;
  font-size: 12px;
}

.toc-level-4,
.toc-level-5,
.toc-level-6 {
  padding-left: 48px;
  font-size: 12px;
  opacity: 0.8;
}

/* 返回顶部组件 */
.back-to-top {
  z-index: 1000 !important;
}

/* 返回顶部按钮样式 */
.back-top-button {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(8, 145, 178, 0.9);
  backdrop-filter: blur(20px) saturate(180%);
  border-radius: 50%;
  box-shadow: 0 4px 16px rgba(8, 145, 178, 0.3);
  color: white;
  cursor: pointer;
  transition: all 0.3s;
}

.back-top-button:hover {
  background: rgba(8, 145, 178, 1);
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(8, 145, 178, 0.5);
}

html.dark .back-top-button {
  background: rgba(56, 189, 248, 0.9);
  box-shadow: 0 4px 16px rgba(56, 189, 248, 0.3);
}

html.dark .back-top-button:hover {
  background: rgba(56, 189, 248, 1);
  box-shadow: 0 8px 24px rgba(56, 189, 248, 0.5);
}
</style>

