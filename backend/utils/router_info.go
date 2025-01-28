package utils

import (
	"bufio"
	"bytes"
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

func GetRouterInfo() ([]RouterInfo, error) {
    var routerInfo []RouterInfo

    // 네트워크 인터페이스 이름 동적 탐색
    interfaceName := ""
    interfaces, err := net.Interfaces()
    if err != nil {
        interfaceName = "Unknown"
    } else {
        for _, iface := range interfaces {
            if strings.Contains(iface.Name, "eth") { // "eth"가 포함된 이름 검색
                interfaceName = iface.Name
                break
            }
        }
        if interfaceName == "" {
            interfaceName = "Unknown"
        }
    }

    // MAC 주소
    mac, err := getMACAddress(interfaceName)
    if err != nil || interfaceName == "Unknown" {
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
	routerInfo = append(routerInfo, RouterInfo{Name: "연결된 기기 수", Value: fmt.Sprintf("%d", connected)})

	download, upload, err := getDownloadUploadSpeed()
	if err != nil {
		download, upload = "Unknown", "Unknown"
	}
	routerInfo = append(routerInfo, RouterInfo{Name: "현재 다운로드 속도", Value: download})
	routerInfo = append(routerInfo, RouterInfo{Name: "현재 업로드 속도", Value: upload})

	return routerInfo, nil
}

func getMACAddress(interfaceName string) (string, error) {
    // 인터페이스 이름이 제공되지 않으면 동적으로 탐색
    if interfaceName == "" {
        interfaces, err := net.Interfaces()
        if err != nil {
            return "", err
        }

        for _, iface := range interfaces {
            if strings.Contains(iface.Name, "eth") { // "eth"가 포함된 이름 검색
                return iface.HardwareAddr.String(), nil
            }
        }
        return "Unknown", nil // 적합한 인터페이스를 찾지 못한 경우
    }

    // 기존 로직
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
	// CPU 사용률을 계산하기 위해 /proc/stat 파일을 읽습니다
	data1, err := os.ReadFile("/proc/stat")
	if err != nil {
		return "", err
	}

	// CPU 통계의 첫 번째 라인을 파싱합니다
	fields1 := strings.Fields(strings.Split(string(data1), "\n")[0])
	if len(fields1) < 5 {
		return "", fmt.Errorf("invalid CPU stats")
	}

	user1, _ := strconv.ParseUint(fields1[1], 10, 64)
	nice1, _ := strconv.ParseUint(fields1[2], 10, 64)
	system1, _ := strconv.ParseUint(fields1[3], 10, 64)
	idle1, _ := strconv.ParseUint(fields1[4], 10, 64)

	// 1초 대기
	time.Sleep(time.Second)

	// 두 번째 측정
	data2, err := os.ReadFile("/proc/stat")
	if err != nil {
		return "", err
	}

	fields2 := strings.Fields(strings.Split(string(data2), "\n")[0])
	if len(fields2) < 5 {
		return "", fmt.Errorf("invalid CPU stats")
	}

	user2, _ := strconv.ParseUint(fields2[1], 10, 64)
	nice2, _ := strconv.ParseUint(fields2[2], 10, 64)
	system2, _ := strconv.ParseUint(fields2[3], 10, 64)
	idle2, _ := strconv.ParseUint(fields2[4], 10, 64)

	// CPU 사용률 계산
	active1 := user1 + nice1 + system1
	active2 := user2 + nice2 + system2
	total1 := active1 + idle1
	total2 := active2 + idle2

	deltaCPU := active2 - active1
	deltaTotal := total2 - total1

	if deltaTotal == 0 {
		return "0.0%", nil
	}

	cpuUsage := float64(deltaCPU) / float64(deltaTotal) * 100
	return fmt.Sprintf("%.1f%%", cpuUsage), nil
}

func getMemoryUsage() (string, error) {
	cmd := exec.Command("sh", "-c", "free | grep Mem | awk '{print $3/$2 * 100.0}'")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	// 숫자만 추출하고 소수점 첫째자리까지만 표시
	memValue := strings.TrimSpace(out.String())
	if memFloat, err := strconv.ParseFloat(memValue, 64); err == nil {
		return fmt.Sprintf("%.1f%%", memFloat), nil
	}
	return "", fmt.Errorf("failed to parse memory usage")
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

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

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
	return count, nil
}

func getDownloadUploadSpeed() (string, string, error) {
    interfaceName := ""
    interfaces, err := net.Interfaces()
    if err != nil {
        return "", "", err
    }

    for _, iface := range interfaces {
        if strings.Contains(iface.Name, "eth") {
            interfaceName = iface.Name
            break
        }
    }

    if interfaceName == "" {
        return "Unknown", "Unknown", fmt.Errorf("no matching network interface found")
    }

    rx1, tx1, err := getNetworkStats(interfaceName)
    if err != nil {
        return "", "", err
    }

    time.Sleep(time.Second)

    rx2, tx2, err := getNetworkStats(interfaceName)
    if err != nil {
        return "", "", err
    }

    downloadSpeed := float64(rx2-rx1) * 8 / 1000000 // bytes to Mbps
    uploadSpeed := float64(tx2-tx1) * 8 / 1000000   // bytes to Mbps

    return fmt.Sprintf("%.1f", downloadSpeed), fmt.Sprintf("%.1f", uploadSpeed), nil
}


func getNetworkStats(interfaceName string) (uint64, uint64, error) {
    // 인터페이스 이름이 제공되지 않으면 동적으로 탐색
    if interfaceName == "" {
        interfaces, err := net.Interfaces()
        if err != nil {
            return 0, 0, err
        }

        for _, iface := range interfaces {
            if strings.Contains(iface.Name, "eth") {
                interfaceName = iface.Name
                break
            }
        }

        if interfaceName == "" {
            return 0, 0, fmt.Errorf("no matching network interface found")
        }
    }

    data, err := os.ReadFile("/proc/net/dev")
    if err != nil {
        return 0, 0, err
    }

    scanner := bufio.NewScanner(strings.NewReader(string(data)))
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, interfaceName) {
            fields := strings.Fields(line)
            if len(fields) < 10 {
                return 0, 0, fmt.Errorf("invalid format in /proc/net/dev")
            }

            rx, err := strconv.ParseUint(fields[1], 10, 64)
            if err != nil {
                return 0, 0, err
            }

            tx, err := strconv.ParseUint(fields[9], 10, 64)
            if err != nil {
                return 0, 0, err
            }

            return rx, tx, nil
        }
    }

    return 0, 0, fmt.Errorf("interface %s not found", interfaceName)
}
