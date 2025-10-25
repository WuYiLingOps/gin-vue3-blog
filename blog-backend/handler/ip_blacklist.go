package handler

import (
	"blog-backend/db"
	"blog-backend/model"
	"blog-backend/util"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type IPBlacklistHandler struct{}

func NewIPBlacklistHandler() *IPBlacklistHandler {
	return &IPBlacklistHandler{}
}

// List 获取IP黑名单列表
func (h *IPBlacklistHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var blacklist []model.IPBlacklist
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	db.DB.Model(&model.IPBlacklist{}).Count(&total)

	// 分页查询
	err := db.DB.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&blacklist).Error

	if err != nil {
		util.Error(c, 500, err.Error())
		return
	}

	util.Success(c, gin.H{
		"list":      blacklist,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Add 添加IP到黑名单
func (h *IPBlacklistHandler) Add(c *gin.Context) {
	var req struct {
		IP       string `json:"ip" binding:"required"`
		Reason   string `json:"reason"`
		Duration int    `json:"duration"` // 封禁时长（小时），0表示永久
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "参数错误")
		return
	}

	// 检查IP格式
	if !util.IsValidIP(req.IP) {
		util.BadRequest(c, "IP地址格式不正确")
		return
	}

	var expireAt *time.Time
	if req.Duration > 0 {
		expire := time.Now().Add(time.Duration(req.Duration) * time.Hour)
		expireAt = &expire
	}

	blacklist := model.IPBlacklist{
		IP:       req.IP,
		Reason:   req.Reason,
		BanType:  2, // 手动封禁
		ExpireAt: expireAt,
	}

	// 检查是否已存在
	var existing model.IPBlacklist
	if err := db.DB.Where("ip = ?", req.IP).First(&existing).Error; err == nil {
		util.Error(c, 400, "该IP已在黑名单中")
		return
	}

	if err := db.DB.Create(&blacklist).Error; err != nil {
		util.Error(c, 500, err.Error())
		return
	}

	util.SuccessWithMessage(c, "添加成功", blacklist)
}

// Delete 从黑名单中删除IP
func (h *IPBlacklistHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.BadRequest(c, "无效的ID")
		return
	}

	if err := db.DB.Delete(&model.IPBlacklist{}, id).Error; err != nil {
		util.Error(c, 500, err.Error())
		return
	}

	util.SuccessWithMessage(c, "删除成功", nil)
}

// Check 检查IP是否在黑名单中
func (h *IPBlacklistHandler) Check(c *gin.Context) {
	ip := c.Query("ip")
	if ip == "" {
		util.BadRequest(c, "请提供IP地址")
		return
	}

	var blacklist model.IPBlacklist
	err := db.DB.Where("ip = ?", ip).First(&blacklist).Error

	if err != nil {
		util.Success(c, gin.H{
			"banned": false,
		})
		return
	}

	// 检查是否过期
	if blacklist.ExpireAt != nil && blacklist.ExpireAt.Before(time.Now()) {
		db.DB.Delete(&blacklist)
		util.Success(c, gin.H{
			"banned": false,
		})
		return
	}

	util.Success(c, gin.H{
		"banned": true,
		"info":   blacklist,
	})
}

// CleanExpired 清理过期的黑名单记录
func (h *IPBlacklistHandler) CleanExpired(c *gin.Context) {
	result := db.DB.Where("expire_at IS NOT NULL AND expire_at < ?", time.Now()).
		Delete(&model.IPBlacklist{})

	if result.Error != nil {
		util.Error(c, 500, result.Error.Error())
		return
	}

	util.SuccessWithMessage(c, "清理成功", gin.H{
		"deleted_count": result.RowsAffected,
	})
}

