/*
 * 项目名称：blog-backend
 * 文件名称：dashboard.go
 * 创建时间：2026-01-31 16:34:35
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：仪表盘业务逻辑层，提供统计数据查询功能，包括文章、用户、评论、访问量等统计
 */
package service

import (
	"context"
	"encoding/json"
	"time"

	"blog-backend/db"
	"blog-backend/repository"
)

// DashboardService 仪表盘业务逻辑层结构体
type DashboardService struct {
	postRepo      *repository.PostRepository
	userRepo      *repository.UserRepository
	commentRepo   *repository.CommentRepository
	categoryRepo  *repository.CategoryRepository
	postViewRepo  *repository.PostViewRepository
	pageViewRepo  *repository.PageViewRepository
}

// NewDashboardService 创建仪表盘业务逻辑层实例
func NewDashboardService() *DashboardService {
	return &DashboardService{
		postRepo:     repository.NewPostRepository(),
		userRepo:     repository.NewUserRepository(),
		commentRepo:  repository.NewCommentRepository(),
		categoryRepo: repository.NewCategoryRepository(),
		postViewRepo: repository.NewPostViewRepository(),
		pageViewRepo: repository.NewPageViewRepository(),
	}
}

// DashboardStats 仪表盘统计数据
type DashboardStats struct {
	Posts    int64 `json:"posts"`
	Users    int64 `json:"users"`
	Comments int64 `json:"comments"`
	Views    int64 `json:"views"`
}

// VisitStat 最近访问统计（按天）
type VisitStat struct {
	Date  string `json:"date"`  // 日期，格式：YYYY-MM-DD
	Count int64  `json:"count"` // 当天访问量
}

// CategoryStats 分类统计数据
type CategoryStats struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
	Color string `json:"color"`
}

// GetStats 获取统计数据
func (s *DashboardService) GetStats() (*DashboardStats, error) {
	ctx := context.Background()
	cacheKey := "dashboard:stats"

	// 1. 先尝试从 Redis 获取缓存
	if cached, err := db.RDB.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
		var stats DashboardStats
		if err := json.Unmarshal([]byte(cached), &stats); err == nil {
			return &stats, nil
		}
	}

	// 2. 缓存未命中，从数据库获取
	stats := &DashboardStats{}

	// 获取文章总数（已发布）
	postsCount, err := s.postRepo.GetPublishedCount()
	if err != nil {
		return nil, err
	}
	stats.Posts = postsCount

	// 获取用户总数
	usersCount, err := s.userRepo.GetTotalCount()
	if err != nil {
		return nil, err
	}
	stats.Users = usersCount

	// 获取评论总数
	commentsCount, err := s.commentRepo.GetTotalCount()
	if err != nil {
		return nil, err
	}
	stats.Comments = commentsCount

	// 获取总浏览量（全站PV）
	totalViews, err := s.pageViewRepo.GetTotalPV()
	if err != nil {
		return nil, err
	}
	stats.Views = totalViews

	// 3. 写入缓存，设置过期时间 2 分钟
	if data, err := json.Marshal(stats); err == nil {
		_ = db.RDB.Set(ctx, cacheKey, string(data), 2*time.Minute).Err()
	}

	return stats, nil
}

// GetCategoryStats 获取分类统计
// 直接从文章表统计每个分类的已发布文章数量，确保统计准确性
func (s *DashboardService) GetCategoryStats() ([]CategoryStats, error) {
	ctx := context.Background()
	cacheKey := "dashboard:category_stats"

	// 1. 先尝试从 Redis 获取缓存
	if cached, err := db.RDB.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
		var stats []CategoryStats
		if err := json.Unmarshal([]byte(cached), &stats); err == nil {
			return stats, nil
		}
	}

	// 2. 缓存未命中，从数据库获取
	categories, err := s.categoryRepo.List()
	if err != nil {
		return nil, err
	}

	var stats []CategoryStats
	for _, category := range categories {
		// 直接从文章表统计该分类的已发布文章数量
		count, err := s.postRepo.GetPublishedCountByCategory(category.ID)
		if err != nil {
			continue // 如果统计失败，跳过该分类
		}

		if count > 0 {
			stats = append(stats, CategoryStats{
				Name:  category.Name,
				Value: int(count),
				Color: category.Color,
			})
		}
	}

	// 3. 写入缓存，设置过期时间 5 分钟
	if data, err := json.Marshal(stats); err == nil {
		_ = db.RDB.Set(ctx, cacheKey, string(data), 5*time.Minute).Err()
	}

	return stats, nil
}

// GetVisitStats 获取最近 N 天全站访问量统计（按天聚合）
// days 最大限制为 30 天，默认 7 天
// 统计全站PV（不限于文章详情页）
func (s *DashboardService) GetVisitStats(days int) ([]VisitStat, error) {
	if days <= 0 || days > 30 {
		days = 7
	}

	ctx := context.Background()
	cacheKey := "dashboard:visit_stats"

	// 1. 先尝试从 Redis 获取缓存
	if cached, err := db.RDB.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
		var result []VisitStat
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return result, nil
		}
	}

	// 2. 缓存未命中，从数据库获取
	// 取 [start, end) 区间，按天统计
	now := time.Now()
	// 以本地时区的当天 00:00 为基准，end 为明天 00:00
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := today.AddDate(0, 0, 1)
	start := end.AddDate(0, 0, -days)

	// 使用全站PV统计（page_views表）
	rawStats, err := s.pageViewRepo.GetVisitStats(start, end)
	if err != nil {
		return nil, err
	}

	// 将查询结果按日期映射，便于补全没有访问记录的日期
	countMap := make(map[string]int64)
	for _, item := range rawStats {
		key := item.Date.Format("2006-01-02")
		countMap[key] = item.Count
	}

	// 从最早的一天到今天，按顺序补全数据
	result := make([]VisitStat, 0, days)
	for i := 0; i < days; i++ {
		day := start.AddDate(0, 0, i)
		key := day.Format("2006-01-02")
		result = append(result, VisitStat{
			Date:  key,
			Count: countMap[key],
		})
	}

	// 3. 写入缓存，设置过期时间 2 分钟
	if data, err := json.Marshal(result); err == nil {
		_ = db.RDB.Set(ctx, cacheKey, string(data), 2*time.Minute).Err()
	}

	return result, nil
}
