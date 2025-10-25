package service

import (
	"errors"

	"blog-backend/model"
	"blog-backend/repository"
	"gorm.io/gorm"
)

type CommentService struct {
	repo     *repository.CommentRepository
	postRepo *repository.PostRepository
}

func NewCommentService() *CommentService {
	return &CommentService{
		repo:     repository.NewCommentRepository(),
		postRepo: repository.NewPostRepository(),
	}
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required"`
	PostID   uint   `json:"post_id" binding:"required"`
	ParentID *uint  `json:"parent_id"`
}

// UpdateCommentRequest 更新评论请求
type UpdateCommentRequest struct {
	Content string `json:"content"`
}

// Create 创建评论
func (s *CommentService) Create(userID uint, req *CreateCommentRequest) (*model.Comment, error) {
	// 检查文章是否存在
	if _, err := s.postRepo.GetByID(req.PostID); err != nil {
		return nil, errors.New("文章不存在")
	}

	// 如果是回复评论，检查父评论是否存在
	if req.ParentID != nil {
		if _, err := s.repo.GetByID(*req.ParentID); err != nil {
			return nil, errors.New("父评论不存在")
		}
	}

	comment := &model.Comment{
		Content:  req.Content,
		PostID:   req.PostID,
		UserID:   userID,
		ParentID: req.ParentID,
		Status:   1,
	}

	if err := s.repo.Create(comment); err != nil {
		return nil, errors.New("评论创建失败")
	}

	return s.repo.GetByID(comment.ID)
}

// GetByID 获取评论详情
func (s *CommentService) GetByID(id uint) (*model.Comment, error) {
	comment, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("评论不存在")
		}
		return nil, errors.New("获取评论失败")
	}
	return comment, nil
}

// Update 更新评论
func (s *CommentService) Update(id, userID uint, role string, req *UpdateCommentRequest) (*model.Comment, error) {
	comment, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("评论不存在")
	}

	// 权限检查：只有作者和管理员可以修改
	if comment.UserID != userID && role != "admin" {
		return nil, errors.New("无权限修改此评论")
	}

	if req.Content != "" {
		comment.Content = req.Content
	}

	if err := s.repo.Update(comment); err != nil {
		return nil, errors.New("评论更新失败")
	}

	return comment, nil
}

// Delete 删除评论
func (s *CommentService) Delete(id, userID uint, role string) error {
	comment, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("评论不存在")
	}

	// 权限检查：只有作者和管理员可以删除
	if comment.UserID != userID && role != "admin" {
		return errors.New("无权限删除此评论")
	}

	return s.repo.Delete(id)
}

// GetByPostID 获取文章的评论列表
func (s *CommentService) GetByPostID(postID uint) ([]model.Comment, error) {
	return s.repo.GetByPostID(postID)
}

// List 获取评论列表（管理后台）
func (s *CommentService) List(page, pageSize int) ([]model.Comment, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.repo.List(page, pageSize)
}

// UpdateStatus 更新评论状态
func (s *CommentService) UpdateStatus(id uint, status int) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return errors.New("评论不存在")
	}

	return s.repo.UpdateStatus(id, status)
}

