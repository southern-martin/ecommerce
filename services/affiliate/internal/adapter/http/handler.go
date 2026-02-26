package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/southern-martin/ecommerce/services/affiliate/internal/domain"
	"github.com/southern-martin/ecommerce/services/affiliate/internal/usecase"
)

// Handler holds all HTTP handlers for the affiliate service.
type Handler struct {
	programUC  *usecase.ProgramUseCase
	linkUC     *usecase.LinkUseCase
	referralUC *usecase.ReferralUseCase
	payoutUC   *usecase.PayoutUseCase
}

// NewHandler creates a new Handler.
func NewHandler(
	programUC *usecase.ProgramUseCase,
	linkUC *usecase.LinkUseCase,
	referralUC *usecase.ReferralUseCase,
	payoutUC *usecase.PayoutUseCase,
) *Handler {
	return &Handler{
		programUC:  programUC,
		linkUC:     linkUC,
		referralUC: referralUC,
		payoutUC:   payoutUC,
	}
}

// Health returns a health check response.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "affiliate"})
}

// --- Link Handlers ---

type createLinkRequest struct {
	TargetURL string `json:"target_url" binding:"required"`
}

func (h *Handler) CreateLink(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	var req createLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	link, err := h.linkUC.CreateLink(c.Request.Context(), usecase.CreateLinkRequest{
		UserID:    userID,
		TargetURL: req.TargetURL,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"link": link})
}

func (h *Handler) ListLinks(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	links, total, err := h.linkUC.ListUserLinks(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"links": links, "total": total, "page": page, "page_size": pageSize})
}

func (h *Handler) GetStats(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	links, _, err := h.linkUC.ListUserLinks(c.Request.Context(), userID, 1, 1000)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var totalClicks, totalConversions, totalEarnings int64
	for _, link := range links {
		totalClicks += link.ClickCount
		totalConversions += link.ConversionCount
		totalEarnings += link.TotalEarningsCents
	}

	c.JSON(http.StatusOK, gin.H{
		"total_clicks":         totalClicks,
		"total_conversions":    totalConversions,
		"total_earnings_cents": totalEarnings,
		"link_count":           len(links),
	})
}

func (h *Handler) TrackClick(c *gin.Context) {
	code := c.Param("code")
	link, err := h.linkUC.TrackClick(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "referral link not found"})
		return
	}

	c.Redirect(http.StatusFound, link.TargetURL)
}

// --- Referral Handlers ---

func (h *Handler) ListReferrals(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	referrals, total, err := h.referralUC.ListReferrals(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"referrals": referrals, "total": total, "page": page, "page_size": pageSize})
}

// --- Payout Handlers ---

type requestPayoutRequest struct {
	AmountCents  int64  `json:"amount_cents" binding:"required"`
	PayoutMethod string `json:"payout_method" binding:"required"`
}

func (h *Handler) RequestPayout(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	var req requestPayoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payout, err := h.payoutUC.RequestPayout(c.Request.Context(), usecase.RequestPayoutRequest{
		UserID:       userID,
		AmountCents:  req.AmountCents,
		PayoutMethod: domain.PayoutMethod(req.PayoutMethod),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"payout": payout})
}

// --- Admin Handlers ---

func (h *Handler) GetProgram(c *gin.Context) {
	program, err := h.programUC.GetProgram(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "program not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"program": program})
}

type updateProgramRequest struct {
	CommissionRate     *float64 `json:"commission_rate"`
	MinPayoutCents     *int64   `json:"min_payout_cents"`
	CookieDays         *int     `json:"cookie_days"`
	ReferrerBonusCents *int64   `json:"referrer_bonus_cents"`
	ReferredBonusCents *int64   `json:"referred_bonus_cents"`
	IsActive           *bool    `json:"is_active"`
}

func (h *Handler) UpdateProgram(c *gin.Context) {
	var req updateProgramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	program, err := h.programUC.UpdateProgram(c.Request.Context(), usecase.UpdateProgramRequest{
		CommissionRate:     req.CommissionRate,
		MinPayoutCents:     req.MinPayoutCents,
		CookieDays:         req.CookieDays,
		ReferrerBonusCents: req.ReferrerBonusCents,
		ReferredBonusCents: req.ReferredBonusCents,
		IsActive:           req.IsActive,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"program": program})
}

func (h *Handler) ListAllPayouts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	payouts, total, err := h.payoutUC.ListAllPayouts(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payouts": payouts, "total": total, "page": page, "page_size": pageSize})
}

type updatePayoutStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func (h *Handler) UpdatePayoutStatus(c *gin.Context) {
	id := c.Param("id")
	var req updatePayoutStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payout, err := h.payoutUC.UpdatePayoutStatus(c.Request.Context(), id, domain.PayoutStatus(req.Status))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payout": payout})
}
