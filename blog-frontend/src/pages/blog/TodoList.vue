<template>
  <div class="todo-page">
    <div class="todo-layout">
      <div class="todo-main">
        <n-card class="todo-card" :bordered="false">
          <template #header>
            <div class="card-header">
              <div class="card-title">
                <n-icon :component="ClipboardOutline" size="22" class="title-icon" />
                <span>任务清单</span>
              </div>
              <div class="card-subtitle">
                仅当前登录用户可见，数据存储在浏览器本地
              </div>
            </div>
          </template>

          <!-- 新任务输入 -->
          <div class="todo-input-row">
            <n-input
              v-model:value="newTitle"
              placeholder="记录一条待办事项，例如：阅读一篇技术文章"
              clearable
              @keyup.enter="handleAdd"
            >
              <template #prefix>
                <n-icon :component="AddOutline" />
              </template>
            </n-input>
            <n-button type="primary" @click="handleAdd">
              新增
            </n-button>
          </div>

          <!-- 顶部统计与筛选 -->
          <div class="todo-toolbar">
            <div class="stats">
              <span>总数：{{ todos.length }}</span>
              <n-divider vertical />
              <span>进行中：{{ ongoingCount }}</span>
              <n-divider vertical />
              <span>已完成：{{ completedCount }}</span>
            </div>
            <div class="filters">
              <n-radio-group v-model:value="filter">
                <n-radio-button value="all">全部</n-radio-button>
                <n-radio-button value="ongoing">进行中</n-radio-button>
                <n-radio-button value="completed">已完成</n-radio-button>
                <n-radio-button value="removed">回收站</n-radio-button>
              </n-radio-group>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="todo-actions" v-if="todos.length">
            <n-space :size="12">
              <n-button quaternary size="small" @click="markAllCompleted" :disabled="!ongoingCount">
                完成全部
              </n-button>
              <n-button quaternary size="small" type="error" @click="clearAll" :disabled="!todos.length">
                清空全部
              </n-button>
            </n-space>
          </div>

          <!-- 列表 -->
          <div v-if="filteredTodos.length" class="todo-list">
            <div
              v-for="item in filteredTodos"
              :key="item.id"
              class="todo-item"
              :class="{ completed: item.completed, removed: item.removed }"
            >
              <div class="todo-main">
                <n-checkbox
                  v-model:checked="item.completed"
                  @update:checked="(val: boolean) => toggleCompleted(item, val)"
                >
                  <span class="todo-title" :title="item.title">
                    {{ item.title }}
                  </span>
                </n-checkbox>
              </div>
              <div class="todo-meta">
                <span class="time">{{ formatTime(item.created_at) }}</span>
                <n-space :size="8">
                  <n-button text size="small" @click="startEdit(item)">
                    编辑
                  </n-button>
                  <n-button
                    v-if="!item.removed"
                    text
                    size="small"
                    type="error"
                    @click="moveToTrash(item)"
                  >
                    删除
                  </n-button>
                  <template v-else>
                    <n-button text size="small" @click="restore(item)">
                      还原
                    </n-button>
                    <n-button text size="small" type="error" @click="deleteForever(item)">
                      彻底删除
                    </n-button>
                  </template>
                </n-space>
              </div>
            </div>
          </div>

          <n-empty v-else description="暂无任务，开始添加一条吧~" class="todo-empty" />
        </n-card>
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

    <!-- 编辑对话框 -->
    <n-modal
      v-model:show="showEditModal"
      preset="dialog"
      title="编辑任务"
      :show-icon="false"
      positive-text="保存"
      negative-text="取消"
      @positive-click="confirmEdit"
    >
      <n-input
        v-model:value="editTitle"
        placeholder="修改任务内容"
        autofocus
        @keyup.enter.stop.prevent="confirmEdit"
      />
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useAuthStore } from '@/stores'
import { useMessage, NIcon } from 'naive-ui'
import { AddOutline, ClipboardOutline } from '@vicons/ionicons5'
import AuthorCard from '@/components/AuthorCard.vue'
import AnnouncementBoard from '@/components/AnnouncementBoard.vue'
import TagCloudWidget from '@/components/TagCloudWidget.vue'
import dayjs from 'dayjs'

interface TodoItem {
  id: number
  title: string
  completed: boolean
  removed: boolean
  archived: boolean
  created_at: string
}

const authStore = useAuthStore()
const message = useMessage()

const newTitle = ref('')
const todos = ref<TodoItem[]>([])
const filter = ref<'all' | 'ongoing' | 'completed' | 'removed'>('all')

const showEditModal = ref(false)
const editingTodo = ref<TodoItem | null>(null)
const editTitle = ref('')

const storageKey = computed(() => {
  const userId = authStore.user?.id || 'guest'
  return `todo-list-${userId}`
})

// 统计
const ongoingCount = computed(() => todos.value.filter(t => !t.completed && !t.removed).length)
const completedCount = computed(() => todos.value.filter(t => t.completed && t.archived).length)

// 过滤后的列表
const filteredTodos = computed(() => {
  switch (filter.value) {
    case 'ongoing':
      return todos.value.filter(t => !t.completed && !t.removed && !t.archived)
    case 'completed':
      return todos.value.filter(t => t.completed && t.archived)
    case 'removed':
      return todos.value.filter(t => t.removed)
    default:
      return todos.value.filter(t => !t.removed && !t.archived)
  }
})

function loadFromStorage() {
  try {
    const raw = localStorage.getItem(storageKey.value)
    if (!raw) return
    const parsed = JSON.parse(raw) as TodoItem[]
    if (Array.isArray(parsed)) {
      todos.value = parsed
    }
  } catch (e) {
    console.error('加载任务清单失败', e)
  }
}

function saveToStorage() {
  try {
    localStorage.setItem(storageKey.value, JSON.stringify(todos.value))
  } catch (e) {
    console.error('保存任务清单失败', e)
  }
}

watch(
  () => storageKey.value,
  () => {
    loadFromStorage()
  }
)

watch(
  todos,
  () => {
    saveToStorage()
  },
  { deep: true }
)

function handleAdd() {
  const title = newTitle.value.trim()
  if (!title) {
    message.warning('请输入任务内容')
    return
  }
  const now = new Date().toISOString()
  const nextId = (todos.value[0]?.id || 0) + 1
  todos.value.unshift({
    id: nextId,
    title,
    completed: false,
    removed: false,
    archived: false,
    created_at: now
  })
  newTitle.value = ''
}

function toggleCompleted(item: TodoItem, value: boolean) {
  item.completed = value
  if (value) {
    // 勾选完成时自动归档，不再显示在「全部/进行中」
    item.archived = true
  } else {
    // 取消完成恢复到进行中
    item.archived = false
    item.removed = false
  }
}

function startEdit(item: TodoItem) {
  editingTodo.value = item
  editTitle.value = item.title
  showEditModal.value = true
}

function confirmEdit() {
  if (!editingTodo.value) return
  const title = editTitle.value.trim()
  if (!title) {
    message.warning('任务内容不能为空')
    return
  }
  editingTodo.value.title = title
  showEditModal.value = false
}

function moveToTrash(item: TodoItem) {
  item.removed = true
}

function restore(item: TodoItem) {
  item.removed = false
}

function deleteForever(item: TodoItem) {
  todos.value = todos.value.filter(t => t.id !== item.id)
}

function markAllCompleted() {
  todos.value.forEach(t => {
    if (!t.removed) {
      t.completed = true
      t.archived = true
    }
  })
}

function clearAll() {
  todos.value = todos.value.map(t => ({ ...t, removed: true }))
}

function formatTime(iso: string) {
  return dayjs(iso).format('YYYY-MM-DD HH:mm')
}

onMounted(() => {
  loadFromStorage()
})
</script>

<style scoped>
.todo-page {
  min-height: calc(100vh - 180px);
  padding: 32px 0;
}

.todo-layout {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 20px;
  display: grid;
  grid-template-columns: minmax(0, 1fr) 360px;
  gap: 32px;
  align-items: start;
}

.todo-main {
  min-width: 0;
}

.todo-card {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.5);
  border-radius: 18px;
  box-shadow: 0 12px 32px rgba(15, 23, 42, 0.08);
}

html.dark .todo-card {
  background: rgba(15, 23, 42, 0.85);
  border: 1px solid rgba(148, 163, 184, 0.2);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.5);
}

.card-header {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 700;
}

.title-icon {
  color: #0891b2;
}

html.dark .title-icon {
  color: #38bdf8;
}

.card-subtitle {
  font-size: 13px;
  color: #64748b;
}

html.dark .card-subtitle {
  color: #94a3b8;
}

.todo-input-row {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}

.todo-input-row :deep(.n-input) {
  flex: 1;
}

.todo-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  flex-wrap: wrap;
  gap: 12px;
}

.stats {
  font-size: 13px;
  color: #64748b;
}

html.dark .stats {
  color: #94a3b8;
}

.todo-actions {
  margin-bottom: 12px;
}

.todo-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.todo-item {
  display: flex;
  flex-direction: column;
  padding: 10px 12px;
  border-radius: 12px;
  border: 1px solid rgba(148, 163, 184, 0.35);
  background: rgba(248, 250, 252, 0.9);
  transition: all 0.2s ease;
}

.todo-item:hover {
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.08);
  transform: translateY(-1px);
}

html.dark .todo-item {
  background: rgba(30, 41, 59, 0.9);
  border-color: rgba(148, 163, 184, 0.5);
}

.todo-item.completed .todo-title {
  text-decoration: line-through;
  color: #9ca3af;
}

.todo-main {
  display: flex;
  align-items: center;
}

.todo-title {
  margin-left: 8px;
  font-size: 14px;
  font-weight: 500;
  max-width: 520px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.todo-meta {
  margin-top: 6px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: #9ca3af;
}

html.dark .todo-meta {
  color: #cbd5e1;
}

.time {
  font-family: 'SF Mono', ui-monospace, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New', monospace;
}

.todo-empty {
  margin-top: 24px;
}

.sidebar-section {
  position: relative;
  z-index: 10;
  margin-left: 12px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

@media (max-width: 1024px) {
  .todo-layout {
    grid-template-columns: 1fr;
    padding: 0 16px;
  }

  .sidebar-section {
    display: none;
  }
}

@media (max-width: 768px) {
  .todo-page {
    max-width: 100%;
  }

  .todo-input-row {
    flex-direction: column;
  }

  .todo-title {
    max-width: 100%;
  }
}
</style>


