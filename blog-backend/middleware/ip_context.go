package middleware

import (
	"blog-backend/util"

	"github.com/gin-gonic/gin"
)

// IPContextMiddleware IP上下文中间件
// 将获取到的客户端IP存入请求上下文，供后续中间件和处理器使用
func IPContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端IP
		ip := util.GetClientIP(c)

		// 存入上下文
		c.Set(util.ClientIPKey, ip)

		c.Next()
	}
}
