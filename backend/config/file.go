package config

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	CustomNetworkPath = "/etc/config/custom_network"
	NetworkConfigPath = "/etc/config/network"
	NetworkBackupPath = "/etc/config/network.bak"
)

// ApplyNetworkConfig 네트워크 설정을 적용하는 함수
func ApplyNetworkConfig() error {
	// 1. 현재 network 설정 백업
	err := copyFile(NetworkConfigPath, NetworkBackupPath)
	if err != nil {
		return fmt.Errorf("failed to backup current network config: %v", err)
	}

	// 2. 새 설정 파일 생성
	newConfig, err := os.Create(NetworkConfigPath)
	if err != nil {
		return fmt.Errorf("failed to create new network config: %v", err)
	}
	defer newConfig.Close()

	// 3. 설정 파일들을 순서대로 합치기
	files := []string{
		filepath.Join(CustomNetworkPath, "another"),
		filepath.Join(CustomNetworkPath, "internet"),
		filepath.Join(CustomNetworkPath, "intranet"),
	}

	for _, file := range files {
		err := appendFile(newConfig, file)
		if err != nil {
			return fmt.Errorf("failed to append %s: %v", file, err)
		}
	}

	// 4. UCI 설정 적용 시도
	err = applyUCIConfig()
	if err != nil {
		// UCI 설정 적용 실패시 백업에서 복원
		fmt.Printf("Failed to apply new network configuration: %v\nRestoring previous configuration...\n", err)
		restoreErr := restoreFromBackup()
		if restoreErr != nil {
			return fmt.Errorf("critical error: failed to restore network config: %v (after failing to apply new config: %v)",
				restoreErr, err)
		}
		return fmt.Errorf("failed to apply new network config (restored to previous config): %v", err)
	}

	return nil
}

// restoreFromBackup 백업에서 복원하는 함수
func restoreFromBackup() error {
	// 백업 파일로 복원
	err := copyFile(NetworkBackupPath, NetworkConfigPath)
	if err != nil {
		return fmt.Errorf("failed to copy backup file: %v", err)
	}

	// uci commit 실행
	commitCmd := exec.Command("uci", "commit", "network")
	if output, err := commitCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to commit restored config: %v (output: %s)", err, output)
	}

	// 네트워크 재시작
	reloadCmd := exec.Command("/etc/init.d/network", "reload")
	if output, err := reloadCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to reload network with restored config: %v (output: %s)", err, output)
	}

	return nil
}

// applyUCIConfig UCI 설정을 적용하는 함수
func applyUCIConfig() error {
	// 1. UCI commit
	commitCmd := exec.Command("uci", "commit", "network")
	if output, err := commitCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to commit network config: %v (output: %s)", err, output)
	}

	// 2. 네트워크 서비스 재시작
	reloadCmd := exec.Command("/etc/init.d/network", "reload")
	if output, err := reloadCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to reload network: %v (output: %s)", err, output)
	}

	// 3. 안정화를 위한 대기
	time.Sleep(2 * time.Second)

	// 4. 네트워크 상태 확인
	checkCmd := exec.Command("ping", "-c", "1", "-W", "3", "8.8.8.8")
	if err := checkCmd.Run(); err != nil {
		return fmt.Errorf("network connectivity check failed")
	}

	// 5. 방화벽 재시작
	fwCmd := exec.Command("service", "firewall", "restart")
	if output, err := fwCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to restart firewall: %v (output: %s)", err, output)
	}

	return nil
}

// copyFile 파일 복사 함수
func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

// appendFile 파일 내용을 추가하는 함수
func appendFile(dst *os.File, src string) error {
	source, err := os.Open(src)
	if err != nil {
		if os.IsNotExist(err) {
			// 파일이 없으면 무시하고 진행
			return nil
		}
		return err
	}
	defer source.Close()

	// 추가할 때 빈 줄 하나 넣기
	_, err = dst.WriteString("\n")
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, source)
	return err
}

// ValidateCustomConfigs 커스텀 설정 파일들의 존재 여부를 확인하는 함수
func ValidateCustomConfigs() error {
	requiredFiles := []string{"internet", "intranet"}

	for _, file := range requiredFiles {
		path := filepath.Join(CustomNetworkPath, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("required config file not found: %s", path)
		}
	}

	return nil
}
