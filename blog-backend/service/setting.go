package service

import (
	"errors"

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

// GetNotificationSettings 获取通知配置
func (s *SettingService) GetNotificationSettings() (map[string]string, error) {
	settings, err := s.repo.GetByGroup("notification")
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, setting := range settings {
		result[setting.Key] = setting.Value
	}

	// 设置默认值
	if result["notify_admin_on_comment"] == "" {
		result["notify_admin_on_comment"] = "0"
	}

	return result, nil
}

// UpdateNotificationSettings 更新通知配置
func (s *SettingService) UpdateNotificationSettings(data map[string]string) error {
	var settings []model.Setting

	// 只允许修改 notify_admin_on_comment
	if notifyAdmin, ok := data["notify_admin_on_comment"]; ok {
		settings = append(settings, model.Setting{
			Group:     "notification",
			Key:       "notify_admin_on_comment",
			Value:     notifyAdmin,
			Type:      "text",
			Label:     "评论时通知管理员",
			UpdatedAt: time.Now(),
		})
	}

	return s.repo.BatchUpsert(settings)
}

// GetRegisterSettings 获取注册配置
func (s *SettingService) GetRegisterSettings() (map[string]string, error) {
	setting, err := s.repo.GetByKey("disable_register")
	if err != nil {
		// 如果配置不存在，返回默认值（允许注册）
		return map[string]string{
			"disable_register": "0",
		}, nil
	}

	return map[string]string{
		"disable_register": setting.Value,
	}, nil
}

// UpdateRegisterSettings 更新注册配置
func (s *SettingService) UpdateRegisterSettings(data map[string]string) error {
	var settings []model.Setting

	// 只允许修改 disable_register
	if disableRegister, ok := data["disable_register"]; ok {
		// 验证值只能是 "0" 或 "1"
		if disableRegister != "0" && disableRegister != "1" {
			return errors.New("disable_register 值只能是 0 或 1")
		}

		settings = append(settings, model.Setting{
			Group:     "site",
			Key:       "disable_register",
			Value:     disableRegister,
			Type:      "text",
			Label:     "禁用用户注册",
			UpdatedAt: time.Now(),
		})
	}

	return s.repo.BatchUpsert(settings)
}

// GetAboutInfo 获取关于我信息
func (s *SettingService) GetAboutInfo() (string, error) {
	setting, err := s.repo.GetByKey("about_content")
	if err != nil {
		// 如果配置不存在，返回空字符串
		return "", nil
	}
	return setting.Value, nil
}

// UpdateAboutInfo 更新关于我信息
func (s *SettingService) UpdateAboutInfo(content string) error {
	setting := model.Setting{
		Group:     "about",
		Key:       "about_content",
		Value:     content,
		Type:      "text",
		Label:     "关于我内容",
		UpdatedAt: time.Now(),
	}
	return s.repo.BatchUpsert([]model.Setting{setting})
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
