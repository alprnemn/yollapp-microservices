package config

import (
	"os"
	"strconv"
	"strings"
)

// GetString gets config as a string using key
func GetString(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return value
}

// GetInt gets config as a int using key
func GetInt(key string, fallback int) int {

	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	valInt, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return valInt
}

// GetBool gets bool env value
func GetBool(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return boolVal
}

// GetSlice gets string slice from .env
func GetSlice(key string, fallback []string) []string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	parts := strings.Split(val, ",")

	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	return parts
}
