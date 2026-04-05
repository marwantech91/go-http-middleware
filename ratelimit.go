package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type client struct {
	tokens    float64
	lastCheck time.Time
}

// RateLimit returns middleware that limits requests per client IP using a
// token bucket algorithm. maxTokens is the bucket capacity; refillPerSec
// is how many tokens are added per second.
func RateLimit(maxTokens int, refillPerSec int) Middleware {
	var mu sync.Mutex
	clients := make(map[string]*client)
	max := float64(maxTokens)
	refill := float64(refillPerSec)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			if ip == "" {
				ip = r.RemoteAddr
			}

			mu.Lock()
			c, exists := clients[ip]
			now := time.Now()

			if !exists {
				c = &client{tokens: max, lastCheck: now}
				clients[ip] = c
			}

			elapsed := now.Sub(c.lastCheck).Seconds()
			c.tokens += elapsed * refill
			if c.tokens > max {
				c.tokens = max
			}
			c.lastCheck = now

			if c.tokens < 1 {
				mu.Unlock()
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}

			c.tokens--
			mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}
