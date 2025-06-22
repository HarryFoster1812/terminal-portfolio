# 🚀 Terminal Portfolio

A beautiful, interactive terminal-based portfolio application built with Go, Bubble Tea, and Glamour. This application showcases projects, blog posts, and personal information in a stunning terminal interface that's fully keyboard navigable and responsive.

## ✨ Features

### 🎨 Beautiful Terminal UI
- **Responsive design** that adapts to any terminal size
- **Sticky navigation bar** that stays visible while scrolling
- **Smooth scrolling** through long content using viewport
- **Beautiful styling** with Lipgloss and consistent color scheme
- **Context-sensitive help** footer with keybindings

### 📝 Markdown Blog System
- **Glamour-powered rendering** with syntax highlighting
- **Frontmatter support** for metadata (title, date, tags, etc.)
- **File-based content management** - just add `.md` files
- **Auto-discovery** of blog posts from `content/blog/` directory
- **Fallback system** when markdown files aren't available

### 🧭 Navigation & Interaction
- **Fully keyboard navigable** - no mouse required
- **Arrow key navigation** between pages and within content
- **Enter to read** blog posts in full-screen mode
- **Backspace to return** to blog list
- **Smooth page transitions** and responsive controls

### 🗂️ Modular Content Structure
- **5 main sections**: Home, Projects, Blog, About, Contact
- **Easy to extend** - add new pages by updating the PageType enum
- **Separated content** from presentation logic
- **Dynamic project loading** from structured data

### 🌐 SSH Integration
- **Wish-powered SSH server** for remote access
- **Built-in SSH key generation** and management
- **Concurrent sessions** support
- **Production-ready** deployment capabilities

## 🛠️ Tech Stack

- **[Go](https://golang.org/)** - Primary programming language
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - TUI framework following Elm architecture
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Styling and layout
- **[Glamour](https://github.com/charmbracelet/glamour)** - Markdown rendering with syntax highlighting
- **[Bubbles](https://github.com/charmbracelet/bubbles)** - Common TUI components (viewport)
- **[Wish](https://github.com/charmbracelet/wish)** - SSH server framework

## 🚀 Quick Start

### Prerequisites
- Go 1.21+ installed
- Terminal with color support
- SSH client (for remote access)

### Installation

1. **Clone the repository:**
```bash
git clone https://github.com/Arpan-206/terminal-portfolio.git
cd terminal-portfolio
```

2. **Install dependencies:**
```bash
go mod tidy
```

3. **Build the application:**
```bash
go build -o portfolio .
```

4. **Run locally:**
```bash
./portfolio
```

5. **Access via SSH (default port 2222):**
```bash
ssh localhost -p 2222
```

## 📁 Project Structure

```
terminal-portfolio/
├── main.go                 # Application entry point
├── tui/                   # Terminal UI package
│   ├── model.go          # Main application model and views
│   ├── middleware.go     # SSH middleware setup
│   └── content.go        # Content loading and markdown processing
├── content/              # Content directory
│   └── blog/            # Blog posts in markdown format
│       ├── README.md    # Blog documentation
│       ├── *.md         # Individual blog posts
└── .ssh/                # SSH keys (generated on first run)
```

## 📝 Managing Blog Posts

### Creating a New Blog Post

1. **Create a new markdown file** in `content/blog/`:
```bash
touch content/blog/my-awesome-post.md
```

2. **Add frontmatter and content:**
```markdown
---
title: "My Awesome Blog Post"
summary: "A brief description of what this post is about."
date: "2024-01-15"
tags: ["go", "programming", "tutorial"]
readTime: "5 min read"
author: "Your Name"
published: true
---

# My Awesome Blog Post

Your markdown content goes here...

## Section 1
- Bullet points
- **Bold text**
- *Italic text*
- `Code snippets`

```go
// Code blocks with syntax highlighting
func main() {
    fmt.Println("Hello, World!")
}
```

More content...
```

3. **The post will automatically appear** in the blog section!

### Frontmatter Fields

| Field | Required | Description |
|-------|----------|-------------|
| `title` | ✅ | Post title displayed in cards and headers |
| `summary` | ✅ | Brief description shown in blog list |
| `date` | ✅ | Publication date (YYYY-MM-DD format) |
| `tags` | ❌ | Array of tags for categorization |
| `readTime` | ❌ | Estimated reading time |
| `author` | ❌ | Author name |
| `published` | ✅ | Set to `true` to make post visible |

## ⌨️ Keyboard Controls

### Global Navigation
- **← →** Navigate between main pages (Home, Projects, Blog, About, Contact)
- **q** or **Ctrl+C** Quit application

### Blog Section
- **↑ ↓** Navigate between blog posts
- **Enter** Open selected blog post
- **Backspace** Return to blog list from post view

### Within Blog Posts & Long Content
- **↑ ↓** Scroll line by line
- **Page Up/Down** Scroll by page
- **Home/End** Go to top/bottom

## 🎨 Customization

### Adding New Pages

1. **Add new page type** to `PageType` enum in `model.go`:
```go
const (
    HomePage PageType = iota
    ProjectsPage
    BlogPage
    AboutPage
    ContactPage
    NewPage  // Add your new page here
)
```

2. **Add page name** to pages slice in `NewModel()`:
```go
pages: []string{"Home", "Projects", "Blog", "About", "Contact", "New Page"},
```

3. **Add case** to `getPageContent()` method:
```go
case NewPage:
    return m.getNewPageContent()
```

4. **Implement the content method:**
```go
func (m Model) getNewPageContent() string {
    content := `# My New Page

Content goes here...`

    renderedContent := renderMarkdownForDisplay(content, m.width)
    return contentStyle.Render(renderedContent)
}
```

### Styling Customization

Modify the styles in `model.go`:
```go
var (
    navbarStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#FAFAFA")).
        Background(lipgloss.Color("#7D56F4")). // Change colors here
        // ... other style properties
)
```

## 🚀 Deployment

### Local Development
```bash
go run main.go
```

### Production Build
```bash
go build -ldflags="-s -w" -o portfolio .
```

### Docker Deployment
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o portfolio .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/portfolio .
COPY --from=builder /app/content ./content
EXPOSE 2222
CMD ["./portfolio"]
```

### SSH Server Configuration

The application automatically generates SSH keys on first run. For production:

1. **Use custom SSH keys:**
```bash
ssh-keygen -t ed25519 -f .ssh/id_ed25519
```

2. **Configure host and port** in `main.go`:
```go
const (
    host = "0.0.0.0"  // Change for production
    port = "2222"     // Change port as needed
)
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **[Charm](https://charm.sh/)** - For the amazing TUI framework and components
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - Elegant terminal UI framework
- **[Glamour](https://github.com/charmbracelet/glamour)** - Beautiful markdown rendering
- **Go Community** - For the powerful and simple programming language

## 🔗 Links

- **GitHub:** [https://github.com/Arpan-206](https://github.com/Arpan-206)
- **LinkedIn:** [https://www.linkedin.com/in/arpan-pandey/](https://www.linkedin.com/in/arpan-pandey/)
- **Bubble Tea:** [https://github.com/charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea)
- **Charm:** [https://charm.sh/](https://charm.sh/)

---

**Made with ❤️ using Go and Bubble Tea**
