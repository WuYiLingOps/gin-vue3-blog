/*
 * 项目名称：blog-backend
 * 文件名称：ip_whitelist.go
 * 创建时间：2026-01-31 16:05:15
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：IP白名单管理处理器，提供IP白名单添加、查询、过期清理等功能，支持CIDR格式
 */
package handler

import (
	"blog-backend/db"
	"blog-backend/model"
	"blog-backend/util"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// IPWhitelistHandler IP白名单处理器结构体
type IPWhitelistHandler struct{}

// NewIPWhitelistHandler 创建IP白名单处理器实例
func NewIPWhitelistHandler() *IPWhitelistHandler {
	return &IPWhitelistHandler{}
}

// List 获取IP白名单列表
func (h *IPWhitelistHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var whitelist []model.IPWhitelist
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	db.DB.Model(&model.IPWhitelist{}).Count(&total)

	// 分页查询
	err := db.DB.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&whitelist).Error

	if err != nil {
		util.Error(c, 500, err.Error())
		return
	}

	util.Success(c, gin.H{
		"list":      whitelist,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Add 添加IP到白名单
func (h *IPWhitelistHandler) Add(c *gin.Context) {
	var req struct {
		IP       string `json:"ip" binding:"required"`
		Reason   string `json:"reason"`
		Duration int    `json:"duration"` // 有效期（小时），0表示永久
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "参数错误")
		return
	}

	// 检查IP格式（支持 CIDR）
	if !isValidIPOrCIDR(req.IP) {
		util.BadRequest(c, "IP地址格式不正确，支持 IPv4、IPv6 和 CIDR 格式")
		return
	}

	var expireAt *time.Time
	if req.Duration > 0 {
		expire := time.Now().Add(time.Duration(req.Duration) * time.Hour)
		expireAt = &expire
	}

	whitelist := model.IPWhitelist{
		IP:       req.IP,
		Reason:   req.Reason,
		ExpireAt: expireAt,
	}

	// 检查是否已存在
	var existing model.IPWhitelist
	if err := db.DB.Where("ip = ?", req.IP).First(&existing).Error; err == nil {
		util.Error(c, 400, "该IP已在白名单中")
		return
	}

	if err := db.DB.Create(&whitelist).Error; err != nil {
		util.Error(c, 500, err.Error())
		return
	}

	util.SuccessWithMessage(c, "添加成功", whitelist)
}

// Delete 从白名单中删除IP
func (h *IPWhitelistHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.BadRequest(c, "无效的ID")
		return
	}

	if err := db.DB.Delete(&model.IPWhitelist{}, id).Error; err != nil {
		util.Error(c, 500, err.Error())
		return
	}

	util.SuccessWithMessage(c, "删除成功", nil)
}

// Check 检查IP是否在白名单中
func (h *IPWhitelistHandler) Check(c *gin.Context) {
	ip := c.Query("ip")
	if ip == "" {
		util.BadRequest(c, "请提供IP地址")
		return
	}

	// 检查数据库中的白名单
	var whitelist []model.IPWhitelist
	db.DB.Find(&whitelist)

	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		util.Success(c, gin.H{
			"whitelisted": false,
		})
		return
	}

	// 检查是否匹配
	for _, wl := range whitelist {
		// 检查是否过期
		if wl.ExpireAt != nil && wl.ExpireAt.Before(time.Now()) {
			continue
		}

		// 支持 CIDR 格式
		if strings.Contains(wl.IP, "/") {
			_, ipNet, err := net.ParseCIDR(wl.IP)
			if err == nil && ipNet.Contains(clientIP) {
				util.Success(c, gin.H{
					"whitelisted": true,
					"info":        wl,
				})
				return
			}
		} else {
			// 精确匹配
			wlIP := net.ParseIP(wl.IP)
			if wlIP != nil && wlIP.Equal(clientIP) {
				util.Success(c, gin.H{
					"whitelisted": true,
					"info":        wl,
				})
				return
			}
		}
	}

	util.Success(c, gin.H{
		"whitelisted": false,
	})
}

// CleanExpired 清理过期的白名单记录
func (h *IPWhitelistHandler) CleanExpired(c *gin.Context) {
	result := db.DB.Where("expire_at IS NOT NULL AND expire_at < ?", time.Now()).
		Delete(&model.IPWhitelist{})

	if result.Error != nil {
		util.Error(c, 500, result.Error.Error())
		return
	}

	util.SuccessWithMessage(c, "清理成功", gin.H{
		"deleted_count": result.RowsAffected,
	})
}

// isValidIPOrCIDR 验证IP或CIDR格式
func isValidIPOrCIDR(ip string) bool {
	// 尝试解析为 CIDR
	if strings.Contains(ip, "/") {
		_, _, err := net.ParseCIDR(ip)
		return err == nil
	}

	// 尝试解析为普通IP
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}


