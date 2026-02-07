/*
 * 项目名称：blog-backend
 * 文件名称：operation_log.go
 * 创建时间：2026-02-06 22:00:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：操作日志工具函数，提供便捷的操作日志记录功能
 */
package util

import (
	"blog-backend/constant"
	"blog-backend/model"
	"blog-backend/repository"

	"github.com/gin-gonic/gin"
)

// LogOperation 记录操作日志
// 参数:
//   - c: Gin上下文
//   - action: 操作类型 (create, update, delete)
//   - module: 操作模块 (post, category, tag, user, comment等)
//   - targetID: 目标ID（可为nil）
//   - targetName: 目标名称（如文章标题、分类名称等）
//   - description: 操作描述
func LogOperation(c *gin.Context, action, module string, targetID *uint, targetName, description string) {
	// 获取用户信息
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return // 未登录用户不记录日志
	}

	userID := userIDVal.(uint)

	// 获取用户角色
	roleVal, exists := c.Get("role")
	if !exists {
		return
	}

	role := roleVal.(string)

	// 只记录管理员和超级管理员的操作
	if !constant.IsAdminRole(role) {
		return
	}

	// 获取用户名
	usernameVal, exists := c.Get("username")
	username := ""
	if exists {
		username = usernameVal.(string)
	} else {
		// 如果上下文中没有用户名，从数据库查询
		userRepo := repository.NewUserRepository()
		user, err := userRepo.GetByID(userID)
		if err == nil {
			username = user.Username
		}
	}

	// 获取IP和UserAgent
	ip := GetClientIP(c)
	userAgent := c.GetHeader("User-Agent")

	// 创建操作日志
	log := &model.OperationLog{
		UserID:      userID,
		Username:    username,
		Action:      action,
		Module:      module,
		TargetType:  module,
		TargetID:    targetID,
		TargetName:  targetName,
		Description: description,
		IP:          ip,
		UserAgent:   userAgent,
	}

	// 异步记录日志（避免影响主流程性能）
	go func() {
		logRepo := repository.NewOperationLogRepository()
		_ = logRepo.Create(log)
	}()
}
