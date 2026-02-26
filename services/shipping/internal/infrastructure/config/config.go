package config

import (
	"fmt"
	"os"
)

// Config holds all configuration for the shipping service.
type Config struct {
	HTTPPort string
	GRPCPort string
	Postgres PostgresConfig
	NATS     NATSConfig
	LogLevel string
}

// PostgresConfig holds Postgres connection configuration.
type PostgresConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

// NATSConfig holds NATS connection configuration.
type NATSConfig struct {
	URL string
}

// DSN returns the Postgres connection string.
func (c PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		c.Host, c.User, c.Password, c.DBName, c.Port,
	)
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		HTTPPort: getEnv("HTTP_PORT", "8085"),
		GRPCPort: getEnv("GRPC_PORT", "9085"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		Postgres: PostgresConfig{
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			DBName:   getEnv("DB_NAME", "ecommerce_shipping"),
		},
		NATS: NATSConfig{
			URL: getEnv("NATS_URL", "nats://localhost:4222"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
