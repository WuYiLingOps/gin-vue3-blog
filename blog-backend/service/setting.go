package service

import (
	"blog-backend/model"
	"blog-backend/repository"
	"time"
)

type SettingService struct {
	repo *repository.SettingRepository
}

func NewSettingService() *SettingService {
	return &SettingService{
		repo: repository.NewSettingRepository(),
	}
}

// GetSiteSettings 获取网站配置
func (s *SettingService) GetSiteSettings() (map[string]string, error) {
	settings, err := s.repo.GetByGroup("site")
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, setting := range settings {
		result[setting.Key] = setting.Value
	}

	return result, nil
}

// UpdateSiteSettings 更新网站配置
func (s *SettingService) UpdateSiteSettings(data map[string]string) error {
	var settings []model.Setting

	for key, value := range data {
		settings = append(settings, model.Setting{
			Key:       key,
			Value:     value,
			UpdatedAt: time.Now(),
		})
	}

	return s.repo.BatchUpdate(settings)
}

// GetPublicSettings 获取公开的网站配置（前端用）
func (s *SettingService) GetPublicSettings() (map[string]string, error) {
	settings, err := s.repo.GetByGroup("site")
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, setting := range settings {
		result[setting.Key] = setting.Value
	}

	return result, nil
}

