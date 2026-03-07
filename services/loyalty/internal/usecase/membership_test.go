package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/loyalty/internal/domain"
)

// ===========================================================================
// CreateMembership tests
// ===========================================================================

func TestCreateMembership_Success(t *testing.T) {
	mRepo, _, tRepo, pub := defaultLoyaltyMocks()

	var saved *domain.Membership
	mRepo.createFn = func(_ context.Context, m *domain.Membership) error {
		saved = m
		return nil
	}

	uc := NewMembershipUseCase(mRepo, tRepo, pub)
	membership, err := uc.CreateMembership(context.Background(), "user-1")

	require.NoError(t, err)
	require.NotNil(t, membership)
	assert.Equal(t, "user-1", membership.UserID)
	assert.Equal(t, domain.TierBronze, membership.Tier)
	assert.Equal(t, int64(0), membership.PointsBalance)
	assert.Equal(t, int64(0), membership.LifetimePoints)
	assert.False(t, membership.JoinedAt.IsZero())
	assert.NotNil(t, saved)
}

func TestCreateMembership_RepoError(t *testing.T) {
	mRepo, _, tRepo, pub := defaultLoyaltyMocks()

	mRepo.createFn = func(_ context.Context, _ *domain.Membership) error {
		return errors.New("duplicate key")
	}

	uc := NewMembershipUseCase(mRepo, tRepo, pub)
	_, err := uc.CreateMembership(context.Background(), "user-1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create membership")
}

// ===========================================================================
// GetMembership tests
// ===========================================================================

func TestGetMembership_Success(t *testing.T) {
	mRepo, _, tRepo, pub := defaultLoyaltyMocks()

	mRepo.getByUserIDFn = func(_ context.Context, userID string) (*domain.Membership, error) {
		return &domain.Membership{
			UserID:        userID,
			Tier:          domain.TierSilver,
			PointsBalance: 500,
		}, nil
	}

	uc := NewMembershipUseCase(mRepo, tRepo, pub)
	membership, err := uc.GetMembership(context.Background(), "user-1")

	require.NoError(t, err)
	require.NotNil(t, membership)
	assert.Equal(t, "user-1", membership.UserID)
	assert.Equal(t, domain.TierSilver, membership.Tier)
	assert.Equal(t, int64(500), membership.PointsBalance)
}

func TestGetMembership_NotFound(t *testing.T) {
	mRepo, _, tRepo, pub := defaultLoyaltyMocks()

	mRepo.getByUserIDFn = func(_ context.Context, _ string) (*domain.Membership, error) {
		return nil, errors.New("not found")
	}

	uc := NewMembershipUseCase(mRepo, tRepo, pub)
	_, err := uc.GetMembership(context.Background(), "user-missing")

	require.Error(t, err)
}

// ===========================================================================
// CheckAndUpgradeTier tests
// ===========================================================================

func TestCheckAndUpgradeTier_UpgradeOccurs(t *testing.T) {
	mRepo, _, tRepo, pub := defaultLoyaltyMocks()

	mRepo.getByUserIDFn = func(_ context.Context, _ string) (*domain.Membership, error) {
		return &domain.Membership{
			UserID:         "user-1",
			Tier:           domain.TierBronze,
			LifetimePoints: 5000,
		}, nil
	}

	tRepo.getTierForPointsFn = func(_ context.Context, _ int64) (*domain.Tier, error) {
		return &domain.Tier{Name: "silver", MinPoints: 1000}, nil
	}

	var updatedTier domain.MemberTier
	mRepo.updateTierFn = func(_ context.Context, _ string, tier domain.MemberTier) error {
		updatedTier = tier
		return nil
	}

	uc := NewMembershipUseCase(mRepo, tRepo, pub)
	err := uc.CheckAndUpgradeTier(context.Background(), "user-1")

	require.NoError(t, err)
	assert.Equal(t, domain.MemberTier("silver"), updatedTier)
}

func TestCheckAndUpgradeTier_NoUpgradeNeeded(t *testing.T) {
	mRepo, _, tRepo, pub := defaultLoyaltyMocks()

	mRepo.getByUserIDFn = func(_ context.Context, _ string) (*domain.Membership, error) {
		return &domain.Membership{
			UserID:         "user-1",
			Tier:           domain.TierSilver,
			LifetimePoints: 2000,
		}, nil
	}

	// Return the same tier the user already has
	tRepo.getTierForPointsFn = func(_ context.Context, _ int64) (*domain.Tier, error) {
		return &domain.Tier{Name: "silver", MinPoints: 1000}, nil
	}

	updateTierCalled := false
	mRepo.updateTierFn = func(_ context.Context, _ string, _ domain.MemberTier) error {
		updateTierCalled = true
		return nil
	}

	uc := NewMembershipUseCase(mRepo, tRepo, pub)
	err := uc.CheckAndUpgradeTier(context.Background(), "user-1")

	require.NoError(t, err)
	assert.False(t, updateTierCalled, "should not update tier when already at correct tier")
}

func TestCheckAndUpgradeTier_MembershipNotFound(t *testing.T) {
	mRepo, _, tRepo, pub := defaultLoyaltyMocks()

	mRepo.getByUserIDFn = func(_ context.Context, _ string) (*domain.Membership, error) {
		return nil, errors.New("not found")
	}

	uc := NewMembershipUseCase(mRepo, tRepo, pub)
	err := uc.CheckAndUpgradeTier(context.Background(), "user-missing")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get membership")
}

func TestCheckAndUpgradeTier_TierLookupError(t *testing.T) {
	mRepo, _, tRepo, pub := defaultLoyaltyMocks()

	mRepo.getByUserIDFn = func(_ context.Context, _ string) (*domain.Membership, error) {
		return &domain.Membership{UserID: "user-1", Tier: domain.TierBronze, LifetimePoints: 500}, nil
	}
	tRepo.getTierForPointsFn = func(_ context.Context, _ int64) (*domain.Tier, error) {
		return nil, errors.New("db error")
	}

	uc := NewMembershipUseCase(mRepo, tRepo, pub)
	err := uc.CheckAndUpgradeTier(context.Background(), "user-1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get tier for points")
}

func TestCheckAndUpgradeTier_NilTierReturned(t *testing.T) {
	mRepo, _, tRepo, pub := defaultLoyaltyMocks()

	mRepo.getByUserIDFn = func(_ context.Context, _ string) (*domain.Membership, error) {
		return &domain.Membership{UserID: "user-1", Tier: domain.TierBronze, LifetimePoints: 10}, nil
	}
	tRepo.getTierForPointsFn = func(_ context.Context, _ int64) (*domain.Tier, error) {
		return nil, nil
	}

	updateTierCalled := false
	mRepo.updateTierFn = func(_ context.Context, _ string, _ domain.MemberTier) error {
		updateTierCalled = true
		return nil
	}

	uc := NewMembershipUseCase(mRepo, tRepo, pub)
	err := uc.CheckAndUpgradeTier(context.Background(), "user-1")

	require.NoError(t, err)
	assert.False(t, updateTierCalled, "should not update tier when GetTierForPoints returns nil")
}
