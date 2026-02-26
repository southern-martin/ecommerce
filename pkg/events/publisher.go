package events

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

// Publisher defines the interface for publishing events.
type Publisher interface {
	Publish(subject string, data interface{}) error
}

// NATSPublisher implements Publisher using NATS JetStream.
type NATSPublisher struct {
	js nats.JetStreamContext
}

// NewNATSPublisher creates a new NATSPublisher with the given JetStream context.
func NewNATSPublisher(js nats.JetStreamContext) *NATSPublisher {
	return &NATSPublisher{js: js}
}

// Publish marshals the data as JSON and publishes it to the given NATS subject.
func (p *NATSPublisher) Publish(subject string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	_, err = p.js.Publish(subject, payload)
	if err != nil {
		return fmt.Errorf("failed to publish event to %s: %w", subject, err)
	}

	return nil
}
