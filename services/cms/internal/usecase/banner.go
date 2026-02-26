package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/southern-martin/ecommerce/services/cms/internal/domain"
)

// BannerUseCase handles banner operations.
type BannerUseCase struct {
	bannerRepo domain.BannerRepository
	publisher  domain.EventPublisher
}

// NewBannerUseCase creates a new BannerUseCase.
func NewBannerUseCase(bannerRepo domain.BannerRepository, publisher domain.EventPublisher) *BannerUseCase {
	return &BannerUseCase{
		bannerRepo: bannerRepo,
		publisher:  publisher,
	}
}

// CreateBanner creates a new banner.
func (uc *BannerUseCase) CreateBanner(ctx context.Context, banner *domain.Banner) error {
	banner.ID = uuid.New().String()

	if err := uc.bannerRepo.Create(ctx, banner); err != nil {
		return fmt.Errorf("failed to create banner: %w", err)
	}

	_ = uc.publisher.Publish(ctx, "cms.banner.created", map[string]interface{}{
		"banner_id": banner.ID,
		"title":     banner.Title,
		"position":  banner.Position,
	})

	return nil
}

// UpdateBanner updates an existing banner.
func (uc *BannerUseCase) UpdateBanner(ctx context.Context, banner *domain.Banner) error {
	if _, err := uc.bannerRepo.GetByID(ctx, banner.ID); err != nil {
		return fmt.Errorf("banner not found: %w", err)
	}

	if err := uc.bannerRepo.Update(ctx, banner); err != nil {
		return fmt.Errorf("failed to update banner: %w", err)
	}

	_ = uc.publisher.Publish(ctx, "cms.banner.updated", map[string]interface{}{
		"banner_id": banner.ID,
	})

	return nil
}

// DeleteBanner deletes a banner by ID.
func (uc *BannerUseCase) DeleteBanner(ctx context.Context, id string) error {
	if _, err := uc.bannerRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("banner not found: %w", err)
	}

	if err := uc.bannerRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete banner: %w", err)
	}

	_ = uc.publisher.Publish(ctx, "cms.banner.deleted", map[string]interface{}{
		"banner_id": id,
	})

	return nil
}

// ListActiveBanners returns active banners, optionally filtered by position.
func (uc *BannerUseCase) ListActiveBanners(ctx context.Context, position string) ([]domain.Banner, error) {
	return uc.bannerRepo.ListActive(ctx, position)
}

// ListAllBanners returns all banners with pagination.
func (uc *BannerUseCase) ListAllBanners(ctx context.Context, page, pageSize int) ([]domain.Banner, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return uc.bannerRepo.ListAll(ctx, page, pageSize)
}
