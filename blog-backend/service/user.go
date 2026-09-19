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
	"blog-backend/db"
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

	// 目标：删除前把该用户的全部文章转移给超级管理员，内容不随账号消失
	// （找不到可用的超级管理员时跳过转移，文章随外键级联删除）
	superAdmin, err := s.repo.GetSuperAdmin()
	if err != nil {
		superAdmin = nil
	}

	// 事务内清理关联数据后删除用户，保证删除行为一致：
	// 1. 文章归属先转移给超级管理员（保留已发布文章与草稿）
	// 2. post_revisions 的 editor_id/reviewer_id 外键无级联规则，需显式删除其提交的修订记录
	// 3. moments/moment_likes/password_reset_tokens 无外键保护，随用户删除避免孤儿数据
	// 4. page_views/post_views 为统计数据，置空用户标识保留记录
	// 5. 其余关联表（comments/likes/operation_logs 等）由外键级联清理
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if superAdmin != nil && superAdmin.ID != id {
			if err := tx.Model(&model.Post{}).Where("user_id = ?", id).Update("user_id", superAdmin.ID).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("editor_id = ?", id).Delete(&model.PostRevision{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&model.Moment{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&model.MomentLike{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&model.PasswordResetToken{}).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.PageView{}).Where("user_id = ?", id).Update("user_id", nil).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.PostView{}).Where("user_id = ?", id).Update("user_id", nil).Error; err != nil {
			return err
		}
		return tx.Delete(&model.User{}, id).Error
	})
}
