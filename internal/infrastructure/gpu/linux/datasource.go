package linux

import (
	applicationGPU "Gyscope/internal/application/gpu"
)

type LinuxGPUDataSource struct {
	backends []GPUBackend
}

func NewLinuxGPUDataSource() *LinuxGPUDataSource {
	return &LinuxGPUDataSource{
		backends: []GPUBackend{
			newNVIDIABackend(),
			newAMDBackend(),
			newIntelBackend(),
		},
	}
}

func (s *LinuxGPUDataSource) Read() ([]applicationGPU.RawGPUState, error) {
	devices, err := discoverGPUDevices()
	if err != nil {
		return nil, err
	}

	result := make([]applicationGPU.RawGPUState, 0)

	for _, backend := range s.backends {
		states, err := backend.Read(devices)
		if err != nil {
			return nil, err
		}

		result = append(result, states...)
	}

	return result, nil
}