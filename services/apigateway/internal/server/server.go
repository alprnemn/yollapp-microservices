package server

import (
	"github.com/alprnemn/yollapp-microservices/services/apigateway/internal/config"
	"github.com/alprnemn/yollapp-microservices/services/apigateway/internal/proxy"
	"log"
	"net/http"
	"time"
)

type Server struct {
	Config config.Config
}

func New(cfg config.Config) *Server {
	return &Server{
		Config: cfg,
	}
}

func (s *Server) Run() error {

	log.Printf("\033[38;5;226m Starting HTTP server on %v \033[0m", s.Config.ServerConfig.Port)

	handler := proxy.NewHandler(s.Config.ClientConfig, s.Config.CircuitBreaker)

	handler.RegisterRoutes()

	sv := &http.Server{
		Addr:         s.Config.ServerConfig.Port,
		Handler:      handler,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	return sv.ListenAndServe()
}
