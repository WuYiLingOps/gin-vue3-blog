/*
 * 项目名称：blog-backend
 * 文件名称：friendlink.go
 * 创建时间：2026-01-31 16:29:06
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：友链数据访问层，提供友情链接的数据库操作功能
 */
package repository

import (
	"blog-backend/db"
	"blog-backend/model"
)

// FriendLinkRepository 友链数据访问层结构体
type FriendLinkRepository struct{}

// NewFriendLinkRepository 创建友链数据访问层实例
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
func (r *FriendLinkRepository) List(page, pageSize int, keyword, status, accessible string) ([]model.FriendLink, int64, error) {
	var friendLinks []model.FriendLink
	var total int64

	offset := (page - 1) * pageSize

	// 构建查询
	query := db.DB.Model(&model.FriendLink{})

	// 关键词筛选
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 链接状态筛选
	switch accessible {
	case "0": // 正常
		query = query.Where("accessible = 0 AND is_invalid = false")
	case "1": // 异常
		query = query.Where("accessible > 0")
	case "2": // 已失效
		query = query.Where("is_invalid = true")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Category").
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

// GetAllForCheck 获取所有需要检测的友链（排除忽略检查的）
func (r *FriendLinkRepository) GetAllForCheck() ([]model.FriendLink, error) {
	var friendLinks []model.FriendLink
	err := db.DB.Where("status = ? AND accessible != ?", 1, -1).
		Find(&friendLinks).Error
	return friendLinks, err
}

// UpdateCheckStatus 更新友链检测状态
func (r *FriendLinkRepository) UpdateCheckStatus(id uint, accessible int) error {
	return db.DB.Model(&model.FriendLink{}).
		Where("id = ?", id).
		Update("accessible", accessible).Error
}

// UpdateInvalidStatus 更新失效状态
func (r *FriendLinkRepository) UpdateInvalidStatus(id uint, isInvalid bool) error {
	return db.DB.Model(&model.FriendLink{}).
		Where("id = ?", id).
		Update("is_invalid", isInvalid).Error
}
