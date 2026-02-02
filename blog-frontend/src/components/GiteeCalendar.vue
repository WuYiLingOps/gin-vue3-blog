<!--
  项目名称：blog-frontend
  文件名称：GiteeCalendar.vue
  创建时间：2026-02-01 20:03:19

  系统用户：Administrator
  作　　者：無以菱
  联系邮箱：huangjing510@126.com
  功能描述：Gitee贡献度热力图组件，展示用户在Gitee上的代码提交活动热力图，包括过去一年、最近一月、最近一周的提交统计数据。
-->
<template>
  <div class="hexo-calendar-card">
    <n-card size="small" class="calendar-card" :bordered="false">
      <template #header>
        <div class="card-header">
          <span class="title">贡献度热力图</span>
          <a
            v-if="username"
            :href="profileUrl"
            target="_blank"
            rel="noopener"
            class="subtitle"
          >
            @{{ username }}
          </a>
        </div>
      </template>

      <div v-if="loading" class="calendar-loading">加载中...</div>

      <div v-else>
        <div v-if="error" class="calendar-error">{{ error }}</div>

        <div v-if="flatDays.length" class="graph-body-outer" ref="outerRef">
          <div class="graph-body-inner" :style="{ '--scale': scale }" ref="innerRef">
            <div class="weekdays-col">
              <div>日</div>
              <div class="hidden-label">一</div>
              <div>二</div>
              <div class="hidden-label">三</div>
              <div>四</div>
              <div class="hidden-label">五</div>
              <div>六</div>
            </div>

            <div class="graph-content-col">
              <div class="months-row">
                <template v-for="(pos, idx) in monthPositions" :key="idx">
                  <span
                    v-if="pos !== null"
                    class="month-label"
                    :style="{ left: pos + 'px' }"
                  >
                    {{ monthLabels[idx] }}
                  </span>
                </template>
              </div>

              <div class="grid">
                <div
                  v-for="day in flatDays"
                  :key="day.date"
                  class="cell"
                  :class="[cellLevelClass(day.count), { 'cell-empty': day.count < 0 }]"
                  :data-date="day.date"
                  :data-count="day.count >= 0 ? day.count : 0"
                />
              </div>
            </div>
          </div>
        </div>

        <div v-if="flatDays.length" class="meta-info">
          <div class="source">
            数据来源
            <a
              v-if="username"
              :href="profileUrl"
              target="_blank"
              rel="noopener"
            >
              @{{ username }}
            </a>
          </div>
          <div class="legend">
            <span>Less</span>
            <div class="legend-box level-0" />
            <div class="legend-box level-1" />
            <div class="legend-box level-2" />
            <div class="legend-box level-3" />
            <div class="legend-box level-4" />
            <span>More</span>
          </div>
        </div>

        <hr v-if="flatDays.length" class="divider" />

        <div v-if="flatDays.length" class="stats">
          <div class="stat-item">
            <p>过去一年提交</p>
            <h3>{{ total }}</h3>
            <span class="range">{{ dateRange }}</span>
          </div>
          <div class="stat-item">
            <p>最近一月提交</p>
            <h3>{{ lastMonthStats.total }}</h3>
            <span class="range">{{ lastMonthStats.from }} - {{ lastMonthStats.to }}</span>
          </div>
          <div class="stat-item">
            <p>最近一周提交</p>
            <h3>{{ lastWeekStats.total }}</h3>
            <span class="range">{{ lastWeekStats.from }} - {{ lastWeekStats.to }}</span>
          </div>
        </div>

        <div v-if="!error && !flatDays.length" class="calendar-empty">
          暂无贡献数据
        </div>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref, computed, watch, nextTick } from 'vue'
import { useMessage } from 'naive-ui'
import { getPublicSettings } from '@/api/setting'

const message = useMessage()

interface CalendarDay {
  date: string
  count: number
}

type CalendarWeek = CalendarDay[]

const username = ref('')
const profileUrl = ref('')

const loading = ref(false)
const error = ref('')

const total = ref(0)
const weeks = ref<CalendarWeek[]>([])

const flatDays = computed(() => weeks.value.flat())

// 自适应缩放相关
const outerRef = ref<HTMLElement | null>(null)
const innerRef = ref<HTMLElement | null>(null)
const scale = ref(1)
// 是否桌面端，用于控制桌面/移动端不同的视觉对齐策略（响应式，支持尺寸切换）
const isDesktop = ref(true)
let resizeObserver: ResizeObserver | null = null

function updateScale() {
  const outer = outerRef.value
  const inner = innerRef.value
  if (!outer || !inner) return

  if (typeof window !== 'undefined') {
    isDesktop.value = window.innerWidth >= 769
  }

  // 移动端（< 769px）保留原始尺寸，改用横向滚动，并默认对齐到最新数据（最右侧）
  if (!isDesktop.value) {
    scale.value = 1
    // 下一帧将滚动条滚动到最右，优先展示最近的提交数据
    requestAnimationFrame(() => {
      const el = outerRef.value
      if (!el) return
      const maxScrollLeft = el.scrollWidth - el.clientWidth
      if (maxScrollLeft > 0) {
        el.scrollLeft = maxScrollLeft
      }
    })
    return
  }

  const outerWidth = outer.clientWidth
  // 获取未缩放时的原始宽度
  const innerWidth = inner.scrollWidth || inner.clientWidth

  if (!innerWidth || !outerWidth) {
    scale.value = 1
    return
  }

  const ratio = outerWidth / innerWidth
  const next = Math.min(1, ratio)
  
  // 完全自适应缩放，确保所有数据都完整显示在容器内
  // 移除下限限制，允许适当缩小以完全适配容器宽度
  // 这样就不需要横向滚动，所有内容都能一次性看到
  scale.value = next
}

const dateRange = computed(() => {
  if (!flatDays.value.length) return ''
  const first = flatDays.value[0].date
  const last = flatDays.value[flatDays.value.length - 1].date
  return `${first} - ${last}`
})

const lastWeekStats = computed(() => {
  if (!flatDays.value.length) return { total: 0, from: '', to: '' }
  
  // 过滤掉无效日期（count < 0 表示超出范围的空白天数）
  const validDays = flatDays.value.filter((d) => d.count >= 0)
  if (!validDays.length) return { total: 0, from: '', to: '' }
  
  // 获取最后一个有效日期作为结束日期
  const last = validDays[validDays.length - 1]
  const to = last.date
  
  // 计算7天前的日期（包含今天，所以是6天前）
  const toDate = new Date(to + 'T00:00:00')
  const fromDate = new Date(toDate)
  fromDate.setDate(toDate.getDate() - 6)
  const fromStr = fromDate.toISOString().slice(0, 10)
  
  // 筛选出最近7天的有效数据（基于日期范围，而不是数组索引）
  const slice = validDays.filter((d) => d.date >= fromStr && d.date <= to)
  // 只累加有效的提交次数（确保 count >= 0）
  const sum = slice.reduce((s, d) => s + Math.max(0, d.count), 0)
  
  return { 
    total: sum, 
    from: slice.length ? slice[0].date : fromStr, 
    to 
  }
})

const lastMonthStats = computed(() => {
  if (!flatDays.value.length) return { total: 0, from: '', to: '' }
  
  // 过滤掉无效日期（count < 0 表示超出范围的空白天数）
  const validDays = flatDays.value.filter((d) => d.count >= 0)
  if (!validDays.length) return { total: 0, from: '', to: '' }
  
  // 获取最后一个有效日期作为结束日期
  const last = validDays[validDays.length - 1]
  const to = last.date
  
  // 计算30天前的日期（包含今天，所以是29天前）
  const fromDate = new Date(to + 'T00:00:00')
  fromDate.setDate(fromDate.getDate() - 29)
  const fromStr = fromDate.toISOString().slice(0, 10)
  
  // 筛选出最近30天的有效数据（基于日期范围）
  const slice = validDays.filter((d) => d.date >= fromStr && d.date <= to)
  // 只累加有效的提交次数（确保 count >= 0）
  const sum = slice.reduce((s, d) => s + Math.max(0, d.count), 0)
  
  return { 
    total: sum, 
    from: slice.length ? slice[0].date : fromStr, 
    to 
  }
})

// 网格尺寸（需与下方 CSS 中 .grid 的宽高 / gap 保持一致）
const CELL_SIZE = 12
const CELL_GAP = 5
// 月份标签最小间距，避免重叠
const MIN_MONTH_GAP = (CELL_SIZE + CELL_GAP) * 2

// 月份标签（中文）
const monthLabelsArray = ['一月', '二月', '三月', '四月', '五月', '六月', '七月', '八月', '九月', '十月', '十一月', '十二月']

// 构造「需要展示的月份」时间轴：以当前月份为结束，回溯 12 个月
const monthTimeline = computed(() => {
  const months: Array<{ year: number; month: number }> = []
  const end = new Date()
  end.setDate(1)
  const start = new Date(end)
  start.setMonth(start.getMonth() - 11)
  for (let i = 0; i < 12; i++) {
    const cur = new Date(start)
    cur.setMonth(start.getMonth() + i)
    months.push({ year: cur.getFullYear(), month: cur.getMonth() })
  }
  return months
})

// 收集各月份出现位置与有效日期范围
const monthMeta = computed(() => {
  const info = new Map<
    string,
    {
      firstCol: number
      firstDate: Date
    }
  >()
  let earliest: Date | null = null
  let latest: Date | null = null

  weeks.value.forEach((week, colIndex) => {
    week.forEach((d) => {
      if (!d.date || d.count < 0) return
      const date = new Date(d.date + 'T00:00:00')
      if (Number.isNaN(date.getTime())) return

      if (!earliest || date < earliest) earliest = date
      if (!latest || date > latest) latest = date

      const key = `${date.getFullYear()}-${date.getMonth()}`
      if (!info.has(key)) {
        info.set(key, { firstCol: colIndex, firstDate: date })
      }
    })
  })

  const maxPos = Math.max(0, (weeks.value.length - 1) * (CELL_SIZE + CELL_GAP))
  return { info, earliest, latest, maxPos }
})

// 动态月份标签（支持跨年、缺失月份插值、最小间距防重叠）
const monthLabels = computed(() => monthTimeline.value.map((m) => monthLabelsArray[m.month]))

// 每个月份对应的像素位置（相对于网格左侧），支持缺失月份插值与最小间距防重叠
// 位置以「月份首列单元格的中心点」为基准，确保缩放/压缩后仍与网格精确对齐
const monthPositions = computed(() => {
  if (!weeks.value.length) return [] as (number | null)[]

  const { info, maxPos } = monthMeta.value
  const timeline = monthTimeline.value
  if (!timeline.length) return []

  // 基于现有数据的基准位置（使用该月第一列所在列的中心点）
  const basePositions = timeline.map(({ year, month }) => {
    const key = `${year}-${month}`
    const meta = info.get(key)
    if (!meta) return null
    return meta.firstCol * (CELL_SIZE + CELL_GAP) + CELL_SIZE / 2
  })
  const anchors = basePositions.map((p) => p !== null)

  // 兜底参考间距（均匀分布，避免全空导致无法插值），同样以中心点为单位
  const uniformSpacing =
    timeline.length > 1
      ? (maxPos + CELL_SIZE / 2) / Math.max(1, timeline.length - 1)
      : maxPos + CELL_SIZE / 2

  // 缺失月份插值：向左右寻找最近已知位置，线性推断
  const interpolated = [...basePositions]
  for (let i = 0; i < interpolated.length; i++) {
    if (interpolated[i] !== null) continue

    // 找到左侧最近的已知位置
    let leftIdx = i - 1
    while (leftIdx >= 0 && interpolated[leftIdx] === null) leftIdx--
    // 找到右侧最近的已知位置
    let rightIdx = i + 1
    while (rightIdx < interpolated.length && interpolated[rightIdx] === null) rightIdx++

    if (leftIdx >= 0 && rightIdx < interpolated.length && interpolated[leftIdx] !== null && interpolated[rightIdx] !== null) {
      const leftPos = interpolated[leftIdx] as number
      const rightPos = interpolated[rightIdx] as number
      const gap = rightIdx - leftIdx
      const ratio = (i - leftIdx) / gap
      interpolated[i] = leftPos + (rightPos - leftPos) * ratio
    } else if (leftIdx >= 0 && interpolated[leftIdx] !== null) {
      // 只有左侧：沿用均匀步进
      interpolated[i] = (interpolated[leftIdx] as number) + uniformSpacing * (i - leftIdx)
    } else if (rightIdx < interpolated.length && interpolated[rightIdx] !== null) {
      // 只有右侧：反向均匀步进
      interpolated[i] = (interpolated[rightIdx] as number) - uniformSpacing * (rightIdx - i)
    } else {
      // 全部为空：完全均匀分布
      interpolated[i] = uniformSpacing * i
    }
  }

  // 防重叠：只移动非锚点，保持已有月份严格对齐网格列
  const adjusted = [...interpolated] as number[]
  for (let i = 1; i < adjusted.length; i++) {
    if (adjusted[i] === null) continue
    if (anchors[i]) continue
    if (adjusted[i] < adjusted[i - 1] + MIN_MONTH_GAP) {
      adjusted[i] = adjusted[i - 1] + MIN_MONTH_GAP
    }
  }

  // 末端限制：只回拉非锚点，避免整体平移导致锚点偏离网格
  for (let i = adjusted.length - 2; i >= 0; i--) {
    if (adjusted[i] === null) continue
    if (anchors[i]) continue
    if (adjusted[i] > adjusted[i + 1] - MIN_MONTH_GAP) {
      adjusted[i] = adjusted[i + 1] - MIN_MONTH_GAP
    }
  }

  // 边界收缩：非锚点超过右界则回拉到 (maxPos + CELL_SIZE/2)，保证中心点仍在网格内
  for (let i = 0; i < adjusted.length; i++) {
    if (adjusted[i] === null) continue
    if (anchors[i]) continue
    const maxCenter = maxPos + CELL_SIZE / 2
    if (adjusted[i] > maxCenter) adjusted[i] = maxCenter
  }

  return adjusted.map((pos) => {
    if (pos === null) return null
    // 中心点不能小于第一列中心，也不能大于最后一列中心
    const minCenter = CELL_SIZE / 2
    const maxCenter = maxPos + CELL_SIZE / 2
    let center = Math.min(Math.max(pos, minCenter), maxCenter)

    // 仅在「桌面端」时做视觉微调；移动端横向滚动模式不受影响
    if (isDesktop.value) {
      center += CELL_SIZE * 4.5
      center = Math.min(Math.max(center, minCenter), maxCenter)
    }

    // month-label 使用 translateX(-50%)，此处直接返回中心点 x 即可
    return center
  })
})

// 颜色等级
function cellLevelClass(count: number) {
  if (count < 0) return '' // 超出范围的日期不显示
  if (count === 0) return ''
  if (count <= 5) return 'l1'
  if (count <= 10) return 'l2'
  if (count <= 15) return 'l3'
  return 'l4'
}

function extractUsername(url?: string) {
  if (!url) return ''
  const match = url.match(/https?:\/\/(?:www\.)?gitee\.com\/([^\/\s?#]+)/i)
  return match ? match[1] : ''
}

// 使用后端缓存接口（带Redis缓存，20分钟过期）
function getCalendarApiBase() {
  let baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
  
  // 移除末尾的斜杠
  baseURL = baseURL.replace(/\/+$/, '')
  
  // 如果 baseURL 不包含 /api，则添加
  if (!baseURL.endsWith('/api')) {
    baseURL = `${baseURL}/api`
  }
  
  return `${baseURL}/calendar/gitee`
}

// 规范化周数据，确保每一周都从周日开始，不足7天的用空白天数补齐
function normalizeWeeks(weeks: CalendarWeek[]): CalendarWeek[] {
  if (!weeks.length) return []
  
  const normalized: CalendarWeek[] = []
  
  // 处理第一周：确保从周日开始
  const firstWeek = weeks[0]
  if (firstWeek && firstWeek.length > 0) {
    const firstDate = new Date(firstWeek[0].date + 'T00:00:00') // 明确指定时区，避免时区问题
    const firstDayOfWeek = firstDate.getDay() // 0=周日, 1=周一, ..., 6=周六
    
    // 如果第一天不是周日，需要在前面补全空白天数
    if (firstDayOfWeek !== 0) {
      const paddedWeek: CalendarDay[] = []
      // 补全前面的空白天数（从周日开始）
      for (let i = 0; i < firstDayOfWeek; i++) {
        const emptyDate = new Date(firstDate)
        emptyDate.setDate(firstDate.getDate() - (firstDayOfWeek - i))
        paddedWeek.push({ date: emptyDate.toISOString().slice(0, 10), count: -1 }) // 用-1标记超出范围的日期
      }
      // 添加原有的天数
      paddedWeek.push(...firstWeek)
      normalized.push(paddedWeek)
    } else {
      normalized.push([...firstWeek])
    }
  }
  
  // 处理中间周：确保每周都有7天
  for (let i = 1; i < weeks.length - 1; i++) {
    const week = weeks[i]
    if (week && week.length === 7) {
      normalized.push([...week])
    } else if (week && week.length > 0) {
      // 如果周数据不足7天，补齐到7天
      const paddedWeek: CalendarDay[] = [...week]
      const lastDate = new Date(week[week.length - 1].date + 'T00:00:00')
      while (paddedWeek.length < 7) {
        const nextDate = new Date(lastDate)
        nextDate.setDate(lastDate.getDate() + (paddedWeek.length - week.length + 1))
        paddedWeek.push({ date: nextDate.toISOString().slice(0, 10), count: -1 }) // 用-1标记超出范围的日期
      }
      normalized.push(paddedWeek)
    }
  }
  
  // 处理最后一周：确保补齐到7天
  if (weeks.length > 1) {
    const lastWeek = weeks[weeks.length - 1]
    if (lastWeek && lastWeek.length > 0) {
      const paddedWeek: CalendarDay[] = [...lastWeek]
      const lastDate = new Date(lastWeek[lastWeek.length - 1].date + 'T00:00:00')
      const today = new Date()
      today.setHours(0, 0, 0, 0)
      
      while (paddedWeek.length < 7) {
        const nextDate = new Date(lastDate)
        nextDate.setDate(lastDate.getDate() + (paddedWeek.length - lastWeek.length + 1))
        // 如果日期超过今天，标记为超出范围
        const isFuture = nextDate > today
        paddedWeek.push({ 
          date: nextDate.toISOString().slice(0, 10), 
          count: isFuture ? -1 : 0 
        })
      }
      normalized.push(paddedWeek)
    }
  } else if (weeks.length === 1) {
    // 如果只有一周，也需要补齐到7天
    const lastWeek = weeks[0]
    if (lastWeek && lastWeek.length > 0 && lastWeek.length < 7) {
      const paddedWeek: CalendarDay[] = [...lastWeek]
      const lastDate = new Date(lastWeek[lastWeek.length - 1].date + 'T00:00:00')
      const today = new Date()
      today.setHours(0, 0, 0, 0)
      
      while (paddedWeek.length < 7) {
        const nextDate = new Date(lastDate)
        nextDate.setDate(lastDate.getDate() + (paddedWeek.length - lastWeek.length + 1))
        const isFuture = nextDate > today
        paddedWeek.push({ 
          date: nextDate.toISOString().slice(0, 10), 
          count: isFuture ? -1 : 0 
        })
      }
      normalized[0] = paddedWeek
    }
  }
  
  return normalized
}

// 构造一整年的「空」数据，用于无数据 / 请求失败时仍然展示完整网格样式
function buildEmptyWeeks(): CalendarWeek[] {
  const today = new Date()
  today.setHours(0, 0, 0, 0) // 设置为当天0点

  // 计算一年前的日期（365天前）
  const oneYearAgo = new Date(today)
  oneYearAgo.setDate(today.getDate() - 364) // 包含今天，所以是364天前
  
  // 找到一年前那个日期所在周的周日（作为起始日期）
  const startDate = new Date(oneYearAgo)
  const dayOfWeek = startDate.getDay() // 0=周日, 1=周一, ..., 6=周六
  startDate.setDate(startDate.getDate() - dayOfWeek) // 回退到周日

  // 从起始周日开始，生成数据直到今天所在周的周六（补齐到完整的周）
  const endDate = new Date(today)
  const endDayOfWeek = endDate.getDay() // 今天是一周的第几天
  // 计算到今天所在周的周六需要多少天
  const daysToEndOfWeek = 6 - endDayOfWeek
  const finalDate = new Date(today)
  finalDate.setDate(today.getDate() + daysToEndOfWeek)
  
  // 计算需要多少周（从起始周日到最终周六）
  const daysDiff = Math.ceil((finalDate.getTime() - startDate.getTime()) / (1000 * 60 * 60 * 24)) + 1
  const weeksNeeded = Math.ceil(daysDiff / 7)
  
  const weeksArr: CalendarWeek[] = []
  for (let weekIdx = 0; weekIdx < weeksNeeded; weekIdx++) {
    const week: CalendarDay[] = []
    for (let dayIdx = 0; dayIdx < 7; dayIdx++) {
      const date = new Date(startDate)
      date.setDate(startDate.getDate() + weekIdx * 7 + dayIdx)
      const dateStr = date.toISOString().slice(0, 10)
      // 只包含一年前到今天之间的日期，其他日期标记为空
      const isInRange = date >= oneYearAgo && date <= today
      week.push({ date: dateStr, count: isInRange ? 0 : -1 }) // 用-1标记超出范围的日期
    }
    weeksArr.push(week)
  }
  
  return weeksArr
}

async function loadUsername() {
  const res = await getPublicSettings()
  const data = res.data || {}
  const giteeUrl = data.social_gitee
  const giteeName = extractUsername(giteeUrl)

  username.value = giteeName

  if (giteeUrl) {
    profileUrl.value = giteeUrl
  } else if (giteeName) {
    profileUrl.value = `https://gitee.com/${giteeName}`
  } else {
    profileUrl.value = ''
  }
}

async function fetchData() {
  loading.value = true
  error.value = ''
  try {
    const targetUser = username.value

    // 没有用户名时，不请求接口，直接渲染空网格
    if (!targetUser) {
      total.value = 0
      weeks.value = buildEmptyWeeks()
      error.value = ''
      return
    }

    const url = `${getCalendarApiBase()}?user=${encodeURIComponent(targetUser)}`
    const res = await fetch(url)
    if (!res.ok) {
      throw new Error(`接口返回错误状态：${res.status}`)
    }
    const response = (await res.json()) as { code: number; message: string; data: { total: number; contributions: CalendarWeek[] } }
    
    // 检查后端返回格式
    if (response.code !== 200 || !response.data) {
      throw new Error(response.message || '接口返回数据格式不正确')
    }
    
    const data = response.data
    if (!data || !Array.isArray(data.contributions)) {
      throw new Error('接口返回数据格式不正确')
    }
    if (!data.contributions.length) {
      total.value = 0
      weeks.value = buildEmptyWeeks()
    } else {
      total.value = data.total || 0
      // 先映射数据，然后规范化周数据以确保周几对齐
      const mappedWeeks = data.contributions.map((week) =>
        week.map((d) => ({
          date: d.date,
          count: d.count
        }))
      )
      weeks.value = normalizeWeeks(mappedWeeks)
    }
  } catch (e: any) {
    console.error(e)
    error.value = e.message || '加载失败，请稍后重试'
    total.value = 0
    // 接口异常时也回退为空网格样式
    weeks.value = buildEmptyWeeks()
    message.error(error.value)
  } finally {
    loading.value = false
    // 数据加载或回退完成后，等待 DOM 更新再计算缩放
    nextTick(() => {
      updateScale()
    })
  }
}

onMounted(async () => {
  try {
    await loadUsername()
    await fetchData()
  } catch {
    // 错误已在内部处理
  }

  if (typeof ResizeObserver !== 'undefined') {
    resizeObserver = new ResizeObserver(() => {
      updateScale()
    })
    if (outerRef.value) {
      resizeObserver.observe(outerRef.value)
    }
  } else {
    window.addEventListener('resize', updateScale)
  }
})

onBeforeUnmount(() => {
  if (resizeObserver) {
    resizeObserver.disconnect()
    resizeObserver = null
  } else {
    window.removeEventListener('resize', updateScale)
  }
})

watch(
  () => flatDays.value.length,
  () => {
    nextTick(() => {
      updateScale()
    })
  }
)
</script>

<style scoped>
.hexo-calendar-card {
  width: 100%;
  overflow: visible; /* 允许提示框超出容器边界显示 */
}

.calendar-card {
  height: 100%;
  overflow: visible; /* 允许提示框超出卡片边界显示 */
}

/* 确保 n-card 内部的所有容器都允许提示框显示 */
.calendar-card :deep(.n-card__content) {
  overflow: visible !important;
}

.calendar-card :deep(.n-card) {
  overflow: visible !important;
}

.card-header {
  display: flex;
  align-items: baseline;
  gap: 8px;
}

.title {
  font-weight: 600;
}

.subtitle {
  font-size: 12px;
  color: #64748b;
}

/* 夜间模式标题优化 */
html.dark .title {
  color: #e5e7eb;
}

html.dark .subtitle {
  color: #9ca3af;
}

.calendar-loading {
  padding: 24px;
  text-align: center;
  font-size: 14px;
  color: #64748b;
}

/* 夜间模式加载文字 */
html.dark .calendar-loading {
  color: #9ca3af;
}

.calendar-error {
  margin-top: 8px;
  padding: 8px 12px;
  background: #fff3cd;
  border-radius: 4px;
  border: 1px solid #ffc107;
  color: #856404;
  font-size: 12px;
}

.calendar-empty {
  padding: 16px;
  text-align: center;
  font-size: 13px;
  color: #9ca3af;
}

/* 夜间模式空数据文字 */
html.dark .calendar-empty {
  color: #6b7280;
}

/* --- 核心布局，参考 go-code-calendar-api/web/index.html --- */
.graph-body-outer {
  margin-top: 12px;
  /* 允许提示框显示，缩放逻辑会确保内容适配容器 */
  overflow: visible;
  position: relative;
  /* 为提示框留出显示空间 */
  padding-top: 30px;
  margin-top: calc(12px - 30px);
}

.graph-body-inner {
  display: flex;
  align-items: flex-end;
  transform-origin: left top;
  transform: scale(var(--scale, 1));
  overflow: visible; /* 允许提示框显示 */
  position: relative;
}

/* 左侧：星期标签柱 */
.weekdays-col {
  display: flex;
  flex-direction: column;
  height: calc((12px * 7) + (5px * 6));
  justify-content: space-between;
  margin-right: 8px;
  padding-bottom: 2px;
  font-size: 12px;
  color: #6e7781;
  text-align: right;
  line-height: 12px;
}

/* 夜间模式星期标签 */
html.dark .weekdays-col {
  color: #9ca3af;
}

/* 右侧：包含月份和网格的内容区 */
.graph-content-col {
  display: flex;
  flex-direction: column;
  overflow: visible; /* 允许提示框超出边界显示 */
  position: relative; /* 为提示框提供定位上下文 */
}

/* 顶部月份栏 */
.months-row {
  position: relative;
  height: 20px; /* 预留固定高度，避免与方块重叠 */
  margin-bottom: 8px;
  font-size: 12px;
  color: #6e7781;
  box-sizing: border-box;
}

/* 夜间模式月份标签 */
html.dark .months-row {
  color: #9ca3af;
}

.month-label {
  position: absolute;
  top: 0;
  transform: translateX(-50%);
  white-space: nowrap;
}

/* 网格区域 */
.grid {
  display: grid;
  grid-template-rows: repeat(7, 12px);
  grid-auto-flow: column;
  grid-auto-columns: 12px;
  gap: 5px;
}

.cell {
  width: 12px;
  height: 12px;
  border-radius: 2px;
  background-color: #ebedf0;
  position: relative;
  cursor: pointer;
  box-sizing: border-box;
  border: 1px solid transparent;
}

.cell-empty {
  background-color: transparent !important;
  cursor: default;
  pointer-events: none;
}

.cell.l1 {
  background-color: #9be9a8;
}
.cell.l2 {
  background-color: #40c463;
}
.cell.l3 {
  background-color: #30a14e;
}
.cell.l4 {
  background-color: #216e39;
}

.cell:hover {
  /* 仅高亮边框，不改变几何尺寸，避免任何布局抖动 */
  border-color: #0969da;
}

.cell {
  position: relative; /* 为提示框提供定位上下文 */
}

.cell:hover::after {
  content: attr(data-date) ' · ' attr(data-count) ' 次提交';
  position: absolute;
  bottom: calc(100% + 5px);
  left: 50%;
  transform: translateX(-50%);
  padding: 4px 8px;
  font-size: 11px;
  color: #fff;
  background: rgba(0, 0, 0, 0.9);
  border-radius: 4px;
  white-space: nowrap;
  pointer-events: none;
  z-index: 1000; /* 提高层级，确保显示在最上层 */
  /* 确保提示框不会被裁剪 */
  min-width: max-content;
  /* 防止提示框超出容器，智能调整位置 */
  max-width: 200px;
  word-break: keep-all;
}

/* 优化边缘方块的提示框位置 - 使用更精确的选择器 */
.grid {
  position: relative;
  overflow: visible; /* 允许提示框超出网格显示 */
}

/* 检测并调整最右侧列的提示框位置 */
/* 由于网格使用 grid-auto-flow: column，最右侧列是最后7个方块 */
.grid > .cell:nth-last-child(-n+7):hover::after {
  /* 最右侧列（最后7个方块），提示框左对齐：从方块中心向左完全延伸 */
  left: 50%;
  right: auto;
  transform: translateX(-100%); /* 提示框完全在方块左侧，右边缘对齐到方块中心 */
}

/* 检测并调整最左侧列的提示框位置 */
/* 由于网格使用 grid-auto-flow: column，最左侧列是前7个方块 */
.grid > .cell:nth-child(-n+7):hover::after {
  /* 最左侧列（前7个方块），提示框右对齐：从方块中心向右完全延伸 */
  left: 50%;
  right: auto;
  transform: translateX(0); /* 提示框完全在方块右侧，左边缘对齐到方块中心 */
}

/* --- 底部信息栏 --- */
.meta-info {
  margin-top: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: #6e7781;
}

.source {
  color: #57606a;
}

/* 夜间模式数据来源 */
html.dark .source {
  color: #9ca3af;
}

.source a {
  color: #0969da;
  text-decoration: none;
}

.source a:hover {
  text-decoration: underline;
}

.legend {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
}

/* 夜间模式图例文字 */
html.dark .legend {
  color: #9ca3af;
}

.legend-box {
  width: 10px;
  height: 10px;
  border-radius: 2px;
  background-color: #ebedf0;
}

.legend-box.level-1 {
  background-color: #9be9a8;
}
.legend-box.level-2 {
  background-color: #40c463;
}
.legend-box.level-3 {
  background-color: #30a14e;
}
.legend-box.level-4 {
  background-color: #216e39;
}

/* 分割线 */
.divider {
  height: 0;
  margin: 24px 0;
  border: none;
  border-top: 1px dashed #e1e4e8;
}

/* 夜间模式分割线 */
html.dark .divider {
  border-top-color: #374151;
}

/* 底部统计数据 */
.stats {
  display: flex;
  justify-content: space-between;
  padding: 0 40px;
  text-align: center;
}

.stat-item h3 {
  margin: 8px 0;
  font-weight: 400;
  font-size: 26px;
  color: #24292f;
}

.stat-item p {
  margin: 0;
  font-size: 13px;
  font-weight: 600;
  color: #6b7280;
}

.stat-item .range {
  font-size: 11px;
  color: #9ca3af;
  display: block;
  margin-top: 4px;
}

/* 夜间模式优化 */
html.dark .stat-item h3 {
  color: #e5e7eb;
  font-weight: 500;
}

html.dark .stat-item p {
  color: #d1d5db;
}

html.dark .stat-item .range {
  color: #9ca3af;
}

.hidden-label {
  visibility: hidden;
}

@media (max-width: 1100px) {
  .months-row {
    font-size: 11px;
  }
}

@media (max-width: 768px) {
  /* 移动端卡片内边距优化 */
  .calendar-card :deep(.n-card__content) {
    padding: 16px !important;
  }

  /* 移动端标题优化 */
  .card-header {
    flex-wrap: wrap;
    gap: 4px;
  }

  .title {
    font-size: 14px;
  }

  .subtitle {
    font-size: 11px;
  }

  /* 移动端月份标签优化 */
  .months-row {
    font-size: 10px;
    height: 18px;
    margin-bottom: 6px;
  }

  /* 移动端星期标签优化 */
  .weekdays-col {
    font-size: 10px;
    margin-right: 6px;
  }

  /* 移动端统计信息优化 */
  .stats {
    flex-direction: column;
    gap: 8px;
    padding: 0 8px 4px;
  }

  .stat-item {
    border-top: 1px dashed #e5e7eb;
    padding-top: 8px;
  }

  .stat-item:first-child {
    border-top: none;
    padding-top: 0;
  }

  .stat-item h3 {
    font-size: 22px;
  }

  .stat-item p {
    font-size: 12px;
  }

  /* 移动端优化：通过缩放适配容器，不显示滚动条 */
  .graph-body-outer {
    overflow: hidden; /* 移除滚动，通过缩放完全适配 */
    margin-top: 8px;
  }

  /* 移动端单元格触摸优化 */
  .cell {
    /* 增加触摸区域，提升移动端交互体验 */
    min-width: 12px;
    min-height: 12px;
    touch-action: manipulation; /* 禁用双击缩放 */
  }

  /* 移动端网格优化 */
  .grid {
    gap: 4px; /* 移动端稍微减小间距，节省空间 */
  }

  /* 移动端图例优化 */
  .legend {
    font-size: 10px;
    gap: 4px;
  }

  .legend-box {
    width: 10px;
    height: 10px;
  }

  /* 移动端数据来源优化 */
  .source {
    font-size: 11px;
  }

  /* 移动端：保留原始尺寸，使用横向滚动展示，避免过度缩放 */
  .graph-body-outer {
    overflow-x: auto;
    overflow-y: visible;
    padding-bottom: 8px;
  }

  .graph-body-inner {
    min-width: max-content;
  }

  .months-row {
    min-width: max-content;
  }

  .grid {
    width: max-content;
  }
}
</style>


