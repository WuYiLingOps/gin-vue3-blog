package handler

import (
	"strconv"

	"blog-backend/service"
	"blog-backend/util"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		service: service.NewUserService(),
	}
}

// GetByID 获取用户详情
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的用户ID")
		return
	}

	user, err := h.service.GetByID(uint(id))
	if err != nil {
		util.Error(c, 404, err.Error())
		return
	}

	util.Success(c, user)
}

// List 获取用户列表
func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, err := h.service.List(page, pageSize)
	if err != nil {
		util.ServerError(c, "获取用户列表失败")
		return
	}

	util.PageSuccess(c, users, total, page, pageSize)
}

// UpdateStatus 更新用户状态
func (h *UserHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的用户ID")
		return
	}

	var req struct {
		Status *int `json:"status" binding:"required"` // 使用指针类型，避免 0 值被认为是空
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	if req.Status == nil {
		util.BadRequest(c, "status 参数不能为空")
		return
	}

	// 验证 status 值（0:禁用 1:启用）
	if *req.Status != 0 && *req.Status != 1 {
		util.BadRequest(c, "status 参数值无效，必须为 0 或 1")
		return
	}

	if err := h.service.UpdateStatus(uint(id), *req.Status); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "状态更新成功", nil)
}

// Delete 删除用户
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的用户ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "用户删除成功", nil)
}
