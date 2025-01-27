package config

import (
	"errors"
	"fmt"
	"net"
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

const interfaceName = "wan"

var (
	InvalidConfigError = errors.New("invalid config data")
	LoadConfigError    = errors.New("failed to load config")
)

func LoadInternetConfig() (*InternetConfig, error) {
	config := &InternetConfig{}

	// Load protocol type
	proto, err := getUCIValue(fmt.Sprintf("network.%s.proto", interfaceName))
	if err != nil {
		return nil, errors.Join(LoadConfigError, err)
	}
	config.Proto = proto

	// Load IP address, netmask, and gateway for static configuration
	if proto == "static" {
		ipaddr, err := getUCIValue(fmt.Sprintf("network.%s.ipaddr", interfaceName))
		if err != nil {
			return nil, errors.Join(LoadConfigError, err)
		}
		config.IPAddr = ipaddr

		netmask, err := getUCIValue(fmt.Sprintf("network.%s.netmask", interfaceName))
		if err != nil {
			return nil, errors.Join(LoadConfigError, err)
		}
		config.Netmask = netmask

		gateway, err := getUCIValue(fmt.Sprintf("network.%s.gateway", interfaceName))
		if err != nil {
			return nil, errors.Join(LoadConfigError, err)
		}
		config.Gateway = gateway
	}

	// Load DNS settings
	dns, err := getUCIValue(fmt.Sprintf("network.%s.dns", interfaceName))
	if err == nil && dns != "" {
		config.IsCustomDNS = true
		config.DNSList = strings.Fields(dns)
	}

	// Load MAC address
	macaddr, err := getUCIValue(fmt.Sprintf("network.%s.macaddr", interfaceName))
	if err == nil && macaddr != "" {
		config.IsCustomMAC = true
		config.MACAddr = macaddr
	}

	// Load MTU
	mtu, err := getUCIValue(fmt.Sprintf("network.%s.mtu", interfaceName))
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

	applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.proto='%s'", interfaceName, config.Proto))

	if config.Proto == "static" {
		if config.IPAddr != "" {
			err := checkValidIP(config.IPAddr)

			if err != nil {
				return errors.Join(InvalidConfigError, err)
			}

			applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.ipaddr='%s'", interfaceName, config.IPAddr))
		} else {
			return errors.Join(InvalidConfigError, fmt.Errorf("ip address value is required"))
		}

		if config.Netmask != "" {
			err := checkValidIP(config.Netmask)

			if err != nil {
				return errors.Join(InvalidConfigError, err)
			}

			applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.netmask='%s'", interfaceName, config.Netmask))
		} else {
			return errors.Join(InvalidConfigError, fmt.Errorf("netmask value is required"))
		}

		if config.Gateway != "" {
			err := checkValidIP(config.Gateway)

			if err != nil {
				return errors.Join(InvalidConfigError, err)
			}

			applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.gateway='%s'", interfaceName, config.Gateway))
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
		applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.dns=%s", interfaceName, dnsValues))
	}

	if config.IsCustomMAC && config.MACAddr != "" {
		err := checkValidMAC(config.MACAddr)

		if err != nil {
			return errors.Join(InvalidConfigError, err)
		}

		applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.macaddr=%s", interfaceName, config.MACAddr))
	}

	if config.MTU > 0 {
		applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.mtu=%d", interfaceName, config.MTU))
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

func getUCIValue(path string) (string, error) {
	cmd := exec.Command("uci", "get", path)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get UCI value for %s: %v", path, err)
	}
	return strings.TrimSpace(string(output)), nil
}

func checkValidIP(ip string) error {
	parsedIP := net.ParseIP(ip)

	if parsedIP != nil && parsedIP.To4() != nil {
		return nil
	}

	return fmt.Errorf("invalid IP: %s", ip)
}

func checkValidMAC(mac string) error {
	parsedMAC, err := net.ParseMAC(mac)

	if err != nil {
		return nil
	}

	if len(parsedMAC) == 6 {
		return nil
	}

	return fmt.Errorf("invalid MAC: %s", mac)
}
