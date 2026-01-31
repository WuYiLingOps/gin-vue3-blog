/*
 * 项目名称：blog-backend
 * 文件名称：password_reset.go
 * 创建时间：2026-01-31 16:29:06
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：密码重置令牌数据访问层，提供密码重置和注册验证令牌的数据库操作功能
 */
package repository

import (
	"time"

	"blog-backend/db"
	"blog-backend/model"
)

// PasswordResetRepository 密码重置令牌数据访问层结构体
type PasswordResetRepository struct{}

// NewPasswordResetRepository 创建密码重置令牌数据访问层实例
func NewPasswordResetRepository() *PasswordResetRepository {
	return &PasswordResetRepository{}
}

// Create 创建密码重置令牌
func (r *PasswordResetRepository) Create(token *model.PasswordResetToken) error {
	return db.DB.Create(token).Error
}

// GetValidToken 获取有效的重置令牌
func (r *PasswordResetRepository) GetValidToken(email, code string) (*model.PasswordResetToken, error) {
	var token model.PasswordResetToken
	err := db.DB.Where("email = ? AND code = ? AND is_used = ? AND expire_at > ?",
		email, code, false, time.Now()).
		Order("created_at DESC").
		First(&token).Error
	return &token, err
}

// Update 更新令牌
func (r *PasswordResetRepository) Update(token *model.PasswordResetToken) error {
	return db.DB.Save(token).Error
}

// DeleteExpired 删除过期的令牌（定期清理）
func (r *PasswordResetRepository) DeleteExpired() error {
	return db.DB.Where("expire_at < ?", time.Now()).Delete(&model.PasswordResetToken{}).Error
}

// GetRecentByEmail 获取指定邮箱最近的令牌（用于防刷）
func (r *PasswordResetRepository) GetRecentByEmail(email string, duration time.Duration) (*model.PasswordResetToken, error) {
	var token model.PasswordResetToken
	err := db.DB.Where("email = ? AND created_at > ?",
		email, time.Now().Add(-duration)).
		Order("created_at DESC").
		First(&token).Error
	return &token, err
}
