package config

import (
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
		viper.SetConfigFile(path)
		viper.SetConfigType("yaml")
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			loadErr = err
			return
		}

		appConfig = &Config{}
		if err := viper.Unmarshal(appConfig); err != nil {
			loadErr = err
			return
		}

		setDefaults()
	})
	return loadErr
}

func setDefaults() {
	if appConfig.Server.Port == "" {
		appConfig.Server.Port = "8080"
	}
	if appConfig.Server.Mode == "" {
		appConfig.Server.Mode = "debug"
	}
	if appConfig.Server.Timezone == "" {
		appConfig.Server.Timezone = "Asia/Shanghai"
	}
	if len(appConfig.Server.AllowedOrigins) == 0 {
		appConfig.Server.AllowedOrigins = []string{"http://localhost:3000"}
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
}

func Get() *Config {
	return appConfig
}
