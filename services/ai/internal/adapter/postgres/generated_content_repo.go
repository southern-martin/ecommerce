package postgres

import (
	"context"

	"github.com/southern-martin/ecommerce/services/ai/internal/domain"
	"gorm.io/gorm"
)

// GeneratedContentRepo implements domain.GeneratedContentRepository.
type GeneratedContentRepo struct {
	db *gorm.DB
}

// NewGeneratedContentRepo creates a new GeneratedContentRepo.
func NewGeneratedContentRepo(db *gorm.DB) *GeneratedContentRepo {
	return &GeneratedContentRepo{db: db}
}

func (r *GeneratedContentRepo) GetByEntity(ctx context.Context, entityType domain.ContentType, entityID string) (*domain.GeneratedContent, error) {
	var model GeneratedContentModel
	if err := r.db.WithContext(ctx).Where("entity_type = ? AND entity_id = ?", string(entityType), entityID).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *GeneratedContentRepo) Create(ctx context.Context, content *domain.GeneratedContent) error {
	model := ToGeneratedContentModel(content)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *GeneratedContentRepo) Update(ctx context.Context, content *domain.GeneratedContent) error {
	return r.db.WithContext(ctx).Model(&GeneratedContentModel{}).Where("id = ?", content.ID).Updates(map[string]interface{}{
		"content":          content.Content,
		"model":            content.Model,
		"prompt_tokens":    content.PromptTokens,
		"completion_tokens": content.CompletionTokens,
	}).Error
}
