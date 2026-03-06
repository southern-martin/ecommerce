package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/southern-martin/ecommerce/services/loyalty/internal/usecase"
	"gorm.io/gorm"
)

// Handler holds all HTTP handlers for the loyalty service.
type Handler struct {
	membershipUC *usecase.MembershipUseCase
	pointsUC     *usecase.PointsUseCase
	tierUC       *usecase.TierUseCase
	db           *gorm.DB
}

// NewHandler creates a new Handler.
func NewHandler(
	membershipUC *usecase.MembershipUseCase,
	pointsUC *usecase.PointsUseCase,
	tierUC *usecase.TierUseCase,
	db *gorm.DB,
) *Handler {
	return &Handler{
		membershipUC: membershipUC,
		pointsUC:     pointsUC,
		tierUC:       tierUC,
		db:           db,
	}
}

// Health returns a health check response.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "loyalty"})
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

// GetMembership godoc
// @Summary      Get or create loyalty membership
// @Tags         Loyalty
// @Produce      json
// @Param        X-User-ID  header  string  true  "User ID"
// @Success      200  {object}  object{membership=object{user_id=string,tier=string,points_balance=int,lifetime_points=int,tier_expires_at=string,joined_at=string}}
// @Failure      401  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /loyalty/membership [get]
// @Security     BearerAuth
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

// GetPointsBalance godoc
// @Summary      Get user points balance
// @Tags         Loyalty
// @Produce      json
// @Param        X-User-ID  header  string  true  "User ID"
// @Success      200  {object}  object{user_id=string,points_balance=int64}
// @Failure      401  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Router       /loyalty/points [get]
// @Security     BearerAuth
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

// ListTransactions godoc
// @Summary      List user points transactions
// @Tags         Loyalty
// @Produce      json
// @Param        X-User-ID   header  string  true   "User ID"
// @Param        page        query   int     false  "Page number"   default(1)
// @Param        page_size   query   int     false  "Page size"     default(20)
// @Success      200  {object}  object{transactions=[]object{id=string,user_id=string,type=string,points=int,source=string,reference_id=string,description=string,created_at=string},total=int,page=int,page_size=int}
// @Failure      401  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /loyalty/transactions [get]
// @Security     BearerAuth
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

// RedeemPoints godoc
// @Summary      Redeem loyalty points
// @Tags         Loyalty
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string         true  "User ID"
// @Param        body       body    redeemRequest  true  "Redeem request"
// @Success      200  {object}  object{transaction=object{id=string,user_id=string,type=string,points=int,source=string,reference_id=string,description=string,created_at=string}}
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Router       /loyalty/redeem [post]
// @Security     BearerAuth
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

// ListTiers godoc
// @Summary      List all loyalty tiers
// @Tags         Loyalty
// @Produce      json
// @Success      200  {object}  object{tiers=[]object{name=string,min_points=int,cashback_rate=number,points_multiplier=number,free_shipping=bool}}
// @Failure      500  {object}  object{error=string}
// @Router       /loyalty/tiers [get]
func (h *Handler) ListTiers(c *gin.Context) {
	tiers, err := h.tierUC.GetAllTiers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tiers": tiers})
}
