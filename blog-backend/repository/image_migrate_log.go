/*
 * 项目名称：blog-backend
 * 文件名称：image_migrate_log.go
 * 创建时间：2026-11-20 10:00:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：图片迁移日志数据访问层，提供迁移日志的持久化与查询功能
 */
package repository

import (
	"blog-backend/model"

	"blog-backend/db"
)

// ImageMigrateLogRepository 图片迁移日志数据访问层结构体
type ImageMigrateLogRepository struct{}

// NewImageMigrateLogRepository 创建图片迁移日志数据访问层实例
func NewImageMigrateLogRepository() *ImageMigrateLogRepository {
	return &ImageMigrateLogRepository{}
}

// Create 写入一条迁移日志
func (r *ImageMigrateLogRepository) Create(log *model.ImageMigrateLog) error {
	return db.DB.Create(log).Error
}

// List 分页查询迁移日志（按完成时间倒序）
func (r *ImageMigrateLogRepository) List(page, pageSize int) ([]model.ImageMigrateLog, int64, error) {
	var logs []model.ImageMigrateLog
	var total int64

	offset := (page - 1) * pageSize

	if err := db.DB.Model(&model.ImageMigrateLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.DB.Order("created_at DESC, id DESC").
		Offset(offset).Limit(pageSize).
		Find(&logs).Error

	return logs, total, err
}

// Delete 删除单条迁移日志
func (r *ImageMigrateLogRepository) Delete(id uint) error {
	return db.DB.Delete(&model.ImageMigrateLog{}, id).Error
}

// DeleteBatch 批量删除迁移日志
func (r *ImageMigrateLogRepository) DeleteBatch(ids []uint) error {
	return db.DB.Where("id IN ?", ids).Delete(&model.ImageMigrateLog{}).Error
}

// DeleteAll 清空迁移日志
func (r *ImageMigrateLogRepository) DeleteAll() error {
	return db.DB.Where("1 = 1").Delete(&model.ImageMigrateLog{}).Error
}
