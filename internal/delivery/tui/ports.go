package tui

import (
	"Gyscope/internal/domain/cpu"
	"Gyscope/internal/domain/disk"
	"Gyscope/internal/domain/gpu"
	"Gyscope/internal/domain/memory"
	"Gyscope/internal/domain/process"
)

type CPUReader interface {
	GetState() (cpu.CPU, error)
}

type MemoryReader interface {
	GetState() (memory.Memory, error)
}

type DiskReader interface {
	GetState() (disk.Disk, error)
}

type ProcessReader interface {
	GetState() ([]process.Process, error)
}

type GPUReader interface {
	GetState() ([]gpu.GPU, error)
}
