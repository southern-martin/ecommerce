package nats

import (
	"github.com/southern-martin/ecommerce/pkg/events"
	"github.com/southern-martin/ecommerce/services/auth/internal/domain"
)

// EventPublisher wraps the shared events.Publisher to provide auth-specific publishing methods.
type EventPublisher struct {
	publisher events.Publisher
}

// NewEventPublisher creates a new EventPublisher wrapping the given Publisher.
func NewEventPublisher(publisher events.Publisher) *EventPublisher {
	return &EventPublisher{publisher: publisher}
}

// PublishUserRegistered publishes a user.registered event.
func (ep *EventPublisher) PublishUserRegistered(evt domain.UserRegisteredEvent) error {
	return ep.publisher.Publish(events.SubjectUserRegistered, evt)
}

// PublishPasswordResetRequested publishes a password.reset.requested event.
func (ep *EventPublisher) PublishPasswordResetRequested(evt domain.PasswordResetRequestedEvent) error {
	return ep.publisher.Publish("password.reset.requested", evt)
}
