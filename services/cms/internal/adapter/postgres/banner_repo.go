package postgres

import (
	"context"
	"time"

	"github.com/southern-martin/ecommerce/services/cms/internal/domain"
	"gorm.io/gorm"
)

// BannerRepo implements domain.BannerRepository.
type BannerRepo struct {
	db *gorm.DB
}

// NewBannerRepo creates a new BannerRepo.
func NewBannerRepo(db *gorm.DB) *BannerRepo {
	return &BannerRepo{db: db}
}

func (r *BannerRepo) GetByID(ctx context.Context, id string) (*domain.Banner, error) {
	var model BannerModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *BannerRepo) ListActive(ctx context.Context, position string) ([]domain.Banner, error) {
	now := time.Now()
	query := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Where("starts_at <= ?", now).
		Where("ends_at IS NULL OR ends_at > ?", now)

	if position != "" {
		query = query.Where("position = ?", position)
	}

	var models []BannerModel
	if err := query.Order("sort_order ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	banners := make([]domain.Banner, len(models))
	for i, m := range models {
		banners[i] = *m.ToDomain()
	}
	return banners, nil
}

func (r *BannerRepo) ListAll(ctx context.Context, page, pageSize int) ([]domain.Banner, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&BannerModel{}).Count(&total)

	var models []BannerModel
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	banners := make([]domain.Banner, len(models))
	for i, m := range models {
		banners[i] = *m.ToDomain()
	}
	return banners, total, nil
}

func (r *BannerRepo) Create(ctx context.Context, banner *domain.Banner) error {
	model := ToBannerModel(banner)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *BannerRepo) Update(ctx context.Context, banner *domain.Banner) error {
	return r.db.WithContext(ctx).Model(&BannerModel{}).Where("id = ?", banner.ID).Updates(map[string]interface{}{
		"title":           banner.Title,
		"image_url":       banner.ImageURL,
		"link_url":        banner.LinkURL,
		"position":        banner.Position,
		"sort_order":      banner.SortOrder,
		"target_audience": banner.TargetAudience,
		"starts_at":       banner.StartsAt,
		"ends_at":         banner.EndsAt,
		"is_active":       banner.IsActive,
	}).Error
}

func (r *BannerRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&BannerModel{}).Error
}
