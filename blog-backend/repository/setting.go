package repository

import (
	"blog-backend/db"
	"blog-backend/model"
)

type SettingRepository struct{}

func NewSettingRepository() *SettingRepository {
	return &SettingRepository{}
}

// GetByGroup 获取指定分组的所有配置
func (r *SettingRepository) GetByGroup(group string) ([]model.Setting, error) {
	var settings []model.Setting
	err := db.DB.Where("\"group\" = ?", group).Find(&settings).Error
	return settings, err
}

// GetByKey 根据 key 获取配置
func (r *SettingRepository) GetByKey(key string) (*model.Setting, error) {
	var setting model.Setting
	err := db.DB.Where("\"key\" = ?", key).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// Update 更新配置
func (r *SettingRepository) Update(setting *model.Setting) error {
	return db.DB.Model(setting).Updates(setting).Error
}

// BatchUpdate 批量更新配置
func (r *SettingRepository) BatchUpdate(settings []model.Setting) error {
	tx := db.DB.Begin()
	for _, setting := range settings {
		if err := tx.Model(&model.Setting{}).Where("\"key\" = ?", setting.Key).
			Updates(map[string]interface{}{
				"value":      setting.Value,
				"updated_at": setting.UpdatedAt,
			}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

