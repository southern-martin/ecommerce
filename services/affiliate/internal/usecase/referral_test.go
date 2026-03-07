package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/affiliate/internal/domain"
)

// ===========================================================================
// TrackConversion tests
// ===========================================================================

func TestTrackConversion_Success(t *testing.T) {
	pRepo, lRepo, rRepo, _, pub := defaultAffiliateMocks()

	lRepo.getByIDFn = func(_ context.Context, id string) (*domain.AffiliateLink, error) {
		return &domain.AffiliateLink{
			ID:     id,
			UserID: "referrer-1",
			Code:   "abc12345",
		}, nil
	}

	pRepo.getFn = func(_ context.Context) (*domain.AffiliateProgram, error) {
		return &domain.AffiliateProgram{
			ID:             "prog-1",
			CommissionRate: 0.10, // 10%
			IsActive:       true,
		}, nil
	}

	var savedReferral *domain.Referral
	rRepo.createFn = func(_ context.Context, r *domain.Referral) error {
		savedReferral = r
		return nil
	}

	uc := NewReferralUseCase(rRepo, lRepo, pRepo, pub)
	referral, err := uc.TrackConversion(context.Background(), TrackConversionRequest{
		LinkID:          "link-1",
		ReferredID:      "referred-1",
		OrderID:         "order-1",
		OrderTotalCents: 10000, // $100.00
	})

	require.NoError(t, err)
	require.NotNil(t, referral)
	assert.NotEmpty(t, referral.ID)
	assert.Equal(t, "referrer-1", referral.ReferrerID)
	assert.Equal(t, "referred-1", referral.ReferredID)
	assert.Equal(t, "order-1", referral.OrderID)
	assert.Equal(t, int64(10000), referral.OrderTotalCents)
	assert.Equal(t, int64(1000), referral.CommissionCents) // 10% of 10000
	assert.Equal(t, domain.ReferralStatusPending, referral.Status)
	assert.NotNil(t, savedReferral)
}

func TestTrackConversion_CommissionCalculation(t *testing.T) {
	pRepo, lRepo, rRepo, _, pub := defaultAffiliateMocks()

	lRepo.getByIDFn = func(_ context.Context, _ string) (*domain.AffiliateLink, error) {
		return &domain.AffiliateLink{ID: "link-1", UserID: "referrer-1"}, nil
	}

	pRepo.getFn = func(_ context.Context) (*domain.AffiliateProgram, error) {
		return &domain.AffiliateProgram{
			CommissionRate: 0.05, // 5%
		}, nil
	}

	rRepo.createFn = func(_ context.Context, _ *domain.Referral) error { return nil }

	uc := NewReferralUseCase(rRepo, lRepo, pRepo, pub)
	referral, err := uc.TrackConversion(context.Background(), TrackConversionRequest{
		LinkID:          "link-1",
		ReferredID:      "referred-1",
		OrderID:         "order-1",
		OrderTotalCents: 7999, // $79.99
	})

	require.NoError(t, err)
	assert.Equal(t, int64(400), referral.CommissionCents) // 5% of 7999 = 399.95, rounded to 400
}

func TestTrackConversion_LinkNotFound(t *testing.T) {
	pRepo, lRepo, rRepo, _, pub := defaultAffiliateMocks()

	lRepo.getByIDFn = func(_ context.Context, _ string) (*domain.AffiliateLink, error) {
		return nil, errors.New("not found")
	}

	uc := NewReferralUseCase(rRepo, lRepo, pRepo, pub)
	_, err := uc.TrackConversion(context.Background(), TrackConversionRequest{
		LinkID: "link-missing",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "affiliate link not found")
}

func TestTrackConversion_ProgramGetError(t *testing.T) {
	pRepo, lRepo, rRepo, _, pub := defaultAffiliateMocks()

	lRepo.getByIDFn = func(_ context.Context, _ string) (*domain.AffiliateLink, error) {
		return &domain.AffiliateLink{ID: "link-1", UserID: "referrer-1"}, nil
	}
	pRepo.getFn = func(_ context.Context) (*domain.AffiliateProgram, error) {
		return nil, errors.New("db error")
	}

	uc := NewReferralUseCase(rRepo, lRepo, pRepo, pub)
	_, err := uc.TrackConversion(context.Background(), TrackConversionRequest{
		LinkID: "link-1",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get affiliate program")
}

func TestTrackConversion_CreateError(t *testing.T) {
	pRepo, lRepo, rRepo, _, pub := defaultAffiliateMocks()

	lRepo.getByIDFn = func(_ context.Context, _ string) (*domain.AffiliateLink, error) {
		return &domain.AffiliateLink{ID: "link-1", UserID: "referrer-1"}, nil
	}
	pRepo.getFn = func(_ context.Context) (*domain.AffiliateProgram, error) {
		return &domain.AffiliateProgram{CommissionRate: 0.1}, nil
	}
	rRepo.createFn = func(_ context.Context, _ *domain.Referral) error {
		return errors.New("db error")
	}

	uc := NewReferralUseCase(rRepo, lRepo, pRepo, pub)
	_, err := uc.TrackConversion(context.Background(), TrackConversionRequest{
		LinkID:          "link-1",
		OrderTotalCents: 5000,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create referral")
}

// ===========================================================================
// ListReferrals tests
// ===========================================================================

func TestListReferrals_Success(t *testing.T) {
	pRepo, lRepo, rRepo, _, pub := defaultAffiliateMocks()

	rRepo.listByReferrerFn = func(_ context.Context, referrerID string, _, _ int) ([]domain.Referral, int64, error) {
		return []domain.Referral{
			{ID: "ref-1", ReferrerID: referrerID, CommissionCents: 500},
			{ID: "ref-2", ReferrerID: referrerID, CommissionCents: 300},
		}, 2, nil
	}

	uc := NewReferralUseCase(rRepo, lRepo, pRepo, pub)
	referrals, total, err := uc.ListReferrals(context.Background(), "referrer-1", 1, 20)

	require.NoError(t, err)
	assert.Len(t, referrals, 2)
	assert.Equal(t, int64(2), total)
}

func TestListReferrals_DefaultsPagination(t *testing.T) {
	pRepo, lRepo, rRepo, _, pub := defaultAffiliateMocks()

	var capturedPage, capturedPageSize int
	rRepo.listByReferrerFn = func(_ context.Context, _ string, page, pageSize int) ([]domain.Referral, int64, error) {
		capturedPage = page
		capturedPageSize = pageSize
		return nil, 0, nil
	}

	uc := NewReferralUseCase(rRepo, lRepo, pRepo, pub)

	_, _, _ = uc.ListReferrals(context.Background(), "referrer-1", -1, 20)
	assert.Equal(t, 1, capturedPage)

	_, _, _ = uc.ListReferrals(context.Background(), "referrer-1", 1, 0)
	assert.Equal(t, 20, capturedPageSize)

	_, _, _ = uc.ListReferrals(context.Background(), "referrer-1", 1, 150)
	assert.Equal(t, 20, capturedPageSize)
}
