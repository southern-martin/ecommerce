package usecase

import (
	"context"

	"github.com/southern-martin/ecommerce/services/loyalty/internal/domain"
)

// TierUseCase handles tier operations.
type TierUseCase struct {
	tierRepo domain.TierRepository
}

// NewTierUseCase creates a new TierUseCase.
func NewTierUseCase(tierRepo domain.TierRepository) *TierUseCase {
	return &TierUseCase{tierRepo: tierRepo}
}

// GetAllTiers retrieves all loyalty tiers.
func (uc *TierUseCase) GetAllTiers(ctx context.Context) ([]domain.Tier, error) {
	return uc.tierRepo.GetAll(ctx)
}
