package util

import (
	"net"

	"github.com/gin-gonic/gin"
)

// GetClientIP 获取客户端真实IP地址
func GetClientIP(c *gin.Context) string {
	// 优先从 X-Forwarded-For 获取（适用于使用了代理的情况）
	ip := c.GetHeader("EO-Connecting-IP")
	//ip := c.GetHeader("X-Forwarded-For")
	//fmt.Println("X-Forwarded-For:", ip)
	//if ip != "" {
	//	// X-Forwarded-For 可能包含多个IP，取第一个
	//	ips := strings.Split(ip, ",")
	//	if len(ips) > 0 {
	//		ip = strings.TrimSpace(ips[0])
	//		if isValidIP(ip) {
	//			return ip
	//		}
	//	}
	//}

	return ip
}

// isValidIP 验证IP地址是否有效（内部使用）
func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// IsValidIP 验证IP地址是否有效（导出供外部使用）
func IsValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}
