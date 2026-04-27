package handlers

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func sessionSecret() []byte {
	s := os.Getenv("SESSION_SECRET")
	if s == "" {
		s = "dev-secret-change-in-production"
	}
	return []byte(s)
}

func makeSessionToken() string {
	expiry := strconv.FormatInt(time.Now().Add(24*time.Hour).Unix(), 10)
	payload := "authenticated|" + expiry
	mac := hmac.New(sha256.New, sessionSecret())
	mac.Write([]byte(payload))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return payload + "|" + sig
}

func verifySessionToken(token string) bool {
	parts := strings.SplitN(token, "|", 3)
	if len(parts) != 3 || parts[0] != "authenticated" {
		return false
	}
	expiry, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || time.Now().Unix() > expiry {
		return false
	}
	payload := parts[0] + "|" + parts[1]
	mac := hmac.New(sha256.New, sessionSecret())
	mac.Write([]byte(payload))
	expected := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(parts[2]), []byte(expected))
}

func setSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    makeSessionToken(),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(24 * time.Hour.Seconds()),
	})
	// Non-HttpOnly flag so JS can detect logged-in state for UI
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_logged_in",
		Value:    "1",
		Path:     "/",
		HttpOnly: false,
		MaxAge:   int(24 * time.Hour.Seconds()),
	})
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_logged_in",
		Value:    "",
		Path:     "/",
		HttpOnly: false,
		MaxAge:   -1,
	})
}

func randomHex(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}

const loginHTML = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8"/>
<meta name="viewport" content="width=device-width,initial-scale=1"/>
<title>admin login</title>
<script src="https://cdn.tailwindcss.com"></script>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;600&display=swap" rel="stylesheet">
<style>
  body { background:#0f041f; font-family:"JetBrains Mono",monospace; }
  ::selection { background:rgba(255,92,214,.25); }
</style>
</head>
<body class="min-h-screen flex items-center justify-center">
  <div class="w-full max-w-sm">
    <div class="border border-[#2a1040] bg-[#140828] p-8" style="box-shadow:0 0 40px rgba(255,92,214,.08);">
      <p class="text-[#ff5cd6] text-xs mb-6 tracking-wider">$ sudo admin --login</p>
      <h1 class="text-[#e8d6f0] text-xl font-semibold mb-8">admin access</h1>
      %s
      <form method="POST" action="/admin/login">
        <input type="hidden" name="csrf_token" value="%s"/>
        <div class="mb-4">
          <label class="block text-[#6b5a7e] text-xs mb-2">password</label>
          <input type="password" name="password" autofocus
            class="w-full bg-[#0f041f] border border-[#2a1040] text-[#e8d6f0] px-3 py-2 text-sm focus:outline-none focus:border-[#ff5cd6]"
            placeholder="••••••••"/>
        </div>
        <button type="submit"
          class="w-full bg-[#ff5cd6] text-[#0f041f] font-semibold py-2 text-sm hover:opacity-90 transition-opacity mt-2">
          enter →
        </button>
      </form>
    </div>
  </div>
</body>
</html>`

func LoginForm(w http.ResponseWriter, r *http.Request) {
	if IsAuthenticated(r) {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	csrf := randomHex(16)
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrf,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		HttpOnly: false,
		MaxAge:   300,
	})
	errMsg := ""
	if r.URL.Query().Get("err") == "1" {
		errMsg = `<p class="text-red-400 text-xs mb-4 border border-red-900 bg-red-950/30 px-3 py-2">incorrect password</p>`
	}
	fmt.Fprintf(w, loginHTML, errMsg, csrf)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// verify CSRF
	cookie, err := r.Cookie("csrf_token")
	if err != nil || r.FormValue("csrf_token") != cookie.Value {
		http.Error(w, "invalid request", http.StatusForbidden)
		return
	}

	password := r.FormValue("password")
	expected := os.Getenv("ADMIN_PASSWORD")
	if expected == "" || password != expected {
		http.Redirect(w, r, "/admin/login?err=1", http.StatusSeeOther)
		return
	}

	setSessionCookie(w)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	clearSessionCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// IsAuthenticated is exported for middleware use
func IsAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("admin_session")
	if err != nil {
		return false
	}
	return verifySessionToken(cookie.Value)
}
