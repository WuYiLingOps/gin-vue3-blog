/*
 * 项目名称：blog-backend
 * 文件名称：post_view.go
 * 创建时间：2026-01-31 16:29:06
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：文章阅读记录数据访问层，提供文章访问记录和统计的数据库操作功能
 */
package repository

import (
	"blog-backend/db"
	"blog-backend/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

// PostViewRepository 文章阅读记录数据访问层结构体
type PostViewRepository struct{}

// NewPostViewRepository 创建文章阅读记录数据访问层实例
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

// VisitStat 按天统计访问量结果
type VisitStat struct {
	Date  time.Time
	Count int64
}

// GetVisitStats 获取指定时间范围内按天聚合的访问量统计
func (r *PostViewRepository) GetVisitStats(start, end time.Time) ([]VisitStat, error) {
	var results []VisitStat

	err := db.DB.Model(&model.PostView{}).
		Select("DATE(created_at) AS date, COUNT(*) AS count").
		Where("created_at >= ? AND created_at < ?", start, end).
		Group("DATE(created_at)").
		Order("DATE(created_at)").
		Scan(&results).Error

	return results, err
}
