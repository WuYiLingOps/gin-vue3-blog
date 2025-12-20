package handler

import (
	"blog-backend/service"
	"blog-backend/util"

	"github.com/gin-gonic/gin"
)

type CalendarHandler struct {
	service *service.CalendarService
}

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
