package service

import (
	"errors"
	"time"

	"blog-backend/model"
	"blog-backend/repository"
	"gorm.io/gorm"
)

type PostService struct {
	postRepo     *repository.PostRepository
	categoryRepo *repository.CategoryRepository
	tagRepo      *repository.TagRepository
	postViewRepo *repository.PostViewRepository
}

func NewPostService() *PostService {
	return &PostService{
		postRepo:     repository.NewPostRepository(),
		categoryRepo: repository.NewCategoryRepository(),
		tagRepo:      repository.NewTagRepository(),
		postViewRepo: repository.NewPostViewRepository(),
	}
}

// CreatePostRequest 创建文章请求
type CreatePostRequest struct {
	Title      string   `json:"title" binding:"required"`
	Content    string   `json:"content" binding:"required"`
	Summary    string   `json:"summary"`
	Cover      string   `json:"cover"`
	CategoryID uint     `json:"category_id" binding:"required"`
	TagIDs     []uint   `json:"tag_ids"`
	Status     int      `json:"status"` // 0:草稿 1:发布
	IsTop      bool     `json:"is_top"`
}

// UpdatePostRequest 更新文章请求
type UpdatePostRequest struct {
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Summary    string   `json:"summary"`
	Cover      string   `json:"cover"`
	CategoryID uint     `json:"category_id"`
	TagIDs     []uint   `json:"tag_ids"`
	Status     int      `json:"status"`
	IsTop      bool     `json:"is_top"`
}

// Create 创建文章
func (s *PostService) Create(userID uint, req *CreatePostRequest) (*model.Post, error) {
	// 检查分类是否存在
	if _, err := s.categoryRepo.GetByID(req.CategoryID); err != nil {
		return nil, errors.New("分类不存在")
	}

	post := &model.Post{
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		Cover:      req.Cover,
		CategoryID: req.CategoryID,
		Status:     req.Status,
		IsTop:      req.IsTop,
		UserID:     userID,
	}

	// 如果是发布状态，设置发布时间
	if req.Status == 1 {
		now := time.Now()
		post.PublishedAt = &now
	}

	if err := s.postRepo.Create(post); err != nil {
		return nil, errors.New("文章创建失败")
	}

	// 更新标签关联
	if len(req.TagIDs) > 0 {
		if err := s.postRepo.UpdateTags(post.ID, req.TagIDs); err != nil {
			return nil, errors.New("标签关联失败")
		}
		
		// 增加标签文章数（仅发布状态）
		if req.Status == 1 {
			for _, tagID := range req.TagIDs {
				s.tagRepo.IncrementPostCount(tagID)
			}
		}
	}

	// 增加分类文章数
	if req.Status == 1 {
		s.categoryRepo.IncrementPostCount(req.CategoryID)
	}

	return s.postRepo.GetByID(post.ID)
}

// GetByID 获取文章详情
func (s *PostService) GetByID(id uint, userID *uint, ip string) (*model.Post, error) {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, errors.New("获取文章失败")
	}

	// 检查是否已阅读，如果没有则记录并增加浏览量
	if ip != "" && ip != "unknown" {
		hasViewed, _ := s.postViewRepo.HasViewed(id, userID, ip)
		if !hasViewed {
			// 记录阅读
			if err := s.postViewRepo.RecordView(id, userID, ip); err == nil {
				// 增加浏览量
				s.postViewRepo.IncrementViewCount(id)
				post.ViewCount++
			}
		}
	}

	// 检查是否已点赞
	liked, _ := s.postRepo.CheckLiked(id, userID, ip)
	post.Liked = liked

	return post, nil
}

// Update 更新文章
func (s *PostService) Update(id, userID uint, role string, req *UpdatePostRequest) (*model.Post, error) {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("文章不存在")
	}

	// 权限检查：只有作者和管理员可以修改
	if post.UserID != userID && role != "admin" {
		return nil, errors.New("无权限修改此文章")
	}

	oldCategoryID := post.CategoryID
	oldStatus := post.Status
	
	// 获取旧的标签列表（在更新之前）
	oldTagIDs := make([]uint, 0)
	if len(post.Tags) > 0 {
		for _, tag := range post.Tags {
			oldTagIDs = append(oldTagIDs, tag.ID)
		}
	}

	// 更新字段
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if req.Summary != "" {
		post.Summary = req.Summary
	}
	if req.Cover != "" {
		post.Cover = req.Cover
	}
	if req.CategoryID > 0 {
		// 检查新分类是否存在
		if _, err := s.categoryRepo.GetByID(req.CategoryID); err != nil {
			return nil, errors.New("分类不存在")
		}
		post.CategoryID = req.CategoryID
	}
	post.Status = req.Status
	post.IsTop = req.IsTop

	// 如果从草稿变为发布，设置发布时间
	if oldStatus == 0 && req.Status == 1 {
		now := time.Now()
		post.PublishedAt = &now
	}

	if err := s.postRepo.Update(post); err != nil {
		return nil, errors.New("文章更新失败")
	}

	// 更新标签关联
	if len(req.TagIDs) > 0 {
		if err := s.postRepo.UpdateTags(post.ID, req.TagIDs); err != nil {
			return nil, errors.New("标签关联失败")
		}
		
		// 更新标签文章数
		if oldStatus == 1 || post.Status == 1 {
			// 找出需要减少计数的标签（旧标签中有但新标签中没有的）
			for _, oldTagID := range oldTagIDs {
				found := false
				for _, newTagID := range req.TagIDs {
					if oldTagID == newTagID {
						found = true
						break
					}
				}
				if !found && oldStatus == 1 {
					s.tagRepo.DecrementPostCount(oldTagID)
				}
			}
			
			// 找出需要增加计数的标签（新标签中有但旧标签中没有的）
			for _, newTagID := range req.TagIDs {
				found := false
				for _, oldTagID := range oldTagIDs {
					if newTagID == oldTagID {
						found = true
						break
					}
				}
				if !found && post.Status == 1 {
					s.tagRepo.IncrementPostCount(newTagID)
				}
			}
			
			// 如果状态从草稿变为发布，所有新标签都要增加计数
			if oldStatus == 0 && post.Status == 1 {
				for _, tagID := range req.TagIDs {
					alreadyCounted := false
					for _, oldTagID := range oldTagIDs {
						if tagID == oldTagID {
							alreadyCounted = true
							break
						}
					}
					if alreadyCounted {
						s.tagRepo.IncrementPostCount(tagID)
					}
				}
			}
			
			// 如果状态从发布变为草稿，所有旧标签都要减少计数
			if oldStatus == 1 && post.Status == 0 {
				for _, oldTagID := range oldTagIDs {
					s.tagRepo.DecrementPostCount(oldTagID)
				}
			}
		}
	}

	// 更新分类文章数
	if oldCategoryID != post.CategoryID {
		if oldStatus == 1 {
			s.categoryRepo.DecrementPostCount(oldCategoryID)
		}
		if post.Status == 1 {
			s.categoryRepo.IncrementPostCount(post.CategoryID)
		}
	} else if oldStatus != post.Status {
		if post.Status == 1 {
			s.categoryRepo.IncrementPostCount(post.CategoryID)
		} else {
			s.categoryRepo.DecrementPostCount(post.CategoryID)
		}
	}

	return s.postRepo.GetByID(post.ID)
}

// Delete 删除文章
func (s *PostService) Delete(id, userID uint, role string) error {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return errors.New("文章不存在")
	}

	// 权限检查
	if post.UserID != userID && role != "admin" {
		return errors.New("无权限删除此文章")
	}

	// 减少分类文章数
	if post.Status == 1 {
		s.categoryRepo.DecrementPostCount(post.CategoryID)
		
		// 减少标签文章数
		if len(post.Tags) > 0 {
			for _, tag := range post.Tags {
				s.tagRepo.DecrementPostCount(tag.ID)
			}
		}
	}

	return s.postRepo.Delete(id)
}

// List 获取文章列表
func (s *PostService) List(page, pageSize int, categoryID uint, keyword string, status int) ([]model.Post, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.postRepo.List(page, pageSize, categoryID, keyword, status)
}

// GetByTag 根据标签获取文章
func (s *PostService) GetByTag(tagID uint, page, pageSize int) ([]model.Post, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.postRepo.GetByTag(tagID, page, pageSize)
}

// Like 点赞/取消点赞文章
func (s *PostService) Like(id uint, userID *uint, ip string) (bool, error) {
	// 检查文章是否存在
	_, err := s.postRepo.GetByID(id)
	if err != nil {
		return false, errors.New("文章不存在")
	}

	// 检查是否已点赞
	liked, err := s.postRepo.CheckLiked(id, userID, ip)
	if err != nil {
		return false, err
	}

	if liked {
		// 已点赞，执行取消点赞
		if err := s.postRepo.DeleteLike(id, userID, ip); err != nil {
			return false, err
		}
		
		// 减少点赞数
		if err := s.postRepo.DecrementLikeCount(id); err != nil {
			return false, err
		}
		
		return false, nil // 返回 false 表示取消点赞
	} else {
		// 未点赞，执行点赞
		like := &model.PostLike{
			PostID: id,
			UserID: userID,
			IP:     ip,
		}
		if err := s.postRepo.CreateLike(like); err != nil {
			return false, err
		}

		// 增加点赞数
		if err := s.postRepo.IncrementLikeCount(id); err != nil {
			return false, err
		}
		
		return true, nil // 返回 true 表示点赞成功
	}
}

// GetArchives 获取归档
func (s *PostService) GetArchives() ([]map[string]interface{}, error) {
	return s.postRepo.GetArchives()
}

// GetHotPosts 获取热门文章
func (s *PostService) GetHotPosts(limit int) ([]model.Post, error) {
	if limit < 1 || limit > 50 {
		limit = 10
	}
	return s.postRepo.GetHotPosts(limit)
}

// GetRecentPosts 获取最新文章
func (s *PostService) GetRecentPosts(limit int) ([]model.Post, error) {
	if limit < 1 || limit > 50 {
		limit = 10
	}
	return s.postRepo.GetRecentPosts(limit)
}

