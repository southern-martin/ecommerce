package config

import (
	"os"
)

// Config holds all configuration for the auth service.
type Config struct {
	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     string
	DBName           string
	RedisURL         string
	NATSURL          string
	JWTSecret        string
	JWTAccessExpiry  string
	JWTRefreshExpiry string
	HTTPPort         string
	GRPCPort         string
	LogLevel         string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		PostgresUser:     getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		DBName:           getEnv("DB_NAME", "ecommerce_auth"),
		RedisURL:         getEnv("REDIS_URL", "localhost:6379"),
		NATSURL:          getEnv("NATS_URL", "nats://localhost:4222"),
		JWTSecret:        getEnv("JWT_SECRET", "default-secret-change-me"),
		JWTAccessExpiry:  getEnv("JWT_ACCESS_EXPIRY", "15m"),
		JWTRefreshExpiry: getEnv("JWT_REFRESH_EXPIRY", "168h"),
		HTTPPort:         getEnv("HTTP_PORT", "8090"),
		GRPCPort:         getEnv("GRPC_PORT", "9090"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
