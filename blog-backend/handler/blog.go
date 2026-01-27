package handler

import (
	"context"
	"encoding/json"
	"time"

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
	// 先尝试从 Redis 获取缓存
	ctx := context.Background()
	cacheKey := "blog:author_profile"
	if cached, err := db.RDB.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
		var resp AuthorProfileResponse
		if err := json.Unmarshal([]byte(cached), &resp); err == nil {
			util.Success(c, resp)
			return
		}
		// 如果解析失败则继续从数据库获取
	}

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

	// 将结果写入 Redis，适当设置过期时间（例如 10 分钟）
	if data, err := json.Marshal(response); err == nil {
		_ = db.RDB.Set(ctx, cacheKey, string(data), 15*time.Minute).Err()
	}

	util.Success(c, response)
}

// TagStat 标签统计数据
type TagStat struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
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
	ctx := context.Background()
	cacheKey := "tag:stats:top10"

	// 1. 先尝试从 Redis 获取缓存
	if cached, err := db.RDB.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
		var stats []TagStat
		if err := json.Unmarshal([]byte(cached), &stats); err == nil {
			util.Success(c, stats)
			return
		}
	}

	var stats []TagStat

	// 2. 缓存未命中，从数据库获取
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

	// 3. 写入缓存，设置过期时间（例如 10 分钟）
	if data, err := json.Marshal(stats); err == nil {
		_ = db.RDB.Set(ctx, cacheKey, string(data), 10*time.Minute).Err()
	}

	util.Success(c, stats)
}