package repository

import (
	"time"

	"blog-backend/db"
	"blog-backend/model"
)

type EmailChangeRepository struct{}

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

