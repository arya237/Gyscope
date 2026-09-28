package main

import (
	applicationcpu "Gyscope/internal/application/cpu"
	applicationdisk "Gyscope/internal/application/disk"
	applicationmemory "Gyscope/internal/application/memory"
	applicationprocess "Gyscope/internal/application/process"

	"Gyscope/internal/delivery/tui"
	"fmt"

	Linuxcpu "Gyscope/internal/infrastructure/cpu/linux"
	linuxdisk "Gyscope/internal/infrastructure/disk/linux"
	Linuxmemory "Gyscope/internal/infrastructure/memory/linux"
	Linuxprocess "Gyscope/internal/infrastructure/process/linux"

	tea "charm.land/bubbletea/v2"
)

func main() {

	fmt.Print("\033[2J\033[H")

	cpuDataSource := Linuxcpu.NewLinuxCPUDataSource()
	memoryDataSource := Linuxmemory.NewLinuxMemoryDataSoruce()
	diskDatasource := linuxdisk.NewLinuxDiskDataSource()
	processDatasource := Linuxprocess.NewLinuxProcessDataSource()

	cpuUsecase := applicationcpu.NewCpuUseCase(cpuDataSource)
	memoryUsecase := applicationmemory.NewMemoryUseCase(memoryDataSource)
	diskUsecase := applicationdisk.NewDiskUseCase(diskDatasource)
	processUsecase := applicationprocess.NewProcessUseCase(processDatasource)


	model := tui.NewModel(cpuUsecase, memoryUsecase, diskUsecase, processUsecase)

	program := tea.NewProgram(model)

	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
