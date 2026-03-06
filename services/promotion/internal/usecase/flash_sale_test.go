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

type mockFlashSaleRepo struct {
	getByIDFn  func(ctx context.Context, id string) (*domain.FlashSale, error)
	listAllFn  func(ctx context.Context, page, pageSize int) ([]*domain.FlashSale, int64, error)
	listActiveFn func(ctx context.Context) ([]*domain.FlashSale, error)
	createFn   func(ctx context.Context, fs *domain.FlashSale) error
	updateFn   func(ctx context.Context, fs *domain.FlashSale) error
}

func (m *mockFlashSaleRepo) GetByID(ctx context.Context, id string) (*domain.FlashSale, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockFlashSaleRepo) ListAll(ctx context.Context, page, pageSize int) ([]*domain.FlashSale, int64, error) {
	return m.listAllFn(ctx, page, pageSize)
}
func (m *mockFlashSaleRepo) ListActive(ctx context.Context) ([]*domain.FlashSale, error) {
	return m.listActiveFn(ctx)
}
func (m *mockFlashSaleRepo) Create(ctx context.Context, fs *domain.FlashSale) error {
	return m.createFn(ctx, fs)
}
func (m *mockFlashSaleRepo) Update(ctx context.Context, fs *domain.FlashSale) error {
	return m.updateFn(ctx, fs)
}

type mockFlashSaleItemRepo struct {
	getByFlashSaleIDFn  func(ctx context.Context, flashSaleID string) ([]*domain.FlashSaleItem, error)
	createFn            func(ctx context.Context, item *domain.FlashSaleItem) error
	incrementSoldCountFn func(ctx context.Context, id string) error
}

func (m *mockFlashSaleItemRepo) GetByFlashSaleID(ctx context.Context, flashSaleID string) ([]*domain.FlashSaleItem, error) {
	return m.getByFlashSaleIDFn(ctx, flashSaleID)
}
func (m *mockFlashSaleItemRepo) Create(ctx context.Context, item *domain.FlashSaleItem) error {
	return m.createFn(ctx, item)
}
func (m *mockFlashSaleItemRepo) IncrementSoldCount(ctx context.Context, id string) error {
	return m.incrementSoldCountFn(ctx, id)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func newTestFlashSaleUseCase(
	fsr *mockFlashSaleRepo,
	fsir *mockFlashSaleItemRepo,
	pub *mockEventPublisher,
) *FlashSaleUseCase {
	return NewFlashSaleUseCase(fsr, fsir, pub)
}

func validCreateFlashSaleInput() CreateFlashSaleInput {
	return CreateFlashSaleInput{
		Name:     "Summer Flash Sale",
		StartsAt: time.Now().Add(-1 * time.Hour),
		EndsAt:   time.Now().Add(24 * time.Hour),
	}
}

// ---------------------------------------------------------------------------
// CreateFlashSale tests
// ---------------------------------------------------------------------------

func TestCreateFlashSale(t *testing.T) {
	tests := []struct {
		name       string
		input      CreateFlashSaleInput
		repoErr    error
		wantErr    string
		wantPublish bool
		checkRes   func(t *testing.T, fs *domain.FlashSale)
	}{
		{
			name:        "success",
			input:       validCreateFlashSaleInput(),
			wantPublish: true, // currently active
			checkRes: func(t *testing.T, fs *domain.FlashSale) {
				assert.NotEmpty(t, fs.ID)
				assert.Equal(t, "Summer Flash Sale", fs.Name)
				assert.True(t, fs.IsActive)
			},
		},
		{
			name:    "missing name",
			input:   func() CreateFlashSaleInput { i := validCreateFlashSaleInput(); i.Name = ""; return i }(),
			wantErr: "flash sale name is required",
		},
		{
			name: "missing startsAt",
			input: func() CreateFlashSaleInput {
				i := validCreateFlashSaleInput()
				i.StartsAt = time.Time{}
				return i
			}(),
			wantErr: "starts_at is required",
		},
		{
			name: "missing endsAt",
			input: func() CreateFlashSaleInput {
				i := validCreateFlashSaleInput()
				i.EndsAt = time.Time{}
				return i
			}(),
			wantErr: "ends_at is required",
		},
		{
			name: "endsAt before startsAt",
			input: func() CreateFlashSaleInput {
				i := validCreateFlashSaleInput()
				i.StartsAt = time.Now().Add(24 * time.Hour)
				i.EndsAt = time.Now().Add(1 * time.Hour)
				return i
			}(),
			wantErr: "ends_at must be after starts_at",
		},
		{
			name: "with items",
			input: func() CreateFlashSaleInput {
				i := validCreateFlashSaleInput()
				i.Items = []CreateFlashSaleItemInput{
					{ProductID: "prod-1", VariantID: "var-1", SalePriceCents: 999, QuantityLimit: 50},
					{ProductID: "prod-2", VariantID: "var-2", SalePriceCents: 1999, QuantityLimit: 20},
				}
				return i
			}(),
			wantPublish: true,
			checkRes: func(t *testing.T, fs *domain.FlashSale) {
				assert.Len(t, fs.Items, 2)
				assert.Equal(t, "prod-1", fs.Items[0].ProductID)
				assert.Equal(t, "var-1", fs.Items[0].VariantID)
				assert.Equal(t, int64(999), fs.Items[0].SalePriceCents)
				assert.Equal(t, 50, fs.Items[0].QuantityLimit)
				assert.Equal(t, "prod-2", fs.Items[1].ProductID)
			},
		},
		{
			name: "publishes event when currently active",
			input: func() CreateFlashSaleInput {
				i := validCreateFlashSaleInput()
				i.StartsAt = time.Now().Add(-1 * time.Hour)
				i.EndsAt = time.Now().Add(1 * time.Hour)
				return i
			}(),
			wantPublish: true,
		},
		{
			name: "no publish when future sale",
			input: func() CreateFlashSaleInput {
				i := validCreateFlashSaleInput()
				i.StartsAt = time.Now().Add(24 * time.Hour)
				i.EndsAt = time.Now().Add(48 * time.Hour)
				return i
			}(),
			wantPublish: false,
		},
		{
			name:    "repo error",
			input:   validCreateFlashSaleInput(),
			repoErr: errors.New("db write failed"),
			wantErr: "db write failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var published bool
			fsr := &mockFlashSaleRepo{
				createFn: func(_ context.Context, _ *domain.FlashSale) error {
					return tt.repoErr
				},
			}
			fsir := &mockFlashSaleItemRepo{}
			pub := &mockEventPublisher{
				publishFn: func(_ context.Context, subject string, _ interface{}) error {
					if subject == domain.EventFlashSaleStarted {
						published = true
					}
					return nil
				},
			}
			uc := newTestFlashSaleUseCase(fsr, fsir, pub)

			result, err := uc.CreateFlashSale(context.Background(), tt.input)

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				assert.Nil(t, result)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, tt.wantPublish, published, "publish event mismatch")
			if tt.checkRes != nil {
				tt.checkRes(t, result)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// GetFlashSale tests
// ---------------------------------------------------------------------------

func TestGetFlashSale(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		repo    func() *mockFlashSaleRepo
		wantErr string
	}{
		{
			name: "success",
			id:   "fs-1",
			repo: func() *mockFlashSaleRepo {
				return &mockFlashSaleRepo{
					getByIDFn: func(_ context.Context, id string) (*domain.FlashSale, error) {
						return &domain.FlashSale{ID: id, Name: "Sale"}, nil
					},
				}
			},
		},
		{
			name:    "missing ID",
			id:      "",
			repo:    func() *mockFlashSaleRepo { return &mockFlashSaleRepo{} },
			wantErr: "flash sale id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := newTestFlashSaleUseCase(tt.repo(), &mockFlashSaleItemRepo{}, noopPublisher())

			result, err := uc.GetFlashSale(context.Background(), tt.id)

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
// ListFlashSales tests
// ---------------------------------------------------------------------------

func TestListFlashSales(t *testing.T) {
	tests := []struct {
		name         string
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
	}{
		{name: "normal values", page: 2, pageSize: 15, wantPage: 2, wantPageSize: 15},
		{name: "page 0 defaults to 1", page: 0, pageSize: 10, wantPage: 1, wantPageSize: 10},
		{name: "negative page defaults to 1", page: -5, pageSize: 10, wantPage: 1, wantPageSize: 10},
		{name: "pageSize 0 defaults to 20", page: 1, pageSize: 0, wantPage: 1, wantPageSize: 20},
		{name: "pageSize over 100 capped", page: 1, pageSize: 150, wantPage: 1, wantPageSize: 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedPage, capturedPageSize int
			fsr := &mockFlashSaleRepo{
				listAllFn: func(_ context.Context, page, pageSize int) ([]*domain.FlashSale, int64, error) {
					capturedPage = page
					capturedPageSize = pageSize
					return []*domain.FlashSale{{ID: "fs-1"}}, 1, nil
				},
			}
			uc := newTestFlashSaleUseCase(fsr, &mockFlashSaleItemRepo{}, noopPublisher())

			sales, total, err := uc.ListFlashSales(context.Background(), tt.page, tt.pageSize)

			require.NoError(t, err)
			assert.Len(t, sales, 1)
			assert.Equal(t, int64(1), total)
			assert.Equal(t, tt.wantPage, capturedPage)
			assert.Equal(t, tt.wantPageSize, capturedPageSize)
		})
	}
}

// ---------------------------------------------------------------------------
// ListActiveFlashSales tests
// ---------------------------------------------------------------------------

func TestListActiveFlashSales(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		fsr := &mockFlashSaleRepo{
			listActiveFn: func(_ context.Context) ([]*domain.FlashSale, error) {
				return []*domain.FlashSale{
					{ID: "fs-1", Name: "Active Sale", IsActive: true},
				}, nil
			},
		}
		uc := newTestFlashSaleUseCase(fsr, &mockFlashSaleItemRepo{}, noopPublisher())

		sales, err := uc.ListActiveFlashSales(context.Background())

		require.NoError(t, err)
		assert.Len(t, sales, 1)
		assert.Equal(t, "Active Sale", sales[0].Name)
	})
}

// ---------------------------------------------------------------------------
// UpdateFlashSale tests
// ---------------------------------------------------------------------------

func TestUpdateFlashSale(t *testing.T) {
	tests := []struct {
		name          string
		flashSale     *domain.FlashSale
		repoErr       error
		wantErr       string
		wantSubject   string
		wantPublished bool
	}{
		{
			name: "publishes ended event when deactivated",
			flashSale: &domain.FlashSale{
				ID:       "fs-1",
				Name:     "Old Sale",
				IsActive: false,
				StartsAt: time.Now().Add(-2 * time.Hour),
				EndsAt:   time.Now().Add(1 * time.Hour),
			},
			wantSubject:   domain.EventFlashSaleEnded,
			wantPublished: true,
		},
		{
			name: "publishes started event when in time range",
			flashSale: &domain.FlashSale{
				ID:       "fs-2",
				Name:     "Current Sale",
				IsActive: true,
				StartsAt: time.Now().Add(-1 * time.Hour),
				EndsAt:   time.Now().Add(2 * time.Hour),
			},
			wantSubject:   domain.EventFlashSaleStarted,
			wantPublished: true,
		},
		{
			name: "no event when active but future",
			flashSale: &domain.FlashSale{
				ID:       "fs-3",
				Name:     "Future Sale",
				IsActive: true,
				StartsAt: time.Now().Add(24 * time.Hour),
				EndsAt:   time.Now().Add(48 * time.Hour),
			},
			wantPublished: false,
		},
		{
			name: "repo error",
			flashSale: &domain.FlashSale{
				ID:       "fs-4",
				Name:     "Fail Sale",
				IsActive: true,
				StartsAt: time.Now().Add(-1 * time.Hour),
				EndsAt:   time.Now().Add(2 * time.Hour),
			},
			repoErr: errors.New("update failed"),
			wantErr: "update failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var publishedSubject string
			var published bool
			fsr := &mockFlashSaleRepo{
				updateFn: func(_ context.Context, _ *domain.FlashSale) error {
					return tt.repoErr
				},
			}
			pub := &mockEventPublisher{
				publishFn: func(_ context.Context, subject string, _ interface{}) error {
					published = true
					publishedSubject = subject
					return nil
				},
			}
			uc := newTestFlashSaleUseCase(fsr, &mockFlashSaleItemRepo{}, pub)

			err := uc.UpdateFlashSale(context.Background(), tt.flashSale)

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantPublished, published, "publish call mismatch")
			if tt.wantPublished {
				assert.Equal(t, tt.wantSubject, publishedSubject)
			}
		})
	}
}
