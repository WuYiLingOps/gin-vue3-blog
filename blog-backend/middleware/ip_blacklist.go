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
		// 排除静态文件路径和 WebSocket 路径，不进行频率限制和黑名单检查
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/uploads") || strings.HasPrefix(path, "/api/chat/ws") {
			c.Next()
			return
		}

		ip := util.GetClientIP(c)

		// 0. 优先检查：如果是管理员或白名单IP，跳过所有限制
		if shouldSkipRateLimit(c, ip) {
			c.Next()
			return
		}

		// 1. 检查是否在黑名单中
		if isIPBanned(ip) {
			util.Error(c, 403, "您的IP已被封禁，请联系管理员")
			c.Abort()
			return
		}

		// 2. 检查访问频率
		if shouldBanIP(ip) {
			// 自动封禁IP
			banIP(ip, "访问频率过高，自动封禁", 1)
			util.Error(c, 429, "访问过于频繁，您的IP已被临时封禁")
			c.Abort()
			return
		}

		// 3. 记录访问
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
