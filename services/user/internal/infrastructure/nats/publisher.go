package nats

import (
	"github.com/nats-io/nats.go"
	"github.com/southern-martin/ecommerce/pkg/events"
)

// NewPublisher creates a new NATS JetStream publisher.
func NewPublisher(js nats.JetStreamContext) events.Publisher {
	return events.NewNATSPublisher(js)
}
