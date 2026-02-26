package domain

import "time"

// EntityType represents the type of entity an embedding is for.
type EntityType string

const (
	EntityTypeProduct  EntityType = "product"
	EntityTypeCategory EntityType = "category"
	EntityTypeQuery    EntityType = "query"
)

// ContentType represents the type of generated content.
type ContentType string

const (
	ContentTypeProductDescription ContentType = "product_description"
	ContentTypeSEOMeta            ContentType = "seo_meta"
	ContentTypeReviewSummary      ContentType = "review_summary"
)

// ChatMessage represents a single message in a conversation.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Embedding stores vector embedding metadata for an entity.
type Embedding struct {
	ID              string
	EntityType      EntityType
	EntityID        string
	EmbeddingVector []float64 // stored as JSON text
	ModelVersion    string
	Dimensions      int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Recommendation represents a product recommendation for a user.
type Recommendation struct {
	ID        string
	UserID    string
	ProductID string
	Score     float64
	Reason    string // e.g. "similar_to_viewed", "frequently_bought_together", "trending"
	IsViewed  bool
	CreatedAt time.Time
}

// AIConversation represents an AI chat conversation.
type AIConversation struct {
	ID           string
	UserID       string
	Title        string
	MessagesJSON string // JSON array of ChatMessage
	Model        string // e.g. "gpt-4"
	TokenCount   int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// GeneratedContent represents AI-generated content for an entity.
type GeneratedContent struct {
	ID               string
	EntityType       ContentType
	EntityID         string
	Content          string
	Model            string
	PromptTokens     int
	CompletionTokens int
	CreatedAt        time.Time
}
