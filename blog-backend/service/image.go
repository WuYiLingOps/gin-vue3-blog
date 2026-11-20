/*
 * 项目名称：blog-backend
 * 文件名称：image.go
 * 创建时间：2026-11-20 10:00:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：图片管理业务逻辑层，提供全站图片盘点和图床间迁移功能
 */
package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"blog-backend/config"
	"blog-backend/model"
	"blog-backend/repository"
	"blog-backend/util"

	"gorm.io/datatypes"
	"go.uber.org/zap"
)

// ImageUsageLocation 图片在某处的使用明细
type ImageUsageLocation struct {
	Table       string `json:"table"`        // 所在表名
	Field       string `json:"field"`        // 所在字段（settings 表为配置 key）
	TableLabel  string `json:"table_label"`  // 位置中文说明（如：文章、用户头像）
	RowID       int64  `json:"row_id"`       // 所在记录 ID（可用于跳转定位）
	Kind        string `json:"kind"`         // 字段类型：url/markdown/json_array/settings
}

// ImageAsset 盘点结果中的单个图片资产
type ImageAsset struct {
	URL          string                `json:"url"`             // 图片 URL（库内原样）
	StorageType  string                `json:"storage_type"`    // 来源存储：local/oss/cos/easyimage/external
	Locations    []ImageUsageLocation  `json:"locations"`       // 使用位置明细
	RefCount     int                   `json:"ref_count"`       // 引用次数
	Deletable    bool                  `json:"deletable"`       // 源文件是否支持 API 删除
}

// MigrateImageRequest 发起图片迁移请求
type MigrateImageRequest struct {
	URLs         []string `json:"urls" binding:"required,min=1"`     // 要迁移的图片 URL 列表
	TargetStorag string   `json:"target_storage" binding:"required"` // 目标存储类型：local/oss/cos/easyimage
	DeleteSource bool     `json:"delete_source"`                     // 迁移成功后是否删除源文件
	OperatorID   *uint    `json:"-"`                                 // 操作用户ID（handler 注入，用于迁移日志）
	OperatorName string   `json:"-"`                                 // 操作用户名（handler 注入）
}

// MigrateFailedItem 迁移失败条目
type MigrateFailedItem struct {
	URL   string `json:"url"`
	Error string `json:"error"`
}

// MigrateMappingItem 迁移映射条目
type MigrateMappingItem struct {
	OldURL string `json:"old_url"`
	NewURL string `json:"new_url"`
}

// MigrateStatus 迁移任务状态（供前端轮询）
type MigrateStatus struct {
	Running     bool                `json:"running"`            // 是否执行中
	Total       int                 `json:"total"`              // 总数
	Migrated    int                 `json:"migrated"`           // 已完成数
	Failed      []MigrateFailedItem `json:"failed"`             // 失败列表
	Success     []MigrateMappingItem `json:"success"`           // 成功映射列表
	Finished    bool                `json:"finished"`           // 本次任务是否已结束（用于区分从未运行和已结束）
	Message     string              `json:"message"`            // 结束时的汇总信息
	CurrentURL  string              `json:"current_url"`        // 正在处理的图片 URL（running 时可见）
	CurrentStep string              `json:"current_step"`       // 当前步骤描述（下载中/上传中/替换数据库/删除源文件）
	ElapsedSec  int64               `json:"elapsed_sec"`        // 任务已进行秒数
}

// migrateTask 迁移任务（内存单例）
type migrateTask struct {
	running     bool
	total       int
	migrated    int
	failed      []MigrateFailedItem
	success     []MigrateMappingItem
	finished    bool
	message     string
	currentURL  string // 正在处理的 URL
	currentStep string // 当前步骤描述
	startedAt   time.Time
	mu          sync.Mutex
}

// ImageService 图片管理业务逻辑层结构体
type ImageService struct {
	repo    *repository.ImageRepository
	logRepo *repository.ImageMigrateLogRepository
}

// 全局迁移任务状态（进程内单例，同一时间只允许一个迁移任务）
var (
	globalMigrateTask = &migrateTask{}
	migrateTaskMu     sync.Mutex
)

// NewImageService 创建图片管理业务逻辑层实例
func NewImageService() *ImageService {
	return &ImageService{
		repo:    repository.NewImageRepository(),
		logRepo: repository.NewImageMigrateLogRepository(),
	}
}

// storageTypeLabel 存储类型的中文标签
var storageTypeLabel = map[string]string{
	"local":     "本地存储",
	"oss":       "阿里云OSS",
	"cos":       "腾讯云COS",
	"easyimage": "EasyImage",
	"external":  "外链图片",
}

// tableFieldLabel 图片字段位置的中文说明（表 + 字段组合），未命中时回退为"表名.字段名"
var tableFieldLabel = map[string]string{
	"users.avatar":                "用户头像",
	"posts.cover":                 "文章封面",
	"posts.content":               "文章正文插图",
	"post_revisions.cover":        "文章修订封面",
	"post_revisions.content":      "文章修订正文插图",
	"moments.images":              "说说配图",
	"albums.image_url":            "相册照片",
	"friend_links.icon":           "友链图标",
	"friend_links.screenshot":     "友链截图",
	"chat_messages.avatar":        "聊天室头像",
	"settings.cover_bg_images":    "首页封面背景图",
	"settings.friendlink_info":    "我的友链卡片图",
	"settings.avatar":             "友链资料头像",
	"settings.screenshot":         "友链资料截图",
	"settings.reward_wechat":      "微信打赏码",
	"settings.reward_alipay":      "支付宝打赏码",
	"settings.about_skills":       "关于页-技能卡",
	"settings.about_intro":        "关于页-作者介绍",
	"settings.about_self_info":    "关于页-个人信息",
	"settings.about_personality":  "关于页-MBTI性格卡",
	"settings.about_careers":      "关于页-职业生涯",
	"settings.about_maxim":        "关于页-座右铭",
	"settings.about_map":          "关于页-地理位置",
	"settings.about_site_tips":    "关于页-打字词轮播",
}

// locationLabel 生成使用位置的中文说明
func locationLabel(table, field string) string {
	if label, ok := tableFieldLabel[table+"."+field]; ok {
		return label
	}
	if table == "settings" {
		return "系统配置-" + field
	}
	return table + "." + field
}

// ScanImageAssets 全站图片盘点：扫描所有含图片 URL 的记录，按 URL 去重聚合
// 盘点全集缓存到 Redis（TTL 10 分钟），命中时不再查询数据库；refresh=true 强制刷新缓存
// storage/keyword 为兼容参数：传入时跳过缓存走实时扫描并按条件筛选（前端已改为本地过滤，正常不再传）
func (s *ImageService) ScanImageAssets(storage, keyword string, refresh bool) ([]ImageAsset, error) {
	// 无筛选条件且不强制刷新时优先读缓存
	if storage == "" && keyword == "" && !refresh {
		if assets, ok := loadScanCache(); ok {
			return assets, nil
		}
	}

	records, err := s.repo.LoadImageRecords()
	if err != nil {
		return nil, err
	}

	assetMap := make(map[string]*ImageAsset)
	var order []string

	for _, rec := range records {
		var urls []string
		switch rec.Kind {
		case "url":
			if rec.Value != "" {
				urls = []string{rec.Value}
			}
		case "markdown":
			urls = util.ExtractImageURLs(rec.Value)
		case "json_array":
			urls = extractJSONArrayStrings(rec.Value)
		case "settings":
			// settings 的 value 形态多样：markdown 图片语法、JSON 数组、JSON 对象内嵌图片字段、纯 URL
			// 仅收录明确的图片类 URL，避免把社交链接（social_gitee）、站点地址（site_url）等非图片配置误判为图片
			urls = util.ExtractImageURLs(rec.Value)
			for _, u := range extractJSONArrayStrings(rec.Value) {
				if looksLikeImageURL(u) {
					urls = append(urls, u)
				}
			}
			for _, u := range extractJSONImageFields(rec.Value) {
				if looksLikeImageURL(u) {
					urls = append(urls, u)
				}
			}
			if u := extractPlainURL(rec.Value); u != "" && looksLikeImageURL(u) {
				urls = append(urls, u)
			}
		}

		for _, rawURL := range urls {
			u := strings.TrimSpace(rawURL)
			if u == "" {
				continue
			}

			asset, exists := assetMap[u]
			if !exists {
				asset = &ImageAsset{
					URL:         u,
					StorageType: string(util.ClassifyImageURL(u)),
					// 可删源 = 当前运行配置可管理的存储（结构性识别出的历史 bucket 无法 API 删除）
					Deletable: util.IsManagedStorageURL(u),
				}
				assetMap[u] = asset
				order = append(order, u)
			}
			asset.RefCount++
			loc := ImageUsageLocation{
				Table:      rec.Table,
				Field:      rec.Field,
				TableLabel: locationLabel(rec.Table, rec.Field),
				RowID:      rec.RowID,
				Kind:       rec.Kind,
			}
			// 同一处引用（表+字段+行）只记一条明细
			duplicated := false
			for _, existing := range asset.Locations {
				if existing.Table == loc.Table && existing.Field == loc.Field && existing.RowID == loc.RowID {
					duplicated = true
					break
				}
			}
			if !duplicated {
				asset.Locations = append(asset.Locations, loc)
			}
		}
	}

	// 组装结果（显式空切片，避免空结果序列化为 JSON null）
	result := make([]ImageAsset, 0)
	for _, u := range order {
		asset := assetMap[u]
		if storage != "" && asset.StorageType != storage {
			continue
		}
		if keyword != "" && !strings.Contains(asset.URL, keyword) {
			continue
		}
		result = append(result, *asset)
	}

	// 无筛选条件的全集结果写入缓存（带筛选条件的结果不缓存）
	if storage == "" && keyword == "" {
		saveScanCache(result)
	}

	return result, nil
}

// GetMigrateStatus 获取迁移任务状态
func (s *ImageService) GetMigrateStatus() MigrateStatus {
	// 先在全局锁下取当前任务快照，避免与 StartMigrate 的任务重置产生竞态
	migrateTaskMu.Lock()
	task := globalMigrateTask
	migrateTaskMu.Unlock()

	task.mu.Lock()
	defer task.mu.Unlock()

	status := MigrateStatus{
		Running:     task.running,
		Total:       task.total,
		Migrated:    task.migrated,
		Finished:    task.finished,
		Message:     task.message,
		CurrentURL:  task.currentURL,
		CurrentStep: task.currentStep,
		// 显式初始化为空切片：Go nil 切片会序列化为 JSON null，前端对 null 取 length 会抛错
		Failed:  []MigrateFailedItem{},
		Success: []MigrateMappingItem{},
	}
	if task.running && !task.startedAt.IsZero() {
		status.ElapsedSec = int64(time.Since(task.startedAt).Seconds())
	}
	for _, f := range task.failed {
		status.Failed = append(status.Failed, f)
	}
	for _, m := range task.success {
		status.Success = append(status.Success, m)
	}
	return status
}

// StartMigrate 启动图片迁移任务（同一时间只允许一个任务）
func (s *ImageService) StartMigrate(req *MigrateImageRequest) error {
	// 校验目标存储类型
	switch req.TargetStorag {
	case string(util.StorageLocal), string(util.StorageOSS), string(util.StorageCOS), string(util.StorageEasyImage):
	default:
		return errors.New("目标存储类型只能是 local、oss、cos 或 easyimage")
	}

	// 目标存储配置校验（避免启动后才失败）
	if err := validateStorageConfig(util.StorageType(req.TargetStorag)); err != nil {
		return err
	}

	// URL 去重
	urlSet := make(map[string]bool)
	var urls []string
	for _, u := range req.URLs {
		u = strings.TrimSpace(u)
		if u == "" || urlSet[u] {
			continue
		}
		urlSet[u] = true
		urls = append(urls, u)
	}
	if len(urls) == 0 {
		return errors.New("没有有效的图片 URL")
	}

	migrateTaskMu.Lock()
	defer migrateTaskMu.Unlock()

	if globalMigrateTask.running {
		return errors.New("已有迁移任务正在执行，请等待完成后再试")
	}

	// 重置任务状态
	globalMigrateTask = &migrateTask{
		running:   true,
		total:     len(urls),
		startedAt: time.Now(),
	}

	go s.runMigrate(urls, req)

	return nil
}

// migrateConcurrency 迁移任务的并发数
const migrateConcurrency = 5

// runMigrate 执行迁移任务（后台 goroutine，5 并发处理）
func (s *ImageService) runMigrate(urls []string, req *MigrateImageRequest) {
	defer func() {
		if r := recover(); r != nil {
			globalMigrateTask.mu.Lock()
			globalMigrateTask.running = false
			globalMigrateTask.finished = true
			globalMigrateTask.message = fmt.Sprintf("迁移任务异常终止: %v", r)
			globalMigrateTask.mu.Unlock()
		}
	}()

	target := util.StorageType(req.TargetStorag)
	taskStart := time.Now()

	type migrateResult struct {
		OldURL string
		NewURL string
		Err    error
	}

	var mu sync.Mutex
	results := make([]migrateResult, 0, len(urls))
	sem := make(chan struct{}, migrateConcurrency)

	var wg sync.WaitGroup
	for _, imgURL := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			newURL, err := s.migrateOneImage(u, target)

			mu.Lock()
			results = append(results, migrateResult{OldURL: u, NewURL: newURL, Err: err})
			mu.Unlock()

			// 实时更新进度（无论成败都推进计数）
			globalMigrateTask.mu.Lock()
			globalMigrateTask.migrated++
			if err != nil {
				globalMigrateTask.failed = append(globalMigrateTask.failed, MigrateFailedItem{URL: u, Error: err.Error()})
			} else {
				globalMigrateTask.success = append(globalMigrateTask.success, MigrateMappingItem{OldURL: u, NewURL: newURL})
			}
			globalMigrateTask.mu.Unlock()
		}(imgURL)
	}
	wg.Wait()

	// 汇总成功映射与失败列表
	mappings := make(map[string]string)
	var failed []MigrateFailedItem
	for _, r := range results {
		if r.Err != nil {
			failed = append(failed, MigrateFailedItem{URL: r.OldURL, Error: r.Err.Error()})
		} else {
			mappings[r.OldURL] = r.NewURL
		}
	}

	// 图片传输阶段结束，进入收尾步骤
	setCurrentStep("", "替换数据库")

	// 汇总结果
	message := fmt.Sprintf("迁移完成：成功 %d 张，失败 %d 张", len(mappings), len(failed))
	updated := 0
	deleted := 0
	if len(mappings) > 0 {
		// 数据库替换（单事务）
		replaced, err := s.repo.ReplaceImageURLs(mappings)
		if err != nil {
			message = fmt.Sprintf("数据库替换失败（图片已上传到目标图床，URL 映射见成功列表）: %v", err)
			if config.Cfg != nil {
				zap.L().Error("图片迁移数据库替换失败", zap.Error(err))
			}
		} else {
			updated = int(replaced)
			message = fmt.Sprintf("迁移完成：成功 %d 张（更新 %d 处数据库记录），失败 %d 张", len(mappings), updated, len(failed))
			// 数据库已变更，失效盘点缓存，下次扫描重新查库
			InvalidateScanCache()
		}

		// 删除源文件（EasyImage 与外链无法删除，自动跳过）
		if req.DeleteSource {
			// 持锁拷贝成功列表后再遍历删除
			globalMigrateTask.mu.Lock()
			successList := make([]MigrateMappingItem, len(globalMigrateTask.success))
			copy(successList, globalMigrateTask.success)
			globalMigrateTask.mu.Unlock()

			setCurrentStep("", "删除源文件")
			for _, m := range successList {
				// 非当前配置可管理的存储（EasyImage/外链/历史 bucket）无法 API 删除，跳过
				if !util.IsManagedStorageURL(m.OldURL) {
					continue
				}
				if err := util.DeleteFileByURL(m.OldURL); err != nil {
					if config.Cfg != nil {
						zap.L().Warn("迁移后删除源文件失败", zap.String("url", m.OldURL), zap.Error(err))
					}
					continue
				}
				deleted++
			}
			if deleted > 0 {
				message += fmt.Sprintf("，已删除源文件 %d 个", deleted)
			}
		}
	}

	globalMigrateTask.mu.Lock()
	globalMigrateTask.running = false
	globalMigrateTask.finished = true
	globalMigrateTask.failed = failed
	globalMigrateTask.message = message
	globalMigrateTask.currentURL = ""
	globalMigrateTask.currentStep = ""
	globalMigrateTask.mu.Unlock()

	if config.Cfg != nil {
		zap.L().Info("图片迁移任务结束", zap.String("result", message))
	}

	// 写入持久化迁移日志（服务重启后仍可查历史；映射按迁移完成顺序存储为数组）
	globalMigrateTask.mu.Lock()
	successSnapshot := make([]MigrateMappingItem, len(globalMigrateTask.success))
	copy(successSnapshot, globalMigrateTask.success)
	globalMigrateTask.mu.Unlock()
	s.saveMigrateLog(req, successSnapshot, failed, updated, deleted, message, time.Since(taskStart), req.OperatorID, req.OperatorName)
}

// migrateOneImage 迁移单张图片：下载 → 上传到目标存储（实时上报当前步骤）
func (s *ImageService) migrateOneImage(imgURL string, target util.StorageType) (string, error) {
	setCurrentStep(imgURL, "下载中")

	data, contentType, filename, err := downloadImageForMigrate(imgURL)
	if err != nil {
		return "", err
	}

	// 上传到目标存储；目录按来源归类（沿用本地存储的目录规范）
	dir := migrateTargetDir(imgURL)
	setCurrentStep(imgURL, "上传到"+storageTypeLabel[string(target)])
	newURL, err := util.UploadBytesToStorage(data, contentType, filename, dir, target)
	if err != nil {
		return "", err
	}

	return newURL, nil
}

// setCurrentStep 更新任务当前处理的 URL 与步骤（并发下仅作展示用途，最后写入者胜出）
func setCurrentStep(imgURL, step string) {
	globalMigrateTask.mu.Lock()
	globalMigrateTask.currentURL = imgURL
	globalMigrateTask.currentStep = step
	globalMigrateTask.mu.Unlock()
}

// migrateImageTimeout 单张图片迁移的下载超时
const migrateImageTimeout = 60 * time.Second

// downloadImageForMigrate 下载图片内容
// 本地存储 URL 直接读本地文件（同机部署时无需 HTTP），远程 URL 走 HTTP（GitHub raw 走代理）
func downloadImageForMigrate(imgURL string) ([]byte, string, string, error) {
	// 本地相对路径直接读盘
	if strings.HasPrefix(imgURL, "/uploads/") {
		localPath := strings.TrimPrefix(imgURL, "/")
		data, err := os.ReadFile(localPath)
		if err != nil {
			return nil, "", "", fmt.Errorf("读取本地文件失败: %w", err)
		}
		contentType := guessContentType(localPath)
		return data, contentType, filepath.Base(localPath), nil
	}

	// 远程 URL：GitHub raw 图片走代理
	requestURL := util.ResolveExportImageURL(imgURL, config.Cfg.ImageProxy.URL)
	client := &http.Client{Timeout: migrateImageTimeout}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, "", "", fmt.Errorf("构建下载请求失败: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", "", fmt.Errorf("下载图片失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", "", fmt.Errorf("下载图片失败: HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", "", fmt.Errorf("读取图片内容失败: %w", err)
	}

	// Content-Type 优先取响应头，取不到或非法时按扩展名推断
	contentType := resp.Header.Get("Content-Type")
	if idx := strings.Index(contentType, ";"); idx != -1 {
		contentType = strings.TrimSpace(contentType[:idx])
	}
	if !util_IsAllowedImageType(contentType) {
		contentType = guessContentType(imgURL)
	}
	if !util_IsAllowedImageType(contentType) {
		return nil, "", "", fmt.Errorf("不支持的图片类型: %s", contentType)
	}

	filename := filepath.Base(imgURL)
	if idx := strings.Index(filename, "?"); idx != -1 {
		filename = filename[:idx]
	}
	if filename == "" || filename == "/" {
		filename = "image" + guessExt(contentType)
	}
	// 文件名必须带扩展名，便于本地存储和图床识别
	if filepath.Ext(filename) == "" {
		filename += guessExt(contentType)
	}

	return data, contentType, filename, nil
}

// migrateTargetDir 迁移文件的目录归类（沿用本地存储目录规范，按 URL 路径特征推断用途）
func migrateTargetDir(imgURL string) string {
	u := imgURL
	if idx := strings.Index(u, "?"); idx != -1 {
		u = u[:idx]
	}
	lower := strings.ToLower(u)

	switch {
	case strings.Contains(lower, "/avatars/"):
		return util.AvatarDir
	case strings.Contains(lower, "/photoalbum/"):
		return util.PhotoAlbumDir
	case strings.Contains(lower, "/friendlinks/"):
		return util.FriendLinkDir
	case strings.Contains(lower, "/sitecovers/"):
		return util.SiteCoverDir
	case strings.Contains(lower, "/postcovers/"):
		return util.PostCoverDir
	default:
		return util.UploadDir
	}
}

// guessContentType 按文件扩展名推断 Content-Type
func guessContentType(name string) string {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return ""
	}
}

// guessExt 按 Content-Type 推断扩展名
func guessExt(contentType string) string {
	switch contentType {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ".jpg"
	}
}

// util_IsAllowedImageType 判断 Content-Type 是否为允许的图片类型
func util_IsAllowedImageType(contentType string) bool {
	switch contentType {
	case "image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp":
		return true
	}
	return false
}

// extractJSONArrayStrings 从 JSON 数组字符串中提取所有字符串元素
func extractJSONArrayStrings(value string) []string {
	value = strings.TrimSpace(value)
	if !strings.HasPrefix(value, "[") {
		return nil
	}

	var arr []string
	if err := json.Unmarshal([]byte(value), &arr); err != nil {
		return nil
	}
	return arr
}

// imageExtPattern 图片扩展名特征（用于判断 URL 是否指向图片文件）
var imageExtPattern = regexp.MustCompile(`\.(jpg|jpeg|png|gif|webp|svg|ico|avif)(\?.*)?$`)

// looksLikeImageURL 判断 URL 是否具有图片特征（图片扩展名）
// markdown 图片语法提取的 URL 已由语法保证是图片，无需再走此判断
func looksLikeImageURL(rawURL string) bool {
	u := strings.ToLower(strings.TrimSpace(rawURL))
	if u == "" {
		return false
	}
	// 去掉查询参数后按扩展名判断（imageExtPattern 已兼容带参数的 URL）
	return imageExtPattern.MatchString(u)
}

// extractJSONImageFields 从 JSON 对象字符串中提取图片类字段的值
// settings 里大量配置是 JSON 对象（friendlink_info.avatar/screenshot、about_* 的 img 等）
// 仅提取字段名含图片关键词的字符串值（字段名即图片证据，值只需是 URL 形态，不要求扩展名）
func extractJSONImageFields(value string) []string {
	value = strings.TrimSpace(value)
	if !strings.HasPrefix(value, "{") {
		return nil
	}

	var obj map[string]any
	if err := json.Unmarshal([]byte(value), &obj); err != nil {
		return nil
	}

	var urls []string
	for key, val := range obj {
		s, ok := val.(string)
		if !ok || s == "" || !containsImageFieldKeyword(key) {
			continue
		}
		// 值必须是 URL 形态（完整 http(s) 或本地 /uploads/ 路径），过滤图标名、纯文本等非 URL 值
		if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") || strings.HasPrefix(s, "/uploads/") {
			urls = append(urls, s)
		}
	}
	return urls
}

// containsImageFieldKeyword 字段名含图片关键词即视为图片字段（如 cover_bg_image、reward_wechat_img）
func containsImageFieldKeyword(key string) bool {
	lower := strings.ToLower(key)
	return strings.Contains(lower, "image") ||
		strings.Contains(lower, "img") ||
		strings.Contains(lower, "avatar") ||
		strings.Contains(lower, "cover") ||
		strings.Contains(lower, "logo") ||
		strings.Contains(lower, "screenshot") ||
		strings.Contains(lower, "icon")
}

// extractPlainURL 提取纯文本 URL（整段就是一个 URL 的情况）
func extractPlainURL(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || len(value) > 500 {
		return ""
	}
	// 含空白或引号的不算纯 URL
	if strings.ContainsAny(value, " \t\r\n\"'") {
		return ""
	}
	u, err := url.ParseRequestURI(value)
	if err != nil || u.Host == "" {
		// 本地相对路径
		if strings.HasPrefix(value, "/uploads/") {
			return value
		}
		return ""
	}
	return value
}

// containsString 判断字符串切片是否包含指定值
func containsString(list []string, target string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}

// ImageMigrateLogServiceItem 迁移日志列表项（前端展示用）
type ImageMigrateLogItem = model.ImageMigrateLog

// ListMigrateLogs 分页查询迁移日志
func (s *ImageService) ListMigrateLogs(page, pageSize int) ([]model.ImageMigrateLog, int64, error) {
	return s.logRepo.List(page, pageSize)
}

// DeleteMigrateLogs 删除迁移日志（ids 为空时清空全部）
func (s *ImageService) DeleteMigrateLogs(ids []uint) error {
	if len(ids) == 0 {
		return s.logRepo.DeleteAll()
	}
	return s.logRepo.DeleteBatch(ids)
}

// saveMigrateLog 将迁移任务结果写入持久化日志（失败不影响主流程，仅记录错误日志）
// mappings 传有序的成功映射切片（按迁移完成顺序），序列化后为 [{old_url,new_url}] 数组
func (s *ImageService) saveMigrateLog(req *MigrateImageRequest, mappings []MigrateMappingItem, failed []MigrateFailedItem, updatedRows, deletedFiles int, message string, duration time.Duration, operatorID *uint, operatorName string) {
	var mappingsJSON, failedJSON datatypes.JSON

	// mappings 统一序列化为数组格式 [{old_url,new_url}]，与前端 ImageMigrateLog.mappings 类型约定一致
	if len(mappings) > 0 {
		if data, err := json.Marshal(mappings); err == nil {
			mappingsJSON = datatypes.JSON(data)
		}
	} else {
		mappingsJSON = datatypes.JSON("[]")
	}
	if len(failed) > 0 {
		if data, err := json.Marshal(failed); err == nil {
			failedJSON = datatypes.JSON(data)
		}
	} else {
		failedJSON = datatypes.JSON("[]")
	}

	log := &model.ImageMigrateLog{
		OperatorID:    operatorID,
		OperatorName:  operatorName,
		TargetStorage: req.TargetStorag,
		DeleteSource:  req.DeleteSource,
		TotalCount:    len(mappings) + len(failed),
		SuccessCount:  len(mappings),
		FailedCount:   len(failed),
		UpdatedRows:   updatedRows,
		DeletedFiles:  deletedFiles,
		Mappings:      mappingsJSON,
		FailedItems:   failedJSON,
		Message:       message,
		DurationSec:   int(duration.Seconds()),
	}

	if err := s.logRepo.Create(log); err != nil {
		if config.Cfg != nil {
			zap.L().Error("保存图片迁移日志失败", zap.Error(err))
		}
	}
}

// validateStorageConfig 校验目标存储类型的配置是否可用
func validateStorageConfig(storageType util.StorageType) error {
	switch storageType {
	case util.StorageOSS:
		return util.ValidateOSSConfig(
			config.Cfg.OSS.Endpoint,
			config.Cfg.OSS.AccessKeyID,
			config.Cfg.OSS.AccessKeySecret,
			config.Cfg.OSS.BucketName,
		)
	case util.StorageCOS:
		return util.ValidateCOSConfig(
			config.Cfg.COS.BucketURL,
			config.Cfg.COS.SecretID,
			config.Cfg.COS.SecretKey,
		)
	case util.StorageEasyImage:
		return util.ValidateEasyImageConfig(
			config.Cfg.EasyImage.ApiURL,
			config.Cfg.EasyImage.Token,
		)
	}
	return nil // local 无需外部配置
}
