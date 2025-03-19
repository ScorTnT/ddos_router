package router_control

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
	"strings"
)

type ConnectionType string

const (
	DHCP   ConnectionType = "dhcp"
	Static ConnectionType = "static"
)

type WANConfig struct {
	ConnectionType ConnectionType `json:"connection_type"`
	IpAddress      string         `json:"ip_address"`
	Netmask        string         `json:"netmask"`
	Gateway        string         `json:"gateway"`
	DNSList        []string       `json:"dns_list"`
	MACAddress     string         `json:"mac_address"`
	MTUSize        int            `json:"mtu_size"`
}

func NewWANConfig() *WANConfig {
	return &WANConfig{
		ConnectionType: DHCP,
		IpAddress:      "",
		Netmask:        "",
		Gateway:        "",
		DNSList:        []string{},
		MACAddress:     "",
		MTUSize:        1500,
	}
}

func (w *WANConfig) SaveConfig(interfaceName string) error {
	if err := validateWANConfig(w); err != nil {
		return err
	}

	var applyCommands []string

	applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.proto=%s", interfaceName, w.ConnectionType))

	if w.ConnectionType == Static {
		applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.ipaddr=%s", interfaceName, w.IpAddress))
		applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.netmask=%s", interfaceName, w.Netmask))
		applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.gateway=%s", interfaceName, w.Gateway))
	}

	if len(w.DNSList) > 0 {
		applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.dns=%s", interfaceName, strings.Join(w.DNSList, " ")))
	}

	if len(w.MACAddress) > 0 {
		applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.macaddr=%s", interfaceName, w.MACAddress))
	}

	applyCommands = append(applyCommands, fmt.Sprintf("set network.%s.mtu=%d", interfaceName, w.MTUSize))

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

func (w *WANConfig) LoadConfig(interfaceName string) error {
	proto, err := getUCIValue(fmt.Sprintf("network.%s.proto", interfaceName))
	if err != nil {
		return err
	}
	w.ConnectionType = ConnectionType(proto)

	if w.ConnectionType == Static {
		ip, err := getUCIValue(fmt.Sprintf("network.%s.ip", interfaceName))
		if err != nil {
			return err
		}
		w.IpAddress = ip

		netmask, err := getUCIValue(fmt.Sprintf("network.%s.netmask", interfaceName))
		if err != nil {
			return err
		}
		w.Netmask = netmask

		gateway, err := getUCIValue(fmt.Sprintf("network.%s.gateway", interfaceName))
		if err != nil {
			return err
		}
		w.Gateway = gateway
	}

	dns, err := getUCIValue(fmt.Sprintf("network.%s.dns", interfaceName))
	if err != nil {
		return err
	}
	if len(dns) > 0 {
		w.DNSList = strings.Split(dns, ",")
	}

	mac, err := getUCIValue(fmt.Sprintf("network.%s.mac", interfaceName))
	if err != nil {
		return err
	}
	if len(mac) > 0 {
		w.MACAddress = mac
	}

	mtu, err := getUCIValue(fmt.Sprintf("network.%s.mtu", interfaceName))
	if err != nil {
		return err
	}
	if len(mtu) > 0 {
		w.MTUSize = len(mtu)
	}

	return nil
}

func validateWANConfig(config *WANConfig) error {
	if config.ConnectionType != Static && config.ConnectionType != DHCP {
		return errors.New("invalid connection type")
	}

	ip := net.ParseIP(config.IpAddress)
	if ip == nil {
		return fmt.Errorf("invalid ip address: %s", config.IpAddress)
	}

	netmask := net.ParseIP(config.Netmask)
	if netmask == nil {
		return fmt.Errorf("invalid netmask: %s", config.Netmask)
	}

	gateway := net.ParseIP(config.Gateway)
	if gateway == nil {
		return fmt.Errorf("invalid gateway: %s", config.Gateway)
	}

	mac, err := net.ParseMAC(config.MACAddress)
	if err != nil {
		return fmt.Errorf("invalid mac address: %s", config.MACAddress)
	}
	if len(mac) != 6 {
		return fmt.Errorf("invalid mac address: %s", config.MACAddress)
	}

	return nil
}
