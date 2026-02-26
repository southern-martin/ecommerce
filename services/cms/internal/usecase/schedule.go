package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/southern-martin/ecommerce/services/cms/internal/domain"
)

// ScheduleUseCase handles content scheduling operations.
type ScheduleUseCase struct {
	scheduleRepo domain.ScheduleRepository
	publisher    domain.EventPublisher
}

// NewScheduleUseCase creates a new ScheduleUseCase.
func NewScheduleUseCase(scheduleRepo domain.ScheduleRepository, publisher domain.EventPublisher) *ScheduleUseCase {
	return &ScheduleUseCase{
		scheduleRepo: scheduleRepo,
		publisher:    publisher,
	}
}

// ScheduleContent creates a new content schedule entry.
func (uc *ScheduleUseCase) ScheduleContent(ctx context.Context, schedule *domain.ContentSchedule) error {
	schedule.ID = uuid.New().String()
	schedule.Executed = false

	if err := uc.scheduleRepo.Create(ctx, schedule); err != nil {
		return fmt.Errorf("failed to create schedule: %w", err)
	}

	_ = uc.publisher.Publish(ctx, "cms.content.scheduled", map[string]interface{}{
		"schedule_id":  schedule.ID,
		"content_type": schedule.ContentType,
		"content_id":   schedule.ContentID,
		"action":       schedule.Action,
		"scheduled_at": schedule.ScheduledAt,
	})

	return nil
}

// GetPendingSchedules returns all pending (unexecuted) schedules that are due.
func (uc *ScheduleUseCase) GetPendingSchedules(ctx context.Context) ([]domain.ContentSchedule, error) {
	return uc.scheduleRepo.GetPending(ctx)
}

// ExecutePendingSchedules processes pending schedules (stub implementation).
func (uc *ScheduleUseCase) ExecutePendingSchedules(ctx context.Context) error {
	schedules, err := uc.scheduleRepo.GetPending(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pending schedules: %w", err)
	}

	for _, s := range schedules {
		log.Info().
			Str("schedule_id", s.ID).
			Str("content_type", s.ContentType).
			Str("content_id", s.ContentID).
			Str("action", s.Action).
			Msg("executing scheduled action (stub)")

		if err := uc.scheduleRepo.MarkExecuted(ctx, s.ID); err != nil {
			log.Error().Err(err).Str("schedule_id", s.ID).Msg("failed to mark schedule as executed")
			continue
		}
	}

	return nil
}
