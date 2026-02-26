package usecase

import (
	"context"
	"errors"

	"github.com/southern-martin/ecommerce/services/promotion/internal/domain"
)

// CreateBundleInput represents the input for creating a bundle.
type CreateBundleInput struct {
	Name             string
	SellerID         string
	ProductIDs       []string
	BundlePriceCents int64
	SavingsCents     int64
}

// BundleUseCase handles bundle business logic.
type BundleUseCase struct {
	bundleRepo domain.BundleRepository
}

// NewBundleUseCase creates a new BundleUseCase instance.
func NewBundleUseCase(bundleRepo domain.BundleRepository) *BundleUseCase {
	return &BundleUseCase{bundleRepo: bundleRepo}
}

// CreateBundle creates a new bundle.
func (uc *BundleUseCase) CreateBundle(ctx context.Context, input CreateBundleInput) (*domain.Bundle, error) {
	if input.Name == "" {
		return nil, errors.New("bundle name is required")
	}
	if input.SellerID == "" {
		return nil, errors.New("seller_id is required")
	}
	if len(input.ProductIDs) < 2 {
		return nil, errors.New("at least two products are required for a bundle")
	}
	if input.BundlePriceCents <= 0 {
		return nil, errors.New("bundle price must be greater than 0")
	}

	bundle := domain.NewBundle(
		input.Name,
		input.SellerID,
		input.ProductIDs,
		input.BundlePriceCents,
		input.SavingsCents,
	)

	if err := uc.bundleRepo.Create(ctx, bundle); err != nil {
		return nil, err
	}

	return bundle, nil
}

// GetBundle retrieves a bundle by ID.
func (uc *BundleUseCase) GetBundle(ctx context.Context, id string) (*domain.Bundle, error) {
	if id == "" {
		return nil, errors.New("bundle id is required")
	}
	return uc.bundleRepo.GetByID(ctx, id)
}

// ListBundles retrieves a paginated list of bundles by seller.
func (uc *BundleUseCase) ListBundles(ctx context.Context, sellerID string, page, pageSize int) ([]*domain.Bundle, int64, error) {
	if sellerID == "" {
		return nil, 0, errors.New("seller id is required")
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return uc.bundleRepo.ListBySeller(ctx, sellerID, page, pageSize)
}

// ListActiveBundles retrieves a paginated list of active bundles.
func (uc *BundleUseCase) ListActiveBundles(ctx context.Context, page, pageSize int) ([]*domain.Bundle, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return uc.bundleRepo.ListActive(ctx, page, pageSize)
}

// UpdateBundle updates an existing bundle.
func (uc *BundleUseCase) UpdateBundle(ctx context.Context, bundle *domain.Bundle) error {
	return uc.bundleRepo.Update(ctx, bundle)
}
