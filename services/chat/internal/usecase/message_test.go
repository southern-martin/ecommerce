package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/chat/internal/domain"
)

// ---------------------------------------------------------------------------
// Hand-written function-field mocks (MessageRepository)
// ---------------------------------------------------------------------------

type mockMessageRepo struct {
	getByIDFn            func(ctx context.Context, id string) (*domain.Message, error)
	listByConversationFn func(ctx context.Context, conversationID string, page, pageSize int) ([]domain.Message, int64, error)
	createFn             func(ctx context.Context, message *domain.Message) error
	markAsReadFn         func(ctx context.Context, conversationID, userID string) error
	countUnreadFn        func(ctx context.Context, conversationID, userID string) (int64, error)
}

func (m *mockMessageRepo) GetByID(ctx context.Context, id string) (*domain.Message, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}
func (m *mockMessageRepo) ListByConversation(ctx context.Context, conversationID string, page, pageSize int) ([]domain.Message, int64, error) {
	if m.listByConversationFn != nil {
		return m.listByConversationFn(ctx, conversationID, page, pageSize)
	}
	return nil, 0, nil
}
func (m *mockMessageRepo) Create(ctx context.Context, message *domain.Message) error {
	if m.createFn != nil {
		return m.createFn(ctx, message)
	}
	return nil
}
func (m *mockMessageRepo) MarkAsRead(ctx context.Context, conversationID, userID string) error {
	if m.markAsReadFn != nil {
		return m.markAsReadFn(ctx, conversationID, userID)
	}
	return nil
}
func (m *mockMessageRepo) CountUnread(ctx context.Context, conversationID, userID string) (int64, error) {
	if m.countUnreadFn != nil {
		return m.countUnreadFn(ctx, conversationID, userID)
	}
	return 0, nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func newMessageUseCase(
	msgRepo *mockMessageRepo,
	convRepo *mockConversationRepo,
	partRepo *mockParticipantRepo,
	pub *mockEventPublisher,
) *MessageUseCase {
	return NewMessageUseCase(msgRepo, convRepo, partRepo, pub)
}

func defaultMessageMocks() (*mockMessageRepo, *mockConversationRepo, *mockParticipantRepo, *mockEventPublisher) {
	return &mockMessageRepo{}, &mockConversationRepo{}, &mockParticipantRepo{}, &mockEventPublisher{}
}

// ===========================================================================
// SendMessage tests
// ===========================================================================

func TestSendMessage_Success(t *testing.T) {
	msgRepo, convRepo, partRepo, pub := defaultMessageMocks()

	partRepo.getByConversationAndUserFn = func(_ context.Context, _, _ string) (*domain.ConversationParticipant, error) {
		return &domain.ConversationParticipant{ID: "part-1", UserID: "user-1"}, nil
	}

	var savedMsg *domain.Message
	msgRepo.createFn = func(_ context.Context, m *domain.Message) error {
		savedMsg = m
		return nil
	}

	uc := newMessageUseCase(msgRepo, convRepo, partRepo, pub)
	msg, err := uc.SendMessage(context.Background(), SendMessageRequest{
		ConversationID: "conv-1",
		SenderID:       "user-1",
		SenderRole:     "buyer",
		Content:        "Hello there!",
		MessageType:    "text",
	})

	require.NoError(t, err)
	require.NotNil(t, msg)
	assert.NotEmpty(t, msg.ID)
	assert.Equal(t, "conv-1", msg.ConversationID)
	assert.Equal(t, "user-1", msg.SenderID)
	assert.Equal(t, domain.SenderRoleBuyer, msg.SenderRole)
	assert.Equal(t, "Hello there!", msg.Content)
	assert.Equal(t, domain.MessageTypeText, msg.MessageType)
	assert.False(t, msg.IsRead)
	assert.NotNil(t, savedMsg)
}

func TestSendMessage_DefaultMessageTypeText(t *testing.T) {
	msgRepo, convRepo, partRepo, pub := defaultMessageMocks()

	partRepo.getByConversationAndUserFn = func(_ context.Context, _, _ string) (*domain.ConversationParticipant, error) {
		return &domain.ConversationParticipant{}, nil
	}
	msgRepo.createFn = func(_ context.Context, _ *domain.Message) error { return nil }

	uc := newMessageUseCase(msgRepo, convRepo, partRepo, pub)
	msg, err := uc.SendMessage(context.Background(), SendMessageRequest{
		ConversationID: "conv-1",
		SenderID:       "user-1",
		Content:        "Hi",
		MessageType:    "", // empty => should default to text
	})

	require.NoError(t, err)
	assert.Equal(t, domain.MessageTypeText, msg.MessageType)
}

func TestSendMessage_DefaultSenderRoleBuyer(t *testing.T) {
	msgRepo, convRepo, partRepo, pub := defaultMessageMocks()

	partRepo.getByConversationAndUserFn = func(_ context.Context, _, _ string) (*domain.ConversationParticipant, error) {
		return &domain.ConversationParticipant{}, nil
	}
	msgRepo.createFn = func(_ context.Context, _ *domain.Message) error { return nil }

	uc := newMessageUseCase(msgRepo, convRepo, partRepo, pub)
	msg, err := uc.SendMessage(context.Background(), SendMessageRequest{
		ConversationID: "conv-1",
		SenderID:       "user-1",
		Content:        "Hi",
		SenderRole:     "", // empty => should default to buyer
	})

	require.NoError(t, err)
	assert.Equal(t, domain.SenderRoleBuyer, msg.SenderRole)
}

func TestSendMessage_SenderNotParticipant(t *testing.T) {
	msgRepo, convRepo, partRepo, pub := defaultMessageMocks()

	partRepo.getByConversationAndUserFn = func(_ context.Context, _, _ string) (*domain.ConversationParticipant, error) {
		return nil, errors.New("not found")
	}

	uc := newMessageUseCase(msgRepo, convRepo, partRepo, pub)
	_, err := uc.SendMessage(context.Background(), SendMessageRequest{
		ConversationID: "conv-1",
		SenderID:       "stranger",
		Content:        "Hello",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "sender is not a participant")
}

func TestSendMessage_RepoCreateError(t *testing.T) {
	msgRepo, convRepo, partRepo, pub := defaultMessageMocks()

	partRepo.getByConversationAndUserFn = func(_ context.Context, _, _ string) (*domain.ConversationParticipant, error) {
		return &domain.ConversationParticipant{}, nil
	}
	msgRepo.createFn = func(_ context.Context, _ *domain.Message) error {
		return errors.New("db write failed")
	}

	uc := newMessageUseCase(msgRepo, convRepo, partRepo, pub)
	_, err := uc.SendMessage(context.Background(), SendMessageRequest{
		ConversationID: "conv-1",
		SenderID:       "user-1",
		Content:        "Hello",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create message")
}

func TestSendMessage_WithAttachments(t *testing.T) {
	msgRepo, convRepo, partRepo, pub := defaultMessageMocks()

	partRepo.getByConversationAndUserFn = func(_ context.Context, _, _ string) (*domain.ConversationParticipant, error) {
		return &domain.ConversationParticipant{}, nil
	}
	msgRepo.createFn = func(_ context.Context, _ *domain.Message) error { return nil }

	uc := newMessageUseCase(msgRepo, convRepo, partRepo, pub)
	msg, err := uc.SendMessage(context.Background(), SendMessageRequest{
		ConversationID: "conv-1",
		SenderID:       "user-1",
		Content:        "See attachments",
		MessageType:    "image",
		Attachments:    []string{"img1.jpg", "img2.jpg"},
	})

	require.NoError(t, err)
	assert.Equal(t, domain.MessageTypeImage, msg.MessageType)
	assert.Len(t, msg.Attachments, 2)
	assert.Equal(t, "img1.jpg", msg.Attachments[0])
}

func TestSendMessage_PublishEventFailDoesNotBreak(t *testing.T) {
	msgRepo, convRepo, partRepo, pub := defaultMessageMocks()

	partRepo.getByConversationAndUserFn = func(_ context.Context, _, _ string) (*domain.ConversationParticipant, error) {
		return &domain.ConversationParticipant{}, nil
	}
	msgRepo.createFn = func(_ context.Context, _ *domain.Message) error { return nil }
	pub.publishFn = func(_ context.Context, _ string, _ interface{}) error {
		return errors.New("nats down")
	}

	uc := newMessageUseCase(msgRepo, convRepo, partRepo, pub)
	msg, err := uc.SendMessage(context.Background(), SendMessageRequest{
		ConversationID: "conv-1",
		SenderID:       "user-1",
		Content:        "Hello",
	})

	require.NoError(t, err)
	require.NotNil(t, msg)
}

// ===========================================================================
// ListMessages tests
// ===========================================================================

func TestListMessages_Success(t *testing.T) {
	msgRepo, convRepo, partRepo, pub := defaultMessageMocks()
	expectedMsgs := []domain.Message{
		{ID: "msg-1", Content: "Hello"},
		{ID: "msg-2", Content: "Hi back"},
	}
	msgRepo.listByConversationFn = func(_ context.Context, convID string, page, pageSize int) ([]domain.Message, int64, error) {
		assert.Equal(t, "conv-1", convID)
		assert.Equal(t, 1, page)
		assert.Equal(t, 50, pageSize)
		return expectedMsgs, 2, nil
	}

	uc := newMessageUseCase(msgRepo, convRepo, partRepo, pub)
	msgs, total, err := uc.ListMessages(context.Background(), "conv-1", 1, 50)

	require.NoError(t, err)
	assert.Len(t, msgs, 2)
	assert.Equal(t, int64(2), total)
}

func TestListMessages_DefaultPagination(t *testing.T) {
	msgRepo, convRepo, partRepo, pub := defaultMessageMocks()

	var capturedPage, capturedPageSize int
	msgRepo.listByConversationFn = func(_ context.Context, _ string, page, pageSize int) ([]domain.Message, int64, error) {
		capturedPage = page
		capturedPageSize = pageSize
		return nil, 0, nil
	}

	uc := newMessageUseCase(msgRepo, convRepo, partRepo, pub)

	// page < 1 should default to 1
	_, _, _ = uc.ListMessages(context.Background(), "conv-1", 0, 50)
	assert.Equal(t, 1, capturedPage)
	assert.Equal(t, 50, capturedPageSize)

	// pageSize > 100 should default to 50
	_, _, _ = uc.ListMessages(context.Background(), "conv-1", 1, 200)
	assert.Equal(t, 50, capturedPageSize)

	// pageSize < 1 should default to 50
	_, _, _ = uc.ListMessages(context.Background(), "conv-1", 1, -1)
	assert.Equal(t, 50, capturedPageSize)
}

// ===========================================================================
// MarkAsRead tests
// ===========================================================================

func TestMarkAsRead_Success(t *testing.T) {
	msgRepo, convRepo, partRepo, pub := defaultMessageMocks()

	var markedConvID, markedUserID string
	msgRepo.markAsReadFn = func(_ context.Context, convID, userID string) error {
		markedConvID = convID
		markedUserID = userID
		return nil
	}

	uc := newMessageUseCase(msgRepo, convRepo, partRepo, pub)
	err := uc.MarkAsRead(context.Background(), "conv-1", "user-1")

	require.NoError(t, err)
	assert.Equal(t, "conv-1", markedConvID)
	assert.Equal(t, "user-1", markedUserID)
}

func TestMarkAsRead_RepoError(t *testing.T) {
	msgRepo, convRepo, partRepo, pub := defaultMessageMocks()
	msgRepo.markAsReadFn = func(_ context.Context, _, _ string) error {
		return errors.New("db error")
	}

	uc := newMessageUseCase(msgRepo, convRepo, partRepo, pub)
	err := uc.MarkAsRead(context.Background(), "conv-1", "user-1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to mark messages as read")
}

func TestMarkAsRead_UpdatesParticipantLastRead(t *testing.T) {
	msgRepo, convRepo, partRepo, pub := defaultMessageMocks()
	msgRepo.markAsReadFn = func(_ context.Context, _, _ string) error { return nil }

	var lastReadUpdated bool
	partRepo.updateLastReadFn = func(_ context.Context, _, _ string) error {
		lastReadUpdated = true
		return nil
	}

	uc := newMessageUseCase(msgRepo, convRepo, partRepo, pub)
	err := uc.MarkAsRead(context.Background(), "conv-1", "user-1")

	require.NoError(t, err)
	assert.True(t, lastReadUpdated)
}

// ===========================================================================
// GetUnreadCount tests
// ===========================================================================

func TestGetUnreadCount_Success(t *testing.T) {
	msgRepo, convRepo, partRepo, pub := defaultMessageMocks()
	msgRepo.countUnreadFn = func(_ context.Context, convID, userID string) (int64, error) {
		assert.Equal(t, "conv-1", convID)
		assert.Equal(t, "user-1", userID)
		return 5, nil
	}

	uc := newMessageUseCase(msgRepo, convRepo, partRepo, pub)
	count, err := uc.GetUnreadCount(context.Background(), "conv-1", "user-1")

	require.NoError(t, err)
	assert.Equal(t, int64(5), count)
}

func TestGetUnreadCount_RepoError(t *testing.T) {
	msgRepo, convRepo, partRepo, pub := defaultMessageMocks()
	msgRepo.countUnreadFn = func(_ context.Context, _, _ string) (int64, error) {
		return 0, errors.New("db error")
	}

	uc := newMessageUseCase(msgRepo, convRepo, partRepo, pub)
	_, err := uc.GetUnreadCount(context.Background(), "conv-1", "user-1")

	require.Error(t, err)
}
