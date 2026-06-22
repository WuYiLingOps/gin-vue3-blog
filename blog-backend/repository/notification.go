/*
 * 项目名称：blog-backend
 * 文件名称：notification.go
 * 创建时间：2026-06-21 10:00:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：通知数据访问层，提供通知的数据库操作功能，支持通知事件和用户通知关联
 */
package repository

import (
	"blog-backend/db"
	"blog-backend/model"

	"gorm.io/gorm"
)

// NotificationRepository 通知数据访问层结构体
type NotificationRepository struct{}

// NewNotificationRepository 创建通知数据访问层实例
func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{}
}

// CreateNotification 创建通知（事务）
func (r *NotificationRepository) CreateNotification(notification *model.Notification, recipientIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// 创建通知事件
		if err := tx.Create(notification).Error; err != nil {
			return err
		}

		// 批量创建用户通知关联
		if len(recipientIDs) > 0 {
			userNotifications := make([]model.UserNotification, len(recipientIDs))
			for i, userID := range recipientIDs {
				userNotifications[i] = model.UserNotification{
					NotificationID: notification.ID,
					UserID:         userID,
					IsRead:         false,
					EmailSent:      false,
				}
			}
			if err := tx.Create(&userNotifications).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetByUserID 根据用户ID获取通知列表（分页）
func (r *NotificationRepository) GetByUserID(userID uint, page, pageSize int) ([]model.UserNotification, int64, error) {
	var userNotifications []model.UserNotification
	var total int64

	offset := (page - 1) * pageSize

	// 统计总数
	if err := db.DB.Model(&model.UserNotification{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询通知列表，预加载通知事件和发送者
	err := db.DB.Preload("Notification").Preload("Notification.Sender").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&userNotifications).Error

	return userNotifications, total, err
}

// GetUnreadCount 获取用户未读通知数量
func (r *NotificationRepository) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := db.DB.Model(&model.UserNotification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

// MarkAsRead 标记单条通知为已读
func (r *NotificationRepository) MarkAsRead(notificationID, userID uint) error {
	return db.DB.Model(&model.UserNotification{}).
		Where("notification_id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true).Error
}

// MarkAllAsRead 标记用户所有通知为已读
func (r *NotificationRepository) MarkAllAsRead(userID uint) error {
	return db.DB.Model(&model.UserNotification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}

// DeleteNotification 软删除通知（设置status=0）
func (r *NotificationRepository) DeleteNotification(id uint) error {
	return db.DB.Model(&model.Notification{}).Where("id = ?", id).Update("status", 0).Error
}

// GetByID 根据ID获取通知
func (r *NotificationRepository) GetByID(id uint) (*model.Notification, error) {
	var notification model.Notification
	err := db.DB.Preload("Sender").First(&notification, id).Error
	return &notification, err
}

// UpdateEmailSent 更新邮件发送状态
func (r *NotificationRepository) UpdateEmailSent(notificationID, userID uint) error {
	return db.DB.Model(&model.UserNotification{}).
		Where("notification_id = ? AND user_id = ?", notificationID, userID).
		Update("email_sent", true).Error
}
