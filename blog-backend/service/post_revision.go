/*
 * 项目名称：blog-backend
 * 文件名称：post_revision.go
 * 创建时间：2026-04-28 17:48:56
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：文章修订版本业务逻辑层，负责审批流程的核心业务逻辑
 */
package service

import (
	"blog-backend/config"
	"blog-backend/model"
	"blog-backend/repository"
	"blog-backend/util"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PostRevisionService 文章修订版本服务
type PostRevisionService struct {
	revisionRepo        *repository.PostRevisionRepository
	postRepo            *repository.PostRepository
	collaboratorRepo    *repository.PostCollaboratorRepository
	notificationService *NotificationService
}

// NewPostRevisionService 创建文章修订版本服务实例
func NewPostRevisionService() *PostRevisionService {
	return &PostRevisionService{
		revisionRepo:     repository.NewPostRevisionRepository(),
		postRepo:         repository.NewPostRepository(),
		collaboratorRepo: repository.NewPostCollaboratorRepository(),
	}
}

// SetNotificationService 设置通知服务
func (s *PostRevisionService) SetNotificationService(notificationService *NotificationService) {
	s.notificationService = notificationService
}

// GetPendingList 获取待审批列表
func (s *PostRevisionService) GetPendingList(page, pageSize int) ([]*model.PostRevision, int64, error) {
	ctx := context.Background()
	return s.revisionRepo.GetPendingList(ctx, page, pageSize)
}

// GetRevisionDetail 获取修订详情
func (s *PostRevisionService) GetRevisionDetail(id uint) (*model.PostRevision, error) {
	ctx := context.Background()
	return s.revisionRepo.GetByID(ctx, id)
}

// ApproveRevision 审批通过（合并修改到原文章，并自动添加编辑者为协作者）
func (s *PostRevisionService) ApproveRevision(id, reviewerID uint) error {
	var editorEmail, editorName, reviewerName, postTitle string
	var postID, editorID uint

	// 使用事务合并修改
	err := s.postRepo.Transaction(func(tx *gorm.DB) error {
		// 在事务内重新查询并加行锁,防止并发审批
		var revision model.PostRevision
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("Editor").
			Preload("Post").
			First(&revision, id).Error; err != nil {
			return errors.New("修订版本不存在")
		}

		// 检查状态
		if revision.Status != "pending" {
			return errors.New("该修订已被审批")
		}

		// 保存编辑者和文章信息用于发送邮件
		if revision.Editor.Email != "" {
			editorEmail = revision.Editor.Email
			editorName = revision.Editor.Username
		}
		editorID = revision.EditorID
		postTitle = revision.Post.Title
		postID = revision.PostID

		// 获取审批人信息
		var reviewer model.User
		if err := tx.First(&reviewer, reviewerID).Error; err == nil {
			reviewerName = reviewer.Username
		}

		// 获取原文章
		post, err := s.postRepo.GetByID(revision.PostID)
		if err != nil {
			return errors.New("原文章不存在")
		}

		// 合并修改（只更新非NULL的字段）
		if revision.Title != nil {
			post.Title = *revision.Title
		}
		if revision.Content != nil {
			post.Content = *revision.Content
		}
		if revision.Summary != nil {
			post.Summary = *revision.Summary
		}
		if revision.Cover != nil {
			post.Cover = *revision.Cover
		}
		if revision.CategoryID != nil {
			post.CategoryID = *revision.CategoryID
		}
		if revision.Visibility != nil {
			post.Visibility = *revision.Visibility
		}
		if revision.IsTop != nil {
			post.IsTop = *revision.IsTop
		}

		// 更新文章
		if err := s.postRepo.UpdateTx(tx, post); err != nil {
			return err
		}

		// 更新修订状态
		now := time.Now()
		revision.Status = "approved"
		revision.ReviewerID = &reviewerID
		revision.ReviewedAt = &now

		return s.revisionRepo.UpdateStatusTx(tx, &revision)
	})

	// 事务成功后：自动添加编辑者为协作者 + 发送邮件通知
	if err == nil {
		// 自动添加编辑者为协作者（如果编辑者不是文章作者，且从未成为过协作者）
		if editorID > 0 && postID > 0 {
			// 获取文章检查作者是否为编辑者本身
			if post, getErr := s.postRepo.GetByID(postID); getErr == nil && post.UserID != editorID {
				// 检查是否曾经是协作者（包括被移除的）
				if everExisted, existErr := s.collaboratorRepo.EverExisted(postID, editorID); existErr == nil && !everExisted {
					_ = s.collaboratorRepo.Add(postID, editorID, 1)
				}
			}
		}

		// 发送邮件通知
		if editorEmail != "" {
			go func() {
				// 从数据库获取网站名称
				settingRepo := repository.NewSettingRepository()
				siteName := config.Cfg.Email.SiteName // 默认使用配置文件中的值
				if siteNameSetting, err := settingRepo.GetByKey("site_name"); err == nil && siteNameSetting.Value != "" {
					siteName = siteNameSetting.Value
				}

				emailConfig := util.EmailConfig{
					Host:     config.Cfg.Email.Host,
					Port:     config.Cfg.Email.Port,
					Username: config.Cfg.Email.Username,
					Password: config.Cfg.Email.Password,
					FromName: config.Cfg.Email.FromName,
					SiteName: siteName,
				}

				postURL := fmt.Sprintf("%s/post/%d", config.Cfg.App.BlogURL, postID)
				_ = util.SendRevisionApprovedEmail(emailConfig, editorEmail, editorName, postTitle, reviewerName, postURL)
			}()
		}

		// 发送应用内通知
		if s.notificationService != nil {
			userRepo := repository.NewUserRepository()
			reviewer, getErr := userRepo.GetByID(reviewerID)
			if getErr == nil {
				go s.notificationService.NotifyRevisionApproved(reviewer, editorID, postTitle, postID, id)
			}
		}
	}

	return err
}

// RejectRevision 审批拒绝
func (s *PostRevisionService) RejectRevision(id, reviewerID uint, reason string) error {
	var editorEmail, editorName, reviewerName, postTitle string
	var postID, editorID uint

	// 使用事务更新状态
	err := s.postRepo.Transaction(func(tx *gorm.DB) error {
		// 在事务内重新查询并加行锁,防止并发审批
		var revision model.PostRevision
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("Editor").
			Preload("Post").
			First(&revision, id).Error; err != nil {
			return errors.New("修订版本不存在")
		}

		// 检查状态
		if revision.Status != "pending" {
			return errors.New("该修订已被审批")
		}

		// 保存编辑者和文章信息用于发送邮件
		if revision.Editor.Email != "" {
			editorEmail = revision.Editor.Email
			editorName = revision.Editor.Username
		}
		editorID = revision.EditorID
		postTitle = revision.Post.Title
		postID = revision.PostID

		// 获取审批人信息
		var reviewer model.User
		if err := tx.First(&reviewer, reviewerID).Error; err == nil {
			reviewerName = reviewer.Username
		}

		// 更新修订状态
		now := time.Now()
		revision.Status = "rejected"
		revision.ReviewerID = &reviewerID
		revision.ReviewedAt = &now
		revision.RejectReason = reason

		return s.revisionRepo.UpdateStatusTx(tx, &revision)
	})

	// 事务成功后发送邮件通知
	if err == nil && editorEmail != "" {
		go func() {
			// 从数据库获取网站名称
			settingRepo := repository.NewSettingRepository()
			siteName := config.Cfg.Email.SiteName // 默认使用配置文件中的值
			if siteNameSetting, err := settingRepo.GetByKey("site_name"); err == nil && siteNameSetting.Value != "" {
				siteName = siteNameSetting.Value
			}

			emailConfig := util.EmailConfig{
				Host:     config.Cfg.Email.Host,
				Port:     config.Cfg.Email.Port,
				Username: config.Cfg.Email.Username,
				Password: config.Cfg.Email.Password,
				FromName: config.Cfg.Email.FromName,
				SiteName: siteName,
			}

			postURL := fmt.Sprintf("%s/admin/post/edit/%d", config.Cfg.App.BlogURL, postID)
			_ = util.SendRevisionRejectedEmail(emailConfig, editorEmail, editorName, postTitle, reviewerName, reason, postURL)
		}()
	}

	// 发送应用内通知
	if err == nil && s.notificationService != nil {
		userRepo := repository.NewUserRepository()
		reviewer, getErr := userRepo.GetByID(reviewerID)
		if getErr == nil {
			go s.notificationService.NotifyRevisionRejected(reviewer, editorID, postTitle, postID, id, reason)
		}
	}

	return err
}

// GetRevisionDiff 获取修改对比
func (s *PostRevisionService) GetRevisionDiff(id uint) (map[string]interface{}, error) {
	ctx := context.Background()

	// 获取修订版本
	revision, err := s.revisionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("修订版本不存在")
	}

	// 获取原文章
	post, err := s.postRepo.GetByID(revision.PostID)
	if err != nil {
		return nil, errors.New("原文章不存在")
	}

	// 构建对比数据
	diff := make(map[string]interface{})
	diff["post_id"] = post.ID
	diff["post_title"] = post.Title
	diff["editor"] = revision.Editor
	diff["editor_comment"] = revision.EditorComment
	diff["created_at"] = revision.CreatedAt

	changes := make([]map[string]interface{}, 0)

	if revision.Title != nil && *revision.Title != post.Title {
		changes = append(changes, map[string]interface{}{
			"field": "title",
			"old":   post.Title,
			"new":   *revision.Title,
		})
	}
	if revision.Content != nil && *revision.Content != post.Content {
		changes = append(changes, map[string]interface{}{
			"field": "content",
			"old":   post.Content,
			"new":   *revision.Content,
		})
	}
	if revision.Summary != nil && *revision.Summary != post.Summary {
		changes = append(changes, map[string]interface{}{
			"field": "summary",
			"old":   post.Summary,
			"new":   *revision.Summary,
		})
	}
	if revision.Cover != nil && *revision.Cover != post.Cover {
		changes = append(changes, map[string]interface{}{
			"field": "cover",
			"old":   post.Cover,
			"new":   *revision.Cover,
		})
	}
	if revision.CategoryID != nil && *revision.CategoryID != post.CategoryID {
		changes = append(changes, map[string]interface{}{
			"field": "category_id",
			"old":   post.CategoryID,
			"new":   *revision.CategoryID,
		})
	}
	if revision.Visibility != nil && *revision.Visibility != post.Visibility {
		changes = append(changes, map[string]interface{}{
			"field": "visibility",
			"old":   post.Visibility,
			"new":   *revision.Visibility,
		})
	}
	if revision.IsTop != nil && *revision.IsTop != post.IsTop {
		changes = append(changes, map[string]interface{}{
			"field": "is_top",
			"old":   post.IsTop,
			"new":   *revision.IsTop,
		})
	}

	// 处理标签对比
	if revision.TagIDs != nil {
		var newTagIDs []uint
		if err := json.Unmarshal([]byte(*revision.TagIDs), &newTagIDs); err == nil {
			oldTagIDs := make([]uint, 0)
			if post.Tags != nil {
				for _, tag := range post.Tags {
					oldTagIDs = append(oldTagIDs, tag.ID)
				}
			}
			changes = append(changes, map[string]interface{}{
				"field": "tag_ids",
				"old":   oldTagIDs,
				"new":   newTagIDs,
			})
		}
	}

	diff["changes"] = changes
	diff["changes_count"] = len(changes)

	return diff, nil
}

// GetPendingCount 获取待审批数量
func (s *PostRevisionService) GetPendingCount() (int64, error) {
	ctx := context.Background()
	return s.revisionRepo.GetPendingCount(ctx)
}

// GetPostRevisions 获取文章修订历史
func (s *PostRevisionService) GetPostRevisions(postID uint) ([]*model.PostRevision, error) {
	ctx := context.Background()
	return s.revisionRepo.GetByPostID(ctx, postID)
}

// WithdrawRevision 撤回修订（仅创建者可撤回）
func (s *PostRevisionService) WithdrawRevision(id, editorID uint) error {
	// 使用事务更新状态
	err := s.postRepo.Transaction(func(tx *gorm.DB) error {
		// 在事务内重新查询并加行锁,防止并发操作
		var revision model.PostRevision
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&revision, id).Error; err != nil {
			return errors.New("修订版本不存在")
		}

		// 检查权限：只有创建者可以撤回
		if revision.EditorID != editorID {
			return errors.New("无权撤回此修订")
		}

		// 检查状态：只有待审批状态可以撤回
		if revision.Status != "pending" {
			return errors.New("只能撤回待审批的修订")
		}

		// 更新修订状态
		revision.Status = "withdrawn"

		return s.revisionRepo.UpdateStatusTx(tx, &revision)
	})

	return err
}

// GetUserRevisions 获取指定用户的修订记录
func (s *PostRevisionService) GetUserRevisions(editorID uint, page, pageSize int, status string) ([]*model.PostRevision, int64, error) {
	ctx := context.Background()
	return s.revisionRepo.GetByEditorID(ctx, editorID, page, pageSize, status)
}
