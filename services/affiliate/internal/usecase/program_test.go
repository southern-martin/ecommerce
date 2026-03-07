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

// ---------------------------------------------------------------------------
// Hand-written function-field mocks
// ---------------------------------------------------------------------------

// --- AffiliateProgramRepository mock ---

type mockProgramRepo struct {
	getFn    func(ctx context.Context) (*domain.AffiliateProgram, error)
	createFn func(ctx context.Context, program *domain.AffiliateProgram) error
	updateFn func(ctx context.Context, program *domain.AffiliateProgram) error
}

func (m *mockProgramRepo) Get(ctx context.Context) (*domain.AffiliateProgram, error) {
	if m.getFn != nil {
		return m.getFn(ctx)
	}
	return nil, errors.New("not found")
}
func (m *mockProgramRepo) Create(ctx context.Context, program *domain.AffiliateProgram) error {
	if m.createFn != nil {
		return m.createFn(ctx, program)
	}
	return nil
}
func (m *mockProgramRepo) Update(ctx context.Context, program *domain.AffiliateProgram) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, program)
	}
	return nil
}

// --- AffiliateLinkRepository mock ---

type mockLinkRepo struct {
	getByIDFn            func(ctx context.Context, id string) (*domain.AffiliateLink, error)
	getByCodeFn          func(ctx context.Context, code string) (*domain.AffiliateLink, error)
	listByUserFn         func(ctx context.Context, userID string, page, pageSize int) ([]domain.AffiliateLink, int64, error)
	createFn             func(ctx context.Context, link *domain.AffiliateLink) error
	incrementClicksFn    func(ctx context.Context, id string) error
	incrementConversionsFn func(ctx context.Context, id string) error
	addEarningsFn        func(ctx context.Context, id string, amountCents int64) error
}

func (m *mockLinkRepo) GetByID(ctx context.Context, id string) (*domain.AffiliateLink, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}
func (m *mockLinkRepo) GetByCode(ctx context.Context, code string) (*domain.AffiliateLink, error) {
	if m.getByCodeFn != nil {
		return m.getByCodeFn(ctx, code)
	}
	return nil, errors.New("not found")
}
func (m *mockLinkRepo) ListByUser(ctx context.Context, userID string, page, pageSize int) ([]domain.AffiliateLink, int64, error) {
	if m.listByUserFn != nil {
		return m.listByUserFn(ctx, userID, page, pageSize)
	}
	return nil, 0, nil
}
func (m *mockLinkRepo) Create(ctx context.Context, link *domain.AffiliateLink) error {
	if m.createFn != nil {
		return m.createFn(ctx, link)
	}
	return nil
}
func (m *mockLinkRepo) IncrementClicks(ctx context.Context, id string) error {
	if m.incrementClicksFn != nil {
		return m.incrementClicksFn(ctx, id)
	}
	return nil
}
func (m *mockLinkRepo) IncrementConversions(ctx context.Context, id string) error {
	if m.incrementConversionsFn != nil {
		return m.incrementConversionsFn(ctx, id)
	}
	return nil
}
func (m *mockLinkRepo) AddEarnings(ctx context.Context, id string, amountCents int64) error {
	if m.addEarningsFn != nil {
		return m.addEarningsFn(ctx, id, amountCents)
	}
	return nil
}

// --- ReferralRepository mock ---

type mockReferralRepo struct {
	getByIDFn       func(ctx context.Context, id string) (*domain.Referral, error)
	listByReferrerFn func(ctx context.Context, referrerID string, page, pageSize int) ([]domain.Referral, int64, error)
	listByReferredFn func(ctx context.Context, referredID string) ([]domain.Referral, error)
	createFn        func(ctx context.Context, referral *domain.Referral) error
	updateStatusFn  func(ctx context.Context, id string, status domain.ReferralStatus) error
}

func (m *mockReferralRepo) GetByID(ctx context.Context, id string) (*domain.Referral, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}
func (m *mockReferralRepo) ListByReferrer(ctx context.Context, referrerID string, page, pageSize int) ([]domain.Referral, int64, error) {
	if m.listByReferrerFn != nil {
		return m.listByReferrerFn(ctx, referrerID, page, pageSize)
	}
	return nil, 0, nil
}
func (m *mockReferralRepo) ListByReferred(ctx context.Context, referredID string) ([]domain.Referral, error) {
	if m.listByReferredFn != nil {
		return m.listByReferredFn(ctx, referredID)
	}
	return nil, nil
}
func (m *mockReferralRepo) Create(ctx context.Context, referral *domain.Referral) error {
	if m.createFn != nil {
		return m.createFn(ctx, referral)
	}
	return nil
}
func (m *mockReferralRepo) UpdateStatus(ctx context.Context, id string, status domain.ReferralStatus) error {
	if m.updateStatusFn != nil {
		return m.updateStatusFn(ctx, id, status)
	}
	return nil
}

// --- PayoutRepository mock ---

type mockPayoutRepo struct {
	getByIDFn      func(ctx context.Context, id string) (*domain.AffiliatePayout, error)
	listByUserFn   func(ctx context.Context, userID string, page, pageSize int) ([]domain.AffiliatePayout, int64, error)
	listAllFn      func(ctx context.Context, page, pageSize int) ([]domain.AffiliatePayout, int64, error)
	createFn       func(ctx context.Context, payout *domain.AffiliatePayout) error
	updateStatusFn func(ctx context.Context, id string, status domain.PayoutStatus, completedAt *time.Time) error
}

func (m *mockPayoutRepo) GetByID(ctx context.Context, id string) (*domain.AffiliatePayout, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}
func (m *mockPayoutRepo) ListByUser(ctx context.Context, userID string, page, pageSize int) ([]domain.AffiliatePayout, int64, error) {
	if m.listByUserFn != nil {
		return m.listByUserFn(ctx, userID, page, pageSize)
	}
	return nil, 0, nil
}
func (m *mockPayoutRepo) ListAll(ctx context.Context, page, pageSize int) ([]domain.AffiliatePayout, int64, error) {
	if m.listAllFn != nil {
		return m.listAllFn(ctx, page, pageSize)
	}
	return nil, 0, nil
}
func (m *mockPayoutRepo) Create(ctx context.Context, payout *domain.AffiliatePayout) error {
	if m.createFn != nil {
		return m.createFn(ctx, payout)
	}
	return nil
}
func (m *mockPayoutRepo) UpdateStatus(ctx context.Context, id string, status domain.PayoutStatus, completedAt *time.Time) error {
	if m.updateStatusFn != nil {
		return m.updateStatusFn(ctx, id, status, completedAt)
	}
	return nil
}

// --- EventPublisher mock ---

type mockAffiliateEventPub struct {
	publishFn func(ctx context.Context, subject string, data interface{}) error
}

func (m *mockAffiliateEventPub) Publish(ctx context.Context, subject string, data interface{}) error {
	if m.publishFn != nil {
		return m.publishFn(ctx, subject, data)
	}
	return nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func defaultAffiliateMocks() (*mockProgramRepo, *mockLinkRepo, *mockReferralRepo, *mockPayoutRepo, *mockAffiliateEventPub) {
	return &mockProgramRepo{},
		&mockLinkRepo{},
		&mockReferralRepo{},
		&mockPayoutRepo{},
		&mockAffiliateEventPub{}
}

// ===========================================================================
// GetProgram tests
// ===========================================================================

func TestGetProgram_Success(t *testing.T) {
	pRepo, _, _, _, _ := defaultAffiliateMocks()

	pRepo.getFn = func(_ context.Context) (*domain.AffiliateProgram, error) {
		return &domain.AffiliateProgram{
			ID:             "prog-1",
			CommissionRate: 0.1,
			MinPayoutCents: 5000,
			IsActive:       true,
		}, nil
	}

	uc := NewProgramUseCase(pRepo)
	program, err := uc.GetProgram(context.Background())

	require.NoError(t, err)
	require.NotNil(t, program)
	assert.Equal(t, "prog-1", program.ID)
	assert.Equal(t, 0.1, program.CommissionRate)
	assert.True(t, program.IsActive)
}

func TestGetProgram_NotFound(t *testing.T) {
	pRepo, _, _, _, _ := defaultAffiliateMocks()

	pRepo.getFn = func(_ context.Context) (*domain.AffiliateProgram, error) {
		return nil, errors.New("not found")
	}

	uc := NewProgramUseCase(pRepo)
	_, err := uc.GetProgram(context.Background())

	require.Error(t, err)
}

// ===========================================================================
// UpdateProgram tests
// ===========================================================================

func TestUpdateProgram_Success(t *testing.T) {
	pRepo, _, _, _, _ := defaultAffiliateMocks()

	pRepo.getFn = func(_ context.Context) (*domain.AffiliateProgram, error) {
		return &domain.AffiliateProgram{
			ID:             "prog-1",
			CommissionRate: 0.1,
			MinPayoutCents: 5000,
			CookieDays:     30,
			IsActive:       true,
		}, nil
	}

	var updated *domain.AffiliateProgram
	pRepo.updateFn = func(_ context.Context, p *domain.AffiliateProgram) error {
		updated = p
		return nil
	}

	uc := NewProgramUseCase(pRepo)
	newRate := 0.15
	newMinPayout := int64(10000)
	program, err := uc.UpdateProgram(context.Background(), UpdateProgramRequest{
		ID:             "prog-1",
		CommissionRate: &newRate,
		MinPayoutCents: &newMinPayout,
	})

	require.NoError(t, err)
	require.NotNil(t, program)
	assert.Equal(t, 0.15, program.CommissionRate)
	assert.Equal(t, int64(10000), program.MinPayoutCents)
	assert.Equal(t, 30, program.CookieDays) // unchanged
	assert.NotNil(t, updated)
}

func TestUpdateProgram_PartialUpdate(t *testing.T) {
	pRepo, _, _, _, _ := defaultAffiliateMocks()

	pRepo.getFn = func(_ context.Context) (*domain.AffiliateProgram, error) {
		return &domain.AffiliateProgram{
			ID:             "prog-1",
			CommissionRate: 0.1,
			MinPayoutCents: 5000,
			CookieDays:     30,
			IsActive:       true,
		}, nil
	}
	pRepo.updateFn = func(_ context.Context, _ *domain.AffiliateProgram) error { return nil }

	uc := NewProgramUseCase(pRepo)
	isActive := false
	program, err := uc.UpdateProgram(context.Background(), UpdateProgramRequest{
		IsActive: &isActive,
	})

	require.NoError(t, err)
	assert.False(t, program.IsActive)
	assert.Equal(t, 0.1, program.CommissionRate)  // unchanged
	assert.Equal(t, int64(5000), program.MinPayoutCents) // unchanged
	assert.Equal(t, 30, program.CookieDays)        // unchanged
}

func TestUpdateProgram_GetError(t *testing.T) {
	pRepo, _, _, _, _ := defaultAffiliateMocks()

	pRepo.getFn = func(_ context.Context) (*domain.AffiliateProgram, error) {
		return nil, errors.New("not found")
	}

	uc := NewProgramUseCase(pRepo)
	_, err := uc.UpdateProgram(context.Background(), UpdateProgramRequest{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get program")
}

func TestUpdateProgram_UpdateError(t *testing.T) {
	pRepo, _, _, _, _ := defaultAffiliateMocks()

	pRepo.getFn = func(_ context.Context) (*domain.AffiliateProgram, error) {
		return &domain.AffiliateProgram{ID: "prog-1"}, nil
	}
	pRepo.updateFn = func(_ context.Context, _ *domain.AffiliateProgram) error {
		return errors.New("db error")
	}

	uc := NewProgramUseCase(pRepo)
	_, err := uc.UpdateProgram(context.Background(), UpdateProgramRequest{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update program")
}
