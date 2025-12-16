package config

import (
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
	Port           string   `mapstructure:"port"`
	Mode           string   `mapstructure:"mode"`
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	Timezone       string   `mapstructure:"timezone"`
}

type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	ExpireHour int    `mapstructure:"expire_hour"`
}

var (
	appConfig *Config
	once      sync.Once
)

func Load(path string) error {
	var loadErr error
	once.Do(func() {
		appConfig = &Config{}

		// 尝试读取配置文件
		if path != "" {
			viper.SetConfigFile(path)
			viper.SetConfigType("yaml")
			if err := viper.ReadInConfig(); err == nil {
				viper.Unmarshal(appConfig)
			}
		}

		// 环境变量覆盖配置文件
		loadFromEnv()

		// 设置默认值
		setDefaults()
	})
	return loadErr
}

func loadFromEnv() {
	// 数据库配置
	if v := os.Getenv("DB_HOST"); v != "" {
		appConfig.Database.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		appConfig.Database.Port = v
	}
	if v := os.Getenv("DB_USER"); v != "" {
		appConfig.Database.User = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		appConfig.Database.Password = v
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		appConfig.Database.DBName = v
	}

	// JWT 配置
	if v := os.Getenv("JWT_SECRET"); v != "" {
		appConfig.JWT.Secret = v
	}
	if v := os.Getenv("JWT_EXPIRE_HOUR"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			appConfig.JWT.ExpireHour = i
		}
	}

	// 服务器配置
	if v := os.Getenv("SERVER_PORT"); v != "" {
		appConfig.Server.Port = v
	}
	if v := os.Getenv("SERVER_MODE"); v != "" {
		appConfig.Server.Mode = v
	}
	if v := os.Getenv("ALLOWED_ORIGINS"); v != "" {
		appConfig.Server.AllowedOrigins = strings.Split(v, ",")
	}
	if v := os.Getenv("TIMEZONE"); v != "" {
		appConfig.Server.Timezone = v
	}
}

func setDefaults() {
	if appConfig.Server.Port == "" {
		appConfig.Server.Port = "8080"
	}
	if appConfig.Server.Mode == "" {
		appConfig.Server.Mode = "release"
	}
	if appConfig.Server.Timezone == "" {
		appConfig.Server.Timezone = "Asia/Shanghai"
	}
	if len(appConfig.Server.AllowedOrigins) == 0 {
		appConfig.Server.AllowedOrigins = []string{"http://localhost:3000"}
	}
	if appConfig.Database.Host == "" {
		appConfig.Database.Host = "localhost"
	}
	if appConfig.Database.Port == "" {
		appConfig.Database.Port = "3306"
	}
	if appConfig.Database.MaxIdleConns == 0 {
		appConfig.Database.MaxIdleConns = 10
	}
	if appConfig.Database.MaxOpenConns == 0 {
		appConfig.Database.MaxOpenConns = 100
	}
	if appConfig.Database.ConnMaxLifetime == 0 {
		appConfig.Database.ConnMaxLifetime = 3600
	}
	if appConfig.JWT.ExpireHour == 0 {
		appConfig.JWT.ExpireHour = 168
	}
	if appConfig.JWT.Secret == "" {
		appConfig.JWT.Secret = "default-secret-please-change"
	}
}

func Get() *Config {
	return appConfig
}
