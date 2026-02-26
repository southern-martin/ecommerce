package usecase

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/southern-martin/ecommerce/pkg/events"
	apperrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/user/internal/domain"
)

// CreateSellerInput holds the fields needed to create a seller profile.
type CreateSellerInput struct {
	StoreName   string `json:"store_name" validate:"required"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
}

// UpdateSellerInput holds the fields that can be updated on a seller profile.
type UpdateSellerInput struct {
	StoreName   *string `json:"store_name"`
	Description *string `json:"description"`
	LogoURL     *string `json:"logo_url"`
}

// SellerApprovedEvent is published when a seller profile is approved.
type SellerApprovedEvent struct {
	SellerID string `json:"seller_id"`
	UserID   string `json:"user_id"`
}

// SellerUseCase handles seller profile business logic.
type SellerUseCase struct {
	repo      domain.SellerProfileRepository
	publisher events.Publisher
	logger    zerolog.Logger
}

// NewSellerUseCase creates a new SellerUseCase.
func NewSellerUseCase(repo domain.SellerProfileRepository, publisher events.Publisher, logger zerolog.Logger) *SellerUseCase {
	return &SellerUseCase{
		repo:      repo,
		publisher: publisher,
		logger:    logger,
	}
}

// CreateSeller creates a new seller profile for the user with status=pending.
func (uc *SellerUseCase) CreateSeller(ctx context.Context, userID string, input CreateSellerInput) (*domain.SellerProfile, error) {
	existing, _ := uc.repo.GetByUserID(ctx, userID)
	if existing != nil {
		return nil, apperrors.NewConflictError("SELLER_EXISTS", "seller profile already exists for this user")
	}

	seller := &domain.SellerProfile{
		UserID:      userID,
		StoreName:   input.StoreName,
		Description: input.Description,
		LogoURL:     input.LogoURL,
		Status:      "pending",
	}

	if err := uc.repo.Create(ctx, seller); err != nil {
		uc.logger.Error().Err(err).Str("user_id", userID).Msg("failed to create seller profile")
		return nil, err
	}

	return seller, nil
}

// GetSeller retrieves a seller profile by seller ID.
func (uc *SellerUseCase) GetSeller(ctx context.Context, sellerID string) (*domain.SellerProfile, error) {
	return uc.repo.GetByID(ctx, sellerID)
}

// GetSellerByUserID retrieves a seller profile by user ID.
func (uc *SellerUseCase) GetSellerByUserID(ctx context.Context, userID string) (*domain.SellerProfile, error) {
	return uc.repo.GetByUserID(ctx, userID)
}

// UpdateSeller updates a seller profile after verifying ownership.
func (uc *SellerUseCase) UpdateSeller(ctx context.Context, userID string, input UpdateSellerInput) (*domain.SellerProfile, error) {
	seller, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if input.StoreName != nil {
		seller.StoreName = *input.StoreName
	}
	if input.Description != nil {
		seller.Description = *input.Description
	}
	if input.LogoURL != nil {
		seller.LogoURL = *input.LogoURL
	}

	if err := uc.repo.Update(ctx, seller); err != nil {
		return nil, err
	}

	return seller, nil
}

// ApproveSeller approves a seller profile and publishes a seller.approved event.
func (uc *SellerUseCase) ApproveSeller(ctx context.Context, sellerID string) (*domain.SellerProfile, error) {
	seller, err := uc.repo.GetByID(ctx, sellerID)
	if err != nil {
		return nil, err
	}

	seller.Status = "approved"

	if err := uc.repo.Update(ctx, seller); err != nil {
		return nil, err
	}

	evt := SellerApprovedEvent{
		SellerID: seller.ID,
		UserID:   seller.UserID,
	}

	if err := uc.publisher.Publish("seller.approved", evt); err != nil {
		uc.logger.Error().Err(err).Str("seller_id", sellerID).Msg("failed to publish seller.approved event")
	}

	return seller, nil
}
