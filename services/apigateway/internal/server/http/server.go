package http

import (
	"context"
	"errors"
	"github.com/alprnemn/yollapp-microservices/services/apigateway/internal/config"
	"github.com/alprnemn/yollapp-microservices/services/apigateway/internal/proxy"
	rl "github.com/alprnemn/yollapp-microservices/services/apigateway/internal/ratelimiter/slidingwindow"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server represents the HTTP server with its configuration and root handler.
type Server struct {
	Config      config.Config
	Handler     http.Handler
	RateLimiter *rl.SlidingWindowRateLimiter
}

// New creates a new Server instance and initializes routes/middleware.
func New(cfg config.Config, rateLimiter *rl.SlidingWindowRateLimiter) *Server {
	s := &Server{
		Config:      cfg,
		RateLimiter: rateLimiter,
	}
	s.Mount()
	return s
}

// Mount sets up the router, middleware stack, and mounts all routes.
func (s *Server) Mount() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   s.Config.CORSConfig.AllowedOrigins,
		AllowedMethods:   s.Config.CORSConfig.AllowedMethods,
		AllowedHeaders:   s.Config.CORSConfig.AllowedHeaders,
		ExposedHeaders:   s.Config.CORSConfig.ExposedHeaders,
		AllowCredentials: s.Config.CORSConfig.AllowCredentials,
		MaxAge:           s.Config.CORSConfig.MaxAge,
	}))

	r.Use(s.RateLimiter.Middleware)

	proxyHandler := proxy.NewHandler(s.Config.ClientConfig, s.Config.CircuitBreaker)

	proxyHandler.RegisterRoutes()

	r.Mount("/", proxyHandler)

	s.Handler = r
}

// Run starts the HTTP server and handles graceful shutdown.
func (s *Server) Run() error {

	srv := &http.Server{
		Addr:         s.Config.ServerConfig.Port,
		Handler:      s.Handler,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("\033[38;5;226m Starting HTTP server on %v \033[0m", s.Config.ServerConfig.Port)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
	log.Printf("\033[38;5;214m Stopped HTTP server on %v \033[0m", s.Config.ServerConfig.Port)

	return nil
}
