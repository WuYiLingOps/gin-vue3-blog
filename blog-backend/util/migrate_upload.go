/*
 * 项目名称：blog-backend
 * 文件名称：migrate_upload.go
 * 创建时间：2026-11-20 10:00:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：图片迁移工具函数，提供字节流上传到指定存储、图片来源识别等功能
 */
package util

import (
	"blog-backend/config"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
)

// 腾讯云 COS / 阿里云 OSS 的标准访问域名特征（用于不依赖运行时配置的结构化识别）
var (
	cosURLPattern = regexp.MustCompile(`^https?://[^/]+\.cos[^.]*\.[^/]+\.myqcloud\.com(/|$)`)
	ossURLPattern = regexp.MustCompile(`^https?://[^/]+\.oss[a-z0-9-]*\.aliyuncs\.com(/|$)`)
)

// MigrationMaxFileSize 图片迁移允许的最大文件大小（50MB，迁移场景放宽普通上传的 5MB 限制）
const MigrationMaxFileSize = 50 << 20

// StorageExternal 外链图片（不属于本站配置的任何存储）
const StorageExternal StorageType = "external"

// ClassifyImageURL 识别图片 URL 所属的存储类型
// 判断规则：本地 /uploads/ 前缀 → local；配置域名前缀精确匹配 → 对应存储；
// 标准厂商域名格式（*.cos.*.myqcloud.com / *.oss*.aliyuncs.com）结构化识别 → 对应存储；都不匹配 → external 外链
func ClassifyImageURL(fileURL string) StorageType {
	if fileURL == "" {
		return StorageExternal
	}

	// 本地存储：相对路径或站点 /uploads/ 路径
	if strings.HasPrefix(fileURL, "/uploads/") {
		return StorageLocal
	}

	// 当前配置的域名前缀精确匹配（覆盖自定义域名场景），config 未加载时跳过
	if config.Cfg != nil {
		if isOSSFile(fileURL) {
			return StorageOSS
		}

		if isCOSFile(fileURL) {
			return StorageCOS
		}

		if isEasyImageFile(fileURL) {
			return StorageEasyImage
		}
	}

	// 标准厂商域名的结构化识别（历史 bucket 与当前配置不一致时仍能识别归属）
	if cosURLPattern.MatchString(fileURL) {
		return StorageCOS
	}
	if ossURLPattern.MatchString(fileURL) {
		return StorageOSS
	}

	return StorageExternal
}

// IsManagedStorageURL 判断 URL 是否属于当前运行配置可管理的存储（能通过 API 删除源文件）
// 与 ClassifyImageURL 的区别：结构化识别出的历史 bucket 不在当前配置内，无法 API 删除，返回 false
func IsManagedStorageURL(fileURL string) bool {
	if fileURL == "" {
		return false
	}

	if strings.HasPrefix(fileURL, "/uploads/") {
		// 本地文件可管理，删除时文件不存在会被视为成功
		return true
	}

	if config.Cfg != nil {
		if isOSSFile(fileURL) || isCOSFile(fileURL) || isEasyImageFile(fileURL) {
			return true
		}
	}

	return false
}

// UploadBytesToStorage 将字节流上传到指定的存储类型（图片迁移专用）
// 与 UploadFile 不同：不依赖全局存储配置，由 target 显式指定目标存储
func UploadBytesToStorage(data []byte, contentType, filename, dir string, target StorageType) (string, error) {
	if len(data) == 0 {
		return "", errors.New("文件内容为空")
	}

	if len(data) > MigrationMaxFileSize {
		return "", errors.New("文件大小超过迁移限制（最大 50MB）")
	}

	// 文件类型校验（与普通上传保持一致，仅支持图片格式）
	if !allowedImageTypes[contentType] {
		return "", errors.New("不支持的文件类型（仅支持 jpg, png, gif, webp）")
	}

	switch target {
	case StorageOSS:
		return uploadBytesToOSS(data, contentType, filename, dir)
	case StorageCOS:
		return uploadBytesToCOS(data, contentType, filename, dir)
	case StorageEasyImage:
		return uploadBytesToEasyImage(data, contentType, filename)
	default:
		// 本地存储
		filePath, err := saveMigratedFileLocal(data, dir, filepath.Ext(filename))
		if err != nil {
			return "", err
		}
		return GetFileURL(filePath), nil
	}
}

// uploadBytesToOSS 上传字节流到阿里云 OSS
func uploadBytesToOSS(data []byte, contentType, filename, dir string) (string, error) {
	endpoint := config.Cfg.OSS.Endpoint
	accessKeyID := config.Cfg.OSS.AccessKeyID
	accessKeySecret := config.Cfg.OSS.AccessKeySecret
	bucketName := config.Cfg.OSS.BucketName
	domain := config.Cfg.OSS.Domain

	if endpoint == "" || accessKeyID == "" || accessKeySecret == "" || bucketName == "" {
		return "", errors.New("OSS 配置不完整，请先在配置文件中设置 OSS 参数")
	}

	client, err := oss.New(endpoint, accessKeyID, accessKeySecret, oss.Timeout(10, 30))
	if err != nil {
		return "", fmt.Errorf("创建 OSS 客户端失败: %w", err)
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return "", fmt.Errorf("获取 OSS Bucket 失败: %w", err)
	}

	objectKey := buildMigratedObjectKey(dir, generateMigratedFilename(filename))

	if err = bucket.PutObject(objectKey, bytes.NewReader(data), oss.ContentType(contentType)); err != nil {
		return "", fmt.Errorf("上传到 OSS 失败: %w", err)
	}

	return getOSSFileURL(objectKey, endpoint, bucketName, domain), nil
}

// uploadBytesToCOS 上传字节流到腾讯云 COS
func uploadBytesToCOS(data []byte, contentType, filename, dir string) (string, error) {
	bucketURL := config.Cfg.COS.BucketURL
	secretID := config.Cfg.COS.SecretID
	secretKey := config.Cfg.COS.SecretKey
	domain := config.Cfg.COS.Domain

	if bucketURL == "" || secretID == "" || secretKey == "" {
		return "", errors.New("COS 配置不完整，请先在配置文件中设置 COS 参数")
	}

	u, err := url.Parse(bucketURL)
	if err != nil {
		return "", fmt.Errorf("COS BucketURL 无效: %w", err)
	}

	baseURL := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(baseURL, &http.Client{
		Timeout: 60 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})

	objectKey := buildMigratedObjectKey(dir, generateMigratedFilename(filename))

	_, err = client.Object.Put(context.Background(), objectKey, bytes.NewReader(data), &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: contentType,
		},
	})
	if err != nil {
		return "", fmt.Errorf("上传到 COS 失败: %w", err)
	}

	return getCOSFileURL(objectKey, bucketURL, domain), nil
}

// uploadBytesToEasyImage 上传字节流到 EasyImage 简单图床
func uploadBytesToEasyImage(data []byte, contentType, filename string) (string, error) {
	apiURL := config.Cfg.EasyImage.ApiURL
	token := config.Cfg.EasyImage.Token

	if apiURL == "" || token == "" {
		return "", errors.New("EasyImage 配置不完整，请先在配置文件中设置 easyimage 参数")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if err := writer.WriteField("token", token); err != nil {
		return "", fmt.Errorf("构建 EasyImage 表单失败: %w", err)
	}

	part, err := writer.CreateFormFile("image", filepath.Base(filename))
	if err != nil {
		return "", fmt.Errorf("构建 EasyImage 表单失败: %w", err)
	}
	if _, err := io.Copy(part, bytes.NewReader(data)); err != nil {
		return "", fmt.Errorf("读取文件内容失败: %w", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("构建 EasyImage 表单失败: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, body)
	if err != nil {
		return "", fmt.Errorf("创建 EasyImage 请求失败: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("上传到 EasyImage 失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取 EasyImage 响应失败: %w", err)
	}

	var result easyImageResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("解析 EasyImage 响应失败（HTTP %d）: %w", resp.StatusCode, err)
	}

	if result.Result != "success" || result.URL == "" {
		if result.Message != "" {
			return "", fmt.Errorf("上传到 EasyImage 失败: %s", result.Message)
		}
		return "", fmt.Errorf("上传到 EasyImage 失败（HTTP %d）", resp.StatusCode)
	}

	return normalizeEasyImageURL(result.URL), nil
}

// saveMigratedFileLocal 保存迁移文件到本地 uploads 目录
func saveMigratedFileLocal(data []byte, dir, ext string) (string, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	filename := generateMigratedFilename("img" + ext)
	filePath := filepath.Join(dir, filename)

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("保存文件失败: %w", err)
	}

	return filePath, nil
}

// generateMigratedFilename 生成迁移文件的唯一文件名（与普通上传的命名规则一致）
func generateMigratedFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	return fmt.Sprintf("%s_%s%s", time.Now().Format("20060102150405"), uuid.New().String()[:8], ext)
}

// buildMigratedObjectKey 构建迁移文件的对象键（与普通上传保持一致：去掉 uploads/ 前缀）
func buildMigratedObjectKey(dir, filename string) string {
	objectKey := strings.TrimPrefix(dir, "uploads/")
	if objectKey != "" {
		objectKey = objectKey + "/" + filename
	} else {
		objectKey = filename
	}
	return objectKey
}
