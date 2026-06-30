/*
 * 项目名称：blog-backend
 * 文件名称：page_view.go
 * 创建时间：2026-06-30
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：全站PV统计中间件，记录所有页面访问
 */
package middleware

import (
	"blog-backend/repository"
	"blog-backend/util"
	"strings"

	"github.com/gin-gonic/gin"
)

// 需要排除的路径前缀（不计入PV）
var excludePaths = []string{
	"/api/",
	"/admin/dashboard/",
	"/uploads/",
	"/assets/",
	"/favicon.ico",
}

// PageView 全站PV统计中间件
// 功能说明：记录所有页面访问，按天去重（同IP同路径每天只计1次）
func PageView() gin.HandlerFunc {
	pageViewRepo := repository.NewPageViewRepository()

	return func(c *gin.Context) {
		// 继续执行后续中间件和处理器
		c.Next()

		// 只统计GET请求
		if c.Request.Method != "GET" {
			return
		}

		// 只统计200状态码
		if c.Writer.Status() != 200 {
			return
		}

		path := c.Request.URL.Path

		// 排除静态资源和API路径
		for _, exclude := range excludePaths {
			if strings.HasPrefix(path, exclude) {
				return
			}
		}

		// 获取客户端IP
		ip := util.GetClientIP(c)
		if ip == "" || ip == "unknown" {
			return
		}

		// 获取User-Agent
		userAgent := c.Request.UserAgent()
		if len(userAgent) > 500 {
			userAgent = userAgent[:500]
		}

		// 获取用户ID（如果已登录）
		var userID *uint
		if uid, exists := c.Get("user_id"); exists {
			if id, ok := uid.(uint); ok {
				userID = &id
			}
		}

		// 检查今天是否已经访问过
		hasViewed, _ := pageViewRepo.HasViewedToday(path, ip)
		if hasViewed {
			return
		}

		// 记录访问
		_ = pageViewRepo.RecordView(path, userID, ip, userAgent)
	}
}
