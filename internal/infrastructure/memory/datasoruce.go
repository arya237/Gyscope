package memory

import (
	"Gyscope/internal/application/memory"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type LinuxMemoryDataSource struct{}

const (
	procMemInfo = "/proc/meminfo"
)

func NewLinuxMemoryDataSoruce() *LinuxMemoryDataSource {
	return &LinuxMemoryDataSource{}
}

func (l *LinuxMemoryDataSource) Read() (memory.RawMemoryState, error) {

	content, err := os.ReadFile(procMemInfo)
	if err != nil {
		return memory.RawMemoryState{}, err
	}

	return parseMemoryInfo(string(content))
}

func parseMemoryInfo(content string) (memory.RawMemoryState, error) {
	var result memory.RawMemoryState

	for _, line := range strings.Split(content, "\n") {
		fields := strings.Fields(line)

		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "MemTotal:":
			value, err := strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return memory.RawMemoryState{}, fmt.Errorf("parse MemTotal: %w", err)
			}
			result.Total = value * 1024

		case "MemAvailable:":
			value, err := strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return memory.RawMemoryState{}, fmt.Errorf("parse MemAvailable: %w", err)
			}
			result.Available = value * 1024
		}
	}

	return result, nil
}
