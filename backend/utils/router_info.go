package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type RouterInfo struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func GetRouterInfo() ([]byte, error) {
	var routerInfo []RouterInfo

	mac, err := getMACAddress("eth0")
	if err != nil {
		mac = "Unknown"
	}
	routerInfo = append(routerInfo, RouterInfo{Name: "MAC 주소", Value: mac})

	model, err := getModelName()
	if err != nil {
		model = "Unknown"
	}
	routerInfo = append(routerInfo, RouterInfo{Name: "모델명", Value: model})

	firmware, err := getFirmwareVersion()
	if err != nil {
		firmware = "Unknown"
	}
	routerInfo = append(routerInfo, RouterInfo{Name: "펌웨어 버전", Value: firmware})

	cpuUsage, err := getCPUUsage()
	if err != nil {
		cpuUsage = "Unknown"
	}
	routerInfo = append(routerInfo, RouterInfo{Name: "CPU 사용률", Value: cpuUsage})

	memUsage, err := getMemoryUsage()
	if err != nil {
		memUsage = "Unknown"
	}
	routerInfo = append(routerInfo, RouterInfo{Name: "메모리 사용률", Value: memUsage})

	uptime, err := getUptime()
	if err != nil {
		uptime = "Unknown"
	}
	routerInfo = append(routerInfo, RouterInfo{Name: "가동 시간", Value: uptime})

	connected, err := getConnectedDevices()
	if err != nil {
		connected = 0
	}
	routerInfo = append(routerInfo, RouterInfo{Name: "연결된 기기 수", Value: fmt.Sprintf("%d대", connected)})

	download, upload, err := getDownloadUploadSpeed()
	if err != nil {
		download, upload = "Unknown", "Unknown"
	}
	routerInfo = append(routerInfo, RouterInfo{Name: "현재 다운로드 속도", Value: download})
	routerInfo = append(routerInfo, RouterInfo{Name: "현재 업로드 속도", Value: upload})

	jsonData, err := json.Marshal(routerInfo)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func getMACAddress(interfaceName string) (string, error) {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return "", err
	}
	return iface.HardwareAddr.String(), nil
}

func getModelName() (string, error) {
	data, err := os.ReadFile("/proc/device-tree/model")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func getFirmwareVersion() (string, error) {
	data, err := os.ReadFile("/etc/openwrt_release")
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "DISTRIB_RELEASE") {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				return strings.Trim(parts[1], "\""), nil
			}
		}
	}
	return "Unknown", nil
}

func getCPUUsage() (string, error) {
	cmd := exec.Command("sh", "-c", "top -bn1 | grep 'Cpu(s)'")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func getMemoryUsage() (string, error) {
	cmd := exec.Command("free", "-m")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func getUptime() (string, error) {
	uptimeSeconds, err := getUptimeSeconds()
	if err != nil {
		return "", err
	}
	duration := time.Duration(uptimeSeconds) * time.Second
	days := int(duration.Hours()) / 24
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60
	return fmt.Sprintf("%d일 %d시간 %d분", days, hours, minutes), nil
}

func getUptimeSeconds() (int, error) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}
	fields := strings.Fields(string(data))
	if len(fields) < 1 {
		return 0, fmt.Errorf("invalid uptime data")
	}
	uptimeFloat, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, err
	}
	return int(uptimeFloat), nil
}

func getConnectedDevices() (int, error) {
	file, err := os.Open("/proc/net/arp")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, ":") { // MAC 주소에 대한 간단한 체크
			count++
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	// 헤더를 제외하기 위해 1을 빼줍니다.
	return count - 1, nil
}

func getDownloadUploadSpeed() (string, string, error) {
	// 실제 속도 측정을 위한 구현이 필요합니다.
	// 예시에서는 고정된 값을 반환합니다.
	return "50Mbps", "20Mbps", nil
}
