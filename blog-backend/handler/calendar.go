/*
 * 项目名称：blog-backend
 * 文件名称：calendar.go
 * 创建时间：2026-01-31 16:05:15
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：日历处理器，提供Gitee贡献热力图数据查询功能
 */
package handler

import (
	"blog-backend/service"
	"blog-backend/util"

	"github.com/gin-gonic/gin"
)

// CalendarHandler 日历处理器结构体
type CalendarHandler struct {
	service *service.CalendarService
}

// NewCalendarHandler 创建日历处理器实例
func NewCalendarHandler() *CalendarHandler {
	return &CalendarHandler{
		service: service.NewCalendarService(),
	}
}

// GetGiteeCalendar 获取Gitee贡献热力图数据
func (h *CalendarHandler) GetGiteeCalendar(c *gin.Context) {
	username := c.Query("user")
	if username == "" {
		util.BadRequest(c, "用户名不能为空")
		return
	}

	data, err := h.service.GetGiteeCalendar(username)
	if err != nil {
		util.Error(c, 500, err.Error())
		return
	}

	util.Success(c, data)
}
