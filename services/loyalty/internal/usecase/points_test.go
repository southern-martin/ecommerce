package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/loyalty/internal/domain"
)

// ---------------------------------------------------------------------------
// Hand-written function-field mocks
// ---------------------------------------------------------------------------

// --- MembershipRepository mock ---

type mockMembershipRepo struct {
	getByUserIDFn  func(ctx context.Context, userID string) (*domain.Membership, error)
	createFn       func(ctx context.Context, membership *domain.Membership) error
	updateFn       func(ctx context.Context, membership *domain.Membership) error
	updateTierFn   func(ctx context.Context, userID string, tier domain.MemberTier) error
	updatePointsFn func(ctx context.Context, userID string, pointsBalance, lifetimePoints int64) error
}

func (m *mockMembershipRepo) GetByUserID(ctx context.Context, userID string) (*domain.Membership, error) {
	if m.getByUserIDFn != nil {
		return m.getByUserIDFn(ctx, userID)
	}
	return nil, errors.New("not found")
}
func (m *mockMembershipRepo) Create(ctx context.Context, membership *domain.Membership) error {
	if m.createFn != nil {
		return m.createFn(ctx, membership)
	}
	return nil
}
func (m *mockMembershipRepo) Update(ctx context.Context, membership *domain.Membership) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, membership)
	}
	return nil
}
func (m *mockMembershipRepo) UpdateTier(ctx context.Context, userID string, tier domain.MemberTier) error {
	if m.updateTierFn != nil {
		return m.updateTierFn(ctx, userID, tier)
	}
	return nil
}
func (m *mockMembershipRepo) UpdatePoints(ctx context.Context, userID string, pointsBalance, lifetimePoints int64) error {
	if m.updatePointsFn != nil {
		return m.updatePointsFn(ctx, userID, pointsBalance, lifetimePoints)
	}
	return nil
}

// --- PointsTransactionRepository mock ---

type mockPointsTxRepo struct {
	getByIDFn  func(ctx context.Context, id string) (*domain.PointsTransaction, error)
	listByUserFn func(ctx context.Context, userID string, page, pageSize int) ([]domain.PointsTransaction, int64, error)
	createFn   func(ctx context.Context, tx *domain.PointsTransaction) error
}

func (m *mockPointsTxRepo) GetByID(ctx context.Context, id string) (*domain.PointsTransaction, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}
func (m *mockPointsTxRepo) ListByUser(ctx context.Context, userID string, page, pageSize int) ([]domain.PointsTransaction, int64, error) {
	if m.listByUserFn != nil {
		return m.listByUserFn(ctx, userID, page, pageSize)
	}
	return nil, 0, nil
}
func (m *mockPointsTxRepo) Create(ctx context.Context, tx *domain.PointsTransaction) error {
	if m.createFn != nil {
		return m.createFn(ctx, tx)
	}
	return nil
}

// --- TierRepository mock ---

type mockTierRepo struct {
	getAllFn          func(ctx context.Context) ([]domain.Tier, error)
	getByNameFn      func(ctx context.Context, name string) (*domain.Tier, error)
	getTierForPointsFn func(ctx context.Context, lifetimePoints int64) (*domain.Tier, error)
}

func (m *mockTierRepo) GetAll(ctx context.Context) ([]domain.Tier, error) {
	if m.getAllFn != nil {
		return m.getAllFn(ctx)
	}
	return nil, nil
}
func (m *mockTierRepo) GetByName(ctx context.Context, name string) (*domain.Tier, error) {
	if m.getByNameFn != nil {
		return m.getByNameFn(ctx, name)
	}
	return nil, errors.New("not found")
}
func (m *mockTierRepo) GetTierForPoints(ctx context.Context, lifetimePoints int64) (*domain.Tier, error) {
	if m.getTierForPointsFn != nil {
		return m.getTierForPointsFn(ctx, lifetimePoints)
	}
	return nil, nil
}

// --- EventPublisher mock ---

type mockLoyaltyEventPub struct {
	publishFn func(ctx context.Context, subject string, data interface{}) error
}

func (m *mockLoyaltyEventPub) Publish(ctx context.Context, subject string, data interface{}) error {
	if m.publishFn != nil {
		return m.publishFn(ctx, subject, data)
	}
	return nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func defaultLoyaltyMocks() (*mockMembershipRepo, *mockPointsTxRepo, *mockTierRepo, *mockLoyaltyEventPub) {
	return &mockMembershipRepo{},
		&mockPointsTxRepo{},
		&mockTierRepo{},
		&mockLoyaltyEventPub{}
}

func newPointsUseCase(
	membershipRepo *mockMembershipRepo,
	txRepo *mockPointsTxRepo,
	tierRepo *mockTierRepo,
	pub *mockLoyaltyEventPub,
) *PointsUseCase {
	membershipUC := NewMembershipUseCase(membershipRepo, tierRepo, pub)
	return NewPointsUseCase(membershipRepo, txRepo, membershipUC, pub)
}

// ===========================================================================
// EarnPoints tests
// ===========================================================================

func TestEarnPoints_Success(t *testing.T) {
	mRepo, txRepo, tRepo, pub := defaultLoyaltyMocks()

	mRepo.getByUserIDFn = func(_ context.Context, _ string) (*domain.Membership, error) {
		return &domain.Membership{
			UserID:         "user-1",
			Tier:           domain.TierBronze,
			PointsBalance:  100,
			LifetimePoints: 100,
		}, nil
	}

	var savedTx *domain.PointsTransaction
	txRepo.createFn = func(_ context.Context, tx *domain.PointsTransaction) error {
		savedTx = tx
		return nil
	}

	var updatedBalance, updatedLifetime int64
	mRepo.updatePointsFn = func(_ context.Context, _ string, balance, lifetime int64) error {
		updatedBalance = balance
		updatedLifetime = lifetime
		return nil
	}

	// Tier check returns nil (no upgrade needed)
	tRepo.getTierForPointsFn = func(_ context.Context, _ int64) (*domain.Tier, error) {
		return nil, nil
	}

	uc := newPointsUseCase(mRepo, txRepo, tRepo, pub)
	tx, err := uc.EarnPoints(context.Background(), EarnPointsRequest{
		UserID:      "user-1",
		Points:      50,
		Source:      domain.SourceOrder,
		ReferenceID: "order-123",
		Description: "Order purchase",
	})

	require.NoError(t, err)
	require.NotNil(t, tx)
	assert.NotEmpty(t, tx.ID)
	assert.Equal(t, "user-1", tx.UserID)
	assert.Equal(t, domain.TransactionEarn, tx.Type)
	assert.Equal(t, int64(50), tx.Points)
	assert.Equal(t, domain.SourceOrder, tx.Source)
	assert.Equal(t, "order-123", tx.ReferenceID)
	assert.NotNil(t, savedTx)
	assert.Equal(t, int64(150), updatedBalance)
	assert.Equal(t, int64(150), updatedLifetime)
}

func TestEarnPoints_MembershipNotFound(t *testing.T) {
	mRepo, txRepo, tRepo, pub := defaultLoyaltyMocks()

	mRepo.getByUserIDFn = func(_ context.Context, _ string) (*domain.Membership, error) {
		return nil, errors.New("not found")
	}

	uc := newPointsUseCase(mRepo, txRepo, tRepo, pub)
	_, err := uc.EarnPoints(context.Background(), EarnPointsRequest{
		UserID: "user-missing",
		Points: 50,
		Source: domain.SourceOrder,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get membership")
}

func TestEarnPoints_TransactionCreateError(t *testing.T) {
	mRepo, txRepo, tRepo, pub := defaultLoyaltyMocks()

	mRepo.getByUserIDFn = func(_ context.Context, _ string) (*domain.Membership, error) {
		return &domain.Membership{UserID: "user-1", PointsBalance: 100, LifetimePoints: 100}, nil
	}
	txRepo.createFn = func(_ context.Context, _ *domain.PointsTransaction) error {
		return errors.New("db error")
	}

	uc := newPointsUseCase(mRepo, txRepo, tRepo, pub)
	_, err := uc.EarnPoints(context.Background(), EarnPointsRequest{
		UserID: "user-1",
		Points: 50,
		Source: domain.SourceOrder,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create transaction")
}

func TestEarnPoints_UpdatePointsError(t *testing.T) {
	mRepo, txRepo, tRepo, pub := defaultLoyaltyMocks()

	mRepo.getByUserIDFn = func(_ context.Context, _ string) (*domain.Membership, error) {
		return &domain.Membership{UserID: "user-1", PointsBalance: 100, LifetimePoints: 100}, nil
	}
	mRepo.updatePointsFn = func(_ context.Context, _ string, _, _ int64) error {
		return errors.New("db error")
	}

	uc := newPointsUseCase(mRepo, txRepo, tRepo, pub)
	_, err := uc.EarnPoints(context.Background(), EarnPointsRequest{
		UserID: "user-1",
		Points: 50,
		Source: domain.SourceOrder,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update points")
}

// ===========================================================================
// RedeemPoints tests
// ===========================================================================

func TestRedeemPoints_Success(t *testing.T) {
	mRepo, txRepo, tRepo, pub := defaultLoyaltyMocks()

	mRepo.getByUserIDFn = func(_ context.Context, _ string) (*domain.Membership, error) {
		return &domain.Membership{
			UserID:         "user-1",
			PointsBalance:  500,
			LifetimePoints: 1000,
		}, nil
	}

	var savedTx *domain.PointsTransaction
	txRepo.createFn = func(_ context.Context, tx *domain.PointsTransaction) error {
		savedTx = tx
		return nil
	}

	var updatedBalance int64
	mRepo.updatePointsFn = func(_ context.Context, _ string, balance, _ int64) error {
		updatedBalance = balance
		return nil
	}

	uc := newPointsUseCase(mRepo, txRepo, tRepo, pub)
	tx, err := uc.RedeemPoints(context.Background(), RedeemPointsRequest{
		UserID:      "user-1",
		Points:      200,
		OrderID:     "order-456",
		Description: "Order discount",
	})

	require.NoError(t, err)
	require.NotNil(t, tx)
	assert.Equal(t, domain.TransactionRedeem, tx.Type)
	assert.Equal(t, int64(200), tx.Points)
	assert.Equal(t, domain.SourceOrder, tx.Source)
	assert.Equal(t, "order-456", tx.ReferenceID)
	assert.NotNil(t, savedTx)
	assert.Equal(t, int64(300), updatedBalance)
}

func TestRedeemPoints_InsufficientBalance(t *testing.T) {
	mRepo, txRepo, tRepo, pub := defaultLoyaltyMocks()

	mRepo.getByUserIDFn = func(_ context.Context, _ string) (*domain.Membership, error) {
		return &domain.Membership{
			UserID:        "user-1",
			PointsBalance: 100,
		}, nil
	}

	uc := newPointsUseCase(mRepo, txRepo, tRepo, pub)
	_, err := uc.RedeemPoints(context.Background(), RedeemPointsRequest{
		UserID:  "user-1",
		Points:  200,
		OrderID: "order-456",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient points")
}

func TestRedeemPoints_MembershipNotFound(t *testing.T) {
	mRepo, txRepo, tRepo, pub := defaultLoyaltyMocks()

	mRepo.getByUserIDFn = func(_ context.Context, _ string) (*domain.Membership, error) {
		return nil, errors.New("not found")
	}

	uc := newPointsUseCase(mRepo, txRepo, tRepo, pub)
	_, err := uc.RedeemPoints(context.Background(), RedeemPointsRequest{
		UserID:  "user-missing",
		Points:  50,
		OrderID: "order-1",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get membership")
}

func TestRedeemPoints_TransactionCreateError(t *testing.T) {
	mRepo, txRepo, tRepo, pub := defaultLoyaltyMocks()

	mRepo.getByUserIDFn = func(_ context.Context, _ string) (*domain.Membership, error) {
		return &domain.Membership{UserID: "user-1", PointsBalance: 500}, nil
	}
	txRepo.createFn = func(_ context.Context, _ *domain.PointsTransaction) error {
		return errors.New("db error")
	}

	uc := newPointsUseCase(mRepo, txRepo, tRepo, pub)
	_, err := uc.RedeemPoints(context.Background(), RedeemPointsRequest{
		UserID:  "user-1",
		Points:  100,
		OrderID: "order-1",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create transaction")
}

// ===========================================================================
// GetBalance tests
// ===========================================================================

func TestGetBalance_Success(t *testing.T) {
	mRepo, txRepo, tRepo, pub := defaultLoyaltyMocks()

	mRepo.getByUserIDFn = func(_ context.Context, _ string) (*domain.Membership, error) {
		return &domain.Membership{UserID: "user-1", PointsBalance: 750}, nil
	}

	uc := newPointsUseCase(mRepo, txRepo, tRepo, pub)
	balance, err := uc.GetBalance(context.Background(), "user-1")

	require.NoError(t, err)
	assert.Equal(t, int64(750), balance)
}

func TestGetBalance_NotFound(t *testing.T) {
	mRepo, txRepo, tRepo, pub := defaultLoyaltyMocks()

	mRepo.getByUserIDFn = func(_ context.Context, _ string) (*domain.Membership, error) {
		return nil, errors.New("not found")
	}

	uc := newPointsUseCase(mRepo, txRepo, tRepo, pub)
	_, err := uc.GetBalance(context.Background(), "user-missing")

	require.Error(t, err)
}

// ===========================================================================
// ListTransactions tests
// ===========================================================================

func TestListTransactions_Success(t *testing.T) {
	mRepo, txRepo, tRepo, pub := defaultLoyaltyMocks()

	txRepo.listByUserFn = func(_ context.Context, userID string, page, pageSize int) ([]domain.PointsTransaction, int64, error) {
		return []domain.PointsTransaction{
			{ID: "tx-1", UserID: userID, Type: domain.TransactionEarn, Points: 100},
			{ID: "tx-2", UserID: userID, Type: domain.TransactionRedeem, Points: 50},
		}, 2, nil
	}

	uc := newPointsUseCase(mRepo, txRepo, tRepo, pub)
	txs, total, err := uc.ListTransactions(context.Background(), "user-1", 1, 20)

	require.NoError(t, err)
	assert.Len(t, txs, 2)
	assert.Equal(t, int64(2), total)
}

func TestListTransactions_DefaultsPagination(t *testing.T) {
	mRepo, txRepo, tRepo, pub := defaultLoyaltyMocks()

	var capturedPage, capturedPageSize int
	txRepo.listByUserFn = func(_ context.Context, _ string, page, pageSize int) ([]domain.PointsTransaction, int64, error) {
		capturedPage = page
		capturedPageSize = pageSize
		return nil, 0, nil
	}

	uc := newPointsUseCase(mRepo, txRepo, tRepo, pub)

	// Test page < 1 defaults to 1
	_, _, _ = uc.ListTransactions(context.Background(), "user-1", 0, 20)
	assert.Equal(t, 1, capturedPage)

	// Test pageSize < 1 defaults to 20
	_, _, _ = uc.ListTransactions(context.Background(), "user-1", 1, 0)
	assert.Equal(t, 20, capturedPageSize)

	// Test pageSize > 100 defaults to 20
	_, _, _ = uc.ListTransactions(context.Background(), "user-1", 1, 200)
	assert.Equal(t, 20, capturedPageSize)
}
