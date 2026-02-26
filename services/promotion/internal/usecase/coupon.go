package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/southern-martin/ecommerce/services/promotion/internal/domain"
)

// CreateCouponInput represents the input for creating a coupon.
type CreateCouponInput struct {
	Code             string
	Type             string
	DiscountValue    int64
	MinOrderCents    int64
	MaxDiscountCents int64
	UsageLimit       int
	PerUserLimit     int
	Scope            string
	ScopeIDs         []string
	CreatedBy        string
	StartsAt         time.Time
	ExpiresAt        time.Time
}

// ValidateCouponInput represents the input for validating a coupon.
type ValidateCouponInput struct {
	Code        string
	UserID      string
	OrderCents  int64
}

// CouponUseCase handles coupon business logic.
type CouponUseCase struct {
	couponRepo      domain.CouponRepository
	couponUsageRepo domain.CouponUsageRepository
	publisher       domain.EventPublisher
}

// NewCouponUseCase creates a new CouponUseCase instance.
func NewCouponUseCase(
	couponRepo domain.CouponRepository,
	couponUsageRepo domain.CouponUsageRepository,
	publisher domain.EventPublisher,
) *CouponUseCase {
	return &CouponUseCase{
		couponRepo:      couponRepo,
		couponUsageRepo: couponUsageRepo,
		publisher:       publisher,
	}
}

// CreateCoupon creates a new coupon.
func (uc *CouponUseCase) CreateCoupon(ctx context.Context, input CreateCouponInput) (*domain.Coupon, error) {
	if input.Code == "" {
		return nil, errors.New("coupon code is required")
	}
	if input.Type == "" {
		return nil, errors.New("coupon type is required")
	}
	if input.DiscountValue <= 0 {
		return nil, errors.New("discount value must be greater than 0")
	}
	if input.CreatedBy == "" {
		return nil, errors.New("created_by is required")
	}
	if input.StartsAt.IsZero() {
		input.StartsAt = time.Now()
	}
	if input.ExpiresAt.IsZero() {
		return nil, errors.New("expires_at is required")
	}

	coupon := domain.NewCoupon(
		strings.ToUpper(input.Code),
		domain.CouponType(input.Type),
		input.DiscountValue,
		input.CreatedBy,
	)
	coupon.MinOrderCents = input.MinOrderCents
	coupon.MaxDiscountCents = input.MaxDiscountCents
	coupon.UsageLimit = input.UsageLimit
	coupon.PerUserLimit = input.PerUserLimit
	coupon.StartsAt = input.StartsAt
	coupon.ExpiresAt = input.ExpiresAt

	if input.Scope != "" {
		coupon.Scope = domain.CouponScope(input.Scope)
	}
	if len(input.ScopeIDs) > 0 {
		coupon.ScopeIDs = input.ScopeIDs
	}

	if err := uc.couponRepo.Create(ctx, coupon); err != nil {
		return nil, err
	}

	return coupon, nil
}

// GetCoupon retrieves a coupon by ID.
func (uc *CouponUseCase) GetCoupon(ctx context.Context, id string) (*domain.Coupon, error) {
	if id == "" {
		return nil, errors.New("coupon id is required")
	}
	return uc.couponRepo.GetByID(ctx, id)
}

// ListCoupons retrieves a paginated list of coupons.
func (uc *CouponUseCase) ListCoupons(ctx context.Context, page, pageSize int) ([]*domain.Coupon, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return uc.couponRepo.ListAll(ctx, page, pageSize)
}

// ListCouponsBySeller retrieves a paginated list of coupons created by a seller.
func (uc *CouponUseCase) ListCouponsBySeller(ctx context.Context, sellerID string, page, pageSize int) ([]*domain.Coupon, int64, error) {
	if sellerID == "" {
		return nil, 0, errors.New("seller id is required")
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return uc.couponRepo.ListBySeller(ctx, sellerID, page, pageSize)
}

// ValidateCoupon checks if a coupon is valid for use.
func (uc *CouponUseCase) ValidateCoupon(ctx context.Context, input ValidateCouponInput) (*domain.Coupon, int64, error) {
	if input.Code == "" {
		return nil, 0, errors.New("coupon code is required")
	}
	if input.UserID == "" {
		return nil, 0, errors.New("user_id is required")
	}

	coupon, err := uc.couponRepo.GetByCode(ctx, strings.ToUpper(input.Code))
	if err != nil {
		return nil, 0, err
	}

	// Check if coupon is active
	if !coupon.IsActive {
		return nil, 0, errors.New("coupon is not active")
	}

	// Check if coupon has started
	now := time.Now()
	if now.Before(coupon.StartsAt) {
		return nil, 0, errors.New("coupon is not yet active")
	}

	// Check if coupon has expired
	if now.After(coupon.ExpiresAt) {
		return nil, 0, errors.New("coupon has expired")
	}

	// Check usage limit
	if coupon.UsageLimit > 0 && coupon.UsageCount >= coupon.UsageLimit {
		return nil, 0, errors.New("coupon usage limit reached")
	}

	// Check per-user limit
	if coupon.PerUserLimit > 0 {
		count, err := uc.couponUsageRepo.CountByUser(ctx, input.UserID, coupon.ID)
		if err != nil {
			return nil, 0, err
		}
		if count >= int64(coupon.PerUserLimit) {
			return nil, 0, errors.New("per-user coupon usage limit reached")
		}
	}

	// Check minimum order amount
	if coupon.MinOrderCents > 0 && input.OrderCents < coupon.MinOrderCents {
		return nil, 0, errors.New("order total does not meet minimum requirement")
	}

	// Calculate discount
	discountCents := uc.calculateDiscount(coupon, input.OrderCents)

	return coupon, discountCents, nil
}

// RedeemCoupon redeems a coupon for a user and order.
func (uc *CouponUseCase) RedeemCoupon(ctx context.Context, code, userID, orderID string, orderCents int64) (*domain.CouponUsage, error) {
	// Validate first
	coupon, discountCents, err := uc.ValidateCoupon(ctx, ValidateCouponInput{
		Code:       code,
		UserID:     userID,
		OrderCents: orderCents,
	})
	if err != nil {
		return nil, err
	}

	// Increment usage count
	if err := uc.couponRepo.IncrementUsageCount(ctx, coupon.ID); err != nil {
		return nil, err
	}

	// Create usage record
	usage := domain.NewCouponUsage(coupon.ID, userID, orderID, discountCents)
	if err := uc.couponUsageRepo.Create(ctx, usage); err != nil {
		return nil, err
	}

	// Publish coupon.redeemed event
	event := domain.CouponRedeemedEvent{
		CouponID:      coupon.ID,
		CouponCode:    coupon.Code,
		UserID:        userID,
		OrderID:       orderID,
		DiscountCents: discountCents,
	}
	_ = uc.publisher.Publish(ctx, domain.EventCouponRedeemed, event)

	return usage, nil
}

// UpdateCoupon updates an existing coupon.
func (uc *CouponUseCase) UpdateCoupon(ctx context.Context, coupon *domain.Coupon) error {
	return uc.couponRepo.Update(ctx, coupon)
}

// calculateDiscount computes the discount amount based on coupon type.
func (uc *CouponUseCase) calculateDiscount(coupon *domain.Coupon, orderCents int64) int64 {
	var discount int64

	switch coupon.Type {
	case domain.CouponTypePercentage:
		// DiscountValue is percentage * 100 (e.g., 1000 = 10%)
		discount = orderCents * coupon.DiscountValue / 10000
	case domain.CouponTypeFixedAmount:
		discount = coupon.DiscountValue
	case domain.CouponTypeFreeShipping:
		// Free shipping is handled externally; discount is 0 here
		discount = 0
	}

	// Cap at max discount
	if coupon.MaxDiscountCents > 0 && discount > coupon.MaxDiscountCents {
		discount = coupon.MaxDiscountCents
	}

	// Don't exceed order total
	if discount > orderCents {
		discount = orderCents
	}

	return discount
}
