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
		fields := strings.Fields(line)
		if len(fields) < 7 {
			continue
		}
		entry := ConntrackEntry{
			Protocol:    fields[0],
			SourceIP:    strings.Split(fields[3], "=")[1],
			DestIP:      strings.Split(fields[4], "=")[1],
			SourcePort:  strings.Split(fields[5], "=")[1],
			DestPort:    strings.Split(fields[6], "=")[1],
			PacketCount: fields[len(fields)-3],
			ByteCount:   fields[len(fields)-2],
		}
		entries = append(entries, entry)
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
