/*
 * 项目名称：blog-backend
 * 文件名称：main.go
 * 创建时间：2026-01-31 16:00:22
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：博客后端服务器主程序，负责初始化配置、数据库、路由等并启动HTTP服务
 */
package main

import (
	"fmt"
	"log"

	"blog-backend/config"
	"blog-backend/db"
	"blog-backend/logger"
	"blog-backend/router"
	"blog-backend/service"
	"blog-backend/util"

	"github.com/gin-gonic/gin"
)

func main() {
	// 根据环境变量加载配置
	if err := config.LoadConfigByEnv(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	// 根据 config.yml 中的 env 判断是否为开发环境
	isDev := config.Cfg.Env == "dev"
	if err := logger.InitLogger(config.Cfg.Log.Level, isDev); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting blog backend server...")

	// 初始化数据库
	if err := db.InitDB(); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to init database: %v", err))
	}

	// 初始化Redis
	if err := db.InitRedis(); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to init redis: %v", err))
	}

	// 初始化上传目录
	if err := util.InitUploadDirs(); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to init upload directories: %v", err))
	}

	// 启动定期清理任务
	cleanupService := service.NewCleanupService()
	cleanupService.StartCleanupTasks()
	logger.Info("Cleanup tasks started")

	// 设置 Gin 模式
	gin.SetMode(config.Cfg.Server.Mode)

	// 配置路由
	r := router.SetupRouter()

	// 启动服务器
	addr := fmt.Sprintf(":%d", config.Cfg.App.Port)
	logger.Info(fmt.Sprintf("Server is running on http://localhost%s", addr))
	if err := r.Run(addr); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to start server: %v", err))
	}
}
