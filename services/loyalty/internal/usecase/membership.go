package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/southern-martin/ecommerce/services/loyalty/internal/domain"
)

// MembershipUseCase handles membership operations.
type MembershipUseCase struct {
	membershipRepo domain.MembershipRepository
	tierRepo       domain.TierRepository
	publisher      domain.EventPublisher
}

// NewMembershipUseCase creates a new MembershipUseCase.
func NewMembershipUseCase(membershipRepo domain.MembershipRepository, tierRepo domain.TierRepository, publisher domain.EventPublisher) *MembershipUseCase {
	return &MembershipUseCase{
		membershipRepo: membershipRepo,
		tierRepo:       tierRepo,
		publisher:      publisher,
	}
}

// GetMembership retrieves a user's loyalty membership.
func (uc *MembershipUseCase) GetMembership(ctx context.Context, userID string) (*domain.Membership, error) {
	return uc.membershipRepo.GetByUserID(ctx, userID)
}

// CreateMembership creates a new loyalty membership with default bronze tier and 0 points.
func (uc *MembershipUseCase) CreateMembership(ctx context.Context, userID string) (*domain.Membership, error) {
	membership := &domain.Membership{
		UserID:         userID,
		Tier:           domain.TierBronze,
		PointsBalance:  0,
		LifetimePoints: 0,
		JoinedAt:       time.Now(),
	}

	if err := uc.membershipRepo.Create(ctx, membership); err != nil {
		return nil, fmt.Errorf("failed to create membership: %w", err)
	}

	return membership, nil
}

// CheckAndUpgradeTier compares lifetime points to tier thresholds and upgrades if eligible.
func (uc *MembershipUseCase) CheckAndUpgradeTier(ctx context.Context, userID string) error {
	membership, err := uc.membershipRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get membership: %w", err)
	}

	newTier, err := uc.tierRepo.GetTierForPoints(ctx, membership.LifetimePoints)
	if err != nil {
		return fmt.Errorf("failed to get tier for points: %w", err)
	}

	if newTier != nil && domain.MemberTier(newTier.Name) != membership.Tier {
		oldTier := membership.Tier
		if err := uc.membershipRepo.UpdateTier(ctx, userID, domain.MemberTier(newTier.Name)); err != nil {
			return fmt.Errorf("failed to update tier: %w", err)
		}

		_ = uc.publisher.Publish(ctx, "loyalty.tier.upgraded", map[string]interface{}{
			"user_id":  userID,
			"old_tier": string(oldTier),
			"new_tier": newTier.Name,
		})
	}

	return nil
}
