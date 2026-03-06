package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/southern-martin/ecommerce/services/payment/internal/domain"
	"github.com/southern-martin/ecommerce/services/payment/internal/usecase"
)

// Handler holds HTTP handlers for the payment service.
type Handler struct {
	paymentRepo    domain.PaymentRepository
	createPayment  *usecase.CreatePaymentUseCase
	confirmPayment *usecase.ConfirmPaymentUseCase
	wallet         *usecase.WalletUseCase
	payout         *usecase.PayoutUseCase
	refund         *usecase.RefundUseCase
	db             *gorm.DB
}

// NewHandler creates a new Handler.
func NewHandler(
	paymentRepo domain.PaymentRepository,
	createPayment *usecase.CreatePaymentUseCase,
	confirmPayment *usecase.ConfirmPaymentUseCase,
	wallet *usecase.WalletUseCase,
	payout *usecase.PayoutUseCase,
	refund *usecase.RefundUseCase,
	db *gorm.DB,
) *Handler {
	return &Handler{
		paymentRepo:    paymentRepo,
		createPayment:  createPayment,
		confirmPayment: confirmPayment,
		wallet:         wallet,
		payout:         payout,
		refund:         refund,
		db:             db,
	}
}

// Health returns a health check response.
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

// CreatePaymentIntent godoc
// @Summary      Create a payment intent
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        body  body  usecase.CreatePaymentInput  true  "Payment intent payload"
// @Success      201  {object}  usecase.CreatePaymentOutput
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /payments/create-intent [post]
// @Security     BearerAuth
func (h *Handler) CreatePaymentIntent(c *gin.Context) {
	var input usecase.CreatePaymentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract buyer ID from context (set by auth middleware) or from body.
	if buyerID, exists := c.Get("user_id"); exists {
		input.BuyerID = buyerID.(string)
	}

	output, err := h.createPayment.Execute(c.Request.Context(), input)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create payment intent")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create payment intent"})
		return
	}

	c.JSON(http.StatusCreated, output)
}

// ListPayments godoc
// @Summary      List payments for the authenticated buyer
// @Tags         Payments
// @Produce      json
// @Param        page       query  int  false  "Page number"   default(1)
// @Param        page_size  query  int  false  "Page size"     default(20)
// @Success      200  {object}  object{payments=[]domain.Payment,total=int64,page=int,page_size=int}
// @Failure      500  {object}  object{error=string}
// @Router       /payments [get]
// @Security     BearerAuth
func (h *Handler) ListPayments(c *gin.Context) {
	buyerID := c.GetString("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	payments, total, err := h.paymentRepo.List(c.Request.Context(), buyerID, page, pageSize)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list payments")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list payments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payments":  payments,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetPayment godoc
// @Summary      Get a payment by ID
// @Tags         Payments
// @Produce      json
// @Param        id  path  string  true  "Payment ID"
// @Success      200  {object}  domain.Payment
// @Failure      404  {object}  object{error=string}
// @Router       /payments/{id} [get]
func (h *Handler) GetPayment(c *gin.Context) {
	paymentID := c.Param("id")

	payment, err := h.paymentRepo.GetByID(c.Request.Context(), paymentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// HandleStripeWebhook godoc
// @Summary      Handle Stripe webhook event
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        body  body  usecase.WebhookEvent  true  "Stripe webhook event"
// @Success      200  {object}  object{received=bool}
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /payments/webhooks/stripe [post]
func (h *Handler) HandleStripeWebhook(c *gin.Context) {
	var event usecase.WebhookEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.confirmPayment.Execute(c.Request.Context(), event); err != nil {
		log.Error().Err(err).Msg("Failed to process webhook event")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process webhook"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}

// GetWalletBalance godoc
// @Summary      Get wallet balance for the authenticated seller
// @Tags         Payments
// @Produce      json
// @Param        seller_id  query  string  false  "Seller ID (fallback)"
// @Success      200  {object}  domain.SellerWallet
// @Failure      500  {object}  object{error=string}
// @Router       /payments/wallet [get]
// @Security     BearerAuth
func (h *Handler) GetWalletBalance(c *gin.Context) {
	sellerID := c.GetString("user_id")
	if sellerID == "" {
		sellerID = c.Query("seller_id")
	}

	wallet, err := h.wallet.GetBalance(c.Request.Context(), sellerID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get wallet balance")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get wallet balance"})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

// ListWalletTransactions godoc
// @Summary      List wallet transactions for the authenticated seller
// @Tags         Payments
// @Produce      json
// @Param        seller_id  query  string  false  "Seller ID (fallback)"
// @Param        page       query  int     false  "Page number"   default(1)
// @Param        page_size  query  int     false  "Page size"     default(20)
// @Success      200  {object}  object{transactions=[]domain.WalletTransaction,total=int64,page=int,page_size=int}
// @Failure      500  {object}  object{error=string}
// @Router       /payments/wallet/transactions [get]
// @Security     BearerAuth
func (h *Handler) ListWalletTransactions(c *gin.Context) {
	sellerID := c.GetString("user_id")
	if sellerID == "" {
		sellerID = c.Query("seller_id")
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	transactions, total, err := h.wallet.ListTransactions(c.Request.Context(), sellerID, page, pageSize)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list wallet transactions")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"total":        total,
		"page":         page,
		"page_size":    pageSize,
	})
}

// RequestPayout godoc
// @Summary      Request a seller payout
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        body  body  usecase.RequestPayoutInput  true  "Payout request payload"
// @Success      201  {object}  domain.Payout
// @Failure      400  {object}  object{error=string}
// @Router       /payments/payouts [post]
// @Security     BearerAuth
func (h *Handler) RequestPayout(c *gin.Context) {
	var input usecase.RequestPayoutInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if sellerID, exists := c.Get("user_id"); exists {
		input.SellerID = sellerID.(string)
	}

	payout, err := h.payout.RequestPayout(c.Request.Context(), input)
	if err != nil {
		log.Error().Err(err).Msg("Failed to request payout")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, payout)
}

// ListPayouts godoc
// @Summary      List payouts for the authenticated seller
// @Tags         Payments
// @Produce      json
// @Param        seller_id  query  string  false  "Seller ID (fallback)"
// @Param        page       query  int     false  "Page number"   default(1)
// @Param        page_size  query  int     false  "Page size"     default(20)
// @Success      200  {object}  object{payouts=[]domain.Payout,total=int64,page=int,page_size=int}
// @Failure      500  {object}  object{error=string}
// @Router       /payments/payouts [get]
// @Security     BearerAuth
func (h *Handler) ListPayouts(c *gin.Context) {
	sellerID := c.GetString("user_id")
	if sellerID == "" {
		sellerID = c.Query("seller_id")
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	payouts, total, err := h.payout.ListPayouts(c.Request.Context(), sellerID, page, pageSize)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list payouts")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list payouts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payouts":   payouts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
