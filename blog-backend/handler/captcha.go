package handler

import (
	"blog-backend/util"

	"github.com/gin-gonic/gin"
)

type CaptchaHandler struct{}

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

