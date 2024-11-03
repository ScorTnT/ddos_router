package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const IntranetConfigPath string = "/etc/config/custom_network/intranet"

// IntranetConfig OpenWrt 내부 네트워크 연결 설정을 위한 구조체
type IntranetConfig struct {
	// IP 설정
	IpAddr  string // xxx.xxx.xxx.xxx 형식
	Netmask string // xxx.xxx.xxx.xxx 형식
}

// NewIntranetConfig 기본값으로 초기화된 내부 네트워크 설정 구조체 생성
func NewIntranetConfig() *IntranetConfig {
	return &IntranetConfig{
		IpAddr:  "192.168.0.1",
		Netmask: "255.255.255.0",
	}
}

// SaveToFile OpenWrt UCI 형식으로 설정을 파일에 저장
func (i *IntranetConfig) SaveToFile() error {
	f, err := os.Create(IntranetConfigPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %v", err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	// 기본 인터페이스 설정
	_, err = writer.WriteString("config interface 'lan'\n")
	if err != nil {
		return err
	}

	// IP 주소 설정
	if i.IpAddr != "" {
		_, err = writer.WriteString(fmt.Sprintf("\toption ipaddr '%s'\n", i.IpAddr))
		if err != nil {
			return err
		}
	}

	// 서브넷 마스크 설정
	if i.Netmask != "" {
		_, err = writer.WriteString(fmt.Sprintf("\toption netmask '%s'\n", i.Netmask))
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

// LoadFromFile OpenWrt UCI 형식의 설정 파일에서 설정 읽기
func (i *IntranetConfig) LoadFromFile() (*IntranetConfig, error) {
	f, err := os.Open(IntranetConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer f.Close()

	config := NewIntranetConfig()
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
			case "ipaddr":
				config.IpAddr = value
			case "netmask":
				config.Netmask = value
			}
		}
	}

	return config, scanner.Err()
}

// SetIPAddress IP 주소 설정 (xxx.xxx.xxx.xxx 형식)
func (i *IntranetConfig) SetIPAddress(ip string) error {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return fmt.Errorf("invalid IP address format")
	}
	i.IpAddr = ip
	return nil
}

// SetNetmask 서브넷 마스크 설정 (xxx.xxx.xxx.xxx 형식)
func (i *IntranetConfig) SetNetmask(netmask string) error {
	parts := strings.Split(netmask, ".")
	if len(parts) != 4 {
		return fmt.Errorf("invalid netmask format")
	}
	i.Netmask = netmask
	return nil
}
