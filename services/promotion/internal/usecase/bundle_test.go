package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/southern-martin/ecommerce/services/promotion/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Hand-written function-field mock
// ---------------------------------------------------------------------------

type mockBundleRepo struct {
	getByIDFn    func(ctx context.Context, id string) (*domain.Bundle, error)
	listBySellerFn func(ctx context.Context, sellerID string, page, pageSize int) ([]*domain.Bundle, int64, error)
	listActiveFn func(ctx context.Context, page, pageSize int) ([]*domain.Bundle, int64, error)
	createFn     func(ctx context.Context, bundle *domain.Bundle) error
	updateFn     func(ctx context.Context, bundle *domain.Bundle) error
}

func (m *mockBundleRepo) GetByID(ctx context.Context, id string) (*domain.Bundle, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockBundleRepo) ListBySeller(ctx context.Context, sellerID string, page, pageSize int) ([]*domain.Bundle, int64, error) {
	return m.listBySellerFn(ctx, sellerID, page, pageSize)
}
func (m *mockBundleRepo) ListActive(ctx context.Context, page, pageSize int) ([]*domain.Bundle, int64, error) {
	return m.listActiveFn(ctx, page, pageSize)
}
func (m *mockBundleRepo) Create(ctx context.Context, bundle *domain.Bundle) error {
	return m.createFn(ctx, bundle)
}
func (m *mockBundleRepo) Update(ctx context.Context, bundle *domain.Bundle) error {
	return m.updateFn(ctx, bundle)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func validCreateBundleInput() CreateBundleInput {
	return CreateBundleInput{
		Name:             "Starter Pack",
		SellerID:         "seller-1",
		ProductIDs:       []string{"prod-1", "prod-2", "prod-3"},
		BundlePriceCents: 4999,
		SavingsCents:     1500,
	}
}

// ---------------------------------------------------------------------------
// CreateBundle tests
// ---------------------------------------------------------------------------

func TestCreateBundle(t *testing.T) {
	tests := []struct {
		name    string
		input   CreateBundleInput
		repoErr error
		wantErr string
		checkRes func(t *testing.T, b *domain.Bundle)
	}{
		{
			name:  "success",
			input: validCreateBundleInput(),
			checkRes: func(t *testing.T, b *domain.Bundle) {
				assert.NotEmpty(t, b.ID)
				assert.Equal(t, "Starter Pack", b.Name)
				assert.Equal(t, "seller-1", b.SellerID)
				assert.Equal(t, []string{"prod-1", "prod-2", "prod-3"}, b.ProductIDs)
				assert.Equal(t, int64(4999), b.BundlePriceCents)
				assert.Equal(t, int64(1500), b.SavingsCents)
				assert.True(t, b.IsActive)
			},
		},
		{
			name:    "missing name",
			input:   func() CreateBundleInput { i := validCreateBundleInput(); i.Name = ""; return i }(),
			wantErr: "bundle name is required",
		},
		{
			name:    "missing sellerID",
			input:   func() CreateBundleInput { i := validCreateBundleInput(); i.SellerID = ""; return i }(),
			wantErr: "seller_id is required",
		},
		{
			name: "less than 2 products",
			input: func() CreateBundleInput {
				i := validCreateBundleInput()
				i.ProductIDs = []string{"prod-1"}
				return i
			}(),
			wantErr: "at least two products are required",
		},
		{
			name: "zero products",
			input: func() CreateBundleInput {
				i := validCreateBundleInput()
				i.ProductIDs = nil
				return i
			}(),
			wantErr: "at least two products are required",
		},
		{
			name: "zero price",
			input: func() CreateBundleInput {
				i := validCreateBundleInput()
				i.BundlePriceCents = 0
				return i
			}(),
			wantErr: "bundle price must be greater than 0",
		},
		{
			name: "negative price",
			input: func() CreateBundleInput {
				i := validCreateBundleInput()
				i.BundlePriceCents = -100
				return i
			}(),
			wantErr: "bundle price must be greater than 0",
		},
		{
			name:    "repo error",
			input:   validCreateBundleInput(),
			repoErr: errors.New("db insert failed"),
			wantErr: "db insert failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockBundleRepo{
				createFn: func(_ context.Context, _ *domain.Bundle) error {
					return tt.repoErr
				},
			}
			uc := NewBundleUseCase(repo)

			result, err := uc.CreateBundle(context.Background(), tt.input)

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
// GetBundle tests
// ---------------------------------------------------------------------------

func TestGetBundle(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		repo    func() *mockBundleRepo
		wantErr string
	}{
		{
			name: "success",
			id:   "bundle-1",
			repo: func() *mockBundleRepo {
				return &mockBundleRepo{
					getByIDFn: func(_ context.Context, id string) (*domain.Bundle, error) {
						return &domain.Bundle{ID: id, Name: "Pack"}, nil
					},
				}
			},
		},
		{
			name:    "missing ID",
			id:      "",
			repo:    func() *mockBundleRepo { return &mockBundleRepo{} },
			wantErr: "bundle id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewBundleUseCase(tt.repo())

			result, err := uc.GetBundle(context.Background(), tt.id)

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
// ListBundles tests
// ---------------------------------------------------------------------------

func TestListBundles(t *testing.T) {
	tests := []struct {
		name         string
		sellerID     string
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
		wantErr      string
	}{
		{name: "success", sellerID: "seller-1", page: 2, pageSize: 15, wantPage: 2, wantPageSize: 15},
		{name: "missing sellerID", sellerID: "", page: 1, pageSize: 10, wantErr: "seller id is required"},
		{name: "page 0 defaults to 1", sellerID: "seller-1", page: 0, pageSize: 10, wantPage: 1, wantPageSize: 10},
		{name: "pageSize 0 defaults to 20", sellerID: "seller-1", page: 1, pageSize: 0, wantPage: 1, wantPageSize: 20},
		{name: "pageSize over 100 capped", sellerID: "seller-1", page: 1, pageSize: 200, wantPage: 1, wantPageSize: 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedPage, capturedPageSize int
			repo := &mockBundleRepo{
				listBySellerFn: func(_ context.Context, _ string, page, pageSize int) ([]*domain.Bundle, int64, error) {
					capturedPage = page
					capturedPageSize = pageSize
					return []*domain.Bundle{{ID: "b1"}}, 1, nil
				},
			}
			uc := NewBundleUseCase(repo)

			bundles, total, err := uc.ListBundles(context.Background(), tt.sellerID, tt.page, tt.pageSize)

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Len(t, bundles, 1)
			assert.Equal(t, int64(1), total)
			assert.Equal(t, tt.wantPage, capturedPage)
			assert.Equal(t, tt.wantPageSize, capturedPageSize)
		})
	}
}

// ---------------------------------------------------------------------------
// ListActiveBundles tests
// ---------------------------------------------------------------------------

func TestListActiveBundles(t *testing.T) {
	tests := []struct {
		name         string
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
	}{
		{name: "normal values", page: 3, pageSize: 25, wantPage: 3, wantPageSize: 25},
		{name: "page 0 defaults to 1", page: 0, pageSize: 10, wantPage: 1, wantPageSize: 10},
		{name: "pageSize 0 defaults to 20", page: 1, pageSize: 0, wantPage: 1, wantPageSize: 20},
		{name: "pageSize over 100 capped", page: 1, pageSize: 999, wantPage: 1, wantPageSize: 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedPage, capturedPageSize int
			repo := &mockBundleRepo{
				listActiveFn: func(_ context.Context, page, pageSize int) ([]*domain.Bundle, int64, error) {
					capturedPage = page
					capturedPageSize = pageSize
					return []*domain.Bundle{{ID: "b1", IsActive: true}}, 1, nil
				},
			}
			uc := NewBundleUseCase(repo)

			bundles, total, err := uc.ListActiveBundles(context.Background(), tt.page, tt.pageSize)

			require.NoError(t, err)
			assert.Len(t, bundles, 1)
			assert.Equal(t, int64(1), total)
			assert.Equal(t, tt.wantPage, capturedPage)
			assert.Equal(t, tt.wantPageSize, capturedPageSize)
		})
	}
}
