package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/cms/internal/domain"
)

// ---------------------------------------------------------------------------
// Hand-written function-field mocks
// ---------------------------------------------------------------------------

// --- BannerRepository mock ---

type mockBannerRepo struct {
	getByIDFn   func(ctx context.Context, id string) (*domain.Banner, error)
	listActiveFn func(ctx context.Context, position string) ([]domain.Banner, error)
	listAllFn   func(ctx context.Context, page, pageSize int) ([]domain.Banner, int64, error)
	createFn    func(ctx context.Context, banner *domain.Banner) error
	updateFn    func(ctx context.Context, banner *domain.Banner) error
	deleteFn    func(ctx context.Context, id string) error
}

func (m *mockBannerRepo) GetByID(ctx context.Context, id string) (*domain.Banner, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}
func (m *mockBannerRepo) ListActive(ctx context.Context, position string) ([]domain.Banner, error) {
	if m.listActiveFn != nil {
		return m.listActiveFn(ctx, position)
	}
	return nil, nil
}
func (m *mockBannerRepo) ListAll(ctx context.Context, page, pageSize int) ([]domain.Banner, int64, error) {
	if m.listAllFn != nil {
		return m.listAllFn(ctx, page, pageSize)
	}
	return nil, 0, nil
}
func (m *mockBannerRepo) Create(ctx context.Context, banner *domain.Banner) error {
	if m.createFn != nil {
		return m.createFn(ctx, banner)
	}
	return nil
}
func (m *mockBannerRepo) Update(ctx context.Context, banner *domain.Banner) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, banner)
	}
	return nil
}
func (m *mockBannerRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

// --- PageRepository mock ---

type mockPageRepo struct {
	getByIDFn      func(ctx context.Context, id string) (*domain.Page, error)
	getBySlugFn    func(ctx context.Context, slug string) (*domain.Page, error)
	listPublishedFn func(ctx context.Context, page, pageSize int) ([]domain.Page, int64, error)
	listAllFn      func(ctx context.Context, page, pageSize int) ([]domain.Page, int64, error)
	createFn       func(ctx context.Context, pg *domain.Page) error
	updateFn       func(ctx context.Context, pg *domain.Page) error
	deleteFn       func(ctx context.Context, id string) error
}

func (m *mockPageRepo) GetByID(ctx context.Context, id string) (*domain.Page, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}
func (m *mockPageRepo) GetBySlug(ctx context.Context, slug string) (*domain.Page, error) {
	if m.getBySlugFn != nil {
		return m.getBySlugFn(ctx, slug)
	}
	return nil, errors.New("not found")
}
func (m *mockPageRepo) ListPublished(ctx context.Context, page, pageSize int) ([]domain.Page, int64, error) {
	if m.listPublishedFn != nil {
		return m.listPublishedFn(ctx, page, pageSize)
	}
	return nil, 0, nil
}
func (m *mockPageRepo) ListAll(ctx context.Context, page, pageSize int) ([]domain.Page, int64, error) {
	if m.listAllFn != nil {
		return m.listAllFn(ctx, page, pageSize)
	}
	return nil, 0, nil
}
func (m *mockPageRepo) Create(ctx context.Context, pg *domain.Page) error {
	if m.createFn != nil {
		return m.createFn(ctx, pg)
	}
	return nil
}
func (m *mockPageRepo) Update(ctx context.Context, pg *domain.Page) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, pg)
	}
	return nil
}
func (m *mockPageRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

// --- ScheduleRepository mock ---

type mockScheduleRepo struct {
	getPendingFn   func(ctx context.Context) ([]domain.ContentSchedule, error)
	createFn       func(ctx context.Context, schedule *domain.ContentSchedule) error
	markExecutedFn func(ctx context.Context, id string) error
}

func (m *mockScheduleRepo) GetPending(ctx context.Context) ([]domain.ContentSchedule, error) {
	if m.getPendingFn != nil {
		return m.getPendingFn(ctx)
	}
	return nil, nil
}
func (m *mockScheduleRepo) Create(ctx context.Context, schedule *domain.ContentSchedule) error {
	if m.createFn != nil {
		return m.createFn(ctx, schedule)
	}
	return nil
}
func (m *mockScheduleRepo) MarkExecuted(ctx context.Context, id string) error {
	if m.markExecutedFn != nil {
		return m.markExecutedFn(ctx, id)
	}
	return nil
}

// --- EventPublisher mock ---

type mockCMSEventPub struct {
	publishFn func(ctx context.Context, subject string, data interface{}) error
}

func (m *mockCMSEventPub) Publish(ctx context.Context, subject string, data interface{}) error {
	if m.publishFn != nil {
		return m.publishFn(ctx, subject, data)
	}
	return nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func defaultCMSMocks() (*mockBannerRepo, *mockPageRepo, *mockScheduleRepo, *mockCMSEventPub) {
	return &mockBannerRepo{},
		&mockPageRepo{},
		&mockScheduleRepo{},
		&mockCMSEventPub{}
}

// ===========================================================================
// CreateBanner tests
// ===========================================================================

func TestCreateBanner_Success(t *testing.T) {
	bRepo, _, _, pub := defaultCMSMocks()

	var saved *domain.Banner
	bRepo.createFn = func(_ context.Context, b *domain.Banner) error {
		saved = b
		return nil
	}

	uc := NewBannerUseCase(bRepo, pub)
	banner := &domain.Banner{
		Title:    "Summer Sale",
		ImageURL: "https://cdn.example.com/summer.jpg",
		LinkURL:  "https://example.com/sale",
		Position: "homepage_top",
		IsActive: true,
	}
	err := uc.CreateBanner(context.Background(), banner)

	require.NoError(t, err)
	assert.NotEmpty(t, banner.ID)
	assert.Equal(t, "Summer Sale", banner.Title)
	assert.NotNil(t, saved)
}

func TestCreateBanner_SanitizesFields(t *testing.T) {
	bRepo, _, _, pub := defaultCMSMocks()
	bRepo.createFn = func(_ context.Context, _ *domain.Banner) error { return nil }

	uc := NewBannerUseCase(bRepo, pub)
	banner := &domain.Banner{
		Title:   "<script>alert('xss')</script>Summer Sale",
		LinkURL: "<b>https://example.com</b>",
	}
	err := uc.CreateBanner(context.Background(), banner)

	require.NoError(t, err)
	// sanitizeText strips all HTML
	assert.NotContains(t, banner.Title, "<script>")
	assert.NotContains(t, banner.LinkURL, "<b>")
}

func TestCreateBanner_RepoError(t *testing.T) {
	bRepo, _, _, pub := defaultCMSMocks()

	bRepo.createFn = func(_ context.Context, _ *domain.Banner) error {
		return errors.New("db error")
	}

	uc := NewBannerUseCase(bRepo, pub)
	err := uc.CreateBanner(context.Background(), &domain.Banner{Title: "Banner"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create banner")
}

// ===========================================================================
// UpdateBanner tests
// ===========================================================================

func TestUpdateBanner_Success(t *testing.T) {
	bRepo, _, _, pub := defaultCMSMocks()

	bRepo.getByIDFn = func(_ context.Context, id string) (*domain.Banner, error) {
		return &domain.Banner{ID: id, Title: "Old Title"}, nil
	}

	var updated *domain.Banner
	bRepo.updateFn = func(_ context.Context, b *domain.Banner) error {
		updated = b
		return nil
	}

	uc := NewBannerUseCase(bRepo, pub)
	err := uc.UpdateBanner(context.Background(), &domain.Banner{
		ID:    "banner-1",
		Title: "New Title",
	})

	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "New Title", updated.Title)
}

func TestUpdateBanner_NotFound(t *testing.T) {
	bRepo, _, _, pub := defaultCMSMocks()

	bRepo.getByIDFn = func(_ context.Context, _ string) (*domain.Banner, error) {
		return nil, errors.New("not found")
	}

	uc := NewBannerUseCase(bRepo, pub)
	err := uc.UpdateBanner(context.Background(), &domain.Banner{ID: "banner-missing"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "banner not found")
}

func TestUpdateBanner_RepoError(t *testing.T) {
	bRepo, _, _, pub := defaultCMSMocks()

	bRepo.getByIDFn = func(_ context.Context, id string) (*domain.Banner, error) {
		return &domain.Banner{ID: id}, nil
	}
	bRepo.updateFn = func(_ context.Context, _ *domain.Banner) error {
		return errors.New("db error")
	}

	uc := NewBannerUseCase(bRepo, pub)
	err := uc.UpdateBanner(context.Background(), &domain.Banner{ID: "banner-1"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update banner")
}

// ===========================================================================
// DeleteBanner tests
// ===========================================================================

func TestDeleteBanner_Success(t *testing.T) {
	bRepo, _, _, pub := defaultCMSMocks()

	bRepo.getByIDFn = func(_ context.Context, id string) (*domain.Banner, error) {
		return &domain.Banner{ID: id}, nil
	}

	var deletedID string
	bRepo.deleteFn = func(_ context.Context, id string) error {
		deletedID = id
		return nil
	}

	uc := NewBannerUseCase(bRepo, pub)
	err := uc.DeleteBanner(context.Background(), "banner-1")

	require.NoError(t, err)
	assert.Equal(t, "banner-1", deletedID)
}

func TestDeleteBanner_NotFound(t *testing.T) {
	bRepo, _, _, pub := defaultCMSMocks()

	bRepo.getByIDFn = func(_ context.Context, _ string) (*domain.Banner, error) {
		return nil, errors.New("not found")
	}

	uc := NewBannerUseCase(bRepo, pub)
	err := uc.DeleteBanner(context.Background(), "banner-missing")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "banner not found")
}

func TestDeleteBanner_RepoError(t *testing.T) {
	bRepo, _, _, pub := defaultCMSMocks()

	bRepo.getByIDFn = func(_ context.Context, id string) (*domain.Banner, error) {
		return &domain.Banner{ID: id}, nil
	}
	bRepo.deleteFn = func(_ context.Context, _ string) error {
		return errors.New("db error")
	}

	uc := NewBannerUseCase(bRepo, pub)
	err := uc.DeleteBanner(context.Background(), "banner-1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete banner")
}

// ===========================================================================
// ListActiveBanners tests
// ===========================================================================

func TestListActiveBanners_Success(t *testing.T) {
	bRepo, _, _, pub := defaultCMSMocks()

	bRepo.listActiveFn = func(_ context.Context, position string) ([]domain.Banner, error) {
		return []domain.Banner{
			{ID: "b-1", Position: position, IsActive: true},
			{ID: "b-2", Position: position, IsActive: true},
		}, nil
	}

	uc := NewBannerUseCase(bRepo, pub)
	banners, err := uc.ListActiveBanners(context.Background(), "homepage_top")

	require.NoError(t, err)
	assert.Len(t, banners, 2)
}

// ===========================================================================
// ListAllBanners tests
// ===========================================================================

func TestListAllBanners_Success(t *testing.T) {
	bRepo, _, _, pub := defaultCMSMocks()

	bRepo.listAllFn = func(_ context.Context, _, _ int) ([]domain.Banner, int64, error) {
		return []domain.Banner{
			{ID: "b-1"},
			{ID: "b-2"},
			{ID: "b-3"},
		}, 3, nil
	}

	uc := NewBannerUseCase(bRepo, pub)
	banners, total, err := uc.ListAllBanners(context.Background(), 1, 20)

	require.NoError(t, err)
	assert.Len(t, banners, 3)
	assert.Equal(t, int64(3), total)
}

func TestListAllBanners_DefaultsPagination(t *testing.T) {
	bRepo, _, _, pub := defaultCMSMocks()

	var capturedPage, capturedPageSize int
	bRepo.listAllFn = func(_ context.Context, page, pageSize int) ([]domain.Banner, int64, error) {
		capturedPage = page
		capturedPageSize = pageSize
		return nil, 0, nil
	}

	uc := NewBannerUseCase(bRepo, pub)

	_, _, _ = uc.ListAllBanners(context.Background(), 0, 20)
	assert.Equal(t, 1, capturedPage)

	_, _, _ = uc.ListAllBanners(context.Background(), 1, 0)
	assert.Equal(t, 20, capturedPageSize)

	_, _, _ = uc.ListAllBanners(context.Background(), 1, 200)
	assert.Equal(t, 20, capturedPageSize)
}
