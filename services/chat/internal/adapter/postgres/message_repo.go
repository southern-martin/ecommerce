package postgres

import (
	"context"
	"time"

	"github.com/southern-martin/ecommerce/services/chat/internal/domain"
	"gorm.io/gorm"
)

// MessageRepo implements domain.MessageRepository.
type MessageRepo struct {
	db *gorm.DB
}

// NewMessageRepo creates a new MessageRepo.
func NewMessageRepo(db *gorm.DB) *MessageRepo {
	return &MessageRepo{db: db}
}

func (r *MessageRepo) GetByID(ctx context.Context, id string) (*domain.Message, error) {
	var model MessageModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *MessageRepo) ListByConversation(ctx context.Context, conversationID string, page, pageSize int) ([]domain.Message, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&MessageModel{}).Where("conversation_id = ?", conversationID).Count(&total)

	var models []MessageModel
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Where("conversation_id = ?", conversationID).
		Order("created_at ASC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	messages := make([]domain.Message, len(models))
	for i, m := range models {
		messages[i] = *m.ToDomain()
	}
	return messages, total, nil
}

func (r *MessageRepo) Create(ctx context.Context, message *domain.Message) error {
	model := ToMessageModel(message)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *MessageRepo) MarkAsRead(ctx context.Context, conversationID, userID string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&MessageModel{}).
		Where("conversation_id = ? AND sender_id != ? AND is_read = false", conversationID, userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

func (r *MessageRepo) CountUnread(ctx context.Context, conversationID, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&MessageModel{}).
		Where("conversation_id = ? AND sender_id != ? AND is_read = false", conversationID, userID).
		Count(&count).Error
	return count, err
}
