package service

import (
	"blog-backend/model"
	"blog-backend/repository"
)

type FriendLinkCategoryService struct {
	repo *repository.FriendLinkCategoryRepository
}

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

	return category, nil
}

// GetByID 根据ID获取友链分类
func (s *FriendLinkCategoryService) GetByID(id uint) (*model.FriendLinkCategory, error) {
	return s.repo.GetByID(id)
}

// List 获取所有友链分类列表
func (s *FriendLinkCategoryService) List() ([]model.FriendLinkCategory, error) {
	return s.repo.List()
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

	return category, nil
}

// Delete 删除友链分类
func (s *FriendLinkCategoryService) Delete(id uint) error {
	return s.repo.Delete(id)
}
