<!--
 * @ProjectName: go-vue3-blog
 * @FileName: FriendLinks.vue
 * @CreateTime: 2026-02-02 11:43:17
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 友情链接页面组件，展示所有友情链接
 -->
<template>
  <div class="friendlinks-page">
    <div class="friendlinks-layout">
      <div class="friendlinks-main">
        <n-spin :show="loading">
          <div class="page-header">
            <h1 class="page-title">🔗 友情链接</h1>
            <p class="page-description">感谢这些优秀的网站和博客，让我们一起成长</p>
          </div>

          <template v-if="!loading">
            <div v-if="categories.length === 0" class="empty-container">
              <n-empty description="暂无友链分类" />
            </div>

            <!-- 按分类展示友链 -->
            <div v-else class="categories-container">
              <div v-for="category in sortedCategories" :key="category.id" class="category-section">
                <div class="category-header">
                  <h2 class="category-title">
                    {{ category.name }} ({{ getCategoryLinksCount(category.id) }})
                  </h2>
                  <p v-if="category.description" class="category-description">
                    {{ category.description }}
                  </p>
                </div>

                <div v-if="getCategoryLinks(category.id).length === 0" class="empty-category">
                  <n-empty description="该分类下暂无友链" size="small" />
                </div>

                <div v-else class="friendlinks-grid">
                  <div
                    v-for="link in getCategoryLinks(category.id)"
                    :key="link.id"
                    class="friendlink-card"
                    @click="handleClick(link.url)"
                  >
                    <div class="card-header">
                      <div class="link-icon">
                        <n-image
                          v-if="link.icon"
                          :src="link.icon"
                          :alt="link.name"
                          width="36"
                          height="36"
                          object-fit="cover"
                          :preview-disabled="true"
                          fallback-src="/logo.jpg"
                        />
                        <div v-else class="icon-placeholder">
                          <n-icon :component="LinkOutline" size="20" />
                        </div>
                      </div>
                      <div class="link-info">
                        <h3 class="link-name">
                          {{ link.name }}
                          <n-tag
                            v-if="link.is_invalid || link.accessible > 0"
                            :type="getLinkStatusType(link)"
                            size="small"
                            style="margin-left: 8px"
                          >
                            {{ getLinkStatusText(link) }}
                          </n-tag>
                        </h3>
                        <p v-if="link.description" class="link-description">
                          {{ link.description }}
                        </p>
                      </div>
                    </div>

                    <div v-if="link.screenshot" class="card-screenshot">
                      <n-image
                        :src="link.screenshot"
                        :alt="link.name"
                        width="100%"
                        height="120"
                        object-fit="cover"
                        :preview-disabled="true"
                        class="screenshot-image"
                      />
                    </div>

                    <div class="card-footer">
                      <a
                        :href="link.url"
                        target="_blank"
                        rel="noopener noreferrer"
                        class="link-url"
                        @click.stop
                      >
                        {{ formatURL(link.url) }}
                      </a>
                      <span v-if="link.atom_url" class="rss-badge" title="支持 RSS 订阅">
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          viewBox="0 0 24 24"
                          fill="currentColor"
                          width="14"
                          height="14"
                        >
                          <circle cx="6.18" cy="17.82" r="2.18" />
                          <path
                            d="M4 4.44v2.83c7.03 0 12.73 5.7 12.73 12.73h2.83c0-8.59-6.97-15.56-15.56-15.56zm0 5.66v2.83c3.9 0 7.07 3.17 7.07 7.07h2.83c0-5.47-4.43-9.9-9.9-9.9z"
                          />
                        </svg>
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </template>
        </n-spin>

        <!-- 友情链接申请 -->
        <div class="application-section">
          <n-card class="application-card">
            <h2 class="section-title">友情链接申请</h2>

            <div class="intro-text">
              <p>
                很高兴能和各位优秀的朋友们交流，本站友链目前采用手动添加，如果你想加入友链，可以在下方留言，我会在不忙的时候统一添加。
              </p>
            </div>

            <div class="rules-section">
              <h3 class="rules-title">申请友链</h3>
              <ul class="rules-list">
                <li>
                  我已添加
                  <strong v-if="myFriendLinkInfo.url">
                    <a
                      :href="myFriendLinkInfo.url"
                      target="_blank"
                      rel="noopener noreferrer"
                      class="site-link"
                    >
                      {{ myFriendLinkInfo.name || '菱风叙' }}
                    </a>
                  </strong>
                  <strong v-else>{{ myFriendLinkInfo.name || '菱风叙' }}</strong>
                  的友情链接。
                </li>
                <li>
                  <strong>请多多进行互动后再来进行友链添加</strong
                  >，若为首次评论直接申请友链，将不会通过。
                </li>
                <li>
                  本站不添加<strong>采集站、纯搬运站点、论坛类站点</strong>等<strong>非个人博客</strong>类型的站点。
                </li>
                <li>站点目前可以在<strong>中国大陆区域</strong>正常访问。</li>
                <li>
                  需要是独立域名，不接受 <code>github.io</code>、<code>vercel.app</code>
                  等第三方域名。
                </li>
                <li>网站内容符合<strong>中国大陆法律法规</strong>。</li>
              </ul>
              <p class="disclaimer">
                若申请时或日后有违反上述规定的站点，博主有权自行删除且不进行通知！
              </p>
            </div>

            <div class="my-info-section">
              <h3 class="info-title">我的友链信息</h3>
              <div class="info-grid">
                <div class="info-item">
                  <span class="info-label">名称：</span>
                  <span class="info-value">{{ myFriendLinkInfo.name || '菱风叙' }}</span>
                </div>
                <div class="info-item" v-if="myFriendLinkInfo.desc">
                  <span class="info-label">描述：</span>
                  <span class="info-value">{{ myFriendLinkInfo.desc }}</span>
                </div>
                <div class="info-item" v-if="myFriendLinkInfo.url">
                  <span class="info-label">地址：</span>
                  <span class="info-value">{{ myFriendLinkInfo.url }}</span>
                </div>
                <div class="info-item" v-if="myFriendLinkInfo.avatar">
                  <span class="info-label">头像：</span>
                  <span class="info-value">{{ myFriendLinkInfo.avatar }}</span>
                </div>
                <div class="info-item" v-if="myFriendLinkInfo.screenshot">
                  <span class="info-label">站点图片：</span>
                  <span class="info-value">{{ myFriendLinkInfo.screenshot }}</span>
                </div>
                <div class="info-item" v-if="myFriendLinkInfo.rss">
                  <span class="info-label">订阅：</span>
                  <span class="info-value">{{ myFriendLinkInfo.rss }}</span>
                </div>
              </div>
            </div>

            <div class="yaml-section">
              <h3 class="yaml-title">YAML 格式</h3>
              <div class="code-wrapper">
                <n-code :code="yamlCode" language="yaml" :show-line-numbers="true" />
                <n-button class="copy-code-btn" size="small" quaternary @click="handleCopyYaml">
                  <template #icon>
                    <n-icon :component="CopyOutline" />
                  </template>
                  {{ copyYamlSuccess ? '已复制' : '复制' }}
                </n-button>
              </div>
            </div>
          </n-card>
        </div>

        <!-- 评论区 -->
        <div class="comments-section">
          <n-card class="comments-card">
            <h2 class="section-title">评论区 ({{ comments.length }})</h2>

            <!-- 评论表单 -->
            <n-card v-if="authStore.isLoggedIn" class="comment-form">
              <!-- 回复提示 -->
              <n-alert
                v-if="replyToComment"
                type="info"
                closable
                style="margin-bottom: 12px"
                @close="
                  replyToComment = null
                  replyToUser = null
                  commentContent = ''
                "
              >
                正在回复
                <strong>@{{ (replyToUser || replyToComment).user.nickname }}</strong> 的评论
              </n-alert>

              <CommentMarkdownEditor v-model="commentContent" height="250px" :max-length="5000" />
              <div style="margin-top: 12px; text-align: right">
                <n-button type="primary" :loading="submitting" @click="handleSubmitComment">
                  {{ replyToComment ? '发表回复' : '发表评论' }}
                </n-button>
              </div>
            </n-card>

            <n-alert v-else type="info" style="margin-bottom: 16px">
              请
              <n-button text type="primary" @click="router.push('/auth/login')">登录</n-button>
              后发表评论
            </n-alert>

            <!-- 评论列表 -->
            <div class="comments-list">
              <div v-if="comments.length === 0" class="empty-comments">
                <n-empty description="暂无评论，快来抢沙发吧~" size="small" />
              </div>
              <div v-for="comment in comments" :key="comment.id" class="comment-item">
                <n-space align="start">
                  <n-avatar :src="comment.user.avatar" round />
                  <div class="comment-content">
                    <div class="comment-header">
                      <strong>{{ comment.user.nickname }}</strong>
                      <span class="comment-time">{{
                        formatDate(comment.created_at, 'YYYY年MM月DD日 HH:mm')
                      }}</span>
                    </div>
                    <CommentContent :content="comment.content" />
                    <div class="comment-actions">
                      <n-button
                        v-if="authStore.isLoggedIn"
                        text
                        size="small"
                        @click="handleReply(comment)"
                      >
                        回复
                      </n-button>
                      <n-button
                        v-if="comment.children && comment.children.length > 0"
                        text
                        size="small"
                        @click="toggleExpand(comment.id)"
                      >
                        {{
                          expandedComments.has(comment.id)
                            ? '收起'
                            : `展开 ${comment.children.length} 条回复`
                        }}
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
                    <div
                      v-if="
                        comment.children &&
                        comment.children.length > 0 &&
                        expandedComments.has(comment.id)
                      "
                      class="reply-list"
                    >
                      <div v-for="reply in comment.children" :key="reply.id" class="reply-item">
                        <n-space align="start">
                          <n-avatar :src="reply.user.avatar" round size="small" />
                          <div class="reply-content">
                            <div class="reply-header">
                              <strong>{{ reply.user.nickname }}</strong>
                              <span class="reply-to"
                                >回复 @{{ getReplyTargetName(reply, comment) }}</span
                              >
                              <span class="comment-time">{{
                                formatDate(reply.created_at, 'YYYY年MM月DD日 HH:mm')
                              }}</span>
                            </div>
                            <CommentContent :content="removeAtMention(reply.content)" />
                            <div class="comment-actions">
                              <n-button
                                v-if="authStore.isLoggedIn"
                                text
                                size="small"
                                @click="handleReply(comment, reply)"
                              >
                                回复
                              </n-button>
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
          </n-card>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage, NIcon } from 'naive-ui'
import { LinkOutline, CopyOutline } from '@vicons/ionicons5'
import { getFriendLinks, getFriendLinkCategories } from '@/api/friendlink'
import type { FriendLink, FriendLinkCategory } from '@/api/friendlink'
import { getCommentsByTypeAndTarget, createComment, deleteComment } from '@/api/comment'
import { formatDate } from '@/utils/format'
import { useAuthStore } from '@/stores'
import type { Comment } from '@/types/blog'
import { getFriendLinkInfo, type FriendLinkInfo } from '@/api/setting'
import CommentMarkdownEditor from '@/components/CommentMarkdownEditor.vue'
import CommentContent from '@/components/CommentContent.vue'

const router = useRouter()
const message = useMessage()
const authStore = useAuthStore()

const loading = ref(false)
const submitting = ref(false)
const friendLinks = ref<FriendLink[]>([])
const categories = ref<FriendLinkCategory[]>([])
const comments = ref<Comment[]>([])
const commentContent = ref('')
const replyToComment = ref<Comment | null>(null)
const replyToUser = ref<Comment | null>(null)
const expandedComments = ref<Set<number>>(new Set())
const copyYamlSuccess = ref(false)
const myFriendLinkInfo = ref<FriendLinkInfo>({
  name: '菱风叙',
  desc: '',
  url: '',
  avatar: '',
  screenshot: '',
  rss: ''
})

// 友链评论类型
const FRIENDLINK_COMMENT_TYPE = 'friendlink'
const FRIENDLINK_TARGET_ID = 0 // 友链页面的target_id固定为0

const yamlCode = computed(() => {
  const info = myFriendLinkInfo.value
  const lines = [
    `name: ${info.name || '菱风叙'}`,
    `desc: ${info.desc || ''}`,
    `url: ${info.url || ''}`,
    `avatar: ${info.avatar || ''}`,
    `screenshot: ${info.screenshot || ''}`
  ]
  // RSS订阅可为空，如果为空则不显示
  if (info.rss) {
    lines.push(`rss: ${info.rss}`)
  }
  return lines.join('\n')
})

function formatURL(url: string): string {
  try {
    const urlObj = new URL(url)
    return urlObj.hostname.replace('www.', '')
  } catch {
    return url
  }
}

function handleClick(url: string) {
  window.open(url, '_blank', 'noopener,noreferrer')
}

// 获取友链分类列表
async function fetchCategories() {
  try {
    const res = await getFriendLinkCategories()
    if (res && res.data) {
      categories.value = Array.isArray(res.data) ? res.data : []
    } else {
      categories.value = []
    }
  } catch (e: any) {
    console.error('获取友链分类失败:', e)
  }
}

async function fetchFriendLinks() {
  loading.value = true
  try {
    const res = await getFriendLinks()
    if (res && res.data) {
      friendLinks.value = Array.isArray(res.data) ? res.data : []
    } else {
      friendLinks.value = []
    }
  } catch (e: any) {
    message.error(e.message || '获取友链列表失败')
  } finally {
    loading.value = false
  }
}

// 按排序顺序获取分类列表
const sortedCategories = computed(() => {
  return [...categories.value].sort((a, b) => b.sort_order - a.sort_order)
})

// 获取指定分类下的友链
function getCategoryLinks(categoryId: number): FriendLink[] {
  return friendLinks.value.filter(link => link.category_id === categoryId)
}

// 获取指定分类下的友链数量
function getCategoryLinksCount(categoryId: number): number {
  return getCategoryLinks(categoryId).length
}

// 获取我的友链信息
async function fetchMyFriendLinkInfo() {
  try {
    const res = await getFriendLinkInfo()
    if (res.data) {
      myFriendLinkInfo.value = {
        name: res.data.name || '菱风叙',
        desc: res.data.desc || '',
        url: res.data.url || '',
        avatar: res.data.avatar || '',
        screenshot: res.data.screenshot || '',
        rss: res.data.rss || ''
      }
    }
  } catch (error: any) {
    console.error('获取我的友链信息失败:', error)
  }
}

// 复制YAML代码
async function handleCopyYaml() {
  try {
    await navigator.clipboard.writeText(yamlCode.value)
    copyYamlSuccess.value = true
    message.success('已复制到剪贴板')
    setTimeout(() => {
      copyYamlSuccess.value = false
    }, 2000)
  } catch (error) {
    message.error('复制失败，请手动复制')
  }
}

// 获取评论列表
async function fetchComments() {
  try {
    const res = await getCommentsByTypeAndTarget(FRIENDLINK_COMMENT_TYPE, FRIENDLINK_TARGET_ID)
    if (res.data) {
      comments.value = res.data
    }
  } catch (error: any) {
    console.error('获取评论失败:', error)
  }
}

// 提交评论
async function handleSubmitComment() {
  if (!authStore.isLoggedIn) {
    message.warning('请先登录')
    return
  }

  if (!commentContent.value.trim()) {
    message.warning('请输入评论内容')
    return
  }

  try {
    submitting.value = true
    const commentData: any = {
      content: commentContent.value,
      comment_type: FRIENDLINK_COMMENT_TYPE,
      target_id: FRIENDLINK_TARGET_ID
    }

    // 如果是回复评论，添加 parent_id
    if (replyToComment.value) {
      commentData.parent_id = replyToComment.value.id
    }

    await createComment(commentData)
    message.success(replyToComment.value ? '回复成功' : '评论成功')
    commentContent.value = ''
    replyToComment.value = null
    replyToUser.value = null
    fetchComments()
  } catch (error: any) {
    message.error(error.message || '评论失败')
  } finally {
    submitting.value = false
  }
}

// 回复评论
function handleReply(parentComment: Comment, targetUser?: Comment) {
  if (!authStore.isLoggedIn) {
    message.warning('请先登录')
    return
  }

  replyToComment.value = parentComment
  replyToUser.value = targetUser || parentComment
  commentContent.value = `@${(targetUser || parentComment).user.nickname} `

  // 滚动到评论框
  nextTick(() => {
    const commentForm = document.querySelector('.comment-form textarea')
    if (commentForm) {
      ;(commentForm as HTMLElement).focus()
      commentForm.scrollIntoView({ behavior: 'smooth', block: 'center' })
    }
  })
}

// 获取回复目标的名称
function getReplyTargetName(reply: Comment, parentComment: Comment): string {
  const match = reply.content.match(/^@(\S+)\s/)
  if (match) {
    return match[1]
  }
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

// 获取友链状态类型
function getLinkStatusType(link: FriendLink) {
  if (link.is_invalid) return 'error'
  if (link.accessible > 0) return 'warning'
  return 'success'
}

// 获取友链状态文本
function getLinkStatusText(link: FriendLink) {
  if (link.is_invalid) return '已失效'
  if (link.accessible > 0) return '访问异常'
  return ''
}

onMounted(() => {
  fetchCategories()
  fetchFriendLinks()
  fetchComments()
  fetchMyFriendLinkInfo()
})
</script>

<style scoped>
.friendlinks-page {
  min-height: calc(100vh - 180px);
  padding: 32px 0;
}

.friendlinks-layout {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.friendlinks-main {
  min-width: 0;
}

.page-header {
  text-align: center;
  margin-bottom: 32px;
}

.page-title {
  font-size: 32px;
  font-weight: 700;
  margin: 0 0 12px 0;
  color: #fff;
  text-shadow: 0 2px 12px rgba(0, 0, 0, 0.5);
}

html.dark .page-title {
  color: #fff;
}

.page-description {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.8);
  text-shadow: 0 1px 8px rgba(0, 0, 0, 0.4);
  margin: 0;
}

html.dark .page-description {
  color: rgba(255, 255, 255, 0.7);
}

.empty-container {
  padding: 80px 20px;
  text-align: center;
}

.categories-container {
  display: flex;
  flex-direction: column;
  gap: 48px;
}

.category-section {
  margin-bottom: 32px;
}

.category-header {
  margin-bottom: 24px;
}

.category-title {
  font-size: 24px;
  font-weight: 700;
  margin: 0 0 8px 0;
  color: #fff;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.5);
}

html.dark .category-title {
  color: #fff;
}

.category-description {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.75);
  text-shadow: 0 1px 6px rgba(0, 0, 0, 0.4);
  margin: 0;
}

html.dark .category-description {
  color: rgba(255, 255, 255, 0.65);
}

.empty-category {
  padding: 40px 20px;
  text-align: center;
}

.friendlinks-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 16px;
}

.friendlink-card {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 16px;
  border: 1px solid rgba(8, 145, 178, 0.1);
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

html.dark .friendlink-card {
  background: rgba(30, 41, 59, 0.8);
  border-color: rgba(56, 189, 248, 0.1);
}

.friendlink-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(8, 145, 178, 0.15);
  border-color: #0891b2;
}

html.dark .friendlink-card:hover {
  box-shadow: 0 8px 24px rgba(56, 189, 248, 0.2);
  border-color: #38bdf8;
}

.card-header {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.link-icon {
  flex-shrink: 0;
  width: 36px;
  height: 36px;
  border-radius: 8px;
  overflow: hidden;
  background: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
}

html.dark .link-icon {
  background: rgba(255, 255, 255, 0.1);
}

.icon-placeholder {
  color: #94a3b8;
}

.link-info {
  flex: 1;
  min-width: 0;
}

.link-name {
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 6px 0;
  color: #1a202c;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

html.dark .link-name {
  color: #e5e5e5;
}

.link-description {
  font-size: 13px;
  color: #64748b;
  line-height: 1.5;
  margin: 0;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  overflow: hidden;
}

html.dark .link-description {
  color: #94a3b8;
}

.card-screenshot {
  border-radius: 8px;
  overflow: hidden;
  background: #f5f5f5;
  margin-top: 4px;
}

html.dark .card-screenshot {
  background: rgba(255, 255, 255, 0.05);
}

.screenshot-image {
  width: 100%;
  display: block;
}

.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 10px;
  margin-top: auto;
  border-top: 1px solid rgba(8, 145, 178, 0.1);
}

html.dark .card-footer {
  border-top-color: rgba(56, 189, 248, 0.1);
}

.link-url {
  font-size: 13px;
  color: #0891b2;
  text-decoration: none;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

html.dark .link-url {
  color: #38bdf8;
}

.link-url:hover {
  text-decoration: underline;
}

.rss-badge {
  flex-shrink: 0;
  margin-left: 8px;
  display: inline-flex;
  align-items: center;
  color: #0891b2;
  cursor: help;
}

html.dark .rss-badge {
  color: #38bdf8;
}

.rss-badge svg {
  display: block;
}

/* 友情链接申请部分 */
.application-section {
  margin-top: 48px;
}

.application-card {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  border: 1px solid rgba(8, 145, 178, 0.1);
}

html.dark .application-card {
  background: rgba(30, 41, 59, 0.8);
  border-color: rgba(56, 189, 248, 0.1);
}

.section-title {
  font-size: 24px;
  font-weight: 700;
  margin: 0 0 20px 0;
  color: #1a202c;
}

html.dark .section-title {
  color: #e5e5e5;
}

.intro-text {
  margin-bottom: 32px;
}

.intro-text p {
  font-size: 15px;
  line-height: 1.8;
  color: #64748b;
  margin: 0;
}

html.dark .intro-text p {
  color: #94a3b8;
}

.rules-section {
  margin-bottom: 32px;
}

.rules-title,
.info-title,
.yaml-title {
  font-size: 18px;
  font-weight: 600;
  margin: 0 0 16px 0;
  color: #1a202c;
}

html.dark .rules-title,
html.dark .info-title,
html.dark .yaml-title {
  color: #e5e5e5;
}

.rules-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.rules-list li {
  position: relative;
  padding-left: 24px;
  margin-bottom: 12px;
  font-size: 14px;
  line-height: 1.8;
  color: #64748b;
}

html.dark .rules-list li {
  color: #94a3b8;
}

.rules-list li::before {
  content: '•';
  position: absolute;
  left: 8px;
  color: #0891b2;
  font-weight: bold;
  font-size: 18px;
}

html.dark .rules-list li::before {
  color: #38bdf8;
}

.rules-list li strong {
  color: #0891b2;
  font-weight: 600;
}

html.dark .rules-list li strong {
  color: #38bdf8;
}

.rules-list li strong .site-link {
  color: inherit;
  text-decoration: none;
  transition: all 0.2s ease;
  border-bottom: 1px solid transparent;
}

.rules-list li strong .site-link:hover {
  border-bottom-color: currentColor;
  opacity: 0.8;
}

.rules-list li code {
  background: rgba(8, 145, 178, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  color: #0891b2;
}

html.dark .rules-list li code {
  background: rgba(56, 189, 248, 0.15);
  color: #38bdf8;
}

.disclaimer {
  margin-top: 16px;
  padding: 12px;
  background: rgba(255, 193, 7, 0.1);
  border-left: 3px solid #ffc107;
  border-radius: 4px;
  font-size: 14px;
  color: #856404;
}

html.dark .disclaimer {
  background: rgba(255, 193, 7, 0.15);
  color: #ffc107;
}

.disclaimer strong {
  color: #ff9800;
}

html.dark .disclaimer strong {
  color: #ffc107;
}

.my-info-section {
  margin-bottom: 32px;
}

.info-grid {
  display: grid;
  gap: 12px;
}

.info-item {
  display: flex;
  align-items: flex-start;
  font-size: 14px;
  line-height: 1.8;
}

.info-label {
  flex-shrink: 0;
  width: 100px;
  font-weight: 600;
  color: #1a202c;
}

html.dark .info-label {
  color: #e5e5e5;
}

.info-value {
  flex: 1;
  color: #64748b;
  word-break: break-all;
}

html.dark .info-value {
  color: #94a3b8;
}

.yaml-section {
  margin-top: 32px;
}

.code-wrapper {
  position: relative;
}

.code-wrapper :deep(.n-code) {
  border-radius: 8px;
  overflow: hidden;
  background: #f5f5f5;
}

html.dark .code-wrapper :deep(.n-code) {
  background: rgba(15, 23, 42, 0.8);
}

.code-wrapper :deep(.n-code pre) {
  margin: 0;
  padding: 16px;
  background: transparent;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.6;
}

.copy-code-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 10;
  opacity: 0;
  transition: opacity 0.3s;
}

.code-wrapper:hover .copy-code-btn {
  opacity: 1;
}

/* 评论区样式 */
.comments-section {
  margin-top: 48px;
}

.comments-card {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  border: 1px solid rgba(8, 145, 178, 0.1);
}

html.dark .comments-card {
  background: rgba(30, 41, 59, 0.8);
  border-color: rgba(56, 189, 248, 0.1);
}

.comments-card .section-title {
  margin-bottom: 24px;
}

.comment-form {
  margin-bottom: 24px;
  background: rgba(255, 255, 255, 0.5);
}

html.dark .comment-form {
  background: rgba(30, 41, 59, 0.5);
}

.empty-comments {
  padding: 40px 20px;
  text-align: center;
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.comment-item {
  padding: 16px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 8px;
  border: 1px solid rgba(8, 145, 178, 0.1);
}

html.dark .comment-item {
  background: rgba(30, 41, 59, 0.5);
  border-color: rgba(56, 189, 248, 0.1);
}

.comment-content {
  flex: 1;
  min-width: 0;
  max-width: 100%;
  /* 允许代码块水平滚动，但不裁剪内容 */
  overflow-x: visible;
  box-sizing: border-box;
}

.comment-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.comment-header strong {
  font-size: 15px;
  font-weight: 600;
  color: #1a202c;
}

html.dark .comment-header strong {
  color: #e5e5e5;
}

.comment-time {
  font-size: 12px;
  color: #94a3b8;
}

.comment-content p {
  margin: 8px 0;
  font-size: 14px;
  line-height: 1.6;
  color: #64748b;
  word-break: break-word;
}

html.dark .comment-content p {
  color: #94a3b8;
}

.comment-actions {
  display: flex;
  gap: 12px;
  margin-top: 8px;
}

.reply-list {
  margin-top: 16px;
  padding-left: 16px;
  border-left: 2px solid rgba(8, 145, 178, 0.2);
}

html.dark .reply-list {
  border-left-color: rgba(56, 189, 248, 0.2);
}

.reply-item {
  padding: 12px;
  margin-bottom: 12px;
  background: rgba(255, 255, 255, 0.3);
  border-radius: 6px;
}

html.dark .reply-item {
  background: rgba(30, 41, 59, 0.3);
}

.reply-content {
  flex: 1;
  min-width: 0;
}

.reply-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
  flex-wrap: wrap;
}

.reply-header strong {
  font-size: 14px;
  font-weight: 600;
  color: #1a202c;
}

html.dark .reply-header strong {
  color: #e5e5e5;
}

.reply-to {
  font-size: 12px;
  color: #0891b2;
}

html.dark .reply-to {
  color: #38bdf8;
}

.reply-content p {
  margin: 6px 0;
  font-size: 13px;
  line-height: 1.6;
  color: #64748b;
}

html.dark .reply-content p {
  color: #94a3b8;
}

/* 响应式 */
@media (max-width: 1024px) {
  .friendlinks-layout {
    padding: 0 16px;
  }
}

@media (max-width: 768px) {
  .friendlinks-page {
    padding: 24px 0;
  }

  .friendlinks-layout {
    padding: 0 16px;
  }

  .page-title {
    font-size: 28px;
  }

  .page-description {
    font-size: 14px;
  }

  .friendlinks-grid {
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .friendlink-card {
    padding: 14px;
  }

  .comment-item {
    padding: 12px;
  }
}

/* 小屏幕移动端优化（小于420px） */
@media (max-width: 420px) {
  .friendlinks-layout {
    padding: 0 12px;
  }

  .comment-item {
    padding: 10px 0;
    margin: 0;
    /* 允许代码块溢出以便滚动查看 */
    overflow-x: auto;
    width: 100%;
    max-width: 100%;
  }

  .comment-content {
    padding: 0 8px;
    /* 允许代码块溢出以便滚动查看 */
    overflow-x: auto;
    width: 100%;
    max-width: 100%;
  }

  .reply-item {
    padding: 8px 0;
    margin: 0;
    /* 允许代码块溢出以便滚动查看 */
    overflow-x: auto;
    width: 100%;
    max-width: 100%;
  }

  .reply-content {
    padding: 0 6px;
    /* 允许代码块溢出以便滚动查看 */
    overflow-x: auto;
    width: 100%;
    max-width: 100%;
  }
}
</style>
