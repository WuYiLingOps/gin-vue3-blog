package service

import (
	"context"
	"encoding/json"
	"time"

	"blog-backend/db"
	"blog-backend/model"
	"blog-backend/repository"
)

type FriendLinkService struct {
	repo *repository.FriendLinkRepository
}

func NewFriendLinkService() *FriendLinkService {
	return &FriendLinkService{
		repo: repository.NewFriendLinkRepository(),
	}
}

// CreateFriendLinkRequest 创建友链请求
type CreateFriendLinkRequest struct {
	Name        string `json:"name" binding:"required"`
	URL         string `json:"url" binding:"required"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Screenshot  string `json:"screenshot"`
	AtomURL     string `json:"atom_url"`
	CategoryID  uint   `json:"category_id" binding:"required"` // 分类ID（必选）
	SortOrder   int    `json:"sort_order"`
	Status      int    `json:"status"`
}

// UpdateFriendLinkRequest 更新友链请求
type UpdateFriendLinkRequest struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Screenshot  string `json:"screenshot"`
	AtomURL     string `json:"atom_url"`
	CategoryID  *uint  `json:"category_id"` // 分类ID（使用指针以区分未传递和传递了0）
	SortOrder   *int   `json:"sort_order"`
	Status      *int   `json:"status"`
}

// Create 创建友链
func (s *FriendLinkService) Create(req *CreateFriendLinkRequest) (*model.FriendLink, error) {
	// 验证分类是否存在
	categoryRepo := repository.NewFriendLinkCategoryRepository()
	if _, err := categoryRepo.GetByID(req.CategoryID); err != nil {
		return nil, err
	}

	friendLink := &model.FriendLink{
		Name:        req.Name,
		URL:         req.URL,
		Icon:        req.Icon,
		Description: req.Description,
		Screenshot:  req.Screenshot,
		AtomURL:     req.AtomURL,
		CategoryID:  req.CategoryID,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
	}

	if friendLink.Status == 0 {
		friendLink.Status = 1 // 默认启用
	}

	if err := s.repo.Create(friendLink); err != nil {
		return nil, err
	}

	// 创建成功后，清理前台友链列表缓存
	go func() {
		ctx := context.Background()
		db.RDB.Del(ctx, "friend_links:public:list")
	}()

	return s.repo.GetByID(friendLink.ID)
}

// GetByID 根据ID获取友链
func (s *FriendLinkService) GetByID(id uint) (*model.FriendLink, error) {
	return s.repo.GetByID(id)
}

// List 获取友链列表（管理员用）
func (s *FriendLinkService) List(page, pageSize int) ([]model.FriendLink, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.repo.List(page, pageSize)
}

// ListPublic 获取公开的友链列表（前端用）
func (s *FriendLinkService) ListPublic() ([]model.FriendLink, error) {
	ctx := context.Background()
	cacheKey := "friend_links:public:list"

	// 1. 先尝试从 Redis 获取缓存
	if cached, err := db.RDB.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
		var list []model.FriendLink
		if err := json.Unmarshal([]byte(cached), &list); err == nil {
			return list, nil
		}
	}

	// 2. 缓存未命中，从数据库获取
	list, err := s.repo.ListPublic()
	if err != nil {
		return nil, err
	}

	// 3. 写入缓存，设置过期时间 1 小时
	if data, err := json.Marshal(list); err == nil {
		_ = db.RDB.Set(ctx, cacheKey, string(data), time.Hour).Err()
	}

	return list, nil
}

// Update 更新友链
func (s *FriendLinkService) Update(id uint, req *UpdateFriendLinkRequest) (*model.FriendLink, error) {
	friendLink, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		friendLink.Name = req.Name
	}
	if req.URL != "" {
		friendLink.URL = req.URL
	}
	if req.Icon != "" {
		friendLink.Icon = req.Icon
	}
	if req.Description != "" {
		friendLink.Description = req.Description
	}
	if req.Screenshot != "" {
		friendLink.Screenshot = req.Screenshot
	}
	if req.AtomURL != "" {
		friendLink.AtomURL = req.AtomURL
	}
	// 如果传递了 category_id，则更新分类
	// 注意：JSON 绑定中，如果字段存在且值为数字，指针会被设置为指向该数字
	// 如果字段不存在，指针为 nil
	if req.CategoryID != nil {
		categoryID := *req.CategoryID
		if categoryID > 0 {
			// 验证分类是否存在
			categoryRepo := repository.NewFriendLinkCategoryRepository()
			if _, err := categoryRepo.GetByID(categoryID); err != nil {
				return nil, err
			}
			friendLink.CategoryID = categoryID
		} else {
			// category_id 为 0 或负数，视为无效值，不更新
			// 这种情况不应该发生，因为前端验证要求 category_id > 0
		}
	}
	// 如果 req.CategoryID == nil，说明前端没有传递 category_id 字段，不更新分类
	if req.SortOrder != nil {
		friendLink.SortOrder = *req.SortOrder
	}
	if req.Status != nil {
		friendLink.Status = *req.Status
	}

	if err := s.repo.Update(friendLink); err != nil {
		return nil, err
	}

	// 更新成功后，清理前台友链列表缓存
	go func() {
		ctx := context.Background()
		db.RDB.Del(ctx, "friend_links:public:list")
	}()

	return s.repo.GetByID(friendLink.ID)
}

// Delete 删除友链
func (s *FriendLinkService) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// 删除成功后，清理前台友链列表缓存
	go func() {
		ctx := context.Background()
		db.RDB.Del(ctx, "friend_links:public:list")
	}()

	return nil
}
