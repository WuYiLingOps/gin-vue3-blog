package util

import (
	"context"
	"errors"
	"fmt"
	"time"

	"blog-backend/db"

	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore

// CaptchaResponse 验证码响应
type CaptchaResponse struct {
	CaptchaID string `json:"captcha_id"`
	ImageData string `json:"image_data"`
}

const (
	// 验证码过期时间（2分钟）
	CaptchaExpireTime = 2 * time.Minute
	// IP限流：每个IP每分钟最多获取验证码次数
	CaptchaIPLimit = 10
	// IP限流时间窗口
	CaptchaIPLimitWindow = 1 * time.Minute
	// 验证码错误次数限制（防暴力破解）
	CaptchaMaxRetries = 5
	// 错误次数限制时间窗口
	CaptchaRetryWindow = 5 * time.Minute
)

// CheckCaptchaIPLimit 检查IP是否超过验证码获取限制
func CheckCaptchaIPLimit(ip string) error {
	ctx := context.Background()
	key := fmt.Sprintf("captcha:ip_limit:%s", ip)

	// 获取当前计数
	count, err := db.RDB.Get(ctx, key).Int()
	if err != nil && err.Error() != "redis: nil" {
		return err
	}

	// 检查是否超过限制
	if count >= CaptchaIPLimit {
		return errors.New("验证码获取过于频繁，请稍后再试")
	}

	// 增加计数
	pipe := db.RDB.Pipeline()
	pipe.Incr(ctx, key)
	// 如果是第一次，设置过期时间
	if count == 0 {
		pipe.Expire(ctx, key, CaptchaIPLimitWindow)
	}
	_, err = pipe.Exec(ctx)

	return err
}

// CheckCaptchaRetryLimit 检查验证码错误次数限制
func CheckCaptchaRetryLimit(ip string) error {
	ctx := context.Background()
	key := fmt.Sprintf("captcha:retry:%s", ip)

	count, err := db.RDB.Get(ctx, key).Int()
	if err != nil && err.Error() != "redis: nil" {
		return err
	}

	if count >= CaptchaMaxRetries {
		return errors.New("验证码错误次数过多，请稍后再试")
	}

	return nil
}

// IncrCaptchaRetryCount 增加验证码错误次数
func IncrCaptchaRetryCount(ip string) {
	ctx := context.Background()
	key := fmt.Sprintf("captcha:retry:%s", ip)

	count, _ := db.RDB.Get(ctx, key).Int()
	
	pipe := db.RDB.Pipeline()
	pipe.Incr(ctx, key)
	if count == 0 {
		pipe.Expire(ctx, key, CaptchaRetryWindow)
	}
	pipe.Exec(ctx)
}

// GenerateCaptcha 生成验证码（带IP限流）
func GenerateCaptcha(ip string) (*CaptchaResponse, error) {
	// 检查IP限流
	if err := CheckCaptchaIPLimit(ip); err != nil {
		return nil, err
	}

	// 配置验证码参数
	driver := base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(driver, store)

	// 生成验证码
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		return nil, err
	}

	// 获取验证码答案并存储到Redis
	answer := store.Get(id, false)
	if answer != "" {
		ctx := context.Background()
		// 存储到Redis，设置2分钟过期
		key := fmt.Sprintf("captcha:%s", id)
		err = db.RDB.Set(ctx, key, answer, CaptchaExpireTime).Err()
		if err != nil {
			return nil, err
		}
	}

	return &CaptchaResponse{
		CaptchaID: id,
		ImageData: b64s,
	}, nil
}

// VerifyCaptcha 验证验证码（带错误次数限制）
func VerifyCaptcha(captchaID, answer, ip string) error {
	if captchaID == "" || answer == "" {
		return errors.New("验证码ID和答案不能为空")
	}

	// 检查错误次数限制
	if err := CheckCaptchaRetryLimit(ip); err != nil {
		return err
	}

	ctx := context.Background()
	key := fmt.Sprintf("captcha:%s", captchaID)

	// 从Redis获取验证码答案
	correctAnswer, err := db.RDB.Get(ctx, key).Result()
	if err != nil {
		return errors.New("验证码已过期或不存在")
	}

	// 比较答案
	if correctAnswer != answer {
		// 记录错误次数
		IncrCaptchaRetryCount(ip)
		// 验证失败也删除验证码（防止暴力破解）
		db.RDB.Del(ctx, key)
		return errors.New("验证码错误")
	}

	// 验证成功后删除验证码（一次性使用）
	db.RDB.Del(ctx, key)

	return nil
}

