package repository

import (
	"blog-backend/db"
	"blog-backend/model"

	"gorm.io/gorm"
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
	if err != nil {
		return nil, err
	}

	// 实时计算每个标签的已发布文章数，确保数据准确性
	for i := range tags {
		var count int64
		err := db.DB.Model(&model.Post{}).
			Joins("JOIN post_tags ON posts.id = post_tags.post_id").
			Where("post_tags.tag_id = ? AND posts.status = 1", tags[i].ID).
			Count(&count).Error
		if err == nil {
			// 如果计算出的数量与数据库中的不一致，更新返回的数据和数据库
			if int64(tags[i].PostCount) != count {
				tags[i].PostCount = int(count)
				// 同步更新数据库中的 post_count，确保数据一致性
				db.DB.Model(&model.Tag{}).Where("id = ?", tags[i].ID).Update("post_count", int(count))
			}
		}
	}

	return tags, nil
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

// IncrementPostCountTx 在事务中增加文章数量
func (r *TagRepository) IncrementPostCountTx(tx *gorm.DB, id uint) error {
	return tx.Model(&model.Tag{}).Where("id = ?", id).UpdateColumn("post_count", gorm.Expr("post_count + 1")).Error
}

// DecrementPostCountTx 在事务中减少文章数量
func (r *TagRepository) DecrementPostCountTx(tx *gorm.DB, id uint) error {
	return tx.Model(&model.Tag{}).Where("id = ?", id).UpdateColumn("post_count", gorm.Expr("CASE WHEN post_count > 0 THEN post_count - 1 ELSE 0 END")).Error
}
