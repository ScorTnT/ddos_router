package firewall

import (
	"fmt"
	"os/exec"
)

func InitFirewall() error {
	var cmd *exec.Cmd
	var err error

	cmd = exec.Command("nft", "add", "table", "inet", "ddos_filter")
	err = cmd.Run()

	if err != nil {
		return err
	}

	cmd = exec.Command(
		"nft",
		"add",
		"chain",
		"inet",
		"ddos_filter",
		"forward",
		"{ type filter hook forward priority 30 \\; policy accept \\; }")
	err = cmd.Run()

	if err != nil {
		return err
	}

	cmd = exec.Command("nft", "add", "set", "inet", "ddos_filter", "ban_set", "{ type ipv4_addr\\; }")
	err = cmd.Run()

	if err != nil {
		return err
	}

	cmd = exec.Command("nft", "add", "rule", "inet", "ddos_filter", "forward", "iifname", "\"eth0\"", "oifname", "\"eth1\"", "ip", "saddr", "@ban_set", "drop")

	return nil
}

func CleanupFirewall() error {
	var cmd *exec.Cmd
	var err error

	cmd = exec.Command("nft", "delete", "table", "inet", "ddos_filter")
	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

// BlockIP adds a nftables rule to block traffic to the specified IP address
func BlockIP(ip string) error {
	cmd := exec.Command("nft", "add", "element", "inet", "ddos_filter", "ban_set", fmt.Sprintf("{ %s }", ip))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to block IP %s: %v\nOutput: %s", ip, err, output)
	}
	fmt.Printf("Successfully blocked IP: %s\n", ip)
	return nil
}

// UnblockIP removes the nftables rule that blocks traffic to the specified IP address
func UnblockIP(ip string) error {
	cmd := exec.Command("nft", "delete", "element", "inet", "ddos_filter", "ban_set", fmt.Sprintf("{ %s }", ip))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to unblock IP %s: %v\nOutput: %s", ip, err, output)
	}
	fmt.Printf("Successfully unblocked IP: %s\n", ip)
	return nil
}
