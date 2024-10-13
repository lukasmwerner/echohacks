package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	itemStyle         = lipgloss.NewStyle()
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
	positiveNumStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	negativeNumStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
)

type Post struct {
	Title     string
	Rank      int
	RankDelta int
	Username  string
}

func (p Post) FilterValue() string { return p.Title }

type postDelegate struct {
}

func (d postDelegate) Height() int                               { return 1 }
func (d postDelegate) Spacing() int                              { return 0 }
func (d postDelegate) Update(msg tea.Msg, l *list.Model) tea.Cmd { return nil }
func (d postDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	post, ok := listItem.(Post)
	if !ok {
		return
	}

	postCount := ""
	if post.Rank > 0 {
		postCount = positiveNumStyle.Render(fmt.Sprintf("%d", post.Rank))
	} else {
		postCount = negativeNumStyle.Render(fmt.Sprintf("%d", post.Rank))
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = selectedItemStyle.Render
	}

	fmt.Fprintln(w, fn(fmt.Sprintf("▲ %s ▼ %s - %s\n", postCount, post.Username, post.Title)))

}

// Just a generic tea.Model to demo terminal information of ssh.
type model struct {
	list       list.Model
	input      textinput.Model
	inputStyle lipgloss.Style
	OnNew      func(Post)
	OnUpdate   func(Post)
	OnDelete   func(Post)
	Refresh    func() []Post
	width      int
	height     int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 1)
	case tea.KeyMsg:
		switch msg.String() {
		case "+", "=":
		case "-", "_":
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	if !m.input.Focused() {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {

	list := m.list.View()

	return "echohacks\n" + list
}
