package tui

import (
	domaincpu "Gyscope/internal/domain/cpu"
	domaindisk "Gyscope/internal/domain/disk"
	domainmemory "Gyscope/internal/domain/memory"
	domainprocess "Gyscope/internal/domain/process"
	domaingpu "Gyscope/internal/domain/gpu"

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

	topPanelWidth := innerWidth / 3
	lastTopPanelWidth := innerWidth - 2*topPanelWidth

	if topPanelWidth <= 0 {
		return altScreenView("Terminal is too small")
	}

	cpuPanel := renderCPU(m.cpu, topPanelWidth, panelHeight)
	memoryPanel := renderMemory(m.memory, topPanelWidth, panelHeight)
	diskPanel := renderDisk(m.disk, lastTopPanelWidth, panelHeight)

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

	row2Height := m.height - fixedOverhead

	if row2Height <= panelFrame+3 {
		return altScreenView("Terminal is too small")
	}

	row2ContentHeight := row2Height - panelFrame

	gpuPanelWidth := topPanelWidth
	processPanelWidth := innerWidth - gpuPanelWidth

	gpuPanel := renderGPU(m.gpus, gpuPanelWidth, row2ContentHeight)

	viewportHeight := row2ContentHeight - 2

	m.processViewport.SetWidth(processPanelWidth - panelStyle.GetHorizontalFrameSize())
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
		processPanelWidth,
		row2ContentHeight,
	)

	row2 := lipgloss.JoinHorizontal(
		lipgloss.Top,
		gpuPanel,
		processPanel,
	)

	dashboard := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		panels,
		"",
		row2,
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
	lines := strings.Split(content, "\n")

	switch {
	case len(lines) > height:
		lines = lines[:height]
	case len(lines) < height:
		lines = append(lines, make([]string, height-len(lines))...)
	}

	return panelStyle.
		Width(width).
		Height(height).
		Render(strings.Join(lines, "\n"))
}

func usageStyle(percent float64) lipgloss.Style {
	switch {
	case percent < 20:
		return usageGreenStyle

	case percent >= 20 && percent < 40:
		return usageGreenYellowStyle

	case percent >= 40 && percent < 60:
		return usageYellowStyle

	case percent >= 60 && percent < 80:
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
	style := usageStyle(percent)

	var bar strings.Builder

	bar.WriteString(
		style.Render(strings.Repeat("█", filled)),
	)

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
		fmt.Sprintf(
			"%-8s %-16s %7s   %-13s",
			"PID",
			"NAME",
			"CPU",
			"MEMORY",
		),
	}

	for _, process := range processes {
		lines = append(
			lines,
			fmt.Sprintf(
				"%-8d %-16s %7.1f%%   %-13s",
				process.PID,
				process.Name,
				process.CPUUsage,
				formatBytes(process.Memory),
			),
		)
	}

	return strings.Join(lines, "\n")
}

func renderGPU(gpus []domaingpu.GPU, width, height int) string {
	if len(gpus) == 0 {
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			labelStyle.Render("GPU"),
			"",
			"No GPU detected",
		)
		return renderPanel(content, width, height)
	}

	barWidth := width - 6

	lines := []string{labelStyle.Render("GPU")}

	for i, g := range gpus {
		name := g.Name
		if name == "" {
			name = fmt.Sprintf("GPU %d", i)
		}

		usage := usageStyle(g.Usage).Render(
			fmt.Sprintf("%.1f%%", g.Usage),
		)

		lines = append(
			lines,
			"",
			valueStyle.Render(name),
			fmt.Sprintf("Usage  %s", usage),
			renderProgressBar(g.Usage, barWidth),
			fmt.Sprintf(
				"Mem    %s / %s",
				formatBytes(g.MemoryUsed),
				formatBytes(g.MemoryTotal),
			),
		)
	}

	return renderPanel(strings.Join(lines, "\n"), width, height)
}

func formatBytes(bytes uint64) string {
	const gb = 1024 * 1024 * 1024
	const mb = 1024 * 1024
	const kb = 1024

	switch {
		case bytes > gb:

			return fmt.Sprintf("%4.2f GB", float64(bytes)/gb)

		case bytes > mb && bytes < gb:
			return fmt.Sprintf("%4.0f MB", float64(bytes)/mb)

		default:
			return fmt.Sprintf("%4.0f KB", float64(bytes)/kb)
	}
}


func altScreenView(s string) tea.View {
	v := tea.NewView(s)
	v.AltScreen = true
	return v
}
