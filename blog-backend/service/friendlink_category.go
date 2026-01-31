/*
 * 项目名称：blog-backend
 * 文件名称：friendlink_category.go
 * 创建时间：2026-01-31 16:34:35
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：友链分类业务逻辑层，提供友链分类的增删改查业务处理，支持Redis缓存
 */
package service

import (
	"context"
	"encoding/json"
	"time"

	"blog-backend/db"
	"blog-backend/model"
	"blog-backend/repository"
)

// FriendLinkCategoryService 友链分类业务逻辑层结构体
type FriendLinkCategoryService struct {
	repo *repository.FriendLinkCategoryRepository
}

// NewFriendLinkCategoryService 创建友链分类业务逻辑层实例
func NewFriendLinkCategoryService() *FriendLinkCategoryService {
	return &FriendLinkCategoryService{
		repo: repository.NewFriendLinkCategoryRepository(),
	}
}

// CreateFriendLinkCategoryRequest 创建友链分类请求
type CreateFriendLinkCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// UpdateFriendLinkCategoryRequest 更新友链分类请求
type UpdateFriendLinkCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// Create 创建友链分类
func (s *FriendLinkCategoryService) Create(req *CreateFriendLinkCategoryRequest) (*model.FriendLinkCategory, error) {
	category := &model.FriendLinkCategory{
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}

	if err := s.repo.Create(category); err != nil {
		return nil, err
	}

	// 创建成功后，清理友链分类列表和前台友链列表缓存
	go func() {
		ctx := context.Background()
		db.RDB.Del(ctx, "friendlink_category:list")
		db.RDB.Del(ctx, "friend_links:public:list")
	}()

	return category, nil
}

// GetByID 根据ID获取友链分类
func (s *FriendLinkCategoryService) GetByID(id uint) (*model.FriendLinkCategory, error) {
	return s.repo.GetByID(id)
}

// List 获取所有友链分类列表
func (s *FriendLinkCategoryService) List() ([]model.FriendLinkCategory, error) {
	ctx := context.Background()
	cacheKey := "friendlink_category:list"

	// 1. 先尝试从 Redis 获取缓存
	if cached, err := db.RDB.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
		var list []model.FriendLinkCategory
		if err := json.Unmarshal([]byte(cached), &list); err == nil {
			return list, nil
		}
	}

	// 2. 缓存未命中，从数据库获取
	list, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	// 3. 写入缓存，设置过期时间 1 小时
	if data, err := json.Marshal(list); err == nil {
		_ = db.RDB.Set(ctx, cacheKey, string(data), time.Hour).Err()
	}

	return list, nil
}

// Update 更新友链分类
func (s *FriendLinkCategoryService) Update(id uint, req *UpdateFriendLinkCategoryRequest) (*model.FriendLinkCategory, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	category.Description = req.Description
	category.SortOrder = req.SortOrder

	if err := s.repo.Update(category); err != nil {
		return nil, err
	}

	// 更新成功后，清理友链分类列表和前台友链列表缓存
	go func() {
		ctx := context.Background()
		db.RDB.Del(ctx, "friendlink_category:list")
		db.RDB.Del(ctx, "friend_links:public:list")
	}()

	return category, nil
}

// Delete 删除友链分类
func (s *FriendLinkCategoryService) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// 删除成功后，清理友链分类列表和前台友链列表缓存
	go func() {
		ctx := context.Background()
		db.RDB.Del(ctx, "friendlink_category:list")
		db.RDB.Del(ctx, "friend_links:public:list")
	}()

	return nil
}
