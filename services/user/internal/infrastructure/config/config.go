package config

import (
	"os"
)

// Config holds all configuration for the user service.
type Config struct {
	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     string
	DBName           string
	NATSURL          string
	HTTPPort         string
	GRPCPort         string
	LogLevel         string
	AuthGRPCAddr     string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		PostgresUser:     getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		DBName:           getEnv("DB_NAME", "ecommerce_users"),
		NATSURL:          getEnv("NATS_URL", "nats://localhost:4222"),
		HTTPPort:         getEnv("HTTP_PORT", "8091"),
		GRPCPort:         getEnv("GRPC_PORT", "9091"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		AuthGRPCAddr:     getEnv("AUTH_GRPC_ADDR", "localhost:9090"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
