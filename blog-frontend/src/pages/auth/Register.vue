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

      <n-form-item path="captcha" label="验证码">
        <captcha-input
          ref="captchaRef"
          v-model:captcha-id="formData.captcha_id"
          v-model:captcha="formData.captcha"
          @enter="handleRegister"
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
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui'
import { useAuthStore } from '@/stores'
import { validateEmail, validateUsername, validatePassword } from '@/utils/validator'
import type { RegisterForm } from '@/types/auth'
import CaptchaInput from '@/components/CaptchaInput.vue'

const router = useRouter()
const message = useMessage()
const authStore = useAuthStore()

const formRef = ref<FormInst | null>(null)
const captchaRef = ref<InstanceType<typeof CaptchaInput> | null>(null)
const loading = ref(false)

const formData = reactive<RegisterForm>({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  captcha_id: '',
  captcha: ''
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
  ],
  captcha: [{ required: true, message: '请输入验证码', trigger: 'blur' }]
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
    // 注册失败后刷新验证码
    captchaRef.value?.refresh()
  } finally {
    loading.value = false
  }
}
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

.footer-links {
  margin-top: 24px;
  text-align: center;
  color: #666;
}

.footer-links span {
  margin-right: 8px;
}
</style>

