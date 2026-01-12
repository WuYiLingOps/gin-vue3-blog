package repository

import (
	"blog-backend/db"
	"blog-backend/model"
)

type AlbumRepository struct{}

func NewAlbumRepository() *AlbumRepository {
	return &AlbumRepository{}
}

// Create 创建相册照片
func (r *AlbumRepository) Create(album *model.Album) error {
	return db.DB.Create(album).Error
}

// GetByID 根据ID获取相册照片
func (r *AlbumRepository) GetByID(id uint) (*model.Album, error) {
	var album model.Album
	err := db.DB.First(&album, id).Error
	return &album, err
}

// List 获取相册列表（管理员用）
func (r *AlbumRepository) List(page, pageSize int) ([]model.Album, int64, error) {
	var albums []model.Album
	var total int64

	offset := (page - 1) * pageSize

	if err := db.DB.Model(&model.Album{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.DB.Order("sort_order DESC, id DESC").
		Offset(offset).Limit(pageSize).
		Find(&albums).Error

	return albums, total, err
}

// ListPublic 获取公开的相册列表（前端用）
func (r *AlbumRepository) ListPublic() ([]model.Album, error) {
	var albums []model.Album
	err := db.DB.Order("sort_order DESC, id DESC").
		Find(&albums).Error
	return albums, err
}

// Update 更新相册照片
func (r *AlbumRepository) Update(album *model.Album) error {
	return db.DB.Save(album).Error
}

// Delete 删除相册照片
func (r *AlbumRepository) Delete(id uint) error {
	return db.DB.Delete(&model.Album{}, id).Error
}
