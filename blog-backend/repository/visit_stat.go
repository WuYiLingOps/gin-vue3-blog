package repository

import (
	"time"

	"blog-backend/db"
	"blog-backend/model"
)

type VisitStatRepository struct{}

func NewVisitStatRepository() *VisitStatRepository {
	return &VisitStatRepository{}
}

// RecordVisit 记录访问（基于 IP 的 UV 统计）
func (r *VisitStatRepository) RecordVisit(ip string) error {
	today := time.Now().Truncate(24 * time.Hour)

	// 先尝试插入访问记录，如果该 IP 今天已经访问过，则会因为唯一索引而失败
	visitRecord := model.VisitRecord{
		Date:      today,
		IP:        ip,
		CreatedAt: time.Now(),
	}

	// 使用事务确保数据一致性
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 尝试插入访问记录
	result := tx.Create(&visitRecord)
	
	// 如果插入成功，说明这是该 IP 今天第一次访问，需要增加 UV 计数
	if result.Error == nil {
		// 查找今天的统计记录
		var stat model.VisitStat
		err := tx.Where("date = ?", today).First(&stat).Error

		if err != nil {
			// 如果记录不存在，创建新记录
			stat = model.VisitStat{
				Date:      today,
				ViewCount: 1,
			}
			if err := tx.Create(&stat).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// 如果记录存在，增加计数
			if err := tx.Model(&stat).UpdateColumn("view_count", db.DB.Raw("view_count + 1")).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	// 如果插入失败（唯一索引冲突），说明该 IP 今天已经访问过，不增加 UV 计数，忽略错误

	return tx.Commit().Error
}

// GetLast7DaysStats 获取最近7天的访问量统计
func (r *VisitStatRepository) GetLast7DaysStats() ([]model.VisitStat, error) {
	var stats []model.VisitStat
	
	// 获取7天前的日期（截断到天）
	sevenDaysAgo := time.Now().AddDate(0, 0, -6).Truncate(24 * time.Hour)

	err := db.DB.Where("date >= ?", sevenDaysAgo).
		Order("date ASC").
		Find(&stats).Error

	return stats, err
}

// EnsureLast7DaysRecords 确保最近7天都有记录（没有访问量的日期记录为0）
func (r *VisitStatRepository) EnsureLast7DaysRecords() error {
	now := time.Now()
	
	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i).Truncate(24 * time.Hour)
		
		var count int64
		db.DB.Model(&model.VisitStat{}).Where("date = ?", date).Count(&count)
		
		// 如果不存在，创建记录
		if count == 0 {
			stat := model.VisitStat{
				Date:      date,
				ViewCount: 0,
			}
			if err := db.DB.Create(&stat).Error; err != nil {
				return err
			}
		}
	}
	
	return nil
}

