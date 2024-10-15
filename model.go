package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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

type keyMap struct {
	Compose  key.Binding
	Uprank   key.Binding
	Downrank key.Binding
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Help     key.Binding
	Quit     key.Binding
	Submit   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Compose, k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Compose, k.Uprank, k.Downrank},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	Compose:  key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new idea")),
	Uprank:   key.NewBinding(key.WithKeys("+", "="), key.WithHelp("+", "uprank")),
	Downrank: key.NewBinding(key.WithKeys("-", "_"), key.WithHelp("-", "downrank")),
	Up:       key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "move up")),
	Down:     key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "move down")),
	Left:     key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "move left")),
	Right:    key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "move right")),
	Help:     key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
	Quit:     key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
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
	} else {
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
	help            help.Model
	input           textinput.Model
	inputStyle      lipgloss.Style
	OnNew           func(Post)
	OnUpdate        func(Post)
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

	// Make sure to update the list if we aren't doing input
	if !m.input.Focused() {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.list.SetWidth(msg.Width)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Compose):
			m.input.Focus()
			cmds = append(cmds, textinput.Blink)

			keys.Quit.SetEnabled(false)
			keys.Downrank.SetEnabled(false)
			keys.Uprank.SetEnabled(false)
			keys.Help.SetEnabled(false)

		case key.Matches(msg, keys.Uprank):
			post := m.list.SelectedItem().(Post)
			if post.RankDelta == 0 || post.RankDelta == -1 {
				post.RankDelta += 1
			}
			m.list.SetItem(m.list.Index(), post)
			if m.OnUpdate != nil {
				m.OnUpdate(post)
			}
		case key.Matches(msg, keys.Downrank):
			post := m.list.SelectedItem().(Post)
			if post.RankDelta == 1 || post.RankDelta == 0 {
				post.RankDelta -= 1
			}
			m.list.SetItem(m.list.Index(), post)
			if m.OnUpdate != nil {
				m.OnUpdate(post)
			}
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))) && m.input.Focused():
			m.input.Blur()
			keys.Quit.SetEnabled(true)
			keys.Downrank.SetEnabled(true)
			keys.Uprank.SetEnabled(true)
			keys.Help.SetEnabled(true)

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))) && m.input.Focused():
			if m.input.Value() == "" {
				break
			}
			post := Post{
				Title:     m.input.Value(),
				Rank:      0,
				RankDelta: 0,
				Username:  m.currentUsername,
			}
			m.input.SetValue("")
			m.input.Blur()
			keys.Quit.SetEnabled(true)
			keys.Downrank.SetEnabled(true)
			keys.Uprank.SetEnabled(true)
			keys.Help.SetEnabled(true)

			cmd = m.list.InsertItem(len(m.list.Items()), post)
			cmds = append(cmds, cmd)
			if m.OnNew != nil {
				m.OnNew(post)
			}

		case key.Matches(msg, keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.input.Focused() {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			m.inputStyle.Render(m.input.View()),
		)
	}

	list := m.list.View()

	app := m.inputStyle.Render(list + "\n" + m.help.View(keys))
	app = lipgloss.PlaceVertical(m.height, lipgloss.Center, app)

	return lipgloss.PlaceHorizontal(m.width, lipgloss.Center, app)
}
