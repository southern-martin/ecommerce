package config

import (
	"os"
)

// Config holds all configuration for the cart service.
type Config struct {
	HTTPPort         string
	GRPCPort         string
	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     string
	DBName           string
	RedisURL         string
	NATSURL          string
	LogLevel         string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		HTTPPort:         getEnv("HTTP_PORT", "8082"),
		GRPCPort:         getEnv("GRPC_PORT", "9082"),
		PostgresUser:     getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		DBName:           getEnv("DB_NAME", "ecommerce_cart"),
		RedisURL:         getEnv("REDIS_URL", "localhost:6379"),
		NATSURL:          getEnv("NATS_URL", "nats://localhost:4222"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
	}
}

// PostgresDSN returns the PostgreSQL connection string.
func (c *Config) PostgresDSN() string {
	return "host=" + c.PostgresHost +
		" user=" + c.PostgresUser +
		" password=" + c.PostgresPassword +
		" dbname=" + c.DBName +
		" port=" + c.PostgresPort +
		" sslmode=disable TimeZone=UTC"
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
