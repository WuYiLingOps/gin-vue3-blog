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
	subject := "【情迁阁】密码重置验证码"
	body := fmt.Sprintf(`您好！

您正在进行密码重置操作，请使用以下验证码完成验证：

验证码：%s

重要提示：
• 验证码有效期为 15 分钟，请尽快使用
• 请勿将验证码告诉他人，以保护账号安全

如果这不是您本人的操作，请忽略此邮件，您的账号仍然是安全的。

---
此邮件由系统自动发送，请勿直接回复
© 2025 情迁阁`, code)

	return sendEmailText(config, to, subject, body)
}

// SendRegisterVerificationEmail 发送注册验证码邮件
func SendRegisterVerificationEmail(config EmailConfig, to string, code string) error {
	subject := "【情迁阁】注册验证码"
	body := fmt.Sprintf(`您好！

欢迎注册情迁阁，请使用以下验证码完成注册：

验证码：%s

重要提示：
• 验证码有效期为 15 分钟，请尽快使用
• 请勿将验证码告诉他人，以保护账号安全

如果这不是您本人的操作，请忽略此邮件。

---
此邮件由系统自动发送，请勿直接回复
© 2025 情迁阁`, code)

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
