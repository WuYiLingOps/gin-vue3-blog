<template>
  <div class="dashboard-page">
    <h1 style="margin-bottom: 24px">仪表盘</h1>

    <!-- 统计卡片 -->
    <n-grid :cols="statsColumns" :x-gap="16" :y-gap="16" responsive="screen">
      <n-gi>
        <n-card>
          <n-statistic label="文章总数" :value="stats.posts">
            <template #prefix>
              <n-icon :component="DocumentTextOutline" />
            </template>
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card>
          <n-statistic label="用户总数" :value="stats.users">
            <template #prefix>
              <n-icon :component="PeopleOutline" />
            </template>
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card>
          <n-statistic label="评论总数" :value="stats.comments">
            <template #prefix>
              <n-icon :component="ChatbubblesOutline" />
            </template>
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card>
          <n-statistic label="总浏览量" :value="stats.views">
            <template #prefix>
              <n-icon :component="EyeOutline" />
            </template>
          </n-statistic>
        </n-card>
      </n-gi>
    </n-grid>

    <!-- 图表区域 -->
    <n-grid :cols="isMobile ? 1 : 2" :x-gap="16" :y-gap="16" style="margin-top: 24px">
      <n-gi>
        <n-card title="文章分类统计">
          <div ref="categoryChartRef" :style="chartStyle"></div>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card title="最近7天访问量">
          <div ref="visitChartRef" :style="chartStyle"></div>
        </n-card>
      </n-gi>
    </n-grid>

    <!-- 快捷操作 -->
    <n-card title="快捷操作" style="margin-top: 24px">
      <n-space :wrap="true">
        <n-button type="primary" @click="router.push('/admin/posts')">
          <template #icon>
            <n-icon :component="AddOutline" />
          </template>
          写文章
        </n-button>
        <n-button @click="router.push('/admin/categories')">
          <template #icon>
            <n-icon :component="FolderOutline" />
          </template>
          管理分类
        </n-button>
        <n-button @click="router.push('/admin/tags')">
          <template #icon>
            <n-icon :component="PricetagsOutline" />
          </template>
          管理标签
        </n-button>
      </n-space>
    </n-card>

    <!-- 系统信息 -->
    <n-card title="系统信息" style="margin-top: 24px">
      <n-descriptions :column="isMobile ? 1 : 2">
        <n-descriptions-item label="系统版本">v2.0.0</n-descriptions-item>
        <n-descriptions-item label="Vue版本">3.3.4</n-descriptions-item>
        <n-descriptions-item label="数据库">PostgreSQL</n-descriptions-item>
        <n-descriptions-item label="后端框架">Go + Gin</n-descriptions-item>
      </n-descriptions>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import * as echarts from 'echarts'
import type { ECharts } from 'echarts'
import {
  DocumentTextOutline,
  PeopleOutline,
  ChatbubblesOutline,
  EyeOutline,
  AddOutline,
  FolderOutline,
  PricetagsOutline
} from '@vicons/ionicons5'
import { useAppStore } from '@/stores'
import { getDashboardStats, getCategoryStats, getVisitStats } from '@/api'
import type { DashboardStats, CategoryStat, VisitStat } from '@/api'

const router = useRouter()
const appStore = useAppStore()
const message = useMessage()

const stats = ref<DashboardStats>({
  posts: 0,
  users: 0,
  comments: 0,
  views: 0
})

const categoryStats = ref<CategoryStat[]>([])
const visitStats = ref<VisitStat[]>([])

const categoryChartRef = ref<HTMLElement>()
const visitChartRef = ref<HTMLElement>()
let categoryChart: ECharts | null = null
let visitChart: ECharts | null = null

const loading = ref(false)
const isMobile = ref(false)

// 响应式列数
const statsColumns = computed(() => {
  if (window.innerWidth <= 480) return 1
  if (window.innerWidth <= 768) return 2
  return 4
})

const chartStyle = computed(() => ({
  width: '100%',
  height: isMobile.value ? '250px' : '300px'
}))

// 检测移动设备
function checkMobile() {
  isMobile.value = window.innerWidth <= 768
}

// 获取统计数据
async function fetchStats() {
  try {
    loading.value = true
    const res = await getDashboardStats()
    if (res.data) {
      stats.value = res.data
    }
  } catch (error) {
    message.error('获取统计数据失败')
    console.error(error)
  } finally {
    loading.value = false
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
  
  // 获取统计数据
  await fetchStats()
  await fetchCategoryStats()
  await fetchVisitStats()

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
  categoryChart?.resize()
  visitChart?.resize()
}

// 初始化分类统计饼图
function initCategoryChart() {
  if (!categoryChartRef.value) return

  if (!categoryChart) {
    categoryChart = echarts.init(categoryChartRef.value)
  }

  // 默认颜色方案
  const defaultColors = ['#5470c6', '#91cc75', '#fac858', '#ee6666', '#73c0de', '#3ba272', '#fc8452', '#9a60b4']

  // 转换数据格式
  const chartData = categoryStats.value.map((item, index) => ({
    value: item.value,
    name: item.name,
    itemStyle: { 
      color: item.color || defaultColors[index % defaultColors.length]
    }
  }))

  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      right: 10,
      top: 'center',
      textStyle: {
        color: appStore.theme === 'dark' ? '#fff' : '#333'
      }
    },
    series: [
      {
        name: '文章分类',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 20,
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: chartData
      }
    ]
  }

  categoryChart.setOption(option)
}

// 初始化最近访问量折线图
function initVisitChart() {
  if (!visitChartRef.value) return

  if (!visitChart) {
    visitChart = echarts.init(visitChartRef.value)
  }

  const dates = visitStats.value.map((item) => item.date.slice(5)) // 显示 MM-DD
  const counts = visitStats.value.map((item) => item.count)

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

.dashboard-page :deep(.n-card:hover) {
  transform: translateY(-4px);
  box-shadow: 0 16px 40px rgba(0, 0, 0, 0.12);
  border-color: rgba(8, 145, 178, 0.3);
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

/* 快捷操作按钮样式 */
.dashboard-page :deep(.n-button) {
  font-weight: 600;
  transition: all 0.3s;
}

.dashboard-page :deep(.n-button:hover) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(8, 145, 178, 0.3);
}

/* 系统信息卡片 */
.dashboard-page :deep(.n-descriptions-item-label) {
  font-weight: 600;
  color: #64748b;
}

html.dark .dashboard-page :deep(.n-descriptions-item-label) {
  color: #94a3b8;
}
</style>

