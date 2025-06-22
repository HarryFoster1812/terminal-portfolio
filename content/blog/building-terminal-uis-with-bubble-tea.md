---
title: "Building Terminal UIs with Bubble Tea"
summary: "A deep dive into creating beautiful terminal applications using Charm's Bubble Tea framework."
date: "2024-01-15"
tags: ["go", "tui", "bubble-tea", "terminal"]
readTime: "8 min read"
author: "Arpan Pandey"
published: true
---

# Building Terminal UIs with Bubble Tea

Bubble Tea is an amazing framework for building terminal user interfaces in Go. It follows the Elm architecture pattern, making it easy to reason about state management and updates.

## Why Terminal UIs?

Terminal applications have several advantages:
- **Lightning fast**: No rendering overhead of web browsers
- **Universal**: Work on any system with a terminal
- **Lightweight**: Minimal resource usage
- **Keyboard-first**: Optimized for power users
- **SSH-friendly**: Can be accessed remotely

## The Elm Architecture

Bubble Tea is built around three core concepts:

### 1. Model
The Model represents your application's state. It's a struct that contains all the data your application needs.

```go
type Model struct {
    currentPage int
    items       []string
    selected    int
}
```

### 2. Update
The Update function handles all the messages (events) that can occur in your application.

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
        case "up":
            if m.selected > 0 {
                m.selected--
            }
        case "down":
            if m.selected < len(m.items)-1 {
                m.selected++
            }
        }
    }
    return m, nil
}
```

### 3. View
The View function renders your model to a string that will be displayed in the terminal.

```go
func (m Model) View() string {
    s := "My Todo App\n\n"
    for i, item := range m.items {
        cursor := " "
        if m.selected == i {
            cursor = ">"
        }
        s += fmt.Sprintf("%s %s\n", cursor, item)
    }
    s += "\nPress q to quit.\n"
    return s
}
```

## Advanced Features

### Styling with Lipgloss
Lipgloss provides powerful styling capabilities:

```go
var titleStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#FAFAFA")).
    Background(lipgloss.Color("#7D56F4")).
    PaddingTop(1).
    PaddingLeft(4)
```

### Components from Bubbles
The Bubbles library provides ready-made components:
- **List**: For displaying and selecting from lists
- **Textinput**: For text input fields  
- **Textarea**: For multi-line text input
- **Viewport**: For scrollable content
- **Progress**: For progress bars
- **Spinner**: For loading indicators

### SSH Integration with Wish
Wish makes it easy to serve your TUI over SSH:

```go
func main() {
    s, err := wish.NewServer(
        wish.WithAddress("localhost:2222"),
        wish.WithMiddleware(
            bubbletea.Middleware(func() tea.Model {
                return NewModel()
            }),
        ),
    )
    if err != nil {
        log.Fatal(err)
    }
    log.Fatal(s.ListenAndServe())
}
```

## Best Practices

1. **Keep models simple**: Store only what you need
2. **Handle all messages**: Even if you don't act on them
3. **Use commands for side effects**: File I/O, HTTP requests, etc.
4. **Style consistently**: Define your styles once and reuse
5. **Test your components**: Unit test your update logic

## Conclusion

Bubble Tea opens up a world of possibilities for terminal applications. Whether you're building developer tools, system monitors, or interactive CLIs, it provides a solid foundation with excellent developer experience.

The combination of Go's performance, Bubble Tea's architecture, and Lipgloss's styling makes for powerful and beautiful terminal applications.

Happy coding! 🎉