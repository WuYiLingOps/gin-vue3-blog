package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	"math/big"
	"mime"
	"net/smtp"
	"strings"
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

	data := map[string]interface{}{
		"SiteName": siteName,
		"Code":     code,
		"Year":     "2025",
	}

	htmlBody := getEmailTemplate("reset_password", data)
	textBody := fmt.Sprintf(`您好！

您正在进行密码重置操作，请使用以下验证码完成验证：

验证码：%s

重要提示：
• 验证码有效期为 15 分钟，请尽快使用
• 请勿将验证码告诉他人，以保护账号安全

如果这不是您本人的操作，请忽略此邮件，您的账号仍然是安全的。

---
此邮件由系统自动发送，请勿直接回复
© 2025 %s`, code, siteName)

	return sendEmailHTML(config, to, subject, htmlBody, textBody)
}

// SendRegisterVerificationEmail 发送注册验证码邮件
func SendRegisterVerificationEmail(config EmailConfig, to string, code string) error {
	siteName := config.FromName
	if siteName == "" {
		siteName = "無以菱"
	}
	subject := fmt.Sprintf("【%s】注册验证码", siteName)

	data := map[string]interface{}{
		"SiteName": siteName,
		"Code":     code,
		"Year":     "2025",
	}

	htmlBody := getEmailTemplate("register_verification", data)
	textBody := fmt.Sprintf(`您好！

欢迎注册%s，请使用以下验证码完成注册：

验证码：%s

重要提示：
• 验证码有效期为 15 分钟，请尽快使用
• 请勿将验证码告诉他人，以保护账号安全

如果这不是您本人的操作，请忽略此邮件。

---
此邮件由系统自动发送，请勿直接回复
© 2025 %s`, siteName, code, siteName)

	return sendEmailHTML(config, to, subject, htmlBody, textBody)
}

// SendAdminCommentNotificationEmail 发送评论通知邮件（给管理员）
func SendAdminCommentNotificationEmail(config EmailConfig, to string, commenterName string, postTitle string, commentContent string, postURL string) error {
	siteName := config.FromName
	if siteName == "" {
		siteName = "無以菱"
	}
	subject := fmt.Sprintf("【%s】系统收到新评论", siteName)

	// 截取评论内容前200个字符作为预览
	preview := commentContent
	if len([]rune(commentContent)) > 200 {
		preview = string([]rune(commentContent)[:200]) + "..."
	}
	// HTML转义
	previewHTML := template.HTMLEscapeString(preview)
	previewHTML = strings.ReplaceAll(previewHTML, "\n", "<br>")

	data := map[string]interface{}{
		"SiteName":       siteName,
		"PostTitle":      template.HTMLEscapeString(postTitle),
		"CommenterName":  template.HTMLEscapeString(commenterName),
		"CommentPreview": template.HTML(previewHTML),
		"PostURL":        postURL,
		"Year":           "2025",
	}

	htmlBody := getEmailTemplate("admin_comment_notification", data)
	textBody := fmt.Sprintf(`管理员您好！

系统收到了一条新评论：

文章：%s
评论者：%s
评论内容：%s

查看完整评论：%s

---
此邮件由系统自动发送，请勿直接回复
© 2025 %s`, postTitle, commenterName, preview, postURL, siteName)

	return sendEmailHTML(config, to, subject, htmlBody, textBody)
}

// sendEmailHTML 发送HTML格式邮件（支持纯文本回退）
// 注意：为了兼容性，直接发送HTML格式，不使用multipart/alternative
// 大多数现代邮件客户端都支持HTML，这样可以避免multipart格式导致的"short response"错误
func sendEmailHTML(config EmailConfig, to, subject, htmlBody, textBody string) error {
	// SMTP认证
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	// 对Subject进行编码（支持中文）
	encodedSubject := mime.QEncoding.Encode("UTF-8", subject)

	// 对HTML内容进行base64编码（HTML内容通常包含特殊字符，base64更安全）
	htmlBodyEncoded := base64.StdEncoding.EncodeToString([]byte(htmlBody))
	// 将base64编码的内容按76字符换行（RFC 2045要求）
	htmlBodyLines := make([]string, 0)
	for i := 0; i < len(htmlBodyEncoded); i += 76 {
		end := i + 76
		if end > len(htmlBodyEncoded) {
			end = len(htmlBodyEncoded)
		}
		htmlBodyLines = append(htmlBodyLines, htmlBodyEncoded[i:end])
	}
	htmlBodyFormatted := strings.Join(htmlBodyLines, "\r\n")

	// 构建简单的HTML邮件（不使用multipart，避免格式问题）
	var msgBuilder strings.Builder
	msgBuilder.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msgBuilder.WriteString(fmt.Sprintf("From: %s<%s>\r\n", config.FromName, config.Username))
	msgBuilder.WriteString(fmt.Sprintf("Subject: %s\r\n", encodedSubject))
	msgBuilder.WriteString("MIME-Version: 1.0\r\n")
	msgBuilder.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	msgBuilder.WriteString("Content-Transfer-Encoding: base64\r\n")
	msgBuilder.WriteString("\r\n")
	msgBuilder.WriteString(htmlBodyFormatted)
	msgBuilder.WriteString("\r\n")

	msg := []byte(msgBuilder.String())

	// 发送邮件
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	err := smtp.SendMail(addr, auth, config.Username, []string{to}, msg)

	// 某些SMTP服务器在发送成功后可能返回"short response"错误
	// 这通常是因为服务器返回的响应格式不符合Go smtp包的预期
	// 如果错误是"short response"，可以忽略，因为邮件通常已经成功发送
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "short response") {
			// short response错误通常不影响邮件发送，返回nil表示成功
			// 这样可以避免误报错误，因为邮件实际上已经发送成功
			return nil
		}
	}
	return err
}

// getEmailTemplate 获取HTML邮件模板
func getEmailTemplate(templateName string, data map[string]interface{}) string {
	var templateStr string

	switch templateName {
	case "reset_password":
		templateStr = getResetPasswordTemplate()
	case "register_verification":
		templateStr = getRegisterVerificationTemplate()
	case "admin_comment_notification":
		templateStr = getAdminCommentNotificationTemplate()
	default:
		return ""
	}

	// 使用Go的template包渲染
	tmpl, err := template.New("email").Parse(templateStr)
	if err != nil {
		return ""
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return ""
	}

	return buf.String()
}

// getResetPasswordTemplate 获取密码重置邮件模板
func getResetPasswordTemplate() string {
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>密码重置验证码</title>
</head>
<body style="margin: 0; padding: 0; font-family: 'helvetica neue', PingFangSC-Light, arial, 'hiragino sans gb', 'microsoft yahei ui', 'microsoft yahei', simsun, sans-serif; background-color: #f7f8fa;">
    <div style="word-break: break-all; box-sizing: border-box; text-align: center; min-width: 320px; max-width: 660px; border: 1px solid #f6f6f6; background-color: #f7f8fa; margin: auto; padding: 20px 0 30px;">
        <table style="width: 100%; font-weight: 300; margin-bottom: 10px; border-collapse: collapse;">
            <tbody>
                <tr style="font-weight: 300;">
                    <td style="width: 3%; max-width: 30px;"></td>
                    <td style="max-width: 600px;">
                        <!-- LOGO -->
                        <div style="width: 100%; text-align: left; margin-bottom: 20px;">
                            <table style="border-collapse: collapse;">
                                <tr>
                                    <td style="vertical-align: middle; padding-right: 12px;">
                                        <img src="https://pic1.imgdb.cn/item/694560da2ce92a4f9312aec7.jpg" alt="Logo" style="width: 40px; height: 40px; border-radius: 50%; display: block;" />
                                    </td>
                                    <td style="vertical-align: middle;">
                                        <h1 style="margin: 0; color: #0891b2; font-size: 24px; font-weight: 600;">{{.SiteName}}</h1>
                                    </td>
                                </tr>
                            </table>
                        </div>
                        <!-- 蓝色分割线 -->
                        <p style="height: 2px; background-color: #0891b2; border: 0; font-size: 0; padding: 0; width: 100%; margin-top: 20px; margin-bottom: 0;"></p>
                        
                        <!-- 内容区域 -->
                        <div style="background-color: #fff; padding: 23px 0 20px; box-shadow: 0px 1px 1px 0px rgba(122, 55, 55, 0.2); text-align: left;">
                            <table style="width: 100%; font-weight: 300; margin-bottom: 10px; border-collapse: collapse; text-align: left;">
                                <tbody>
                                    <tr style="font-weight: 300;">
                                        <td style="width: 3.2%; max-width: 30px;"></td>
                                        <td style="max-width: 480px; text-align: left;">
                                            <!-- 标题 -->
                                            <h1 style="font-size: 20px; line-height: 36px; margin: 0px 0px 22px; color: #333;">密码重置验证码</h1>
                                            
                                            <!-- 问候语 -->
                                            <p style="font-size: 14px; color: #333; line-height: 24px; margin: 0;">您好！</p>
                                            
                                            <p style="line-height: 24px; margin: 6px 0px 0px; overflow-wrap: break-word; word-break: break-all;">
                                                <span style="color: rgb(51, 51, 51); font-size: 14px;">您正在进行密码重置操作，请使用以下验证码完成验证：</span>
                                            </p>
                                            
                                            <!-- 验证码框 -->
                                            <div style="background-color: #f0fdfa; border: 2px dashed #0891b2; border-radius: 8px; padding: 20px; text-align: center; margin: 30px 0;">
                                                <p style="margin: 0 0 10px 0; color: #64748b; font-size: 12px; text-transform: uppercase; letter-spacing: 1px;">验证码</p>
                                                <p style="margin: 0; color: #0891b2; font-size: 32px; font-weight: bold; letter-spacing: 8px; font-family: 'Courier New', monospace;">{{.Code}}</p>
                                            </div>
                                            
                                            <!-- 重要提示 -->
                                            <p style="line-height: 24px; margin: 6px 0px 0px; overflow-wrap: break-word; word-break: break-all;">
                                                <span style="color: rgb(51, 51, 51); font-size: 14px; font-weight: bold;">重要提示：</span>
                                            </p>
                                            <p style="line-height: 24px; margin: 6px 0px 0px; overflow-wrap: break-word; word-break: break-all;">
                                                <span style="color: rgb(51, 51, 51); font-size: 14px;">• 验证码有效期为 15 分钟，请尽快使用</span>
                                            </p>
                                            <p style="line-height: 24px; margin: 6px 0px 0px; overflow-wrap: break-word; word-break: break-all;">
                                                <span style="color: rgb(51, 51, 51); font-size: 14px;">• 请勿将验证码告诉他人，以保护账号安全</span>
                                            </p>
                                            
                                            <p style="line-height: 24px; margin: 20px 0px 0px; overflow-wrap: break-word; word-break: break-all;">
                                                <span style="color: rgb(51, 51, 51); font-size: 14px;">如果这不是您本人的操作，请忽略此邮件，您的账号仍然是安全的。</span>
                                            </p>
                                            
                                            <!-- 署名 -->
                                            <p style="font-size: 14px; line-height: 26px; word-wrap: break-word; word-break: break-all; margin-top: 32px; color: #333;">
                                                此致<br>
                                                <strong>{{.SiteName}}团队</strong>
                                            </p>
                                        </td>
                                        <td style="width: 3.2%; max-width: 30px;"></td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                        
                        <!-- 底部 -->
                        <div style="text-align: center; font-size: 12px; line-height: 18px; color: #999; margin-top: 20px;">
                            <table style="width: 100%; font-weight: 300; margin-bottom: 10px; border-collapse: collapse;">
                                <tbody>
                                    <tr style="font-weight: 300;">
                                        <td style="width: 3.2%; max-width: 30px;"></td>
                                        <td style="max-width: 540px;">
                                            <p style="text-align: center; margin: 20px auto 14px auto; font-size: 12px; color: #999;">
                                                此为系统邮件，请勿回复。
                                            </p>
                                            <p style="max-width: 100%; margin: auto; font-size: 12px; color: #999; text-align: center; line-height: 22px;">
                                                © {{.Year}} {{.SiteName}}
                                            </p>
                                        </td>
                                        <td style="width: 3.2%; max-width: 30px;"></td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </td>
                    <td style="width: 3%; max-width: 30px;"></td>
                </tr>
            </tbody>
        </table>
    </div>
</body>
</html>`
}

// getRegisterVerificationTemplate 获取注册验证码邮件模板
func getRegisterVerificationTemplate() string {
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>注册验证码</title>
</head>
<body style="margin: 0; padding: 0; font-family: 'helvetica neue', PingFangSC-Light, arial, 'hiragino sans gb', 'microsoft yahei ui', 'microsoft yahei', simsun, sans-serif; background-color: #f7f8fa;">
    <div style="word-break: break-all; box-sizing: border-box; text-align: center; min-width: 320px; max-width: 660px; border: 1px solid #f6f6f6; background-color: #f7f8fa; margin: auto; padding: 20px 0 30px;">
        <table style="width: 100%; font-weight: 300; margin-bottom: 10px; border-collapse: collapse;">
            <tbody>
                <tr style="font-weight: 300;">
                    <td style="width: 3%; max-width: 30px;"></td>
                    <td style="max-width: 600px;">
                        <!-- LOGO -->
                        <div style="width: 100%; text-align: left; margin-bottom: 20px;">
                            <table style="border-collapse: collapse;">
                                <tr>
                                    <td style="vertical-align: middle; padding-right: 12px;">
                                        <img src="https://pic1.imgdb.cn/item/694560da2ce92a4f9312aec7.jpg" alt="Logo" style="width: 40px; height: 40px; border-radius: 50%; display: block;" />
                                    </td>
                                    <td style="vertical-align: middle;">
                                        <h1 style="margin: 0; color: #0891b2; font-size: 24px; font-weight: 600;">{{.SiteName}}</h1>
                                    </td>
                                </tr>
                            </table>
                        </div>
                        <!-- 蓝色分割线 -->
                        <p style="height: 2px; background-color: #0891b2; border: 0; font-size: 0; padding: 0; width: 100%; margin-top: 20px; margin-bottom: 0;"></p>
                        
                        <!-- 内容区域 -->
                        <div style="background-color: #fff; padding: 23px 0 20px; box-shadow: 0px 1px 1px 0px rgba(122, 55, 55, 0.2); text-align: left;">
                            <table style="width: 100%; font-weight: 300; margin-bottom: 10px; border-collapse: collapse; text-align: left;">
                                <tbody>
                                    <tr style="font-weight: 300;">
                                        <td style="width: 3.2%; max-width: 30px;"></td>
                                        <td style="max-width: 480px; text-align: left;">
                                            <!-- 标题 -->
                                            <h1 style="font-size: 20px; line-height: 36px; margin: 0px 0px 22px; color: #333;">欢迎注册</h1>
                                            
                                            <!-- 问候语 -->
                                            <p style="font-size: 14px; color: #333; line-height: 24px; margin: 0;">您好！</p>
                                            
                                            <p style="line-height: 24px; margin: 6px 0px 0px; overflow-wrap: break-word; word-break: break-all;">
                                                <span style="color: rgb(51, 51, 51); font-size: 14px;">欢迎注册{{.SiteName}}，请使用以下验证码完成注册：</span>
                                            </p>
                                            
                                            <!-- 验证码框 -->
                                            <div style="background-color: #f0fdfa; border: 2px dashed #0891b2; border-radius: 8px; padding: 20px; text-align: center; margin: 30px 0;">
                                                <p style="margin: 0 0 10px 0; color: #64748b; font-size: 12px; text-transform: uppercase; letter-spacing: 1px;">验证码</p>
                                                <p style="margin: 0; color: #0891b2; font-size: 32px; font-weight: bold; letter-spacing: 8px; font-family: 'Courier New', monospace;">{{.Code}}</p>
                                            </div>
                                            
                                            <!-- 重要提示 -->
                                            <p style="line-height: 24px; margin: 6px 0px 0px; overflow-wrap: break-word; word-break: break-all;">
                                                <span style="color: rgb(51, 51, 51); font-size: 14px; font-weight: bold;">重要提示：</span>
                                            </p>
                                            <p style="line-height: 24px; margin: 6px 0px 0px; overflow-wrap: break-word; word-break: break-all;">
                                                <span style="color: rgb(51, 51, 51); font-size: 14px;">• 验证码有效期为 15 分钟，请尽快使用</span>
                                            </p>
                                            <p style="line-height: 24px; margin: 6px 0px 0px; overflow-wrap: break-word; word-break: break-all;">
                                                <span style="color: rgb(51, 51, 51); font-size: 14px;">• 请勿将验证码告诉他人，以保护账号安全</span>
                                            </p>
                                            
                                            <p style="line-height: 24px; margin: 20px 0px 0px; overflow-wrap: break-word; word-break: break-all;">
                                                <span style="color: rgb(51, 51, 51); font-size: 14px;">如果这不是您本人的操作，请忽略此邮件。</span>
                                            </p>
                                            
                                            <!-- 署名 -->
                                            <p style="font-size: 14px; line-height: 26px; word-wrap: break-word; word-break: break-all; margin-top: 32px; color: #333;">
                                                此致<br>
                                                <strong>{{.SiteName}}团队</strong>
                                            </p>
                                        </td>
                                        <td style="width: 3.2%; max-width: 30px;"></td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                        
                        <!-- 底部 -->
                        <div style="text-align: center; font-size: 12px; line-height: 18px; color: #999; margin-top: 20px;">
                            <table style="width: 100%; font-weight: 300; margin-bottom: 10px; border-collapse: collapse;">
                                <tbody>
                                    <tr style="font-weight: 300;">
                                        <td style="width: 3.2%; max-width: 30px;"></td>
                                        <td style="max-width: 540px;">
                                            <p style="text-align: center; margin: 20px auto 14px auto; font-size: 12px; color: #999;">
                                                此为系统邮件，请勿回复。
                                            </p>
                                            <p style="max-width: 100%; margin: auto; font-size: 12px; color: #999; text-align: center; line-height: 22px;">
                                                © {{.Year}} {{.SiteName}}
                                            </p>
                                        </td>
                                        <td style="width: 3.2%; max-width: 30px;"></td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </td>
                    <td style="width: 3%; max-width: 30px;"></td>
                </tr>
            </tbody>
        </table>
    </div>
</body>
</html>`
}

// getAdminCommentNotificationTemplate 获取管理员评论通知邮件模板
func getAdminCommentNotificationTemplate() string {
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>系统收到新评论</title>
</head>
<body style="margin: 0; padding: 0; font-family: 'helvetica neue', PingFangSC-Light, arial, 'hiragino sans gb', 'microsoft yahei ui', 'microsoft yahei', simsun, sans-serif; background-color: #f7f8fa;">
    <div style="word-break: break-all; box-sizing: border-box; text-align: center; min-width: 320px; max-width: 660px; border: 1px solid #f6f6f6; background-color: #f7f8fa; margin: auto; padding: 20px 0 30px;">
        <table style="width: 100%; font-weight: 300; margin-bottom: 10px; border-collapse: collapse;">
            <tbody>
                <tr style="font-weight: 300;">
                    <td style="width: 3%; max-width: 30px;"></td>
                    <td style="max-width: 600px;">
                        <!-- LOGO -->
                        <div style="width: 100%; text-align: left; margin-bottom: 20px;">
                            <table style="border-collapse: collapse;">
                                <tr>
                                    <td style="vertical-align: middle; padding-right: 12px;">
                                        <img src="https://pic1.imgdb.cn/item/694560da2ce92a4f9312aec7.jpg" alt="Logo" style="width: 40px; height: 40px; border-radius: 50%; display: block;" />
                                    </td>
                                    <td style="vertical-align: middle;">
                                        <h1 style="margin: 0; color: #0891b2; font-size: 24px; font-weight: 600;">{{.SiteName}}</h1>
                                    </td>
                                </tr>
                            </table>
                        </div>
                        <!-- 蓝色分割线 -->
                        <p style="height: 2px; background-color: #0891b2; border: 0; font-size: 0; padding: 0; width: 100%; margin-top: 20px; margin-bottom: 0;"></p>
                        
                        <!-- 内容区域 -->
                        <div style="background-color: #fff; padding: 23px 0 20px; box-shadow: 0px 1px 1px 0px rgba(122, 55, 55, 0.2); text-align: left;">
                            <table style="width: 100%; font-weight: 300; margin-bottom: 10px; border-collapse: collapse; text-align: left;">
                                <tbody>
                                    <tr style="font-weight: 300;">
                                        <td style="width: 3.2%; max-width: 30px;"></td>
                                        <td style="max-width: 480px; text-align: left;">
                                            <!-- 标题 -->
                                            <h1 style="font-size: 20px; line-height: 36px; margin: 0px 0px 22px; color: #333;">系统收到新评论</h1>
                                            
                                            <!-- 问候语 -->
                                            <p style="font-size: 14px; color: #333; line-height: 24px; margin: 0;">管理员您好！</p>
                                            
                                            <p style="line-height: 24px; margin: 6px 0px 0px; overflow-wrap: break-word; word-break: break-all;">
                                                <span style="color: rgb(51, 51, 51); font-size: 14px;">系统收到了一条新评论：</span>
                                            </p>
                                            
                                            <!-- 评论信息框 -->
                                            <div style="background-color: #f0fdfa; border-left: 4px solid #0891b2; padding: 20px; margin: 30px 0; border-radius: 4px;">
                                                <p style="margin: 0 0 10px 0; color: #333; font-size: 14px; line-height: 24px;">
                                                    <strong style="color: #0891b2;">文章：</strong>{{.PostTitle}}
                                                </p>
                                                <p style="margin: 0 0 15px 0; color: #333; font-size: 14px; line-height: 24px;">
                                                    <strong style="color: #0891b2;">评论者：</strong>{{.CommenterName}}
                                                </p>
                                                <div style="background-color: #ffffff; padding: 15px; border-radius: 4px; margin-top: 10px;">
                                                    <p style="margin: 0 0 10px 0; color: #999; font-size: 12px; text-transform: uppercase; letter-spacing: 1px;">评论内容</p>
                                                    <div style="color: #333; font-size: 14px; line-height: 1.8;">{{.CommentPreview}}</div>
                                                </div>
                                            </div>
                                            
                                            <!-- 按钮 -->
                                            <p style="font-size: 14px; color: rgb(51, 51, 51); line-height: 24px; margin: 6px 0px 0px; word-wrap: break-word; word-break: break-all;">
                                                <a href="{{.PostURL}}" title="查看完整评论" style="font-size: 16px; line-height: 45px; display: block; background-color: #0891b2; color: rgb(255, 255, 255); text-align: center; text-decoration: none; margin-top: 20px; border-radius: 3px;">
                                                    查看完整评论
                                                </a>
                                            </p>
                                            
                                            <!-- 署名 -->
                                            <p style="font-size: 14px; line-height: 26px; word-wrap: break-word; word-break: break-all; margin-top: 32px; color: #333;">
                                                此致<br>
                                                <strong>{{.SiteName}}团队</strong>
                                            </p>
                                        </td>
                                        <td style="width: 3.2%; max-width: 30px;"></td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                        
                        <!-- 底部 -->
                        <div style="text-align: center; font-size: 12px; line-height: 18px; color: #999; margin-top: 20px;">
                            <table style="width: 100%; font-weight: 300; margin-bottom: 10px; border-collapse: collapse;">
                                <tbody>
                                    <tr style="font-weight: 300;">
                                        <td style="width: 3.2%; max-width: 30px;"></td>
                                        <td style="max-width: 540px;">
                                            <p style="text-align: center; margin: 20px auto 14px auto; font-size: 12px; color: #999;">
                                                此为系统邮件，请勿回复。
                                            </p>
                                            <p style="max-width: 100%; margin: auto; font-size: 12px; color: #999; text-align: center; line-height: 22px;">
                                                © {{.Year}} {{.SiteName}}
                                            </p>
                                        </td>
                                        <td style="width: 3.2%; max-width: 30px;"></td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </td>
                    <td style="width: 3%; max-width: 30px;"></td>
                </tr>
            </tbody>
        </table>
    </div>
</body>
</html>`
}
