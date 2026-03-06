package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/return/internal/domain"
)

// --- Helper to build a valid CreateDisputeRequest ---

func validCreateDisputeRequest() CreateDisputeRequest {
	return CreateDisputeRequest{
		OrderID:     "order-1",
		ReturnID:    "return-1",
		BuyerID:     "buyer-1",
		SellerID:    "seller-1",
		Type:        "item_not_received",
		Description: "I never received the item",
	}
}

// --- CreateDispute tests ---

func TestDisputeUseCase_CreateDispute_Success(t *testing.T) {
	var captured *domain.Dispute
	disputeRepo := &mockDisputeRepo{
		createFn: func(_ context.Context, d *domain.Dispute) error {
			captured = d
			return nil
		},
	}
	msgRepo := &mockDisputeMessageRepo{}
	pub := &mockEventPublisher{}
	uc := NewDisputeUseCase(disputeRepo, msgRepo, pub)

	result, err := uc.CreateDispute(context.Background(), validCreateDisputeRequest())
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Len(t, result.ID, 36, "dispute ID should be a UUID")
	assert.Equal(t, domain.DisputeStatusOpen, result.Status)
	assert.Equal(t, "order-1", result.OrderID)
	assert.Equal(t, "return-1", result.ReturnID)
	assert.Equal(t, "buyer-1", result.BuyerID)
	assert.Equal(t, "seller-1", result.SellerID)
	assert.Equal(t, domain.DisputeType("item_not_received"), result.Type)
	assert.Equal(t, "I never received the item", result.Description)

	require.NotNil(t, captured)
	assert.Equal(t, result.ID, captured.ID)

	// Verify event
	require.Len(t, pub.calls, 1)
	assert.Equal(t, "dispute.opened", pub.calls[0].Subject)
	eventData := pub.calls[0].Data.(map[string]interface{})
	assert.Equal(t, result.ID, eventData["dispute_id"])
	assert.Equal(t, "order-1", eventData["order_id"])
	assert.Equal(t, "item_not_received", eventData["type"])
	assert.Equal(t, "buyer-1", eventData["buyer_id"])
	assert.Equal(t, "seller-1", eventData["seller_id"])
}

func TestDisputeUseCase_CreateDispute_ValidationErrors(t *testing.T) {
	disputeRepo := &mockDisputeRepo{}
	msgRepo := &mockDisputeMessageRepo{}
	pub := &mockEventPublisher{}
	uc := NewDisputeUseCase(disputeRepo, msgRepo, pub)

	tests := []struct {
		name    string
		mutate  func(req *CreateDisputeRequest)
		wantErr string
	}{
		{
			name:    "empty order_id",
			mutate:  func(req *CreateDisputeRequest) { req.OrderID = "" },
			wantErr: "order_id, buyer_id, and seller_id are required",
		},
		{
			name:    "empty buyer_id",
			mutate:  func(req *CreateDisputeRequest) { req.BuyerID = "" },
			wantErr: "order_id, buyer_id, and seller_id are required",
		},
		{
			name:    "empty seller_id",
			mutate:  func(req *CreateDisputeRequest) { req.SellerID = "" },
			wantErr: "order_id, buyer_id, and seller_id are required",
		},
		{
			name:    "empty type",
			mutate:  func(req *CreateDisputeRequest) { req.Type = "" },
			wantErr: "type and description are required",
		},
		{
			name:    "empty description",
			mutate:  func(req *CreateDisputeRequest) { req.Description = "" },
			wantErr: "type and description are required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := validCreateDisputeRequest()
			tc.mutate(&req)

			result, err := uc.CreateDispute(context.Background(), req)
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}

	assert.Empty(t, pub.calls)
}

func TestDisputeUseCase_CreateDispute_RepoError(t *testing.T) {
	disputeRepo := &mockDisputeRepo{
		createFn: func(_ context.Context, _ *domain.Dispute) error {
			return errors.New("db error")
		},
	}
	msgRepo := &mockDisputeMessageRepo{}
	pub := &mockEventPublisher{}
	uc := NewDisputeUseCase(disputeRepo, msgRepo, pub)

	result, err := uc.CreateDispute(context.Background(), validCreateDisputeRequest())
	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create dispute")
	assert.Contains(t, err.Error(), "db error")
	assert.Empty(t, pub.calls)
}

// --- GetDispute tests ---

func TestDisputeUseCase_GetDispute_Success(t *testing.T) {
	expected := &domain.Dispute{ID: "disp-1", OrderID: "ord-1"}
	disputeRepo := &mockDisputeRepo{
		getByIDFn: func(_ context.Context, id string) (*domain.Dispute, error) {
			assert.Equal(t, "disp-1", id)
			return expected, nil
		},
	}
	msgRepo := &mockDisputeMessageRepo{}
	pub := &mockEventPublisher{}
	uc := NewDisputeUseCase(disputeRepo, msgRepo, pub)

	result, err := uc.GetDispute(context.Background(), "disp-1")
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

// --- ListAllDisputes pagination tests ---

func TestDisputeUseCase_ListAllDisputes_Pagination(t *testing.T) {
	tests := []struct {
		name             string
		page             int
		pageSize         int
		expectedPage     int
		expectedPageSize int
	}{
		{"valid values", 2, 50, 2, 50},
		{"page clamped to 1", 0, 10, 1, 10},
		{"negative page clamped to 1", -3, 10, 1, 10},
		{"pageSize 0 defaults to 20", 1, 0, 1, 20},
		{"pageSize over 100 defaults to 20", 1, 150, 1, 20},
		{"negative pageSize defaults to 20", 1, -5, 1, 20},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			disputeRepo := &mockDisputeRepo{
				listAllFn: func(_ context.Context, page, pageSize int) ([]domain.Dispute, int64, error) {
					assert.Equal(t, tc.expectedPage, page)
					assert.Equal(t, tc.expectedPageSize, pageSize)
					return []domain.Dispute{}, 0, nil
				},
			}
			msgRepo := &mockDisputeMessageRepo{}
			pub := &mockEventPublisher{}
			uc := NewDisputeUseCase(disputeRepo, msgRepo, pub)

			_, _, err := uc.ListAllDisputes(context.Background(), tc.page, tc.pageSize)
			require.NoError(t, err)
		})
	}
}

// --- ListBuyerDisputes tests ---

func TestDisputeUseCase_ListBuyerDisputes_Pagination(t *testing.T) {
	disputeRepo := &mockDisputeRepo{
		listByBuyerFn: func(_ context.Context, buyerID string, page, pageSize int) ([]domain.Dispute, int64, error) {
			assert.Equal(t, "buyer-1", buyerID)
			assert.Equal(t, 1, page)
			assert.Equal(t, 20, pageSize)
			return []domain.Dispute{{ID: "d-1"}}, 1, nil
		},
	}
	msgRepo := &mockDisputeMessageRepo{}
	pub := &mockEventPublisher{}
	uc := NewDisputeUseCase(disputeRepo, msgRepo, pub)

	results, total, err := uc.ListBuyerDisputes(context.Background(), "buyer-1", 0, 0)
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, int64(1), total)
}

// --- AddMessage tests ---

func TestDisputeUseCase_AddMessage_Success(t *testing.T) {
	dispute := &domain.Dispute{
		ID:     "disp-1",
		Status: domain.DisputeStatusUnderReview,
	}
	var capturedMsg *domain.DisputeMessage
	disputeRepo := &mockDisputeRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Dispute, error) {
			return dispute, nil
		},
	}
	msgRepo := &mockDisputeMessageRepo{
		createFn: func(_ context.Context, msg *domain.DisputeMessage) error {
			capturedMsg = msg
			return nil
		},
	}
	pub := &mockEventPublisher{}
	uc := NewDisputeUseCase(disputeRepo, msgRepo, pub)

	req := AddMessageRequest{
		DisputeID:   "disp-1",
		SenderID:    "buyer-1",
		SenderRole:  "buyer",
		Message:     "Here is more info",
		Attachments: []string{"https://img.example.com/proof.jpg"},
	}

	result, err := uc.AddMessage(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.ID, 36, "message ID should be a UUID")
	assert.Equal(t, "disp-1", result.DisputeID)
	assert.Equal(t, "buyer-1", result.SenderID)
	assert.Equal(t, "buyer", result.SenderRole)
	assert.Equal(t, "Here is more info", result.Message)
	assert.Equal(t, []string{"https://img.example.com/proof.jpg"}, result.Attachments)
	require.NotNil(t, capturedMsg)
}

func TestDisputeUseCase_AddMessage_AutoTransitionsOpenToUnderReview(t *testing.T) {
	dispute := &domain.Dispute{
		ID:     "disp-1",
		Status: domain.DisputeStatusOpen,
	}
	var updatedDispute *domain.Dispute
	disputeRepo := &mockDisputeRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Dispute, error) {
			return dispute, nil
		},
		updateFn: func(_ context.Context, d *domain.Dispute) error {
			updatedDispute = d
			return nil
		},
	}
	msgRepo := &mockDisputeMessageRepo{
		createFn: func(_ context.Context, _ *domain.DisputeMessage) error { return nil },
	}
	pub := &mockEventPublisher{}
	uc := NewDisputeUseCase(disputeRepo, msgRepo, pub)

	req := AddMessageRequest{
		DisputeID:  "disp-1",
		SenderID:   "seller-1",
		SenderRole: "seller",
		Message:    "We are looking into this",
	}

	_, err := uc.AddMessage(context.Background(), req)
	require.NoError(t, err)

	// The dispute should have been updated to under_review
	require.NotNil(t, updatedDispute)
	assert.Equal(t, domain.DisputeStatusUnderReview, updatedDispute.Status)
}

func TestDisputeUseCase_AddMessage_BlocksOnResolvedBuyerDispute(t *testing.T) {
	dispute := &domain.Dispute{
		ID:     "disp-1",
		Status: domain.DisputeStatusResolvedBuyer,
	}
	disputeRepo := &mockDisputeRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Dispute, error) {
			return dispute, nil
		},
	}
	msgRepo := &mockDisputeMessageRepo{}
	pub := &mockEventPublisher{}
	uc := NewDisputeUseCase(disputeRepo, msgRepo, pub)

	req := AddMessageRequest{
		DisputeID:  "disp-1",
		SenderID:   "buyer-1",
		SenderRole: "buyer",
		Message:    "I want to add more",
	}

	result, err := uc.AddMessage(context.Background(), req)
	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot add message to resolved dispute")
}

func TestDisputeUseCase_AddMessage_BlocksOnResolvedSellerDispute(t *testing.T) {
	dispute := &domain.Dispute{
		ID:     "disp-1",
		Status: domain.DisputeStatusResolvedSeller,
	}
	disputeRepo := &mockDisputeRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Dispute, error) {
			return dispute, nil
		},
	}
	msgRepo := &mockDisputeMessageRepo{}
	pub := &mockEventPublisher{}
	uc := NewDisputeUseCase(disputeRepo, msgRepo, pub)

	req := AddMessageRequest{
		DisputeID:  "disp-1",
		SenderID:   "seller-1",
		SenderRole: "seller",
		Message:    "Additional info",
	}

	result, err := uc.AddMessage(context.Background(), req)
	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot add message to resolved dispute")
}

func TestDisputeUseCase_AddMessage_DisputeNotFound(t *testing.T) {
	disputeRepo := &mockDisputeRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Dispute, error) {
			return nil, errors.New("not found")
		},
	}
	msgRepo := &mockDisputeMessageRepo{}
	pub := &mockEventPublisher{}
	uc := NewDisputeUseCase(disputeRepo, msgRepo, pub)

	req := AddMessageRequest{
		DisputeID:  "nonexistent",
		SenderID:   "buyer-1",
		SenderRole: "buyer",
		Message:    "hello",
	}

	result, err := uc.AddMessage(context.Background(), req)
	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "dispute not found")
}

// --- ResolveDispute tests ---

func TestDisputeUseCase_ResolveDispute_Success(t *testing.T) {
	dispute := &domain.Dispute{
		ID:      "disp-1",
		OrderID: "ord-1",
		Status:  domain.DisputeStatusUnderReview,
	}
	var updatedDispute *domain.Dispute
	disputeRepo := &mockDisputeRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Dispute, error) {
			return dispute, nil
		},
		updateFn: func(_ context.Context, d *domain.Dispute) error {
			updatedDispute = d
			return nil
		},
	}
	msgRepo := &mockDisputeMessageRepo{}
	pub := &mockEventPublisher{}
	uc := NewDisputeUseCase(disputeRepo, msgRepo, pub)

	req := ResolveDisputeRequest{
		DisputeID:  "disp-1",
		Resolution: "Full refund to buyer",
		ResolvedBy: "admin-1",
		Status:     "resolved_buyer",
	}

	result, err := uc.ResolveDispute(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, domain.DisputeStatusResolvedBuyer, result.Status)
	assert.Equal(t, "Full refund to buyer", result.Resolution)
	assert.Equal(t, "admin-1", result.ResolvedBy)
	assert.NotNil(t, result.ResolvedAt)

	require.NotNil(t, updatedDispute)

	// Verify event
	require.Len(t, pub.calls, 1)
	assert.Equal(t, "dispute.resolved", pub.calls[0].Subject)
	eventData := pub.calls[0].Data.(map[string]interface{})
	assert.Equal(t, "disp-1", eventData["dispute_id"])
	assert.Equal(t, "ord-1", eventData["order_id"])
	assert.Equal(t, "resolved_buyer", eventData["status"])
	assert.Equal(t, "Full refund to buyer", eventData["resolution"])
}

func TestDisputeUseCase_ResolveDispute_ResolvedSeller(t *testing.T) {
	dispute := &domain.Dispute{
		ID:      "disp-1",
		OrderID: "ord-1",
		Status:  domain.DisputeStatusUnderReview,
	}
	disputeRepo := &mockDisputeRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Dispute, error) {
			return dispute, nil
		},
		updateFn: func(_ context.Context, _ *domain.Dispute) error { return nil },
	}
	msgRepo := &mockDisputeMessageRepo{}
	pub := &mockEventPublisher{}
	uc := NewDisputeUseCase(disputeRepo, msgRepo, pub)

	req := ResolveDisputeRequest{
		DisputeID:  "disp-1",
		Resolution: "Seller was correct",
		ResolvedBy: "admin-1",
		Status:     "resolved_seller",
	}

	result, err := uc.ResolveDispute(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, domain.DisputeStatusResolvedSeller, result.Status)
}

func TestDisputeUseCase_ResolveDispute_AlreadyResolved(t *testing.T) {
	tests := []struct {
		name   string
		status domain.DisputeStatus
	}{
		{"already resolved_buyer", domain.DisputeStatusResolvedBuyer},
		{"already resolved_seller", domain.DisputeStatusResolvedSeller},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dispute := &domain.Dispute{
				ID:     "disp-1",
				Status: tc.status,
			}
			disputeRepo := &mockDisputeRepo{
				getByIDFn: func(_ context.Context, _ string) (*domain.Dispute, error) {
					return dispute, nil
				},
			}
			msgRepo := &mockDisputeMessageRepo{}
			pub := &mockEventPublisher{}
			uc := NewDisputeUseCase(disputeRepo, msgRepo, pub)

			req := ResolveDisputeRequest{
				DisputeID:  "disp-1",
				Resolution: "Another resolution",
				ResolvedBy: "admin-1",
				Status:     "resolved_buyer",
			}

			result, err := uc.ResolveDispute(context.Background(), req)
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "dispute is already resolved")
		})
	}
}

func TestDisputeUseCase_ResolveDispute_InvalidStatus(t *testing.T) {
	dispute := &domain.Dispute{
		ID:     "disp-1",
		Status: domain.DisputeStatusOpen,
	}
	disputeRepo := &mockDisputeRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Dispute, error) {
			return dispute, nil
		},
	}
	msgRepo := &mockDisputeMessageRepo{}
	pub := &mockEventPublisher{}
	uc := NewDisputeUseCase(disputeRepo, msgRepo, pub)

	tests := []struct {
		name   string
		status string
	}{
		{"invalid status open", "open"},
		{"invalid status escalated", "escalated"},
		{"invalid status under_review", "under_review"},
		{"invalid status garbage", "garbage"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := ResolveDisputeRequest{
				DisputeID:  "disp-1",
				Resolution: "Some resolution",
				ResolvedBy: "admin-1",
				Status:     tc.status,
			}

			result, err := uc.ResolveDispute(context.Background(), req)
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid resolution status: must be resolved_buyer or resolved_seller")
		})
	}
}

func TestDisputeUseCase_ResolveDispute_NotFound(t *testing.T) {
	disputeRepo := &mockDisputeRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Dispute, error) {
			return nil, errors.New("not found")
		},
	}
	msgRepo := &mockDisputeMessageRepo{}
	pub := &mockEventPublisher{}
	uc := NewDisputeUseCase(disputeRepo, msgRepo, pub)

	req := ResolveDisputeRequest{
		DisputeID:  "nonexistent",
		Resolution: "N/A",
		ResolvedBy: "admin-1",
		Status:     "resolved_buyer",
	}

	result, err := uc.ResolveDispute(context.Background(), req)
	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "dispute not found")
}
