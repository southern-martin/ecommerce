package usecase

import (
	"context"

	"github.com/rs/zerolog"

	pkgerrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/auth/internal/domain"
)

// UpdateRoleInput holds the input data for updating a user's role.
type UpdateRoleInput struct {
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required,oneof=buyer seller admin"`
}

// UpdateRoleUseCase handles updating a user's role.
type UpdateRoleUseCase struct {
	repo   domain.UserRepository
	logger zerolog.Logger
}

// NewUpdateRoleUseCase creates a new UpdateRoleUseCase.
func NewUpdateRoleUseCase(
	repo domain.UserRepository,
	logger zerolog.Logger,
) *UpdateRoleUseCase {
	return &UpdateRoleUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute performs the role update.
func (uc *UpdateRoleUseCase) Execute(ctx context.Context, input UpdateRoleInput) error {
	user, err := uc.repo.GetByID(ctx, input.UserID)
	if err != nil || user == nil {
		return pkgerrors.NewNotFoundError("AUTH_USER_NOT_FOUND", "user not found")
	}

	if err := uc.repo.UpdateRole(ctx, input.UserID, input.Role); err != nil {
		uc.logger.Error().Err(err).Msg("failed to update role")
		return pkgerrors.NewInternalError("AUTH_UPDATE_ROLE_FAILED", "failed to update role")
	}

	return nil
}
