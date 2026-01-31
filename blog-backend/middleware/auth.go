/*
 * 项目名称：blog-backend
 * 文件名称：auth.go
 * 创建时间：2026-01-31 16:17:29
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：认证中间件，提供JWT认证、可选认证和管理员权限验证功能
 */
package middleware

import (
	"strings"

	"blog-backend/util"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
// 功能说明：验证请求中的JWT Token，解析用户信息并存入上下文
// 要求：请求头必须包含有效的Authorization: Bearer <token>
// 返回:
//   - gin.HandlerFunc: Gin中间件处理函数
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization请求头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			util.Unauthorized(c, "缺少认证信息")
			c.Abort()
			return
		}

		// 验证格式: Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			util.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}

		// 解析Token
		claims, err := util.ParseToken(parts[1])
		if err != nil {
			util.Unauthorized(c, "无效的认证信息")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中，供后续处理器使用
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件（用于WebSocket等场景）
// 功能说明：尝试解析JWT Token，如果存在且有效则存储用户信息，但不强制要求认证
// 适用场景：WebSocket连接、公开接口等需要可选认证的场景
// 返回:
//   - gin.HandlerFunc: Gin中间件处理函数
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		// 优先从Header获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}

		// 如果Header中没有，尝试从查询参数获取（用于WebSocket）
		if token == "" {
			token = c.Query("token")
		}

		// 如果有token，尝试解析
		if token != "" {
			claims, err := util.ParseToken(token)
			if err == nil {
				// 认证成功，存储用户信息
				c.Set("user_id", claims.UserID)
				c.Set("username", claims.Username)
				c.Set("role", claims.Role)
			}
		}

		// 无论是否认证成功都继续执行
		c.Next()
	}
}

// AdminMiddleware 管理员权限中间件
// 功能说明：验证当前用户是否具有管理员权限
// 要求：必须在AuthMiddleware或OptionalAuthMiddleware之后使用
// 返回:
//   - gin.HandlerFunc: Gin中间件处理函数
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取用户角色
		role, exists := c.Get("role")
		if !exists {
			util.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		// 验证是否为管理员
		if role != "admin" {
			util.Forbidden(c, "需要管理员权限")
			c.Abort()
			return
		}

		c.Next()
	}
}
