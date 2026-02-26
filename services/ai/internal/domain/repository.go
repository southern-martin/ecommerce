package domain

import "context"

// EmbeddingRepository defines the interface for embedding persistence.
type EmbeddingRepository interface {
	GetByEntity(ctx context.Context, entityType EntityType, entityID string) (*Embedding, error)
	ListByType(ctx context.Context, entityType EntityType, page, pageSize int) ([]Embedding, int64, error)
	Create(ctx context.Context, embedding *Embedding) error
	Update(ctx context.Context, embedding *Embedding) error
	Delete(ctx context.Context, id string) error
}

// RecommendationRepository defines the interface for recommendation persistence.
type RecommendationRepository interface {
	GetByID(ctx context.Context, id string) (*Recommendation, error)
	ListByUser(ctx context.Context, userID string, page, pageSize int, filterViewed *bool) ([]Recommendation, int64, error)
	Create(ctx context.Context, recommendation *Recommendation) error
	MarkViewed(ctx context.Context, id string) error
}

// AIConversationRepository defines the interface for AI conversation persistence.
type AIConversationRepository interface {
	GetByID(ctx context.Context, id string) (*AIConversation, error)
	ListByUser(ctx context.Context, userID string, page, pageSize int) ([]AIConversation, int64, error)
	Create(ctx context.Context, conversation *AIConversation) error
	Update(ctx context.Context, conversation *AIConversation) error
}

// GeneratedContentRepository defines the interface for generated content persistence.
type GeneratedContentRepository interface {
	GetByEntity(ctx context.Context, entityType ContentType, entityID string) (*GeneratedContent, error)
	Create(ctx context.Context, content *GeneratedContent) error
	Update(ctx context.Context, content *GeneratedContent) error
}
