/*
 * 项目名称：blog-backend
 * 文件名称：upload.go
 * 创建时间：2026-01-31 16:41:24
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：文件上传工具函数，提供本地文件上传、保存和删除功能
 */
package util

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	// 上传目录
	UploadDir = "uploads"
	// 头像目录
	AvatarDir = "uploads/avatars"
	// 最大文件大小 (5MB)
	MaxFileSize = 5 << 20
)

// 允许的图片格式
var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// InitUploadDirs 初始化上传目录
func InitUploadDirs() error {
	dirs := []string{UploadDir, AvatarDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

// SaveUploadedFile 保存上传的文件
func SaveUploadedFile(file *multipart.FileHeader, dir string) (string, error) {
	// 验证文件大小
	if file.Size > MaxFileSize {
		return "", errors.New("文件大小超过限制（最大 5MB）")
	}

	// 验证文件类型
	src, err := file.Open()
	if err != nil {
		return "", errors.New("无法打开文件")
	}
	defer src.Close()

	// 读取文件头判断类型
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return "", errors.New("无法读取文件")
	}

	// 重置读取位置
	src.Seek(0, io.SeekStart)

	contentType := file.Header.Get("Content-Type")
	if !allowedImageTypes[contentType] {
		return "", errors.New("不支持的文件类型（仅支持 jpg, png, gif, webp）")
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%s%s", time.Now().Format("20060102150405"), uuid.New().String()[:8], ext)
	filepath := filepath.Join(dir, filename)

	// 创建目标文件
	dst, err := os.Create(filepath)
	if err != nil {
		return "", errors.New("无法创建文件")
	}
	defer dst.Close()

	// 复制文件内容
	if _, err = io.Copy(dst, src); err != nil {
		return "", errors.New("文件保存失败")
	}

	return filepath, nil
}

// DeleteFile 删除文件
func DeleteFile(filePath string) error {
	// 确保文件在 uploads 目录下
	if !strings.HasPrefix(filePath, UploadDir) {
		return errors.New("invalid file path")
	}

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return nil // 文件不存在，认为删除成功
		}
		return err
	}

	return nil
}

// GetFileURL 获取文件访问 URL
func GetFileURL(filePath string) string {
	// 将本地路径转换为 URL 路径
	return "/" + strings.ReplaceAll(filePath, "\\", "/")
}
