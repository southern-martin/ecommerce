package logger

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type contextKey string

const correlationIDKey contextKey = "correlation_id"

// New creates a new zerolog JSON logger with the given log level and service name.
func New(level, serviceName string) zerolog.Logger {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	return zerolog.New(os.Stdout).
		Level(lvl).
		With().
		Timestamp().
		Str("service", serviceName).
		Logger().
		Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			NoColor:    true,
		})
}

// WithCorrelationID extracts the correlation ID from the context and adds it to the logger.
func WithCorrelationID(ctx context.Context, logger zerolog.Logger) zerolog.Logger {
	if id, ok := ctx.Value(correlationIDKey).(string); ok && id != "" {
		return logger.With().Str("correlation_id", id).Logger()
	}
	return logger
}

// SetCorrelationID stores a correlation ID into the context.
func SetCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlationIDKey, id)
}

// GetCorrelationID retrieves the correlation ID from the context.
func GetCorrelationID(ctx context.Context) string {
	if id, ok := ctx.Value(correlationIDKey).(string); ok {
		return id
	}
	return ""
}
