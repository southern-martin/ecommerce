package config

import (
	"fmt"
	"os"
)

// Config holds all configuration for the tax service.
type Config struct {
	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     string
	DBName           string
	HTTPPort         string
	GRPCPort         string
	LogLevel         string
}

// Load reads configuration from environment variables with defaults.
func Load() *Config {
	return &Config{
		PostgresUser:     getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		DBName:           getEnv("DB_NAME", "ecommerce_tax"),
		HTTPPort:         getEnv("HTTP_PORT", "8098"),
		GRPCPort:         getEnv("GRPC_PORT", "9098"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
	}
}

// DSN returns the PostgreSQL connection string.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		c.PostgresHost, c.PostgresUser, c.PostgresPassword, c.DBName, c.PostgresPort,
	)
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
