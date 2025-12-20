package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"blog-backend/config"
	"blog-backend/db"
)

type CalendarService struct{}

func NewCalendarService() *CalendarService {
	return &CalendarService{}
}

// CalendarResponse 贡献热力图响应结构
type CalendarResponse struct {
	Total         int64           `json:"total"`
	Contributions [][]CalendarDay `json:"contributions"`
}

// CalendarDay 单日贡献数据
type CalendarDay struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// GetGiteeCalendar 获取Gitee贡献热力图数据（带Redis缓存）
func (s *CalendarService) GetGiteeCalendar(username string) (*CalendarResponse, error) {
	if username == "" {
		return nil, fmt.Errorf("用户名不能为空")
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("gitee_calendar:%s", username)

	// 1. 先尝试从Redis缓存获取
	cachedData, err := db.RDB.Get(ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		// 缓存命中，解析并返回
		var response CalendarResponse
		if err := json.Unmarshal([]byte(cachedData), &response); err == nil {
			return &response, nil
		}
		// 如果解析失败，继续从API获取
	}

	// 2. 缓存未命中或解析失败，从gitee-calendar-api获取
	// 优先使用配置的API地址，如果没有配置则使用默认地址
	baseURL := config.Cfg.GiteeCalendar.APIURL
	if baseURL == "" {
		baseURL = "http://127.0.0.1:8081/api"
	}
	apiURL := fmt.Sprintf("%s?user=%s", baseURL, username)
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("请求gitee-calendar-api失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gitee-calendar-api返回错误状态: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var response CalendarResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 3. 将数据存入Redis缓存，过期时间20分钟
	cacheData, err := json.Marshal(response)
	if err == nil {
		db.RDB.Set(ctx, cacheKey, string(cacheData), 20*time.Minute)
	}

	return &response, nil
}
