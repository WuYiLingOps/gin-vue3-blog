package handler

import (
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

