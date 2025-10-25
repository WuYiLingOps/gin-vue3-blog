package handler

import (
	"blog-backend/model"
	"blog-backend/service"
	"blog-backend/util"
	"github.com/gin-gonic/gin"
	"strconv"
)

type MomentHandler struct {
	service *service.MomentService
}

func NewMomentHandler() *MomentHandler {
	return &MomentHandler{
		service: service.NewMomentService(),
	}
}

// Create 创建说说
func (h *MomentHandler) Create(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
		Images  string `json:"images"`
		Status  int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "参数错误")
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		util.Unauthorized(c, "未登录")
		return
	}

	moment := &model.Moment{
		Content: req.Content,
		Images:  req.Images,
		UserID:  userID.(uint),
		Status:  req.Status,
	}

	if moment.Status == 0 {
		moment.Status = 1 // 默认公开
	}

	if err := h.service.Create(moment); err != nil {
		util.Error(c, 500, err.Error())
		return
	}

	util.SuccessWithMessage(c, "说说发布成功", moment)
}

// Update 更新说说
func (h *MomentHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.BadRequest(c, "无效的ID")
		return
	}

	var req struct {
		Content string `json:"content"`
		Images  string `json:"images"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "参数错误")
		return
	}

	if err := h.service.Update(uint(id), req.Content, req.Images); err != nil {
		util.Error(c, 500, err.Error())
		return
	}

	util.SuccessWithMessage(c, "更新成功", nil)
}

// Delete 删除说说
func (h *MomentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.BadRequest(c, "无效的ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		util.Error(c, 500, err.Error())
		return
	}

	util.SuccessWithMessage(c, "删除成功", nil)
}

// GetByID 获取说说详情
func (h *MomentHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.BadRequest(c, "无效的ID")
		return
	}

	moment, err := h.service.GetByID(uint(id))
	if err != nil {
		util.NotFound(c, "说说不存在")
		return
	}

	util.Success(c, moment)
}

// List 获取说说列表
func (h *MomentHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	statusStr := c.Query("status")
	keyword := c.Query("keyword")

	var status *int
	if statusStr != "" {
		statusVal, _ := strconv.Atoi(statusStr)
		status = &statusVal
	} else {
		// 公开接口只返回公开的说说
		publicStatus := 1
		status = &publicStatus
	}

	moments, total, err := h.service.List(page, pageSize, status, keyword)
	if err != nil {
		util.Error(c, 500, err.Error())
		return
	}

	util.Success(c, gin.H{
		"list":      moments,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetRecent 获取最新说说
func (h *MomentHandler) GetRecent(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	moments, err := h.service.GetRecent(limit)
	if err != nil {
		util.Error(c, 500, err.Error())
		return
	}

	util.Success(c, moments)
}

// AdminList 管理员获取说说列表（包含所有状态）
func (h *MomentHandler) AdminList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	statusStr := c.Query("status")
	keyword := c.Query("keyword")

	var status *int
	if statusStr != "" {
		statusVal, _ := strconv.Atoi(statusStr)
		status = &statusVal
	}

	moments, total, err := h.service.List(page, pageSize, status, keyword)
	if err != nil {
		util.Error(c, 500, err.Error())
		return
	}

	util.Success(c, gin.H{
		"list":      moments,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

