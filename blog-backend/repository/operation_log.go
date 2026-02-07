/*
 * 项目名称：blog-backend
 * 文件名称：operation_log.go
 * 创建时间：2026-02-06 22:00:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：操作日志数据访问层，提供操作日志的增删改查功能
 */
package repository

import (
	"blog-backend/db"
	"blog-backend/model"
)

// OperationLogRepository 操作日志数据访问层结构体
type OperationLogRepository struct{}

// NewOperationLogRepository 创建操作日志数据访问层实例
func NewOperationLogRepository() *OperationLogRepository {
	return &OperationLogRepository{}
}

// Create 创建操作日志
func (r *OperationLogRepository) Create(log *model.OperationLog) error {
	return db.DB.Create(log).Error
}

// List 获取操作日志列表（支持分页和筛选）
func (r *OperationLogRepository) List(page, pageSize int, module, action, username string) ([]model.OperationLog, int64, error) {
	var logs []model.OperationLog
	var total int64

	query := db.DB.Model(&model.OperationLog{})

	// 筛选条件
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Preload("User").
		Find(&logs).Error

	return logs, total, err
}

// GetByID 根据ID获取操作日志
func (r *OperationLogRepository) GetByID(id uint) (*model.OperationLog, error) {
	var log model.OperationLog
	err := db.DB.Preload("User").First(&log, id).Error
	return &log, err
}

// GetByUserID 根据用户ID获取操作日志列表
func (r *OperationLogRepository) GetByUserID(userID uint, limit int) ([]model.OperationLog, error) {
	var logs []model.OperationLog
	query := db.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs)
	return logs, query.Error
}

// GetByModule 根据模块获取操作日志列表
func (r *OperationLogRepository) GetByModule(module string, limit int) ([]model.OperationLog, error) {
	var logs []model.OperationLog
	query := db.DB.Where("module = ?", module).
		Order("created_at DESC").
		Limit(limit).
		Preload("User").
		Find(&logs)
	return logs, query.Error
}

// Delete 删除单个操作日志
func (r *OperationLogRepository) Delete(id uint) error {
	return db.DB.Delete(&model.OperationLog{}, id).Error
}

// DeleteBatch 批量删除操作日志
func (r *OperationLogRepository) DeleteBatch(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return db.DB.Where("id IN ?", ids).Delete(&model.OperationLog{}).Error
}

// DeleteOldLogs 删除指定天数之前的日志
func (r *OperationLogRepository) DeleteOldLogs(days int) error {
	return db.DB.Where("created_at < NOW() - INTERVAL '? days'", days).
		Delete(&model.OperationLog{}).Error
}
