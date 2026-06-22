/*
 * 项目名称：blog-backend
 * 文件名称：notification.go
 * 创建时间：2026-06-21 11:00:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：通知业务逻辑层，提供通知调度核心功能，支持应用内通知、邮件通知和WebSocket推送
 */
package service

import (
	"blog-backend/model"
	"blog-backend/pkg/email"
	"blog-backend/repository"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// NotificationService 通知服务
type NotificationService struct {
	repo        *repository.NotificationRepository
	userRepo    *repository.UserRepository
	settingRepo *repository.SettingRepository
	emailClient *email.Client
	hub         *NotificationHub
}

// NewNotificationService 创建通知服务
func NewNotificationService(hub *NotificationHub, emailClient *email.Client) *NotificationService {
	return &NotificationService{
		repo:        repository.NewNotificationRepository(),
		userRepo:    repository.NewUserRepository(),
		settingRepo: repository.NewSettingRepository(),
		emailClient: emailClient,
		hub:         hub,
	}
}

// NotifyCommentReply 评论回复通知
func (s *NotificationService) NotifyCommentReply(commenter *model.User, repliedUserID uint, postTitle string, postID uint, content string) {
	commenterID := commenter.ID
	// 不通知自己
	if commenterID == repliedUserID {
		return
	}

	// 先推送WebSocket（即时反馈）
	s.hub.SendNotificationToUser(repliedUserID, map[string]interface{}{
		"type":       "comment_reply",
		"title":      fmt.Sprintf("%s 回复了您的评论", commenter.Nickname),
		"content":    content,
		"sender":     commenter,
		"created_at": time.Now(),
	})

	// 再异步写库（不阻塞推送）
	go func() {
		notification := &model.Notification{
			Type:       "comment_reply",
			Title:      fmt.Sprintf("%s 回复了您的评论", commenter.Nickname),
			Content:    content,
			SenderID:   &commenterID,
			TargetType: "post",
			TargetID:   &postID,
		}
		extra := map[string]interface{}{"post_id": postID, "post_title": postTitle}
		extraJSON, _ := json.Marshal(extra)
		notification.Extra = stringPtr(string(extraJSON))

		if err := s.repo.CreateNotification(notification, []uint{repliedUserID}); err != nil {
			log.Printf("创建评论回复通知失败: %v", err)
		}
	}()

	// 异步发送邮件
	go func() {
		s.sendCommentReplyEmail(repliedUserID, commenter.Nickname, postTitle, content, postID)
	}()
}

// NotifyCommentNew 新评论通知（管理员）
func (s *NotificationService) NotifyCommentNew(commenter *model.User, postTitle string, postID uint, content string) {
	commenterID := commenter.ID

	// 获取所有管理员
	admins, err := s.userRepo.GetAdmins()
	if err != nil {
		return
	}

	// 过滤掉评论者自己
	var recipientIDs []uint
	for _, admin := range admins {
		if admin.ID != commenterID {
			recipientIDs = append(recipientIDs, admin.ID)
		}
	}

	if len(recipientIDs) == 0 {
		return
	}

	// 先推送WebSocket（即时反馈）
	for _, adminID := range recipientIDs {
		s.hub.SendNotificationToUser(adminID, map[string]interface{}{
			"type":       "comment_new",
			"title":      fmt.Sprintf("%s 在《%s》发表了新评论", commenter.Nickname, postTitle),
			"content":    content,
			"sender":     commenter,
			"created_at": time.Now(),
		})
	}

	// 再异步写库（不阻塞推送）
	go func() {
		notification := &model.Notification{
			Type:       "comment_new",
			Title:      fmt.Sprintf("%s 在《%s》发表了新评论", commenter.Nickname, postTitle),
			Content:    content,
			SenderID:   &commenterID,
			TargetType: "post",
			TargetID:   &postID,
		}
		extra := map[string]interface{}{"post_id": postID, "post_title": postTitle}
		extraJSON, _ := json.Marshal(extra)
		notification.Extra = stringPtr(string(extraJSON))

		if err := s.repo.CreateNotification(notification, recipientIDs); err != nil {
			log.Printf("创建新评论通知失败: %v", err)
		}
	}()

	// 异步发送邮件给所有管理员
	go func() {
		for _, adminID := range recipientIDs {
			s.sendCommentNewEmail(adminID, commenter.Nickname, postTitle, content, postID)
		}
	}()
}

// NotifyArticlePush 文章推送通知
func (s *NotificationService) NotifyArticlePush(articleID uint, articleTitle string) {
	// 获取所有订阅者
	subscribers, err := s.userRepo.GetSubscribers()
	if err != nil {
		return
	}

	if len(subscribers) == 0 {
		return
	}

	var recipientIDs []uint
	for _, subscriber := range subscribers {
		recipientIDs = append(recipientIDs, subscriber.ID)
	}

	// 先推送WebSocket（即时反馈）
	for _, subscriberID := range recipientIDs {
		s.hub.SendNotificationToUser(subscriberID, map[string]interface{}{
			"type":       "article_push",
			"title":      fmt.Sprintf("新文章发布：%s", articleTitle),
			"content":    fmt.Sprintf("文章《%s》已发布，快来阅读吧！", articleTitle),
			"created_at": time.Now(),
		})
	}

	// 再异步写库
	go func() {
		notification := &model.Notification{
			Type:       "article_push",
			Title:      fmt.Sprintf("新文章发布：%s", articleTitle),
			Content:    fmt.Sprintf("文章《%s》已发布，快来阅读吧！", articleTitle),
			TargetType: "post",
			TargetID:   &articleID,
		}
		extra := map[string]interface{}{"post_id": articleID, "post_title": articleTitle}
		extraJSON, _ := json.Marshal(extra)
		notification.Extra = stringPtr(string(extraJSON))

		if err := s.repo.CreateNotification(notification, recipientIDs); err != nil {
			log.Printf("创建文章推送通知失败: %v", err)
		}
	}()
}

// NotifyRevisionSubmitted 文章修订提交通知（超级管理员）
func (s *NotificationService) NotifyRevisionSubmitted(editor *model.User, postTitle string, postID uint, revisionID uint) {
	editorID := editor.ID

	// 获取所有超级管理员
	superAdmins, err := s.userRepo.GetSuperAdmins(context.Background())
	if err != nil {
		return
	}

	// 过滤掉编辑者自己
	var recipientIDs []uint
	for _, admin := range superAdmins {
		if admin.ID != editorID {
			recipientIDs = append(recipientIDs, admin.ID)
		}
	}

	if len(recipientIDs) == 0 {
		return
	}

	// 先推送WebSocket（即时反馈）
	for _, adminID := range recipientIDs {
		s.hub.SendNotificationToUser(adminID, map[string]interface{}{
			"type":       "revision_submitted",
			"title":      fmt.Sprintf("新修订待审批：%s", postTitle),
			"content":    fmt.Sprintf("%s 提交了文章《%s》的修订，请前往审批", editor.Nickname, postTitle),
			"sender":     editor,
			"created_at": time.Now(),
		})
	}

	// 再异步写库（不阻塞推送）
	go func() {
		notification := &model.Notification{
			Type:       "revision_submitted",
			Title:      fmt.Sprintf("新修订待审批：%s", postTitle),
			Content:    fmt.Sprintf("%s 提交了文章《%s》的修订，请前往审批", editor.Nickname, postTitle),
			SenderID:   &editorID,
			TargetType: "post_revision",
			TargetID:   &revisionID,
		}
		extra := map[string]interface{}{"post_id": postID, "post_title": postTitle, "revision_id": revisionID}
		extraJSON, _ := json.Marshal(extra)
		notification.Extra = stringPtr(string(extraJSON))

		if err := s.repo.CreateNotification(notification, recipientIDs); err != nil {
			log.Printf("创建修订提交通知失败: %v", err)
		}
	}()
}

// NotifyRevisionApproved 文章修订通过通知（编辑者）
func (s *NotificationService) NotifyRevisionApproved(reviewer *model.User, editorID uint, postTitle string, postID uint, revisionID uint) {
	// 不通知自己
	if reviewer.ID == editorID {
		return
	}

	// 先推送WebSocket（即时反馈）
	s.hub.SendNotificationToUser(editorID, map[string]interface{}{
		"type":       "revision_approved",
		"title":      fmt.Sprintf("修订已通过：%s", postTitle),
		"content":    fmt.Sprintf("您提交的文章《%s》修订已被 %s 通过", postTitle, reviewer.Nickname),
		"sender":     reviewer,
		"created_at": time.Now(),
	})

	// 再异步写库（不阻塞推送）
	go func() {
		reviewerID := reviewer.ID
		notification := &model.Notification{
			Type:       "revision_approved",
			Title:      fmt.Sprintf("修订已通过：%s", postTitle),
			Content:    fmt.Sprintf("您提交的文章《%s》修订已被 %s 通过", postTitle, reviewer.Nickname),
			SenderID:   &reviewerID,
			TargetType: "post_revision",
			TargetID:   &revisionID,
		}
		extra := map[string]interface{}{"post_id": postID, "post_title": postTitle, "revision_id": revisionID}
		extraJSON, _ := json.Marshal(extra)
		notification.Extra = stringPtr(string(extraJSON))

		if err := s.repo.CreateNotification(notification, []uint{editorID}); err != nil {
			log.Printf("创建修订通过通知失败: %v", err)
		}
	}()
}

// NotifyRevisionRejected 文章修订拒绝通知（编辑者）
func (s *NotificationService) NotifyRevisionRejected(reviewer *model.User, editorID uint, postTitle string, postID uint, revisionID uint, reason string) {
	// 不通知自己
	if reviewer.ID == editorID {
		return
	}

	// 先推送WebSocket（即时反馈）
	s.hub.SendNotificationToUser(editorID, map[string]interface{}{
		"type":       "revision_rejected",
		"title":      fmt.Sprintf("修订被拒绝：%s", postTitle),
		"content":    fmt.Sprintf("您提交的文章《%s》修订已被 %s 拒绝，原因：%s", postTitle, reviewer.Nickname, reason),
		"sender":     reviewer,
		"created_at": time.Now(),
	})

	// 再异步写库（不阻塞推送）
	go func() {
		reviewerID := reviewer.ID
		notification := &model.Notification{
			Type:       "revision_rejected",
			Title:      fmt.Sprintf("修订被拒绝：%s", postTitle),
			Content:    fmt.Sprintf("您提交的文章《%s》修订已被 %s 拒绝，原因：%s", postTitle, reviewer.Nickname, reason),
			SenderID:   &reviewerID,
			TargetType: "post_revision",
			TargetID:   &revisionID,
		}
		extra := map[string]interface{}{"post_id": postID, "post_title": postTitle, "revision_id": revisionID, "reason": reason}
		extraJSON, _ := json.Marshal(extra)
		notification.Extra = stringPtr(string(extraJSON))

		if err := s.repo.CreateNotification(notification, []uint{editorID}); err != nil {
			log.Printf("创建修订拒绝通知失败: %v", err)
		}
	}()
}

// NotifyFriendLinkApply 友链申请通知（管理员）
func (s *NotificationService) NotifyFriendLinkApply(applicantName string, siteName string, siteURL string, friendLinkID uint) {
	// 获取所有管理员
	admins, err := s.userRepo.GetAdmins()
	if err != nil {
		return
	}

	var recipientIDs []uint
	for _, admin := range admins {
		recipientIDs = append(recipientIDs, admin.ID)
	}

	if len(recipientIDs) == 0 {
		return
	}

	// 先推送WebSocket（即时反馈）
	for _, adminID := range recipientIDs {
		s.hub.SendNotificationToUser(adminID, map[string]interface{}{
			"type":       "friend_link_apply",
			"title":      fmt.Sprintf("新友链申请：%s", siteName),
			"content":    fmt.Sprintf("%s 申请了友情链接：%s (%s)", applicantName, siteName, siteURL),
			"created_at": time.Now(),
		})
	}

	// 再异步写库（不阻塞推送）
	go func() {
		notification := &model.Notification{
			Type:       "friend_link_apply",
			Title:      fmt.Sprintf("新友链申请：%s", siteName),
			Content:    fmt.Sprintf("%s 申请了友情链接：%s (%s)", applicantName, siteName, siteURL),
			TargetType: "friend_link",
			TargetID:   &friendLinkID,
		}
		extra := map[string]interface{}{"friend_link_id": friendLinkID, "site_name": siteName, "site_url": siteURL}
		extraJSON, _ := json.Marshal(extra)
		notification.Extra = stringPtr(string(extraJSON))

		if err := s.repo.CreateNotification(notification, recipientIDs); err != nil {
			log.Printf("创建友链申请通知失败: %v", err)
		}
	}()
}

// NotifyFriendAbnormal 友链异常通知（管理员）
func (s *NotificationService) NotifyFriendAbnormal(friendLinkID uint, friendLinkName string, friendLinkURL string, failCount int) {
	// 获取所有管理员
	admins, err := s.userRepo.GetAdmins()
	if err != nil {
		return
	}

	var recipientIDs []uint
	for _, admin := range admins {
		recipientIDs = append(recipientIDs, admin.ID)
	}

	if len(recipientIDs) == 0 {
		return
	}

	// 先推送WebSocket（即时反馈）
	for _, adminID := range recipientIDs {
		s.hub.SendNotificationToUser(adminID, map[string]interface{}{
			"type":       "friend_abnormal",
			"title":      fmt.Sprintf("友链异常：%s", friendLinkName),
			"content":    fmt.Sprintf("友链 %s (%s) 连续访问失败 %d 次", friendLinkName, friendLinkURL, failCount),
			"created_at": time.Now(),
		})
	}

	// 再异步写库
	go func() {
		notification := &model.Notification{
			Type:       "friend_abnormal",
			Title:      fmt.Sprintf("友链异常：%s", friendLinkName),
			Content:    fmt.Sprintf("友链 %s (%s) 连续访问失败 %d 次", friendLinkName, friendLinkURL, failCount),
			TargetType: "friend_link",
			TargetID:   &friendLinkID,
		}
		extra := map[string]interface{}{
			"friend_link_id":   friendLinkID,
			"friend_link_name": friendLinkName,
			"friend_link_url":  friendLinkURL,
			"fail_count":       failCount,
		}
		extraJSON, _ := json.Marshal(extra)
		notification.Extra = stringPtr(string(extraJSON))

		if err := s.repo.CreateNotification(notification, recipientIDs); err != nil {
			log.Printf("创建友链异常通知失败: %v", err)
		}
	}()
}

// GetNotifications 获取用户通知列表
func (s *NotificationService) GetNotifications(userID uint, page, pageSize int) ([]model.UserNotification, int64, error) {
	return s.repo.GetByUserID(userID, page, pageSize)
}

// GetUnreadCount 获取用户未读通知数量
func (s *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	return s.repo.GetUnreadCount(userID)
}

// MarkAsRead 标记单条通知为已读
func (s *NotificationService) MarkAsRead(notificationID, userID uint) error {
	return s.repo.MarkAsRead(notificationID, userID)
}

// MarkAllAsRead 标记用户所有通知为已读
func (s *NotificationService) MarkAllAsRead(userID uint) error {
	return s.repo.MarkAllAsRead(userID)
}

// DeleteNotification 删除通知
func (s *NotificationService) DeleteNotification(id uint) error {
	return s.repo.DeleteNotification(id)
}

// sendCommentReplyEmail 发送评论回复邮件
func (s *NotificationService) sendCommentReplyEmail(userID uint, replierName, postTitle, content string, postID uint) {
	if s.emailClient == nil {
		return
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		log.Printf("获取用户信息失败: %v", err)
		return
	}

	siteName := s.getSiteName()
	articleURL := fmt.Sprintf("%s/post/%d", s.getSiteURL(), postID)

	if err := s.emailClient.SendCommentReplyNotification(user.Email, siteName, replierName, postTitle, content, articleURL); err != nil {
		log.Printf("发送评论回复邮件失败: %v", err)
	}
}

// sendCommentNewEmail 发送新评论邮件
func (s *NotificationService) sendCommentNewEmail(adminID uint, commenterName, postTitle, content string, postID uint) {
	if s.emailClient == nil {
		return
	}

	admin, err := s.userRepo.GetByID(adminID)
	if err != nil {
		log.Printf("获取管理员信息失败: %v", err)
		return
	}

	siteName := s.getSiteName()
	articleURL := fmt.Sprintf("%s/post/%d", s.getSiteURL(), postID)

	if err := s.emailClient.SendCommentNewNotification(admin.Email, siteName, commenterName, postTitle, content, articleURL); err != nil {
		log.Printf("发送新评论邮件失败: %v", err)
	}
}

// getSiteName 获取网站名称
func (s *NotificationService) getSiteName() string {
	setting, err := s.settingRepo.GetByKey("site_name")
	if err != nil || setting == nil {
		return "菱风叙"
	}
	return setting.Value
}

// getSiteURL 获取网站URL
func (s *NotificationService) getSiteURL() string {
	setting, err := s.settingRepo.GetByKey("site_url")
	if err != nil || setting == nil || setting.Value == "" {
		return "http://localhost:3000"
	}
	return setting.Value
}

// stringPtr 返回字符串指针
func stringPtr(s string) *string {
	return &s
}
