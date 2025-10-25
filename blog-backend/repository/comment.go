package repository

import (
	"blog-backend/db"
	"blog-backend/model"
	"gorm.io/gorm"
)

type CommentRepository struct{}

func NewCommentRepository() *CommentRepository {
	return &CommentRepository{}
}

// Create 创建评论
func (r *CommentRepository) Create(comment *model.Comment) error {
	return db.DB.Create(comment).Error
}

// GetByID 根据ID获取评论
func (r *CommentRepository) GetByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	err := db.DB.Preload("User").Preload("Parent").First(&comment, id).Error
	return &comment, err
}

// Update 更新评论
func (r *CommentRepository) Update(comment *model.Comment) error {
	return db.DB.Save(comment).Error
}

// Delete 删除评论
func (r *CommentRepository) Delete(id uint) error {
	return db.DB.Delete(&model.Comment{}, id).Error
}

// GetByPostID 根据文章ID获取评论列表
func (r *CommentRepository) GetByPostID(postID uint) ([]model.Comment, error) {
	var comments []model.Comment
	// 主评论按时间倒序，子评论按时间正序
	err := db.DB.Preload("User").
		Preload("Children", func(tx *gorm.DB) *gorm.DB {
			return tx.Where("status = 1").Order("created_at ASC")
		}).
		Preload("Children.User").
		Where("post_id = ? AND parent_id IS NULL AND status = 1", postID).
		Order("created_at DESC").
		Find(&comments).Error
	return comments, err
}

// List 获取评论列表（管理后台）
func (r *CommentRepository) List(page, pageSize int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	offset := (page - 1) * pageSize

	if err := db.DB.Model(&model.Comment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.DB.Preload("User").Preload("Post").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&comments).Error

	return comments, total, err
}

// UpdateStatus 更新评论状态
func (r *CommentRepository) UpdateStatus(id uint, status int) error {
	return db.DB.Model(&model.Comment{}).Where("id = ?", id).Update("status", status).Error
}

// GetRecentComments 获取最新评论
func (r *CommentRepository) GetRecentComments(limit int) ([]model.Comment, error) {
	var comments []model.Comment
	err := db.DB.Preload("User").Preload("Post").
		Where("status = 1").
		Order("created_at DESC").
		Limit(limit).
		Find(&comments).Error
	return comments, err
}

