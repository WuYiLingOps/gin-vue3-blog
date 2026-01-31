/*
 * 项目名称：blog-backend
 * 文件名称：rate_limit.go
 * 创建时间：2026-01-31 16:17:29
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：限流中间件，基于IP地址实现请求频率限制，防止恶意请求
 */
package middleware

import (
	"sync"
	"time"

	"blog-backend/util"

	"github.com/gin-gonic/gin"
)

// visitor 访客访问记录结构体
type visitor struct {
	lastSeen time.Time // 最后访问时间
	count    int       // 访问次数
}

var (
	// visitors 访客访问记录映射表，key为IP地址
	visitors = make(map[string]*visitor)
	// mu 互斥锁，保护visitors映射的并发访问
	mu sync.Mutex
)

// RateLimit 限流中间件
// 功能说明：基于IP地址实现请求频率限制，在指定时间窗口内限制每个IP的最大请求数
// 参数:
//   - maxRequests: 时间窗口内允许的最大请求数
//   - window: 时间窗口大小
//
// 返回:
//   - gin.HandlerFunc: Gin中间件处理函数
func RateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
	// 启动后台协程，定期清理过期的访客记录
	go func() {
		for {
			time.Sleep(window)
			mu.Lock()
			now := time.Now()
			for ip, v := range visitors {
				// 如果距离最后访问时间超过时间窗口，删除记录
				if now.Sub(v.lastSeen) > window {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		// 获取客户端IP地址
		ip := c.ClientIP()

		mu.Lock()
		v, exists := visitors[ip]
		if !exists {
			// 首次访问，创建新记录
			visitors[ip] = &visitor{lastSeen: time.Now(), count: 1}
			mu.Unlock()
			c.Next()
			return
		}

		// 检查是否超过时间窗口
		if time.Since(v.lastSeen) > window {
			// 超过时间窗口，重置计数
			v.count = 1
			v.lastSeen = time.Now()
			mu.Unlock()
			c.Next()
			return
		}

		// 检查是否超过最大请求数
		if v.count >= maxRequests {
			mu.Unlock()
			util.Error(c, 429, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		// 增加访问计数并更新最后访问时间
		v.count++
		v.lastSeen = time.Now()
		mu.Unlock()

		c.Next()
	}
}
