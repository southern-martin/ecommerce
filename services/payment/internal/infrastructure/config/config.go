package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the payment service.
type Config struct {
	HTTPPort               string
	GRPCPort               string
	PostgresDSN            string
	NatsURL                string
	LogLevel               string
	StripeSecretKey        string
	StripeWebhookSecret    string
	PlatformCommissionRate float64
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	commissionStr := getEnv("PLATFORM_COMMISSION_RATE", "0.10")
	commissionRate, err := strconv.ParseFloat(commissionStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid PLATFORM_COMMISSION_RATE: %w", err)
	}

	pgHost := getEnv("POSTGRES_HOST", "localhost")
	pgPort := getEnv("POSTGRES_PORT", "5432")
	pgUser := getEnv("POSTGRES_USER", "postgres")
	pgPass := getEnv("POSTGRES_PASSWORD", "postgres")
	pgDB := getEnv("DB_NAME", "ecommerce_payments")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		pgHost, pgPort, pgUser, pgPass, pgDB,
	)

	return &Config{
		HTTPPort:               getEnv("HTTP_PORT", "8084"),
		GRPCPort:               getEnv("GRPC_PORT", "9084"),
		PostgresDSN:            dsn,
		NatsURL:                getEnv("NATS_URL", "nats://localhost:4222"),
		LogLevel:               getEnv("LOG_LEVEL", "info"),
		StripeSecretKey:        getEnv("STRIPE_SECRET_KEY", "sk_test_mock"),
		StripeWebhookSecret:    getEnv("STRIPE_WEBHOOK_SECRET", "whsec_mock"),
		PlatformCommissionRate: commissionRate,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
