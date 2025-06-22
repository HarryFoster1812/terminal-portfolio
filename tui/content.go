package tui

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
)

// FrontMatter represents the YAML frontmatter in markdown files
type FrontMatter struct {
	Title     string   `yaml:"title"`
	Summary   string   `yaml:"summary"`
	Date      string   `yaml:"date"`
	Tags      []string `yaml:"tags"`
	ReadTime  string   `yaml:"readTime"`
	Author    string   `yaml:"author"`
	Published bool     `yaml:"published"`
}

// BlogPost represents a blog post with metadata
type BlogPost struct {
	ID          string
	Title       string
	Summary     string
	Content     string
	Date        time.Time
	Tags        []string
	ReadTime    string
	Author      string
	Published   bool
	FilePath    string
}

// parseFrontMatter extracts YAML frontmatter and content from markdown
func parseFrontMatter(content string) (FrontMatter, string, error) {
	var fm FrontMatter
	
	if !strings.HasPrefix(content, "---\n") {
		return fm, content, nil
	}
	
	parts := strings.SplitN(content, "---\n", 3)
	if len(parts) < 3 {
		return fm, content, fmt.Errorf("invalid frontmatter format")
	}
	
	frontmatterContent := parts[1]
	markdownContent := parts[2]
	
	// Simple YAML parsing for our specific use case
	lines := strings.Split(frontmatterContent, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"`)
		
		switch key {
		case "title":
			fm.Title = value
		case "summary":
			fm.Summary = value
		case "date":
			fm.Date = value
		case "readTime":
			fm.ReadTime = value
		case "author":
			fm.Author = value
		case "published":
			fm.Published = value == "true"
		case "tags":
			// Parse tags array - simple implementation
			value = strings.Trim(value, "[]")
			if value != "" {
				tags := strings.Split(value, ",")
				for i, tag := range tags {
					tags[i] = strings.Trim(strings.TrimSpace(tag), `"`)
				}
				fm.Tags = tags
			}
		}
	}
	
	return fm, markdownContent, nil
}

// readMarkdownFile reads and parses a markdown file
func readMarkdownFile(filePath string) (BlogPost, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return BlogPost{}, err
	}
	
	fm, markdownContent, err := parseFrontMatter(string(content))
	if err != nil {
		return BlogPost{}, err
	}
	
	// Parse date
	date, err := time.Parse("2006-01-02", fm.Date)
	if err != nil {
		date = time.Now()
	}
	
	// Generate ID from filename
	filename := filepath.Base(filePath)
	id := strings.TrimSuffix(filename, filepath.Ext(filename))
	
	return BlogPost{
		ID:        id,
		Title:     fm.Title,
		Summary:   fm.Summary,
		Content:   markdownContent,
		Date:      date,
		Tags:      fm.Tags,
		ReadTime:  fm.ReadTime,
		Author:    fm.Author,
		Published: fm.Published,
		FilePath:  filePath,
	}, nil
}

// GetBlogPosts returns all available blog posts from markdown files
func GetBlogPosts() []BlogEntry {
	var posts []BlogPost
	
	// Get the current working directory
	pwd, err := os.Getwd()
	if err != nil {
		// Fallback to hardcoded posts if we can't read files
		return getFallbackBlogPosts()
	}
	
	blogDir := filepath.Join(pwd, "content", "blog")
	
	// Check if blog directory exists
	if _, err := os.Stat(blogDir); os.IsNotExist(err) {
		// Fallback to hardcoded posts if directory doesn't exist
		return getFallbackBlogPosts()
	}
	
	// Read all markdown files in the blog directory
	files, err := ioutil.ReadDir(blogDir)
	if err != nil {
		return getFallbackBlogPosts()
	}
	
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		
		if strings.HasSuffix(file.Name(), ".md") {
			filePath := filepath.Join(blogDir, file.Name())
			post, err := readMarkdownFile(filePath)
			if err != nil {
				continue // Skip files that can't be parsed
			}
			
			if post.Published {
				posts = append(posts, post)
			}
		}
	}
	
	// Sort posts by date (newest first)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})
	
	// Convert to BlogEntry format for compatibility
	var entries []BlogEntry
	for _, post := range posts {
		entries = append(entries, BlogEntry{
			Title:   post.Title,
			Summary: post.Summary,
			Content: post.Content,
			Date:    post.Date.Format("2006-01-02"),
		})
	}
	
	// If no posts were found, return fallback
	if len(entries) == 0 {
		return getFallbackBlogPosts()
	}
	
	return entries
}

// renderMarkdown renders markdown content using Glamour
func renderMarkdown(content string, width int) string {
	// Create a custom glamour renderer with a width that fits the terminal
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width-4), // Leave some padding
	)
	if err != nil {
		// Fallback to plain text if glamour fails
		return content
	}
	
	rendered, err := r.Render(content)
	if err != nil {
		// Fallback to plain text if rendering fails
		return content
	}
	
	return rendered
}

// renderMarkdownForDisplay renders markdown content optimized for display in the viewport
func renderMarkdownForDisplay(content string, width int) string {
	// Use a more conservative width for better readability
	displayWidth := width - 8 // Account for padding and borders
	if displayWidth < 40 {
		displayWidth = 40 // Minimum readable width
	}
	
	return renderMarkdown(content, displayWidth)
}

// getFallbackBlogPosts returns hardcoded blog posts as fallback
func getFallbackBlogPosts() []BlogEntry {
	fallbackContent1 := `# Building Terminal UIs with Bubble Tea

Bubble Tea is an amazing framework for building terminal user interfaces in Go. It follows the Elm architecture pattern, making it easy to reason about state management and updates.

## Key Features

- **Reactive programming model**: Based on the Elm Architecture
- **Built-in support**: Mouse and keyboard input handling
- **Flexible styling**: Integration with Lipgloss for beautiful UIs
- **Component-based**: Modular architecture for reusable components

## Getting Started

To get started with Bubble Tea, you need to understand three core concepts:

1. **Model**: Represents your application state
2. **Update**: Handles messages and updates the model  
3. **View**: Renders the model to the terminal

This is a fallback version of the blog post. Please check that your markdown files are properly configured in the ` + "`content/blog`" + ` directory.`

	fallbackContent2 := `# Modern Web Development Trends

The web development landscape continues to evolve rapidly with new frameworks and tools emerging regularly.

## Key Trends

- **React Server Components**: Zero bundle impact rendering
- **Edge Computing**: Bringing computation closer to users
- **TypeScript Everywhere**: Type safety across the full stack
- **AI Integration**: Copilot and AI-powered development tools

## Popular Frameworks

- **Next.js**: React with server-side rendering
- **Nuxt**: Vue.js full-stack framework
- **SvelteKit**: Svelte's full-stack solution
- **Remix**: Web standards focused React framework

This is a fallback version of the blog post. Please check that your markdown files are properly configured in the ` + "`content/blog`" + ` directory.`

	fallbackContent3 := `# Go vs Rust: A Developer's Perspective

Both Go and Rust have gained significant traction in the systems programming space.

## Go: Simplicity and Productivity

**Strengths:**
- Simple syntax and fast compilation
- Built-in concurrency with goroutines
- Excellent standard library
- Great for web services and microservices

## Rust: Safety and Performance  

**Strengths:**
- Memory safety without garbage collection
- Zero-cost abstractions
- Excellent performance
- Growing ecosystem with Cargo

## When to Choose What

**Choose Go** for rapid development and team productivity.
**Choose Rust** for maximum performance and memory safety.

This is a fallback version of the blog post. Please check that your markdown files are properly configured in the ` + "`content/blog`" + ` directory.`

	return []BlogEntry{
		{
			Title:   "Building Terminal UIs with Bubble Tea",
			Summary: "A deep dive into creating beautiful terminal applications using Charm's Bubble Tea framework.",
			Content: fallbackContent1,
			Date:    "2024-01-15",
		},
		{
			Title:   "Modern Web Development Trends",
			Summary: "Exploring the latest trends and technologies shaping web development in 2024.",
			Content: fallbackContent2,
			Date:    "2024-01-10",
		},
		{
			Title:   "Go vs Rust: A Developer's Perspective",
			Summary: "Comparing two modern systems programming languages from a practical standpoint.",
			Content: fallbackContent3,
			Date:    "2024-01-05",
		},
	}
}

// GetFeaturedProjects returns featured project information
func GetFeaturedProjects() []Project {
	return []Project{
		{
			Name:        "Terminal Portfolio",
			Description: "A beautiful terminal-based portfolio built with Bubble Tea",
			Tech:        []string{"Go", "Bubble Tea", "Lipgloss", "SSH"},
			Features:    []string{"Responsive design", "Keyboard navigation", "Smooth scrolling", "SSH server integration"},
			Status:      "Active",
			URL:         "https://github.com/Arpan-206/terminal-portfolio",
		},
		{
			Name:        "Full-Stack Web Applications",
			Description: "Modern web applications using cutting-edge frameworks",
			Tech:        []string{"React", "Next.js", "Node.js", "PostgreSQL"},
			Features:    []string{"Real-time updates", "Responsive UI", "RESTful APIs", "Authentication"},
			Status:      "Active",
			URL:         "https://github.com/Arpan-206",
		},
		{
			Name:        "DevOps Automation Tools",
			Description: "Scripts and tools for automating development workflows",
			Tech:        []string{"Python", "Bash", "Docker", "GitHub Actions"},
			Features:    []string{"CI/CD pipelines", "Automated testing", "Deployment scripts", "Infrastructure as Code"},
			Status:      "Active",
			URL:         "https://github.com/Arpan-206",
		},
		{
			Name:        "Machine Learning Projects",
			Description: "AI/ML applications solving real-world problems",
			Tech:        []string{"Python", "TensorFlow", "PyTorch", "scikit-learn"},
			Features:    []string{"Natural language processing", "Computer vision", "Data analysis", "Model deployment"},
			Status:      "Active",
			URL:         "https://github.com/Arpan-206",
		},
		{
			Name:        "Mobile Applications",
			Description: "Cross-platform mobile apps with native performance",
			Tech:        []string{"React Native", "Flutter", "Firebase", "SQLite"},
			Features:    []string{"Offline-first", "Real-time sync", "Push notifications", "Cross-platform UI"},
			Status:      "Active",
			URL:         "https://github.com/Arpan-206",
		},
		{
			Name:        "Open Source Contributions",
			Description: "Contributing to various open source projects in the community",
			Tech:        []string{"Go", "JavaScript", "Python", "Rust"},
			Features:    []string{"Bug fixes", "New features", "Documentation", "Community support"},
			Status:      "Ongoing",
			URL:         "https://github.com/Arpan-206",
		},
	}
}

// Project represents a project entry
type Project struct {
	Name        string
	Description string
	Tech        []string
	Features    []string
	Status      string
	URL         string
}