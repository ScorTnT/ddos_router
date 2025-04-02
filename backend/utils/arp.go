package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// ARPEntry 구조체 정의
type ARPEntry struct {
	IP  string `json:"ip"`
	MAC string `json:"mac"`
}

// GetARPTable 함수 정의
func GetARPTable() ([]ARPEntry, error) {
	cmd := exec.Command("ip", "neigh", "show","nud","reachable")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to execute ip command: %v", err)
	}

	// 결과를 파싱하여 ARPEntry 배열 생성
	lines := strings.Split(out.String(), "\n")
	var entries []ARPEntry

	for _, line := range lines {
		fields := strings.Fields(line)
		// ip neigh 출력 형식: <IP> dev <interface> lladdr <MAC> REACHABLE
		// 인터페이스명이 "br-lan"인 경우만 추가
		if len(fields) >= 5 &&
			fields[1] == "dev" &&
			fields[2] == "br-lan" &&
			fields[3] == "lladdr" &&
			isValidMAC(fields[4]) {
			entry := ARPEntry{
				IP:  fields[0],
				MAC: fields[4],
			}
			entries = append(entries, entry)
		}
	}

	return entries, nil
}

// 유효한 MAC 주소인지 확인하는 함수
func isValidMAC(mac string) bool {
	parts := strings.Split(mac, ":")
	return len(parts) == 6
}
