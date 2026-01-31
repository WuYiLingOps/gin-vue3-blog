/*
 * 项目名称：blog-backend
 * 文件名称：category.go
 * 创建时间：2026-01-31 16:29:06
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：分类数据访问层，提供文章分类的数据库操作功能
 */
package repository

import (
	"blog-backend/db"
	"blog-backend/model"

	"gorm.io/gorm"
)

// CategoryRepository 分类数据访问层结构体
type CategoryRepository struct{}

// NewCategoryRepository 创建分类数据访问层实例
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
