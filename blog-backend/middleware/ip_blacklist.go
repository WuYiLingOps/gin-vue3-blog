/*
 * 项目名称：blog-backend
 * 文件名称：ip_blacklist.go
 * 创建时间：2026-01-31 16:17:29
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：IP黑名单检查中间件，实现IP封禁、访问频率限制、白名单管理等功能，支持管理员自动解封
 */
package middleware

import (
	"blog-backend/config"
	"blog-backend/constant"
	"blog-backend/db"
	"blog-backend/model"
	"blog-backend/util"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// ipAccessRecord IP访问记录结构体
type ipAccessRecord struct {
	count     int       // 访问次数
	firstTime time.Time // 首次访问时间
	banned    bool      // 是否已被封禁
}

var (
	// ipAccessMap IP访问记录映射表，key为IP地址
	ipAccessMap = make(map[string]*ipAccessRecord)
	// ipMapMutex 读写互斥锁，保护ipAccessMap的并发访问
	ipMapMutex sync.RWMutex

	// 配置参数
	maxRequestsPerMinute    = 120 // 每分钟最大请求数
	maxRequestsPer10Min     = 600 // 10分钟最大请求数
	banDurationMinutes      = 30  // 自动封禁时长（分钟）
	adminAutoWhitelistHours = 2   // 管理员登录后当前IP自动加入白名单的时长（小时）
)

// IPBlacklistMiddleware IP黑名单检查中间件
// 功能说明：
//  1. 检查IP是否在黑名单中，封禁的IP将被拒绝访问（认证路径除外）
//  2. 监控IP访问频率，超过限制的IP将被自动封禁
//  3. 管理员和白名单IP不受限制
//  4. 管理员登录后，其IP会自动加入临时白名单
//
// 返回:
//   - gin.HandlerFunc: Gin中间件处理函数
func IPBlacklistMiddleware() gin.HandlerFunc {
	// 启动定时清理过期记录的协程
	go cleanupExpiredRecords()

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		// 从上下文获取IP（已通过IPContextMiddleware设置）
		ip := util.GetClientIP(c)

		// 1. 完全排除的路径：静态文件和 WebSocket
		if strings.HasPrefix(path, "/uploads") || strings.HasPrefix(path, "/api/chat/ws") {
			c.Next()
			return
		}

		// 2. 登录相关路径：允许通过黑名单检查，但记录访问（以便管理员能够登录）
		// 这些路径即使IP被封禁也要允许访问，给管理员登录的机会
		isAuthPath := strings.HasPrefix(path, "/api/auth/login") ||
			strings.HasPrefix(path, "/api/captcha") ||
			strings.HasPrefix(path, "/api/auth/register") ||
			strings.HasPrefix(path, "/api/auth/send-register-code") ||
			strings.HasPrefix(path, "/api/auth/forgot-password") ||
			strings.HasPrefix(path, "/api/auth/reset-password") ||
			strings.HasPrefix(path, "/api/auth/refresh")

		// 3. 优先检查：如果是管理员或白名单IP，跳过所有限制
		// 管理员和白名单IP不会被计入访问频率，也不会被封禁
		isAdmin := isAdminUser(c)
		isWhitelisted := isIPInWhitelist(ip)

		// 3.1 如果是管理员用户：即使当前 IP 在黑名单或频率记录中，也自动解除封禁并加入临时白名单
		// 这样从该 IP 发出的后续请求（包括未携带 Token 的公开接口）都不会再被误封
		if isAdmin {
			if isIPBanned(ip) {
				unbanIP(ip)
			}
			addIPToTemporaryWhitelist(ip)
			c.Next()
			return
		}

		// 3.2 如果是白名单 IP：直接放行，不做任何限制
		if isWhitelisted {
			c.Next()
			return
		}

		// 4. 检查是否在黑名单中
		// 关键优化：认证相关路径（登录、验证码等）跳过黑名单检查，给管理员登录的机会
		if !isAuthPath && isIPBanned(ip) {
			util.Error(c, 403, "您的IP已被封禁，请联系管理员")
			c.Abort()
			return
		}

		// 5. 检查访问频率
		// 认证路径也要检查频率，但不会自动封禁（给管理员登录机会）
		if shouldBanIP(ip) {
			// 封禁前最后检查：确保不会误封管理员
			if isAdminUser(c) || isIPInWhitelist(ip) {
				// 管理员和白名单IP即使访问频繁也不封禁
				c.Next()
				return
			}

			// 认证路径不自动封禁，只返回频率限制错误
			if isAuthPath {
				util.Error(c, 429, "访问过于频繁，请稍后再试")
				c.Abort()
				return
			}

			// 非认证路径：自动封禁IP
			banIP(ip, "访问频率过高，自动封禁", 1)
			util.Error(c, 429, "访问过于频繁，您的IP已被临时封禁")
			c.Abort()
			return
		}

		// 6. 记录访问（只记录非管理员、非白名单、非本地IP的访问）
		// 本地开发环境IP（127.0.0.1 和 ::1）不记录访问频率
		if ip != "127.0.0.1" && ip != "::1" {
			recordAccess(ip)
		}

		c.Next()
	}
}

// isIPBanned 检查IP是否在黑名单中
func isIPBanned(ip string) bool {
	var blacklist model.IPBlacklist
	err := db.DB.Where("ip = ?", ip).First(&blacklist).Error

	if err != nil {
		return false
	}

	// 检查是否已过期
	if blacklist.ExpireAt != nil && blacklist.ExpireAt.Before(time.Now()) {
		// 已过期，删除记录
		db.DB.Delete(&blacklist)
		return false
	}

	return true
}

// recordAccess 记录IP访问
func recordAccess(ip string) {
	ipMapMutex.Lock()
	defer ipMapMutex.Unlock()

	now := time.Now()

	if record, exists := ipAccessMap[ip]; exists {
		// 如果距离首次访问超过10分钟，重置计数
		if now.Sub(record.firstTime) > 10*time.Minute {
			ipAccessMap[ip] = &ipAccessRecord{
				count:     1,
				firstTime: now,
				banned:    false,
			}
		} else {
			record.count++
		}
	} else {
		ipAccessMap[ip] = &ipAccessRecord{
			count:     1,
			firstTime: now,
			banned:    false,
		}
	}
}

// shouldBanIP 判断是否应该封禁IP
func shouldBanIP(ip string) bool {
	ipMapMutex.RLock()
	defer ipMapMutex.RUnlock()

	record, exists := ipAccessMap[ip]
	if !exists {
		return false
	}

	now := time.Now()
	duration := now.Sub(record.firstTime)

	// 1分钟内请求超过限制
	if duration <= 1*time.Minute && record.count > maxRequestsPerMinute {
		return true
	}

	// 10分钟内请求超过限制
	if duration <= 10*time.Minute && record.count > maxRequestsPer10Min {
		return true
	}

	return false
}

// banIP 封禁IP
func banIP(ip string, reason string, banType int) {
	// 如果IP为空，不进行封禁
	if ip == "" {
		return
	}

	// 本地开发环境IP不封禁（127.0.0.1 和 ::1）
	if ip == "127.0.0.1" || ip == "::1" {
		return
	}

	// 再次检查是否在白名单中（防止配置未加载的情况）
	if isIPInWhitelist(ip) {
		return
	}

	expireAt := time.Now().Add(time.Duration(banDurationMinutes) * time.Minute)

	blacklist := model.IPBlacklist{
		IP:       ip,
		Reason:   reason,
		BanType:  banType,
		ExpireAt: &expireAt,
	}

	// 检查是否已存在
	var existing model.IPBlacklist
	err := db.DB.Where("ip = ?", ip).First(&existing).Error

	if err != nil {
		// 不存在，创建新记录
		if err := db.DB.Create(&blacklist).Error; err != nil {
			// 创建失败，记录错误但不影响主流程
			return
		}
	} else {
		// 已存在，更新记录（更新过期时间和原因）
		existing.Reason = reason
		existing.BanType = banType
		existing.ExpireAt = &expireAt
		db.DB.Save(&existing)
	}

	// 标记为已封禁
	ipMapMutex.Lock()
	if record, exists := ipAccessMap[ip]; exists {
		record.banned = true
	}
	ipMapMutex.Unlock()
}

// unbanIP 解除 IP 封禁（供管理员身份自动解封使用）
// 注意：仅在当前请求已被确认是管理员用户时调用
func unbanIP(ip string) {
	if ip == "" {
		return
	}

	// 删除数据库中的黑名单记录
	db.DB.Where("ip = ?", ip).Delete(&model.IPBlacklist{})

	// 清理内存中的访问记录封禁标记
	ipMapMutex.Lock()
	if record, exists := ipAccessMap[ip]; exists {
		record.banned = false
	}
	ipMapMutex.Unlock()
}

// addIPToTemporaryWhitelist 管理员登录后，将当前 IP 加入一个有限期的白名单
// 目的是：管理员确认无异常后，该 IP 在一段时间内不再因高频访问被自动封禁
func addIPToTemporaryWhitelist(ip string) {
	if ip == "" {
		return
	}

	// 本地开发 IP 无需加入白名单
	if ip == "127.0.0.1" || ip == "::1" {
		return
	}

	// 计算新的过期时间
	now := time.Now()
	newExpire := now.Add(time.Duration(adminAutoWhitelistHours) * time.Hour)

	var existing model.IPWhitelist
	if err := db.DB.Where("ip = ?", ip).First(&existing).Error; err == nil {
		// 已存在白名单记录：如果是永久白名单或过期时间晚于新的时间，则不修改
		if existing.ExpireAt == nil || existing.ExpireAt.After(newExpire) {
			return
		}

		// 否则延长白名单有效期
		existing.ExpireAt = &newExpire
		db.DB.Save(&existing)
		return
	}

	// 不存在白名单记录，创建临时白名单
	whitelist := model.IPWhitelist{
		IP:       ip,
		Reason:   "管理员登录自动加入临时白名单",
		ExpireAt: &newExpire,
	}
	_ = db.DB.Create(&whitelist).Error
}

// shouldSkipRateLimit 判断是否应该跳过频率限制
// 返回 true 表示应该跳过（管理员或白名单IP）
// 注意：为了提高性能和多次检查的可靠性，中间件中已内联此逻辑
func shouldSkipRateLimit(c *gin.Context, ip string) bool {
	// 方案一：检查是否是管理员用户（通过解析 Token）
	if isAdminUser(c) {
		return true
	}

	// 方案二：检查 IP 是否在配置的白名单中
	if isIPInWhitelist(ip) {
		return true
	}

	return false
}

// isAdminUser 检查当前用户是否是管理员
// 通过解析 JWT Token 获取用户角色
func isAdminUser(c *gin.Context) bool {
	// 优先从上下文获取（如果 OptionalAuthMiddleware 已执行）
	if role, exists := c.Get("role"); exists {
		if r, ok := role.(string); ok {
			return constant.IsAdminRole(r)
		}
		return false
	}

	// 如果上下文中没有，尝试从 Header 解析 Token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return false
	}

	// 格式: Bearer <token>
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return false
	}

	// 解析 Token
	claims, err := util.ParseToken(parts[1])
	if err != nil {
		return false
	}

	// 检查是否是具备管理员权限的角色
	return constant.IsAdminRole(claims.Role)
}

// isIPInWhitelist 检查 IP 是否在白名单中（配置文件 + 数据库）
func isIPInWhitelist(ip string) bool {
	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		return false
	}

	// 1. 检查配置文件中的白名单
	if config.Cfg != nil && len(config.Cfg.Security.AdminIPWhitelist) > 0 {
		for _, whitelistIP := range config.Cfg.Security.AdminIPWhitelist {
			// 支持 CIDR 格式（如：192.168.1.0/24）
			if strings.Contains(whitelistIP, "/") {
				_, ipNet, err := net.ParseCIDR(whitelistIP)
				if err == nil && ipNet.Contains(clientIP) {
					return true
				}
			} else {
				// 精确匹配
				whitelistIPParsed := net.ParseIP(whitelistIP)
				if whitelistIPParsed != nil && whitelistIPParsed.Equal(clientIP) {
					return true
				}
			}
		}
	}

	// 2. 检查数据库中的白名单
	var whitelist []model.IPWhitelist
	if err := db.DB.Find(&whitelist).Error; err == nil {
		for _, wl := range whitelist {
			// 检查是否过期
			if wl.ExpireAt != nil && wl.ExpireAt.Before(time.Now()) {
				continue
			}

			// 支持 CIDR 格式
			if strings.Contains(wl.IP, "/") {
				_, ipNet, err := net.ParseCIDR(wl.IP)
				if err == nil && ipNet.Contains(clientIP) {
					return true
				}
			} else {
				// 精确匹配
				wlIP := net.ParseIP(wl.IP)
				if wlIP != nil && wlIP.Equal(clientIP) {
					return true
				}
			}
		}
	}

	return false
}

// cleanupExpiredRecords 定时清理过期的访问记录和黑名单
func cleanupExpiredRecords() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		// 清理内存中的过期访问记录
		ipMapMutex.Lock()
		now := time.Now()
		for ip, record := range ipAccessMap {
			if now.Sub(record.firstTime) > 15*time.Minute {
				delete(ipAccessMap, ip)
			}
		}
		ipMapMutex.Unlock()

		// 清理数据库中的过期黑名单
		db.DB.Where("expire_at IS NOT NULL AND expire_at < ?", now).Delete(&model.IPBlacklist{})
	}
}
