package config

import (
	"github.com/alprnemn/yollapp-microservices/pkg/config"
	"github.com/joho/godotenv"
	"log"
	"time"
)

type Config struct {
	ServerConfig   ServerConfig
	ClientConfig   ClientConfig
	RLConfig       RLConfig
	JWTConfig      JWTConfig
	CircuitBreaker CircuitBreakerConfig
	CORSConfig     CORSConfig
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

type ServerConfig struct {
	Host string
	Port string
}

func (s ServerConfig) GetFullAddr() string {
	return s.Host + s.Port
}

type RLConfig struct {
	WindowSize time.Duration
	Limit      int
}

type JWTConfig struct {
	Secret string
	Exp    time.Duration
	Issuer string
}

type ClientConfig struct {
	Timeout         time.Duration
	TransportConfig TransportConfig
}

type TransportConfig struct {
	MaxIdleConns          int
	MaxIdleConnsPerHost   int
	MaxConnsPerHost       int
	IdleConnTimeout       time.Duration
	DisableCompression    bool
	DisableKeepAlives     bool
	TLSHandshakeTimeout   time.Duration
	ExpectContinueTimeout time.Duration
}

// CircuitBreakerConfig holds configuration for a CircuitBreaker
type CircuitBreakerConfig struct {
	FailureThreshold int           // Number of consecutive failures to open the circuit
	ResetTimeout     time.Duration // Duration to wait before allowing requests again (half-open state)
}

const envPath = "services/apigateway/.env"

func Load() Config {
	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("error occurred while getting envs")
	}

	return Config{
		ServerConfig: ServerConfig{
			Host: config.GetString("PUBLIC_HOST", "http://127.0.0.1"),
			Port: config.GetString("ADDRESS", ":8080"),
		},
		ClientConfig: ClientConfig{
			Timeout: time.Duration(config.GetInt("CLIENT_TIMEOUT_SEC", 5)) * time.Second,
			TransportConfig: TransportConfig{
				MaxIdleConns:          config.GetInt("MAX_IDLE_CONNS", 100),
				MaxIdleConnsPerHost:   config.GetInt("MAX_IDLE_CONNS_PER_HOST", 10),
				MaxConnsPerHost:       config.GetInt("MAX_CONNS_PER_HOST", 100),
				IdleConnTimeout:       time.Duration(config.GetInt("IDLE_CONN_TIMEOUT_SEC", 90)) * time.Second,
				DisableCompression:    config.GetBool("DISABLE_COMPRESSION", false),
				DisableKeepAlives:     config.GetBool("DISABLE_KEEP_ALIVES", false),
				TLSHandshakeTimeout:   time.Duration(config.GetInt("TLS_HANDSHAKE_TIMEOUT_SEC", 10)) * time.Second,
				ExpectContinueTimeout: time.Duration(config.GetInt("EXPECT_CONTINUE_TIMEOUT_SEC", 1)) * time.Second,
			},
		},
		CircuitBreaker: CircuitBreakerConfig{
			FailureThreshold: config.GetInt("CB_FAILURE_THRESHOLD", 5),
			ResetTimeout:     time.Duration(config.GetInt("CB_RESET_TIMEOUT_SEC", 10)) * time.Second,
		},
		CORSConfig: CORSConfig{
			AllowedOrigins:   config.GetSlice("CORS_ALLOWED_ORIGINS", []string{"http://127.0.0.1:3000"}),
			AllowedMethods:   config.GetSlice("CORS_ALLOWED_METHODS", []string{"*"}),
			AllowedHeaders:   config.GetSlice("CORS_ALLOWED_HEADERS", []string{"*"}),
			ExposedHeaders:   config.GetSlice("CORS_EXPOSED_HEADERS", []string{"*"}),
			AllowCredentials: config.GetBool("CORS_ALLOWED_CREDENTIALS", true),
			MaxAge:           config.GetInt("CORS_MAX_AGE", 300),
		},
		RLConfig: RLConfig{
			WindowSize: time.Duration(config.GetInt("RL_WINDOW_SIZE", 60)) * time.Second,
			Limit:      config.GetInt("RL_LIMIT", 100),
		},
	}

}
