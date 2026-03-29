package main

import (
	"github.com/alprnemn/yollapp-microservices/services/auth/internal/config"
	"github.com/alprnemn/yollapp-microservices/services/auth/internal/jwt"
	server "github.com/alprnemn/yollapp-microservices/services/auth/internal/server/http"
)

func main() {

	cfg := config.Load()

	authenticator := jwt.New(
		cfg.JWTConfig.Secret,
		cfg.JWTConfig.Issuer,
		cfg.JWTConfig.Exp,
	)

	if err := server.New(cfg, authenticator).Run(); err != nil {
		panic(err)
	}

}
