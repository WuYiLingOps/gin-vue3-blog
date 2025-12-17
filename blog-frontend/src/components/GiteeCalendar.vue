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

        <div v-if="flatDays.length" class="graph-body-container">
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
              <span v-for="(month, idx) in monthLabels" :key="idx">{{ month }}</span>
            </div>

            <div class="grid">
              <div
                v-for="day in flatDays"
                :key="day.date"
                class="cell"
                :class="cellLevelClass(day.count)"
                :data-date="day.date"
                :data-count="day.count"
              />
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
import { onMounted, ref, computed } from 'vue'
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

const dateRange = computed(() => {
  if (!flatDays.value.length) return ''
  const first = flatDays.value[0].date
  const last = flatDays.value[flatDays.value.length - 1].date
  return `${first} - ${last}`
})

const lastWeekStats = computed(() => {
  if (!flatDays.value.length) return { total: 0, from: '', to: '' }
  const len = flatDays.value.length
  const last = flatDays.value[len - 1]
  const to = last.date
  const fromIdx = Math.max(0, len - 7)
  const slice = flatDays.value.slice(fromIdx, len)
  const sum = slice.reduce((s, d) => s + d.count, 0)
  return { total: sum, from: slice[0].date, to }
})

const lastMonthStats = computed(() => {
  if (!flatDays.value.length) return { total: 0, from: '', to: '' }
  const last = flatDays.value[flatDays.value.length - 1]
  const to = last.date
  const fromDate = new Date(to)
  fromDate.setDate(fromDate.getDate() - 29)
  const fromStr = fromDate.toISOString().slice(0, 10)
  const slice = flatDays.value.filter((d) => d.date >= fromStr && d.date <= to)
  const sum = slice.reduce((s, d) => s + d.count, 0)
  return { total: sum, from: slice.length ? slice[0].date : fromStr, to }
})

// 月份标签（静态 12 个月）
const monthLabels = ['一月', '二月', '三月', '四月', '五月', '六月', '七月', '八月', '九月', '十月', '十一月', '十二月']

// 颜色等级
function cellLevelClass(count: number) {
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

// go-code-calendar-api 的后端地址，建议在生产环境用环境变量覆盖
function getCalendarApiBase() {
  return import.meta.env.VITE_GITEE_CALENDAR_API || 'http://localhost:8081/api'
}

// 构造一整年的「空」数据，用于无数据 / 请求失败时仍然展示完整网格样式
function buildEmptyWeeks(): CalendarWeek[] {
  const days: CalendarDay[] = []
  const today = new Date()

  // 过去 365 天
  for (let i = 364; i >= 0; i--) {
    const d = new Date(today)
    d.setDate(today.getDate() - i)
    const dateStr = d.toISOString().slice(0, 10)
    days.push({ date: dateStr, count: 0 })
  }

  const weeksArr: CalendarWeek[] = []
  for (let i = 0; i < days.length; i += 7) {
    weeksArr.push(days.slice(i, i + 7))
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
    const data = (await res.json()) as { total: number; contributions: CalendarWeek[] }
    if (!data || !Array.isArray(data.contributions)) {
      throw new Error('接口返回数据格式不正确')
    }
    if (!data.contributions.length) {
      total.value = 0
      weeks.value = buildEmptyWeeks()
    } else {
      total.value = data.total || 0
      weeks.value = data.contributions.map((week) =>
        week.map((d) => ({
          date: d.date,
          count: d.count
        }))
      )
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
  }
}

onMounted(async () => {
  try {
    await loadUsername()
    await fetchData()
  } catch {
    // 错误已在内部处理
  }
})
</script>

<style scoped>
.hexo-calendar-card {
  width: 100%;
}

.calendar-card {
  height: 100%;
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

.calendar-loading {
  padding: 24px;
  text-align: center;
  font-size: 14px;
  color: #64748b;
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

/* --- 核心布局，参考 go-code-calendar-api/web/index.html --- */
.graph-body-container {
  display: flex;
  align-items: flex-end;
  margin-top: 12px;
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

/* 右侧：包含月份和网格的内容区 */
.graph-content-col {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-width: 0;
}

/* 顶部月份栏 */
.months-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 12px;
  color: #6e7781;
  width: 100%;
  box-sizing: border-box;
  padding: 0 2px;
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

.cell:hover::after {
  content: attr(data-date) ' · ' attr(data-count) ' 次提交';
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
  margin-bottom: 5px;
  padding: 4px 8px;
  font-size: 11px;
  color: #fff;
  background: rgba(0, 0, 0, 0.9);
  border-radius: 4px;
  white-space: nowrap;
  pointer-events: none;
  z-index: 10;
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

.hidden-label {
  visibility: hidden;
}

@media (max-width: 1100px) {
  .grid {
    grid-auto-columns: 10px;
  }

  .months-row span:nth-child(odd) {
    display: none;
  }
}

@media (max-width: 768px) {
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
}
</style>


