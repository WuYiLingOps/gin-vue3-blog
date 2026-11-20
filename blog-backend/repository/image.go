/*
 * 项目名称：blog-backend
 * 文件名称：image.go
 * 创建时间：2026-11-20 10:00:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：图片管理数据访问层，提供全站图片 URL 的扫描与批量替换功能
 */
package repository

import (
	"fmt"
	"strings"

	"blog-backend/db"
	"blog-backend/model"

	"gorm.io/gorm"
)

// imageFieldDef 含图片 URL 的字段定义（白名单，用于扫描和替换，防止动态 SQL 注入）
type imageFieldDef struct {
	Table  string // 表名
	Column string // 字段列名
	Kind   string // url-单URL字段 / markdown-Markdown长文本 / json_array-JSON数组 / settings-KV配置表
}

// imageFieldDefs 图片字段清单（覆盖全站所有存储图片 URL 的位置）
var imageFieldDefs = []imageFieldDef{
	{Table: "users", Column: "avatar", Kind: "url"},
	{Table: "posts", Column: "cover", Kind: "url"},
	{Table: "posts", Column: "content", Kind: "markdown"},
	{Table: "post_revisions", Column: "cover", Kind: "url"},
	{Table: "post_revisions", Column: "content", Kind: "markdown"},
	{Table: "moments", Column: "images", Kind: "json_array"},
	{Table: "albums", Column: "image_url", Kind: "url"},
	{Table: "friend_links", Column: "icon", Kind: "url"},
	{Table: "friend_links", Column: "screenshot", Kind: "url"},
	{Table: "chat_messages", Column: "avatar", Kind: "url"},
	{Table: "settings", Column: "value", Kind: "settings"},
}

// ImageRecord 一条含图片字段的数据库记录
type ImageRecord struct {
	Table string // 表名
	Field string // 字段名（settings 表为配置 key）
	Kind  string // 字段类型（同 imageFieldDef.Kind）
	RowID int64  // 行 ID
	Value string // 字段值原文
}

// ImageRepository 图片管理数据访问层结构体
type ImageRepository struct{}

// NewImageRepository 创建图片管理数据访问层实例
func NewImageRepository() *ImageRepository {
	return &ImageRepository{}
}

// LoadImageRecords 加载所有含图片 URL 的记录原文（使用全局连接）
func (r *ImageRepository) LoadImageRecords() ([]ImageRecord, error) {
	return r.loadImageRecords(db.DB)
}

// loadImageRecords 指定 GORM 会话加载记录（事务内复用，保证替换时读写同一事务）
func (r *ImageRepository) loadImageRecords(gormDB *gorm.DB) ([]ImageRecord, error) {
	var records []ImageRecord

	for _, def := range imageFieldDefs {
		// settings 表：按 key-value 加载
		if def.Kind == "settings" {
			var settings []model.Setting
			if err := gormDB.Find(&settings).Error; err != nil {
				return nil, err
			}
			for _, s := range settings {
				records = append(records, ImageRecord{
					Table: "settings",
					Field: s.Key,
					Kind:  def.Kind,
					RowID: int64(s.ID),
					Value: s.Value,
				})
			}
			continue
		}

		// 常规表：COALESCE 处理可空字段
		var rows []struct {
			ID    int64  `gorm:"column:id"`
			Value string `gorm:"column:value"`
		}
		sql := fmt.Sprintf("SELECT id, COALESCE(%s, '') AS value FROM %s", def.Column, def.Table)
		if err := gormDB.Raw(sql).Scan(&rows).Error; err != nil {
			return nil, err
		}
		for _, row := range rows {
			records = append(records, ImageRecord{
				Table: def.Table,
				Field: def.Column,
				Kind:  def.Kind,
				RowID: row.ID,
				Value: row.Value,
			})
		}
	}

	return records, nil
}

// ReplaceImageURLs 在单个事务中把所有记录里的旧 URL 替换为新 URL
// 替换按旧 URL 长度降序执行，避免前缀 URL 互相误伤
// 返回更新的记录行数
func (r *ImageRepository) ReplaceImageURLs(mappings map[string]string) (int64, error) {
	olds := make([]string, 0, len(mappings))
	for old := range mappings {
		olds = append(olds, old)
	}
	// 按 URL 长度降序，长 URL 优先替换
	for i := 0; i < len(olds); i++ {
		for j := i + 1; j < len(olds); j++ {
			if len(olds[j]) > len(olds[i]) {
				olds[i], olds[j] = olds[j], olds[i]
			}
		}
	}

	var updated int64
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 事务内读取，保证替换时读写同一事务快照
		records, err := r.loadImageRecords(tx)
		if err != nil {
			return err
		}

		for _, rec := range records {
			newVal := rec.Value
			matched := false
			for _, old := range olds {
				if strings.Contains(newVal, old) {
					newVal = strings.ReplaceAll(newVal, old, mappings[old])
					matched = true
				}
			}
			if !matched {
				continue
			}

			sql := fmt.Sprintf("UPDATE %s SET %s = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", rec.Table, rec.Field)
			// settings 表的字段名是 key，实际更新列为 value
			if rec.Kind == "settings" {
				sql = "UPDATE settings SET value = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"
			}
			if err := tx.Exec(sql, newVal, rec.RowID).Error; err != nil {
				return err
			}
			updated++
		}
		return nil
	})

	return updated, err
}
