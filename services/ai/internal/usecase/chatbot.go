package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/southern-martin/ecommerce/services/ai/internal/domain"
	"github.com/southern-martin/ecommerce/services/ai/internal/infrastructure/aiclient"
)

// ChatbotUseCase handles AI chatbot operations.
type ChatbotUseCase struct {
	conversationRepo domain.AIConversationRepository
	aiClient         *aiclient.MockAIClient
}

// NewChatbotUseCase creates a new ChatbotUseCase.
func NewChatbotUseCase(
	conversationRepo domain.AIConversationRepository,
	aiClient *aiclient.MockAIClient,
) *ChatbotUseCase {
	return &ChatbotUseCase{
		conversationRepo: conversationRepo,
		aiClient:         aiClient,
	}
}

// ChatRequest is the input for a chat message.
type ChatRequest struct {
	UserID         string
	ConversationID string // optional, creates new if empty
	Message        string
}

// ChatResponse is the output of a chat message.
type ChatResponse struct {
	ConversationID string `json:"conversation_id"`
	Response       string `json:"response"`
	Model          string `json:"model"`
	TokenCount     int    `json:"token_count"`
}

// Chat sends a message and returns the AI response.
func (uc *ChatbotUseCase) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	var conversation *domain.AIConversation
	var messages []domain.ChatMessage

	if req.ConversationID != "" {
		// Load existing conversation
		existing, err := uc.conversationRepo.GetByID(ctx, req.ConversationID)
		if err != nil {
			return nil, fmt.Errorf("conversation not found: %w", err)
		}
		conversation = existing

		if err := json.Unmarshal([]byte(conversation.MessagesJSON), &messages); err != nil {
			messages = []domain.ChatMessage{}
		}
	} else {
		// Create new conversation
		conversation = &domain.AIConversation{
			ID:     uuid.New().String(),
			UserID: req.UserID,
			Title:  truncate(req.Message, 100),
			Model:  "gpt-4",
		}
	}

	// Add user message
	messages = append(messages, domain.ChatMessage{
		Role:    "user",
		Content: req.Message,
	})

	// Convert to AI client messages
	clientMessages := make([]aiclient.ChatMessage, len(messages))
	for i, m := range messages {
		clientMessages[i] = aiclient.ChatMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}

	// Get AI response
	response, err := uc.aiClient.Chat(clientMessages)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI response: %w", err)
	}

	// Add assistant response
	messages = append(messages, domain.ChatMessage{
		Role:    "assistant",
		Content: response,
	})

	// Serialize messages
	messagesJSON, err := json.Marshal(messages)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize messages: %w", err)
	}

	// Mock token count
	tokenCount := len(req.Message)/4 + len(response)/4

	conversation.MessagesJSON = string(messagesJSON)
	conversation.TokenCount = tokenCount

	if req.ConversationID != "" {
		if err := uc.conversationRepo.Update(ctx, conversation); err != nil {
			return nil, fmt.Errorf("failed to update conversation: %w", err)
		}
	} else {
		if err := uc.conversationRepo.Create(ctx, conversation); err != nil {
			return nil, fmt.Errorf("failed to create conversation: %w", err)
		}
	}

	return &ChatResponse{
		ConversationID: conversation.ID,
		Response:       response,
		Model:          conversation.Model,
		TokenCount:     tokenCount,
	}, nil
}

// GetConversation retrieves a conversation by ID.
func (uc *ChatbotUseCase) GetConversation(ctx context.Context, id string) (*domain.AIConversation, error) {
	return uc.conversationRepo.GetByID(ctx, id)
}

// ListConversations lists conversations for a user.
func (uc *ChatbotUseCase) ListConversations(ctx context.Context, userID string, page, pageSize int) ([]domain.AIConversation, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return uc.conversationRepo.ListByUser(ctx, userID, page, pageSize)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
