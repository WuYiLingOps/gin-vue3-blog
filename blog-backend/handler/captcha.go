/*
 * 项目名称：blog-backend
 * 文件名称：captcha.go
 * 创建时间：2026-01-31 16:05:15
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：验证码处理器，提供图形验证码生成功能
 */
package handler

import (
	"blog-backend/util"

	"github.com/gin-gonic/gin"
)

// CaptchaHandler 验证码处理器结构体
type CaptchaHandler struct{}

// NewCaptchaHandler 创建验证码处理器实例
func NewCaptchaHandler() *CaptchaHandler {
	return &CaptchaHandler{}
}

// GetCaptcha 获取验证码
func (h *CaptchaHandler) GetCaptcha(c *gin.Context) {
	// 获取客户端IP
	ip := util.GetClientIP(c)

	captcha, err := util.GenerateCaptcha(ip)
	if err != nil {
		util.Error(c, 429, err.Error())
		return
	}

	util.Success(c, captcha)
}

