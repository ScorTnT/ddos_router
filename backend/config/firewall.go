package config

import (
	"fmt"
	"os/exec"
	"strings"
)

func InitFirewall() error {
	// Define the NFT commands to setup the firewall
	commands := [][]string{
		{"add", "table", "inet", "ddos_filter"},
		{"add", "chain", "inet", "ddos_filter", "forward", "{ type filter hook forward priority 30 ; policy accept ; }"},
		{"add", "set", "inet", "ddos_filter", "ban_set", "{ type ipv4_addr; }"},
		{"add", "rule", "inet", "ddos_filter", "forward", "iifname", "br-lan", "oifname", "eth0", "ip", "saddr", "@ban_set", "drop"},
	}

	// Execute each command in sequence
	for _, args := range commands {
		cmd := exec.Command("nft", args...)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to execute nft command %v: %w\nOutput: %s", args, err, output)
		}
	}

	return nil
}

func CleanupFirewall() error {
	cmd := exec.Command("nft", "delete", "table", "inet", "ddos_filter")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to clean up firewall: %w\nOutput: %s", err, output)
	}
	return nil
}

// BlockIP adds an IP address to the ban set
func BlockIP(ip string) error {
	// [테스트용] 신기능 테스트를 위해 임시 비활성화
	fmt.Printf("[Firewall] Blocking IP: %s\n", ip)
	return nil

	err := modifyBanSet("add", ip)
	if err != nil {
		return err
	}
	return nil
}

// UnblockIP removes an IP address from the ban set
func UnblockIP(ip string) error {
	// [테스트용] 신기능 테스트를 위해 임시 비활성화
	fmt.Printf("[Firewall] Unblocking IP: %s\n", ip)
	return nil

	err := modifyBanSet("delete", ip)
	if err != nil {
		return err
	}
	return nil
}

// modifyBanSet is a helper function to add or delete IPs from the ban set
func modifyBanSet(action, ip string) error {
	cmd := exec.Command("nft", action, "element", "inet", "ddos_filter", "ban_set", fmt.Sprintf("{ %s }", ip))
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to %s IP %s: %w\nOutput: %s", action, ip, err, output)
	}
	// 직접 출력하지 않고 성공 여부만 반환
	return nil
}

// GetBlockedIPs returns a list of blocked IP addresses
func GetBlockedIPs() ([]string, error) {
	cmd := exec.Command("nft", "list", "set", "inet", "ddos_filter", "ban_set")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get blocked IPs: %w", err)
	}

	return parseBlockedIPs(string(output)), nil
}

// parseBlockedIPs extracts IP addresses from nft command output
func parseBlockedIPs(outputStr string) []string {
	var blockedIPs []string
	lines := strings.Split(outputStr, "\n")

	for _, line := range lines {
		if !strings.Contains(line, "elements =") {
			continue
		}

		start := strings.Index(line, "{")
		end := strings.LastIndex(line, "}")

		if start == -1 || end == -1 || end <= start {
			continue
		}

		elementsStr := line[start+1 : end]
		elements := strings.Split(elementsStr, ",")

		for _, ip := range elements {
			ip = strings.TrimSpace(ip)
			if ip != "" {
				blockedIPs = append(blockedIPs, ip)
			}
		}
		break
	}

	return blockedIPs
}
