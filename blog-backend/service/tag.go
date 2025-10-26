package service

import (
	"errors"

	"blog-backend/model"
	"blog-backend/repository"
	"gorm.io/gorm"
)

type TagService struct {
	repo *repository.TagRepository
}

func NewTagService() *TagService {
	return &TagService{
		repo: repository.NewTagRepository(),
	}
}

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name      string  `json:"name" binding:"required"`
	Color     string  `json:"color"`
	TextColor *string `json:"text_color"`
	FontSize  *int    `json:"font_size"`
}

// UpdateTagRequest 更新标签请求
type UpdateTagRequest struct {
	Name      string  `json:"name"`
	Color     string  `json:"color"`
	TextColor *string `json:"text_color"`
	FontSize  *int    `json:"font_size"`
}

// Create 创建标签
func (s *TagService) Create(req *CreateTagRequest) (*model.Tag, error) {
	// 检查标签名是否已存在
	if _, err := s.repo.GetByName(req.Name); err == nil {
		return nil, errors.New("标签名已存在")
	}

	tag := &model.Tag{
		Name:      req.Name,
		Color:     req.Color,
		TextColor: req.TextColor,
		FontSize:  req.FontSize,
	}

	if err := s.repo.Create(tag); err != nil {
		return nil, errors.New("标签创建失败")
	}

	return tag, nil
}

// GetByID 获取标签详情
func (s *TagService) GetByID(id uint) (*model.Tag, error) {
	tag, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("标签不存在")
		}
		return nil, errors.New("获取标签失败")
	}
	return tag, nil
}

// Update 更新标签
func (s *TagService) Update(id uint, req *UpdateTagRequest) (*model.Tag, error) {
	tag, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("标签不存在")
	}

	// 如果修改了名称，检查新名称是否已存在
	if req.Name != "" && req.Name != tag.Name {
		if _, err := s.repo.GetByName(req.Name); err == nil {
			return nil, errors.New("标签名已存在")
		}
		tag.Name = req.Name
	}

	if req.Color != "" {
		tag.Color = req.Color
	}

	// 更新文字颜色和字体大小（即使是 nil 也要更新，以支持清空操作）
	tag.TextColor = req.TextColor
	tag.FontSize = req.FontSize

	if err := s.repo.Update(tag); err != nil {
		return nil, errors.New("标签更新失败")
	}

	return tag, nil
}

// Delete 删除标签
func (s *TagService) Delete(id uint) error {
	tag, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("标签不存在")
	}

	// 检查标签下是否有文章
	if tag.PostCount > 0 {
		return errors.New("该标签下还有文章，无法删除")
	}

	return s.repo.Delete(id)
}

// List 获取标签列表
func (s *TagService) List() ([]model.Tag, error) {
	return s.repo.List()
}

