/*
 * 项目名称：blog-backend
 * 文件名称：tag.go
 * 创建时间：2026-01-31 16:05:15
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：标签管理处理器，提供文章标签的增删改查功能，支持按标签查询文章
 */
package handler

import (
	"strconv"

	"blog-backend/service"
	"blog-backend/util"
	"github.com/gin-gonic/gin"
)

// TagHandler 标签处理器结构体
type TagHandler struct {
	service *service.TagService
}

// NewTagHandler 创建标签处理器实例
func NewTagHandler() *TagHandler {
	return &TagHandler{
		service: service.NewTagService(),
	}
}

// Create 创建标签
func (h *TagHandler) Create(c *gin.Context) {
	var req service.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	tag, err := h.service.Create(&req)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	// 记录操作日志
	tagID := tag.ID
	util.LogOperation(c, "create", "tag", &tagID, tag.Name, "创建标签："+tag.Name)

	util.SuccessWithMessage(c, "标签创建成功", tag)
}

// GetByID 获取标签详情
func (h *TagHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的标签ID")
		return
	}

	tag, err := h.service.GetByID(uint(id))
	if err != nil {
		util.Error(c, 404, err.Error())
		return
	}

	util.Success(c, tag)
}

// Update 更新标签
func (h *TagHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的标签ID")
		return
	}

	var req service.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	tag, err := h.service.Update(uint(id), &req)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	// 记录操作日志
	tagID := tag.ID
	util.LogOperation(c, "update", "tag", &tagID, tag.Name, "更新标签："+tag.Name)

	util.SuccessWithMessage(c, "标签更新成功", tag)
}

// Delete 删除标签
func (h *TagHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的标签ID")
		return
	}

	// 先获取标签信息用于日志记录
	tag, _ := h.service.GetByID(uint(id))
	var tagName string
	if tag != nil {
		tagName = tag.Name
	}

	if err := h.service.Delete(uint(id)); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	// 记录操作日志
	tagID := uint(id)
	util.LogOperation(c, "delete", "tag", &tagID, tagName, "删除标签："+tagName)

	util.SuccessWithMessage(c, "标签删除成功", nil)
}

// List 获取标签列表
func (h *TagHandler) List(c *gin.Context) {
	tags, err := h.service.List()
	if err != nil {
		util.ServerError(c, "获取标签列表失败")
		return
	}

	util.Success(c, tags)
}

// GetPostsByTag 获取标签下的文章列表
func (h *TagHandler) GetPostsByTag(c *gin.Context) {
	tagID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的标签ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 使用 PostService 获取文章列表
	postService := service.NewPostService()
	posts, total, err := postService.GetByTag(uint(tagID), page, pageSize)
	if err != nil {
		util.ServerError(c, "获取文章列表失败")
		return
	}

	util.PageSuccess(c, posts, total, page, pageSize)
}
