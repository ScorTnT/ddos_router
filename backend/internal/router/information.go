package router

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/ScorTnT/ddos_router/backend/internal"
)

type Information struct {
	MAC                  string `json:"mac"`
	Model                string `json:"model"`
	FirmwareVersion      string `json:"firmware_version"`
	CPUUsage             string `json:"cpu_usage"`
	MemoryUsage          string `json:"memory_usage"`
	Uptime               string `json:"uptime"`
	ConnectedDeviceCount int    `json:"connected_device_count"`
	UploadSpeed          string `json:"upload_speed"`
	DownloadSpeed        string `json:"download_speed"`
}

func GetInformation(config *internal.Config) (*Information, error) {
	interfaceName := config.Interface.WANInterface
	if interfaceName == "" {
		return nil, fmt.Errorf("interface name parameter cannot be empty")
	}

	mac, err := getMACAddress(interfaceName)
	if err != nil {
		return nil, err
	}

	model, err := getModel()
	if err != nil {
		return nil, err
	}

	firmwareVersion, err := getFirmwareVersion()
	if err != nil {
		return nil, err
	}

	cpuUsage, err := getCPUUsage()
	if err != nil {
		return nil, err
	}

	memoryUsage, err := getMemoryUsage()
	if err != nil {
		return nil, err
	}

	uptime, err := getUptime()
	if err != nil {
		return nil, err
	}

	connectedDeviceCount, err := getConnectedDeviceCount(interfaceName)
	if err != nil {
		return nil, err
	}

	uploadSpeed, downloadSpeed, err := getEthernetSpeed(interfaceName)
	if err != nil {
		return nil, err
	}

	return &Information{
		MAC:                  mac,
		Model:                model,
		FirmwareVersion:      firmwareVersion,
		CPUUsage:             cpuUsage,
		MemoryUsage:          memoryUsage,
		Uptime:               uptime,
		ConnectedDeviceCount: connectedDeviceCount,
		UploadSpeed:          uploadSpeed,
		DownloadSpeed:        downloadSpeed,
	}, nil
}

func getMACAddress(interfaceName string) (string, error) {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return "", err
	}

	return iface.HardwareAddr.String(), nil
}

func getModel() (string, error) {
	data, err := os.ReadFile("/proc/device-tree/model")
	if err != nil {
		return "Unknown", err
	}

	return strings.TrimSpace(string(data)), nil
}

func getFirmwareVersion() (string, error) {
	data, err := os.ReadFile("/etc/openwrt_release")
	if err != nil {
		return "Unknown", err
	}

	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "DISTRIB_RELEASE=") {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				return strings.Trim(parts[1], "\""), nil
			}
		}
	}

	return "Unknown", nil
}

func getCPUUsage() (string, error) {
	readAndParseTimes := func() (uint64, uint64, error) {
		file, errFile := os.Open("/proc/stat")
		if errFile != nil {
			return 0, 0, fmt.Errorf("could not open /proc/stat: %w", errFile)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		if !scanner.Scan() { // 첫 번째 줄 ("cpu ...")만 필요합니다.
			if errScan := scanner.Err(); errScan != nil {
				return 0, 0, fmt.Errorf("could not scan /proc/stat: %w", errScan)
			}
			return 0, 0, fmt.Errorf("/proc/stat was empty or unreadable")
		}

		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) == 0 || fields[0] != "cpu" {
			return 0, 0, fmt.Errorf("unexpected format in /proc/stat: line does not start with 'cpu'")
		}

		if len(fields) < 5 {
			return 0, 0, fmt.Errorf("not enough fields in /proc/stat cpu line. Expected at least 5 (cpu + 4 values), got %d", len(fields))
		}

		var times [8]uint64
		var parseFieldErr error

		for i := 0; i < 8; i++ {
			if (i + 1) < len(fields) {
				val, errParse := strconv.ParseUint(fields[i+1], 10, 64)
				if errParse != nil {
					if i < 4 {
						parseFieldErr = fmt.Errorf("failed to parse critical cpu time field '%s' (index %d): %w", fields[i+1], i+1, errParse)
						break
					}
					times[i] = 0
				} else {
					times[i] = val
				}
			} else {
				if i < 4 {
					parseFieldErr = fmt.Errorf("critical cpu time field (index %d) missing from /proc/stat line", i+1)
					break
				}
				times[i] = 0
			}
		}

		if parseFieldErr != nil {
			return 0, 0, fmt.Errorf("error parsing CPU times from /proc/stat: %w", parseFieldErr)
		}

		var currentTotalTime uint64
		for _, t := range times {
			currentTotalTime += t
		}
		currentIdleTime := times[3]

		return currentTotalTime, currentIdleTime, nil
	}

	prevTotal, prevIdle, err := readAndParseTimes()
	if err != nil {
		return "", fmt.Errorf("failed to read initial CPU times: %w. This function is intended for Linux systems", err)
	}

	time.Sleep(1 * time.Second)

	currTotal, currIdle, err := readAndParseTimes()
	if err != nil {
		return "", fmt.Errorf("failed to read current CPU times: %w", err)
	}

	deltaTotal := float64(currTotal) - float64(prevTotal)
	deltaIdle := float64(currIdle) - float64(prevIdle)

	var cpuUsagePercentage float64

	if deltaTotal > 0 {
		usage := 1.0 - (deltaIdle / deltaTotal)
		cpuUsagePercentage = usage * 100.0
	} else if deltaTotal == 0 && deltaIdle == 0 {
		cpuUsagePercentage = 0.0
	} else {
		cpuUsagePercentage = 0.0
	}

	if cpuUsagePercentage < 0 {
		cpuUsagePercentage = 0.0
	}
	if cpuUsagePercentage > 100 {
		cpuUsagePercentage = 100.0
	}

	roundedUsage := int(math.Round(cpuUsagePercentage))

	return fmt.Sprintf("%d%%", roundedUsage), nil
}

func getMemoryUsage() (string, error) {
	cmd := exec.Command("sh", "-c", "free | grep Mem | awk '{print $3/$2 * 100.0}'")
	var out bytes.Buffer

	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}

	memValue := strings.TrimSpace(out.String())
	if memFloat, err := strconv.ParseFloat(memValue, 64); err == nil {
		return fmt.Sprintf("%.1f%%", memFloat), nil
	}

	return "", fmt.Errorf("failed to parse memory usage")
}

func getUptime() (string, error) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return "Unknown", err
	}

	fields := strings.Fields(string(data))
	if len(fields) < 1 {
		return "Unknown", fmt.Errorf("invalid uptime data")
	}

	uptime, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return "Unknown", err
	}

	uptimeInt := int(uptime)

	duration := time.Duration(uptimeInt) * time.Second
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds), nil
}

func getConnectedDeviceCount(interfaceName string) (int, error) {
	neighbors, err := GetOnlineNeighbors()
	if err != nil {
		return 0, err
	}

	count := 0

	for _, neighbor := range neighbors {
		if neighbor.Device == interfaceName {
			count++
		}
	}

	return count, nil
}

func getEthernetSpeed(interfaceName string) (string, string, error) {
	if strings.TrimSpace(interfaceName) == "" {
		return "Unknown", "Unknown", fmt.Errorf("interface name parameter cannot be empty")
	}

	readStatsForInterface := func(ifName string) (rxBytes uint64, txBytes uint64, e error) {
		data, errRead := os.ReadFile("/proc/net/dev")
		if errRead != nil {
			return 0, 0, fmt.Errorf("failed to read /proc/net/dev: %w", errRead)
		}

		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		var foundInterface bool
		for scanner.Scan() {
			line := scanner.Text()
			trimmedLine := strings.TrimSpace(line)
			if strings.HasPrefix(trimmedLine, ifName+":") {
				foundInterface = true
				fields := strings.Fields(trimmedLine)

				if len(fields) >= 10 {
					rx, errParseRx := strconv.ParseUint(fields[1], 10, 64)
					if errParseRx != nil {
						return 0, 0, fmt.Errorf("failed to parse RX bytes for %s from '%s': %w", ifName, fields[1], errParseRx)
					}

					tx, errParseTx := strconv.ParseUint(fields[9], 10, 64)
					if errParseTx != nil {
						return 0, 0, fmt.Errorf("failed to parse TX bytes for %s from '%s': %w", ifName, fields[9], errParseTx)
					}
					return rx, tx, nil
				}

				return 0, 0, fmt.Errorf("data format error for interface %s in /proc/net/dev (expected at least 10 fields, got %d)", ifName, len(fields))
			}
		}
		if errScan := scanner.Err(); errScan != nil {
			return 0, 0, fmt.Errorf("error scanning /proc/net/dev: %w", errScan)
		}
		if !foundInterface {
			return 0, 0, fmt.Errorf("interface %s not found in /proc/net/dev", ifName)
		}

		return 0, 0, fmt.Errorf("unknown error after scanning for interface %s", ifName)
	}

	rx1, tx1, err := readStatsForInterface(interfaceName)
	if err != nil {
		return "Unknown", "Unknown", fmt.Errorf("failed to get initial network stats for %s: %w", interfaceName, err)
	}

	time.Sleep(time.Second)

	rx2, tx2, err := readStatsForInterface(interfaceName)
	if err != nil {
		return "Unknown", "Unknown", fmt.Errorf("failed to get final network stats for %s: %w", interfaceName, err)
	}

	var downloadSpeedMbps, uploadSpeedMbps float64

	if rx2 >= rx1 {
		downloadSpeedMbps = float64(rx2-rx1) * 8.0 / (1000 * 1000.0)
	} else {
		downloadSpeedMbps = 0.0
	}

	if tx2 >= tx1 {
		uploadSpeedMbps = float64(tx2-tx1) * 8.0 / (1000 * 1000.0)
	} else {
		uploadSpeedMbps = 0.0
	}

	downloadSpeedStr := fmt.Sprintf("%.1f", downloadSpeedMbps)
	uploadSpeedStr := fmt.Sprintf("%.1f", uploadSpeedMbps)

	return downloadSpeedStr, uploadSpeedStr, nil
}
