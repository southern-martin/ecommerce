package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/ai/internal/domain"
	"github.com/southern-martin/ecommerce/services/ai/internal/infrastructure/aiclient"
)

// ---------------------------------------------------------------------------
// Hand-written function-field mocks
// ---------------------------------------------------------------------------

// --- AIConversationRepository mock ---

type mockAIConversationRepo struct {
	getByIDFn    func(ctx context.Context, id string) (*domain.AIConversation, error)
	listByUserFn func(ctx context.Context, userID string, page, pageSize int) ([]domain.AIConversation, int64, error)
	createFn     func(ctx context.Context, conversation *domain.AIConversation) error
	updateFn     func(ctx context.Context, conversation *domain.AIConversation) error
}

func (m *mockAIConversationRepo) GetByID(ctx context.Context, id string) (*domain.AIConversation, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}
func (m *mockAIConversationRepo) ListByUser(ctx context.Context, userID string, page, pageSize int) ([]domain.AIConversation, int64, error) {
	if m.listByUserFn != nil {
		return m.listByUserFn(ctx, userID, page, pageSize)
	}
	return nil, 0, nil
}
func (m *mockAIConversationRepo) Create(ctx context.Context, conversation *domain.AIConversation) error {
	if m.createFn != nil {
		return m.createFn(ctx, conversation)
	}
	return nil
}
func (m *mockAIConversationRepo) Update(ctx context.Context, conversation *domain.AIConversation) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, conversation)
	}
	return nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func newChatbotUseCase(
	convRepo *mockAIConversationRepo,
	aiClient *aiclient.MockAIClient,
) *ChatbotUseCase {
	return NewChatbotUseCase(convRepo, aiClient)
}

// ===========================================================================
// Chat tests
// ===========================================================================

func TestChat_NewConversation_Success(t *testing.T) {
	convRepo := &mockAIConversationRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")

	var savedConv *domain.AIConversation
	convRepo.createFn = func(_ context.Context, c *domain.AIConversation) error {
		savedConv = c
		return nil
	}

	uc := newChatbotUseCase(convRepo, client)
	resp, err := uc.Chat(context.Background(), ChatRequest{
		UserID:  "user-1",
		Message: "What are the best running shoes?",
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.ConversationID)
	assert.NotEmpty(t, resp.Response)
	assert.Equal(t, "gpt-4", resp.Model)
	assert.Greater(t, resp.TokenCount, 0)
	assert.NotNil(t, savedConv)
	assert.Equal(t, "user-1", savedConv.UserID)
}

func TestChat_NewConversation_TitleTruncated(t *testing.T) {
	convRepo := &mockAIConversationRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")

	var savedConv *domain.AIConversation
	convRepo.createFn = func(_ context.Context, c *domain.AIConversation) error {
		savedConv = c
		return nil
	}

	// Create a message longer than 100 characters
	longMessage := ""
	for i := 0; i < 120; i++ {
		longMessage += "a"
	}

	uc := newChatbotUseCase(convRepo, client)
	_, err := uc.Chat(context.Background(), ChatRequest{
		UserID:  "user-1",
		Message: longMessage,
	})

	require.NoError(t, err)
	assert.Len(t, savedConv.Title, 100)
}

func TestChat_ExistingConversation_Success(t *testing.T) {
	convRepo := &mockAIConversationRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")

	existingMessages := []domain.ChatMessage{
		{Role: "user", Content: "Hello"},
		{Role: "assistant", Content: "Hi there!"},
	}
	messagesJSON, _ := json.Marshal(existingMessages)

	convRepo.getByIDFn = func(_ context.Context, id string) (*domain.AIConversation, error) {
		return &domain.AIConversation{
			ID:           "conv-1",
			UserID:       "user-1",
			Title:        "Hello",
			MessagesJSON: string(messagesJSON),
			Model:        "gpt-4",
		}, nil
	}

	var updatedConv *domain.AIConversation
	convRepo.updateFn = func(_ context.Context, c *domain.AIConversation) error {
		updatedConv = c
		return nil
	}

	uc := newChatbotUseCase(convRepo, client)
	resp, err := uc.Chat(context.Background(), ChatRequest{
		UserID:         "user-1",
		ConversationID: "conv-1",
		Message:        "Tell me about shoes",
	})

	require.NoError(t, err)
	assert.Equal(t, "conv-1", resp.ConversationID)
	assert.NotEmpty(t, resp.Response)
	assert.NotNil(t, updatedConv)

	// Verify messages were appended
	var updatedMessages []domain.ChatMessage
	err = json.Unmarshal([]byte(updatedConv.MessagesJSON), &updatedMessages)
	require.NoError(t, err)
	assert.Len(t, updatedMessages, 4) // 2 existing + 1 user + 1 assistant
	assert.Equal(t, "user", updatedMessages[2].Role)
	assert.Equal(t, "Tell me about shoes", updatedMessages[2].Content)
	assert.Equal(t, "assistant", updatedMessages[3].Role)
}

func TestChat_ExistingConversation_NotFound(t *testing.T) {
	convRepo := &mockAIConversationRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")

	convRepo.getByIDFn = func(_ context.Context, _ string) (*domain.AIConversation, error) {
		return nil, errors.New("not found")
	}

	uc := newChatbotUseCase(convRepo, client)
	_, err := uc.Chat(context.Background(), ChatRequest{
		UserID:         "user-1",
		ConversationID: "nonexistent",
		Message:        "Hello",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "conversation not found")
}

func TestChat_CreateRepoError(t *testing.T) {
	convRepo := &mockAIConversationRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")

	convRepo.createFn = func(_ context.Context, _ *domain.AIConversation) error {
		return errors.New("db write failed")
	}

	uc := newChatbotUseCase(convRepo, client)
	_, err := uc.Chat(context.Background(), ChatRequest{
		UserID:  "user-1",
		Message: "Hello",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create conversation")
}

func TestChat_UpdateRepoError(t *testing.T) {
	convRepo := &mockAIConversationRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")

	emptyMessages, _ := json.Marshal([]domain.ChatMessage{})
	convRepo.getByIDFn = func(_ context.Context, _ string) (*domain.AIConversation, error) {
		return &domain.AIConversation{
			ID:           "conv-1",
			MessagesJSON: string(emptyMessages),
			Model:        "gpt-4",
		}, nil
	}
	convRepo.updateFn = func(_ context.Context, _ *domain.AIConversation) error {
		return errors.New("db update failed")
	}

	uc := newChatbotUseCase(convRepo, client)
	_, err := uc.Chat(context.Background(), ChatRequest{
		UserID:         "user-1",
		ConversationID: "conv-1",
		Message:        "Hello",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update conversation")
}

func TestChat_ExistingConversation_InvalidJSON(t *testing.T) {
	convRepo := &mockAIConversationRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")

	convRepo.getByIDFn = func(_ context.Context, _ string) (*domain.AIConversation, error) {
		return &domain.AIConversation{
			ID:           "conv-1",
			MessagesJSON: "invalid-json",
			Model:        "gpt-4",
		}, nil
	}
	convRepo.updateFn = func(_ context.Context, _ *domain.AIConversation) error { return nil }

	uc := newChatbotUseCase(convRepo, client)
	resp, err := uc.Chat(context.Background(), ChatRequest{
		UserID:         "user-1",
		ConversationID: "conv-1",
		Message:        "Hello",
	})

	// Should succeed - invalid JSON is handled gracefully by resetting messages
	require.NoError(t, err)
	assert.NotEmpty(t, resp.Response)
}

// ===========================================================================
// GetConversation tests
// ===========================================================================

func TestGetConversation_Success(t *testing.T) {
	convRepo := &mockAIConversationRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")

	expected := &domain.AIConversation{
		ID:     "conv-1",
		UserID: "user-1",
		Title:  "Test Conversation",
		Model:  "gpt-4",
	}
	convRepo.getByIDFn = func(_ context.Context, id string) (*domain.AIConversation, error) {
		assert.Equal(t, "conv-1", id)
		return expected, nil
	}

	uc := newChatbotUseCase(convRepo, client)
	conv, err := uc.GetConversation(context.Background(), "conv-1")

	require.NoError(t, err)
	assert.Equal(t, expected, conv)
}

func TestGetConversation_NotFound(t *testing.T) {
	convRepo := &mockAIConversationRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")

	uc := newChatbotUseCase(convRepo, client)
	_, err := uc.GetConversation(context.Background(), "nonexistent")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// ===========================================================================
// ListConversations tests
// ===========================================================================

func TestListConversations_Success(t *testing.T) {
	convRepo := &mockAIConversationRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")

	expectedConvs := []domain.AIConversation{
		{ID: "conv-1", UserID: "user-1"},
		{ID: "conv-2", UserID: "user-1"},
	}
	convRepo.listByUserFn = func(_ context.Context, userID string, page, pageSize int) ([]domain.AIConversation, int64, error) {
		assert.Equal(t, "user-1", userID)
		assert.Equal(t, 1, page)
		assert.Equal(t, 20, pageSize)
		return expectedConvs, 2, nil
	}

	uc := newChatbotUseCase(convRepo, client)
	convs, total, err := uc.ListConversations(context.Background(), "user-1", 1, 20)

	require.NoError(t, err)
	assert.Len(t, convs, 2)
	assert.Equal(t, int64(2), total)
}

func TestListConversations_DefaultPagination(t *testing.T) {
	convRepo := &mockAIConversationRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")

	var capturedPage, capturedPageSize int
	convRepo.listByUserFn = func(_ context.Context, _ string, page, pageSize int) ([]domain.AIConversation, int64, error) {
		capturedPage = page
		capturedPageSize = pageSize
		return nil, 0, nil
	}

	uc := newChatbotUseCase(convRepo, client)

	// page < 1 should default to 1
	_, _, _ = uc.ListConversations(context.Background(), "user-1", 0, 20)
	assert.Equal(t, 1, capturedPage)

	// pageSize > 100 should default to 20
	_, _, _ = uc.ListConversations(context.Background(), "user-1", 1, 200)
	assert.Equal(t, 20, capturedPageSize)

	// pageSize < 1 should default to 20
	_, _, _ = uc.ListConversations(context.Background(), "user-1", 1, -1)
	assert.Equal(t, 20, capturedPageSize)
}

// ===========================================================================
// truncate tests
// ===========================================================================

func TestTruncate(t *testing.T) {
	assert.Equal(t, "hello", truncate("hello", 10))
	assert.Equal(t, "hel", truncate("hello", 3))
	assert.Equal(t, "", truncate("", 5))
	assert.Equal(t, "exact", truncate("exact", 5))
}
