package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/smtp"
)

// EmailConfig 邮件配置
type EmailConfig struct {
	Host     string // SMTP服务器地址
	Port     int    // 端口
	Username string // 发件人邮箱
	Password string // 授权码
	FromName string // 发件人名称
}

// GenerateVerificationCode 生成6位数字验证码
func GenerateVerificationCode() string {
	code := ""
	for i := 0; i < 6; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		code += fmt.Sprintf("%d", n)
	}
	return code
}

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[n.Int64()]
	}
	return string(result)
}

// SendResetPasswordEmail 发送重置密码邮件
func SendResetPasswordEmail(config EmailConfig, to string, code string) error {
	siteName := config.FromName
	if siteName == "" {
		siteName = "無以菱"
	}
	subject := fmt.Sprintf("【%s】密码重置验证码", siteName)
	body := fmt.Sprintf(`您好！

您正在进行密码重置操作，请使用以下验证码完成验证：

验证码：%s

重要提示：
• 验证码有效期为 15 分钟，请尽快使用
• 请勿将验证码告诉他人，以保护账号安全

如果这不是您本人的操作，请忽略此邮件，您的账号仍然是安全的。

---
此邮件由系统自动发送，请勿直接回复
© 2025 %s`, code, siteName)

	return sendEmailText(config, to, subject, body)
}

// SendRegisterVerificationEmail 发送注册验证码邮件
func SendRegisterVerificationEmail(config EmailConfig, to string, code string) error {
	siteName := config.FromName
	if siteName == "" {
		siteName = "無以菱"
	}
	subject := fmt.Sprintf("【%s】注册验证码", siteName)
	body := fmt.Sprintf(`您好！

欢迎注册%s，请使用以下验证码完成注册：

验证码：%s

重要提示：
• 验证码有效期为 15 分钟，请尽快使用
• 请勿将验证码告诉他人，以保护账号安全

如果这不是您本人的操作，请忽略此邮件。

---
此邮件由系统自动发送，请勿直接回复
© 2025 %s`, siteName, code, siteName)

	return sendEmailText(config, to, subject, body)
}

// SendCommentNotificationEmail 发送评论通知邮件（给文章作者）
func SendCommentNotificationEmail(config EmailConfig, to string, commenterName string, postTitle string, commentContent string, postURL string) error {
	siteName := config.FromName
	if siteName == "" {
		siteName = "無以菱"
	}
	subject := fmt.Sprintf("【%s】您的文章收到新评论", siteName)

	// 截取评论内容前100个字符作为预览
	preview := commentContent
	if len([]rune(commentContent)) > 100 {
		preview = string([]rune(commentContent)[:100]) + "..."
	}

	body := fmt.Sprintf(`您好！

您的文章《%s》收到了一条新评论：

评论者：%s
评论内容：%s

查看完整评论：%s

---
此邮件由系统自动发送，请勿直接回复
© 2025 %s`, postTitle, commenterName, preview, postURL, siteName)

	return sendEmailText(config, to, subject, body)
}

// SendAdminCommentNotificationEmail 发送评论通知邮件（给管理员）
func SendAdminCommentNotificationEmail(config EmailConfig, to string, commenterName string, postTitle string, commentContent string, postURL string) error {
	siteName := config.FromName
	if siteName == "" {
		siteName = "無以菱"
	}
	subject := fmt.Sprintf("【%s】系统收到新评论", siteName)

	// 截取评论内容前100个字符作为预览
	preview := commentContent
	if len([]rune(commentContent)) > 100 {
		preview = string([]rune(commentContent)[:100]) + "..."
	}

	body := fmt.Sprintf(`管理员您好！

系统收到了一条新评论：

文章：%s
评论者：%s
评论内容：%s

查看完整评论：%s

---
此邮件由系统自动发送，请勿直接回复
© 2025 %s`, postTitle, commenterName, preview, postURL, siteName)

	return sendEmailText(config, to, subject, body)
}

// sendEmailText 发送纯文本邮件
func sendEmailText(config EmailConfig, to, subject, body string) error {
	// SMTP认证
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	// 构建邮件内容
	contentType := "Content-Type: text/plain; charset=UTF-8"
	msg := []byte(fmt.Sprintf(
		"To: %s\r\nFrom: %s<%s>\r\nSubject: %s\r\n%s\r\n\r\n%s",
		to, config.FromName, config.Username, subject, contentType, body,
	))

	// 发送邮件
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	return smtp.SendMail(addr, auth, config.Username, []string{to}, msg)
}
