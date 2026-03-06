package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/return/internal/domain"
)

// --- Shared mocks used across usecase test files ---

type mockReturnRepo struct {
	getByIDFn      func(ctx context.Context, id string) (*domain.Return, error)
	getByOrderIDFn func(ctx context.Context, orderID string) ([]domain.Return, error)
	listByBuyerFn  func(ctx context.Context, buyerID string, page, pageSize int) ([]domain.Return, int64, error)
	listBySellerFn func(ctx context.Context, sellerID string, page, pageSize int) ([]domain.Return, int64, error)
	createFn       func(ctx context.Context, ret *domain.Return) error
	updateFn       func(ctx context.Context, ret *domain.Return) error
}

func (m *mockReturnRepo) GetByID(ctx context.Context, id string) (*domain.Return, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockReturnRepo) GetByOrderID(ctx context.Context, orderID string) ([]domain.Return, error) {
	if m.getByOrderIDFn != nil {
		return m.getByOrderIDFn(ctx, orderID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockReturnRepo) ListByBuyer(ctx context.Context, buyerID string, page, pageSize int) ([]domain.Return, int64, error) {
	if m.listByBuyerFn != nil {
		return m.listByBuyerFn(ctx, buyerID, page, pageSize)
	}
	return nil, 0, errors.New("not implemented")
}

func (m *mockReturnRepo) ListBySeller(ctx context.Context, sellerID string, page, pageSize int) ([]domain.Return, int64, error) {
	if m.listBySellerFn != nil {
		return m.listBySellerFn(ctx, sellerID, page, pageSize)
	}
	return nil, 0, errors.New("not implemented")
}

func (m *mockReturnRepo) Create(ctx context.Context, ret *domain.Return) error {
	if m.createFn != nil {
		return m.createFn(ctx, ret)
	}
	return nil
}

func (m *mockReturnRepo) Update(ctx context.Context, ret *domain.Return) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, ret)
	}
	return nil
}

type mockEventPublisher struct {
	publishFn func(ctx context.Context, subject string, data interface{}) error
	calls     []publishCall
}

type publishCall struct {
	Subject string
	Data    interface{}
}

func (m *mockEventPublisher) Publish(ctx context.Context, subject string, data interface{}) error {
	m.calls = append(m.calls, publishCall{Subject: subject, Data: data})
	if m.publishFn != nil {
		return m.publishFn(ctx, subject, data)
	}
	return nil
}

type mockDisputeRepo struct {
	getByIDFn      func(ctx context.Context, id string) (*domain.Dispute, error)
	getByOrderIDFn func(ctx context.Context, orderID string) ([]domain.Dispute, error)
	listAllFn      func(ctx context.Context, page, pageSize int) ([]domain.Dispute, int64, error)
	listByBuyerFn  func(ctx context.Context, buyerID string, page, pageSize int) ([]domain.Dispute, int64, error)
	createFn       func(ctx context.Context, dispute *domain.Dispute) error
	updateFn       func(ctx context.Context, dispute *domain.Dispute) error
}

func (m *mockDisputeRepo) GetByID(ctx context.Context, id string) (*domain.Dispute, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockDisputeRepo) GetByOrderID(ctx context.Context, orderID string) ([]domain.Dispute, error) {
	if m.getByOrderIDFn != nil {
		return m.getByOrderIDFn(ctx, orderID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockDisputeRepo) ListAll(ctx context.Context, page, pageSize int) ([]domain.Dispute, int64, error) {
	if m.listAllFn != nil {
		return m.listAllFn(ctx, page, pageSize)
	}
	return nil, 0, errors.New("not implemented")
}

func (m *mockDisputeRepo) ListByBuyer(ctx context.Context, buyerID string, page, pageSize int) ([]domain.Dispute, int64, error) {
	if m.listByBuyerFn != nil {
		return m.listByBuyerFn(ctx, buyerID, page, pageSize)
	}
	return nil, 0, errors.New("not implemented")
}

func (m *mockDisputeRepo) Create(ctx context.Context, dispute *domain.Dispute) error {
	if m.createFn != nil {
		return m.createFn(ctx, dispute)
	}
	return nil
}

func (m *mockDisputeRepo) Update(ctx context.Context, dispute *domain.Dispute) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, dispute)
	}
	return nil
}

type mockDisputeMessageRepo struct {
	getByDisputeIDFn func(ctx context.Context, disputeID string) ([]domain.DisputeMessage, error)
	createFn         func(ctx context.Context, msg *domain.DisputeMessage) error
}

func (m *mockDisputeMessageRepo) GetByDisputeID(ctx context.Context, disputeID string) ([]domain.DisputeMessage, error) {
	if m.getByDisputeIDFn != nil {
		return m.getByDisputeIDFn(ctx, disputeID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockDisputeMessageRepo) Create(ctx context.Context, msg *domain.DisputeMessage) error {
	if m.createFn != nil {
		return m.createFn(ctx, msg)
	}
	return nil
}

// --- Helper to build a valid CreateReturnRequest ---

func validCreateReturnRequest() CreateReturnRequest {
	return CreateReturnRequest{
		OrderID:           "order-1",
		BuyerID:           "buyer-1",
		SellerID:          "seller-1",
		Reason:            "defective",
		Description:       "The product arrived broken",
		ImageURLs:         []string{"https://img.example.com/1.jpg"},
		RefundAmountCents: 5000,
		RefundMethod:      "original_payment",
		Items: []CreateReturnItemRequest{
			{
				OrderItemID: "oi-1",
				ProductID:   "prod-1",
				VariantID:   "var-1",
				Quantity:    1,
				Reason:      "defective",
			},
		},
	}
}

// --- CreateReturnUseCase tests ---

func TestCreateReturnUseCase_Execute_Success(t *testing.T) {
	var captured *domain.Return
	repo := &mockReturnRepo{
		createFn: func(_ context.Context, ret *domain.Return) error {
			captured = ret
			return nil
		},
	}
	pub := &mockEventPublisher{}

	uc := NewCreateReturnUseCase(repo, pub)
	req := validCreateReturnRequest()

	result, err := uc.Execute(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify UUID format (36 chars with hyphens)
	assert.Len(t, result.ID, 36, "return ID should be a UUID")
	assert.Equal(t, domain.ReturnStatusRequested, result.Status)
	assert.Equal(t, req.OrderID, result.OrderID)
	assert.Equal(t, req.BuyerID, result.BuyerID)
	assert.Equal(t, req.SellerID, result.SellerID)
	assert.Equal(t, domain.ReturnReason(req.Reason), result.Reason)
	assert.Equal(t, req.Description, result.Description)
	assert.Equal(t, req.ImageURLs, result.ImageURLs)
	assert.Equal(t, req.RefundAmountCents, result.RefundAmountCents)
	assert.Equal(t, req.RefundMethod, result.RefundMethod)

	// Verify items
	require.Len(t, result.Items, 1)
	item := result.Items[0]
	assert.Len(t, item.ID, 36, "item ID should be a UUID")
	assert.Equal(t, result.ID, item.ReturnID)
	assert.Equal(t, "oi-1", item.OrderItemID)
	assert.Equal(t, "prod-1", item.ProductID)
	assert.Equal(t, "var-1", item.VariantID)
	assert.Equal(t, 1, item.Quantity)
	assert.Equal(t, "defective", item.Reason)

	// Verify repo was called
	require.NotNil(t, captured)
	assert.Equal(t, result.ID, captured.ID)

	// Verify event published
	require.Len(t, pub.calls, 1)
	assert.Equal(t, "return.requested", pub.calls[0].Subject)
	eventData := pub.calls[0].Data.(map[string]interface{})
	assert.Equal(t, result.ID, eventData["return_id"])
	assert.Equal(t, req.OrderID, eventData["order_id"])
	assert.Equal(t, req.BuyerID, eventData["buyer_id"])
	assert.Equal(t, req.SellerID, eventData["seller_id"])
	assert.Equal(t, "defective", eventData["reason"])
	assert.Equal(t, 1, eventData["item_count"])
}

func TestCreateReturnUseCase_Execute_ValidationErrors(t *testing.T) {
	repo := &mockReturnRepo{}
	pub := &mockEventPublisher{}
	uc := NewCreateReturnUseCase(repo, pub)

	tests := []struct {
		name    string
		mutate  func(req *CreateReturnRequest)
		wantErr string
	}{
		{
			name:    "empty order_id",
			mutate:  func(req *CreateReturnRequest) { req.OrderID = "" },
			wantErr: "order_id, buyer_id, and seller_id are required",
		},
		{
			name:    "empty buyer_id",
			mutate:  func(req *CreateReturnRequest) { req.BuyerID = "" },
			wantErr: "order_id, buyer_id, and seller_id are required",
		},
		{
			name:    "empty seller_id",
			mutate:  func(req *CreateReturnRequest) { req.SellerID = "" },
			wantErr: "order_id, buyer_id, and seller_id are required",
		},
		{
			name:    "empty reason",
			mutate:  func(req *CreateReturnRequest) { req.Reason = "" },
			wantErr: "reason is required",
		},
		{
			name:    "empty items",
			mutate:  func(req *CreateReturnRequest) { req.Items = nil },
			wantErr: "at least one item is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := validCreateReturnRequest()
			tc.mutate(&req)

			result, err := uc.Execute(context.Background(), req)
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}

	// Verify no events were published for validation errors
	assert.Empty(t, pub.calls)
}

func TestCreateReturnUseCase_Execute_RepoError(t *testing.T) {
	repo := &mockReturnRepo{
		createFn: func(_ context.Context, _ *domain.Return) error {
			return errors.New("db connection failed")
		},
	}
	pub := &mockEventPublisher{}
	uc := NewCreateReturnUseCase(repo, pub)

	result, err := uc.Execute(context.Background(), validCreateReturnRequest())
	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create return")
	assert.Contains(t, err.Error(), "db connection failed")

	// Event should not be published on repo error
	assert.Empty(t, pub.calls)
}
