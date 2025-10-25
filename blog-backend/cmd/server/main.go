package main

import (
	"fmt"
	"log"

	"blog-backend/config"
	"blog-backend/db"
	"blog-backend/logger"
	"blog-backend/util"

	"github.com/gin-gonic/gin"
)

func main() {
	// 根据环境变量加载配置
	if err := config.LoadConfigByEnv(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	if err := logger.InitLogger(config.Cfg.Log.Level); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting blog backend server...")

	// 初始化数据库
	if err := db.InitDB(); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to init database: %v", err))
	}

	// 初始化上传目录
	if err := util.InitUploadDirs(); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to init upload directories: %v", err))
	}

	// 设置 Gin 模式
	gin.SetMode(config.Cfg.Server.Mode)

	// 配置路由
	r := setupRouter()

	// 启动服务器
	addr := fmt.Sprintf(":%d", config.Cfg.App.Port)
	logger.Info(fmt.Sprintf("Server is running on http://localhost%s", addr))

	if err := r.Run(addr); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to start server: %v", err))
	}
}
