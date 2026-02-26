package usecase

import (
	"context"

	"github.com/rs/zerolog"

	apperrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/user/internal/domain"
)

const maxAddresses = 10

// CreateAddressInput holds the fields needed to create an address.
type CreateAddressInput struct {
	Label      string `json:"label"`
	FullName   string `json:"full_name" validate:"required"`
	Phone      string `json:"phone" validate:"required"`
	Street     string `json:"street" validate:"required"`
	City       string `json:"city" validate:"required"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code" validate:"required"`
	Country    string `json:"country" validate:"required"`
}

// UpdateAddressInput holds the fields that can be updated on an address.
type UpdateAddressInput struct {
	Label      *string `json:"label"`
	FullName   *string `json:"full_name"`
	Phone      *string `json:"phone"`
	Street     *string `json:"street"`
	City       *string `json:"city"`
	State      *string `json:"state"`
	PostalCode *string `json:"postal_code"`
	Country    *string `json:"country"`
}

// AddressUseCase handles address business logic.
type AddressUseCase struct {
	repo   domain.AddressRepository
	logger zerolog.Logger
}

// NewAddressUseCase creates a new AddressUseCase.
func NewAddressUseCase(repo domain.AddressRepository, logger zerolog.Logger) *AddressUseCase {
	return &AddressUseCase{
		repo:   repo,
		logger: logger,
	}
}

// CreateAddress creates a new address for the user. Validates max 10 addresses.
// If it is the first address, it is automatically set as the default.
func (uc *AddressUseCase) CreateAddress(ctx context.Context, userID string, input CreateAddressInput) (*domain.Address, error) {
	count, err := uc.repo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if count >= maxAddresses {
		return nil, apperrors.NewValidationError("MAX_ADDRESSES", "maximum of 10 addresses allowed")
	}

	addr := &domain.Address{
		UserID:     userID,
		Label:      input.Label,
		FullName:   input.FullName,
		Phone:      input.Phone,
		Street:     input.Street,
		City:       input.City,
		State:      input.State,
		PostalCode: input.PostalCode,
		Country:    input.Country,
		IsDefault:  count == 0, // first address is default
	}

	if err := uc.repo.Create(ctx, addr); err != nil {
		uc.logger.Error().Err(err).Str("user_id", userID).Msg("failed to create address")
		return nil, err
	}

	return addr, nil
}

// ListAddresses lists all addresses for a user.
func (uc *AddressUseCase) ListAddresses(ctx context.Context, userID string) ([]domain.Address, error) {
	return uc.repo.ListByUserID(ctx, userID)
}

// UpdateAddress updates an address after verifying ownership.
func (uc *AddressUseCase) UpdateAddress(ctx context.Context, userID, addrID string, input UpdateAddressInput) (*domain.Address, error) {
	addr, err := uc.repo.GetByID(ctx, addrID)
	if err != nil {
		return nil, err
	}

	if addr.UserID != userID {
		return nil, apperrors.NewForbiddenError("FORBIDDEN", "you do not own this address")
	}

	if input.Label != nil {
		addr.Label = *input.Label
	}
	if input.FullName != nil {
		addr.FullName = *input.FullName
	}
	if input.Phone != nil {
		addr.Phone = *input.Phone
	}
	if input.Street != nil {
		addr.Street = *input.Street
	}
	if input.City != nil {
		addr.City = *input.City
	}
	if input.State != nil {
		addr.State = *input.State
	}
	if input.PostalCode != nil {
		addr.PostalCode = *input.PostalCode
	}
	if input.Country != nil {
		addr.Country = *input.Country
	}

	if err := uc.repo.Update(ctx, addr); err != nil {
		return nil, err
	}

	return addr, nil
}

// DeleteAddress deletes an address after verifying ownership.
func (uc *AddressUseCase) DeleteAddress(ctx context.Context, userID, addrID string) error {
	addr, err := uc.repo.GetByID(ctx, addrID)
	if err != nil {
		return err
	}

	if addr.UserID != userID {
		return apperrors.NewForbiddenError("FORBIDDEN", "you do not own this address")
	}

	return uc.repo.Delete(ctx, addrID)
}

// SetDefault sets a specific address as the default and clears other defaults.
func (uc *AddressUseCase) SetDefault(ctx context.Context, userID, addrID string) error {
	addr, err := uc.repo.GetByID(ctx, addrID)
	if err != nil {
		return err
	}

	if addr.UserID != userID {
		return apperrors.NewForbiddenError("FORBIDDEN", "you do not own this address")
	}

	if err := uc.repo.ClearDefaultByUserID(ctx, userID); err != nil {
		return err
	}

	addr.IsDefault = true
	return uc.repo.Update(ctx, addr)
}
