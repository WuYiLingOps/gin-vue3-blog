/*
 * 项目名称：blog-backend
 * 文件名称：operation_log.go
 * 创建时间：2026-02-06 22:00:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：操作日志业务逻辑层，提供操作日志的查询、统计等业务处理
 */
package service

import (
	"errors"

	"blog-backend/model"
	"blog-backend/repository"
)

// OperationLogService 操作日志业务逻辑层结构体
type OperationLogService struct {
	repo *repository.OperationLogRepository
}

// NewOperationLogService 创建操作日志业务逻辑层实例
func NewOperationLogService() *OperationLogService {
	return &OperationLogService{
		repo: repository.NewOperationLogRepository(),
	}
}

// List 获取操作日志列表（支持分页和筛选）
func (s *OperationLogService) List(page, pageSize int, module, action, username string) ([]model.OperationLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return s.repo.List(page, pageSize, module, action, username)
}

// GetByID 根据ID获取操作日志详情
func (s *OperationLogService) GetByID(id uint) (*model.OperationLog, error) {
	return s.repo.GetByID(id)
}

// GetByUserID 根据用户ID获取操作日志列表
func (s *OperationLogService) GetByUserID(userID uint, limit int) ([]model.OperationLog, error) {
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.GetByUserID(userID, limit)
}

// GetByModule 根据模块获取操作日志列表
func (s *OperationLogService) GetByModule(module string, limit int) ([]model.OperationLog, error) {
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.GetByModule(module, limit)
}

// Delete 删除单个操作日志
func (s *OperationLogService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// DeleteBatch 批量删除操作日志
func (s *OperationLogService) DeleteBatch(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	// 限制批量删除数量，防止一次性删除过多
	if len(ids) > 100 {
		return errors.New("批量删除数量不能超过100条")
	}
	return s.repo.DeleteBatch(ids)
}
