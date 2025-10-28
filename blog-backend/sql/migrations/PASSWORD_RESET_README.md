# 找回密码功能使用说明

## 📋 功能概述

该功能允许用户通过邮箱验证码重置忘记的密码，使用QQ邮箱发送验证码。

## 🗄️ 数据库迁移

### 执行SQL迁移

运行以下SQL文件创建密码重置令牌表：

```bash
psql -U postgres -d blogdb -f blog-backend/sql/migrations/add_password_reset_tokens.sql
```

或者在数据库客户端中直接执行 `add_password_reset_tokens.sql` 文件。

## ⚙️ 配置QQ邮箱

### 1. 获取QQ邮箱授权码

1. 登录QQ邮箱网页版：https://mail.qq.com
2. 点击 **设置** → **账户**
3. 找到 **POP3/IMAP/SMTP/Exchange/CardDAV/CalDAV服务**
4. 开启 **IMAP/SMTP服务** 或 **POP3/SMTP服务**
5. 点击 **生成授权码**，按提示用手机QQ扫码验证
6. 记住生成的授权码（16位字符）

### 2. 配置邮箱信息

编辑 `blog-backend/config/config-dev.yml` 文件：

```yaml
email:
  host: smtp.qq.com
  port: 587
  username: your_email@qq.com     # 替换为你的QQ邮箱
  password: your_auth_code         # 替换为上面获取的授权码
  from_name: 个人博客系统
```

**重要提示：**
- `username` 填写完整的QQ邮箱地址
- `password` 填写的是授权码，不是QQ密码
- 授权码是16位字符（如：`abcdabcdabcdabcd`）

### 3. 生产环境配置

同样在 `blog-backend/config/config-prod.yml` 中配置：

```yaml
email:
  host: smtp.qq.com
  port: 587
  username: your_email@qq.com
  password: your_auth_code
  from_name: 个人博客系统
```

## 🚀 使用流程

### 用户端操作

1. **访问找回密码页面**
   - 在登录页面点击 "忘记密码？"
   - 或直接访问 `/auth/forgot-password`

2. **第一步：发送验证码**
   - 输入注册邮箱
   - 完成图片验证码
   - 点击"发送验证码"

3. **第二步：重置密码**
   - 查收邮箱中的6位数字验证码（有效期15分钟）
   - 输入验证码
   - 输入新密码（至少6位）
   - 确认新密码
   - 点击"重置密码"

4. **完成**
   - 使用新密码登录

## 🔒 安全特性

1. **防刷机制**
   - 同一邮箱1分钟内只能发送一次验证码
   - 需要先通过图片验证码验证

2. **验证码安全**
   - 6位随机数字验证码
   - 15分钟有效期
   - 使用后自动失效

3. **隐私保护**
   - 无论邮箱是否存在，都返回成功（防止邮箱探测）
   - 验证码不会在日志中明文记录

## 📧 邮件模板

发送的邮件包含：
- 精美的HTML格式
- 醒目的6位验证码
- 有效期提醒
- 安全提示

## 🐛 常见问题

### 1. 邮件发送失败

**可能原因：**
- QQ邮箱未开启SMTP服务
- 授权码配置错误
- 授权码已过期或被重置

**解决方案：**
- 重新生成QQ邮箱授权码
- 检查配置文件中的邮箱和授权码是否正确
- 确认防火墙未阻止587端口

### 2. 收不到邮件

**可能原因：**
- 邮件进入了垃圾箱
- 邮箱地址输入错误
- 邮件服务器延迟

**解决方案：**
- 检查垃圾邮件文件夹
- 等待1-2分钟后再检查
- 尝试重新发送

### 3. 验证码无效或已过期

**可能原因：**
- 验证码超过15分钟
- 验证码已被使用
- 输入错误

**解决方案：**
- 重新发送验证码
- 仔细核对验证码（6位数字）

## 🔧 后端API接口

### 1. 发送验证码

```http
POST /api/auth/forgot-password
Content-Type: application/json

{
  "email": "user@example.com",
  "captcha_id": "captcha_id_string",
  "captcha": "captcha_code"
}
```

### 2. 重置密码

```http
POST /api/auth/reset-password
Content-Type: application/json

{
  "email": "user@example.com",
  "code": "123456",
  "new_password": "newpassword123"
}
```

## 📝 数据库清理

建议定期清理过期的令牌（可选）：

```sql
-- 清理30天前的过期令牌
DELETE FROM password_reset_tokens 
WHERE expire_at < NOW() - INTERVAL '30 days';
```

或者在代码中实现定时任务自动清理。

## 💡 提示

- 首次使用前务必先配置好QQ邮箱
- 测试功能前先给自己的邮箱发送一次
- 生产环境建议使用企业邮箱服务（更稳定）
- 可以根据需要调整验证码有效期（目前是15分钟）

