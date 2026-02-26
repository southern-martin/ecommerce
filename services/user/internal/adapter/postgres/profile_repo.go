package postgres

import (
	"context"
	"errors"

	apperrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/user/internal/domain"
	"gorm.io/gorm"
)

// ProfileRepository implements domain.UserProfileRepository using GORM.
type ProfileRepository struct {
	db *gorm.DB
}

// NewProfileRepository creates a new ProfileRepository.
func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// Create inserts a new user profile.
func (r *ProfileRepository) Create(ctx context.Context, profile *domain.UserProfile) error {
	if err := r.db.WithContext(ctx).Create(profile).Error; err != nil {
		return apperrors.NewInternalError("DB_ERROR", "failed to create user profile")
	}
	return nil
}

// GetByID retrieves a user profile by its ID.
func (r *ProfileRepository) GetByID(ctx context.Context, id string) (*domain.UserProfile, error) {
	var profile domain.UserProfile
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("PROFILE_NOT_FOUND", "user profile not found")
		}
		return nil, apperrors.NewInternalError("DB_ERROR", "failed to get user profile")
	}
	return &profile, nil
}

// Update persists changes to an existing user profile.
func (r *ProfileRepository) Update(ctx context.Context, profile *domain.UserProfile) error {
	if err := r.db.WithContext(ctx).Save(profile).Error; err != nil {
		return apperrors.NewInternalError("DB_ERROR", "failed to update user profile")
	}
	return nil
}
