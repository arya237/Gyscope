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
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		return m, nil

	case systemTickMsg:
		snapshot := m.snapshotReader.Snapshot()

		m.cpu = snapshot.CPU
		m.memory = snapshot.Memory
		m.disk = snapshot.Disk
		m.gpus = snapshot.GPUs
		m.processes = snapshot.Processes
		m.network = snapshot.Network

		return m, tickCmd()

	case tea.KeyPressMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	}

	return m, nil
}
