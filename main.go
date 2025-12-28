package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type Post struct {
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
	Image   string    `json:"image"`
	Slug    string    `json:"slug"`
	Content string    `json:"-"`
}

func parsePost(path string) (*Post, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	post := &Post{}
	contentStart := 0

	for i, line := range lines {
		if strings.HasPrefix(line, "title:") {
			post.Title = strings.TrimSpace(strings.TrimPrefix(line, "title:"))
		} else if strings.HasPrefix(line, "date:") {
			dateStr := strings.TrimSpace(strings.TrimPrefix(line, "date:"))
			post.Date, _ = time.Parse("2006-01-02", dateStr)
		} else if strings.HasPrefix(line, "image:") {
			post.Image = strings.TrimSpace(strings.TrimPrefix(line, "image:"))
		} else if strings.TrimSpace(line) == "---" && i > 0 {
			contentStart = i + 1
			break
		}
	}

	post.Content = strings.Join(lines[contentStart:], "\n")
	post.Slug = strings.TrimSuffix(filepath.Base(path), ".md")
	return post, nil
}

func getPosts() ([]*Post, error) {
	files, err := filepath.Glob("posts/*.md")
	if err != nil {
		return nil, err
	}

	var posts []*Post
	for _, file := range files {
		post, err := parsePost(file)
		if err != nil {
			continue
		}
		posts = append(posts, post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	return posts, nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := getPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/post/")
	post, err := parsePost("posts/" + slug + ".md")
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
	var buf strings.Builder
	if err := md.Convert([]byte(post.Content), &buf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>{{.Title}}</title>
	<link rel="stylesheet" href="/static/css/blog.css">
</head>
<body>
	<article>
		<header>
			<h1>{{.Title}}</h1>
			<time>{{.Date.Format "January 2, 2006"}}</time>
			{{if .Image}}<img src="/{{.Image}}" alt="{{.Title}}">{{end}}
		</header>
		<div class="content">{{.HTML}}</div>
		<a href="/" class="back">← Back to posts</a>
	</article>
</body>
</html>`

	t := template.Must(template.New("post").Parse(tmpl))
	data := struct {
		*Post
		HTML template.HTML
	}{post, template.HTML(buf.String())}
	t.Execute(w, data)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/posts", postsHandler)
	http.HandleFunc("/post/", postHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/posts/", http.StripPrefix("/posts/", http.FileServer(http.Dir("posts"))))

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
