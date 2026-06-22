<!--
 * @ProjectName: go-vue3-blog
 * @FileName: Subscribe.vue
 * @CreateTime: 2026-04-18 14:42:07
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 邮件订阅页面组件，提供订阅和退订功能
 -->
<template>
  <div class="subscribe-page">
    <div class="subscribe-layout">
      <div class="subscribe-main">
        <n-spin :show="loading">
          <!-- 页面头部 -->
          <div class="page-header">
            <h1 class="page-title">📬 订阅本站</h1>
            <p class="page-description">选择您喜欢的订阅方式，第一时间获取最新文章推送</p>
          </div>

          <!-- 退订成功状态 -->
          <n-card v-if="unsubscribeSuccess" class="result-card">
            <n-result
              status="success"
              title="退订成功"
              description="您已成功取消订阅，不会再收到新文章推送邮件"
            >
              <template #footer>
                <n-button type="primary" @click="resetForm">返回</n-button>
              </template>
            </n-result>
          </n-card>

          <!-- 订阅内容区 -->
          <template v-else>
            <!-- 订阅方式选择卡片 -->
            <n-card class="intro-card">
              <div class="intro-content">
                <div class="intro-icon">✉️</div>
                <h2 class="intro-title">选择订阅方式</h2>
                <p class="intro-subtitle">两种方式任选其一，都能第一时间获取更新</p>

                <div class="subscribe-methods">
                  <!-- 邮件订阅选项 -->
                  <div
                    class="method-card"
                    :class="{ active: subscribeMethod === 'email' }"
                    @click="subscribeMethod = 'email'"
                  >
                    <div class="method-icon">📧</div>
                    <h3 class="method-title">邮件订阅</h3>
                    <p class="method-desc">新文章发布时自动推送到您的邮箱</p>
                    <div class="method-features">
                      <div class="feature-tag">✓ 自动推送</div>
                      <div class="feature-tag">✓ 随时退订</div>
                    </div>
                  </div>

                  <!-- RSS 订阅选项 -->
                  <div
                    class="method-card"
                    :class="{ active: subscribeMethod === 'rss', disabled: !rssEnabled }"
                    @click="handleRssMethodClick"
                  >
                    <div class="method-icon">📡</div>
                    <h3 class="method-title">RSS 订阅</h3>
                    <p class="method-desc">
                      {{
                        rssEnabled
                          ? '通过 RSS 阅读器订阅，随时随地获取更新'
                          : 'RSS 订阅功能暂未开放'
                      }}
                    </p>
                    <div class="method-features" v-if="rssEnabled">
                      <div class="feature-tag">✓ 无需邮箱</div>
                      <div class="feature-tag">✓ 隐私保护</div>
                    </div>
                    <div v-else class="method-disabled-tag">🔒 暂未开放</div>
                  </div>
                </div>
              </div>
            </n-card>

            <!-- 订阅表单卡片 -->
            <n-card class="subscribe-form-card">
              <!-- 邮件订阅成功 -->
              <n-result
                v-if="subscribeSuccess && subscribeMethod === 'email'"
                status="success"
                title="订阅成功"
                description="欢迎邮件已发送到您的邮箱，请查收确认"
              >
                <template #footer>
                  <n-button type="primary" @click="resetForm">继续浏览</n-button>
                </template>
              </n-result>

              <!-- 邮件订阅表单 -->
              <template v-else-if="subscribeMethod === 'email'">
                <h2 class="form-title">📧 邮件订阅</h2>
                <n-form ref="formRef" :model="formData" :rules="rules" label-placement="top">
                  <n-form-item label="邮箱地址" path="email">
                    <n-input
                      v-model:value="formData.email"
                      placeholder="请输入您的邮箱地址"
                      size="large"
                      :disabled="loading"
                      @keyup.enter="handleSubscribe"
                    >
                      <template #prefix>
                        <n-icon :component="MailOutline" />
                      </template>
                    </n-input>
                  </n-form-item>

                  <n-button
                    type="primary"
                    size="large"
                    block
                    :loading="loading"
                    @click="handleSubscribe"
                    class="subscribe-button"
                  >
                    <template #icon>
                      <n-icon :component="SendOutline" />
                    </template>
                    立即订阅
                  </n-button>
                </n-form>

                <div class="tips-section">
                  <n-alert type="info" :bordered="false">
                    <template #icon>
                      <n-icon :component="InformationCircleOutline" />
                    </template>
                    订阅后，您将在新文章发布时收到邮件通知。如需取消订阅，可点击邮件中的退订链接。
                  </n-alert>
                </div>
              </template>

              <!-- RSS 订阅内容 -->
              <template v-else-if="subscribeMethod === 'rss'">
                <!-- RSS 功能未启用提示 -->
                <n-result
                  v-if="!rssEnabled"
                  status="warning"
                  title="RSS 订阅暂未开放"
                  description="管理员暂未开启 RSS 订阅功能，请使用邮件订阅或稍后再试"
                >
                  <template #footer>
                    <n-button type="primary" @click="subscribeMethod = 'email'">
                      切换到邮件订阅
                    </n-button>
                  </template>
                </n-result>

                <!-- RSS 订阅链接 -->
                <template v-else>
                  <h2 class="form-title">📡 RSS 订阅</h2>
                  <p class="rss-intro">通过 RSS 阅读器订阅本站，随时随地获取更新，无需提供邮箱</p>

                  <div class="rss-links">
                    <div class="rss-link-item">
                      <div class="link-info">
                        <div class="link-icon">🌐</div>
                        <div class="link-text">
                          <div class="link-title">全站订阅</div>
                          <div class="link-desc">订阅所有文章和说说</div>
                        </div>
                      </div>
                      <n-button
                        text
                        type="primary"
                        tag="a"
                        :href="`${baseUrl}/api/feed.xml`"
                        target="_blank"
                        class="rss-button"
                      >
                        <template #icon>
                          <n-icon :component="LinkOutline" />
                        </template>
                        订阅
                      </n-button>
                    </div>

                    <div class="rss-link-item">
                      <div class="link-info">
                        <div class="link-icon">📝</div>
                        <div class="link-text">
                          <div class="link-title">文章订阅</div>
                          <div class="link-desc">仅订阅博客文章</div>
                        </div>
                      </div>
                      <n-button
                        text
                        type="primary"
                        tag="a"
                        :href="`${baseUrl}/api/rss/posts.xml`"
                        target="_blank"
                        class="rss-button"
                      >
                        <template #icon>
                          <n-icon :component="LinkOutline" />
                        </template>
                        订阅
                      </n-button>
                    </div>

                    <div class="rss-link-item">
                      <div class="link-info">
                        <div class="link-icon">💭</div>
                        <div class="link-text">
                          <div class="link-title">说说订阅</div>
                          <div class="link-desc">仅订阅动态说说</div>
                        </div>
                      </div>
                      <n-button
                        text
                        type="primary"
                        tag="a"
                        :href="`${baseUrl}/api/rss/moments.xml`"
                        target="_blank"
                        class="rss-button"
                      >
                        <template #icon>
                          <n-icon :component="LinkOutline" />
                        </template>
                        订阅
                      </n-button>
                    </div>
                  </div>

                  <div class="tips-section">
                    <n-alert type="info" :bordered="false">
                      <template #icon>
                        <n-icon :component="InformationCircleOutline" />
                      </template>
                      推荐使用 Feedly、Inoreader、NetNewsWire 等 RSS 阅读器订阅本站
                    </n-alert>
                  </div>
                </template>
              </template>
            </n-card>

            <!-- 统计信息卡片 -->
            <n-card class="stats-card">
              <div class="stats-content">
                <div class="stat-item">
                  <div class="stat-icon">👥</div>
                  <div class="stat-info">
                    <div class="stat-value">{{ subscriberCount }}</div>
                    <div class="stat-label">订阅用户</div>
                  </div>
                </div>
                <div class="stat-divider"></div>
                <div class="stat-item">
                  <div class="stat-icon">📝</div>
                  <div class="stat-info">
                    <div class="stat-value">{{ postCount }}</div>
                    <div class="stat-label">文章总数</div>
                  </div>
                </div>
                <div class="stat-divider"></div>
                <div class="stat-item">
                  <div class="stat-icon">🔥</div>
                  <div class="stat-info">
                    <div class="stat-value">定期更新</div>
                    <div class="stat-label">持续输出</div>
                  </div>
                </div>
              </div>
            </n-card>
          </template>
        </n-spin>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useMessage, NIcon } from 'naive-ui'
import { MailOutline, SendOutline, InformationCircleOutline, LinkOutline } from '@vicons/ionicons5'
import request from '@/utils/request'
import { getAuthorProfile } from '@/api/blog'
import { getRSSStatus } from '@/api/rss'
import { getSubscribeStats } from '@/api/subscribe'

const route = useRoute()
const message = useMessage()

const formRef = ref(null)
const loading = ref(false)
const subscribeSuccess = ref(false)
const unsubscribeSuccess = ref(false)
const subscriberCount = ref('0')
const postCount = ref('0')
const baseUrl = ref(window.location.origin)
const subscribeMethod = ref('email') // 默认选择邮件订阅
const rssEnabled = ref(true) // RSS 是否启用

const formData = reactive({
  email: ''
})

const rules = {
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: ['blur', 'change'] }
  ]
}

// 订阅邮件
const subscribeEmail = async email => {
  return request({
    url: '/subscribe',
    method: 'post',
    data: { email }
  })
}

// 退订邮件
const unsubscribeEmail = async token => {
  return request({
    url: '/subscribe/unsubscribe',
    method: 'get',
    params: { token }
  })
}

// 获取 RSS 启用状态
const fetchRSSStatus = async () => {
  try {
    const res = await getRSSStatus()
    if (res.data) {
      rssEnabled.value = res.data.enabled
      // 如果 RSS 未启用且当前选择的是 RSS，自动切换到邮件订阅
      if (!rssEnabled.value && subscribeMethod.value === 'rss') {
        subscribeMethod.value = 'email'
      }
    }
  } catch (error) {
    console.error('获取 RSS 状态失败:', error)
    // 失败时默认启用
    rssEnabled.value = true
  }
}

// 获取统计信息
const fetchStats = async () => {
  try {
    // 获取博主资料和统计数据
    const res = await getAuthorProfile()
    if (res.data) {
      postCount.value = res.data.stats.posts.toString()
    }

    // 获取累积订阅者总数（包括已退订的用户）
    const statsRes = await getSubscribeStats()
    if (statsRes.data) {
      const count = statsRes.data.total_count
      subscriberCount.value = count > 0 ? count.toString() : '0'
    }
  } catch (error) {
    console.error('获取统计信息失败:', error)
    // 失败时使用默认值
    subscriberCount.value = '0'
    postCount.value = '0'
  }
}

// 处理 RSS 方式点击
const handleRssMethodClick = () => {
  if (!rssEnabled.value) {
    message.warning('RSS 订阅功能暂未开放，请使用邮件订阅')
    return
  }
  subscribeMethod.value = 'rss'
}

// 处理订阅
const handleSubscribe = async () => {
  try {
    await formRef.value?.validate()
    loading.value = true

    await subscribeEmail(formData.email)
    subscribeSuccess.value = true
    message.success('订阅成功！欢迎邮件已发送到您的邮箱')
  } catch (error) {
    if (error?.message) {
      message.error(error.message)
    }
  } finally {
    loading.value = false
  }
}

// 处理退订
const handleUnsubscribe = async token => {
  try {
    loading.value = true
    await unsubscribeEmail(token)
    unsubscribeSuccess.value = true
    message.success('退订成功')
  } catch (error) {
    message.error(error?.message || '退订失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

// 重置表单
const resetForm = () => {
  subscribeSuccess.value = false
  unsubscribeSuccess.value = false
  formData.email = ''
  formRef.value?.restoreValidation()
}

// 页面加载时检查是否是退订操作
onMounted(() => {
  const action = route.query.action
  const token = route.query.token

  if (action === 'unsubscribe' && token) {
    handleUnsubscribe(token)
  }

  fetchRSSStatus()
  fetchStats()
})
</script>

<style scoped>
.subscribe-page {
  min-height: calc(100vh - 180px);
  padding: 32px 0;
}

.subscribe-layout {
  max-width: 900px;
  margin: 0 auto;
  padding: 0 20px;
}

.subscribe-main {
  min-width: 0;
}

.page-header {
  text-align: center;
  margin-bottom: 32px;
}

.page-title {
  font-size: 32px;
  font-weight: 700;
  margin: 0 0 12px 0;
  color: #fff;
  text-shadow: 0 2px 12px rgba(0, 0, 0, 0.5);
}

html.dark .page-title {
  color: #fff;
}

.page-description {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.8);
  text-shadow: 0 1px 8px rgba(0, 0, 0, 0.4);
  margin: 0;
}

html.dark .page-description {
  color: rgba(255, 255, 255, 0.7);
}

.result-card {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  border: 1px solid rgba(8, 145, 178, 0.1);
  padding: 40px 0;
}

html.dark .result-card {
  background: rgba(30, 41, 59, 0.8);
  border-color: rgba(56, 189, 248, 0.1);
}

/* 介绍卡片 */
.intro-card {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  border: 1px solid rgba(8, 145, 178, 0.1);
  margin-bottom: 24px;
}

html.dark .intro-card {
  background: rgba(30, 41, 59, 0.8);
  border-color: rgba(56, 189, 248, 0.1);
}

.intro-content {
  text-align: center;
}

.intro-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.intro-title {
  font-size: 24px;
  font-weight: 700;
  margin: 0 0 8px 0;
  color: #1a202c;
}

html.dark .intro-title {
  color: #e5e5e5;
}

.intro-subtitle {
  font-size: 14px;
  color: #64748b;
  margin: 0 0 32px 0;
}

html.dark .intro-subtitle {
  color: #94a3b8;
}

/* 订阅方式选择 */
.subscribe-methods {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
  margin-top: 24px;
}

.method-card {
  padding: 24px;
  background: rgba(8, 145, 178, 0.05);
  border: 2px solid transparent;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  text-align: center;
}

.method-card:hover {
  background: rgba(8, 145, 178, 0.1);
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(8, 145, 178, 0.15);
}

.method-card.active {
  background: rgba(8, 145, 178, 0.15);
  border-color: #0891b2;
  box-shadow: 0 8px 24px rgba(8, 145, 178, 0.2);
}

html.dark .method-card {
  background: rgba(56, 189, 248, 0.1);
}

html.dark .method-card:hover {
  background: rgba(56, 189, 248, 0.15);
  box-shadow: 0 8px 24px rgba(56, 189, 248, 0.15);
}

html.dark .method-card.active {
  background: rgba(56, 189, 248, 0.2);
  border-color: #38bdf8;
  box-shadow: 0 8px 24px rgba(56, 189, 248, 0.2);
}

.method-card.disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.method-card.disabled:hover {
  transform: none;
  box-shadow: none;
}

.method-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.method-title {
  font-size: 18px;
  font-weight: 600;
  margin: 0 0 8px 0;
  color: #1a202c;
}

html.dark .method-title {
  color: #e5e5e5;
}

.method-desc {
  font-size: 14px;
  color: #64748b;
  margin: 0 0 16px 0;
  line-height: 1.6;
}

html.dark .method-desc {
  color: #94a3b8;
}

.method-features {
  display: flex;
  justify-content: center;
  gap: 8px;
  flex-wrap: wrap;
}

.feature-tag {
  font-size: 12px;
  padding: 4px 12px;
  background: rgba(8, 145, 178, 0.1);
  border-radius: 12px;
  color: #0891b2;
}

html.dark .feature-tag {
  background: rgba(56, 189, 248, 0.15);
  color: #38bdf8;
}

.method-disabled-tag {
  font-size: 12px;
  padding: 4px 12px;
  background: rgba(239, 68, 68, 0.1);
  border-radius: 12px;
  color: #ef4444;
  margin-top: 8px;
}

html.dark .method-disabled-tag {
  background: rgba(239, 68, 68, 0.15);
  color: #f87171;
}

/* 订阅表单卡片 */
.subscribe-form-card {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  border: 1px solid rgba(8, 145, 178, 0.1);
  margin-bottom: 24px;
}

html.dark .subscribe-form-card {
  background: rgba(30, 41, 59, 0.8);
  border-color: rgba(56, 189, 248, 0.1);
}

.form-title {
  font-size: 20px;
  font-weight: 600;
  margin: 0 0 24px 0;
  color: #1a202c;
  text-align: center;
}

html.dark .form-title {
  color: #e5e5e5;
}

.subscribe-button {
  margin-top: 8px;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
}

.tips-section {
  margin-top: 24px;
}

/* RSS 订阅内容样式 */
.rss-intro {
  font-size: 14px;
  color: #64748b;
  text-align: center;
  margin: 0 0 24px 0;
  line-height: 1.6;
}

html.dark .rss-intro {
  color: #94a3b8;
}

.rss-links {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 24px;
}

.rss-link-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: rgba(8, 145, 178, 0.05);
  border-radius: 8px;
  transition: all 0.3s ease;
}

.rss-link-item:hover {
  background: rgba(8, 145, 178, 0.1);
  transform: translateX(4px);
}

html.dark .rss-link-item {
  background: rgba(56, 189, 248, 0.1);
}

html.dark .rss-link-item:hover {
  background: rgba(56, 189, 248, 0.15);
}

.link-info {
  display: flex;
  align-items: center;
  gap: 16px;
  flex: 1;
}

.link-icon {
  font-size: 32px;
  flex-shrink: 0;
}

.link-text {
  flex: 1;
}

.link-title {
  font-size: 16px;
  font-weight: 600;
  color: #1a202c;
  margin-bottom: 4px;
}

html.dark .link-title {
  color: #e5e5e5;
}

.link-desc {
  font-size: 13px;
  color: #64748b;
}

html.dark .link-desc {
  color: #94a3b8;
}

.rss-button {
  font-size: 14px;
  font-weight: 600;
}

/* 统计信息卡片 */
.stats-card {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  border: 1px solid rgba(8, 145, 178, 0.1);
}

html.dark .stats-card {
  background: rgba(30, 41, 59, 0.8);
  border-color: rgba(56, 189, 248, 0.1);
}

.stats-content {
  display: flex;
  justify-content: space-around;
  align-items: center;
  gap: 24px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 16px;
  flex: 1;
  justify-content: center;
}

.stat-icon {
  font-size: 40px;
}

.stat-info {
  text-align: left;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #0891b2;
  line-height: 1.2;
}

html.dark .stat-value {
  color: #38bdf8;
}

.stat-label {
  font-size: 14px;
  color: #64748b;
  margin-top: 4px;
}

html.dark .stat-label {
  color: #94a3b8;
}

.stat-divider {
  width: 1px;
  height: 50px;
  background: rgba(8, 145, 178, 0.2);
}

html.dark .stat-divider {
  background: rgba(56, 189, 248, 0.2);
}

/* 响应式 */
@media (max-width: 768px) {
  .subscribe-page {
    padding: 24px 0;
  }

  .subscribe-layout {
    padding: 0 16px;
  }

  .page-title {
    font-size: 28px;
  }

  .page-description {
    font-size: 14px;
  }

  .subscribe-methods {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .method-card {
    padding: 20px;
  }

  .stats-content {
    flex-direction: column;
    gap: 20px;
  }

  .stat-divider {
    width: 80%;
    height: 1px;
  }

  .stat-item {
    width: 100%;
  }
}
</style>
