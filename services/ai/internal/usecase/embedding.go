package usecase

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/southern-martin/ecommerce/services/ai/internal/domain"
	"github.com/southern-martin/ecommerce/services/ai/internal/infrastructure/aiclient"
)

// EmbeddingUseCase handles embedding operations.
type EmbeddingUseCase struct {
	embeddingRepo domain.EmbeddingRepository
	aiClient      *aiclient.MockAIClient
	publisher     domain.EventPublisher
}

// NewEmbeddingUseCase creates a new EmbeddingUseCase.
func NewEmbeddingUseCase(
	embeddingRepo domain.EmbeddingRepository,
	aiClient *aiclient.MockAIClient,
	publisher domain.EventPublisher,
) *EmbeddingUseCase {
	return &EmbeddingUseCase{
		embeddingRepo: embeddingRepo,
		aiClient:      aiClient,
		publisher:     publisher,
	}
}

// GenerateEmbeddingRequest is the input for generating an embedding.
type GenerateEmbeddingRequest struct {
	EntityType domain.EntityType
	EntityID   string
	Text       string
}

// GenerateEmbedding generates an embedding for the given text and stores it.
func (uc *EmbeddingUseCase) GenerateEmbedding(ctx context.Context, req GenerateEmbeddingRequest) (*domain.Embedding, error) {
	vector, err := uc.aiClient.GenerateEmbedding(req.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	embedding := &domain.Embedding{
		ID:              uuid.New().String(),
		EntityType:      req.EntityType,
		EntityID:        req.EntityID,
		EmbeddingVector: vector,
		ModelVersion:    "mock-v1",
		Dimensions:      len(vector),
	}

	if err := uc.embeddingRepo.Create(ctx, embedding); err != nil {
		return nil, fmt.Errorf("failed to store embedding: %w", err)
	}

	_ = uc.publisher.Publish(ctx, "ai.embedding.ready", map[string]interface{}{
		"embedding_id": embedding.ID,
		"entity_type":  string(embedding.EntityType),
		"entity_id":    embedding.EntityID,
		"dimensions":   embedding.Dimensions,
	})

	return embedding, nil
}

// GetEmbedding retrieves an embedding by entity type and ID.
func (uc *EmbeddingUseCase) GetEmbedding(ctx context.Context, entityType domain.EntityType, entityID string) (*domain.Embedding, error) {
	return uc.embeddingRepo.GetByEntity(ctx, entityType, entityID)
}

// SimilarResult represents a similar entity search result.
type SimilarResult struct {
	EntityID string  `json:"entity_id"`
	Score    float64 `json:"score"`
}

// SearchSimilar returns mock similar product IDs.
func (uc *EmbeddingUseCase) SearchSimilar(ctx context.Context, entityType domain.EntityType, entityID string, limit int) ([]SimilarResult, error) {
	if limit <= 0 {
		limit = 10
	}

	// Mock: return random product IDs with decreasing scores
	results := make([]SimilarResult, limit)
	for i := 0; i < limit; i++ {
		results[i] = SimilarResult{
			EntityID: uuid.New().String(),
			Score:    1.0 - float64(i)*0.1 - rand.Float64()*0.05,
		}
	}
	return results, nil
}
