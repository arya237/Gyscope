package main

import (
	applicationcpu "Gyscope/internal/application/cpu"
	applicationmemory "Gyscope/internal/application/memory"
	"Gyscope/internal/delivery/tui"
	"fmt"

	Linuxcpu "Gyscope/internal/infrastructure/cpu/linux"
	Linuxmemory "Gyscope/internal/infrastructure/memory/linux"

	tea "charm.land/bubbletea/v2"
)

func main() {

	fmt.Print("\033[2J\033[H")

	cpuDataSource := Linuxcpu.NewLinuxCPUDataSource()
	memoryDataSource := Linuxmemory.NewLinuxMemoryDataSoruce()

	cpuUsecase := applicationcpu.NewCpuUseCase(cpuDataSource)
	memoryUsecase := applicationmemory.NewMemoryUseCase(memoryDataSource)

	model := tui.NewModel(cpuUsecase, memoryUsecase)

	program := tea.NewProgram(model)

	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
