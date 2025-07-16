package middleware

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

// per-client token‑bucket limiter
var (
	limiters = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	l, exists := limiters[ip]
	if !exists {
		l = rate.NewLimiter(1, 5) // 1 req/sec, burst 5
		limiters[ip] = l
	}
	return l
}

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter := getLimiter(ip)
		if !limiter.Allow() {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
