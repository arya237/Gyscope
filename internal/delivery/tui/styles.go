package tui

import "charm.land/lipgloss/v2"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Align(lipgloss.Center)

	subtitleStyle = lipgloss.NewStyle().
			Align(lipgloss.Center)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2)

	labelStyle = lipgloss.NewStyle().
			Bold(true)

	valueStyle = lipgloss.NewStyle().
			Bold(true)

	footerStyle = lipgloss.NewStyle().
			Align(lipgloss.Center)
	progressFilledStyle = lipgloss.NewStyle().
				Bold(true)

	usageLowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	usageMediumStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("214"))

	usageHighStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	usageGreenStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46"))

	usageGreenYellowStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("118"))

	usageYellowStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("226"))

	usageOrangeStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("208"))

	usageRedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	progressEmptyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240"))

	dashboardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2)
)
