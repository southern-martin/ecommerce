package domain

import "time"

// ReviewStatus represents the status of a review.
type ReviewStatus string

const (
	ReviewStatusPending  ReviewStatus = "pending"
	ReviewStatusApproved ReviewStatus = "approved"
	ReviewStatusRejected ReviewStatus = "rejected"
)

// Review represents a product review.
type Review struct {
	ID                 string
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
	HelpfulCount       int
	UnhelpfulCount     int
	Status             ReviewStatus
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// ReviewSummary holds aggregated review data for a product.
type ReviewSummary struct {
	ProductID          string
	AverageRating      float64
	TotalReviews       int
	RatingDistribution map[int]int
}

// ReviewFilter holds filtering and pagination options for listing reviews.
type ReviewFilter struct {
	ProductID string
	UserID    string
	MinRating int
	Status    ReviewStatus
	SortBy    string
	SortOrder string
	Page      int
	PageSize  int
}
