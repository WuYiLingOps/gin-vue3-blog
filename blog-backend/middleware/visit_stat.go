package middleware

import (
	"blog-backend/repository"
	"blog-backend/util"
	"strings"

	"github.com/gin-gonic/gin"
)

// VisitStatMiddleware 访问量统计中间件（统计独立访客 UV）
func VisitStatMiddleware() gin.HandlerFunc {
	visitStatRepo := repository.NewVisitStatRepository()

	return func(c *gin.Context) {
		// 先处理请求
		c.Next()

		// 只统计以下情况：
		// 1. GET 请求
		// 2. 200-299 状态码
		// 3. 非静态资源请求（排除 /uploads, /assets 等）
		// 4. 非 API 请求（排除 /api）
		path := c.Request.URL.Path
		method := c.Request.Method
		status := c.Writer.Status()

		// 过滤条件
		isStaticResource := strings.HasPrefix(path, "/uploads") ||
			strings.HasPrefix(path, "/assets") ||
			strings.HasPrefix(path, "/favicon.ico") ||
			strings.HasSuffix(path, ".js") ||
			strings.HasSuffix(path, ".css") ||
			strings.HasSuffix(path, ".png") ||
			strings.HasSuffix(path, ".jpg") ||
			strings.HasSuffix(path, ".jpeg") ||
			strings.HasSuffix(path, ".gif") ||
			strings.HasSuffix(path, ".svg") ||
			strings.HasSuffix(path, ".ico")

		isAPIRequest := strings.HasPrefix(path, "/api")

		// 只统计：GET 请求 + 成功状态 + 非静态资源 + 非 API
		if method == "GET" && status >= 200 && status < 300 && !isStaticResource && !isAPIRequest {
			// 获取访客 IP 地址
			ip := util.GetClientIP(c)
			
			// 异步记录访问量（基于 IP 的 UV 统计），不影响响应速度
			go func() {
				visitStatRepo.RecordVisit(ip)
			}()
		}
	}
}

