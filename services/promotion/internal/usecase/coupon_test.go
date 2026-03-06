package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/southern-martin/ecommerce/services/promotion/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Hand-written function-field mocks
// ---------------------------------------------------------------------------

type mockCouponRepo struct {
	getByIDFn        func(ctx context.Context, id string) (*domain.Coupon, error)
	getByCodeFn      func(ctx context.Context, code string) (*domain.Coupon, error)
	listAllFn        func(ctx context.Context, page, pageSize int) ([]*domain.Coupon, int64, error)
	listBySellerFn   func(ctx context.Context, sellerID string, page, pageSize int) ([]*domain.Coupon, int64, error)
	createFn         func(ctx context.Context, coupon *domain.Coupon) error
	updateFn         func(ctx context.Context, coupon *domain.Coupon) error
	incrementUsageFn func(ctx context.Context, id string) error
}

func (m *mockCouponRepo) GetByID(ctx context.Context, id string) (*domain.Coupon, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockCouponRepo) GetByCode(ctx context.Context, code string) (*domain.Coupon, error) {
	return m.getByCodeFn(ctx, code)
}
func (m *mockCouponRepo) ListAll(ctx context.Context, page, pageSize int) ([]*domain.Coupon, int64, error) {
	return m.listAllFn(ctx, page, pageSize)
}
func (m *mockCouponRepo) ListBySeller(ctx context.Context, sellerID string, page, pageSize int) ([]*domain.Coupon, int64, error) {
	return m.listBySellerFn(ctx, sellerID, page, pageSize)
}
func (m *mockCouponRepo) Create(ctx context.Context, coupon *domain.Coupon) error {
	return m.createFn(ctx, coupon)
}
func (m *mockCouponRepo) Update(ctx context.Context, coupon *domain.Coupon) error {
	return m.updateFn(ctx, coupon)
}
func (m *mockCouponRepo) IncrementUsageCount(ctx context.Context, id string) error {
	return m.incrementUsageFn(ctx, id)
}

type mockCouponUsageRepo struct {
	getByUserAndCouponFn func(ctx context.Context, userID, couponID string) ([]*domain.CouponUsage, error)
	countByUserFn        func(ctx context.Context, userID, couponID string) (int64, error)
	createFn             func(ctx context.Context, usage *domain.CouponUsage) error
}

func (m *mockCouponUsageRepo) GetByUserAndCoupon(ctx context.Context, userID, couponID string) ([]*domain.CouponUsage, error) {
	return m.getByUserAndCouponFn(ctx, userID, couponID)
}
func (m *mockCouponUsageRepo) CountByUser(ctx context.Context, userID, couponID string) (int64, error) {
	return m.countByUserFn(ctx, userID, couponID)
}
func (m *mockCouponUsageRepo) Create(ctx context.Context, usage *domain.CouponUsage) error {
	return m.createFn(ctx, usage)
}

type mockEventPublisher struct {
	publishFn func(ctx context.Context, subject string, data interface{}) error
}

func (m *mockEventPublisher) Publish(ctx context.Context, subject string, data interface{}) error {
	return m.publishFn(ctx, subject, data)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func newTestCouponUseCase(
	cr *mockCouponRepo,
	cur *mockCouponUsageRepo,
	pub *mockEventPublisher,
) *CouponUseCase {
	return NewCouponUseCase(cr, cur, pub)
}

func validCreateCouponInput() CreateCouponInput {
	return CreateCouponInput{
		Code:          "SAVE10",
		Type:          "percentage",
		DiscountValue: 1000, // 10%
		CreatedBy:     "seller-1",
		ExpiresAt:     time.Now().Add(24 * time.Hour),
	}
}

// activeCoupon returns a coupon that is currently valid.
func activeCoupon() *domain.Coupon {
	return &domain.Coupon{
		ID:            "coupon-1",
		Code:          "SAVE10",
		Type:          domain.CouponTypePercentage,
		DiscountValue: 1000, // 10 %
		IsActive:      true,
		StartsAt:      time.Now().Add(-1 * time.Hour),
		ExpiresAt:     time.Now().Add(24 * time.Hour),
	}
}

// noopPublisher returns a publisher that always succeeds.
func noopPublisher() *mockEventPublisher {
	return &mockEventPublisher{
		publishFn: func(_ context.Context, _ string, _ interface{}) error { return nil },
	}
}

// ---------------------------------------------------------------------------
// CreateCoupon tests
// ---------------------------------------------------------------------------

func TestCreateCoupon(t *testing.T) {
	tests := []struct {
		name      string
		input     CreateCouponInput
		repoErr   error
		wantErr   string
		checkRes  func(t *testing.T, c *domain.Coupon)
	}{
		{
			name:  "success",
			input: validCreateCouponInput(),
			checkRes: func(t *testing.T, c *domain.Coupon) {
				assert.NotEmpty(t, c.ID)
				assert.Equal(t, "SAVE10", c.Code)
				assert.Equal(t, domain.CouponTypePercentage, c.Type)
				assert.Equal(t, int64(1000), c.DiscountValue)
				assert.True(t, c.IsActive)
			},
		},
		{
			name:    "missing code",
			input:   func() CreateCouponInput { i := validCreateCouponInput(); i.Code = ""; return i }(),
			wantErr: "coupon code is required",
		},
		{
			name:    "missing type",
			input:   func() CreateCouponInput { i := validCreateCouponInput(); i.Type = ""; return i }(),
			wantErr: "coupon type is required",
		},
		{
			name:    "zero discount value",
			input:   func() CreateCouponInput { i := validCreateCouponInput(); i.DiscountValue = 0; return i }(),
			wantErr: "discount value must be greater than 0",
		},
		{
			name:    "missing createdBy",
			input:   func() CreateCouponInput { i := validCreateCouponInput(); i.CreatedBy = ""; return i }(),
			wantErr: "created_by is required",
		},
		{
			name:    "missing expiresAt",
			input:   func() CreateCouponInput { i := validCreateCouponInput(); i.ExpiresAt = time.Time{}; return i }(),
			wantErr: "expires_at is required",
		},
		{
			name: "defaults_StartsAt set to now when zero",
			input: func() CreateCouponInput {
				i := validCreateCouponInput()
				i.StartsAt = time.Time{}
				return i
			}(),
			checkRes: func(t *testing.T, c *domain.Coupon) {
				assert.False(t, c.StartsAt.IsZero(), "StartsAt should default to now")
			},
		},
		{
			name: "defaults_Scope set to all when empty",
			input: func() CreateCouponInput {
				i := validCreateCouponInput()
				i.Scope = ""
				return i
			}(),
			checkRes: func(t *testing.T, c *domain.Coupon) {
				assert.Equal(t, domain.CouponScopeAll, c.Scope)
			},
		},
		{
			name:    "repo error",
			input:   validCreateCouponInput(),
			repoErr: errors.New("db connection lost"),
			wantErr: "db connection lost",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := &mockCouponRepo{
				createFn: func(_ context.Context, _ *domain.Coupon) error {
					return tt.repoErr
				},
			}
			uc := newTestCouponUseCase(cr, &mockCouponUsageRepo{}, noopPublisher())

			result, err := uc.CreateCoupon(context.Background(), tt.input)

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				assert.Nil(t, result)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, result)
			if tt.checkRes != nil {
				tt.checkRes(t, result)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// GetCoupon tests
// ---------------------------------------------------------------------------

func TestGetCoupon(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		repo    func() *mockCouponRepo
		wantErr string
	}{
		{
			name: "success",
			id:   "coupon-1",
			repo: func() *mockCouponRepo {
				return &mockCouponRepo{
					getByIDFn: func(_ context.Context, id string) (*domain.Coupon, error) {
						return &domain.Coupon{ID: id, Code: "SAVE10"}, nil
					},
				}
			},
		},
		{
			name:    "missing ID",
			id:      "",
			repo:    func() *mockCouponRepo { return &mockCouponRepo{} },
			wantErr: "coupon id is required",
		},
		{
			name: "not found",
			id:   "nonexistent",
			repo: func() *mockCouponRepo {
				return &mockCouponRepo{
					getByIDFn: func(_ context.Context, _ string) (*domain.Coupon, error) {
						return nil, errors.New("not found")
					},
				}
			},
			wantErr: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := newTestCouponUseCase(tt.repo(), &mockCouponUsageRepo{}, noopPublisher())

			result, err := uc.GetCoupon(context.Background(), tt.id)

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, tt.id, result.ID)
		})
	}
}

// ---------------------------------------------------------------------------
// ListCoupons tests
// ---------------------------------------------------------------------------

func TestListCoupons(t *testing.T) {
	tests := []struct {
		name         string
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
	}{
		{name: "success", page: 2, pageSize: 10, wantPage: 2, wantPageSize: 10},
		{name: "page 0 defaults to 1", page: 0, pageSize: 10, wantPage: 1, wantPageSize: 10},
		{name: "pageSize 0 defaults to 20", page: 1, pageSize: 0, wantPage: 1, wantPageSize: 20},
		{name: "pageSize over 100 capped to 100", page: 1, pageSize: 200, wantPage: 1, wantPageSize: 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedPage, capturedPageSize int
			cr := &mockCouponRepo{
				listAllFn: func(_ context.Context, page, pageSize int) ([]*domain.Coupon, int64, error) {
					capturedPage = page
					capturedPageSize = pageSize
					return []*domain.Coupon{{ID: "c1"}}, 1, nil
				},
			}
			uc := newTestCouponUseCase(cr, &mockCouponUsageRepo{}, noopPublisher())

			coupons, total, err := uc.ListCoupons(context.Background(), tt.page, tt.pageSize)

			require.NoError(t, err)
			assert.Len(t, coupons, 1)
			assert.Equal(t, int64(1), total)
			assert.Equal(t, tt.wantPage, capturedPage)
			assert.Equal(t, tt.wantPageSize, capturedPageSize)
		})
	}
}

// ---------------------------------------------------------------------------
// ListCouponsBySeller tests
// ---------------------------------------------------------------------------

func TestListCouponsBySeller(t *testing.T) {
	tests := []struct {
		name     string
		sellerID string
		page     int
		pageSize int
		wantErr  string
	}{
		{name: "success", sellerID: "seller-1", page: 1, pageSize: 10},
		{name: "missing sellerID", sellerID: "", page: 1, pageSize: 10, wantErr: "seller id is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := &mockCouponRepo{
				listBySellerFn: func(_ context.Context, _ string, _, _ int) ([]*domain.Coupon, int64, error) {
					return []*domain.Coupon{{ID: "c1"}}, 1, nil
				},
			}
			uc := newTestCouponUseCase(cr, &mockCouponUsageRepo{}, noopPublisher())

			coupons, total, err := uc.ListCouponsBySeller(context.Background(), tt.sellerID, tt.page, tt.pageSize)

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Len(t, coupons, 1)
			assert.Equal(t, int64(1), total)
		})
	}
}

// ---------------------------------------------------------------------------
// ValidateCoupon tests
// ---------------------------------------------------------------------------

func TestValidateCoupon(t *testing.T) {
	tests := []struct {
		name         string
		input        ValidateCouponInput
		coupon       *domain.Coupon
		repoErr      error
		userCount    int64
		wantErr      string
		wantDiscount int64
	}{
		{
			name: "success percentage",
			input: ValidateCouponInput{
				Code:       "SAVE10",
				UserID:     "user-1",
				OrderCents: 10000, // $100
			},
			coupon:       activeCoupon(), // 10% => 1000
			wantDiscount: 1000,
		},
		{
			name: "success fixed_amount",
			input: ValidateCouponInput{
				Code:       "FLAT20",
				UserID:     "user-1",
				OrderCents: 10000,
			},
			coupon: func() *domain.Coupon {
				c := activeCoupon()
				c.Code = "FLAT20"
				c.Type = domain.CouponTypeFixedAmount
				c.DiscountValue = 2000
				return c
			}(),
			wantDiscount: 2000,
		},
		{
			name:    "missing code",
			input:   ValidateCouponInput{Code: "", UserID: "user-1", OrderCents: 5000},
			wantErr: "coupon code is required",
		},
		{
			name: "inactive coupon",
			input: ValidateCouponInput{
				Code:       "SAVE10",
				UserID:     "user-1",
				OrderCents: 5000,
			},
			coupon: func() *domain.Coupon {
				c := activeCoupon()
				c.IsActive = false
				return c
			}(),
			wantErr: "coupon is not active",
		},
		{
			name: "expired coupon",
			input: ValidateCouponInput{
				Code:       "SAVE10",
				UserID:     "user-1",
				OrderCents: 5000,
			},
			coupon: func() *domain.Coupon {
				c := activeCoupon()
				c.ExpiresAt = time.Now().Add(-1 * time.Hour)
				return c
			}(),
			wantErr: "coupon has expired",
		},
		{
			name: "not started",
			input: ValidateCouponInput{
				Code:       "SAVE10",
				UserID:     "user-1",
				OrderCents: 5000,
			},
			coupon: func() *domain.Coupon {
				c := activeCoupon()
				c.StartsAt = time.Now().Add(24 * time.Hour)
				c.ExpiresAt = time.Now().Add(48 * time.Hour)
				return c
			}(),
			wantErr: "coupon is not yet active",
		},
		{
			name: "usage limit exceeded",
			input: ValidateCouponInput{
				Code:       "SAVE10",
				UserID:     "user-1",
				OrderCents: 5000,
			},
			coupon: func() *domain.Coupon {
				c := activeCoupon()
				c.UsageLimit = 5
				c.UsageCount = 5
				return c
			}(),
			wantErr: "coupon usage limit reached",
		},
		{
			name: "per-user limit exceeded",
			input: ValidateCouponInput{
				Code:       "SAVE10",
				UserID:     "user-1",
				OrderCents: 5000,
			},
			coupon: func() *domain.Coupon {
				c := activeCoupon()
				c.PerUserLimit = 1
				return c
			}(),
			userCount: 1,
			wantErr:   "per-user coupon usage limit reached",
		},
		{
			name: "min order not met",
			input: ValidateCouponInput{
				Code:       "SAVE10",
				UserID:     "user-1",
				OrderCents: 1000, // $10
			},
			coupon: func() *domain.Coupon {
				c := activeCoupon()
				c.MinOrderCents = 5000 // $50
				return c
			}(),
			wantErr: "order total does not meet minimum requirement",
		},
		{
			name: "maxDiscountCents cap",
			input: ValidateCouponInput{
				Code:       "SAVE10",
				UserID:     "user-1",
				OrderCents: 100000, // $1000 => 10% = $100
			},
			coupon: func() *domain.Coupon {
				c := activeCoupon()
				c.MaxDiscountCents = 2000 // cap at $20
				return c
			}(),
			wantDiscount: 2000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := &mockCouponRepo{
				getByCodeFn: func(_ context.Context, _ string) (*domain.Coupon, error) {
					if tt.repoErr != nil {
						return nil, tt.repoErr
					}
					return tt.coupon, nil
				},
			}
			cur := &mockCouponUsageRepo{
				countByUserFn: func(_ context.Context, _, _ string) (int64, error) {
					return tt.userCount, nil
				},
			}
			uc := newTestCouponUseCase(cr, cur, noopPublisher())

			coupon, discount, err := uc.ValidateCoupon(context.Background(), tt.input)

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, coupon)
			assert.Equal(t, tt.wantDiscount, discount)
		})
	}
}

// ---------------------------------------------------------------------------
// RedeemCoupon tests
// ---------------------------------------------------------------------------

func TestRedeemCoupon(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		userID  string
		orderID string
		cents   int64
		coupon  *domain.Coupon
		wantErr string
	}{
		{
			name:    "success",
			code:    "SAVE10",
			userID:  "user-1",
			orderID: "order-1",
			cents:   10000,
			coupon:  activeCoupon(),
		},
		{
			name:    "validation fails - missing code",
			code:    "",
			userID:  "user-1",
			orderID: "order-1",
			cents:   10000,
			wantErr: "coupon code is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var publishedSubject string
			cr := &mockCouponRepo{
				getByCodeFn: func(_ context.Context, _ string) (*domain.Coupon, error) {
					if tt.coupon == nil {
						return nil, errors.New("not found")
					}
					return tt.coupon, nil
				},
				incrementUsageFn: func(_ context.Context, _ string) error {
					return nil
				},
			}
			cur := &mockCouponUsageRepo{
				countByUserFn: func(_ context.Context, _, _ string) (int64, error) {
					return 0, nil
				},
				createFn: func(_ context.Context, _ *domain.CouponUsage) error {
					return nil
				},
			}
			pub := &mockEventPublisher{
				publishFn: func(_ context.Context, subject string, _ interface{}) error {
					publishedSubject = subject
					return nil
				},
			}
			uc := newTestCouponUseCase(cr, cur, pub)

			usage, err := uc.RedeemCoupon(context.Background(), tt.code, tt.userID, tt.orderID, tt.cents)

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				assert.Nil(t, usage)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, usage)
			assert.Equal(t, tt.coupon.ID, usage.CouponID)
			assert.Equal(t, tt.userID, usage.UserID)
			assert.Equal(t, tt.orderID, usage.OrderID)
			assert.Equal(t, domain.EventCouponRedeemed, publishedSubject)
		})
	}
}

// ---------------------------------------------------------------------------
// calculateDiscount tests
// ---------------------------------------------------------------------------

func TestCalculateDiscount(t *testing.T) {
	tests := []struct {
		name       string
		coupon     *domain.Coupon
		orderCents int64
		want       int64
	}{
		{
			name: "percentage",
			coupon: &domain.Coupon{
				Type:          domain.CouponTypePercentage,
				DiscountValue: 1000, // 10%
			},
			orderCents: 10000, // $100 => $10
			want:       1000,
		},
		{
			name: "fixed_amount",
			coupon: &domain.Coupon{
				Type:          domain.CouponTypeFixedAmount,
				DiscountValue: 2500,
			},
			orderCents: 10000,
			want:       2500,
		},
		{
			name: "free_shipping returns 0",
			coupon: &domain.Coupon{
				Type:          domain.CouponTypeFreeShipping,
				DiscountValue: 500,
			},
			orderCents: 10000,
			want:       0,
		},
		{
			name: "cap at maxDiscountCents",
			coupon: &domain.Coupon{
				Type:             domain.CouponTypePercentage,
				DiscountValue:    5000, // 50%
				MaxDiscountCents: 1000,
			},
			orderCents: 10000, // 50% = 5000, capped at 1000
			want:       1000,
		},
		{
			name: "cap at orderCents",
			coupon: &domain.Coupon{
				Type:          domain.CouponTypeFixedAmount,
				DiscountValue: 20000,
			},
			orderCents: 5000, // discount 20000 capped at order 5000
			want:       5000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &CouponUseCase{}
			got := uc.calculateDiscount(tt.coupon, tt.orderCents)
			assert.Equal(t, tt.want, got)
		})
	}
}
