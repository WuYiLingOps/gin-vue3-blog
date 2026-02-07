/*
 * 项目名称：blog-backend
 * 文件名称：user.go
 * 创建时间：2026-01-31 16:34:35
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：用户业务逻辑层，提供用户信息查询、状态更新、删除等业务处理
 */
package service

import (
	"errors"

	"blog-backend/constant"
	"blog-backend/model"
	"blog-backend/repository"

	"gorm.io/gorm"
)

// UserService 用户业务逻辑层结构体
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService 创建用户业务逻辑层实例
func NewUserService() *UserService {
	return &UserService{
		repo: repository.NewUserRepository(),
	}
}

// GetByID 获取用户详情
func (s *UserService) GetByID(id uint) (*model.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, errors.New("获取用户失败")
	}
	return user, nil
}

// List 获取用户列表
func (s *UserService) List(page, pageSize int) ([]model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.repo.List(page, pageSize)
}

// UpdateStatus 更新用户状态
func (s *UserService) UpdateStatus(id uint, status int) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return errors.New("用户不存在")
	}

	return s.repo.UpdateStatus(id, status)
}

// UpdateRole 更新用户角色
func (s *UserService) UpdateRole(id uint, role string) error {
	// 验证角色值
	if role != constant.RoleSuperAdmin && role != constant.RoleAdmin && role != constant.RoleUser {
		return errors.New("无效的角色值")
	}

	// 检查用户是否存在
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return errors.New("获取用户失败")
	}

	// 禁止将 super_admin 降级为其他角色
	if user.Role == constant.RoleSuperAdmin && role != constant.RoleSuperAdmin {
		return errors.New("禁止修改超级管理员的角色")
	}

	// 禁止将普通用户直接升级为 super_admin（只能通过数据库手动设置）
	if role == constant.RoleSuperAdmin && user.Role != constant.RoleSuperAdmin {
		return errors.New("禁止将用户升级为超级管理员，请通过数据库手动设置")
	}

	return s.repo.UpdateRole(id, role)
}

// Delete 删除用户
func (s *UserService) Delete(id uint) error {
	// 检查用户是否存在
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return errors.New("获取用户失败")
	}

	// 禁止删除 super_admin 账号
	if user.Role == constant.RoleSuperAdmin {
		return errors.New("禁止删除超级管理员账号")
	}

	return s.repo.Delete(id)
}
