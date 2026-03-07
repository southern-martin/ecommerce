package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/ai/internal/domain"
)

// ---------------------------------------------------------------------------
// Hand-written function-field mocks
// ---------------------------------------------------------------------------

// --- RecommendationRepository mock ---

type mockRecommendationRepo struct {
	getByIDFn    func(ctx context.Context, id string) (*domain.Recommendation, error)
	listByUserFn func(ctx context.Context, userID string, page, pageSize int, filterViewed *bool) ([]domain.Recommendation, int64, error)
	createFn     func(ctx context.Context, recommendation *domain.Recommendation) error
	markViewedFn func(ctx context.Context, id string) error
}

func (m *mockRecommendationRepo) GetByID(ctx context.Context, id string) (*domain.Recommendation, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}
func (m *mockRecommendationRepo) ListByUser(ctx context.Context, userID string, page, pageSize int, filterViewed *bool) ([]domain.Recommendation, int64, error) {
	if m.listByUserFn != nil {
		return m.listByUserFn(ctx, userID, page, pageSize, filterViewed)
	}
	return nil, 0, nil
}
func (m *mockRecommendationRepo) Create(ctx context.Context, recommendation *domain.Recommendation) error {
	if m.createFn != nil {
		return m.createFn(ctx, recommendation)
	}
	return nil
}
func (m *mockRecommendationRepo) MarkViewed(ctx context.Context, id string) error {
	if m.markViewedFn != nil {
		return m.markViewedFn(ctx, id)
	}
	return nil
}

// ===========================================================================
// GetRecommendations tests
// ===========================================================================

func TestGetRecommendations_Success(t *testing.T) {
	recRepo := &mockRecommendationRepo{}
	pub := &mockAIEventPublisher{}

	expectedRecs := []domain.Recommendation{
		{ID: "rec-1", UserID: "user-1", ProductID: "prod-1", Score: 0.95},
		{ID: "rec-2", UserID: "user-1", ProductID: "prod-2", Score: 0.80},
	}
	recRepo.listByUserFn = func(_ context.Context, userID string, page, pageSize int, filterViewed *bool) ([]domain.Recommendation, int64, error) {
		assert.Equal(t, "user-1", userID)
		assert.Equal(t, 1, page)
		assert.Equal(t, 20, pageSize)
		assert.Nil(t, filterViewed)
		return expectedRecs, 2, nil
	}

	uc := NewRecommendationUseCase(recRepo, pub)
	recs, total, err := uc.GetRecommendations(context.Background(), "user-1", 1, 20)

	require.NoError(t, err)
	assert.Len(t, recs, 2)
	assert.Equal(t, int64(2), total)
}

func TestGetRecommendations_DefaultPagination(t *testing.T) {
	recRepo := &mockRecommendationRepo{}
	pub := &mockAIEventPublisher{}

	var capturedPage, capturedPageSize int
	recRepo.listByUserFn = func(_ context.Context, _ string, page, pageSize int, _ *bool) ([]domain.Recommendation, int64, error) {
		capturedPage = page
		capturedPageSize = pageSize
		return nil, 0, nil
	}

	uc := NewRecommendationUseCase(recRepo, pub)

	// page < 1 should default to 1
	_, _, _ = uc.GetRecommendations(context.Background(), "user-1", 0, 20)
	assert.Equal(t, 1, capturedPage)

	// pageSize > 100 should default to 20
	_, _, _ = uc.GetRecommendations(context.Background(), "user-1", 1, 200)
	assert.Equal(t, 20, capturedPageSize)

	// pageSize < 1 should default to 20
	_, _, _ = uc.GetRecommendations(context.Background(), "user-1", 1, -5)
	assert.Equal(t, 20, capturedPageSize)
}

func TestGetRecommendations_RepoError(t *testing.T) {
	recRepo := &mockRecommendationRepo{}
	pub := &mockAIEventPublisher{}

	recRepo.listByUserFn = func(_ context.Context, _ string, _, _ int, _ *bool) ([]domain.Recommendation, int64, error) {
		return nil, 0, errors.New("db error")
	}

	uc := NewRecommendationUseCase(recRepo, pub)
	_, _, err := uc.GetRecommendations(context.Background(), "user-1", 1, 20)

	require.Error(t, err)
}

// ===========================================================================
// GenerateRecommendations tests
// ===========================================================================

func TestGenerateRecommendations_Success(t *testing.T) {
	recRepo := &mockRecommendationRepo{}
	pub := &mockAIEventPublisher{}

	var createdCount int
	recRepo.createFn = func(_ context.Context, r *domain.Recommendation) error {
		createdCount++
		assert.Equal(t, "user-1", r.UserID)
		assert.NotEmpty(t, r.ProductID)
		assert.NotEmpty(t, r.Reason)
		assert.False(t, r.IsViewed)
		assert.Greater(t, r.Score, 0.0)
		return nil
	}

	uc := NewRecommendationUseCase(recRepo, pub)
	recs, err := uc.GenerateRecommendations(context.Background(), "user-1")

	require.NoError(t, err)
	assert.Len(t, recs, 5)
	assert.Equal(t, 5, createdCount)

	// All recommendations should have valid fields
	for _, rec := range recs {
		assert.NotEmpty(t, rec.ID)
		assert.Equal(t, "user-1", rec.UserID)
		assert.NotEmpty(t, rec.ProductID)
		assert.NotEmpty(t, rec.Reason)
	}
}

func TestGenerateRecommendations_CreateRepoError(t *testing.T) {
	recRepo := &mockRecommendationRepo{}
	pub := &mockAIEventPublisher{}

	recRepo.createFn = func(_ context.Context, _ *domain.Recommendation) error {
		return errors.New("db write failed")
	}

	uc := NewRecommendationUseCase(recRepo, pub)
	_, err := uc.GenerateRecommendations(context.Background(), "user-1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create recommendation")
}

func TestGenerateRecommendations_PublishEventFailDoesNotBreak(t *testing.T) {
	recRepo := &mockRecommendationRepo{}
	pub := &mockAIEventPublisher{}

	recRepo.createFn = func(_ context.Context, _ *domain.Recommendation) error { return nil }
	pub.publishFn = func(_ context.Context, _ string, _ interface{}) error {
		return errors.New("nats down")
	}

	uc := NewRecommendationUseCase(recRepo, pub)
	recs, err := uc.GenerateRecommendations(context.Background(), "user-1")

	require.NoError(t, err)
	assert.Len(t, recs, 5)
}

func TestGenerateRecommendations_ValidReasons(t *testing.T) {
	recRepo := &mockRecommendationRepo{}
	pub := &mockAIEventPublisher{}
	recRepo.createFn = func(_ context.Context, _ *domain.Recommendation) error { return nil }

	validReasons := map[string]bool{
		"similar_to_viewed":          true,
		"frequently_bought_together": true,
		"trending":                   true,
		"based_on_history":           true,
	}

	uc := NewRecommendationUseCase(recRepo, pub)
	recs, err := uc.GenerateRecommendations(context.Background(), "user-1")

	require.NoError(t, err)
	for _, rec := range recs {
		assert.True(t, validReasons[rec.Reason], "unexpected reason: %s", rec.Reason)
	}
}

// ===========================================================================
// MarkViewed tests
// ===========================================================================

func TestMarkViewed_Success(t *testing.T) {
	recRepo := &mockRecommendationRepo{}
	pub := &mockAIEventPublisher{}

	var viewedID string
	recRepo.markViewedFn = func(_ context.Context, id string) error {
		viewedID = id
		return nil
	}

	uc := NewRecommendationUseCase(recRepo, pub)
	err := uc.MarkViewed(context.Background(), "rec-1")

	require.NoError(t, err)
	assert.Equal(t, "rec-1", viewedID)
}

func TestMarkViewed_RepoError(t *testing.T) {
	recRepo := &mockRecommendationRepo{}
	pub := &mockAIEventPublisher{}

	recRepo.markViewedFn = func(_ context.Context, _ string) error {
		return errors.New("db error")
	}

	uc := NewRecommendationUseCase(recRepo, pub)
	err := uc.MarkViewed(context.Background(), "rec-1")

	require.Error(t, err)
}
