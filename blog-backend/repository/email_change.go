/*
 * 项目名称：blog-backend
 * 文件名称：email_change.go
 * 创建时间：2026-01-31 16:29:06
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：邮箱修改记录数据访问层，提供邮箱修改记录的数据库操作功能
 */
package repository

import (
	"time"

	"blog-backend/db"
	"blog-backend/model"
)

// EmailChangeRepository 邮箱修改记录数据访问层结构体
type EmailChangeRepository struct{}

// NewEmailChangeRepository 创建邮箱修改记录数据访问层实例
func NewEmailChangeRepository() *EmailChangeRepository {
	return &EmailChangeRepository{}
}

// Create 创建邮箱修改记录
func (r *EmailChangeRepository) Create(record *model.EmailChangeRecord) error {
	return db.DB.Create(record).Error
}

// CountByUserIDInYear 统计用户一年内的邮箱修改次数
func (r *EmailChangeRepository) CountByUserIDInYear(userID uint) (int64, error) {
	var count int64
	oneYearAgo := time.Now().AddDate(-1, 0, 0) // 一年前
	err := db.DB.Model(&model.EmailChangeRecord{}).
		Where("user_id = ? AND changed_at >= ?", userID, oneYearAgo).
		Count(&count).Error
	return count, err
}

// GetRecordsByUserID 获取用户的邮箱修改记录
func (r *EmailChangeRepository) GetRecordsByUserID(userID uint, limit int) ([]model.EmailChangeRecord, error) {
	var records []model.EmailChangeRecord
	err := db.DB.Where("user_id = ?", userID).
		Order("changed_at DESC").
		Limit(limit).
		Find(&records).Error
	return records, err
}
