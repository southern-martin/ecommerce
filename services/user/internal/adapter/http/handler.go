package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/pkg/middleware"
	"github.com/southern-martin/ecommerce/pkg/pagination"
	"github.com/southern-martin/ecommerce/pkg/validator"
	"github.com/southern-martin/ecommerce/services/user/internal/usecase"
)

// Handler holds all HTTP handlers for the user service.
type Handler struct {
	profile *usecase.ProfileUseCase
	address *usecase.AddressUseCase
	seller  *usecase.SellerUseCase
	follow  *usecase.FollowUseCase
}

// NewHandler creates a new Handler with the given use cases.
func NewHandler(
	profile *usecase.ProfileUseCase,
	address *usecase.AddressUseCase,
	seller *usecase.SellerUseCase,
	follow *usecase.FollowUseCase,
) *Handler {
	return &Handler{
		profile: profile,
		address: address,
		seller:  seller,
		follow:  follow,
	}
}

// getUserID extracts the user ID from the Gin context.
func getUserID(c *gin.Context) string {
	id, _ := c.Get(middleware.ContextKeyUserID)
	if s, ok := id.(string); ok {
		return s
	}
	return ""
}

// handleError maps domain errors to HTTP responses.
func handleError(c *gin.Context, err error) {
	status := apperrors.ToHTTPStatus(err)
	c.JSON(status, gin.H{"error": err.Error()})
}

// Health returns a simple health check response.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetProfile handles GET /api/v1/users/me
func (h *Handler) GetProfile(c *gin.Context) {
	userID := getUserID(c)
	profile, err := h.profile.GetProfile(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, profile)
}

// UpdateProfile handles PATCH /api/v1/users/me
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID := getUserID(c)

	var input usecase.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	profile, err := h.profile.UpdateProfile(c.Request.Context(), userID, input)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, profile)
}

// CreateAddress handles POST /api/v1/users/me/addresses
func (h *Handler) CreateAddress(c *gin.Context) {
	userID := getUserID(c)

	var input usecase.CreateAddressInput
	if err := validator.BindAndValidate(c, &input); err != nil {
		handleError(c, err)
		return
	}

	addr, err := h.address.CreateAddress(c.Request.Context(), userID, input)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, addr)
}

// ListAddresses handles GET /api/v1/users/me/addresses
func (h *Handler) ListAddresses(c *gin.Context) {
	userID := getUserID(c)

	addresses, err := h.address.ListAddresses(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"addresses": addresses})
}

// UpdateAddress handles PATCH /api/v1/users/me/addresses/:id
func (h *Handler) UpdateAddress(c *gin.Context) {
	userID := getUserID(c)
	addrID := c.Param("id")

	var input usecase.UpdateAddressInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	addr, err := h.address.UpdateAddress(c.Request.Context(), userID, addrID, input)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, addr)
}

// DeleteAddress handles DELETE /api/v1/users/me/addresses/:id
func (h *Handler) DeleteAddress(c *gin.Context) {
	userID := getUserID(c)
	addrID := c.Param("id")

	if err := h.address.DeleteAddress(c.Request.Context(), userID, addrID); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "address deleted"})
}

// SetDefaultAddress handles PATCH /api/v1/users/me/addresses/:id/default
func (h *Handler) SetDefaultAddress(c *gin.Context) {
	userID := getUserID(c)
	addrID := c.Param("id")

	if err := h.address.SetDefault(c.Request.Context(), userID, addrID); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "default address updated"})
}

// CreateSeller handles POST /api/v1/sellers
func (h *Handler) CreateSeller(c *gin.Context) {
	userID := getUserID(c)

	var input usecase.CreateSellerInput
	if err := validator.BindAndValidate(c, &input); err != nil {
		handleError(c, err)
		return
	}

	seller, err := h.seller.CreateSeller(c.Request.Context(), userID, input)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, seller)
}

// GetSeller handles GET /api/v1/sellers/:id
func (h *Handler) GetSeller(c *gin.Context) {
	sellerID := c.Param("id")

	seller, err := h.seller.GetSeller(c.Request.Context(), sellerID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, seller)
}

// UpdateSeller handles PATCH /api/v1/sellers/me
func (h *Handler) UpdateSeller(c *gin.Context) {
	userID := getUserID(c)

	var input usecase.UpdateSellerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	seller, err := h.seller.UpdateSeller(c.Request.Context(), userID, input)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, seller)
}

// ApproveSeller handles POST /api/v1/admin/sellers/:id/approve
func (h *Handler) ApproveSeller(c *gin.Context) {
	sellerID := c.Param("id")

	seller, err := h.seller.ApproveSeller(c.Request.Context(), sellerID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, seller)
}

// FollowSeller handles POST /api/v1/users/:id/follow
func (h *Handler) FollowSeller(c *gin.Context) {
	userID := getUserID(c)
	sellerID := c.Param("id")

	if err := h.follow.Follow(c.Request.Context(), userID, sellerID); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "followed"})
}

// UnfollowSeller handles DELETE /api/v1/users/:id/follow
func (h *Handler) UnfollowSeller(c *gin.Context) {
	userID := getUserID(c)
	sellerID := c.Param("id")

	if err := h.follow.Unfollow(c.Request.Context(), userID, sellerID); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "unfollowed"})
}

// ListFollowed handles GET /api/v1/users/me/following
func (h *Handler) ListFollowed(c *gin.Context) {
	userID := getUserID(c)
	params := pagination.ParseFromGin(c)

	sellers, total, err := h.follow.ListFollowed(c.Request.Context(), userID, params.Page, params.PageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	result := pagination.NewPaginatedResult(sellers, total, params)
	c.JSON(http.StatusOK, result)
}

// GetFollowerCount handles GET /api/v1/sellers/:id/followers/count
func (h *Handler) GetFollowerCount(c *gin.Context) {
	sellerID := c.Param("id")

	count, err := h.follow.GetFollowerCount(c.Request.Context(), sellerID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"seller_id": sellerID, "follower_count": count})
}
