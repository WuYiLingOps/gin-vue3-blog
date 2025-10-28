<template>
  <div class="forgot-password-page">
    <h2>找回密码</h2>
    <p class="subtitle">请输入您注册时使用的邮箱</p>
    
    <n-form ref="formRef" :model="formData" :rules="rules" size="large">
      <n-form-item path="email">
        <n-input
          v-model:value="formData.email"
          placeholder="邮箱地址"
          size="large"
          @keyup.enter="handleSubmit"
        >
          <template #prefix>
            <n-icon :component="MailOutline" />
          </template>
        </n-input>
      </n-form-item>

      <n-form-item path="code">
        <n-input
          v-model:value="formData.code"
          placeholder="验证码"
          maxlength="6"
          size="large"
          @keyup.enter="handleSubmit"
        >
          <template #prefix>
            <n-icon :component="ShieldCheckmarkOutline" />
          </template>
          <template #suffix>
            <n-button 
              type="primary" 
              :disabled="countdown > 0"
              :loading="sending"
              @click="handleSendCode"
            >
              {{ countdown > 0 ? `${countdown}秒后重发` : '获取验证码' }}
            </n-button>
          </template>
        </n-input>
      </n-form-item>

      <n-form-item path="new_password">
        <n-input
          v-model:value="formData.new_password"
          type="password"
          show-password-on="click"
          placeholder="新密码"
          size="large"
          @keyup.enter="handleSubmit"
        >
          <template #prefix>
            <n-icon :component="LockClosedOutline" />
          </template>
        </n-input>
      </n-form-item>

      <n-form-item path="confirm_password" class="last-item">
        <n-input
          v-model:value="formData.confirm_password"
          type="password"
          show-password-on="click"
          placeholder="确认密码"
          size="large"
          @keyup.enter="handleSubmit"
        >
          <template #prefix>
            <n-icon :component="LockClosedOutline" />
          </template>
        </n-input>
      </n-form-item>

      <n-button 
        type="primary" 
        block 
        size="large" 
        :loading="submitting" 
        @click="handleSubmit"
      >
        重置密码
      </n-button>
    </n-form>

    <div class="footer-links">
      <n-button text type="primary" @click="router.push('/auth/login')">
        返回登录
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage, NIcon } from 'naive-ui'
import type { FormInst, FormRules, FormItemRule } from 'naive-ui'
import { LockClosedOutline, ShieldCheckmarkOutline, MailOutline } from '@vicons/ionicons5'
import { forgotPassword, resetPassword } from '@/api/auth'

const router = useRouter()
const message = useMessage()

const sending = ref(false)
const submitting = ref(false)
const countdown = ref(0)
let countdownTimer: number | null = null

const formRef = ref<FormInst | null>(null)

const formData = reactive({
  email: '',
  code: '',
  new_password: '',
  confirm_password: ''
})

const validatePasswordConfirm = (_rule: FormItemRule, value: string): boolean | Error => {
  if (value !== formData.new_password) {
    return new Error('两次输入的密码不一致')
  }
  return true
}

const rules: FormRules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '邮箱格式不正确', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 6, message: '验证码为6位数字', trigger: 'blur' }
  ],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码至少6个字符', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validatePasswordConfirm, trigger: ['blur', 'input'] }
  ]
}

async function handleSendCode() {
  if (countdown.value > 0) return
  
  // 验证邮箱字段
  if (!formData.email) {
    message.error('请先输入邮箱地址')
    return
  }
  
  // 验证邮箱格式
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(formData.email)) {
    message.error('邮箱格式不正确')
    return
  }
  
  try {
    sending.value = true
    await forgotPassword({
      email: formData.email
    })
    message.success('验证码已发送到您的邮箱，请查收')
    startCountdown()
  } catch (error: any) {
    message.error(error.message || '发送失败')
  } finally {
    sending.value = false
  }
}

function startCountdown() {
  countdown.value = 60
  countdownTimer = window.setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      stopCountdown()
    }
  }, 1000)
}

function stopCountdown() {
  if (countdownTimer) {
    clearInterval(countdownTimer)
    countdownTimer = null
  }
  countdown.value = 0
}

async function handleSubmit() {
  try {
    await formRef.value?.validate()
    submitting.value = true
    
    await resetPassword({
      email: formData.email,
      code: formData.code,
      new_password: formData.new_password
    })
    
    message.success('密码重置成功，请使用新密码登录')
    router.push('/auth/login')
  } catch (error: any) {
    // 如果是表单验证错误，不显示错误提示
    if (!error?.errors) {
      message.error(error.message || '重置失败')
    }
  } finally {
    submitting.value = false
  }
}

onUnmounted(() => {
  stopCountdown()
})
</script>

<style scoped>
.forgot-password-page {
  width: 100%;
}

h2 {
  text-align: center;
  margin-bottom: 6px;
  color: #333;
  font-size: 26px;
  font-weight: 600;
}

.subtitle {
  text-align: center;
  margin-bottom: 24px;
  color: #666;
  font-size: 14px;
}

.footer-links {
  margin-top: 24px;
  text-align: center;
}

.last-item {
  margin-bottom: 9px;
}

@media (max-width: 768px) {
  h2 {
    font-size: 24px;
    margin-bottom: 8px;
  }
  
  .subtitle {
    margin-bottom: 20px;
  }
  
  .forgot-password-page :deep(.n-form-item) {
    margin-bottom: 16px;
  }
  
  .last-item {
    margin-bottom: 12px;
  }
}
</style>

