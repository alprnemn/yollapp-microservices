package http

import (
	"context"
	"errors"
	"github.com/alprnemn/yollapp-microservices/services/auth/internal/config"
	"github.com/alprnemn/yollapp-microservices/services/auth/internal/db"
	userGateway "github.com/alprnemn/yollapp-microservices/services/auth/internal/gateway/http/user"
	handler "github.com/alprnemn/yollapp-microservices/services/auth/internal/handler/http"
	"github.com/alprnemn/yollapp-microservices/services/auth/internal/jwt"
	"github.com/alprnemn/yollapp-microservices/services/auth/internal/mailer/resend"
	"github.com/alprnemn/yollapp-microservices/services/auth/internal/repository"
	"github.com/alprnemn/yollapp-microservices/services/auth/internal/service"
	"github.com/alprnemn/yollapp-microservices/shared/discovery"
	"github.com/alprnemn/yollapp-microservices/shared/discovery/consul"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	router := http.NewServeMux()

	authDB, err := db.Init()
	if err != nil {
		log.Fatalf("error starting database: %s", err.Error())
	}

	repo := repository.New(authDB)

	usergway := userGateway.New(":8081")

	authenticator := jwt.New(
		s.Config.JWTConfig.Secret,
		s.Config.JWTConfig.Issuer,
		s.Config.JWTConfig.Exp,
	)

	resendMailer := resend.NewResendMailer(
		s.Config.MailerConfig.FromMail,
		s.Config.MailerConfig.ApiKey,
	)

	svc := service.New(usergway, authenticator, resendMailer, repo)

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

func reportHealthy(instanceID, serviceName string, registry discovery.Registry) {
	for {
		if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
			log.Println("Failed to report healthy state: " + err.Error())
		}
		time.Sleep(1 * time.Second)
	}
}
