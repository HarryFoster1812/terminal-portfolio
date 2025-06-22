# Blog Posts

This directory contains blog posts in Markdown format with YAML frontmatter.

## File Structure

Each blog post should be a `.md` file with the following frontmatter structure:

```markdown
---
title: "Your Blog Post Title"
summary: "A brief summary of your blog post that will appear in the card view."
date: "2024-01-15"
tags: ["tag1", "tag2", "tag3"]
readTime: "5 min read"
author: "Arpan Pandey"
published: true
---

# Your Blog Post Title

Your markdown content goes here...

## Sections

You can use all standard markdown features:

- Lists
- **Bold text**
- *Italic text*
- `Code snippets`

```go
// Code blocks with syntax highlighting
func main() {
    fmt.Println("Hello, World!")
}
```

## More content...
```

## Frontmatter Fields

- `title`: The title of your blog post (required)
- `summary`: A brief description shown in the blog list (required)
- `date`: Publication date in YYYY-MM-DD format (required)
- `tags`: Array of tags for categorization (optional)
- `readTime`: Estimated reading time (optional)
- `author`: Author name (optional)
- `published`: Set to `true` to make the post visible (required)

## File Naming

Use kebab-case for filenames:
- `my-awesome-blog-post.md`
- `building-terminal-uis.md`
- `web-development-trends-2024.md`

## Adding New Posts

1. Create a new `.md` file in this directory
2. Add the frontmatter with all required fields
3. Write your content in Markdown
4. Set `published: true` when ready to publish
5. The application will automatically load and display your post

Posts are automatically sorted by date (newest first) in the blog interface.