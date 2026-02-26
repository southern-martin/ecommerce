package domain

import "context"

// CouponRepository defines the interface for coupon persistence.
type CouponRepository interface {
	GetByID(ctx context.Context, id string) (*Coupon, error)
	GetByCode(ctx context.Context, code string) (*Coupon, error)
	ListAll(ctx context.Context, page, pageSize int) ([]*Coupon, int64, error)
	ListBySeller(ctx context.Context, sellerID string, page, pageSize int) ([]*Coupon, int64, error)
	Create(ctx context.Context, coupon *Coupon) error
	Update(ctx context.Context, coupon *Coupon) error
	IncrementUsageCount(ctx context.Context, id string) error
}

// CouponUsageRepository defines the interface for coupon usage persistence.
type CouponUsageRepository interface {
	GetByUserAndCoupon(ctx context.Context, userID, couponID string) ([]*CouponUsage, error)
	CountByUser(ctx context.Context, userID, couponID string) (int64, error)
	Create(ctx context.Context, usage *CouponUsage) error
}

// FlashSaleRepository defines the interface for flash sale persistence.
type FlashSaleRepository interface {
	GetByID(ctx context.Context, id string) (*FlashSale, error)
	ListActive(ctx context.Context) ([]*FlashSale, error)
	ListAll(ctx context.Context, page, pageSize int) ([]*FlashSale, int64, error)
	Create(ctx context.Context, flashSale *FlashSale) error
	Update(ctx context.Context, flashSale *FlashSale) error
}

// FlashSaleItemRepository defines the interface for flash sale item persistence.
type FlashSaleItemRepository interface {
	GetByFlashSaleID(ctx context.Context, flashSaleID string) ([]*FlashSaleItem, error)
	Create(ctx context.Context, item *FlashSaleItem) error
	IncrementSoldCount(ctx context.Context, id string) error
}

// BundleRepository defines the interface for bundle persistence.
type BundleRepository interface {
	GetByID(ctx context.Context, id string) (*Bundle, error)
	ListBySeller(ctx context.Context, sellerID string, page, pageSize int) ([]*Bundle, int64, error)
	ListActive(ctx context.Context, page, pageSize int) ([]*Bundle, int64, error)
	Create(ctx context.Context, bundle *Bundle) error
	Update(ctx context.Context, bundle *Bundle) error
}
