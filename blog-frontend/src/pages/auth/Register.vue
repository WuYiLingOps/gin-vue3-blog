<template>
  <div class="register-page">
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

const router = useRouter()
const message = useMessage()
const authStore = useAuthStore()

const formRef = ref<FormInst | null>(null)
const loading = ref(false)

const formData = reactive<RegisterForm>({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
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
  ]
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
</script>

<style scoped>
.register-page {
  width: 100%;
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

