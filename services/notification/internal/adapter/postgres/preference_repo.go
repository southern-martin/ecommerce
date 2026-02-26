package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/southern-martin/ecommerce/services/notification/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PreferenceRepo implements domain.PreferenceRepository using GORM.
type PreferenceRepo struct {
	db *gorm.DB
}

// NewPreferenceRepo creates a new PreferenceRepo.
func NewPreferenceRepo(db *gorm.DB) *PreferenceRepo {
	return &PreferenceRepo{db: db}
}

// GetByUser retrieves all notification preferences for a user.
func (r *PreferenceRepo) GetByUser(ctx context.Context, userID string) ([]domain.NotificationPreference, error) {
	var models []PreferenceModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil, err
	}

	preferences := make([]domain.NotificationPreference, len(models))
	for i, m := range models {
		preferences[i] = *m.ToDomain()
	}

	return preferences, nil
}

// Upsert creates or updates a notification preference.
func (r *PreferenceRepo) Upsert(ctx context.Context, preference *domain.NotificationPreference) error {
	if preference.ID == "" {
		preference.ID = uuid.New().String()
	}

	model := ToPreferenceModel(preference)
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "channel"}},
			DoUpdates: clause.AssignmentColumns([]string{"enabled"}),
		}).
		Create(model).Error
}
