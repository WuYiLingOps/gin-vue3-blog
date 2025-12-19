package handler

import (
	"blog-backend/util"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct{}

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
	role, exists := c.Get("role")
	if exists && role == "admin" {
		// 管理员使用配置的存储方式（本地/OSS/COS）
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
	role, exists := c.Get("role")
	if exists && role == "admin" {
		// 管理员使用配置的存储方式（本地/OSS/COS）
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
