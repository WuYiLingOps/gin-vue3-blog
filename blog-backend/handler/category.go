package handler

import (
	"strconv"

	"blog-backend/service"
	"blog-backend/util"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		service: service.NewCategoryService(),
	}
}

// Create 创建分类
func (h *CategoryHandler) Create(c *gin.Context) {
	var req service.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	category, err := h.service.Create(&req)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "分类创建成功", category)
}

// GetByID 获取分类详情
func (h *CategoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的分类ID")
		return
	}

	category, err := h.service.GetByID(uint(id))
	if err != nil {
		util.Error(c, 404, err.Error())
		return
	}

	util.Success(c, category)
}

// Update 更新分类
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的分类ID")
		return
	}

	var req service.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	category, err := h.service.Update(uint(id), &req)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "分类更新成功", category)
}

// Delete 删除分类
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的分类ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "分类删除成功", nil)
}

// List 获取分类列表
func (h *CategoryHandler) List(c *gin.Context) {
	categories, err := h.service.List()
	if err != nil {
		util.ServerError(c, "获取分类列表失败")
		return
	}

	util.Success(c, categories)
}

