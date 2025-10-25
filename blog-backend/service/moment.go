package service

import (
	"blog-backend/model"
	"blog-backend/repository"
	"errors"
)

type MomentService struct {
	repo *repository.MomentRepository
}

func NewMomentService() *MomentService {
	return &MomentService{
		repo: repository.NewMomentRepository(),
	}
}

// Create 创建说说
func (s *MomentService) Create(moment *model.Moment) error {
	if moment.Content == "" {
		return errors.New("说说内容不能为空")
	}
	return s.repo.Create(moment)
}

// Update 更新说说
func (s *MomentService) Update(id uint, content, images string) error {
	moment, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if content != "" {
		moment.Content = content
	}
	moment.Images = images

	return s.repo.Update(moment)
}

// Delete 删除说说
func (s *MomentService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// GetByID 获取说说详情
func (s *MomentService) GetByID(id uint) (*model.Moment, error) {
	return s.repo.GetByID(id)
}

// List 获取说说列表
func (s *MomentService) List(page, pageSize int, status *int, keyword string) ([]model.Moment, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return s.repo.List(page, pageSize, status, keyword)
}

// GetRecent 获取最新说说
func (s *MomentService) GetRecent(limit int) ([]model.Moment, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	return s.repo.GetRecent(limit)
}

