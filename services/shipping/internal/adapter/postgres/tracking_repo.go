package postgres

import (
	"context"

	"github.com/southern-martin/ecommerce/services/shipping/internal/domain"
	"gorm.io/gorm"
)

// TrackingEventRepo implements domain.TrackingEventRepository.
type TrackingEventRepo struct {
	db *gorm.DB
}

// NewTrackingEventRepo creates a new TrackingEventRepo.
func NewTrackingEventRepo(db *gorm.DB) *TrackingEventRepo {
	return &TrackingEventRepo{db: db}
}

func (r *TrackingEventRepo) GetByShipmentID(ctx context.Context, shipmentID string) ([]domain.TrackingEvent, error) {
	var models []TrackingEventModel
	if err := r.db.WithContext(ctx).Where("shipment_id = ?", shipmentID).
		Order("event_at ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	events := make([]domain.TrackingEvent, len(models))
	for i, m := range models {
		events[i] = *m.ToDomain()
	}
	return events, nil
}

func (r *TrackingEventRepo) Create(ctx context.Context, event *domain.TrackingEvent) error {
	model := ToTrackingEventModel(event)
	return r.db.WithContext(ctx).Create(model).Error
}
