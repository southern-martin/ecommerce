package domain

import "context"

// UserProfileRepository defines the interface for user profile persistence.
type UserProfileRepository interface {
	Create(ctx context.Context, profile *UserProfile) error
	GetByID(ctx context.Context, id string) (*UserProfile, error)
	Update(ctx context.Context, profile *UserProfile) error
}

// AddressRepository defines the interface for address persistence.
type AddressRepository interface {
	Create(ctx context.Context, addr *Address) error
	GetByID(ctx context.Context, id string) (*Address, error)
	ListByUserID(ctx context.Context, userID string) ([]Address, error)
	Update(ctx context.Context, addr *Address) error
	Delete(ctx context.Context, id string) error
	CountByUserID(ctx context.Context, userID string) (int64, error)
	ClearDefaultByUserID(ctx context.Context, userID string) error
}

// SellerProfileRepository defines the interface for seller profile persistence.
type SellerProfileRepository interface {
	Create(ctx context.Context, seller *SellerProfile) error
	GetByID(ctx context.Context, id string) (*SellerProfile, error)
	GetByUserID(ctx context.Context, userID string) (*SellerProfile, error)
	Update(ctx context.Context, seller *SellerProfile) error
	List(ctx context.Context, page, size int) ([]SellerProfile, int64, error)
}

// FollowRepository defines the interface for user follow persistence.
type FollowRepository interface {
	Create(ctx context.Context, follow *UserFollow) error
	Delete(ctx context.Context, followerID, sellerID string) error
	ListByFollowerID(ctx context.Context, followerID string, page, size int) ([]SellerProfile, int64, error)
	CountBySellerID(ctx context.Context, sellerID string) (int64, error)
	Exists(ctx context.Context, followerID, sellerID string) (bool, error)
}
