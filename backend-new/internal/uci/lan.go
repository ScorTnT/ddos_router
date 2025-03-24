package config

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

type LANConfig struct {
	IPAddress string `json:"IPAddress"`
	Netmask   string `json:"Netmask"`
}

func NewLANConfig() *LANConfig {
	return &LANConfig{
		IPAddress: "192.168.0.1",
		Netmask:   "255.255.255.0",
	}
}

func (l *LANConfig) SaveConfig(interfaceName string) error {
	if err := validateLANConfig(l); err != nil {
		return err
	}

	var applyCommands []string

	applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.ipaddr=%s", interfaceName, l.IPAddress))
	applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.netmask=%s", interfaceName, l.Netmask))

	for _, command := range applyCommands {
		cmd := exec.Command("uci", strings.Split(command, " ")...)

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	if err := exec.Command("uci", "commit", "network").Run(); err != nil {
		return err
	}

	if err := exec.Command("/etc/init.d/network", "restart").Run(); err != nil {
		return err
	}

	return nil
}

func (l *LANConfig) LoadConfig(interfaceName string) error {
	ip, err := getUCIValue(fmt.Sprintf("network.%s.ipaddr", interfaceName))
	if err != nil {
		return err
	}
	l.IPAddress = ip

	netmask, err := getUCIValue(fmt.Sprintf("network.%s.netmask", interfaceName))
	if err != nil {
		return err
	}
	l.Netmask = netmask

	return nil
}

func validateLANConfig(config *LANConfig) error {
	ip := net.ParseIP(config.IPAddress)

	if ip == nil {
		return fmt.Errorf("invalid IP address: %s", config.IPAddress)
	}

	netmask := net.ParseIP(config.Netmask)
	if netmask == nil {
		return fmt.Errorf("invalid netmask: %s", config.Netmask)
	}

	return nil
}
