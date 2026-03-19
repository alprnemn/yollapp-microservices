package slidingwindow

import (
	"github.com/alprnemn/yollapp-microservices/pkg/json"
	"net/http"
	"sync"
	"time"
)

// SlidingWindowRateLimiter implements a sliding window rate limiting algorithm.
// It keeps track of request timestamps per key (e.g., IP address)
// and allows only a certain number of requests within a time window.
type SlidingWindowRateLimiter struct {
	windowSize time.Duration
	limit      int
	requests   map[string][]time.Time
	mu         sync.RWMutex
}

// NewSlidingWindowRateLimiter creates a new rate limiter instance.
// It also starts a background cleanup goroutine to prevent memory leaks.
func NewSlidingWindowRateLimiter(windowSize time.Duration, limit int) *SlidingWindowRateLimiter {
	limiter := &SlidingWindowRateLimiter{
		windowSize: windowSize,
		limit:      limit,
		requests:   make(map[string][]time.Time),
	}

	// Start cleanup routine
	go limiter.cleanup()

	return limiter
}

// Allow checks whether a request for the given key is allowed.
// It:
// 1. Removes outdated timestamps (outside the window)
// 2. Checks if the request count exceeds the limit
// 3. Adds the current request if allowed
func (l *SlidingWindowRateLimiter) Allow(key string) bool {
	now := time.Now()
	windowStart := now.Add(-l.windowSize)

	l.mu.Lock()
	defer l.mu.Unlock()

	// Remove old requests outside window
	times := l.requests[key]
	valid := 0
	for _, t := range times {
		if t.After(windowStart) {
			times[valid] = t
			valid++
		}
	}
	times = times[:valid]

	// Check if limit is reached
	if len(times) >= l.limit {
		return false
	}

	// Add new request
	l.requests[key] = append(times, now)
	return true
}

// cleanup runs periodically and removes old timestamps
// from all keys to prevent unbounded memory growth.
func (l *SlidingWindowRateLimiter) cleanup() {
	ticker := time.NewTicker(l.windowSize)
	for range ticker.C {
		l.mu.Lock()
		for key, times := range l.requests {
			windowStart := time.Now().Add(-l.windowSize)
			valid := 0
			for _, t := range times {
				if t.After(windowStart) {
					times[valid] = t
					valid++
				}
			}
			if valid == 0 {
				delete(l.requests, key)
			} else {
				l.requests[key] = times[:valid]
			}
		}
		l.mu.Unlock()
	}
}

// Middleware wraps an HTTP handler with rate limiting logic.
// It extracts a key (currently client IP) and checks if the request is allowed.
// If not allowed, it returns HTTP 429 Too Many Requests.
func (l *SlidingWindowRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.RemoteAddr // Or use a more sophisticated key

		if !l.Allow(key) {
			json.WriteError(w, http.StatusTooManyRequests, "rate limit exceeded")
			return
		}

		next.ServeHTTP(w, r)
	})
}
