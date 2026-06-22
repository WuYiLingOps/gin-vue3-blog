/*
 * 项目名称：blog-backend
 * 文件名称：post_revision.go
 * 创建时间：2026-04-28 17:48:56
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：文章修订版本处理器，提供审批相关的API接口
 */
package handler

import (
	"blog-backend/constant"
	"blog-backend/service"
	"blog-backend/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PostRevisionHandler 文章修订版本处理器结构体
type PostRevisionHandler struct {
	service *service.PostRevisionService
}

// NewPostRevisionHandler 创建文章修订版本处理器实例
func NewPostRevisionHandler(notificationService *service.NotificationService) *PostRevisionHandler {
	revisionService := service.NewPostRevisionService()
	revisionService.SetNotificationService(notificationService)
	return &PostRevisionHandler{
		service: revisionService,
	}
}

// GetPendingList 获取待审批列表
func (h *PostRevisionHandler) GetPendingList(c *gin.Context) {
	// 验证是否为超级管理员
	role, exists := c.Get("role")
	if !exists || role != constant.RoleSuperAdmin {
		util.Forbidden(c, "仅超级管理员可查看")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	revisions, total, err := h.service.GetPendingList(page, pageSize)
	if err != nil {
		util.Error(c, 500, "获取列表失败")
		return
	}

	util.Success(c, gin.H{
		"list":      revisions,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetPendingCount 获取待审批数量
func (h *PostRevisionHandler) GetPendingCount(c *gin.Context) {
	// 验证是否为超级管理员
	role, exists := c.Get("role")
	if !exists || role != constant.RoleSuperAdmin {
		util.Forbidden(c, "仅超级管理员可查看")
		return
	}

	count, err := h.service.GetPendingCount()
	if err != nil {
		util.Error(c, 500, "获取数量失败")
		return
	}

	util.Success(c, gin.H{
		"count": count,
	})
}

// GetRevisionDetail 获取修订详情
func (h *PostRevisionHandler) GetRevisionDetail(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		util.Unauthorized(c, "未登录")
		return
	}

	role, _ := c.Get("role")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "参数错误")
		return
	}

	revision, err := h.service.GetRevisionDetail(uint(id))
	if err != nil {
		util.Error(c, 404, "修订不存在")
		return
	}

	// 权限检查：超级管理员可以查看所有，普通管理员只能查看自己的
	if role != constant.RoleSuperAdmin && revision.EditorID != userID.(uint) {
		util.Forbidden(c, "无权查看此修订")
		return
	}

	util.Success(c, revision)
}

// GetRevisionDiff 获取修改对比
func (h *PostRevisionHandler) GetRevisionDiff(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		util.Unauthorized(c, "未登录")
		return
	}

	role, _ := c.Get("role")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "参数错误")
		return
	}

	// 先获取修订信息以检查权限
	revision, err := h.service.GetRevisionDetail(uint(id))
	if err != nil {
		util.Error(c, 404, "修订不存在")
		return
	}

	// 权限检查：超级管理员可以查看所有，普通管理员只能查看自己的
	if role != constant.RoleSuperAdmin && revision.EditorID != userID.(uint) {
		util.Forbidden(c, "无权查看此修订")
		return
	}

	diff, err := h.service.GetRevisionDiff(uint(id))
	if err != nil {
		util.Error(c, 404, "修订不存在")
		return
	}

	util.Success(c, diff)
}

// Approve 审批通过
func (h *PostRevisionHandler) Approve(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		util.Unauthorized(c, "未登录")
		return
	}

	// 验证是否为超级管理员
	role, exists := c.Get("role")
	if !exists || role != constant.RoleSuperAdmin {
		util.Forbidden(c, "仅超级管理员可审批")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的修订ID")
		return
	}

	if err := h.service.ApproveRevision(uint(id), userID.(uint)); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	// 记录操作日志
	revisionID := uint(id)
	util.LogOperation(c, "approve", "post_revision", &revisionID, "", "审批通过文章修订")

	util.Success(c, "审批通过")
}

// Reject 审批拒绝
func (h *PostRevisionHandler) Reject(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		util.Unauthorized(c, "未登录")
		return
	}

	// 验证是否为超级管理员
	role, exists := c.Get("role")
	if !exists || role != constant.RoleSuperAdmin {
		util.Forbidden(c, "仅超级管理员可审批")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的修订ID")
		return
	}

	var req struct {
		RejectReason string `json:"reject_reason" binding:"max=500"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "拒绝原因最多500字")
		return
	}

	if err := h.service.RejectRevision(uint(id), userID.(uint), req.RejectReason); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	// 记录操作日志
	revisionID := uint(id)
	util.LogOperation(c, "reject", "post_revision", &revisionID, req.RejectReason, "审批拒绝文章修订")

	util.Success(c, "审批拒绝")
}

// GetPostRevisions 获取文章修订历史
func (h *PostRevisionHandler) GetPostRevisions(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的文章ID")
		return
	}

	revisions, err := h.service.GetPostRevisions(uint(postID))
	if err != nil {
		util.Error(c, 500, "获取修订历史失败")
		return
	}

	util.Success(c, revisions)
}

// Withdraw 撤回修订
func (h *PostRevisionHandler) Withdraw(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		util.Unauthorized(c, "未登录")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的修订ID")
		return
	}

	if err := h.service.WithdrawRevision(uint(id), userID.(uint)); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	// 记录操作日志
	revisionID := uint(id)
	util.LogOperation(c, "withdraw", "post_revision", &revisionID, "", "撤回文章修订")

	util.Success(c, "撤回成功")
}

// GetMyRevisions 获取当前用户的修订记录
func (h *PostRevisionHandler) GetMyRevisions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		util.Unauthorized(c, "未登录")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status") // 可选：pending, approved, rejected, withdrawn

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	revisions, total, err := h.service.GetUserRevisions(userID.(uint), page, pageSize, status)
	if err != nil {
		util.Error(c, 500, "获取列表失败")
		return
	}

	util.Success(c, gin.H{
		"list":      revisions,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
