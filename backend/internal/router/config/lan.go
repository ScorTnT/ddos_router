package config

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type LANConfig struct {
	IPAddress string
	Netmask   string
}

func NewLANConfig() *LANConfig {
	return &LANConfig{
		IPAddress: "192.168.0.1",
		Netmask:   "255.255.255.0",
	}
}

func ApplyLANConfig(config *LANConfig) error {
	var applyCommands []string

	if config.IPAddress != "" {
		err := checkValidIP(config.IPAddress)

		if err != nil {
			return errors.Join(ErrInvalidConfig, err)
		}

		applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.ipaddr=%s", LANInterfaceName, config.IPAddress))
	} else {
		return errors.Join(ErrInvalidConfig, fmt.Errorf("ip address value is required"))
	}

	if config.Netmask != "" {
		err := checkValidIP(config.Netmask)

		if err != nil {
			return errors.Join(ErrInvalidConfig, err)
		}

		applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.netmask=%s", LANInterfaceName, config.Netmask))
	} else {
		return errors.Join(ErrInvalidConfig, fmt.Errorf("netmask value is required"))
	}

	for _, command := range applyCommands {
		cmd := exec.Command("uci", strings.Split(command, " ")...)
		if err := cmd.Run(); err != nil {
			return errors.Join(ErrInvalidConfig, fmt.Errorf("failed to run uci command, %s: %v", cmd, err))
		}
	}

	if err := exec.Command("uci", "commit", "network").Run(); err != nil {
		return errors.Join(ErrInvalidConfig, fmt.Errorf("failed to commit changes: %v", err))
	}

	if err := exec.Command("/etc/init.d/network", "restart").Run(); err != nil {
		return errors.Join(ErrInvalidConfig, fmt.Errorf("failed to restart network: %v", err))
	}

	return nil
}

func LoadLANConfig() (*LANConfig, error) {
	config := NewLANConfig()

	// IP 주소 가져오기
	ipAddress, err := getUCIValue(fmt.Sprintf("network.%s.ipaddr", LANInterfaceName))
	if err == nil {
		config.IPAddress = ipAddress
	}

	// 넷마스크 가져오기
	netmask, err := getUCIValue(fmt.Sprintf("network.%s.netmask", LANInterfaceName))
	if err == nil {
		config.Netmask = netmask
	}

	return config, nil
}
