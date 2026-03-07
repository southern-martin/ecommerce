package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/cms/internal/domain"
)

// ===========================================================================
// CreatePage tests
// ===========================================================================

func TestCreatePage_Success(t *testing.T) {
	_, pgRepo, _, pub := defaultCMSMocks()

	var saved *domain.Page
	pgRepo.createFn = func(_ context.Context, pg *domain.Page) error {
		saved = pg
		return nil
	}

	uc := NewPageUseCase(pgRepo, pub)
	page := &domain.Page{
		Title:       "About Us",
		ContentHTML: "<p>We are a company.</p>",
	}
	err := uc.CreatePage(context.Background(), page)

	require.NoError(t, err)
	assert.NotEmpty(t, page.ID)
	assert.Equal(t, "about-us", page.Slug)
	assert.Equal(t, domain.PageStatusDraft, page.Status)
	assert.NotNil(t, saved)
}

func TestCreatePage_GeneratesSlug(t *testing.T) {
	_, pgRepo, _, pub := defaultCMSMocks()
	pgRepo.createFn = func(_ context.Context, _ *domain.Page) error { return nil }

	uc := NewPageUseCase(pgRepo, pub)
	page := &domain.Page{Title: "Terms & Conditions!!!"}
	err := uc.CreatePage(context.Background(), page)

	require.NoError(t, err)
	assert.Equal(t, "terms-conditions", page.Slug)
}

func TestCreatePage_SanitizesHTML(t *testing.T) {
	_, pgRepo, _, pub := defaultCMSMocks()
	pgRepo.createFn = func(_ context.Context, _ *domain.Page) error { return nil }

	uc := NewPageUseCase(pgRepo, pub)
	page := &domain.Page{
		Title:       "Test Page",
		ContentHTML: "<script>alert('xss')</script><p>Safe content</p>",
	}
	err := uc.CreatePage(context.Background(), page)

	require.NoError(t, err)
	assert.NotContains(t, page.ContentHTML, "<script>")
	assert.Contains(t, page.ContentHTML, "<p>Safe content</p>")
}

func TestCreatePage_RepoError(t *testing.T) {
	_, pgRepo, _, pub := defaultCMSMocks()

	pgRepo.createFn = func(_ context.Context, _ *domain.Page) error {
		return errors.New("db error")
	}

	uc := NewPageUseCase(pgRepo, pub)
	err := uc.CreatePage(context.Background(), &domain.Page{Title: "Test"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create page")
}

// ===========================================================================
// GetPage tests
// ===========================================================================

func TestGetPage_Success(t *testing.T) {
	_, pgRepo, _, pub := defaultCMSMocks()

	pgRepo.getByIDFn = func(_ context.Context, id string) (*domain.Page, error) {
		return &domain.Page{ID: id, Title: "About Us", Slug: "about-us"}, nil
	}

	uc := NewPageUseCase(pgRepo, pub)
	page, err := uc.GetPage(context.Background(), "page-1")

	require.NoError(t, err)
	require.NotNil(t, page)
	assert.Equal(t, "page-1", page.ID)
	assert.Equal(t, "About Us", page.Title)
}

func TestGetPage_NotFound(t *testing.T) {
	_, pgRepo, _, pub := defaultCMSMocks()

	uc := NewPageUseCase(pgRepo, pub)
	_, err := uc.GetPage(context.Background(), "page-missing")

	require.Error(t, err)
}

// ===========================================================================
// GetPageBySlug tests
// ===========================================================================

func TestGetPageBySlug_Success(t *testing.T) {
	_, pgRepo, _, pub := defaultCMSMocks()

	pgRepo.getBySlugFn = func(_ context.Context, slug string) (*domain.Page, error) {
		return &domain.Page{ID: "page-1", Title: "About Us", Slug: slug}, nil
	}

	uc := NewPageUseCase(pgRepo, pub)
	page, err := uc.GetPageBySlug(context.Background(), "about-us")

	require.NoError(t, err)
	require.NotNil(t, page)
	assert.Equal(t, "about-us", page.Slug)
}

func TestGetPageBySlug_NotFound(t *testing.T) {
	_, pgRepo, _, pub := defaultCMSMocks()

	uc := NewPageUseCase(pgRepo, pub)
	_, err := uc.GetPageBySlug(context.Background(), "nonexistent-slug")

	require.Error(t, err)
}

// ===========================================================================
// ListPages tests
// ===========================================================================

func TestListPublishedPages_Success(t *testing.T) {
	_, pgRepo, _, pub := defaultCMSMocks()

	pgRepo.listPublishedFn = func(_ context.Context, _, _ int) ([]domain.Page, int64, error) {
		return []domain.Page{
			{ID: "page-1", Status: domain.PageStatusPublished},
		}, 1, nil
	}

	uc := NewPageUseCase(pgRepo, pub)
	pages, total, err := uc.ListPublishedPages(context.Background(), 1, 20)

	require.NoError(t, err)
	assert.Len(t, pages, 1)
	assert.Equal(t, int64(1), total)
}

func TestListPublishedPages_DefaultsPagination(t *testing.T) {
	_, pgRepo, _, pub := defaultCMSMocks()

	var capturedPage, capturedPageSize int
	pgRepo.listPublishedFn = func(_ context.Context, page, pageSize int) ([]domain.Page, int64, error) {
		capturedPage = page
		capturedPageSize = pageSize
		return nil, 0, nil
	}

	uc := NewPageUseCase(pgRepo, pub)

	_, _, _ = uc.ListPublishedPages(context.Background(), 0, 20)
	assert.Equal(t, 1, capturedPage)

	_, _, _ = uc.ListPublishedPages(context.Background(), 1, 0)
	assert.Equal(t, 20, capturedPageSize)
}

func TestListAllPages_Success(t *testing.T) {
	_, pgRepo, _, pub := defaultCMSMocks()

	pgRepo.listAllFn = func(_ context.Context, _, _ int) ([]domain.Page, int64, error) {
		return []domain.Page{
			{ID: "page-1", Status: domain.PageStatusDraft},
			{ID: "page-2", Status: domain.PageStatusPublished},
		}, 2, nil
	}

	uc := NewPageUseCase(pgRepo, pub)
	pages, total, err := uc.ListAllPages(context.Background(), 1, 20)

	require.NoError(t, err)
	assert.Len(t, pages, 2)
	assert.Equal(t, int64(2), total)
}

func TestListAllPages_DefaultsPagination(t *testing.T) {
	_, pgRepo, _, pub := defaultCMSMocks()

	var capturedPage, capturedPageSize int
	pgRepo.listAllFn = func(_ context.Context, page, pageSize int) ([]domain.Page, int64, error) {
		capturedPage = page
		capturedPageSize = pageSize
		return nil, 0, nil
	}

	uc := NewPageUseCase(pgRepo, pub)

	_, _, _ = uc.ListAllPages(context.Background(), -1, 20)
	assert.Equal(t, 1, capturedPage)

	_, _, _ = uc.ListAllPages(context.Background(), 1, 200)
	assert.Equal(t, 20, capturedPageSize)
}

// ===========================================================================
// UpdatePage tests
// ===========================================================================

func TestUpdatePage_Success(t *testing.T) {
	_, pgRepo, _, pub := defaultCMSMocks()

	pgRepo.getByIDFn = func(_ context.Context, id string) (*domain.Page, error) {
		return &domain.Page{
			ID:          id,
			Title:       "Old Title",
			Slug:        "old-title",
			ContentHTML: "<p>Old content</p>",
		}, nil
	}

	var updated *domain.Page
	pgRepo.updateFn = func(_ context.Context, pg *domain.Page) error {
		updated = pg
		return nil
	}

	uc := NewPageUseCase(pgRepo, pub)
	err := uc.UpdatePage(context.Background(), &domain.Page{
		ID:          "page-1",
		Title:       "New Title",
		ContentHTML: "<p>New content</p>",
	})

	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "new-title", updated.Slug) // slug regenerated
	assert.Contains(t, updated.ContentHTML, "New content")
}

func TestUpdatePage_KeepsSlugWhenTitleUnchanged(t *testing.T) {
	_, pgRepo, _, pub := defaultCMSMocks()

	pgRepo.getByIDFn = func(_ context.Context, id string) (*domain.Page, error) {
		return &domain.Page{
			ID:    id,
			Title: "Same Title",
			Slug:  "same-title",
		}, nil
	}

	var updated *domain.Page
	pgRepo.updateFn = func(_ context.Context, pg *domain.Page) error {
		updated = pg
		return nil
	}

	uc := NewPageUseCase(pgRepo, pub)
	err := uc.UpdatePage(context.Background(), &domain.Page{
		ID:    "page-1",
		Title: "Same Title", // same title
	})

	require.NoError(t, err)
	assert.Equal(t, "same-title", updated.Slug) // slug unchanged
}

func TestUpdatePage_KeepsSlugWhenTitleEmpty(t *testing.T) {
	_, pgRepo, _, pub := defaultCMSMocks()

	pgRepo.getByIDFn = func(_ context.Context, id string) (*domain.Page, error) {
		return &domain.Page{
			ID:    id,
			Title: "Existing Title",
			Slug:  "existing-title",
		}, nil
	}

	var updated *domain.Page
	pgRepo.updateFn = func(_ context.Context, pg *domain.Page) error {
		updated = pg
		return nil
	}

	uc := NewPageUseCase(pgRepo, pub)
	err := uc.UpdatePage(context.Background(), &domain.Page{
		ID:    "page-1",
		Title: "", // empty title => keep existing slug
	})

	require.NoError(t, err)
	assert.Equal(t, "existing-title", updated.Slug)
}

func TestUpdatePage_NotFound(t *testing.T) {
	_, pgRepo, _, pub := defaultCMSMocks()

	pgRepo.getByIDFn = func(_ context.Context, _ string) (*domain.Page, error) {
		return nil, errors.New("not found")
	}

	uc := NewPageUseCase(pgRepo, pub)
	err := uc.UpdatePage(context.Background(), &domain.Page{ID: "page-missing"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "page not found")
}

func TestUpdatePage_RepoError(t *testing.T) {
	_, pgRepo, _, pub := defaultCMSMocks()

	pgRepo.getByIDFn = func(_ context.Context, id string) (*domain.Page, error) {
		return &domain.Page{ID: id, Title: "Title", Slug: "title"}, nil
	}
	pgRepo.updateFn = func(_ context.Context, _ *domain.Page) error {
		return errors.New("db error")
	}

	uc := NewPageUseCase(pgRepo, pub)
	err := uc.UpdatePage(context.Background(), &domain.Page{ID: "page-1"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update page")
}

// ===========================================================================
// DeletePage tests
// ===========================================================================

func TestDeletePage_Success(t *testing.T) {
	_, pgRepo, _, pub := defaultCMSMocks()

	pgRepo.getByIDFn = func(_ context.Context, id string) (*domain.Page, error) {
		return &domain.Page{ID: id}, nil
	}

	var deletedID string
	pgRepo.deleteFn = func(_ context.Context, id string) error {
		deletedID = id
		return nil
	}

	uc := NewPageUseCase(pgRepo, pub)
	err := uc.DeletePage(context.Background(), "page-1")

	require.NoError(t, err)
	assert.Equal(t, "page-1", deletedID)
}

func TestDeletePage_NotFound(t *testing.T) {
	_, pgRepo, _, pub := defaultCMSMocks()

	pgRepo.getByIDFn = func(_ context.Context, _ string) (*domain.Page, error) {
		return nil, errors.New("not found")
	}

	uc := NewPageUseCase(pgRepo, pub)
	err := uc.DeletePage(context.Background(), "page-missing")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "page not found")
}

func TestDeletePage_RepoError(t *testing.T) {
	_, pgRepo, _, pub := defaultCMSMocks()

	pgRepo.getByIDFn = func(_ context.Context, id string) (*domain.Page, error) {
		return &domain.Page{ID: id}, nil
	}
	pgRepo.deleteFn = func(_ context.Context, _ string) error {
		return errors.New("db error")
	}

	uc := NewPageUseCase(pgRepo, pub)
	err := uc.DeletePage(context.Background(), "page-1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete page")
}

// ===========================================================================
// PublishPage tests
// ===========================================================================

func TestPublishPage_Success(t *testing.T) {
	_, pgRepo, _, pub := defaultCMSMocks()

	pgRepo.getByIDFn = func(_ context.Context, id string) (*domain.Page, error) {
		return &domain.Page{
			ID:     id,
			Title:  "About Us",
			Slug:   "about-us",
			Status: domain.PageStatusDraft,
		}, nil
	}

	var updated *domain.Page
	pgRepo.updateFn = func(_ context.Context, pg *domain.Page) error {
		updated = pg
		return nil
	}

	uc := NewPageUseCase(pgRepo, pub)
	page, err := uc.PublishPage(context.Background(), "page-1")

	require.NoError(t, err)
	require.NotNil(t, page)
	assert.Equal(t, domain.PageStatusPublished, page.Status)
	assert.NotNil(t, page.PublishedAt)
	assert.NotNil(t, updated)
}

func TestPublishPage_NotFound(t *testing.T) {
	_, pgRepo, _, pub := defaultCMSMocks()

	pgRepo.getByIDFn = func(_ context.Context, _ string) (*domain.Page, error) {
		return nil, errors.New("not found")
	}

	uc := NewPageUseCase(pgRepo, pub)
	_, err := uc.PublishPage(context.Background(), "page-missing")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "page not found")
}

func TestPublishPage_UpdateError(t *testing.T) {
	_, pgRepo, _, pub := defaultCMSMocks()

	pgRepo.getByIDFn = func(_ context.Context, id string) (*domain.Page, error) {
		return &domain.Page{ID: id, Status: domain.PageStatusDraft}, nil
	}
	pgRepo.updateFn = func(_ context.Context, _ *domain.Page) error {
		return errors.New("db error")
	}

	uc := NewPageUseCase(pgRepo, pub)
	_, err := uc.PublishPage(context.Background(), "page-1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to publish page")
}

// ===========================================================================
// generateSlug tests (CMS version)
// ===========================================================================

func TestCMSGenerateSlug_Basic(t *testing.T) {
	slug := generateSlug("About Us")
	assert.Equal(t, "about-us", slug)
}

func TestCMSGenerateSlug_SpecialChars(t *testing.T) {
	slug := generateSlug("Terms & Conditions!!!")
	assert.Equal(t, "terms-conditions", slug)
}

func TestCMSGenerateSlug_PreservesNumbers(t *testing.T) {
	slug := generateSlug("Top 10 Products")
	assert.Equal(t, "top-10-products", slug)
}
