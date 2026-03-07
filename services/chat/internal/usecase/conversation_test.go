package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/chat/internal/domain"
)

// ---------------------------------------------------------------------------
// Hand-written function-field mocks
// ---------------------------------------------------------------------------

// --- ConversationRepository mock ---

type mockConversationRepo struct {
	getByIDFn           func(ctx context.Context, id string) (*domain.Conversation, error)
	listByUserFn        func(ctx context.Context, userID string, status string, page, pageSize int) ([]domain.Conversation, int64, error)
	listByParticipantsFn func(ctx context.Context, participantIDs []string) ([]domain.Conversation, error)
	createFn            func(ctx context.Context, conversation *domain.Conversation) error
	updateFn            func(ctx context.Context, conversation *domain.Conversation) error
	updateLastMessageFn func(ctx context.Context, id string, lastMessageAt *time.Time) error
}

func (m *mockConversationRepo) GetByID(ctx context.Context, id string) (*domain.Conversation, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}
func (m *mockConversationRepo) ListByUser(ctx context.Context, userID string, status string, page, pageSize int) ([]domain.Conversation, int64, error) {
	if m.listByUserFn != nil {
		return m.listByUserFn(ctx, userID, status, page, pageSize)
	}
	return nil, 0, nil
}
func (m *mockConversationRepo) ListByParticipants(ctx context.Context, participantIDs []string) ([]domain.Conversation, error) {
	if m.listByParticipantsFn != nil {
		return m.listByParticipantsFn(ctx, participantIDs)
	}
	return nil, nil
}
func (m *mockConversationRepo) Create(ctx context.Context, conversation *domain.Conversation) error {
	if m.createFn != nil {
		return m.createFn(ctx, conversation)
	}
	return nil
}
func (m *mockConversationRepo) Update(ctx context.Context, conversation *domain.Conversation) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, conversation)
	}
	return nil
}
func (m *mockConversationRepo) UpdateLastMessage(ctx context.Context, id string, lastMessageAt *time.Time) error {
	if m.updateLastMessageFn != nil {
		return m.updateLastMessageFn(ctx, id, lastMessageAt)
	}
	return nil
}

// --- ParticipantRepository mock ---

type mockParticipantRepo struct {
	getByConversationAndUserFn func(ctx context.Context, conversationID, userID string) (*domain.ConversationParticipant, error)
	listByConversationFn       func(ctx context.Context, conversationID string) ([]domain.ConversationParticipant, error)
	createFn                   func(ctx context.Context, participant *domain.ConversationParticipant) error
	updateLastReadFn           func(ctx context.Context, conversationID, userID string) error
}

func (m *mockParticipantRepo) GetByConversationAndUser(ctx context.Context, conversationID, userID string) (*domain.ConversationParticipant, error) {
	if m.getByConversationAndUserFn != nil {
		return m.getByConversationAndUserFn(ctx, conversationID, userID)
	}
	return nil, errors.New("not found")
}
func (m *mockParticipantRepo) ListByConversation(ctx context.Context, conversationID string) ([]domain.ConversationParticipant, error) {
	if m.listByConversationFn != nil {
		return m.listByConversationFn(ctx, conversationID)
	}
	return nil, nil
}
func (m *mockParticipantRepo) Create(ctx context.Context, participant *domain.ConversationParticipant) error {
	if m.createFn != nil {
		return m.createFn(ctx, participant)
	}
	return nil
}
func (m *mockParticipantRepo) UpdateLastRead(ctx context.Context, conversationID, userID string) error {
	if m.updateLastReadFn != nil {
		return m.updateLastReadFn(ctx, conversationID, userID)
	}
	return nil
}

// --- EventPublisher mock ---

type mockEventPublisher struct {
	publishFn func(ctx context.Context, subject string, data interface{}) error
}

func (m *mockEventPublisher) Publish(ctx context.Context, subject string, data interface{}) error {
	if m.publishFn != nil {
		return m.publishFn(ctx, subject, data)
	}
	return nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func newConversationUseCase(
	convRepo *mockConversationRepo,
	partRepo *mockParticipantRepo,
	pub *mockEventPublisher,
) *ConversationUseCase {
	return NewConversationUseCase(convRepo, partRepo, pub)
}

func defaultConversationMocks() (*mockConversationRepo, *mockParticipantRepo, *mockEventPublisher) {
	return &mockConversationRepo{}, &mockParticipantRepo{}, &mockEventPublisher{}
}

// ===========================================================================
// CreateConversation tests
// ===========================================================================

func TestCreateConversation_Success(t *testing.T) {
	convRepo, partRepo, pub := defaultConversationMocks()

	var savedConv *domain.Conversation
	convRepo.createFn = func(_ context.Context, c *domain.Conversation) error {
		savedConv = c
		return nil
	}

	var createdParticipants []*domain.ConversationParticipant
	partRepo.createFn = func(_ context.Context, p *domain.ConversationParticipant) error {
		createdParticipants = append(createdParticipants, p)
		return nil
	}

	uc := newConversationUseCase(convRepo, partRepo, pub)
	conv, err := uc.CreateConversation(context.Background(), "user-1", CreateConversationRequest{
		BuyerID:  "buyer-1",
		SellerID: "seller-1",
		OrderID:  "order-1",
		Subject:  "Product question",
	})

	require.NoError(t, err)
	require.NotNil(t, conv)
	assert.NotEmpty(t, conv.ID)
	assert.Equal(t, "buyer-1", conv.BuyerID)
	assert.Equal(t, "seller-1", conv.SellerID)
	assert.Equal(t, "order-1", conv.OrderID)
	assert.Equal(t, "Product question", conv.Subject)
	assert.Equal(t, domain.ConversationStatusActive, conv.Status)
	assert.Equal(t, domain.ConversationTypeBuyerSeller, conv.Type)
	assert.NotNil(t, savedConv)
	assert.Len(t, createdParticipants, 2)
}

func TestCreateConversation_DefaultTypeBuyerSeller(t *testing.T) {
	convRepo, partRepo, pub := defaultConversationMocks()
	convRepo.createFn = func(_ context.Context, _ *domain.Conversation) error { return nil }

	uc := newConversationUseCase(convRepo, partRepo, pub)
	conv, err := uc.CreateConversation(context.Background(), "user-1", CreateConversationRequest{
		Type:     "", // empty => should default to buyer_seller
		BuyerID:  "buyer-1",
		SellerID: "seller-1",
	})

	require.NoError(t, err)
	assert.Equal(t, domain.ConversationTypeBuyerSeller, conv.Type)
}

func TestCreateConversation_SupportType(t *testing.T) {
	convRepo, partRepo, pub := defaultConversationMocks()
	convRepo.createFn = func(_ context.Context, _ *domain.Conversation) error { return nil }

	uc := newConversationUseCase(convRepo, partRepo, pub)
	conv, err := uc.CreateConversation(context.Background(), "user-1", CreateConversationRequest{
		Type:     "support",
		BuyerID:  "buyer-1",
		SellerID: "seller-1",
	})

	require.NoError(t, err)
	assert.Equal(t, domain.ConversationTypeSupport, conv.Type)
}

func TestCreateConversation_RepoError(t *testing.T) {
	convRepo, partRepo, pub := defaultConversationMocks()
	convRepo.createFn = func(_ context.Context, _ *domain.Conversation) error {
		return errors.New("db connection lost")
	}

	uc := newConversationUseCase(convRepo, partRepo, pub)
	_, err := uc.CreateConversation(context.Background(), "user-1", CreateConversationRequest{
		BuyerID:  "buyer-1",
		SellerID: "seller-1",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create conversation")
}

func TestCreateConversation_ParticipantCreateError(t *testing.T) {
	convRepo, partRepo, pub := defaultConversationMocks()
	convRepo.createFn = func(_ context.Context, _ *domain.Conversation) error { return nil }
	partRepo.createFn = func(_ context.Context, _ *domain.ConversationParticipant) error {
		return errors.New("participant insert failed")
	}

	uc := newConversationUseCase(convRepo, partRepo, pub)
	_, err := uc.CreateConversation(context.Background(), "user-1", CreateConversationRequest{
		BuyerID:  "buyer-1",
		SellerID: "seller-1",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create buyer participant")
}

func TestCreateConversation_ParticipantRoles(t *testing.T) {
	convRepo, partRepo, pub := defaultConversationMocks()
	convRepo.createFn = func(_ context.Context, _ *domain.Conversation) error { return nil }

	var participants []*domain.ConversationParticipant
	partRepo.createFn = func(_ context.Context, p *domain.ConversationParticipant) error {
		participants = append(participants, p)
		return nil
	}

	uc := newConversationUseCase(convRepo, partRepo, pub)
	_, err := uc.CreateConversation(context.Background(), "user-1", CreateConversationRequest{
		BuyerID:  "buyer-1",
		SellerID: "seller-1",
	})

	require.NoError(t, err)
	require.Len(t, participants, 2)
	assert.Equal(t, domain.ParticipantRoleBuyer, participants[0].Role)
	assert.Equal(t, "buyer-1", participants[0].UserID)
	assert.Equal(t, domain.ParticipantRoleSeller, participants[1].Role)
	assert.Equal(t, "seller-1", participants[1].UserID)
}

func TestCreateConversation_PublishEventFailDoesNotBreak(t *testing.T) {
	convRepo, partRepo, pub := defaultConversationMocks()
	convRepo.createFn = func(_ context.Context, _ *domain.Conversation) error { return nil }
	pub.publishFn = func(_ context.Context, _ string, _ interface{}) error {
		return errors.New("nats down")
	}

	uc := newConversationUseCase(convRepo, partRepo, pub)
	conv, err := uc.CreateConversation(context.Background(), "user-1", CreateConversationRequest{
		BuyerID:  "buyer-1",
		SellerID: "seller-1",
	})

	// Should succeed even if event publishing fails
	require.NoError(t, err)
	require.NotNil(t, conv)
}

// ===========================================================================
// GetConversation tests
// ===========================================================================

func TestGetConversation_Success(t *testing.T) {
	convRepo, partRepo, pub := defaultConversationMocks()
	expected := &domain.Conversation{
		ID:      "conv-1",
		BuyerID: "buyer-1",
		Subject: "Test",
		Status:  domain.ConversationStatusActive,
	}
	convRepo.getByIDFn = func(_ context.Context, id string) (*domain.Conversation, error) {
		assert.Equal(t, "conv-1", id)
		return expected, nil
	}

	uc := newConversationUseCase(convRepo, partRepo, pub)
	conv, err := uc.GetConversation(context.Background(), "conv-1")

	require.NoError(t, err)
	assert.Equal(t, expected, conv)
}

func TestGetConversation_NotFound(t *testing.T) {
	convRepo, partRepo, pub := defaultConversationMocks()
	convRepo.getByIDFn = func(_ context.Context, _ string) (*domain.Conversation, error) {
		return nil, errors.New("not found")
	}

	uc := newConversationUseCase(convRepo, partRepo, pub)
	_, err := uc.GetConversation(context.Background(), "nonexistent")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// ===========================================================================
// ListUserConversations tests
// ===========================================================================

func TestListUserConversations_Success(t *testing.T) {
	convRepo, partRepo, pub := defaultConversationMocks()
	expectedConvs := []domain.Conversation{
		{ID: "conv-1", BuyerID: "user-1"},
		{ID: "conv-2", BuyerID: "user-1"},
	}
	convRepo.listByUserFn = func(_ context.Context, userID string, status string, page, pageSize int) ([]domain.Conversation, int64, error) {
		assert.Equal(t, "user-1", userID)
		assert.Equal(t, "active", status)
		assert.Equal(t, 1, page)
		assert.Equal(t, 20, pageSize)
		return expectedConvs, 2, nil
	}

	uc := newConversationUseCase(convRepo, partRepo, pub)
	convs, total, err := uc.ListUserConversations(context.Background(), "user-1", "active", 1, 20)

	require.NoError(t, err)
	assert.Len(t, convs, 2)
	assert.Equal(t, int64(2), total)
}

func TestListUserConversations_DefaultPagination(t *testing.T) {
	convRepo, partRepo, pub := defaultConversationMocks()

	var capturedPage, capturedPageSize int
	convRepo.listByUserFn = func(_ context.Context, _ string, _ string, page, pageSize int) ([]domain.Conversation, int64, error) {
		capturedPage = page
		capturedPageSize = pageSize
		return nil, 0, nil
	}

	uc := newConversationUseCase(convRepo, partRepo, pub)

	// page < 1 should default to 1
	_, _, _ = uc.ListUserConversations(context.Background(), "user-1", "", 0, 50)
	assert.Equal(t, 1, capturedPage)
	assert.Equal(t, 50, capturedPageSize)

	// pageSize > 100 should default to 20
	_, _, _ = uc.ListUserConversations(context.Background(), "user-1", "", 1, 200)
	assert.Equal(t, 20, capturedPageSize)

	// pageSize < 1 should default to 20
	_, _, _ = uc.ListUserConversations(context.Background(), "user-1", "", 1, 0)
	assert.Equal(t, 20, capturedPageSize)
}

// ===========================================================================
// ArchiveConversation tests
// ===========================================================================

func TestArchiveConversation_Success(t *testing.T) {
	convRepo, partRepo, pub := defaultConversationMocks()
	convRepo.getByIDFn = func(_ context.Context, _ string) (*domain.Conversation, error) {
		return &domain.Conversation{
			ID:     "conv-1",
			Status: domain.ConversationStatusActive,
		}, nil
	}

	var updatedConv *domain.Conversation
	convRepo.updateFn = func(_ context.Context, c *domain.Conversation) error {
		updatedConv = c
		return nil
	}

	uc := newConversationUseCase(convRepo, partRepo, pub)
	err := uc.ArchiveConversation(context.Background(), "conv-1")

	require.NoError(t, err)
	require.NotNil(t, updatedConv)
	assert.Equal(t, domain.ConversationStatusArchived, updatedConv.Status)
}

func TestArchiveConversation_NotFound(t *testing.T) {
	convRepo, partRepo, pub := defaultConversationMocks()
	convRepo.getByIDFn = func(_ context.Context, _ string) (*domain.Conversation, error) {
		return nil, errors.New("not found")
	}

	uc := newConversationUseCase(convRepo, partRepo, pub)
	err := uc.ArchiveConversation(context.Background(), "nonexistent")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "conversation not found")
}
