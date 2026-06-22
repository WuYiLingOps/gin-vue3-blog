<!--
 * @ProjectName: go-vue3-blog
 * @FileName: Register.vue
 * @CreateTime: 2026-02-02 11:43:17
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 注册页面组件，提供用户注册功能
 -->
<template>
  <div class="register-page">
    <h2 class="page-title">
      <span class="gradient-text">创建账号</span>
    </h2>
    <p class="page-subtitle">创建账号，开启互动之旅</p>

    <n-form ref="formRef" :model="formData" :rules="rules" size="large">
      <n-form-item path="username" label="用户名">
        <n-input v-model:value="formData.username" placeholder="3-20个字符，字母数字下划线">
          <template #prefix>
            <n-icon :component="PersonOutline" />
          </template>
        </n-input>
      </n-form-item>

      <n-form-item path="email" label="邮箱">
        <n-input v-model:value="formData.email" placeholder="请输入邮箱">
          <template #prefix>
            <n-icon :component="MailOutline" />
          </template>
        </n-input>
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
        >
          <template #prefix>
            <n-icon :component="LockClosedOutline" />
          </template>
        </n-input>
      </n-form-item>

      <n-form-item path="confirmPassword" label="确认密码">
        <n-input
          v-model:value="formData.confirmPassword"
          type="password"
          show-password-on="click"
          placeholder="请再次输入密码"
        >
          <template #prefix>
            <n-icon :component="LockClosedOutline" />
          </template>
        </n-input>
      </n-form-item>

      <n-button type="primary" block size="large" :loading="loading" @click="handleRegister">
        注册
      </n-button>
    </n-form>

    <div class="footer-links">
      <span>已有账号？</span>
      <n-button text type="primary" @click="router.push('/auth/login')"> 立即登录 </n-button>
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
import { sendRegisterCode, VERIFICATION_CODE_RESEND_SECONDS } from '@/api/auth'
import { PersonOutline, MailOutline, LockClosedOutline } from '@vicons/ionicons5'

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
    {
      validator: (_rule, value) => validateEmail(value),
      message: '邮箱格式不正确',
      trigger: 'blur'
    }
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
  if (!formData.username) {
    message.warning('请先输入用户名')
    return
  }

  if (!validateUsername(formData.username)) {
    message.warning('请输入正确的用户名格式')
    return
  }

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
    await sendRegisterCode({ email: formData.email, username: formData.username })
    message.success('验证码已发送，请查收邮箱')

    // 开始倒计时
    countdown.value = VERIFICATION_CODE_RESEND_SECONDS
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

.register-page :deep(.n-form-item) {
  margin-bottom: 14px;
}

.register-page :deep(.n-form-item-label) {
  padding-bottom: 6px;
}

.page-title {
  text-align: center;
  margin-bottom: 8px;
  font-size: 28px;
  font-weight: 700;
}

.page-subtitle {
  text-align: center;
  margin-bottom: 28px;
  color: #666;
  font-size: 14px;
}

html.dark .page-subtitle {
  color: #94a3b8;
}

h2 {
  text-align: center;
  margin-bottom: 20px;
  color: #333;
  font-size: 26px;
  font-weight: 600;
}

html.dark h2 {
  color: #e5e5e5;
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
  margin-top: 20px;
  text-align: center;
  color: #666;
}

.footer-links span {
  margin-right: 8px;
}

@media (max-width: 768px) {
  h2 {
    font-size: 24px;
    margin-bottom: 20px;
  }

  .register-page :deep(.n-form-item) {
    margin-bottom: 16px;
  }

  .code-input-wrapper .n-button {
    min-width: 100px;
    font-size: 13px;
    padding: 0 12px;
  }
}
</style>
