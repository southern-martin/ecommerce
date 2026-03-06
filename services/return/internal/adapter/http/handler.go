package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/southern-martin/ecommerce/services/return/internal/domain"
	"github.com/southern-martin/ecommerce/services/return/internal/usecase"
)

// Handler holds all HTTP handlers for the return service.
type Handler struct {
	createReturnUC *usecase.CreateReturnUseCase
	manageReturnUC *usecase.ManageReturnUseCase
	disputeUC      *usecase.DisputeUseCase
	db             *gorm.DB
}

// NewHandler creates a new Handler.
func NewHandler(
	createReturnUC *usecase.CreateReturnUseCase,
	manageReturnUC *usecase.ManageReturnUseCase,
	disputeUC *usecase.DisputeUseCase,
	db *gorm.DB,
) *Handler {
	return &Handler{
		createReturnUC: createReturnUC,
		manageReturnUC: manageReturnUC,
		disputeUC:      disputeUC,
		db:             db,
	}
}

// Health returns a health check response.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "return"})
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

// --- Return Handlers ---

type createReturnRequest struct {
	OrderID           string                      `json:"order_id" binding:"required"`
	SellerID          string                      `json:"seller_id" binding:"required"`
	Reason            string                      `json:"reason" binding:"required"`
	Description       string                      `json:"description"`
	ImageURLs         []string                    `json:"image_urls"`
	RefundAmountCents int64                       `json:"refund_amount_cents"`
	RefundMethod      string                      `json:"refund_method"`
	Items             []createReturnItemRequest   `json:"items" binding:"required"`
}

type createReturnItemRequest struct {
	OrderItemID string `json:"order_item_id" binding:"required"`
	ProductID   string `json:"product_id" binding:"required"`
	VariantID   string `json:"variant_id"`
	Quantity    int    `json:"quantity" binding:"required"`
	Reason      string `json:"reason"`
}

// CreateReturn godoc
// @Summary      Create a return request
// @Tags         Returns
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header    string               true  "User ID"
// @Param        body       body      createReturnRequest  true  "Return request details"
// @Success      201        {object}  object{return=object}
// @Failure      400        {object}  object{error=string}
// @Failure      401        {object}  object{error=string}
// @Router       /returns [post]
// @Security     BearerAuth
func (h *Handler) CreateReturn(c *gin.Context) {
	buyerID := c.GetHeader("X-User-ID")
	if buyerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	var req createReturnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	items := make([]usecase.CreateReturnItemRequest, len(req.Items))
	for i, item := range req.Items {
		items[i] = usecase.CreateReturnItemRequest{
			OrderItemID: item.OrderItemID,
			ProductID:   item.ProductID,
			VariantID:   item.VariantID,
			Quantity:    item.Quantity,
			Reason:      item.Reason,
		}
	}

	ret, err := h.createReturnUC.Execute(c.Request.Context(), usecase.CreateReturnRequest{
		OrderID:           req.OrderID,
		BuyerID:           buyerID,
		SellerID:          req.SellerID,
		Reason:            req.Reason,
		Description:       req.Description,
		ImageURLs:         req.ImageURLs,
		RefundAmountCents: req.RefundAmountCents,
		RefundMethod:      req.RefundMethod,
		Items:             items,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"return": ret})
}

// ListBuyerReturns godoc
// @Summary      List returns for the authenticated buyer
// @Tags         Returns
// @Produce      json
// @Param        X-User-ID  header    string  true   "User ID"
// @Param        page       query     int     false  "Page number"
// @Param        page_size  query     int     false  "Page size"
// @Success      200        {object}  object{returns=[]object,total=int,page=int,page_size=int}
// @Failure      401        {object}  object{error=string}
// @Failure      500        {object}  object{error=string}
// @Router       /returns [get]
// @Security     BearerAuth
func (h *Handler) ListBuyerReturns(c *gin.Context) {
	buyerID := c.GetHeader("X-User-ID")
	if buyerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	returns, total, err := h.manageReturnUC.ListBuyerReturns(c.Request.Context(), buyerID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"returns": returns, "total": total, "page": page, "page_size": pageSize})
}

// GetReturn godoc
// @Summary      Get a return by ID
// @Tags         Returns
// @Produce      json
// @Param        id   path      string  true  "Return ID"
// @Success      200  {object}  object{return=object}
// @Failure      404  {object}  object{error=string}
// @Router       /returns/{id} [get]
func (h *Handler) GetReturn(c *gin.Context) {
	id := c.Param("id")
	ret, err := h.manageReturnUC.GetReturn(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "return not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"return": ret})
}

// ListSellerReturns godoc
// @Summary      List returns for the authenticated seller
// @Tags         Returns
// @Produce      json
// @Param        X-User-ID  header    string  true   "User ID"
// @Param        page       query     int     false  "Page number"
// @Param        page_size  query     int     false  "Page size"
// @Success      200        {object}  object{returns=[]object,total=int,page=int,page_size=int}
// @Failure      401        {object}  object{error=string}
// @Failure      500        {object}  object{error=string}
// @Router       /seller/returns [get]
// @Security     BearerAuth
func (h *Handler) ListSellerReturns(c *gin.Context) {
	sellerID := c.GetHeader("X-User-ID")
	if sellerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	returns, total, err := h.manageReturnUC.ListSellerReturns(c.Request.Context(), sellerID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"returns": returns, "total": total, "page": page, "page_size": pageSize})
}

type approveReturnRequest struct {
	RefundAmountCents int64 `json:"refund_amount_cents"`
}

// ApproveReturn godoc
// @Summary      Approve a return request (seller)
// @Tags         Returns
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header    string                true  "User ID"
// @Param        id         path      string                true  "Return ID"
// @Param        body       body      approveReturnRequest  false "Optional refund amount override"
// @Success      200        {object}  object{return=object}
// @Failure      400        {object}  object{error=string}
// @Failure      401        {object}  object{error=string}
// @Router       /seller/returns/{id}/approve [patch]
// @Security     BearerAuth
func (h *Handler) ApproveReturn(c *gin.Context) {
	sellerID := c.GetHeader("X-User-ID")
	if sellerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	id := c.Param("id")
	var req approveReturnRequest
	_ = c.ShouldBindJSON(&req)

	ret, err := h.manageReturnUC.ApproveReturn(c.Request.Context(), id, sellerID, req.RefundAmountCents)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"return": ret})
}

// RejectReturn godoc
// @Summary      Reject a return request (seller)
// @Tags         Returns
// @Produce      json
// @Param        X-User-ID  header    string  true  "User ID"
// @Param        id         path      string  true  "Return ID"
// @Success      200        {object}  object{return=object}
// @Failure      400        {object}  object{error=string}
// @Failure      401        {object}  object{error=string}
// @Router       /seller/returns/{id}/reject [patch]
// @Security     BearerAuth
func (h *Handler) RejectReturn(c *gin.Context) {
	sellerID := c.GetHeader("X-User-ID")
	if sellerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	id := c.Param("id")
	ret, err := h.manageReturnUC.RejectReturn(c.Request.Context(), id, sellerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"return": ret})
}

type updateReturnStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// UpdateReturnStatus godoc
// @Summary      Update return status (seller)
// @Tags         Returns
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header    string                     true  "User ID"
// @Param        id         path      string                     true  "Return ID"
// @Param        body       body      updateReturnStatusRequest  true  "New status"
// @Success      200        {object}  object{return=object}
// @Failure      400        {object}  object{error=string}
// @Failure      401        {object}  object{error=string}
// @Router       /seller/returns/{id}/status [patch]
// @Security     BearerAuth
func (h *Handler) UpdateReturnStatus(c *gin.Context) {
	sellerID := c.GetHeader("X-User-ID")
	if sellerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	id := c.Param("id")
	var req updateReturnStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ret, err := h.manageReturnUC.UpdateReturnStatus(c.Request.Context(), id, sellerID, domain.ReturnStatus(req.Status))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"return": ret})
}

// --- Dispute Handlers ---

type createDisputeRequest struct {
	OrderID     string `json:"order_id" binding:"required"`
	ReturnID    string `json:"return_id"`
	SellerID    string `json:"seller_id" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// CreateDispute godoc
// @Summary      Create a dispute
// @Tags         Returns
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header    string                true  "User ID"
// @Param        body       body      createDisputeRequest  true  "Dispute details"
// @Success      201        {object}  object{dispute=object}
// @Failure      400        {object}  object{error=string}
// @Failure      401        {object}  object{error=string}
// @Router       /disputes [post]
// @Security     BearerAuth
func (h *Handler) CreateDispute(c *gin.Context) {
	buyerID := c.GetHeader("X-User-ID")
	if buyerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	var req createDisputeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dispute, err := h.disputeUC.CreateDispute(c.Request.Context(), usecase.CreateDisputeRequest{
		OrderID:     req.OrderID,
		ReturnID:    req.ReturnID,
		BuyerID:     buyerID,
		SellerID:    req.SellerID,
		Type:        req.Type,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"dispute": dispute})
}

// ListBuyerDisputes godoc
// @Summary      List disputes for the authenticated buyer
// @Tags         Returns
// @Produce      json
// @Param        X-User-ID  header    string  true   "User ID"
// @Param        page       query     int     false  "Page number"
// @Param        page_size  query     int     false  "Page size"
// @Success      200        {object}  object{disputes=[]object,total=int,page=int,page_size=int}
// @Failure      401        {object}  object{error=string}
// @Failure      500        {object}  object{error=string}
// @Router       /disputes [get]
// @Security     BearerAuth
func (h *Handler) ListBuyerDisputes(c *gin.Context) {
	buyerID := c.GetHeader("X-User-ID")
	if buyerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	disputes, total, err := h.disputeUC.ListBuyerDisputes(c.Request.Context(), buyerID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"disputes": disputes, "total": total, "page": page, "page_size": pageSize})
}

// GetDispute godoc
// @Summary      Get a dispute by ID
// @Tags         Returns
// @Produce      json
// @Param        id   path      string  true  "Dispute ID"
// @Success      200  {object}  object{dispute=object}
// @Failure      404  {object}  object{error=string}
// @Router       /disputes/{id} [get]
func (h *Handler) GetDispute(c *gin.Context) {
	id := c.Param("id")
	dispute, err := h.disputeUC.GetDispute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "dispute not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dispute": dispute})
}

type addMessageRequest struct {
	Message     string   `json:"message" binding:"required"`
	Attachments []string `json:"attachments"`
}

// AddMessage godoc
// @Summary      Add a message to a dispute
// @Tags         Returns
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header    string             true   "User ID"
// @Param        id         path      string             true   "Dispute ID"
// @Param        role       query     string             false  "Sender role (buyer or seller)"
// @Param        body       body      addMessageRequest  true   "Message details"
// @Success      201        {object}  object{message=object}
// @Failure      400        {object}  object{error=string}
// @Failure      401        {object}  object{error=string}
// @Router       /disputes/{id}/messages [post]
// @Security     BearerAuth
func (h *Handler) AddMessage(c *gin.Context) {
	senderID := c.GetHeader("X-User-ID")
	if senderID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	disputeID := c.Param("id")
	senderRole := c.DefaultQuery("role", "buyer")

	var req addMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg, err := h.disputeUC.AddMessage(c.Request.Context(), usecase.AddMessageRequest{
		DisputeID:   disputeID,
		SenderID:    senderID,
		SenderRole:  senderRole,
		Message:     req.Message,
		Attachments: req.Attachments,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": msg})
}

// ListAllDisputes godoc
// @Summary      List all disputes (admin)
// @Tags         Returns
// @Produce      json
// @Param        page       query     int  false  "Page number"
// @Param        page_size  query     int  false  "Page size"
// @Success      200        {object}  object{disputes=[]object,total=int,page=int,page_size=int}
// @Failure      500        {object}  object{error=string}
// @Router       /admin/disputes [get]
// @Security     BearerAuth
func (h *Handler) ListAllDisputes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	disputes, total, err := h.disputeUC.ListAllDisputes(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"disputes": disputes, "total": total, "page": page, "page_size": pageSize})
}

type resolveDisputeRequest struct {
	Resolution string `json:"resolution" binding:"required"`
	Status     string `json:"status" binding:"required"`
}

// ResolveDispute godoc
// @Summary      Resolve a dispute (admin)
// @Tags         Returns
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header    string                  true  "User ID"
// @Param        id         path      string                  true  "Dispute ID"
// @Param        body       body      resolveDisputeRequest   true  "Resolution details"
// @Success      200        {object}  object{dispute=object}
// @Failure      400        {object}  object{error=string}
// @Failure      401        {object}  object{error=string}
// @Router       /admin/disputes/{id}/resolve [patch]
// @Security     BearerAuth
func (h *Handler) ResolveDispute(c *gin.Context) {
	adminID := c.GetHeader("X-User-ID")
	if adminID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	disputeID := c.Param("id")
	var req resolveDisputeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dispute, err := h.disputeUC.ResolveDispute(c.Request.Context(), usecase.ResolveDisputeRequest{
		DisputeID:  disputeID,
		Resolution: req.Resolution,
		ResolvedBy: adminID,
		Status:     req.Status,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dispute": dispute})
}
