package main

import (
	"github.com/alprnemn/yollapp-microservices/services/apigateway/internal/auth/jwt"
	"github.com/alprnemn/yollapp-microservices/services/apigateway/internal/config"
	"github.com/alprnemn/yollapp-microservices/services/apigateway/internal/ratelimiter/slidingwindow"
	server "github.com/alprnemn/yollapp-microservices/services/apigateway/internal/server/http"
	"log"
)

func main() {

	cfg := config.Load()

	rateLimiter := slidingwindow.NewSlidingWindowRateLimiter(
		cfg.RLConfig.WindowSize,
		cfg.RLConfig.Limit,
	)

	authenticator := jwt.NewJWTAuthenticator(
		cfg.JWTConfig.Exp,
		cfg.JWTConfig.Secret,
		cfg.JWTConfig.Issuer,
	)

	if err := server.New(
		cfg,
		rateLimiter,
		authenticator,
	).Run(); err != nil {
		log.Fatal(err)
	}
}
