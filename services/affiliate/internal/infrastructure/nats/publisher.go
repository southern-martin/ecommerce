package nats

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

// Publisher implements domain.EventPublisher using NATS.
type Publisher struct {
	conn *nats.Conn
}

// NewPublisher creates a new NATS Publisher.
func NewPublisher(url string) (*Publisher, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	log.Info().Str("url", url).Msg("connected to NATS")
	return &Publisher{conn: conn}, nil
}

// Publish publishes a domain event to the given NATS subject.
func (p *Publisher) Publish(ctx context.Context, subject string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Str("subject", subject).Msg("failed to marshal event")
		return err
	}

	if err := p.conn.Publish(subject, payload); err != nil {
		log.Error().Err(err).Str("subject", subject).Msg("failed to publish event")
		return err
	}

	log.Info().Str("subject", subject).Msg("event published")
	return nil
}

// Close closes the NATS connection.
func (p *Publisher) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}
