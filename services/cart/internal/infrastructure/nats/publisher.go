package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/southern-martin/ecommerce/services/cart/internal/domain"
)

// publisher implements domain.EventPublisher using NATS.
type publisher struct {
	conn   *nats.Conn
	logger zerolog.Logger
}

// NewEventPublisher creates a new NATS-backed event publisher.
func NewEventPublisher(conn *nats.Conn, logger zerolog.Logger) domain.EventPublisher {
	return &publisher{
		conn:   conn,
		logger: logger.With().Str("component", "nats_publisher").Logger(),
	}
}

// Publish publishes an event to the given NATS subject.
func (p *publisher) Publish(_ context.Context, subject string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal event payload: %w", err)
	}

	if err := p.conn.Publish(subject, data); err != nil {
		return fmt.Errorf("nats publish: %w", err)
	}

	p.logger.Debug().Str("subject", subject).Msg("event published")
	return nil
}

// Connect establishes a connection to the NATS server.
func Connect(url string, logger zerolog.Logger) (*nats.Conn, error) {
	conn, err := nats.Connect(url,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			logger.Warn().Err(err).Msg("NATS disconnected")
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			logger.Info().Msg("NATS reconnected")
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("nats connect: %w", err)
	}

	logger.Info().Str("url", url).Msg("NATS connected")
	return conn, nil
}
