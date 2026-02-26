package postgres

import (
	"context"

	"github.com/southern-martin/ecommerce/services/ai/internal/domain"
	"gorm.io/gorm"
)

// EmbeddingRepo implements domain.EmbeddingRepository.
type EmbeddingRepo struct {
	db *gorm.DB
}

// NewEmbeddingRepo creates a new EmbeddingRepo.
func NewEmbeddingRepo(db *gorm.DB) *EmbeddingRepo {
	return &EmbeddingRepo{db: db}
}

func (r *EmbeddingRepo) GetByEntity(ctx context.Context, entityType domain.EntityType, entityID string) (*domain.Embedding, error) {
	var model EmbeddingModel
	if err := r.db.WithContext(ctx).Where("entity_type = ? AND entity_id = ?", string(entityType), entityID).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *EmbeddingRepo) ListByType(ctx context.Context, entityType domain.EntityType, page, pageSize int) ([]domain.Embedding, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&EmbeddingModel{}).Where("entity_type = ?", string(entityType)).Count(&total)

	var models []EmbeddingModel
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Where("entity_type = ?", string(entityType)).
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	embeddings := make([]domain.Embedding, len(models))
	for i, m := range models {
		embeddings[i] = *m.ToDomain()
	}
	return embeddings, total, nil
}

func (r *EmbeddingRepo) Create(ctx context.Context, embedding *domain.Embedding) error {
	model := ToEmbeddingModel(embedding)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *EmbeddingRepo) Update(ctx context.Context, embedding *domain.Embedding) error {
	return r.db.WithContext(ctx).Model(&EmbeddingModel{}).Where("id = ?", embedding.ID).Updates(map[string]interface{}{
		"embedding_vector": Float64Array(embedding.EmbeddingVector),
		"model_version":    embedding.ModelVersion,
		"dimensions":       embedding.Dimensions,
	}).Error
}

func (r *EmbeddingRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&EmbeddingModel{}).Error
}
