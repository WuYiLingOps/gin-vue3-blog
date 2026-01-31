/*
 * 项目名称：blog-backend
 * 文件名称：oss.go
 * 创建时间：2026-01-31 16:41:24
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：对象存储工具函数，提供本地存储、阿里云OSS和腾讯云COS的统一文件上传接口
 */
package util

import (
	"blog-backend/config"
	"blog-backend/repository"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
)

// StorageType 存储类型枚举
type StorageType string

const (
	StorageLocal StorageType = "local" // 本地存储
	StorageOSS   StorageType = "oss"   // 阿里云OSS对象存储
	StorageCOS   StorageType = "cos"   // 腾讯云COS对象存储
)

// getUploadSettings 从数据库获取上传配置（仅存储类型）
func getUploadSettings() map[string]string {
	settingRepo := repository.NewSettingRepository()
	settings, err := settingRepo.GetByGroup("upload")
	if err != nil {
		return map[string]string{"storage_type": "local"}
	}

	result := make(map[string]string)
	for _, setting := range settings {
		result[setting.Key] = setting.Value
	}

	if result["storage_type"] == "" {
		result["storage_type"] = "local"
	}

	return result
}

// GetStorageType 从数据库获取当前存储类型
func GetStorageType() StorageType {
	settings := getUploadSettings()
	switch settings["storage_type"] {
	case string(StorageOSS):
		return StorageOSS
	case string(StorageCOS):
		return StorageCOS
	}
	return StorageLocal
}

// UploadFile 上传文件（根据配置自动选择存储方式）
func UploadFile(file *multipart.FileHeader, dir string) (string, error) {
	storageType := GetStorageType()

	switch storageType {
	case StorageOSS:
		return UploadToOSS(file, dir)
	case StorageCOS:
		return UploadToCOS(file, dir)
	default:
		// 本地存储
		filePath, err := SaveUploadedFile(file, dir)
		if err != nil {
			return "", err
		}
		// 将本地路径转换为 URL
		return GetFileURL(filePath), nil
	}
}

// UploadToOSS 上传文件到阿里云 OSS
func UploadToOSS(file *multipart.FileHeader, dir string) (string, error) {
	// 验证文件大小
	if file.Size > MaxFileSize {
		return "", errors.New("文件大小超过限制（最大 5MB）")
	}

	// 验证文件类型
	contentType := file.Header.Get("Content-Type")
	if !allowedImageTypes[contentType] {
		return "", errors.New("不支持的文件类型（仅支持 jpg, png, gif, webp）")
	}

	// 从配置文件获取 OSS 配置
	endpoint := config.Cfg.OSS.Endpoint
	accessKeyID := config.Cfg.OSS.AccessKeyID
	accessKeySecret := config.Cfg.OSS.AccessKeySecret
	bucketName := config.Cfg.OSS.BucketName
	domain := config.Cfg.OSS.Domain

	// 检查 OSS 配置
	if endpoint == "" || accessKeyID == "" || accessKeySecret == "" || bucketName == "" {
		return "", errors.New("OSS 配置不完整，请先在配置文件中设置 OSS 参数")
	}

	// 创建 OSS 客户端
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return "", fmt.Errorf("创建 OSS 客户端失败: %w", err)
	}

	// 获取 Bucket
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return "", fmt.Errorf("获取 OSS Bucket 失败: %w", err)
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		return "", errors.New("无法打开文件")
	}
	defer src.Close()

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%s%s", time.Now().Format("20060102150405"), uuid.New().String()[:8], ext)

	// 构建对象路径（去掉本地路径前缀）
	objectKey := strings.TrimPrefix(dir, "uploads/")
	if objectKey != "" {
		objectKey = objectKey + "/" + filename
	} else {
		objectKey = filename
	}

	// 上传文件到 OSS
	err = bucket.PutObject(objectKey, src, oss.ContentType(contentType))
	if err != nil {
		return "", fmt.Errorf("上传到 OSS 失败: %w", err)
	}

	// 返回访问 URL
	fileURL := getOSSFileURL(objectKey, endpoint, bucketName, domain)
	return fileURL, nil
}

// UploadToCOS 上传文件到腾讯云 COS
func UploadToCOS(file *multipart.FileHeader, dir string) (string, error) {
	// 验证文件大小
	if file.Size > MaxFileSize {
		return "", errors.New("文件大小超过限制（最大 5MB）")
	}

	// 验证文件类型
	contentType := file.Header.Get("Content-Type")
	if !allowedImageTypes[contentType] {
		return "", errors.New("不支持的文件类型（仅支持 jpg, png, gif, webp）")
	}

	// 从配置文件获取 COS 配置
	bucketURL := config.Cfg.COS.BucketURL
	secretID := config.Cfg.COS.SecretID
	secretKey := config.Cfg.COS.SecretKey
	domain := config.Cfg.COS.Domain

	// 检查 COS 配置
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

	// 打开文件
	src, err := file.Open()
	if err != nil {
		return "", errors.New("无法打开文件")
	}
	defer src.Close()

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%s%s", time.Now().Format("20060102150405"), uuid.New().String()[:8], ext)

	// 构建对象路径（去掉本地路径前缀）
	objectKey := strings.TrimPrefix(dir, "uploads/")
	if objectKey != "" {
		objectKey = objectKey + "/" + filename
	} else {
		objectKey = filename
	}

	// 上传文件到 COS
	_, err = client.Object.Put(context.Background(), objectKey, src, &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: contentType,
		},
	})
	if err != nil {
		return "", fmt.Errorf("上传到 COS 失败: %w", err)
	}

	// 返回访问 URL
	fileURL := getCOSFileURL(objectKey, bucketURL, domain)
	return fileURL, nil
}

// getOSSFileURL 获取 OSS 文件访问 URL
func getOSSFileURL(objectKey, endpoint, bucketName, domain string) string {
	// 如果配置了自定义域名，使用自定义域名
	if domain != "" {
		return strings.TrimSuffix(domain, "/") + "/" + objectKey
	}

	// 否则使用默认的 OSS 域名
	bucketDomain := fmt.Sprintf("https://%s.%s", bucketName, endpoint)
	return bucketDomain + "/" + objectKey
}

// GetOSSFileURL 获取 OSS 文件访问 URL（向后兼容）
func GetOSSFileURL(objectKey string) string {
	return getOSSFileURL(objectKey, config.Cfg.OSS.Endpoint, config.Cfg.OSS.BucketName, config.Cfg.OSS.Domain)
}

// getCOSFileURL 获取 COS 文件访问 URL
func getCOSFileURL(objectKey, bucketURL, domain string) string {
	// 如果配置了自定义域名，使用自定义域名
	if domain != "" {
		return strings.TrimSuffix(domain, "/") + "/" + objectKey
	}

	// 否则使用 BucketURL（例如：https://<bucket>.cos.<region>.myqcloud.com）
	return strings.TrimSuffix(bucketURL, "/") + "/" + objectKey
}

// DeleteOSSFile 删除 OSS 文件
func DeleteOSSFile(fileURL string) error {
	// 从配置文件获取 OSS 配置
	endpoint := config.Cfg.OSS.Endpoint
	accessKeyID := config.Cfg.OSS.AccessKeyID
	accessKeySecret := config.Cfg.OSS.AccessKeySecret
	bucketName := config.Cfg.OSS.BucketName
	domain := config.Cfg.OSS.Domain

	// 检查 OSS 配置
	if endpoint == "" || accessKeyID == "" || accessKeySecret == "" || bucketName == "" {
		return errors.New("OSS 配置不完整")
	}

	// 从 URL 提取对象键
	objectKey := extractObjectKeyFromURL(fileURL, endpoint, bucketName, domain)
	if objectKey == "" {
		return errors.New("无效的文件 URL")
	}

	// 创建 OSS 客户端
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return err
	}

	// 获取 Bucket
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	// 删除对象
	err = bucket.DeleteObject(objectKey)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCOSFile 删除 COS 文件
func DeleteCOSFile(fileURL string) error {
	bucketURL := config.Cfg.COS.BucketURL
	secretID := config.Cfg.COS.SecretID
	secretKey := config.Cfg.COS.SecretKey
	domain := config.Cfg.COS.Domain

	if bucketURL == "" || secretID == "" || secretKey == "" {
		return errors.New("COS 配置不完整")
	}

	objectKey := extractCOSObjectKeyFromURL(fileURL, bucketURL, domain)
	if objectKey == "" {
		return errors.New("无效的文件 URL")
	}

	u, err := url.Parse(bucketURL)
	if err != nil {
		return fmt.Errorf("COS BucketURL 无效: %w", err)
	}

	baseURL := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(baseURL, &http.Client{
		Timeout: 60 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})

	_, err = client.Object.Delete(context.Background(), objectKey)
	if err != nil {
		return fmt.Errorf("删除 COS 对象失败: %w", err)
	}

	return nil
}

// extractObjectKeyFromURL 从 URL 中提取对象键
func extractObjectKeyFromURL(fileURL, endpoint, bucketName, domain string) string {
	// 支持自定义域名
	if domain != "" {
		domainClean := strings.TrimSuffix(domain, "/")
		return strings.TrimPrefix(fileURL, domainClean+"/")
	}

	// 支持默认 OSS 域名
	bucketDomain := fmt.Sprintf("https://%s.%s/", bucketName, endpoint)
	return strings.TrimPrefix(fileURL, bucketDomain)
}

// DeleteFileByURL 根据 URL 删除文件（自动判断存储类型）
func DeleteFileByURL(fileURL string) error {
	if fileURL == "" {
		return nil
	}

	// 判断是否为 OSS 文件
	if isOSSFile(fileURL) {
		return DeleteOSSFile(fileURL)
	}

	// 判断是否为 COS 文件
	if isCOSFile(fileURL) {
		return DeleteCOSFile(fileURL)
	}

	// 本地文件
	// 将 URL 转换为本地路径
	filePath := strings.TrimPrefix(fileURL, "/")
	return DeleteFile(filePath)
}

// isOSSFile 判断是否为 OSS 文件
func isOSSFile(fileURL string) bool {
	// 检查是否包含 OSS 自定义域名
	if config.Cfg.OSS.Domain != "" && strings.HasPrefix(fileURL, config.Cfg.OSS.Domain) {
		return true
	}

	// 检查是否包含 OSS 默认域名
	if config.Cfg.OSS.BucketName != "" && config.Cfg.OSS.Endpoint != "" {
		bucketDomain := fmt.Sprintf("https://%s.%s", config.Cfg.OSS.BucketName, config.Cfg.OSS.Endpoint)
		if strings.HasPrefix(fileURL, bucketDomain) {
			return true
		}
	}

	return false
}

// extractCOSObjectKeyFromURL 从 COS URL 中提取对象键
func extractCOSObjectKeyFromURL(fileURL, bucketURL, domain string) string {
	// 支持自定义域名
	if domain != "" {
		domainClean := strings.TrimSuffix(domain, "/")
		return strings.TrimPrefix(fileURL, domainClean+"/")
	}

	// 使用 BucketURL 前缀
	bucketPrefix := strings.TrimSuffix(bucketURL, "/") + "/"
	return strings.TrimPrefix(fileURL, bucketPrefix)
}

// isCOSFile 判断是否为 COS 文件
func isCOSFile(fileURL string) bool {
	// 检查是否包含 COS 自定义域名
	if config.Cfg.COS.Domain != "" && strings.HasPrefix(fileURL, config.Cfg.COS.Domain) {
		return true
	}

	// 检查是否以 BucketURL 开头
	if config.Cfg.COS.BucketURL != "" && strings.HasPrefix(fileURL, strings.TrimSuffix(config.Cfg.COS.BucketURL, "/")) {
		return true
	}

	return false
}

// ValidateOSSConfig 验证 OSS 配置是否有效
func ValidateOSSConfig(endpoint, accessKeyID, accessKeySecret, bucketName string) error {
	if endpoint == "" || accessKeyID == "" || accessKeySecret == "" || bucketName == "" {
		return errors.New("OSS 配置不完整")
	}

	// 尝试创建客户端并访问 Bucket
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return fmt.Errorf("OSS 配置无效: %w", err)
	}

	_, err = client.Bucket(bucketName)
	if err != nil {
		return fmt.Errorf("无法访问 OSS Bucket: %w", err)
	}

	return nil
}

// ValidateCOSConfig 验证 COS 配置是否有效
func ValidateCOSConfig(bucketURL, secretID, secretKey string) error {
	if bucketURL == "" || secretID == "" || secretKey == "" {
		return errors.New("COS 配置不完整")
	}

	u, err := url.Parse(bucketURL)
	if err != nil {
		return fmt.Errorf("COS BucketURL 无效: %w", err)
	}

	baseURL := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(baseURL, &http.Client{
		Timeout: 30 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})

	// 简单调用 Get Bucket 检查配置是否可用
	_, _, err = client.Bucket.Get(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("无法访问 COS Bucket: %w", err)
	}

	return nil
}

// CopyFileReader 复制文件读取器（用于重复读取）
func CopyFileReader(src io.Reader) (io.Reader, io.Reader, error) {
	var buf []byte
	var err error

	buf, err = io.ReadAll(src)
	if err != nil {
		return nil, nil, err
	}

	return io.NopCloser(strings.NewReader(string(buf))),
		io.NopCloser(strings.NewReader(string(buf))),
		nil
}
