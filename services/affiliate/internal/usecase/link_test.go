package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/affiliate/internal/domain"
)

// ===========================================================================
// CreateLink tests
// ===========================================================================

func TestCreateLink_Success(t *testing.T) {
	_, lRepo, _, _, pub := defaultAffiliateMocks()

	var saved *domain.AffiliateLink
	lRepo.createFn = func(_ context.Context, link *domain.AffiliateLink) error {
		saved = link
		return nil
	}

	uc := NewLinkUseCase(lRepo, pub)
	link, err := uc.CreateLink(context.Background(), CreateLinkRequest{
		UserID:    "user-1",
		TargetURL: "https://example.com/product/123",
	})

	require.NoError(t, err)
	require.NotNil(t, link)
	assert.NotEmpty(t, link.ID)
	assert.Equal(t, "user-1", link.UserID)
	assert.Equal(t, "https://example.com/product/123", link.TargetURL)
	assert.Len(t, link.Code, codeLength)
	assert.NotNil(t, saved)
}

func TestCreateLink_RepoError(t *testing.T) {
	_, lRepo, _, _, pub := defaultAffiliateMocks()

	lRepo.createFn = func(_ context.Context, _ *domain.AffiliateLink) error {
		return errors.New("db error")
	}

	uc := NewLinkUseCase(lRepo, pub)
	_, err := uc.CreateLink(context.Background(), CreateLinkRequest{
		UserID:    "user-1",
		TargetURL: "https://example.com",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create affiliate link")
}

// ===========================================================================
// GetLink tests
// ===========================================================================

func TestGetLink_Success(t *testing.T) {
	_, lRepo, _, _, pub := defaultAffiliateMocks()

	lRepo.getByIDFn = func(_ context.Context, id string) (*domain.AffiliateLink, error) {
		return &domain.AffiliateLink{
			ID:     id,
			UserID: "user-1",
			Code:   "abc12345",
		}, nil
	}

	uc := NewLinkUseCase(lRepo, pub)
	link, err := uc.GetLink(context.Background(), "link-1")

	require.NoError(t, err)
	require.NotNil(t, link)
	assert.Equal(t, "link-1", link.ID)
}

func TestGetLink_NotFound(t *testing.T) {
	_, lRepo, _, _, pub := defaultAffiliateMocks()

	uc := NewLinkUseCase(lRepo, pub)
	_, err := uc.GetLink(context.Background(), "link-missing")

	require.Error(t, err)
}

// ===========================================================================
// TrackClick tests
// ===========================================================================

func TestTrackClick_Success(t *testing.T) {
	_, lRepo, _, _, pub := defaultAffiliateMocks()

	lRepo.getByCodeFn = func(_ context.Context, code string) (*domain.AffiliateLink, error) {
		return &domain.AffiliateLink{
			ID:         "link-1",
			UserID:     "user-1",
			Code:       code,
			ClickCount: 10,
		}, nil
	}

	var clickedID string
	lRepo.incrementClicksFn = func(_ context.Context, id string) error {
		clickedID = id
		return nil
	}

	uc := NewLinkUseCase(lRepo, pub)
	link, err := uc.TrackClick(context.Background(), "abc12345")

	require.NoError(t, err)
	require.NotNil(t, link)
	assert.Equal(t, "link-1", clickedID)
	assert.Equal(t, int64(11), link.ClickCount) // incremented in-memory
}

func TestTrackClick_LinkNotFound(t *testing.T) {
	_, lRepo, _, _, pub := defaultAffiliateMocks()

	lRepo.getByCodeFn = func(_ context.Context, _ string) (*domain.AffiliateLink, error) {
		return nil, errors.New("not found")
	}

	uc := NewLinkUseCase(lRepo, pub)
	_, err := uc.TrackClick(context.Background(), "badcode")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "affiliate link not found")
}

func TestTrackClick_IncrementError(t *testing.T) {
	_, lRepo, _, _, pub := defaultAffiliateMocks()

	lRepo.getByCodeFn = func(_ context.Context, _ string) (*domain.AffiliateLink, error) {
		return &domain.AffiliateLink{ID: "link-1", Code: "abc12345"}, nil
	}
	lRepo.incrementClicksFn = func(_ context.Context, _ string) error {
		return errors.New("db error")
	}

	uc := NewLinkUseCase(lRepo, pub)
	_, err := uc.TrackClick(context.Background(), "abc12345")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to increment clicks")
}

// ===========================================================================
// ListUserLinks tests
// ===========================================================================

func TestListUserLinks_Success(t *testing.T) {
	_, lRepo, _, _, pub := defaultAffiliateMocks()

	lRepo.listByUserFn = func(_ context.Context, userID string, page, pageSize int) ([]domain.AffiliateLink, int64, error) {
		return []domain.AffiliateLink{
			{ID: "link-1", UserID: userID},
			{ID: "link-2", UserID: userID},
		}, 2, nil
	}

	uc := NewLinkUseCase(lRepo, pub)
	links, total, err := uc.ListUserLinks(context.Background(), "user-1", 1, 20)

	require.NoError(t, err)
	assert.Len(t, links, 2)
	assert.Equal(t, int64(2), total)
}

func TestListUserLinks_DefaultsPagination(t *testing.T) {
	_, lRepo, _, _, pub := defaultAffiliateMocks()

	var capturedPage, capturedPageSize int
	lRepo.listByUserFn = func(_ context.Context, _ string, page, pageSize int) ([]domain.AffiliateLink, int64, error) {
		capturedPage = page
		capturedPageSize = pageSize
		return nil, 0, nil
	}

	uc := NewLinkUseCase(lRepo, pub)

	// page < 1 defaults to 1
	_, _, _ = uc.ListUserLinks(context.Background(), "user-1", 0, 20)
	assert.Equal(t, 1, capturedPage)

	// pageSize < 1 defaults to 20
	_, _, _ = uc.ListUserLinks(context.Background(), "user-1", 1, 0)
	assert.Equal(t, 20, capturedPageSize)

	// pageSize > 100 defaults to 20
	_, _, _ = uc.ListUserLinks(context.Background(), "user-1", 1, 200)
	assert.Equal(t, 20, capturedPageSize)
}
