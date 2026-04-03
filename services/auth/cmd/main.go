package main

import (
	"github.com/alprnemn/yollapp-microservices/services/auth/internal/config"
	server "github.com/alprnemn/yollapp-microservices/services/auth/internal/server/http"
)

func main() {

	cfg := config.Load()

	if err := server.New(cfg).Run(); err != nil {
		panic(err)
	}

}
