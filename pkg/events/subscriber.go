package events

import (
	"fmt"
	"strings"

	"github.com/nats-io/nats.go"
)

// Subscriber handles subscribing to NATS JetStream subjects.
type Subscriber struct {
	js nats.JetStreamContext
}

// NewSubscriber creates a new Subscriber with the given JetStream context.
func NewSubscriber(js nats.JetStreamContext) *Subscriber {
	return &Subscriber{js: js}
}

// Subscribe creates a durable push subscription on the given subject and processes messages
// using the provided handler function. It ensures the stream exists before subscribing.
func (s *Subscriber) Subscribe(subject, durable string, handler func(data []byte)) error {
	if err := s.ensureStream(subject); err != nil {
		return fmt.Errorf("failed to ensure stream for %s: %w", subject, err)
	}

	_, err := s.js.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg.Data)
		if ackErr := msg.Ack(); ackErr != nil {
			fmt.Printf("failed to ack message on %s: %v\n", subject, ackErr)
		}
	}, nats.Durable(durable), nats.ManualAck())
	if err != nil {
		// If the consumer is already bound (stale from a previous run), delete it and retry.
		if strings.Contains(err.Error(), "already bound") {
			streamName := strings.ToUpper(strings.Split(subject, ".")[0])
			_ = s.js.DeleteConsumer(streamName, durable)
			_, err = s.js.Subscribe(subject, func(msg *nats.Msg) {
				handler(msg.Data)
				if ackErr := msg.Ack(); ackErr != nil {
					fmt.Printf("failed to ack message on %s: %v\n", subject, ackErr)
				}
			}, nats.Durable(durable), nats.ManualAck())
		}
		if err != nil {
			return fmt.Errorf("failed to subscribe to %s: %w", subject, err)
		}
	}

	return nil
}

// ensureStream creates the JetStream stream for the given subject if it doesn't already exist.
func (s *Subscriber) ensureStream(subject string) error {
	streamName := strings.ToUpper(strings.Split(subject, ".")[0])
	subjectPrefix := strings.Split(subject, ".")[0] + ".>"

	_, err := s.js.StreamInfo(streamName)
	if err == nil {
		return nil
	}

	_, err = s.js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{subjectPrefix},
	})
	if err != nil {
		return fmt.Errorf("failed to create stream %s: %w", streamName, err)
	}

	return nil
}
