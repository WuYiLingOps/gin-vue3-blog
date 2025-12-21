package service

import (
	"blog-backend/repository"
	"time"
)

type DashboardService struct {
	postRepo     *repository.PostRepository
	userRepo     *repository.UserRepository
	commentRepo  *repository.CommentRepository
	categoryRepo *repository.CategoryRepository
	postViewRepo *repository.PostViewRepository
}

func NewDashboardService() *DashboardService {
	return &DashboardService{
		postRepo:     repository.NewPostRepository(),
		userRepo:     repository.NewUserRepository(),
		commentRepo:  repository.NewCommentRepository(),
		categoryRepo: repository.NewCategoryRepository(),
		postViewRepo: repository.NewPostViewRepository(),
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

	// 获取总浏览量
	totalViews, err := s.postRepo.GetTotalViews()
	if err != nil {
		return nil, err
	}
	stats.Views = totalViews

	return stats, nil
}

// GetCategoryStats 获取分类统计
// 直接从文章表统计每个分类的已发布文章数量，确保统计准确性
func (s *DashboardService) GetCategoryStats() ([]CategoryStats, error) {
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

	return stats, nil
}

// GetVisitStats 获取最近 N 天访问量统计（按天聚合）
// days 最大限制为 30 天，默认 7 天
func (s *DashboardService) GetVisitStats(days int) ([]VisitStat, error) {
	if days <= 0 || days > 30 {
		days = 7
	}

	// 取 [start, end) 区间，按天统计
	now := time.Now()
	// 以本地时区的当天 00:00 为基准，end 为明天 00:00
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := today.AddDate(0, 0, 1)
	start := end.AddDate(0, 0, -days)

	rawStats, err := s.postViewRepo.GetVisitStats(start, end)
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

	return result, nil
}
