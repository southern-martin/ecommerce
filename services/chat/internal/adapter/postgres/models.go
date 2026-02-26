package postgres

import (
	"time"

	"github.com/lib/pq"
	"github.com/southern-martin/ecommerce/services/chat/internal/domain"
)

// ConversationModel is the GORM model for the conversations table.
type ConversationModel struct {
	ID             string         `gorm:"type:uuid;primaryKey"`
	Type           string         `gorm:"type:varchar(20);not null;default:'buyer_seller'"`
	ParticipantIDs pq.StringArray `gorm:"type:text[]"`
	BuyerID        string         `gorm:"type:uuid;index;not null"`
	SellerID       string         `gorm:"type:uuid;index;not null"`
	OrderID        string         `gorm:"type:uuid;index"`
	Subject        string         `gorm:"type:varchar(500)"`
	Status         string         `gorm:"type:varchar(20);default:'active';index"`
	LastMessageAt  *time.Time
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func (ConversationModel) TableName() string { return "conversations" }

func (m *ConversationModel) ToDomain() *domain.Conversation {
	return &domain.Conversation{
		ID:             m.ID,
		Type:           domain.ConversationType(m.Type),
		ParticipantIDs: []string(m.ParticipantIDs),
		BuyerID:        m.BuyerID,
		SellerID:       m.SellerID,
		OrderID:        m.OrderID,
		Subject:        m.Subject,
		Status:         domain.ConversationStatus(m.Status),
		LastMessageAt:  m.LastMessageAt,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func ToConversationModel(c *domain.Conversation) *ConversationModel {
	return &ConversationModel{
		ID:             c.ID,
		Type:           string(c.Type),
		ParticipantIDs: pq.StringArray(c.ParticipantIDs),
		BuyerID:        c.BuyerID,
		SellerID:       c.SellerID,
		OrderID:        c.OrderID,
		Subject:        c.Subject,
		Status:         string(c.Status),
		LastMessageAt:  c.LastMessageAt,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
}

// MessageModel is the GORM model for the messages table.
type MessageModel struct {
	ID             string         `gorm:"type:uuid;primaryKey"`
	ConversationID string         `gorm:"type:uuid;index;not null"`
	SenderID       string         `gorm:"type:uuid;not null"`
	SenderRole     string         `gorm:"type:varchar(20);not null"`
	Content        string         `gorm:"type:text;not null"`
	MessageType    string         `gorm:"type:varchar(20);default:'text'"`
	Attachments    pq.StringArray `gorm:"type:text[]"`
	IsRead         bool           `gorm:"default:false"`
	ReadAt         *time.Time
	CreatedAt      time.Time `gorm:"autoCreateTime"`
}

func (MessageModel) TableName() string { return "messages" }

func (m *MessageModel) ToDomain() *domain.Message {
	return &domain.Message{
		ID:             m.ID,
		ConversationID: m.ConversationID,
		SenderID:       m.SenderID,
		SenderRole:     domain.SenderRole(m.SenderRole),
		Content:        m.Content,
		MessageType:    domain.MessageType(m.MessageType),
		Attachments:    []string(m.Attachments),
		IsRead:         m.IsRead,
		ReadAt:         m.ReadAt,
		CreatedAt:      m.CreatedAt,
	}
}

func ToMessageModel(msg *domain.Message) *MessageModel {
	return &MessageModel{
		ID:             msg.ID,
		ConversationID: msg.ConversationID,
		SenderID:       msg.SenderID,
		SenderRole:     string(msg.SenderRole),
		Content:        msg.Content,
		MessageType:    string(msg.MessageType),
		Attachments:    pq.StringArray(msg.Attachments),
		IsRead:         msg.IsRead,
		ReadAt:         msg.ReadAt,
		CreatedAt:      msg.CreatedAt,
	}
}

// ConversationParticipantModel is the GORM model for the conversation_participants table.
type ConversationParticipantModel struct {
	ID             string `gorm:"type:uuid;primaryKey"`
	ConversationID string `gorm:"type:uuid;index;not null;uniqueIndex:idx_conv_user"`
	UserID         string `gorm:"type:uuid;index;not null;uniqueIndex:idx_conv_user"`
	Role           string `gorm:"type:varchar(20);not null"`
	JoinedAt       time.Time
	LastReadAt     *time.Time
}

func (ConversationParticipantModel) TableName() string { return "conversation_participants" }

func (m *ConversationParticipantModel) ToDomain() *domain.ConversationParticipant {
	return &domain.ConversationParticipant{
		ID:             m.ID,
		ConversationID: m.ConversationID,
		UserID:         m.UserID,
		Role:           domain.ParticipantRole(m.Role),
		JoinedAt:       m.JoinedAt,
		LastReadAt:     m.LastReadAt,
	}
}

func ToConversationParticipantModel(p *domain.ConversationParticipant) *ConversationParticipantModel {
	return &ConversationParticipantModel{
		ID:             p.ID,
		ConversationID: p.ConversationID,
		UserID:         p.UserID,
		Role:           string(p.Role),
		JoinedAt:       p.JoinedAt,
		LastReadAt:     p.LastReadAt,
	}
}
