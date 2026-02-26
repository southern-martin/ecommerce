package nats

import (
	"encoding/json"

	"github.com/rs/zerolog"

	"github.com/southern-martin/ecommerce/pkg/events"
	"github.com/southern-martin/ecommerce/services/user/internal/usecase"
)

// UserRegisteredEvent matches the auth service's event payload.
type UserRegisteredEvent struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

// StartSubscriber subscribes to the user.registered subject and creates user profiles.
func StartSubscriber(sub *events.Subscriber, profileUC *usecase.ProfileUseCase, logger zerolog.Logger) error {
	return sub.Subscribe(events.SubjectUserRegistered, "user-service-registered", func(data []byte) {
		var evt UserRegisteredEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal user.registered event")
			return
		}

		logger.Info().
			Str("user_id", evt.UserID).
			Str("email", evt.Email).
			Msg("received user.registered event")

		if err := profileUC.CreateFromEvent(evt.UserID, evt.Email, evt.Role); err != nil {
			logger.Error().Err(err).Str("user_id", evt.UserID).Msg("failed to create profile from event")
		}
	})
}
