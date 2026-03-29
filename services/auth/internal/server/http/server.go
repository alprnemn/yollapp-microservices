package http

import (
	"context"
	"errors"
	"github.com/alprnemn/yollapp-microservices/services/auth/internal/config"
	handler "github.com/alprnemn/yollapp-microservices/services/auth/internal/handler/http"
	"github.com/alprnemn/yollapp-microservices/services/auth/internal/jwt"
	"github.com/alprnemn/yollapp-microservices/services/auth/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Config        config.Config
	Authenticator *jwt.Authenticator
}

func New(cfg config.Config, authenticator *jwt.Authenticator) *Server {
	return &Server{
		Config:        cfg,
		Authenticator: authenticator,
	}
}

func (s *Server) Run() error {

	router := http.NewServeMux()

	svc := service.New()

	h := handler.New(svc)

	h.RegisterRoutes(router)

	srv := &http.Server{
		Addr:         s.Config.ServerConfig.Port,
		Handler:      router,
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

	log.Printf("\033[38;5;226m Starting HTTP Auth server on %v \033[0m", s.Config.ServerConfig.Port)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
	log.Printf("\033[38;5;214m Stopped HTTP server on %v \033[0m", s.Config.ServerConfig.Port)

	return nil
}
