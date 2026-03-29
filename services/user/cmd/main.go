package main

import (
	"github.com/alprnemn/yollapp-microservices/services/user/internal/config"
	httpServer "github.com/alprnemn/yollapp-microservices/services/user/internal/server/http"
)

func main() {

	cfg := config.Load()

	sv := httpServer.New(cfg)

	if err := sv.Run(); err != nil {
		panic(err)
	}

}
