package config

import (
	"fmt"
	"os"
)

// Config holds all configuration for the notification service.
type Config struct {
	HTTPPort string
	GRPCPort string
	Postgres PostgresConfig
	NATS     NATSConfig
	SMTP     SMTPConfig
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

// SMTPConfig holds SMTP connection configuration.
type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
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
		HTTPPort: getEnv("HTTP_PORT", "8092"),
		GRPCPort: getEnv("GRPC_PORT", "9092"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		Postgres: PostgresConfig{
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			DBName:   getEnv("DB_NAME", "ecommerce_notifications"),
		},
		NATS: NATSConfig{
			URL: getEnv("NATS_URL", "nats://localhost:4222"),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "localhost"),
			Port:     getEnv("SMTP_PORT", "1025"),
			User:     getEnv("SMTP_USER", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
