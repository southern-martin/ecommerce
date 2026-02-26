package config

import (
	"fmt"
	"os"
)

// Config holds all configuration for the media service.
type Config struct {
	HTTPPort string
	GRPCPort string
	Postgres PostgresConfig
	NATS     NATSConfig
	S3       S3Config
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

// S3Config holds S3/MinIO storage configuration.
type S3Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
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
		HTTPPort: getEnv("HTTP_PORT", "8089"),
		GRPCPort: getEnv("GRPC_PORT", "9089"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		Postgres: PostgresConfig{
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			DBName:   getEnv("DB_NAME", "ecommerce_media"),
		},
		NATS: NATSConfig{
			URL: getEnv("NATS_URL", "nats://localhost:4222"),
		},
		S3: S3Config{
			Endpoint:  getEnv("S3_ENDPOINT", "localhost:19000"),
			AccessKey: getEnv("S3_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("S3_SECRET_KEY", "minioadmin"),
			Bucket:    getEnv("S3_BUCKET", "ecommerce-media"),
			Region:    getEnv("S3_REGION", "us-east-1"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
