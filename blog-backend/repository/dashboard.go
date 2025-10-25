package repository

import (
	"blog-backend/db"
	"blog-backend/model"
)

// GetPublishedCount 获取已发布文章总数
func (r *PostRepository) GetPublishedCount() (int64, error) {
	var count int64
	err := db.DB.Model(&model.Post{}).Where("status = 1").Count(&count).Error
	return count, err
}

// GetTotalViews 获取总浏览量
func (r *PostRepository) GetTotalViews() (int64, error) {
	var total int64
	err := db.DB.Model(&model.Post{}).Where("status = 1").Select("COALESCE(SUM(view_count), 0)").Scan(&total).Error
	return total, err
}

// GetTotalCount 获取用户总数
func (r *UserRepository) GetTotalCount() (int64, error) {
	var count int64
	err := db.DB.Model(&model.User{}).Count(&count).Error
	return count, err
}

// GetTotalCount 获取评论总数
func (r *CommentRepository) GetTotalCount() (int64, error) {
	var count int64
	err := db.DB.Model(&model.Comment{}).Where("status = 1").Count(&count).Error
	return count, err
}

