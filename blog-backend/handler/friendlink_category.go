package handler

import (
	"strconv"

	"blog-backend/service"
	"blog-backend/util"

	"github.com/gin-gonic/gin"
)

type FriendLinkCategoryHandler struct {
	service *service.FriendLinkCategoryService
}

func NewFriendLinkCategoryHandler() *FriendLinkCategoryHandler {
	return &FriendLinkCategoryHandler{
		service: service.NewFriendLinkCategoryService(),
	}
}

// Create 创建友链分类
func (h *FriendLinkCategoryHandler) Create(c *gin.Context) {
	var req service.CreateFriendLinkCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	category, err := h.service.Create(&req)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "分类创建成功", category)
}

// GetByID 获取友链分类详情
func (h *FriendLinkCategoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的分类ID")
		return
	}

	category, err := h.service.GetByID(uint(id))
	if err != nil {
		util.Error(c, 404, "分类不存在")
		return
	}

	util.Success(c, category)
}

// List 获取所有友链分类列表
func (h *FriendLinkCategoryHandler) List(c *gin.Context) {
	categories, err := h.service.List()
	if err != nil {
		util.ServerError(c, "获取分类列表失败")
		return
	}

	util.Success(c, categories)
}

// Update 更新友链分类
func (h *FriendLinkCategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的分类ID")
		return
	}

	var req service.UpdateFriendLinkCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	category, err := h.service.Update(uint(id), &req)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "分类更新成功", category)
}

// Delete 删除友链分类
func (h *FriendLinkCategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的分类ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "分类删除成功", nil)
}
