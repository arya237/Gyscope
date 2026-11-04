package tui

import (
	"Gyscope/internal/domain/cpu"
	"Gyscope/internal/domain/disk"
	"Gyscope/internal/domain/gpu"
	"Gyscope/internal/domain/memory"
	"Gyscope/internal/domain/network"
	"Gyscope/internal/domain/process"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

type Model struct {
	cpuReader     CPUReader
	memoryReader  MemoryReader
	diskReader    DiskReader
	gpuReader     GPUReader
	processReader ProcessReader
	netowrkReader NetowrkReader

	cpu       cpu.CPU
	memory    memory.Memory
	disk      disk.Disk
	gpus      []gpu.GPU
	processes []process.Process
	network   network.Network

	memoryErr       error
	cpuErr          error
	diskErr         error
	gpuErr          error
	processErr      error
	netwrokErr      error
	processViewport viewport.Model
	width           int
	height          int
}

func NewModel(
	cpuReader CPUReader,
	memoryReader MemoryReader,
	diskReader DiskReader,
	gpuReader GPUReader,
	process ProcessReader,
	Network NetowrkReader,
) Model {
	return Model{
		cpuReader:       cpuReader,
		memoryReader:    memoryReader,
		diskReader:      diskReader,
		gpuReader:       gpuReader,
		processReader:   process,
		netowrkReader:   Network,
		processViewport: viewport.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}
