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
              <n-button text tertiary size="tiny" @click="clearField('site_name')" type="error">
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
              <n-button text tertiary size="tiny" @click="clearField('site_icp')" type="error">
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
              <n-button text tertiary size="tiny" @click="clearField('site_police')" type="error">
                清空
              </n-button>
            </template>
          </n-input>
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
              <n-button text tertiary size="tiny" @click="clearField('social_github')" type="error">
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
              <n-button text tertiary size="tiny" @click="clearField('social_gitee')" type="error">
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
              <n-button text tertiary size="tiny" @click="clearField('social_email')" type="error">
                清空
              </n-button>
            </template>
          </n-input>
        </n-form-item>

        <n-form-item label="RSS链接" path="social_rss">
          <n-input
            v-model:value="formData.social_rss"
            placeholder="例如：https://example.com/rss.xml"
            maxlength="200"
            clearable
          >
            <template #suffix>
              <n-button text tertiary size="tiny" @click="clearField('social_rss')" type="error">
                清空
              </n-button>
            </template>
          </n-input>
        </n-form-item>

        <n-form-item label="QQ号" path="social_qq">
          <n-input
            v-model:value="formData.social_qq"
            placeholder="例如：123456789"
            maxlength="20"
            clearable
          >
            <template #suffix>
              <n-button text tertiary size="tiny" @click="clearField('social_qq')" type="error">
                清空
              </n-button>
            </template>
          </n-input>
        </n-form-item>

        <n-form-item label="微信二维码" path="social_wechat">
          <n-input
            v-model:value="formData.social_wechat"
            placeholder="微信二维码图片URL，例如：https://example.com/wechat-qr.png"
            maxlength="500"
            clearable
          >
            <template #suffix>
              <n-button text tertiary size="tiny" @click="clearField('social_wechat')" type="error">
                清空
              </n-button>
            </template>
          </n-input>
          <template #feedback>
            <span style="font-size: 12px; color: #999;">上传微信二维码图片后，将图片URL填入此处</span>
          </template>
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
            </n-space>
          </n-radio-group>
        </n-form-item>

        <n-alert v-if="uploadFormData.storage_type === 'oss'" type="info" style="margin-bottom: 16px;">
          OSS 配置需要在服务器配置文件中设置（config/config-dev.yml 或 config/config-prod.yml）
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

    <n-card title="设置说明" style="margin-top: 24px;">
      <n-space vertical>
        <p><strong>网站名称：</strong>显示在网站底部的名称</p>
        <p><strong>ICP备案号：</strong>工信部备案号，点击后会跳转到 beian.miit.gov.cn</p>
        <p><strong>公安备案号：</strong>公安部备案号，点击后会跳转到 www.beian.gov.cn</p>
        <n-divider />
        <p><strong>社交链接：</strong>配置个人信息卡片中显示的社交链接（未配置的链接不会显示）</p>
        <p><strong>GitHub链接：</strong>您的GitHub主页地址</p>
        <p><strong>Gitee链接：</strong>您的Gitee主页地址</p>
        <p><strong>邮箱地址：</strong>联系邮箱</p>
        <p><strong>RSS链接：</strong>RSS订阅地址</p>
        <p><strong>QQ号：</strong>QQ号码，点击后会打开QQ聊天窗口</p>
        <p><strong>微信二维码：</strong>微信二维码图片的URL地址</p>
        <n-divider />
        <p><strong>存储方式：</strong>选择文件上传的存储方式</p>
        <p><strong>本地存储：</strong>文件保存在服务器本地，适合小型网站或开发环境</p>
        <p><strong>阿里云 OSS：</strong>文件保存到阿里云对象存储，适合生产环境</p>
        <p style="color: #f90; font-size: 13px;">
          ⚠️ 重要：使用 OSS 存储前，请先在服务器配置文件中填写 OSS 相关参数（endpoint、access_key_id、access_key_secret、bucket_name）
        </p>
      </n-space>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useMessage, type FormInst } from 'naive-ui'
import { getSiteSettings, updateSiteSettings, getUploadSettings, updateUploadSettings } from '@/api/setting'

const message = useMessage()

const formRef = ref<FormInst | null>(null)
const uploadFormRef = ref<FormInst | null>(null)

const defaultFormData = {
  site_name: '',
  site_icp: '',
  site_police: '',
  social_github: '',
  social_gitee: '',
  social_email: '',
  social_rss: '',
  social_qq: '',
  social_wechat: ''
}

const formData = ref({
  site_name: '',
  site_icp: '',
  site_police: '',
  social_github: '',
  social_gitee: '',
  social_email: '',
  social_rss: '',
  social_qq: '',
  social_wechat: ''
})

const uploadFormData = ref({
  storage_type: 'local'
})

const originalData = ref({ ...formData.value })
const originalUploadData = ref({ ...uploadFormData.value })
const loading = ref(false)
const uploadLoading = ref(false)
const isMobile = ref(false)

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
        social_rss: res.data.social_rss || '',
        social_qq: res.data.social_qq || '',
        social_wechat: res.data.social_wechat || ''
      }
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

// 提交表单
async function handleSubmit() {
  loading.value = true
  try {
    // 提交全部字段（包含空字符串），以便支持清空配置
    const dataToSubmit: Record<string, string> = {}
    Object.keys(formData.value).forEach(key => {
      const value = (formData.value as any)[key]
      dataToSubmit[key] = value === null || value === undefined ? '' : String(value).trim()
    })

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

// 单字段清除
function clearField(key: keyof typeof formData.value) {
  formData.value[key] = ''
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  fetchSettings()
  fetchUploadSettings()
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

/* 移动端样式 */
@media (max-width: 768px) {
  .site-settings-page :deep(.n-form-item) {
    margin-bottom: 20px;
  }
  
  .site-settings-page p {
    font-size: 14px;
  }
}
</style>

