package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"

	domaincpu "Gyscope/internal/domain/cpu"
	domainmemory "Gyscope/internal/domain/memory"
)

func (m Model) View() tea.View {
	content := renderCPU(m.cpu) +
		"\n\n" +
		renderMemory(m.memory)

	return tea.NewView(content)
}

func renderCPU(cpu domaincpu.CPU) string {
	return fmt.Sprintf(
		"CPU\n"+
			"  Usage: %.2f%%\n"+
			"  Cores: %d\n"+
			"  Load: %.2f %.2f %.2f",
		cpu.Usage,
		cpu.LogicalCores,
		cpu.LoadAverage.OneMinute,
		cpu.LoadAverage.FiveMinutes,
		cpu.LoadAverage.FifteenMinutes,
	)
}

func renderMemory(memory domainmemory.Memory) string {
	return fmt.Sprintf(
		"Memory\n"+
			"  Usage: %.2f%%\n"+
			"  Used: %s\n"+
			"  Available: %s\n"+
			"  Total: %s",
		memory.Usage,
		formatBytes(memory.Used),
		formatBytes(memory.Available),
		formatBytes(memory.Total),
	)
}

func formatBytes(bytes uint64) string {
	const gb = 1024 * 1024 * 1024

	return fmt.Sprintf("%.2f GB", float64(bytes)/gb)
}
