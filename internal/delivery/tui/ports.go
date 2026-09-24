package tui

import (
	"Gyscope/internal/domain/cpu"
	"Gyscope/internal/domain/memory"
)

type CPUReader interface {
	GetState() (cpu.CPU, error)
}

type MemoryReader interface {
	GetState() (memory.Memory, error)
}
