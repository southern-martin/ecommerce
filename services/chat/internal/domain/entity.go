package domain

import "time"

// ConversationType represents the type of a conversation.
type ConversationType string

const (
	ConversationTypeBuyerSeller ConversationType = "buyer_seller"
	ConversationTypeSupport     ConversationType = "support"
)

// ConversationStatus represents the status of a conversation.
type ConversationStatus string

const (
	ConversationStatusActive   ConversationStatus = "active"
	ConversationStatusArchived ConversationStatus = "archived"
	ConversationStatusClosed   ConversationStatus = "closed"
)

// SenderRole represents the role of a message sender.
type SenderRole string

const (
	SenderRoleBuyer  SenderRole = "buyer"
	SenderRoleSeller SenderRole = "seller"
	SenderRoleAdmin  SenderRole = "admin"
	SenderRoleSystem SenderRole = "system"
)

// MessageType represents the type of a message.
type MessageType string

const (
	MessageTypeText   MessageType = "text"
	MessageTypeImage  MessageType = "image"
	MessageTypeSystem MessageType = "system"
)

// ParticipantRole represents a participant's role in a conversation.
type ParticipantRole string

const (
	ParticipantRoleBuyer  ParticipantRole = "buyer"
	ParticipantRoleSeller ParticipantRole = "seller"
	ParticipantRoleAdmin  ParticipantRole = "admin"
)

// Conversation represents a chat conversation.
type Conversation struct {
	ID             string
	Type           ConversationType
	ParticipantIDs []string
	BuyerID        string
	SellerID       string
	OrderID        string
	Subject        string
	Status         ConversationStatus
	LastMessageAt  *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Message represents a chat message.
type Message struct {
	ID             string
	ConversationID string
	SenderID       string
	SenderRole     SenderRole
	Content        string
	MessageType    MessageType
	Attachments    []string
	IsRead         bool
	ReadAt         *time.Time
	CreatedAt      time.Time
}

// ConversationParticipant represents a participant in a conversation.
type ConversationParticipant struct {
	ID             string
	ConversationID string
	UserID         string
	Role           ParticipantRole
	JoinedAt       time.Time
	LastReadAt     *time.Time
}
