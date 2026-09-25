package tui

import (
	domaincpu "Gyscope/internal/domain/cpu"
	domaindisk "Gyscope/internal/domain/disk"
	domainmemory "Gyscope/internal/domain/memory"

	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m Model) View() tea.View {
	if m.width <= 0 {
		return tea.NewView("Starting Gyscope...")
	}

	gap := 2
	panelWidth := (m.width - 2*gap - 6) / 3
	panelHeight := 9

	if panelWidth <= 0 {
		return tea.NewView("Terminal is too small")
	}

	cpuPanel := renderCPU(
		m.cpu,
		panelWidth,
		panelHeight,
	)

	memoryPanel := renderMemory(
		m.memory,
		panelWidth,
		panelHeight,
	)

	diskPanel := renderDisk(
		m.disk,
		panelWidth,
		panelHeight,
	)

	panels := lipgloss.JoinHorizontal(
		lipgloss.Top,
		cpuPanel,
		memoryPanel,
		diskPanel,
	)

	header := renderHeader(m.width)

	footer := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render("Refreshing every 1s  •  q Quit")

	dashboard := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		panels,
		"",
		footer,
	)

	content := dashboardStyle.
		Width(m.width).
		Render(dashboard)

	return tea.NewView(content)
}

func renderHeader(width int) string {
	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				titleStyle.Render("GYSCOPE"),
				subtitleStyle.Render("Go based Linux System Monitor"),
			),
		)
}

func renderPanel(content string, width, height int) string {
	return panelStyle.
		Width(width).
		Height(height).
		Render(
			lipgloss.Place(
				width,
				height,
				lipgloss.Left,
				lipgloss.Top,
				content,
			),
		)
}

func usageStyle(percent float64) lipgloss.Style {
	switch {
	case percent < 25:
		return usageGreenStyle

	case percent >= 25 && percent < 50:
		return usageGreenYellowStyle

	case percent >= 50 && percent < 75:
		return usageYellowStyle

	case percent >= 75 && percent < 85:
		return usageOrangeStyle

	default:
		return usageRedStyle
	}
}

func renderProgressBar(percent float64, width int) string {
	if width <= 0 {
		return ""
	}

	if percent < 0 {
		percent = 0
	}

	if percent > 100 {
		percent = 100
	}

	filled := int(percent / 100 * float64(width))

	var bar strings.Builder

	for i := 0; i < filled; i++ {
		position := float64(i) / float64(width) * 100

		bar.WriteString(usageStyle(position).Render("█"))
	}

	bar.WriteString(
		progressEmptyStyle.Render(
			strings.Repeat("░", width-filled),
		),
	)

	return bar.String()
}

func renderCPU(
	cpu domaincpu.CPU,
	width int,
	height int,
) string {
	barWidth := width - 6

	usage := usageStyle(cpu.Usage).Render(
		fmt.Sprintf("%.1f%%", cpu.Usage),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,

		labelStyle.Render("CPU"),

		fmt.Sprintf(
			"Usage      %s",
			usage,
		),

		renderProgressBar(
			cpu.Usage,
			barWidth,
		),

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

	return renderPanel(
		content,
		width,
		height,
	)
}

func renderMemory(
	memory domainmemory.Memory,
	width int,
	height int,
) string {
	barWidth := width - 6

	usage := usageStyle(memory.Usage).Render(
		fmt.Sprintf("%.1f%%", memory.Usage),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,

		labelStyle.Render("MEMORY"),

		fmt.Sprintf(
			"Usage      %s",
			usage,
		),

		renderProgressBar(
			memory.Usage,
			barWidth,
		),

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

	return renderPanel(
		content,
		width,
		height,
	)
}

func renderDisk(disk domaindisk.Disk, width int, height int) string {
	barWidth := width - 6

	usage := usageStyle(disk.Usage).Render(
		fmt.Sprintf("%.1f%%", disk.Usage),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		labelStyle.Render("DISK"),
		fmt.Sprintf("Usage      %s", usage),
		renderProgressBar(disk.Usage, barWidth),
		"",
		fmt.Sprintf("Used       %s", formatBytes(disk.Used)),
		fmt.Sprintf("Free       %s", formatBytes(disk.Free)),
		fmt.Sprintf("Total      %s", formatBytes(disk.Total)),
	)

	return renderPanel(content, width, height)
}

func formatBytes(bytes uint64) string {
	const gb = 1024 * 1024 * 1024

	return fmt.Sprintf("%.2f GB", float64(bytes)/gb)
}
