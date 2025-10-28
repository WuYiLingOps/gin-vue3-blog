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

// GenerateCaptcha 生成验证码
func GenerateCaptcha() (*CaptchaResponse, error) {
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
		// 存储到Redis，设置5分钟过期
		key := fmt.Sprintf("captcha:%s", id)
		err = db.RDB.Set(ctx, key, answer, 5*time.Minute).Err()
		if err != nil {
			return nil, err
		}
	}

	return &CaptchaResponse{
		CaptchaID: id,
		ImageData: b64s,
	}, nil
}

// VerifyCaptcha 验证验证码
func VerifyCaptcha(captchaID, answer string) error {
	if captchaID == "" || answer == "" {
		return errors.New("验证码ID和答案不能为空")
	}

	ctx := context.Background()
	key := fmt.Sprintf("captcha:%s", captchaID)

	// 从Redis获取验证码答案
	correctAnswer, err := db.RDB.Get(ctx, key).Result()
	if err != nil {
		return errors.New("验证码已过期或不存在")
	}

	// 验证成功后删除验证码（一次性使用）
	db.RDB.Del(ctx, key)

	// 比较答案（不区分大小写）
	if correctAnswer != answer {
		return errors.New("验证码错误")
	}

	return nil
}

