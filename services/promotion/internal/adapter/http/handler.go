package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/southern-martin/ecommerce/services/promotion/internal/domain"
	"github.com/southern-martin/ecommerce/services/promotion/internal/usecase"
)

// Handler holds the HTTP handlers for the promotion service.
type Handler struct {
	couponUC    *usecase.CouponUseCase
	flashSaleUC *usecase.FlashSaleUseCase
	bundleUC    *usecase.BundleUseCase
}

// NewHandler creates a new Handler instance.
func NewHandler(
	couponUC *usecase.CouponUseCase,
	flashSaleUC *usecase.FlashSaleUseCase,
	bundleUC *usecase.BundleUseCase,
) *Handler {
	return &Handler{
		couponUC:    couponUC,
		flashSaleUC: flashSaleUC,
		bundleUC:    bundleUC,
	}
}

// --- Request/Response DTOs ---

type createCouponRequest struct {
	Code             string   `json:"code" binding:"required"`
	Type             string   `json:"type" binding:"required"`
	DiscountValue    int64    `json:"discount_value" binding:"required"`
	MinOrderCents    int64    `json:"min_order_cents"`
	MaxDiscountCents int64    `json:"max_discount_cents"`
	UsageLimit       int      `json:"usage_limit"`
	PerUserLimit     int      `json:"per_user_limit"`
	Scope            string   `json:"scope"`
	ScopeIDs         []string `json:"scope_ids"`
	StartsAt         string   `json:"starts_at"`
	ExpiresAt        string   `json:"expires_at" binding:"required"`
}

type validateCouponRequest struct {
	Code       string `json:"code" binding:"required"`
	OrderCents int64  `json:"order_cents"`
}

type updateCouponRequest struct {
	IsActive         *bool   `json:"is_active"`
	UsageLimit       *int    `json:"usage_limit"`
	PerUserLimit     *int    `json:"per_user_limit"`
	MaxDiscountCents *int64  `json:"max_discount_cents"`
	ExpiresAt        *string `json:"expires_at"`
}

type createFlashSaleRequest struct {
	Name     string                    `json:"name" binding:"required"`
	StartsAt string                    `json:"starts_at" binding:"required"`
	EndsAt   string                    `json:"ends_at" binding:"required"`
	Items    []flashSaleItemRequestDTO `json:"items"`
}

type flashSaleItemRequestDTO struct {
	ProductID      string `json:"product_id" binding:"required"`
	VariantID      string `json:"variant_id"`
	SalePriceCents int64  `json:"sale_price_cents" binding:"required"`
	QuantityLimit  int    `json:"quantity_limit"`
}

type updateFlashSaleRequest struct {
	IsActive *bool   `json:"is_active"`
	Name     *string `json:"name"`
}

type createBundleRequest struct {
	Name             string   `json:"name" binding:"required"`
	ProductIDs       []string `json:"product_ids" binding:"required,min=2"`
	BundlePriceCents int64    `json:"bundle_price_cents" binding:"required"`
	SavingsCents     int64    `json:"savings_cents"`
}

type updateBundleRequest struct {
	IsActive         *bool  `json:"is_active"`
	Name             *string `json:"name"`
	BundlePriceCents *int64  `json:"bundle_price_cents"`
	SavingsCents     *int64  `json:"savings_cents"`
}

type couponResponse struct {
	ID               string   `json:"id"`
	Code             string   `json:"code"`
	Type             string   `json:"type"`
	DiscountValue    int64    `json:"discount_value"`
	MinOrderCents    int64    `json:"min_order_cents"`
	MaxDiscountCents int64    `json:"max_discount_cents"`
	UsageLimit       int      `json:"usage_limit"`
	UsageCount       int      `json:"usage_count"`
	PerUserLimit     int      `json:"per_user_limit"`
	Scope            string   `json:"scope"`
	ScopeIDs         []string `json:"scope_ids"`
	CreatedBy        string   `json:"created_by"`
	StartsAt         string   `json:"starts_at"`
	ExpiresAt        string   `json:"expires_at"`
	IsActive         bool     `json:"is_active"`
	CreatedAt        string   `json:"created_at"`
}

type flashSaleResponse struct {
	ID        string                  `json:"id"`
	Name      string                  `json:"name"`
	StartsAt  string                  `json:"starts_at"`
	EndsAt    string                  `json:"ends_at"`
	IsActive  bool                    `json:"is_active"`
	Items     []flashSaleItemResponse `json:"items"`
	CreatedAt string                  `json:"created_at"`
}

type flashSaleItemResponse struct {
	ID             string `json:"id"`
	FlashSaleID    string `json:"flash_sale_id"`
	ProductID      string `json:"product_id"`
	VariantID      string `json:"variant_id"`
	SalePriceCents int64  `json:"sale_price_cents"`
	QuantityLimit  int    `json:"quantity_limit"`
	SoldCount      int    `json:"sold_count"`
}

type bundleResponse struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	SellerID         string   `json:"seller_id"`
	ProductIDs       []string `json:"product_ids"`
	BundlePriceCents int64    `json:"bundle_price_cents"`
	SavingsCents     int64    `json:"savings_cents"`
	IsActive         bool     `json:"is_active"`
	CreatedAt        string   `json:"created_at"`
}

type listResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int64       `json:"total_pages"`
}

// --- Coupon Handlers ---

// ValidateCoupon handles POST /api/v1/coupons/validate
func (h *Handler) ValidateCoupon(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "X-User-ID header is required"})
		return
	}

	var req validateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coupon, discountCents, err := h.couponUC.ValidateCoupon(c.Request.Context(), usecase.ValidateCouponInput{
		Code:       req.Code,
		UserID:     userID,
		OrderCents: req.OrderCents,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"coupon":         toCouponResponse(coupon),
			"discount_cents": discountCents,
		},
	})
}

// ListActiveCoupons handles GET /api/v1/coupons
func (h *Handler) ListActiveCoupons(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	coupons, total, err := h.couponUC.ListCoupons(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var resp []couponResponse
	for _, cp := range coupons {
		resp = append(resp, toCouponResponse(cp))
	}

	totalPages := total / int64(pageSize)
	if total%int64(pageSize) != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, listResponse{
		Data:       resp,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// CreateSellerCoupon handles POST /api/v1/seller/coupons
func (h *Handler) CreateSellerCoupon(c *gin.Context) {
	sellerID := c.GetHeader("X-User-ID")
	if sellerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "X-User-ID header is required"})
		return
	}

	var req createCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var startsAt time.Time
	if req.StartsAt != "" {
		var err error
		startsAt, err = time.Parse(time.RFC3339, req.StartsAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid starts_at format"})
			return
		}
	}

	expiresAt, err := time.Parse(time.RFC3339, req.ExpiresAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expires_at format"})
		return
	}

	coupon, err := h.couponUC.CreateCoupon(c.Request.Context(), usecase.CreateCouponInput{
		Code:             req.Code,
		Type:             req.Type,
		DiscountValue:    req.DiscountValue,
		MinOrderCents:    req.MinOrderCents,
		MaxDiscountCents: req.MaxDiscountCents,
		UsageLimit:       req.UsageLimit,
		PerUserLimit:     req.PerUserLimit,
		Scope:            req.Scope,
		ScopeIDs:         req.ScopeIDs,
		CreatedBy:        sellerID,
		StartsAt:         startsAt,
		ExpiresAt:        expiresAt,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": toCouponResponse(coupon)})
}

// ListSellerCoupons handles GET /api/v1/seller/coupons
func (h *Handler) ListSellerCoupons(c *gin.Context) {
	sellerID := c.GetHeader("X-User-ID")
	if sellerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "X-User-ID header is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	coupons, total, err := h.couponUC.ListCouponsBySeller(c.Request.Context(), sellerID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var resp []couponResponse
	for _, cp := range coupons {
		resp = append(resp, toCouponResponse(cp))
	}

	totalPages := total / int64(pageSize)
	if total%int64(pageSize) != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, listResponse{
		Data:       resp,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// GetSellerCoupon handles GET /api/v1/seller/coupons/:id
func (h *Handler) GetSellerCoupon(c *gin.Context) {
	id := c.Param("id")
	coupon, err := h.couponUC.GetCoupon(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": toCouponResponse(coupon)})
}

// UpdateSellerCoupon handles PATCH /api/v1/seller/coupons/:id
func (h *Handler) UpdateSellerCoupon(c *gin.Context) {
	id := c.Param("id")

	coupon, err := h.couponUC.GetCoupon(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var req updateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.IsActive != nil {
		coupon.IsActive = *req.IsActive
	}
	if req.UsageLimit != nil {
		coupon.UsageLimit = *req.UsageLimit
	}
	if req.PerUserLimit != nil {
		coupon.PerUserLimit = *req.PerUserLimit
	}
	if req.MaxDiscountCents != nil {
		coupon.MaxDiscountCents = *req.MaxDiscountCents
	}
	if req.ExpiresAt != nil {
		expiresAt, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expires_at format"})
			return
		}
		coupon.ExpiresAt = expiresAt
	}

	if err := h.couponUC.UpdateCoupon(c.Request.Context(), coupon); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toCouponResponse(coupon)})
}

// DeleteSellerCoupon handles DELETE /api/v1/seller/coupons/:id (soft delete via deactivation)
func (h *Handler) DeleteSellerCoupon(c *gin.Context) {
	id := c.Param("id")

	coupon, err := h.couponUC.GetCoupon(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	coupon.IsActive = false
	if err := h.couponUC.UpdateCoupon(c.Request.Context(), coupon); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "coupon deactivated"})
}

// --- Admin Coupon Handlers ---

// AdminCreateCoupon handles POST /api/v1/admin/promotions/coupons
func (h *Handler) AdminCreateCoupon(c *gin.Context) {
	var req createCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var startsAt time.Time
	if req.StartsAt != "" {
		var err error
		startsAt, err = time.Parse(time.RFC3339, req.StartsAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid starts_at format"})
			return
		}
	}

	expiresAt, err := time.Parse(time.RFC3339, req.ExpiresAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expires_at format"})
		return
	}

	coupon, err := h.couponUC.CreateCoupon(c.Request.Context(), usecase.CreateCouponInput{
		Code:             req.Code,
		Type:             req.Type,
		DiscountValue:    req.DiscountValue,
		MinOrderCents:    req.MinOrderCents,
		MaxDiscountCents: req.MaxDiscountCents,
		UsageLimit:       req.UsageLimit,
		PerUserLimit:     req.PerUserLimit,
		Scope:            req.Scope,
		ScopeIDs:         req.ScopeIDs,
		CreatedBy:        "platform",
		StartsAt:         startsAt,
		ExpiresAt:        expiresAt,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": toCouponResponse(coupon)})
}

// AdminListCoupons handles GET /api/v1/admin/promotions/coupons
func (h *Handler) AdminListCoupons(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	coupons, total, err := h.couponUC.ListCoupons(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var resp []couponResponse
	for _, cp := range coupons {
		resp = append(resp, toCouponResponse(cp))
	}

	totalPages := total / int64(pageSize)
	if total%int64(pageSize) != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, listResponse{
		Data:       resp,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// AdminGetCoupon handles GET /api/v1/admin/promotions/coupons/:id
func (h *Handler) AdminGetCoupon(c *gin.Context) {
	id := c.Param("id")
	coupon, err := h.couponUC.GetCoupon(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": toCouponResponse(coupon)})
}

// AdminUpdateCoupon handles PATCH /api/v1/admin/promotions/coupons/:id
func (h *Handler) AdminUpdateCoupon(c *gin.Context) {
	id := c.Param("id")

	coupon, err := h.couponUC.GetCoupon(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var req updateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.IsActive != nil {
		coupon.IsActive = *req.IsActive
	}
	if req.UsageLimit != nil {
		coupon.UsageLimit = *req.UsageLimit
	}
	if req.PerUserLimit != nil {
		coupon.PerUserLimit = *req.PerUserLimit
	}
	if req.MaxDiscountCents != nil {
		coupon.MaxDiscountCents = *req.MaxDiscountCents
	}
	if req.ExpiresAt != nil {
		expiresAt, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expires_at format"})
			return
		}
		coupon.ExpiresAt = expiresAt
	}

	if err := h.couponUC.UpdateCoupon(c.Request.Context(), coupon); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toCouponResponse(coupon)})
}

// AdminDeleteCoupon handles DELETE /api/v1/admin/promotions/coupons/:id
func (h *Handler) AdminDeleteCoupon(c *gin.Context) {
	id := c.Param("id")

	coupon, err := h.couponUC.GetCoupon(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	coupon.IsActive = false
	if err := h.couponUC.UpdateCoupon(c.Request.Context(), coupon); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "coupon deactivated"})
}

// --- Flash Sale Handlers ---

// ListActiveFlashSales handles GET /api/v1/flash-sales
func (h *Handler) ListActiveFlashSales(c *gin.Context) {
	flashSales, err := h.flashSaleUC.ListActiveFlashSales(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var resp []flashSaleResponse
	for _, fs := range flashSales {
		resp = append(resp, toFlashSaleResponse(fs))
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// AdminCreateFlashSale handles POST /api/v1/admin/promotions/flash-sales
func (h *Handler) AdminCreateFlashSale(c *gin.Context) {
	var req createFlashSaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startsAt, err := time.Parse(time.RFC3339, req.StartsAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid starts_at format"})
		return
	}

	endsAt, err := time.Parse(time.RFC3339, req.EndsAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ends_at format"})
		return
	}

	var items []usecase.CreateFlashSaleItemInput
	for _, item := range req.Items {
		items = append(items, usecase.CreateFlashSaleItemInput{
			ProductID:      item.ProductID,
			VariantID:      item.VariantID,
			SalePriceCents: item.SalePriceCents,
			QuantityLimit:  item.QuantityLimit,
		})
	}

	flashSale, err := h.flashSaleUC.CreateFlashSale(c.Request.Context(), usecase.CreateFlashSaleInput{
		Name:     req.Name,
		StartsAt: startsAt,
		EndsAt:   endsAt,
		Items:    items,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": toFlashSaleResponse(flashSale)})
}

// AdminListFlashSales handles GET /api/v1/admin/promotions/flash-sales
func (h *Handler) AdminListFlashSales(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	flashSales, total, err := h.flashSaleUC.ListFlashSales(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var resp []flashSaleResponse
	for _, fs := range flashSales {
		resp = append(resp, toFlashSaleResponse(fs))
	}

	totalPages := total / int64(pageSize)
	if total%int64(pageSize) != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, listResponse{
		Data:       resp,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// AdminGetFlashSale handles GET /api/v1/admin/promotions/flash-sales/:id
func (h *Handler) AdminGetFlashSale(c *gin.Context) {
	id := c.Param("id")
	flashSale, err := h.flashSaleUC.GetFlashSale(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": toFlashSaleResponse(flashSale)})
}

// AdminUpdateFlashSale handles PATCH /api/v1/admin/promotions/flash-sales/:id
func (h *Handler) AdminUpdateFlashSale(c *gin.Context) {
	id := c.Param("id")

	flashSale, err := h.flashSaleUC.GetFlashSale(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var req updateFlashSaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.IsActive != nil {
		flashSale.IsActive = *req.IsActive
	}
	if req.Name != nil {
		flashSale.Name = *req.Name
	}

	if err := h.flashSaleUC.UpdateFlashSale(c.Request.Context(), flashSale); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toFlashSaleResponse(flashSale)})
}

// --- Bundle Handlers ---

// ListActiveBundles handles GET /api/v1/bundles
func (h *Handler) ListActiveBundles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	bundles, total, err := h.bundleUC.ListActiveBundles(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var resp []bundleResponse
	for _, b := range bundles {
		resp = append(resp, toBundleResponse(b))
	}

	totalPages := total / int64(pageSize)
	if total%int64(pageSize) != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, listResponse{
		Data:       resp,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// AdminCreateBundle handles POST /api/v1/admin/promotions/bundles
func (h *Handler) AdminCreateBundle(c *gin.Context) {
	var req createBundleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sellerID := c.GetHeader("X-User-ID")
	if sellerID == "" {
		sellerID = "platform"
	}

	bundle, err := h.bundleUC.CreateBundle(c.Request.Context(), usecase.CreateBundleInput{
		Name:             req.Name,
		SellerID:         sellerID,
		ProductIDs:       req.ProductIDs,
		BundlePriceCents: req.BundlePriceCents,
		SavingsCents:     req.SavingsCents,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": toBundleResponse(bundle)})
}

// AdminListBundles handles GET /api/v1/admin/promotions/bundles
func (h *Handler) AdminListBundles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	bundles, total, err := h.bundleUC.ListActiveBundles(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var resp []bundleResponse
	for _, b := range bundles {
		resp = append(resp, toBundleResponse(b))
	}

	totalPages := total / int64(pageSize)
	if total%int64(pageSize) != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, listResponse{
		Data:       resp,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// AdminGetBundle handles GET /api/v1/admin/promotions/bundles/:id
func (h *Handler) AdminGetBundle(c *gin.Context) {
	id := c.Param("id")
	bundle, err := h.bundleUC.GetBundle(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": toBundleResponse(bundle)})
}

// AdminUpdateBundle handles PATCH /api/v1/admin/promotions/bundles/:id
func (h *Handler) AdminUpdateBundle(c *gin.Context) {
	id := c.Param("id")

	bundle, err := h.bundleUC.GetBundle(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var req updateBundleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.IsActive != nil {
		bundle.IsActive = *req.IsActive
	}
	if req.Name != nil {
		bundle.Name = *req.Name
	}
	if req.BundlePriceCents != nil {
		bundle.BundlePriceCents = *req.BundlePriceCents
	}
	if req.SavingsCents != nil {
		bundle.SavingsCents = *req.SavingsCents
	}

	if err := h.bundleUC.UpdateBundle(c.Request.Context(), bundle); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toBundleResponse(bundle)})
}

// AdminDeleteBundle handles DELETE /api/v1/admin/promotions/bundles/:id
func (h *Handler) AdminDeleteBundle(c *gin.Context) {
	id := c.Param("id")

	bundle, err := h.bundleUC.GetBundle(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	bundle.IsActive = false
	if err := h.bundleUC.UpdateBundle(c.Request.Context(), bundle); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "bundle deactivated"})
}

// Health handles GET /health
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// --- Converters ---

func toCouponResponse(c *domain.Coupon) couponResponse {
	scopeIDs := c.ScopeIDs
	if scopeIDs == nil {
		scopeIDs = []string{}
	}
	return couponResponse{
		ID:               c.ID,
		Code:             c.Code,
		Type:             string(c.Type),
		DiscountValue:    c.DiscountValue,
		MinOrderCents:    c.MinOrderCents,
		MaxDiscountCents: c.MaxDiscountCents,
		UsageLimit:       c.UsageLimit,
		UsageCount:       c.UsageCount,
		PerUserLimit:     c.PerUserLimit,
		Scope:            string(c.Scope),
		ScopeIDs:         scopeIDs,
		CreatedBy:        c.CreatedBy,
		StartsAt:         c.StartsAt.Format("2006-01-02T15:04:05Z"),
		ExpiresAt:        c.ExpiresAt.Format("2006-01-02T15:04:05Z"),
		IsActive:         c.IsActive,
		CreatedAt:        c.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toFlashSaleResponse(fs *domain.FlashSale) flashSaleResponse {
	resp := flashSaleResponse{
		ID:        fs.ID,
		Name:      fs.Name,
		StartsAt:  fs.StartsAt.Format("2006-01-02T15:04:05Z"),
		EndsAt:    fs.EndsAt.Format("2006-01-02T15:04:05Z"),
		IsActive:  fs.IsActive,
		CreatedAt: fs.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	for _, item := range fs.Items {
		resp.Items = append(resp.Items, flashSaleItemResponse{
			ID:             item.ID,
			FlashSaleID:    item.FlashSaleID,
			ProductID:      item.ProductID,
			VariantID:      item.VariantID,
			SalePriceCents: item.SalePriceCents,
			QuantityLimit:  item.QuantityLimit,
			SoldCount:      item.SoldCount,
		})
	}
	if resp.Items == nil {
		resp.Items = []flashSaleItemResponse{}
	}
	return resp
}

func toBundleResponse(b *domain.Bundle) bundleResponse {
	productIDs := b.ProductIDs
	if productIDs == nil {
		productIDs = []string{}
	}
	return bundleResponse{
		ID:               b.ID,
		Name:             b.Name,
		SellerID:         b.SellerID,
		ProductIDs:       productIDs,
		BundlePriceCents: b.BundlePriceCents,
		SavingsCents:     b.SavingsCents,
		IsActive:         b.IsActive,
		CreatedAt:        b.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
