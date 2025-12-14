package handler

import (
	"blog-backend/db"
	"blog-backend/model"
	"blog-backend/util"

	"github.com/gin-gonic/gin"
)

type BlogHandler struct{}

func NewBlogHandler() *BlogHandler {
	return &BlogHandler{}
}

// AuthorInfo 博主信息
type AuthorInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
}

// BlogStats 博客统计数据
type BlogStats struct {
	Posts      int64 `json:"posts"`      // 文章数
	Tags       int64 `json:"tags"`       // 标签数
	Categories int64 `json:"categories"` // 分类数
}

// AuthorProfileResponse 博主资料响应
type AuthorProfileResponse struct {
	Author AuthorInfo `json:"author"`
	Stats  BlogStats  `json:"stats"`
}

// GetAuthorProfile 获取博主资料和统计数据（公开接口）
func (h *BlogHandler) GetAuthorProfile(c *gin.Context) {
	// 获取管理员用户（博主）
	var author model.User
	err := db.DB.Where("role = ? AND status = ?", "admin", 1).First(&author).Error
	if err != nil {
		// 如果没有管理员，获取第一个启用的用户
		err = db.DB.Where("status = ?", 1).Order("id ASC").First(&author).Error
		if err != nil {
			util.Error(c, 404, "未找到博主信息")
			return
		}
	}

	// 获取统计数据
	var stats BlogStats

	// 文章数（已发布）
	if err := db.DB.Model(&model.Post{}).Where("status = ?", 1).Count(&stats.Posts).Error; err != nil {
		stats.Posts = 0
	}

	// 标签数
	if err := db.DB.Model(&model.Tag{}).Count(&stats.Tags).Error; err != nil {
		stats.Tags = 0
	}

	// 分类数
	if err := db.DB.Model(&model.Category{}).Count(&stats.Categories).Error; err != nil {
		stats.Categories = 0
	}

	response := AuthorProfileResponse{
		Author: AuthorInfo{
			ID:       author.ID,
			Username: author.Username,
			Nickname: author.Nickname,
			Avatar:   author.Avatar,
			Bio:      author.Bio,
		},
		Stats: stats,
	}

	util.Success(c, response)
}
