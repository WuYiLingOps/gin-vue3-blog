/*
 * 项目名称：blog-backend
 * 文件名称：operation_log.go
 * 创建时间：2026-02-06 22:00:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：操作日志处理器，提供操作日志的查询接口
 */
package handler

import (
	"strconv"

	"blog-backend/service"
	"blog-backend/util"

	"github.com/gin-gonic/gin"
)

// OperationLogHandler 操作日志处理器结构体
type OperationLogHandler struct {
	service *service.OperationLogService
}

// NewOperationLogHandler 创建操作日志处理器实例
func NewOperationLogHandler() *OperationLogHandler {
	return &OperationLogHandler{
		service: service.NewOperationLogService(),
	}
}

// List 获取操作日志列表
func (h *OperationLogHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	module := c.Query("module")
	action := c.Query("action")
	username := c.Query("username")

	logs, total, err := h.service.List(page, pageSize, module, action, username)
	if err != nil {
		util.ServerError(c, "获取操作日志列表失败")
		return
	}

	util.PageSuccess(c, logs, total, page, pageSize)
}

// GetByID 获取操作日志详情
func (h *OperationLogHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的日志ID")
		return
	}

	log, err := h.service.GetByID(uint(id))
	if err != nil {
		util.Error(c, 404, "操作日志不存在")
		return
	}

	util.Success(c, log)
}

// Delete 删除单个操作日志
func (h *OperationLogHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的日志ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		util.Error(c, 400, "删除操作日志失败")
		return
	}

	util.SuccessWithMessage(c, "删除成功", nil)
}

// DeleteBatch 批量删除操作日志
func (h *OperationLogHandler) DeleteBatch(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	if len(req.IDs) == 0 {
		util.BadRequest(c, "请选择要删除的日志")
		return
	}

	if err := h.service.DeleteBatch(req.IDs); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "批量删除成功", nil)
}
