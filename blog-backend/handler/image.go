/*
 * 项目名称：blog-backend
 * 文件名称：image.go
 * 创建时间：2026-11-20 10:00:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：图片管理处理器，提供全站图片盘点和图床迁移功能
 */
package handler

import (
	"fmt"
	"strconv"

	"blog-backend/service"
	"blog-backend/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ImageHandler 图片管理处理器结构体
type ImageHandler struct {
	service *service.ImageService
}

// NewImageHandler 创建图片管理处理器实例
func NewImageHandler() *ImageHandler {
	return &ImageHandler{
		service: service.NewImageService(),
	}
}

// ScanAssets 全站图片盘点
// GET /api/admin/images/assets?storage=xxx&keyword=xxx&refresh=true
// refresh=true 强制刷新 Redis 缓存并重新扫描数据库
func (h *ImageHandler) ScanAssets(c *gin.Context) {
	storage := c.Query("storage")
	keyword := c.Query("keyword")
	refresh := c.Query("refresh") == "true"

	assets, err := h.service.ScanImageAssets(storage, keyword, refresh)
	if err != nil {
		zap.L().Error("扫描图片资源失败", zap.Error(err))
		util.ServerError(c, "扫描图片资源失败")
		return
	}

	util.Success(c, assets)
}

// Migrate 发起图片迁移任务
// POST /api/admin/images/assets/migrate
func (h *ImageHandler) Migrate(c *gin.Context) {
	var req service.MigrateImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	// 注入操作者信息（用于迁移日志记录）
	if idVal, exists := c.Get("user_id"); exists {
		if id, ok := idVal.(uint); ok {
			req.OperatorID = &id
		}
	}
	if nameVal, exists := c.Get("username"); exists {
		if name, ok := nameVal.(string); ok {
			req.OperatorName = name
		}
	}

	if err := h.service.StartMigrate(&req); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	// 记录操作日志（迁移为异步任务，此处记录发起动作，结果在迁移日志模块查看）
	util.LogOperation(
		c, "update", "image", nil, targetStorageLabel(req.TargetStorag),
		fmt.Sprintf("发起图床迁移：%d 张图片迁移到 %s（删源：%s）",
			len(req.URLs), targetStorageLabel(req.TargetStorag), deleteSourceText(req.DeleteSource)),
	)

	util.SuccessWithMessage(c, "迁移任务已启动", nil)
}

// targetStorageLabel 目标存储类型的中文标签
func targetStorageLabel(storage string) string {
	switch storage {
	case "local":
		return "本地存储"
	case "oss":
		return "阿里云OSS"
	case "cos":
		return "腾讯云COS"
	case "easyimage":
		return "EasyImage 图床"
	default:
		return storage
	}
}

// deleteSourceText 是否删除源文件的文案
func deleteSourceText(deleteSource bool) string {
	if deleteSource {
		return "是"
	}
	return "否"
}

// MigrateStatus 获取迁移任务进度
// GET /api/admin/images/assets/migrate/status
func (h *ImageHandler) MigrateStatus(c *gin.Context) {
	util.Success(c, h.service.GetMigrateStatus())
}

// MigrateLogs 分页查询迁移日志
// GET /api/admin/images/migrate-logs?page=1&page_size=10
func (h *ImageHandler) MigrateLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	logs, total, err := h.service.ListMigrateLogs(page, pageSize)
	if err != nil {
		zap.L().Error("获取迁移日志失败", zap.Error(err))
		util.ServerError(c, "获取迁移日志失败")
		return
	}

	util.PageSuccess(c, logs, total, page, pageSize)
}

// DeleteMigrateLogs 删除迁移日志（body 为空或 ids 为空数组时清空全部）
// POST /api/admin/images/migrate-logs/delete
func (h *ImageHandler) DeleteMigrateLogs(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}
	// body 可为空（清空全部场景）
	_ = c.ShouldBindJSON(&req)

	if err := h.service.DeleteMigrateLogs(req.IDs); err != nil {
		util.Error(c, 400, "删除迁移日志失败")
		return
	}

	util.SuccessWithMessage(c, "迁移日志删除成功", nil)
}
