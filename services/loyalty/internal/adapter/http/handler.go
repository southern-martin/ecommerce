package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/southern-martin/ecommerce/services/loyalty/internal/usecase"
)

// Handler holds all HTTP handlers for the loyalty service.
type Handler struct {
	membershipUC *usecase.MembershipUseCase
	pointsUC     *usecase.PointsUseCase
	tierUC       *usecase.TierUseCase
}

// NewHandler creates a new Handler.
func NewHandler(
	membershipUC *usecase.MembershipUseCase,
	pointsUC *usecase.PointsUseCase,
	tierUC *usecase.TierUseCase,
) *Handler {
	return &Handler{
		membershipUC: membershipUC,
		pointsUC:     pointsUC,
		tierUC:       tierUC,
	}
}

// Health returns a health check response.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "loyalty"})
}

// GetMembership returns a user's loyalty membership.
func (h *Handler) GetMembership(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	membership, err := h.membershipUC.GetMembership(c.Request.Context(), userID)
	if err != nil {
		// If not found, create one
		membership, err = h.membershipUC.CreateMembership(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"membership": membership})
}

// GetPointsBalance returns a user's points balance.
func (h *Handler) GetPointsBalance(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	balance, err := h.pointsUC.GetBalance(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "membership not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": userID, "points_balance": balance})
}

// ListTransactions returns a user's points transactions.
func (h *Handler) ListTransactions(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	transactions, total, err := h.pointsUC.ListTransactions(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions, "total": total, "page": page, "page_size": pageSize})
}

type redeemRequest struct {
	Points  int64  `json:"points" binding:"required"`
	OrderID string `json:"order_id" binding:"required"`
}

// RedeemPoints redeems loyalty points.
func (h *Handler) RedeemPoints(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	var req redeemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.pointsUC.RedeemPoints(c.Request.Context(), usecase.RedeemPointsRequest{
		UserID:  userID,
		Points:  req.Points,
		OrderID: req.OrderID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction": tx})
}

// ListTiers returns all loyalty tiers.
func (h *Handler) ListTiers(c *gin.Context) {
	tiers, err := h.tierUC.GetAllTiers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tiers": tiers})
}
