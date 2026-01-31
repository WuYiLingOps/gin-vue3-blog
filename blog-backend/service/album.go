/*
 * 项目名称：blog-backend
 * 文件名称：album.go
 * 创建时间：2026-01-31 16:34:35
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：相册业务逻辑层，提供相册照片的增删改查业务处理
 */
package service

import (
	"blog-backend/model"
	"blog-backend/repository"
)

// AlbumService 相册业务逻辑层结构体
type AlbumService struct {
	repo *repository.AlbumRepository
}

// NewAlbumService 创建相册业务逻辑层实例
func NewAlbumService() *AlbumService {
	return &AlbumService{
		repo: repository.NewAlbumRepository(),
	}
}

// CreateAlbumRequest 创建相册照片请求
type CreateAlbumRequest struct {
	ImageURL    string `json:"image_url" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// UpdateAlbumRequest 更新相册照片请求
type UpdateAlbumRequest struct {
	ImageURL    *string `json:"image_url"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	SortOrder   *int    `json:"sort_order"`
}

// Create 创建相册照片
func (s *AlbumService) Create(req *CreateAlbumRequest) (*model.Album, error) {
	album := &model.Album{
		ImageURL:    req.ImageURL,
		Title:       req.Title,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}

	if err := s.repo.Create(album); err != nil {
		return nil, err
	}

	return s.repo.GetByID(album.ID)
}

// GetByID 根据ID获取相册照片
func (s *AlbumService) GetByID(id uint) (*model.Album, error) {
	return s.repo.GetByID(id)
}

// List 获取相册列表（管理员用）
func (s *AlbumService) List(page, pageSize int) ([]model.Album, int64, error) {
	return s.repo.List(page, pageSize)
}

// ListPublic 获取公开的相册列表（前端用）
func (s *AlbumService) ListPublic() ([]model.Album, error) {
	return s.repo.ListPublic()
}

// Update 更新相册照片
func (s *AlbumService) Update(id uint, req *UpdateAlbumRequest) (*model.Album, error) {
	album, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.ImageURL != nil {
		album.ImageURL = *req.ImageURL
	}
	if req.Title != nil {
		album.Title = *req.Title
	}
	if req.Description != nil {
		album.Description = *req.Description
	}
	if req.SortOrder != nil {
		album.SortOrder = *req.SortOrder
	}

	if err := s.repo.Update(album); err != nil {
		return nil, err
	}

	return s.repo.GetByID(id)
}

// Delete 删除相册照片
func (s *AlbumService) Delete(id uint) error {
	return s.repo.Delete(id)
}
