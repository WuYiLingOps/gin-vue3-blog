package middleware

import (
	"strings"

	"blog-backend/util"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			util.Unauthorized(c, "缺少认证信息")
			c.Abort()
			return
		}

		// 格式: Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			util.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}

		// 解析 Token
		claims, err := util.ParseToken(parts[1])
		if err != nil {
			util.Unauthorized(c, "无效的认证信息")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件（用于WebSocket等场景）
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
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			util.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		if role != "admin" {
			util.Forbidden(c, "需要管理员权限")
			c.Abort()
			return
		}

		c.Next()
	}
}
