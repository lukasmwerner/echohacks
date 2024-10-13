package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	positiveNumStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	negativeNumStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
)

type Post struct {
	Title    string
	Rank     int
	RankDelta     int
	Username string
}
func (p Post) FilterValue() string { return p.Title }

// Just a generic tea.Model to demo terminal information of ssh.
type model struct {
	list       list.Model
	input      textinput.Model
	inputStyle lipgloss.Style
	OnNew      func(item)
	OnUpdate   func(item)
	OnDelete   func(item)
	Refresh    func() []item
	width	  int
	height	  int
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
		case "k", "up":
			selected -= 1
		case "j", "down":
			selected += 1
		case "+", "=":
			//if RankDelta = 0 {
				//m.posts[selected] += 1
			//}
		case "-", "_":
			//if RankDelta = 0 {
				//m.posts[selected] -= 1
			//}
			//m.postcount -= 1
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {

	s := ""
	postCount := ""

	if m.postcount > 0 {
		postCount = positiveNumStyle.Render(fmt.Sprintf("%d", m.postcount))
	} else {
		postCount = negativeNumStyle.Render(fmt.Sprintf("%d", m.postcount))
	}

	s += fmt.Sprintf("▲ %s ▼ %s- %s\n", postCount, "blaster", "tiktok for lecture videos")

	return s
}
