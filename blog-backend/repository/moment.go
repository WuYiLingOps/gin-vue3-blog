/*
 * 项目名称：blog-backend
 * 文件名称：moment.go
 * 创建时间：2026-01-31 16:29:06
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：说说数据访问层，提供说说的数据库操作功能，支持全文搜索和点赞管理
 */
package repository

import (
	"blog-backend/db"
	"blog-backend/model"

	"gorm.io/gorm"
)

// MomentRepository 说说数据访问层结构体
type MomentRepository struct{}

// NewMomentRepository 创建说说数据访问层实例
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
	} else {
		// 默认排除已删除的说说（status = -1）
		query = query.Where("status != ?", -1)
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

// CreateLike 创建点赞记录
func (r *MomentRepository) CreateLike(like *model.MomentLike) error {
	return db.DB.Create(like).Error
}

// DeleteLike 删除点赞记录
func (r *MomentRepository) DeleteLike(momentID uint, userID *uint, ip string) error {
	query := db.DB.Where("moment_id = ?", momentID)

	if userID != nil && *userID > 0 {
		query = query.Where("user_id = ?", *userID)
	} else {
		query = query.Where("ip = ? AND user_id IS NULL", ip)
	}

	return query.Delete(&model.MomentLike{}).Error
}

// CheckLiked 检查是否已点赞（通过用户ID或IP）
func (r *MomentRepository) CheckLiked(momentID uint, userID *uint, ip string) (bool, error) {
	var count int64
	query := db.DB.Model(&model.MomentLike{}).Where("moment_id = ?", momentID)

	if userID != nil && *userID > 0 {
		// 已登录用户，通过用户ID查询
		query = query.Where("user_id = ?", *userID)
	} else {
		// 匿名用户，通过IP查询
		query = query.Where("ip = ? AND user_id IS NULL", ip)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// GetLikedMomentIDs 获取用户点赞过的说说ID列表
func (r *MomentRepository) GetLikedMomentIDs(momentIDs []uint, userID *uint, ip string) ([]uint, error) {
	var likes []model.MomentLike
	query := db.DB.Where("moment_id IN ?", momentIDs)

	if userID != nil && *userID > 0 {
		query = query.Where("user_id = ?", *userID)
	} else {
		query = query.Where("ip = ? AND user_id IS NULL", ip)
	}

	err := query.Find(&likes).Error
	if err != nil {
		return nil, err
	}

	likedIDs := make([]uint, 0, len(likes))
	for _, like := range likes {
		likedIDs = append(likedIDs, like.MomentID)
	}

	return likedIDs, nil
}

// Transaction 执行事务
func (r *MomentRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return db.DB.Transaction(fn)
}

// CreateLikeTx 在事务中创建点赞记录
func (r *MomentRepository) CreateLikeTx(tx *gorm.DB, like *model.MomentLike) error {
	return tx.Create(like).Error
}

// DeleteLikeTx 在事务中删除点赞记录
func (r *MomentRepository) DeleteLikeTx(tx *gorm.DB, momentID uint, userID *uint, ip string) error {
	query := tx.Where("moment_id = ?", momentID)

	if userID != nil && *userID > 0 {
		query = query.Where("user_id = ?", *userID)
	} else {
		query = query.Where("ip = ? AND user_id IS NULL", ip)
	}

	return query.Delete(&model.MomentLike{}).Error
}

// UpdateTx 在事务中更新说说
func (r *MomentRepository) UpdateTx(tx *gorm.DB, moment *model.Moment) error {
	err := tx.Save(moment).Error
	if err != nil {
		return err
	}

	// 更新全文搜索向量
	tx.Exec(
		"UPDATE moments SET content_tsv = to_tsvector('english', coalesce(content, '')) WHERE id = ?",
		moment.ID,
	)

	return nil
}
