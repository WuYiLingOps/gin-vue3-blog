package middleware

import (
	"sync"
	"time"

	"blog-backend/util"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	lastSeen time.Time
	count    int
}

var visitors = make(map[string]*visitor)
var mu sync.Mutex

// RateLimit 限流中间件
func RateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
	// 清理过期的访客记录
	go func() {
		for {
			time.Sleep(window)
			mu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastSeen) > window {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		v, exists := visitors[ip]
		if !exists {
			visitors[ip] = &visitor{lastSeen: time.Now(), count: 1}
			mu.Unlock()
			c.Next()
			return
		}

		if time.Since(v.lastSeen) > window {
			v.count = 1
			v.lastSeen = time.Now()
			mu.Unlock()
			c.Next()
			return
		}

		if v.count >= maxRequests {
			mu.Unlock()
			util.Error(c, 429, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		v.count++
		v.lastSeen = time.Now()
		mu.Unlock()

		c.Next()
	}
}
