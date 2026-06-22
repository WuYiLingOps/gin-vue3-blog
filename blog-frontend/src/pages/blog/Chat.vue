<!--
 * @ProjectName: go-vue3-blog
 * @FileName: Chat.vue
 * @CreateTime: 2026-02-02 11:43:17
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 聊天室页面组件，提供实时聊天功能
 -->
<template>
  <div class="chat-page">
    <n-card title="聊天室" class="chat-container" :theme-overrides="cardThemeOverrides">
      <template #header-extra>
        <n-space align="center">
          <n-badge :value="onlineCount" :max="99" type="success">
            <n-icon size="20">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                <path
                  fill="currentColor"
                  d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"
                />
              </svg>
            </n-icon>
          </n-badge>
          <n-text depth="3">{{ onlineCount }} 人在线</n-text>
        </n-space>
      </template>

      <div class="chat-layout">
        <!-- 聊天消息区域 -->
        <div class="chat-messages" ref="messagesContainer">
          <div v-if="messages.length === 0" class="empty-messages">
            <n-empty description="暂无消息" />
          </div>

          <div v-else class="messages-list">
            <div
              v-for="msg in messages"
              :key="msg.id"
              class="message-item"
              :class="{ 'own-message': isOwnMessage(msg) }"
            >
              <!-- 头像右键菜单 -->
              <n-dropdown
                v-if="authStore.isAdmin && !isOwnMessage(msg) && msg.client_id"
                trigger="manual"
                placement="bottom-start"
                :show="avatarDropdownShow === msg.id"
                :options="avatarMenuOptions"
                @select="key => handleAvatarMenuSelect(key, msg)"
                @clickoutside="avatarDropdownShow = null"
              >
                <div class="message-avatar" @contextmenu.prevent="e => showAvatarDropdown(e, msg)">
                  <n-avatar v-if="msg.avatar" :size="40" :src="msg.avatar" />
                  <n-avatar
                    v-else
                    :size="40"
                    :style="{
                      backgroundColor: getAvatarColor(msg.username),
                      color: 'white',
                      fontWeight: 'bold'
                    }"
                  >
                    {{ getAvatarText(msg.username) }}
                  </n-avatar>
                </div>
              </n-dropdown>
              <div v-else class="message-avatar">
                <n-avatar v-if="msg.avatar" :size="40" :src="msg.avatar" />
                <n-avatar
                  v-else
                  :size="40"
                  :style="{
                    backgroundColor: getAvatarColor(msg.username),
                    color: 'white',
                    fontWeight: 'bold'
                  }"
                >
                  {{ getAvatarText(msg.username) }}
                </n-avatar>
              </div>

              <!-- 消息内容右键菜单 -->
              <n-dropdown
                v-if="authStore.isAdmin && !isOwnMessage(msg)"
                trigger="manual"
                placement="bottom-start"
                :show="messageDropdownShow === msg.id"
                :options="messageMenuOptions"
                @select="key => handleMessageMenuSelect(key, msg)"
                @clickoutside="messageDropdownShow = null"
              >
                <div
                  class="message-content"
                  @contextmenu.prevent="e => showMessageDropdown(e, msg)"
                >
                  <div class="message-header">
                    <span class="message-username">{{ msg.username }}</span>
                    <span class="message-time">{{ formatTime(msg.created_at) }}</span>
                  </div>
                  <div class="message-text">{{ msg.content }}</div>
                </div>
              </n-dropdown>
              <div v-else class="message-content">
                <div class="message-header">
                  <span class="message-username">{{ msg.username }}</span>
                  <span class="message-time">{{ formatTime(msg.created_at) }}</span>
                </div>
                <div class="message-text">{{ msg.content }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 输入区域 -->
        <div class="chat-input">
          <n-space vertical size="small" style="width: 100%">
            <div class="input-wrapper">
              <n-input
                v-model:value="messageInput"
                type="textarea"
                :placeholder="
                  isChatMutedForUser ? '已开启全员禁言，只有管理员可发言' : '输入消息...'
                "
                :disabled="isChatMutedForUser"
                :autosize="{ minRows: 2, maxRows: 4 }"
                @keydown.enter.prevent="handleSendMessage"
              />
              <n-popover trigger="click" placement="top-start">
                <template #trigger>
                  <n-button text class="emoji-btn" title="选择表情">
                    <template #icon>
                      <n-icon :size="28">
                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
                          <path
                            fill="currentColor"
                            d="M256 48C141.1 48 48 141.1 48 256s93.1 208 208 208 208-93.1 208-208S370.9 48 256 48zm0 374c-91.4 0-166-74.6-166-166S164.6 90 256 90s166 74.6 166 166-74.6 166-166 166z"
                          />
                          <circle cx="192" cy="208" r="24" fill="currentColor" />
                          <circle cx="320" cy="208" r="24" fill="currentColor" />
                          <path
                            fill="currentColor"
                            d="M256 304c-33.1 0-62.6 17.5-79.2 43.8-3.5 5.6-.9 12.9 5.3 15.1 2 .7 4.1.9 6.1.9 4.4 0 8.6-2.3 10.9-6.3 12.3-19.6 33.6-31.3 57-31.3s44.7 11.7 57 31.3c3.5 5.6 10.9 7.3 16.5 3.8s7.3-10.9 3.8-16.5C318.6 321.5 289.1 304 256 304z"
                          />
                        </svg>
                      </n-icon>
                    </template>
                  </n-button>
                </template>
                <div class="emoji-picker">
                  <div
                    v-for="emoji in emojis"
                    :key="emoji"
                    class="emoji-item"
                    @click="insertEmoji(emoji)"
                  >
                    {{ emoji }}
                  </div>
                </div>
              </n-popover>
            </div>
            <n-space justify="space-between">
              <n-text depth="3" style="font-size: 12px">
                按 Enter 发送，Shift + Enter 换行
                <span v-if="isChatMutedForUser" style="color: #f59e0b; margin-left: 8px"
                  >已开启全员禁言</span
                >
              </n-text>
              <n-button
                type="primary"
                :disabled="!messageInput.trim() || !isConnected || isChatMutedForUser"
                @click="handleSendMessage"
              >
                发送
              </n-button>
            </n-space>
          </n-space>
        </div>
      </div>
    </n-card>

    <!-- 用户设置对话框 -->
    <n-modal v-model:show="showUserSetup" :mask-closable="false">
      <n-card title="设置用户信息" :bordered="false" size="large" style="max-width: 500px">
        <n-form>
          <n-form-item label="昵称">
            <n-input v-model:value="userSetup.username" placeholder="请输入昵称" maxlength="20" />
          </n-form-item>
        </n-form>
        <template #footer>
          <n-space justify="end">
            <n-button type="primary" @click="confirmUserSetup"> 确定 </n-button>
          </n-space>
        </template>
      </n-card>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick, h, computed } from 'vue'
import {
  NCard,
  NSpace,
  NIcon,
  NBadge,
  NText,
  NEmpty,
  NAvatar,
  NInput,
  NButton,
  NModal,
  NForm,
  NFormItem,
  NPopover,
  NDropdown,
  useMessage,
  useDialog
} from 'naive-ui'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { createChatWebSocket, ChatWebSocket } from '@/utils/websocket'
import { adminDeleteMessage, adminKickUser, type ChatMessage, type OnlineUser } from '@/api/chat'
import { formatDistanceToNow } from '@/utils/format'
import request from '@/utils/request'

const authStore = useAuthStore()
const appStore = useAppStore()
const message = useMessage()
const dialog = useDialog()

// 卡片暗色主题覆盖
const cardThemeOverrides = computed(() => {
  if (appStore.theme !== 'dark') return {}
  return {
    color: '#1a2332',
    colorEmbedded: '#1a2332',
    borderColor: '#243447',
    titleTextColor: '#e2e8f0',
    textColor: '#e5e7eb'
  }
})

// WebSocket连接
let ws: ChatWebSocket | null = null
const isConnected = ref(false)

// 聊天数据
const messages = ref<ChatMessage[]>([])
const messageInput = ref('')
const onlineCount = ref(0)
const onlineUsers = ref<OnlineUser[]>([])
type ChatSettingState = { chat_mute_all: string }
const chatSettings = ref<ChatSettingState>({ chat_mute_all: '0' })
const isChatMutedForUser = computed(
  () => chatSettings.value.chat_mute_all === '1' && !authStore.isAdmin
)

// 用户设置
const showUserSetup = ref(false)
const userSetup = ref({
  username: ''
})

// 消息容器引用
const messagesContainer = ref<HTMLElement>()

// 右键菜单状态
const messageDropdownShow = ref<number | null>(null)
const avatarDropdownShow = ref<number | null>(null)

// 消息右键菜单选项
const messageMenuOptions = [
  {
    label: '删除消息',
    key: 'delete',
    icon: () =>
      h(NIcon, null, {
        default: () =>
          h(
            'svg',
            {
              xmlns: 'http://www.w3.org/2000/svg',
              viewBox: '0 0 24 24'
            },
            [
              h('path', {
                fill: 'currentColor',
                d: 'M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z'
              })
            ]
          )
      })
  }
]

// 头像右键菜单选项
const avatarMenuOptions = [
  {
    label: '踢出用户',
    key: 'kick',
    icon: () =>
      h(NIcon, null, {
        default: () =>
          h(
            'svg',
            {
              xmlns: 'http://www.w3.org/2000/svg',
              viewBox: '0 0 24 24'
            },
            [
              h('path', {
                fill: 'currentColor',
                d: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 11c-.55 0-1-.45-1-1V8c0-.55.45-1 1-1s1 .45 1 1v4c0 .55-.45 1-1 1zm1 4h-2v-2h2v2z'
              })
            ]
          )
      })
  }
]

// 常用表情列表
const emojis = [
  '😀',
  '😃',
  '😄',
  '😁',
  '😆',
  '😅',
  '🤣',
  '😂',
  '🙂',
  '🙃',
  '😉',
  '😊',
  '😇',
  '🥰',
  '😍',
  '🤩',
  '😘',
  '😗',
  '😚',
  '😙',
  '😋',
  '😛',
  '😜',
  '🤪',
  '😝',
  '🤑',
  '🤗',
  '🤭',
  '🤫',
  '🤔',
  '🤐',
  '🤨',
  '😐',
  '😑',
  '😶',
  '😏',
  '😒',
  '🙄',
  '😬',
  '🤥',
  '😌',
  '😔',
  '😪',
  '🤤',
  '😴',
  '😷',
  '🤒',
  '🤕',
  '🤢',
  '🤮',
  '🤧',
  '🥵',
  '🥶',
  '😶‍🌫️',
  '😵',
  '🤯',
  '🤠',
  '🥳',
  '😎',
  '🤓',
  '🧐',
  '😕',
  '😟',
  '🙁',
  '☹️',
  '😮',
  '😯',
  '😲',
  '😳',
  '🥺',
  '😦',
  '😧',
  '😨',
  '😰',
  '😥',
  '😢',
  '😭',
  '😱',
  '😖',
  '😣',
  '😞',
  '😓',
  '😩',
  '😫',
  '🥱',
  '😤',
  '😡',
  '😠',
  '🤬',
  '😈',
  '👿',
  '💀',
  '☠️',
  '💩',
  '🤡',
  '👹',
  '👺',
  '👻',
  '👽',
  '👾',
  '🤖',
  '😺',
  '😸',
  '😹',
  '👍',
  '👎',
  '👌',
  '✌️',
  '🤞',
  '🤟',
  '🤘',
  '🤙',
  '👋',
  '🤚',
  '🖐️',
  '✋',
  '🖖',
  '👏',
  '🙌',
  '👐',
  '🤝',
  '🙏',
  '✍️',
  '💪',
  '🦾',
  '🦿',
  '🦵',
  '🦶',
  '❤️',
  '🧡',
  '💛',
  '💚',
  '💙',
  '💜',
  '🖤',
  '🤍',
  '💔',
  '❣️',
  '💕',
  '💞',
  '💓',
  '💗',
  '💖',
  '💘',
  '💝',
  '💟',
  '☮️',
  '✝️',
  '☪️',
  '🕉️',
  '☸️',
  '✡️',
  '🔯',
  '🕎',
  '☯️',
  '☦️',
  '🛐',
  '⛎',
  '♈',
  '♉'
]

// 获取头像文字（用户名首字母）
const getAvatarText = (username: string) => {
  if (!username) return '?'
  // 如果是中文，取第一个字；如果是英文，取首字母大写
  const firstChar = username.charAt(0)
  return /[\u4e00-\u9fa5]/.test(firstChar) ? firstChar : firstChar.toUpperCase()
}

// 根据用户名生成头像颜色（类似微信的柔和配色）
const getAvatarColor = (username: string) => {
  const colors = [
    '#7DB9DE', // 浅蓝
    '#FF9999', // 浅红
    '#95D5B2', // 浅绿
    '#FFB366', // 浅橙
    '#C499BA', // 浅紫
    '#74C0C2', // 青色
    '#F4A9B8', // 粉色
    '#FFD93D', // 金黄
    '#B8B8D1', // 淡紫
    '#6CB4EE', // 天蓝
    '#98D8C8', // 薄荷绿
    '#F4B4A5' // 珊瑚色
  ]

  // 使用用户名生成一个稳定的哈希值
  let hash = 0
  for (let i = 0; i < username.length; i++) {
    hash = username.charCodeAt(i) + ((hash << 5) - hash)
    hash = hash & hash // Convert to 32bit integer
  }

  return colors[Math.abs(hash) % colors.length]
}

// 判断是否是自己的消息
const isOwnMessage = (msg: ChatMessage) => {
  // 已登录用户只用 user_id 比较，避免同名用户误判
  if (authStore.isLoggedIn && authStore.user?.id) {
    return (
      msg.user_id !== null &&
      msg.user_id !== undefined &&
      Number(msg.user_id) === Number(authStore.user.id)
    )
  }
  // 匿名用户比较用户名
  if (userSetup.value.username) {
    return msg.username === userSetup.value.username
  }
  return false
}

// 格式化时间
const formatTime = (time: string) => {
  return formatDistanceToNow(new Date(time))
}

// 滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

// 获取聊天室配置（用于禁言提示）
const fetchChatSettingsData = async () => {
  try {
    const res = await request.get<ChatSettingState>('/chat/settings')
    chatSettings.value = res.data || { chat_mute_all: '0' }
  } catch (error) {
    console.error('获取聊天室配置失败', error)
  }
}

// 连接WebSocket
const connectWebSocket = () => {
  const username = authStore.isLoggedIn ? authStore.user?.username : userSetup.value.username

  const avatar = authStore.isLoggedIn ? authStore.user?.avatar : undefined

  // 如果用户已登录，传递token
  const token = authStore.isLoggedIn && authStore.token ? authStore.token : undefined

  ws = createChatWebSocket(username, avatar, token)

  // 连接成功
  ws.on('open', () => {
    isConnected.value = true
    message.success('已连接到聊天室')
  })

  // 连接关闭
  ws.on('close', () => {
    isConnected.value = false
    message.warning('已断开连接')
  })

  // 连接错误
  ws.on('error', () => {
    message.error('连接失败')
  })

  // 接收历史消息
  ws.on('history', (data: ChatMessage[]) => {
    // 过滤掉仅投递到公告栏的系统广播
    messages.value = data.filter(msg => !(msg.is_broadcast && msg.target === 'announcement'))
    scrollToBottom()
  })

  // 接收新消息
  ws.on('message', (data: ChatMessage) => {
    // 普通消息或投递到聊天室/同时的广播才显示
    if (data.is_broadcast && data.target === 'announcement') {
      return
    }
    messages.value.push(data)
    scrollToBottom()
  })

  // 用户加入（不再自己维护计数，等待 user_list 更新）
  ws.on('user_join', () => {
    // 不做任何操作，等待后端发送完整的 user_list
  })

  // 用户离开（不再自己维护计数，等待 user_list 更新）
  ws.on('user_leave', () => {
    // 不做任何操作，等待后端发送完整的 user_list
  })

  // 在线用户列表（唯一的在线数据来源）
  ws.on('user_list', (data: OnlineUser[]) => {
    onlineUsers.value = data
    onlineCount.value = data.length
  })

  // 系统消息
  ws.on('system', (data: any) => {
    // 如果仅投递到公告栏，则不在聊天室展示
    if (data?.is_broadcast && data?.target === 'announcement') {
      return
    }
    const content = data?.content || data?.message
    if (content) {
      message.info('系统消息: ' + content)
    }
    if (data && data.content) {
      messages.value.push(data as ChatMessage)
      scrollToBottom()
    }
  })

  // 被踢出
  ws.on('kick', (data: any) => {
    message.error(data.reason || '您已被踢出聊天室')
    setTimeout(() => {
      ws?.close()
      isConnected.value = false
    }, 1000)
  })

  ws.connect().catch(err => {
    console.error('连接WebSocket失败:', err)
  })
}

// 显示消息右键菜单
const showMessageDropdown = (e: MouseEvent, msg: ChatMessage) => {
  e.preventDefault()
  messageDropdownShow.value = msg.id
  avatarDropdownShow.value = null
}

// 显示头像右键菜单
const showAvatarDropdown = (e: MouseEvent, msg: ChatMessage) => {
  e.preventDefault()
  avatarDropdownShow.value = msg.id
  messageDropdownShow.value = null
}

// 处理消息菜单选择
const handleMessageMenuSelect = async (key: string, msg: ChatMessage) => {
  messageDropdownShow.value = null

  if (key === 'delete') {
    dialog.warning({
      title: '删除消息',
      content: `确定要删除 ${msg.username} 的消息吗？`,
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: async () => {
        try {
          await adminDeleteMessage(msg.id)
          message.success('删除成功')
          // 从本地消息列表中移除
          const index = messages.value.findIndex(m => m.id === msg.id)
          if (index > -1) {
            messages.value.splice(index, 1)
          }
        } catch (error) {
          message.error('删除失败')
        }
      }
    })
  }
}

// 处理头像菜单选择
const handleAvatarMenuSelect = async (key: string, msg: ChatMessage) => {
  avatarDropdownShow.value = null

  if (key === 'kick' && msg.client_id) {
    dialog.warning({
      title: '踢出用户',
      content: `确定要踢出用户 ${msg.username} 吗？`,
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: async () => {
        try {
          await adminKickUser(msg.client_id!, '违反聊天室规则')
          message.success('已踢出用户')
        } catch (error) {
          message.error('踢出失败')
        }
      }
    })
  }
}

// 插入表情
const insertEmoji = (emoji: string) => {
  messageInput.value += emoji
}

// 发送消息
const handleSendMessage = () => {
  if (!messageInput.value.trim()) {
    return
  }

  if (isChatMutedForUser.value) {
    message.warning('当前已开启全员禁言，只有管理员可发言')
    return
  }

  if (!isConnected.value) {
    message.error('未连接到聊天室')
    return
  }

  ws?.sendMessage(messageInput.value.trim())
  messageInput.value = ''
}

// 确认用户设置
const confirmUserSetup = () => {
  if (!userSetup.value.username.trim()) {
    message.error('请输入昵称')
    return
  }
  showUserSetup.value = false
  connectWebSocket()
}

// 初始化
onMounted(() => {
  fetchChatSettingsData()
  // 如果已登录，直接连接
  if (authStore.isLoggedIn) {
    connectWebSocket()
  } else {
    // 否则显示用户设置对话框
    showUserSetup.value = true
  }
})

// 清理
onUnmounted(() => {
  ws?.close()
})
</script>

<style scoped>
.chat-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  /* 使用 DefaultLayout 注入的 CSS 变量，自动响应 header 显隐 */
  height: calc(100vh - var(--header-height, 72px) - 32px);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.chat-container {
  height: 100%;
  display: flex;
  flex-direction: column;
}

/* 穿透 Naive UI n-card 所有内部包裹层，强制 flex 高度链 */
.chat-container :deep(.n-card) {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.chat-container :deep(.n-card__content) {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
  overflow: hidden;
  padding: 0 !important;
}

/* n-card__content 内部可能还有 .n-scrollbar 包裹层 */
.chat-container :deep(.n-card__content > .n-scrollbar),
.chat-container :deep(.n-card__content > .n-scrollbar > .n-scrollbar-container),
.chat-container
  :deep(.n-card__content > .n-scrollbar > .n-scrollbar-container > .n-scrollbar-content) {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
  height: 100%;
}

.chat-layout {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
  margin-bottom: 16px;
  min-height: 0;
  /* header(72px) + main-padding(32px) + card-header(55px) + card-padding(48px) + chat-input(130px) + margin(16px) + chat-page-padding(40px) */
  max-height: calc(100vh - var(--header-height, 72px) - 32px - 289px);
}

.empty-messages {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.messages-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  /* 重置 global.css 的 max-width:100% 对 flex 容器的干扰 */
  max-width: none;
  width: 100%;
}

.message-item {
  display: flex;
  gap: 12px;
  animation: fadeIn 0.3s ease-in;
  /* 防止 global.css max-width 压缩 flex 子项布局 */
  max-width: none;
  width: 100%;
  box-sizing: border-box;
}

.message-item.own-message {
  flex-direction: row-reverse;
  padding-right: 0;
}

.message-item.own-message .message-content {
  align-items: flex-end;
}

.message-item.own-message .message-text {
  background: #18a058;
  color: white;
}

.message-avatar {
  flex-shrink: 0;
  max-width: none;
}

.message-avatar :deep(.n-avatar) {
  max-width: none;
}

.message-avatar :deep(.n-avatar img),
.message-avatar :deep(.n-avatar svg) {
  max-width: none;
  width: 100%;
  height: 100%;
}

.message-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
  max-width: 70%;
  min-width: 0; /* 防止 flex 子项撑破父容器 */
}

.message-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
}

.message-username {
  font-weight: 500;
  color: #333;
}

.message-time {
  color: #999;
}

.message-text {
  padding: 8px 12px;
  background: white;
  border-radius: 8px;
  word-wrap: break-word;
  white-space: pre-wrap;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.chat-input {
  padding: 16px;
  background: white;
  border-radius: 8px;
  border: 1px solid #e0e0e0;
  flex-shrink: 0;
  overflow: hidden;
  box-sizing: border-box;
  width: 100%;
}

.input-wrapper {
  position: relative;
  display: flex;
  gap: 8px;
  align-items: flex-end;
}

.input-wrapper :deep(.n-input) {
  flex: 1;
}

.emoji-btn {
  padding: 6px 10px;
  cursor: pointer;
  transition: all 0.2s;
  color: #666;
}

.emoji-btn:hover {
  color: #18a058;
  transform: scale(1.15);
}

.emoji-picker {
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: 8px;
  max-width: 320px;
  max-height: 300px;
  overflow-y: auto;
  padding: 8px;
}

.emoji-item {
  font-size: 24px;
  cursor: pointer;
  text-align: center;
  padding: 4px;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.emoji-item:hover {
  background-color: #f5f5f5;
  transform: scale(1.2);
}

/* ==== 暗色模式适配 ==== */
html.dark .chat-page {
  color: #e5e7eb;
}

html.dark .chat-container {
  border-color: #243447;
}

html.dark .chat-container :deep(.n-card) {
  background-color: #1a2332 !important;
  border-color: #243447 !important;
  --n-border-color: #243447 !important;
  --n-color: #1a2332 !important;
}

html.dark .chat-container :deep(.n-card-header) {
  background-color: #1a2332 !important;
  border-bottom-color: #243447 !important;
  --n-border-color: #243447 !important;
}

html.dark .chat-container :deep(.n-card-header__main) {
  color: #e2e8f0;
}

html.dark .chat-container :deep(.n-card__content) {
  background-color: #1a2332 !important;
  padding: 0 !important;
}

/* 消除卡片底部白色 padding */
html.dark .chat-container :deep(.n-card__footer) {
  background-color: #1a2332 !important;
  border-top-color: #243447;
}

html.dark .chat-messages {
  background: #0f172a;
  border: 1px solid #1e2d3d;
}

html.dark .message-username {
  color: #e2e8f0;
}

html.dark .message-time {
  color: #94a3b8;
}

html.dark .message-text {
  background: #1a2332;
  color: #e5e7eb;
  box-shadow: none;
  border: 1px solid #243447;
}

html.dark .message-item.own-message .message-text {
  background: #0ea5e9;
  color: #f8fafc;
  border-color: #0ea5e9;
}

html.dark .chat-input {
  background: #1a2332;
  border-color: #243447;
  border: 1px solid #243447;
  border-radius: 8px;
}

html.dark .emoji-btn {
  color: #cbd5e1;
}

html.dark .emoji-btn:hover {
  color: #38bdf8;
}

html.dark .emoji-item:hover {
  background-color: #1f2937;
}

html.dark .chat-messages::-webkit-scrollbar-track {
  background: #111827 !important;
}

html.dark .chat-messages::-webkit-scrollbar-thumb {
  background: #374151 !important;
}

html.dark .chat-messages::-webkit-scrollbar-thumb:hover {
  background: #4b5563 !important;
}

html.dark .input-wrapper :deep(.n-input) {
  background-color: #0f172a;
  border-color: #243447;
}

html.dark .input-wrapper :deep(.n-input__textarea-el) {
  color: #e5e7eb;
  background: transparent;
}

html.dark .input-wrapper :deep(.n-input__textarea-el::placeholder) {
  color: #94a3b8;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 滚动条样式（!important 覆盖 global.css 的 10px 全局设定） */
.chat-messages::-webkit-scrollbar {
  width: 6px !important;
}

.chat-messages::-webkit-scrollbar-track {
  background: #f1f1f1 !important;
  border-radius: 3px !important;
}

.chat-messages::-webkit-scrollbar-thumb {
  background: #888 !important;
  border-radius: 3px !important;
  border: none !important;
}

.chat-messages::-webkit-scrollbar-thumb:hover {
  background: #555 !important;
}

/* ===== 移动端响应式适配 ===== */
@media screen and (max-width: 768px) {
  .chat-page {
    padding: 0;
    max-width: 100%;
    height: calc(100vh - var(--header-height, 60px));
  }

  .chat-container {
    border-radius: 0;
    height: 100%;
  }

  .chat-container {
    border-radius: 0;
  }

  .chat-container :deep(.n-card__content) {
    padding: 12px;
    height: calc(100vh - 60px);
  }

  .chat-messages {
    padding: 12px;
    border-radius: 0;
    margin-bottom: 12px;
  }

  .message-item {
    gap: 8px;
  }

  .message-avatar :deep(.n-avatar) {
    width: 32px !important;
    height: 32px !important;
    font-size: 14px;
  }

  .message-content {
    max-width: 75%;
  }

  .message-header {
    font-size: 11px;
  }

  .message-text {
    padding: 8px 10px;
    font-size: 14px;
  }

  .chat-input {
    padding: 12px;
    border-radius: 0;
  }

  .input-wrapper :deep(.n-input) {
    font-size: 14px;
  }

  .emoji-btn {
    padding: 4px 8px;
  }

  .emoji-btn :deep(.n-icon) {
    font-size: 24px !important;
  }

  .emoji-picker {
    grid-template-columns: repeat(6, 1fr);
    max-width: 280px;
  }

  .emoji-item {
    font-size: 20px;
  }

  /* 在线人数徽章 */
  .chat-container :deep(.n-card-header) {
    padding: 12px;
    font-size: 16px;
  }

  .chat-container :deep(.n-card-header__extra) {
    font-size: 12px;
  }
}

/* 小屏手机优化 (iPhone SE 等) */
@media screen and (max-width: 375px) {
  .message-avatar :deep(.n-avatar) {
    width: 28px !important;
    height: 28px !important;
    font-size: 12px;
  }

  .message-content {
    max-width: 70%;
  }

  .message-text {
    padding: 6px 8px;
    font-size: 13px;
  }

  .emoji-picker {
    grid-template-columns: repeat(5, 1fr);
    max-width: 240px;
  }
}
</style>
