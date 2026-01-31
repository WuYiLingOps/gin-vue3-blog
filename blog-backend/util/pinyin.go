/*
 * 项目名称：blog-backend
 * 文件名称：pinyin.go
 * 创建时间：2026-01-31 16:41:24
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：拼音转换工具函数，提供中文转拼音生成URL友好的slug功能
 */
package util

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mozillazg/go-pinyin"
)

// GenerateSlug 根据标题生成URL友好的slug（拼音）
// 支持混合中英文：中文转换为拼音，英文和数字直接保留
func GenerateSlug(title string) string {
	if title == "" {
		return ""
	}

	// 拼音转换参数
	a := pinyin.NewArgs()
	a.Style = pinyin.Normal // 不带声调
	a.Separator = ""        // 不使用分隔符，我们手动处理

	var result strings.Builder
	runes := []rune(title)

	// 标记上一个有效字符的类型：0=无, 1=中文/拼音, 2=英文数字
	lastType := 0
	// 标记是否需要添加分隔符（遇到空格/标点后）
	needSeparator := false

	for i := 0; i < len(runes); i++ {
		char := runes[i]

		// 判断是否为中文字符
		if isChinese(char) {
			// 中文字符转换为拼音
			pinyinSlice := pinyin.Pinyin(string(char), a)
			if len(pinyinSlice) > 0 && len(pinyinSlice[0]) > 0 {
				pinyinStr := pinyinSlice[0][0]
				// 如果需要分隔符、之前是英文数字，或者之前也是中文（连续中文之间需要分隔），添加连字符
				if needSeparator || (result.Len() > 0 && (lastType == 2 || lastType == 1)) {
					result.WriteString("-")
					needSeparator = false
				}
				result.WriteString(pinyinStr)
				lastType = 1
			}
		} else if isAlphanumeric(char) {
			// 英文字母或数字，直接保留
			// 如果需要分隔符或之前是中文/拼音，添加连字符
			if needSeparator || (result.Len() > 0 && lastType == 1) {
				result.WriteString("-")
				needSeparator = false
			}
			result.WriteRune(char)
			lastType = 2
		} else {
			// 其他字符（空格、标点等），标记需要分隔符
			if result.Len() > 0 {
				needSeparator = true
			}
		}
	}

	slug := strings.ToLower(result.String())

	// 移除连续的连字符
	reg := regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")

	// 移除开头和结尾的连字符
	slug = strings.Trim(slug, "-")

	// 如果结果为空，返回默认值
	if slug == "" {
		slug = "post"
	}

	return slug
}

// isChinese 判断是否为中文字符
func isChinese(r rune) bool {
	return r >= 0x4e00 && r <= 0x9fff
}

// isAlphanumeric 判断是否为字母或数字
func isAlphanumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

// GenerateUniqueSlug 生成唯一的slug，如果已存在则添加数字后缀
func GenerateUniqueSlug(baseSlug string, checkExists func(string) bool) string {
	slug := baseSlug
	counter := 1

	for checkExists(slug) {
		slug = baseSlug + "-" + fmt.Sprintf("%d", counter)
		counter++
		if counter > 100 { // 防止无限循环
			break
		}
	}

	return slug
}
