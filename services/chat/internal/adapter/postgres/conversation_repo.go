package postgres

import (
	"context"
	"time"

	"github.com/southern-martin/ecommerce/services/chat/internal/domain"
	"gorm.io/gorm"
)

// ConversationRepo implements domain.ConversationRepository.
type ConversationRepo struct {
	db *gorm.DB
}

// NewConversationRepo creates a new ConversationRepo.
func NewConversationRepo(db *gorm.DB) *ConversationRepo {
	return &ConversationRepo{db: db}
}

func (r *ConversationRepo) GetByID(ctx context.Context, id string) (*domain.Conversation, error) {
	var model ConversationModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *ConversationRepo) ListByUser(ctx context.Context, userID string, status string, page, pageSize int) ([]domain.Conversation, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&ConversationModel{}).Where("? = ANY(participant_ids)", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	var models []ConversationModel
	offset := (page - 1) * pageSize
	if err := query.Order("COALESCE(last_message_at, created_at) DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	conversations := make([]domain.Conversation, len(models))
	for i, m := range models {
		conversations[i] = *m.ToDomain()
	}
	return conversations, total, nil
}

func (r *ConversationRepo) ListByParticipants(ctx context.Context, participantIDs []string) ([]domain.Conversation, error) {
	var models []ConversationModel
	query := r.db.WithContext(ctx)
	for _, pid := range participantIDs {
		query = query.Where("? = ANY(participant_ids)", pid)
	}
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	conversations := make([]domain.Conversation, len(models))
	for i, m := range models {
		conversations[i] = *m.ToDomain()
	}
	return conversations, nil
}

func (r *ConversationRepo) Create(ctx context.Context, conversation *domain.Conversation) error {
	model := ToConversationModel(conversation)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *ConversationRepo) Update(ctx context.Context, conversation *domain.Conversation) error {
	return r.db.WithContext(ctx).Model(&ConversationModel{}).Where("id = ?", conversation.ID).Updates(map[string]interface{}{
		"subject":         conversation.Subject,
		"status":          string(conversation.Status),
		"last_message_at": conversation.LastMessageAt,
	}).Error
}

func (r *ConversationRepo) UpdateLastMessage(ctx context.Context, id string, lastMessageAt *time.Time) error {
	return r.db.WithContext(ctx).Model(&ConversationModel{}).Where("id = ?", id).Update("last_message_at", lastMessageAt).Error
}
