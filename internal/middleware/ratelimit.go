package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements a simple token bucket algorithm
type RateLimiter struct {
	rate       float64
	bucketSize float64
	mu         sync.Mutex
	tokens     map[string]float64
	lastRefill map[string]time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate, bucketSize float64) *RateLimiter {
	return &RateLimiter{
		rate:       rate,
		bucketSize: bucketSize,
		tokens:     make(map[string]float64),
		lastRefill: make(map[string]time.Time),
	}
}

// RateLimit middleware implements rate limiting
func (rl *RateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		rl.mu.Lock()
		now := time.Now()
		last, exists := rl.lastRefill[ip]

		if !exists {
			rl.tokens[ip] = rl.bucketSize
			rl.lastRefill[ip] = now
		} else {
			elapsed := now.Sub(last).Seconds()
			newTokens := elapsed * rl.rate

			if tokens := rl.tokens[ip] + newTokens; tokens > rl.bucketSize {
				rl.tokens[ip] = rl.bucketSize
			} else {
				rl.tokens[ip] = tokens
			}
			rl.lastRefill[ip] = now
		}

		if rl.tokens[ip] < 1 {
			rl.mu.Unlock()
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		rl.tokens[ip]--
		rl.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
