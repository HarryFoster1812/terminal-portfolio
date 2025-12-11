package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// PageType represents different page types
type PageType int

const (
	HomePage PageType = iota
	ProjectsPage
	BlogPage
	AboutPage
	ContactPage
)

// BlogEntry represents a blog post
type BlogEntry struct {
	Title   string
	Summary string
	Content string
	Date    string
}

// Model represents the terminal UI state
type Model struct {
	width             int
	height            int
	currentPage       PageType
	pages             []string
	blogEntries       []BlogEntry
	selectedBlogEntry int
	viewingBlogEntry  bool
	viewport          viewport.Model
	ready             bool
	lastKey 		  string
}

// Styles
var (
	navbarStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		PaddingTop(1).
		PaddingBottom(1).
		PaddingLeft(2).
		PaddingRight(2).
		Width(100).
		MaxWidth(200)

	activeNavStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Background(lipgloss.Color("#FAFAFA")).
		PaddingLeft(1).
		PaddingRight(1)

	inactiveNavStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		PaddingLeft(1).
		PaddingRight(1)

	contentStyle = lipgloss.NewStyle().
		Padding(1, 2)

	footerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Background(lipgloss.Color("#1a1a1a")).
		PaddingLeft(2).
		PaddingRight(2).
		Width(100).
		MaxWidth(200)

	cardStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1).
		Margin(1, 0).
		Width(70)

	selectedCardStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#F25D94")).
		Background(lipgloss.Color("#2a2a2a")).
		Padding(1).
		Margin(1, 0).
		Width(70)
)

// NewModel creates a new Model instance
func NewModel(width, height int) Model {
	blogEntries := GetBlogPosts()

	vp := viewport.New(width, height-4) // Reserve space for navbar and footer
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.HiddenBorder()).
		PaddingLeft(0).
		PaddingRight(0)

	return Model{
		width:             width,
		height:            height,
		currentPage:       HomePage,
		pages:             []string{"Home", "Projects", "Blog", "About", "Contact"},
		blogEntries:       blogEntries,
		selectedBlogEntry: 0,
		viewingBlogEntry:  false,
		viewport:          vp,
		ready:             false,
		lastKey: 			"",
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles model updates
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-4)
			m.viewport.Style = lipgloss.NewStyle().
				BorderStyle(lipgloss.HiddenBorder()).
				PaddingLeft(0).
				PaddingRight(0)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 4
		}
		
		// Update navbar and footer styles with new width
		navbarStyle = navbarStyle.Width(msg.Width)
		footerStyle = footerStyle.Width(msg.Width)
		
		m.width = msg.Width
		m.height = msg.Height
		
		// Update viewport content
		m.updateViewportContent()

	case tea.KeyMsg:
		key := msg.String()

		// Detect capital letters (Shift+key)
		var shifted rune
		if len(msg.Runes) == 1 {
			
			shifted = msg.Runes[0]
			fmt.Println(shifted)
		}

		if key == "ctrl+c" || key == "q" {
			return m, tea.Quit
		}

		switch key {
			case "left", "h":
				if !m.viewingBlogEntry {
					if m.currentPage > 0 {
						m.currentPage--
						m.updateViewportContent()
					}
				}

			case "right", "l":
				if !m.viewingBlogEntry {
					if int(m.currentPage) < len(m.pages)-1 {
						m.currentPage++
						m.updateViewportContent()
					}
				}
		}


		if shifted == 'G' {
			// Shift+g = G = jump bottom
			if m.currentPage == BlogPage && !m.viewingBlogEntry {
				m.selectedBlogEntry = len(m.blogEntries) - 1
			}

        m.viewport.GotoBottom()

		content := m.getPageContent()
		m.viewport.SetContent(content)
        return m, nil
		}

		if key == "g" {
			if m.lastKey == "g" {
				// GG detected
				m.viewport.GotoTop()
				if m.currentPage == BlogPage && !m.viewingBlogEntry {
					m.selectedBlogEntry = 0
				}
				m.updateViewportContent()
				m.lastKey = ""
				return m, nil
			}

			// Single g: wait for next key
			m.lastKey = "g"
			return m, nil
		}

		// Any other key cancels "g" state
		m.lastKey = ""
	
		switch key {
			case "ctrl+u":
				m.viewport.LineUp(10)
				return m, nil

			case "ctrl+d":
				m.viewport.LineDown(10)
				return m, nil

			case "up", "k":
				if m.currentPage == BlogPage && !m.viewingBlogEntry {
					if m.selectedBlogEntry > 0 {
						m.selectedBlogEntry--
						m.updateViewportContent()
					}
				} else if m.viewingBlogEntry || (m.currentPage != BlogPage) {
					m.viewport, cmd = m.viewport.Update(msg)
					return m, cmd
				}
				return m, nil

			case "down", "j":
				if m.currentPage == BlogPage && !m.viewingBlogEntry {
					if m.selectedBlogEntry < len(m.blogEntries)-1 {
						m.selectedBlogEntry++
						m.updateViewportContent()
					}
				} else if m.viewingBlogEntry || (m.currentPage != BlogPage) {
					m.viewport, cmd = m.viewport.Update(msg)
					return m, cmd
				}
				return m, nil


			case "enter":
				if m.currentPage == BlogPage && !m.viewingBlogEntry {
					m.viewingBlogEntry = true
					m.updateViewportContent()
				}
				return m, nil

			case "backspace":
				if m.viewingBlogEntry {
					m.viewingBlogEntry = false
					m.updateViewportContent()
				}			
				return m, nil

		}
	}

	// Update viewport
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// updateViewportContent updates the viewport content based on current state
func (m *Model) updateViewportContent() {
	content := m.getPageContent()
	m.viewport.SetContent(content)
	m.viewport.GotoTop()
}

// View renders the model
func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	// Build navbar - always visible at top
	navbar := m.renderNavbar()
	
	// Build footer - always visible at bottom
	footer := m.renderFooter()

	// Viewport handles the scrollable content
	view := lipgloss.JoinVertical(
		lipgloss.Left,
		navbar,
		m.viewport.View(),
		footer,
	)

	return view
}

func (m Model) renderNavbar() string {
	var navItems []string
	
	for i, page := range m.pages {
		if PageType(i) == m.currentPage {
			navItems = append(navItems, activeNavStyle.Render(page))
		} else {
			navItems = append(navItems, inactiveNavStyle.Render(page))
		}
	}

	navbar := strings.Join(navItems, " ")
	return navbarStyle.Render("📍 Arpan's Portfolio  |  " + navbar)
}

func (m Model) getPageContent() string {
	switch m.currentPage {
	case HomePage:
		return m.getHomeContent()
	case ProjectsPage:
		return m.getProjectsContent()
	case BlogPage:
		if m.viewingBlogEntry {
			return m.getBlogEntryContent()
		}
		return m.getBlogContent()
	case AboutPage:
		return m.getAboutContent()
	case ContactPage:
		return m.getContactContent()
	default:
		return "Page not found"
	}
}

func (m Model) getHomeContent() string {
	markdownContent := `# 🚀 Welcome to Arpan's Terminal Portfolio!

Hi there! I'm **Arpan Pandey**, a passionate tech enthusiast and developer.

## 🔗 Connect with me

- **GitHub:** [https://github.com/Arpan-206](https://github.com/Arpan-206)
- **LinkedIn:** [https://www.linkedin.com/in/arpan-pandey/](https://www.linkedin.com/in/arpan-pandey/)

✨ Navigate using arrow keys to explore my projects and blog posts!

💡 This portfolio is built with **Go**, **Bubble Tea**, and lots of ❤️

## 🎯 What you'll find here

- My latest projects and open source contributions
- Technical blog posts and tutorials
- Information about my skills and experience  
- Ways to get in touch and collaborate

## 🌟 Features of this terminal portfolio

- **Fully keyboard navigable** interface
- **Responsive design** that adapts to your terminal size
- **Beautiful styling** with Lipgloss and Glamour
- **Smooth scrolling** for long content
- **Interactive blog post viewer** with markdown rendering

## ⚡ Quick navigation tips

- Use **← →** arrows to switch between main sections
- Use **↑ ↓** arrows to navigate within sections
- Press **Enter** to open blog posts
- Press **Backspace** to go back from blog posts
- Press **q** or **Ctrl+C** to quit

---

**Happy exploring!** 🎉`

	// Render the markdown content using Glamour
	renderedContent := renderMarkdownForDisplay(markdownContent, m.width)
	
	return contentStyle.Render(renderedContent)
}

func (m Model) getProjectsContent() string {
	projects := GetFeaturedProjects()
	
	var markdownContent strings.Builder
	markdownContent.WriteString("# 🛠️ Featured Projects\n\n")
	markdownContent.WriteString("Here are some of my notable projects and contributions:\n\n")
	
	for i, project := range projects {
		markdownContent.WriteString(fmt.Sprintf("## %d. %s\n\n", i+1, project.Name))
		markdownContent.WriteString(fmt.Sprintf("**Description:** %s\n\n", project.Description))
		
		// Tech stack
		markdownContent.WriteString("**Technologies:**\n")
		for _, tech := range project.Tech {
			markdownContent.WriteString(fmt.Sprintf("- %s\n", tech))
		}
		markdownContent.WriteString("\n")
		
		// Features
		if len(project.Features) > 0 {
			markdownContent.WriteString("**Key Features:**\n")
			for _, feature := range project.Features {
				markdownContent.WriteString(fmt.Sprintf("- %s\n", feature))
			}
			markdownContent.WriteString("\n")
		}
		
		// Status and URL
		markdownContent.WriteString(fmt.Sprintf("**Status:** %s\n", project.Status))
		if project.URL != "" {
			markdownContent.WriteString(fmt.Sprintf("**Repository:** [%s](%s)\n", project.URL, project.URL))
		}
		
		markdownContent.WriteString("\n---\n\n")
	}
	
	markdownContent.WriteString("## 🔗 Links\n\n")
	markdownContent.WriteString("Check out my GitHub profile for complete project listings, source code, and detailed documentation:\n\n")
	markdownContent.WriteString("**GitHub:** [https://github.com/Arpan-206](https://github.com/Arpan-206)\n\n")
	markdownContent.WriteString("🚀 Always working on something new and exciting!\n")
	
	// Render the markdown content using Glamour
	renderedContent := renderMarkdownForDisplay(markdownContent.String(), m.width)
	
	return contentStyle.Render(renderedContent)
}

func (m Model) getBlogContent() string {
	var cards []string
	
	for i, entry := range m.blogEntries {
		cardContent := fmt.Sprintf("📝 %s\n\n%s\n\n📅 %s", entry.Title, entry.Summary, entry.Date)
		
		if i == m.selectedBlogEntry {
			cards = append(cards, selectedCardStyle.Render(cardContent))
		} else {
			cards = append(cards, cardStyle.Render(cardContent))
		}
	}
	
	header := contentStyle.Render("📚 Blog Posts\n\nUse ↑/↓ to navigate posts, Enter to read, scroll within posts\n")
	content := header + "\n" + strings.Join(cards, "\n")
	
	return content
}

func (m Model) getBlogEntryContent() string {
	if m.selectedBlogEntry >= len(m.blogEntries) {
		return contentStyle.Render("Blog entry not found")
	}
	
	entry := m.blogEntries[m.selectedBlogEntry]
	
	// Render the markdown content using Glamour
	renderedContent := renderMarkdownForDisplay(entry.Content, m.width)
	
	// Create header with title and date
	header := fmt.Sprintf("📝 %s\n📅 %s\n\n", entry.Title, entry.Date)
	
	return contentStyle.Render(header + renderedContent)
}

func (m Model) getAboutContent() string {
	markdownContent := `# 👋 About Arpan

I'm a passionate software developer and tech enthusiast with a love for creating beautiful, functional applications. My journey in technology spans across various domains, always driven by curiosity and the desire to solve complex problems.

## 🎓 Background & Education
- **Computer Science and Engineering**
- **Full-stack development** experience across multiple technologies
- **Continuous learner**, always exploring new technologies and methodologies
- **Active contributor** to open source projects and tech communities

## 💻 Technical Expertise

### 🌐 Frontend Development
- **React, Next.js, Vue.js** - Modern JavaScript frameworks
- **TypeScript** for type-safe development
- **HTML5, CSS3, Sass/SCSS** for styling
- **Responsive design** and accessibility best practices
- **State management** with Redux, Zustand, Context API

### ⚙️ Backend Development
- **Go** - Systems programming and web services
- **Node.js, Express.js** - JavaScript backend development
- **Python, Django, FastAPI** - Rapid development and data processing
- **RESTful APIs and GraphQL**
- **Microservices architecture** and distributed systems

### 🗄️ Database Technologies
- **PostgreSQL, MySQL** - Relational databases
- **MongoDB, Redis** - NoSQL solutions
- **Database design** and optimization
- **Data modeling** and migration strategies

### ☁️ Cloud & DevOps
- **AWS, Google Cloud Platform** - Cloud infrastructure
- **Docker, Kubernetes** - Containerization and orchestration
- **CI/CD pipelines** with GitHub Actions, GitLab CI
- **Infrastructure as Code** with Terraform
- **Monitoring and logging** solutions

### 🛠️ Development Tools & Practices
- **Git version control** and collaborative workflows
- **Test-driven development (TDD)** and automated testing
- **Code review processes** and pair programming
- **Agile methodologies** and project management
- **Performance optimization** and security best practices

## 🌟 Philosophy & Approach

I believe in writing clean, maintainable code that not only solves problems but is also a joy to work with. Every project, whether it's a complex enterprise application or a simple CLI tool, deserves attention to detail and thoughtful architecture.

**My approach emphasizes:**
- User-centered design and experience
- Scalable and maintainable code architecture
- Collaborative development and knowledge sharing
- Continuous learning and adaptation to new technologies
- Open source contribution and community building

## 🚀 Current Focus & Interests
- Building developer tools that improve productivity
- Exploring systems programming with Go and Rust
- Contributing to open source projects
- Terminal applications and command-line interfaces
- Modern web technologies and frameworks
- Machine learning and AI applications
- Mentoring junior developers and sharing knowledge

## 🎯 Goals & Aspirations
- Create impactful software that solves real-world problems
- Build and maintain high-quality open source projects
- Foster inclusive and collaborative development communities
- Continue learning and staying current with technology trends
- Share knowledge through writing, speaking, and mentoring

---

When I'm not coding, you might find me exploring new technologies, contributing to open source projects, writing technical blogs, or engaging with the developer community. I'm always excited to learn something new and share that knowledge with others.

**Let's build something amazing together!** 🎉`

	// Render the markdown content using Glamour
	renderedContent := renderMarkdownForDisplay(markdownContent, m.width)
	
	return contentStyle.Render(renderedContent)
}

func (m Model) getContactContent() string {
	markdownContent := `# 📬 Get In Touch

I'm always excited to connect with fellow developers, potential collaborators, or anyone interested in technology! Whether you want to discuss a project, share ideas, or just have a friendly chat about development, I'd love to hear from you.

## 🔗 Find Me Online

### 🐙 GitHub
**[https://github.com/Arpan-206](https://github.com/Arpan-206)**
- Check out my repositories and contributions
- See my latest projects and code samples
- Contribute to open source projects together
- Star repositories you find interesting!

### 💼 LinkedIn
**[https://www.linkedin.com/in/arpan-pandey/](https://www.linkedin.com/in/arpan-pandey/)**
- Professional background and experience
- Connect for networking and opportunities
- Endorse skills and get recommendations
- Stay updated with my professional journey

### 📧 Email
**Best reached via LinkedIn**  
For direct communication, please connect with me on LinkedIn first. I'm responsive and check messages regularly!

---

## 💼 Professional Opportunities

### 🤝 Open to Collaboration
- Open source project contributions
- Technical writing and documentation
- Code reviews and pair programming sessions
- Speaking at tech events and conferences
- Mentoring and knowledge sharing

### 💻 Freelance & Contract Work
- Full-stack web application development
- API design and backend services
- Terminal applications and CLI tools
- Code audits and technical consulting
- DevOps and infrastructure setup

### 🏢 Full-time Positions
- Software Engineer / Senior Software Engineer
- Full-stack Developer positions
- Backend/Systems Engineer roles
- DevOps Engineer opportunities
- Technical Lead positions

---

## 🎯 Areas of Interest

### 🛠️ Technology Domains
- Go, Rust, and systems programming
- Modern JavaScript/TypeScript ecosystems
- Cloud-native applications and microservices
- Developer tooling and CLI applications
- Database design and optimization
- API development and integration

### 🌍 Industry Sectors
- Developer tools and productivity software
- Financial technology (FinTech)
- Healthcare technology solutions
- Educational technology platforms
- Open source and community-driven projects
- Startups and innovative tech companies

### 💡 Project Types
- Greenfield projects with modern tech stacks
- Legacy system modernization and migration
- Performance optimization and scalability improvements
- Integration projects and API development
- Automation and workflow improvement tools

---

## 🤔 What I'm Looking For

### 🎯 In Collaborations
- Passionate and skilled team members
- Projects that make a positive impact
- Opportunities to learn and grow
- Respectful and inclusive work environments
- Clear communication and shared goals

### 💪 In Roles
- Challenging technical problems to solve
- Opportunities for professional growth
- Mentorship and knowledge sharing culture
- Work-life balance and flexibility
- Competitive compensation and benefits

---

## 📅 Let's Connect!

Whether you're interested in:
- Discussing potential collaborations
- Exploring job opportunities
- Getting technical advice or mentorship
- Sharing ideas about technology and development
- Just having a friendly chat about coding

I'm always happy to connect! The best way to reach me is through LinkedIn, where I'm active and responsive. Let's build something amazing together!

🚀 **Looking forward to hearing from you!** 🎉

---

**P.S.** If you enjoyed this terminal portfolio, feel free to star it on GitHub or share it with others who might appreciate terminal-based applications. Your support means a lot! ⭐`

	// Render the markdown content using Glamour
	renderedContent := renderMarkdownForDisplay(markdownContent, m.width)
	
	return contentStyle.Render(renderedContent)
}

func (m Model) renderFooter() string {
	var helpText string
	
	if m.viewingBlogEntry {
		helpText = "📖 Reading blog post • ↑/↓ PgUp/PgDn to scroll • Backspace to return • q/Ctrl+C to quit"
	} else {
		switch m.currentPage {
		case BlogPage:
			helpText = "📚 Blog posts • ↑/↓ navigate • Enter to read • ←/→ change page • q/Ctrl+C to quit"
		default:
			helpText = "🧭 Portfolio navigation • ←/→ navigate pages • ↑/↓ scroll content • q/Ctrl+C to quit"
		}
	}
	
	return footerStyle.Render(helpText)
}
