package config

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type IntranetConfig struct {
	IPAddress string
	Netmask   string
}

func NewIntranetConfig() *IntranetConfig {
	return &IntranetConfig{
		IPAddress: "192.168.0.1",
		Netmask:   "255.255.255.0",
	}
}

func ApplyIntranetConfig(config *IntranetConfig) error {
	var applyCommands []string

	if config.IPAddress != "" {
		err := checkValidIP(config.IPAddress)

		if err != nil {
			return errors.Join(InvalidConfigError, err)
		}

		applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.ipaddr=%s", LanInterfaceName, config.IPAddress))
	} else {
		return errors.Join(InvalidConfigError, fmt.Errorf("ip address value is required"))
	}

	if config.Netmask != "" {
		err := checkValidIP(config.Netmask)

		if err != nil {
			return errors.Join(InvalidConfigError, err)
		}

		applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.netmask=%s", LanInterfaceName, config.Netmask))
	} else {
		return errors.Join(InvalidConfigError, fmt.Errorf("netmask value is required"))
	}

	for _, command := range applyCommands {
		cmd := exec.Command("uci", strings.Split(command, " ")...)
		if err := cmd.Run(); err != nil {
			return errors.Join(InvalidConfigError, fmt.Errorf("failed to run uci command, %s: %v", cmd, err))
		}
	}

	if err := exec.Command("uci", "commit", "network").Run(); err != nil {
		return errors.Join(InvalidConfigError, fmt.Errorf("failed to commit changes: %v", err))
	}

	if err := exec.Command("/etc/init.d/network", "restart").Run(); err != nil {
		return errors.Join(InvalidConfigError, fmt.Errorf("failed to restart network: %v", err))
	}

	return nil
}

func LoadIntranetConfig() (*IntranetConfig, error) {
	config := NewIntranetConfig()

	// IP 주소 가져오기
	ipAddress, err := getUCIValue(fmt.Sprintf("network.%s.ipaddr", LanInterfaceName))
	if err == nil {
		config.IPAddress = ipAddress
	}

	// 넷마스크 가져오기
	netmask, err := getUCIValue(fmt.Sprintf("network.%s.netmask", LanInterfaceName))
	if err == nil {
		config.Netmask = netmask
	}

	return config, nil
}
