package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"

	pkgerrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/auth/internal/domain"
)

// UserRepository implements domain.UserRepository using GORM and PostgreSQL.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new AuthUser into the database.
func (r *UserRepository) Create(ctx context.Context, user *domain.AuthUser) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return pkgerrors.NewInternalError("REPO_CREATE_FAILED", err.Error())
	}
	return nil
}

// GetByID retrieves an AuthUser by its ID.
func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.AuthUser, error) {
	var user domain.AuthUser
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pkgerrors.NewNotFoundError("REPO_NOT_FOUND", "user not found")
		}
		return nil, pkgerrors.NewInternalError("REPO_QUERY_FAILED", err.Error())
	}
	return &user, nil
}

// GetByEmail retrieves an AuthUser by email address.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.AuthUser, error) {
	var user domain.AuthUser
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pkgerrors.NewNotFoundError("REPO_NOT_FOUND", "user not found")
		}
		return nil, pkgerrors.NewInternalError("REPO_QUERY_FAILED", err.Error())
	}
	return &user, nil
}

// GetByOAuthProvider retrieves an AuthUser by OAuth provider and provider ID.
func (r *UserRepository) GetByOAuthProvider(ctx context.Context, provider, providerID string) (*domain.AuthUser, error) {
	var user domain.AuthUser
	if err := r.db.WithContext(ctx).
		Where("oauth_provider = ? AND oauth_provider_id = ?", provider, providerID).
		First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pkgerrors.NewNotFoundError("REPO_NOT_FOUND", "user not found")
		}
		return nil, pkgerrors.NewInternalError("REPO_QUERY_FAILED", err.Error())
	}
	return &user, nil
}

// UpdateRefreshToken updates the refresh token for a user.
func (r *UserRepository) UpdateRefreshToken(ctx context.Context, id, token string) error {
	result := r.db.WithContext(ctx).Model(&domain.AuthUser{}).Where("id = ?", id).Update("refresh_token", token)
	if result.Error != nil {
		return pkgerrors.NewInternalError("REPO_UPDATE_FAILED", result.Error.Error())
	}
	return nil
}

// UpdatePassword updates the password hash for a user.
func (r *UserRepository) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	result := r.db.WithContext(ctx).Model(&domain.AuthUser{}).Where("id = ?", id).Update("password_hash", passwordHash)
	if result.Error != nil {
		return pkgerrors.NewInternalError("REPO_UPDATE_FAILED", result.Error.Error())
	}
	return nil
}

// UpdateResetToken sets the reset token and its expiry for a user.
func (r *UserRepository) UpdateResetToken(ctx context.Context, id, token string, exp time.Time) error {
	result := r.db.WithContext(ctx).Model(&domain.AuthUser{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"reset_token":     token,
			"reset_token_exp": exp,
		})
	if result.Error != nil {
		return pkgerrors.NewInternalError("REPO_UPDATE_FAILED", result.Error.Error())
	}
	return nil
}

// ClearResetToken removes the reset token for a user.
func (r *UserRepository) ClearResetToken(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Model(&domain.AuthUser{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"reset_token":     "",
			"reset_token_exp": nil,
		})
	if result.Error != nil {
		return pkgerrors.NewInternalError("REPO_UPDATE_FAILED", result.Error.Error())
	}
	return nil
}

// UpdateRole updates the role for a user.
func (r *UserRepository) UpdateRole(ctx context.Context, id, role string) error {
	result := r.db.WithContext(ctx).Model(&domain.AuthUser{}).Where("id = ?", id).Update("role", role)
	if result.Error != nil {
		return pkgerrors.NewInternalError("REPO_UPDATE_FAILED", result.Error.Error())
	}
	return nil
}
