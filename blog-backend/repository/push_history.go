/*
 * 项目名称：blog-backend
 * 文件名称：push_history.go
 * 创建时间：2026-04-28 13:10:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：推送历史记录数据访问层，提供推送历史的增删改查操作
 */
package repository

import (
	"blog-backend/db"
	"blog-backend/model"
	"context"

	"gorm.io/gorm"
)

// PushHistoryRepository 推送历史记录仓库结构体
type PushHistoryRepository struct{}

// NewPushHistoryRepository 创建推送历史记录仓库实例
func NewPushHistoryRepository() *PushHistoryRepository {
	return &PushHistoryRepository{}
}

// Create 创建推送历史记录
func (r *PushHistoryRepository) Create(ctx context.Context, history *model.PushHistory) error {
	return db.DB.WithContext(ctx).Create(history).Error
}

// Update 更新推送历史记录
func (r *PushHistoryRepository) Update(ctx context.Context, history *model.PushHistory) error {
	return db.DB.WithContext(ctx).Save(history).Error
}

// GetByID 根据ID获取推送历史记录
func (r *PushHistoryRepository) GetByID(ctx context.Context, id uint) (*model.PushHistory, error) {
	var history model.PushHistory
	err := db.DB.WithContext(ctx).Preload("Details").First(&history, id).Error
	return &history, err
}

// List 获取推送历史记录列表（分页）
func (r *PushHistoryRepository) List(ctx context.Context, page, pageSize int) ([]model.PushHistory, int64, error) {
	var histories []model.PushHistory
	var total int64

	offset := (page - 1) * pageSize

	// 统计总数
	if err := db.DB.WithContext(ctx).Model(&model.PushHistory{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询列表
	err := db.DB.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&histories).Error

	return histories, total, err
}

// GetByPostID 根据文章ID获取推送历史记录
func (r *PushHistoryRepository) GetByPostID(ctx context.Context, postID uint) ([]model.PushHistory, error) {
	var histories []model.PushHistory
	err := db.DB.WithContext(ctx).
		Where("post_id = ?", postID).
		Order("created_at DESC").
		Find(&histories).Error
	return histories, err
}

// Delete 删除推送历史记录
func (r *PushHistoryRepository) Delete(ctx context.Context, id uint) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除推送详情
		if err := tx.Where("push_history_id = ?", id).Delete(&model.PushDetail{}).Error; err != nil {
			return err
		}
		// 删除推送历史
		return tx.Delete(&model.PushHistory{}, id).Error
	})
}

// CreateDetail 创建推送详情
func (r *PushHistoryRepository) CreateDetail(ctx context.Context, detail *model.PushDetail) error {
	return db.DB.WithContext(ctx).Create(detail).Error
}

// UpdateDetail 更新推送详情
func (r *PushHistoryRepository) UpdateDetail(ctx context.Context, detail *model.PushDetail) error {
	return db.DB.WithContext(ctx).Save(detail).Error
}

// GetDetailsByHistoryID 根据推送历史ID获取推送详情列表
func (r *PushHistoryRepository) GetDetailsByHistoryID(ctx context.Context, historyID uint, page, pageSize int) ([]model.PushDetail, int64, error) {
	var details []model.PushDetail
	var total int64

	offset := (page - 1) * pageSize

	// 统计总数
	if err := db.DB.WithContext(ctx).Model(&model.PushDetail{}).Where("push_history_id = ?", historyID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询列表
	err := db.DB.WithContext(ctx).
		Where("push_history_id = ?", historyID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&details).Error

	return details, total, err
}

// GetStats 获取推送统计信息
func (r *PushHistoryRepository) GetStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总推送次数
	var totalCount int64
	db.DB.WithContext(ctx).Model(&model.PushHistory{}).Count(&totalCount)
	stats["total_count"] = totalCount

	// 总推送成功数
	var totalSuccess int64
	db.DB.WithContext(ctx).Model(&model.PushHistory{}).Select("COALESCE(SUM(success_count), 0)").Scan(&totalSuccess)
	stats["total_success"] = totalSuccess

	// 总推送失败数
	var totalFailed int64
	db.DB.WithContext(ctx).Model(&model.PushHistory{}).Select("COALESCE(SUM(failed_count), 0)").Scan(&totalFailed)
	stats["total_failed"] = totalFailed

	// 最近一次推送时间
	var lastPush model.PushHistory
	if err := db.DB.WithContext(ctx).Order("created_at DESC").First(&lastPush).Error; err == nil {
		stats["last_push_at"] = lastPush.CreatedAt
	}

	return stats, nil
}
