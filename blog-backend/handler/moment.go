/*
 * 项目名称：blog-backend
 * 文件名称：moment.go
 * 创建时间：2026-01-31 16:05:15
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：说说管理处理器，提供说说的增删改查、点赞等功能，支持公开和私密说说
 */
package handler

import (
	"blog-backend/model"
	"blog-backend/service"
	"blog-backend/util"
	"strconv"
	"strings"

	"blog-backend/constant"

	"github.com/gin-gonic/gin"
)

// MomentHandler 说说处理器结构体
type MomentHandler struct {
	service *service.MomentService
}

// NewMomentHandler 创建说说处理器实例
func NewMomentHandler() *MomentHandler {
	return &MomentHandler{
		service: service.NewMomentService(),
	}
}

// Create 创建说说
func (h *MomentHandler) Create(c *gin.Context) {
	// 使用通用 map 接收原始 JSON，避免整型零值与默认值干扰
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		util.BadRequest(c, "参数错误")
		return
	}

	// 内容校验
	contentVal, ok := body["content"].(string)
	if !ok || strings.TrimSpace(contentVal) == "" {
		util.BadRequest(c, "说说内容不能为空")
		return
	}

	// 图片（可选）
	imagesVal, _ := body["images"].(string)

	// 状态处理：仅允许 0 或 1，未传则默认公开(1)
	status := 1
	if raw, exists := body["status"]; exists {
		switch v := raw.(type) {
		case float64:
			// JSON 数字（如 0 或 1）
			if v != 0 && v != 1 {
				util.BadRequest(c, "状态参数错误")
				return
			}
			status = int(v)
		case string:
			// JSON 字符串（"0" 或 "1"）
			if v != "0" && v != "1" {
				util.BadRequest(c, "状态参数错误")
				return
			}
			if v == "0" {
				status = 0
			} else {
				status = 1
			}
		default:
			// 其他类型一律视为参数错误
			util.BadRequest(c, "状态参数错误")
			return
		}
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		util.Unauthorized(c, "未登录")
		return
	}

	moment := &model.Moment{
		Content: contentVal,
		Images:  imagesVal,
		UserID:  userID.(uint),
		Status:  status,
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
		Status  *int   `json:"status"` // 1:公开 0:私密
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "参数错误")
		return
	}

	if err := h.service.Update(uint(id), req.Content, req.Images, req.Status); err != nil {
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

	// 权限校验：私密说说仅作者或具备管理员权限的用户可见
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	var userID *uint
	if uid, exists := c.Get("user_id"); exists {
		uidVal := uid.(uint)
		userID = &uidVal
	}
	if moment.Status == 0 { // 私密
		if !constant.IsAdminRole(role) && (userID == nil || *userID != moment.UserID) {
			util.Forbidden(c, "无权查看该说说")
			return
		}
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
	}

	// 权限：具备管理员权限的用户可查看全部（含私密），普通用户/游客仅公开
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	if !constant.IsAdminRole(role) {
		publicStatus := 1
		status = &publicStatus
	}

	// 获取用户ID和IP
	var userID *uint
	if uid, exists := c.Get("user_id"); exists {
		id := uid.(uint)
		userID = &id
	}
	ip := util.GetClientIP(c)

	moments, total, err := h.service.List(page, pageSize, status, keyword, userID, ip)
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

	// 获取用户ID和IP
	var userID *uint
	if uid, exists := c.Get("user_id"); exists {
		id := uid.(uint)
		userID = &id
	}
	ip := util.GetClientIP(c)

	moments, total, err := h.service.List(page, pageSize, status, keyword, userID, ip)
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

// Like 点赞/取消点赞说说
func (h *MomentHandler) Like(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.BadRequest(c, "无效的ID")
		return
	}

	// 获取用户ID和IP
	var userID *uint
	if uid, exists := c.Get("user_id"); exists {
		id := uid.(uint)
		userID = &id
	}
	ip := util.GetClientIP(c)

	liked, err := h.service.Like(uint(id), userID, ip)
	if err != nil {
		util.Error(c, 500, err.Error())
		return
	}

	// 返回点赞状态
	util.Success(c, gin.H{
		"liked": liked,
	})
}
