package postgres

import (
	"context"

	"github.com/southern-martin/ecommerce/services/ai/internal/domain"
	"gorm.io/gorm"
)

// RecommendationRepo implements domain.RecommendationRepository.
type RecommendationRepo struct {
	db *gorm.DB
}

// NewRecommendationRepo creates a new RecommendationRepo.
func NewRecommendationRepo(db *gorm.DB) *RecommendationRepo {
	return &RecommendationRepo{db: db}
}

func (r *RecommendationRepo) GetByID(ctx context.Context, id string) (*domain.Recommendation, error) {
	var model RecommendationModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *RecommendationRepo) ListByUser(ctx context.Context, userID string, page, pageSize int, filterViewed *bool) ([]domain.Recommendation, int64, error) {
	query := r.db.WithContext(ctx).Model(&RecommendationModel{}).Where("user_id = ?", userID)
	if filterViewed != nil {
		query = query.Where("is_viewed = ?", *filterViewed)
	}

	var total int64
	query.Count(&total)

	var models []RecommendationModel
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if filterViewed != nil {
				return db.Where("is_viewed = ?", *filterViewed)
			}
			return db
		}).
		Order("score DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	recommendations := make([]domain.Recommendation, len(models))
	for i, m := range models {
		recommendations[i] = *m.ToDomain()
	}
	return recommendations, total, nil
}

func (r *RecommendationRepo) Create(ctx context.Context, recommendation *domain.Recommendation) error {
	model := ToRecommendationModel(recommendation)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *RecommendationRepo) MarkViewed(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&RecommendationModel{}).Where("id = ?", id).Update("is_viewed", true).Error
}
