package db

import (
	"context"
	"fmt"

	"blog-backend/config"
	blogLogger "blog-backend/logger"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

// InitRedis 初始化Redis连接
func InitRedis() error {
	cfg := config.Cfg.Redis
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	ctx := context.Background()
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		return err
	}

	blogLogger.Info("Redis connected successfully")
	return nil
}
