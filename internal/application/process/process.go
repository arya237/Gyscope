package process

import Process "Gyscope/internal/domain/process"

type ProcessUseCase struct {
	dataSource ProcessDataSource
}

func NewProcessUseCase(dataSource ProcessDataSource) *ProcessUseCase {
	return &ProcessUseCase{
		dataSource: dataSource,
	}
}

func (u *ProcessUseCase) GetProcesses() ([]Process.Process, error) {
	rawProcesses, err := u.dataSource.Read()
	if err != nil {
		return nil, err
	}

	processes := make([]Process.Process, 0, len(rawProcesses))

	for _, raw := range rawProcesses {
		processes = append(processes, Process.Process{
			PID:    raw.PID,
			Name:   raw.Name,
			Memory: raw.Memory,
		})
	}

	return processes, nil
}
