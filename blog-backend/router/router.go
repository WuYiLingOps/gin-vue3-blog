/*
 * 项目名称：blog-backend
 * 文件名称：router.go
 * 创建时间：2026-01-31 16:45:02
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：路由配置模块，负责配置所有API路由、中间件和静态文件服务
 */
package router

import (
	"blog-backend/constant"
	"blog-backend/handler"
	"blog-backend/middleware"
	"blog-backend/service"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置并返回Gin路由引擎
// 功能说明：
//  1. 初始化Gin引擎并配置中间件
//  2. 配置静态文件服务
//  3. 初始化WebSocket Hub
//  4. 初始化所有处理器
//  5. 配置所有API路由组
//
// 返回:
//   - *gin.Engine: 配置好的Gin路由引擎
func SetupRouter() *gin.Engine {
	r := gin.New()

	// 使用中间件（按顺序执行）
	r.Use(gin.Recovery())                     // Gin内置恢复中间件，捕获panic
	r.Use(middleware.IPContextMiddleware())   // IP上下文中间件（最先执行，确保IP可用）
	r.Use(middleware.Logger())                // HTTP请求日志中间件
	r.Use(middleware.CORS())                  // 跨域资源共享中间件
	r.Use(middleware.IPBlacklistMiddleware()) // IP黑名单和频率限制中间件

	// 静态文件服务（用于访问上传的文件）
	// 使用绝对路径，确保无论从哪个目录运行都能找到 uploads 目录
	uploadsPath, _ := filepath.Abs("./uploads")
	r.Static("/uploads", uploadsPath)

	// 初始化WebSocket Hub（用于实时聊天功能）
	chatHub := service.NewHub()
	go chatHub.Run() // 启动Hub，在后台goroutine中运行

	// 初始化所有业务处理器
	authHandler := handler.NewAuthHandler()
	postHandler := handler.NewPostHandler()
	categoryHandler := handler.NewCategoryHandler()
	tagHandler := handler.NewTagHandler()
	commentHandler := handler.NewCommentHandler()
	userHandler := handler.NewUserHandler()
	uploadHandler := handler.NewUploadHandler()
	settingHandler := handler.NewSettingHandler()
	dashboardHandler := handler.NewDashboardHandler()
	momentHandler := handler.NewMomentHandler()
	ipBlacklistHandler := handler.NewIPBlacklistHandler()
	ipWhitelistHandler := handler.NewIPWhitelistHandler()
	captchaHandler := handler.NewCaptchaHandler()
	chatHandler := handler.NewChatHandler(chatHub)
	blogHandler := handler.NewBlogHandler()
	announcementHandler := handler.NewAnnouncementHandler()
	friendLinkHandler := handler.NewFriendLinkHandler()
	friendLinkCategoryHandler := handler.NewFriendLinkCategoryHandler()
	calendarHandler := handler.NewCalendarHandler()
	albumHandler := handler.NewAlbumHandler()

	// 健康检查接口（用于服务监控和负载均衡器健康检查）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API 路由组（所有API接口都以 /api 为前缀）
	api := r.Group("/api")
	{
		// 配置各个功能模块的路由
		setupAuthRoutes(api, authHandler)                                                                                                                                                                                                 // 认证相关路由
		setupCaptchaRoutes(api, captchaHandler)                                                                                                                                                                                           // 验证码路由
		setupBlogRoutes(api, blogHandler, announcementHandler, friendLinkHandler, friendLinkCategoryHandler, albumHandler)                                                                                                                // 博客公开接口路由
		setupCalendarRoutes(api, calendarHandler)                                                                                                                                                                                         // 日历路由
		setupPostRoutes(api, postHandler)                                                                                                                                                                                                 // 文章路由
		setupCategoryRoutes(api, categoryHandler)                                                                                                                                                                                         // 分类路由
		setupTagRoutes(api, tagHandler)                                                                                                                                                                                                   // 标签路由
		setupCommentRoutes(api, commentHandler)                                                                                                                                                                                           // 评论路由
		setupUploadRoutes(api, uploadHandler)                                                                                                                                                                                             // 文件上传路由
		setupSettingRoutes(api, settingHandler)                                                                                                                                                                                           // 系统设置路由
		setupMomentRoutes(api, momentHandler)                                                                                                                                                                                             // 说说路由
		setupChatRoutes(api, chatHandler)                                                                                                                                                                                                 // 聊天室路由
		setupAdminRoutes(api, userHandler, postHandler, commentHandler, dashboardHandler, momentHandler, ipBlacklistHandler, ipWhitelistHandler, chatHandler, friendLinkHandler, friendLinkCategoryHandler, settingHandler, albumHandler) // 管理后台路由
	}

	return r
}

// setupAuthRoutes 配置认证相关路由
// 功能说明：配置用户注册、登录、登出、密码重置、个人信息管理等路由
// 参数:
//   - api: API路由组
//   - h: 认证处理器实例
func setupAuthRoutes(api *gin.RouterGroup, h *handler.AuthHandler) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/send-register-code", h.SendRegisterCode) // 发送注册验证码
		auth.POST("/login", h.Login)
		auth.POST("/logout", h.Logout)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/forgot-password", h.ForgotPassword) // 忘记密码 - 发送验证码
		auth.POST("/reset-password", h.ResetPassword)   // 重置密码

		// 需要认证的接口
		authRequired := auth.Group("")
		authRequired.Use(middleware.AuthMiddleware())
		{
			authRequired.GET("/profile", h.GetProfile)
			authRequired.PUT("/profile", h.UpdateProfile)
			authRequired.PUT("/password", h.UpdatePassword)
			authRequired.PUT("/email", h.UpdateEmail)                    // 修改邮箱
			authRequired.GET("/email-change-info", h.GetEmailChangeInfo) // 获取邮箱修改信息
		}
	}
}

// setupCaptchaRoutes 配置验证码路由
// 功能说明：配置图形验证码获取接口
// 参数:
//   - api: API路由组
//   - h: 验证码处理器实例
func setupCaptchaRoutes(api *gin.RouterGroup, h *handler.CaptchaHandler) {
	captcha := api.Group("/captcha")
	{
		captcha.GET("", h.GetCaptcha)
	}
}

// setupBlogRoutes 配置博客公开接口路由
// 功能说明：配置博主信息、网站统计、公告、友链、相册等公开接口
// 参数:
//   - api: API路由组
//   - h: 博客处理器实例
//   - a: 公告处理器实例
//   - fl: 友链处理器实例
//   - flc: 友链分类处理器实例
//   - al: 相册处理器实例
func setupBlogRoutes(api *gin.RouterGroup, h *handler.BlogHandler, a *handler.AnnouncementHandler, fl *handler.FriendLinkHandler, flc *handler.FriendLinkCategoryHandler, al *handler.AlbumHandler) {
	blog := api.Group("/blog")
	{
		// 获取博主资料和统计数据
		blog.GET("/author", h.GetAuthorProfile)
		// 关于我信息（公开）
		blog.GET("/about", h.GetAboutInfo)
		// 统计接口（公开）
		blog.GET("/tag-stats", h.GetTagStats)
		// 网站资讯（公开）
		blog.GET("/website-info", h.GetWebsiteInfo)
		// 公告/系统广播
		blog.GET("/announcements", a.GetAnnouncements)
		blog.GET("/announcements/:id", a.GetAnnouncementDetail)
		// 友链（公开接口）
		blog.GET("/friend-links", fl.ListPublic)
		blog.GET("/friend-link-categories", flc.List) // 公开获取分类列表
		// 相册（公开接口）
		blog.GET("/albums", al.ListPublic)
	}
}

// setupCalendarRoutes 配置贡献热力图路由（公开接口）
// 功能说明：配置Gitee贡献热力图数据查询接口
// 参数:
//   - api: API路由组
//   - h: 日历处理器实例
func setupCalendarRoutes(api *gin.RouterGroup, h *handler.CalendarHandler) {
	calendar := api.Group("/calendar")
	{
		calendar.GET("/gitee", h.GetGiteeCalendar)
	}
}

// setupPostRoutes 配置文章路由
// 功能说明：配置文章的增删改查、点赞、归档等路由，支持公开接口和需要认证的接口
// 参数:
//   - api: API路由组
//   - h: 文章处理器实例
func setupPostRoutes(api *gin.RouterGroup, h *handler.PostHandler) {
	posts := api.Group("/posts")
	// 前台获取文章相关接口允许携带可选认证信息，用于区分管理员和普通用户
	posts.Use(middleware.OptionalAuthMiddleware())
	{
		// 公开接口
		posts.GET("", h.List)
		posts.GET("/:id", h.GetByID)
		posts.GET("/archives", h.GetArchives)
		posts.GET("/hot", h.GetHotPosts)
		posts.GET("/recent", h.GetRecentPosts)
		posts.POST("/:id/like", h.Like)

		// 需要认证的接口
		postsAuth := posts.Group("")
		postsAuth.Use(middleware.AuthMiddleware())
		{
			postsAuth.POST("", h.Create)
			postsAuth.PUT("/:id", h.Update)
			postsAuth.DELETE("/:id", h.Delete)
		}
	}
}

// setupCategoryRoutes 配置分类路由
// 功能说明：配置文章分类的增删改查路由，查询接口公开，管理接口需要管理员权限
// 参数:
//   - api: API路由组
//   - h: 分类处理器实例
func setupCategoryRoutes(api *gin.RouterGroup, h *handler.CategoryHandler) {
	categories := api.Group("/categories")
	{
		// 公开接口
		categories.GET("", h.List)
		categories.GET("/:id", h.GetByID)

		// 需要管理员权限的接口
		categoriesAdmin := categories.Group("")
		categoriesAdmin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			categoriesAdmin.POST("", h.Create)
			categoriesAdmin.PUT("/:id", h.Update)
			categoriesAdmin.DELETE("/:id", h.Delete)
		}
	}
}

// setupTagRoutes 配置标签路由
// 功能说明：配置文章标签的增删改查路由，查询接口公开，管理接口需要认证
// 参数:
//   - api: API路由组
//   - h: 标签处理器实例
func setupTagRoutes(api *gin.RouterGroup, h *handler.TagHandler) {
	tags := api.Group("/tags")
	{
		// 公开接口
		tags.GET("", h.List)
		tags.GET("/:id", h.GetByID)
		tags.GET("/:id/posts", h.GetPostsByTag)

		// 需要认证的接口
		tagsAuth := tags.Group("")
		tagsAuth.Use(middleware.AuthMiddleware())
		{
			tagsAuth.POST("", h.Create)
			tagsAuth.PUT("/:id", h.Update)
			tagsAuth.DELETE("/:id", h.Delete)
		}
	}
}

// setupCommentRoutes 配置评论路由
// 功能说明：配置评论的增删改查路由，查询接口公开，创建和管理接口需要认证
// 参数:
//   - api: API路由组
//   - h: 评论处理器实例
func setupCommentRoutes(api *gin.RouterGroup, h *handler.CommentHandler) {
	comments := api.Group("/comments")
	{
		// 公开接口
		comments.GET("/post/:id", h.GetByPostID)
		comments.GET("/type", h.GetByTypeAndTarget) // 根据类型和目标ID获取评论（用于友链等特殊页面）

		// 需要认证的接口
		commentsAuth := comments.Group("")
		commentsAuth.Use(middleware.AuthMiddleware())
		{
			commentsAuth.POST("", h.Create)
			commentsAuth.PUT("/:id", h.Update)
			commentsAuth.DELETE("/:id", h.Delete)
		}
	}
}

// setupUploadRoutes 配置文件上传路由
// 功能说明：配置头像和图片上传接口，需要认证
// 参数:
//   - api: API路由组
//   - h: 文件上传处理器实例
func setupUploadRoutes(api *gin.RouterGroup, h *handler.UploadHandler) {
	upload := api.Group("/upload")
	upload.Use(middleware.AuthMiddleware())
	{
		upload.POST("/avatar", h.UploadAvatar)
		upload.POST("/image", h.UploadImage)
	}
}

// setupSettingRoutes 配置系统设置路由
// 功能说明：配置系统设置的查询和更新路由，部分接口公开，管理接口需要管理员权限
// 参数:
//   - api: API路由组
//   - h: 系统设置处理器实例
func setupSettingRoutes(api *gin.RouterGroup, h *handler.SettingHandler) {
	settings := api.Group("/settings")
	{
		// 公开接口
		settings.GET("/public", h.GetPublicSettings)
		settings.GET("/friendlink-info", h.GetFriendLinkInfo)

		// 需要超级管理员（super_admin）权限：系统级配置
		settingsAdmin := settings.Group("")
		settingsAdmin.Use(middleware.AuthMiddleware(), middleware.RoleRequiredMiddleware(constant.RoleSuperAdmin))
		{
			settingsAdmin.GET("/site", h.GetSiteSettings)
			settingsAdmin.PUT("/site", h.UpdateSiteSettings)
			settingsAdmin.GET("/upload", h.GetUploadSettings)
			settingsAdmin.PUT("/upload", h.UpdateUploadSettings)
			settingsAdmin.GET("/notification", h.GetNotificationSettings)
			settingsAdmin.PUT("/notification", h.UpdateNotificationSettings)
			settingsAdmin.GET("/register", h.GetRegisterSettings)
			settingsAdmin.PUT("/register", h.UpdateRegisterSettings)
			settingsAdmin.PUT("/friendlink-info", h.UpdateFriendLinkInfo)
		}
	}
}

// setupMomentRoutes 配置说说路由
// 功能说明：配置说说的增删改查、点赞等路由，查询接口支持可选认证，管理接口需要认证
// 参数:
//   - api: API路由组
//   - h: 说说处理器实例
func setupMomentRoutes(api *gin.RouterGroup, h *handler.MomentHandler) {
	moments := api.Group("/moments")
	{
		// 公开接口
		moments.GET("", middleware.OptionalAuthMiddleware(), h.List)
		moments.GET("/:id", middleware.OptionalAuthMiddleware(), h.GetByID)
		moments.GET("/recent", h.GetRecent)
		moments.POST("/:id/like", h.Like)

		// 需要认证的接口
		momentsAuth := moments.Group("")
		momentsAuth.Use(middleware.AuthMiddleware())
		{
			momentsAuth.POST("", h.Create)
			momentsAuth.PUT("/:id", h.Update)
			momentsAuth.DELETE("/:id", h.Delete)
		}
	}
}

// setupChatRoutes 配置聊天室路由
// 功能说明：配置WebSocket连接、消息查询、在线信息等路由，支持认证和匿名访问
// 参数:
//   - api: API路由组
//   - h: 聊天室处理器实例
func setupChatRoutes(api *gin.RouterGroup, h *handler.ChatHandler) {
	chat := api.Group("/chat")
	{
		// WebSocket连接（支持认证和匿名，使用可选认证中间件）
		chat.GET("/ws", middleware.OptionalAuthMiddleware(), h.HandleWebSocket)

		// 公开接口
		chat.GET("/messages", h.GetMessages)
		chat.GET("/online", h.GetOnlineInfo)
		chat.GET("/settings", h.GetChatSettings)
	}
}

// setupAdminRoutes 配置管理后台路由
// 功能说明：配置所有管理后台功能的路由，包括仪表盘、用户管理、文章管理、评论管理、IP管理、聊天室管理等
// 所有接口都需要管理员权限（AuthMiddleware + AdminMiddleware）
// 参数:
//   - api: API路由组
//   - userHandler: 用户处理器实例
//   - postHandler: 文章处理器实例
//   - commentHandler: 评论处理器实例
//   - dashboardHandler: 仪表盘处理器实例
//   - momentHandler: 说说处理器实例
//   - ipBlacklistHandler: IP黑名单处理器实例
//   - ipWhitelistHandler: IP白名单处理器实例
//   - chatHandler: 聊天室处理器实例
//   - friendLinkHandler: 友链处理器实例
//   - friendLinkCategoryHandler: 友链分类处理器实例
//   - settingHandler: 系统设置处理器实例
//   - albumHandler: 相册处理器实例
func setupAdminRoutes(api *gin.RouterGroup, userHandler *handler.UserHandler, postHandler *handler.PostHandler, commentHandler *handler.CommentHandler, dashboardHandler *handler.DashboardHandler, momentHandler *handler.MomentHandler, ipBlacklistHandler *handler.IPBlacklistHandler, ipWhitelistHandler *handler.IPWhitelistHandler, chatHandler *handler.ChatHandler, friendLinkHandler *handler.FriendLinkHandler, friendLinkCategoryHandler *handler.FriendLinkCategoryHandler, settingHandler *handler.SettingHandler, albumHandler *handler.AlbumHandler) {
	admin := api.Group("/admin")
	// admin 路由基础权限：admin 或 super_admin
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		// 仪表盘
		admin.GET("/dashboard/stats", dashboardHandler.GetStats)
		admin.GET("/dashboard/category-stats", dashboardHandler.GetCategoryStats)
		admin.GET("/dashboard/visit-stats", dashboardHandler.GetVisitStats)

		// super_admin 专属路由组：系统级高危操作（用户管理/关于我/友链/相册等）
		super := admin.Group("")
		super.Use(middleware.RoleRequiredMiddleware(constant.RoleSuperAdmin))
		{
			// 用户管理（仅超级管理员）
			super.GET("/users", userHandler.List)
			super.GET("/users/:id", userHandler.GetByID)
			super.PUT("/users/:id/status", userHandler.UpdateStatus)
			super.PUT("/users/:id/role", userHandler.UpdateRole) // 更新用户角色
			super.DELETE("/users/:id", userHandler.Delete)

			// 关于我信息管理（仅超级管理员）
			super.GET("/about", settingHandler.GetAboutInfo)
			super.PUT("/about", settingHandler.UpdateAboutInfo)

			// 友链管理（仅超级管理员）
			super.GET("/friend-links", friendLinkHandler.List)
			super.GET("/friend-links/:id", friendLinkHandler.GetByID)
			super.POST("/friend-links", friendLinkHandler.Create)
			super.PUT("/friend-links/:id", friendLinkHandler.Update)
			super.DELETE("/friend-links/:id", friendLinkHandler.Delete)

			// 友链分类管理（仅超级管理员）
			super.GET("/friend-link-categories", friendLinkCategoryHandler.List)
			super.GET("/friend-link-categories/:id", friendLinkCategoryHandler.GetByID)
			super.POST("/friend-link-categories", friendLinkCategoryHandler.Create)
			super.PUT("/friend-link-categories/:id", friendLinkCategoryHandler.Update)
			super.DELETE("/friend-link-categories/:id", friendLinkCategoryHandler.Delete)

			// 相册管理（仅超级管理员）
			super.GET("/albums", albumHandler.List)
			super.GET("/albums/:id", albumHandler.GetByID)
			super.POST("/albums", albumHandler.Create)
			super.PUT("/albums/:id", albumHandler.Update)
			super.DELETE("/albums/:id", albumHandler.Delete)

			// 注册设置管理（仅超级管理员）
			super.GET("/settings/register", settingHandler.GetRegisterSettings)
			super.PUT("/settings/register", settingHandler.UpdateRegisterSettings)
		}

		// 文章管理
		admin.GET("/posts", postHandler.List)
		admin.GET("/posts/:id/export", postHandler.Export)

		// 评论管理
		admin.GET("/comments", commentHandler.List)
		admin.PUT("/comments/:id/status", commentHandler.UpdateStatus)

		// 说说管理
		admin.GET("/moments", momentHandler.AdminList)

		// IP黑名单管理
		admin.GET("/ip-blacklist", ipBlacklistHandler.List)
		admin.POST("/ip-blacklist", ipBlacklistHandler.Add)
		admin.DELETE("/ip-blacklist/:id", ipBlacklistHandler.Delete)
		admin.GET("/ip-blacklist/check", ipBlacklistHandler.Check)
		admin.POST("/ip-blacklist/clean-expired", ipBlacklistHandler.CleanExpired)

		// IP白名单管理
		admin.GET("/ip-whitelist", ipWhitelistHandler.List)
		admin.POST("/ip-whitelist", ipWhitelistHandler.Add)
		admin.DELETE("/ip-whitelist/:id", ipWhitelistHandler.Delete)
		admin.GET("/ip-whitelist/check", ipWhitelistHandler.Check)
		admin.POST("/ip-whitelist/clean-expired", ipWhitelistHandler.CleanExpired)

		// 聊天室管理
		admin.GET("/chat/messages", chatHandler.AdminListMessages)
		admin.DELETE("/chat/messages/:id", chatHandler.DeleteMessage)
		admin.POST("/chat/broadcast", chatHandler.BroadcastSystemMessage)
		admin.POST("/chat/kick", chatHandler.KickUser) // 踢出用户
		admin.POST("/chat/ban", chatHandler.BanIP)     // 封禁IP
		admin.GET("/chat/settings", chatHandler.GetChatSettings)
		admin.PUT("/chat/settings", chatHandler.UpdateChatSettings)
	}
}
