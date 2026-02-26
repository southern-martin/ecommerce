package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/southern-martin/ecommerce/services/chat/internal/domain"
)

// ConversationUseCase handles conversation operations.
type ConversationUseCase struct {
	conversationRepo domain.ConversationRepository
	participantRepo  domain.ParticipantRepository
	publisher        domain.EventPublisher
}

// NewConversationUseCase creates a new ConversationUseCase.
func NewConversationUseCase(
	conversationRepo domain.ConversationRepository,
	participantRepo domain.ParticipantRepository,
	publisher domain.EventPublisher,
) *ConversationUseCase {
	return &ConversationUseCase{
		conversationRepo: conversationRepo,
		participantRepo:  participantRepo,
		publisher:        publisher,
	}
}

// CreateConversationRequest is the input for creating a conversation.
type CreateConversationRequest struct {
	Type     string `json:"type"`
	BuyerID  string `json:"buyer_id"`
	SellerID string `json:"seller_id"`
	OrderID  string `json:"order_id"`
	Subject  string `json:"subject"`
}

// CreateConversation creates a new conversation.
func (uc *ConversationUseCase) CreateConversation(ctx context.Context, userID string, req CreateConversationRequest) (*domain.Conversation, error) {
	convType := domain.ConversationType(req.Type)
	if convType == "" {
		convType = domain.ConversationTypeBuyerSeller
	}

	participantIDs := []string{req.BuyerID, req.SellerID}

	conversation := &domain.Conversation{
		ID:             uuid.New().String(),
		Type:           convType,
		ParticipantIDs: participantIDs,
		BuyerID:        req.BuyerID,
		SellerID:       req.SellerID,
		OrderID:        req.OrderID,
		Subject:        req.Subject,
		Status:         domain.ConversationStatusActive,
	}

	if err := uc.conversationRepo.Create(ctx, conversation); err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	// Create participant records
	buyerParticipant := &domain.ConversationParticipant{
		ID:             uuid.New().String(),
		ConversationID: conversation.ID,
		UserID:         req.BuyerID,
		Role:           domain.ParticipantRoleBuyer,
	}
	if err := uc.participantRepo.Create(ctx, buyerParticipant); err != nil {
		return nil, fmt.Errorf("failed to create buyer participant: %w", err)
	}

	sellerParticipant := &domain.ConversationParticipant{
		ID:             uuid.New().String(),
		ConversationID: conversation.ID,
		UserID:         req.SellerID,
		Role:           domain.ParticipantRoleSeller,
	}
	if err := uc.participantRepo.Create(ctx, sellerParticipant); err != nil {
		return nil, fmt.Errorf("failed to create seller participant: %w", err)
	}

	// Publish event
	_ = uc.publisher.Publish(ctx, "chat.conversation.new", map[string]interface{}{
		"conversation_id": conversation.ID,
		"type":            string(conversation.Type),
		"buyer_id":        conversation.BuyerID,
		"seller_id":       conversation.SellerID,
		"order_id":        conversation.OrderID,
	})

	return conversation, nil
}

// GetConversation retrieves a conversation by ID.
func (uc *ConversationUseCase) GetConversation(ctx context.Context, id string) (*domain.Conversation, error) {
	return uc.conversationRepo.GetByID(ctx, id)
}

// ListUserConversations lists conversations for a user with pagination.
func (uc *ConversationUseCase) ListUserConversations(ctx context.Context, userID string, status string, page, pageSize int) ([]domain.Conversation, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return uc.conversationRepo.ListByUser(ctx, userID, status, page, pageSize)
}

// ArchiveConversation archives a conversation.
func (uc *ConversationUseCase) ArchiveConversation(ctx context.Context, id string) error {
	conv, err := uc.conversationRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("conversation not found: %w", err)
	}

	conv.Status = domain.ConversationStatusArchived
	return uc.conversationRepo.Update(ctx, conv)
}
