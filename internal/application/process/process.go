package process

import (
	Process "Gyscope/internal/domain/process"
	"sort"
)

type ProcessUseCase struct {
	dataSource       ProcessDataSource
	previousSnapshot *RawProcessSnapshot
}

func NewProcessUseCase(dataSource ProcessDataSource) *ProcessUseCase {
	return &ProcessUseCase{
		dataSource: dataSource,
	}
}

func (u *ProcessUseCase) GetState() ([]Process.Process, error) {
	currentSnapshot, err := u.dataSource.Read()
	if err != nil {
		return nil, err
	}

	if u.previousSnapshot == nil {
		u.previousSnapshot = &currentSnapshot

		processes := u.buildProcesses(currentSnapshot, nil)
		processes = sortAndLimitProcesses(processes)
		
		u.previousSnapshot = &currentSnapshot
		return processes, nil
	}

	processes := u.buildProcesses(
		currentSnapshot,
		u.previousSnapshot,
	)

	processes = sortAndLimitProcesses(processes)
	
	u.previousSnapshot = &currentSnapshot

	return processes, nil
}

func (u *ProcessUseCase) buildProcesses(
	current RawProcessSnapshot,
	previous *RawProcessSnapshot,
) []Process.Process {
	processes := make([]Process.Process, 0, len(current.Processes))

	for _, currentProcess := range current.Processes {
		process := Process.Process{
			PID:    currentProcess.PID,
			Name:   currentProcess.Name,
			Memory: currentProcess.Memory,
		}

		if previous != nil {
			if previousProcess, ok := findProcess(
				previous.Processes,
				currentProcess.PID,
			); ok {
				process.CPUUsage = calculateCPUUsage(
					currentProcess,
					previousProcess,
					current.TotalCPUTime,
					previous.TotalCPUTime,
				)
			}
		}

		processes = append(processes, process)
	}

	return processes
}

func findProcess(
	processes []RawProcessState,
	pid int,
) (RawProcessState, bool) {
	for _, process := range processes {
		if process.PID == pid {
			return process, true
		}
	}

	return RawProcessState{}, false
}

func calculateCPUUsage(
	current RawProcessState,
	previous RawProcessState,
	currentTotalCPU uint64,
	previousTotalCPU uint64,
) float64 {
	if current.CPUTime < previous.CPUTime {
		return 0
	}

	if currentTotalCPU < previousTotalCPU {
		return 0
	}

	processDelta := current.CPUTime - previous.CPUTime
	systemDelta := currentTotalCPU - previousTotalCPU

	if systemDelta == 0 {
		return 0
	}

	return float64(processDelta) / float64(systemDelta) * 100
}

func sortAndLimitProcesses(processes []Process.Process) []Process.Process {
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].Memory > processes[j].Memory
	})

	if len(processes) > 15 {
		processes = processes[:15]
	}

	return processes
}