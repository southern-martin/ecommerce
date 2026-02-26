package domain

import (
	"time"

	"github.com/google/uuid"
)

// CouponType represents the type of discount a coupon provides.
type CouponType string

const (
	CouponTypePercentage   CouponType = "percentage"
	CouponTypeFixedAmount  CouponType = "fixed_amount"
	CouponTypeFreeShipping CouponType = "free_shipping"
)

// CouponScope represents the scope of a coupon.
type CouponScope string

const (
	CouponScopeAll      CouponScope = "all"
	CouponScopeCategory CouponScope = "category"
	CouponScopeProduct  CouponScope = "product"
	CouponScopeSeller   CouponScope = "seller"
)

// Coupon represents a discount coupon.
type Coupon struct {
	ID              string
	Code            string
	Type            CouponType
	DiscountValue   int64 // cents or percentage*100
	MinOrderCents   int64
	MaxDiscountCents int64
	UsageLimit      int
	UsageCount      int
	PerUserLimit    int
	Scope           CouponScope
	ScopeIDs        []string
	CreatedBy       string // seller_id or "platform"
	StartsAt        time.Time
	ExpiresAt       time.Time
	IsActive        bool
	CreatedAt       time.Time
}

// CouponUsage records when a coupon is redeemed.
type CouponUsage struct {
	ID            string
	CouponID      string
	UserID        string
	OrderID       string
	DiscountCents int64
	CreatedAt     time.Time
}

// FlashSale represents a time-limited sale event.
type FlashSale struct {
	ID        string
	Name      string
	StartsAt  time.Time
	EndsAt    time.Time
	IsActive  bool
	Items     []FlashSaleItem
	CreatedAt time.Time
}

// FlashSaleItem represents a product in a flash sale.
type FlashSaleItem struct {
	ID             string
	FlashSaleID    string
	ProductID      string
	VariantID      string
	SalePriceCents int64
	QuantityLimit  int
	SoldCount      int
}

// Bundle represents a product bundle offering.
type Bundle struct {
	ID              string
	Name            string
	SellerID        string
	ProductIDs      []string
	BundlePriceCents int64
	SavingsCents    int64
	IsActive        bool
	CreatedAt       time.Time
}

// NewCoupon creates a new Coupon with a generated ID.
func NewCoupon(code string, couponType CouponType, discountValue int64, createdBy string) *Coupon {
	now := time.Now()
	return &Coupon{
		ID:            uuid.New().String(),
		Code:          code,
		Type:          couponType,
		DiscountValue: discountValue,
		Scope:         CouponScopeAll,
		CreatedBy:     createdBy,
		IsActive:      true,
		CreatedAt:     now,
	}
}

// NewCouponUsage creates a new CouponUsage record.
func NewCouponUsage(couponID, userID, orderID string, discountCents int64) *CouponUsage {
	return &CouponUsage{
		ID:            uuid.New().String(),
		CouponID:      couponID,
		UserID:        userID,
		OrderID:       orderID,
		DiscountCents: discountCents,
		CreatedAt:     time.Now(),
	}
}

// NewFlashSale creates a new FlashSale with a generated ID.
func NewFlashSale(name string, startsAt, endsAt time.Time) *FlashSale {
	return &FlashSale{
		ID:        uuid.New().String(),
		Name:      name,
		StartsAt:  startsAt,
		EndsAt:    endsAt,
		IsActive:  true,
		CreatedAt: time.Now(),
	}
}

// NewFlashSaleItem creates a new FlashSaleItem with a generated ID.
func NewFlashSaleItem(flashSaleID, productID, variantID string, salePriceCents int64, quantityLimit int) *FlashSaleItem {
	return &FlashSaleItem{
		ID:             uuid.New().String(),
		FlashSaleID:    flashSaleID,
		ProductID:      productID,
		VariantID:      variantID,
		SalePriceCents: salePriceCents,
		QuantityLimit:  quantityLimit,
		SoldCount:      0,
	}
}

// NewBundle creates a new Bundle with a generated ID.
func NewBundle(name, sellerID string, productIDs []string, bundlePriceCents, savingsCents int64) *Bundle {
	return &Bundle{
		ID:               uuid.New().String(),
		Name:             name,
		SellerID:         sellerID,
		ProductIDs:       productIDs,
		BundlePriceCents: bundlePriceCents,
		SavingsCents:     savingsCents,
		IsActive:         true,
		CreatedAt:        time.Now(),
	}
}
