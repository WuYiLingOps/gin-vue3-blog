/*
 * 项目名称：blog-backend
 * 文件名称：ip_context.go
 * 创建时间：2026-01-31 16:17:29
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：IP上下文中间件，将客户端IP存入请求上下文，供后续中间件和处理器使用
 */
package middleware

import (
	"blog-backend/util"

	"github.com/gin-gonic/gin"
)

// IPContextMiddleware IP上下文中间件
// 功能说明：将获取到的客户端IP存入请求上下文，供后续中间件和处理器使用
// 返回:
//   - gin.HandlerFunc: Gin中间件处理函数
func IPContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端IP
		ip := util.GetClientIP(c)

		// 存入上下文
		c.Set(util.ClientIPKey, ip)

		c.Next()
	}
}
