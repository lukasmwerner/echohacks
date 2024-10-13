package main

import (
	"fmt"


"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
)
var(
	positiveNumStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("00FF00"));
	negativeNumStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("FF0000"));
	
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
	postCount := ""
	if m.postcount > 0{
		postCount = positiveNumStyle.Render(fmt.Sprintf("%d", m.postcount))


	} else {
		postCount = negativeNumStyle.Render(fmt.Sprintf("%d", m.postcount))
	}

	s += fmt.Sprintf("▲ %s ▼ %s- %s\n", postCount, "blaster", "tiktok for lecture videos")

	return s
}
