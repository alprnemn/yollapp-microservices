package config

import (
	"github.com/alprnemn/yollapp-microservices/shared/config"
	"github.com/joho/godotenv"
	"log"
	"time"
)

type Config struct {
	ServerConfig ServerConfig
	JWTConfig    JWTConfig
}

type ServerConfig struct {
	Host string
	Port string
}

func (s ServerConfig) GetFullAddr() string {
	return s.Host + s.Port
}

type JWTConfig struct {
	Secret string
	Exp    time.Duration
	Issuer string
}

func Load() Config {
	if err := godotenv.Load("services/auth/.env"); err != nil {
		log.Printf(err.Error())
		log.Fatal("error occurred while getting envs")
	}

	return Config{
		ServerConfig: ServerConfig{
			Host: config.GetString("PUBLIC_HOST", "http://127.0.0.1"),
			Port: config.GetString("ADDRESS", ":8082"),
		},
		JWTConfig: JWTConfig{
			Secret: config.GetString("JWT_SECRET_KEY", "asdasd"),
			Issuer: config.GetString("JWT_ISSUER", "asdfg"),
			Exp:    time.Duration(config.GetInt("JWT_EXP_SECOND", 15)),
		},
	}

}
