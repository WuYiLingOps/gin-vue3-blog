/*
 * 项目名称：blog-backend
 * 文件名称：post_collaborator.go
 * 创建时间：2026-06-16
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：文章协作者处理器，提供文章协作者的增删改查功能
 */
package handler

import (
	"blog-backend/constant"
	"blog-backend/repository"
	"blog-backend/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PostCollaboratorHandler 文章协作者处理器结构体
type PostCollaboratorHandler struct {
	repo            *repository.PostCollaboratorRepository
	postRepo        *repository.PostRepository
	userRepo        *repository.UserRepository
}

// NewPostCollaboratorHandler 创建文章协作者处理器实例
func NewPostCollaboratorHandler() *PostCollaboratorHandler {
	return &PostCollaboratorHandler{
		repo:     repository.NewPostCollaboratorRepository(),
		postRepo: repository.NewPostRepository(),
		userRepo: repository.NewUserRepository(),
	}
}

// GetCollaborators 获取文章协作者列表
func (h *PostCollaboratorHandler) GetCollaborators(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的文章ID")
		return
	}

	// 检查文章是否存在
	_, err = h.postRepo.GetByID(uint(postID))
	if err != nil {
		util.Error(c, 404, "文章不存在")
		return
	}

	users, err := h.repo.GetByPostID(uint(postID))
	if err != nil {
		util.ServerError(c, "获取协作者列表失败")
		return
	}

	util.Success(c, users)
}

// AddCollaborator 添加文章协作者
func (h *PostCollaboratorHandler) AddCollaborator(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的文章ID")
		return
	}

	// 获取当前用户
	userID, exists := c.Get("user_id")
	if !exists {
		util.Unauthorized(c, "未登录")
		return
	}
	role, _ := c.Get("role")

	// 检查文章是否存在
	post, err := h.postRepo.GetByID(uint(postID))
	if err != nil {
		util.Error(c, 404, "文章不存在")
		return
	}

	// 检查权限：只有文章作者和管理员可以添加协作者
	if post.UserID != userID.(uint) && !util.IsAdminRole(role.(string)) {
		util.Forbidden(c, "无权操作此文章")
		return
	}

	// 检查协作者数量限制（最多1个）
	count, err := h.repo.CountByPostID(uint(postID))
	if err != nil {
		util.ServerError(c, "获取协作者数量失败")
		return
	}
	if count >= 1 {
		util.Error(c, 400, "每篇文章最多添加1个协作者")
		return
	}

	// 解析请求体
	var req struct {
		UserID    uint `json:"user_id" binding:"required"`
		SortOrder int  `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	// 检查用户是否存在
	_, err = h.userRepo.GetByID(req.UserID)
	if err != nil {
		util.Error(c, 404, "用户不存在")
		return
	}

	// 检查是否已经是协作者
	exists, err = h.repo.Exists(uint(postID), req.UserID)
	if err != nil {
		util.ServerError(c, "检查协作者状态失败")
		return
	}
	if exists {
		util.Error(c, 400, "该用户已经是文章协作者")
		return
	}

	// 检查不能添加文章作者为协作者
	if post.UserID == req.UserID {
		util.Error(c, 400, "不能将文章作者添加为协作者")
		return
	}

	// 设置排序权重
	sortOrder := req.SortOrder
	if sortOrder == 0 {
		sortOrder = int(count) + 1
	}

	// 添加协作者
	err = h.repo.Add(uint(postID), req.UserID, sortOrder)
	if err != nil {
		util.ServerError(c, "添加协作者失败")
		return
	}

	util.SuccessWithMessage(c, "添加协作者成功", nil)
}

// RemoveCollaborator 移除文章协作者
func (h *PostCollaboratorHandler) RemoveCollaborator(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的文章ID")
		return
	}

	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的用户ID")
		return
	}

	// 获取当前用户
	currentUserID, exists := c.Get("user_id")
	if !exists {
		util.Unauthorized(c, "未登录")
		return
	}
	role, _ := c.Get("role")

	// 检查文章是否存在
	post, err := h.postRepo.GetByID(uint(postID))
	if err != nil {
		util.Error(c, 404, "文章不存在")
		return
	}

	// 检查权限：只有文章作者和超级管理员可以移除协作者
	if post.UserID != currentUserID.(uint) && role != constant.RoleSuperAdmin {
		util.Forbidden(c, "无权操作此文章")
		return
	}

	// 移除协作者
	err = h.repo.Remove(uint(postID), uint(userID))
	if err != nil {
		util.ServerError(c, "移除协作者失败")
		return
	}

	util.SuccessWithMessage(c, "移除协作者成功", nil)
}

// SearchUsers 搜索用户（用于添加协作者时选择用户）
func (h *PostCollaboratorHandler) SearchUsers(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		util.BadRequest(c, "请输入搜索关键词")
		return
	}

	// 限制搜索结果数量
	limit := 10

	users, err := h.userRepo.Search(keyword, limit)
	if err != nil {
		util.ServerError(c, "搜索用户失败")
		return
	}

	// 过敏敏感信息
	type UserResponse struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}

	var result []UserResponse
	for _, user := range users {
		result = append(result, UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
		})
	}

	util.Success(c, result)
}
