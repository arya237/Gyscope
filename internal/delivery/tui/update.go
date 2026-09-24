package tui

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

type systemTickMsg struct{}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return systemTickMsg{}
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case systemTickMsg:
		cpuState, err := m.cpuReader.GetState()
		if err != nil {
			m.cpuErr = err
		} else {
			m.cpu = cpuState
			m.cpuErr = nil
		}

		memoryState, err := m.memoryReader.GetState()
		if err != nil {
			m.memoryErr = err
		} else {
			m.memory = memoryState
			m.memoryErr = nil
		}

		return m, tickCmd()
	}

	return m, nil
}
