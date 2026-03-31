package config

import (
	"github.com/alprnemn/yollapp-microservices/shared/config"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	ServerConfig ServerConfig
	DBConfig     DBConfig
}

type DBConfig struct {
	Address      string
	Port         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type ServerConfig struct {
	Host string
	Port string
}

func (s ServerConfig) GetFullAddr() string {
	return s.Host + s.Port
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error occurred while getting envs: %s", err.Error())
	}

	return Config{
		ServerConfig: ServerConfig{
			Host: config.GetString("PUBLIC_HOST", "http://127.0.0.1"),
			Port: config.GetString("ADDRESS", ":8080"),
		},
		DBConfig: DBConfig{
			Address:      config.GetString("DB_ADDR", "postgres://user:adminpassword@localhost/yollapi?sslmode=disable"),
			MaxOpenConns: config.GetInt("DB_MAXOPENCONNS", 3),
			MaxIdleConns: config.GetInt("DB_MAXIDLECONNS", 3),
			MaxIdleTime:  config.GetString("DB_MAXIDLETIME", "15min"),
		},
	}

}
