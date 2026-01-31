/*
 * 项目名称：blog-backend
 * 文件名称：dashboard.go
 * 创建时间：2026-01-31 16:05:15
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：仪表盘处理器，提供统计数据查询功能
 */
package handler

import (
	"blog-backend/service"
	"blog-backend/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DashboardHandler 仪表盘处理器结构体
type DashboardHandler struct {
	service *service.DashboardService
}

// NewDashboardHandler 创建仪表盘处理器实例
func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{
		service: service.NewDashboardService(),
	}
}

// GetStats 获取仪表盘统计数据
func (h *DashboardHandler) GetStats(c *gin.Context) {
	stats, err := h.service.GetStats()
	if err != nil {
		util.ServerError(c, "获取统计数据失败")
		return
	}

	util.Success(c, stats)
}

// GetCategoryStats 获取分类统计
func (h *DashboardHandler) GetCategoryStats(c *gin.Context) {
	stats, err := h.service.GetCategoryStats()
	if err != nil {
		util.ServerError(c, "获取分类统计失败")
		return
	}

	util.Success(c, stats)
}

// GetVisitStats 获取最近 N 天访问量统计（按天）
func (h *DashboardHandler) GetVisitStats(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		days = 7
	}

	stats, err := h.service.GetVisitStats(days)
	if err != nil {
		util.ServerError(c, "获取访问统计失败")
		return
	}

	util.Success(c, stats)
}
