package tui

import (
	"Gyscope/internal/application/process"
	"Gyscope/internal/domain/cpu"
	"Gyscope/internal/domain/disk"
	"Gyscope/internal/domain/memory"

	tea "charm.land/bubbletea/v2"
)

type Model struct {
	cpuReader    CPUReader
	memoryReader MemoryReader
	diskReader   DiskReader
	processReader ProcessReader
	cpu          cpu.CPU
	memory       memory.Memory
	disk         disk.Disk
	process      process.RawProcessState
	memoryErr    error
	cpuErr       error
	diskErr      error
	processErr   error
	width        int
	height       int
}

func NewModel(
	cpuReader CPUReader,
	memoryReader MemoryReader,
	diskReader DiskReader,
	process processReader
) Model {
	return Model{
		cpuReader:    cpuReader,
		memoryReader: memoryReader,
		diskReader:   diskReader,
	}
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}
