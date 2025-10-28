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
          />
        </n-form-item>

        <n-form-item label="ICP备案号" path="site_icp">
          <n-input
            v-model:value="formData.site_icp"
            placeholder="例如：京ICP备12345678号"
            maxlength="50"
          />
        </n-form-item>

        <n-form-item label="公安备案号" path="site_police">
          <n-input
            v-model:value="formData.site_police"
            placeholder="例如：京公网安备11010802012345号"
            maxlength="50"
          />
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

    <n-card title="设置说明" style="margin-top: 24px;">
      <n-space vertical>
        <p><strong>网站名称：</strong>显示在网站底部的名称</p>
        <p><strong>ICP备案号：</strong>工信部备案号，点击后会跳转到 beian.miit.gov.cn</p>
        <p><strong>公安备案号：</strong>公安部备案号，点击后会跳转到 www.beian.gov.cn</p>
        <p style="color: #999; font-size: 13px;">
          提示：如果不需要显示备案信息，可以留空
        </p>
      </n-space>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useMessage } from 'naive-ui'
import { getSiteSettings, updateSiteSettings } from '@/api/setting'

const message = useMessage()

const formData = ref({
  site_name: '',
  site_icp: '',
  site_police: ''
})

const originalData = ref({...formData.value})
const loading = ref(false)
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
        site_police: res.data.site_police || ''
      }
      originalData.value = {...formData.value}
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '获取设置失败')
  }
}

// 提交表单
async function handleSubmit() {
  loading.value = true
  try {
    await updateSiteSettings(formData.value)
    message.success('设置保存成功')
    originalData.value = {...formData.value}
  } catch (error: any) {
    message.error(error.response?.data?.message || '保存失败')
  } finally {
    loading.value = false
  }
}

// 重置表单
function handleReset() {
  formData.value = {...originalData.value}
  message.info('已重置')
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  fetchSettings()
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

