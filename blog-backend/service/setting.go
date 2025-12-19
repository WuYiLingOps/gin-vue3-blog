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
			Group:     "site",
			Key:       key,
			Value:     value,
			UpdatedAt: time.Now(),
		})
	}

	return s.repo.BatchUpsert(settings)
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

// GetUploadSettings 获取上传配置（仅返回存储类型开关）
func (s *SettingService) GetUploadSettings() (map[string]string, error) {
	settings, err := s.repo.GetByGroup("upload")
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, setting := range settings {
		result[setting.Key] = setting.Value
	}

	// 设置默认值（兼容老数据）
	if result["storage_type"] == "" {
		result["storage_type"] = "local"
	}

	return result, nil
}

// UpdateUploadSettings 更新上传配置（仅更新存储类型开关）
func (s *SettingService) UpdateUploadSettings(data map[string]string) error {
	var settings []model.Setting

	// 只保存 storage_type，其他 OSS 配置从配置文件读取
	if storageType, ok := data["storage_type"]; ok {
		settings = append(settings, model.Setting{
			Group:     "upload",
			Key:       "storage_type",
			Value:     storageType,
			UpdatedAt: time.Now(),
		})
	}

	return s.repo.BatchUpdate(settings)
}

// GetFriendLinkInfo 获取我的友链信息
func (s *SettingService) GetFriendLinkInfo() (map[string]string, error) {
	settings, err := s.repo.GetByGroup("friendlink_info")
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, setting := range settings {
		result[setting.Key] = setting.Value
	}

	// 设置默认值
	if result["name"] == "" {
		result["name"] = "無以菱"
	}
	if result["desc"] == "" {
		result["desc"] = "分享技术与科技生活"
	}
	if result["url"] == "" {
		result["url"] = "https://xxxxx.cn/"
	}
	if result["avatar"] == "" {
		result["avatar"] = "https://pic.imgdb.cn/xxxx/xxxx.png"
	}
	if result["screenshot"] == "" {
		result["screenshot"] = "https://pic.imgdb.cn/xxxx/xxxx.png"
	}
	// RSS订阅可为空，不设置默认值

	return result, nil
}

// UpdateFriendLinkInfo 更新我的友链信息
func (s *SettingService) UpdateFriendLinkInfo(data map[string]string) error {
	var settings []model.Setting

	// 定义允许的字段
	allowedKeys := map[string]bool{
		"name":       true,
		"desc":       true,
		"url":        true,
		"avatar":     true,
		"screenshot": true,
		"rss":        true,
	}

	for key, value := range data {
		if allowedKeys[key] {
			settings = append(settings, model.Setting{
				Group:     "friendlink_info",
				Key:       key,
				Value:     value,
				Type:      "text",
				Label:     getFriendLinkInfoLabel(key),
				UpdatedAt: time.Now(),
			})
		}
	}

	return s.repo.BatchUpsert(settings)
}

// getFriendLinkInfoLabel 获取字段标签
func getFriendLinkInfoLabel(key string) string {
	labels := map[string]string{
		"name":       "名称",
		"desc":       "描述",
		"url":        "地址",
		"avatar":     "头像",
		"screenshot": "站点图片",
		"rss":        "订阅",
	}
	if label, ok := labels[key]; ok {
		return label
	}
	return key
}
