package router

import (
	"blog-backend/handler"
	"blog-backend/middleware"
	"blog-backend/service"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	r := gin.New()

	// 使用中间件
	r.Use(gin.Recovery())
	r.Use(middleware.IPContextMiddleware()) // IP上下文中间件（最先执行，确保IP可用）
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(middleware.IPBlacklistMiddleware()) // IP黑名单和频率限制

	// 静态文件服务（用于访问上传的文件）
	// 使用绝对路径，确保无论从哪个目录运行都能找到 uploads 目录
	uploadsPath, _ := filepath.Abs("./uploads")
	r.Static("/uploads", uploadsPath)

	// 初始化WebSocket Hub
	chatHub := service.NewHub()
	go chatHub.Run() // 启动Hub

	// 初始化处理器
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

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API 路由组
	api := r.Group("/api")
	{
		setupAuthRoutes(api, authHandler)
		setupCaptchaRoutes(api, captchaHandler)
		setupBlogRoutes(api, blogHandler, announcementHandler, friendLinkHandler, friendLinkCategoryHandler)
		setupCalendarRoutes(api, calendarHandler)
		setupPostRoutes(api, postHandler)
		setupCategoryRoutes(api, categoryHandler)
		setupTagRoutes(api, tagHandler)
		setupCommentRoutes(api, commentHandler)
		setupUploadRoutes(api, uploadHandler)
		setupSettingRoutes(api, settingHandler)
		setupMomentRoutes(api, momentHandler)
		setupChatRoutes(api, chatHandler)
		setupAdminRoutes(api, userHandler, postHandler, commentHandler, dashboardHandler, momentHandler, ipBlacklistHandler, ipWhitelistHandler, chatHandler, friendLinkHandler, friendLinkCategoryHandler, settingHandler)
	}

	return r
}

// setupAuthRoutes 认证路由
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

// setupCaptchaRoutes 验证码路由
func setupCaptchaRoutes(api *gin.RouterGroup, h *handler.CaptchaHandler) {
	captcha := api.Group("/captcha")
	{
		captcha.GET("", h.GetCaptcha)
	}
}

// setupBlogRoutes 博客路由（公开接口）
func setupBlogRoutes(api *gin.RouterGroup, h *handler.BlogHandler, a *handler.AnnouncementHandler, fl *handler.FriendLinkHandler, flc *handler.FriendLinkCategoryHandler) {
	blog := api.Group("/blog")
	{
		// 获取博主资料和统计数据
		blog.GET("/author", h.GetAuthorProfile)
		// 关于我信息（公开）
		blog.GET("/about", h.GetAboutInfo)
		// 统计接口（公开）
		blog.GET("/category-stats", h.GetCategoryStats)
		blog.GET("/tag-stats", h.GetTagStats)
		// 公告/系统广播
		blog.GET("/announcements", a.GetAnnouncements)
		blog.GET("/announcements/:id", a.GetAnnouncementDetail)
		// 友链（公开接口）
		blog.GET("/friend-links", fl.ListPublic)
		blog.GET("/friend-link-categories", flc.List) // 公开获取分类列表
	}
}

// setupCalendarRoutes 贡献热力图路由（公开接口）
func setupCalendarRoutes(api *gin.RouterGroup, h *handler.CalendarHandler) {
	calendar := api.Group("/calendar")
	{
		calendar.GET("/gitee", h.GetGiteeCalendar)
	}
}

// setupPostRoutes 文章路由
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

// setupCategoryRoutes 分类路由
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

// setupTagRoutes 标签路由
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

// setupCommentRoutes 评论路由
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

// setupUploadRoutes 上传路由
func setupUploadRoutes(api *gin.RouterGroup, h *handler.UploadHandler) {
	upload := api.Group("/upload")
	upload.Use(middleware.AuthMiddleware())
	{
		upload.POST("/avatar", h.UploadAvatar)
		upload.POST("/image", h.UploadImage)
	}
}

// setupSettingRoutes 配置路由
func setupSettingRoutes(api *gin.RouterGroup, h *handler.SettingHandler) {
	settings := api.Group("/settings")
	{
		// 公开接口
		settings.GET("/public", h.GetPublicSettings)
		settings.GET("/friendlink-info", h.GetFriendLinkInfo)

		// 需要管理员权限
		settingsAdmin := settings.Group("")
		settingsAdmin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
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

// setupMomentRoutes 说说路由
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

// setupChatRoutes 聊天室路由
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

// setupAdminRoutes 管理后台路由
func setupAdminRoutes(api *gin.RouterGroup, userHandler *handler.UserHandler, postHandler *handler.PostHandler, commentHandler *handler.CommentHandler, dashboardHandler *handler.DashboardHandler, momentHandler *handler.MomentHandler, ipBlacklistHandler *handler.IPBlacklistHandler, ipWhitelistHandler *handler.IPWhitelistHandler, chatHandler *handler.ChatHandler, friendLinkHandler *handler.FriendLinkHandler, friendLinkCategoryHandler *handler.FriendLinkCategoryHandler, settingHandler *handler.SettingHandler) {
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		// 仪表盘
		admin.GET("/dashboard/stats", dashboardHandler.GetStats)
		admin.GET("/dashboard/category-stats", dashboardHandler.GetCategoryStats)
		admin.GET("/dashboard/visit-stats", dashboardHandler.GetVisitStats)

		// 用户管理
		admin.GET("/users", userHandler.List)
		admin.GET("/users/:id", userHandler.GetByID)
		admin.PUT("/users/:id/status", userHandler.UpdateStatus)
		admin.DELETE("/users/:id", userHandler.Delete)

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

		// 友链管理
		// 友链管理
		admin.GET("/friend-links", friendLinkHandler.List)
		admin.GET("/friend-links/:id", friendLinkHandler.GetByID)
		admin.POST("/friend-links", friendLinkHandler.Create)
		admin.PUT("/friend-links/:id", friendLinkHandler.Update)
		admin.DELETE("/friend-links/:id", friendLinkHandler.Delete)

		// 友链分类管理
		admin.GET("/friend-link-categories", friendLinkCategoryHandler.List)
		admin.GET("/friend-link-categories/:id", friendLinkCategoryHandler.GetByID)
		admin.POST("/friend-link-categories", friendLinkCategoryHandler.Create)
		admin.PUT("/friend-link-categories/:id", friendLinkCategoryHandler.Update)
		admin.DELETE("/friend-link-categories/:id", friendLinkCategoryHandler.Delete)

		// 关于我信息管理
		admin.GET("/about", settingHandler.GetAboutInfo)
		admin.PUT("/about", settingHandler.UpdateAboutInfo)
	}
}
