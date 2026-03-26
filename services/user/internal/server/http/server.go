package http

import (
	"context"
	"database/sql"
	"errors"
	httpConfig "github.com/alprnemn/yollapp-microservices/services/user/internal/config/http"
	userHandler "github.com/alprnemn/yollapp-microservices/services/user/internal/handler/http"
	"github.com/alprnemn/yollapp-microservices/services/user/internal/repository"
	userService "github.com/alprnemn/yollapp-microservices/services/user/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server represents the HTTP server with its configuration and root handler.
type Server struct {
	Config httpConfig.Config
}

func New(config httpConfig.Config) *Server {
	sv := &Server{
		Config: config,
	}
	return sv
}

// Run starts the HTTP server and handles graceful shutdown.
func (s *Server) Run() error {
	router := http.NewServeMux()

	db := &sql.DB{}

	repo := repository.NewRepository(db)

	service := userService.NewService(repo)

	handler := userHandler.New(service)

	handler.RegisterRoutes(router)

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

	log.Printf("\033[38;5;226m Starting HTTP server on %v \033[0m", s.Config.ServerConfig.Port)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
	log.Printf("\033[38;5;214m Stopped HTTP server on %v \033[0m", s.Config.ServerConfig.Port)

	return nil
}
