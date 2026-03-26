package http

import (
	"github.com/alprnemn/yollapp-microservices/pkg/config"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	ServerConfig ServerConfig
}

type ServerConfig struct {
	Host        string
	Port        string
	ServiceName string
}

func (s ServerConfig) GetFullAddr() string {
	return s.Host + s.Port
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
