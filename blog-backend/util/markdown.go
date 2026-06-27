/*
 * 项目名称：blog-backend
 * 文件名称：markdown.go
 * 创建时间：2026-06-27 15:30:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：Markdown工具函数，提供图片URL提取、GitHub图片代理解析、文件名清理等功能
 */
package util

import (
	"regexp"
	"strings"
)

// ExtractImageURLs 从 Markdown 内容中提取所有图片 URL，自动去重
func ExtractImageURLs(content string) []string {
	re := regexp.MustCompile(`!\[.*?\]\((.*?)\)`)
	matches := re.FindAllStringSubmatch(content, -1)

	seen := make(map[string]bool)
	var urls []string
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		url := strings.TrimSpace(match[1])
		if url == "" || seen[url] {
			continue
		}
		seen[url] = true
		urls = append(urls, url)
	}
	return urls
}

// IsGitHubRawURL 判断是否为 GitHub Raw 图片 URL
func IsGitHubRawURL(imgURL string) bool {
	return strings.Contains(imgURL, "raw.githubusercontent.com") ||
		strings.Contains(imgURL, "github.com/") && strings.Contains(imgURL, "/raw/")
}

// ResolveExportImageURL 处理导出时的图片 URL，GitHub 图片走代理
func ResolveExportImageURL(imgURL string, proxy string) string {
	if proxy == "" {
		return imgURL
	}
	if IsGitHubRawURL(imgURL) {
		proxy = strings.TrimRight(proxy, "/")
		return proxy + "/" + imgURL
	}
	return imgURL
}

// SanitizeExportFilename 清理文件名中的非法字符
func SanitizeExportFilename(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	re := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)
	name = re.ReplaceAllString(name, "-")
	name = strings.Trim(name, ". ")
	if len(name) > 200 {
		name = name[:200]
	}
	return name
}
