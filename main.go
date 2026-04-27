package main

import (
	"bufio"
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/jmwai80/personal-website/handlers"
	"github.com/jmwai80/personal-website/middleware"
	"github.com/jmwai80/personal-website/models"
)

//go:embed templates
var templateFiles embed.FS

var tmpl = template.Must(
	template.New("").Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"not": func(v interface{}) bool {
			if v == nil {
				return true
			}
			if s, ok := v.([]*models.Post); ok {
				return len(s) == 0
			}
			return false
		},
	}).ParseFS(templateFiles, "templates/*.html"),
)

func init() {
	for _, path := range []string{"/app/.env", ".env"} {
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			if key != "" && os.Getenv(key) == "" {
				os.Setenv(key, val)
			}
		}
		break
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// Public routes
	r.Get("/", homeHandler)
	r.Get("/health", healthHandler)
	r.Get("/blog", blogHandler)
	r.Get("/blog/posts/{slug}", postHandler)

	// Auth routes (public)
	r.Get("/admin/login", handlers.LoginForm)
	r.With(middleware.RateLimit(10, 1*time.Minute)).Post("/admin/login", handlers.LoginHandler)
	r.Post("/admin/logout", handlers.LogoutHandler)

	// Protected admin routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAdmin)
		r.Use(middleware.CSRFProtect)

		r.Get("/admin", handlers.AdminIndex)
		r.Get("/admin/new", handlers.NewPostForm)
		r.Post("/admin/new", handlers.CreatePost)
		r.Get("/admin/edit/{slug}", handlers.EditPostForm)
		r.Post("/admin/edit/{slug}", handlers.UpdatePost)
		r.Post("/admin/delete/{slug}", handlers.DeletePost)
	})

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.ExecuteTemplate(w, "base", nil); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("template error: %v", err)
	}
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := models.ListPosts()
	if err != nil {
		log.Printf("ListPosts error: %v", err)
		posts = []*models.Post{}
	}
	// Collect unique tags across all posts
	seen := map[string]bool{}
	var tags []string
	for _, p := range posts {
		for _, t := range p.Tags {
			if !seen[t] {
				seen[t] = true
				tags = append(tags, t)
			}
		}
	}
	data := struct {
		Posts []*models.Post
		Tags  []string
	}{Posts: posts, Tags: tags}
	if err := tmpl.ExecuteTemplate(w, "blog", data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("blog template error: %v", err)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	post, err := models.LoadPost(slug)
	if err != nil || post.Draft {
		http.NotFound(w, r)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "post", post); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("post template error: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
