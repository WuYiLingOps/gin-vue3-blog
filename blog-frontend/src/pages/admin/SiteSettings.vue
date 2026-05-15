<!--
 * @ProjectName: go-vue3-blog
 * @FileName: SiteSettings.vue
 * @CreateTime: 2026-02-02 11:43:17
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 网站设置页面组件，提供网站基本信息的配置功能
 -->
<template>
  <div class="site-settings-page">
    <n-card title="网站设置">
      <n-form
        ref="formRef"
        :model="formData"
        :label-placement="isMobile ? 'top' : 'left'"
        :label-width="isMobile ? 'auto' : '120'"
        require-mark-placement="right-hanging"
      >
        <n-form-item label="网站名称" path="site_name">
          <n-input
            v-model:value="formData.site_name"
            placeholder="请输入网站名称"
            maxlength="50"
            show-count
            clearable
          >
            <template #suffix>
              <n-button text size="tiny" @click="clearField('site_name')" type="error">
                清空
              </n-button>
            </template>
          </n-input>
        </n-form-item>

        <n-form-item label="ICP备案号" path="site_icp">
          <n-input
            v-model:value="formData.site_icp"
            placeholder="例如：京ICP备12345678号"
            maxlength="50"
            clearable
          >
            <template #suffix>
              <n-button text size="tiny" @click="clearField('site_icp')" type="error">
                清空
              </n-button>
            </template>
          </n-input>
        </n-form-item>

        <n-form-item label="公安备案号" path="site_police">
          <n-input
            v-model:value="formData.site_police"
            placeholder="例如：京公网安备11010802012345号"
            maxlength="50"
            clearable
          >
            <template #suffix>
              <n-button text size="tiny" @click="clearField('site_police')" type="error">
                清空
              </n-button>
            </template>
          </n-input>
        </n-form-item>

        <n-form-item label="网站成立时间" path="site_start_date">
          <n-input
            v-model:value="formData.site_start_date"
            placeholder="例如：2025-10-13 00:00:00"
            maxlength="30"
            clearable
          >
            <template #suffix>
              <n-button text size="tiny" @click="clearField('site_start_date')" type="error">
                清空
              </n-button>
            </template>
          </n-input>
          <template #feedback>
            <span style="font-size: 12px; color: #999;">用于计算网站底部"已运行时间"，格式：YYYY-MM-DD HH:mm:ss，未配置时默认 2025-10-13</span>
          </template>
        </n-form-item>

        <n-divider />

        <n-form-item label="GitHub链接" path="social_github">
          <n-input
            v-model:value="formData.social_github"
            placeholder="例如：https://github.com/username"
            maxlength="200"
            clearable
          >
            <template #suffix>
              <n-button text size="tiny" @click="clearField('social_github')" type="error">
                清空
              </n-button>
            </template>
          </n-input>
        </n-form-item>

        <n-form-item label="Gitee链接" path="social_gitee">
          <n-input
            v-model:value="formData.social_gitee"
            placeholder="例如：https://gitee.com/username"
            maxlength="200"
            clearable
          >
            <template #suffix>
              <n-button text size="tiny" @click="clearField('social_gitee')" type="error">
                清空
              </n-button>
            </template>
          </n-input>
        </n-form-item>

        <n-form-item label="邮箱地址" path="social_email">
          <n-input
            v-model:value="formData.social_email"
            placeholder="例如：example@email.com"
            maxlength="100"
            clearable
          >
            <template #suffix>
              <n-button text size="tiny" @click="clearField('social_email')" type="error">
                清空
              </n-button>
            </template>
          </n-input>
        </n-form-item>

        <n-form-item label="CSDN链接" path="social_csdn">
          <n-input
            v-model:value="formData.social_csdn"
            placeholder="例如：https://blog.csdn.net/username"
            maxlength="200"
            clearable
          >
            <template #suffix>
              <n-button text size="tiny" @click="clearField('social_csdn')" type="error">
                清空
              </n-button>
            </template>
          </n-input>
        </n-form-item>

        <n-form-item label="QQ二维码" path="social_qq">
          <div style="width: 100%;">
            <ImageUpload
              v-model="formData.social_qq"
              :width="280"
              :height="160"
              alt="QQ二维码"
              compact
            />
            <n-input
              v-model:value="formData.social_qq"
              placeholder="QQ二维码图片URL，也可通过上方直接上传"
              maxlength="500"
              clearable
              style="margin-top: 8px;"
            >
              <template #suffix>
                <n-button text size="tiny" @click="clearField('social_qq')" type="error">
                  清空
                </n-button>
              </template>
            </n-input>
          </div>
        </n-form-item>

        <n-form-item label="微信二维码" path="social_wechat">
          <div style="width: 100%;">
            <ImageUpload
              v-model="formData.social_wechat"
              :width="280"
              :height="160"
              alt="微信二维码"
              compact
            />
            <n-input
              v-model:value="formData.social_wechat"
              placeholder="微信二维码图片URL，也可通过上方直接上传"
              maxlength="500"
              clearable
              style="margin-top: 8px;"
            >
              <template #suffix>
                <n-button text size="tiny" @click="clearField('social_wechat')" type="error">
                  清空
                </n-button>
              </template>
            </n-input>
          </div>
        </n-form-item>

        <n-divider />

        <n-form-item label="社交链接排序">
          <div class="social-link-sort-container">
            <div class="sort-tip">
              <n-text depth="3" style="font-size: 12px;">
                拖拽下方项目可调整显示顺序（仅显示已配置的链接，未配置的链接不会显示。新配置的链接会自动添加到末尾）
              </n-text>
              <div style="margin-top: 8px;">
                <n-text depth="3" style="font-size: 12px; color: #f90;">
                  ⚠️ 注意：前端个人名片最多只显示前 5 个已配置的社交链接，请将重要的链接排在前面
                </n-text>
              </div>
            </div>
            <div class="sortable-list">
              <div
                v-for="(item, index) in sortedSocialLinks"
                :key="item.type"
                class="sortable-item"
                :class="{ dragging: draggedIndex === index }"
                :draggable="true"
                @dragstart="handleDragStart(index, $event)"
                @dragover.prevent="handleDragOver(index, $event)"
                @drop="handleDrop(index, $event)"
                @dragend="handleDragEnd"
              >
                <div class="drag-handle">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" style="width: 18px; height: 18px; color: #999; cursor: move;">
                    <path d="M9 5h2v2H9V5zm0 6h2v2H9v-2zm0 6h2v2H9v-2zm4-12h2v2h-2V5zm0 6h2v2h-2v-2zm0 6h2v2h-2v-2z"/>
                  </svg>
                </div>
                <div class="item-content">
                  <span class="item-label">{{ item.label }}</span>
                  <n-tag v-if="formData[item.key as keyof typeof formData]?.trim()" size="small" type="success">已配置</n-tag>
                  <n-tag v-else size="small" type="default">未配置</n-tag>
                </div>
              </div>
            </div>
          </div>
        </n-form-item>

        <n-form-item>
          <n-space>
            <n-button type="primary" @click="handleSubmit" :loading="loading">
              保存设置
            </n-button>
            <n-button @click="handleReset">
              重置
            </n-button>
          </n-space>
        </n-form-item>
      </n-form>
    </n-card>

    <n-card title="上传存储配置" style="margin-top: 24px;">
      <n-form
        ref="uploadFormRef"
        :model="uploadFormData"
        :label-placement="isMobile ? 'top' : 'left'"
        :label-width="isMobile ? 'auto' : '120'"
        require-mark-placement="right-hanging"
      >
        <n-form-item label="存储方式" path="storage_type">
          <n-radio-group v-model:value="uploadFormData.storage_type">
            <n-space>
              <n-radio value="local">本地存储</n-radio>
              <n-radio value="oss">阿里云 OSS</n-radio>
              <n-radio value="cos">腾讯云 COS</n-radio>
            </n-space>
          </n-radio-group>
        </n-form-item>

        <n-alert v-if="uploadFormData.storage_type === 'oss'" type="info" style="margin-bottom: 16px;">
          OSS 配置需要在服务器配置文件中设置（config/config-dev.yml 或 config/config-prod.yml 的 oss 节点）
        </n-alert>

        <n-alert v-if="uploadFormData.storage_type === 'cos'" type="info" style="margin-bottom: 16px;">
          COS 配置需要在服务器配置文件中设置（config/config-dev.yml 或 config/config-prod.yml 的 cos 节点）
        </n-alert>

        <n-form-item>
          <n-space>
            <n-button type="primary" @click="handleUploadSubmit" :loading="uploadLoading">
              保存配置
            </n-button>
            <n-button @click="handleUploadReset">
              重置
            </n-button>
          </n-space>
        </n-form-item>
      </n-form>
    </n-card>

    <n-card title="通知配置" style="margin-top: 24px;">
      <n-form
        ref="notificationFormRef"
        :model="notificationFormData"
        :label-placement="isMobile ? 'top' : 'left'"
        :label-width="isMobile ? 'auto' : '120'"
        require-mark-placement="right-hanging"
      >
        <n-form-item label="管理员通知" path="notify_admin_on_comment">
          <n-switch v-model:value="notificationFormData.notify_admin_on_comment" />
          <template #feedback>
            <span style="font-size: 12px; color: #999">
              开启后，当有用户评论文章时，所有管理员都会收到评论通知邮件
            </span>
          </template>
        </n-form-item>

        <n-alert type="info" style="margin-bottom: 16px;">
          <template #header>通知说明</template>
          <ul style="margin: 8px 0; padding-left: 20px;">
            <li><strong>管理员通知</strong>：需要在此处开启，开启后所有管理员都会收到评论通知邮件</li>
            <li><strong>说明</strong>：由于普通用户没有权限写文章，文章作者只能是管理员，因此统一通过管理员通知处理</li>
            <li><strong>邮件配置</strong>：需要在服务器配置文件中设置邮件参数（config/config-dev.yml 或 config/config-prod.yml 的 email 节点）</li>
          </ul>
        </n-alert>

        <n-form-item>
          <n-space>
            <n-button type="primary" @click="handleNotificationSubmit" :loading="notificationLoading">
              保存配置
            </n-button>
            <n-button @click="handleNotificationReset">
              重置
            </n-button>
          </n-space>
        </n-form-item>
      </n-form>
    </n-card>

    <n-card title="打赏配置" style="margin-top: 24px;">
      <n-form
        :model="formData"
        :label-placement="isMobile ? 'top' : 'left'"
        :label-width="isMobile ? 'auto' : '120'"
      >
        <n-form-item label="微信收款码" path="reward_wechat">
          <div style="width: 100%;">
            <ImageUpload
              v-model="formData.reward_wechat"
              :width="280"
              :height="160"
              alt="微信收款码"
              compact
            />
            <n-input
              v-model:value="formData.reward_wechat"
              placeholder="微信收款码图片URL，也可通过上方直接上传"
              maxlength="500"
              clearable
              style="margin-top: 8px;"
            >
              <template #suffix>
                <n-button text size="tiny" @click="clearField('reward_wechat')" type="error">清空</n-button>
              </template>
            </n-input>
          </div>
        </n-form-item>

        <n-form-item label="支付宝收款码" path="reward_alipay">
          <div style="width: 100%;">
            <ImageUpload
              v-model="formData.reward_alipay"
              :width="280"
              :height="160"
              alt="支付宝收款码"
              compact
            />
            <n-input
              v-model:value="formData.reward_alipay"
              placeholder="支付宝收款码图片URL，也可通过上方直接上传"
              maxlength="500"
              clearable
              style="margin-top: 8px;"
            >
              <template #suffix>
                <n-button text size="tiny" @click="clearField('reward_alipay')" type="error">清空</n-button>
              </template>
            </n-input>
          </div>
        </n-form-item>

        <n-alert type="info" style="margin-bottom: 16px;">
          上传收款码图片后将图片URL填入对应字段，保存后导航栏打赏按钮将展示对应收款码。若两个字段均为空，点击打赏按钮时将提示"暂未开放打赏功能"。
        </n-alert>

        <n-form-item>
          <n-button type="primary" @click="handleSubmit" :loading="loading">保存设置</n-button>
        </n-form-item>
      </n-form>
    </n-card>

    <!-- 封面设置 -->
    <n-card title="首页封面设置" style="margin-top: 24px;">
      <n-form
        :model="formData"
        :label-placement="isMobile ? 'top' : 'left'"
        :label-width="isMobile ? 'auto' : '120'"
      >
        <n-form-item label="封面副标题" path="cover_subtitle">
          <n-input
            v-model:value="formData.cover_subtitle"
            type="textarea"
            placeholder="每行一条副标题，打字机效果轮播显示&#10;例如：&#10;记录技术，分享生活&#10;每一行代码都是一次对话"
            :autosize="{ minRows: 3, maxRows: 6 }"
            maxlength="500"
            show-count
            clearable
          />
        </n-form-item>

        <n-form-item label="封面背景图" path="cover_bg_images">
          <div style="width: 100%;">
            <MultiImageUpload
              v-model="formData.cover_bg_images"
              :max-count="3"
              :upload-fn="uploadSiteCover"
            />
            <n-text depth="3" style="font-size: 12px; margin-top: 8px; display: block;">
              建议上传 1920×1080 或更高分辨率的横版图片，最多上传 3 张，每次刷新页面随机显示其中一张
            </n-text>
            <n-text depth="3" style="font-size: 12px; color: #f90; display: block; margin-top: 4px;">
              ⚠️ 未上传时将使用默认背景图
            </n-text>
          </div>
        </n-form-item>

        <n-form-item>
          <n-button type="primary" @click="handleSubmit" :loading="loading">保存设置</n-button>
        </n-form-item>
      </n-form>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
import { useMessage, type FormInst } from 'naive-ui'
import { getSiteSettings, updateSiteSettings, getUploadSettings, updateUploadSettings, getNotificationSettings, updateNotificationSettings } from '@/api/setting'
import { uploadSiteCover } from '@/api/upload'
import ImageUpload from '@/components/ImageUpload.vue'
import MultiImageUpload from '@/components/MultiImageUpload.vue'

const message = useMessage()

const formRef = ref<FormInst | null>(null)
const uploadFormRef = ref<FormInst | null>(null)
const notificationFormRef = ref<FormInst | null>(null)

const defaultFormData = {
  site_name: '',
  site_icp: '',
  site_police: '',
  social_github: '',
  social_gitee: '',
  social_email: '',
  social_csdn: '',
  social_qq: '',
  social_wechat: '',
  reward_wechat: '',
  reward_alipay: '',
  cover_subtitle: '',
  cover_bg_images: '',
  site_start_date: ''
}

const formData = ref({
  site_name: '',
  site_icp: '',
  site_police: '',
  social_github: '',
  social_gitee: '',
  social_email: '',
  social_csdn: '',
  social_qq: '',
  social_wechat: '',
  reward_wechat: '',
  reward_alipay: '',
  cover_subtitle: '',
  cover_bg_images: '',
  site_start_date: ''
})

const uploadFormData = ref({
  storage_type: 'local'
})

const notificationFormData = ref({
  notify_admin_on_comment: false
})

const originalData = ref({ ...formData.value })
const originalUploadData = ref({ ...uploadFormData.value })
const originalNotificationData = ref({ ...notificationFormData.value })
const loading = ref(false)
const uploadLoading = ref(false)
const notificationLoading = ref(false)
const isMobile = ref(false)

// 社交链接配置
const socialLinkConfig = [
  { type: 'github', key: 'social_github', label: 'GitHub' },
  { type: 'gitee', key: 'social_gitee', label: 'Gitee' },
  { type: 'email', key: 'social_email', label: '邮箱' },
  { type: 'csdn', key: 'social_csdn', label: 'CSDN' },
  { type: 'qq', key: 'social_qq', label: 'QQ' },
  { type: 'wechat', key: 'social_wechat', label: '微信' }
]

// 社交链接排序
const socialLinkOrder = ref<string[]>([])
const draggedIndex = ref<number | null>(null)

// 计算排序后的社交链接列表（只显示已配置的链接）
const sortedSocialLinks = computed(() => {
  // 先筛选出已配置的链接
  const configuredLinks = socialLinkConfig.filter(item => {
    const value = formData.value[item.key as keyof typeof formData.value]
    return value && String(value).trim()
  })
  
  if (configuredLinks.length === 0) {
    return []
  }
  
  // 如果没有排序配置，返回默认顺序的已配置链接
  if (socialLinkOrder.value.length === 0) {
    return configuredLinks
  }
  
  // 按照保存的顺序排序已配置的链接
  const ordered = socialLinkOrder.value
    .map(type => configuredLinks.find(item => item.type === type))
    .filter(Boolean) as typeof socialLinkConfig
  
  // 添加未在排序中的已配置链接（新配置的链接）到末尾
  const orderedTypes = new Set(socialLinkOrder.value)
  const unordered = configuredLinks.filter(item => !orderedTypes.has(item.type))
  
  return [...ordered, ...unordered]
})

// 拖拽开始
function handleDragStart(index: number, event: DragEvent) {
  draggedIndex.value = index
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'move'
    event.dataTransfer.setData('text/html', '')
  }
}

// 拖拽悬停
function handleDragOver(_index: number, event: DragEvent) {
  event.preventDefault()
  if (event.dataTransfer) {
    event.dataTransfer.dropEffect = 'move'
  }
}

// 放置
function handleDrop(index: number, event: DragEvent) {
  event.preventDefault()
  if (draggedIndex.value === null || draggedIndex.value === index) {
    return
  }
  
  const items = [...sortedSocialLinks.value]
  const draggedItem = items[draggedIndex.value]
  items.splice(draggedIndex.value, 1)
  items.splice(index, 0, draggedItem)
  
  // 更新排序顺序
  socialLinkOrder.value = items.map(item => item.type)
  
  draggedIndex.value = null
}

// 拖拽结束
function handleDragEnd() {
  draggedIndex.value = null
}

// 更新社交链接排序：移除未配置的，添加新配置的到末尾
function updateSocialLinkOrder() {
  // 获取当前已配置的链接类型
  const configuredTypes = socialLinkConfig
    .filter(item => {
      const value = formData.value[item.key as keyof typeof formData.value]
      return value && String(value).trim()
    })
    .map(item => item.type)
  
  if (configuredTypes.length === 0) {
    socialLinkOrder.value = []
    return
  }
  
  // 保留排序中已配置的链接
  const existingOrdered = socialLinkOrder.value.filter(type => configuredTypes.includes(type))
  
  // 添加新配置的链接到末尾
  const existingOrderedSet = new Set(existingOrdered)
  const newConfigured = configuredTypes.filter(type => !existingOrderedSet.has(type))
  
  socialLinkOrder.value = [...existingOrdered, ...newConfigured]
}

// 检测移动设备
function checkMobile() {
  isMobile.value = window.innerWidth <= 768
}

// 获取网站设置
async function fetchSettings() {
  try {
    const res = await getSiteSettings()
    if (res.data) {
      formData.value = {
        site_name: res.data.site_name || '',
        site_icp: res.data.site_icp || '',
        site_police: res.data.site_police || '',
        social_github: res.data.social_github || '',
        social_gitee: res.data.social_gitee || '',
        social_email: res.data.social_email || '',
        social_csdn: res.data.social_csdn || '',
        social_qq: res.data.social_qq || '',
        social_wechat: res.data.social_wechat || '',
        reward_wechat: res.data.reward_wechat || '',
        reward_alipay: res.data.reward_alipay || '',
        cover_subtitle: res.data.cover_subtitle || '',
        cover_bg_images: res.data.cover_bg_images || '',
        site_start_date: res.data.site_start_date || ''
      }
      
      // 初始化社交链接排序
      if (res.data.social_link_order) {
        socialLinkOrder.value = res.data.social_link_order.split(',').filter(Boolean)
      } else {
        // 默认顺序（只包含已配置的）
        socialLinkOrder.value = []
      }
      
      // 更新排序列表：移除未配置的，添加新配置的到末尾
      updateSocialLinkOrder()
      
      originalData.value = { ...formData.value }
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '获取设置失败')
  }
}

// 获取上传配置
async function fetchUploadSettings() {
  try {
    const res = await getUploadSettings()
    if (res.data) {
      uploadFormData.value = {
        storage_type: res.data.storage_type || 'local'
      }
      originalUploadData.value = { ...uploadFormData.value }
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '获取上传配置失败')
  }
}

// 获取通知配置
async function fetchNotificationSettings() {
  try {
    const res = await getNotificationSettings()
    if (res.data) {
      notificationFormData.value = {
        notify_admin_on_comment: res.data.notify_admin_on_comment === '1' || res.data.notify_admin_on_comment === 'true'
      }
      originalNotificationData.value = { ...notificationFormData.value }
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '获取通知配置失败')
  }
}

// 提交表单
async function handleSubmit() {
  loading.value = true
  try {
    // 更新社交链接排序（移除未配置的，添加新配置的到末尾）
    updateSocialLinkOrder()
    
    // 提交全部字段（包含空字符串），以便支持清空配置
    const dataToSubmit: Record<string, string> = {}
    Object.keys(formData.value).forEach(key => {
      const value = (formData.value as any)[key]
      dataToSubmit[key] = value === null || value === undefined ? '' : String(value).trim()
    })
    
    // 添加社交链接排序顺序
    dataToSubmit.social_link_order = socialLinkOrder.value.join(',')

    await updateSiteSettings(dataToSubmit)
    message.success('设置保存成功')
    originalData.value = { ...formData.value }
    
    // 提示用户刷新页面查看效果
    message.info('社交链接已更新，请刷新首页查看效果')
  } catch (error: any) {
    message.error(error.response?.data?.message || '保存失败')
  } finally {
    loading.value = false
  }
}

// 提交上传配置
async function handleUploadSubmit() {
  uploadLoading.value = true
  try {
    await updateUploadSettings(uploadFormData.value)
    message.success('上传配置保存成功')
    originalUploadData.value = { ...uploadFormData.value }
  } catch (error: any) {
    message.error(error.response?.data?.message || '保存失败')
  } finally {
    uploadLoading.value = false
  }
}

// 重置表单
function handleReset() {
  // 重置为初始默认值（清空所有字段）
  formData.value = { ...defaultFormData }
  originalData.value = { ...defaultFormData }
  formRef.value?.restoreValidation()
  message.info('已重置为初始默认值，请保存后生效')
}

// 重置上传配置
function handleUploadReset() {
  uploadFormData.value = { ...originalUploadData.value }
  uploadFormRef.value?.restoreValidation()
  message.info('已重置为上次保存的数据')
}

// 提交通知配置
async function handleNotificationSubmit() {
  notificationLoading.value = true
  try {
    await updateNotificationSettings({
      notify_admin_on_comment: notificationFormData.value.notify_admin_on_comment ? '1' : '0'
    })
    message.success('通知配置保存成功')
    originalNotificationData.value = { ...notificationFormData.value }
  } catch (error: any) {
    message.error(error.response?.data?.message || '保存失败')
  } finally {
    notificationLoading.value = false
  }
}

// 重置通知配置
function handleNotificationReset() {
  notificationFormData.value = { ...originalNotificationData.value }
  notificationFormRef.value?.restoreValidation()
  message.info('已重置为上次保存的数据')
}

// 单字段清除
function clearField(key: keyof typeof formData.value) {
  formData.value[key] = ''
}

// 监听表单数据变化，自动更新排序列表
watch(
  () => [
    formData.value.social_github,
    formData.value.social_gitee,
    formData.value.social_email,
    formData.value.social_csdn,
    formData.value.social_qq,
    formData.value.social_wechat
  ],
  () => {
    // 延迟更新，避免频繁触发
    setTimeout(() => {
      updateSocialLinkOrder()
    }, 100)
  },
  { deep: true }
)

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  fetchSettings()
  fetchUploadSettings()
  fetchNotificationSettings()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped>
.site-settings-page {
  padding: 0;
}

.site-settings-page :deep(.n-card) {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(8, 145, 178, 0.1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  transition: all 0.3s;
}

html.dark .site-settings-page :deep(.n-card) {
  background: rgba(30, 41, 59, 0.85);
  border: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

.site-settings-page :deep(.n-card:hover) {
  transform: translateY(-2px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.12);
}

.site-settings-page :deep(.n-form-item) {
  margin-bottom: 24px;
}

.site-settings-page p {
  margin: 8px 0;
  line-height: 1.6;
}

/* 社交链接排序样式 */
.social-link-sort-container {
  width: 100%;
}

.sort-tip {
  margin-bottom: 12px;
}

.sortable-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.sortable-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f5f5f5;
  border-radius: 6px;
  cursor: move;
  transition: all 0.2s;
  border: 1px solid transparent;
}

.sortable-item:hover {
  background: #e8e8e8;
  border-color: #d9d9d9;
}

.sortable-item.dragging {
  opacity: 0.5;
  background: #e0e0e0;
}

.drag-handle {
  display: flex;
  align-items: center;
  cursor: grab;
  user-select: none;
}

.drag-handle:active {
  cursor: grabbing;
}

.item-content {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.item-label {
  font-weight: 500;
  color: #333;
}

html.dark .sortable-item {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.1);
}

html.dark .sortable-item:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.2);
}

html.dark .sortable-item.dragging {
  background: rgba(255, 255, 255, 0.08);
}

html.dark .item-label {
  color: #e5e5e5;
}

/* 移动端样式 */
@media (max-width: 768px) {
  .site-settings-page :deep(.n-form-item) {
    margin-bottom: 20px;
  }
  
  .site-settings-page p {
    font-size: 14px;
  }
  
  .sortable-item {
    padding: 10px;
  }
  
  .item-content {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}
</style>

