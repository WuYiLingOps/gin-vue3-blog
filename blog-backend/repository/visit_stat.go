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

// IncrementTodayViewCount 增加今天的浏览量
func (r *VisitStatRepository) IncrementTodayViewCount() error {
	today := time.Now().Truncate(24 * time.Hour)

	// 查找今天的记录
	var stat model.VisitStat
	err := db.DB.Where("date = ?", today).First(&stat).Error

	if err != nil {
		// 如果记录不存在，创建新记录
		stat = model.VisitStat{
			Date:      today,
			ViewCount: 1,
		}
		return db.DB.Create(&stat).Error
	}

	// 如果记录存在，增加计数
	return db.DB.Model(&stat).UpdateColumn("view_count", db.DB.Raw("view_count + 1")).Error
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

