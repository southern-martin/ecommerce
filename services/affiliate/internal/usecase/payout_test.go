package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/affiliate/internal/domain"
)

// ===========================================================================
// RequestPayout tests
// ===========================================================================

func TestRequestPayout_Success(t *testing.T) {
	pRepo, lRepo, _, payRepo, pub := defaultAffiliateMocks()

	pRepo.getFn = func(_ context.Context) (*domain.AffiliateProgram, error) {
		return &domain.AffiliateProgram{
			ID:             "prog-1",
			MinPayoutCents: 5000, // $50 minimum
		}, nil
	}

	lRepo.listByUserFn = func(_ context.Context, _ string, _, _ int) ([]domain.AffiliateLink, int64, error) {
		return []domain.AffiliateLink{
			{ID: "link-1", TotalEarningsCents: 8000},
			{ID: "link-2", TotalEarningsCents: 2000},
		}, 2, nil
	}

	var saved *domain.AffiliatePayout
	payRepo.createFn = func(_ context.Context, p *domain.AffiliatePayout) error {
		saved = p
		return nil
	}

	uc := NewPayoutUseCase(payRepo, pRepo, lRepo, pub)
	payout, err := uc.RequestPayout(context.Background(), RequestPayoutRequest{
		UserID:       "user-1",
		AmountCents:  7000,
		PayoutMethod: domain.PayoutMethodBankTransfer,
	})

	require.NoError(t, err)
	require.NotNil(t, payout)
	assert.NotEmpty(t, payout.ID)
	assert.Equal(t, "user-1", payout.UserID)
	assert.Equal(t, int64(7000), payout.AmountCents)
	assert.Equal(t, domain.PayoutStatusRequested, payout.Status)
	assert.Equal(t, domain.PayoutMethodBankTransfer, payout.PayoutMethod)
	assert.NotNil(t, saved)
}

func TestRequestPayout_BelowMinimum(t *testing.T) {
	pRepo, lRepo, _, payRepo, pub := defaultAffiliateMocks()

	pRepo.getFn = func(_ context.Context) (*domain.AffiliateProgram, error) {
		return &domain.AffiliateProgram{MinPayoutCents: 5000}, nil
	}

	uc := NewPayoutUseCase(payRepo, pRepo, lRepo, pub)
	_, err := uc.RequestPayout(context.Background(), RequestPayoutRequest{
		UserID:      "user-1",
		AmountCents: 3000, // below 5000 minimum
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "below minimum")
}

func TestRequestPayout_ExceedsEarnings(t *testing.T) {
	pRepo, lRepo, _, payRepo, pub := defaultAffiliateMocks()

	pRepo.getFn = func(_ context.Context) (*domain.AffiliateProgram, error) {
		return &domain.AffiliateProgram{MinPayoutCents: 1000}, nil
	}

	lRepo.listByUserFn = func(_ context.Context, _ string, _, _ int) ([]domain.AffiliateLink, int64, error) {
		return []domain.AffiliateLink{
			{ID: "link-1", TotalEarningsCents: 3000},
		}, 1, nil
	}

	uc := NewPayoutUseCase(payRepo, pRepo, lRepo, pub)
	_, err := uc.RequestPayout(context.Background(), RequestPayoutRequest{
		UserID:      "user-1",
		AmountCents: 5000, // exceeds 3000 total earnings
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds total earnings")
}

func TestRequestPayout_ProgramGetError(t *testing.T) {
	pRepo, lRepo, _, payRepo, pub := defaultAffiliateMocks()

	pRepo.getFn = func(_ context.Context) (*domain.AffiliateProgram, error) {
		return nil, errors.New("db error")
	}

	uc := NewPayoutUseCase(payRepo, pRepo, lRepo, pub)
	_, err := uc.RequestPayout(context.Background(), RequestPayoutRequest{
		UserID:      "user-1",
		AmountCents: 5000,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get affiliate program")
}

func TestRequestPayout_CreateError(t *testing.T) {
	pRepo, lRepo, _, payRepo, pub := defaultAffiliateMocks()

	pRepo.getFn = func(_ context.Context) (*domain.AffiliateProgram, error) {
		return &domain.AffiliateProgram{MinPayoutCents: 1000}, nil
	}
	lRepo.listByUserFn = func(_ context.Context, _ string, _, _ int) ([]domain.AffiliateLink, int64, error) {
		return []domain.AffiliateLink{{TotalEarningsCents: 10000}}, 1, nil
	}
	payRepo.createFn = func(_ context.Context, _ *domain.AffiliatePayout) error {
		return errors.New("db error")
	}

	uc := NewPayoutUseCase(payRepo, pRepo, lRepo, pub)
	_, err := uc.RequestPayout(context.Background(), RequestPayoutRequest{
		UserID:      "user-1",
		AmountCents: 5000,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create payout")
}

// ===========================================================================
// UpdatePayoutStatus tests
// ===========================================================================

func TestUpdatePayoutStatus_Success(t *testing.T) {
	pRepo, lRepo, _, payRepo, pub := defaultAffiliateMocks()

	var capturedStatus domain.PayoutStatus
	var capturedCompletedAt *time.Time
	payRepo.updateStatusFn = func(_ context.Context, _ string, status domain.PayoutStatus, completedAt *time.Time) error {
		capturedStatus = status
		capturedCompletedAt = completedAt
		return nil
	}

	payRepo.getByIDFn = func(_ context.Context, id string) (*domain.AffiliatePayout, error) {
		return &domain.AffiliatePayout{
			ID:     id,
			Status: domain.PayoutStatusProcessing,
		}, nil
	}

	uc := NewPayoutUseCase(payRepo, pRepo, lRepo, pub)
	payout, err := uc.UpdatePayoutStatus(context.Background(), "payout-1", domain.PayoutStatusProcessing)

	require.NoError(t, err)
	require.NotNil(t, payout)
	assert.Equal(t, domain.PayoutStatusProcessing, capturedStatus)
	assert.Nil(t, capturedCompletedAt) // not completed, no completedAt
}

func TestUpdatePayoutStatus_Completed(t *testing.T) {
	pRepo, lRepo, _, payRepo, pub := defaultAffiliateMocks()

	var capturedCompletedAt *time.Time
	payRepo.updateStatusFn = func(_ context.Context, _ string, _ domain.PayoutStatus, completedAt *time.Time) error {
		capturedCompletedAt = completedAt
		return nil
	}
	payRepo.getByIDFn = func(_ context.Context, id string) (*domain.AffiliatePayout, error) {
		return &domain.AffiliatePayout{ID: id, Status: domain.PayoutStatusCompleted}, nil
	}

	uc := NewPayoutUseCase(payRepo, pRepo, lRepo, pub)
	_, err := uc.UpdatePayoutStatus(context.Background(), "payout-1", domain.PayoutStatusCompleted)

	require.NoError(t, err)
	assert.NotNil(t, capturedCompletedAt, "completedAt should be set when status is completed")
}

func TestUpdatePayoutStatus_UpdateError(t *testing.T) {
	pRepo, lRepo, _, payRepo, pub := defaultAffiliateMocks()

	payRepo.updateStatusFn = func(_ context.Context, _ string, _ domain.PayoutStatus, _ *time.Time) error {
		return errors.New("db error")
	}

	uc := NewPayoutUseCase(payRepo, pRepo, lRepo, pub)
	_, err := uc.UpdatePayoutStatus(context.Background(), "payout-1", domain.PayoutStatusFailed)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update payout status")
}

// ===========================================================================
// GetPayout tests
// ===========================================================================

func TestGetPayout_Success(t *testing.T) {
	pRepo, lRepo, _, payRepo, pub := defaultAffiliateMocks()

	payRepo.getByIDFn = func(_ context.Context, id string) (*domain.AffiliatePayout, error) {
		return &domain.AffiliatePayout{ID: id, UserID: "user-1", AmountCents: 5000}, nil
	}

	uc := NewPayoutUseCase(payRepo, pRepo, lRepo, pub)
	payout, err := uc.GetPayout(context.Background(), "payout-1")

	require.NoError(t, err)
	require.NotNil(t, payout)
	assert.Equal(t, "payout-1", payout.ID)
	assert.Equal(t, int64(5000), payout.AmountCents)
}

func TestGetPayout_NotFound(t *testing.T) {
	pRepo, lRepo, _, payRepo, pub := defaultAffiliateMocks()

	uc := NewPayoutUseCase(payRepo, pRepo, lRepo, pub)
	_, err := uc.GetPayout(context.Background(), "payout-missing")

	require.Error(t, err)
}

// ===========================================================================
// ListPayouts tests
// ===========================================================================

func TestListPayouts_Success(t *testing.T) {
	pRepo, lRepo, _, payRepo, pub := defaultAffiliateMocks()

	payRepo.listByUserFn = func(_ context.Context, userID string, _, _ int) ([]domain.AffiliatePayout, int64, error) {
		return []domain.AffiliatePayout{
			{ID: "pay-1", UserID: userID, AmountCents: 5000},
		}, 1, nil
	}

	uc := NewPayoutUseCase(payRepo, pRepo, lRepo, pub)
	payouts, total, err := uc.ListPayouts(context.Background(), "user-1", 1, 20)

	require.NoError(t, err)
	assert.Len(t, payouts, 1)
	assert.Equal(t, int64(1), total)
}

func TestListPayouts_DefaultsPagination(t *testing.T) {
	pRepo, lRepo, _, payRepo, pub := defaultAffiliateMocks()

	var capturedPage, capturedPageSize int
	payRepo.listByUserFn = func(_ context.Context, _ string, page, pageSize int) ([]domain.AffiliatePayout, int64, error) {
		capturedPage = page
		capturedPageSize = pageSize
		return nil, 0, nil
	}

	uc := NewPayoutUseCase(payRepo, pRepo, lRepo, pub)

	_, _, _ = uc.ListPayouts(context.Background(), "user-1", 0, 20)
	assert.Equal(t, 1, capturedPage)

	_, _, _ = uc.ListPayouts(context.Background(), "user-1", 1, 0)
	assert.Equal(t, 20, capturedPageSize)
}

// ===========================================================================
// ListAllPayouts tests
// ===========================================================================

func TestListAllPayouts_Success(t *testing.T) {
	pRepo, lRepo, _, payRepo, pub := defaultAffiliateMocks()

	payRepo.listAllFn = func(_ context.Context, _, _ int) ([]domain.AffiliatePayout, int64, error) {
		return []domain.AffiliatePayout{
			{ID: "pay-1", AmountCents: 5000},
			{ID: "pay-2", AmountCents: 3000},
		}, 2, nil
	}

	uc := NewPayoutUseCase(payRepo, pRepo, lRepo, pub)
	payouts, total, err := uc.ListAllPayouts(context.Background(), 1, 20)

	require.NoError(t, err)
	assert.Len(t, payouts, 2)
	assert.Equal(t, int64(2), total)
}

func TestListAllPayouts_DefaultsPagination(t *testing.T) {
	pRepo, lRepo, _, payRepo, pub := defaultAffiliateMocks()

	var capturedPage, capturedPageSize int
	payRepo.listAllFn = func(_ context.Context, page, pageSize int) ([]domain.AffiliatePayout, int64, error) {
		capturedPage = page
		capturedPageSize = pageSize
		return nil, 0, nil
	}

	uc := NewPayoutUseCase(payRepo, pRepo, lRepo, pub)

	_, _, _ = uc.ListAllPayouts(context.Background(), -5, 20)
	assert.Equal(t, 1, capturedPage)

	_, _, _ = uc.ListAllPayouts(context.Background(), 1, 999)
	assert.Equal(t, 20, capturedPageSize)
}
