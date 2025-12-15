package handler

import (
	"blog-backend/repository"
	"blog-backend/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AnnouncementHandler 公告/系统广播相关接口
type AnnouncementHandler struct {
	repo *repository.ChatRepository
}

func NewAnnouncementHandler() *AnnouncementHandler {
	return &AnnouncementHandler{
		repo: repository.NewChatRepository(),
	}
}

// GetAnnouncements 获取最新公告列表
func (h *AnnouncementHandler) GetAnnouncements(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "3")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 3
	}
	if limit > 20 {
		limit = 20 // 防止一次拉取过多
	}

	list, err := h.repo.GetBroadcasts(limit)
	if err != nil {
		util.ServerError(c, "获取公告失败")
		return
	}

	util.Success(c, list)
}

// GetAnnouncementDetail 获取公告详情
func (h *AnnouncementHandler) GetAnnouncementDetail(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		util.BadRequest(c, "公告ID不合法")
		return
	}

	announcement, err := h.repo.GetBroadcastByID(uint(id64))
	if err != nil {
		util.Error(c, 404, "公告不存在或已删除")
		return
	}

	util.Success(c, announcement)
}
