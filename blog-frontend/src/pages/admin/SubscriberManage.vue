<!--
 * @ProjectName: go-vue3-blog
 * @FileName: SubscriberManage.vue
 * @CreateTime: 2026-04-15 10:00:00
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 订阅管理页面组件，提供订阅者管理和 RSS 配置功能
 -->
<template>
  <div class="subscriber-manage-page">
    <div class="page-header">
      <h1 class="page-title">订阅管理</h1>
    </div>

    <!-- Tab 切换 -->
    <n-tabs v-model:value="activeTab" type="line" animated>
      <!-- 订阅者管理 Tab -->
      <n-tab-pane name="subscribers" tab="订阅者管理">
        <div class="tab-header">
          <!-- 视图切换按钮（仅桌面端显示） -->
          <n-button-group v-if="!isMobile" size="small" class="view-toggle-group">
            <n-button :type="viewMode === 'table' ? 'primary' : 'default'" @click="viewMode = 'table'">
              <template #icon>
                <n-icon :component="GridOutline" />
              </template>
              表格
            </n-button>
            <n-button :type="viewMode === 'card' ? 'primary' : 'default'" @click="viewMode = 'card'">
              <template #icon>
                <n-icon :component="AppsOutline" />
              </template>
              卡片
            </n-button>
          </n-button-group>
        </div>

        <!-- 订阅者列表内容区域 -->
        <div class="content-area">
      <div v-if="isMobile || viewMode === 'card'" class="card-list">
        <n-card v-for="subscriber in subscribers" :key="subscriber.id" class="list-card" :size="isMobile ? 'small' : 'medium'">
          <template #header>
            <div class="card-header-content">
              <span class="email">{{ subscriber.email }}</span>
              <n-tag :type="subscriber.is_active ? 'success' : 'default'" :size="isMobile ? 'tiny' : 'small'">
                {{ subscriber.is_active ? '活跃' : '已退订' }}
              </n-tag>
            </div>
          </template>
          <div class="card-content">
            <div class="info-item">
              <span class="label">订阅时间：</span>
              <span class="value">{{ formatDate(subscriber.subscribed_at, 'YYYY-MM-DD HH:mm') }}</span>
            </div>
            <div v-if="!subscriber.is_active && subscriber.unsubscribed_at" class="info-item">
              <span class="label">退订时间：</span>
              <span class="value">{{ formatDate(subscriber.unsubscribed_at, 'YYYY-MM-DD HH:mm') }}</span>
            </div>
          </div>
          <template #footer>
            <n-space justify="end" :size="isMobile ? 'small' : 'medium'">
              <n-button :size="isMobile ? 'tiny' : 'small'" type="error" @click="handleDelete(subscriber.id)">
                删除
              </n-button>
            </n-space>
          </template>
        </n-card>
      </div>

      <n-data-table
        v-else-if="viewMode === 'table'"
        :columns="columns"
        :data="subscribers"
        :loading="loading"
        :single-line="false"
      />

      <!-- 分页 - 位于表格右下角 -->
      <div class="pagination-wrapper">
        <n-pagination
          v-if="total > 0"
          v-model:page="currentPage"
          :page-count="totalPages"
          :page-size="pageSize"
          :page-slot="isMobile ? 3 : 7"
          :simple="isMobile"
          @update:page="handlePageChange"
        />
      </div>
        </div>
      </n-tab-pane>

      <!-- 推送历史 Tab -->
      <n-tab-pane name="push-history" tab="推送历史">
        <PushHistoryTab />
      </n-tab-pane>

      <!-- RSS 配置 Tab -->
      <n-tab-pane name="rss" tab="RSS 配置">
        <div class="rss-config-content">
          <!-- RSS 配置表单 -->
          <n-card title="RSS 订阅配置" class="config-card">
            <n-form
              ref="formRef"
              :model="rssFormData"
              :label-placement="isMobile ? 'top' : 'left'"
              :label-width="isMobile ? 'auto' : '120'"
              require-mark-placement="right-hanging"
            >
              <n-form-item label="启用 RSS" path="enabled">
                <n-switch v-model:value="rssFormData.enabled" />
                <template #feedback>
                  <span style="font-size: 12px; color: #999;">
                    关闭后，所有 RSS Feed 接口将返回"功能未启用"错误
                  </span>
                </template>
              </n-form-item>

              <n-form-item label="RSS 标题" path="title">
                <n-input
                  v-model:value="rssFormData.title"
                  placeholder="例如：我的博客"
                  maxlength="100"
                  show-count
                  clearable
                />
              </n-form-item>

              <n-form-item label="RSS 描述" path="description">
                <n-input
                  v-model:value="rssFormData.description"
                  type="textarea"
                  placeholder="例如：分享技术与生活"
                  :autosize="{ minRows: 2, maxRows: 4 }"
                  maxlength="500"
                  show-count
                  clearable
                />
              </n-form-item>

              <n-form-item label="作者名称" path="author_name">
                <n-input
                  v-model:value="rssFormData.author_name"
                  placeholder="例如：張三"
                  maxlength="50"
                  clearable
                />
              </n-form-item>

              <n-form-item label="作者邮箱" path="author_email">
                <n-input
                  v-model:value="rssFormData.author_email"
                  placeholder="例如：author@example.com"
                  maxlength="100"
                  clearable
                />
              </n-form-item>

              <n-form-item label="网站链接" path="link">
                <n-input
                  v-model:value="rssFormData.link"
                  placeholder="例如：https://blog.example.com"
                  maxlength="200"
                  clearable
                />
                <template #feedback>
                  <span style="font-size: 12px; color: #999;">
                    RSS Feed 中的网站主页链接
                  </span>
                </template>
              </n-form-item>

              <n-form-item label="语言" path="language">
                <n-input
                  v-model:value="rssFormData.language"
                  placeholder="例如：zh-CN"
                  maxlength="10"
                  clearable
                />
                <template #feedback>
                  <span style="font-size: 12px; color: #999;">
                    RSS Feed 语言代码，如 zh-CN（简体中文）、en-US（英语）
                  </span>
                </template>
              </n-form-item>

              <n-form-item label="版权信息" path="copyright">
                <n-input
                  v-model:value="rssFormData.copyright"
                  placeholder="例如：Copyright © 2024 我的博客"
                  maxlength="200"
                  clearable
                />
              </n-form-item>

              <n-form-item label="文章数量限制" path="item_limit">
                <n-input-number
                  v-model:value="rssFormData.item_limit"
                  :min="5"
                  :max="100"
                  style="width: 100%"
                >
                  <template #suffix>篇</template>
                </n-input-number>
                <template #feedback>
                  <span style="font-size: 12px; color: #999;">
                    每个 RSS Feed 包含的最大文章数量，建议 10-30 篇
                  </span>
                </template>
              </n-form-item>

              <n-form-item label="缓存时长" path="cache_duration">
                <n-input-number
                  v-model:value="rssFormData.cache_duration"
                  :min="1"
                  :max="1440"
                  style="width: 100%"
                >
                  <template #suffix>分钟</template>
                </n-input-number>
                <template #feedback>
                  <span style="font-size: 12px; color: #999;">
                    RSS Feed 缓存时长，建议 30-60 分钟
                  </span>
                </template>
              </n-form-item>

              <n-form-item>
                <n-space>
                  <n-button type="primary" @click="handleRSSSubmit" :loading="rssLoading">
                    保存配置
                  </n-button>
                  <n-button @click="handleRSSReset">
                    重置
                  </n-button>
                </n-space>
              </n-form-item>
            </n-form>
          </n-card>

          <!-- RSS 订阅地址 -->
          <n-card title="RSS 订阅地址" class="config-card" style="margin-top: 16px;">
            <n-alert type="success">
              <template #header>公开订阅地址</template>
              <ul style="margin: 8px 0; padding-left: 20px; line-height: 1.8;">
                <li><strong>全站 RSS：</strong> {{ baseUrl }}/api/feed.xml</li>
                <li><strong>文章 RSS：</strong> {{ baseUrl }}/api/rss/posts.xml</li>
                <li><strong>说说 RSS：</strong> {{ baseUrl }}/api/rss/moments.xml</li>
                <li><strong>分类 RSS：</strong> {{ baseUrl }}/api/rss/category/{分类ID}.xml</li>
                <li><strong>标签 RSS：</strong> {{ baseUrl }}/api/rss/tag/{标签ID}.xml</li>
              </ul>
            </n-alert>

            <n-space style="margin-top: 16px;">
              <n-button @click="copyFeedUrl(`${baseUrl}/api/feed.xml`)">
                复制全站 RSS 地址
              </n-button>
              <n-button @click="openFeedUrl(`${baseUrl}/api/feed.xml`)">
                在新窗口打开
              </n-button>
              <n-button @click="handleClearRSSCache" :loading="clearCacheLoading">
                清除 RSS 缓存
              </n-button>
            </n-space>
          </n-card>
        </div>
      </n-tab-pane>
    </n-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, h, watch, reactive } from 'vue'
import { useMessage, useDialog, NButton, NButtonGroup, NIcon, NTag, NSpace } from 'naive-ui'
import type { DataTableColumns, FormInst } from 'naive-ui'
import { GridOutline, AppsOutline } from '@vicons/ionicons5'
import request from '@/utils/request'
import { formatDate } from '@/utils/format'
import { getRSSConfig, updateRSSConfig, clearRSSCache } from '@/api/rss'
import type { RSSConfig } from '@/api/rss'
import PushHistoryTab from './components/PushHistoryTab.vue'

interface Subscriber {
  id: number
  email: string
  is_active: boolean
  subscribed_at: string
  unsubscribed_at?: string
}

const message = useMessage()
const dialog = useDialog()

// 订阅者管理相关状态
const loading = ref(false)
const subscribers = ref<Subscriber[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const viewMode = ref<'table' | 'card'>('table')
const isMobile = ref(false)
const activeTab = ref('subscribers')

// RSS 配置相关状态
const formRef = ref<FormInst | null>(null)
const rssLoading = ref(false)
const clearCacheLoading = ref(false)
const baseUrl = window.location.origin

const rssFormData = reactive<RSSConfig>({
  enabled: true,
  title: '',
  description: '',
  link: '',
  author_name: '',
  author_email: '',
  language: 'zh-CN',
  copyright: '',
  item_limit: 20,
  cache_duration: 60
})

const originalRSSData = ref<RSSConfig>({})

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

// 获取订阅者列表
const getSubscribers = async (page = 1) => {
  try {
    loading.value = true
    const response = await request({
      url: '/admin/subscribers',
      method: 'get',
      params: {
        page,
        page_size: pageSize.value
      }
    })
    subscribers.value = response.data.list || []
    total.value = response.data.total || 0
  } catch (error: any) {
    message.error(error?.message || '获取订阅者列表失败')
  } finally {
    loading.value = false
  }
}

// 删除订阅者
const deleteSubscriber = async (id: number) => {
  try {
    await request({
      url: `/admin/subscribers/${id}`,
      method: 'delete'
    })
    message.success('删除成功')
    await getSubscribers(currentPage.value)
  } catch (error: any) {
    message.error(error?.message || '删除失败')
  }
}

// 处理删除
const handleDelete = (id: number) => {
  dialog.warning({
    title: '确认删除',
    content: '确定要删除这个订阅者吗？此操作不可恢复。',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: () => {
      deleteSubscriber(id)
    }
  })
}

// 处理分页变化
const handlePageChange = (page: number) => {
  currentPage.value = page
  getSubscribers(page)
}

// 检测屏幕尺寸
const checkMobile = () => {
  isMobile.value = window.innerWidth < 768
}

// RSS 配置相关方法
const fetchRSSConfig = async () => {
  try {
    rssLoading.value = true
    const res = await getRSSConfig()
    if (res.data) {
      Object.assign(rssFormData, {
        enabled: res.data.enabled ?? true,
        title: res.data.title || '',
        description: res.data.description || '',
        link: res.data.link || '',
        author_name: res.data.author_name || '',
        author_email: res.data.author_email || '',
        language: res.data.language || 'zh-CN',
        copyright: res.data.copyright || '',
        item_limit: res.data.item_limit || 20,
        cache_duration: res.data.cache_duration || 60
      })
      originalRSSData.value = { ...rssFormData }
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '获取 RSS 配置失败')
  } finally {
    rssLoading.value = false
  }
}

const handleRSSSubmit = async () => {
  rssLoading.value = true
  try {
    await updateRSSConfig(rssFormData)
    message.success('RSS 配置保存成功')
    originalRSSData.value = { ...rssFormData }
    message.info('RSS 缓存已自动清除，新配置立即生效')
  } catch (error: any) {
    message.error(error.response?.data?.message || '保存失败')
  } finally {
    rssLoading.value = false
  }
}

const handleRSSReset = () => {
  Object.assign(rssFormData, originalRSSData.value)
  formRef.value?.restoreValidation()
  message.info('已重置为上次保存的数据')
}

const handleClearRSSCache = async () => {
  try {
    clearCacheLoading.value = true
    await clearRSSCache()
    message.success('RSS 缓存已清除')
  } catch (error: any) {
    message.error(error.response?.data?.message || '清除缓存失败')
  } finally {
    clearCacheLoading.value = false
  }
}

const copyFeedUrl = (url: string) => {
  navigator.clipboard.writeText(url).then(() => {
    message.success('RSS 地址已复制到剪贴板')
  }).catch(() => {
    message.error('复制失败，请手动复制')
  })
}

const openFeedUrl = (url: string) => {
  window.open(url, '_blank')
}

// 表格列定义
const columns = computed<DataTableColumns>(() => [
  {
    title: 'ID',
    key: 'id',
    width: 80
  },
  {
    title: '邮箱',
    key: 'email',
    minWidth: 200
  },
  {
    title: '状态',
    key: 'is_active',
    width: 100,
    render: (row: any) => {
      return h(
        NTag,
        {
          type: row.is_active ? 'success' : 'default',
          size: 'small'
        },
        { default: () => (row.is_active ? '活跃' : '已退订') }
      )
    }
  },
  {
    title: '订阅时间',
    key: 'subscribed_at',
    width: 180,
    render: (row: any) => formatDate(row.subscribed_at, 'YYYY-MM-DD HH:mm')
  },
  {
    title: '退订时间',
    key: 'unsubscribed_at',
    width: 180,
    render: (row: any) => row.unsubscribed_at ? formatDate(row.unsubscribed_at, 'YYYY-MM-DD HH:mm') : '-'
  },
  {
    title: '操作',
    key: 'actions',
    width: 120,
    render: (row: any) => {
      return h(
        NSpace,
        { size: 'small' },
        {
          default: () => [
            h(
              NButton,
              {
                size: 'small',
                type: 'error',
                onClick: () => handleDelete(row.id)
              },
              { default: () => '删除' }
            )
          ]
        }
      )
    }
  }
])

// 监听视图模式变化，保存到本地存储
watch(viewMode, (newMode) => {
  localStorage.setItem('subscriber-view-mode', newMode)
})

// 监听 Tab 切换，加载对应数据
watch(activeTab, (newTab) => {
  if (newTab === 'rss') {
    fetchRSSConfig()
  }
})

onMounted(() => {
  getSubscribers()
  checkMobile()
  window.addEventListener('resize', checkMobile)

  // 从本地存储恢复视图模式
  const savedMode = localStorage.getItem('subscriber-view-mode')
  if (savedMode === 'card' || savedMode === 'table') {
    viewMode.value = savedMode
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped lang="scss">
.subscriber-manage-page {
  padding: 20px;

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    .page-title {
      font-size: 24px;
      font-weight: 600;
      margin: 0;
    }
  }

  .tab-header {
    display: flex;
    justify-content: flex-end;
    margin-bottom: 16px;

    .view-toggle-group {
      display: flex;
      gap: 0;
    }
  }

  .content-area {
    background: var(--card-color);
    border-radius: 8px;
    padding: 20px;

    .card-list {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
      gap: 16px;

      .list-card {
        .card-header-content {
          display: flex;
          justify-content: space-between;
          align-items: center;

          .email {
            font-weight: 500;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }
        }

        .card-content {
          .info-item {
            margin-bottom: 8px;
            font-size: 14px;

            .label {
              color: var(--text-color-3);
            }

            .value {
              color: var(--text-color-2);
            }
          }
        }
      }
    }

    .pagination-wrapper {
      display: flex;
      justify-content: flex-end;
      margin-top: 20px;
    }
  }

  .rss-config-content {
    .config-card {
      background: rgba(255, 255, 255, 0.85);
      backdrop-filter: blur(20px) saturate(180%);
      border: 1px solid rgba(8, 145, 178, 0.1);
      box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
      transition: all 0.3s;

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 12px 40px rgba(0, 0, 0, 0.12);
      }
    }
  }
}

html.dark .subscriber-manage-page {
  .rss-config-content {
    .config-card {
      background: rgba(30, 41, 59, 0.85);
      border: 1px solid rgba(255, 255, 255, 0.08);
      box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
    }
  }
}

@media (max-width: 768px) {
  .subscriber-manage-page {
    padding: 12px;

    .page-header {
      .page-title {
        font-size: 20px;
      }
    }

    .content-area {
      padding: 12px;

      .card-list {
        grid-template-columns: 1fr;
        gap: 12px;
      }
    }
  }
}
</style>
