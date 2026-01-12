package handler

import (
	"blog-backend/db"
	"blog-backend/model"
	"blog-backend/service"
	"blog-backend/util"

	"github.com/gin-gonic/gin"
)

type BlogHandler struct {
	settingService *service.SettingService
}

func NewBlogHandler() *BlogHandler {
	return &BlogHandler{
		settingService: service.NewSettingService(),
	}
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

// CategoryStat 分类统计数据
type CategoryStat struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
	Color string `json:"color"`
}

// TagStat 标签统计数据
type TagStat struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// GetCategoryStats 获取分类统计（公开接口）
func (h *BlogHandler) GetCategoryStats(c *gin.Context) {
	var stats []CategoryStat
	
	// 获取所有分类
	var categories []model.Category
	if err := db.DB.Find(&categories).Error; err != nil {
		util.ServerError(c, "获取分类失败")
		return
	}

	// 统计每个分类的已发布文章数
	for _, category := range categories {
		var count int64
		if err := db.DB.Model(&model.Post{}).
			Where("category_id = ? AND status = ?", category.ID, 1).
			Count(&count).Error; err == nil && count > 0 {
			stats = append(stats, CategoryStat{
				Name:  category.Name,
				Value: int(count),
				Color: category.Color,
			})
		}
	}

	util.Success(c, stats)
}

// GetAboutInfo 获取关于我信息（公开接口）
func (h *BlogHandler) GetAboutInfo(c *gin.Context) {
	content, err := h.settingService.GetAboutInfo()
	if err != nil {
		util.ServerError(c, "获取关于我信息失败")
		return
	}

	util.Success(c, map[string]string{
		"content": content,
	})
}

// GetTagStats 获取标签统计（TOP10，公开接口）
func (h *BlogHandler) GetTagStats(c *gin.Context) {
	var stats []TagStat
	
	// 获取所有标签，按文章数排序
	var tags []model.Tag
	if err := db.DB.Order("post_count DESC").Limit(10).Find(&tags).Error; err != nil {
		util.ServerError(c, "获取标签失败")
		return
	}

	// 实时计算每个标签的已发布文章数
	for _, tag := range tags {
		var count int64
		if err := db.DB.Model(&model.Post{}).
			Joins("JOIN post_tags ON posts.id = post_tags.post_id").
			Where("post_tags.tag_id = ? AND posts.status = ?", tag.ID, 1).
			Count(&count).Error; err == nil && count > 0 {
			stats = append(stats, TagStat{
				Name:  tag.Name,
				Value: int(count),
			})
		}
	}

	util.Success(c, stats)
}