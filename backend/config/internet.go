package config

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type InternetConfig struct {
	Proto       string   `json:"connection_type"`
	IPAddr      string   `json:"ip_addr"`
	Netmask     string   `json:"netmask"`
	Gateway     string   `json:"gateway"`
	IsCustomDNS bool     `json:"is_custom_dns"`
	DNSList     []string `json:"dns_list"`
	IsCustomMAC bool     `json:"is_custom_mac"`
	MACAddr     string   `json:"mac_addr"`
	MTU         int      `json:"mtu"`
}

func NewInternetConfig() *InternetConfig {
	return &InternetConfig{
		Proto:       "dhcp",
		IPAddr:      "",
		Netmask:     "",
		Gateway:     "",
		DNSList:     []string{},
		IsCustomDNS: false,
		MACAddr:     "",
		IsCustomMAC: false,
		MTU:         1500,
	}
}

func LoadInternetConfig() (*InternetConfig, error) {
	config := &InternetConfig{}

	// Load protocol type
	proto, err := getUCIValue(fmt.Sprintf("network.%s.proto", WanInterfaceName))
	if err != nil {
		return nil, errors.Join(LoadConfigError, err)
	}
	config.Proto = proto

	// Load IP address, netmask, and gateway for static configuration
	if proto == "static" {
		ipaddr, err := getUCIValue(fmt.Sprintf("network.%s.ipaddr", WanInterfaceName))
		if err != nil {
			return nil, errors.Join(LoadConfigError, err)
		}
		config.IPAddr = ipaddr

		netmask, err := getUCIValue(fmt.Sprintf("network.%s.netmask", WanInterfaceName))
		if err != nil {
			return nil, errors.Join(LoadConfigError, err)
		}
		config.Netmask = netmask

		gateway, err := getUCIValue(fmt.Sprintf("network.%s.gateway", WanInterfaceName))
		if err != nil {
			return nil, errors.Join(LoadConfigError, err)
		}
		config.Gateway = gateway
	}

	// Load DNS settings
	dns, err := getUCIValue(fmt.Sprintf("network.%s.dns", WanInterfaceName))
	if err == nil && dns != "" {
		config.IsCustomDNS = true
		config.DNSList = strings.Fields(dns)
	}

	// Load MAC address
	macAddr, err := getUCIValue(fmt.Sprintf("network.%s.macAddr", WanInterfaceName))
	if err == nil && macAddr != "" {
		config.IsCustomMAC = true
		config.MACAddr = macAddr
	}

	// Load MTU
	mtu, err := getUCIValue(fmt.Sprintf("network.%s.mtu", WanInterfaceName))
	if err == nil && mtu != "" {
		mtuValue, err := strconv.Atoi(mtu)
		if err == nil {
			config.MTU = mtuValue
		}
	}

	return config, nil
}

func ApplyInternetConfig(config *InternetConfig) error {
	if config.Proto != "dhcp" && config.Proto != "static" {
		return errors.Join(InvalidConfigError, fmt.Errorf("invalid internet protocol: %s", config.Proto))
	}

	var applyCommands []string

	applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.proto=%s", WanInterfaceName, config.Proto))

	if config.Proto == "static" {
		if config.IPAddr != "" {
			err := checkValidIP(config.IPAddr)

			if err != nil {
				return errors.Join(InvalidConfigError, err)
			}

			applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.ipaddr=%s", WanInterfaceName, config.IPAddr))
		} else {
			return errors.Join(InvalidConfigError, fmt.Errorf("ip address value is required"))
		}

		if config.Netmask != "" {
			err := checkValidIP(config.Netmask)

			if err != nil {
				return errors.Join(InvalidConfigError, err)
			}

			applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.netmask=%s", WanInterfaceName, config.Netmask))
		} else {
			return errors.Join(InvalidConfigError, fmt.Errorf("netmask value is required"))
		}

		if config.Gateway != "" {
			err := checkValidIP(config.Gateway)

			if err != nil {
				return errors.Join(InvalidConfigError, err)
			}

			applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.gateway=%s", WanInterfaceName, config.Gateway))
		} else {
			return errors.Join(InvalidConfigError, fmt.Errorf("gateway value is required"))
		}
	}

	if config.IsCustomDNS && len(config.DNSList) > 0 {
		for _, dns := range config.DNSList {
			err := checkValidIP(dns)

			if err != nil {
				return errors.Join(InvalidConfigError, err)
			}
		}

		dnsValues := strings.Join(config.DNSList, " ")
		applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.dns='%s'", WanInterfaceName, dnsValues))
	}

	if config.IsCustomMAC && config.MACAddr != "" {
		err := checkValidMAC(config.MACAddr)

		if err != nil {
			return errors.Join(InvalidConfigError, err)
		}

		applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.macaddr=%s", WanInterfaceName, config.MACAddr))
	}

	if config.MTU > 0 {
		applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.mtu=%d", WanInterfaceName, config.MTU))
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
