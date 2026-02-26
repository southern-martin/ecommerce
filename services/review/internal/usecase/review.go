package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/southern-martin/ecommerce/services/review/internal/domain"
)

// ReviewUseCase handles review operations.
type ReviewUseCase struct {
	reviewRepo domain.ReviewRepository
	publisher  domain.EventPublisher
}

// NewReviewUseCase creates a new ReviewUseCase.
func NewReviewUseCase(reviewRepo domain.ReviewRepository, publisher domain.EventPublisher) *ReviewUseCase {
	return &ReviewUseCase{
		reviewRepo: reviewRepo,
		publisher:  publisher,
	}
}

// CreateReviewRequest is the input for creating a review.
type CreateReviewRequest struct {
	ProductID          string
	UserID             string
	UserName           string
	Rating             int
	Title              string
	Content            string
	Pros               []string
	Cons               []string
	Images             []string
	IsVerifiedPurchase bool
}

// CreateReview creates a new review.
func (uc *ReviewUseCase) CreateReview(ctx context.Context, req CreateReviewRequest) (*domain.Review, error) {
	if req.Rating < 1 || req.Rating > 5 {
		return nil, errors.New("rating must be between 1 and 5")
	}
	if req.ProductID == "" {
		return nil, errors.New("product_id is required")
	}
	if req.UserID == "" {
		return nil, errors.New("user_id is required")
	}

	review := &domain.Review{
		ID:                 uuid.New().String(),
		ProductID:          req.ProductID,
		UserID:             req.UserID,
		UserName:           req.UserName,
		Rating:             req.Rating,
		Title:              req.Title,
		Content:            req.Content,
		Pros:               req.Pros,
		Cons:               req.Cons,
		Images:             req.Images,
		IsVerifiedPurchase: req.IsVerifiedPurchase,
		Status:             domain.ReviewStatusPending,
	}

	if err := uc.reviewRepo.Create(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	// Publish event
	_ = uc.publisher.Publish(ctx, "review.created", map[string]interface{}{
		"review_id":  review.ID,
		"product_id": review.ProductID,
		"user_id":    review.UserID,
		"rating":     review.Rating,
	})

	return review, nil
}

// GetReview retrieves a review by ID.
func (uc *ReviewUseCase) GetReview(ctx context.Context, id string) (*domain.Review, error) {
	return uc.reviewRepo.GetByID(ctx, id)
}

// ListProductReviews lists reviews for a product with filtering.
func (uc *ReviewUseCase) ListProductReviews(ctx context.Context, productID string, filter domain.ReviewFilter) ([]domain.Review, int64, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 || filter.PageSize > 100 {
		filter.PageSize = 20
	}
	return uc.reviewRepo.ListByProduct(ctx, productID, filter)
}

// ListUserReviews lists reviews by a user with pagination.
func (uc *ReviewUseCase) ListUserReviews(ctx context.Context, userID string, page, pageSize int) ([]domain.Review, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return uc.reviewRepo.ListByUser(ctx, userID, page, pageSize)
}

// UpdateReviewRequest is the input for updating a review.
type UpdateReviewRequest struct {
	Rating  *int
	Title   *string
	Content *string
	Pros    []string
	Cons    []string
	Images  []string
}

// UpdateReview updates an existing review.
func (uc *ReviewUseCase) UpdateReview(ctx context.Context, id string, req UpdateReviewRequest) (*domain.Review, error) {
	review, err := uc.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("review not found: %w", err)
	}

	if req.Rating != nil {
		if *req.Rating < 1 || *req.Rating > 5 {
			return nil, errors.New("rating must be between 1 and 5")
		}
		review.Rating = *req.Rating
	}
	if req.Title != nil {
		review.Title = *req.Title
	}
	if req.Content != nil {
		review.Content = *req.Content
	}
	if req.Pros != nil {
		review.Pros = req.Pros
	}
	if req.Cons != nil {
		review.Cons = req.Cons
	}
	if req.Images != nil {
		review.Images = req.Images
	}

	if err := uc.reviewRepo.Update(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to update review: %w", err)
	}

	return review, nil
}

// DeleteReview deletes a review by ID.
func (uc *ReviewUseCase) DeleteReview(ctx context.Context, id string) error {
	return uc.reviewRepo.Delete(ctx, id)
}

// ApproveReview approves a review.
func (uc *ReviewUseCase) ApproveReview(ctx context.Context, id string) (*domain.Review, error) {
	review, err := uc.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("review not found: %w", err)
	}

	review.Status = domain.ReviewStatusApproved

	if err := uc.reviewRepo.Update(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to approve review: %w", err)
	}

	// Publish event
	_ = uc.publisher.Publish(ctx, "review.approved", map[string]interface{}{
		"review_id":  review.ID,
		"product_id": review.ProductID,
		"user_id":    review.UserID,
		"rating":     review.Rating,
	})

	return review, nil
}

// GetProductSummary returns the review summary for a product.
func (uc *ReviewUseCase) GetProductSummary(ctx context.Context, productID string) (*domain.ReviewSummary, error) {
	return uc.reviewRepo.GetSummary(ctx, productID)
}
