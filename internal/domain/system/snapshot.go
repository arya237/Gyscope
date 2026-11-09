package system

import (
	"time"

	"Gyscope/internal/domain/cpu"
	"Gyscope/internal/domain/disk"
	"Gyscope/internal/domain/gpu"
	"Gyscope/internal/domain/memory"
	"Gyscope/internal/domain/network"
	"Gyscope/internal/domain/process"
)

type Snapshot struct {
	CollectedAt time.Time
	CPU         cpu.CPU
	Memory      memory.Memory
	Disk        disk.Disk
	GPUs        []gpu.GPU
	Processes   []process.Process
	Network     network.Network
}


