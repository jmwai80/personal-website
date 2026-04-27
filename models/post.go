package models

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

type Post struct {
	Slug        string
	Title       string
	Date        time.Time
	Tags        []string
	Draft       bool
	Description string
	ContentHTML template.HTML
}

var slugPattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

type frontmatter struct {
	Title       string   `yaml:"title"`
	Date        string   `yaml:"date"`
	Tags        []string `yaml:"tags"`
	Draft       bool     `yaml:"draft"`
	Description string   `yaml:"description"`
}

func LoadPost(slug string) (*Post, error) {
	if !slugPattern.MatchString(slug) {
		return nil, fmt.Errorf("invalid slug: %s", slug)
	}
	path := filepath.Join("content", "posts", slug+".md")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	parts := strings.SplitN(string(data), "---", 3)
	if len(parts) < 3 || strings.TrimSpace(parts[0]) != "" {
		return nil, fmt.Errorf("invalid frontmatter in %s", path)
	}

	var fm frontmatter
	if err := yaml.Unmarshal([]byte(parts[1]), &fm); err != nil {
		return nil, fmt.Errorf("parsing frontmatter: %w", err)
	}

	date, err := time.Parse("2006-01-02", fm.Date)
	if err != nil {
		return nil, fmt.Errorf("parsing date: %w", err)
	}

	var buf strings.Builder
	if err := goldmark.Convert([]byte(strings.TrimSpace(parts[2])), &buf); err != nil {
		return nil, fmt.Errorf("rendering markdown: %w", err)
	}

	return &Post{
		Slug:        slug,
		Title:       fm.Title,
		Date:        date,
		Tags:        fm.Tags,
		Draft:       fm.Draft,
		Description: fm.Description,
		ContentHTML: template.HTML(buf.String()),
	}, nil
}

func SavePost(slug, title, description string, tags []string, draft bool, content string) error {
	slug = strings.ToLower(slug)
	if !slugPattern.MatchString(slug) {
		return fmt.Errorf("invalid slug: %s", slug)
	}
	base, _ := filepath.Abs("content/posts")
	fullPath := filepath.Clean(filepath.Join(base, slug+".md"))
	rel, err := filepath.Rel(base, fullPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return fmt.Errorf("path escapes base directory")
	}

	type fm struct {
		Title       string   `yaml:"title"`
		Date        string   `yaml:"date"`
		Tags        []string `yaml:"tags"`
		Draft       bool     `yaml:"draft"`
		Description string   `yaml:"description"`
	}
	// Preserve existing date if file exists
	date := time.Now().Format("2006-01-02")
	if existing, err := LoadPost(slug); err == nil {
		date = existing.Date.Format("2006-01-02")
	}

	front, err := yaml.Marshal(fm{
		Title:       title,
		Date:        date,
		Tags:        tags,
		Draft:       draft,
		Description: description,
	})
	if err != nil {
		return fmt.Errorf("marshaling frontmatter: %w", err)
	}

	out := "---\n" + string(front) + "---\n\n" + strings.TrimSpace(content) + "\n"

	tmpPath := fullPath + ".tmp"
	if err := os.WriteFile(tmpPath, []byte(out), 0644); err != nil {
		return fmt.Errorf("writing tmp file: %w", err)
	}
	if err := os.Rename(tmpPath, fullPath); err != nil {
		return fmt.Errorf("renaming tmp file: %w", err)
	}
	return nil
}

func DeletePost(slug string) error {
	if !slugPattern.MatchString(slug) {
		return fmt.Errorf("invalid slug: %s", slug)
	}
	base, _ := filepath.Abs("content/posts")
	fullPath := filepath.Clean(filepath.Join(base, slug+".md"))
	rel, err := filepath.Rel(base, fullPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return fmt.Errorf("path escapes base directory")
	}
	return os.Remove(fullPath)
}

func listPosts(includeDrafts bool) ([]*Post, error) {
	entries, err := os.ReadDir("content/posts")
	if err != nil {
		return nil, fmt.Errorf("reading posts dir: %w", err)
	}
	var posts []*Post
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		slug := strings.TrimSuffix(e.Name(), ".md")
		post, err := LoadPost(slug)
		if err != nil {
			log.Printf("skipping %s: %v", e.Name(), err)
			continue
		}
		if includeDrafts || !post.Draft {
			posts = append(posts, post)
		}
	}
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})
	return posts, nil
}

// ListPosts returns published posts only, sorted by date descending.
func ListPosts() ([]*Post, error) { return listPosts(false) }

// ListAllPosts returns all posts including drafts (for admin).
func ListAllPosts() ([]*Post, error) { return listPosts(true) }
