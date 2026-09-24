package tui

import (
	"Gyscope/internal/domain/cpu"
	"Gyscope/internal/domain/memory"

	tea "charm.land/bubbletea/v2"
)

type Model struct {
	cpuReader    CPUReader
	memoryReader MemoryReader
	cpu          cpu.CPU
	memory       memory.Memory
	memoryErr    error
	cpuErr       error
}

func NewModel(cpuReader CPUReader, memoryReader MemoryReader) Model {
	return Model{
		cpuReader:    cpuReader,
		memoryReader: memoryReader,
	}
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}
