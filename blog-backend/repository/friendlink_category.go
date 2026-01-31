/*
 * 项目名称：blog-backend
 * 文件名称：friendlink_category.go
 * 创建时间：2026-01-31 16:29:06
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：友链分类数据访问层，提供友链分类的数据库操作功能
 */
package repository

import (
	"blog-backend/db"
	"blog-backend/model"
)

// FriendLinkCategoryRepository 友链分类数据访问层结构体
type FriendLinkCategoryRepository struct{}

// NewFriendLinkCategoryRepository 创建友链分类数据访问层实例
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
