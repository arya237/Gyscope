package linux

import (
	"os"
	"path/filepath"
	"strconv"

	applicationProcess "Gyscope/internal/application/process"
)

type LinuxProcessDataSource struct {
	procPath string
}

func NewLinuxProcessDataSource() *LinuxProcessDataSource {
	return &LinuxProcessDataSource{
		procPath: "/proc",
	}
}

func (s *LinuxProcessDataSource) Read() ([]applicationProcess.RawProcessState, error) {
	entries, err := os.ReadDir(s.procPath)
	if err != nil {
		return nil, err
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

		processDir := filepath.Join(s.procPath, entry.Name())

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

	return processes, nil
}