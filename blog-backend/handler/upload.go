/*
 * 项目名称：blog-backend
 * 文件名称：upload.go
 * 创建时间：2026-01-31 16:05:15
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：文件上传处理器，提供头像和图片上传功能，支持本地存储和云存储（OSS/COS）
 */
package handler

import (
	"blog-backend/constant"
	"blog-backend/util"

	"github.com/gin-gonic/gin"
)

// UploadHandler 文件上传处理器结构体
type UploadHandler struct{}

// NewUploadHandler 创建文件上传处理器实例
func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

// UploadAvatar 上传头像
func (h *UploadHandler) UploadAvatar(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		util.BadRequest(c, "请选择要上传的文件")
		return
	}

	var fileURL string
	// 获取用户角色，普通用户强制使用本地存储
	roleVal, exists := c.Get("role")
	role, _ := roleVal.(string)
	if exists && constant.IsAdminRole(role) {
		// 具备管理员权限的用户使用配置的存储方式（本地/OSS/COS）
		fileURL, err = util.UploadFile(file, util.AvatarDir)
		if err != nil {
			util.Error(c, 400, err.Error())
			return
		}
	} else {
		// 普通用户强制使用本地存储
		filePath, err := util.SaveUploadedFile(file, util.AvatarDir)
		if err != nil {
			util.Error(c, 400, err.Error())
			return
		}
		fileURL = util.GetFileURL(filePath)
	}

	// 返回文件 URL
	util.SuccessWithMessage(c, "上传成功", gin.H{
		"url": fileURL,
	})
}

// UploadImage 上传图片（通用）
func (h *UploadHandler) UploadImage(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		util.BadRequest(c, "请选择要上传的文件")
		return
	}

	var fileURL string
	// 获取用户角色，普通用户强制使用本地存储
	roleVal, exists := c.Get("role")
	role, _ := roleVal.(string)
	if exists && constant.IsAdminRole(role) {
		// 具备管理员权限的用户使用配置的存储方式（本地/OSS/COS）
		fileURL, err = util.UploadFile(file, util.UploadDir)
		if err != nil {
			util.Error(c, 400, err.Error())
			return
		}
	} else {
		// 普通用户强制使用本地存储
		filePath, err := util.SaveUploadedFile(file, util.UploadDir)
		if err != nil {
			util.Error(c, 400, err.Error())
			return
		}
		fileURL = util.GetFileURL(filePath)
	}

	// 返回文件 URL
	util.SuccessWithMessage(c, "上传成功", gin.H{
		"url": fileURL,
	})
}
