package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	WanInterfaceName = "wan"
	LanInterfaceName = "lan"
)

var (
	InvalidConfigError = errors.New("invalid config data")
	LoadConfigError    = errors.New("failed to load config")
)

type Config struct {
	InitialAdmin struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"initial_admin"`
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
