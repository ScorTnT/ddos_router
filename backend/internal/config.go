package internal

import (
	"os"

	"gopkg.in/yaml.v3"
)

type WebAPIConfig struct {
	ListenAddr string `yaml:"listen_addr"`
}

type Config struct {
	WebAPI WebAPIConfig `yaml:"web_api"`
}

const configPath = "config.yaml"

func LoadConfig() (*Config, error) {
	f, err := os.Open(configPath)

	if err != nil {
		return nil, err
	}

	var config Config
	decoder := yaml.NewDecoder(f)

	if err = decoder.Decode(&config); err != nil {
		return nil, err
	}

	if err = f.Close(); err != nil {
		return nil, err
	}

	return &config, nil
}
