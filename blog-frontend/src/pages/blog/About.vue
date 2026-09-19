<!--
 * @ProjectName: go-vue3-blog
 * @FileName: About.vue
 * @CreateTime: 2026-02-02 11:43:17
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 关于我页面组件（重构版），设计蓝本：anheyu「关于本站」。
 *               纯单栏卡片流布局：头像框 + 大标题 + 基础介绍/打字词 + 技能/生涯
 *               + 座右铭/地理位置 + Markdown 正文 + 文章统计图 + 相册 + 评论区。
 -->
<template>
  <div class="about-page">
    <n-spin :show="loading">
      <!-- 英雄区：玻璃底板（头像框 + 大标题） -->
      <div class="hero-panel">
        <!-- 头像框：居中头像 + 左右浮动技能标签 -->
        <AuthorHero
          :avatar="authorProfile?.author.avatar || ''"
          :fallback-text="fallbackText"
          :left-tags="skills.left"
          :right-tags="skills.right"
        />

        <!-- 页面大标题 -->
        <h1 class="page-title">关于我</h1>
      </div>

      <!-- 作者介绍渐变卡 + 右列（个人信息卡 / MBTI 性格卡上下堆叠，总高与介绍卡相等） -->
      <div class="card-row intro-row">
        <AuthorIntroCard
          :config="intro"
          :author-name="authorName"
          :site-tips="siteTips"
          class="intro-flex"
        />
        <div class="side-stack">
          <SelfInfoCard :config="selfInfo" />
          <PersonalityCard :config="personality" />
        </div>
      </div>

      <!-- 技能 + 职业生涯 -->
      <div v-if="skills.list.length || careers.list.length" class="card-row">
        <SkillsCard v-if="skills.list.length" :config="skills" />
        <CareersCard v-if="careers.list.length" :config="careers" />
      </div>

      <!-- 追求 + 地理位置 -->
      <div class="card-row">
        <MaximCard :config="maxim" />
        <MapLocationCard :config="mapConfig" />
      </div>

      <!-- Markdown 正文 -->
      <AboutCard
        v-if="personalIntroMarkdown"
        tips="关于博主"
        title="自我介绍"
        class="block-card"
      >
        <div class="markdown-body-wrap">
          <MarkdownPreview :content="personalIntroMarkdown" />
        </div>
      </AboutCard>

      <!-- 文章统计图 -->
      <AboutCard tips="坚持记录与输出" title="文章统计" class="block-card">
        <div class="charts-container">
          <div class="chart-item">
            <div class="chart-title">文章发布统计</div>
            <div ref="postPublishChartRef" class="chart-wrapper"></div>
          </div>
          <div class="chart-item">
            <div class="chart-title">TOP10 标签统计</div>
            <div ref="tagChartRef" class="chart-wrapper"></div>
          </div>
        </div>
      </AboutCard>

      <!-- 相册 -->
      <AboutCard v-if="albums.length > 0" tips="用镜头记录生活" title="相册" class="block-card">
        <div class="album-grid">
          <div
            v-for="album in albums"
            :key="album.id"
            class="album-item"
            @click="handleImageClick(album)"
          >
            <n-image
              :src="album.image_url"
              :alt="album.title || '相册照片'"
              object-fit="cover"
              preview-disabled
              class="album-image"
            >
              <template #placeholder>
                <div class="image-placeholder">
                  <n-spin size="small" />
                </div>
              </template>
            </n-image>
            <div v-if="album.title" class="album-title">{{ album.title }}</div>
          </div>
        </div>
      </AboutCard>

      <!-- 图片预览 -->
      <n-image-preview v-model:show="showImagePreview" :src="previewImageUrl" />

      <!-- 评论区 -->
      <AboutCard tips="欢迎留下你的想法" :title="`评论区 (${comments.length})`" class="block-card">
        <!-- 评论表单 -->
        <n-card v-if="authStore.isLoggedIn" class="comment-form">
          <n-alert
            v-if="replyToComment"
            type="info"
            closable
            style="margin-bottom: 12px"
            @close="
              replyToComment = null;
              replyToUser = null;
              commentContent = '';
            "
          >
            正在回复
            <strong>@{{ (replyToUser || replyToComment).user.nickname }}</strong> 的评论
          </n-alert>

          <CommentMarkdownEditor v-model="commentContent" height="250px" :max-length="5000" />
          <div class="comment-submit">
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
      </AboutCard>
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'
// echarts 按需导入
import * as echarts from 'echarts/core'
import { LineChart, BarChart } from 'echarts/charts'
import { TooltipComponent, LegendComponent, GridComponent, MarkLineComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import type { ECharts } from 'echarts/core'

// 注册必需的组件
echarts.use([
  LineChart,
  BarChart,
  TooltipComponent,
  LegendComponent,
  GridComponent,
  MarkLineComponent,
  CanvasRenderer
])
import {
  getAuthorProfile,
  type AuthorProfile,
  getPublicTagStats,
  type TagStat
} from '@/api/blog'
import { getArchives } from '@/api/post'
import { getPublicAboutInfo, getPublicSettings } from '@/api/setting'
import { getPublicAlbums, type Album } from '@/api/album'
import { getCommentsByTypeAndTarget, createComment, deleteComment } from '@/api/comment'
import { formatDate } from '@/utils/format'
import { useAppStore, useAuthStore } from '@/stores'
import MarkdownPreview from '@/components/MarkdownPreview.vue'
import { useMessage } from 'naive-ui'
import AuthorHero from '@/components/about/AuthorHero.vue'
import AboutCard from '@/components/about/AboutCard.vue'
import SkillsCard from '@/components/about/SkillsCard.vue'
import CareersCard from '@/components/about/CareersCard.vue'
import MaximCard from '@/components/about/MaximCard.vue'
import MapLocationCard from '@/components/about/MapLocationCard.vue'
import AuthorIntroCard from '@/components/about/AuthorIntroCard.vue'
import SelfInfoCard from '@/components/about/SelfInfoCard.vue'
import PersonalityCard from '@/components/about/PersonalityCard.vue'
import CommentMarkdownEditor from '@/components/CommentMarkdownEditor.vue'
import CommentContent from '@/components/CommentContent.vue'
import type { Comment } from '@/types/blog'
import {
  parseAboutJson,
  DEFAULT_ABOUT_SITE_TIPS,
  DEFAULT_ABOUT_SKILLS,
  DEFAULT_ABOUT_CAREERS,
  DEFAULT_ABOUT_MAXIM,
  DEFAULT_ABOUT_MAP,
  DEFAULT_ABOUT_INTRO,
  DEFAULT_ABOUT_SELF_INFO,
  DEFAULT_ABOUT_PERSONALITY,
  type AboutSiteTips,
  type AboutSkills,
  type AboutCareers,
  type AboutQuote,
  type AboutMapConfig,
  type AboutIntro,
  type AboutSelfInfo,
  type AboutPersonality
} from '@/types/about'

const router = useRouter()
const appStore = useAppStore()
const authStore = useAuthStore()
const message = useMessage()
const authorProfile = ref<AuthorProfile | null>(null)
const loading = ref(false)

// 评论相关
const comments = ref<Comment[]>([])
const commentContent = ref('')
const replyToComment = ref<Comment | null>(null)
const replyToUser = ref<Comment | null>(null)
const expandedComments = ref<Set<number>>(new Set())
const submitting = ref(false)

// 关于我页面的评论类型
const ABOUT_COMMENT_TYPE = 'about'
const ABOUT_TARGET_ID = 0 // 关于我页面的target_id固定为0

// 关于页卡片配置（settings site 组，JSON 解析失败时回退默认值）
const siteTips = ref<AboutSiteTips>(cloneDefault(DEFAULT_ABOUT_SITE_TIPS))
const skills = ref<AboutSkills>(cloneDefault(DEFAULT_ABOUT_SKILLS))
const careers = ref<AboutCareers>(cloneDefault(DEFAULT_ABOUT_CAREERS))
const maxim = ref<AboutQuote>(cloneDefault(DEFAULT_ABOUT_MAXIM))
const mapConfig = ref<AboutMapConfig>(cloneDefault(DEFAULT_ABOUT_MAP))
const intro = ref<AboutIntro>(cloneDefault(DEFAULT_ABOUT_INTRO))
const selfInfo = ref<AboutSelfInfo>(cloneDefault(DEFAULT_ABOUT_SELF_INFO))
const personality = ref<AboutPersonality>(cloneDefault(DEFAULT_ABOUT_PERSONALITY))

function cloneDefault<T>(value: T): T {
  return JSON.parse(JSON.stringify(value)) as T
}

// 图表相关
const postPublishChartRef = ref<HTMLElement>()
const tagChartRef = ref<HTMLElement>()
let postPublishChart: ECharts | null = null
let tagChart: ECharts | null = null
// 记录上一次是否为移动端布局，用于判断是否跨越断点
let lastIsMobile = window.innerWidth <= 1024

// 统计数据
const archiveStats = ref<Array<{ month: string; count: number }>>([])
const tagStats = ref<TagStat[]>([])

// 个人介绍详情（从API获取，Markdown格式）
const personalIntroMarkdown = ref<string>('')

// 相册列表
const albums = ref<Album[]>([])
const previewImageUrl = ref<string>('')
const showImagePreview = ref(false)

// 计算属性：博主展示信息
const authorName = computed(
  () => authorProfile.value?.author.nickname || authorProfile.value?.author.username || '博主'
)
const fallbackText = computed(() => authorName.value.charAt(0).toUpperCase())

// 获取博主信息
async function fetchAuthorProfile() {
  try {
    loading.value = true
    const res = await getAuthorProfile()
    if (res.data) {
      authorProfile.value = res.data
    }
  } catch (error: any) {
    console.error('获取博主信息失败:', error)
    // 如果获取失败，设置默认值避免显示错误
    authorProfile.value = {
      author: {
        id: 0,
        username: '博主',
        nickname: '博主',
        avatar: '',
        bio: ''
      },
      stats: {
        posts: 0,
        tags: 0,
        categories: 0
      }
    }
  } finally {
    loading.value = false
  }
}

// 获取关于页卡片配置（site 组设置）
async function fetchSiteSettings() {
  try {
    const res = await getPublicSettings()
    if (res.data) {
      siteTips.value = parseAboutJson(res.data.about_site_tips, DEFAULT_ABOUT_SITE_TIPS)
      skills.value = parseAboutJson(res.data.about_skills, DEFAULT_ABOUT_SKILLS)
      careers.value = parseAboutJson(res.data.about_careers, DEFAULT_ABOUT_CAREERS)
      maxim.value = parseAboutJson(res.data.about_maxim, DEFAULT_ABOUT_MAXIM)
      mapConfig.value = parseAboutJson(res.data.about_map, DEFAULT_ABOUT_MAP)
      intro.value = parseAboutJson(res.data.about_intro, DEFAULT_ABOUT_INTRO)
      selfInfo.value = parseAboutJson(res.data.about_self_info, DEFAULT_ABOUT_SELF_INFO)
      personality.value = parseAboutJson(res.data.about_personality, DEFAULT_ABOUT_PERSONALITY)
    }
  } catch (error) {
    console.error('获取关于页配置失败:', error)
  }
}

// 获取关于我信息
async function fetchAboutInfo() {
  try {
    const res = await getPublicAboutInfo()
    if (res.data && res.data.content) {
      personalIntroMarkdown.value = res.data.content
    } else {
      // 如果没有内容，设置为空
      personalIntroMarkdown.value = ''
    }
  } catch (error) {
    console.error('获取关于我信息失败:', error)
    // 如果获取失败，使用空字符串
    personalIntroMarkdown.value = ''
  }
}

// 获取相册数据
async function fetchAlbums() {
  try {
    const res = await getPublicAlbums()
    if (res.data) {
      albums.value = res.data
    }
  } catch (error) {
    console.error('获取相册数据失败:', error)
    albums.value = []
  }
}

// 处理图片点击
function handleImageClick(album: Album) {
  previewImageUrl.value = album.image_url
  showImagePreview.value = true
}

// 获取归档统计数据
async function fetchArchiveStats() {
  try {
    const res = await getArchives()
    if (res.data) {
      archiveStats.value = res.data.map((item: any) => ({
        month: item.month,
        count: Number(item.count)
      }))
      nextTick(() => {
        // 延迟初始化，确保DOM完全渲染
        setTimeout(() => {
          initPostPublishChart()
        }, 100)
      })
    }
  } catch (error) {
    console.error('获取归档统计失败:', error)
  }
}

// 获取评论列表
async function fetchComments() {
  try {
    const res = await getCommentsByTypeAndTarget(ABOUT_COMMENT_TYPE, ABOUT_TARGET_ID)
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
      comment_type: ABOUT_COMMENT_TYPE,
      target_id: ABOUT_TARGET_ID
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

// 获取标签统计数据
async function fetchTagStats() {
  try {
    const res = await getPublicTagStats()
    if (res.data) {
      tagStats.value = res.data.slice(0, 10) // 只取TOP10
      nextTick(() => {
        // 延迟初始化，确保DOM完全渲染
        setTimeout(() => {
          initTagChart()
        }, 100)
      })
    }
  } catch (error) {
    console.error('获取标签统计失败:', error)
  }
}

// ===== 图表配置生成（折线/柱状共用部分抽取，减少重复） =====

interface ChartViewport {
  isMobile: boolean
  isSmallMobile: boolean
}

function getViewport(): ChartViewport {
  return {
    isMobile: window.innerWidth <= 1024,
    isSmallMobile: window.innerWidth <= 767
  }
}

// 坐标轴主题（明暗两套）
function buildAxisTheme(isDark: boolean) {
  return {
    axisLine: { lineStyle: { color: isDark ? '#64748b' : '#cbd5e1' } },
    axisTick: { show: false },
    axisLabel: { color: isDark ? '#e5e7eb' : '#64748b' }
  }
}

// 平均值标记线（紫色虚线徽标）
function buildAverageMarkLine(average: number | string, isSmallMobile: boolean) {
  return {
    silent: true,
    data: [
      {
        yAxis: average,
        name: '平均值',
        label: {
          formatter: `平均: ${average}`,
          position: isSmallMobile ? 'insideEndTop' : 'end',
          backgroundColor: 'rgba(154, 96, 180, 0.8)',
          color: '#fff',
          padding: isSmallMobile ? [3, 8] : [4, 10],
          borderRadius: 4,
          fontSize: isSmallMobile ? 10 : 11,
          distance: isSmallMobile ? [0, -5] : [10, 0]
        },
        lineStyle: {
          type: 'dashed',
          color: '#9a60b4',
          width: 2
        }
      }
    ]
  }
}

// 通用 tooltip / grid
function buildTooltipAndGrid(viewport: ChartViewport, gridBottom: number) {
  return {
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        const param = params[0]
        return `${param.name}<br/>${param.seriesName}: ${param.value}`
      }
    },
    grid: {
      left: viewport.isSmallMobile ? 45 : viewport.isMobile ? 50 : 60,
      right: viewport.isSmallMobile ? 18 : viewport.isMobile ? 28 : 80,
      top: 50,
      bottom: gridBottom,
      containLabel: false
    }
  }
}

// 初始化 echarts 实例（复用已有实例则先 resize，option 全量替换）
function initChartInstance(
  el: HTMLElement | undefined,
  chart: ECharts | null,
  option: echarts.EChartsCoreOption
): ECharts | null {
  if (!el) return chart
  let instance = chart
  if (!instance) {
    instance = echarts.init(el)
  } else {
    instance.resize()
  }
  instance.setOption(option, true)
  return instance
}

// 初始化文章发布统计图（折线图）
function initPostPublishChart() {
  const viewport = getViewport()
  const countsReady = archiveStats.value.length > 0

  // 处理数据：格式化月份，补全缺失的月份
  const dataMap = new Map<string, number>()
  archiveStats.value.forEach(item => {
    const month = item.month.substring(0, 7) // YYYY-MM
    dataMap.set(month, Number(item.count))
  })

  // 生成从最早到最晚的连续月份数组
  const months: string[] = []
  if (countsReady) {
    const sortedMonths = Array.from(dataMap.keys()).sort()
    const [startYear, startMonthNum] = sortedMonths[0].split('-').map(Number)
    const [endYear, endMonthNum] = sortedMonths[sortedMonths.length - 1].split('-').map(Number)

    let currentYear = startYear
    let currentMonth = startMonthNum
    while (currentYear < endYear || (currentYear === endYear && currentMonth <= endMonthNum)) {
      const monthStr = `${currentYear}-${String(currentMonth).padStart(2, '0')}`
      months.push(monthStr)
      currentMonth++
      if (currentMonth > 12) {
        currentMonth = 1
        currentYear++
      }
    }
  }

  const counts = months.map(month => dataMap.get(month) || 0)

  // 计算平均值
  const total = counts.reduce((sum, count) => sum + count, 0)
  const average = counts.length > 0 ? (total / counts.length).toFixed(2) : 0

  const isDark = appStore.theme === 'dark'

  // 根据数据点数量判断是否需要旋转标签（"2025-10"标签较长，超过3个点桌面端也旋转）
  const needRotate = viewport.isSmallMobile ? true : months.length > 3
  const rotateAngle = viewport.isSmallMobile ? 45 : needRotate ? 60 : 0
  const gridBottom = viewport.isSmallMobile
    ? 50
    : viewport.isMobile
      ? needRotate
        ? 80
        : 55
      : needRotate
        ? 90
        : 60

  const axisTheme = buildAxisTheme(isDark)

  const option: echarts.EChartsCoreOption = {
    ...buildTooltipAndGrid(viewport, gridBottom),
    xAxis: {
      type: 'category',
      data: months,
      ...axisTheme,
      axisLabel: {
        ...axisTheme.axisLabel,
        rotate: rotateAngle,
        fontSize: viewport.isSmallMobile ? 10 : needRotate ? 11 : 12,
        interval: 0,
        overflow: 'break',
        width: viewport.isSmallMobile ? 50 : needRotate ? 70 : 60,
        margin: viewport.isSmallMobile ? 8 : needRotate ? 15 : 0
      }
    },
    yAxis: {
      type: 'value',
      minInterval: 1,
      axisLine: { show: false },
      axisTick: { show: false },
      splitLine: {
        lineStyle: {
          color: isDark ? '#1e293b' : '#e5e7eb'
        }
      },
      axisLabel: { color: isDark ? '#e5e7eb' : '#64748b' }
    },
    series: [
      {
        name: '文章数',
        type: 'line',
        data: counts,
        smooth: true,
        symbol: 'circle',
        symbolSize: 8,
        lineStyle: {
          color: '#0891b2',
          width: 2
        },
        itemStyle: {
          color: '#0891b2'
        },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(8, 145, 178, 0.3)' },
              { offset: 1, color: 'rgba(8, 145, 178, 0.05)' }
            ]
          }
        },
        markLine: buildAverageMarkLine(average, viewport.isSmallMobile)
      }
    ]
  }

  postPublishChart = initChartInstance(postPublishChartRef.value, postPublishChart, option)
}

// 初始化标签统计图（柱状图）
function initTagChart() {
  const viewport = getViewport()

  // 按数量降序排序
  const sortedTags = [...tagStats.value].sort((a, b) => b.value - a.value)
  const tagNames = sortedTags.map(t => t.name)
  const tagCounts = sortedTags.map(t => t.value)

  // 计算平均值
  const total = tagCounts.reduce((sum, count) => sum + count, 0)
  const average = tagCounts.length > 0 ? (total / tagCounts.length).toFixed(1) : 0

  const isDark = appStore.theme === 'dark'

  // 标签名称通常较长，统一使用60度旋转避免重叠
  const gridBottom = viewport.isSmallMobile ? 70 : viewport.isMobile ? 95 : 90

  const axisTheme = buildAxisTheme(isDark)

  const option: echarts.EChartsCoreOption = {
    ...buildTooltipAndGrid(viewport, gridBottom),
    xAxis: {
      type: 'category',
      data: tagNames,
      ...axisTheme,
      axisLabel: {
        ...axisTheme.axisLabel,
        rotate: 60,
        fontSize: viewport.isSmallMobile ? 9 : 11,
        interval: 0,
        overflow: 'break',
        width: viewport.isSmallMobile ? 40 : 55,
        margin: viewport.isSmallMobile ? 10 : 15
      }
    },
    yAxis: {
      type: 'value',
      minInterval: 1,
      axisLine: { show: false },
      axisTick: { show: false },
      splitLine: {
        lineStyle: {
          color: isDark ? '#1e293b' : '#e5e7eb'
        }
      },
      axisLabel: { color: isDark ? '#e5e7eb' : '#64748b' }
    },
    series: [
      {
        name: '文章数',
        type: 'bar',
        data: tagCounts,
        itemStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: '#0891b2' },
              { offset: 1, color: '#06b6d4' }
            ]
          },
          borderRadius: [4, 4, 0, 0]
        },
        markLine: buildAverageMarkLine(average, viewport.isSmallMobile)
      }
    ]
  }

  tagChart = initChartInstance(tagChartRef.value, tagChart, option)
}

// 监听主题变化，重新渲染图表
watch(
  () => appStore.theme,
  () => {
    nextTick(() => {
      if (archiveStats.value.length > 0) initPostPublishChart()
      if (tagStats.value.length > 0) initTagChart()
    })
  }
)

// 防抖函数
function debounce(func: Function, wait: number) {
  let timeout: ReturnType<typeof setTimeout> | null = null
  return function executedFunction(...args: any[]) {
    const later = () => {
      timeout = null
      func(...args)
    }
    if (timeout) clearTimeout(timeout)
    timeout = setTimeout(later, wait)
  }
}

// 窗口大小变化时调整图表
const handleResize = debounce(() => {
  const nowIsMobile = window.innerWidth <= 1024

  nextTick(() => {
    // 延迟执行，确保DOM和布局完全更新
    setTimeout(() => {
      const crossedBreakpoint = lastIsMobile !== nowIsMobile

      if (crossedBreakpoint) {
        // 从移动端切到桌面端（或反之）时，彻底销毁并重建图表，避免尺寸计算异常
        if (postPublishChart) {
          postPublishChart.dispose()
          postPublishChart = null
        }
        if (tagChart) {
          tagChart.dispose()
          tagChart = null
        }

        // 布局稳定后，根据当前数据重新初始化
        if (archiveStats.value.length > 0) {
          initPostPublishChart()
        }
        if (tagStats.value.length > 0) {
          initTagChart()
        }
      } else {
        // 同一断点内只做自适应
        if (postPublishChart) {
          postPublishChart.resize()
        }
        if (tagChart) {
          tagChart.resize()
        }
      }

      // 记录当前断点状态
      lastIsMobile = nowIsMobile
    }, 200)
  })
}, 150)

// 监听媒体查询变化（响应式断点）
let mediaQueryList: MediaQueryList | null = null

function handleMediaChange() {
  nextTick(() => {
    handleResize()
  })
}

onMounted(() => {
  fetchAuthorProfile()
  fetchSiteSettings()
  fetchAboutInfo()
  fetchAlbums()
  fetchArchiveStats()
  fetchTagStats()
  fetchComments()

  // 监听窗口resize
  window.addEventListener('resize', handleResize)

  // 监听媒体查询变化（用于响应式断点）
  if (window.matchMedia) {
    mediaQueryList = window.matchMedia('(max-width: 1024px)')
    if (mediaQueryList.addEventListener) {
      mediaQueryList.addEventListener('change', handleMediaChange)
    } else {
      // 兼容旧版浏览器
      mediaQueryList.addListener(handleMediaChange)
    }
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  if (mediaQueryList) {
    if (mediaQueryList.removeEventListener) {
      mediaQueryList.removeEventListener('change', handleMediaChange)
    } else {
      // 兼容旧版浏览器
      mediaQueryList.removeListener(handleMediaChange)
    }
  }
  postPublishChart?.dispose()
  tagChart?.dispose()
})
</script>

<style scoped>
.about-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px 20px;
  position: relative;
  z-index: 1;
}

/* 英雄区玻璃底板：头像框 + 标题 + 介绍区，保证壁纸上文字对比度 */
.hero-panel {
  padding: 32px 24px 28px;
  margin-bottom: 24px;
  text-align: center;
  background: rgba(255, 255, 255, 0.62);
  backdrop-filter: blur(24px) saturate(180%);
  -webkit-backdrop-filter: blur(24px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.5);
  border-radius: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
}

html.dark .hero-panel {
  background: rgba(30, 41, 59, 0.6);
  border-color: rgba(255, 255, 255, 0.08);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

/* 头像框与标题 */
.page-title {
  margin: 0.625rem 0 0.25rem;
  font-size: 2.5rem;
  font-weight: 700;
  line-height: 1;
  text-align: center;
  color: #1a202c;
}

html.dark .page-title {
  color: #e5e5e5;
}

/* 成对卡片行 */
.card-row {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  align-items: stretch;
  margin-bottom: 24px;
}

.card-row > * {
  flex: 1;
  min-width: 0;
}

/* 介绍卡略宽于个人信息卡 */
.intro-row .intro-flex {
  flex: 1.5;
}

/* 右列上下堆叠：两张卡片等分高度，总和与左侧打招呼卡相等 */
.side-stack {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
}

.side-stack > * {
  flex: 1;
  min-height: 0;
}

/* 独占一行卡片 */
.block-card {
  margin-bottom: 24px;
}

/* Markdown 正文 */
.markdown-body-wrap {
  font-size: 15px;
  line-height: 1.9;
  color: #475569;
}

html.dark .markdown-body-wrap {
  color: #cbd5e1;
}

/* 文章统计图 */
.charts-container {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
  width: 100%;
  box-sizing: border-box;
}

.chart-item {
  display: flex;
  flex-direction: column;
  gap: 12px;
  background: rgba(8, 145, 178, 0.03);
  border: 1px solid rgba(8, 145, 178, 0.08);
  border-radius: 12px;
  padding: 16px;
  transition: all 0.3s ease;
  overflow: hidden;
  box-sizing: border-box;
  width: 100%;
}

.chart-item:hover {
  background: rgba(8, 145, 178, 0.06);
  transform: translateY(-2px);
}

html.dark .chart-item {
  background: rgba(56, 189, 248, 0.05);
  border-color: rgba(56, 189, 248, 0.12);
}

html.dark .chart-item:hover {
  background: rgba(56, 189, 248, 0.1);
}

.chart-title {
  font-size: 16px;
  font-weight: 600;
  color: #1a202c;
  text-align: center;
  letter-spacing: 0.5px;
  padding-bottom: 10px;
  border-bottom: 2px solid rgba(8, 145, 178, 0.1);
}

html.dark .chart-title {
  color: #e5e5e5;
  border-bottom-color: rgba(56, 189, 248, 0.2);
}

.chart-wrapper {
  width: 100%;
  height: 340px;
  min-height: 340px;
  overflow: hidden;
  position: relative;
  box-sizing: border-box;
}

/* 相册 */
.album-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
}

.album-item {
  position: relative;
  cursor: pointer;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s ease;
  background: rgba(0, 0, 0, 0.02);
}

.album-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

html.dark .album-item {
  background: rgba(255, 255, 255, 0.05);
}

html.dark .album-item:hover {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.5);
}

.album-image {
  width: 100%;
  height: 200px;
  object-fit: cover;
  border-radius: 12px;
}

.album-image :deep(img) {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-placeholder {
  width: 100%;
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.05);
  border-radius: 12px;
}

html.dark .image-placeholder {
  background: rgba(255, 255, 255, 0.05);
}

.album-title {
  padding: 12px;
  font-size: 14px;
  font-weight: 500;
  color: #1a202c;
  text-align: center;
  background: rgba(255, 255, 255, 0.9);
}

html.dark .album-title {
  color: #e5e5e5;
  background: rgba(30, 41, 59, 0.9);
}

/* 评论区 */
.comment-form {
  margin-bottom: 24px;
  background: rgba(255, 255, 255, 0.5);
}

html.dark .comment-form {
  background: rgba(30, 41, 59, 0.5);
}

.comment-submit {
  margin-top: 12px;
  text-align: right;
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

.reply-item:last-child {
  margin-bottom: 0;
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

/* 响应式：平板与移动端 */
@media (max-width: 1024px) {
  .charts-container {
    grid-template-columns: 1fr;
    gap: 20px;
  }

  .chart-wrapper {
    height: 380px;
  }
}

@media (max-width: 767px) {
  .about-page {
    padding: 16px 12px;
  }

  .hero-panel {
    padding: 24px 14px 20px;
    margin-bottom: 16px;
  }

  .page-title {
    font-size: 2rem;
  }

  .author-name {
    font-size: 28px;
  }

  .card-row {
    flex-direction: column;
    gap: 16px;
    margin-bottom: 16px;
  }

  .card-row > * {
    width: 100%;
    flex: none;
  }

  .block-card {
    margin-bottom: 16px;
  }

  .charts-container {
    gap: 16px;
  }

  .chart-item {
    padding: 12px 8px;
  }

  .chart-title {
    font-size: 15px;
  }

  .chart-wrapper {
    height: 320px;
  }

  .album-grid {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 12px;
  }

  .album-image {
    height: 150px;
  }

  .image-placeholder {
    height: 150px;
  }

  .comment-item {
    padding: 12px;
  }
}

/* 小屏幕移动端优化（小于420px） */
@media (max-width: 420px) {
  .comment-item {
    padding: 10px 0;
    margin: 0;
    overflow-x: auto;
    width: 100%;
    max-width: 100%;
  }

  .comment-content {
    padding: 0 8px;
    overflow-x: auto;
    width: 100%;
    max-width: 100%;
  }

  .reply-item {
    padding: 8px 0;
    margin: 0;
    overflow-x: auto;
    width: 100%;
    max-width: 100%;
  }

  .reply-content {
    padding: 0 6px;
    overflow-x: auto;
    width: 100%;
    max-width: 100%;
  }
}
</style>
