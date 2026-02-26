package usecase

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/southern-martin/ecommerce/services/ai/internal/domain"
)

// RecommendationUseCase handles recommendation operations.
type RecommendationUseCase struct {
	recommendationRepo domain.RecommendationRepository
	publisher          domain.EventPublisher
}

// NewRecommendationUseCase creates a new RecommendationUseCase.
func NewRecommendationUseCase(
	recommendationRepo domain.RecommendationRepository,
	publisher domain.EventPublisher,
) *RecommendationUseCase {
	return &RecommendationUseCase{
		recommendationRepo: recommendationRepo,
		publisher:          publisher,
	}
}

// GetRecommendations retrieves recommendations for a user.
func (uc *RecommendationUseCase) GetRecommendations(ctx context.Context, userID string, page, pageSize int) ([]domain.Recommendation, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return uc.recommendationRepo.ListByUser(ctx, userID, page, pageSize, nil)
}

// GenerateRecommendations creates mock recommendations for a user.
func (uc *RecommendationUseCase) GenerateRecommendations(ctx context.Context, userID string) ([]domain.Recommendation, error) {
	reasons := []string{"similar_to_viewed", "frequently_bought_together", "trending", "based_on_history"}

	var recommendations []domain.Recommendation
	for i := 0; i < 5; i++ {
		rec := domain.Recommendation{
			ID:        uuid.New().String(),
			UserID:    userID,
			ProductID: uuid.New().String(),
			Score:     1.0 - float64(i)*0.15 - rand.Float64()*0.05,
			Reason:    reasons[rand.Intn(len(reasons))],
			IsViewed:  false,
		}

		if err := uc.recommendationRepo.Create(ctx, &rec); err != nil {
			return nil, fmt.Errorf("failed to create recommendation: %w", err)
		}
		recommendations = append(recommendations, rec)
	}

	_ = uc.publisher.Publish(ctx, "ai.recommendation.ready", map[string]interface{}{
		"user_id": userID,
		"count":   len(recommendations),
	})

	return recommendations, nil
}

// MarkViewed marks a recommendation as viewed.
func (uc *RecommendationUseCase) MarkViewed(ctx context.Context, id string) error {
	return uc.recommendationRepo.MarkViewed(ctx, id)
}
