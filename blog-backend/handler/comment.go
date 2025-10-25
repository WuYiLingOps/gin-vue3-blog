package handler

import (
	"strconv"

	"blog-backend/service"
	"blog-backend/util"
	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	service *service.CommentService
}

func NewCommentHandler() *CommentHandler {
	return &CommentHandler{
		service: service.NewCommentService(),
	}
}

// Create 创建评论
func (h *CommentHandler) Create(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		util.Unauthorized(c, "未登录")
		return
	}

	var req service.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	comment, err := h.service.Create(userID.(uint), &req)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "评论成功", comment)
}

// GetByID 获取评论详情
func (h *CommentHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的评论ID")
		return
	}

	comment, err := h.service.GetByID(uint(id))
	if err != nil {
		util.Error(c, 404, err.Error())
		return
	}

	util.Success(c, comment)
}

// Update 更新评论
func (h *CommentHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的评论ID")
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var req service.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	comment, err := h.service.Update(uint(id), userID.(uint), role.(string), &req)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "评论更新成功", comment)
}

// Delete 删除评论
func (h *CommentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的评论ID")
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	if err := h.service.Delete(uint(id), userID.(uint), role.(string)); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "评论删除成功", nil)
}

// GetByPostID 获取文章的评论列表
func (h *CommentHandler) GetByPostID(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的文章ID")
		return
	}

	comments, err := h.service.GetByPostID(uint(postID))
	if err != nil {
		util.ServerError(c, "获取评论列表失败")
		return
	}

	util.Success(c, comments)
}

// List 获取评论列表（管理后台）
func (h *CommentHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	comments, total, err := h.service.List(page, pageSize)
	if err != nil {
		util.ServerError(c, "获取评论列表失败")
		return
	}

	util.PageSuccess(c, comments, total, page, pageSize)
}

// UpdateStatus 更新评论状态
func (h *CommentHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的评论ID")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.service.UpdateStatus(uint(id), req.Status); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "状态更新成功", nil)
}

