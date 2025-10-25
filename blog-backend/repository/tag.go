package repository

import (
	"blog-backend/db"
	"blog-backend/model"
)

type TagRepository struct{}

func NewTagRepository() *TagRepository {
	return &TagRepository{}
}

// Create 创建标签
func (r *TagRepository) Create(tag *model.Tag) error {
	return db.DB.Create(tag).Error
}

// GetByID 根据ID获取标签
func (r *TagRepository) GetByID(id uint) (*model.Tag, error) {
	var tag model.Tag
	err := db.DB.First(&tag, id).Error
	return &tag, err
}

// GetByName 根据名称获取标签
func (r *TagRepository) GetByName(name string) (*model.Tag, error) {
	var tag model.Tag
	err := db.DB.Where("name = ?", name).First(&tag).Error
	return &tag, err
}

// Update 更新标签
func (r *TagRepository) Update(tag *model.Tag) error {
	return db.DB.Save(tag).Error
}

// Delete 删除标签
func (r *TagRepository) Delete(id uint) error {
	return db.DB.Delete(&model.Tag{}, id).Error
}

// List 获取标签列表
func (r *TagRepository) List() ([]model.Tag, error) {
	var tags []model.Tag
	err := db.DB.Order("post_count DESC, created_at DESC").Find(&tags).Error
	return tags, err
}

// GetOrCreate 获取或创建标签
func (r *TagRepository) GetOrCreate(name string) (*model.Tag, error) {
	var tag model.Tag
	err := db.DB.Where("name = ?", name).First(&tag).Error
	if err == nil {
		return &tag, nil
	}

	tag = model.Tag{Name: name}
	if err := db.DB.Create(&tag).Error; err != nil {
		return nil, err
	}

	return &tag, nil
}

// IncrementPostCount 增加文章数量
func (r *TagRepository) IncrementPostCount(id uint) error {
	return db.DB.Model(&model.Tag{}).Where("id = ?", id).UpdateColumn("post_count", db.DB.Raw("post_count + 1")).Error
}

// DecrementPostCount 减少文章数量
func (r *TagRepository) DecrementPostCount(id uint) error {
	return db.DB.Model(&model.Tag{}).Where("id = ?", id).UpdateColumn("post_count", db.DB.Raw("post_count - 1")).Error
}

