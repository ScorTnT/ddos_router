package internal

import (
	"os"

	"gopkg.in/yaml.v3"
)

type WebAPIConfig struct {
	ListenAddr string `yaml:"listen_addr"`
}

type InferfaceConfig struct {
	WANInterfaceName string `yaml:"wan_interface_name"`
	LANInterfaceName string `yaml:"lan_interface_name"`
}

type SnortAlertScannerConfig struct {
	SnortLogPath string `yaml:"snort_log_path"`
	BufferSize   int    `yaml:"buffer_size"`
}

type ProtectionConfig struct {
	ProtectionTTL int `yaml:"protection_ttl"`
	RefreshTick   int `yaml:"refresh_tick"`
}

type Config struct {
	WebAPI     WebAPIConfig            `yaml:"web_api"`
	Interface  InferfaceConfig         `yaml:"interface"`
	Snort      SnortAlertScannerConfig `yaml:"snort"`
	Protection ProtectionConfig        `yaml:"protection"`
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
