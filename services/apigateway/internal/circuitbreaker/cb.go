package circuitbreaker

import (
	"fmt"
	cl "github.com/alprnemn/yollapp-microservices/services/apigateway/internal/client/http"
	"github.com/alprnemn/yollapp-microservices/services/apigateway/internal/config"
	"net/http"
	"sync"
	"time"
)

// CircuitBreaker implements a simple circuit breaker pattern
// to prevent cascading failures when a backend service is failing.
type CircuitBreaker struct {
	failureThreshold int
	resetTimeout     time.Duration
	failures         int
	lastFailure      time.Time
	mu               sync.RWMutex
}

// NewCircuitBreaker creates a new CircuitBreaker with the given configuration
func NewCircuitBreaker(cfg config.CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: cfg.FailureThreshold,
		resetTimeout:     cfg.ResetTimeout,
		failures:         0,
		lastFailure:      time.Time{},
	}
}

// AllowRequest determines whether a new request should be allowed
// based on the current state of the circuit breaker (closed, open, half-open)
func (cb *CircuitBreaker) AllowRequest() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	if cb.failures < cb.failureThreshold {
		return true
	}

	if time.Since(cb.lastFailure) > cb.resetTimeout {
		return true
	}

	return false
}

// Execute performs the HTTP request using the default client
// and updates the circuit breaker state on failure
func (cb *CircuitBreaker) Execute(req *http.Request, client *cl.Client) (*http.Response, error) {
	if !cb.AllowRequest() {
		return nil, fmt.Errorf("circuit breaker open")
	}

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		cb.RecordFailure()
		return nil, err
	}

	return resp, nil
}

// RecordFailure increments the failure counter and updates the last failure timestamp
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailure = time.Now()
}
