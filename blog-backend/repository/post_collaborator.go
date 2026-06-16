/*
 * 项目名称：blog-backend
 * 文件名称：post_collaborator.go
 * 创建时间：2026-06-16
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：文章协作者数据访问层，提供文章协作者的数据库操作功能
 */
package repository

import (
	"blog-backend/db"
	"blog-backend/model"
)

// PostCollaboratorRepository 文章协作者数据访问层结构体
type PostCollaboratorRepository struct{}

// NewPostCollaboratorRepository 创建文章协作者数据访问层实例
func NewPostCollaboratorRepository() *PostCollaboratorRepository {
	return &PostCollaboratorRepository{}
}

// GetByPostID 根据文章ID获取协作者列表（仅返回未移除的）
func (r *PostCollaboratorRepository) GetByPostID(postID uint) ([]model.User, error) {
	var users []model.User
	err := db.DB.
		Joins("JOIN post_collaborators ON post_collaborators.user_id = users.id").
		Where("post_collaborators.post_id = ? AND post_collaborators.removed = false", postID).
		Order("post_collaborators.sort_order ASC").
		Find(&users).Error
	return users, err
}

// Add 添加文章协作者（如果之前被移除则恢复，否则新建）
func (r *PostCollaboratorRepository) Add(postID, userID uint, sortOrder int) error {
	var existing model.PostCollaborator
	err := db.DB.
		Where("post_id = ? AND user_id = ?", postID, userID).
		First(&existing).Error

	if err == nil {
		// 记录存在，更新为未移除状态
		return db.DB.Model(&existing).Update("removed", false).Error
	}

	// 记录不存在，创建新记录
	collaborator := model.PostCollaborator{
		PostID:    postID,
		UserID:    userID,
		SortOrder: sortOrder,
		Removed:   false,
	}
	return db.DB.Create(&collaborator).Error
}

// Remove 移除文章协作者（软删除）
func (r *PostCollaboratorRepository) Remove(postID, userID uint) error {
	return db.DB.
		Model(&model.PostCollaborator{}).
		Where("post_id = ? AND user_id = ?", postID, userID).
		Update("removed", true).Error
}

// RemoveByPostID 移除文章所有协作者（软删除）
func (r *PostCollaboratorRepository) RemoveByPostID(postID uint) error {
	return db.DB.
		Model(&model.PostCollaborator{}).
		Where("post_id = ?", postID).
		Update("removed", true).Error
}

// CountByPostID 获取文章协作者数量（仅统计未移除的）
func (r *PostCollaboratorRepository) CountByPostID(postID uint) (int64, error) {
	var count int64
	err := db.DB.
		Model(&model.PostCollaborator{}).
		Where("post_id = ? AND removed = false", postID).
		Count(&count).Error
	return count, err
}

// Exists 检查用户是否为文章协作者（未移除状态）
func (r *PostCollaboratorRepository) Exists(postID, userID uint) (bool, error) {
	var count int64
	err := db.DB.
		Model(&model.PostCollaborator{}).
		Where("post_id = ? AND user_id = ? AND removed = false", postID, userID).
		Count(&count).Error
	return count > 0, err
}

// EverExisted 检查用户是否曾经是文章协作者（包括被移除的）
func (r *PostCollaboratorRepository) EverExisted(postID, userID uint) (bool, error) {
	var count int64
	err := db.DB.
		Model(&model.PostCollaborator{}).
		Where("post_id = ? AND user_id = ?", postID, userID).
		Count(&count).Error
	return count > 0, err
}

// UpdateSortOrder 更新协作者排序
func (r *PostCollaboratorRepository) UpdateSortOrder(postID, userID uint, sortOrder int) error {
	return db.DB.
		Model(&model.PostCollaborator{}).
		Where("post_id = ? AND user_id = ?", postID, userID).
		Update("sort_order", sortOrder).Error
}
