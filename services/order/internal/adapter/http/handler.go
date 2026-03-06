package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/southern-martin/ecommerce/services/order/internal/domain"
	"github.com/southern-martin/ecommerce/services/order/internal/usecase"
)

// Handler holds the HTTP handlers for the order service.
type Handler struct {
	createOrder   *usecase.CreateOrderUseCase
	getOrder      *usecase.GetOrderUseCase
	updateStatus  *usecase.UpdateOrderStatusUseCase
	cancelOrder   *usecase.CancelOrderUseCase
	db            *gorm.DB
}

// NewHandler creates a new Handler instance.
func NewHandler(
	createOrder *usecase.CreateOrderUseCase,
	getOrder *usecase.GetOrderUseCase,
	updateStatus *usecase.UpdateOrderStatusUseCase,
	cancelOrder *usecase.CancelOrderUseCase,
	db *gorm.DB,
) *Handler {
	return &Handler{
		createOrder:  createOrder,
		getOrder:     getOrder,
		updateStatus: updateStatus,
		cancelOrder:  cancelOrder,
		db:           db,
	}
}

// --- Request/Response DTOs ---

type createOrderRequest struct {
	BuyerID         string             `json:"buyer_id" binding:"required"`
	BuyerEmail      string             `json:"buyer_email"`
	Currency        string             `json:"currency"`
	ShippingAddress addressDTO         `json:"shipping_address" binding:"required"`
	Items           []orderItemRequest `json:"items" binding:"required,min=1"`
}

type addressDTO struct {
	FullName    string `json:"full_name"`
	Line1       string `json:"line1"`
	Line2       string `json:"line2"`
	City        string `json:"city"`
	State       string `json:"state"`
	PostalCode  string `json:"postal_code"`
	CountryCode string `json:"country_code"`
	Phone       string `json:"phone"`
}

type orderItemRequest struct {
	ProductID      string `json:"product_id" binding:"required"`
	VariantID      string `json:"variant_id"`
	ProductName    string `json:"product_name" binding:"required"`
	VariantName    string `json:"variant_name"`
	SKU            string `json:"sku"`
	Quantity       int    `json:"quantity" binding:"required,min=1"`
	UnitPriceCents int64  `json:"unit_price_cents" binding:"required,min=1"`
	SellerID       string `json:"seller_id" binding:"required"`
	ImageURL       string `json:"image_url"`
}

type updateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type orderResponse struct {
	ID              string              `json:"id"`
	OrderNumber     string              `json:"order_number"`
	BuyerID         string              `json:"buyer_id"`
	Status          string              `json:"status"`
	SubtotalCents   int64               `json:"subtotal_cents"`
	ShippingCents   int64               `json:"shipping_cents"`
	TaxCents        int64               `json:"tax_cents"`
	DiscountCents   int64               `json:"discount_cents"`
	TotalCents      int64               `json:"total_cents"`
	Currency        string              `json:"currency"`
	ShippingAddress addressDTO          `json:"shipping_address"`
	Items           []orderItemResponse `json:"items"`
	SellerOrders    []sellerOrderResponse `json:"seller_orders"`
	CreatedAt       string              `json:"created_at"`
	UpdatedAt       string              `json:"updated_at"`
}

type orderItemResponse struct {
	ID             string `json:"id"`
	OrderID        string `json:"order_id"`
	ProductID      string `json:"product_id"`
	VariantID      string `json:"variant_id"`
	ProductName    string `json:"product_name"`
	VariantName    string `json:"variant_name"`
	SKU            string `json:"sku"`
	Quantity       int    `json:"quantity"`
	UnitPriceCents int64  `json:"unit_price_cents"`
	TotalCents     int64  `json:"total_cents"`
	SellerID       string `json:"seller_id"`
	ImageURL       string `json:"image_url"`
}

type sellerOrderResponse struct {
	ID            string `json:"id"`
	OrderID       string `json:"order_id"`
	SellerID      string `json:"seller_id"`
	Status        string `json:"status"`
	SubtotalCents int64  `json:"subtotal_cents"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type listResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int64       `json:"total_pages"`
}

// --- Handlers ---

// CreateOrder godoc
// @Summary      Create a new order
// @Tags         Buyer Orders
// @Accept       json
// @Produce      json
// @Param        body  body      createOrderRequest  true  "Order creation payload"
// @Success      201   {object}  object{data=orderResponse}
// @Failure      400   {object}  object{error=string}
// @Router       /orders [post]
func (h *Handler) CreateOrder(c *gin.Context) {
	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var items []usecase.CreateOrderItemInput
	for _, item := range req.Items {
		items = append(items, usecase.CreateOrderItemInput{
			ProductID:      item.ProductID,
			VariantID:      item.VariantID,
			ProductName:    item.ProductName,
			VariantName:    item.VariantName,
			SKU:            item.SKU,
			Quantity:       item.Quantity,
			UnitPriceCents: item.UnitPriceCents,
			SellerID:       item.SellerID,
			ImageURL:       item.ImageURL,
		})
	}

	input := usecase.CreateOrderInput{
		BuyerID:    req.BuyerID,
		BuyerEmail: req.BuyerEmail,
		Currency:   req.Currency,
		ShippingAddress: domain.Address{
			FullName:    req.ShippingAddress.FullName,
			Line1:       req.ShippingAddress.Line1,
			Line2:       req.ShippingAddress.Line2,
			City:        req.ShippingAddress.City,
			State:       req.ShippingAddress.State,
			PostalCode:  req.ShippingAddress.PostalCode,
			CountryCode: req.ShippingAddress.CountryCode,
			Phone:       req.ShippingAddress.Phone,
		},
		Items: items,
	}

	order, err := h.createOrder.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": toOrderResponse(order)})
}

// GetOrder godoc
// @Summary      Get an order by ID
// @Tags         Buyer Orders
// @Produce      json
// @Param        id   path      string  true  "Order ID"
// @Success      200  {object}  object{data=orderResponse}
// @Failure      404  {object}  object{error=string}
// @Router       /orders/{id} [get]
func (h *Handler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := h.getOrder.GetOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": toOrderResponse(order)})
}

// ListOrders godoc
// @Summary      List orders for a buyer
// @Tags         Buyer Orders
// @Produce      json
// @Param        X-User-ID  header  string  false  "User ID"
// @Param        buyer_id   query   string  false  "Buyer ID"
// @Param        status     query   string  false  "Order status filter"
// @Param        page       query   int     false  "Page number"  default(1)
// @Param        page_size  query   int     false  "Page size"    default(20)
// @Success      200  {object}  listResponse
// @Failure      500  {object}  object{error=string}
// @Router       /orders [get]
// @Security     BearerAuth
func (h *Handler) ListOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	buyerID := c.Query("buyer_id")
	if buyerID == "" {
		buyerID = c.GetHeader("X-User-ID")
	}
	status := c.Query("status")

	filter := domain.OrderFilter{
		BuyerID:  buyerID,
		Status:   status,
		Page:     page,
		PageSize: pageSize,
	}

	orders, total, err := h.getOrder.ListOrders(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var resp []orderResponse
	for _, o := range orders {
		resp = append(resp, toOrderResponse(o))
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

// CancelOrder godoc
// @Summary      Cancel an order
// @Tags         Buyer Orders
// @Produce      json
// @Param        id        path   string  true  "Order ID"
// @Param        buyer_id  query  string  true  "Buyer ID"
// @Success      200  {object}  object{data=orderResponse}
// @Failure      400  {object}  object{error=string}
// @Router       /orders/{id}/cancel [post]
func (h *Handler) CancelOrder(c *gin.Context) {
	id := c.Param("id")
	buyerID := c.Query("buyer_id")
	if buyerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "buyer_id query parameter is required"})
		return
	}

	order, err := h.cancelOrder.Execute(c.Request.Context(), id, buyerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toOrderResponse(order)})
}

// ListSellerOrders godoc
// @Summary      List seller orders
// @Tags         Seller Orders
// @Produce      json
// @Param        seller_id  query  string  true   "Seller ID"
// @Param        page       query  int     false  "Page number"  default(1)
// @Param        page_size  query  int     false  "Page size"    default(20)
// @Success      200  {object}  listResponse
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /seller/orders [get]
func (h *Handler) ListSellerOrders(c *gin.Context) {
	sellerID := c.Query("seller_id")
	if sellerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "seller_id query parameter is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	sellerOrders, total, err := h.getOrder.ListSellerOrders(c.Request.Context(), sellerID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var resp []sellerOrderResponse
	for _, so := range sellerOrders {
		resp = append(resp, toSellerOrderResponse(so))
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

// GetSellerOrder godoc
// @Summary      Get a seller order by ID
// @Tags         Seller Orders
// @Produce      json
// @Param        id   path      string  true  "Seller Order ID"
// @Success      200  {object}  object{data=sellerOrderResponse}
// @Failure      404  {object}  object{error=string}
// @Router       /seller/orders/{id} [get]
func (h *Handler) GetSellerOrder(c *gin.Context) {
	id := c.Param("id")
	sellerOrder, err := h.getOrder.GetSellerOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": toSellerOrderResponse(sellerOrder)})
}

// UpdateSellerOrderStatus godoc
// @Summary      Update seller order status
// @Tags         Seller Orders
// @Accept       json
// @Produce      json
// @Param        id    path      string              true  "Seller Order ID"
// @Param        body  body      updateStatusRequest  true  "New status"
// @Success      200   {object}  object{data=sellerOrderResponse}
// @Failure      400   {object}  object{error=string}
// @Router       /seller/orders/{id}/status [patch]
func (h *Handler) UpdateSellerOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var req updateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sellerOrder, err := h.updateStatus.Execute(c.Request.Context(), id, domain.OrderStatus(req.Status))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toSellerOrderResponse(sellerOrder)})
}

// Health handles GET /health
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

// --- Converters ---

func toOrderResponse(o *domain.Order) orderResponse {
	resp := orderResponse{
		ID:            o.ID,
		OrderNumber:   o.OrderNumber,
		BuyerID:       o.BuyerID,
		Status:        string(o.Status),
		SubtotalCents: o.SubtotalCents,
		ShippingCents: o.ShippingCents,
		TaxCents:      o.TaxCents,
		DiscountCents: o.DiscountCents,
		TotalCents:    o.TotalCents,
		Currency:      o.Currency,
		ShippingAddress: addressDTO{
			FullName:    o.ShippingAddress.FullName,
			Line1:       o.ShippingAddress.Line1,
			Line2:       o.ShippingAddress.Line2,
			City:        o.ShippingAddress.City,
			State:       o.ShippingAddress.State,
			PostalCode:  o.ShippingAddress.PostalCode,
			CountryCode: o.ShippingAddress.CountryCode,
			Phone:       o.ShippingAddress.Phone,
		},
		CreatedAt: o.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: o.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	for _, item := range o.Items {
		resp.Items = append(resp.Items, orderItemResponse{
			ID:             item.ID,
			OrderID:        item.OrderID,
			ProductID:      item.ProductID,
			VariantID:      item.VariantID,
			ProductName:    item.ProductName,
			VariantName:    item.VariantName,
			SKU:            item.SKU,
			Quantity:       item.Quantity,
			UnitPriceCents: item.UnitPriceCents,
			TotalCents:     item.TotalCents,
			SellerID:       item.SellerID,
			ImageURL:       item.ImageURL,
		})
	}

	for _, so := range o.SellerOrders {
		resp.SellerOrders = append(resp.SellerOrders, sellerOrderResponse{
			ID:            so.ID,
			OrderID:       so.OrderID,
			SellerID:      so.SellerID,
			Status:        string(so.Status),
			SubtotalCents: so.SubtotalCents,
			CreatedAt:     so.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:     so.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return resp
}

func toSellerOrderResponse(so *domain.SellerOrder) sellerOrderResponse {
	return sellerOrderResponse{
		ID:            so.ID,
		OrderID:       so.OrderID,
		SellerID:      so.SellerID,
		Status:        string(so.Status),
		SubtotalCents: so.SubtotalCents,
		CreatedAt:     so.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     so.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
