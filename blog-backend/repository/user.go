/*
 * 项目名称：blog-backend
 * 文件名称：user.go
 * 创建时间：2026-01-31 16:29:06
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：用户数据访问层，提供用户信息的数据库操作功能
 */
package repository

import (
	"blog-backend/db"
	"blog-backend/model"
)

// UserRepository 用户数据访问层结构体
type UserRepository struct{}

// NewUserRepository 创建用户数据访问层实例
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return db.DB.Create(user).Error
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := db.DB.First(&user, id).Error
	return &user, err
}

// GetByUsername 根据用户名获取用户
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := db.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
	return db.DB.Save(user).Error
}

// Delete 删除用户
func (r *UserRepository) Delete(id uint) error {
	return db.DB.Delete(&model.User{}, id).Error
}

// List 获取用户列表
func (r *UserRepository) List(page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	offset := (page - 1) * pageSize

	if err := db.DB.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.DB.Offset(offset).Limit(pageSize).Find(&users).Error
	return users, total, err
}

// UpdateStatus 更新用户状态
func (r *UserRepository) UpdateStatus(id uint, status int) error {
	return db.DB.Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

// GetAdmins 获取所有管理员用户
func (r *UserRepository) GetAdmins() ([]model.User, error) {
	var admins []model.User
	err := db.DB.Where("role = ? AND status = ?", "admin", 1).Find(&admins).Error
	return admins, err
}
