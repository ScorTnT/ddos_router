package config

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

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
