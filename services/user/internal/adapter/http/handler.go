package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	apperrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/pkg/middleware"
	"github.com/southern-martin/ecommerce/pkg/pagination"
	"github.com/southern-martin/ecommerce/pkg/validator"
	"github.com/southern-martin/ecommerce/services/user/internal/usecase"
)

// Handler holds all HTTP handlers for the user service.
type Handler struct {
	profile  *usecase.ProfileUseCase
	address  *usecase.AddressUseCase
	seller   *usecase.SellerUseCase
	follow   *usecase.FollowUseCase
	wishlist *usecase.WishlistUseCase
	db       *gorm.DB
}

// NewHandler creates a new Handler with the given use cases.
func NewHandler(
	profile *usecase.ProfileUseCase,
	address *usecase.AddressUseCase,
	seller *usecase.SellerUseCase,
	follow *usecase.FollowUseCase,
	wishlist *usecase.WishlistUseCase,
	db *gorm.DB,
) *Handler {
	return &Handler{
		profile:  profile,
		address:  address,
		seller:   seller,
		follow:   follow,
		wishlist: wishlist,
		db:       db,
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

// GetProfile godoc
// @Summary      Get current user profile
// @Tags         User Profile
// @Produce      json
// @Param        X-User-ID  header  string  true  "User ID"
// @Success      200  {object}  object{id=string,email=string,first_name=string,last_name=string,display_name=string,phone=string,avatar_url=string,role=string,created_at=string}
// @Failure      404  {object}  object{error=string}
// @Router       /users/me [get]
// @Security     BearerAuth
func (h *Handler) GetProfile(c *gin.Context) {
	userID := getUserID(c)
	profile, err := h.profile.GetProfile(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, profile)
}

// UpdateProfile godoc
// @Summary      Update current user profile
// @Tags         User Profile
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string                    true  "User ID"
// @Param        body       body    usecase.UpdateProfileInput true  "Profile fields to update"
// @Success      200  {object}  object{id=string,email=string,first_name=string,last_name=string,display_name=string,phone=string,avatar_url=string,role=string,created_at=string}
// @Failure      400  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Router       /users/me [patch]
// @Security     BearerAuth
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

// CreateAddress godoc
// @Summary      Create a new address
// @Tags         Addresses
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string                     true  "User ID"
// @Param        body       body    usecase.CreateAddressInput  true  "Address details"
// @Success      201  {object}  object{id=string,user_id=string,label=string,full_name=string,phone=string,street=string,city=string,state=string,postal_code=string,country=string,is_default=bool}
// @Failure      400  {object}  object{error=string}
// @Router       /users/me/addresses [post]
// @Security     BearerAuth
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

// ListAddresses godoc
// @Summary      List all addresses for current user
// @Tags         Addresses
// @Produce      json
// @Param        X-User-ID  header  string  true  "User ID"
// @Success      200  {object}  object{addresses=[]object{id=string,user_id=string,label=string,full_name=string,phone=string,street=string,city=string,country=string,is_default=bool}}
// @Router       /users/me/addresses [get]
// @Security     BearerAuth
func (h *Handler) ListAddresses(c *gin.Context) {
	userID := getUserID(c)

	addresses, err := h.address.ListAddresses(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"addresses": addresses})
}

// UpdateAddress godoc
// @Summary      Update an address
// @Tags         Addresses
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string                     true  "User ID"
// @Param        id         path    string                     true  "Address ID"
// @Param        body       body    usecase.UpdateAddressInput  true  "Address fields to update"
// @Success      200  {object}  object{id=string,user_id=string,label=string,full_name=string,phone=string,street=string,city=string,state=string,postal_code=string,country=string,is_default=bool}
// @Failure      400  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Router       /users/me/addresses/{id} [patch]
// @Security     BearerAuth
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

// DeleteAddress godoc
// @Summary      Delete an address
// @Tags         Addresses
// @Produce      json
// @Param        X-User-ID  header  string  true  "User ID"
// @Param        id         path    string  true  "Address ID"
// @Success      200  {object}  object{message=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Router       /users/me/addresses/{id} [delete]
// @Security     BearerAuth
func (h *Handler) DeleteAddress(c *gin.Context) {
	userID := getUserID(c)
	addrID := c.Param("id")

	if err := h.address.DeleteAddress(c.Request.Context(), userID, addrID); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "address deleted"})
}

// SetDefaultAddress godoc
// @Summary      Set an address as default
// @Tags         Addresses
// @Produce      json
// @Param        X-User-ID  header  string  true  "User ID"
// @Param        id         path    string  true  "Address ID"
// @Success      200  {object}  object{message=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Router       /users/me/addresses/{id}/default [patch]
// @Security     BearerAuth
func (h *Handler) SetDefaultAddress(c *gin.Context) {
	userID := getUserID(c)
	addrID := c.Param("id")

	if err := h.address.SetDefault(c.Request.Context(), userID, addrID); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "default address updated"})
}

// CreateSeller godoc
// @Summary      Create a seller profile
// @Tags         Sellers
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string                     true  "User ID"
// @Param        body       body    usecase.CreateSellerInput   true  "Seller profile details"
// @Success      201  {object}  object{id=string,user_id=string,store_name=string,description=string,logo_url=string,rating=number,total_sales=int,status=string,created_at=string}
// @Failure      400  {object}  object{error=string}
// @Failure      409  {object}  object{error=string}
// @Router       /sellers [post]
// @Security     BearerAuth
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

// GetSeller godoc
// @Summary      Get a seller profile by ID
// @Tags         Sellers
// @Produce      json
// @Param        id  path  string  true  "Seller ID"
// @Success      200  {object}  object{id=string,user_id=string,store_name=string,description=string,logo_url=string,rating=number,total_sales=int,status=string,created_at=string}
// @Failure      404  {object}  object{error=string}
// @Router       /sellers/{id} [get]
// @Security     BearerAuth
func (h *Handler) GetSeller(c *gin.Context) {
	sellerID := c.Param("id")

	seller, err := h.seller.GetSeller(c.Request.Context(), sellerID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, seller)
}

// UpdateSeller godoc
// @Summary      Update current user's seller profile
// @Tags         Sellers
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string                     true  "User ID"
// @Param        body       body    usecase.UpdateSellerInput   true  "Seller fields to update"
// @Success      200  {object}  object{id=string,user_id=string,store_name=string,description=string,logo_url=string,rating=number,total_sales=int,status=string,created_at=string}
// @Failure      400  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Router       /sellers/me [patch]
// @Security     BearerAuth
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

// ApproveSeller godoc
// @Summary      Approve a seller profile (admin only)
// @Tags         Admin Users
// @Produce      json
// @Param        id  path  string  true  "Seller ID"
// @Success      200  {object}  object{id=string,user_id=string,store_name=string,description=string,logo_url=string,rating=number,total_sales=int,status=string,created_at=string}
// @Failure      404  {object}  object{error=string}
// @Router       /admin/sellers/{id}/approve [post]
// @Security     BearerAuth
func (h *Handler) ApproveSeller(c *gin.Context) {
	sellerID := c.Param("id")

	seller, err := h.seller.ApproveSeller(c.Request.Context(), sellerID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, seller)
}

// FollowSeller godoc
// @Summary      Follow a seller
// @Tags         Following
// @Produce      json
// @Param        X-User-ID  header  string  true  "User ID"
// @Param        id         path    string  true  "Seller ID"
// @Success      200  {object}  object{message=string}
// @Failure      409  {object}  object{error=string}
// @Router       /users/{id}/follow [post]
// @Security     BearerAuth
func (h *Handler) FollowSeller(c *gin.Context) {
	userID := getUserID(c)
	sellerID := c.Param("id")

	if err := h.follow.Follow(c.Request.Context(), userID, sellerID); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "followed"})
}

// UnfollowSeller godoc
// @Summary      Unfollow a seller
// @Tags         Following
// @Produce      json
// @Param        X-User-ID  header  string  true  "User ID"
// @Param        id         path    string  true  "Seller ID"
// @Success      200  {object}  object{message=string}
// @Failure      404  {object}  object{error=string}
// @Router       /users/{id}/follow [delete]
// @Security     BearerAuth
func (h *Handler) UnfollowSeller(c *gin.Context) {
	userID := getUserID(c)
	sellerID := c.Param("id")

	if err := h.follow.Unfollow(c.Request.Context(), userID, sellerID); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "unfollowed"})
}

// ListFollowed godoc
// @Summary      List sellers the current user follows
// @Tags         Following
// @Produce      json
// @Param        X-User-ID  header  string  true   "User ID"
// @Param        page       query   int     false  "Page number"
// @Param        page_size  query   int     false  "Page size"
// @Success      200  {object}  object{data=[]object{id=string,user_id=string,store_name=string,description=string,rating=number,status=string,created_at=string},total=int,page=int,page_size=int}
// @Router       /users/me/following [get]
// @Security     BearerAuth
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

// GetFollowerCount godoc
// @Summary      Get follower count for a seller
// @Tags         Sellers
// @Produce      json
// @Param        id  path  string  true  "Seller ID"
// @Success      200  {object}  object{seller_id=string,follower_count=int}
// @Failure      404  {object}  object{error=string}
// @Router       /sellers/{id}/followers/count [get]
// @Security     BearerAuth
func (h *Handler) GetFollowerCount(c *gin.Context) {
	sellerID := c.Param("id")

	count, err := h.follow.GetFollowerCount(c.Request.Context(), sellerID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"seller_id": sellerID, "follower_count": count})
}

// GetWishlist godoc
// @Summary      Get current user's wishlist
// @Tags         Wishlist
// @Produce      json
// @Param        X-User-ID  header  string  true  "User ID"
// @Success      200  {object}  object{data=[]object{id=string,user_id=string,product_id=string,created_at=string}}
// @Router       /wishlist [get]
// @Security     BearerAuth
func (h *Handler) GetWishlist(c *gin.Context) {
	userID := getUserID(c)
	items, err := h.wishlist.ListItems(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// AddToWishlist godoc
// @Summary      Add a product to wishlist
// @Tags         Wishlist
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string                           true  "User ID"
// @Param        body       body    object{product_id=string}        true  "Product to add"
// @Success      201  {object}  object{message=string}
// @Failure      400  {object}  object{error=string}
// @Failure      409  {object}  object{error=string}
// @Router       /wishlist [post]
// @Security     BearerAuth
func (h *Handler) AddToWishlist(c *gin.Context) {
	userID := getUserID(c)

	var input struct {
		ProductID string `json:"product_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_id is required"})
		return
	}

	if err := h.wishlist.AddItem(c.Request.Context(), userID, input.ProductID); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "added to wishlist"})
}

// RemoveFromWishlist godoc
// @Summary      Remove a product from wishlist
// @Tags         Wishlist
// @Produce      json
// @Param        X-User-ID   header  string  true  "User ID"
// @Param        productId   path    string  true  "Product ID"
// @Success      200  {object}  object{message=string}
// @Failure      404  {object}  object{error=string}
// @Router       /wishlist/{productId} [delete]
// @Security     BearerAuth
func (h *Handler) RemoveFromWishlist(c *gin.Context) {
	userID := getUserID(c)
	productID := c.Param("productId")

	if err := h.wishlist.RemoveItem(c.Request.Context(), userID, productID); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "removed from wishlist"})
}
