package repository

import (
	"blog-backend/db"
	"blog-backend/model"

	"gorm.io/gorm"
)

type PostRepository struct{}

func NewPostRepository() *PostRepository {
	return &PostRepository{}
}

// Create 创建文章
func (r *PostRepository) Create(post *model.Post) error {
	// 创建文章时，同时更新search_tsv字段（如果列存在）
	err := db.DB.Create(post).Error
	if err != nil {
		return err
	}

	// 更新全文搜索向量
	db.DB.Exec(
		"UPDATE posts SET search_tsv = setweight(to_tsvector('english', coalesce(title, '')), 'A') || setweight(to_tsvector('english', coalesce(content, '')), 'B') WHERE id = ?",
		post.ID,
	)

	return nil
}

// GetByID 根据ID获取文章
func (r *PostRepository) GetByID(id uint) (*model.Post, error) {
	var post model.Post
	err := db.DB.Preload("User").Preload("Category").Preload("Tags").First(&post, id).Error
	return &post, err
}

// Update 更新文章
func (r *PostRepository) Update(post *model.Post) error {
	err := db.DB.Save(post).Error
	if err != nil {
		return err
	}

	// 更新全文搜索向量
	db.DB.Exec(
		"UPDATE posts SET search_tsv = setweight(to_tsvector('english', coalesce(title, '')), 'A') || setweight(to_tsvector('english', coalesce(content, '')), 'B') WHERE id = ?",
		post.ID,
	)

	return nil
}

// Delete 删除文章
func (r *PostRepository) Delete(id uint) error {
	return db.DB.Delete(&model.Post{}, id).Error
}

// GetPublishedCountByCategory 获取指定分类的已发布文章数量
func (r *PostRepository) GetPublishedCountByCategory(categoryID uint) (int64, error) {
	var count int64
	err := db.DB.Model(&model.Post{}).Where("status = 1 AND category_id = ?", categoryID).Count(&count).Error
	return count, err
}

// List 获取文章列表
func (r *PostRepository) List(page, pageSize int, categoryID uint, keyword string, status int, visibility *int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	offset := (page - 1) * pageSize
	query := db.DB.Model(&model.Post{})

	// 筛选条件
	if status >= 0 {
		query = query.Where("status = ?", status)
	}
	if visibility != nil {
		query = query.Where("visibility = ?", *visibility)
	}
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		// 使用PostgreSQL全文搜索（优先）+ ILIKE后备
		// 如果search_tsv字段不存在或为NULL，查询会忽略该条件，只使用ILIKE
		query = query.Where(
			"(search_tsv IS NOT NULL AND search_tsv @@ plainto_tsquery('english', ?)) OR title ILIKE ? OR content ILIKE ?",
			keyword, "%"+keyword+"%", "%"+keyword+"%",
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("User").Preload("Category").Preload("Tags").
		Order("is_top DESC, created_at DESC").
		Offset(offset).Limit(pageSize).Find(&posts).Error

	return posts, total, err
}

// GetByTag 根据标签获取文章列表
func (r *PostRepository) GetByTag(tagID uint, page, pageSize int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	offset := (page - 1) * pageSize

	// 通过多对多关系查询
	if err := db.DB.Model(&model.Post{}).
		Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Where("post_tags.tag_id = ? AND posts.status = 1", tagID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.DB.Preload("User").Preload("Category").Preload("Tags").
		Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Where("post_tags.tag_id = ? AND posts.status = 1", tagID).
		Order("posts.created_at DESC").
		Offset(offset).Limit(pageSize).Find(&posts).Error

	return posts, total, err
}

// IncrementViewCount 增加浏览量
func (r *PostRepository) IncrementViewCount(id uint) error {
	return db.DB.Model(&model.Post{}).Where("id = ?", id).UpdateColumn("view_count", db.DB.Raw("view_count + 1")).Error
}

// IncrementLikeCount 增加点赞数
func (r *PostRepository) IncrementLikeCount(id uint) error {
	return db.DB.Model(&model.Post{}).Where("id = ?", id).UpdateColumn("like_count", db.DB.Raw("like_count + 1")).Error
}

// DecrementLikeCount 减少点赞数
func (r *PostRepository) DecrementLikeCount(id uint) error {
	return db.DB.Model(&model.Post{}).Where("id = ?", id).UpdateColumn("like_count", db.DB.Raw("CASE WHEN like_count > 0 THEN like_count - 1 ELSE 0 END")).Error
}

// CreateLike 创建点赞记录
func (r *PostRepository) CreateLike(like *model.PostLike) error {
	return db.DB.Create(like).Error
}

// DeleteLike 删除点赞记录
func (r *PostRepository) DeleteLike(postID uint, userID *uint, ip string) error {
	query := db.DB.Where("post_id = ?", postID)

	if userID != nil && *userID > 0 {
		query = query.Where("user_id = ?", *userID)
	} else {
		query = query.Where("ip = ? AND user_id IS NULL", ip)
	}

	return query.Delete(&model.PostLike{}).Error
}

// CheckLiked 检查是否已点赞
func (r *PostRepository) CheckLiked(postID uint, userID *uint, ip string) (bool, error) {
	var count int64
	query := db.DB.Model(&model.PostLike{}).Where("post_id = ?", postID)

	if userID != nil && *userID > 0 {
		query = query.Where("user_id = ?", *userID)
	} else {
		query = query.Where("ip = ? AND user_id IS NULL", ip)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// GetArchives 获取归档列表
func (r *PostRepository) GetArchives() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := db.DB.Model(&model.Post{}).
		Select("DATE_TRUNC('month', created_at) as month, COUNT(*) as count").
		Where("status = 1").
		Group("month").
		Order("month DESC").
		Find(&results).Error
	return results, err
}

// GetHotPosts 获取热门文章
func (r *PostRepository) GetHotPosts(limit int) ([]model.Post, error) {
	var posts []model.Post
	err := db.DB.Preload("User").Preload("Category").
		Where("status = 1").
		Order("view_count DESC").
		Limit(limit).Find(&posts).Error
	return posts, err
}

// GetRecentPosts 获取最新文章
// - 普通用户/游客：仅公开
// - 管理员：公开 + 自己的私密
func (r *PostRepository) GetRecentPosts(limit int, userID *uint, role string) ([]model.Post, error) {
	var posts []model.Post

	query := db.DB.Preload("User").Preload("Category").Where("status = 1")

	if role == "admin" && userID != nil {
		// 管理员可见自己的私密文章，其余仍需公开
		query = query.Where("visibility = 1 OR user_id = ?", *userID)
	} else {
		// 普通用户/游客仅公开文章
		query = query.Where("visibility = 1")
	}

	err := query.Order("created_at DESC").Limit(limit).Find(&posts).Error
	return posts, err
}

// UpdateTags 更新文章标签
func (r *PostRepository) UpdateTags(postID uint, tagIDs []uint) error {
	var post model.Post
	if err := db.DB.First(&post, postID).Error; err != nil {
		return err
	}

	var tags []model.Tag
	if err := db.DB.Find(&tags, tagIDs).Error; err != nil {
		return err
	}

	return db.DB.Model(&post).Association("Tags").Replace(tags)
}

// Transaction 执行事务
func (r *PostRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return db.DB.Transaction(fn)
}

// DeleteTx 在事务中删除文章
func (r *PostRepository) DeleteTx(tx *gorm.DB, id uint) error {
	return tx.Delete(&model.Post{}, id).Error
}

// CreateTx 在事务中创建文章
func (r *PostRepository) CreateTx(tx *gorm.DB, post *model.Post) error {
	// 保存原始的 status 和 visibility 值
	originalStatus := post.Status
	originalVisibility := post.Visibility

	// 使用原生 SQL 插入，完全绕过 GORM 的 default 标签逻辑
	// 这样可以确保零值（0）能正确保存
	var publishedAtValue interface{}
	if post.PublishedAt != nil {
		publishedAtValue = post.PublishedAt
	} else {
		publishedAtValue = nil
	}

	// 使用原生 SQL INSERT，明确指定所有字段的值
	// PostgreSQL 使用 $1, $2... 占位符，GORM 会自动转换
	err := tx.Raw(`
		INSERT INTO posts (
			title, content, summary, cover,
			status, visibility, is_top,
			user_id, category_id, published_at,
			view_count, like_count,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW())
		RETURNING id
	`,
		post.Title,
		post.Content,
		post.Summary,
		post.Cover,
		originalStatus,     // 直接使用原始值，确保 0 能正确保存
		originalVisibility, // 直接使用原始值，确保 0 能正确保存
		post.IsTop,
		post.UserID,
		post.CategoryID,
		publishedAtValue,
		post.ViewCount,
		post.LikeCount,
	).Scan(&post.ID).Error

	if err != nil {
		return err
	}

	// 更新 post 对象的值，确保后续代码使用正确的值
	post.Status = originalStatus
	post.Visibility = originalVisibility

	// 更新全文搜索向量
	tx.Exec(
		"UPDATE posts SET search_tsv = setweight(to_tsvector('english', coalesce(title, '')), 'A') || setweight(to_tsvector('english', coalesce(content, '')), 'B') WHERE id = ?",
		post.ID,
	)

	return nil
}

// UpdateTx 在事务中更新文章
func (r *PostRepository) UpdateTx(tx *gorm.DB, post *model.Post) error {
	err := tx.Save(post).Error
	if err != nil {
		return err
	}

	// 更新全文搜索向量
	tx.Exec(
		"UPDATE posts SET search_tsv = setweight(to_tsvector('english', coalesce(title, '')), 'A') || setweight(to_tsvector('english', coalesce(content, '')), 'B') WHERE id = ?",
		post.ID,
	)

	return nil
}

// UpdateTagsTx 在事务中更新文章标签
func (r *PostRepository) UpdateTagsTx(tx *gorm.DB, postID uint, tagIDs []uint) error {
	var post model.Post
	if err := tx.First(&post, postID).Error; err != nil {
		return err
	}

	var tags []model.Tag
	if err := tx.Find(&tags, tagIDs).Error; err != nil {
		return err
	}

	return tx.Model(&post).Association("Tags").Replace(tags)
}

// CreateLikeTx 在事务中创建点赞记录
func (r *PostRepository) CreateLikeTx(tx *gorm.DB, like *model.PostLike) error {
	return tx.Create(like).Error
}

// DeleteLikeTx 在事务中删除点赞记录
func (r *PostRepository) DeleteLikeTx(tx *gorm.DB, postID uint, userID *uint, ip string) error {
	query := tx.Where("post_id = ?", postID)

	if userID != nil && *userID > 0 {
		query = query.Where("user_id = ?", *userID)
	} else {
		query = query.Where("ip = ? AND user_id IS NULL", ip)
	}

	return query.Delete(&model.PostLike{}).Error
}

// IncrementLikeCountTx 在事务中增加点赞数
func (r *PostRepository) IncrementLikeCountTx(tx *gorm.DB, id uint) error {
	return tx.Model(&model.Post{}).Where("id = ?", id).UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
}

// DecrementLikeCountTx 在事务中减少点赞数
func (r *PostRepository) DecrementLikeCountTx(tx *gorm.DB, id uint) error {
	return tx.Model(&model.Post{}).Where("id = ?", id).UpdateColumn("like_count", gorm.Expr("CASE WHEN like_count > 0 THEN like_count - 1 ELSE 0 END")).Error
}
