package middleware

import (
	"blog-backend/config"
	"blog-backend/db"
	"blog-backend/model"
	"blog-backend/util"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// IP访问记录
type ipAccessRecord struct {
	count     int
	firstTime time.Time
	banned    bool
}

var (
	// IP访问记录映射
	ipAccessMap = make(map[string]*ipAccessRecord)
	ipMapMutex  sync.RWMutex

	// 配置参数
	maxRequestsPerMinute = 60  // 每分钟最大请求数
	maxRequestsPer10Min  = 300 // 10分钟最大请求数
	banDuration          = 1   // 自动封禁时长（小时）
)

// IPBlacklistMiddleware IP黑名单检查中间件
func IPBlacklistMiddleware() gin.HandlerFunc {
	// 启动定时清理过期记录的协程
	go cleanupExpiredRecords()

	return func(c *gin.Context) {
		path := c.Request.URL.Path
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

		if isAdmin || isWhitelisted {
			// 管理员和白名单IP完全豁免，不记录访问，不检查黑名单
			c.Next()
			return
		}

		// 4. 检查是否在黑名单中
		// 关键优化：认证相关路径（登录、验证码等）跳过黑名单检查，给管理员登录的机会
		if !isAuthPath && isIPBanned(ip) {
			// 再次检查是否是管理员或白名单（防止Token在后续处理中才被正确解析）
			if isAdminUser(c) || isIPInWhitelist(ip) {
				c.Next()
				return
			}
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

		// 6. 记录访问（只记录非管理员、非白名单的访问）
		recordAccess(ip)

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
	expireAt := time.Now().Add(time.Duration(banDuration) * time.Hour)

	blacklist := model.IPBlacklist{
		IP:       ip,
		Reason:   reason,
		BanType:  banType,
		ExpireAt: &expireAt,
	}

	// 使用 FirstOrCreate 避免重复插入
	db.DB.Where("ip = ?", ip).FirstOrCreate(&blacklist)

	// 标记为已封禁
	ipMapMutex.Lock()
	if record, exists := ipAccessMap[ip]; exists {
		record.banned = true
	}
	ipMapMutex.Unlock()
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
		return role == "admin"
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

	// 检查是否是管理员
	return claims.Role == "admin"
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
