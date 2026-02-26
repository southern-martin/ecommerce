package postgres

import (
	"time"

	"github.com/lib/pq"
	"github.com/southern-martin/ecommerce/services/review/internal/domain"
)

// ReviewModel is the GORM model for the reviews table.
type ReviewModel struct {
	ID                 string         `gorm:"type:uuid;primaryKey"`
	ProductID          string         `gorm:"type:uuid;index;not null"`
	UserID             string         `gorm:"type:uuid;index;not null"`
	UserName           string         `gorm:"type:varchar(255);not null"`
	Rating             int            `gorm:"not null;check:rating >= 1 AND rating <= 5"`
	Title              string         `gorm:"type:varchar(500)"`
	Content            string         `gorm:"type:text"`
	Pros               pq.StringArray `gorm:"type:text[]"`
	Cons               pq.StringArray `gorm:"type:text[]"`
	Images             pq.StringArray `gorm:"type:text[]"`
	IsVerifiedPurchase bool           `gorm:"default:false"`
	HelpfulCount       int            `gorm:"default:0"`
	UnhelpfulCount     int            `gorm:"default:0"`
	Status             string         `gorm:"type:varchar(20);default:'pending';index"`
	CreatedAt          time.Time      `gorm:"autoCreateTime"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime"`
}

// TableName returns the table name for ReviewModel.
func (ReviewModel) TableName() string { return "reviews" }

// ToDomain converts a ReviewModel to a domain Review.
func (m *ReviewModel) ToDomain() *domain.Review {
	return &domain.Review{
		ID:                 m.ID,
		ProductID:          m.ProductID,
		UserID:             m.UserID,
		UserName:           m.UserName,
		Rating:             m.Rating,
		Title:              m.Title,
		Content:            m.Content,
		Pros:               []string(m.Pros),
		Cons:               []string(m.Cons),
		Images:             []string(m.Images),
		IsVerifiedPurchase: m.IsVerifiedPurchase,
		HelpfulCount:       m.HelpfulCount,
		UnhelpfulCount:     m.UnhelpfulCount,
		Status:             domain.ReviewStatus(m.Status),
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

// ToReviewModel converts a domain Review to a ReviewModel.
func ToReviewModel(r *domain.Review) *ReviewModel {
	return &ReviewModel{
		ID:                 r.ID,
		ProductID:          r.ProductID,
		UserID:             r.UserID,
		UserName:           r.UserName,
		Rating:             r.Rating,
		Title:              r.Title,
		Content:            r.Content,
		Pros:               pq.StringArray(r.Pros),
		Cons:               pq.StringArray(r.Cons),
		Images:             pq.StringArray(r.Images),
		IsVerifiedPurchase: r.IsVerifiedPurchase,
		HelpfulCount:       r.HelpfulCount,
		UnhelpfulCount:     r.UnhelpfulCount,
		Status:             string(r.Status),
		CreatedAt:          r.CreatedAt,
		UpdatedAt:          r.UpdatedAt,
	}
}
