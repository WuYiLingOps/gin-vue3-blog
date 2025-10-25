package middleware

import (
	"blog-backend/repository"
	"github.com/gin-gonic/gin"
)

// VisitStatMiddleware 访问量统计中间件
func VisitStatMiddleware() gin.HandlerFunc {
	visitStatRepo := repository.NewVisitStatRepository()

	return func(c *gin.Context) {
		// 先处理请求
		c.Next()

		// 只在成功响应时统计（200-299状态码）
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			// 异步记录访问量，不影响响应速度
			go func() {
				visitStatRepo.IncrementTodayViewCount()
			}()
		}
	}
}

