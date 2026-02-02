<!--
 * @ProjectName: go-vue3-blog
 * @FileName: Moments.vue
 * @CreateTime: 2026-02-02 11:43:17
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 说说页面组件，展示所有说说动态
 -->
<template>
  <div class="moments-page">
    <div class="moments-layout">
      <div class="moments-main">
        <n-spin :show="loading">
          <div v-if="moments.length > 0" class="moments-list">
            <div v-for="moment in moments" :key="moment.id" class="moment-card">
              <!-- 左侧头像 -->
              <div class="moment-avatar">
                <n-avatar :src="moment.user.avatar" :size="56" round />
              </div>

              <!-- 右侧内容区 -->
              <div class="moment-main">
                <div class="moment-header">
                  <div class="user-info">
                    <span class="username">{{ moment.user.nickname || moment.user.username }}</span>
                    <span class="time">{{ formatDate(moment.created_at, 'YYYY-MM-DD HH:mm') }}</span>
                  </div>
                </div>

                <div class="moment-content">
                  <p>{{ moment.content }}</p>
                </div>

                <div v-if="moment.images" class="moment-images">
                  <n-image-group>
                    <n-space :size="12">
                      <n-image
                        v-for="(img, index) in parseImages(moment.images)"
                        :key="index"
                        :src="img"
                        object-fit="cover"
                        class="moment-image"
                      />
                    </n-space>
                  </n-image-group>
                </div>

                <div class="moment-footer">
                  <n-space :size="16">
                    <span class="stat-item like-button" @click="handleLike(moment)">
                      <n-icon :component="moment.liked ? Heart : HeartOutline" :class="{ liked: moment.liked }" />
                      <span v-if="moment.like_count > 0" class="like-text">{{ moment.like_count }}</span>
                    </span>
                    <span class="stat-item comment-button" @click="toggleCommentInput(moment.id)">
                      <n-icon :component="ChatbubbleOutline" />
                      <span v-if="getCommentCount(moment.id) > 0" class="comment-text">{{ getCommentCount(moment.id) }}</span>
                    </span>
                  </n-space>
                </div>

                <!-- 评论区域 -->
                <div v-if="momentComments[moment.id] && momentComments[moment.id].length > 0" class="comments-section">
                  <div v-for="comment in momentComments[moment.id]" :key="comment.id" class="comment-item">
                    <div class="comment-content">
                      <span class="comment-author">{{ comment.user.nickname || comment.user.username }}</span>
                      <template v-if="comment.parent">
                        <span class="comment-reply-to">回复</span>
                        <span class="comment-reply-target">{{ comment.parent.user.nickname || comment.parent.user.username }}</span>
                      </template>
                      <span class="comment-colon">:</span>
                      <span class="comment-text">{{ comment.content }}</span>
                    </div>
                    <div class="comment-actions">
                      <span class="comment-time">{{ formatRelativeTime(comment.created_at) }}</span>
                      <n-button
                        v-if="authStore.isLoggedIn"
                        text
                        size="tiny"
                        @click="handleReply(moment.id, comment)"
                      >
                        回复
                      </n-button>
                      <n-popconfirm
                        v-if="canDeleteComment(comment)"
                        @positive-click="handleDeleteComment(comment.id, moment.id)"
                      >
                        <template #trigger>
                          <n-button text size="tiny" type="error">删除</n-button>
                        </template>
                        确定要删除这条评论吗？
                      </n-popconfirm>
                    </div>
                    <!-- 子评论 -->
                    <div v-if="comment.children && comment.children.length > 0" class="reply-list">
                      <div v-for="reply in comment.children" :key="reply.id" class="reply-item">
                        <span class="comment-author">{{ reply.user.nickname || reply.user.username }}</span>
                        <span class="comment-reply-to">回复</span>
                        <span class="comment-reply-target">{{ comment.user.nickname || comment.user.username }}</span>
                        <span class="comment-colon">：</span>
                        <span class="comment-text">{{ reply.content }}</span>
                        <div class="comment-actions">
                          <span class="comment-time">{{ formatRelativeTime(reply.created_at) }}</span>
                          <n-button
                            v-if="authStore.isLoggedIn"
                            text
                            size="tiny"
                            @click="handleReply(moment.id, comment, reply)"
                          >
                            回复
                          </n-button>
                          <n-popconfirm
                            v-if="canDeleteComment(reply)"
                            @positive-click="handleDeleteComment(reply.id, moment.id)"
                          >
                            <template #trigger>
                              <n-button text size="tiny" type="error">删除</n-button>
                            </template>
                            确定要删除这条回复吗？
                          </n-popconfirm>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- 评论输入框 -->
                <div v-if="activeCommentInput === moment.id" class="comment-input-section">
                  <div v-if="replyToComment[moment.id]" class="reply-hint">
                    回复 <strong>{{ getReplyTargetName(moment.id) }}</strong>
                    <n-button text size="tiny" @click="cancelReply(moment.id)">取消回复</n-button>
                  </div>
                  <n-input
                    v-model:value="commentInputs[moment.id]"
                    type="textarea"
                    :placeholder="replyToComment[moment.id] ? `回复 ${getReplyTargetName(moment.id)}...` : '写评论...'"
                    :rows="2"
                    :maxlength="500"
                    show-count
                    @keydown.enter.ctrl="handleSubmitComment(moment.id)"
                  />
                  <div class="comment-input-actions">
                    <n-space justify="end">
                      <n-button size="small" @click="cancelComment(moment.id)">取消</n-button>
                      <n-button
                        type="primary"
                        size="small"
                        :loading="submittingComments[moment.id]"
                        :disabled="!commentInputs[moment.id]?.trim()"
                        @click="handleSubmitComment(moment.id)"
                      >
                        发送
                      </n-button>
                    </n-space>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <n-empty v-else description="还没有发布说说" style="margin: 60px 0" />
        </n-spin>

        <!-- 分页 -->
        <div v-if="total > pageSize" class="pagination-wrapper">
          <n-pagination
            v-model:page="currentPage"
            :page-count="Math.ceil(total / pageSize)"
            :page-size="pageSize"
            show-size-picker
            :page-sizes="[10, 20, 30]"
            @update:page="handlePageChange"
            @update:page-size="handlePageSizeChange"
          />
        </div>
      </div>

      <div class="sidebar-section">
        <div class="sidebar-card-wrapper">
          <AuthorCard />
        </div>
        <div class="sidebar-card-wrapper">
          <AnnouncementBoard :limit="3" />
        </div>
        <div class="sidebar-card-wrapper">
          <TagCloudWidget />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { useMessage } from 'naive-ui'
import { HeartOutline, Heart, ChatbubbleOutline } from '@vicons/ionicons5'
import { getMoments, likeMoment } from '@/api/moment'
import type { Moment, MomentParams } from '@/api/moment'
import { formatDate, formatRelativeTime } from '@/utils/format'
import { useAuthStore } from '@/stores/auth'
import { getCommentsByTypeAndTarget, createComment, deleteComment } from '@/api/comment'
import type { Comment } from '@/types/blog'
import AuthorCard from '@/components/AuthorCard.vue'
import AnnouncementBoard from '@/components/AnnouncementBoard.vue'
import TagCloudWidget from '@/components/TagCloudWidget.vue'

const message = useMessage()
const authStore = useAuthStore()

const loading = ref(false)
const moments = ref<Moment[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 评论相关
const momentComments = reactive<Record<number, Comment[]>>({})
const activeCommentInput = ref<number | null>(null)
const commentInputs = reactive<Record<number, string>>({})
const submittingComments = reactive<Record<number, boolean>>({})
const replyToComment = reactive<Record<number, { parent: Comment; target?: Comment } | null>>({})

// 说说评论类型
const MOMENT_COMMENT_TYPE = 'moment'

// 解析图片JSON
function parseImages(images: string): string[] {
  try {
    return images ? JSON.parse(images) : []
  } catch {
    return []
  }
}

// 获取说说列表
async function fetchMoments() {
  try {
    loading.value = true
    const isAdmin = authStore.isLoggedIn && authStore.user?.role === 'admin'

    const params: MomentParams = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (!isAdmin) {
      params.status = 1
    }

    const res = await getMoments(params) as any

    if (res && res.data) {
      moments.value = res.data.list || []
      total.value = res.data.total || 0
      
      // 获取每个说说的评论
      for (const moment of moments.value) {
        await fetchComments(moment.id)
      }
    }
  } catch (error) {
    console.error('获取说说列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取说说的评论列表
async function fetchComments(momentId: number) {
  try {
    const res = await getCommentsByTypeAndTarget(MOMENT_COMMENT_TYPE, momentId)
    if (res.data) {
      momentComments[momentId] = res.data
    }
  } catch (error) {
    console.error('获取评论失败:', error)
    momentComments[momentId] = []
  }
}

// 获取评论数量
function getCommentCount(momentId: number): number {
  const comments = momentComments[momentId] || []
  let count = comments.length
  comments.forEach(comment => {
    if (comment.children) {
      count += comment.children.length
    }
  })
  return count
}

// 切换评论输入框
function toggleCommentInput(momentId: number) {
  if (!authStore.isLoggedIn) {
    message.warning('请先登录')
    return
  }
  
  if (activeCommentInput.value === momentId) {
    activeCommentInput.value = null
    commentInputs[momentId] = ''
    replyToComment[momentId] = null
  } else {
    activeCommentInput.value = momentId
    if (!commentInputs[momentId]) {
      commentInputs[momentId] = ''
    }
  }
}

// 取消评论
function cancelComment(momentId: number) {
  activeCommentInput.value = null
  commentInputs[momentId] = ''
  replyToComment[momentId] = null
}

// 取消回复
function cancelReply(momentId: number) {
  replyToComment[momentId] = null
  commentInputs[momentId] = ''
}

// 获取回复目标名称
function getReplyTargetName(momentId: number): string {
  const replyInfo = replyToComment[momentId]
  if (!replyInfo) {
    return ''
  }
  if (replyInfo.target) {
    return replyInfo.target.user.nickname || replyInfo.target.user.username
  }
  return replyInfo.parent.user.nickname || replyInfo.parent.user.username
}

// 回复评论
function handleReply(momentId: number, parent: Comment, target?: Comment) {
  if (!authStore.isLoggedIn) {
    message.warning('请先登录')
    return
  }
  
  activeCommentInput.value = momentId
  replyToComment[momentId] = { parent, target }
  if (!commentInputs[momentId]) {
    commentInputs[momentId] = ''
  }
  
  // 聚焦输入框
  setTimeout(() => {
    const textarea = document.querySelector(`textarea[placeholder="写评论..."]`) as HTMLTextAreaElement
    if (textarea) {
      textarea.focus()
    }
  }, 100)
}

// 提交评论
async function handleSubmitComment(momentId: number) {
  if (!authStore.isLoggedIn) {
    message.warning('请先登录')
    return
  }

  const content = commentInputs[momentId]?.trim()
  if (!content) {
    message.warning('请输入评论内容')
    return
  }

  try {
    submittingComments[momentId] = true
    const commentData: any = {
      content,
      comment_type: MOMENT_COMMENT_TYPE,
      target_id: momentId
    }
    
    // 如果是回复评论，添加 parent_id
    const replyInfo = replyToComment[momentId]
    if (replyInfo && replyInfo.parent) {
      commentData.parent_id = replyInfo.parent.id
    }
    
    await createComment(commentData)
    message.success(replyInfo ? '回复成功' : '评论成功')
    commentInputs[momentId] = ''
    replyToComment[momentId] = null
    activeCommentInput.value = null
    
    // 重新获取评论列表
    await fetchComments(momentId)
  } catch (error: any) {
    message.error(error.message || '评论失败')
  } finally {
    submittingComments[momentId] = false
  }
}

// 删除评论
async function handleDeleteComment(commentId: number, momentId: number) {
  try {
    await deleteComment(commentId)
    message.success('删除成功')
    await fetchComments(momentId)
  } catch (error: any) {
    message.error(error.message || '删除失败')
  }
}

// 检查是否可以删除评论
function canDeleteComment(comment: Comment): boolean {
  if (!authStore.isLoggedIn) {
    return false
  }
  const user = authStore.user
  if (!user) {
    return false
  }
  // 管理员或评论作者可以删除
  return user.role === 'admin' || user.id === comment.user_id
}

// 页码变化
function handlePageChange(page: number) {
  currentPage.value = page
  fetchMoments()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

// 每页数量变化
function handlePageSizeChange(size: number) {
  pageSize.value = size
  currentPage.value = 1
  fetchMoments()
}

// 点赞/取消点赞说说
async function handleLike(moment: Moment) {
  const wasLiked = moment.liked
  
  try {
    // 乐观更新 UI
    if (wasLiked) {
      // 取消点赞
      moment.like_count--
      moment.liked = false
    } else {
      // 点赞
      moment.like_count++
      moment.liked = true
    }
    
    // 调用后端 API
    await likeMoment(moment.id)
  } catch (error) {
    // 失败时回滚
    if (wasLiked) {
      moment.like_count++
      moment.liked = true
    } else {
      moment.like_count--
      moment.liked = false
    }
    console.error('操作失败:', error)
    message.error('操作失败，请重试')
  }
}

onMounted(() => {
  fetchMoments()
})
</script>

<style scoped>
.moments-page {
  min-height: calc(100vh - 180px);
  padding: 32px 0;
}

.moments-layout {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 20px;
  display: grid;
  grid-template-columns: minmax(0, 1fr) 360px;
  gap: 32px;
  align-items: start;
}

.moments-main {
  min-width: 0;
}

.moments-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* 卡片布局 - 左侧头像，右侧内容 */
.moment-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  transition: all 0.3s;
  display: flex;
  gap: 16px;
}

.moment-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

html.dark .moment-card {
  background: #1f1f1f;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

html.dark .moment-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.4);
}

/* 左侧头像区域 */
.moment-avatar {
  flex-shrink: 0;
}

.moment-avatar :deep(.n-avatar) {
  transition: all 0.3s;
}

.moment-card:hover .moment-avatar :deep(.n-avatar) {
  transform: scale(1.05);
}

/* 右侧主内容区域 */
.moment-main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.moment-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.user-info {
  display: flex;
  align-items: baseline;
  gap: 12px;
  flex-wrap: wrap;
}

.username {
  font-size: 16px;
  font-weight: 700;
  color: #1a202c;
  letter-spacing: -0.01em;
}

html.dark .username {
  color: #e5e5e5;
}

.time {
  font-size: 13px;
  color: #94a3b8;
  font-weight: 400;
}

html.dark .time {
  color: #64748b;
}

/* 说说内容 */
.moment-content {
  line-height: 1.8;
}

.moment-content p {
  font-size: 15px;
  line-height: 1.8;
  color: #334155;
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
}

html.dark .moment-content p {
  color: #cbd5e1;
}

/* 图片展示区域 */
.moment-images {
  margin-top: 8px;
}

.moment-image {
  width: 200px;
  height: 200px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  border: 1px solid rgba(0, 0, 0, 0.06);
}

.moment-image:hover {
  transform: scale(1.02);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

html.dark .moment-image {
  border-color: rgba(255, 255, 255, 0.1);
}

/* 底部互动区域 */
.moment-footer {
  padding-top: 12px;
  border-top: 1px solid #e5e7eb;
  margin-top: 4px;
}

.comment-button {
  cursor: pointer;
  user-select: none;
  padding: 6px 12px;
  border-radius: 6px;
  background: transparent;
  transition: all 0.3s;
}

.comment-button:hover {
  color: #3b82f6;
  background: rgba(59, 130, 246, 0.08);
  transform: translateY(-1px);
}

html.dark .comment-button:hover {
  color: #60a5fa;
  background: rgba(96, 165, 250, 0.12);
}

.comment-text {
  font-weight: 500;
}

/* 评论区域 */
.comments-section {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #e5e7eb;
}

html.dark .comments-section {
  border-top-color: #374151;
}

.comment-item {
  margin-bottom: 8px;
  padding: 8px 12px;
  background: #f8fafc;
  border-radius: 8px;
  transition: all 0.2s;
}

html.dark .comment-item {
  background: #1e293b;
}

.comment-content {
  display: flex;
  flex-wrap: wrap;
  align-items: baseline;
  gap: 4px;
  line-height: 1.6;
  font-size: 14px;
}

.comment-author {
  font-weight: 600;
  color: #1e40af;
  cursor: pointer;
}

html.dark .comment-author {
  color: #60a5fa;
}

.comment-author:hover {
  text-decoration: underline;
}

.comment-reply-to {
  color: #64748b;
  font-size: 13px;
}

html.dark .comment-reply-to {
  color: #94a3b8;
}

.comment-reply-target {
  font-weight: 600;
  color: #1e40af;
  cursor: pointer;
}

html.dark .comment-reply-target {
  color: #60a5fa;
}

.comment-reply-target:hover {
  text-decoration: underline;
}

.comment-colon {
  color: #334155;
  margin: 0 2px;
}

html.dark .comment-colon {
  color: #cbd5e1;
}

.comment-text {
  color: #334155;
  word-break: break-word;
}

html.dark .comment-text {
  color: #cbd5e1;
}

.comment-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 4px;
  font-size: 12px;
}

.comment-time {
  color: #94a3b8;
  font-size: 12px;
}

html.dark .comment-time {
  color: #64748b;
}

/* 回复列表 */
.reply-list {
  margin-top: 8px;
  padding-left: 12px;
  border-left: 2px solid #e2e8f0;
}

html.dark .reply-list {
  border-left-color: #475569;
}

.reply-item {
  margin-top: 6px;
  padding: 6px 10px;
  background: #f1f5f9;
  border-radius: 6px;
  display: flex;
  flex-wrap: wrap;
  align-items: baseline;
  gap: 4px;
  font-size: 13px;
  line-height: 1.5;
}

html.dark .reply-item {
  background: #0f172a;
}

/* 评论输入区域 */
.comment-input-section {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #e5e7eb;
}

html.dark .comment-input-section {
  border-top-color: #374151;
}

.comment-input-actions {
  margin-top: 8px;
}

.reply-hint {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 12px;
  margin-bottom: 8px;
  background: #eff6ff;
  border-radius: 6px;
  font-size: 13px;
  color: #1e40af;
}

html.dark .reply-hint {
  background: #1e3a8a;
  color: #93c5fd;
}

.reply-hint strong {
  font-weight: 600;
}

html.dark .moment-footer {
  border-top-color: #374151;
}

.stat-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: #64748b;
  transition: all 0.3s;
}

.like-button {
  cursor: pointer;
  user-select: none;
  padding: 6px 12px;
  border-radius: 6px;
  background: transparent;
  transition: all 0.3s;
}

.like-button:hover {
  color: #ef4444;
  background: rgba(239, 68, 68, 0.08);
  transform: translateY(-1px);
}

.like-button:active {
  transform: translateY(0) scale(0.98);
}

.like-button .liked {
  color: #ef4444;
  animation: likeAnimation 0.3s ease;
}

.like-text {
  font-weight: 500;
}

@keyframes likeAnimation {
  0% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.3);
  }
  100% {
    transform: scale(1);
  }
}

html.dark .stat-item {
  color: #94a3b8;
}

html.dark .like-button:hover {
  color: #f87171;
  background: rgba(248, 113, 113, 0.12);
}

/* 分页 */
.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 40px;
  padding: 24px 0;
}

.sidebar-section {
  position: relative;
  z-index: 10;
  margin-left: 12px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* 统一侧边栏卡片圆角样式，与首页保持一致 */
.moments-page :deep(.sidebar-card-wrapper .n-card) {
  border-radius: 16px;
}

@media (max-width: 1024px) {
  .moments-layout {
    grid-template-columns: 1fr;
    padding: 0 16px;
  }

  .sidebar-section {
    display: none;
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .moment-card {
    padding: 16px;
    gap: 12px;
  }

  .moment-avatar :deep(.n-avatar) {
    width: 44px !important;
    height: 44px !important;
  }

  .username {
    font-size: 15px;
  }

  .time {
    font-size: 12px;
  }

  .moment-content p {
    font-size: 14px;
  }

  .moment-image {
    width: 150px;
    height: 150px;
  }
}

@media (max-width: 480px) {
  .moments-container {
    padding: 0 16px;
  }

  .moment-card {
    padding: 12px;
  }

  .moment-image {
    width: 120px;
    height: 120px;
  }

  .user-info {
    flex-direction: column;
    gap: 4px;
    align-items: flex-start;
  }
}
</style>

