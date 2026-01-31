/*
 * 项目名称：blog-backend
 * 文件名称：logger.go
 * 创建时间：2026-01-31 16:17:29
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：HTTP请求日志中间件，记录请求方法、路径、状态码、耗时等信息
 */
package middleware

import (
	"time"

	"blog-backend/logger"
	"blog-backend/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger HTTP请求日志中间件
// 功能说明：记录每个HTTP请求的详细信息，包括请求方法、路径、查询参数、状态码、客户端IP、耗时和User-Agent
// 返回:
//   - gin.HandlerFunc: Gin中间件处理函数
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 继续执行后续中间件和处理器
		c.Next()

		// 计算请求耗时
		cost := time.Since(start)

		// 从上下文获取IP，如果不存在则使用工具函数获取
		ip := util.GetClientIP(c)

		// 记录请求日志
		logger.Info("HTTP Request",
			zap.String("method", c.Request.Method),          // 请求方法
			zap.String("path", path),                        // 请求路径
			zap.String("query", query),                      // 查询参数
			zap.Int("status", c.Writer.Status()),            // HTTP状态码
			zap.String("ip", ip),                            // 客户端IP
			zap.Duration("cost", cost),                      // 请求耗时
			zap.String("user-agent", c.Request.UserAgent()), // User-Agent
		)
	}
}
