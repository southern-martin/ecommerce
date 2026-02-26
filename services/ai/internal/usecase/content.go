package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/southern-martin/ecommerce/services/ai/internal/domain"
	"github.com/southern-martin/ecommerce/services/ai/internal/infrastructure/aiclient"
)

// ContentUseCase handles AI content generation operations.
type ContentUseCase struct {
	contentRepo domain.GeneratedContentRepository
	aiClient    *aiclient.MockAIClient
	publisher   domain.EventPublisher
}

// NewContentUseCase creates a new ContentUseCase.
func NewContentUseCase(
	contentRepo domain.GeneratedContentRepository,
	aiClient *aiclient.MockAIClient,
	publisher domain.EventPublisher,
) *ContentUseCase {
	return &ContentUseCase{
		contentRepo: contentRepo,
		aiClient:    aiClient,
		publisher:   publisher,
	}
}

// GenerateDescriptionRequest is the input for generating a product description.
type GenerateDescriptionRequest struct {
	ProductID   string
	ProductName string
	Category    string
}

// GenerateDescription generates a product description using the AI client.
func (uc *ContentUseCase) GenerateDescription(ctx context.Context, req GenerateDescriptionRequest) (*domain.GeneratedContent, error) {
	description, err := uc.aiClient.GenerateDescription(req.ProductName, req.Category)
	if err != nil {
		return nil, fmt.Errorf("failed to generate description: %w", err)
	}

	// Mock token counts
	promptTokens := (len(req.ProductName) + len(req.Category)) / 4
	completionTokens := len(description) / 4

	content := &domain.GeneratedContent{
		ID:               uuid.New().String(),
		EntityType:       domain.ContentTypeProductDescription,
		EntityID:         req.ProductID,
		Content:          description,
		Model:            "gpt-4",
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
	}

	if err := uc.contentRepo.Create(ctx, content); err != nil {
		return nil, fmt.Errorf("failed to store generated content: %w", err)
	}

	_ = uc.publisher.Publish(ctx, "ai.description.generated", map[string]interface{}{
		"content_id": content.ID,
		"product_id": content.EntityID,
		"model":      content.Model,
	})

	return content, nil
}

// GetGeneratedContent retrieves generated content by entity type and ID.
func (uc *ContentUseCase) GetGeneratedContent(ctx context.Context, entityType domain.ContentType, entityID string) (*domain.GeneratedContent, error) {
	return uc.contentRepo.GetByEntity(ctx, entityType, entityID)
}
