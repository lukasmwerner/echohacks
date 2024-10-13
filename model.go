package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Just a generic tea.Model to demo terminal information of ssh.
type model struct {
	height    int
	width     int
	postcount int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "+", "=":
			m.postcount += 1
		case "-", "_":
			m.postcount -= 1
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := ""

	s += fmt.Sprintf("▲ %d ▼ %s- %s\n", m.postcount, "blaster", "tiktok for lecture videos")

	return s
}
