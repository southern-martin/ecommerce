package postgres

import (
	"context"

	"github.com/southern-martin/ecommerce/services/media/internal/domain"
	"gorm.io/gorm"
)

// MediaRepo implements domain.MediaRepository.
type MediaRepo struct {
	db *gorm.DB
}

// NewMediaRepo creates a new MediaRepo.
func NewMediaRepo(db *gorm.DB) *MediaRepo {
	return &MediaRepo{db: db}
}

func (r *MediaRepo) GetByID(ctx context.Context, id string) (*domain.Media, error) {
	var model MediaModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *MediaRepo) ListByOwner(ctx context.Context, ownerID, ownerType string, page, pageSize int) ([]domain.Media, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&MediaModel{})

	if ownerID != "" {
		query = query.Where("owner_id = ?", ownerID)
	}
	if ownerType != "" {
		query = query.Where("owner_type = ?", ownerType)
	}

	query.Count(&total)

	var models []MediaModel
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	media := make([]domain.Media, len(models))
	for i, m := range models {
		media[i] = *m.ToDomain()
	}
	return media, total, nil
}

func (r *MediaRepo) Create(ctx context.Context, media *domain.Media) error {
	model := ToMediaModel(media)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *MediaRepo) Update(ctx context.Context, media *domain.Media) error {
	return r.db.WithContext(ctx).Model(&MediaModel{}).Where("id = ?", media.ID).Updates(map[string]interface{}{
		"file_name":     media.FileName,
		"original_name": media.OriginalName,
		"content_type":  media.ContentType,
		"size_bytes":    media.SizeBytes,
		"url":           media.URL,
		"thumbnail_url": media.ThumbnailURL,
		"width":         media.Width,
		"height":        media.Height,
		"status":        string(media.Status),
	}).Error
}

func (r *MediaRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&MediaModel{}).Error
}
