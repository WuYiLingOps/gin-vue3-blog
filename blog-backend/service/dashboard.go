package service

import (
	"blog-backend/repository"
)

type DashboardService struct {
	postRepo     *repository.PostRepository
	userRepo     *repository.UserRepository
	commentRepo  *repository.CommentRepository
	categoryRepo *repository.CategoryRepository
}

func NewDashboardService() *DashboardService {
	return &DashboardService{
		postRepo:     repository.NewPostRepository(),
		userRepo:     repository.NewUserRepository(),
		commentRepo:  repository.NewCommentRepository(),
		categoryRepo: repository.NewCategoryRepository(),
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
