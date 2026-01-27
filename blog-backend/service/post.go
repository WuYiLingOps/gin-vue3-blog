package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"blog-backend/db"
	"blog-backend/model"
	"blog-backend/repository"
	"blog-backend/util"

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
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Summary    string `json:"summary"`
	Cover      string `json:"cover"`
	CategoryID uint   `json:"category_id" binding:"required"`
	TagIDs     []uint `json:"tag_ids"`
	Status     int    `json:"status"`     // 0:草稿 1:发布
	Visibility int    `json:"visibility"` // 1:公开 0:私密（默认值 1 在前端设置）
	IsTop      bool   `json:"is_top"`
}

// UpdatePostRequest 更新文章请求
type UpdatePostRequest struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	Summary    string `json:"summary"`
	Cover      string `json:"cover"`
	CategoryID *uint  `json:"category_id"` // 使用指针类型，nil 表示不修改
	TagIDs     []uint `json:"tag_ids"`
	Status     int    `json:"status"`
	Visibility *int   `json:"visibility"` // 1:公开 0:私密（nil 表示不修改）
	IsTop      bool   `json:"is_top"`
}

// Create 创建文章
func (s *PostService) Create(userID uint, req *CreatePostRequest) (*model.Post, error) {
	// 检查分类是否存在
	if _, err := s.categoryRepo.GetByID(req.CategoryID); err != nil {
		return nil, errors.New("分类不存在")
	}

	// 处理状态：确保草稿状态（0）能正确保存
	// 如果 status 为 0（草稿），需要明确设置，避免被默认值覆盖
	postStatus := req.Status
	if postStatus != 0 && postStatus != 1 {
		// 如果状态值无效，默认为发布（1）
		postStatus = 1
	}

	// 处理可见性：
	// - 如果状态为草稿（0），自动设置为私密（0）
	// - 如果状态为发布（1），使用用户选择的可见性，默认公开（1）
	var visibility int
	if postStatus == 0 {
		// 草稿状态，强制设置为私密
		visibility = 0
	} else {
		// 发布状态，使用用户选择的可见性
		visibility = req.Visibility
		if visibility != 0 && visibility != 1 {
			// 如果可见性值无效，默认为公开（1）
			visibility = 1
		}
	}

	// 生成slug
	baseSlug := util.GenerateSlug(req.Title)
	if baseSlug == "" {
		return nil, errors.New("无法生成文章slug，请检查标题")
	}

	postRepo := s.postRepo // 在闭包外捕获
	slug := util.GenerateUniqueSlug(baseSlug, func(slug string) bool {
		return postRepo.CheckSlugExists(slug, 0)
	})

	if slug == "" {
		return nil, errors.New("生成的slug为空")
	}

	post := &model.Post{
		Title:      req.Title,
		Slug:       slug,
		Content:    req.Content,
		Summary:    req.Summary,
		Cover:      req.Cover,
		CategoryID: req.CategoryID,
		Status:     postStatus,
		Visibility: visibility,
		IsTop:      req.IsTop,
		UserID:     userID,
	}

	// 如果是发布状态，设置发布时间
	if req.Status == 1 {
		now := time.Now()
		post.PublishedAt = &now
	}

	// 使用事务确保数据一致性
	err := s.postRepo.Transaction(func(tx *gorm.DB) error {
		// 创建文章
		if err := s.postRepo.CreateTx(tx, post); err != nil {
			// 如果是唯一性约束错误，返回更友好的错误信息
			errStr := err.Error()
			if strings.Contains(errStr, "duplicate key") ||
				strings.Contains(errStr, "unique constraint") ||
				strings.Contains(errStr, "violates unique constraint") {
				return errors.New("文章slug已存在，请修改标题后重试")
			}
			// 返回原始错误，包含更多调试信息
			return errors.New("创建文章失败: " + errStr)
		}

		// 更新标签关联
		if len(req.TagIDs) > 0 {
			if err := s.postRepo.UpdateTagsTx(tx, post.ID, req.TagIDs); err != nil {
				return err
			}

			// 增加标签文章数（仅发布状态）
			if req.Status == 1 {
				for _, tagID := range req.TagIDs {
					if err := s.tagRepo.IncrementPostCountTx(tx, tagID); err != nil {
						return err
					}
				}
			}
		}

		// 增加分类文章数
		if req.Status == 1 {
			if err := s.categoryRepo.IncrementPostCountTx(tx, req.CategoryID); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, errors.New("文章创建失败: " + err.Error())
	}

	// 写操作成功后，删除与文章列表相关的缓存（最新文章等）
	go func() {
		ctx := context.Background()
		// 删除常用的最新文章缓存（不同 limit 可以按需扩展）
		for _, limit := range []int{5, 10, 20} {
			key := fmt.Sprintf("post:recent:%d", limit)
			db.RDB.Del(ctx, key)
		}
		// 文章数、标签统计等也会受影响，清理相关缓存
		db.RDB.Del(ctx, "blog:author_profile")
		db.RDB.Del(ctx, "tag:stats:top10")
	}()

	return s.postRepo.GetByID(post.ID)
}

// GetByID 获取文章详情（含权限校验）
func (s *PostService) GetByID(id uint, userID *uint, role string, ip string) (*model.Post, error) {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, errors.New("获取文章失败")
	}

	return s.checkPostPermission(post, userID, role, ip)
}

// GetBySlug 根据slug获取文章详情（含权限校验）
func (s *PostService) GetBySlug(slug string, userID *uint, role string, ip string) (*model.Post, error) {
	post, err := s.postRepo.GetBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, errors.New("获取文章失败")
	}

	return s.checkPostPermission(post, userID, role, ip)
}

// checkPostPermission 检查文章权限并记录浏览
func (s *PostService) checkPostPermission(post *model.Post, userID *uint, role string, ip string) (*model.Post, error) {

	// 私密/草稿仅作者或管理员可见
	if (post.Visibility == 0 || post.Status == 0) && role != "admin" {
		if userID == nil || *userID != post.UserID {
			return nil, errors.New("无权限查看")
		}
	}

	// 检查是否已阅读，如果没有则记录并增加浏览量
	if ip != "" && ip != "unknown" {
		hasViewed, _ := s.postViewRepo.HasViewed(post.ID, userID, ip)
		if !hasViewed {
			// 记录阅读
			if err := s.postViewRepo.RecordView(post.ID, userID, ip); err == nil {
				// 增加浏览量
				s.postViewRepo.IncrementViewCount(post.ID)
				post.ViewCount++
			}
		}
	}

	// 检查是否已点赞
	liked, _ := s.postRepo.CheckLiked(post.ID, userID, ip)
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
		// 如果标题改变，重新生成slug
		baseSlug := util.GenerateSlug(req.Title)
		postRepo := s.postRepo // 在闭包外捕获
		postID := post.ID      // 在闭包外捕获
		post.Slug = util.GenerateUniqueSlug(baseSlug, func(slug string) bool {
			return postRepo.CheckSlugExists(slug, postID)
		})
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
	if req.CategoryID != nil {
		// 检查新分类是否存在
		if _, err := s.categoryRepo.GetByID(*req.CategoryID); err != nil {
			return nil, errors.New("分类不存在")
		}
		post.CategoryID = *req.CategoryID
	}
	post.Status = req.Status
	post.IsTop = req.IsTop

	// 更新可见性：
	// - 如果状态为草稿（0），自动设置为私密（0）
	// - 如果状态为发布（1），使用用户选择的可见性（如果传入）
	if req.Status == 0 {
		// 草稿状态，强制设置为私密
		post.Visibility = 0
	} else if req.Visibility != nil {
		// 发布状态，使用用户选择的可见性
		if *req.Visibility != 0 && *req.Visibility != 1 {
			return nil, errors.New("可见性参数错误")
		}
		post.Visibility = *req.Visibility
	}
	// 如果状态为发布且未传入可见性，保持原有可见性不变

	// 如果从草稿变为发布，设置发布时间
	if oldStatus == 0 && req.Status == 1 {
		now := time.Now()
		post.PublishedAt = &now
	}

	// 使用事务确保数据一致性
	err = s.postRepo.Transaction(func(tx *gorm.DB) error {
		// 更新文章
		if err := s.postRepo.UpdateTx(tx, post); err != nil {
			return err
		}

		// 更新标签关联
		if len(req.TagIDs) > 0 {
			if err := s.postRepo.UpdateTagsTx(tx, post.ID, req.TagIDs); err != nil {
				return err
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
						if err := s.tagRepo.DecrementPostCountTx(tx, oldTagID); err != nil {
							return err
						}
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
						if err := s.tagRepo.IncrementPostCountTx(tx, newTagID); err != nil {
							return err
						}
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
							if err := s.tagRepo.IncrementPostCountTx(tx, tagID); err != nil {
								return err
							}
						}
					}
				}

				// 如果状态从发布变为草稿，所有旧标签都要减少计数
				if oldStatus == 1 && post.Status == 0 {
					for _, oldTagID := range oldTagIDs {
						if err := s.tagRepo.DecrementPostCountTx(tx, oldTagID); err != nil {
							return err
						}
					}
				}
			}
		}

		// 更新分类文章数
		if oldCategoryID != post.CategoryID {
			if oldStatus == 1 {
				if err := s.categoryRepo.DecrementPostCountTx(tx, oldCategoryID); err != nil {
					return err
				}
			}
			if post.Status == 1 {
				if err := s.categoryRepo.IncrementPostCountTx(tx, post.CategoryID); err != nil {
					return err
				}
			}
		} else if oldStatus != post.Status {
			if post.Status == 1 {
				if err := s.categoryRepo.IncrementPostCountTx(tx, post.CategoryID); err != nil {
					return err
				}
			} else {
				if err := s.categoryRepo.DecrementPostCountTx(tx, post.CategoryID); err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, errors.New("文章更新失败")
	}

	// 写操作成功后，删除与文章列表相关的缓存（最新文章等）
	go func() {
		ctx := context.Background()
		for _, limit := range []int{5, 10, 20} {
			key := fmt.Sprintf("post:recent:%d", limit)
			db.RDB.Del(ctx, key)
		}
		db.RDB.Del(ctx, "blog:author_profile")
		db.RDB.Del(ctx, "tag:stats:top10")
	}()

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

	// 使用事务确保数据一致性
	return s.postRepo.Transaction(func(tx *gorm.DB) error {
		// 减少分类文章数
		if post.Status == 1 {
			if err := s.categoryRepo.DecrementPostCountTx(tx, post.CategoryID); err != nil {
				return err
			}

			// 减少标签文章数
			if len(post.Tags) > 0 {
				for _, tag := range post.Tags {
					if err := s.tagRepo.DecrementPostCountTx(tx, tag.ID); err != nil {
						return err
					}

					// 删除成功后，清理与文章列表相关的缓存
					go func() {
						ctx := context.Background()
						for _, limit := range []int{5, 10, 20} {
							key := fmt.Sprintf("post:recent:%d", limit)
							db.RDB.Del(ctx, key)
						}
						db.RDB.Del(ctx, "blog:author_profile")
						db.RDB.Del(ctx, "tag:stats:top10")
					}()

					return nil
				}
			}
		}

		// 删除文章
		return s.postRepo.DeleteTx(tx, id)
	})
}

// List 获取文章列表
func (s *PostService) List(page, pageSize int, categoryID uint, keyword string, status int, visibility *int) ([]model.Post, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.postRepo.List(page, pageSize, categoryID, keyword, status, visibility)
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

	// 使用事务确保数据一致性
	var isLiked bool
	err = s.postRepo.Transaction(func(tx *gorm.DB) error {
		if liked {
			// 已点赞，执行取消点赞
			if err := s.postRepo.DeleteLikeTx(tx, id, userID, ip); err != nil {
				return err
			}

			// 减少点赞数
			if err := s.postRepo.DecrementLikeCountTx(tx, id); err != nil {
				return err
			}

			isLiked = false
		} else {
			// 未点赞，执行点赞
			like := &model.PostLike{
				PostID: id,
				UserID: userID,
				IP:     ip,
			}
			if err := s.postRepo.CreateLikeTx(tx, like); err != nil {
				return err
			}

			// 增加点赞数
			if err := s.postRepo.IncrementLikeCountTx(tx, id); err != nil {
				return err
			}

			isLiked = true
		}
		return nil
	})

	if err != nil {
		return false, err
	}

	return isLiked, nil
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
func (s *PostService) GetRecentPosts(limit int, userID *uint, role string) ([]model.Post, error) {
	if limit < 1 || limit > 50 {
		limit = 10
	}
	// 对公开接口的最新文章列表做缓存（管理员视图通常不缓存）
	if role == "" || role == "user" {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("post:recent:%d", limit)

		// 1. 先尝试从 Redis 获取缓存
		if cached, err := db.RDB.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
			var posts []model.Post
			if err := json.Unmarshal([]byte(cached), &posts); err == nil {
				return posts, nil
			}
			// 解析失败则继续走数据库查询
		}

		// 2. 缓存未命中，从数据库获取
		posts, err := s.postRepo.GetRecentPosts(limit, userID, role)
		if err != nil {
			return nil, err
		}

		// 3. 将结果写入 Redis，设置适当过期时间（例如 10 分钟）
		if data, err := json.Marshal(posts); err == nil {
			_ = db.RDB.Set(ctx, cacheKey, string(data), 10*time.Minute).Err()
		}

		return posts, nil
	}

	// 管理员等角色直接走数据库，避免缓存带来的视图差异
	return s.postRepo.GetRecentPosts(limit, userID, role)
}

// GetByIDForAdmin 管理端获取文章（不计浏览、无权限限制）
func (s *PostService) GetByIDForAdmin(id uint) (*model.Post, error) {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, errors.New("获取文章失败")
	}
	return post, nil
}
