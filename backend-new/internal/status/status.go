package status

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"openwrt-ips/internal/config"
)

type RouterStatus struct {
	MACAddress       string `json:"mac_address"`
	Model            string `json:"model"`
	Firmware         string `json:"firmware"`
	CPUUsage         string `json:"cpu_usage"`
	MemoryUsage      string `json:"memory_usage"`
	Uptime           string `json:"uptime"`
	ConnectedDevices int    `json:"connected_devices"`
	DownloadSpeed    string `json:"download_speed"`
	UploadSpeed      string `json:"upload_speed"`
}

func GetRouterStatus(cfg *config.Config) (*RouterStatus, error) {
	iface := cfg.Ethernet.WANInterface

	if iface == "" {
		iface = findEthInterface()
	}

	status := &RouterStatus{
		MACAddress:       getOrUnknown(getMACAddress(iface)),
		Model:            getOrUnknown(getFileTrim("/proc/device-tree/model")),
		Firmware:         getOrUnknown(parseKey("/etc/openwrt_release", "DISTRIB_RELEASE")),
		CPUUsage:         getOrUnknown(getCPUUsage()),
		MemoryUsage:      getOrUnknown(getMemoryUsage()),
		Uptime:           getOrUnknown(getUptime()),
		ConnectedDevices: getOrZero(getConnectedDevices()),
	}

	download, upload, _ := getDownloadUploadSpeed(iface)
	status.DownloadSpeed = defaultIfEmpty(download, "Unknown")
	status.UploadSpeed = defaultIfEmpty(upload, "Unknown")

	return status, nil
}

func findEthInterface() string {
	for _, iface := range must(net.Interfaces()) {
		if strings.Contains(iface.Name, "eth") {
			return iface.Name
		}
	}
	return ""
}

func getOrUnknown(val string, err error) string {
	if err != nil || val == "" {
		return "Unknown"
	}
	return val
}

func getOrZero(count int, err error) int {
	if err != nil {
		return 0
	}
	return count
}

func defaultIfEmpty(val, def string) string {
	if val == "" {
		return def
	}
	return val
}

func getMACAddress(iface string) (string, error) {
	if iface == "" {
		return "", nil
	}
	i, err := net.InterfaceByName(iface)
	if err != nil {
		return "", err
	}
	return i.HardwareAddr.String(), nil
}

func getFileTrim(path string) (string, error) {
	data, err := os.ReadFile(path)
	return strings.TrimSpace(string(data)), err
}

func parseKey(path, key string) (string, error) {
	for _, line := range readLines(path) {
		if strings.HasPrefix(line, key+"=") {
			return strings.Trim(strings.SplitN(line, "=", 2)[1], `"'`), nil
		}
	}
	return "", nil
}

func getCPUUsage() (string, error) {
	parse := func() (active, total uint64) {
		fields := strings.Fields(readLines("/proc/stat")[0])
		user, _ := strconv.ParseUint(fields[1], 10, 64)
		nice, _ := strconv.ParseUint(fields[2], 10, 64)
		sys, _ := strconv.ParseUint(fields[3], 10, 64)
		idle, _ := strconv.ParseUint(fields[4], 10, 64)
		active = user + nice + sys
		total = active + idle
		return
	}
	a1, t1 := parse()
	time.Sleep(time.Second)
	a2, t2 := parse()

	if delta := t2 - t1; delta > 0 {
		return fmt.Sprintf("%.1f%%", float64(a2-a1)/float64(delta)*100), nil
	}
	return "0.0%", nil
}

func getMemoryUsage() (string, error) {
	var total, avail float64
	for _, line := range readLines("/proc/meminfo") {
		fields := strings.Fields(line)
		switch fields[0] {
		case "MemTotal:":
			total, _ = strconv.ParseFloat(fields[1], 64)
		case "MemAvailable:":
			avail, _ = strconv.ParseFloat(fields[1], 64)
		}
	}
	if total == 0 {
		return "", fmt.Errorf("invalid meminfo")
	}
	return fmt.Sprintf("%.1f%%", (total-avail)/total*100), nil
}

func getUptime() (string, error) {
	fields := strings.Fields(readLines("/proc/uptime")[0])
	secs, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return "", err
	}
	d := time.Duration(secs) * time.Second
	return fmt.Sprintf("%d일 %d시간 %d분", int(d.Hours()/24), int(d.Hours())%24, int(d.Minutes())%60), nil
}

func getConnectedDevices() (int, error) {
	count := 0
	for _, line := range readLines("/proc/net/arp")[1:] {
		if strings.Count(line, ":") >= 5 {
			count++
		}
	}
	return count, nil
}

func getDownloadUploadSpeed(iface string) (string, string, error) {
	stats := func() (rx, tx uint64) {
		for _, line := range readLines("/proc/net/dev")[2:] {
			if strings.Contains(line, iface) {
				fields := strings.Fields(line)
				rx, _ = strconv.ParseUint(fields[1], 10, 64)
				tx, _ = strconv.ParseUint(fields[9], 10, 64)
				return
			}
		}
		return
	}
	r1, t1 := stats()
	time.Sleep(time.Second)
	r2, t2 := stats()

	toMbps := func(bytes uint64) string { return fmt.Sprintf("%.1f", float64(bytes*8)/1e6) }
	return toMbps(r2 - r1), toMbps(t2 - t1), nil
}

func readLines(path string) []string {
	f, _ := os.Open(path)
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}
