package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

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
	db         *gorm.DB
}

// NewHandler creates a new Handler.
func NewHandler(
	rateUC *usecase.RateUseCase,
	shipmentUC *usecase.ShipmentUseCase,
	labelUC *usecase.LabelUseCase,
	trackingUC *usecase.TrackingUseCase,
	carrierUC *usecase.CarrierUseCase,
	db *gorm.DB,
) *Handler {
	return &Handler{
		rateUC:     rateUC,
		shipmentUC: shipmentUC,
		labelUC:    labelUC,
		trackingUC: trackingUC,
		carrierUC:  carrierUC,
		db:         db,
	}
}

// Health returns a health check response.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "shipping"})
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

// --- Rate Handlers ---

type getRatesRequest struct {
	Origin      domain.Address `json:"origin" binding:"required"`
	Destination domain.Address `json:"destination" binding:"required"`
	WeightGrams int            `json:"weight_grams" binding:"required"`
	Currency    string         `json:"currency"`
}

// GetShippingRates godoc
// @Summary      Get shipping rates
// @Tags         Shipping
// @Accept       json
// @Produce      json
// @Param        body  body      getRatesRequest  true  "Rate calculation request"
// @Success      200   {object}  object{rates=[]object}
// @Failure      400   {object}  object{error=string}
// @Failure      500   {object}  object{error=string}
// @Router       /shipping/rates [post]
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

// CreateShipment godoc
// @Summary      Create a shipment
// @Tags         Shipping
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header    string                true  "User ID"
// @Param        body       body      createShipmentRequest true  "Shipment details"
// @Success      201        {object}  object{shipment=object}
// @Failure      400        {object}  object{error=string}
// @Failure      401        {object}  object{error=string}
// @Failure      500        {object}  object{error=string}
// @Router       /shipments [post]
// @Security     BearerAuth
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

// GetShipment godoc
// @Summary      Get a shipment by ID
// @Tags         Shipping
// @Produce      json
// @Param        id   path      string  true  "Shipment ID"
// @Success      200  {object}  object{shipment=object}
// @Failure      404  {object}  object{error=string}
// @Router       /shipments/{id} [get]
func (h *Handler) GetShipment(c *gin.Context) {
	id := c.Param("id")
	shipment, err := h.shipmentUC.GetShipment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shipment not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"shipment": shipment})
}

// GetTracking godoc
// @Summary      Get tracking info by tracking number
// @Tags         Shipping
// @Produce      json
// @Param        tracking_number  path      string  true  "Tracking number"
// @Success      200              {object}  object{tracking_number=string,carrier_code=string,status=string,events=[]object}
// @Failure      404              {object}  object{error=string}
// @Router       /tracking/{tracking_number} [get]
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

// GenerateLabel godoc
// @Summary      Generate a shipping label
// @Tags         Shipping
// @Produce      json
// @Param        X-User-ID  header    string  true  "User ID"
// @Param        id         path      string  true  "Shipment ID"
// @Success      200        {object}  object{shipment=object,tracking_number=string,label_url=string}
// @Failure      400        {object}  object{error=string}
// @Failure      401        {object}  object{error=string}
// @Router       /shipments/{id}/label [post]
// @Security     BearerAuth
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

// AddTrackingEvent godoc
// @Summary      Add a tracking event to a shipment
// @Tags         Shipping
// @Accept       json
// @Produce      json
// @Param        id    path      string                   true  "Shipment ID"
// @Param        body  body      addTrackingEventRequest  true  "Tracking event details"
// @Success      201   {object}  object{event=object}
// @Failure      400   {object}  object{error=string}
// @Router       /shipments/{id}/tracking [post]
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

// ListSellerShipments godoc
// @Summary      List shipments for the authenticated seller
// @Tags         Shipping
// @Produce      json
// @Param        X-User-ID  header    string  true   "User ID"
// @Param        page       query     int     false  "Page number"
// @Param        page_size  query     int     false  "Page size"
// @Success      200        {object}  object{shipments=[]object,total=int,page=int,page_size=int}
// @Failure      401        {object}  object{error=string}
// @Failure      500        {object}  object{error=string}
// @Router       /seller/shipments [get]
// @Security     BearerAuth
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

// SetupSellerCarrier godoc
// @Summary      Set up a carrier for the authenticated seller
// @Tags         Shipping
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header    string                     true  "User ID"
// @Param        body       body      setupSellerCarrierRequest  true  "Carrier setup details"
// @Success      200        {object}  object{credential=object}
// @Failure      400        {object}  object{error=string}
// @Failure      401        {object}  object{error=string}
// @Router       /seller/carriers [post]
// @Security     BearerAuth
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

// GetSellerCarriers godoc
// @Summary      List carriers configured for the authenticated seller
// @Tags         Shipping
// @Produce      json
// @Param        X-User-ID  header    string  true  "User ID"
// @Success      200        {object}  object{carriers=[]object}
// @Failure      401        {object}  object{error=string}
// @Failure      500        {object}  object{error=string}
// @Router       /seller/carriers [get]
// @Security     BearerAuth
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

// CreateCarrier godoc
// @Summary      Create a new carrier (admin)
// @Tags         Shipping
// @Accept       json
// @Produce      json
// @Param        body  body      createCarrierRequest  true  "Carrier details"
// @Success      201   {object}  object{carrier=object}
// @Failure      400   {object}  object{error=string}
// @Router       /admin/carriers [post]
// @Security     BearerAuth
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

// UpdateCarrier godoc
// @Summary      Update a carrier (admin)
// @Tags         Shipping
// @Accept       json
// @Produce      json
// @Param        code  path      string                true  "Carrier code"
// @Param        body  body      updateCarrierRequest  true  "Carrier update details"
// @Success      200   {object}  object{carrier=object}
// @Failure      400   {object}  object{error=string}
// @Router       /admin/carriers/{code} [patch]
// @Security     BearerAuth
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

// ListCarriers godoc
// @Summary      List all carriers (admin)
// @Tags         Shipping
// @Produce      json
// @Success      200  {object}  object{carriers=[]object}
// @Failure      500  {object}  object{error=string}
// @Router       /admin/carriers [get]
// @Security     BearerAuth
func (h *Handler) ListCarriers(c *gin.Context) {
	carriers, err := h.carrierUC.ListCarriers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"carriers": carriers})
}
