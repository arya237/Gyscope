package tui

import (
	domaincpu "Gyscope/internal/domain/cpu"
	domaindisk "Gyscope/internal/domain/disk"
	domainmemory "Gyscope/internal/domain/memory"
	domainprocess "Gyscope/internal/domain/process"

	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m Model) View() tea.View {
	if m.width <= 0 {
		return altScreenView("Starting Gyscope...")
	}

	panelHeight := 9

	innerWidth := m.width - dashboardStyle.GetHorizontalFrameSize()

	panelWidth := innerWidth / 3
	lastPanelWidth := innerWidth - 2*panelWidth

	if panelWidth <= 0 {
		return altScreenView("Terminal is too small")
	}

	cpuPanel := renderCPU(m.cpu, panelWidth, panelHeight)
	memoryPanel := renderMemory(m.memory, panelWidth, panelHeight)
	diskPanel := renderDisk(m.disk, lastPanelWidth, panelHeight)

	panels := lipgloss.JoinHorizontal(
		lipgloss.Top,
		cpuPanel,
		memoryPanel,
		diskPanel,
	)

	header := renderHeader(innerWidth)

	footer := lipgloss.NewStyle().
		Width(innerWidth).
		Align(lipgloss.Center).
		Render("Refreshing every 1s  •  q Quit")

	blankLines := 3

	dashboardFrame := dashboardStyle.GetVerticalFrameSize()
	panelFrame := panelStyle.GetVerticalFrameSize()

	fixedOverhead := lipgloss.Height(header) +
		lipgloss.Height(panels) +
		lipgloss.Height(footer) +
		dashboardFrame +
		blankLines

	processHeight := m.height - fixedOverhead

	if processHeight <= panelFrame+3 {
		return altScreenView("Terminal is too small")
	}

	processContentHeight := processHeight - panelFrame
	viewportHeight := processContentHeight - 2

	m.processViewport.SetWidth(innerWidth - panelStyle.GetHorizontalFrameSize())
	m.processViewport.SetHeight(viewportHeight)
	m.processViewport.SetContent(renderProcesses(m.processes))

	processContent := lipgloss.JoinVertical(
		lipgloss.Left,
		labelStyle.Render("PROCESSES"),
		"",
		m.processViewport.View(),
	)

	processPanel := renderPanel(
		processContent,
		innerWidth,
		processContentHeight,
	)

	dashboard := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		panels,
		"",
		processPanel,
		"",
		footer,
	)

	content := dashboardStyle.
		Width(m.width).
		Render(dashboard)

	return altScreenView(content)
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
		Render(content)
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

	temperature := usageStyle(cpu.Temperature).Render(
		fmt.Sprintf("%.1f°C", cpu.Temperature),
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

		fmt.Sprintf("Temp	   %s", temperature),
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

func renderProcesses(
	processes []domainprocess.Process,
) string {
	lines := []string{
		"PID      NAME                 CPU       MEMORY",
	}

	for _, process := range processes {
		lines = append(
			lines,
			fmt.Sprintf(
				"%-8d %-16s %7.1f%% %11.1f%%",
				process.PID,
				process.Name,
				process.CPUUsage,
				process.MemoryUsage,
			),
		)
	}

	return strings.Join(lines, "\n")
}

func formatBytes(bytes uint64) string {
	const gb = 1024 * 1024 * 1024

	return fmt.Sprintf("%.2f GB", float64(bytes)/gb)
}

func altScreenView(s string) tea.View {
	v := tea.NewView(s)
	v.AltScreen = true
	return v
}
