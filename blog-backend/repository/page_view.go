/*
 * 项目名称：blog-backend
 * 文件名称：page_view.go
 * 创建时间：2026-06-30
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：全站页面访问记录数据访问层，提供页面访问记录和统计的数据库操作功能
 */
package repository

import (
	"blog-backend/db"
	"blog-backend/model"
	"time"
)

// PageViewRepository 全站页面访问记录数据访问层结构体
type PageViewRepository struct{}

// NewPageViewRepository 创建全站页面访问记录数据访问层实例
func NewPageViewRepository() *PageViewRepository {
	return &PageViewRepository{}
}

// HasViewedToday 检查今天是否已经访问过该路径
// 按天去重：同IP对同路径每天只计1次
func (r *PageViewRepository) HasViewedToday(path string, ip string) (bool, error) {
	var count int64
	today := time.Now().Format("2006-01-02")

	err := db.DB.Model(&model.PageView{}).
		Where("path = ? AND ip = ? AND DATE(created_at) = ?", path, ip, today).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// RecordView 记录页面访问
func (r *PageViewRepository) RecordView(path string, userID *uint, ip string, userAgent string) error {
	pageView := model.PageView{
		Path:      path,
		UserID:    userID,
		IP:        ip,
		UserAgent: userAgent,
	}

	return db.DB.Create(&pageView).Error
}

// VisitStat 按天统计访问量结果
type PageVisitStat struct {
	Date  time.Time
	Count int64
}

// GetVisitStats 获取指定时间范围内按天聚合的全站访问量统计
func (r *PageViewRepository) GetVisitStats(start, end time.Time) ([]PageVisitStat, error) {
	var results []PageVisitStat

	err := db.DB.Model(&model.PageView{}).
		Select("DATE(created_at) AS date, COUNT(*) AS count").
		Where("created_at >= ? AND created_at < ?", start, end).
		Group("DATE(created_at)").
		Order("DATE(created_at)").
		Scan(&results).Error

	return results, err
}

// GetTotalPV 获取总PV数
func (r *PageViewRepository) GetTotalPV() (int64, error) {
	var count int64
	err := db.DB.Model(&model.PageView{}).Count(&count).Error
	return count, err
}
