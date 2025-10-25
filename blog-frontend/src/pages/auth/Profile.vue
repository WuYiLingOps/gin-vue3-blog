<template>
  <div class="profile-page">
    <n-grid :cols="1" :x-gap="24" :y-gap="24">
      <!-- 基本信息 -->
      <n-gi>
        <n-card title="基本信息">
          <n-form ref="profileFormRef" :model="profileForm" label-width="80px">
            <n-form-item label="头像">
              <avatar-upload
                v-model="profileForm.avatar"
                :size="100"
                :default-text="authStore.user?.nickname?.charAt(0) || 'U'"
                @success="handleAvatarSuccess"
              />
            </n-form-item>

            <n-form-item label="用户名">
              <n-input v-model:value="authStore.user!.username" disabled />
            </n-form-item>

            <n-form-item label="昵称" path="nickname">
              <n-input v-model:value="profileForm.nickname" placeholder="请输入昵称" />
            </n-form-item>

            <n-form-item label="邮箱">
              <n-input v-model:value="authStore.user!.email" disabled />
            </n-form-item>

            <n-form-item label="个人简介" path="bio">
              <n-input
                v-model:value="profileForm.bio"
                type="textarea"
                :rows="4"
                placeholder="介绍一下自己吧"
              />
            </n-form-item>

            <n-form-item>
              <n-button type="primary" :loading="updating" @click="handleUpdateProfile">
                保存修改
              </n-button>
            </n-form-item>
          </n-form>
        </n-card>
      </n-gi>

      <!-- 修改密码 -->
      <n-gi>
        <n-card title="修改密码">
          <n-form
            ref="passwordFormRef"
            :model="passwordForm"
            :rules="passwordRules"
            label-width="100px"
          >
            <n-form-item label="当前密码" path="old_password">
              <n-input
                v-model:value="passwordForm.old_password"
                type="password"
                show-password-on="click"
                placeholder="请输入当前密码"
              />
            </n-form-item>

            <n-form-item label="新密码" path="new_password">
              <n-input
                v-model:value="passwordForm.new_password"
                type="password"
                show-password-on="click"
                placeholder="请输入新密码（至少6个字符）"
              />
            </n-form-item>

            <n-form-item label="确认新密码" path="confirm_password">
              <n-input
                v-model:value="passwordForm.confirm_password"
                type="password"
                show-password-on="click"
                placeholder="请再次输入新密码"
              />
            </n-form-item>

            <n-form-item>
              <n-button type="primary" :loading="changingPassword" @click="handleChangePassword">
                修改密码
              </n-button>
            </n-form-item>
          </n-form>
        </n-card>
      </n-gi>

      <!-- 账号信息 -->
      <n-gi>
        <n-card title="账号信息">
          <n-descriptions :column="1">
            <n-descriptions-item label="用户ID">
              {{ authStore.user?.id }}
            </n-descriptions-item>
            <n-descriptions-item label="角色">
              <n-tag :type="authStore.isAdmin ? 'error' : 'info'">
                {{ authStore.isAdmin ? '管理员' : '普通用户' }}
              </n-tag>
            </n-descriptions-item>
            <n-descriptions-item label="注册时间">
              {{ formatDate(authStore.user!.created_at) }}
            </n-descriptions-item>
            <n-descriptions-item label="最后更新">
              {{ formatDate(authStore.user!.updated_at) }}
            </n-descriptions-item>
          </n-descriptions>
        </n-card>
      </n-gi>
    </n-grid>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui'
import { useAuthStore } from '@/stores'
import { updateProfile, updatePassword } from '@/api/auth'
import { formatDate } from '@/utils/format'
import type { ProfileForm, PasswordForm } from '@/types/auth'
import AvatarUpload from '@/components/AvatarUpload.vue'

const message = useMessage()
const authStore = useAuthStore()

const profileFormRef = ref<FormInst | null>(null)
const passwordFormRef = ref<FormInst | null>(null)
const updating = ref(false)
const changingPassword = ref(false)

const profileForm = reactive<ProfileForm>({
  nickname: '',
  avatar: '',
  bio: ''
})

const passwordForm = reactive<PasswordForm>({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const passwordRules: FormRules = {
  old_password: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码至少6个字符', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    {
      validator: (_rule, value) => value === passwordForm.new_password,
      message: '两次密码不一致',
      trigger: ['blur', 'input']
    }
  ]
}

onMounted(() => {
  if (authStore.user) {
    profileForm.nickname = authStore.user.nickname
    profileForm.avatar = authStore.user.avatar
    profileForm.bio = authStore.user.bio
  }
})

async function handleAvatarSuccess(url: string) {
  profileForm.avatar = url
  
  // 上传头像后自动保存到用户资料
  await updateProfile({ avatar: url })
  await authStore.fetchUserInfo()
  message.success('头像更新成功')
}

async function handleUpdateProfile() {
  await profileFormRef.value?.validate()
  updating.value = true

  await updateProfile(profileForm)
  await authStore.fetchUserInfo()
  
  // 重新加载表单数据（确保显示最新的头像）
  if (authStore.user) {
    profileForm.nickname = authStore.user.nickname
    profileForm.avatar = authStore.user.avatar
    profileForm.bio = authStore.user.bio
  }
  
  updating.value = false
  message.success('个人信息更新成功')
}

async function handleChangePassword() {
  await passwordFormRef.value?.validate()
  changingPassword.value = true

  await updatePassword(passwordForm)
  message.success('密码修改成功，请重新登录')

  // 清空表单
  passwordForm.old_password = ''
  passwordForm.new_password = ''
  passwordForm.confirm_password = ''

  changingPassword.value = false

  // 退出登录
  setTimeout(() => {
    authStore.logout()
  }, 1500)
}
</script>

<style scoped>
.profile-page {
  max-width: 800px;
  margin: 0 auto;
}
</style>

