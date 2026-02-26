package domain

import "context"

// MembershipRepository defines the interface for membership persistence.
type MembershipRepository interface {
	GetByUserID(ctx context.Context, userID string) (*Membership, error)
	Create(ctx context.Context, membership *Membership) error
	Update(ctx context.Context, membership *Membership) error
	UpdateTier(ctx context.Context, userID string, tier MemberTier) error
	UpdatePoints(ctx context.Context, userID string, pointsBalance, lifetimePoints int64) error
}

// PointsTransactionRepository defines the interface for points transaction persistence.
type PointsTransactionRepository interface {
	GetByID(ctx context.Context, id string) (*PointsTransaction, error)
	ListByUser(ctx context.Context, userID string, page, pageSize int) ([]PointsTransaction, int64, error)
	Create(ctx context.Context, tx *PointsTransaction) error
}

// TierRepository defines the interface for tier persistence.
type TierRepository interface {
	GetAll(ctx context.Context) ([]Tier, error)
	GetByName(ctx context.Context, name string) (*Tier, error)
	GetTierForPoints(ctx context.Context, lifetimePoints int64) (*Tier, error)
}
