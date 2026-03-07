package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/review/internal/domain"
)

// ---------------------------------------------------------------------------
// Hand-written function-field mocks
// ---------------------------------------------------------------------------

// --- ReviewRepository mock ---

type mockReviewRepo struct {
	getByIDFn     func(ctx context.Context, id string) (*domain.Review, error)
	listByProductFn func(ctx context.Context, productID string, filter domain.ReviewFilter) ([]domain.Review, int64, error)
	listByUserFn  func(ctx context.Context, userID string, page, pageSize int) ([]domain.Review, int64, error)
	createFn      func(ctx context.Context, review *domain.Review) error
	updateFn      func(ctx context.Context, review *domain.Review) error
	deleteFn      func(ctx context.Context, id string) error
	getSummaryFn  func(ctx context.Context, productID string) (*domain.ReviewSummary, error)
}

func (m *mockReviewRepo) GetByID(ctx context.Context, id string) (*domain.Review, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}
func (m *mockReviewRepo) ListByProduct(ctx context.Context, productID string, filter domain.ReviewFilter) ([]domain.Review, int64, error) {
	if m.listByProductFn != nil {
		return m.listByProductFn(ctx, productID, filter)
	}
	return nil, 0, nil
}
func (m *mockReviewRepo) ListByUser(ctx context.Context, userID string, page, pageSize int) ([]domain.Review, int64, error) {
	if m.listByUserFn != nil {
		return m.listByUserFn(ctx, userID, page, pageSize)
	}
	return nil, 0, nil
}
func (m *mockReviewRepo) Create(ctx context.Context, review *domain.Review) error {
	if m.createFn != nil {
		return m.createFn(ctx, review)
	}
	return nil
}
func (m *mockReviewRepo) Update(ctx context.Context, review *domain.Review) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, review)
	}
	return nil
}
func (m *mockReviewRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}
func (m *mockReviewRepo) GetSummary(ctx context.Context, productID string) (*domain.ReviewSummary, error) {
	if m.getSummaryFn != nil {
		return m.getSummaryFn(ctx, productID)
	}
	return nil, errors.New("not found")
}

// --- EventPublisher mock ---

type mockEventPublisher struct {
	publishFn func(ctx context.Context, subject string, data interface{}) error
}

func (m *mockEventPublisher) Publish(ctx context.Context, subject string, data interface{}) error {
	if m.publishFn != nil {
		return m.publishFn(ctx, subject, data)
	}
	return nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func newReviewUseCase(repo *mockReviewRepo, pub *mockEventPublisher) *ReviewUseCase {
	return NewReviewUseCase(repo, pub)
}

func defaultReviewMocks() (*mockReviewRepo, *mockEventPublisher) {
	return &mockReviewRepo{}, &mockEventPublisher{}
}

// ===========================================================================
// CreateReview tests
// ===========================================================================

func TestCreateReview_Success(t *testing.T) {
	repo, pub := defaultReviewMocks()
	var saved *domain.Review
	repo.createFn = func(_ context.Context, r *domain.Review) error {
		saved = r
		return nil
	}

	uc := newReviewUseCase(repo, pub)
	review, err := uc.CreateReview(context.Background(), CreateReviewRequest{
		ProductID:          "prod-1",
		UserID:             "user-1",
		UserName:           "John",
		Rating:             5,
		Title:              "Great product",
		Content:            "Love it!",
		IsVerifiedPurchase: true,
	})

	require.NoError(t, err)
	require.NotNil(t, review)
	assert.NotEmpty(t, review.ID)
	assert.Equal(t, "prod-1", review.ProductID)
	assert.Equal(t, "user-1", review.UserID)
	assert.Equal(t, "John", review.UserName)
	assert.Equal(t, 5, review.Rating)
	assert.Equal(t, "Great product", review.Title)
	assert.Equal(t, "Love it!", review.Content)
	assert.True(t, review.IsVerifiedPurchase)
	assert.Equal(t, domain.ReviewStatusPending, review.Status)
	assert.NotNil(t, saved)
}

func TestCreateReview_WithProsConsImages(t *testing.T) {
	repo, pub := defaultReviewMocks()
	repo.createFn = func(_ context.Context, _ *domain.Review) error { return nil }

	uc := newReviewUseCase(repo, pub)
	review, err := uc.CreateReview(context.Background(), CreateReviewRequest{
		ProductID: "prod-1",
		UserID:    "user-1",
		Rating:    4,
		Pros:      []string{"durable", "affordable"},
		Cons:      []string{"heavy"},
		Images:    []string{"img1.jpg", "img2.jpg"},
	})

	require.NoError(t, err)
	assert.Equal(t, []string{"durable", "affordable"}, review.Pros)
	assert.Equal(t, []string{"heavy"}, review.Cons)
	assert.Equal(t, []string{"img1.jpg", "img2.jpg"}, review.Images)
}

func TestCreateReview_ValidationRatingTooLow(t *testing.T) {
	repo, pub := defaultReviewMocks()
	uc := newReviewUseCase(repo, pub)

	_, err := uc.CreateReview(context.Background(), CreateReviewRequest{
		ProductID: "prod-1",
		UserID:    "user-1",
		Rating:    0,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "rating must be between 1 and 5")
}

func TestCreateReview_ValidationRatingTooHigh(t *testing.T) {
	repo, pub := defaultReviewMocks()
	uc := newReviewUseCase(repo, pub)

	_, err := uc.CreateReview(context.Background(), CreateReviewRequest{
		ProductID: "prod-1",
		UserID:    "user-1",
		Rating:    6,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "rating must be between 1 and 5")
}

func TestCreateReview_ValidationMissingProductID(t *testing.T) {
	repo, pub := defaultReviewMocks()
	uc := newReviewUseCase(repo, pub)

	_, err := uc.CreateReview(context.Background(), CreateReviewRequest{
		ProductID: "",
		UserID:    "user-1",
		Rating:    3,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "product_id is required")
}

func TestCreateReview_ValidationMissingUserID(t *testing.T) {
	repo, pub := defaultReviewMocks()
	uc := newReviewUseCase(repo, pub)

	_, err := uc.CreateReview(context.Background(), CreateReviewRequest{
		ProductID: "prod-1",
		UserID:    "",
		Rating:    3,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "user_id is required")
}

func TestCreateReview_RepoError(t *testing.T) {
	repo, pub := defaultReviewMocks()
	repo.createFn = func(_ context.Context, _ *domain.Review) error {
		return errors.New("db connection failed")
	}

	uc := newReviewUseCase(repo, pub)
	_, err := uc.CreateReview(context.Background(), CreateReviewRequest{
		ProductID: "prod-1",
		UserID:    "user-1",
		Rating:    4,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create review")
}

func TestCreateReview_PublishEventErrorDoesNotFail(t *testing.T) {
	repo, pub := defaultReviewMocks()
	repo.createFn = func(_ context.Context, _ *domain.Review) error { return nil }
	pub.publishFn = func(_ context.Context, _ string, _ interface{}) error {
		return errors.New("event bus down")
	}

	uc := newReviewUseCase(repo, pub)
	review, err := uc.CreateReview(context.Background(), CreateReviewRequest{
		ProductID: "prod-1",
		UserID:    "user-1",
		Rating:    4,
	})

	// Event publish failure should not cause CreateReview to fail
	require.NoError(t, err)
	require.NotNil(t, review)
}

// ===========================================================================
// GetReview tests
// ===========================================================================

func TestGetReview_Success(t *testing.T) {
	repo, pub := defaultReviewMocks()
	expected := &domain.Review{
		ID:        "rev-1",
		ProductID: "prod-1",
		UserID:    "user-1",
		Rating:    5,
		Title:     "Great",
	}
	repo.getByIDFn = func(_ context.Context, id string) (*domain.Review, error) {
		assert.Equal(t, "rev-1", id)
		return expected, nil
	}

	uc := newReviewUseCase(repo, pub)
	review, err := uc.GetReview(context.Background(), "rev-1")

	require.NoError(t, err)
	assert.Equal(t, expected, review)
}

func TestGetReview_NotFound(t *testing.T) {
	repo, pub := defaultReviewMocks()
	// Default getByIDFn returns "not found" error

	uc := newReviewUseCase(repo, pub)
	_, err := uc.GetReview(context.Background(), "nonexistent")

	require.Error(t, err)
}

// ===========================================================================
// ListProductReviews tests
// ===========================================================================

func TestListProductReviews_Success(t *testing.T) {
	repo, pub := defaultReviewMocks()
	reviews := []domain.Review{
		{ID: "rev-1", ProductID: "prod-1", Rating: 5},
		{ID: "rev-2", ProductID: "prod-1", Rating: 4},
	}
	repo.listByProductFn = func(_ context.Context, productID string, filter domain.ReviewFilter) ([]domain.Review, int64, error) {
		assert.Equal(t, "prod-1", productID)
		return reviews, 2, nil
	}

	uc := newReviewUseCase(repo, pub)
	result, total, err := uc.ListProductReviews(context.Background(), "prod-1", domain.ReviewFilter{
		Page:     1,
		PageSize: 20,
	})

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
}

func TestListProductReviews_DefaultsPagination(t *testing.T) {
	repo, pub := defaultReviewMocks()
	var capturedFilter domain.ReviewFilter
	repo.listByProductFn = func(_ context.Context, _ string, filter domain.ReviewFilter) ([]domain.Review, int64, error) {
		capturedFilter = filter
		return nil, 0, nil
	}

	uc := newReviewUseCase(repo, pub)
	// Pass zero/invalid page/pageSize to trigger defaults
	_, _, err := uc.ListProductReviews(context.Background(), "prod-1", domain.ReviewFilter{
		Page:     0,
		PageSize: 0,
	})

	require.NoError(t, err)
	assert.Equal(t, 1, capturedFilter.Page)
	assert.Equal(t, 20, capturedFilter.PageSize)
}

func TestListProductReviews_CapsPageSize(t *testing.T) {
	repo, pub := defaultReviewMocks()
	var capturedFilter domain.ReviewFilter
	repo.listByProductFn = func(_ context.Context, _ string, filter domain.ReviewFilter) ([]domain.Review, int64, error) {
		capturedFilter = filter
		return nil, 0, nil
	}

	uc := newReviewUseCase(repo, pub)
	_, _, err := uc.ListProductReviews(context.Background(), "prod-1", domain.ReviewFilter{
		Page:     1,
		PageSize: 200, // exceeds max of 100
	})

	require.NoError(t, err)
	assert.Equal(t, 20, capturedFilter.PageSize) // should be reset to 20
}

// ===========================================================================
// ListUserReviews tests
// ===========================================================================

func TestListUserReviews_Success(t *testing.T) {
	repo, pub := defaultReviewMocks()
	reviews := []domain.Review{
		{ID: "rev-1", UserID: "user-1", Rating: 5},
	}
	repo.listByUserFn = func(_ context.Context, userID string, page, pageSize int) ([]domain.Review, int64, error) {
		assert.Equal(t, "user-1", userID)
		assert.Equal(t, 1, page)
		assert.Equal(t, 10, pageSize)
		return reviews, 1, nil
	}

	uc := newReviewUseCase(repo, pub)
	result, total, err := uc.ListUserReviews(context.Background(), "user-1", 1, 10)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
}

func TestListUserReviews_DefaultsPagination(t *testing.T) {
	repo, pub := defaultReviewMocks()
	var capturedPage, capturedPageSize int
	repo.listByUserFn = func(_ context.Context, _ string, page, pageSize int) ([]domain.Review, int64, error) {
		capturedPage = page
		capturedPageSize = pageSize
		return nil, 0, nil
	}

	uc := newReviewUseCase(repo, pub)
	_, _, err := uc.ListUserReviews(context.Background(), "user-1", 0, 0)

	require.NoError(t, err)
	assert.Equal(t, 1, capturedPage)
	assert.Equal(t, 20, capturedPageSize)
}

// ===========================================================================
// UpdateReview tests
// ===========================================================================

func TestUpdateReview_Success(t *testing.T) {
	repo, pub := defaultReviewMocks()
	existing := &domain.Review{
		ID:        "rev-1",
		ProductID: "prod-1",
		UserID:    "user-1",
		Rating:    3,
		Title:     "Okay product",
		Content:   "It's fine",
	}
	repo.getByIDFn = func(_ context.Context, _ string) (*domain.Review, error) {
		return existing, nil
	}
	var updated *domain.Review
	repo.updateFn = func(_ context.Context, r *domain.Review) error {
		updated = r
		return nil
	}

	uc := newReviewUseCase(repo, pub)
	newRating := 5
	newTitle := "Amazing product"
	review, err := uc.UpdateReview(context.Background(), "rev-1", UpdateReviewRequest{
		Rating: &newRating,
		Title:  &newTitle,
	})

	require.NoError(t, err)
	assert.Equal(t, 5, review.Rating)
	assert.Equal(t, "Amazing product", review.Title)
	assert.Equal(t, "It's fine", review.Content) // unchanged
	assert.NotNil(t, updated)
}

func TestUpdateReview_PartialUpdate(t *testing.T) {
	repo, pub := defaultReviewMocks()
	existing := &domain.Review{
		ID:      "rev-1",
		Rating:  4,
		Title:   "Good",
		Content: "Nice product",
	}
	repo.getByIDFn = func(_ context.Context, _ string) (*domain.Review, error) {
		return existing, nil
	}
	repo.updateFn = func(_ context.Context, _ *domain.Review) error { return nil }

	uc := newReviewUseCase(repo, pub)
	newContent := "Updated content"
	review, err := uc.UpdateReview(context.Background(), "rev-1", UpdateReviewRequest{
		Content: &newContent,
	})

	require.NoError(t, err)
	assert.Equal(t, 4, review.Rating)         // unchanged
	assert.Equal(t, "Good", review.Title)      // unchanged
	assert.Equal(t, "Updated content", review.Content)
}

func TestUpdateReview_InvalidRating(t *testing.T) {
	repo, pub := defaultReviewMocks()
	repo.getByIDFn = func(_ context.Context, _ string) (*domain.Review, error) {
		return &domain.Review{ID: "rev-1", Rating: 3}, nil
	}

	uc := newReviewUseCase(repo, pub)
	badRating := 0
	_, err := uc.UpdateReview(context.Background(), "rev-1", UpdateReviewRequest{
		Rating: &badRating,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "rating must be between 1 and 5")
}

func TestUpdateReview_NotFound(t *testing.T) {
	repo, pub := defaultReviewMocks()
	// Default getByIDFn returns "not found" error

	uc := newReviewUseCase(repo, pub)
	_, err := uc.UpdateReview(context.Background(), "nonexistent", UpdateReviewRequest{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "review not found")
}

func TestUpdateReview_UpdateProsConsImages(t *testing.T) {
	repo, pub := defaultReviewMocks()
	existing := &domain.Review{
		ID:    "rev-1",
		Pros:  []string{"old-pro"},
		Cons:  []string{"old-con"},
		Images: []string{"old-img.jpg"},
	}
	repo.getByIDFn = func(_ context.Context, _ string) (*domain.Review, error) {
		return existing, nil
	}
	repo.updateFn = func(_ context.Context, _ *domain.Review) error { return nil }

	uc := newReviewUseCase(repo, pub)
	review, err := uc.UpdateReview(context.Background(), "rev-1", UpdateReviewRequest{
		Pros:   []string{"new-pro1", "new-pro2"},
		Cons:   []string{"new-con"},
		Images: []string{"new-img1.jpg", "new-img2.jpg"},
	})

	require.NoError(t, err)
	assert.Equal(t, []string{"new-pro1", "new-pro2"}, review.Pros)
	assert.Equal(t, []string{"new-con"}, review.Cons)
	assert.Equal(t, []string{"new-img1.jpg", "new-img2.jpg"}, review.Images)
}

// ===========================================================================
// DeleteReview tests
// ===========================================================================

func TestDeleteReview_Success(t *testing.T) {
	repo, pub := defaultReviewMocks()
	var deletedID string
	repo.deleteFn = func(_ context.Context, id string) error {
		deletedID = id
		return nil
	}

	uc := newReviewUseCase(repo, pub)
	err := uc.DeleteReview(context.Background(), "rev-1")

	require.NoError(t, err)
	assert.Equal(t, "rev-1", deletedID)
}

func TestDeleteReview_RepoError(t *testing.T) {
	repo, pub := defaultReviewMocks()
	repo.deleteFn = func(_ context.Context, _ string) error {
		return errors.New("delete failed")
	}

	uc := newReviewUseCase(repo, pub)
	err := uc.DeleteReview(context.Background(), "rev-1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "delete failed")
}

// ===========================================================================
// GetProductSummary (GetProductRatingSummary) tests
// ===========================================================================

func TestGetProductSummary_Success(t *testing.T) {
	repo, pub := defaultReviewMocks()
	expected := &domain.ReviewSummary{
		ProductID:     "prod-1",
		AverageRating: 4.5,
		TotalReviews:  10,
		RatingDistribution: map[int]int{
			1: 0, 2: 1, 3: 1, 4: 3, 5: 5,
		},
	}
	repo.getSummaryFn = func(_ context.Context, productID string) (*domain.ReviewSummary, error) {
		assert.Equal(t, "prod-1", productID)
		return expected, nil
	}

	uc := newReviewUseCase(repo, pub)
	summary, err := uc.GetProductSummary(context.Background(), "prod-1")

	require.NoError(t, err)
	assert.Equal(t, "prod-1", summary.ProductID)
	assert.Equal(t, 4.5, summary.AverageRating)
	assert.Equal(t, 10, summary.TotalReviews)
	assert.Equal(t, 5, summary.RatingDistribution[5])
}

func TestGetProductSummary_NotFound(t *testing.T) {
	repo, pub := defaultReviewMocks()
	// Default getSummaryFn returns "not found" error

	uc := newReviewUseCase(repo, pub)
	_, err := uc.GetProductSummary(context.Background(), "nonexistent")

	require.Error(t, err)
}
