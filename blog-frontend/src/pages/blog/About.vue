<template>
  <div class="about-page">
    <div class="about-layout">
      <!-- 关于我内容 -->
      <div class="content-section">
        <n-spin :show="loading">
          <n-space vertical :size="10">
            <!-- 个人介绍卡片 -->
            <n-card class="intro-card" :bordered="false">
              <div class="intro-content">
                <!-- 头像和基本信息区域 -->
                <div class="profile-header">
                  <div class="avatar-section">
                    <n-avatar
                      :src="authorProfile?.author.avatar || ''"
                      :size="120"
                      round
                      :fallback-src="defaultAvatar"
                    >
                      <template v-if="!authorProfile?.author.avatar">
                        {{ (authorProfile?.author.nickname || authorProfile?.author.username || '博主').charAt(0).toUpperCase() }}
                      </template>
                    </n-avatar>
                  </div>
                  <div class="profile-info">
                    <h1 class="author-name">{{ authorProfile?.author.nickname || authorProfile?.author.username || '博主' }}</h1>
                    <p v-if="authorProfile?.author.bio" class="author-bio">{{ authorProfile.author.bio }}</p>
                  </div>
                </div>
                
                <!-- 关于我内容区域 -->
                <div v-if="personalIntroMarkdown" class="philosophy-section">
                  <h3 class="section-title">关于我</h3>
                  <div class="intro-detail-content">
                    <MarkdownPreview :content="personalIntroMarkdown" />
                  </div>
                </div>
              </div>
            </n-card>

            <!-- 文章统计图 -->
            <n-card class="stats-chart-card" title="文章统计图" :bordered="false">
              <div class="charts-container">
                <!-- 文章发布统计图（折线图） -->
                <div class="chart-item">
                  <div class="chart-title">文章发布统计图</div>
                  <div ref="postPublishChartRef" class="chart-wrapper"></div>
                </div>

                <!-- 文章分类统计图（饼图） -->
                <div class="chart-item">
                  <div class="chart-title">文章分类统计图</div>
                  <div ref="categoryChartRef" class="chart-wrapper"></div>
                </div>

                <!-- TOP10 标签统计图（柱状图） -->
                <div class="chart-item">
                  <div class="chart-title">TOP10 标签统计图</div>
                  <div ref="tagChartRef" class="chart-wrapper"></div>
                </div>
              </div>
            </n-card>

            <!-- 相册 -->
            <n-card v-if="albums.length > 0" class="album-card" title="相册" :bordered="false">
              <div class="album-grid">
                <div v-for="album in albums" :key="album.id" class="album-item" @click="handleImageClick(album)">
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
            </n-card>

            <!-- 图片预览 -->
            <n-image-preview
              v-if="showImagePreview"
              v-model:show="showImagePreview"
              :src="previewImageUrl"
            />
          </n-space>
        </n-spin>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue'
import * as echarts from 'echarts'
import type { ECharts } from 'echarts'
import { getAuthorProfile, type AuthorProfile, getPublicCategoryStats, getPublicTagStats, type CategoryStat, type TagStat } from '@/api/blog'
import { getArchives } from '@/api/post'
import { getPublicAboutInfo } from '@/api/setting'
import { getPublicAlbums, type Album } from '@/api/album'
import { useAppStore } from '@/stores'
import MarkdownPreview from '@/components/MarkdownPreview.vue'
import { useMessage } from 'naive-ui'


const appStore = useAppStore()
const message = useMessage()
const authorProfile = ref<AuthorProfile | null>(null)
const loading = ref(false)
const defaultAvatar = '/default-avatar.png'

// 图表相关
const postPublishChartRef = ref<HTMLElement>()
const categoryChartRef = ref<HTMLElement>()
const tagChartRef = ref<HTMLElement>()
let postPublishChart: ECharts | null = null
let categoryChart: ECharts | null = null
let tagChart: ECharts | null = null
// 记录上一次是否为移动端布局，用于判断是否跨越断点
let lastIsMobile = window.innerWidth <= 1024

// 统计数据
const archiveStats = ref<Array<{ month: string; count: number }>>([])
const categoryStats = ref<CategoryStat[]>([])
const tagStats = ref<TagStat[]>([])

// 个人介绍详情（从API获取，Markdown格式）
const personalIntroMarkdown = ref<string>('')

// 相册列表
const albums = ref<Album[]>([])
const previewImageUrl = ref<string>('')
const showImagePreview = ref(false)


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

// 获取分类统计数据
async function fetchCategoryStats() {
  try {
    const res = await getPublicCategoryStats()
    if (res.data) {
      categoryStats.value = res.data
      nextTick(() => {
        // 延迟初始化，确保DOM完全渲染
        setTimeout(() => {
          initCategoryChart()
        }, 100)
      })
    }
  } catch (error) {
    console.error('获取分类统计失败:', error)
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

// 初始化文章发布统计图（折线图）
function initPostPublishChart() {
  if (!postPublishChartRef.value) return

  // 检测是否为移动端
  const isMobile = window.innerWidth <= 1024
  const isSmallMobile = window.innerWidth <= 767
  
  if (!postPublishChart) {
    postPublishChart = echarts.init(postPublishChartRef.value)
  } else {
    // 如果图表已存在，先resize确保尺寸正确
    postPublishChart.resize()
  }

  // 处理数据：格式化月份，补全缺失的月份
  const dataMap = new Map<string, number>()
  archiveStats.value.forEach(item => {
    const month = item.month.substring(0, 7) // YYYY-MM
    dataMap.set(month, Number(item.count))
  })

  // 获取最早和最晚的月份
  const months: string[] = []
  if (archiveStats.value.length > 0) {
    const sortedMonths = Array.from(dataMap.keys()).sort()
    const startMonth = sortedMonths[0]
    const endMonth = sortedMonths[sortedMonths.length - 1]
    
    // 生成月份数组
    const [startYear, startMonthNum] = startMonth.split('-').map(Number)
    const [endYear, endMonthNum] = endMonth.split('-').map(Number)
    
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
  
  // 根据屏幕尺寸动态调整grid配置
  const gridLeft = isSmallMobile ? 45 : isMobile ? 50 : 60
  const gridRight = isSmallMobile ? 20 : isMobile ? 30 : 90
  const gridBottom = isSmallMobile ? 50 : isMobile ? 55 : 60

  const option = {
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        const param = params[0]
        return `${param.name}<br/>${param.seriesName}: ${param.value}`
      }
    },
    grid: {
      left: gridLeft,
      right: gridRight,
      top: 50,
      bottom: gridBottom,
      containLabel: false
    },
    xAxis: {
      type: 'category',
      data: months.map(m => m.replace('-', '-')),
      axisLine: {
        lineStyle: {
          color: isDark ? '#64748b' : '#cbd5e1'
        }
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        color: isDark ? '#e5e7eb' : '#64748b',
        rotate: isSmallMobile ? 45 : 0,
        fontSize: isSmallMobile ? 10 : 12,
        interval: isSmallMobile ? 'auto' : 0,
        overflow: 'break',
        width: isSmallMobile ? 50 : 60,
        margin: isSmallMobile ? 8 : 0
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
        markLine: {
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
    ]
  }

  postPublishChart.setOption(option)
}

// 初始化分类统计图（饼图）
function initCategoryChart() {
  if (!categoryChartRef.value) return

  // 检测是否为移动端
  const isMobile = window.innerWidth <= 1024
  const isSmallMobile = window.innerWidth <= 767

  if (!categoryChart) {
    categoryChart = echarts.init(categoryChartRef.value)
  } else {
    // 如果图表已存在，先resize确保尺寸正确
    categoryChart.resize()
  }

  const defaultColors = ['#5470c6', '#91cc75', '#fac858', '#ee6666', '#73c0de', '#3ba272', '#fc8452', '#9a60b4', '#ea7ccc', '#ff9f7f']
  
  const chartData = categoryStats.value.map((item, index) => ({
    value: item.value,
    name: item.name,
    itemStyle: {
      color: item.color || defaultColors[index % defaultColors.length]
    }
  }))

  const isDark = appStore.theme === 'dark'
  
  // 根据屏幕尺寸动态调整饼图配置
  const radius = isSmallMobile ? ['30%', '60%'] : isMobile ? ['32%', '62%'] : ['35%', '65%']
  const center = isSmallMobile ? ['50%', '50%'] : isMobile ? ['48%', '50%'] : ['45%', '50%']

  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    // 隐藏右侧图例，直接在扇区外侧显示标签
    legend: {
      show: false
    },
    series: [
      {
        name: '文章分类',
        type: 'pie',
        radius: radius,
        center: center,
        avoidLabelOverlap: true,
        itemStyle: {
          borderRadius: 8,
          borderColor: isDark ? '#1e293b' : '#fff',
          borderWidth: 2
        },
        label: {
          show: true,
          formatter: (params: any) => {
            // 根据百分比决定显示方式
            if (params.percent < 5) {
              // 小百分比只显示名称，不截断
              return params.name
            } else {
              // 大百分比显示名称和百分比，不截断名称
              return `${params.name}: ${params.percent.toFixed(1)}%`
            }
          },
          color: isDark ? '#e5e7eb' : '#475569',
          fontSize: isSmallMobile ? 10 : 12,
          fontWeight: 500,
          overflow: 'break',
          width: isSmallMobile ? 100 : isMobile ? 120 : 150,
          lineHeight: isSmallMobile ? 14 : 16,
          rich: {
            name: {
              fontSize: isSmallMobile ? 10 : 12,
              fontWeight: 500
            },
            percent: {
              fontSize: isSmallMobile ? 9 : 11,
              fontWeight: 600
            }
          }
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 14,
            fontWeight: 'bold',
            formatter: '{b}: {d}%'
          }
        },
        labelLine: {
          show: true,
          length: isSmallMobile ? 15 : isMobile ? 20 : 30,
          length2: isSmallMobile ? 10 : isMobile ? 15 : 20,
          lineStyle: {
            color: isDark ? '#94a3b8' : '#94a3b8',
            width: 1
          },
          smooth: 0.2
        },
        data: chartData
      }
    ]
  }

  categoryChart.setOption(option)
}

// 初始化标签统计图（柱状图）
function initTagChart() {
  if (!tagChartRef.value) return

  // 检测是否为移动端
  const isMobile = window.innerWidth <= 1024
  const isSmallMobile = window.innerWidth <= 767

  if (!tagChart) {
    tagChart = echarts.init(tagChartRef.value)
  } else {
    // 如果图表已存在，先resize确保尺寸正确
    tagChart.resize()
  }

  // 按数量降序排序
  const sortedTags = [...tagStats.value].sort((a, b) => b.value - a.value)
  const tagNames = sortedTags.map(t => t.name)
  const tagCounts = sortedTags.map(t => t.value)
  
  // 计算平均值
  const total = tagCounts.reduce((sum, count) => sum + count, 0)
  const average = tagCounts.length > 0 ? (total / tagCounts.length).toFixed(1) : 0

  const isDark = appStore.theme === 'dark'
  
  // 根据屏幕尺寸动态调整grid配置
  const gridLeft = isSmallMobile ? 50 : isMobile ? 60 : 70
  const gridRight = isSmallMobile ? 15 : isMobile ? 25 : 90
  const gridBottom = isSmallMobile ? 70 : isMobile ? 90 : 80

  const option = {
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        const param = params[0]
        return `${param.name}<br/>${param.seriesName}: ${param.value}`
      }
    },
    grid: {
      left: gridLeft,
      right: gridRight,
      top: 50,
      bottom: gridBottom,
      containLabel: false
    },
    xAxis: {
      type: 'category',
      data: tagNames,
      axisLine: {
        lineStyle: {
          color: isDark ? '#64748b' : '#cbd5e1'
        }
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        color: isDark ? '#e5e7eb' : '#64748b',
        rotate: isSmallMobile ? 60 : 45,
        fontSize: isSmallMobile ? 9 : 11,
        interval: 0,
        overflow: 'break',
        width: isSmallMobile ? 40 : 50,
        margin: isSmallMobile ? 10 : 12
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
        markLine: {
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
                fontSize: isSmallMobile ? 9 : 11,
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
    ]
  }

  tagChart.setOption(option)
}

// 监听主题变化，重新渲染图表
watch(() => appStore.theme, () => {
  nextTick(() => {
    if (archiveStats.value.length > 0) initPostPublishChart()
    if (categoryStats.value.length > 0) initCategoryChart()
    if (tagStats.value.length > 0) initTagChart()
  })
})

// 防抖函数
function debounce(func: Function, wait: number) {
  let timeout: NodeJS.Timeout | null = null
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
        if (categoryChart) {
          categoryChart.dispose()
          categoryChart = null
        }
        if (tagChart) {
          tagChart.dispose()
          tagChart = null
        }

        // 布局稳定后，根据当前数据重新初始化
        if (archiveStats.value.length > 0) {
          initPostPublishChart()
        }
        if (categoryStats.value.length > 0) {
          initCategoryChart()
        }
        if (tagStats.value.length > 0) {
          initTagChart()
        }
      } else {
        // 同一断点内只做自适应
        if (postPublishChart) {
          postPublishChart.resize()
        }
        if (categoryChart) {
          categoryChart.resize()
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
  fetchAboutInfo()
  fetchAlbums()
  fetchArchiveStats()
  fetchCategoryStats()
  fetchTagStats()
  
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
  categoryChart?.dispose()
  tagChart?.dispose()
})
</script>

<style scoped>
.about-page {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 20px;
  position: relative;
  z-index: 1;
}

/* 桌面端显示，移动端隐藏 */
.desktop-only {
  display: block;
}

.about-layout {
  display: block;
}

.content-section {
  min-width: 0;
}

/* 个人介绍卡片 */
.intro-card {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.5);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  border-radius: 16px;
  margin-bottom: 24px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.intro-card:hover {
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.12);
  border-color: rgba(8, 145, 178, 0.3);
}

html.dark .intro-card {
  background: rgba(30, 41, 59, 0.7);
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

html.dark .intro-card:hover {
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.5);
  border-color: rgba(56, 189, 248, 0.3);
}

.intro-content {
  display: flex;
  flex-direction: column;
  gap: 0;
}

/* 头像和基本信息区域 */
.profile-header {
  display: flex;
  gap: 32px;
  align-items: flex-start;
  margin-bottom: 32px;
  padding-bottom: 32px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
}

html.dark .profile-header {
  border-bottom-color: rgba(255, 255, 255, 0.1);
}

.avatar-section {
  flex-shrink: 0;
}

.avatar-section :deep(.n-avatar) {
  box-shadow: 0 8px 24px rgba(8, 145, 178, 0.15);
  transition: all 0.3s;
  border: 2px solid rgba(8, 145, 178, 0.1);
}

.avatar-section :deep(.n-avatar:hover) {
  transform: translateY(-3px) scale(1.05);
  box-shadow: 0 14px 32px rgba(8, 145, 178, 0.25);
  border-color: rgba(8, 145, 178, 0.3);
}

.profile-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.author-name {
  margin: 0 0 12px 0;
  font-size: 32px;
  font-weight: 700;
  color: #1a202c;
  background: linear-gradient(135deg, #0891b2 0%, #059669 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  line-height: 1.3;
}

html.dark .author-name {
  background: linear-gradient(135deg, #38bdf8 0%, #4ade80 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.author-bio {
  margin: 0;
  font-size: 16px;
  color: #64748b;
  line-height: 1.8;
  font-style: italic;
}

html.dark .author-bio {
  color: #94a3b8;
}

/* 关于我内容区域 */
.philosophy-section {
  margin-top: 0;
  padding-top: 0;
  border-top: none;
}

.section-title {
  margin: 0 0 20px 0;
  font-size: 22px;
  font-weight: 600;
  color: #1a202c;
  padding-bottom: 0;
  border-bottom: none;
}

html.dark .section-title {
  color: #e5e5e5;
  border-bottom: none;
}

.intro-detail-content {
  font-size: 15px;
  line-height: 1.9;
  color: #475569;
}

html.dark .intro-detail-content {
  color: #cbd5e1;
}

.intro-detail-content p {
  margin: 0 0 16px 0;
  text-align: justify;
}

.intro-detail-content p:last-child {
  margin-bottom: 0;
}

.intro-detail-content p strong {
  color: #0891b2;
  font-weight: 600;
}

html.dark .intro-detail-content p strong {
  color: #38bdf8;
}

/* Markdown样式 */
.intro-detail-content.markdown-body {
  font-size: 15px;
  line-height: 1.9;
}

.intro-detail-content.markdown-body :deep(h1),
.intro-detail-content.markdown-body :deep(h2),
.intro-detail-content.markdown-body :deep(h3),
.intro-detail-content.markdown-body :deep(h4),
.intro-detail-content.markdown-body :deep(h5),
.intro-detail-content.markdown-body :deep(h6) {
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  line-height: 1.25;
}

.intro-detail-content.markdown-body :deep(h1) {
  font-size: 2em;
}

.intro-detail-content.markdown-body :deep(h2) {
  font-size: 1.5em;
}

.intro-detail-content.markdown-body :deep(h3) {
  font-size: 1.25em;
}

.intro-detail-content.markdown-body :deep(code) {
  padding: 2px 6px;
  background: rgba(8, 145, 178, 0.1);
  border-radius: 4px;
  font-size: 0.9em;
}

html.dark .intro-detail-content.markdown-body :deep(code) {
  background: rgba(56, 189, 248, 0.2);
}

.intro-detail-content.markdown-body :deep(pre) {
  padding: 16px;
  background: rgba(0, 0, 0, 0.05);
  border-radius: 8px;
  overflow-x: auto;
}

html.dark .intro-detail-content.markdown-body :deep(pre) {
  background: rgba(255, 255, 255, 0.05);
}

.intro-detail-content.markdown-body :deep(blockquote) {
  padding: 0 16px;
  border-left: 4px solid #0891b2;
  color: #64748b;
}

html.dark .intro-detail-content.markdown-body :deep(blockquote) {
  border-left-color: #38bdf8;
  color: #94a3b8;
}

.intro-detail-content.markdown-body :deep(ul),
.intro-detail-content.markdown-body :deep(ol) {
  padding-left: 24px;
  margin-bottom: 16px;
}

.intro-detail-content.markdown-body :deep(li) {
  margin-bottom: 8px;
}

.intro-detail-content.markdown-body :deep(a) {
  color: #0891b2;
  text-decoration: none;
}

.intro-detail-content.markdown-body :deep(a):hover {
  text-decoration: underline;
}

html.dark .intro-detail-content.markdown-body :deep(a) {
  color: #38bdf8;
}

/* 相册卡片 */
.album-card {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.5);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  border-radius: 16px;
  margin-bottom: 24px;
}

html.dark .album-card {
  background: rgba(30, 41, 59, 0.7);
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

.album-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
  margin-top: 16px;
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
  backdrop-filter: blur(10px);
}

html.dark .album-title {
  color: #e5e5e5;
  background: rgba(30, 41, 59, 0.9);
}

/* 文章统计图卡片 */
.stats-chart-card {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.5);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  border-radius: 16px;
  margin-bottom: 24px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: visible;
}

.stats-chart-card:hover {
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.12);
  border-color: rgba(8, 145, 178, 0.2);
}

html.dark .stats-chart-card {
  background: rgba(30, 41, 59, 0.7);
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

html.dark .stats-chart-card:hover {
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.5);
  border-color: rgba(56, 189, 248, 0.2);
}

.charts-container {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 32px;
  padding: 8px 16px;
  overflow: hidden;
  width: 100%;
  box-sizing: border-box;
}

.chart-item {
  display: flex;
  flex-direction: column;
  gap: 16px;
  background: rgba(8, 145, 178, 0.02);
  border-radius: 12px;
  padding: 20px 16px;
  transition: all 0.3s ease;
  overflow: hidden;
  box-sizing: border-box;
  width: 100%;
}

.chart-item:hover {
  background: rgba(8, 145, 178, 0.05);
  transform: translateY(-2px);
}

html.dark .chart-item {
  background: rgba(56, 189, 248, 0.05);
}

html.dark .chart-item:hover {
  background: rgba(56, 189, 248, 0.1);
}

.chart-title {
  font-size: 17px;
  font-weight: 600;
  color: #1a202c;
  text-align: center;
  margin-bottom: 4px;
  letter-spacing: 0.5px;
  padding-bottom: 12px;
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

/* 移动端响应式 */
@media (max-width: 1024px) {
  .charts-container {
    grid-template-columns: 1fr;
    gap: 28px;
    padding: 8px 12px;
  }

  .chart-item {
    padding: 16px 12px;
    margin: 0;
    width: 100%;
    max-width: 100%;
  }

  .chart-wrapper {
    height: 380px;
    width: 100%;
    max-width: 100%;
  }
}

@media (max-width: 767px) {
  .charts-container {
    gap: 24px;
    padding: 8px 8px;
  }

  .chart-item {
    padding: 12px 8px;
    margin: 0;
    width: 100%;
    max-width: 100%;
  }

  .chart-wrapper {
    height: 320px;
    width: 100%;
    max-width: 100%;
  }

  .chart-title {
    font-size: 15px;
  }
}




/* 移动端布局（< 768px） */
@media (max-width: 767px) {
  .profile-header {
    flex-direction: column;
    align-items: center;
    text-align: center;
    gap: 20px;
    margin-bottom: 24px;
    padding-bottom: 24px;
  }

  .profile-info {
    align-items: center;
  }

  .author-name {
    font-size: 24px;
  }

  .section-title {
    font-size: 20px;
    text-align: left;
  }

  .intro-detail-content {
    text-align: left;
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
}

/* 玻璃态卡片效果 */
.about-page :deep(.n-card) {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.5);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  border-radius: 16px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.about-page :deep(.n-card):hover {
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.12);
  border-color: rgba(8, 145, 178, 0.3);
}

.about-page :deep(.n-card .n-card__content) {
  padding: 20px !important;
}

/* 深色模式卡片 */
html.dark .about-page :deep(.n-card) {
  background: rgba(30, 41, 59, 0.7);
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

html.dark .about-page :deep(.n-card):hover {
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.5);
  border-color: rgba(56, 189, 248, 0.3);
}
</style>
