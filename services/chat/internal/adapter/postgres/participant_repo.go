package postgres

import (
	"context"
	"time"

	"github.com/southern-martin/ecommerce/services/chat/internal/domain"
	"gorm.io/gorm"
)

// ParticipantRepo implements domain.ParticipantRepository.
type ParticipantRepo struct {
	db *gorm.DB
}

// NewParticipantRepo creates a new ParticipantRepo.
func NewParticipantRepo(db *gorm.DB) *ParticipantRepo {
	return &ParticipantRepo{db: db}
}

func (r *ParticipantRepo) GetByConversationAndUser(ctx context.Context, conversationID, userID string) (*domain.ConversationParticipant, error) {
	var model ConversationParticipantModel
	if err := r.db.WithContext(ctx).Where("conversation_id = ? AND user_id = ?", conversationID, userID).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *ParticipantRepo) ListByConversation(ctx context.Context, conversationID string) ([]domain.ConversationParticipant, error) {
	var models []ConversationParticipantModel
	if err := r.db.WithContext(ctx).Where("conversation_id = ?", conversationID).Find(&models).Error; err != nil {
		return nil, err
	}

	participants := make([]domain.ConversationParticipant, len(models))
	for i, m := range models {
		participants[i] = *m.ToDomain()
	}
	return participants, nil
}

func (r *ParticipantRepo) Create(ctx context.Context, participant *domain.ConversationParticipant) error {
	model := ToConversationParticipantModel(participant)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *ParticipantRepo) UpdateLastRead(ctx context.Context, conversationID, userID string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&ConversationParticipantModel{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Update("last_read_at", now).Error
}
