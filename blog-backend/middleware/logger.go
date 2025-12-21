package middleware

import (
	"time"

	"blog-backend/logger"
	"blog-backend/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		cost := time.Since(start)

		// 从上下文获取IP，如果不存在则使用工具函数获取
		ip := util.GetClientIP(c)

		logger.Info("HTTP Request",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", c.Writer.Status()),
			zap.String("ip", ip),
			zap.Duration("cost", cost),
			zap.String("user-agent", c.Request.UserAgent()),
		)
	}
}
