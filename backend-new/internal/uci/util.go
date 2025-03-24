package config

import (
	"fmt"
	"os/exec"
	"strings"
)

func getUCIValue(value string) (string, error) {
	cmd := exec.Command("uci", "get", value)
	output, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("failed to get uci value for %s: %v", value, err)
	}

	return strings.TrimSpace(string(output)), nil
}
