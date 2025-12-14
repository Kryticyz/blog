# Minimalist Go Blog

Simple blog engine using Go and React.
Initial structure generated using Claude Code.

## Structure

```
.
├── main.go              # Go server
├── posts/               # Markdown blog posts
│   ├── *.md            # Blog post files
│   └── images/         # Blog images
├── static/
│   ├── index.html      # React home page
│   ├── js/app.js       # React component
│   └── css/            # Stylesheets
│       ├── home.css    # Home page styles
│       └── blog.css    # Blog post styles
└── go.mod

```

## Blog Post Format

Each `.md` file in `posts/` should have this front matter:

```markdown
title: Your Post Title
date: 2025-01-15
image: posts/images/blog_page/your-image.jpg
---

Your markdown content here...
```

## Running
You will have to install the goldmark dependency first.

```bash
go run main.go
```

Visit http://localhost:8080

## Customizing Styles

Edit `static/css/home.css` for the home page and `static/css/blog.css` for blog posts.
