package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Crypto   CryptoConfig   `mapstructure:"crypto"`
	Chromium ChromiumConfig `mapstructure:"chromium"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

type CryptoConfig struct {
	SecretKey string `mapstructure:"secret_key"`
}

type ChromiumConfig struct {
	Path          string `mapstructure:"path"`
	MaxConcurrent int    `mapstructure:"max_concurrent"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
	Path  string `mapstructure:"path"`
}

var AppConfig *Config

func Load(path string) error {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	AppConfig = &Config{}
	return viper.Unmarshal(AppConfig)
}

func Get() *Config {
	if AppConfig == nil {
		AppConfig = &Config{
			Server: ServerConfig{Port: 8080, Mode: "release"},
			Database: DatabaseConfig{Path: "./data/monitor.db"},
			Crypto: CryptoConfig{SecretKey: ""},
			Chromium: ChromiumConfig{Path: "", MaxConcurrent: 3},
			Log: LogConfig{Level: "info", Path: "./logs/app.log"},
		}
	}
	return AppConfig
}
