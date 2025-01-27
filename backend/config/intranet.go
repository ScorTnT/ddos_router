package config

import (
	"fmt"
	"os/exec"
	"strings"
)

type IntranetConfig struct {
	IpAddr  string
	Netmask string
}

func NewIntranetConfig() *IntranetConfig {
	return &IntranetConfig{
		IpAddr:  "192.168.0.1",
		Netmask: "255.255.255.0",
	}
}

func (i *IntranetConfig) SaveIntranetConfig() error {
	// UCI 설정을 위한 명령어들
	commands := [][]string{
		{"uci", "set", "network.lan.device=br-lan"},
		{"uci", "set", "network.lan.proto=static"},
		{"uci", "set", fmt.Sprintf("network.lan.ipaddr=%s", i.IpAddr)},
		{"uci", "set", fmt.Sprintf("network.lan.netmask=%s", i.Netmask)},
		{"uci", "commit", "network"},
	}

	// 각 명령어 실행
	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err != nil {
			return fmt.Errorf("failed to execute UCI command %v: %v", cmd, err)
		}
	}

	return nil
}

func LoadIntranetConfig() (*IntranetConfig, error) {
	config := NewIntranetConfig()

	// IP 주소 가져오기
	ipCmd := exec.Command("uci", "get", "network.lan.ipaddr")
	ipOut, err := ipCmd.Output()
	if err == nil {
		config.IpAddr = strings.TrimSpace(string(ipOut))
	}

	// 넷마스크 가져오기
	netmaskCmd := exec.Command("uci", "get", "network.lan.netmask")
	netmaskOut, err := netmaskCmd.Output()
	if err == nil {
		config.Netmask = strings.TrimSpace(string(netmaskOut))
	}

	return config, nil
}

func (i *IntranetConfig) SetIPAddress(ip string) error {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return fmt.Errorf("invalid IP address format")
	}
	i.IpAddr = ip
	return nil
}

func (i *IntranetConfig) SetNetmask(netmask string) error {
	parts := strings.Split(netmask, ".")
	if len(parts) != 4 {
		return fmt.Errorf("invalid netmask format")
	}
	i.Netmask = netmask
	return nil
}
