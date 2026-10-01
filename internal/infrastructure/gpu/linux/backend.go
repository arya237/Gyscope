package linux

import applicationGPU "Gyscope/internal/application/gpu"

type GPUBackend interface {
	Read(devices []gpuDevice) ([]applicationGPU.RawGPUState, error)
}
