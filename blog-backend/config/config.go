package config

import (
	"github.com/spf13/viper"
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
	JWT struct {
		Secret      string `mapstructure:"secret"`
		ExpireHours int    `mapstructure:"expire_hours"`
	} `mapstructure:"jwt"`
	Log struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"log"`
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
	return nil
}
