/*
 * 项目名称：blog-backend
 * 文件名称：friendlink.go
 * 创建时间：2026-01-31 16:05:15
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：友链管理处理器，提供友链的增删改查功能
 */
package handler

import (
	"strconv"

	"blog-backend/service"
	"blog-backend/util"

	"github.com/gin-gonic/gin"
)

// FriendLinkHandler 友链处理器结构体
type FriendLinkHandler struct {
	service *service.FriendLinkService
}

// NewFriendLinkHandler 创建友链处理器实例
func NewFriendLinkHandler() *FriendLinkHandler {
	return &FriendLinkHandler{
		service: service.NewFriendLinkService(),
	}
}

// Create 创建友链
func (h *FriendLinkHandler) Create(c *gin.Context) {
	var req service.CreateFriendLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	friendLink, err := h.service.Create(&req)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "友链创建成功", friendLink)
}

// GetByID 获取友链详情
func (h *FriendLinkHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的友链ID")
		return
	}

	friendLink, err := h.service.GetByID(uint(id))
	if err != nil {
		util.Error(c, 404, "友链不存在")
		return
	}

	util.Success(c, friendLink)
}

// List 获取友链列表（管理员用）
func (h *FriendLinkHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	friendLinks, total, err := h.service.List(page, pageSize)
	if err != nil {
		util.ServerError(c, "获取友链列表失败")
		return
	}

	util.PageSuccess(c, friendLinks, total, page, pageSize)
}

// ListPublic 获取公开的友链列表（前端用）
func (h *FriendLinkHandler) ListPublic(c *gin.Context) {
	friendLinks, err := h.service.ListPublic()
	if err != nil {
		util.ServerError(c, "获取友链列表失败")
		return
	}

	util.Success(c, friendLinks)
}

// Update 更新友链
func (h *FriendLinkHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的友链ID")
		return
	}

	var req service.UpdateFriendLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	// 调试：打印接收到的请求数据
	if req.CategoryID != nil {
		// 这里可以添加日志输出，但为了不影响生产环境，暂时注释
		// fmt.Printf("接收到 category_id: %d\n", *req.CategoryID)
	}

	friendLink, err := h.service.Update(uint(id), &req)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "友链更新成功", friendLink)
}

// Delete 删除友链
func (h *FriendLinkHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的友链ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "友链删除成功", nil)
}
