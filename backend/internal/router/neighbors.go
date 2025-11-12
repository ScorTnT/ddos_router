package router

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

type Neighbor struct {
	IP     string `json:"ip"`
	Device string `json:"device"`
	MAC    string `json:"mac"`
	State  string `json:"state"`
}

func GetOnlineNeighbors() ([]Neighbor, error) {
	cmd := exec.Command("ip", "neigh", "show", "|", "grep", "br-lan")
	outputBytes, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute 'ip neigh show': %w", err)
	}

	var neighbors []Neighbor
	scanner := bufio.NewScanner(strings.NewReader(string(outputBytes)))

	knownStates := map[string]bool{
		"PERMANENT": true, "NOARP": true, "REACHABLE": true, "STALE": true,
		"NONE": true, "INCOMPLETE": true, "DELAY": true, "PROBE": true, "FAILED": true,
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		entry := Neighbor{
			IP: fields[0],
		}

		var macAddressFound bool
		var currentState string

		for i := 1; i < len(fields); i++ {
			fieldValue := fields[i]

			if _, isKnownState := knownStates[strings.ToUpper(fieldValue)]; isKnownState {
				currentState = fieldValue
			}

			switch fieldValue {
			case "dev":
				if i+1 < len(fields) {
					entry.Device = fields[i+1]
					i++
				}
			case "lladdr":
				if i+1 < len(fields) {
					mac := fields[i+1]
					if mac != "(incomplete)" {
						entry.MAC = mac
						macAddressFound = true
					}
					i++
				}
			}
		}

		entry.State = currentState

		if macAddressFound && entry.Device != "" && entry.State != "" {
			switch strings.ToUpper(entry.State) {
			case "REACHABLE", "STALE", "DELAY", "PROBE", "PERMANENT":
				neighbors = append(neighbors, entry)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning 'ip neigh' output: %w", err)
	}

	return neighbors, nil
}
