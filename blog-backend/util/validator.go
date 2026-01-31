/*
 * 项目名称：blog-backend
 * 文件名称：validator.go
 * 创建时间：2026-01-31 16:41:24
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：数据验证工具函数，提供邮箱、用户名、密码等格式验证功能
 */
package util

import (
	"regexp"
	"unicode/utf8"
)

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}

// ValidateUsername 验证用户名格式（3-20个字符，只能包含字母、数字、下划线）
func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	pattern := `^[a-zA-Z0-9_]+$`
	match, _ := regexp.MatchString(pattern, username)
	return match
}

// ValidatePassword 验证密码强度（至少6个字符）
func ValidatePassword(password string) bool {
	return utf8.RuneCountInString(password) >= 6
}

// ValidateRequired 验证必填字段
func ValidateRequired(value string) bool {
	return len(value) > 0
}

// ValidateLength 验证字符串长度
func ValidateLength(value string, min, max int) bool {
	length := utf8.RuneCountInString(value)
	return length >= min && length <= max
}
