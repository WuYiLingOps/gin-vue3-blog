/*
 * 项目名称：blog-backend
 * 文件名称：database.go
 * 创建时间：2026-01-31 16:03:44
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：数据库连接管理模块，负责初始化和管理PostgreSQL数据库连接
 */
package db

import (
	"fmt"

	"blog-backend/config"
	blogLogger "blog-backend/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库连接实例
var DB *gorm.DB

// InitDB 初始化PostgreSQL数据库连接
// 功能说明：
//  1. 从配置中读取数据库连接信息
//  2. 构建PostgreSQL连接字符串（DSN）
//  3. 使用GORM打开数据库连接
//  4. 设置时区为Asia/Shanghai
//  5. 关闭SQL日志输出（生产环境推荐）
//
// 返回:
//   - error: 连接失败时返回错误
func InitDB() error {
	// 获取数据库配置
	cfg := config.Cfg.DB

	// 构建PostgreSQL连接字符串（DSN）
	// 格式：host=xxx user=xxx password=xxx dbname=xxx port=xxx sslmode=xxx TimeZone=Asia/Shanghai
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s "+
		"TimeZone=Asia/Shanghai",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	// 使用GORM打开数据库连接
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 关闭SQL日志输出，减少日志噪音
	})

	if err != nil {
		return err
	}

	// 连接成功，记录日志
	blogLogger.Info("Database connected successfully")
	return nil
}
