package proxy

import (
	"context"
	"errors"
	auth "github.com/alprnemn/yollapp-microservices/services/apigateway/internal/auth/jwt"
	"github.com/alprnemn/yollapp-microservices/services/apigateway/internal/circuitbreaker"
	cl "github.com/alprnemn/yollapp-microservices/services/apigateway/internal/client/http"
	"github.com/alprnemn/yollapp-microservices/services/apigateway/internal/config"
	e "github.com/alprnemn/yollapp-microservices/services/apigateway/pkg/errors"
	"github.com/alprnemn/yollapp-microservices/shared/utils"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ServiceURL struct {
	Url         *url.URL
	RequireAuth bool
}

// Handler represents a reverse proxy handler that receives HTTP requests
// and forwards them to appropriate backend services.
type Handler struct {
	Routes         map[string]*ServiceURL
	client         *cl.Client
	timeout        time.Duration
	rewriteMap     map[string]string
	circuitBreaker *circuitbreaker.CircuitBreaker
	TokenValidator *auth.TokenValidator
}

// NewHandler creates Handler instance using timeout and cfg
func NewHandler(cfg config.ClientConfig, cbCfg config.CircuitBreakerConfig, tokenValidator *auth.TokenValidator) *Handler {
	return &Handler{
		Routes:         make(map[string]*ServiceURL),
		client:         cl.NewClient(cfg),
		timeout:        cfg.Timeout,
		rewriteMap:     make(map[string]string),
		circuitBreaker: circuitbreaker.NewCircuitBreaker(cbCfg),
		TokenValidator: tokenValidator,
	}
}

// AddRoute registers a new routing rule in the handler.
// It maps an incoming request path prefix to a backend service URL.
func (h *Handler) AddRoute(prefix string, serviceAddr string, requireAuth bool) error {

	serviceURL, err := url.Parse(serviceAddr)
	if err != nil {
		return err
	}

	URL := &ServiceURL{
		Url:         serviceURL,
		RequireAuth: requireAuth,
	}

	h.Routes[prefix] = URL
	return nil
}

// RegisterRoutes registers all route prefixes and maps them to backend services.
func (h *Handler) RegisterRoutes() {

	err := h.AddRoute("/user", "http://127.0.0.1:8081", true)
	if err != nil {
		panic(err)
	}

	err = h.AddRoute("/auth/register", "http://127.0.0.1:8082", false)
	if err != nil {
		panic(err)
	}

	err = h.AddRoute("/auth/login", "http://127.0.0.1:8082", false)
	if err != nil {
		panic(err)
	}

	err = h.AddRoute("/auth/activate", "http://127.0.0.1:8082", false)
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

	var reqAuth bool
	var targetURL *url.URL
	var longestPrefix string

	for p, s := range h.Routes {
		if len(p) > len(longestPrefix) && strings.HasPrefix(r.URL.Path, p) {
			longestPrefix = p
			targetURL = s.Url
			reqAuth = s.RequireAuth
		}
	}

	if targetURL == nil {
		utils.WriteError(w, http.StatusNotFound, "service not found")
		return
	}

	if reqAuth {

		tkn, err := h.ValidateAuthHeader(r)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, err.Error())
			return
		}

		token, err := h.TokenValidator.Validate(tkn)
		if err != nil {
			utils.UnauthorizedError(w, r, err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx := context.WithValue(r.Context(), "claims", claims)
			r = r.WithContext(ctx)
		}
	}

	// Create backend request
	outReq := h.createProxyRequest(r, targetURL)

	// Execute request with circuit breaker and client
	resp, err := h.circuitBreaker.Execute(outReq, h.client)
	if err != nil {
		utils.WriteError(w, http.StatusBadGateway, "backend error")
		log.Printf("Backend error: %v", err)
		return
	}
	defer resp.Body.Close()

	if err := h.CreateNewResponse(w, resp); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		log.Printf("Backend error: %v", err)
		return
	}

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

func (h *Handler) ValidateAuthHeader(r *http.Request) (string, error) {

	authHeaders := r.Header["Authorization"]

	// Prevent multiple headers
	if len(authHeaders) == 0 {
		return "", e.ErrMissingAuthHeader
	}

	if len(authHeaders) > 1 {
		return "", e.ErrWrongAuthHeader
	}

	authHeader := strings.TrimSpace(authHeaders[0])

	parts := strings.SplitN(authHeader, " ", 2)

	if len(parts) != 2 {
		return "", e.ErrWrongAuthHeader
	}

	scheme := strings.ToLower(strings.TrimSpace(parts[0]))
	if scheme != "bearer" {
		return "", e.ErrWrongAuthHeader
	}

	token := strings.TrimSpace(parts[1])

	if token == "" {
		return "", e.ErrWrongAuthHeader
	}

	// Optional security check
	if len(token) > 4096 {
		return "", e.ErrWrongAuthHeader
	}

	return token, nil
}

func (h *Handler) CreateNewResponse(w http.ResponseWriter, resp *http.Response) error {

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Copy status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	_, err := io.Copy(w, resp.Body)
	if err != nil {
		return errors.New("backend error: io copy resp.Body")

	}
	return nil
}
