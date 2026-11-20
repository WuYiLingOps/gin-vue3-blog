/*
 * 项目名称：blog-backend
 * 文件名称：image_cache.go
 * 创建时间：2026-11-20 10:00:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：图片盘点缓存层，用 Redis 缓存全站图片盘点结果，减轻数据库压力
 */
package service

import (
	"context"
	"encoding/json"
	"time"

	"blog-backend/db"

	"go.uber.org/zap"
)

const (
	// imageScanCacheKey 图片盘点结果的 Redis 缓存键
	imageScanCacheKey = "image:assets:v1"
	// imageScanCacheTTL 盘点缓存过期时间（兜底其他业务改图后的数据新鲜度）
	imageScanCacheTTL = 10 * time.Minute
)

// loadScanCache 从 Redis 读取盘点缓存（Redis 不可用或未命中返回 false，静默降级为查库）
func loadScanCache() ([]ImageAsset, bool) {
	if db.RDB == nil {
		return nil, false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	data, err := db.RDB.Get(ctx, imageScanCacheKey).Bytes()
	if err != nil || len(data) == 0 {
		return nil, false
	}

	var assets []ImageAsset
	if err := json.Unmarshal(data, &assets); err != nil {
		// 缓存内容异常，删除后降级查库
		_ = db.RDB.Del(ctx, imageScanCacheKey).Err()
		return nil, false
	}

	return assets, true
}

// saveScanCache 将盘点结果写入 Redis（写入失败不影响主流程）
func saveScanCache(assets []ImageAsset) {
	if db.RDB == nil {
		return
	}

	data, err := json.Marshal(assets)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.RDB.Set(ctx, imageScanCacheKey, data, imageScanCacheTTL).Err(); err != nil {
		zap.L().Warn("图片盘点缓存写入失败", zap.Error(err))
	}
}

// InvalidateScanCache 删除图片盘点缓存（迁移任务替换数据库后调用，保证盘点结果与库一致）
func InvalidateScanCache() {
	if db.RDB == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.RDB.Del(ctx, imageScanCacheKey).Err(); err != nil {
		zap.L().Warn("图片盘点缓存删除失败", zap.Error(err))
	}
}
