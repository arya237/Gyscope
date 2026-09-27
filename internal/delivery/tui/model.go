package tui

import (
	"Gyscope/internal/domain/cpu"
	"Gyscope/internal/domain/disk"
	"Gyscope/internal/domain/memory"
	"Gyscope/internal/domain/process"

	tea "charm.land/bubbletea/v2"
)

type Model struct {
	cpuReader     CPUReader
	memoryReader  MemoryReader
	diskReader    DiskReader
	processReader ProcessReader
	cpu           cpu.CPU
	memory        memory.Memory
	disk          disk.Disk
	processes     []process.Process
	memoryErr     error
	cpuErr        error
	diskErr       error
	processErr    error
	width         int
	height        int
}

func NewModel(
	cpuReader CPUReader,
	memoryReader MemoryReader,
	diskReader DiskReader,
	process ProcessReader,
) Model {
	return Model{
		cpuReader:     cpuReader,
		memoryReader:  memoryReader,
		diskReader:    diskReader,
		processReader: process,
	}
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}
