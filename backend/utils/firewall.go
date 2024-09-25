package utils

import (
	"fmt"
	"os/exec"
)

// BlockIP adds a nftables rule to block traffic to the specified IP address
func BlockIP(ip string) error {
	cmd := exec.Command("nft", "add", "rule", "inet", "filter", "input", "ip", "saddr", ip, "drop")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to block IP %s: %v\nOutput: %s", ip, err, output)
	}
	fmt.Printf("Successfully blocked IP: %s\n", ip)
	return nil
}

// UnblockIP removes the nftables rule that blocks traffic to the specified IP address
func UnblockIP(ip string) error {
	cmd := exec.Command("nft", "delete", "rule", "inet", "filter", "input", "ip", "saddr", ip, "drop")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to unblock IP %s: %v\nOutput: %s", ip, err, output)
	}
	fmt.Printf("Successfully unblocked IP: %s\n", ip)
	return nil
}
