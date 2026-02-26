package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	pkgerrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/auth/internal/usecase"
)

// Handler holds all use cases and provides HTTP handler methods.
type Handler struct {
	register       *usecase.RegisterUseCase
	login          *usecase.LoginUseCase
	refreshToken   *usecase.RefreshTokenUseCase
	logout         *usecase.LogoutUseCase
	forgotPassword *usecase.ForgotPasswordUseCase
	resetPassword  *usecase.ResetPasswordUseCase
	oauthLogin     *usecase.OAuthLoginUseCase
}

// NewHandler creates a new Handler with all use cases.
func NewHandler(
	register *usecase.RegisterUseCase,
	login *usecase.LoginUseCase,
	refreshToken *usecase.RefreshTokenUseCase,
	logout *usecase.LogoutUseCase,
	forgotPassword *usecase.ForgotPasswordUseCase,
	resetPassword *usecase.ResetPasswordUseCase,
	oauthLogin *usecase.OAuthLoginUseCase,
) *Handler {
	return &Handler{
		register:       register,
		login:          login,
		refreshToken:   refreshToken,
		logout:         logout,
		forgotPassword: forgotPassword,
		resetPassword:  resetPassword,
		oauthLogin:     oauthLogin,
	}
}

// Register handles POST /api/v1/auth/register
func (h *Handler) Register(c *gin.Context) {
	var input usecase.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.register.Execute(c.Request.Context(), input)
	if err != nil {
		status := pkgerrors.ToHTTPStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, output)
}

// Login handles POST /api/v1/auth/login
func (h *Handler) Login(c *gin.Context) {
	var input usecase.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.login.Execute(c.Request.Context(), input)
	if err != nil {
		status := pkgerrors.ToHTTPStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}

// RefreshToken handles POST /api/v1/auth/refresh
func (h *Handler) RefreshToken(c *gin.Context) {
	var input usecase.RefreshTokenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.refreshToken.Execute(c.Request.Context(), input)
	if err != nil {
		status := pkgerrors.ToHTTPStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}

// Logout handles POST /api/v1/auth/logout
func (h *Handler) Logout(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-ID header is required"})
		return
	}

	var input usecase.LogoutInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.UserID = userID

	if err := h.logout.Execute(c.Request.Context(), input); err != nil {
		status := pkgerrors.ToHTTPStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// ForgotPassword handles POST /api/v1/auth/forgot-password
func (h *Handler) ForgotPassword(c *gin.Context) {
	var input usecase.ForgotPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.forgotPassword.Execute(c.Request.Context(), input); err != nil {
		status := pkgerrors.ToHTTPStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Always return success to avoid leaking whether the email exists
	c.JSON(http.StatusOK, gin.H{"message": "if that email exists, a reset link has been sent"})
}

// ResetPassword handles POST /api/v1/auth/reset-password
func (h *Handler) ResetPassword(c *gin.Context) {
	var input usecase.ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.resetPassword.Execute(c.Request.Context(), input); err != nil {
		status := pkgerrors.ToHTTPStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}

// OAuthLogin handles POST /api/v1/auth/oauth/:provider
func (h *Handler) OAuthLogin(c *gin.Context) {
	provider := c.Param("provider")

	var input usecase.OAuthLoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.Provider = provider

	output, err := h.oauthLogin.Execute(c.Request.Context(), input)
	if err != nil {
		status := pkgerrors.ToHTTPStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}

// Health handles GET /health
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
