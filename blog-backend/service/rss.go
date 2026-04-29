/*
 * 项目名称：blog-backend
 * 文件名称：rss.go
 * 创建时间：2026-04-23
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：RSS Feed 生成业务逻辑层，提供文章、说说等内容的 RSS 订阅功能
 */
package service

import (
	"blog-backend/config"
	"blog-backend/logger"
	"blog-backend/model"
	"blog-backend/repository"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/feeds"
)

// RSSService RSS 业务逻辑层结构体
type RSSService struct {
	postRepo     *repository.PostRepository
	momentRepo   *repository.MomentRepository
	categoryRepo *repository.CategoryRepository
	tagRepo      *repository.TagRepository
	settingRepo  *repository.SettingRepository
	userRepo     *repository.UserRepository
	config       *config.Config
	cache        *RSSCache
}

// RSSCache RSS 缓存结构
type RSSCache struct {
	feeds map[string]*CachedFeed
	mu    sync.RWMutex
}

// CachedFeed 缓存的 Feed 数据
type CachedFeed struct {
	Content   string
	ExpiresAt time.Time
}

// NewRSSService 创建 RSS 业务逻辑层实例
func NewRSSService(cfg *config.Config) *RSSService {
	return &RSSService{
		postRepo:     repository.NewPostRepository(),
		momentRepo:   repository.NewMomentRepository(),
		categoryRepo: repository.NewCategoryRepository(),
		tagRepo:      repository.NewTagRepository(),
		settingRepo:  repository.NewSettingRepository(),
		userRepo:     repository.NewUserRepository(),
		config:       cfg,
		cache: &RSSCache{
			feeds: make(map[string]*CachedFeed),
		},
	}
}

// RSSConfig RSS 配置结构
type RSSConfig struct {
	Enabled       bool   `json:"enabled"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Link          string `json:"link"`
	AuthorName    string `json:"author_name"`
	AuthorEmail   string `json:"author_email"`
	Language      string `json:"language"`
	Copyright     string `json:"copyright"`
	ItemLimit     int    `json:"item_limit"`
	CacheDuration int    `json:"cache_duration"` // 秒（内部使用）
}

// GetRSSConfig 获取 RSS 配置
func (s *RSSService) GetRSSConfig(ctx context.Context) (*RSSConfig, error) {
	// 获取超级管理员信息作为默认值
	superAdmin, err := s.userRepo.GetSuperAdmin()
	if err != nil {
		logger.Warn("获取超级管理员信息失败: " + err.Error())
	}

	// 从超级管理员信息中提取默认值
	siteName := "菱风叙"
	siteDescription := "个人技术博客"
	authorName := "無以菱"
	authorEmail := "huangjing510@126.com"
	blogURL := s.config.App.BlogURL

	// 使用超级管理员的信息
	if superAdmin != nil && superAdmin.ID > 0 {
		logger.Info(fmt.Sprintf("超级管理员信息: ID=%d, Nickname=%s, Bio=%s, Email=%s",
			superAdmin.ID, superAdmin.Nickname, superAdmin.Bio, superAdmin.Email))

		if superAdmin.Nickname != "" {
			authorName = superAdmin.Nickname
		}
		if superAdmin.Bio != "" {
			siteDescription = superAdmin.Bio // RSS 描述使用超级管理员个人简介
		}
		if superAdmin.Email != "" {
			authorEmail = superAdmin.Email
		}
	} else {
		logger.Warn("未找到超级管理员或超级管理员信息为空")
	}

	// 获取网站基本设置（如果配置了则覆盖默认值）
	siteSettings, err := s.settingRepo.GetByGroup("site")
	if err != nil {
		return nil, fmt.Errorf("获取网站设置失败: %w", err)
	}

	for _, setting := range siteSettings {
		if setting.Key == "site_name" && setting.Value != "" {
			logger.Info(fmt.Sprintf("从网站设置读取 site_name: %s", setting.Value))
			siteName = setting.Value // RSS 标题使用网站设置中的网站名称
			break
		}
	}

	logger.Info(fmt.Sprintf("最终默认值: siteName=%s, siteDescription=%s, authorName=%s, authorEmail=%s",
		siteName, siteDescription, authorName, authorEmail))

	// 初始化 RSS 配置，使用超级管理员信息和网站设置作为默认值
	config := &RSSConfig{
		Enabled:       true,
		Title:         siteName,        // 使用网站名称
		Description:   siteDescription, // 使用超级管理员个人简介
		Link:          blogURL,
		AuthorName:    authorName,    // 使用超级管理员昵称
		AuthorEmail:   authorEmail,   // 使用超级管理员邮箱
		Language:      "zh-CN",
		Copyright:     fmt.Sprintf("Copyright © %s", siteName),
		ItemLimit:     20,
		CacheDuration: 3600,
	}

	// 获取 RSS 专属配置（如果已设置则覆盖默认值）
	rssSettings, err := s.settingRepo.GetByGroup("rss")
	if err != nil {
		return nil, fmt.Errorf("获取 RSS 配置失败: %w", err)
	}

	// 从 RSS 配置中读取自定义值
	for _, setting := range rssSettings {
		switch setting.Key {
		case "rss_enabled":
			config.Enabled = setting.Value == "true"
		case "rss_title":
			if setting.Value != "" {
				config.Title = setting.Value
			}
		case "rss_description":
			if setting.Value != "" {
				config.Description = setting.Value
			}
		case "rss_link":
			if setting.Value != "" {
				config.Link = setting.Value
			}
		case "rss_author_name":
			if setting.Value != "" {
				config.AuthorName = setting.Value
			}
		case "rss_author_email":
			if setting.Value != "" {
				config.AuthorEmail = setting.Value
			}
		case "rss_language":
			if setting.Value != "" {
				config.Language = setting.Value
			}
		case "rss_copyright":
			if setting.Value != "" {
				config.Copyright = setting.Value
			}
		case "rss_item_limit":
			if setting.Value != "" {
				fmt.Sscanf(setting.Value, "%d", &config.ItemLimit)
			}
		case "rss_cache_duration":
			if setting.Value != "" {
				fmt.Sscanf(setting.Value, "%d", &config.CacheDuration)
			}
		}
	}

	return config, nil
}

// UpdateRSSConfig 更新 RSS 配置
func (s *RSSService) UpdateRSSConfig(ctx context.Context, config *RSSConfig) error {
	settingsData := []struct {
		key   string
		value string
	}{
		{"rss_enabled", fmt.Sprintf("%t", config.Enabled)},
		{"rss_title", config.Title},
		{"rss_description", config.Description},
		{"rss_link", config.Link},
		{"rss_author_name", config.AuthorName},
		{"rss_author_email", config.AuthorEmail},
		{"rss_language", config.Language},
		{"rss_copyright", config.Copyright},
		{"rss_item_limit", fmt.Sprintf("%d", config.ItemLimit)},
		{"rss_cache_duration", fmt.Sprintf("%d", config.CacheDuration)},
	}

	var settings []model.Setting
	for _, item := range settingsData {
		settings = append(settings, model.Setting{
			Key:   item.key,
			Value: item.value,
			Type:  "text",
			Group: "rss",
			Label: item.key,
		})
	}

	if err := s.settingRepo.BatchUpsert(settings); err != nil {
		return fmt.Errorf("更新 RSS 配置失败: %w", err)
	}

	// 清除缓存
	s.ClearCache()

	return nil
}

// GeneratePostsFeed 生成文章 RSS Feed
func (s *RSSService) GeneratePostsFeed(ctx context.Context) (string, string, string, error) {
	cacheKey := "rss:posts"

	// 检查缓存
	if cached := s.getCache(cacheKey); cached != "" {
		etag, lastModified := s.generateCacheHeaders(cached)
		return cached, etag, lastModified, nil
	}

	config, err := s.GetRSSConfig(ctx)
	if err != nil {
		return "", "", "", err
	}

	if !config.Enabled {
		return "", "", "", fmt.Errorf("RSS 功能未启用")
	}

	// 获取最新文章（只获取已发布的公开文章）
	visibility := 1
	posts, _, err := s.postRepo.List(1, config.ItemLimit, 0, "", 1, &visibility)
	if err != nil {
		return "", "", "", fmt.Errorf("获取文章列表失败: %w", err)
	}

	feed := s.createBaseFeed(config, "最新文章")
	feed.Link = &feeds.Link{Href: s.config.App.BlogURL + "/api/rss/posts.xml"}

	for _, post := range posts {
		item := s.postToFeedItem(&post)
		feed.Items = append(feed.Items, item)
	}

	rss, err := feed.ToRss()
	if err != nil {
		return "", "", "", fmt.Errorf("生成 RSS Feed 失败: %w", err)
	}

	// 缓存结果
	s.setCache(cacheKey, rss, config.CacheDuration)

	// 生成缓存头
	etag, lastModified := s.generateCacheHeaders(rss)

	return rss, etag, lastModified, nil
}

// GenerateMomentsFeed 生成说说 RSS Feed
func (s *RSSService) GenerateMomentsFeed(ctx context.Context) (string, string, string, error) {
	cacheKey := "rss:moments"

	// 检查缓存
	if cached := s.getCache(cacheKey); cached != "" {
		etag, lastModified := s.generateCacheHeaders(cached)
		return cached, etag, lastModified, nil
	}

	config, err := s.GetRSSConfig(ctx)
	if err != nil {
		return "", "", "", err
	}

	if !config.Enabled {
		return "", "", "", fmt.Errorf("RSS 功能未启用")
	}

	// 获取最新说说（只获取公开的）
	status := 1
	moments, _, err := s.momentRepo.List(1, config.ItemLimit, &status, "")
	if err != nil {
		return "", "", "", fmt.Errorf("获取说说列表失败: %w", err)
	}

	feed := s.createBaseFeed(config, "最新说说")
	feed.Link = &feeds.Link{Href: s.config.App.BlogURL + "/api/rss/moments.xml"}

	for _, moment := range moments {
		item := s.momentToFeedItem(&moment)
		feed.Items = append(feed.Items, item)
	}

	rss, err := feed.ToRss()
	if err != nil {
		return "", "", "", fmt.Errorf("生成 RSS Feed 失败: %w", err)
	}

	// 缓存结果
	s.setCache(cacheKey, rss, config.CacheDuration)

	etag, lastModified := s.generateCacheHeaders(rss)
	return rss, etag, lastModified, nil
}

// GenerateAllFeed 生成全站 RSS Feed（文章 + 说说混合）
func (s *RSSService) GenerateAllFeed(ctx context.Context) (string, string, string, error) {
	cacheKey := "rss:all"

	// 检查缓存
	if cached := s.getCache(cacheKey); cached != "" {
		etag, lastModified := s.generateCacheHeaders(cached)
		return cached, etag, lastModified, nil
	}

	config, err := s.GetRSSConfig(ctx)
	if err != nil {
		return "", "", "", err
	}

	if !config.Enabled {
		return "", "", "", fmt.Errorf("RSS 功能未启用")
	}

	feed := s.createBaseFeed(config, "全站更新")
	feed.Link = &feeds.Link{Href: s.config.App.BlogURL + "/feed.xml"}

	// 获取文章和说说
	visibility := 1
	posts, _, _ := s.postRepo.List(1, config.ItemLimit/2, 0, "", 1, &visibility)
	status := 1
	moments, _, _ := s.momentRepo.List(1, config.ItemLimit/2, &status, "")

	// 转换为 Feed Items
	var items []*feeds.Item
	for _, post := range posts {
		items = append(items, s.postToFeedItem(&post))
	}
	for _, moment := range moments {
		items = append(items, s.momentToFeedItem(&moment))
	}

	// 按时间排序（最新的在前）
	for i := 0; i < len(items)-1; i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i].Created.Before(items[j].Created) {
				items[i], items[j] = items[j], items[i]
			}
		}
	}

	// 限制总数
	if len(items) > config.ItemLimit {
		items = items[:config.ItemLimit]
	}

	feed.Items = items

	rss, err := feed.ToRss()
	if err != nil {
		return "", "", "", fmt.Errorf("生成 RSS Feed 失败: %w", err)
	}

	// 缓存结果
	s.setCache(cacheKey, rss, config.CacheDuration)

	etag, lastModified := s.generateCacheHeaders(rss)
	return rss, etag, lastModified, nil
}

// GenerateCategoryFeed 生成分类 RSS Feed
func (s *RSSService) GenerateCategoryFeed(ctx context.Context, categoryID uint) (string, string, string, error) {
	cacheKey := fmt.Sprintf("rss:category:%d", categoryID)

	// 检查缓存
	if cached := s.getCache(cacheKey); cached != "" {
		etag, lastModified := s.generateCacheHeaders(cached)
		return cached, etag, lastModified, nil
	}

	config, err := s.GetRSSConfig(ctx)
	if err != nil {
		return "", "", "", err
	}

	if !config.Enabled {
		return "", "", "", fmt.Errorf("RSS 功能未启用")
	}

	// 获取分类信息
	category, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return "", "", "", fmt.Errorf("获取分类信息失败: %w", err)
	}

	// 获取该分类的文章
	visibility := 1
	posts, _, err := s.postRepo.List(1, config.ItemLimit, categoryID, "", 1, &visibility)
	if err != nil {
		return "", "", "", fmt.Errorf("获取文章列表失败: %w", err)
	}

	feed := s.createBaseFeed(config, fmt.Sprintf("%s - %s", category.Name, config.Title))
	feed.Description = category.Description
	feed.Link = &feeds.Link{Href: fmt.Sprintf("%s/api/rss/category/%d.xml", s.config.App.BlogURL, categoryID)}

	for _, post := range posts {
		item := s.postToFeedItem(&post)
		feed.Items = append(feed.Items, item)
	}

	rss, err := feed.ToRss()
	if err != nil {
		return "", "", "", fmt.Errorf("生成 RSS Feed 失败: %w", err)
	}

	// 缓存结果
	s.setCache(cacheKey, rss, config.CacheDuration)

	etag, lastModified := s.generateCacheHeaders(rss)
	return rss, etag, lastModified, nil
}

// GenerateTagFeed 生成标签 RSS Feed
func (s *RSSService) GenerateTagFeed(ctx context.Context, tagID uint) (string, string, string, error) {
	cacheKey := fmt.Sprintf("rss:tag:%d", tagID)

	// 检查缓存
	if cached := s.getCache(cacheKey); cached != "" {
		etag, lastModified := s.generateCacheHeaders(cached)
		return cached, etag, lastModified, nil
	}

	config, err := s.GetRSSConfig(ctx)
	if err != nil {
		return "", "", "", err
	}

	if !config.Enabled {
		return "", "", "", fmt.Errorf("RSS 功能未启用")
	}

	// 获取标签信息
	tag, err := s.tagRepo.GetByID(tagID)
	if err != nil {
		return "", "", "", fmt.Errorf("获取标签信息失败: %w", err)
	}

	// 获取该标签的已发布文章
	posts, err := s.postRepo.GetPublishedPostsByTag(tagID, config.ItemLimit)
	if err != nil {
		return "", "", "", fmt.Errorf("获取文章列表失败: %w", err)
	}

	feed := s.createBaseFeed(config, fmt.Sprintf("#%s - %s", tag.Name, config.Title))
	feed.Link = &feeds.Link{Href: fmt.Sprintf("%s/api/rss/tag/%d.xml", s.config.App.BlogURL, tagID)}

	for _, post := range posts {
		item := s.postToFeedItem(&post)
		feed.Items = append(feed.Items, item)
	}

	rss, err := feed.ToRss()
	if err != nil {
		return "", "", "", fmt.Errorf("生成 RSS Feed 失败: %w", err)
	}

	// 缓存结果
	s.setCache(cacheKey, rss, config.CacheDuration)

	etag, lastModified := s.generateCacheHeaders(rss)
	return rss, etag, lastModified, nil
}

// ClearCache 清除所有 RSS 缓存
func (s *RSSService) ClearCache() {
	s.cache.mu.Lock()
	defer s.cache.mu.Unlock()
	s.cache.feeds = make(map[string]*CachedFeed)
	logger.Info("RSS 缓存已清除")
}

// createBaseFeed 创建基础 Feed 结构
func (s *RSSService) createBaseFeed(config *RSSConfig, title string) *feeds.Feed {
	now := time.Now()
	return &feeds.Feed{
		Title:       title,
		Link:        &feeds.Link{Href: s.config.App.BlogURL},
		Description: config.Description,
		Author: &feeds.Author{
			Name:  config.AuthorName,
			Email: config.AuthorEmail,
		},
		Created: now,
		Updated: now,
	}
}

// postToFeedItem 将文章转换为 Feed Item
func (s *RSSService) postToFeedItem(post *model.Post) *feeds.Item {
	link := fmt.Sprintf("%s/post/%s", s.config.App.BlogURL, post.Slug)
	
	description := post.Summary
	if description == "" && len(post.Content) > 200 {
		description = post.Content[:200] + "..."
	}

	item := &feeds.Item{
		Title:       post.Title,
		Link:        &feeds.Link{Href: link},
		Description: description,
		Id:          link,
		Created:     post.CreatedAt,
	}

	if post.PublishedAt != nil {
		item.Created = *post.PublishedAt
	}

	return item
}

// momentToFeedItem 将说说转换为 Feed Item
func (s *RSSService) momentToFeedItem(moment *model.Moment) *feeds.Item {
	link := fmt.Sprintf("%s/moments/%d", s.config.App.BlogURL, moment.ID)
	
	title := "说说"
	if len(moment.Content) > 30 {
		title = moment.Content[:30] + "..."
	} else {
		title = moment.Content
	}

	return &feeds.Item{
		Title:       title,
		Link:        &feeds.Link{Href: link},
		Description: moment.Content,
		Id:          link,
		Created:     moment.CreatedAt,
	}
}

// getCache 获取缓存
func (s *RSSCache) getCache(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if cached, ok := s.feeds[key]; ok {
		if time.Now().Before(cached.ExpiresAt) {
			return cached.Content
		}
		// 缓存过期，删除
		delete(s.feeds, key)
	}

	return ""
}

// setCache 设置缓存
func (s *RSSCache) setCache(key, content string, duration int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.feeds[key] = &CachedFeed{
		Content:   content,
		ExpiresAt: time.Now().Add(time.Duration(duration) * time.Second),
	}
}

// getCache 获取缓存（RSSService 方法）
func (s *RSSService) getCache(key string) string {
	return s.cache.getCache(key)
}

// setCache 设置缓存（RSSService 方法）
func (s *RSSService) setCache(key, content string, duration int) {
	s.cache.setCache(key, content, duration)
}

// generateCacheHeaders 生成缓存头（ETag 和 Last-Modified）
func (s *RSSService) generateCacheHeaders(content string) (string, string) {
	// 生成 ETag：使用内容的 MD5 哈希
	hash := md5.Sum([]byte(content))
	etag := `"` + hex.EncodeToString(hash[:]) + `"`

	// 生成 Last-Modified：使用当前时间（RFC1123 格式）
	lastModified := time.Now().UTC().Format(time.RFC1123)

	return etag, lastModified
}
