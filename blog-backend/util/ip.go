package util

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// ClientIPKey 上下文中的客户端IP键名
	ClientIPKey = "client_ip"
)

// GetClientIP 获取客户端真实IP地址
// 优先级：X-Forwarded-For -> X-Real-IP -> RemoteAddr -> 127.0.0.1（本地开发环境）
func GetClientIP(c *gin.Context) string {
	// 1. 优先从上下文获取（如果已通过中间件设置）
	if ip, exists := c.Get(ClientIPKey); exists {
		if ipStr, ok := ip.(string); ok && ipStr != "" {
			return ipStr
		}
	}

	// 2. 从 X-Forwarded-For 获取（适用于使用了代理的情况）
	ip := c.GetHeader("X-Forwarded-For")
	if ip != "" {
		// X-Forwarded-For 可能包含多个IP，取第一个
		ips := strings.Split(ip, ",")
		if len(ips) > 0 {
			ip = strings.TrimSpace(ips[0])
			if isValidIP(ip) {
				return ip
			}
		}
	}

	// 3. 从 X-Real-IP 获取
	ip = c.GetHeader("X-Real-IP")
	if ip != "" {
		ip = strings.TrimSpace(ip)
		if isValidIP(ip) {
			return ip
		}
	}

	// 4. 从 EO-Connecting-IP 获取（特定代理）
	ip = c.GetHeader("EO-Connecting-IP")
	if ip != "" {
		ip = strings.TrimSpace(ip)
		if isValidIP(ip) {
			return ip
		}
	}

	// 5. 从 RemoteAddr 获取（直接连接）
	if c.Request.RemoteAddr != "" {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err == nil && ip != "" {
			if isValidIP(ip) {
				return ip
			}
			// 即使验证失败，也返回IP（可能是IPv6或其他格式）
			return ip
		}
		// 如果 SplitHostPort 失败，可能是没有端口号，直接使用 RemoteAddr
		if isValidIP(c.Request.RemoteAddr) {
			return c.Request.RemoteAddr
		}
	}

	// 6. 本地开发环境默认返回 127.0.0.1
	return "127.0.0.1"
}

// isValidIP 验证IP地址是否有效（内部使用）
func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// IsValidIP 验证IP地址是否有效（导出供外部使用）
func IsValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}
