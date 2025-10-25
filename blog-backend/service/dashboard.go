package service

import (
	"blog-backend/repository"
)

type DashboardService struct {
	postRepo      *repository.PostRepository
	userRepo      *repository.UserRepository
	commentRepo   *repository.CommentRepository
	categoryRepo  *repository.CategoryRepository
	visitStatRepo *repository.VisitStatRepository
}

func NewDashboardService() *DashboardService {
	return &DashboardService{
		postRepo:      repository.NewPostRepository(),
		userRepo:      repository.NewUserRepository(),
		commentRepo:   repository.NewCommentRepository(),
		categoryRepo:  repository.NewCategoryRepository(),
		visitStatRepo: repository.NewVisitStatRepository(),
	}
}

// DashboardStats 仪表盘统计数据
type DashboardStats struct {
	Posts    int64 `json:"posts"`
	Users    int64 `json:"users"`
	Comments int64 `json:"comments"`
	Views    int64 `json:"views"`
}

// CategoryStats 分类统计数据
type CategoryStats struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
	Color string `json:"color"`
}

// VisitStatData 访问量统计数据
type VisitStatData struct {
	Date  string `json:"date"`  // 日期，格式：MM-DD
	Count int    `json:"count"` // 访问量
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
func (s *DashboardService) GetCategoryStats() ([]CategoryStats, error) {
	categories, err := s.categoryRepo.List()
	if err != nil {
		return nil, err
	}

	var stats []CategoryStats
	for _, category := range categories {
		if category.PostCount > 0 {
			stats = append(stats, CategoryStats{
				Name:  category.Name,
				Value: category.PostCount,
				Color: category.Color,
			})
		}
	}

	return stats, nil
}

// GetLast7DaysVisitStats 获取最近7天的访问量统计
func (s *DashboardService) GetLast7DaysVisitStats() ([]VisitStatData, error) {
	// 确保最近7天都有记录
	if err := s.visitStatRepo.EnsureLast7DaysRecords(); err != nil {
		return nil, err
	}

	// 获取最近7天的统计数据
	stats, err := s.visitStatRepo.GetLast7DaysStats()
	if err != nil {
		return nil, err
	}

	// 转换为前端需要的格式
	var result []VisitStatData
	for _, stat := range stats {
		result = append(result, VisitStatData{
			Date:  stat.Date.Format("01-02"), // 格式化为 MM-DD
			Count: stat.ViewCount,
		})
	}

	return result, nil
}
