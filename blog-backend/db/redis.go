/*
 * 项目名称：blog-backend
 * 文件名称：redis.go
 * 创建时间：2026-01-31 16:03:44
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：Redis缓存连接管理模块，负责初始化和管理Redis连接
 */
package db

import (
	"context"
	"fmt"

	"blog-backend/config"
	blogLogger "blog-backend/logger"

	"github.com/redis/go-redis/v9"
)

// RDB 全局Redis客户端实例
var RDB *redis.Client

// InitRedis 初始化Redis连接
// 功能说明：
//  1. 从配置中读取Redis连接信息
//  2. 创建Redis客户端实例
//  3. 通过Ping命令测试连接是否正常
//  4. 连接成功后记录日志
//
// 返回:
//   - error: 连接失败时返回错误
func InitRedis() error {
	// 获取Redis配置
	cfg := config.Cfg.Redis

	// 创建Redis客户端实例
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), // Redis服务器地址和端口
		Password: cfg.Password,                             // Redis密码（如果设置了）
		DB:       cfg.DB,                                   // Redis数据库编号（默认0）
	})

	// 测试连接是否正常
	ctx := context.Background()
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		return err
	}

	// 连接成功，记录日志
	blogLogger.Info("Redis connected successfully")
	return nil
}
