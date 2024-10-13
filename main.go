package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

const (
	host = "localhost"
	port = "23234"
)

func main() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

// You can wire any Bubble Tea model up to the middleware with a function that
// handles the incoming ssh.Session. Here we just grab the terminal info and
// pass it to the new model. You can also return tea.ProgramOptions (such as
// tea.WithAltScreen) on a session by session basis.
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// This should never fail, as we are using the activeterm middleware.
	pty, _, _ := s.Pty()

	listItems := []list.Item{}

	listItems = append(listItems, Post{
		Title:     "blasterhacks in rust btw",
		Rank:      0,
		RankDelta: 0,
		Username:  "dogwater",
	})

	listItems = append(listItems, Post{
		Title:     "wargames but in pascal",
		Rank:      0,
		RankDelta: 0,
		Username:  "redhot",
	})

	l := list.New(listItems, postDelegate{}, 20, 12)
	l.Title = "echohacks"
	//inputStyle := lipgloss.NewStyle().
	//	Border(lipgloss.RoundedBorder()).
	//	BorderForeground(lipgloss.Color("#874BFD")).
	//	Padding(1).
	//	BorderTop(true).
	//	BorderLeft(true).
	//	BorderRight(true).
	//	BorderBottom(true)

	//input := textinput.New()
	//input.Placeholder = "New Hackathon Idea"
	//input.Prompt = ""
	//input.Width = 40

	m := model{
		width:  pty.Window.Width,
		height: pty.Window.Height,
		list:   l,
	}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
