package main

import (
	"github.com/alprnemn/yollapp-microservices/services/apigateway/internal/config"
	"github.com/alprnemn/yollapp-microservices/services/apigateway/internal/server"
	"log"
)

func main() {

	cfg := config.Load()

	if err := server.New(cfg).Run(); err != nil {
		log.Fatal(err)
	}
}
