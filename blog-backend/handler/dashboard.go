package handler

import (
	"blog-backend/service"
	"blog-backend/util"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	service *service.DashboardService
}

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

// GetVisitStats 获取最近7天访问量统计
func (h *DashboardHandler) GetVisitStats(c *gin.Context) {
	stats, err := h.service.GetLast7DaysVisitStats()
	if err != nil {
		util.ServerError(c, "获取访问量统计失败")
		return
	}

	util.Success(c, stats)
}

