package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	domaincpu "Gyscope/internal/domain/cpu"
	domainmemory "Gyscope/internal/domain/memory"
)

func (m Model) View() tea.View {
	if m.width <= 0 {
		return tea.NewView("Starting Gyscope...")
	}

	gap := 2
	panelWidth := (m.width - gap) / 2

	if panelWidth <= 0 {
		return tea.NewView("Terminal is too small")
	}

	cpuPanel := renderCPU(m.cpu, panelWidth)
	memoryPanel := renderMemory(m.memory, panelWidth)

	panels := lipgloss.JoinHorizontal(
		lipgloss.Top,
		cpuPanel,
		memoryPanel,
	)

	header := renderHeader()

	footer := footerStyle.Render(
		"Refreshing every 1s  •  q Quit",
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		panels,
		"",
		footer,
	)

	return tea.NewView(content)
}

func renderHeader() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render("GYSCOPE"),
		subtitleStyle.Render("Linux System Monitor"),
	)
}

func renderProgressBar(percent float64, width int) string {
	if percent < 0 {
		percent = 0
	}

	if percent > 100 {
		percent = 100
	}

	filled := int(percent / 100 * float64(width))

	return strings.Repeat("█", filled) +
		strings.Repeat("░", width-filled)
}

func renderCPU(cpu domaincpu.CPU, width int) string {
	barWidth := width - 6

	content := lipgloss.JoinVertical(
		lipgloss.Left,

		labelStyle.Render("CPU"),

		fmt.Sprintf(
			"Usage      %.1f%%",
			cpu.Usage,
		),

		renderProgressBar(cpu.Usage, barWidth),

		"",

		fmt.Sprintf(
			"Cores      %d",
			cpu.LogicalCores,
		),

		fmt.Sprintf(
			"Load       %.2f  %.2f  %.2f",
			cpu.LoadAverage.OneMinute,
			cpu.LoadAverage.FiveMinutes,
			cpu.LoadAverage.FifteenMinutes,
		),
	)

	return panelStyle.
		Width(width).
		Render(content)
}

func renderMemory(memory domainmemory.Memory, width int) string {
	barWidth := width - 6

	content := lipgloss.JoinVertical(
		lipgloss.Left,

		labelStyle.Render("MEMORY"),

		fmt.Sprintf(
			"Usage      %.1f%%",
			memory.Usage,
		),

		renderProgressBar(memory.Usage, barWidth),

		"",

		fmt.Sprintf(
			"Used       %s",
			formatBytes(memory.Used),
		),

		fmt.Sprintf(
			"Available  %s",
			formatBytes(memory.Available),
		),

		fmt.Sprintf(
			"Total      %s",
			formatBytes(memory.Total),
		),
	)

	return panelStyle.
		Width(width).
		Render(content)
}

func formatBytes(bytes uint64) string {
	const gb = 1024 * 1024 * 1024

	return fmt.Sprintf("%.2f GB", float64(bytes)/gb)
}
