package main

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "modernc.org/sqlite"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

const (
	host = "0.0.0.0"
	port = "23234"
)

func main() {

	db, err := sql.Open("sqlite", "./app.db")
	if err != nil {
		log.Error("Unable to read sqlite database", "error", err)
	}

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(SqliteBubbleHandler(db)),
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
func SqliteBubbleHandler(db *sql.DB) func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		// This should never fail, as we are using the activeterm middleware.
		pty, _, _ := s.Pty()

		listItems := []list.Item{}

		rows, err := db.Query(`
		SELECT title, rank, username
		FROM posts
		ORDER BY rowid DESC
	`)
		if err != nil {
			return nil, nil
		}
		defer rows.Close()
		i := 0
		for rows.Next() {
			i += 1
			var p Post
			if i >= 20 {
				break
			}
			if err := rows.Scan(&p.Title, &p.Rank, &p.Username); err != nil {
				return nil, nil
			}
			listItems = append(listItems, Post{
				Title:     p.Title,
				Rank:      p.Rank,
				RankDelta: 0,
				Username:  p.Username,
			})
		}

		if err := rows.Err(); err != nil {
			return nil, nil
		}

		l := list.New(listItems, postDelegate{}, 20, 12)
		l.Title = "echohacks: @" + s.User()
		l.SetShowHelp(false)

		inputStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

		input := textinput.New()
		input.Placeholder = "New Hackathon Idea"
		input.Prompt = ""
		input.Width = 40

		m := model{
			currentUsername: s.User(),
			width:           pty.Window.Width,
			height:          pty.Window.Height,
			list:            l,
			input:           input,
			inputStyle:      inputStyle,
		}

		m.OnNew = func(p Post) {
			_, err := db.Exec(`
			INSERT INTO posts (title, rank, username) 
			VALUES (?, ?, ?)
		`, p.Title, p.Rank, p.Username)
			if err != nil {
				log.Printf("Error inserting new post: %v", err)
			}
		}
		m.OnUpdate = func(p Post) {
			_, err := db.Exec(`
			UPDATE posts 
			SET rank = ?
			WHERE title = ? AND username = ?
		`, p.Rank+p.RankDelta, p.Title, p.Username)
			if err != nil {
				log.Printf("Error updating post: %v", err)
			}
		}

		return m, []tea.ProgramOption{tea.WithAltScreen()}
	}
}
