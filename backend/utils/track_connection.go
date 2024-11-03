package utils

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

// ConntrackEntry represents a single conntrack entry
type ConntrackEntry struct {
	Protocol    string `json:"protocol"`
	SourceIP    string `json:"source_ip"`
	DestIP      string `json:"dest_ip"`
	SourcePort  string `json:"source_port"`
	DestPort    string `json:"dest_port"`
	PacketCount string `json:"packet_count"`
	ByteCount   string `json:"byte_count"`
}

// ParseConntrackOutput parses the output of the conntrack command
// and returns a slice of ConntrackEntry
func ParseConntrackOutput(output string) []ConntrackEntry {
	var entries []ConntrackEntry
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		fields := strings.Fields(line)

		// Basic validation of minimum fields
		if len(fields) < 7 {
			continue
		}

		entry := ConntrackEntry{}

		// Safely parse protocol
		entry.Protocol = fields[0]

		// Safely parse IP addresses and ports
		for _, field := range fields {
			parts := strings.Split(field, "=")
			if len(parts) != 2 {
				continue
			}

			key := parts[0]
			value := parts[1]

			switch key {
			case "src":
				entry.SourceIP = value
			case "dst":
				entry.DestIP = value
			case "sport":
				entry.SourcePort = value
			case "dport":
				entry.DestPort = value
			}
		}

		// Safely parse packet and byte counts
		// Usually these are the last fields, but we'll search for them specifically
		for _, field := range fields {
			if strings.HasPrefix(field, "packets=") {
				parts := strings.Split(field, "=")
				if len(parts) == 2 {
					entry.PacketCount = parts[1]
				}
			}
			if strings.HasPrefix(field, "bytes=") {
				parts := strings.Split(field, "=")
				if len(parts) == 2 {
					entry.ByteCount = parts[1]
				}
			}
		}

		// Only append if we have the minimum required fields
		if entry.SourceIP != "" && entry.DestIP != "" {
			entries = append(entries, entry)
		}
	}
	return entries
}

// GetConntrackEntries executes the conntrack command and returns
// a slice of ConntrackEntry and an error if any
func GetConntrackEntries() ([]ConntrackEntry, error) {
	cmd := exec.Command("conntrack", "-L")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error executing conntrack: %w", err)
	}

	return ParseConntrackOutput(string(output)), nil
}
