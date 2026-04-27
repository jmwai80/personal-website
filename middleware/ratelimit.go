package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type visitor struct {
	count       int
	windowStart time.Time
}

// RateLimit returns middleware that limits each IP to maxReqs per window.
func RateLimit(maxReqs int, window time.Duration) func(http.Handler) http.Handler {
	var mu sync.Mutex
	visitors := map[string]*visitor{}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			mu.Lock()
			v, ok := visitors[ip]
			now := time.Now()
			if !ok || now.Sub(v.windowStart) > window {
				visitors[ip] = &visitor{count: 1, windowStart: now}
				mu.Unlock()
				next.ServeHTTP(w, r)
				return
			}
			v.count++
			if v.count > maxReqs {
				mu.Unlock()
				http.Error(w, "Too many requests. Try again later.", http.StatusTooManyRequests)
				return
			}
			mu.Unlock()
			next.ServeHTTP(w, r)
		})
	}
}
