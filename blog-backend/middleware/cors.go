/*
 * 项目名称：blog-backend
 * 文件名称：cors.go
 * 创建时间：2026-01-31 16:17:29
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：跨域资源共享（CORS）中间件，处理跨域请求和预检请求
 */
package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS 跨域资源共享中间件
// 功能说明：设置CORS响应头，允许跨域请求，处理预检请求（OPTIONS）
// 返回:
//   - gin.HandlerFunc: Gin中间件处理函数
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许的源（*表示允许所有源）
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// 允许携带凭证（cookies等）
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// 允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Token")
		// 允许的请求方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		// 处理预检请求（OPTIONS），直接返回204 No Content
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
