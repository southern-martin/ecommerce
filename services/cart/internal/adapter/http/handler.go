package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/southern-martin/ecommerce/services/cart/internal/domain"
	"github.com/southern-martin/ecommerce/services/cart/internal/usecase"
)

// CartHandler handles HTTP requests for cart operations.
type CartHandler struct {
	cartUC *usecase.CartUseCase
	logger zerolog.Logger
	db     *gorm.DB
}

// NewCartHandler creates a new CartHandler.
func NewCartHandler(cartUC *usecase.CartUseCase, logger zerolog.Logger, db *gorm.DB) *CartHandler {
	return &CartHandler{
		cartUC: cartUC,
		logger: logger.With().Str("component", "cart_handler").Logger(),
		db:     db,
	}
}

// cartResponse is the standard response for cart endpoints.
type cartResponse struct {
	UserID       string            `json:"user_id"`
	Items        []domain.CartItem `json:"items"`
	TotalItems   int               `json:"total_items"`
	SubtotalCents int64            `json:"subtotal_cents"`
}

func toCartResponse(cart *domain.Cart) cartResponse {
	items := cart.Items
	if items == nil {
		items = []domain.CartItem{}
	}
	return cartResponse{
		UserID:        cart.UserID,
		Items:         items,
		TotalItems:    cart.TotalItems(),
		SubtotalCents: cart.SubtotalCents(),
	}
}

// getUserID extracts the user ID from the X-User-ID header.
func getUserID(c *gin.Context) string {
	return c.GetHeader("X-User-ID")
}

// GetCart godoc
// @Summary      Get current user's cart
// @Tags         Cart
// @Produce      json
// @Param        X-User-ID  header  string  true  "User ID"
// @Success      200  {object}  cartResponse
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /cart [get]
// @Security     BearerAuth
func (h *CartHandler) GetCart(c *gin.Context) {
	userID := getUserID(c)
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-ID header is required"})
		return
	}

	cart, err := h.cartUC.GetCart(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", userID).Msg("failed to get cart")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cart"})
		return
	}

	c.JSON(http.StatusOK, toCartResponse(cart))
}

// addItemRequest is the request body for adding an item to the cart.
type addItemRequest struct {
	ProductID   string `json:"product_id" binding:"required"`
	VariantID   string `json:"variant_id"`
	ProductName string `json:"product_name"`
	VariantName string `json:"variant_name"`
	SKU         string `json:"sku"`
	PriceCents  int64  `json:"price_cents" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required,min=1"`
	ImageURL    string `json:"image_url"`
	SellerID    string `json:"seller_id"`
}

// AddItem godoc
// @Summary      Add an item to the cart
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string          true  "User ID"
// @Param        body       body    addItemRequest  true  "Item to add"
// @Success      200  {object}  cartResponse
// @Failure      400  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /cart/items [post]
// @Security     BearerAuth
func (h *CartHandler) AddItem(c *gin.Context) {
	userID := getUserID(c)
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-ID header is required"})
		return
	}

	var req addItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := domain.CartItem{
		ProductID:   req.ProductID,
		VariantID:   req.VariantID,
		ProductName: req.ProductName,
		VariantName: req.VariantName,
		SKU:         req.SKU,
		PriceCents:  req.PriceCents,
		Quantity:    req.Quantity,
		ImageURL:    req.ImageURL,
		SellerID:    req.SellerID,
	}

	cart, err := h.cartUC.AddItem(c.Request.Context(), userID, item)
	if err != nil {
		h.handleUseCaseError(c, err, "failed to add item")
		return
	}

	c.JSON(http.StatusOK, toCartResponse(cart))
}

// updateQuantityRequest is the request body for updating an item's quantity.
type updateQuantityRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	VariantID string `json:"variant_id"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

// UpdateQuantity godoc
// @Summary      Update item quantity in the cart
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string                 true  "User ID"
// @Param        body       body    updateQuantityRequest  true  "Quantity update"
// @Success      200  {object}  cartResponse
// @Failure      400  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /cart/items [patch]
// @Security     BearerAuth
func (h *CartHandler) UpdateQuantity(c *gin.Context) {
	userID := getUserID(c)
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-ID header is required"})
		return
	}

	var req updateQuantityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.cartUC.UpdateQuantity(c.Request.Context(), userID, req.ProductID, req.VariantID, req.Quantity)
	if err != nil {
		h.handleUseCaseError(c, err, "failed to update quantity")
		return
	}

	c.JSON(http.StatusOK, toCartResponse(cart))
}

// removeItemRequest is the request body for removing an item from the cart.
type removeItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	VariantID string `json:"variant_id"`
}

// RemoveItem godoc
// @Summary      Remove an item from the cart
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string             true  "User ID"
// @Param        body       body    removeItemRequest  true  "Item to remove"
// @Success      200  {object}  cartResponse
// @Failure      400  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /cart/items [delete]
// @Security     BearerAuth
func (h *CartHandler) RemoveItem(c *gin.Context) {
	userID := getUserID(c)
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-ID header is required"})
		return
	}

	var req removeItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.cartUC.RemoveItem(c.Request.Context(), userID, req.ProductID, req.VariantID)
	if err != nil {
		h.handleUseCaseError(c, err, "failed to remove item")
		return
	}

	c.JSON(http.StatusOK, toCartResponse(cart))
}

// ClearCart godoc
// @Summary      Clear all items from the cart
// @Tags         Cart
// @Produce      json
// @Param        X-User-ID  header  string  true  "User ID"
// @Success      200  {object}  object{message=string}
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /cart [delete]
// @Security     BearerAuth
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID := getUserID(c)
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-ID header is required"})
		return
	}

	if err := h.cartUC.ClearCart(c.Request.Context(), userID); err != nil {
		h.logger.Error().Err(err).Str("user_id", userID).Msg("failed to clear cart")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cart cleared"})
}

// mergeCartRequest is the request body for merging guest cart items.
type mergeCartRequest struct {
	Items []domain.CartItem `json:"items" binding:"required"`
}

// MergeCart godoc
// @Summary      Merge guest cart items into authenticated user's cart
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string            true  "User ID"
// @Param        body       body    mergeCartRequest  true  "Guest cart items to merge"
// @Success      200  {object}  cartResponse
// @Failure      400  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /cart/merge [post]
// @Security     BearerAuth
func (h *CartHandler) MergeCart(c *gin.Context) {
	userID := getUserID(c)
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-ID header is required"})
		return
	}

	var req mergeCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.cartUC.MergeCart(c.Request.Context(), userID, req.Items)
	if err != nil {
		h.handleUseCaseError(c, err, "failed to merge cart")
		return
	}

	c.JSON(http.StatusOK, toCartResponse(cart))
}

// Health handles GET /health
func (h *CartHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Ready handles GET /ready — deep health check including database connectivity.
func (h *CartHandler) Ready(c *gin.Context) {
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

// handleUseCaseError maps use case errors to appropriate HTTP responses.
func (h *CartHandler) handleUseCaseError(c *gin.Context, err error, msg string) {
	switch {
	case errors.Is(err, usecase.ErrInvalidUserID),
		errors.Is(err, usecase.ErrInvalidProduct),
		errors.Is(err, usecase.ErrInvalidQuantity):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, usecase.ErrItemNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		h.logger.Error().Err(err).Msg(msg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
	}
}
