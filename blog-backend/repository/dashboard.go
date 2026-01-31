/*
 * 项目名称：blog-backend
 * 文件名称：dashboard.go
 * 创建时间：2026-01-31 16:29:06
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：仪表盘数据访问层，提供统计数据查询功能
 */
package repository

import (
	"blog-backend/db"
	"blog-backend/model"
)

// GetPublishedCount 获取已发布文章总数
func (r *PostRepository) GetPublishedCount() (int64, error) {
	var count int64
	err := db.DB.Model(&model.Post{}).Where("status = 1").Count(&count).Error
	return count, err
}

// GetTotalViews 获取总浏览量
func (r *PostRepository) GetTotalViews() (int64, error) {
	var total int64
	err := db.DB.Model(&model.Post{}).Where("status = 1").Select("COALESCE(SUM(view_count), 0)").Scan(&total).Error
	return total, err
}

// GetTotalCount 获取用户总数
func (r *UserRepository) GetTotalCount() (int64, error) {
	var count int64
	err := db.DB.Model(&model.User{}).Count(&count).Error
	return count, err
}

// GetTotalCount 获取评论总数
func (r *CommentRepository) GetTotalCount() (int64, error) {
	var count int64
	err := db.DB.Model(&model.Comment{}).Where("status = 1").Count(&count).Error
	return count, err
}
