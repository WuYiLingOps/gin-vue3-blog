package config

import (
	"os"
	"strconv"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

type Config struct {
	Env string // 当前环境：dev 或 prod
	App struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"app"`
	Server struct {
		Mode string `mapstructure:"mode"`
	} `mapstructure:"server"`
	DB struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
	} `mapstructure:"db"`
	Redis struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`
	GiteeCalendar struct {
		APIURL string `mapstructure:"api_url"` // gitee-calendar-api 地址
	} `mapstructure:"gitee_calendar"`
	JWT struct {
		Secret      string `mapstructure:"secret"`
		ExpireHours int    `mapstructure:"expire_hours"`
	} `mapstructure:"jwt"`
	Email struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		FromName string `mapstructure:"from_name"`
	} `mapstructure:"email"`
	Log struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"log"`
	OSS struct {
		Endpoint        string `mapstructure:"endpoint"`
		AccessKeyID     string `mapstructure:"access_key_id"`
		AccessKeySecret string `mapstructure:"access_key_secret"`
		BucketName      string `mapstructure:"bucket_name"`
		Domain          string `mapstructure:"domain"` // 自定义域名（可选）
	} `mapstructure:"oss"`
	COS struct {
		// BucketURL 形如：https://<bucket>.cos.<region>.myqcloud.com
		BucketURL string `mapstructure:"bucket_url"`
		SecretID  string `mapstructure:"secret_id"`
		SecretKey string `mapstructure:"secret_key"`
		Domain    string `mapstructure:"domain"` // 自定义域名（可选）
	} `mapstructure:"cos"`
	Security struct {
		AdminIPWhitelist []string `mapstructure:"admin_ip_whitelist"` // 管理员IP白名单
	} `mapstructure:"security"`
}

var Cfg *Config

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
func loadEnvOverrides(env string) {
	// 尝试加载同级目录下的 .env.config.<env> 文件（不存在则忽略）
	_ = gotenv.Load(".env.config." + env)

	// 数据库
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

	// Redis
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

	// Gitee Calendar API
	if v := os.Getenv("GITEE_CALENDAR_API_URL"); v != "" {
		Cfg.GiteeCalendar.APIURL = v
	}

	// JWT
	if v := os.Getenv("JWT_SECRET"); v != "" {
		Cfg.JWT.Secret = v
	}
	if v := os.Getenv("JWT_EXPIRE_HOURS"); v != "" {
		if h, err := strconv.Atoi(v); err == nil {
			Cfg.JWT.ExpireHours = h
		}
	}

	// 邮件
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

	// OSS
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

	// COS
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
func LoadConfigByEnv() error {
	// 先读取 config.yml 获取环境配置
	v := viper.New()
	v.SetConfigFile("./config/config.yml")
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	env := v.GetString("env")
	if env == "" {
		env = "dev" // 默认开发环境
	}

	// 根据环境加载对应的配置文件
	configPath := "./config/config-" + env + ".yml"
	if err := LoadConfig(configPath); err != nil {
		return err
	}

	// 保存环境变量到配置中
	Cfg.Env = env

	// 允许通过环境变量 / .env.config.<env> 覆盖敏感信息
	loadEnvOverrides(env)
	return nil
}
