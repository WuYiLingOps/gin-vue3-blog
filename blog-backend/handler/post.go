package handler

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"blog-backend/model"
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

	// 确保 status 字段正确传递（0 或 1）
	// 如果前端没有传递 status，默认为 1（发布）
	if req.Status != 0 && req.Status != 1 {
		req.Status = 1
	}

	// 处理可见性：如果前端传递了 visibility 字段（包括 0），确保正确接收
	// 由于 JSON 中 0 是有效值，如果前端传递了 visibility: 0，应该能正确接收
	// 但如果前端没有传递 visibility 字段，req.Visibility 会是 nil，此时使用默认值 1
	// 这里不需要额外处理，因为 service 层会处理 nil 的情况

	post, err := h.service.Create(userID.(uint), &req)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	if post == nil {
		util.Error(c, 500, "文章创建失败：返回数据为空")
		return
	}

	util.SuccessWithMessage(c, "文章创建成功", post)
}

// GetByID 获取文章详情（支持ID或slug）
func (h *PostHandler) GetByID(c *gin.Context) {
	identifier := c.Param("id")

	// 获取用户ID（如果已登录）
	var userID *uint
	var role string
	if uid, exists := c.Get("user_id"); exists {
		id := uid.(uint)
		userID = &id
	}
	if r, exists := c.Get("role"); exists {
		role = r.(string)
	}

	// 获取客户端IP
	ip := util.GetClientIP(c)

	var post *model.Post
	var err error

	// 尝试解析为数字ID，如果失败则作为slug处理
	if id, parseErr := strconv.ParseUint(identifier, 10, 32); parseErr == nil {
		// 是数字ID
		post, err = h.service.GetByID(uint(id), userID, role, ip)
	} else {
		// 是slug
		post, err = h.service.GetBySlug(identifier, userID, role, ip)
	}

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

	// 默认只返回公开文章；管理员则可以查看所有可见性
	var visibility *int
	if role, exists := c.Get("role"); !exists || role != "admin" {
		v := 1
		visibility = &v
	}

	posts, total, err := h.service.List(page, pageSize, uint(categoryID), keyword, status, visibility)
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

// Like 点赞/取消点赞文章
func (h *PostHandler) Like(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的文章ID")
		return
	}

	// 获取用户ID和IP
	var userID *uint
	if uid, exists := c.Get("user_id"); exists {
		id := uid.(uint)
		userID = &id
	}
	ip := util.GetClientIP(c)

	liked, err := h.service.Like(uint(id), userID, ip)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	// 返回点赞状态
	message := "点赞成功"
	if !liked {
		message = "取消点赞"
	}
	util.Success(c, gin.H{
		"liked":   liked,
		"message": message,
	})
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

	var userID *uint
	var role string
	if uid, exists := c.Get("user_id"); exists {
		id := uid.(uint)
		userID = &id
	}
	if r, exists := c.Get("role"); exists {
		role = r.(string)
	}

	posts, err := h.service.GetRecentPosts(limit, userID, role)
	if err != nil {
		util.ServerError(c, "获取最新文章失败")
		return
	}

	util.Success(c, posts)
}

// Export 导出文章为 Markdown
func (h *PostHandler) Export(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的文章ID")
		return
	}

	post, err := h.service.GetByIDForAdmin(uint(id))
	if err != nil {
		util.Error(c, 404, "文章不存在")
		return
	}

	var buf bytes.Buffer
	buf.WriteString("---\n")
	buf.WriteString(fmt.Sprintf("title: \"%s\"\n", escapeYAML(post.Title)))
	if post.PublishedAt != nil {
		buf.WriteString(fmt.Sprintf("date: %s\n", post.PublishedAt.Format(time.RFC3339)))
	} else {
		buf.WriteString(fmt.Sprintf("date: %s\n", post.CreatedAt.Format(time.RFC3339)))
	}
	buf.WriteString(fmt.Sprintf("status: %d\n", post.Status))
	buf.WriteString(fmt.Sprintf("category: \"%s\"\n", escapeYAML(post.Category.Name)))

	if len(post.Tags) > 0 {
		buf.WriteString("tags:\n")
		for _, t := range post.Tags {
			buf.WriteString(fmt.Sprintf("  - \"%s\"\n", escapeYAML(t.Name)))
		}
	} else {
		buf.WriteString("tags: []\n")
	}

	buf.WriteString("---\n\n")
	buf.WriteString(post.Content)

	filename := sanitizeFilename(post.Title)
	if filename == "" {
		filename = fmt.Sprintf("post-%d", post.ID)
	}

	c.Header("Content-Type", "text/markdown; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.md\"", filename))
	c.String(200, buf.String())
}

// escapeYAML 简单转义引号
func escapeYAML(s string) string {
	return strings.ReplaceAll(s, "\"", "\\\"")
}

// sanitizeFilename 将标题转为安全文件名
func sanitizeFilename(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	re := regexp.MustCompile(`[^a-zA-Z0-9\\-_.]+`)
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = re.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}
