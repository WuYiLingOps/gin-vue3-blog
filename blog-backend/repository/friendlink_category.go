package repository

import (
	"blog-backend/db"
	"blog-backend/model"
)

type FriendLinkCategoryRepository struct{}

func NewFriendLinkCategoryRepository() *FriendLinkCategoryRepository {
	return &FriendLinkCategoryRepository{}
}

// Create 创建友链分类
func (r *FriendLinkCategoryRepository) Create(category *model.FriendLinkCategory) error {
	return db.DB.Create(category).Error
}

// GetByID 根据ID获取友链分类
func (r *FriendLinkCategoryRepository) GetByID(id uint) (*model.FriendLinkCategory, error) {
	var category model.FriendLinkCategory
	err := db.DB.First(&category, id).Error
	return &category, err
}

// List 获取所有友链分类列表
func (r *FriendLinkCategoryRepository) List() ([]model.FriendLinkCategory, error) {
	var categories []model.FriendLinkCategory
	err := db.DB.Order("sort_order DESC, id DESC").Find(&categories).Error
	return categories, err
}

// Update 更新友链分类
func (r *FriendLinkCategoryRepository) Update(category *model.FriendLinkCategory) error {
	return db.DB.Save(category).Error
}

// Delete 删除友链分类
func (r *FriendLinkCategoryRepository) Delete(id uint) error {
	return db.DB.Delete(&model.FriendLinkCategory{}, id).Error
}
