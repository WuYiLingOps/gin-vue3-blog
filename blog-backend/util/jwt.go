/*
 * 项目名称：blog-backend
 * 文件名称：jwt.go
 * 创建时间：2026-01-31 16:41:24
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：JWT工具函数，提供Token生成、解析和刷新功能
 */
package util

import (
	"errors"
	"time"

	"blog-backend/config"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT声明结构体
// 包含用户ID、用户名、角色等自定义声明，以及JWT标准声明（过期时间、签发时间等）
type Claims struct {
	UserID               uint   `json:"user_id"`  // 用户ID
	Username             string `json:"username"` // 用户名
	Role                 string `json:"role"`     // 用户角色：admin或user
	jwt.RegisteredClaims        // JWT标准声明（过期时间、签发时间、签发者等）
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID uint, username, role string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(config.Cfg.JWT.ExpireHours) * time.Hour)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "blog-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Cfg.JWT.Secret))
}

// ParseToken 解析 JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken 刷新 Token
func RefreshToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	// 如果 Token 还有超过 1 小时的有效期，不需要刷新
	if time.Until(claims.ExpiresAt.Time) > time.Hour {
		return tokenString, nil
	}

	// 生成新的 Token
	return GenerateToken(claims.UserID, claims.Username, claims.Role)
}

// GetTimeAfterMinutes 获取指定分钟后的时间
func GetTimeAfterMinutes(minutes int) time.Time {
	return time.Now().Add(time.Duration(minutes) * time.Minute)
}
