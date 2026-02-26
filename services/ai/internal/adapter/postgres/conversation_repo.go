package postgres

import (
	"context"

	"github.com/southern-martin/ecommerce/services/ai/internal/domain"
	"gorm.io/gorm"
)

// AIConversationRepo implements domain.AIConversationRepository.
type AIConversationRepo struct {
	db *gorm.DB
}

// NewAIConversationRepo creates a new AIConversationRepo.
func NewAIConversationRepo(db *gorm.DB) *AIConversationRepo {
	return &AIConversationRepo{db: db}
}

func (r *AIConversationRepo) GetByID(ctx context.Context, id string) (*domain.AIConversation, error) {
	var model AIConversationModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *AIConversationRepo) ListByUser(ctx context.Context, userID string, page, pageSize int) ([]domain.AIConversation, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&AIConversationModel{}).Where("user_id = ?", userID).Count(&total)

	var models []AIConversationModel
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	conversations := make([]domain.AIConversation, len(models))
	for i, m := range models {
		conversations[i] = *m.ToDomain()
	}
	return conversations, total, nil
}

func (r *AIConversationRepo) Create(ctx context.Context, conversation *domain.AIConversation) error {
	model := ToAIConversationModel(conversation)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *AIConversationRepo) Update(ctx context.Context, conversation *domain.AIConversation) error {
	return r.db.WithContext(ctx).Model(&AIConversationModel{}).Where("id = ?", conversation.ID).Updates(map[string]interface{}{
		"title":         conversation.Title,
		"messages_json": conversation.MessagesJSON,
		"model":         conversation.Model,
		"token_count":   conversation.TokenCount,
	}).Error
}
