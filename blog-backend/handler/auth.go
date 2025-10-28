package handler

import (
	"blog-backend/service"
	"blog-backend/util"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		service: service.NewAuthService(),
	}
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	// 获取客户端IP
	ip := util.GetClientIP(c)

	user, err := h.service.Register(&req, ip)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "注册成功", user)
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	// 获取客户端IP
	ip := util.GetClientIP(c)

	resp, err := h.service.Login(&req, ip)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "登录成功", resp)
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	util.SuccessWithMessage(c, "登出成功", nil)
}

// GetProfile 获取用户信息
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		util.Unauthorized(c, "未登录")
		return
	}

	user, err := h.service.GetProfile(userID.(uint))
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.Success(c, user)
}

// UpdateProfile 更新用户信息
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		util.Unauthorized(c, "未登录")
		return
	}

	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	user, err := h.service.UpdateProfile(userID.(uint), &req)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "更新成功", user)
}

// UpdatePassword 修改密码
func (h *AuthHandler) UpdatePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		util.Unauthorized(c, "未登录")
		return
	}

	var req service.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.service.UpdatePassword(userID.(uint), &req); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "密码修改成功", nil)
}

// RefreshToken 刷新Token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		util.Unauthorized(c, "缺少认证信息")
		return
	}

	// 去掉 "Bearer " 前缀
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	newToken, err := util.RefreshToken(token)
	if err != nil {
		util.Unauthorized(c, "Token刷新失败")
		return
	}

	util.Success(c, gin.H{"token": newToken})
}

// ForgotPassword 忘记密码 - 发送验证码
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req service.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	ip := util.GetClientIP(c)

	if err := h.service.ForgotPassword(&req, ip); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "验证码已发送到您的邮箱，请查收（有效期15分钟）", nil)
}

// ResetPassword 重置密码
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req service.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.service.ResetPassword(&req); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "密码重置成功，请使用新密码登录", nil)
}

// UpdateEmail 修改邮箱
func (h *AuthHandler) UpdateEmail(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		util.Unauthorized(c, "未登录")
		return
	}

	var req service.UpdateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.service.UpdateEmail(userID.(uint), &req); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "邮箱修改成功", nil)
}

// GetEmailChangeInfo 获取邮箱修改信息
func (h *AuthHandler) GetEmailChangeInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		util.Unauthorized(c, "未登录")
		return
	}

	count, err := h.service.GetEmailChangeCount(userID.(uint))
	if err != nil {
		util.ServerError(c, "获取信息失败")
		return
	}

	util.Success(c, gin.H{
		"change_count":     count,
		"remaining_times":  2 - count,
		"can_change":       count < 2,
	})
}

// SendRegisterCode 发送注册验证码
func (h *AuthHandler) SendRegisterCode(c *gin.Context) {
	var req service.SendRegisterCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	ip := util.GetClientIP(c)

	if err := h.service.SendRegisterCode(&req, ip); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "验证码已发送到您的邮箱，请查收（有效期15分钟）", nil)
}

