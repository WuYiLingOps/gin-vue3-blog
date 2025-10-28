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
	captcha, err := util.GenerateCaptcha()
	if err != nil {
		util.Error(c, 500, "生成验证码失败")
		return
	}

	util.Success(c, captcha)
}

