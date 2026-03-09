package events

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

// RetryConfig configures the retry behavior.
type RetryConfig struct {
	MaxRetries int           // Maximum retry attempts (default: 3)
	BaseDelay  time.Duration // Base delay for backoff (default: 1s)
	MaxDelay   time.Duration // Maximum delay cap (default: 30s)
	DLQSubject string        // Dead-letter queue subject prefix (default: "dlq.")
}

// DefaultRetryConfig returns a RetryConfig with sensible defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries: 3,
		BaseDelay:  1 * time.Second,
		MaxDelay:   30 * time.Second,
		DLQSubject: SubjectDLQPrefix,
	}
}

// RetryMessage wraps an event with retry metadata.
type RetryMessage struct {
	Subject     string          `json:"subject"`
	Payload     json.RawMessage `json:"payload"`
	Attempt     int             `json:"attempt"`
	MaxRetries  int             `json:"max_retries"`
	FirstError  string          `json:"first_error"`
	LastError   string          `json:"last_error"`
	CreatedAt   time.Time       `json:"created_at"`
	NextRetryAt time.Time       `json:"next_retry_at"`
}

// RetryPublisher wraps a JetStream context with retry capabilities.
type RetryPublisher struct {
	js     nats.JetStreamContext
	config RetryConfig
	logger zerolog.Logger
}

// NewRetryPublisher creates a new RetryPublisher.
func NewRetryPublisher(js nats.JetStreamContext, config RetryConfig, logger zerolog.Logger) *RetryPublisher {
	if config.MaxRetries <= 0 {
		config.MaxRetries = 3
	}
	if config.BaseDelay <= 0 {
		config.BaseDelay = 1 * time.Second
	}
	if config.MaxDelay <= 0 {
		config.MaxDelay = 30 * time.Second
	}
	if config.DLQSubject == "" {
		config.DLQSubject = SubjectDLQPrefix
	}
	return &RetryPublisher{
		js:     js,
		config: config,
		logger: logger,
	}
}

// PublishToRetry publishes a failed event to the retry subject with backoff metadata.
// If max retries are exhausted, it publishes to the dead-letter queue instead.
func (p *RetryPublisher) PublishToRetry(subject string, payload []byte, handlerErr error) error {
	// Try to decode as an existing RetryMessage to increment attempt count.
	var msg RetryMessage
	if err := json.Unmarshal(payload, &msg); err != nil || msg.Subject == "" {
		// First failure — wrap the original payload.
		msg = RetryMessage{
			Subject:    subject,
			Payload:    payload,
			Attempt:    1,
			MaxRetries: p.config.MaxRetries,
			FirstError: handlerErr.Error(),
			LastError:  handlerErr.Error(),
			CreatedAt:  time.Now().UTC(),
		}
	} else {
		msg.Attempt++
		msg.LastError = handlerErr.Error()
	}

	// If retries exhausted, send to DLQ.
	if msg.Attempt > p.config.MaxRetries {
		return p.publishToDLQ(msg)
	}

	delay := p.calculateDelay(msg.Attempt)
	msg.NextRetryAt = time.Now().UTC().Add(delay)

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal retry message: %w", err)
	}

	retrySubject := SubjectRetryPrefix + subject
	if _, err := p.js.Publish(retrySubject, data); err != nil {
		return fmt.Errorf("failed to publish retry message to %s: %w", retrySubject, err)
	}

	p.logger.Warn().
		Str("subject", subject).
		Int("attempt", msg.Attempt).
		Int("max_retries", msg.MaxRetries).
		Dur("delay", delay).
		Str("error", handlerErr.Error()).
		Msg("event scheduled for retry")

	return nil
}

// publishToDLQ sends a message to the dead-letter queue.
func (p *RetryPublisher) publishToDLQ(msg RetryMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal DLQ message: %w", err)
	}

	dlqSubject := p.config.DLQSubject + msg.Subject
	if _, err := p.js.Publish(dlqSubject, data); err != nil {
		return fmt.Errorf("failed to publish to DLQ %s: %w", dlqSubject, err)
	}

	p.logger.Error().
		Str("subject", msg.Subject).
		Int("attempts", msg.Attempt).
		Str("first_error", msg.FirstError).
		Str("last_error", msg.LastError).
		Time("created_at", msg.CreatedAt).
		Msg("event moved to dead-letter queue after exhausting retries")

	return nil
}

// StartRetryConsumer listens on retry.> subjects, waits for the backoff delay,
// then re-publishes to the original subject. Run this as a goroutine.
func (p *RetryPublisher) StartRetryConsumer(ctx context.Context) error {
	if err := p.ensureStream(StreamRetry, SubjectRetryPrefix+">"); err != nil {
		return err
	}

	handler := func(msg *nats.Msg) {
		var retryMsg RetryMessage
		if unmarshalErr := json.Unmarshal(msg.Data, &retryMsg); unmarshalErr != nil {
			p.logger.Error().Err(unmarshalErr).Msg("failed to unmarshal retry message")
			_ = msg.Ack()
			return
		}

		// Wait until the scheduled retry time.
		delay := time.Until(retryMsg.NextRetryAt)
		if delay > 0 {
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				// Context cancelled; NAK so the message can be redelivered.
				_ = msg.Nak()
				return
			}
		}

		// Re-publish the original payload to the original subject.
		if _, pubErr := p.js.Publish(retryMsg.Subject, retryMsg.Payload); pubErr != nil {
			p.logger.Error().Err(pubErr).
				Str("subject", retryMsg.Subject).
				Int("attempt", retryMsg.Attempt).
				Msg("failed to re-publish retry message")
			_ = msg.Nak()
			return
		}

		p.logger.Info().
			Str("subject", retryMsg.Subject).
			Int("attempt", retryMsg.Attempt).
			Msg("retry message re-published to original subject")

		_ = msg.Ack()
	}

	sub, err := p.js.Subscribe(SubjectRetryPrefix+">", handler,
		nats.Durable("retry-consumer"), nats.ManualAck())
	if err != nil {
		// Handle stale consumer binding.
		if strings.Contains(err.Error(), "already bound") {
			_ = p.js.DeleteConsumer(StreamRetry, "retry-consumer")
			sub, err = p.js.Subscribe(SubjectRetryPrefix+">", handler,
				nats.Durable("retry-consumer"), nats.ManualAck())
		}
		if err != nil {
			return fmt.Errorf("failed to subscribe to retry subjects: %w", err)
		}
	}

	// Block until context is done, then drain.
	<-ctx.Done()
	_ = sub.Drain()
	return nil
}

// calculateDelay returns exponential backoff delay: baseDelay * 2^(attempt-1), capped at maxDelay.
func (p *RetryPublisher) calculateDelay(attempt int) time.Duration {
	delay := time.Duration(float64(p.config.BaseDelay) * math.Pow(2, float64(attempt-1)))
	if delay > p.config.MaxDelay {
		delay = p.config.MaxDelay
	}
	return delay
}

// ensureStream creates a JetStream stream if it does not already exist.
func (p *RetryPublisher) ensureStream(name, subjects string) error {
	_, err := p.js.StreamInfo(name)
	if err == nil {
		return nil
	}
	_, err = p.js.AddStream(&nats.StreamConfig{
		Name:     name,
		Subjects: []string{subjects},
	})
	if err != nil {
		return fmt.Errorf("failed to create stream %s: %w", name, err)
	}
	return nil
}
