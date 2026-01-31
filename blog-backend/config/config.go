/*
 * 项目名称：blog-backend
 * 文件名称：config.go
 * 创建时间：2026-01-31 16:01:41
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：系统配置管理模块，负责加载和管理应用配置，支持多环境配置和环境变量覆盖
 */
package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

// Config 系统配置结构体，包含应用的所有配置项
type Config struct {
	Env string // 当前环境：dev 或 prod

	// App 应用基础配置
	App struct {
		Port int `mapstructure:"port"` // 应用端口号
	} `mapstructure:"app"`

	// Server 服务器配置
	Server struct {
		Mode string `mapstructure:"mode"` // Gin运行模式：debug、release、test
	} `mapstructure:"server"`

	// DB 数据库配置
	DB struct {
		Host     string `mapstructure:"host"`     // 数据库主机地址
		Port     int    `mapstructure:"port"`     // 数据库端口号
		User     string `mapstructure:"user"`     // 数据库用户名
		Password string `mapstructure:"password"` // 数据库密码
		DBName   string `mapstructure:"dbname"`   // 数据库名称
		SSLMode  string `mapstructure:"sslmode"`  // SSL模式
	} `mapstructure:"db"`

	// Redis Redis缓存配置
	Redis struct {
		Host     string `mapstructure:"host"`     // Redis主机地址
		Port     int    `mapstructure:"port"`     // Redis端口号
		Password string `mapstructure:"password"` // Redis密码
		DB       int    `mapstructure:"db"`       // Redis数据库编号
	} `mapstructure:"redis"`

	// GiteeCalendar Gitee日历API配置
	GiteeCalendar struct {
		APIURL string `mapstructure:"api_url"` // gitee-calendar-api 服务地址
	} `mapstructure:"gitee_calendar"`

	// JWT JWT令牌配置
	JWT struct {
		Secret      string `mapstructure:"secret"`       // JWT密钥
		ExpireHours int    `mapstructure:"expire_hours"` // JWT过期时间（小时）
	} `mapstructure:"jwt"`

	// Email 邮件服务配置
	Email struct {
		Host     string `mapstructure:"host"`      // SMTP服务器地址
		Port     int    `mapstructure:"port"`      // SMTP服务器端口
		Username string `mapstructure:"username"`  // 发件人邮箱
		Password string `mapstructure:"password"`  // 邮箱密码或授权码
		FromName string `mapstructure:"from_name"` // 发件人名称
	} `mapstructure:"email"`

	// Log 日志配置
	Log struct {
		Level string `mapstructure:"level"` // 日志级别：debug、info、warn、error
	} `mapstructure:"log"`

	// OSS 阿里云OSS对象存储配置
	OSS struct {
		Endpoint        string `mapstructure:"endpoint"`          // OSS服务端点
		AccessKeyID     string `mapstructure:"access_key_id"`     // 访问密钥ID
		AccessKeySecret string `mapstructure:"access_key_secret"` // 访问密钥Secret
		BucketName      string `mapstructure:"bucket_name"`       // 存储桶名称
		Domain          string `mapstructure:"domain"`            // 自定义域名（可选）
	} `mapstructure:"oss"`

	// COS 腾讯云COS对象存储配置
	COS struct {
		BucketURL string `mapstructure:"bucket_url"` // 存储桶URL，形如：https://<bucket>.cos.<region>.myqcloud.com
		SecretID  string `mapstructure:"secret_id"`  // 访问密钥ID
		SecretKey string `mapstructure:"secret_key"` // 访问密钥Key
		Domain    string `mapstructure:"domain"`     // 自定义域名（可选）
	} `mapstructure:"cos"`

	// Security 安全配置
	Security struct {
		AdminIPWhitelist []string `mapstructure:"admin_ip_whitelist"` // 管理员IP白名单列表
	} `mapstructure:"security"`
}

// Cfg 全局配置实例
var Cfg *Config

// LoadConfig 从指定路径加载配置文件
// 参数:
//   - path: 配置文件路径
//
// 返回:
//   - error: 加载失败时返回错误
func LoadConfig(path string) error {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return err
	}
	Cfg = &cfg

	return nil
}

// loadEnvOverrides 从环境变量（或 .env.config.<env> 文件）中覆盖敏感配置
// 该方法允许通过环境变量或 .env.config.<env> 文件来覆盖配置文件中的敏感信息（如密码、密钥等）
// 优先级：环境变量 > .env.config.<env> 文件 > 配置文件
// 参数:
//   - env: 当前环境名称（dev 或 prod）
func loadEnvOverrides(env string) {
	// 如果配置未加载，直接返回
	if Cfg == nil {
		return
	}

	// 尝试加载同级目录下的 .env.config.<env> 文件（不存在则忽略）
	// 例如：.env.config.dev 或 .env.config.prod
	_ = gotenv.Load(".env.config." + env)

	// 数据库配置覆盖
	if v := os.Getenv("DB_HOST"); v != "" {
		Cfg.DB.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			Cfg.DB.Port = p
		}
	}
	if v := os.Getenv("DB_USER"); v != "" {
		Cfg.DB.User = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		Cfg.DB.Password = v
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		Cfg.DB.DBName = v
	}

	// Redis配置覆盖
	if v := os.Getenv("REDIS_HOST"); v != "" {
		Cfg.Redis.Host = v
	}
	if v := os.Getenv("REDIS_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			Cfg.Redis.Port = p
		}
	}
	if v := os.Getenv("REDIS_PASSWORD"); v != "" {
		Cfg.Redis.Password = v
	}

	// Gitee Calendar API配置覆盖
	if v := os.Getenv("GITEE_CALENDAR_API_URL"); v != "" {
		Cfg.GiteeCalendar.APIURL = v
	}

	// JWT配置覆盖
	if v := os.Getenv("JWT_SECRET"); v != "" {
		Cfg.JWT.Secret = v
	}
	if v := os.Getenv("JWT_EXPIRE_HOURS"); v != "" {
		if h, err := strconv.Atoi(v); err == nil {
			Cfg.JWT.ExpireHours = h
		}
	}

	// 邮件服务配置覆盖
	if v := os.Getenv("EMAIL_HOST"); v != "" {
		Cfg.Email.Host = v
	}
	if v := os.Getenv("EMAIL_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			Cfg.Email.Port = p
		}
	}
	if v := os.Getenv("EMAIL_USERNAME"); v != "" {
		Cfg.Email.Username = v
	}
	if v := os.Getenv("EMAIL_PASSWORD"); v != "" {
		Cfg.Email.Password = v
	}

	// 阿里云OSS配置覆盖
	if v := os.Getenv("OSS_ENDPOINT"); v != "" {
		Cfg.OSS.Endpoint = v
	}
	if v := os.Getenv("OSS_ACCESS_KEY_ID"); v != "" {
		Cfg.OSS.AccessKeyID = v
	}
	if v := os.Getenv("OSS_ACCESS_KEY_SECRET"); v != "" {
		Cfg.OSS.AccessKeySecret = v
	}
	if v := os.Getenv("OSS_BUCKET_NAME"); v != "" {
		Cfg.OSS.BucketName = v
	}
	if v := os.Getenv("OSS_DOMAIN"); v != "" {
		Cfg.OSS.Domain = v
	}

	// 腾讯云COS配置覆盖
	if v := os.Getenv("COS_BUCKET_URL"); v != "" {
		Cfg.COS.BucketURL = v
	}
	if v := os.Getenv("COS_SECRET_ID"); v != "" {
		Cfg.COS.SecretID = v
	}
	if v := os.Getenv("COS_SECRET_KEY"); v != "" {
		Cfg.COS.SecretKey = v
	}
	if v := os.Getenv("COS_DOMAIN"); v != "" {
		Cfg.COS.Domain = v
	}
}

// LoadConfigByEnv 根据 config.yml 中的 env 字段加载对应环境的配置
// 加载流程：
//  1. 读取 config.yml 获取环境标识（env）
//  2. 根据环境标识加载对应的配置文件（config-dev.yml 或 config-prod.yml）
//  3. 允许通过环境变量或 .env.config.<env> 文件覆盖敏感配置
//
// 返回:
//   - error: 加载失败时返回错误
func LoadConfigByEnv() error {
	// 先读取 config.yml 获取环境配置
	v := viper.New()
	v.SetConfigFile("./config/config.yml")
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	// 获取环境标识，默认为 dev（开发环境）
	env := v.GetString("env")
	if env == "" {
		env = "dev" // 默认开发环境
	}

	// 根据环境加载对应的配置文件
	// 例如：config-dev.yml 或 config-prod.yml
	configPath := "./config/config-" + env + ".yml"
	if err := LoadConfig(configPath); err != nil {
		return err
	}

	// 检查配置是否成功加载
	if Cfg == nil {
		return fmt.Errorf("配置加载失败：Cfg 为 nil")
	}

	// 保存环境变量到配置中，供其他模块使用
	Cfg.Env = env

	// 允许通过环境变量或 .env.config.<env> 文件覆盖敏感信息
	// 这样可以避免将敏感信息（如密码、密钥）提交到版本控制系统
	loadEnvOverrides(env)
	return nil
}
