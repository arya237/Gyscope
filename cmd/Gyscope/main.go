package main

import (
	applicationcpu "Gyscope/internal/application/cpu"
	applicationdisk "Gyscope/internal/application/disk"
	applicationgpu "Gyscope/internal/application/gpu"
	applicationmemory "Gyscope/internal/application/memory"
	applicationmonitor "Gyscope/internal/application/monitor"
	applicationnetwork "Gyscope/internal/application/network"
	applicationprocess "Gyscope/internal/application/process"
	"context"
	"time"

	"Gyscope/internal/delivery/tui"
	"fmt"

	Linuxcpu "Gyscope/internal/infrastructure/cpu/linux"
	Linuxdisk "Gyscope/internal/infrastructure/disk/linux"
	Linuxgpu "Gyscope/internal/infrastructure/gpu/linux"
	Linuxmemory "Gyscope/internal/infrastructure/memory/linux"
	Linuxnetwork "Gyscope/internal/infrastructure/network/linux"
	Linuxprocess "Gyscope/internal/infrastructure/process/linux"

	tea "charm.land/bubbletea/v2"
)

func main() {

	fmt.Print("\033[2J\033[H")

	cpuDataSource := Linuxcpu.NewLinuxCPUDataSource()
	memoryDataSource := Linuxmemory.NewLinuxMemoryDataSoruce()
	diskDatasource := Linuxdisk.NewLinuxDiskDataSource()
	processDatasource := Linuxprocess.NewLinuxProcessDataSource()
	gpuDatasource := Linuxgpu.NewLinuxGPUDataSource()
	networkDatasource := Linuxnetwork.NewLinuxNetworkDatasource()

	cpuUsecase := applicationcpu.NewCpuUseCase(cpuDataSource)
	memoryUsecase := applicationmemory.NewMemoryUseCase(memoryDataSource)
	diskUsecase := applicationdisk.NewDiskUseCase(diskDatasource)
	processUsecase := applicationprocess.NewProcessUseCase(processDatasource)
	gpuUsecase := applicationgpu.NewGPUUsecase(gpuDatasource)
	networkUsecase := applicationnetwork.NewNetworkUseCase(networkDatasource)

	collector := applicationmonitor.NewCollector(
		applicationmonitor.Readers{
			CPU:     cpuUsecase,
			GPU:     gpuUsecase,
			Memory:  memoryUsecase,
			Process: processUsecase,
			Network: networkUsecase,
			Disk:    diskUsecase,
		},
		time.Second,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go collector.Run(ctx)

	model := tui.NewModel(collector)

	program := tea.NewProgram(model)

	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
