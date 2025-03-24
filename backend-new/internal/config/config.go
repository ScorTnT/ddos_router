package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int    `mapstructure:"port"`
		Host string `mapstructure:"host"`
	}
	Ethernet struct {
		WANInterface string `mapstructure:"wan_interface"`
		LANInterface string `mapstructure:"lan_interface"`
	}
}

var (
	instance *Config
	once     sync.Once
	mu       sync.RWMutex
)

func GetConfig() *Config {
	once.Do(func() {
		cfg, err := loadConfig()
		if err != nil {
			panic(
				fmt.Errorf("failed to load configuration: %v", err),
			)
		}
		instance = cfg
	})

	return instance
}

func ReloadConfig() error {
	mu.Lock()
	defer mu.Unlock()

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	instance = cfg

	return nil
}

func loadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetEnvPrefix("OPENWRT-IPS")
	v.AutomaticEnv()

	// 기본값 설정
	v.SetDefault("server.port", 8000)
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("ethernet.wan_interface", "eth0")
	v.SetDefault("ethernet.lan_interface", "br-lan")

	// 설정 파일 읽기 시도
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 파일이 없으면 기본값으로 config.yaml 생성
			if writeErr := v.SafeWriteConfigAs("config.yaml"); writeErr != nil {
				return nil, fmt.Errorf("failed to create default config file: %w", writeErr)
			}
		} else {
			return nil, err
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
