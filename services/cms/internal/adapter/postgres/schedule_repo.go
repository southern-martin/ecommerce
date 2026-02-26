package postgres

import (
	"context"
	"time"

	"github.com/southern-martin/ecommerce/services/cms/internal/domain"
	"gorm.io/gorm"
)

// ScheduleRepo implements domain.ScheduleRepository.
type ScheduleRepo struct {
	db *gorm.DB
}

// NewScheduleRepo creates a new ScheduleRepo.
func NewScheduleRepo(db *gorm.DB) *ScheduleRepo {
	return &ScheduleRepo{db: db}
}

func (r *ScheduleRepo) GetPending(ctx context.Context) ([]domain.ContentSchedule, error) {
	now := time.Now()
	var models []ContentScheduleModel
	if err := r.db.WithContext(ctx).
		Where("scheduled_at <= ? AND executed = ?", now, false).
		Order("scheduled_at ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	schedules := make([]domain.ContentSchedule, len(models))
	for i, m := range models {
		schedules[i] = *m.ToDomain()
	}
	return schedules, nil
}

func (r *ScheduleRepo) Create(ctx context.Context, schedule *domain.ContentSchedule) error {
	model := ToContentScheduleModel(schedule)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *ScheduleRepo) MarkExecuted(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&ContentScheduleModel{}).Where("id = ?", id).Update("executed", true).Error
}
