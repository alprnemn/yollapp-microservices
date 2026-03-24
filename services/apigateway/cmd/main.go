package main

import (
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

	if err := server.New(
		cfg,
		rateLimiter,
	).Run(); err != nil {
		log.Fatal(err)
	}
}
