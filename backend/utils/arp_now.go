package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type ARPEntry struct {
	IP  string `json:"ip"`
	MAC string `json:"mac"`
}

func GetARPNowTable() ([]ARPEntry, error) {
	cmd := exec.Command("awk", "{print $2, $3}", "/tmp/dhcp.leases")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to execute awk command: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	var entries []ARPEntry

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}

		mac := fields[0]
		ip := fields[1]

		if isValidMAC(mac) {
			entry := ARPEntry{
				IP:  ip,
				MAC: mac,
			}
			entries = append(entries, entry)
		}
	}

	return entries, nil
}

func isValidMAC(mac string) bool {
	parts := strings.Split(mac, ":")
	return len(parts) == 6
}
