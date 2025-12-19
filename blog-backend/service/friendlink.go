package service

import (
	"blog-backend/model"
	"blog-backend/repository"
)

type FriendLinkService struct {
	repo *repository.FriendLinkRepository
}

func NewFriendLinkService() *FriendLinkService {
	return &FriendLinkService{
		repo: repository.NewFriendLinkRepository(),
	}
}

// CreateFriendLinkRequest 创建友链请求
type CreateFriendLinkRequest struct {
	Name        string `json:"name" binding:"required"`
	URL         string `json:"url" binding:"required"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Screenshot  string `json:"screenshot"`
	AtomURL     string `json:"atom_url"`
	SortOrder   int    `json:"sort_order"`
	Status      int    `json:"status"`
}

// UpdateFriendLinkRequest 更新友链请求
type UpdateFriendLinkRequest struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Screenshot  string `json:"screenshot"`
	AtomURL     string `json:"atom_url"`
	SortOrder   int    `json:"sort_order"`
	Status      int    `json:"status"`
}

// Create 创建友链
func (s *FriendLinkService) Create(req *CreateFriendLinkRequest) (*model.FriendLink, error) {
	friendLink := &model.FriendLink{
		Name:        req.Name,
		URL:         req.URL,
		Icon:        req.Icon,
		Description: req.Description,
		Screenshot:  req.Screenshot,
		AtomURL:     req.AtomURL,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
	}

	if friendLink.Status == 0 {
		friendLink.Status = 1 // 默认启用
	}

	if err := s.repo.Create(friendLink); err != nil {
		return nil, err
	}

	return friendLink, nil
}

// GetByID 根据ID获取友链
func (s *FriendLinkService) GetByID(id uint) (*model.FriendLink, error) {
	return s.repo.GetByID(id)
}

// List 获取友链列表（管理员用）
func (s *FriendLinkService) List(page, pageSize int) ([]model.FriendLink, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.repo.List(page, pageSize)
}

// ListPublic 获取公开的友链列表（前端用）
func (s *FriendLinkService) ListPublic() ([]model.FriendLink, error) {
	return s.repo.ListPublic()
}

// Update 更新友链
func (s *FriendLinkService) Update(id uint, req *UpdateFriendLinkRequest) (*model.FriendLink, error) {
	friendLink, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		friendLink.Name = req.Name
	}
	if req.URL != "" {
		friendLink.URL = req.URL
	}
	friendLink.Icon = req.Icon
	friendLink.Description = req.Description
	friendLink.Screenshot = req.Screenshot
	friendLink.AtomURL = req.AtomURL
	friendLink.SortOrder = req.SortOrder
	if req.Status >= 0 {
		friendLink.Status = req.Status
	}

	if err := s.repo.Update(friendLink); err != nil {
		return nil, err
	}

	return friendLink, nil
}

// Delete 删除友链
func (s *FriendLinkService) Delete(id uint) error {
	return s.repo.Delete(id)
}
