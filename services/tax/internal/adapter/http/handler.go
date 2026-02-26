package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/southern-martin/ecommerce/services/tax/internal/domain"
	"github.com/southern-martin/ecommerce/services/tax/internal/usecase"
)

// Handler handles HTTP requests for the tax service.
type Handler struct {
	calculateTax *usecase.CalculateTaxUseCase
	manageRules  *usecase.ManageRulesUseCase
	manageZones  *usecase.ManageZonesUseCase
}

// NewHandler creates a new HTTP handler.
func NewHandler(
	calculateTax *usecase.CalculateTaxUseCase,
	manageRules *usecase.ManageRulesUseCase,
	manageZones *usecase.ManageZonesUseCase,
) *Handler {
	return &Handler{
		calculateTax: calculateTax,
		manageRules:  manageRules,
		manageZones:  manageZones,
	}
}

// Health returns service health status.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// ListZones returns all tax zones.
func (h *Handler) ListZones(c *gin.Context) {
	zones, err := h.manageZones.ListZones(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list zones"})
		return
	}

	response := make([]zoneResponse, len(zones))
	for i, z := range zones {
		response[i] = toZoneResponse(z)
	}

	c.JSON(http.StatusOK, gin.H{"zones": response})
}

// ListRules returns all active tax rules.
func (h *Handler) ListRules(c *gin.Context) {
	zoneID := c.Query("zone_id")

	var rules []*domain.TaxRule
	var err error

	if zoneID != "" {
		rules, err = h.manageRules.ListRulesByZone(c.Request.Context(), zoneID)
	} else {
		rules, err = h.manageRules.ListRules(c.Request.Context())
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list rules"})
		return
	}

	response := make([]ruleResponse, len(rules))
	for i, r := range rules {
		response[i] = toRuleResponse(r)
	}

	c.JSON(http.StatusOK, gin.H{"rules": response})
}

// CreateRule creates a new tax rule.
func (h *Handler) CreateRule(c *gin.Context) {
	var req createRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := usecase.CreateRuleInput{
		ZoneID:    req.ZoneID,
		TaxName:   req.TaxName,
		Rate:      req.Rate,
		Category:  req.Category,
		Inclusive: req.Inclusive,
		StartsAt:  req.StartsAt,
		ExpiresAt: req.ExpiresAt,
	}

	rule, err := h.manageRules.CreateRule(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create rule"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"rule": toRuleResponse(rule)})
}

// UpdateRule updates an existing tax rule.
func (h *Handler) UpdateRule(c *gin.Context) {
	id := c.Param("id")

	var req updateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := usecase.UpdateRuleInput{
		TaxName:   req.TaxName,
		Rate:      req.Rate,
		Category:  req.Category,
		Inclusive: req.Inclusive,
		IsActive:  req.IsActive,
		ExpiresAt: req.ExpiresAt,
	}

	rule, err := h.manageRules.UpdateRule(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update rule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rule": toRuleResponse(rule)})
}

// DeleteRule deletes a tax rule.
func (h *Handler) DeleteRule(c *gin.Context) {
	id := c.Param("id")

	if err := h.manageRules.DeleteRule(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete rule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rule deleted"})
}

// CalculateTax calculates tax for the given items and address.
func (h *Handler) CalculateTax(c *gin.Context) {
	var req calculateTaxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domainReq := &domain.TaxCalculationRequest{
		Items: make([]domain.TaxItem, len(req.Items)),
		ShippingAddress: domain.TaxAddress{
			CountryCode: req.ShippingAddress.CountryCode,
			StateCode:   req.ShippingAddress.StateCode,
			City:        req.ShippingAddress.City,
			PostalCode:  req.ShippingAddress.PostalCode,
		},
	}

	for i, item := range req.Items {
		domainReq.Items[i] = domain.TaxItem{
			ProductID:  item.ProductID,
			VariantID:  item.VariantID,
			Category:   item.Category,
			PriceCents: item.PriceCents,
			Quantity:   item.Quantity,
		}
	}

	result, err := h.calculateTax.Execute(c.Request.Context(), domainReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate tax"})
		return
	}

	breakdown := make([]taxBreakdownResponse, len(result.Breakdown))
	for i, b := range result.Breakdown {
		breakdown[i] = taxBreakdownResponse{
			TaxName:      b.TaxName,
			Rate:         b.Rate,
			AmountCents:  b.AmountCents,
			Jurisdiction: b.Jurisdiction,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"subtotal_cents":   result.SubtotalCents,
		"tax_amount_cents": result.TaxAmountCents,
		"breakdown":        breakdown,
	})
}

// --- Request / Response types ---

type createRuleRequest struct {
	ZoneID    string     `json:"zone_id" binding:"required"`
	TaxName   string     `json:"tax_name" binding:"required"`
	Rate      float64    `json:"rate" binding:"required"`
	Category  string     `json:"category"`
	Inclusive bool       `json:"inclusive"`
	StartsAt  *time.Time `json:"starts_at"`
	ExpiresAt *time.Time `json:"expires_at"`
}

type updateRuleRequest struct {
	TaxName   *string    `json:"tax_name"`
	Rate      *float64   `json:"rate"`
	Category  *string    `json:"category"`
	Inclusive *bool      `json:"inclusive"`
	IsActive  *bool      `json:"is_active"`
	ExpiresAt *time.Time `json:"expires_at"`
}

type calculateTaxRequest struct {
	Items           []taxItemRequest    `json:"items" binding:"required"`
	ShippingAddress taxAddressRequest   `json:"shipping_address" binding:"required"`
}

type taxItemRequest struct {
	ProductID  string `json:"product_id"`
	VariantID  string `json:"variant_id"`
	Category   string `json:"category"`
	PriceCents int64  `json:"price_cents" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required"`
}

type taxAddressRequest struct {
	CountryCode string `json:"country_code" binding:"required"`
	StateCode   string `json:"state_code"`
	City        string `json:"city"`
	PostalCode  string `json:"postal_code"`
}

type zoneResponse struct {
	ID          string `json:"id"`
	CountryCode string `json:"country_code"`
	StateCode   string `json:"state_code"`
	Name        string `json:"name"`
}

type ruleResponse struct {
	ID        string     `json:"id"`
	ZoneID    string     `json:"zone_id"`
	TaxName   string     `json:"tax_name"`
	Rate      float64    `json:"rate"`
	Category  string     `json:"category"`
	Inclusive bool       `json:"inclusive"`
	StartsAt  time.Time  `json:"starts_at"`
	ExpiresAt *time.Time `json:"expires_at"`
	IsActive  bool       `json:"is_active"`
}

type taxBreakdownResponse struct {
	TaxName      string  `json:"tax_name"`
	Rate         float64 `json:"rate"`
	AmountCents  int64   `json:"amount_cents"`
	Jurisdiction string  `json:"jurisdiction"`
}

func toZoneResponse(z *domain.TaxZone) zoneResponse {
	return zoneResponse{
		ID:          z.ID,
		CountryCode: z.CountryCode,
		StateCode:   z.StateCode,
		Name:        z.Name,
	}
}

func toRuleResponse(r *domain.TaxRule) ruleResponse {
	return ruleResponse{
		ID:        r.ID,
		ZoneID:    r.ZoneID,
		TaxName:   r.TaxName,
		Rate:      r.Rate,
		Category:  r.Category,
		Inclusive:  r.Inclusive,
		StartsAt:  r.StartsAt,
		ExpiresAt: r.ExpiresAt,
		IsActive:  r.IsActive,
	}
}
