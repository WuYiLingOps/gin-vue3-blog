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

	user, err := h.service.Register(&req)
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

	resp, err := h.service.Login(&req)
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

