<!--
 * @ProjectName: go-vue3-blog
 * @FileName: ChatManage.vue
 * @CreateTime: 2026-02-02 11:43:17
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 聊天室管理页面组件，提供聊天消息的管理功能
 -->
<template>
  <div class="chat-manage">
    <n-card title="聊天室管理">
      <template #header-extra>
        <n-space>
          <n-button type="primary" @click="showBroadcastModal = true">
            发送系统广播
          </n-button>
          <n-button @click="fetchMessages">
            <template #icon>
              <n-icon>
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                  <path fill="currentColor" d="M17.65 6.35A7.958 7.958 0 0 0 12 4c-4.42 0-7.99 3.58-7.99 8s3.57 8 7.99 8c3.73 0 6.84-2.55 7.73-6h-2.08A5.99 5.99 0 0 1 12 18c-3.31 0-6-2.69-6-6s2.69-6 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4l-2.35 2.35z"/>
                </svg>
              </n-icon>
            </template>
            刷新
          </n-button>
        </n-space>
      </template>

      <!-- 全员禁言开关 -->
      <n-card size="small" style="margin-bottom: 16px">
        <n-space align="center" justify="space-between">
          <div>
            <div style="font-weight: 600;">全员禁言</div>
            <n-text depth="3">开启后仅管理员可发言，防止刷屏攻击</n-text>
          </div>
          <n-switch
            :value="chatSettings.chat_mute_all === '1'"
            :loading="chatSettingLoading"
            @update:value="handleToggleChatMute"
          />
        </n-space>
      </n-card>

      <!-- 统计信息 -->
      <n-space class="stats-section" size="large">
        <n-statistic label="在线人数" :value="onlineInfo.online_count">
          <template #prefix>
            <n-icon>
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                <path fill="currentColor" d="M16 11c1.66 0 2.99-1.34 2.99-3S17.66 5 16 5c-1.66 0-3 1.34-3 3s1.34 3 3 3zm-8 0c1.66 0 2.99-1.34 2.99-3S9.66 5 8 5C6.34 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5c0-2.33-4.67-3.5-7-3.5zm8 0c-.29 0-.62.02-.97.05 1.16.84 1.97 1.97 1.97 3.45V19h6v-2.5c0-2.33-4.67-3.5-7-3.5z"/>
              </svg>
            </n-icon>
          </template>
        </n-statistic>
        <n-statistic label="消息总数" :value="total">
          <template #prefix>
            <n-icon>
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                <path fill="currentColor" d="M20 2H4c-1.1 0-1.99.9-1.99 2L2 22l4-4h14c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zM6 9h12v2H6V9zm8 5H6v-2h8v2zm4-6H6V6h12v2z"/>
              </svg>
            </n-icon>
          </template>
        </n-statistic>
      </n-space>

      <n-divider />

      <!-- 在线用户列表 -->
      <n-card title="在线用户" size="small" style="margin-bottom: 20px">
        <n-space v-if="onlineInfo.online_users && onlineInfo.online_users.length > 0" vertical>
          <n-space
            v-for="user in onlineInfo.online_users"
            :key="user.id"
            justify="space-between"
            align="center"
            class="online-user-item"
          >
            <n-space align="center">
              <n-avatar :size="32" :src="user.avatar" />
              <n-text>{{ user.username }}</n-text>
            </n-space>
            <n-space>
              <n-button
                size="small"
                type="warning"
                @click="handleKickUser(user)"
              >
                踢出
              </n-button>
              <n-button
                size="small"
                type="error"
                @click="handleBanUser(user)"
              >
                封禁
              </n-button>
            </n-space>
          </n-space>
        </n-space>
        <n-empty v-else description="暂无在线用户" />
      </n-card>

      <n-divider />

      <!-- 批量操作 -->
      <n-space v-if="selectedRowKeys.length > 0" style="margin-bottom: 16px">
        <n-button type="error" @click="handleBatchDelete">
          批量删除 ({{ selectedRowKeys.length }})
        </n-button>
        <n-button @click="selectedRowKeys = []">
          取消选择
        </n-button>
      </n-space>

      <!-- 消息列表 -->
      <div class="content-area">
        <n-data-table
          :columns="columns"
          :data="messages"
          :loading="loading"
          :row-key="(row: ChatMessage) => row.id"
          v-model:checked-row-keys="selectedRowKeys"
          @update:checked-row-keys="handleCheckedRowKeysChange"
        />
        
        <!-- 分页 - 位于表格右下角 -->
        <div class="pagination-wrapper">
          <n-pagination
            v-if="total > 0"
            v-model:page="currentPage"
            :page-count="totalPages"
            :page-size="pageSize"
            :page-slot="7"
            @update:page="handlePageChange"
          />
        </div>
      </div>
    </n-card>

    <!-- 系统广播对话框 -->
    <n-modal v-model:show="showBroadcastModal">
      <n-card
        title="发送系统广播"
        :bordered="false"
        size="large"
        style="max-width: 600px"
        closable
        @close="showBroadcastModal = false"
      >
        <n-form>
          <n-form-item label="优先级">
            <n-select
              v-model:value="broadcastPriority"
              :options="[
                { label: '置顶', value: 1 },
                { label: '普通', value: 0 }
              ]"
            />
          </n-form-item>
          <n-form-item label="投递到">
            <n-select
              v-model:value="broadcastTarget"
              :options="[
                { label: '公告栏', value: 'announcement' },
                { label: '聊天室', value: 'chat' },
                { label: '同时', value: 'both' }
              ]"
            />
          </n-form-item>
          <n-form-item label="广播内容">
            <n-input
              v-model:value="broadcastContent"
              type="textarea"
              placeholder="请输入广播内容"
              :autosize="{ minRows: 3, maxRows: 6 }"
            />
          </n-form-item>
        </n-form>
        <template #footer>
          <n-space justify="end">
            <n-button @click="showBroadcastModal = false">取消</n-button>
            <n-button type="primary" :loading="broadcasting" @click="handleBroadcast">
              发送
            </n-button>
          </n-space>
        </template>
      </n-card>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, h } from 'vue'
import {
  NCard,
  NSpace,
  NButton,
  NIcon,
  NStatistic,
  NDivider,
  NDataTable,
  NModal,
  NForm,
  NFormItem,
  NInput,
  NSelect,
  NPopconfirm,
  NTime,
  NAvatar,
  NText,
  NEmpty,
  NTag,
  useMessage,
  useDialog,
  type DataTableColumns
} from 'naive-ui'
import { 
  adminGetMessages, 
  adminDeleteMessage, 
  adminBroadcastMessage, 
  getOnlineInfo,
  adminKickUser,
  adminBanIP,
  getChatSettings,
  updateChatSettings
} from '@/api/chat'
import type { ChatMessage, OnlineUser, OnlineInfo } from '@/api/chat'

const message = useMessage()
const dialog = useDialog()

// 数据
const messages = ref<ChatMessage[]>([])
const loading = ref(false)
const onlineInfo = ref<OnlineInfo>({
  online_count: 0,
  online_users: []
})

// 分页 - 固定每页显示15条消息
const currentPage = ref(1)
const pageSize = 15 // 固定每页显示15条消息
const total = ref(0)

// 计算总页数
const totalPages = computed(() => Math.ceil(total.value / pageSize))

// 系统广播
const showBroadcastModal = ref(false)
const broadcastContent = ref('')
const broadcastPriority = ref(0)
const broadcastTarget = ref<'announcement' | 'chat' | 'both'>('both')
const broadcasting = ref(false)
const chatSettings = ref({ chat_mute_all: '0' })
const chatSettingLoading = ref(false)

// 批量删除
const selectedRowKeys = ref<Array<string | number>>([])

// 表格列
const columns: DataTableColumns<ChatMessage> = [
  {
    type: 'selection'
  },
  {
    title: 'ID',
    key: 'id',
    width: 80
  },
  {
    title: '用户名',
    key: 'username',
    width: 150,
    render: (row) => {
      if (row.is_broadcast) {
        return h(NSpace, { align: 'center', size: 'small' }, {
          default: () => [
            h(NTag, { type: 'error', size: 'small' }, { default: () => '系统' }),
            h('span', row.username)
          ]
        })
      }
      return row.username
    }
  },
  {
    title: '消息内容',
    key: 'content',
    ellipsis: {
      tooltip: true
    }
  },
  {
    title: '投递目标',
    key: 'target',
    width: 140,
    render: (row) => {
      if (!row.is_broadcast) {
        return h('span', { style: { color: '#999' } }, '普通消息')
      }
      const targetMap: Record<string, { label: string; type: 'default' | 'info' | 'success' | 'warning' | 'error' }> = {
        'announcement': { label: '公告栏', type: 'info' },
        'chat': { label: '聊天室', type: 'success' },
        'both': { label: '同时', type: 'warning' }
      }
      const target = row.target || 'announcement'
      const config = targetMap[target] || { label: '未知', type: 'default' }
      return h(NTag, { type: config.type, size: 'small' }, { default: () => config.label })
    }
  },
  {
    title: 'IP地址',
    key: 'ip',
    width: 140
  },
  {
    title: '发送时间',
    key: 'created_at',
    width: 180,
    render: (row) => h(NTime, { time: new Date(row.created_at) })
  },
  {
    title: '操作',
    key: 'actions',
    width: 120,
    render: (row) => {
      return h(
        NPopconfirm,
        {
          onPositiveClick: () => handleDelete(row.id)
        },
        {
          trigger: () =>
            h(
              NButton,
              {
                size: 'small',
                type: 'error'
              },
              { default: () => '删除' }
            ),
          default: () => '确定删除这条消息吗？'
        }
      )
    }
  }
]

// 获取消息列表
const fetchMessages = async () => {
  loading.value = true
  try {
    const res = await adminGetMessages({
      page: currentPage.value,
      page_size: pageSize
    })
    messages.value = res.data.list || []
    total.value = res.data.total || 0
    // 确保页码不超过最大页数
    const maxPage = Math.ceil(total.value / pageSize) || 1
    if (currentPage.value > maxPage && maxPage > 0) {
      currentPage.value = maxPage
    }
  } catch (error) {
    message.error('获取消息列表失败')
  } finally {
    loading.value = false
  }
}

// 获取在线信息
const fetchOnlineInfo = async () => {
  try {
    const res = await getOnlineInfo()
    onlineInfo.value = res.data
  } catch (error) {
    console.error('获取在线信息失败:', error)
  }
}

// 获取聊天室配置
const fetchChatSettingsData = async () => {
  try {
    const res = await getChatSettings()
    if (res.data) {
      chatSettings.value = res.data
    }
  } catch (error) {
    console.error('获取聊天室配置失败:', error)
  }
}

// 切换全员禁言
const handleToggleChatMute = async (val: boolean) => {
  chatSettingLoading.value = true
  try {
    await updateChatSettings({ chat_mute_all: val ? '1' : '0' })
    chatSettings.value.chat_mute_all = val ? '1' : '0'
    message.success(val ? '已开启全员禁言' : '已关闭全员禁言')
  } catch (error: any) {
    console.error('更新聊天室配置失败', error)
    message.error(error?.message || '更新失败')
  } finally {
    chatSettingLoading.value = false
  }
}

// 踢出用户
const handleKickUser = (user: OnlineUser) => {
  dialog.warning({
    title: '踢出用户',
    content: `确定要踢出用户 ${user.username} 吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await adminKickUser(user.id, '违反聊天室规则')
        message.success('踢出成功')
        fetchOnlineInfo()
      } catch (error) {
        message.error('踢出失败')
      }
    }
  })
}

// 封禁用户
const handleBanUser = (user: OnlineUser) => {
  dialog.error({
    title: '封禁用户',
    content: `确定要封禁用户 ${user.username} 的IP吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await adminBanIP(user.id, '违反聊天室规则', 24) // 封禁24小时
        message.success('封禁成功')
        fetchOnlineInfo()
      } catch (error) {
        message.error('封禁失败')
      }
    }
  })
}

// 删除消息
const handleDelete = async (id: number) => {
  try {
    await adminDeleteMessage(id)
    message.success('删除成功')
    fetchMessages()
  } catch (error) {
    message.error('删除失败')
  }
}

// 批量删除
const handleBatchDelete = () => {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请先选择要删除的消息')
    return
  }

  dialog.warning({
    title: '批量删除',
    content: `确定要删除选中的 ${selectedRowKeys.value.length} 条消息吗？此操作不可恢复。`,
    positiveText: '确定删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      loading.value = true
      try {
        // 批量删除所有选中的消息
        await Promise.all(selectedRowKeys.value.map(id => adminDeleteMessage(Number(id))))
        message.success(`成功删除 ${selectedRowKeys.value.length} 条消息`)
        selectedRowKeys.value = []
        fetchMessages()
      } catch (error) {
        message.error('批量删除失败')
      } finally {
        loading.value = false
      }
    }
  })
}

// 处理选择变化
const handleCheckedRowKeysChange = (keys: Array<string | number>) => {
  selectedRowKeys.value = keys
}

// 发送系统广播
const handleBroadcast = async () => {
  if (!broadcastContent.value.trim()) {
    message.error('请输入广播内容')
    return
  }

  broadcasting.value = true
  try {
    await adminBroadcastMessage(broadcastContent.value.trim(), broadcastPriority.value, broadcastTarget.value)
    message.success('广播发送成功')
    showBroadcastModal.value = false
    broadcastContent.value = ''
    broadcastPriority.value = 0
    broadcastTarget.value = 'both'
    fetchMessages()
  } catch (error) {
    message.error('广播发送失败')
  } finally {
    broadcasting.value = false
  }
}

// 分页处理
const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchMessages()
}

// 定时器引用
let onlineInfoTimer: number | null = null

// 初始化
onMounted(() => {
  fetchMessages()
  fetchOnlineInfo()
  fetchChatSettingsData()
  
  // 定时刷新在线人数（每10秒）
  onlineInfoTimer = window.setInterval(() => {
    // 只在页面可见时刷新
    if (!document.hidden) {
      fetchOnlineInfo()
    }
  }, 10000)
})

// 组件卸载时清理定时器
onBeforeUnmount(() => {
  if (onlineInfoTimer !== null) {
    clearInterval(onlineInfoTimer)
    onlineInfoTimer = null
  }
})
</script>

<style scoped>
.chat-manage {
  padding: 20px;
  position: relative;
}

.stats-section {
  margin-bottom: 20px;
}

.online-user-item {
  padding: 12px;
  background: #f9f9f9;
  border-radius: 8px;
  transition: background-color 0.2s;
}

.online-user-item:hover {
  background: #f0f0f0;
}

/* 内容区域 */
.content-area {
  position: relative;
  padding-bottom: 50px; /* 为分页器预留空间 */
}

/* 分页样式 - 位于表格右下角 */
.pagination-wrapper {
  position: absolute;
  bottom: 10px; /* 距离表格底部 */
  right: 20px; /* 距离右侧 */
  z-index: 10;
  display: flex;
  justify-content: flex-end;
  align-items: center;
  box-sizing: border-box;
}

/* 移动端样式 */
@media (max-width: 768px) {
  .pagination-wrapper {
    bottom: 8px;
    right: 16px;
  }
}
</style>

