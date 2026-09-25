package tui

import (
	"Gyscope/internal/domain/cpu"
	"Gyscope/internal/domain/disk"
	"Gyscope/internal/domain/memory"

	tea "charm.land/bubbletea/v2"
)

type Model struct {
	cpuReader    CPUReader
	memoryReader MemoryReader
	diskReader   DiskReader
	cpu          cpu.CPU
	memory       memory.Memory
	disk         disk.Disk
	memoryErr    error
	cpuErr       error
	diskErr      error
	width        int
	height       int
}

func NewModel(cpuReader CPUReader, memoryReader MemoryReader, diskReader DiskReader) Model {
	return Model{
		cpuReader:    cpuReader,
		memoryReader: memoryReader,
		diskReader:   diskReader,
	}
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}
