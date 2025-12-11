---
applyTo: "**"
---

# DisGo Wiki - Copilot Instructions

## Project Overview

This is the official wiki for **DisGo** - a Discord bot library for the Go programming language. The wiki is built using:
- **Hugo** static site generator (v0.147.9+)
- **Hextra** theme for documentation
- Deployed via Netlify/GitHub Pages
- Content written in Markdown with Hugo front matter

## Key Technologies

- **Hugo**: Static site generator
- **Go Modules**: For theme management
- **Hextra Theme**: Documentation theme imported as Hugo module
- **Markdown**: Content format with YAML front matter
- **Git**: Version control

## Project Structure

```
/workspaces/disgo-wiki/
├── content/              # All wiki content (Markdown files)
│   ├── _index.md        # Homepage
│   └── guide/           # Main guide sections
│       ├── introduction.md
│       ├── getting-started/
│       └── disgo-bot/
├── static/              # Static assets (images, etc.)
├── layouts/             # Custom Hugo layouts/templates
├── assets/              # Asset pipeline files (CSS, JS)
├── public/              # Generated static site (DO NOT EDIT)
├── i18n/                # Internationalization files
├── hugo.yaml            # Main Hugo configuration
├── go.mod               # Go module dependencies
└── netlify.toml         # Netlify build configuration
```

## Content Guidelines

### Creating New Content

1. **File Location**: All content goes in `content/` directory
2. **File Format**: Use Markdown (`.md`) with YAML front matter
3. **Front Matter Requirements**:
   ```yaml
   ---
   title: Page Title
   linkTitle: Short Title  # For navigation
   weight: 1              # Lower numbers appear first
   next: path/to/next     # Optional navigation
   prev: path/to/prev     # Optional navigation
   breadcrumbs: true      # Default, set false to hide
   ---
   ```

### Writing Style

- **Audience**: Go developers learning to build Discord bots
- **Tone**: Friendly, educational, encouraging
- **Prerequisites**: Assume basic Go knowledge; link to resources for beginners
- **Code Examples**: Use Go syntax highlighting with \`\`\`go blocks
- **Links**: Use Hugo's `pageRef` for internal links or full URLs for external

### Content Organization

- Use `_index.md` for section landing pages
- Organize related content in subdirectories
- Use meaningful slugs (lowercase, hyphenated)
- Follow existing weight/ordering system

## Development Commands

### Local Development
```bash
# Start Hugo development server
hugo server -D --disableFastRender

# Build for production
hugo --gc --minify
```

### Testing
- Preview changes locally before committing
- Check for broken links
- Verify images load correctly
- Test navigation flow

## Configuration

### Main Config (`hugo.yaml`)
- Site title, copyright, and metadata
- Menu structure (navbar, footer)
- Hextra theme settings
- Markdown rendering options
- Banner messages for announcements

### Important Settings
- `enableGitInfo: true` - Shows last updated dates from Git
- `displayUpdatedDate: true` - Display update timestamps
- External links get decoration automatically

## Code Standards

### Markdown
- Use proper heading hierarchy (h1 for title, h2/h3 for sections)
- Use code fences with language identifiers
- Keep line length reasonable (~100 chars)
- Use relative links for internal content

### Front Matter
- Always include `title` and `weight`
- Use `linkTitle` for long titles in navigation
- Set `next`/`prev` to guide users through tutorials
- Use `breadcrumbs: false` only when necessary

### File Naming
- Use lowercase with hyphens: `project-setup.md`
- Match directory structure to URL structure
- Keep names descriptive but concise

## Common Tasks

### Adding a New Guide Page
1. Create file in appropriate `content/guide/` subdirectory
2. Add front matter with title and weight
3. Write content following style guidelines
4. Add to parent `_index.md` if needed
5. Test locally with `hugo server`

### Adding Images
1. Place in `static/images/` with descriptive subdirectory
2. Reference in Markdown: `![Alt text](/images/category/image.png)`
3. Use descriptive alt text
4. Optimize images before adding

### Updating Navigation
1. Edit `menu` section in `hugo.yaml`
2. Adjust `weight` to reorder items
3. Use `pageRef` for internal pages, `url` for external

## Important Notes

- **DO NOT EDIT** the `public/` directory - it's auto-generated
- The wiki is a work in progress (see banner in hugo.yaml)
- Refer to DisGo documentation at https://pkg.go.dev/github.com/disgoorg/disgo
- Discord community: https://discord.gg/9tKpqXjYVC
- GitHub repo: https://github.com/disgoorg/disgo

## Hextra Theme Features

- Built-in search functionality
- Dark/light theme toggle
- Responsive design
- Git-based update timestamps
- Breadcrumb navigation
- Next/previous page navigation
- Syntax highlighting
- External link decoration

## Deployment

- **Netlify**: Auto-deploys from main branch
- **GitHub Pages**: Via GitHub Actions workflow
- Build command: `hugo --gc --minify`
- Hugo version: 0.147.9 (configurable in netlify.toml)

## When Making Changes

1. Always test locally first
2. Check Hugo build succeeds without errors
3. Verify links work correctly
4. Ensure proper front matter formatting
5. Follow existing content structure
6. Keep tone consistent with existing content
7. Add helpful code examples for Go/DisGo concepts
8. Link to official DisGo docs when appropriate
