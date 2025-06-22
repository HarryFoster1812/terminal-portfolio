package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/muesli/termenv"
)

// CustomBubbleteaMiddleware creates a custom Bubble Tea middleware
// that wraps tea.Program with SSH session integration
func CustomBubbleteaMiddleware() wish.Middleware {
	newProg := func(m tea.Model, opts ...tea.ProgramOption) *tea.Program {
		p := tea.NewProgram(m, opts...)
		go func() {
			for {
				<-time.After(1 * time.Second)
				p.Send(TimeMsg(time.Now()))
			}
		}()
		return p
	}

	teaHandler := func(s ssh.Session) *tea.Program {
		pty, _, active := s.Pty()
		if !active {
			wish.Fatalln(s, "no active terminal, skipping")
			return nil
		}

		m := NewModel(
			pty.Term,
			pty.Window.Width,
			pty.Window.Height,
		)

		return newProg(m, append(bubbletea.MakeOptions(s), tea.WithAltScreen())...)
	}

	return bubbletea.MiddlewareWithProgramHandler(teaHandler, termenv.ANSI256)
}
