package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the terminal UI state
type Model struct {
	term   string
	width  int
	height int
	time   time.Time
}

// TimeMsg represents a time update message
type TimeMsg time.Time

// NewModel creates a new Model instance
func NewModel(term string, width, height int) Model {
	return Model{
		term:   term,
		width:  width,
		height: height,
		time:   time.Now(),
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles model updates based on received messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TimeMsg:
		m.time = time.Time(msg)
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the model
func (m Model) View() string {
	s := "Your term is %s\n"
	s += "Your window size is x: %d y: %d\n"
	s += "Time: " + m.time.Format(time.RFC1123) + "\n\n"
	s += "Press 'q' to quit\n"
	return fmt.Sprintf(s, m.term, m.width, m.height)
}
