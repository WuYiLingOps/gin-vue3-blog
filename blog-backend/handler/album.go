/*
 * 项目名称：blog-backend
 * 文件名称：album.go
 * 创建时间：2026-01-31 16:05:15
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：相册管理处理器，提供相册照片的增删改查功能
 */
package handler

import (
	"strconv"

	"blog-backend/service"
	"blog-backend/util"

	"github.com/gin-gonic/gin"
)

// AlbumHandler 相册处理器结构体
type AlbumHandler struct {
	service *service.AlbumService
}

// NewAlbumHandler 创建相册处理器实例
func NewAlbumHandler() *AlbumHandler {
	return &AlbumHandler{
		service: service.NewAlbumService(),
	}
}

// Create 创建相册照片
func (h *AlbumHandler) Create(c *gin.Context) {
	var req service.CreateAlbumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	album, err := h.service.Create(&req)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "相册照片创建成功", album)
}

// GetByID 获取相册照片详情
func (h *AlbumHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的相册照片ID")
		return
	}

	album, err := h.service.GetByID(uint(id))
	if err != nil {
		util.Error(c, 404, "相册照片不存在")
		return
	}

	util.Success(c, album)
}

// List 获取相册列表（管理员用）
func (h *AlbumHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	albums, total, err := h.service.List(page, pageSize)
	if err != nil {
		util.ServerError(c, "获取相册列表失败")
		return
	}

	util.PageSuccess(c, albums, total, page, pageSize)
}

// ListPublic 获取公开的相册列表（前端用）
func (h *AlbumHandler) ListPublic(c *gin.Context) {
	albums, err := h.service.ListPublic()
	if err != nil {
		util.ServerError(c, "获取相册列表失败")
		return
	}

	util.Success(c, albums)
}

// Update 更新相册照片
func (h *AlbumHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的相册照片ID")
		return
	}

	var req service.UpdateAlbumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	album, err := h.service.Update(uint(id), &req)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "相册照片更新成功", album)
}

// Delete 删除相册照片
func (h *AlbumHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的相册照片ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "相册照片删除成功", nil)
}
