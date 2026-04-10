package http

import (
	"context"
	"errors"
	httpConfig "github.com/alprnemn/yollapp-microservices/services/user/internal/config"
	"github.com/alprnemn/yollapp-microservices/services/user/internal/db"
	userHandler "github.com/alprnemn/yollapp-microservices/services/user/internal/handler/http"
	"github.com/alprnemn/yollapp-microservices/services/user/internal/repository"
	userService "github.com/alprnemn/yollapp-microservices/services/user/internal/service"
	"github.com/alprnemn/yollapp-microservices/shared/discovery"
	"github.com/alprnemn/yollapp-microservices/shared/discovery/consul"
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

	registry, err := consul.New("localhost:8500")
	if err != nil {
		log.Fatalf("err: %s", err.Error())
	}

	ctx := context.Background()

	instanceID := discovery.GenerateInstanceID(s.Config.ServerConfig.Name)

	if err := registry.Register(ctx, instanceID, s.Config.ServerConfig.Name, s.Config.ServerConfig.GetFullAddr()); err != nil {
		panic(err)
	}
	defer registry.Deregister(ctx, instanceID, s.Config.ServerConfig.Name)

	go reportHealthy(instanceID, s.Config.ServerConfig.Name, registry)

	database, err := db.New(
		s.Config.DBConfig.Address,
		s.Config.DBConfig.MaxOpenConns,
		s.Config.DBConfig.MaxIdleConns,
		s.Config.DBConfig.MaxIdleTime,
	)
	if err != nil {
		log.Fatalf(err.Error())
	}

	router := http.NewServeMux()

	repo := repository.NewRepository(database)

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

func reportHealthy(instanceID, serviceName string, registry discovery.Registry) {
	for {
		if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
			log.Println("Failed to report healthy state: " + err.Error())
		}
		time.Sleep(1 * time.Second)
	}
}
