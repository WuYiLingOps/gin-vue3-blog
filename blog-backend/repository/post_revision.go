/*
 * 项目名称：blog-backend
 * 文件名称：post_revision.go
 * 创建时间：2026-04-28 17:48:56
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：文章修订版本数据访问层，负责修订版本的CRUD操作
 */
package repository

import (
	"blog-backend/db"
	"blog-backend/model"
	"context"

	"gorm.io/gorm"
)

// PostRevisionRepository 文章修订版本仓库结构体
type PostRevisionRepository struct{}

// NewPostRevisionRepository 创建文章修订版本仓库实例
func NewPostRevisionRepository() *PostRevisionRepository {
	return &PostRevisionRepository{}
}

// Create 创建修订版本
func (r *PostRevisionRepository) Create(ctx context.Context, revision *model.PostRevision) error {
	return db.DB.WithContext(ctx).Create(revision).Error
}

// GetByID 根据ID获取修订版本
func (r *PostRevisionRepository) GetByID(ctx context.Context, id uint) (*model.PostRevision, error) {
	var revision model.PostRevision
	err := db.DB.WithContext(ctx).
		Preload("Post").
		Preload("Editor").
		Preload("Reviewer").
		First(&revision, id).Error
	return &revision, err
}

// GetPendingList 获取待审批列表（分页）
func (r *PostRevisionRepository) GetPendingList(ctx context.Context, page, pageSize int) ([]*model.PostRevision, int64, error) {
	var revisions []*model.PostRevision
	var total int64

	offset := (page - 1) * pageSize

	// 统计总数
	if err := db.DB.WithContext(ctx).Model(&model.PostRevision{}).Where("status = ?", "pending").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询列表
	err := db.DB.WithContext(ctx).
		Preload("Post").
		Preload("Editor").
		Where("status = ?", "pending").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&revisions).Error

	return revisions, total, err
}

// GetByPostID 获取文章的修订历史
func (r *PostRevisionRepository) GetByPostID(ctx context.Context, postID uint) ([]*model.PostRevision, error) {
	var revisions []*model.PostRevision
	err := db.DB.WithContext(ctx).
		Preload("Editor").
		Preload("Reviewer").
		Where("post_id = ?", postID).
		Order("created_at DESC").
		Find(&revisions).Error
	return revisions, err
}

// UpdateStatus 更新修订状态
func (r *PostRevisionRepository) UpdateStatus(ctx context.Context, revision *model.PostRevision) error {
	return db.DB.WithContext(ctx).Save(revision).Error
}

// UpdateStatusTx 更新修订状态(事务版本)
func (r *PostRevisionRepository) UpdateStatusTx(tx *gorm.DB, revision *model.PostRevision) error {
	return tx.Save(revision).Error
}

// GetPendingCount 获取待审批数量
func (r *PostRevisionRepository) GetPendingCount(ctx context.Context) (int64, error) {
	var count int64
	err := db.DB.WithContext(ctx).Model(&model.PostRevision{}).Where("status = ?", "pending").Count(&count).Error
	return count, err
}

// GetByEditorID 获取指定用户的修订记录（分页）
func (r *PostRevisionRepository) GetByEditorID(ctx context.Context, editorID uint, page, pageSize int, status string) ([]*model.PostRevision, int64, error) {
	var revisions []*model.PostRevision
	var total int64

	offset := (page - 1) * pageSize

	query := db.DB.WithContext(ctx).Model(&model.PostRevision{}).Where("editor_id = ?", editorID)

	// 如果指定了状态，添加状态过滤
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询列表
	err := query.
		Preload("Post").
		Preload("Reviewer").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&revisions).Error

	return revisions, total, err
}
