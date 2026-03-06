package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/return/internal/domain"
)

// --- GetReturn tests ---

func TestManageReturnUseCase_GetReturn_Success(t *testing.T) {
	expected := &domain.Return{ID: "ret-1", OrderID: "ord-1", BuyerID: "b-1", SellerID: "s-1"}
	repo := &mockReturnRepo{
		getByIDFn: func(_ context.Context, id string) (*domain.Return, error) {
			assert.Equal(t, "ret-1", id)
			return expected, nil
		},
	}
	pub := &mockEventPublisher{}
	uc := NewManageReturnUseCase(repo, pub)

	result, err := uc.GetReturn(context.Background(), "ret-1")
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestManageReturnUseCase_GetReturn_NotFound(t *testing.T) {
	repo := &mockReturnRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Return, error) {
			return nil, errors.New("not found")
		},
	}
	pub := &mockEventPublisher{}
	uc := NewManageReturnUseCase(repo, pub)

	result, err := uc.GetReturn(context.Background(), "nonexistent")
	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// --- ListBuyerReturns pagination tests ---

func TestManageReturnUseCase_ListBuyerReturns_Pagination(t *testing.T) {
	tests := []struct {
		name             string
		page             int
		pageSize         int
		expectedPage     int
		expectedPageSize int
	}{
		{"valid values", 2, 10, 2, 10},
		{"page clamped to 1", 0, 10, 1, 10},
		{"negative page clamped to 1", -5, 10, 1, 10},
		{"pageSize 0 defaults to 20", 1, 0, 1, 20},
		{"pageSize over 100 defaults to 20", 1, 101, 1, 20},
		{"negative pageSize defaults to 20", 1, -1, 1, 20},
		{"pageSize exactly 1", 1, 1, 1, 1},
		{"pageSize exactly 100", 1, 100, 1, 100},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mockReturnRepo{
				listByBuyerFn: func(_ context.Context, buyerID string, page, pageSize int) ([]domain.Return, int64, error) {
					assert.Equal(t, "buyer-1", buyerID)
					assert.Equal(t, tc.expectedPage, page)
					assert.Equal(t, tc.expectedPageSize, pageSize)
					return []domain.Return{}, 0, nil
				},
			}
			pub := &mockEventPublisher{}
			uc := NewManageReturnUseCase(repo, pub)

			_, _, err := uc.ListBuyerReturns(context.Background(), "buyer-1", tc.page, tc.pageSize)
			require.NoError(t, err)
		})
	}
}

// --- ListSellerReturns pagination tests ---

func TestManageReturnUseCase_ListSellerReturns_Pagination(t *testing.T) {
	tests := []struct {
		name             string
		page             int
		pageSize         int
		expectedPage     int
		expectedPageSize int
	}{
		{"valid values", 3, 25, 3, 25},
		{"page clamped to 1", -1, 10, 1, 10},
		{"pageSize over 100 defaults to 20", 1, 200, 1, 20},
		{"pageSize 0 defaults to 20", 1, 0, 1, 20},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mockReturnRepo{
				listBySellerFn: func(_ context.Context, sellerID string, page, pageSize int) ([]domain.Return, int64, error) {
					assert.Equal(t, "seller-1", sellerID)
					assert.Equal(t, tc.expectedPage, page)
					assert.Equal(t, tc.expectedPageSize, pageSize)
					return []domain.Return{{ID: "ret-1"}}, 1, nil
				},
			}
			pub := &mockEventPublisher{}
			uc := NewManageReturnUseCase(repo, pub)

			results, total, err := uc.ListSellerReturns(context.Background(), "seller-1", tc.page, tc.pageSize)
			require.NoError(t, err)
			assert.Len(t, results, 1)
			assert.Equal(t, int64(1), total)
		})
	}
}

// --- ApproveReturn tests ---

func newRequestedReturn() *domain.Return {
	return &domain.Return{
		ID:                "ret-1",
		OrderID:           "ord-1",
		BuyerID:           "buyer-1",
		SellerID:          "seller-1",
		Status:            domain.ReturnStatusRequested,
		RefundAmountCents: 5000,
	}
}

func TestManageReturnUseCase_ApproveReturn_Success(t *testing.T) {
	ret := newRequestedReturn()
	var updated *domain.Return
	repo := &mockReturnRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Return, error) {
			return ret, nil
		},
		updateFn: func(_ context.Context, r *domain.Return) error {
			updated = r
			return nil
		},
	}
	pub := &mockEventPublisher{}
	uc := NewManageReturnUseCase(repo, pub)

	result, err := uc.ApproveReturn(context.Background(), "ret-1", "seller-1", 0)
	require.NoError(t, err)
	assert.Equal(t, domain.ReturnStatusApproved, result.Status)
	assert.Equal(t, int64(5000), result.RefundAmountCents) // unchanged when refundAmountCents=0
	require.NotNil(t, updated)

	require.Len(t, pub.calls, 1)
	assert.Equal(t, "return.approved", pub.calls[0].Subject)
}

func TestManageReturnUseCase_ApproveReturn_RefundAmountOverride(t *testing.T) {
	ret := newRequestedReturn()
	repo := &mockReturnRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Return, error) {
			return ret, nil
		},
		updateFn: func(_ context.Context, _ *domain.Return) error { return nil },
	}
	pub := &mockEventPublisher{}
	uc := NewManageReturnUseCase(repo, pub)

	result, err := uc.ApproveReturn(context.Background(), "ret-1", "seller-1", 3000)
	require.NoError(t, err)
	assert.Equal(t, int64(3000), result.RefundAmountCents)

	eventData := pub.calls[0].Data.(map[string]interface{})
	assert.Equal(t, int64(3000), eventData["refund_amount_cents"])
}

func TestManageReturnUseCase_ApproveReturn_WrongSeller(t *testing.T) {
	ret := newRequestedReturn()
	repo := &mockReturnRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Return, error) {
			return ret, nil
		},
	}
	pub := &mockEventPublisher{}
	uc := NewManageReturnUseCase(repo, pub)

	result, err := uc.ApproveReturn(context.Background(), "ret-1", "wrong-seller", 0)
	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized: return belongs to different seller")
	assert.Empty(t, pub.calls)
}

func TestManageReturnUseCase_ApproveReturn_InvalidTransition(t *testing.T) {
	ret := newRequestedReturn()
	ret.Status = domain.ReturnStatusApproved // already approved
	repo := &mockReturnRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Return, error) {
			return ret, nil
		},
	}
	pub := &mockEventPublisher{}
	uc := NewManageReturnUseCase(repo, pub)

	result, err := uc.ApproveReturn(context.Background(), "ret-1", "seller-1", 0)
	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot approve return in status approved")
	assert.Empty(t, pub.calls)
}

// --- RejectReturn tests ---

func TestManageReturnUseCase_RejectReturn_Success(t *testing.T) {
	ret := newRequestedReturn()
	repo := &mockReturnRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Return, error) {
			return ret, nil
		},
		updateFn: func(_ context.Context, _ *domain.Return) error { return nil },
	}
	pub := &mockEventPublisher{}
	uc := NewManageReturnUseCase(repo, pub)

	result, err := uc.RejectReturn(context.Background(), "ret-1", "seller-1")
	require.NoError(t, err)
	assert.Equal(t, domain.ReturnStatusRejected, result.Status)

	require.Len(t, pub.calls, 1)
	assert.Equal(t, "return.rejected", pub.calls[0].Subject)
	eventData := pub.calls[0].Data.(map[string]interface{})
	assert.Equal(t, "ret-1", eventData["return_id"])
	assert.Equal(t, "ord-1", eventData["order_id"])
}

func TestManageReturnUseCase_RejectReturn_WrongSeller(t *testing.T) {
	ret := newRequestedReturn()
	repo := &mockReturnRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Return, error) {
			return ret, nil
		},
	}
	pub := &mockEventPublisher{}
	uc := NewManageReturnUseCase(repo, pub)

	result, err := uc.RejectReturn(context.Background(), "ret-1", "wrong-seller")
	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized: return belongs to different seller")
}

func TestManageReturnUseCase_RejectReturn_InvalidTransition(t *testing.T) {
	ret := newRequestedReturn()
	ret.Status = domain.ReturnStatusApproved // can't reject after approval
	repo := &mockReturnRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Return, error) {
			return ret, nil
		},
	}
	pub := &mockEventPublisher{}
	uc := NewManageReturnUseCase(repo, pub)

	result, err := uc.RejectReturn(context.Background(), "ret-1", "seller-1")
	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot reject return in status approved")
}

// --- UpdateReturnStatus tests ---

func TestManageReturnUseCase_UpdateReturnStatus_Success(t *testing.T) {
	ret := &domain.Return{
		ID:       "ret-1",
		OrderID:  "ord-1",
		SellerID: "seller-1",
		Status:   domain.ReturnStatusApproved,
	}
	repo := &mockReturnRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Return, error) {
			return ret, nil
		},
		updateFn: func(_ context.Context, _ *domain.Return) error { return nil },
	}
	pub := &mockEventPublisher{}
	uc := NewManageReturnUseCase(repo, pub)

	result, err := uc.UpdateReturnStatus(context.Background(), "ret-1", "seller-1", domain.ReturnStatusShippedBack)
	require.NoError(t, err)
	assert.Equal(t, domain.ReturnStatusShippedBack, result.Status)

	// Non-refunded transitions should not publish return.completed
	assert.Empty(t, pub.calls)
}

func TestManageReturnUseCase_UpdateReturnStatus_RefundedPublishesCompletedEvent(t *testing.T) {
	ret := &domain.Return{
		ID:                "ret-1",
		OrderID:           "ord-1",
		SellerID:          "seller-1",
		Status:            domain.ReturnStatusReceived,
		RefundAmountCents: 7500,
	}
	repo := &mockReturnRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Return, error) {
			return ret, nil
		},
		updateFn: func(_ context.Context, _ *domain.Return) error { return nil },
	}
	pub := &mockEventPublisher{}
	uc := NewManageReturnUseCase(repo, pub)

	result, err := uc.UpdateReturnStatus(context.Background(), "ret-1", "seller-1", domain.ReturnStatusRefunded)
	require.NoError(t, err)
	assert.Equal(t, domain.ReturnStatusRefunded, result.Status)

	require.Len(t, pub.calls, 1)
	assert.Equal(t, "return.completed", pub.calls[0].Subject)
	eventData := pub.calls[0].Data.(map[string]interface{})
	assert.Equal(t, "ret-1", eventData["return_id"])
	assert.Equal(t, "ord-1", eventData["order_id"])
	assert.Equal(t, int64(7500), eventData["refund_amount_cents"])
}

func TestManageReturnUseCase_UpdateReturnStatus_WrongSeller(t *testing.T) {
	ret := &domain.Return{
		ID:       "ret-1",
		SellerID: "seller-1",
		Status:   domain.ReturnStatusApproved,
	}
	repo := &mockReturnRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Return, error) {
			return ret, nil
		},
	}
	pub := &mockEventPublisher{}
	uc := NewManageReturnUseCase(repo, pub)

	result, err := uc.UpdateReturnStatus(context.Background(), "ret-1", "other-seller", domain.ReturnStatusShippedBack)
	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized: return belongs to different seller")
}

func TestManageReturnUseCase_UpdateReturnStatus_InvalidTransition(t *testing.T) {
	ret := &domain.Return{
		ID:       "ret-1",
		SellerID: "seller-1",
		Status:   domain.ReturnStatusRequested,
	}
	repo := &mockReturnRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Return, error) {
			return ret, nil
		},
	}
	pub := &mockEventPublisher{}
	uc := NewManageReturnUseCase(repo, pub)

	result, err := uc.UpdateReturnStatus(context.Background(), "ret-1", "seller-1", domain.ReturnStatusRefunded)
	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid transition from requested to refunded")
}
