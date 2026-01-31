/*
 * 项目名称：blog-backend
 * 文件名称：tag.go
 * 创建时间：2026-01-31 16:34:35
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：标签业务逻辑层，提供文章标签的增删改查业务处理，支持Redis缓存
 */
package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"blog-backend/db"
	"blog-backend/model"
	"blog-backend/repository"

	"gorm.io/gorm"
)

// TagService 标签业务逻辑层结构体
type TagService struct {
	repo *repository.TagRepository
}

// NewTagService 创建标签业务逻辑层实例
func NewTagService() *TagService {
	return &TagService{
		repo: repository.NewTagRepository(),
	}
}

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name      string  `json:"name" binding:"required"`
	Color     string  `json:"color"`
	TextColor *string `json:"text_color"`
	FontSize  *int    `json:"font_size"`
}

// UpdateTagRequest 更新标签请求
type UpdateTagRequest struct {
	Name      string  `json:"name"`
	Color     string  `json:"color"`
	TextColor *string `json:"text_color"`
	FontSize  *int    `json:"font_size"`
}

// Create 创建标签
func (s *TagService) Create(req *CreateTagRequest) (*model.Tag, error) {
	// 检查标签名是否已存在
	if _, err := s.repo.GetByName(req.Name); err == nil {
		return nil, errors.New("标签名已存在")
	}

	tag := &model.Tag{
		Name:      req.Name,
		Color:     req.Color,
		TextColor: req.TextColor,
		FontSize:  req.FontSize,
	}

	if err := s.repo.Create(tag); err != nil {
		return nil, errors.New("标签创建失败")
	}

	// 写操作成功后，删除标签相关缓存
	go func() {
		ctx := context.Background()
		db.RDB.Del(ctx, "tag:list")
		db.RDB.Del(ctx, "tag:stats:top10")
		db.RDB.Del(ctx, "blog:author_profile")
	}()

	return tag, nil
}

// GetByID 获取标签详情
func (s *TagService) GetByID(id uint) (*model.Tag, error) {
	tag, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("标签不存在")
		}
		return nil, errors.New("获取标签失败")
	}
	return tag, nil
}

// Update 更新标签
func (s *TagService) Update(id uint, req *UpdateTagRequest) (*model.Tag, error) {
	tag, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("标签不存在")
	}

	// 如果修改了名称，检查新名称是否已存在
	if req.Name != "" && req.Name != tag.Name {
		if _, err := s.repo.GetByName(req.Name); err == nil {
			return nil, errors.New("标签名已存在")
		}
		tag.Name = req.Name
	}

	if req.Color != "" {
		tag.Color = req.Color
	}

	// 更新文字颜色和字体大小（即使是 nil 也要更新，以支持清空操作）
	tag.TextColor = req.TextColor
	tag.FontSize = req.FontSize

	if err := s.repo.Update(tag); err != nil {
		return nil, errors.New("标签更新失败")
	}

	// 写操作成功后，删除标签相关缓存
	go func() {
		ctx := context.Background()
		db.RDB.Del(ctx, "tag:list")
		db.RDB.Del(ctx, "tag:stats:top10")
		db.RDB.Del(ctx, "blog:author_profile")
	}()

	return tag, nil
}

// Delete 删除标签
func (s *TagService) Delete(id uint) error {
	tag, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("标签不存在")
	}

	// 检查标签下是否有文章
	if tag.PostCount > 0 {
		return errors.New("该标签下还有文章，无法删除")
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// 删除成功后，清理标签相关缓存
	go func() {
		ctx := context.Background()
		db.RDB.Del(ctx, "tag:list")
		db.RDB.Del(ctx, "tag:stats:top10")
		db.RDB.Del(ctx, "blog:author_profile")
	}()

	return nil
}

// List 获取标签列表
func (s *TagService) List() ([]model.Tag, error) {
	ctx := context.Background()
	cacheKey := "tag:list"

	// 1. 先尝试从 Redis 获取缓存
	if cached, err := db.RDB.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
		var list []model.Tag
		if err := json.Unmarshal([]byte(cached), &list); err == nil {
			return list, nil
		}
	}

	// 2. 缓存未命中，从数据库获取
	list, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	// 3. 写入缓存，设置过期时间（例如 30 分钟）
	if data, err := json.Marshal(list); err == nil {
		_ = db.RDB.Set(ctx, cacheKey, string(data), 30*time.Minute).Err()
	}

	return list, nil
}
