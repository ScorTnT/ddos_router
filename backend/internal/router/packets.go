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
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		entry := parseConntrackLine(line)
		if entry == nil {
			continue
		}

		// LAN 내부에서 시작된 아웃바운드 트래픽만 필터링
		if entry.SourceIP != "" && cidr.Contains(net.ParseIP(entry.SourceIP)) {
			entries = append(entries, *entry)
		}
	}

	return entries, nil
}

func parseConntrackLine(line string) *PacketEntry {
	fields := strings.Fields(line)
	if len(fields) < 7 {
		return nil
	}

	entry := &PacketEntry{Protocol: fields[0]}

	// 정규식으로 더 정교한 파싱도 가능하지만, 현재 방식이 더 직관적
	fieldMap := make(map[string]string)

	// 모든 key=value 파싱
	for _, field := range fields {
		if key, value, ok := parseKeyValue(field); ok {
			// 첫 번째 값만 저장 (아웃바운드 방향)
			if _, exists := fieldMap[key]; !exists {
				fieldMap[key] = value
			}
		}
	}

	// 필드 할당
	entry.SourceIP = fieldMap["src"]
	entry.DestIP = fieldMap["dst"]
	entry.SourcePort = fieldMap["sport"]
	entry.DestPort = fieldMap["dport"]
	entry.PacketCount = fieldMap["packets"]
	entry.ByteCount = fieldMap["bytes"]

	// 필수 필드 검증
	if entry.SourceIP == "" || entry.DestIP == "" {
		return nil
	}

	return entry
}

func parseKeyValue(field string) (key, value string, ok bool) {
	parts := strings.SplitN(field, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func getCIDR() (*net.IPNet, error) {
	ipCmd := exec.Command("uci", "get", "network.lan.ipaddr")
	ipOutput, err := ipCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get IP: %v", err)
	}

	maskCmd := exec.Command("uci", "get", "network.lan.netmask")
	maskOutput, err := maskCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get netmask: %v", err)
	}

	ipAddr := strings.TrimSpace(string(ipOutput))
	netmask := strings.TrimSpace(string(maskOutput))

	cidr := ipAddr + "/" + netmaskToCIDR(netmask)

	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CIDR: %v", err)
	}

	return ipNet, nil
}

func netmaskToCIDR(netmask string) string {
	// 255.255.255.0 -> 24 변환 로직
	switch netmask {
	case "255.255.255.0":
		return "24"
	case "255.255.0.0":
		return "16"
	case "255.0.0.0":
		return "8"
	default:
		return "24" // 기본값
	}
}
