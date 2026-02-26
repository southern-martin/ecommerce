package domain

// UserRegisteredEvent is published when a new user registers.
type UserRegisteredEvent struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

// PasswordResetRequestedEvent is published when a user requests a password reset.
type PasswordResetRequestedEvent struct {
	UserID     string `json:"user_id"`
	Email      string `json:"email"`
	ResetToken string `json:"reset_token"`
}
