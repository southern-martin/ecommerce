package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

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
	db             *gorm.DB
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
	db *gorm.DB,
) *Handler {
	return &Handler{
		register:       register,
		login:          login,
		refreshToken:   refreshToken,
		logout:         logout,
		forgotPassword: forgotPassword,
		resetPassword:  resetPassword,
		oauthLogin:     oauthLogin,
		db:             db,
	}
}

// Register godoc
// @Summary      Register a new user
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        body  body      usecase.RegisterInput   true  "Registration details"
// @Success      201   {object}  usecase.RegisterOutput
// @Failure      400   {object}  object{error=string}
// @Failure      409   {object}  object{error=string}
// @Router       /auth/register [post]
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

// Login godoc
// @Summary      Log in with email and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        body  body      usecase.LoginInput   true  "Login credentials"
// @Success      200   {object}  usecase.LoginOutput
// @Failure      400   {object}  object{error=string}
// @Failure      401   {object}  object{error=string}
// @Router       /auth/login [post]
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

// RefreshToken godoc
// @Summary      Refresh an access token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        body  body      usecase.RefreshTokenInput   true  "Refresh token"
// @Success      200   {object}  usecase.RefreshTokenOutput
// @Failure      400   {object}  object{error=string}
// @Failure      401   {object}  object{error=string}
// @Router       /auth/refresh [post]
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

// Logout godoc
// @Summary      Log out the current user
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        X-User-ID  header    string              true  "Authenticated user ID"
// @Param        body       body      usecase.LogoutInput  true  "Access token to blacklist"
// @Success      200        {object}  object{message=string}
// @Failure      400        {object}  object{error=string}
// @Router       /auth/logout [post]
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

// ForgotPassword godoc
// @Summary      Request a password reset email
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        body  body      usecase.ForgotPasswordInput  true  "Email address"
// @Success      200   {object}  object{message=string}
// @Failure      400   {object}  object{error=string}
// @Router       /auth/forgot-password [post]
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

// ResetPassword godoc
// @Summary      Reset password with a reset token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        body  body      usecase.ResetPasswordInput  true  "Reset token and new password"
// @Success      200   {object}  object{message=string}
// @Failure      400   {object}  object{error=string}
// @Router       /auth/reset-password [post]
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

// OAuthLogin godoc
// @Summary      Authenticate via OAuth provider
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        provider  path      string                  true  "OAuth provider (e.g. google, github)"
// @Param        body      body      usecase.OAuthLoginInput  true  "OAuth token payload"
// @Success      200       {object}  usecase.OAuthLoginOutput
// @Failure      400       {object}  object{error=string}
// @Failure      401       {object}  object{error=string}
// @Router       /auth/oauth/{provider} [post]
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

// Ready handles GET /ready — deep health check including database connectivity.
func (h *Handler) Ready(c *gin.Context) {
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "error": "db connection lost"})
		return
	}
	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "error": "db ping failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
