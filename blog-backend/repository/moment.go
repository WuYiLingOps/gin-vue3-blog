package repository

import (
	"blog-backend/db"
	"blog-backend/model"
)

type MomentRepository struct{}

func NewMomentRepository() *MomentRepository {
	return &MomentRepository{}
}

// Create 创建说说
func (r *MomentRepository) Create(moment *model.Moment) error {
	err := db.DB.Create(moment).Error
	if err != nil {
		return err
	}
	
	// 更新全文搜索向量
	db.DB.Exec(
		"UPDATE moments SET content_tsv = to_tsvector('english', coalesce(content, '')) WHERE id = ?",
		moment.ID,
	)
	
	return nil
}

// Update 更新说说
func (r *MomentRepository) Update(moment *model.Moment) error {
	err := db.DB.Save(moment).Error
	if err != nil {
		return err
	}
	
	// 更新全文搜索向量
	db.DB.Exec(
		"UPDATE moments SET content_tsv = to_tsvector('english', coalesce(content, '')) WHERE id = ?",
		moment.ID,
	)
	
	return nil
}

// Delete 删除说说（软删除）
func (r *MomentRepository) Delete(id uint) error {
	return db.DB.Model(&model.Moment{}).Where("id = ?", id).Update("status", -1).Error
}

// GetByID 根据ID获取说说
func (r *MomentRepository) GetByID(id uint) (*model.Moment, error) {
	var moment model.Moment
	err := db.DB.Preload("User").Where("id = ?", id).First(&moment).Error
	if err != nil {
		return nil, err
	}
	return &moment, nil
}

// List 获取说说列表
func (r *MomentRepository) List(page, pageSize int, status *int, keyword string) ([]model.Moment, int64, error) {
	var moments []model.Moment
	var total int64

	query := db.DB.Model(&model.Moment{})
	
	// 如果指定了状态，则过滤
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 关键词搜索
	if keyword != "" {
		// 使用PostgreSQL全文搜索（优先）+ ILIKE后备
		// 如果content_tsv字段不存在或为NULL，查询会忽略该条件，只使用ILIKE
		query = query.Where(
			"(content_tsv IS NOT NULL AND content_tsv @@ plainto_tsquery('english', ?)) OR content ILIKE ?",
			keyword, "%"+keyword+"%",
		)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Preload("User").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&moments).Error

	if err != nil {
		return nil, 0, err
	}

	return moments, total, nil
}

// GetRecent 获取最新的说说
func (r *MomentRepository) GetRecent(limit int) ([]model.Moment, error) {
	var moments []model.Moment
	err := db.DB.Preload("User").
		Where("status = ?", 1).
		Order("created_at DESC").
		Limit(limit).
		Find(&moments).Error
	
	return moments, err
}

