package proxy

import (
	"github.com/alprnemn/yollapp-microservices/pkg/utils"
	"github.com/alprnemn/yollapp-microservices/services/apigateway/internal/circuitbreaker"
	cl "github.com/alprnemn/yollapp-microservices/services/apigateway/internal/client/http"
	"github.com/alprnemn/yollapp-microservices/services/apigateway/internal/config"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Handler represents a reverse proxy handler that receives HTTP requests
// and forwards them to appropriate backend services.
type Handler struct {
	routes         map[string]*url.URL
	client         *cl.Client
	timeout        time.Duration
	rewriteMap     map[string]string
	circuitBreaker *circuitbreaker.CircuitBreaker
}

// NewHandler creates Handler instance using timeout and cfg
func NewHandler(cfg config.ClientConfig, cbCfg config.CircuitBreakerConfig) *Handler {
	return &Handler{
		routes:         make(map[string]*url.URL),
		client:         cl.NewClient(cfg),
		timeout:        cfg.Timeout,
		rewriteMap:     make(map[string]string),
		circuitBreaker: circuitbreaker.NewCircuitBreaker(cbCfg),
	}
}

// AddRoute registers a new routing rule in the handler.
// It maps an incoming request path prefix to a backend service URL.
func (h *Handler) AddRoute(prefix string, backend string) error {

	backendURL, err := url.Parse(backend)
	if err != nil {
		return err
	}

	h.routes[prefix] = backendURL
	return nil
}

// RegisterRoutes registers all route prefixes and maps them to backend services.
func (h *Handler) RegisterRoutes() {
	err := h.AddRoute("/users", "http://127.0.0.1:8082")
	if err != nil {
		panic(err)
	}
}

// ServeHTTP implements the http.Handler interface.
// It acts as a reverse proxy entry point:
//  1. Matches incoming request path to a backend using the longest prefix match
//  2. Creates a new request targeting the selected backend
//  3. Sends the request via a circuit breaker for resilience
//  4. Copies the backend response back to the client
//
// It also logs request details such as method, path, target service,
// response status, and latency for observability.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	startTime := time.Now()

	var targetURL *url.URL
	var longestPrefix string

	for prefix, backend := range h.routes {
		if len(prefix) > len(longestPrefix) && strings.HasPrefix(r.URL.Path, prefix) {
			longestPrefix = prefix
			targetURL = backend
		}
	}

	if targetURL == nil {
		utils.WriteError(w, http.StatusNotFound, "service not found")
		return
	}

	// Create backend request
	outReq := h.createProxyRequest(r, targetURL)

	// Execute request with circuit breaker
	resp, err := h.circuitBreaker.Execute(outReq)
	if err != nil {
		utils.WriteError(w, http.StatusBadGateway, "backend error")
		log.Printf("Backend error: %v", err)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Copy status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	io.Copy(w, resp.Body)

	// Log request details
	log.Printf(
		"Proxy request: %s %s -> %s, status: %d, latency: %v",
		r.Method,
		r.URL.Path,
		targetURL.String(),
		resp.StatusCode,
		time.Since(startTime),
	)
}

// createProxyRequest builds a new HTTP request to sent to the target backend.
// It:
//  1. Clones the incoming request URL and replaces scheme + host
//  2. Merges query parameters from both incoming request and target
//  3. Copies all headers to preserve metadata (auth, content-type, etc.)
//  4. Adds X-Forwarded-* headers to provide client context to the backend
//
// The returned request shares the original context for proper timeout and cancellation handling.
func (h *Handler) createProxyRequest(req *http.Request, target *url.URL) *http.Request {

	targetQuery := target.RawQuery
	outURL := *req.URL
	outURL.Scheme = target.Scheme
	outURL.Host = target.Host

	if targetQuery == "" || req.URL.RawQuery == "" {
		outURL.RawQuery = req.URL.RawQuery + targetQuery
	} else {
		outURL.RawQuery = req.URL.RawQuery + "&" + targetQuery
	}

	// Create new request
	outReq, err := http.NewRequestWithContext(
		req.Context(),
		req.Method,
		outURL.String(),
		req.Body,
	)
	if err != nil {
		log.Printf("Error creating proxy request: %v", err)
		return nil
	}

	// Copy original headers
	for key, values := range req.Header {
		for _, value := range values {
			outReq.Header.Add(key, value)
		}
	}

	// Add X-Forwarded headers
	outReq.Header.Set("X-Forwarded-For", req.RemoteAddr)
	outReq.Header.Set("X-Forwarded-Host", req.Host)
	outReq.Header.Set("X-Forwarded-Proto", req.URL.Scheme)

	return outReq
}
