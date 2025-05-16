package router

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os/exec"
	"strings"

	"github.com/ScorTnT/ddos_router/backend/internal"
)

type PacketEntry struct {
	Protocol    string `json:"protocol"`
	SourceIP    string `json:"source_ip"`
	DestIP      string `json:"dest_ip"`
	SourcePort  string `json:"source_port"`
	DestPort    string `json:"dest_port"`
	PacketCount string `json:"packet_count"`
	ByteCount   string `json:"byte_count"`
}

func GetPackets(config *internal.Config) ([]PacketEntry, error) {
	var entries []PacketEntry

	cidr, err := getCIDR()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("conntrack", "-L")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get packets: %v", err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()

		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		fields := strings.Fields(line)

		if len(fields) < 7 {
			continue
		}

		entry := PacketEntry{}

		entry.Protocol = fields[0]

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

		if !cidr.Contains(net.ParseIP(entry.SourceIP)) {
			continue
		}

		if entry.SourceIP == "" || entry.DestIP == "" {
			entries = append(entries, entry)
		}
	}

	return entries, nil
}

func getCIDR() (*net.IPNet, error) {
	cmd := exec.Command("uci", "get", "network.lan.ipaddr")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get CIDR: %v", err)
	}

	ipAddr := strings.TrimSpace(string(output))
	_, ipNet, err := net.ParseCIDR(ipAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CIDR: %v", err)
	}

	return ipNet, nil
}
