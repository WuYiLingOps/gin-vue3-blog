package repository

import (
	"blog-backend/db"
	"blog-backend/model"
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

// List 获取文章列表
func (r *PostRepository) List(page, pageSize int, categoryID uint, keyword string, status int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	offset := (page - 1) * pageSize
	query := db.DB.Model(&model.Post{})

	// 筛选条件
	if status >= 0 {
		query = query.Where("status = ?", status)
	}
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		// 使用PostgreSQL全文搜索
		query = query.Where("search_tsv @@ plainto_tsquery('english', ?)", keyword)
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
func (r *PostRepository) GetRecentPosts(limit int) ([]model.Post, error) {
	var posts []model.Post
	err := db.DB.Preload("User").Preload("Category").
		Where("status = 1").
		Order("created_at DESC").
		Limit(limit).Find(&posts).Error
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

