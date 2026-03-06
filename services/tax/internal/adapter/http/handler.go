package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/southern-martin/ecommerce/services/tax/internal/domain"
	"github.com/southern-martin/ecommerce/services/tax/internal/usecase"
)

// Handler handles HTTP requests for the tax service.
type Handler struct {
	calculateTax *usecase.CalculateTaxUseCase
	manageRules  *usecase.ManageRulesUseCase
	manageZones  *usecase.ManageZonesUseCase
	db           *gorm.DB
}

// NewHandler creates a new HTTP handler.
func NewHandler(
	calculateTax *usecase.CalculateTaxUseCase,
	manageRules *usecase.ManageRulesUseCase,
	manageZones *usecase.ManageZonesUseCase,
	db *gorm.DB,
) *Handler {
	return &Handler{
		calculateTax: calculateTax,
		manageRules:  manageRules,
		manageZones:  manageZones,
		db:           db,
	}
}

// Health returns service health status.
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

// ListZones godoc
// @Summary      List all tax zones
// @Tags         Tax
// @Produce      json
// @Success      200  {object}  object{zones=[]zoneResponse}
// @Failure      500  {object}  object{error=string}
// @Router       /tax/zones [get]
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

// ListRules godoc
// @Summary      List tax rules
// @Tags         Tax
// @Produce      json
// @Param        zone_id  query  string  false  "Filter by zone ID"
// @Success      200  {object}  object{rules=[]ruleResponse}
// @Failure      500  {object}  object{error=string}
// @Router       /admin/tax/rules [get]
// @Security     BearerAuth
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

// CreateRule godoc
// @Summary      Create a tax rule
// @Tags         Tax
// @Accept       json
// @Produce      json
// @Param        body  body  createRuleRequest  true  "Tax rule to create"
// @Success      201  {object}  object{rule=ruleResponse}
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /admin/tax/rules [post]
// @Security     BearerAuth
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

// UpdateRule godoc
// @Summary      Update a tax rule
// @Tags         Tax
// @Accept       json
// @Produce      json
// @Param        id    path  string             true  "Rule ID"
// @Param        body  body  updateRuleRequest  true  "Fields to update"
// @Success      200  {object}  object{rule=ruleResponse}
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /admin/tax/rules/{id} [patch]
// @Security     BearerAuth
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

// DeleteRule godoc
// @Summary      Delete a tax rule
// @Tags         Tax
// @Produce      json
// @Param        id  path  string  true  "Rule ID"
// @Success      200  {object}  object{message=string}
// @Failure      500  {object}  object{error=string}
// @Router       /admin/tax/rules/{id} [delete]
// @Security     BearerAuth
func (h *Handler) DeleteRule(c *gin.Context) {
	id := c.Param("id")

	if err := h.manageRules.DeleteRule(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete rule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rule deleted"})
}

// CalculateTax godoc
// @Summary      Calculate tax for items
// @Tags         Tax
// @Accept       json
// @Produce      json
// @Param        body  body  calculateTaxRequest  true  "Items and shipping address"
// @Success      200  {object}  object{subtotal_cents=int64,tax_amount_cents=int64,breakdown=[]taxBreakdownResponse}
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /tax/calculate [post]
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
