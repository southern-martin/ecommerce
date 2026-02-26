package postgres

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/southern-martin/ecommerce/services/ai/internal/domain"
)

// Float64Array is a GORM-compatible type for storing []float64 as JSON text.
type Float64Array []float64

// Value implements the driver.Valuer interface for JSON storage.
func (a Float64Array) Value() (driver.Value, error) {
	if a == nil {
		return "[]", nil
	}
	b, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// Scan implements the sql.Scanner interface for JSON retrieval.
func (a *Float64Array) Scan(value interface{}) error {
	if value == nil {
		*a = Float64Array{}
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("failed to scan Float64Array: unsupported type")
	}
	return json.Unmarshal(bytes, a)
}

// EmbeddingModel is the GORM model for the embeddings table.
type EmbeddingModel struct {
	ID              string       `gorm:"type:uuid;primaryKey"`
	EntityType      string       `gorm:"type:varchar(50);not null;uniqueIndex:idx_entity"`
	EntityID        string       `gorm:"type:varchar(255);not null;uniqueIndex:idx_entity"`
	EmbeddingVector Float64Array `gorm:"type:text"`
	ModelVersion    string       `gorm:"type:varchar(100)"`
	Dimensions      int          `gorm:"not null"`
	CreatedAt       time.Time    `gorm:"autoCreateTime"`
	UpdatedAt       time.Time    `gorm:"autoUpdateTime"`
}

func (EmbeddingModel) TableName() string { return "embeddings" }

func (m *EmbeddingModel) ToDomain() *domain.Embedding {
	return &domain.Embedding{
		ID:              m.ID,
		EntityType:      domain.EntityType(m.EntityType),
		EntityID:        m.EntityID,
		EmbeddingVector: []float64(m.EmbeddingVector),
		ModelVersion:    m.ModelVersion,
		Dimensions:      m.Dimensions,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func ToEmbeddingModel(e *domain.Embedding) *EmbeddingModel {
	return &EmbeddingModel{
		ID:              e.ID,
		EntityType:      string(e.EntityType),
		EntityID:        e.EntityID,
		EmbeddingVector: Float64Array(e.EmbeddingVector),
		ModelVersion:    e.ModelVersion,
		Dimensions:      e.Dimensions,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}

// RecommendationModel is the GORM model for the recommendations table.
type RecommendationModel struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	UserID    string    `gorm:"type:uuid;index;not null"`
	ProductID string    `gorm:"type:uuid;not null"`
	Score     float64   `gorm:"not null;default:0"`
	Reason    string    `gorm:"type:varchar(100)"`
	IsViewed  bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (RecommendationModel) TableName() string { return "recommendations" }

func (m *RecommendationModel) ToDomain() *domain.Recommendation {
	return &domain.Recommendation{
		ID:        m.ID,
		UserID:    m.UserID,
		ProductID: m.ProductID,
		Score:     m.Score,
		Reason:    m.Reason,
		IsViewed:  m.IsViewed,
		CreatedAt: m.CreatedAt,
	}
}

func ToRecommendationModel(r *domain.Recommendation) *RecommendationModel {
	return &RecommendationModel{
		ID:        r.ID,
		UserID:    r.UserID,
		ProductID: r.ProductID,
		Score:     r.Score,
		Reason:    r.Reason,
		IsViewed:  r.IsViewed,
		CreatedAt: r.CreatedAt,
	}
}

// AIConversationModel is the GORM model for the ai_conversations table.
type AIConversationModel struct {
	ID           string    `gorm:"type:uuid;primaryKey"`
	UserID       string    `gorm:"type:uuid;index;not null"`
	Title        string    `gorm:"type:varchar(255)"`
	MessagesJSON string    `gorm:"type:text"`
	Model        string    `gorm:"type:varchar(50)"`
	TokenCount   int       `gorm:"default:0"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (AIConversationModel) TableName() string { return "ai_conversations" }

func (m *AIConversationModel) ToDomain() *domain.AIConversation {
	return &domain.AIConversation{
		ID:           m.ID,
		UserID:       m.UserID,
		Title:        m.Title,
		MessagesJSON: m.MessagesJSON,
		Model:        m.Model,
		TokenCount:   m.TokenCount,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func ToAIConversationModel(c *domain.AIConversation) *AIConversationModel {
	return &AIConversationModel{
		ID:           c.ID,
		UserID:       c.UserID,
		Title:        c.Title,
		MessagesJSON: c.MessagesJSON,
		Model:        c.Model,
		TokenCount:   c.TokenCount,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

// GeneratedContentModel is the GORM model for the generated_contents table.
type GeneratedContentModel struct {
	ID               string    `gorm:"type:uuid;primaryKey"`
	EntityType       string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_gen_entity"`
	EntityID         string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_gen_entity"`
	Content          string    `gorm:"type:text"`
	Model            string    `gorm:"type:varchar(50)"`
	PromptTokens     int       `gorm:"default:0"`
	CompletionTokens int       `gorm:"default:0"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
}

func (GeneratedContentModel) TableName() string { return "generated_contents" }

func (m *GeneratedContentModel) ToDomain() *domain.GeneratedContent {
	return &domain.GeneratedContent{
		ID:               m.ID,
		EntityType:       domain.ContentType(m.EntityType),
		EntityID:         m.EntityID,
		Content:          m.Content,
		Model:            m.Model,
		PromptTokens:     m.PromptTokens,
		CompletionTokens: m.CompletionTokens,
		CreatedAt:        m.CreatedAt,
	}
}

func ToGeneratedContentModel(c *domain.GeneratedContent) *GeneratedContentModel {
	return &GeneratedContentModel{
		ID:               c.ID,
		EntityType:       string(c.EntityType),
		EntityID:         c.EntityID,
		Content:          c.Content,
		Model:            c.Model,
		PromptTokens:     c.PromptTokens,
		CompletionTokens: c.CompletionTokens,
		CreatedAt:        c.CreatedAt,
	}
}
