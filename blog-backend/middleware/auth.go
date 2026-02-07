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
	"errors"
	"strings"

	"blog-backend/constant"
	"blog-backend/repository"
	"blog-backend/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

// RoleRequiredMiddleware 角色权限中间件
// 功能说明：验证当前用户是否属于指定角色之一
// 要求：必须在 AuthMiddleware 或 OptionalAuthMiddleware 之后使用
// 参数:
//   - roles: 允许访问的角色列表
//
// 返回:
//   - gin.HandlerFunc: Gin中间件处理函数
// 注意：此中间件会从数据库验证用户的实际角色，确保角色变更后旧token立即失效
func RoleRequiredMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID（从AuthMiddleware或OptionalAuthMiddleware设置）
		userIDVal, exists := c.Get("user_id")
		if !exists {
			util.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		userID, ok := userIDVal.(uint)
		if !ok {
			util.Forbidden(c, "用户信息异常")
			c.Abort()
			return
		}

		// 从数据库查询用户的实际角色和状态（确保使用最新数据）
		userRepo := repository.NewUserRepository()
		user, err := userRepo.GetByID(userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				util.Unauthorized(c, "用户不存在，请重新登录")
				c.Abort()
				return
			}
			util.Error(c, 500, "获取用户信息失败")
			c.Abort()
			return
		}

		// 检查用户状态（如果被禁用，拒绝访问）
		if user.Status != 1 {
			util.Unauthorized(c, "账号已被禁用，请重新登录")
			c.Abort()
			return
		}

		// 获取token中的角色（用于对比）
		tokenRoleVal, exists := c.Get("role")
		if !exists {
			util.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		tokenRole, ok := tokenRoleVal.(string)
		if !ok {
			util.Forbidden(c, "权限信息异常")
			c.Abort()
			return
		}

		// 验证token中的角色与数据库中的角色是否一致
		// 如果不一致，说明用户角色已被修改，要求重新登录
		if tokenRole != user.Role {
			util.Unauthorized(c, "用户权限已变更，请重新登录")
			c.Abort()
			return
		}

		// 使用数据库中的实际角色进行权限检查
		currentRole := user.Role
		for _, allowed := range roles {
			if currentRole == allowed {
				// 更新上下文中的角色为数据库中的最新角色（确保后续使用最新数据）
				c.Set("role", currentRole)
				c.Next()
				return
			}
		}

		util.Forbidden(c, "需要更高权限")
		c.Abort()
	}
}

// AdminMiddleware 兼容保留：管理员权限中间件
// 目前视为具备后台管理权限的角色（包含 super_admin 和 admin）
func AdminMiddleware() gin.HandlerFunc {
	return RoleRequiredMiddleware(constant.RoleSuperAdmin, constant.RoleAdmin)
}
