package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"

	"github.com/southern-martin/ecommerce/services/payment/internal/domain"
)

// Publisher implements domain.EventPublisher using NATS.
type Publisher struct {
	conn *nats.Conn
}

// NewPublisher creates a new NATS publisher.
func NewPublisher(url string) (*Publisher, error) {
	nc, err := nats.Connect(url,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			log.Warn().Err(err).Msg("NATS disconnected")
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			log.Info().Msg("NATS reconnected")
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	log.Info().Str("url", url).Msg("Connected to NATS")
	return &Publisher{conn: nc}, nil
}

// Publish publishes an event to the given subject.
func (p *Publisher) Publish(_ context.Context, subject string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	if err := p.conn.Publish(subject, payload); err != nil {
		return fmt.Errorf("failed to publish event to %s: %w", subject, err)
	}

	log.Debug().Str("subject", subject).Msg("Published event")
	return nil
}

// Subscribe subscribes to a NATS subject with a handler function.
func (p *Publisher) Subscribe(subject string, handler func(data []byte)) (*nats.Subscription, error) {
	sub, err := p.conn.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg.Data)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to %s: %w", subject, err)
	}

	log.Info().Str("subject", subject).Msg("Subscribed to subject")
	return sub, nil
}

// Close closes the NATS connection.
func (p *Publisher) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}

// Ensure Publisher implements domain.EventPublisher.
var _ domain.EventPublisher = (*Publisher)(nil)
