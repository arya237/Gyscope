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
)
