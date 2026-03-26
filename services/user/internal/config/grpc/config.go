package grpc

import (
	"github.com/alprnemn/yollapp-microservices/shared/config"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	ServerConfig ServerConfig
}

type ServerConfig struct {
	Host string
	Port string
}

const envPath = "services/user/.env"

func Load() Config {
	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("error occurred while getting envs")
	}

	return Config{
		ServerConfig: ServerConfig{
			Host: config.GetString("PUBLIC_HOST", "http://127.0.0.1"),
			Port: config.GetString("ADDRESS", ":8080"),
		},
	}

}
