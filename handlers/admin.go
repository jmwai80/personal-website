package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jmwai80/personal-website/models"
)

var nonAlpha = regexp.MustCompile(`[^a-z0-9]+`)

var sharedHead = `
<meta charset="utf-8"/>
<meta name="viewport" content="width=device-width,initial-scale=1"/>
<script src="https://cdn.tailwindcss.com"></script>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;600&display=swap" rel="stylesheet">
<style>
  body{background:#0f041f;font-family:"JetBrains Mono",monospace;color:#e8d6f0;}
  :root{--accent:#ff5cd6;}
  input,textarea,select{background:#0a0318;border:1px solid #2a1040;color:#e8d6f0;padding:.5rem .75rem;width:100%;font-family:inherit;font-size:.875rem;}
  input:focus,textarea:focus{outline:none;border-color:#ff5cd6;}
  .btn{display:inline-flex;align-items:center;gap:.5rem;padding:.5rem 1.25rem;font-size:.875rem;font-weight:600;cursor:pointer;transition:opacity .2s;border:none;}
  .btn-primary{background:#ff5cd6;color:#0f041f;}
  .btn-primary:hover{opacity:.85;}
  .btn-ghost{background:transparent;color:#e8d6f0;border:1px solid #2a1040;}
  .btn-ghost:hover{border-color:#ff5cd6;color:#ff5cd6;}
  .btn-danger{background:transparent;color:#f87171;border:1px solid rgba(248,113,113,.25);}
  .btn-danger:hover{background:rgba(248,113,113,.08);border-color:#f87171;}
  ::selection{background:rgba(255,92,214,.25);}
</style>`

var adminNav = `
<header class="border-b border-[#1e0f35] px-6 py-4 flex items-center justify-between">
  <a href="/admin" class="text-[#ff5cd6] font-semibold tracking-tight text-sm">{/} admin</a>
  <div class="flex items-center gap-4 text-xs text-[#6b5a7e]">
    <a href="/" class="hover:text-[#e8d6f0] transition-colors">← site</a>
    <form method="POST" action="/admin/logout" style="display:inline">
      <input type="hidden" name="csrf_token" value="%s"/>
      <button type="submit" class="hover:text-[#e8d6f0] transition-colors cursor-pointer">logout</button>
    </form>
  </div>
</header>`

// AdminIndex lists all posts
func AdminIndex(w http.ResponseWriter, r *http.Request) {
	posts, _ := models.ListAllPosts()
	csrf := randomHex(16)
	http.SetCookie(w, &http.Cookie{
		Name: "csrf_token", Value: csrf, Path: "/",
		SameSite: http.SameSiteStrictMode, MaxAge: 3600,
	})

	var rows strings.Builder
	for _, p := range posts {
		draft := ""
		if p.Draft {
			draft = `<span class="text-[10px] border border-yellow-700 text-yellow-500 px-1.5 py-0.5 ml-2 align-middle">draft</span>`
		}
		rows.WriteString(fmt.Sprintf(`
	<tr class="border-b border-[#1e0f35] hover:bg-[#140828] transition-colors">
	  <td class="px-5 py-4">
	    <div class="text-sm text-fg font-medium">%s %s</div>
	    <div class="text-[11px] text-mute mt-1">%s</div>
	  </td>
	  <td class="px-5 py-4 text-xs text-mute">%s</td>
	  <td class="px-5 py-4">
	    <div class="flex items-center gap-3">
	      <a href="/admin/edit/%s" class="btn-primary text-[11px] px-4 py-1.5 rounded" style="background:#ff5cd6;color:#0f041f;display:inline-flex;align-items:center;gap:.35rem;text-decoration:none;">✎ edit</a>
	      <form method="POST" action="/admin/delete/%s" onsubmit="return confirm('Delete «%s»?')" style="display:inline">
	        <input type="hidden" name="csrf_token" value="%s"/>
	        <button type="submit" class="btn-danger text-[11px] px-4 py-1.5 rounded" style="background:transparent;color:#f87171;border:1px solid rgba(248,113,113,.25);cursor:pointer;display:inline-flex;align-items:center;gap:.35rem;">✕ delete</button>
	      </form>
	    </div>
	  </td>
	</tr>`, template.HTMLEscapeString(p.Title), draft,
			template.HTMLEscapeString(p.Slug),
			p.Date.Format("2006-01-02"),
			template.HTMLEscapeString(p.Slug),
			template.HTMLEscapeString(p.Slug),
			template.HTMLEscapeString(p.Title),
			csrf))
	}

	fmt.Fprintf(w, `<!doctype html><html><head><title>posts · admin</title>%s</head>
	<body class="min-h-screen">
	%s
	<main class="max-w-5xl mx-auto px-6 py-10">
	  <div class="flex items-center justify-between mb-8">
	    <div>
	      <p class="text-[#ff5cd6] text-xs mb-1 tracking-wider">$ ls content/posts/</p>
	      <h1 class="text-2xl font-semibold tracking-tight">posts <span class="text-mute text-base font-normal">(%d total)</span></h1>
	    </div>
	    <a href="/admin/new" class="btn-primary text-sm px-5 py-2.5 rounded flex items-center gap-2" style="background:#ff5cd6;color:#0f041f;text-decoration:none;">+ new post</a>
	  </div>
	  <div class="border border-[#1e0f35] rounded overflow-hidden">
	    <table class="w-full text-left">
	      <thead>
	        <tr class="bg-[#140828] border-b border-[#1e0f35]">
	          <th class="px-5 py-3.5 text-xs text-mute font-medium tracking-wider uppercase">post</th>
	          <th class="px-5 py-3.5 text-xs text-mute font-medium tracking-wider uppercase">date</th>
	          <th class="px-5 py-3.5 text-xs text-mute font-medium tracking-wider uppercase">actions</th>
	        </tr>
	      </thead>
	      <tbody>%s</tbody>
	    </table>
	  </div>
	  <p class="mt-6 text-mute text-[11px]">
	    <span style="color:var(--accent);">$</span> find content/posts/ -name "*.md" | wc -l &nbsp;<span style="color:#1e0f35;">→</span>&nbsp; %d files
	  </p>
	</main>
	</body></html>`, sharedHead, fmt.Sprintf(adminNav, csrf), len(posts), rows.String(), len(posts))
}

var postFormHTML = `<!doctype html><html><head><title>%s — admin</title>%s
<script src="https://cdn.jsdelivr.net/npm/marked/marked.min.js"></script>
</head>
<body class="min-h-screen">
%s
<main class="max-w-6xl mx-auto px-6 py-10">
  <div class="mb-6">
    <p class="text-[#ff5cd6] text-xs mb-1 tracking-wider">$ vim content/posts/%s.md</p>
    <h1 class="text-2xl font-semibold tracking-tight">%s</h1>
  </div>
  %s
  <form method="POST" action="%s" id="postForm">
    <input type="hidden" name="csrf_token" value="%s"/>
    <div class="grid grid-cols-2 gap-6 mb-4">
      <div>
        <label class="block text-xs text-mute mb-1 tracking-wider uppercase">slug</label>
        <input type="text" name="slug" id="slug" value="%s" %s placeholder="auto-generated" pattern="[a-zA-Z0-9][a-zA-Z0-9_-]*"/>
        <p class="text-xs text-mute mt-1">url: /blog/posts/{slug}</p>
      </div>
      <div>
        <label class="block text-xs text-mute mb-1 tracking-wider uppercase">title</label>
        <input type="text" name="title" id="title" value="%s" required/>
      </div>
    </div>
    <div class="mb-4">
      <label class="block text-xs text-mute mb-1 tracking-wider uppercase">description</label>
      <input type="text" name="description" value="%s"/>
    </div>
    <div class="grid grid-cols-2 gap-6 mb-4">
      <div>
        <label class="block text-xs text-mute mb-1 tracking-wider uppercase">tags</label>
        <input type="text" name="tags" value="%s" placeholder="kafka, systems, go"/>
      </div>
      <div class="flex items-center gap-3 pt-5">
        <input type="checkbox" name="draft" id="draft" %s class="w-4 h-4 accent-[#ff5cd6]"/>
        <label for="draft" class="text-sm text-fg">save as draft</label>
      </div>
    </div>
    <div class="grid grid-cols-2 gap-6">
      <div>
        <label class="block text-xs text-mute mb-1 tracking-wider uppercase">content (markdown)</label>
        <textarea name="content" id="content" rows="28"
          class="font-mono text-sm resize-none" style="min-height:500px"
          oninput="updatePreview(this.value)">%s</textarea>
      </div>
      <div>
        <label class="block text-xs text-mute mb-1 tracking-wider uppercase">preview</label>
        <div id="preview"
          class="border border-[#2a1040] bg-[#0a0318] p-4 overflow-auto text-sm leading-relaxed prose"
          style="min-height:500px; color:#e8d6f0;"></div>
      </div>
    </div>
    <div class="flex gap-3 mt-6">
      <button type="submit" class="btn-primary text-sm px-6 py-2.5 rounded flex items-center gap-2" style="background:#ff5cd6;color:#0f041f;cursor:pointer;border:none;">save post →</button>
      <a href="/admin" class="btn-ghost text-sm px-6 py-2.5 rounded" style="text-decoration:none;">cancel</a>
    </div>
  </form>
</main>
<style>
  #preview h1,#preview h2,#preview h3{color:#fff0fb;margin:1em 0 .4em;font-weight:600;}
  #preview h2{font-size:1.2em;border-bottom:1px solid #1e0f35;padding-bottom:.3em;}
  #preview p{color:rgba(232,214,240,.8);line-height:1.75;margin:.8em 0;}
  #preview code{background:#1e0f35;padding:.1em .35em;border-radius:3px;color:#ff5cd6;font-size:.9em;}
  #preview pre{background:#140828;border:1px solid #1e0f35;padding:1rem;border-radius:4px;overflow-x:auto;}
  #preview pre code{background:none;color:#e8d6f0;}
  #preview blockquote{border-left:2px solid rgba(255,92,214,.4);padding:.2em 1em;color:rgba(232,214,240,.6);margin:1em 0;}
  #preview ul,#preview ol{margin:.8em 0 .8em 1.4em;}
  #preview li{color:rgba(232,214,240,.8);line-height:1.7;}
  #preview ul li{list-style:disc;}
  #preview strong{color:#fff0fb;}
  #preview a{color:#ff5cd6;}
</style>
<script>
  function updatePreview(md) {
    document.getElementById('preview').innerHTML = marked.parse(md || '');
  }
  updatePreview(document.getElementById('content').value);

  // Auto-generate slug from title, stop when user manually edits slug
  (function() {
    var title = document.getElementById('title');
    var slug = document.getElementById('slug');
    var dirty = slug.value !== '';
    title.addEventListener('input', function() {
      if (!dirty) {
        slug.value = title.value.toLowerCase()
          .replace(/[^a-z0-9]+/g, '-')
          .replace(/^-+|-+$/g, '');
      }
    });
    slug.addEventListener('input', function() {
      dirty = true;
    });
  })();
</script>
</body></html>`

func slugify(title string) string {
	s := strings.ToLower(title)
	s = nonAlpha.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

func NewPostForm(w http.ResponseWriter, r *http.Request) {
	csrf := randomHex(16)
	http.SetCookie(w, &http.Cookie{
		Name: "csrf_token", Value: csrf, Path: "/",
		SameSite: http.SameSiteStrictMode, MaxAge: 3600,
	})
	errMsg := ""
	if r.URL.Query().Get("err") == "taken" {
		errMsg = `<p class="text-red-400 text-xs mb-4 border border-red-900 bg-red-950/30 px-3 py-2">slug already exists — choose a different one</p>`
	}
	fmt.Fprintf(w, postFormHTML,
		"new post", sharedHead, fmt.Sprintf(adminNav, csrf),
		"", "new post",
		errMsg,
		"/admin/new", csrf,
		"", "", // slug (editable), title
		"", // description
		"", // tags
		"", // draft unchecked
		"", // content
	)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	slug := strings.TrimSpace(r.FormValue("slug"))
	if slug == "" {
		slug = slugify(title)
	}
	description := r.FormValue("description")
	tags := parseTags(r.FormValue("tags"))
	draft := r.FormValue("draft") == "on"
	content := r.FormValue("content")

	// Check slug not taken
	if _, err := models.LoadPost(slug); err == nil {
		http.Redirect(w, r, "/admin/new?err=taken", http.StatusSeeOther)
		return
	}

	if err := models.SavePost(slug, title, description, tags, draft, content); err != nil {
		http.Error(w, "failed to save: "+err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func EditPostForm(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	post, err := models.LoadPost(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	csrf := randomHex(16)
	http.SetCookie(w, &http.Cookie{
		Name: "csrf_token", Value: csrf, Path: "/",
		SameSite: http.SameSiteStrictMode, MaxAge: 3600,
	})
	draftChecked := ""
	if post.Draft {
		draftChecked = "checked"
	}
	// Extract raw markdown body from file
	body := rawBody(slug)

	fmt.Fprintf(w, postFormHTML,
		"edit: "+post.Title, sharedHead, fmt.Sprintf(adminNav, csrf),
		slug, "edit post",
		"",
		"/admin/edit/"+slug, csrf,
		slug, "readonly",
		template.HTMLEscapeString(post.Title),
		template.HTMLEscapeString(post.Description),
		template.HTMLEscapeString(strings.Join(post.Tags, ", ")),
		draftChecked,
		template.HTMLEscapeString(body),
	)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	title := r.FormValue("title")
	description := r.FormValue("description")
	tags := parseTags(r.FormValue("tags"))
	draft := r.FormValue("draft") == "on"
	content := r.FormValue("content")

	if err := models.SavePost(slug, title, description, tags, draft, content); err != nil {
		http.Error(w, "failed to save: "+err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	models.DeletePost(slug)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func parseTags(s string) []string {
	var tags []string
	for _, t := range strings.Split(s, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}

// rawBody reads the raw markdown body (after frontmatter) from disk.
func rawBody(slug string) string {
	data, err := os.ReadFile("content/posts/" + slug + ".md")
	if err != nil {
		return ""
	}
	parts := strings.SplitN(string(data), "---", 3)
	if len(parts) < 3 {
		return string(data)
	}
	return strings.TrimSpace(parts[2])
}
