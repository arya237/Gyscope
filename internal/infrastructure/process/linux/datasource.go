package linux

import (
	"os"
	"path/filepath"
	"strconv"

	applicationProcess "Gyscope/internal/application/process"
)


const (
	procPath = "/proc"
)

type LinuxProcessDataSource struct {}

func NewLinuxProcessDataSource() *LinuxProcessDataSource {
	return &LinuxProcessDataSource{}
}

func (s *LinuxProcessDataSource) Read() (applicationProcess.RawProcessSnapshot, error) {
	entries, err := os.ReadDir(procPath)
	if err != nil {
		return applicationProcess.RawProcessSnapshot{}, err
	}

	cpuStatContent, err := os.ReadFile(filepath.Join(procPath, "stat"))
	if err != nil {
		return applicationProcess.RawProcessSnapshot{}, err
	}

	totalCPUTime, err := parseTotalCPUTime(string(cpuStatContent))
	if err != nil {
		return applicationProcess.RawProcessSnapshot{}, err
	}

	processes := make([]applicationProcess.RawProcessState, 0)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		processDir := filepath.Join(procPath, entry.Name())

		statContent, err := os.ReadFile(filepath.Join(processDir, "stat"))
		if err != nil {
			continue
		}

		state, err := parseProcessStat(string(statContent))
		if err != nil {
			continue
		}

		statusContent, err := os.ReadFile(filepath.Join(processDir, "status"))
		if err != nil {
			continue
		}

		memory, err := parseProcessStatus(string(statusContent))
		if err != nil {
			continue
		}

		state.PID = pid
		state.Memory = memory

		processes = append(processes, state)
	}

	return applicationProcess.RawProcessSnapshot{
		Processes:    processes,
		TotalCPUTime: totalCPUTime,
	}, nil
}