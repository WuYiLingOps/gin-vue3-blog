<!--
  项目名称：blog-frontend
  文件名称：CaptchaInput.vue
  创建时间：2026-02-01 20:03:19

  系统用户：Administrator
  作　　者：無以菱
  联系邮箱：huangjing510@126.com
  功能描述：验证码输入组件，包含验证码图片显示和输入框，支持点击刷新验证码，组件挂载时自动获取验证码。
-->
<template>
  <div class="captcha-container">
    <n-input
      v-model:value="captchaValue"
      placeholder="请输入验证码"
      :disabled="loading"
      @input="handleInput"
      @keyup.enter="handleEnter"
    />
    <div class="captcha-image" @click="refreshCaptcha">
      <n-spin :show="loading" size="small">
        <img v-if="imageData" :src="imageData" alt="验证码" />
        <div v-else class="captcha-placeholder">点击获取验证码</div>
      </n-spin>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import { getCaptcha } from '@/api/auth'

const emit = defineEmits<{
  (e: 'update:captchaId', value: string): void
  (e: 'update:captcha', value: string): void
  (e: 'enter'): void
}>()

const message = useMessage()
const loading = ref(false)
const captchaValue = ref('')
const imageData = ref('')
const captchaId = ref('')

// 获取验证码
async function refreshCaptcha() {
  try {
    loading.value = true
    const res = await getCaptcha()
    if (res.data) {
      captchaId.value = res.data.captcha_id
      imageData.value = res.data.image_data
      captchaValue.value = ''
      emit('update:captchaId', captchaId.value)
      emit('update:captcha', '')
    }
  } catch (error: any) {
    message.error(error.message || '获取验证码失败')
  } finally {
    loading.value = false
  }
}

function handleInput() {
  emit('update:captcha', captchaValue.value)
}

function handleEnter() {
  emit('enter')
}

// 组件挂载时自动获取验证码
onMounted(() => {
  refreshCaptcha()
})

// 暴露刷新方法给父组件
defineExpose({
  refresh: refreshCaptcha
})
</script>

<style scoped>
.captcha-container {
  display: flex;
  gap: 12px;
  align-items: center;
}

.captcha-image {
  flex-shrink: 0;
  width: 120px;
  height: 40px;
  cursor: pointer;
  border: 1px solid #e0e0e6;
  border-radius: 4px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
  transition: all 0.3s;
}

.captcha-image:hover {
  border-color: #18a058;
  box-shadow: 0 0 0 2px rgba(24, 160, 88, 0.1);
}

.captcha-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.captcha-placeholder {
  font-size: 12px;
  color: #999;
  text-align: center;
  padding: 0 8px;
}
</style>

