package config

import "github.com/spf13/viper"

type Config struct {
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
