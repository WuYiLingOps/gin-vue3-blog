<template>
  <div class="register-page">
    <h2>注册</h2>
    <n-form ref="formRef" :model="formData" :rules="rules" size="large">
      <n-form-item path="username" label="用户名">
        <n-input v-model:value="formData.username" placeholder="3-20个字符，字母数字下划线" />
      </n-form-item>

      <n-form-item path="email" label="邮箱">
        <n-input v-model:value="formData.email" placeholder="请输入邮箱" />
      </n-form-item>

      <n-form-item path="code" label="邮箱验证码">
        <div class="code-input-wrapper">
          <n-input
            v-model:value="formData.code"
            placeholder="请输入邮箱验证码"
            :maxlength="6"
            @keyup.enter="handleRegister"
          />
          <n-button
            type="primary"
            :disabled="sendCodeDisabled || countdown > 0"
            :loading="sendingCode"
            @click="handleSendCode"
          >
            {{ countdown > 0 ? `${countdown}秒后重试` : '发送验证码' }}
          </n-button>
        </div>
      </n-form-item>

      <n-form-item path="password" label="密码">
        <n-input
          v-model:value="formData.password"
          type="password"
          show-password-on="click"
          placeholder="至少6个字符"
        />
      </n-form-item>

      <n-form-item path="confirmPassword" label="确认密码">
        <n-input
          v-model:value="formData.confirmPassword"
          type="password"
          show-password-on="click"
          placeholder="请再次输入密码"
        />
      </n-form-item>

      <n-button type="primary" block size="large" :loading="loading" @click="handleRegister">
        注册
      </n-button>
    </n-form>

    <div class="footer-links">
      <span>已有账号？</span>
      <n-button text type="primary" @click="router.push('/auth/login')">
        立即登录
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui'
import { useAuthStore } from '@/stores'
import { validateEmail, validateUsername, validatePassword } from '@/utils/validator'
import type { RegisterForm } from '@/types/auth'
import { sendRegisterCode } from '@/api/auth'

const router = useRouter()
const message = useMessage()
const authStore = useAuthStore()

const formRef = ref<FormInst | null>(null)
const loading = ref(false)
const sendingCode = ref(false)
const countdown = ref(0)
let timer: ReturnType<typeof setInterval> | null = null

const formData = reactive<RegisterForm>({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  code: ''
})

const sendCodeDisabled = computed(() => {
  return !formData.email || !validateEmail(formData.email)
})

const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    {
      validator: (_rule, value) => validateUsername(value),
      message: '用户名格式不正确',
      trigger: 'blur'
    }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { validator: (_rule, value) => validateEmail(value), message: '邮箱格式不正确', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入邮箱验证码', trigger: 'blur' },
    { len: 6, message: '验证码为6位数字', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    {
      validator: (_rule, value) => validatePassword(value),
      message: '密码至少6个字符',
      trigger: 'blur'
    }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    {
      validator: (_rule, value) => value === formData.password,
      message: '两次密码不一致',
      trigger: ['blur', 'input']
    }
  ]
}

async function handleSendCode() {
  if (!formData.email) {
    message.warning('请先输入邮箱')
    return
  }

  if (!validateEmail(formData.email)) {
    message.warning('请输入正确的邮箱格式')
    return
  }

  try {
    sendingCode.value = true
    await sendRegisterCode({ email: formData.email })
    message.success('验证码已发送，请查收邮箱')
    
    // 开始倒计时
    countdown.value = 60
    timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0 && timer) {
        clearInterval(timer)
        timer = null
      }
    }, 1000)
  } catch (error: any) {
    message.error(error.message || '发送失败，请稍后重试')
  } finally {
    sendingCode.value = false
  }
}

async function handleRegister() {
  try {
    await formRef.value?.validate()
    loading.value = true

    await authStore.register(formData)
    message.success('注册成功，请登录')
    router.push('/auth/login')
  } catch (error: any) {
    message.error(error.message || '注册失败')
  } finally {
    loading.value = false
  }
}

// 组件卸载时清理定时器
onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
  }
})
</script>

<style scoped>
.register-page {
  width: 100%;
}

h2 {
  text-align: center;
  margin-bottom: 30px;
  color: #333;
  font-size: 24px;
}

.code-input-wrapper {
  display: flex;
  gap: 8px;
  width: 100%;
}

.code-input-wrapper :deep(.n-input) {
  flex: 1;
}

.code-input-wrapper .n-button {
  flex-shrink: 0;
  white-space: nowrap;
}

.footer-links {
  margin-top: 24px;
  text-align: center;
  color: #666;
}

.footer-links span {
  margin-right: 8px;
}
</style>

