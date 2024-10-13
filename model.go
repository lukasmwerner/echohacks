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
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#5F5FDF")).Bold(true)
	positiveNumStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	negativeNumStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	neutralNumStyle   = lipgloss.NewStyle()
)

type Post struct {
	Title     string
	Rank      int
	RankDelta int // This will be used to track the user's vote
	Username  string
}

func (p Post) FilterValue() string { return p.Title }

type postDelegate struct{}

func (d postDelegate) Height() int                               { return 1 }
func (d postDelegate) Spacing() int                              { return 0 }
func (d postDelegate) Update(msg tea.Msg, l *list.Model) tea.Cmd { return nil }
func (d postDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	post, ok := listItem.(Post)
	if !ok {
		return
	}

	postCount := ""
	if post.Rank+post.RankDelta > 0 {
		postCount = positiveNumStyle.Render(fmt.Sprintf("%d", post.Rank+post.RankDelta))
	} else if post.Rank+post.RankDelta < 0 {
		postCount = negativeNumStyle.Render(fmt.Sprintf("%d", post.Rank+post.RankDelta))
	} else{
		postCount = neutralNumStyle.Render(fmt.Sprintf("%d", post.Rank+post.RankDelta))
	}
	fn := itemStyle.Render
	if index == m.Index() {
		fn = selectedItemStyle.Render
	}

	fmt.Fprintln(w, fmt.Sprintf("   ▲ %s ▼ %s \n", postCount, fn(post.Username+" - "+post.Title)))

}

// Just a generic tea.Model to demo terminal information of ssh.
type model struct {
	currentUsername string
	list            list.Model
	input           textinput.Model
	inputStyle      lipgloss.Style
	OnNew           func(Post)
	OnUpdate        func(Post)
	OnDelete        func(Post)
	Refresh         func() []Post
	width           int
	height          int
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
		m.list.SetHeight(msg.Height - 8)
	case tea.KeyMsg:
		switch msg.String() {
		case "+", "=":
			post := m.list.SelectedItem().(Post)
			if post.RankDelta == 0 || post.RankDelta == -1 {
				post.RankDelta += 1
			}
			m.list.SetItem(m.list.Index(), post)
		case "-", "_":
			post := m.list.SelectedItem().(Post)
			if post.RankDelta == 1 || post.RankDelta == 0 {
				post.RankDelta -= 1
			}
			m.list.SetItem(m.list.Index(), post)
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
	return list
}
