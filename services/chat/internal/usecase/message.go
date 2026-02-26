package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/southern-martin/ecommerce/services/chat/internal/domain"
)

// MessageUseCase handles message operations.
type MessageUseCase struct {
	messageRepo      domain.MessageRepository
	conversationRepo domain.ConversationRepository
	participantRepo  domain.ParticipantRepository
	publisher        domain.EventPublisher
}

// NewMessageUseCase creates a new MessageUseCase.
func NewMessageUseCase(
	messageRepo domain.MessageRepository,
	conversationRepo domain.ConversationRepository,
	participantRepo domain.ParticipantRepository,
	publisher domain.EventPublisher,
) *MessageUseCase {
	return &MessageUseCase{
		messageRepo:      messageRepo,
		conversationRepo: conversationRepo,
		participantRepo:  participantRepo,
		publisher:        publisher,
	}
}

// SendMessageRequest is the input for sending a message.
type SendMessageRequest struct {
	ConversationID string   `json:"conversation_id"`
	SenderID       string   `json:"sender_id"`
	SenderRole     string   `json:"sender_role"`
	Content        string   `json:"content"`
	MessageType    string   `json:"message_type"`
	Attachments    []string `json:"attachments"`
}

// SendMessage sends a message in a conversation.
func (uc *MessageUseCase) SendMessage(ctx context.Context, req SendMessageRequest) (*domain.Message, error) {
	// Validate participant belongs to conversation
	_, err := uc.participantRepo.GetByConversationAndUser(ctx, req.ConversationID, req.SenderID)
	if err != nil {
		return nil, fmt.Errorf("sender is not a participant of this conversation: %w", err)
	}

	msgType := domain.MessageType(req.MessageType)
	if msgType == "" {
		msgType = domain.MessageTypeText
	}

	senderRole := domain.SenderRole(req.SenderRole)
	if senderRole == "" {
		senderRole = domain.SenderRoleBuyer
	}

	message := &domain.Message{
		ID:             uuid.New().String(),
		ConversationID: req.ConversationID,
		SenderID:       req.SenderID,
		SenderRole:     senderRole,
		Content:        req.Content,
		MessageType:    msgType,
		Attachments:    req.Attachments,
		IsRead:         false,
	}

	if err := uc.messageRepo.Create(ctx, message); err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	// Update conversation last_message_at
	now := time.Now()
	_ = uc.conversationRepo.UpdateLastMessage(ctx, req.ConversationID, &now)

	// Publish event
	_ = uc.publisher.Publish(ctx, "chat.message.sent", map[string]interface{}{
		"message_id":      message.ID,
		"conversation_id": message.ConversationID,
		"sender_id":       message.SenderID,
		"sender_role":     string(message.SenderRole),
		"content":         message.Content,
		"message_type":    string(message.MessageType),
	})

	return message, nil
}

// ListMessages lists messages in a conversation with pagination.
func (uc *MessageUseCase) ListMessages(ctx context.Context, conversationID string, page, pageSize int) ([]domain.Message, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}
	return uc.messageRepo.ListByConversation(ctx, conversationID, page, pageSize)
}

// MarkAsRead marks messages in a conversation as read for a user.
func (uc *MessageUseCase) MarkAsRead(ctx context.Context, conversationID, userID string) error {
	if err := uc.messageRepo.MarkAsRead(ctx, conversationID, userID); err != nil {
		return fmt.Errorf("failed to mark messages as read: %w", err)
	}

	// Also update participant's last_read_at
	_ = uc.participantRepo.UpdateLastRead(ctx, conversationID, userID)

	return nil
}

// GetUnreadCount returns the unread message count for a user in a conversation.
func (uc *MessageUseCase) GetUnreadCount(ctx context.Context, conversationID, userID string) (int64, error) {
	return uc.messageRepo.CountUnread(ctx, conversationID, userID)
}
