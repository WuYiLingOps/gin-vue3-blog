// 认证相关类型

// 用户信息
export interface User {
  id: number
  username: string
  email: string
  nickname: string
  avatar: string
  bio: string
  role: 'admin' | 'user'
  status: number
  created_at: string
  updated_at: string
}

// 验证码响应
export interface CaptchaResponse {
  captcha_id: string
  image_data: string
}

// 登录表单
export interface LoginForm {
  username: string
  password: string
  captcha_id?: string
  captcha?: string
  remember?: boolean
}

// 注册表单
export interface RegisterForm {
  username: string
  email: string
  password: string
  confirmPassword: string
  code?: string
}

// 登录响应
export interface LoginResponse {
  token: string
  user: User
}

// 个人资料表单
export interface ProfileForm {
  nickname?: string
  avatar?: string
  bio?: string
}

// 修改密码表单
export interface PasswordForm {
  old_password: string
  new_password: string
  confirm_password?: string
}

