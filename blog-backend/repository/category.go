package repository

import (
	"blog-backend/db"
	"blog-backend/model"
	"gorm.io/gorm"
)

type CategoryRepository struct{}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

// Create 创建分类
func (r *CategoryRepository) Create(category *model.Category) error {
	return db.DB.Create(category).Error
}

// GetByID 根据ID获取分类
func (r *CategoryRepository) GetByID(id uint) (*model.Category, error) {
	var category model.Category
	err := db.DB.First(&category, id).Error
	return &category, err
}

// GetByName 根据名称获取分类
func (r *CategoryRepository) GetByName(name string) (*model.Category, error) {
	var category model.Category
	err := db.DB.Where("name = ?", name).First(&category).Error
	return &category, err
}

// Update 更新分类
func (r *CategoryRepository) Update(category *model.Category) error {
	return db.DB.Save(category).Error
}

// Delete 删除分类
func (r *CategoryRepository) Delete(id uint) error {
	return db.DB.Delete(&model.Category{}, id).Error
}

// List 获取分类列表
func (r *CategoryRepository) List() ([]model.Category, error) {
	var categories []model.Category
	err := db.DB.Order("sort ASC, created_at DESC").Find(&categories).Error
	return categories, err
}

// IncrementPostCount 增加文章数量
func (r *CategoryRepository) IncrementPostCount(id uint) error {
	return db.DB.Model(&model.Category{}).Where("id = ?", id).UpdateColumn("post_count", db.DB.Raw("post_count + 1")).Error
}

// DecrementPostCount 减少文章数量
func (r *CategoryRepository) DecrementPostCount(id uint) error {
	return db.DB.Model(&model.Category{}).Where("id = ?", id).UpdateColumn("post_count", db.DB.Raw("post_count - 1")).Error
}

// IncrementPostCountTx 在事务中增加文章数量
func (r *CategoryRepository) IncrementPostCountTx(tx *gorm.DB, id uint) error {
	return tx.Model(&model.Category{}).Where("id = ?", id).UpdateColumn("post_count", gorm.Expr("post_count + 1")).Error
}

// DecrementPostCountTx 在事务中减少文章数量
func (r *CategoryRepository) DecrementPostCountTx(tx *gorm.DB, id uint) error {
	return tx.Model(&model.Category{}).Where("id = ?", id).UpdateColumn("post_count", gorm.Expr("CASE WHEN post_count > 0 THEN post_count - 1 ELSE 0 END")).Error
}

