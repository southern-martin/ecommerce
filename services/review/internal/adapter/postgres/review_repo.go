package postgres

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	"github.com/southern-martin/ecommerce/services/review/internal/domain"
	"gorm.io/gorm"
)

// ReviewRepo implements domain.ReviewRepository.
type ReviewRepo struct {
	db *gorm.DB
}

// NewReviewRepo creates a new ReviewRepo.
func NewReviewRepo(db *gorm.DB) *ReviewRepo {
	return &ReviewRepo{db: db}
}

func (r *ReviewRepo) GetByID(ctx context.Context, id string) (*domain.Review, error) {
	var model ReviewModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *ReviewRepo) ListByProduct(ctx context.Context, productID string, filter domain.ReviewFilter) ([]domain.Review, int64, error) {
	query := r.db.WithContext(ctx).Model(&ReviewModel{}).Where("product_id = ?", productID)

	if filter.MinRating > 0 {
		query = query.Where("rating >= ?", filter.MinRating)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", string(filter.Status))
	}

	var total int64
	query.Count(&total)

	// Sorting
	sortBy := "created_at"
	if filter.SortBy != "" {
		allowedSorts := map[string]bool{
			"created_at":    true,
			"rating":        true,
			"helpful_count": true,
		}
		if allowedSorts[filter.SortBy] {
			sortBy = filter.SortBy
		}
	}
	sortOrder := "DESC"
	if filter.SortOrder == "asc" || filter.SortOrder == "ASC" {
		sortOrder = "ASC"
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))

	// Pagination
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var models []ReviewModel
	if err := query.Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	reviews := make([]domain.Review, len(models))
	for i, m := range models {
		reviews[i] = *m.ToDomain()
	}
	return reviews, total, nil
}

func (r *ReviewRepo) ListByUser(ctx context.Context, userID string, page, pageSize int) ([]domain.Review, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&ReviewModel{}).Where("user_id = ?", userID).Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var models []ReviewModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	reviews := make([]domain.Review, len(models))
	for i, m := range models {
		reviews[i] = *m.ToDomain()
	}
	return reviews, total, nil
}

func (r *ReviewRepo) Create(ctx context.Context, review *domain.Review) error {
	model := ToReviewModel(review)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *ReviewRepo) Update(ctx context.Context, review *domain.Review) error {
	return r.db.WithContext(ctx).Model(&ReviewModel{}).Where("id = ?", review.ID).Updates(map[string]interface{}{
		"rating":              review.Rating,
		"title":               review.Title,
		"content":             review.Content,
		"pros":                pq.StringArray(review.Pros),
		"cons":                pq.StringArray(review.Cons),
		"images":              pq.StringArray(review.Images),
		"is_verified_purchase": review.IsVerifiedPurchase,
		"helpful_count":       review.HelpfulCount,
		"unhelpful_count":     review.UnhelpfulCount,
		"status":              string(review.Status),
	}).Error
}

func (r *ReviewRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&ReviewModel{}).Error
}

func (r *ReviewRepo) GetSummary(ctx context.Context, productID string) (*domain.ReviewSummary, error) {
	summary := &domain.ReviewSummary{
		ProductID:          productID,
		RatingDistribution: make(map[int]int),
	}

	// Get average rating and total reviews
	var result struct {
		AvgRating    float64
		TotalReviews int64
	}
	if err := r.db.WithContext(ctx).Model(&ReviewModel{}).
		Where("product_id = ? AND status = ?", productID, string(domain.ReviewStatusApproved)).
		Select("COALESCE(AVG(rating), 0) as avg_rating, COUNT(*) as total_reviews").
		Scan(&result).Error; err != nil {
		return nil, err
	}
	summary.AverageRating = result.AvgRating
	summary.TotalReviews = int(result.TotalReviews)

	// Get rating distribution
	type RatingCount struct {
		Rating int
		Count  int64
	}
	var ratingCounts []RatingCount
	if err := r.db.WithContext(ctx).Model(&ReviewModel{}).
		Where("product_id = ? AND status = ?", productID, string(domain.ReviewStatusApproved)).
		Select("rating, COUNT(*) as count").
		Group("rating").
		Scan(&ratingCounts).Error; err != nil {
		return nil, err
	}

	for _, rc := range ratingCounts {
		summary.RatingDistribution[rc.Rating] = int(rc.Count)
	}

	return summary, nil
}
