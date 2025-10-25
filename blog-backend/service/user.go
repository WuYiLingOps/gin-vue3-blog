package service

import (
	"errors"

	"blog-backend/model"
	"blog-backend/repository"

	"gorm.io/gorm"
)

type UserService struct {
	repo *repository.UserRepository
}

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

// Delete 删除用户
func (s *UserService) Delete(id uint) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return errors.New("用户不存在")
	}

	return s.repo.Delete(id)
}
