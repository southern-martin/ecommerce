package domain

import (
	"context"
	"time"
)

// ConversationRepository defines the interface for conversation persistence.
type ConversationRepository interface {
	GetByID(ctx context.Context, id string) (*Conversation, error)
	ListByUser(ctx context.Context, userID string, status string, page, pageSize int) ([]Conversation, int64, error)
	ListByParticipants(ctx context.Context, participantIDs []string) ([]Conversation, error)
	Create(ctx context.Context, conversation *Conversation) error
	Update(ctx context.Context, conversation *Conversation) error
	UpdateLastMessage(ctx context.Context, id string, lastMessageAt *time.Time) error
}

// MessageRepository defines the interface for message persistence.
type MessageRepository interface {
	GetByID(ctx context.Context, id string) (*Message, error)
	ListByConversation(ctx context.Context, conversationID string, page, pageSize int) ([]Message, int64, error)
	Create(ctx context.Context, message *Message) error
	MarkAsRead(ctx context.Context, conversationID, userID string) error
	CountUnread(ctx context.Context, conversationID, userID string) (int64, error)
}

// ParticipantRepository defines the interface for conversation participant persistence.
type ParticipantRepository interface {
	GetByConversationAndUser(ctx context.Context, conversationID, userID string) (*ConversationParticipant, error)
	ListByConversation(ctx context.Context, conversationID string) ([]ConversationParticipant, error)
	Create(ctx context.Context, participant *ConversationParticipant) error
	UpdateLastRead(ctx context.Context, conversationID, userID string) error
}
