package repository

import (
	"blog-backend/db"
	"blog-backend/model"
	"errors"

	"gorm.io/gorm"
)

type PostViewRepository struct{}

func NewPostViewRepository() *PostViewRepository {
	return &PostViewRepository{}
}

// HasViewed 检查是否已经阅读过
// 优先检查用户ID，如果没有则检查IP
func (r *PostViewRepository) HasViewed(postID uint, userID *uint, ip string) (bool, error) {
	var count int64
	query := db.DB.Model(&model.PostView{}).Where("post_id = ?", postID)

	// 如果是登录用户，按用户ID查询
	if userID != nil && *userID > 0 {
		query = query.Where("user_id = ?", *userID)
	} else {
		// 匿名用户按IP查询
		query = query.Where("ip = ? AND user_id IS NULL", ip)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// RecordView 记录文章阅读
func (r *PostViewRepository) RecordView(postID uint, userID *uint, ip string) error {
	// 先检查是否已阅读
	hasViewed, err := r.HasViewed(postID, userID, ip)
	if err != nil {
		return err
	}

	if hasViewed {
		return errors.New("already viewed")
	}

	// 记录阅读
	postView := model.PostView{
		PostID: postID,
		UserID: userID,
		IP:     ip,
	}

	return db.DB.Create(&postView).Error
}

// IncrementViewCount 增加文章阅读量
func (r *PostViewRepository) IncrementViewCount(postID uint) error {
	return db.DB.Model(&model.Post{}).
		Where("id = ?", postID).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).
		Error
}

