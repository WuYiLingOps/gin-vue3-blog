package handler

import (
	"blog-backend/config"
	"blog-backend/service"
	"blog-backend/util"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	service *service.SettingService
}

func NewSettingHandler() *SettingHandler {
	return &SettingHandler{
		service: service.NewSettingService(),
	}
}

// GetSiteSettings 获取网站配置（仅管理员）
func (h *SettingHandler) GetSiteSettings(c *gin.Context) {
	settings, err := h.service.GetSiteSettings()
	if err != nil {
		util.Error(c, 500, "获取配置失败")
		return
	}

	util.Success(c, settings)
}

// UpdateSiteSettings 更新网站配置（仅管理员）
func (h *SettingHandler) UpdateSiteSettings(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "参数错误")
		return
	}

	if err := h.service.UpdateSiteSettings(req); err != nil {
		util.Error(c, 500, "更新配置失败")
		return
	}

	util.SuccessWithMessage(c, "更新成功", nil)
}

// GetPublicSettings 获取公开的网站配置（前端用，无需认证）
func (h *SettingHandler) GetPublicSettings(c *gin.Context) {
	settings, err := h.service.GetPublicSettings()
	if err != nil {
		util.Error(c, 500, "获取配置失败")
		return
	}

	util.Success(c, settings)
}

// GetUploadSettings 获取上传配置（仅管理员）
func (h *SettingHandler) GetUploadSettings(c *gin.Context) {
	settings, err := h.service.GetUploadSettings()
	if err != nil {
		util.Error(c, 500, "获取配置失败")
		return
	}

	util.Success(c, settings)
}

// UpdateUploadSettings 更新上传配置（仅管理员，只能切换存储类型）
func (h *SettingHandler) UpdateUploadSettings(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "参数错误")
		return
	}

	// 只允许修改 storage_type
	storageType := req["storage_type"]
	if storageType != "local" && storageType != "oss" && storageType != "cos" {
		util.BadRequest(c, "存储类型只能是 local、oss 或 cos")
		return
	}

	// 如果选择 OSS 存储，检查配置文件中的 OSS 配置是否完整
	if storageType == "oss" {
		if err := util.ValidateOSSConfig(
			config.Cfg.OSS.Endpoint,
			config.Cfg.OSS.AccessKeyID,
			config.Cfg.OSS.AccessKeySecret,
			config.Cfg.OSS.BucketName,
		); err != nil {
			util.BadRequest(c, "OSS 配置不完整，请先在配置文件中设置 OSS 参数")
			return
		}
	}

	// 如果选择 COS 存储，检查配置文件中的 COS 配置是否完整
	if storageType == "cos" {
		if err := util.ValidateCOSConfig(
			config.Cfg.COS.BucketURL,
			config.Cfg.COS.SecretID,
			config.Cfg.COS.SecretKey,
		); err != nil {
			util.BadRequest(c, "COS 配置不完整，请先在配置文件中设置 COS 参数")
			return
		}
	}

	if err := h.service.UpdateUploadSettings(req); err != nil {
		util.Error(c, 500, "更新配置失败")
		return
	}

	util.SuccessWithMessage(c, "更新成功", nil)
}

// GetFriendLinkInfo 获取我的友链信息（公开接口）
func (h *SettingHandler) GetFriendLinkInfo(c *gin.Context) {
	info, err := h.service.GetFriendLinkInfo()
	if err != nil {
		util.Error(c, 500, "获取友链信息失败")
		return
	}

	util.Success(c, info)
}

// UpdateFriendLinkInfo 更新我的友链信息（仅管理员）
func (h *SettingHandler) UpdateFriendLinkInfo(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "参数错误")
		return
	}

	if err := h.service.UpdateFriendLinkInfo(req); err != nil {
		util.Error(c, 500, "更新友链信息失败")
		return
	}

	util.SuccessWithMessage(c, "更新成功", nil)
}

// GetNotificationSettings 获取通知配置（仅管理员）
func (h *SettingHandler) GetNotificationSettings(c *gin.Context) {
	settings, err := h.service.GetNotificationSettings()
	if err != nil {
		util.Error(c, 500, "获取通知配置失败")
		return
	}

	util.Success(c, settings)
}

// UpdateNotificationSettings 更新通知配置（仅管理员）
func (h *SettingHandler) UpdateNotificationSettings(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "参数错误")
		return
	}

	if err := h.service.UpdateNotificationSettings(req); err != nil {
		util.Error(c, 500, "更新通知配置失败")
		return
	}

	util.SuccessWithMessage(c, "更新成功", nil)
}

// GetRegisterSettings 获取注册配置（仅管理员）
func (h *SettingHandler) GetRegisterSettings(c *gin.Context) {
	settings, err := h.service.GetRegisterSettings()
	if err != nil {
		util.Error(c, 500, "获取注册配置失败")
		return
	}

	util.Success(c, settings)
}

// UpdateRegisterSettings 更新注册配置（仅管理员）
func (h *SettingHandler) UpdateRegisterSettings(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "参数错误")
		return
	}

	if err := h.service.UpdateRegisterSettings(req); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "更新成功", nil)
}