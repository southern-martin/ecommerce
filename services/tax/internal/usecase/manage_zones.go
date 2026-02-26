package usecase

import (
	"context"

	"github.com/southern-martin/ecommerce/services/tax/internal/domain"
)

// ManageZonesUseCase handles operations for TaxZones.
type ManageZonesUseCase struct {
	zoneRepo domain.TaxZoneRepository
}

// NewManageZonesUseCase creates a new ManageZonesUseCase.
func NewManageZonesUseCase(zoneRepo domain.TaxZoneRepository) *ManageZonesUseCase {
	return &ManageZonesUseCase{zoneRepo: zoneRepo}
}

// ListZones returns all tax zones.
func (uc *ManageZonesUseCase) ListZones(ctx context.Context) ([]*domain.TaxZone, error) {
	return uc.zoneRepo.List(ctx)
}

// GetZoneByID returns a tax zone by ID.
func (uc *ManageZonesUseCase) GetZoneByID(ctx context.Context, id string) (*domain.TaxZone, error) {
	return uc.zoneRepo.GetByID(ctx, id)
}

// GetZoneByLocation returns a tax zone by country and state code.
func (uc *ManageZonesUseCase) GetZoneByLocation(ctx context.Context, countryCode, stateCode string) (*domain.TaxZone, error) {
	return uc.zoneRepo.GetByLocation(ctx, countryCode, stateCode)
}
