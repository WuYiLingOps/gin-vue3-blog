package repository

import (
	"blog-backend/db"
	"blog-backend/model"
)

type FriendLinkRepository struct{}

func NewFriendLinkRepository() *FriendLinkRepository {
	return &FriendLinkRepository{}
}

// Create 创建友链
func (r *FriendLinkRepository) Create(friendLink *model.FriendLink) error {
	return db.DB.Create(friendLink).Error
}

// GetByID 根据ID获取友链
func (r *FriendLinkRepository) GetByID(id uint) (*model.FriendLink, error) {
	var friendLink model.FriendLink
	err := db.DB.Preload("Category").First(&friendLink, id).Error
	return &friendLink, err
}

// List 获取友链列表（管理员用，包含所有状态）
func (r *FriendLinkRepository) List(page, pageSize int) ([]model.FriendLink, int64, error) {
	var friendLinks []model.FriendLink
	var total int64

	offset := (page - 1) * pageSize

	if err := db.DB.Model(&model.FriendLink{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.DB.Preload("Category").
		Order("category_id ASC, sort_order DESC, id DESC").
		Offset(offset).Limit(pageSize).
		Find(&friendLinks).Error

	return friendLinks, total, err
}

// ListPublic 获取公开的友链列表（前端用，只返回启用的，按分类分组）
func (r *FriendLinkRepository) ListPublic() ([]model.FriendLink, error) {
	var friendLinks []model.FriendLink
	err := db.DB.Preload("Category").
		Where("status = ?", 1).
		Order("category_id ASC, sort_order DESC, id DESC").
		Find(&friendLinks).Error
	return friendLinks, err
}

// ListByCategory 根据分类ID获取友链列表
func (r *FriendLinkRepository) ListByCategory(categoryID uint) ([]model.FriendLink, error) {
	var friendLinks []model.FriendLink
	err := db.DB.Where("category_id = ? AND status = ?", categoryID, 1).
		Order("sort_order DESC, id DESC").
		Find(&friendLinks).Error
	return friendLinks, err
}

// Update 更新友链
func (r *FriendLinkRepository) Update(friendLink *model.FriendLink) error {
	// 使用 Select 明确指定要更新的字段，确保 category_id 被更新
	return db.DB.Model(friendLink).
		Select("name", "url", "icon", "description", "screenshot", "atom_url", "category_id", "sort_order", "status", "updated_at").
		Updates(friendLink).Error
}

// Delete 删除友链
func (r *FriendLinkRepository) Delete(id uint) error {
	return db.DB.Delete(&model.FriendLink{}, id).Error
}
