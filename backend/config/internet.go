package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ConnectionType 인터넷 연결 타입
type ConnectionType string

const (
	DynamicIP ConnectionType = "dhcp"
	StaticIP  ConnectionType = "static"
)

const ConfigPath string = "/etc/config/custom_network/internet"

// InternetConfig OpenWrt 인터넷 연결 설정을 위한 구조체
type InternetConfig struct {
	// 네트워크 프로토콜 설정
	Proto ConnectionType // 'dhcp' 또는 'static'

	// IP 설정
	IpAddr  string // xxx.xxx.xxx.xxx 형식
	Netmask string // xxx.xxx.xxx.xxx 형식
	Gateway string // xxx.xxx.xxx.xxx 형식

	// DNS 설정
	DNS          []string // DNS 서버 목록
	UseCustomDNS bool     // DNS 커스텀 설정 여부

	// MAC 주소 설정
	MACAddr  string // XX:XX:XX:XX:XX:XX 형식
	CloneMAC bool   // MAC 주소 변경 사용 여부

	// MTU 설정
	MTU int // MTU 값
}

// NewInternetConfig 기본값으로 초기화된 인터넷 설정 구조체 생성
func NewInternetConfig() *InternetConfig {
	return &InternetConfig{
		Proto:        DynamicIP,
		IpAddr:       "",
		Netmask:      "",
		Gateway:      "",
		DNS:          []string{},
		UseCustomDNS: false,
		MACAddr:      "",
		CloneMAC:     false,
		MTU:          1500,
	}
}

// SaveToFile OpenWrt UCI 형식으로 설정을 파일에 저장
func (i *InternetConfig) SaveToFile() error {
	var err error
	f, err := os.Create(ConfigPath)

	if err != nil {
		return fmt.Errorf("failed to create config file: %v", err)
	}

	defer f.Close()

	writer := bufio.NewWriter(f)

	// 기본 인터페이스 설정
	_, err = writer.WriteString("config interface 'wan'\n")

	if err != nil {
		return err
	}

	_, err = writer.WriteString(fmt.Sprintf("\toption proto '%s'\n", i.Proto))

	if err != nil {
		return err
	}

	// IP 관련 설정
	if i.Proto == StaticIP {
		if i.IpAddr != "" {
			_, err := writer.WriteString(fmt.Sprintf("\toption ipaddr '%s'\n", i.IpAddr))

			if err != nil {
				return err
			}
		}
		if i.Netmask != "" {
			_, err := writer.WriteString(fmt.Sprintf("\toption netmask '%s'\n", i.Netmask))

			if err != nil {
				return err
			}
		}
		if i.Gateway != "" {
			_, err := writer.WriteString(fmt.Sprintf("\toption gateway '%s'\n", i.Gateway))

			if err != nil {
				return err
			}
		}
	}

	// DNS 설정
	if i.UseCustomDNS {
		_, err := writer.WriteString("\toption peerdns '0'\n")

		if err != nil {
			return err
		}

		for idx, dns := range i.DNS {
			_, err := writer.WriteString(fmt.Sprintf("\toption dns_%d '%s'\n", idx+1, dns))
			if err != nil {
				return err
			}
		}
	}

	// MAC 주소 설정
	if i.CloneMAC && i.MACAddr != "" {
		_, err := writer.WriteString(fmt.Sprintf("\toption macaddr '%s'\n", i.MACAddr))
		if err != nil {
			return err
		}
	}

	// MTU 설정
	if i.MTU != 1500 {
		_, err := writer.WriteString(fmt.Sprintf("\toption mtu '%d'\n", i.MTU))
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

// LoadFromFile OpenWrt UCI 형식의 설정 파일에서 설정 읽기
func LoadFromFile() (*InternetConfig, error) {
	f, err := os.Open(ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer f.Close()

	config := NewInternetConfig()
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// option 라인 파싱
		if strings.HasPrefix(line, "option") {
			parts := strings.Fields(line)
			if len(parts) < 3 {
				continue
			}

			// 따옴표 제거
			value := strings.Trim(strings.Join(parts[2:], " "), "'")

			switch parts[1] {
			case "proto":
				config.Proto = ConnectionType(value)
			case "ipaddr":
				config.IpAddr = value
			case "netmask":
				config.Netmask = value
			case "gateway":
				config.Gateway = value
			case "peerdns":
				config.UseCustomDNS = value != "0"
			case "macaddr":
				config.MACAddr = value
				config.CloneMAC = true
			case "mtu":
				_, err := fmt.Sscanf(value, "%d", &config.MTU)
				if err != nil {
					return nil, err
				}
			}

			// DNS 서버 파싱
			if strings.HasPrefix(parts[1], "dns_") {
				config.DNS = append(config.DNS, value)
			}
		}
	}

	return config, scanner.Err()
}

// SetIPAddress IP 주소 설정 (xxx.xxx.xxx.xxx 형식)
func (i *InternetConfig) SetIPAddress(ip string) error {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return fmt.Errorf("invalid IP address format")
	}
	i.IpAddr = ip
	return nil
}

// SetMACAddress MAC 주소 설정 (XX:XX:XX:XX:XX:XX 형식)
func (i *InternetConfig) SetMACAddress(mac string) error {
	parts := strings.Split(strings.ToUpper(mac), ":")
	if len(parts) != 6 {
		return fmt.Errorf("invalid MAC address format")
	}
	i.MACAddr = mac
	i.CloneMAC = true
	return nil
}

// SetDNSServers DNS 서버 목록 설정
func (i *InternetConfig) SetDNSServers(servers []string) {
	i.DNS = servers
}

// EnableDynamicIP DynamicIP 설정으로 변경
func (i *InternetConfig) EnableDynamicIP() {
	i.Proto = DynamicIP
	i.IpAddr = ""
	i.Gateway = ""
}

// EnableStaticIP 고정 IP 설정으로 변경
func (i *InternetConfig) EnableStaticIP() {
	i.Proto = StaticIP
}
