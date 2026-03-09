package events

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

// DLQHandler processes messages that exhausted all retries.
type DLQHandler struct {
	js     nats.JetStreamContext
	logger zerolog.Logger
}

// NewDLQHandler creates a new DLQHandler.
func NewDLQHandler(js nats.JetStreamContext, logger zerolog.Logger) *DLQHandler {
	return &DLQHandler{
		js:     js,
		logger: logger,
	}
}

// StartDLQConsumer listens on dlq.> subjects and logs failed messages for inspection.
// Run this as a goroutine.
func (h *DLQHandler) StartDLQConsumer(ctx context.Context) error {
	if err := h.ensureStream(StreamDLQ, SubjectDLQPrefix+">"); err != nil {
		return err
	}

	handler := func(msg *nats.Msg) {
		var retryMsg RetryMessage
		if err := json.Unmarshal(msg.Data, &retryMsg); err != nil {
			h.logger.Error().Err(err).
				Str("raw_subject", msg.Subject).
				Msg("failed to unmarshal DLQ message")
			_ = msg.Ack()
			return
		}

		h.logger.Error().
			Str("original_subject", retryMsg.Subject).
			Int("total_attempts", retryMsg.Attempt).
			Str("first_error", retryMsg.FirstError).
			Str("last_error", retryMsg.LastError).
			Time("created_at", retryMsg.CreatedAt).
			RawJSON("payload", retryMsg.Payload).
			Msg("DLQ: message permanently failed")

		_ = msg.Ack()
	}

	sub, err := h.js.Subscribe(SubjectDLQPrefix+">", handler,
		nats.Durable("dlq-consumer"), nats.ManualAck())
	if err != nil {
		if strings.Contains(err.Error(), "already bound") {
			_ = h.js.DeleteConsumer(StreamDLQ, "dlq-consumer")
			sub, err = h.js.Subscribe(SubjectDLQPrefix+">", handler,
				nats.Durable("dlq-consumer"), nats.ManualAck())
		}
		if err != nil {
			return fmt.Errorf("failed to subscribe to DLQ subjects: %w", err)
		}
	}

	<-ctx.Done()
	_ = sub.Drain()
	return nil
}

// ReplayDLQ re-publishes all pending messages on the given DLQ subject back to
// their original subjects for reprocessing. The dlqSubject should be the full
// DLQ subject (e.g. "dlq.notification.email").
func (h *DLQHandler) ReplayDLQ(ctx context.Context, dlqSubject string) error {
	// Create a pull subscription to iterate over existing messages.
	subName := "dlq-replay-" + strings.ReplaceAll(dlqSubject, ".", "-")
	sub, err := h.js.PullSubscribe(dlqSubject, subName)
	if err != nil {
		return fmt.Errorf("failed to pull-subscribe to %s: %w", dlqSubject, err)
	}
	defer func() {
		_ = sub.Unsubscribe()
	}()

	replayed := 0
	for {
		select {
		case <-ctx.Done():
			h.logger.Info().Int("replayed", replayed).Str("subject", dlqSubject).Msg("DLQ replay cancelled")
			return ctx.Err()
		default:
		}

		msgs, err := sub.Fetch(10, nats.MaxWait(2*1e9)) // 2 second timeout
		if err != nil {
			// No more messages available.
			break
		}

		for _, msg := range msgs {
			var retryMsg RetryMessage
			if unmarshalErr := json.Unmarshal(msg.Data, &retryMsg); unmarshalErr != nil {
				h.logger.Error().Err(unmarshalErr).Msg("failed to unmarshal DLQ message during replay")
				_ = msg.Ack()
				continue
			}

			// Re-publish the original payload to the original subject.
			if _, pubErr := h.js.Publish(retryMsg.Subject, retryMsg.Payload); pubErr != nil {
				h.logger.Error().Err(pubErr).
					Str("subject", retryMsg.Subject).
					Msg("failed to replay DLQ message")
				_ = msg.Nak()
				continue
			}

			h.logger.Info().
				Str("original_subject", retryMsg.Subject).
				Int("original_attempts", retryMsg.Attempt).
				Msg("DLQ message replayed to original subject")

			_ = msg.Ack()
			replayed++
		}
	}

	h.logger.Info().Int("replayed", replayed).Str("subject", dlqSubject).Msg("DLQ replay completed")
	return nil
}

// ensureStream creates a JetStream stream if it does not already exist.
func (h *DLQHandler) ensureStream(name, subjects string) error {
	_, err := h.js.StreamInfo(name)
	if err == nil {
		return nil
	}
	_, err = h.js.AddStream(&nats.StreamConfig{
		Name:     name,
		Subjects: []string{subjects},
	})
	if err != nil {
		return fmt.Errorf("failed to create stream %s: %w", name, err)
	}
	return nil
}
