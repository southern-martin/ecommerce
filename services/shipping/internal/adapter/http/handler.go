package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/southern-martin/ecommerce/services/shipping/internal/domain"
	"github.com/southern-martin/ecommerce/services/shipping/internal/usecase"
)

// Handler holds all HTTP handlers for the shipping service.
type Handler struct {
	rateUC     *usecase.RateUseCase
	shipmentUC *usecase.ShipmentUseCase
	labelUC    *usecase.LabelUseCase
	trackingUC *usecase.TrackingUseCase
	carrierUC  *usecase.CarrierUseCase
}

// NewHandler creates a new Handler.
func NewHandler(
	rateUC *usecase.RateUseCase,
	shipmentUC *usecase.ShipmentUseCase,
	labelUC *usecase.LabelUseCase,
	trackingUC *usecase.TrackingUseCase,
	carrierUC *usecase.CarrierUseCase,
) *Handler {
	return &Handler{
		rateUC:     rateUC,
		shipmentUC: shipmentUC,
		labelUC:    labelUC,
		trackingUC: trackingUC,
		carrierUC:  carrierUC,
	}
}

// Health returns a health check response.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "shipping"})
}

// --- Rate Handlers ---

type getRatesRequest struct {
	Origin      domain.Address `json:"origin" binding:"required"`
	Destination domain.Address `json:"destination" binding:"required"`
	WeightGrams int            `json:"weight_grams" binding:"required"`
	Currency    string         `json:"currency"`
}

func (h *Handler) GetShippingRates(c *gin.Context) {
	var req getRatesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rates, err := h.rateUC.GetShippingRates(c.Request.Context(), usecase.GetShippingRatesRequest{
		Origin:      req.Origin,
		Destination: req.Destination,
		WeightGrams: req.WeightGrams,
		Currency:    req.Currency,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rates": rates})
}

// --- Shipment Handlers ---

type createShipmentRequest struct {
	OrderID     string              `json:"order_id" binding:"required"`
	CarrierCode string              `json:"carrier_code" binding:"required"`
	ServiceCode string              `json:"service_code"`
	Origin      domain.Address      `json:"origin" binding:"required"`
	Destination domain.Address      `json:"destination" binding:"required"`
	WeightGrams int                 `json:"weight_grams"`
	RateCents   int64               `json:"rate_cents"`
	Currency    string              `json:"currency"`
	Items       []domain.ShipmentItem `json:"items"`
}

func (h *Handler) CreateShipment(c *gin.Context) {
	sellerID := c.GetHeader("X-User-ID")
	if sellerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	var req createShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shipment, err := h.shipmentUC.CreateShipment(c.Request.Context(), usecase.CreateShipmentRequest{
		OrderID:     req.OrderID,
		SellerID:    sellerID,
		CarrierCode: req.CarrierCode,
		ServiceCode: req.ServiceCode,
		Origin:      req.Origin,
		Destination: req.Destination,
		WeightGrams: req.WeightGrams,
		RateCents:   req.RateCents,
		Currency:    req.Currency,
		Items:       req.Items,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"shipment": shipment})
}

func (h *Handler) GetShipment(c *gin.Context) {
	id := c.Param("id")
	shipment, err := h.shipmentUC.GetShipment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shipment not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"shipment": shipment})
}

func (h *Handler) GetTracking(c *gin.Context) {
	trackingNumber := c.Param("tracking_number")
	shipment, err := h.shipmentUC.GetShipmentByTracking(c.Request.Context(), trackingNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tracking not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"tracking_number": shipment.TrackingNumber,
		"carrier_code":    shipment.CarrierCode,
		"status":          shipment.Status,
		"events":          shipment.TrackingEvents,
	})
}

func (h *Handler) GenerateLabel(c *gin.Context) {
	sellerID := c.GetHeader("X-User-ID")
	if sellerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	id := c.Param("id")
	shipment, err := h.labelUC.GenerateLabel(c.Request.Context(), id, sellerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"shipment":        shipment,
		"tracking_number": shipment.TrackingNumber,
		"label_url":       shipment.LabelURL,
	})
}

type addTrackingEventRequest struct {
	Status      string `json:"status" binding:"required"`
	Description string `json:"description"`
	Location    string `json:"location"`
	EventAt     string `json:"event_at"`
}

func (h *Handler) AddTrackingEvent(c *gin.Context) {
	id := c.Param("id")
	var req addTrackingEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := h.trackingUC.AddTrackingEvent(c.Request.Context(), usecase.AddTrackingEventRequest{
		ShipmentID:  id,
		Status:      req.Status,
		Description: req.Description,
		Location:    req.Location,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"event": event})
}

func (h *Handler) ListSellerShipments(c *gin.Context) {
	sellerID := c.GetHeader("X-User-ID")
	if sellerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	shipments, total, err := h.shipmentUC.ListSellerShipments(c.Request.Context(), sellerID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"shipments": shipments, "total": total, "page": page, "page_size": pageSize})
}

// --- Carrier Handlers ---

type setupSellerCarrierRequest struct {
	CarrierCode string `json:"carrier_code" binding:"required"`
	Credentials string `json:"credentials" binding:"required"`
}

func (h *Handler) SetupSellerCarrier(c *gin.Context) {
	sellerID := c.GetHeader("X-User-ID")
	if sellerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	var req setupSellerCarrierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cred, err := h.carrierUC.SetupSellerCarrier(c.Request.Context(), sellerID, req.CarrierCode, req.Credentials)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"credential": cred})
}

func (h *Handler) GetSellerCarriers(c *gin.Context) {
	sellerID := c.GetHeader("X-User-ID")
	if sellerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	creds, err := h.carrierUC.GetSellerCarriers(c.Request.Context(), sellerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"carriers": creds})
}

type createCarrierRequest struct {
	Code               string   `json:"code" binding:"required"`
	Name               string   `json:"name" binding:"required"`
	SupportedCountries []string `json:"supported_countries"`
	APIBaseURL         string   `json:"api_base_url"`
}

func (h *Handler) CreateCarrier(c *gin.Context) {
	var req createCarrierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	carrier := &domain.Carrier{
		Code:               req.Code,
		Name:               req.Name,
		IsActive:           true,
		SupportedCountries: req.SupportedCountries,
		APIBaseURL:         req.APIBaseURL,
	}

	if err := h.carrierUC.CreateCarrier(c.Request.Context(), carrier); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"carrier": carrier})
}

type updateCarrierRequest struct {
	Name               string   `json:"name"`
	IsActive           *bool    `json:"is_active"`
	SupportedCountries []string `json:"supported_countries"`
	APIBaseURL         string   `json:"api_base_url"`
}

func (h *Handler) UpdateCarrier(c *gin.Context) {
	code := c.Param("code")
	var req updateCarrierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	carrier := &domain.Carrier{
		Code:               code,
		Name:               req.Name,
		IsActive:           isActive,
		SupportedCountries: req.SupportedCountries,
		APIBaseURL:         req.APIBaseURL,
	}

	if err := h.carrierUC.UpdateCarrier(c.Request.Context(), carrier); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"carrier": carrier})
}

func (h *Handler) ListCarriers(c *gin.Context) {
	carriers, err := h.carrierUC.ListCarriers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"carriers": carriers})
}
