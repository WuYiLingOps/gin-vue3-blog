package main

import (
	"blog-backend/handler"
	"blog-backend/middleware"
	"github.com/gin-gonic/gin"
)

// setupRouter 配置路由
func setupRouter() *gin.Engine {
	r := gin.New()

	// 使用中间件
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(middleware.IPBlacklistMiddleware()) // IP黑名单和频率限制
	r.Use(middleware.VisitStatMiddleware())

	// 静态文件服务（用于访问上传的文件）
	r.Static("/uploads", "./uploads")

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

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API 路由组
	api := r.Group("/api")
	{
		setupAuthRoutes(api, authHandler)
		setupPostRoutes(api, postHandler)
		setupCategoryRoutes(api, categoryHandler)
		setupTagRoutes(api, tagHandler)
		setupCommentRoutes(api, commentHandler)
		setupUploadRoutes(api, uploadHandler)
		setupSettingRoutes(api, settingHandler)
		setupMomentRoutes(api, momentHandler)
		setupAdminRoutes(api, userHandler, postHandler, commentHandler, dashboardHandler, momentHandler, ipBlacklistHandler)
	}

	return r
}

// setupAuthRoutes 认证路由
func setupAuthRoutes(api *gin.RouterGroup, h *handler.AuthHandler) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/logout", h.Logout)
		auth.POST("/refresh", h.RefreshToken)

		// 需要认证的接口
		authRequired := auth.Group("")
		authRequired.Use(middleware.AuthMiddleware())
		{
			authRequired.GET("/profile", h.GetProfile)
			authRequired.PUT("/profile", h.UpdateProfile)
			authRequired.PUT("/password", h.UpdatePassword)
		}
	}
}

// setupPostRoutes 文章路由
func setupPostRoutes(api *gin.RouterGroup, h *handler.PostHandler) {
	posts := api.Group("/posts")
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
		settings.GET("/about", h.GetAboutSettings)
		settings.GET("/public", h.GetPublicSettings)

		// 需要管理员权限
		settingsAdmin := settings.Group("")
		settingsAdmin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			settingsAdmin.PUT("/about", h.UpdateAboutSettings)
			settingsAdmin.GET("/site", h.GetSiteSettings)
			settingsAdmin.PUT("/site", h.UpdateSiteSettings)
		}
	}
}

// setupMomentRoutes 说说路由
func setupMomentRoutes(api *gin.RouterGroup, h *handler.MomentHandler) {
	moments := api.Group("/moments")
	{
		// 公开接口
		moments.GET("", h.List)
		moments.GET("/:id", h.GetByID)
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

// setupAdminRoutes 管理后台路由
func setupAdminRoutes(api *gin.RouterGroup, userHandler *handler.UserHandler, postHandler *handler.PostHandler, commentHandler *handler.CommentHandler, dashboardHandler *handler.DashboardHandler, momentHandler *handler.MomentHandler, ipBlacklistHandler *handler.IPBlacklistHandler) {
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
	}
}

