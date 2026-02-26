package domain

import "time"

// UserProfile represents a user's profile information.
type UserProfile struct {
	ID          string    `gorm:"type:uuid;primaryKey" json:"id"`
	Email       string    `gorm:"not null" json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DisplayName string    `json:"display_name"`
	Phone       string    `json:"phone"`
	AvatarURL   string    `json:"avatar_url"`
	Bio         string    `json:"bio"`
	Role        string    `gorm:"default:buyer" json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName returns the database table name for UserProfile.
func (UserProfile) TableName() string {
	return "user_profiles"
}

// Address represents a shipping/billing address for a user.
type Address struct {
	ID         string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID     string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Label      string    `json:"label"`
	FullName   string    `gorm:"not null" json:"full_name"`
	Phone      string    `gorm:"not null" json:"phone"`
	Street     string    `gorm:"not null" json:"street"`
	City       string    `gorm:"not null" json:"city"`
	State      string    `json:"state"`
	PostalCode string    `gorm:"not null" json:"postal_code"`
	Country    string    `gorm:"not null" json:"country"`
	IsDefault  bool      `gorm:"default:false" json:"is_default"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TableName returns the database table name for Address.
func (Address) TableName() string {
	return "addresses"
}

// SellerProfile represents a seller's store profile.
type SellerProfile struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID      string    `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	StoreName   string    `gorm:"not null" json:"store_name"`
	Description string    `json:"description"`
	LogoURL     string    `json:"logo_url"`
	Rating      float64   `gorm:"default:0" json:"rating"`
	TotalSales  int       `gorm:"default:0" json:"total_sales"`
	Status      string    `gorm:"default:pending" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName returns the database table name for SellerProfile.
func (SellerProfile) TableName() string {
	return "seller_profiles"
}

// UserFollow represents a follow relationship between a user and a seller.
type UserFollow struct {
	ID         string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FollowerID string    `gorm:"type:uuid;not null;index" json:"follower_id"`
	SellerID   string    `gorm:"type:uuid;not null;index" json:"seller_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// TableName returns the database table name for UserFollow.
func (UserFollow) TableName() string {
	return "user_follows"
}
