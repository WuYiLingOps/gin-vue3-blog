/*
 * 项目名称：blog-backend
 * 文件名称：rss.go
 * 创建时间：2026-04-23
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：RSS Feed HTTP 处理器，提供 RSS 订阅接口和管理接口
 */
package handler

import (
	"blog-backend/config"
	"blog-backend/logger"
	"blog-backend/service"
	"blog-backend/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RSSHandler RSS 处理器结构体
type RSSHandler struct {
	rssService *service.RSSService
}

// NewRSSHandler 创建 RSS 处理器实例
func NewRSSHandler(cfg *config.Config) *RSSHandler {
	return &RSSHandler{
		rssService: service.NewRSSService(cfg),
	}
}

// GetPostsFeed 获取文章 RSS Feed
// @Summary 获取文章 RSS Feed
// @Description 获取最新文章的 RSS Feed
// @Tags RSS
// @Produce xml
// @Success 200 {string} string "RSS XML"
// @Router /api/rss/posts.xml [get]
func (h *RSSHandler) GetPostsFeed(c *gin.Context) {
	rss, etag, lastModified, err := h.rssService.GeneratePostsFeed(c.Request.Context())
	if err != nil {
		logger.Error("生成文章 RSS Feed 失败: " + err.Error())
		c.String(http.StatusInternalServerError, "生成 RSS Feed 失败")
		return
	}

	// 检查客户端缓存
	if checkClientCache(c, etag, lastModified) {
		c.Status(http.StatusNotModified)
		return
	}

	// 设置缓存头
	setCacheHeaders(c, etag, lastModified)
	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.String(http.StatusOK, rss)
}

// GetMomentsFeed 获取说说 RSS Feed
// @Summary 获取说说 RSS Feed
// @Description 获取最新说说的 RSS Feed
// @Tags RSS
// @Produce xml
// @Success 200 {string} string "RSS XML"
// @Router /api/rss/moments.xml [get]
func (h *RSSHandler) GetMomentsFeed(c *gin.Context) {
	rss, etag, lastModified, err := h.rssService.GenerateMomentsFeed(c.Request.Context())
	if err != nil {
		logger.Error("生成说说 RSS Feed 失败: " + err.Error())
		c.String(http.StatusInternalServerError, "生成 RSS Feed 失败")
		return
	}

	// 检查客户端缓存
	if checkClientCache(c, etag, lastModified) {
		c.Status(http.StatusNotModified)
		return
	}

	// 设置缓存头
	setCacheHeaders(c, etag, lastModified)
	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.String(http.StatusOK, rss)
}

// GetAllFeed 获取全站 RSS Feed
// @Summary 获取全站 RSS Feed
// @Description 获取全站更新（文章 + 说说）的 RSS Feed
// @Tags RSS
// @Produce xml
// @Success 200 {string} string "RSS XML"
// @Router /api/rss/feed.xml [get]
// @Router /feed.xml [get]
func (h *RSSHandler) GetAllFeed(c *gin.Context) {
	rss, etag, lastModified, err := h.rssService.GenerateAllFeed(c.Request.Context())
	if err != nil {
		logger.Error("生成全站 RSS Feed 失败: " + err.Error())
		c.String(http.StatusInternalServerError, "生成 RSS Feed 失败")
		return
	}

	// 检查客户端缓存
	if checkClientCache(c, etag, lastModified) {
		c.Status(http.StatusNotModified)
		return
	}

	// 设置缓存头
	setCacheHeaders(c, etag, lastModified)
	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.String(http.StatusOK, rss)
}

// GetCategoryFeed 获取分类 RSS Feed
// @Summary 获取分类 RSS Feed
// @Description 获取指定分类的文章 RSS Feed
// @Tags RSS
// @Produce xml
// @Param id path int true "分类ID"
// @Success 200 {string} string "RSS XML"
// @Router /api/rss/category/{id}.xml [get]
func (h *RSSHandler) GetCategoryFeed(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.String(http.StatusBadRequest, "无效的分类ID")
		return
	}

	rss, etag, lastModified, err := h.rssService.GenerateCategoryFeed(c.Request.Context(), uint(id))
	if err != nil {
		logger.Error("生成分类 RSS Feed 失败: " + err.Error())
		c.String(http.StatusInternalServerError, "生成 RSS Feed 失败")
		return
	}

	// 检查客户端缓存
	if checkClientCache(c, etag, lastModified) {
		c.Status(http.StatusNotModified)
		return
	}

	// 设置缓存头
	setCacheHeaders(c, etag, lastModified)
	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.String(http.StatusOK, rss)
}

// GetTagFeed 获取标签 RSS Feed
// @Summary 获取标签 RSS Feed
// @Description 获取指定标签的文章 RSS Feed
// @Tags RSS
// @Produce xml
// @Param id path int true "标签ID"
// @Success 200 {string} string "RSS XML"
// @Router /api/rss/tag/{id}.xml [get]
func (h *RSSHandler) GetTagFeed(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.String(http.StatusBadRequest, "无效的标签ID")
		return
	}

	rss, etag, lastModified, err := h.rssService.GenerateTagFeed(c.Request.Context(), uint(id))
	if err != nil {
		logger.Error("生成标签 RSS Feed 失败: " + err.Error())
		c.String(http.StatusInternalServerError, "生成 RSS Feed 失败")
		return
	}

	// 检查客户端缓存
	if checkClientCache(c, etag, lastModified) {
		c.Status(http.StatusNotModified)
		return
	}

	// 设置缓存头
	setCacheHeaders(c, etag, lastModified)
	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.String(http.StatusOK, rss)
}

// checkClientCache 检查客户端缓存是否有效
func checkClientCache(c *gin.Context, etag, lastModified string) bool {
	// 检查 If-None-Match (ETag)
	if clientETag := c.GetHeader("If-None-Match"); clientETag != "" && clientETag == etag {
		return true
	}

	// 检查 If-Modified-Since (Last-Modified)
	if clientLastModified := c.GetHeader("If-Modified-Since"); clientLastModified != "" && clientLastModified == lastModified {
		return true
	}

	return false
}

// setCacheHeaders 设置缓存响应头
func setCacheHeaders(c *gin.Context, etag, lastModified string) {
	if etag != "" {
		c.Header("ETag", etag)
	}
	if lastModified != "" {
		c.Header("Last-Modified", lastModified)
	}
	// 设置缓存控制：公开缓存，最大缓存时间 1 小时
	c.Header("Cache-Control", "public, max-age=3600")
}

// GetConfig 获取 RSS 配置（管理后台）
// @Summary 获取 RSS 配置
// @Description 获取 RSS 配置信息
// @Tags RSS Admin
// @Produce json
// @Success 200 {object} service.RSSConfig
// @Router /api/admin/rss/config [get]
func (h *RSSHandler) GetConfig(c *gin.Context) {
	config, err := h.rssService.GetRSSConfig(c.Request.Context())
	if err != nil {
		logger.Error("获取 RSS 配置失败: " + err.Error())
		util.ServerError(c, "获取 RSS 配置失败")
		return
	}

	// 转换缓存时长：秒 -> 分钟（前端使用分钟）
	cacheDurationMinutes := config.CacheDuration / 60
	if cacheDurationMinutes < 1 {
		cacheDurationMinutes = 1
	}

	util.Success(c, gin.H{
		"enabled":        config.Enabled,
		"title":          config.Title,
		"description":    config.Description,
		"link":           config.Link,
		"author_name":    config.AuthorName,
		"author_email":   config.AuthorEmail,
		"language":       config.Language,
		"copyright":      config.Copyright,
		"item_limit":     config.ItemLimit,
		"cache_duration": cacheDurationMinutes,
	})
}

// UpdateConfig 更新 RSS 配置（管理后台）
// @Summary 更新 RSS 配置
// @Description 更新 RSS 配置信息
// @Tags RSS Admin
// @Accept json
// @Produce json
// @Param config body service.RSSConfig true "RSS 配置"
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/rss/config [put]
func (h *RSSHandler) UpdateConfig(c *gin.Context) {
	var config service.RSSConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		util.BadRequest(c, "参数错误")
		return
	}

	// 转换缓存时长：分钟 -> 秒（前端使用分钟，后端存储秒）
	config.CacheDuration = config.CacheDuration * 60

	if err := h.rssService.UpdateRSSConfig(c.Request.Context(), &config); err != nil {
		logger.Error("更新 RSS 配置失败: " + err.Error())
		util.ServerError(c, "更新 RSS 配置失败")
		return
	}

	util.SuccessWithMessage(c, "更新成功", nil)
}

// Preview 预览 RSS Feed（管理后台）
// @Summary 预览 RSS Feed
// @Description 预览指定类型的 RSS Feed
// @Tags RSS Admin
// @Produce xml
// @Param type query string true "Feed 类型" Enums(posts, moments, all)
// @Success 200 {string} string "RSS XML"
// @Router /api/admin/rss/preview [get]
func (h *RSSHandler) Preview(c *gin.Context) {
	feedType := c.Query("type")
	if feedType == "" {
		feedType = "posts"
	}

	var rss, etag, lastModified string
	var err error

	switch feedType {
	case "posts":
		rss, etag, lastModified, err = h.rssService.GeneratePostsFeed(c.Request.Context())
	case "moments":
		rss, etag, lastModified, err = h.rssService.GenerateMomentsFeed(c.Request.Context())
	case "all":
		rss, etag, lastModified, err = h.rssService.GenerateAllFeed(c.Request.Context())
	default:
		c.String(http.StatusBadRequest, "无效的 Feed 类型")
		return
	}

	// 预览时不需要使用 etag 和 lastModified，用 _ 忽略
	_, _ = etag, lastModified

	if err != nil {
		logger.Error("预览 RSS Feed 失败: " + err.Error())
		c.String(http.StatusInternalServerError, "生成 RSS Feed 失败")
		return
	}

	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.String(http.StatusOK, rss)
}

// ClearCache 清除 RSS 缓存（管理后台）
// @Summary 清除 RSS 缓存
// @Description 清除所有 RSS Feed 缓存
// @Tags RSS Admin
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/rss/clear-cache [post]
func (h *RSSHandler) ClearCache(c *gin.Context) {
	h.rssService.ClearCache()
	util.SuccessWithMessage(c, "缓存已清除", nil)
}

// GetStats 获取 RSS 统计信息（管理后台）
// @Summary 获取 RSS 统计信息
// @Description 获取 RSS 订阅统计信息
// @Tags RSS Admin
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/rss/stats [get]
func (h *RSSHandler) GetStats(c *gin.Context) {
	config, err := h.rssService.GetRSSConfig(c.Request.Context())
	if err != nil {
		logger.Error("获取 RSS 配置失败: " + err.Error())
		util.ServerError(c, "获取统计信息失败")
		return
	}

	stats := gin.H{
		"enabled":        config.Enabled,
		"item_limit":     config.ItemLimit,
		"cache_duration": config.CacheDuration,
		"feeds": []gin.H{
			{
				"name": "文章 RSS",
				"url":  "/api/rss/posts.xml",
				"type": "posts",
			},
			{
				"name": "说说 RSS",
				"url":  "/api/rss/moments.xml",
				"type": "moments",
			},
			{
				"name": "全站 RSS",
				"url":  "/feed.xml",
				"type": "all",
			},
		},
	}

	util.Success(c, stats)
}

// GetStatus 获取 RSS 启用状态（公开接口）
// @Summary 获取 RSS 启用状态
// @Description 获取 RSS 功能是否启用
// @Tags RSS
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/rss/status [get]
func (h *RSSHandler) GetStatus(c *gin.Context) {
	config, err := h.rssService.GetRSSConfig(c.Request.Context())
	if err != nil {
		logger.Error("获取 RSS 配置失败: " + err.Error())
		util.ServerError(c, "获取 RSS 状态失败")
		return
	}

	util.Success(c, gin.H{
		"enabled": config.Enabled,
	})
}