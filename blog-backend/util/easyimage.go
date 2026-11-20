/*
 * 项目名称：blog-backend
 * 文件名称：easyimage.go
 * 创建时间：2026-11-20 10:00:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：EasyImage简单图床工具函数，提供图床上传和配置校验功能
 */
package util

import (
	"blog-backend/config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

// easyImageResponse EasyImage 上传接口响应结构
type easyImageResponse struct {
	Result  string `json:"result"`  // success 或 failed
	Code    int    `json:"code"`    // 状态码，成功为 200
	URL     string `json:"url"`     // 图片访问 URL
	Thumb   string `json:"thumb"`   // 缩略图 URL
	Del     string `json:"del"`     // 删除链接（含一次性 hash）
	Message string `json:"message"` // 失败时的错误信息
}

// UploadToEasyImage 上传文件到 EasyImage 简单图床
func UploadToEasyImage(file *multipart.FileHeader, dir string) (string, error) {
	// 验证文件大小
	if file.Size > MaxFileSize {
		return "", errors.New("文件大小超过限制（最大 5MB）")
	}

	// 验证文件类型
	contentType := file.Header.Get("Content-Type")
	if !allowedImageTypes[contentType] {
		return "", errors.New("不支持的文件类型（仅支持 jpg, png, gif, webp）")
	}

	// 从配置文件获取 EasyImage 配置
	apiURL := config.Cfg.EasyImage.ApiURL
	token := config.Cfg.EasyImage.Token

	// 检查 EasyImage 配置
	if apiURL == "" || token == "" {
		return "", errors.New("EasyImage 配置不完整，请先在配置文件中设置 easyimage 参数")
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		return "", errors.New("无法打开文件")
	}
	defer src.Close()

	// 构建上传表单：文件字段名 image，文本字段 token
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if err := writer.WriteField("token", token); err != nil {
		return "", fmt.Errorf("构建 EasyImage 表单失败: %w", err)
	}
	part, err := writer.CreateFormFile("image", filepath.Base(file.Filename))
	if err != nil {
		return "", fmt.Errorf("构建 EasyImage 表单失败: %w", err)
	}
	if _, err := io.Copy(part, src); err != nil {
		return "", fmt.Errorf("读取文件内容失败: %w", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("构建 EasyImage 表单失败: %w", err)
	}

	// 发起上传请求
	req, err := http.NewRequest(http.MethodPost, apiURL, body)
	if err != nil {
		return "", fmt.Errorf("创建 EasyImage 请求失败: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("上传到 EasyImage 失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取并解析响应
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

// normalizeEasyImageURL 按配置的访问域名改写 EasyImage 返回的图片 URL
// EasyImage 的特性：即使站点配置了 https，上传接口返回的图片链接仍可能是 http（图床程序按请求协议或自身配置生成），
// 会导致 https 页面出现 mixed content。配置 domain 后强制以它为协议+域名返回标准 URL；
// 未配置 domain 时仅做同域名的 http→https 纠正（api_url 为 https 且图片 URL host 与其一致时），
// host 不一致则原样返回（避免改写第三方 CDN 等场景）
func normalizeEasyImageURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {
		return rawURL
	}

	// 配置了自定义访问域名：强制替换协议+域名
	if domain := config.Cfg.EasyImage.Domain; domain != "" {
		d, err := url.Parse(domain)
		if err == nil && d.Scheme != "" && d.Host != "" {
			u.Scheme = d.Scheme
			u.Host = d.Host
			return u.String()
		}
		return rawURL
	}

	// 未配置 domain：api_url 为 https 且图片 URL 与其同 host 且为 http 时，升级为 https
	if apiURL := config.Cfg.EasyImage.ApiURL; apiURL != "" {
		a, err := url.Parse(apiURL)
		if err == nil && a.Scheme == "https" && a.Hostname() != "" &&
			strings.EqualFold(u.Hostname(), a.Hostname()) && u.Scheme == "http" {
			u.Scheme = "https"
			return u.String()
		}
	}

	return rawURL
}

// ValidateEasyImageConfig 验证 EasyImage 配置是否有效（仅本地校验，不做真实请求）
func ValidateEasyImageConfig(apiURL, token string) error {
	if apiURL == "" || token == "" {
		return errors.New("EasyImage 配置不完整")
	}

	u, err := url.Parse(apiURL)
	if err != nil {
		return fmt.Errorf("EasyImage API 地址无效: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("EasyImage API 地址必须以 http:// 或 https:// 开头")
	}
	if u.Host == "" {
		return errors.New("EasyImage API 地址缺少主机名")
	}

	return nil
}
