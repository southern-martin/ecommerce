package domain

import "time"

// AuthUser represents a user in the authentication system.
type AuthUser struct {
	ID              string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email           string     `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash    string     `json:"-"`
	Role            string     `gorm:"default:buyer" json:"role"` // buyer, seller, admin
	OAuthProvider   string     `json:"oauth_provider,omitempty"`
	OAuthProviderID string     `json:"oauth_provider_id,omitempty"`
	RefreshToken    string     `json:"-"`
	ResetToken      string     `json:"-"`
	ResetTokenExp   *time.Time `json:"-"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// TableName overrides the default GORM table name.
func (AuthUser) TableName() string {
	return "auth_users"
}
