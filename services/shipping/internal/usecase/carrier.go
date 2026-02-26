package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/southern-martin/ecommerce/services/shipping/internal/domain"
)

// CarrierUseCase handles carrier management operations.
type CarrierUseCase struct {
	carrierRepo    domain.CarrierRepository
	credentialRepo domain.CarrierCredentialRepository
}

// NewCarrierUseCase creates a new CarrierUseCase.
func NewCarrierUseCase(carrierRepo domain.CarrierRepository, credentialRepo domain.CarrierCredentialRepository) *CarrierUseCase {
	return &CarrierUseCase{
		carrierRepo:    carrierRepo,
		credentialRepo: credentialRepo,
	}
}

// ListCarriers returns all active carriers.
func (uc *CarrierUseCase) ListCarriers(ctx context.Context) ([]domain.Carrier, error) {
	return uc.carrierRepo.GetAll(ctx)
}

// CreateCarrier creates a new carrier (admin).
func (uc *CarrierUseCase) CreateCarrier(ctx context.Context, carrier *domain.Carrier) error {
	if carrier.Code == "" || carrier.Name == "" {
		return fmt.Errorf("carrier code and name are required")
	}
	return uc.carrierRepo.Create(ctx, carrier)
}

// UpdateCarrier updates a carrier (admin).
func (uc *CarrierUseCase) UpdateCarrier(ctx context.Context, carrier *domain.Carrier) error {
	existing, err := uc.carrierRepo.GetByCode(ctx, carrier.Code)
	if err != nil {
		return fmt.Errorf("carrier not found: %w", err)
	}

	if carrier.Name != "" {
		existing.Name = carrier.Name
	}
	existing.IsActive = carrier.IsActive
	if len(carrier.SupportedCountries) > 0 {
		existing.SupportedCountries = carrier.SupportedCountries
	}
	if carrier.APIBaseURL != "" {
		existing.APIBaseURL = carrier.APIBaseURL
	}

	return uc.carrierRepo.Update(ctx, existing)
}

// SetupSellerCarrier configures a carrier for a seller.
func (uc *CarrierUseCase) SetupSellerCarrier(ctx context.Context, sellerID, carrierCode, credentials string) (*domain.CarrierCredential, error) {
	// Verify carrier exists
	if _, err := uc.carrierRepo.GetByCode(ctx, carrierCode); err != nil {
		return nil, fmt.Errorf("carrier not found: %w", err)
	}

	// Check if already exists
	existing, _ := uc.credentialRepo.GetBySellerAndCarrier(ctx, sellerID, carrierCode)
	if existing != nil {
		existing.Credentials = credentials
		existing.IsActive = true
		if err := uc.credentialRepo.Update(ctx, existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	cred := &domain.CarrierCredential{
		ID:          uuid.New().String(),
		SellerID:    sellerID,
		CarrierCode: carrierCode,
		Credentials: credentials,
		IsActive:    true,
	}

	if err := uc.credentialRepo.Create(ctx, cred); err != nil {
		return nil, fmt.Errorf("failed to setup carrier: %w", err)
	}

	return cred, nil
}

// GetSellerCarriers returns all carrier credentials for a seller.
func (uc *CarrierUseCase) GetSellerCarriers(ctx context.Context, sellerID string) ([]domain.CarrierCredential, error) {
	return uc.credentialRepo.ListBySeller(ctx, sellerID)
}
