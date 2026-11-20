<!--
 * @ProjectName: go-vue3-blog
 * @FileName: Dashboard.vue
 * @CreateTime: 2026-02-02 11:43:17
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 管理后台仪表盘页面组件，展示网站统计数据
 -->
<template>
  <div class="dashboard-page">
    <!-- 统计卡片 -->
    <n-grid cols="1 s:2 m:4" :x-gap="16" :y-gap="16" responsive="screen">
      <n-gi v-for="item in statCards" :key="item.key">
        <n-card class="stat-card" :class="`stat-card--${item.key}`" hoverable @click="router.push(item.to)">
          <n-skeleton v-if="loading" width="60%" :sharp="false" />
          <n-statistic v-else :label="item.label" :value="item.value">
            <template #prefix>
              <div class="stat-icon">
                <n-icon :component="item.icon" :size="20" />
              </div>
            </template>
          </n-statistic>
        </n-card>
      </n-gi>
    </n-grid>

    <!-- 图表区域 -->
    <n-grid cols="1 m:2" :x-gap="16" :y-gap="16" style="margin-top: 24px" responsive="screen">
      <n-gi>
        <n-card title="文章分类统计">
          <div v-show="categoryStats.length > 0" ref="categoryChartRef" :style="chartStyle"></div>
          <n-empty
            v-if="categoryStats.length === 0"
            description="暂无分类数据"
            :style="{ height: chartStyle.height }"
          />
        </n-card>
      </n-gi>
      <n-gi>
        <n-card title="最近7天访问量">
          <div v-show="visitStats.length > 0" ref="visitChartRef" :style="chartStyle"></div>
          <n-empty
            v-if="visitStats.length === 0"
            description="暂无访问数据"
            :style="{ height: chartStyle.height }"
          />
        </n-card>
      </n-gi>
    </n-grid>

    <!-- 快捷操作 -->
    <n-card title="快捷操作" style="margin-top: 24px">
      <div class="quick-actions">
        <div
          v-for="action in quickActions"
          :key="action.label"
          class="quick-action-item"
          @click="router.push(action.to)"
        >
          <div class="quick-action-icon" :style="{ background: action.bg, color: action.color }">
            <n-icon :component="action.icon" :size="22" />
          </div>
          <span class="quick-action-label">{{ action.label }}</span>
          <n-badge v-if="action.badge" :value="action.badge" :max="99" />
        </div>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
// echarts 按需导入
import * as echarts from 'echarts/core'
import { PieChart, LineChart } from 'echarts/charts'
import { TooltipComponent, LegendComponent, GridComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import type { ECharts } from 'echarts/core'

// 注册必需的组件
echarts.use([PieChart, LineChart, TooltipComponent, LegendComponent, GridComponent, CanvasRenderer])
import {
  DocumentTextOutline,
  PeopleOutline,
  ChatbubblesOutline,
  EyeOutline,
  AddOutline,
  FolderOutline,
  PricetagsOutline,
  ImagesOutline,
  ChatboxOutline,
  CheckmarkDoneOutline,
  CloudUploadOutline
} from '@vicons/ionicons5'
import { useAppStore, useAuthStore } from '@/stores'
import { getDashboardStats, getCategoryStats, getVisitStats } from '@/api'
import type { DashboardStats, CategoryStat, VisitStat } from '@/api'

const router = useRouter()
const appStore = useAppStore()
const authStore = useAuthStore()
const message = useMessage()

const stats = ref<DashboardStats>({
  posts: 0,
  users: 0,
  comments: 0,
  views: 0
})

// 统计卡片配置：可点击跳转到对应管理页
const statCards = computed(() => [
  {
    key: 'posts',
    label: '文章总数',
    value: stats.value.posts.toLocaleString(),
    to: '/admin/posts',
    icon: DocumentTextOutline
  },
  {
    key: 'users',
    label: '用户总数',
    value: stats.value.users.toLocaleString(),
    to: '/admin/users',
    icon: PeopleOutline
  },
  {
    key: 'comments',
    label: '评论总数',
    value: stats.value.comments.toLocaleString(),
    to: '/admin/comments',
    icon: ChatbubblesOutline
  },
  {
    key: 'views',
    label: '总浏览量',
    value: stats.value.views.toLocaleString(),
    to: '/',
    icon: EyeOutline
  }
])

// 快捷操作入口：根据角色过滤（部分入口仅 super_admin 可见），badge 预留待审数量展示
interface QuickAction {
  label: string
  to: string
  icon: typeof AddOutline
  bg: string
  color: string
  superOnly: boolean
  badge?: number
}

const allQuickActions: QuickAction[] = [
  { label: '写文章', to: '/admin/posts', icon: AddOutline, bg: 'rgba(8, 145, 178, 0.12)', color: '#0891b2', superOnly: false },
  { label: '上传 Markdown', to: '/admin/posts', icon: CloudUploadOutline, bg: 'rgba(5, 150, 105, 0.12)', color: '#059669', superOnly: false },
  { label: '管理分类', to: '/admin/categories', icon: FolderOutline, bg: 'rgba(217, 119, 6, 0.12)', color: '#d97706', superOnly: false },
  { label: '管理标签', to: '/admin/tags', icon: PricetagsOutline, bg: 'rgba(2, 132, 199, 0.12)', color: '#0284c7', superOnly: false },
  { label: '评论管理', to: '/admin/comments', icon: ChatbubblesOutline, bg: 'rgba(124, 58, 237, 0.12)', color: '#7c3aed', superOnly: false },
  { label: '说说管理', to: '/admin/moments', icon: ChatboxOutline, bg: 'rgba(236, 72, 153, 0.12)', color: '#ec4899', superOnly: false },
  { label: '图片管理', to: '/admin/images', icon: ImagesOutline, bg: 'rgba(13, 148, 136, 0.12)', color: '#0d9488', superOnly: true },
  { label: '审批管理', to: '/admin/post-revisions', icon: CheckmarkDoneOutline, bg: 'rgba(239, 68, 68, 0.12)', color: '#ef4444', superOnly: true }
]

const quickActions = computed(() =>
  allQuickActions.filter(action => !action.superOnly || authStore.isSuperAdmin)
)

const categoryStats = ref<CategoryStat[]>([])
const visitStats = ref<VisitStat[]>([])

const categoryChartRef = ref<HTMLElement>()
const visitChartRef = ref<HTMLElement>()
let categoryChart: ECharts | null = null
let visitChart: ECharts | null = null

const loading = ref(false)
const isMobile = ref(false)

const chartStyle = computed(() => ({
  width: '100%',
  height: isMobile.value ? '320px' : '350px'
}))

// 检测移动设备
function checkMobile() {
  isMobile.value = window.innerWidth <= 768
}

// 获取统计数据
async function fetchStats() {
  try {
    const res = await getDashboardStats()
    if (res.data) {
      stats.value = res.data
    }
  } catch (error) {
    message.error('获取统计数据失败')
    console.error(error)
  }
}

// 获取分类统计数据
async function fetchCategoryStats() {
  try {
    const res = await getCategoryStats()
    if (res.data) {
      categoryStats.value = res.data
      // 更新分类图表
      if (res.data.length > 0) {
        nextTick(() => {
          initCategoryChart()
        })
      }
    }
  } catch (error) {
    message.error('获取分类统计失败')
    console.error(error)
  }
}

// 获取最近访问统计数据（最近 7 天）
async function fetchVisitStats() {
  try {
    const res = await getVisitStats(7)
    if (res.data) {
      visitStats.value = res.data
      if (res.data.length > 0) {
        nextTick(() => {
          initVisitChart()
        })
      }
    }
  } catch (error) {
    message.error('获取访问统计失败')
    console.error(error)
  }
}

onMounted(async () => {
  // 检测移动设备
  checkMobile()

  // 并行获取所有统计数据，缩短首屏时间
  loading.value = true
  await Promise.all([fetchStats(), fetchCategoryStats(), fetchVisitStats()])
  loading.value = false

  // 监听窗口大小变化
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  categoryChart?.dispose()
  visitChart?.dispose()
})

function handleResize() {
  checkMobile()
  // 防抖：拖动窗口时避免频繁重绘图表
  if (resizeTimer) window.clearTimeout(resizeTimer)
  resizeTimer = window.setTimeout(() => {
    categoryChart?.resize()
    visitChart?.resize()
  }, 200)
}

let resizeTimer: number | undefined

// 初始化分类统计饼图
function initCategoryChart() {
  if (!categoryChartRef.value) return

  if (!categoryChart) {
    categoryChart = echarts.init(categoryChartRef.value)
  }

  const isDark = appStore.theme === 'dark'

  // 默认颜色方案：感知均匀色板，分类未设色/撞色时兜底，保证扇区间区分度
  const defaultColors = [
    '#5470c6',
    '#91cc75',
    '#fac858',
    '#ee6666',
    '#73c0de',
    '#3ba272',
    '#fc8452',
    '#9a60b4'
  ]

  // 按数量降序排列，图例与扇区顺序一致，小分类沉底
  const sorted = [...categoryStats.value].sort((a, b) => b.value - a.value)

  // 转换数据格式：已设色的沿用分类色；未设色或与前面分类撞色时用兜底色板
  const usedColors = new Set<string>()
  const chartData = sorted.map((item, index) => {
    let color = item.color?.trim()
    if (!color || usedColors.has(color.toLowerCase())) {
      color = defaultColors[index % defaultColors.length]
    }
    usedColors.add(color.toLowerCase())
    return {
      value: item.value,
      name: item.name,
      itemStyle: {
        color
      }
    }
  })

  // 获取容器实际宽度
  const containerWidth = categoryChartRef.value.clientWidth
  // 判断屏幕模式：超小屏 < 380, 窮屏 < 500, 正常
  const isExtraSmall = containerWidth < 380
  const isNarrow = containerWidth < 500

  // 超小屏：图例在底部水平排列；窄屏：图例在右侧垂直排列；正常：同窄屏
  // formatter 带上数量，图例即是迷你统计表
  const legendFormatter = (name: string) => {
    const found = sorted.find(item => item.name === name)
    return found ? `${name}  ${found.value}` : name
  }
  const legendConfig = isExtraSmall
    ? {
        type: 'scroll',
        orient: 'horizontal',
        bottom: 0,
        left: 'center',
        padding: [8, 0, 0, 0],
        formatter: legendFormatter,
        textStyle: {
          color: appStore.theme === 'dark' ? '#fff' : '#333',
          fontSize: 11,
          overflow: 'truncate',
          width: 80
        },
        pageIconColor: appStore.theme === 'dark' ? '#fff' : '#333',
        pageTextStyle: {
          color: appStore.theme === 'dark' ? '#fff' : '#333'
        },
        tooltip: {
          show: true
        }
      }
    : {
        type: 'scroll',
        orient: 'vertical',
        right: '2%',
        top: 'middle',
        padding: [5, 10],
        formatter: legendFormatter,
        textStyle: {
          color: appStore.theme === 'dark' ? '#fff' : '#333',
          fontSize: 12,
          overflow: 'truncate',
          width: isNarrow ? 100 : 160
        },
        pageIconColor: appStore.theme === 'dark' ? '#fff' : '#333',
        pageTextStyle: {
          color: appStore.theme === 'dark' ? '#fff' : '#333'
        },
        tooltip: {
          show: true
        }
      }

  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: legendConfig,
    series: [
      {
        name: '文章分类',
        type: 'pie',
        // 超小屏：更小的饼图居中；窄屏：偏左；正常：同窄屏
        radius: isExtraSmall ? ['20%', '38%'] : isNarrow ? ['25%', '45%'] : ['30%', '55%'],
        center: isExtraSmall ? ['50%', '40%'] : isNarrow ? ['38%', '50%'] : ['32%', '50%'],
        avoidLabelOverlap: true,
        itemStyle: {
          borderRadius: 8,
          // 扇区分隔色跟随主题，避免深色模式下白边刺眼
          borderColor: isDark ? '#1e293b' : '#fff',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'outside',
          formatter: '{b}',
          fontSize: 12,
          overflow: 'truncate',
          width: 120
        },
        emphasis: {
          label: {
            show: false,
            fontSize: 14,
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false,
          length: 8,
          length2: 5
        },
        data: chartData
      }
    ]
  }

  categoryChart.setOption(option, true) // 使用 true 强制重绘，避免旧配置残留
}

// 初始化最近访问量折线图
function initVisitChart() {
  if (!visitChartRef.value) return

  if (!visitChart) {
    visitChart = echarts.init(visitChartRef.value)
  }

  const dates = visitStats.value.map(item => item.date.slice(5)) // 显示 MM-DD
  const counts = visitStats.value.map(item => item.count)

  const isDark = appStore.theme === 'dark'

  const option = {
    tooltip: {
      trigger: 'axis'
    },
    grid: {
      left: 40,
      right: 24,
      top: 24,
      bottom: 32
    },
    xAxis: {
      type: 'category',
      data: dates,
      boundaryGap: false,
      axisLine: {
        lineStyle: {
          color: isDark ? '#64748b' : '#cbd5e1'
        }
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        color: isDark ? '#e5e7eb' : '#64748b'
      }
    },
    yAxis: {
      type: 'value',
      minInterval: 1,
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      splitLine: {
        lineStyle: {
          color: isDark ? '#1e293b' : '#e5e7eb'
        }
      },
      axisLabel: {
        color: isDark ? '#e5e7eb' : '#64748b'
      }
    },
    series: [
      {
        name: '访问量',
        type: 'line',
        data: counts,
        smooth: true,
        symbol: 'circle',
        symbolSize: 6,
        lineStyle: {
          width: 3,
          color: '#0ea5e9'
        },
        itemStyle: {
          color: '#0ea5e9'
        },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(14,165,233,0.35)' },
              { offset: 1, color: 'rgba(14,165,233,0.02)' }
            ]
          }
        }
      }
    ]
  }

  visitChart.setOption(option)
}
</script>

<style scoped>
.dashboard-page {
  padding: 0;
  position: relative;
  z-index: 1;
}

.dashboard-page h1 {
  font-size: 28px;
}

/* 移动端样式 */
@media (max-width: 768px) {
  .dashboard-page h1 {
    font-size: 22px;
    margin-bottom: 16px !important;
  }
}

/* 优化统计卡片样式 */
.dashboard-page :deep(.n-card) {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.5);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  border-radius: 16px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 上浮动效只给统计卡和快捷入口，图表卡片保持稳定避免晃眼 */
.dashboard-page :deep(.n-card:hover) {
  box-shadow: 0 16px 40px rgba(0, 0, 0, 0.12);
  border-color: rgba(8, 145, 178, 0.3);
}

.dashboard-page :deep(.n-card.stat-card:hover) {
  transform: translateY(-4px);
}

/* 深色模式卡片 */
html.dark .dashboard-page :deep(.n-card) {
  background: rgba(30, 41, 59, 0.7);
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

html.dark .dashboard-page :deep(.n-card:hover) {
  box-shadow: 0 16px 40px rgba(0, 0, 0, 0.5);
  border-color: rgba(56, 189, 248, 0.3);
}

/* 统计卡片：可点击跳转，彩色图标块区分指标 */
.stat-card {
  cursor: pointer;
}

.stat-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-card--posts .stat-icon {
  background: rgba(8, 145, 178, 0.12);
  color: #0891b2;
}

.stat-card--users .stat-icon {
  background: rgba(5, 150, 105, 0.12);
  color: #059669;
}

.stat-card--comments .stat-icon {
  background: rgba(217, 119, 6, 0.12);
  color: #d97706;
}

.stat-card--views .stat-icon {
  background: rgba(2, 132, 199, 0.12);
  color: #0284c7;
}

html.dark .stat-card--posts .stat-icon {
  background: rgba(56, 189, 248, 0.15);
  color: #38bdf8;
}

html.dark .stat-card--users .stat-icon {
  background: rgba(74, 222, 128, 0.15);
  color: #4ade80;
}

html.dark .stat-card--comments .stat-icon {
  background: rgba(251, 191, 36, 0.15);
  color: #fbbf24;
}

html.dark .stat-card--views .stat-icon {
  background: rgba(56, 189, 248, 0.15);
  color: #38bdf8;
}

/* 统计数字样式优化 */
.dashboard-page :deep(.n-statistic .n-statistic-value) {
  font-weight: 800;
  background: linear-gradient(135deg, #0891b2 0%, #059669 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

html.dark .dashboard-page :deep(.n-statistic .n-statistic-value) {
  background: linear-gradient(135deg, #38bdf8 0%, #4ade80 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* 快捷操作图标网格 */
.quick-actions {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(110px, 1fr));
  gap: 12px;
}

.quick-action-item {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 16px 8px;
  border-radius: 12px;
  cursor: pointer;
  transition: background 0.2s ease, transform 0.2s ease;
}

.quick-action-item:hover {
  background: rgba(8, 145, 178, 0.06);
  transform: translateY(-2px);
}

html.dark .quick-action-item:hover {
  background: rgba(56, 189, 248, 0.08);
}

.quick-action-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.quick-action-label {
  font-size: 13px;
  font-weight: 500;
  color: #475569;
}

html.dark .quick-action-label {
  color: #cbd5e1;
}

.quick-action-item .n-badge {
  position: absolute;
  top: 8px;
  right: 12px;
}

/* 减少动效偏好时关闭上浮 */
@media (prefers-reduced-motion: reduce) {
  .dashboard-page :deep(.n-card.stat-card:hover),
  .quick-action-item:hover {
    transform: none;
  }
}

</style>
