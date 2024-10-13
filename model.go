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
	if post.Rank > 0 {
		postCount = positiveNumStyle.Render(fmt.Sprintf("%d", post.Rank))
	} else {
		postCount = negativeNumStyle.Render(fmt.Sprintf("%d", post.Rank))
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = selectedItemStyle.Render
	}

	fmt.Fprintln(w, fmt.Sprintf("   ▲ %s ▼ %s \n", postCount, fn(post.Username+" - "+post.Title)))
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
		m.list.SetHeight(msg.Height - 8)
	case tea.KeyMsg:
		index := m.list.Index()
		if index >= 0 {
			if item, ok := m.list.SelectedItem().(Post); ok {
				switch msg.String() {
				case "up": // Upvote
					if item.RankDelta == -1 { // User previously downvoted
						item.Rank++   // Remove the downvote
					}
					if item.RankDelta != 1 { // User has not upvoted yet
						item.Rank++       // Add upvote
						item.RankDelta = 1 // Mark as upvoted
					}
					m.list.SetItem(index, item) // Update the selected item
				case "down": // Downvote
					if item.RankDelta == 1 { // User previously upvoted
						item.Rank--   // Remove the upvote
					}
					if item.RankDelta != -1 { // User has not downvoted yet
						item.Rank--       // Add downvote
						item.RankDelta = -1 // Mark as downvoted
					}
					m.list.SetItem(index, item) // Update the selected item
				case "q", "ctrl+c":
					return m, tea.Quit
				}
			}
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
