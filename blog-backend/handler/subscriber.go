/*
 * 项目名称：blog-backend
 * 文件名称：subscriber.go
 * 创建时间：2026-04-17 23:40:47
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：邮件订阅者处理器，提供订阅、退订、订阅者管理功能
 */
package handler

import (
	"blog-backend/config"
	"blog-backend/service"
	"blog-backend/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SubscriberHandler 订阅者处理器结构体
type SubscriberHandler struct {
	service *service.SubscriberService
}

// NewSubscriberHandler 创建订阅者处理器实例
func NewSubscriberHandler(cfg *config.Config) *SubscriberHandler {
	return &SubscriberHandler{
		service: service.NewSubscriberService(cfg),
	}
}

// SubscribeRequest 订阅请求
type SubscribeRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Source    string `json:"source"`     // 订阅来源：homepage/post/subscribe/other
	SourceURL string `json:"source_url"` // 来源URL
}

// Subscribe 订阅
func (h *SubscriberHandler) Subscribe(c *gin.Context) {
	var req SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请输入有效的邮箱地址")
		return
	}

	// 获取客户端信息
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// 如果前端没有传递来源，默认为 subscribe
	if req.Source == "" {
		req.Source = "subscribe"
	}

	if err := h.service.Subscribe(c.Request.Context(), req.Email, req.Source, req.SourceURL, ip, userAgent); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "订阅成功！欢迎邮件已发送到您的邮箱", nil)
}

// Unsubscribe 退订
func (h *SubscriberHandler) Unsubscribe(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		util.BadRequest(c, "缺少退订令牌")
		return
	}

	if err := h.service.Unsubscribe(c.Request.Context(), token); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "退订成功", nil)
}

// List 获取订阅者列表（管理员用）
func (h *SubscriberHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	subscribers, total, err := h.service.List(c.Request.Context(), page, pageSize)
	if err != nil {
		util.ServerError(c, "获取订阅者列表失败")
		return
	}

	util.PageSuccess(c, subscribers, total, page, pageSize)
}

// Delete 删除订阅者（管理员用）
func (h *SubscriberHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的订阅者ID")
		return
	}

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "删除成功", nil)
}

// GetStats 获取订阅统计信息（公开接口）
// @Summary 获取订阅统计信息
// @Description 获取活跃订阅者数量等统计信息
// @Tags Subscribe
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/subscribe/stats [get]
func (h *SubscriberHandler) GetStats(c *gin.Context) {
	totalCount, err := h.service.GetTotalCount(c.Request.Context())
	if err != nil {
		util.ServerError(c, "获取统计信息失败")
		return
	}

	activeCount, err := h.service.GetActiveCount(c.Request.Context())
	if err != nil {
		util.ServerError(c, "获取统计信息失败")
		return
	}

	util.Success(c, gin.H{
		"total_count":  totalCount,  // 累积订阅总数（包括已退订）
		"active_count": activeCount, // 当前活跃订阅者数量
	})
}

// GetPushHistories 获取推送历史记录列表（管理员）
func (h *SubscriberHandler) GetPushHistories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	histories, total, err := h.service.GetPushHistories(c.Request.Context(), page, pageSize)
	if err != nil {
		util.ServerError(c, "获取推送历史失败")
		return
	}

	util.PageSuccess(c, histories, total, page, pageSize)
}

// GetPushHistoryDetail 获取推送历史详情（管理员）
func (h *SubscriberHandler) GetPushHistoryDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的推送历史ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	history, details, total, err := h.service.GetPushHistoryDetail(c.Request.Context(), uint(id), page, pageSize)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.Success(c, gin.H{
		"history": history,
		"details": details,
		"total":   total,
		"page":    page,
		"page_size": pageSize,
	})
}

// DeletePushHistory 删除推送历史（管理员）
func (h *SubscriberHandler) DeletePushHistory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的推送历史ID")
		return
	}

	if err := h.service.DeletePushHistory(c.Request.Context(), uint(id)); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "删除成功", nil)
}

// GetPushStats 获取推送统计信息（管理员）
func (h *SubscriberHandler) GetPushStats(c *gin.Context) {
	stats, err := h.service.GetPushStats(c.Request.Context())
	if err != nil {
		util.ServerError(c, "获取推送统计失败")
		return
	}

	util.Success(c, stats)
}

