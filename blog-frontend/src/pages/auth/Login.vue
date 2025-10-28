<template>
  <div class="login-page">
    <h2>登录</h2>
    <n-form ref="formRef" :model="formData" :rules="rules" size="large" >
      <n-form-item path="username" label="用户名">
        <n-input
          v-model:value="formData.username"
          placeholder="请输入用户名"
          @keyup.enter="handleLogin"
        />
      </n-form-item>

      <n-form-item path="password" label="密码">
        <n-input
          v-model:value="formData.password"
          type="password"
          show-password-on="click"
          placeholder="请输入密码"
          @keyup.enter="handleLogin"
        />
      </n-form-item>

      <n-form-item path="captcha" label="验证码">
        <captcha-input
          ref="captchaRef"
          v-model:captcha-id="formData.captcha_id"
          v-model:captcha="formData.captcha"
          @enter="handleLogin"
        />
      </n-form-item>

      <n-form-item>
        <div style="display: flex; justify-content: space-between; width: 100%; align-items: center;">
          <n-checkbox v-model:checked="formData.remember">记住我</n-checkbox>
          <n-button text type="primary" @click="router.push('/auth/forgot-password')">
            忘记密码？
          </n-button>
        </div>
      </n-form-item>

      <n-button type="primary" block size="large" :loading="loading" @click="handleLogin">
        登录
      </n-button>
    </n-form>

    <div class="footer-links">
      <span>还没有账号？</span>
      <n-button text type="primary" @click="router.push('/auth/register')">
        立即注册
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useMessage } from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui'
import { useAuthStore } from '@/stores'
import type { LoginForm } from '@/types/auth'
import CaptchaInput from '@/components/CaptchaInput.vue'

const router = useRouter()
const route = useRoute()
const message = useMessage()
const authStore = useAuthStore()

const formRef = ref<FormInst | null>(null)
const captchaRef = ref<InstanceType<typeof CaptchaInput> | null>(null)
const loading = ref(false)

const formData = reactive<LoginForm>({
  username: '',
  password: '',
  captcha_id: '',
  captcha: '',
  remember: false
})

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6个字符', trigger: 'blur' }
  ],
  captcha: [{ required: true, message: '请输入验证码', trigger: 'blur' }]
}

async function handleLogin() {
  try {
    await formRef.value?.validate()
    loading.value = true

    await authStore.login(formData)
    message.success('登录成功')

    // 重定向到来源页面或首页
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (error: any) {
    message.error(error.message || '登录失败')
    // 登录失败后刷新验证码
    captchaRef.value?.refresh()
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  width: 100%;
}

.login-page :deep(.n-form-item) {
  margin-bottom: 18px;
}

h2 {
  text-align: center;
  margin-bottom: 24px;
  color: #333;
  font-size: 26px;
  font-weight: 600;
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

