package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"blog-backend/db"
	"blog-backend/model"
	"blog-backend/repository"

	"gorm.io/gorm"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService() *CategoryService {
	return &CategoryService{
		repo: repository.NewCategoryRepository(),
	}
}

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Sort        int    `json:"sort"`
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Sort        int    `json:"sort"`
}

// Create 创建分类
func (s *CategoryService) Create(req *CreateCategoryRequest) (*model.Category, error) {
	// 检查分类名是否已存在
	if _, err := s.repo.GetByName(req.Name); err == nil {
		return nil, errors.New("分类名已存在")
	}

	category := &model.Category{
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		Sort:        req.Sort,
	}

	if err := s.repo.Create(category); err != nil {
		return nil, errors.New("分类创建失败")
	}

	// 写操作成功后，删除分类列表相关缓存
	go func() {
		ctx := context.Background()
		db.RDB.Del(ctx, "category:list")
		db.RDB.Del(ctx, "blog:author_profile")
	}()

	return category, nil
}

// GetByID 获取分类详情
func (s *CategoryService) GetByID(id uint) (*model.Category, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分类不存在")
		}
		return nil, errors.New("获取分类失败")
	}
	return category, nil
}

// Update 更新分类
func (s *CategoryService) Update(id uint, req *UpdateCategoryRequest) (*model.Category, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("分类不存在")
	}

	// 如果修改了名称，检查新名称是否已存在
	if req.Name != "" && req.Name != category.Name {
		if _, err := s.repo.GetByName(req.Name); err == nil {
			return nil, errors.New("分类名已存在")
		}
		category.Name = req.Name
	}

	if req.Description != "" {
		category.Description = req.Description
	}
	if req.Color != "" {
		category.Color = req.Color
	}
	category.Sort = req.Sort

	if err := s.repo.Update(category); err != nil {
		return nil, errors.New("分类更新失败")
	}

	// 写操作成功后，删除分类列表相关缓存
	go func() {
		ctx := context.Background()
		db.RDB.Del(ctx, "category:list")
		db.RDB.Del(ctx, "blog:author_profile")
	}()

	return category, nil
}

// Delete 删除分类
func (s *CategoryService) Delete(id uint) error {
	category, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("分类不存在")
	}

	// 检查分类下是否有文章
	if category.PostCount > 0 {
		return errors.New("该分类下还有文章，无法删除")
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// 删除成功后，清理分类相关缓存
	go func() {
		ctx := context.Background()
		db.RDB.Del(ctx, "category:list")
		db.RDB.Del(ctx, "blog:author_profile")
	}()

	return nil
}

// List 获取分类列表
func (s *CategoryService) List() ([]model.Category, error) {
	ctx := context.Background()
	cacheKey := "category:list"

	// 1. 先尝试从 Redis 获取缓存
	if cached, err := db.RDB.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
		var list []model.Category
		if err := json.Unmarshal([]byte(cached), &list); err == nil {
			return list, nil
		}
	}

	// 2. 缓存未命中，从数据库获取
	list, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	// 3. 写入缓存，设置过期时间（例如 30 分钟）
	if data, err := json.Marshal(list); err == nil {
		_ = db.RDB.Set(ctx, cacheKey, string(data), 30*time.Minute).Err()
	}

	return list, nil
}
