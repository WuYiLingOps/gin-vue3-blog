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

	// 保存文件
	filePath, err := util.SaveUploadedFile(file, util.AvatarDir)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	// 返回文件 URL
	fileURL := util.GetFileURL(filePath)
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

	// 保存文件
	filePath, err := util.SaveUploadedFile(file, util.UploadDir)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	// 返回文件 URL
	fileURL := util.GetFileURL(filePath)
	util.SuccessWithMessage(c, "上传成功", gin.H{
		"url": fileURL,
	})
}

