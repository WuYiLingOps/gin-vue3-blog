/*
 * 项目名称：blog-backend
 * 文件名称：operation_log.go
 * 创建时间：2026-02-06 22:00:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：操作日志数据模型定义
 */
package model

import (
	"time"
)

// OperationLog 操作日志模型
// 功能说明：记录管理员和超级管理员的所有操作日志，包括文章、分类、标签、用户等管理操作
type OperationLog struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"index;not null"`
	Username    string    `json:"username" gorm:"not null;size:50"`
	Action      string    `json:"action" gorm:"not null;size:50;index"` // create, update, delete
	Module      string    `json:"module" gorm:"not null;size:50;index"` // post, category, tag, user, comment等
	TargetType  string    `json:"target_type" gorm:"size:50;index"`     // 目标类型（与module相同）
	TargetID    *uint     `json:"target_id" gorm:"index"`               // 目标ID（如文章ID、分类ID等）
	TargetName  string    `json:"target_name" gorm:"size:255"`          // 目标名称（如文章标题、分类名称等）
	Description string    `json:"description" gorm:"type:text"`         // 操作描述
	IP          string    `json:"ip" gorm:"size:45"`                    // 操作IP地址
	UserAgent   string    `json:"user_agent" gorm:"type:text"`          // 用户代理
	CreatedAt   time.Time `json:"created_at" gorm:"index"`

	// 关联关系
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
