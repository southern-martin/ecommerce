package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

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
}

// NewHandler creates a new Handler.
func NewHandler(
	paymentRepo domain.PaymentRepository,
	createPayment *usecase.CreatePaymentUseCase,
	confirmPayment *usecase.ConfirmPaymentUseCase,
	wallet *usecase.WalletUseCase,
	payout *usecase.PayoutUseCase,
	refund *usecase.RefundUseCase,
) *Handler {
	return &Handler{
		paymentRepo:    paymentRepo,
		createPayment:  createPayment,
		confirmPayment: confirmPayment,
		wallet:         wallet,
		payout:         payout,
		refund:         refund,
	}
}

// Health returns a health check response.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// CreatePaymentIntent creates a new payment intent.
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

// ListPayments lists payments for the authenticated buyer.
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

// GetPayment retrieves a single payment by ID.
func (h *Handler) GetPayment(c *gin.Context) {
	paymentID := c.Param("id")

	payment, err := h.paymentRepo.GetByID(c.Request.Context(), paymentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// HandleStripeWebhook handles incoming Stripe webhook events.
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

// GetWalletBalance returns the wallet balance for the authenticated seller.
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

// ListWalletTransactions lists wallet transactions for the authenticated seller.
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

// RequestPayout handles a seller's payout request.
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

// ListPayouts lists payouts for the authenticated seller.
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
