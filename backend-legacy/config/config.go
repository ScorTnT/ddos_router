package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	WanInterfaceName = "wan"
	LanInterfaceName = "lan"
)

var (
	ErrInvalidConfig = errors.New("invalid config data")
	ErrLoadConfig    = errors.New("failed to load config")
)

type Config struct {
	InitialAdmin struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"initial_admin"`
	WebAPI struct {
		ListenAddress string `yaml:"listen_address"`
	} `yaml:"web_api"`
}

var configPath = "config.yml"

func ReadConfig() (*Config, error) {
	f, err := os.Open(configPath)

	if err != nil {
		return nil, err
	}

	decoder := yaml.NewDecoder(f)
	var config Config
	err = decoder.Decode(&config)

	if err != nil {
		return nil, err
	}

	err = f.Close()

	if err != nil {
		return nil, err
	}

	return &config, nil
}
