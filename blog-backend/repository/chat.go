package repository

import (
	"blog-backend/db"
	"blog-backend/model"
)

// ChatRepository 聊天仓储
type ChatRepository struct{}

// NewChatRepository 创建聊天仓储
func NewChatRepository() *ChatRepository {
	return &ChatRepository{}
}

// Create 创建消息
func (r *ChatRepository) Create(message *model.ChatMessage) error {
	return db.DB.Create(message).Error
}

// GetMessages 获取消息列表（分页）
// includeAnnouncementOnly 表示是否包含仅投递到公告栏的广播
func (r *ChatRepository) GetMessages(page, pageSize int, includeAnnouncementOnly bool) ([]model.ChatMessage, int64, error) {
	var messages []model.ChatMessage
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	totalQuery := db.DB.Model(&model.ChatMessage{}).Where("status = ?", 1)
	if !includeAnnouncementOnly {
		totalQuery = totalQuery.Where("(is_broadcast = ? OR target IN ? OR target IS NULL OR target = '')",
			false, []string{"chat", "both"})
	}
	if err := totalQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取消息列表
	query := db.DB.Where("status = ?", 1)
	if !includeAnnouncementOnly {
		query = query.Where("(is_broadcast = ? OR target IN ? OR target IS NULL OR target = '')",
			false, []string{"chat", "both"})
	}
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&messages).Error

	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// GetRecentMessages 获取最近的消息
func (r *ChatRepository) GetRecentMessages(limit int) ([]model.ChatMessage, error) {
	var messages []model.ChatMessage

	err := db.DB.Where("status = ? AND (is_broadcast = ? OR target IN ? OR target IS NULL OR target = '')",
		1, false, []string{"chat", "both"}).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error

	if err != nil {
		return nil, err
	}

	// 反转数组，使其按时间正序排列
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// Delete 删除消息（软删除）
func (r *ChatRepository) Delete(id uint) error {
	return db.DB.Model(&model.ChatMessage{}).Where("id = ?", id).Update("status", 0).Error
}

// GetByID 根据ID获取消息
func (r *ChatRepository) GetByID(id uint) (*model.ChatMessage, error) {
	var message model.ChatMessage
	err := db.DB.Where("id = ? AND status = ?", id, 1).First(&message).Error
	return &message, err
}

// GetBroadcasts 获取系统广播（公告）
func (r *ChatRepository) GetBroadcasts(limit int) ([]model.ChatMessage, error) {
	var messages []model.ChatMessage

	err := db.DB.Where("status = ? AND is_broadcast = ? AND (target = ? OR target = ? OR target IS NULL OR target = '')",
		1, true, "announcement", "both").
		Order("priority DESC, created_at DESC").
		Limit(limit).
		Find(&messages).Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}

// GetBroadcastByID 根据ID获取广播详情
func (r *ChatRepository) GetBroadcastByID(id uint) (*model.ChatMessage, error) {
	var message model.ChatMessage
	err := db.DB.
		Where("id = ? AND status = ? AND is_broadcast = ? AND (target = ? OR target = ? OR target IS NULL OR target = '')",
			id, 1, true, "announcement", "both").
		First(&message).Error
	return &message, err
}
