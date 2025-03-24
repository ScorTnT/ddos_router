package analyzer

import (
	"fmt"
	"net"
	"os"
	"strings"

	"openwrt-ips/internal/config"
)

type ARPEntry struct {
	IP  string `json:"ip"`
	MAC string `json:"mac"`
}

func GetARPTable(cfg *config.Config) ([]ARPEntry, error) {
	data, err := os.ReadFile("/proc/net/arp")
	if err != nil {
		return nil, fmt.Errorf("failed to read /proc/net/arp: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var entries []ARPEntry

	for _, line := range lines[1:] { // skip header
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		device := fields[5]
		if device != cfg.Ethernet.LANInterface {
			continue
		}

		mac := fields[3]
		if !isValidMAC(mac) {
			continue
		}

		entries = append(entries, ARPEntry{
			IP:  fields[0],
			MAC: mac,
		})
	}

	return entries, nil
}

func isValidMAC(mac string) bool {
	_, err := net.ParseMAC(mac)
	return err == nil
}
