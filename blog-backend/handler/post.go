package handler

import (
	"strconv"

	"blog-backend/service"
	"blog-backend/util"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	service *service.PostService
}

func NewPostHandler() *PostHandler {
	return &PostHandler{
		service: service.NewPostService(),
	}
}

// Create 创建文章
func (h *PostHandler) Create(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		util.Unauthorized(c, "未登录")
		return
	}

	var req service.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	post, err := h.service.Create(userID.(uint), &req)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "文章创建成功", post)
}

// GetByID 获取文章详情
func (h *PostHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的文章ID")
		return
	}

	// 获取用户ID（如果已登录）
	var userID *uint
	if uid, exists := c.Get("user_id"); exists {
		id := uid.(uint)
		userID = &id
	}

	// 获取客户端IP
	ip := util.GetClientIP(c)

	post, err := h.service.GetByID(uint(id), userID, ip)
	if err != nil {
		util.Error(c, 404, err.Error())
		return
	}

	util.Success(c, post)
}

// Update 更新文章
func (h *PostHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的文章ID")
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var req service.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	post, err := h.service.Update(uint(id), userID.(uint), role.(string), &req)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "文章更新成功", post)
}

// Delete 删除文章
func (h *PostHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的文章ID")
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	if err := h.service.Delete(uint(id), userID.(uint), role.(string)); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "文章删除成功", nil)
}

// List 获取文章列表
func (h *PostHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	categoryID, _ := strconv.ParseUint(c.Query("category_id"), 10, 32)
	keyword := c.Query("keyword")
	status, _ := strconv.Atoi(c.DefaultQuery("status", "1"))

	posts, total, err := h.service.List(page, pageSize, uint(categoryID), keyword, status)
	if err != nil {
		util.ServerError(c, "获取文章列表失败")
		return
	}

	util.PageSuccess(c, posts, total, page, pageSize)
}

// GetByTag 根据标签获取文章
func (h *PostHandler) GetByTag(c *gin.Context) {
	tagID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的标签ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	posts, total, err := h.service.GetByTag(uint(tagID), page, pageSize)
	if err != nil {
		util.ServerError(c, "获取文章列表失败")
		return
	}

	util.PageSuccess(c, posts, total, page, pageSize)
}

// Like 点赞文章
func (h *PostHandler) Like(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的文章ID")
		return
	}

	if err := h.service.Like(uint(id)); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "点赞成功", nil)
}

// GetArchives 获取归档
func (h *PostHandler) GetArchives(c *gin.Context) {
	archives, err := h.service.GetArchives()
	if err != nil {
		util.ServerError(c, "获取归档失败")
		return
	}

	util.Success(c, archives)
}

// GetHotPosts 获取热门文章
func (h *PostHandler) GetHotPosts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	posts, err := h.service.GetHotPosts(limit)
	if err != nil {
		util.ServerError(c, "获取热门文章失败")
		return
	}

	util.Success(c, posts)
}

// GetRecentPosts 获取最新文章
func (h *PostHandler) GetRecentPosts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	posts, err := h.service.GetRecentPosts(limit)
	if err != nil {
		util.ServerError(c, "获取最新文章失败")
		return
	}

	util.Success(c, posts)
}

