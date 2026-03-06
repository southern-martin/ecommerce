package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/southern-martin/ecommerce/services/chat/internal/usecase"
	"gorm.io/gorm"
)

// Handler holds all HTTP handlers for the chat service.
type Handler struct {
	conversationUC *usecase.ConversationUseCase
	messageUC      *usecase.MessageUseCase
	db             *gorm.DB
}

// NewHandler creates a new Handler.
func NewHandler(conversationUC *usecase.ConversationUseCase, messageUC *usecase.MessageUseCase, db *gorm.DB) *Handler {
	return &Handler{
		conversationUC: conversationUC,
		messageUC:      messageUC,
		db:             db,
	}
}

// Health returns a health check response.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "chat"})
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

// --- Conversation Handlers ---

type createConversationRequest struct {
	Type     string `json:"type"`
	BuyerID  string `json:"buyer_id" binding:"required"`
	SellerID string `json:"seller_id" binding:"required"`
	OrderID  string `json:"order_id"`
	Subject  string `json:"subject"`
}

// CreateConversation godoc
// @Summary      Create a new conversation
// @Tags         Chat
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string                      true  "User ID"
// @Param        body       body    createConversationRequest    true  "Conversation creation payload"
// @Success      201  {object}  object{conversation=object{id=string,type=string,buyer_id=string,seller_id=string,order_id=string,subject=string,status=string,created_at=string}}
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /conversations [post]
// @Security     BearerAuth
func (h *Handler) CreateConversation(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	var req createConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conversation, err := h.conversationUC.CreateConversation(c.Request.Context(), userID, usecase.CreateConversationRequest{
		Type:     req.Type,
		BuyerID:  req.BuyerID,
		SellerID: req.SellerID,
		OrderID:  req.OrderID,
		Subject:  req.Subject,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"conversation": conversation})
}

// ListConversations godoc
// @Summary      List conversations for a user
// @Tags         Chat
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string  true   "User ID"
// @Param        status     query   string  false  "Filter by status"
// @Param        page       query   int     false  "Page number"         default(1)
// @Param        page_size  query   int     false  "Items per page"      default(20)
// @Success      200  {object}  object{conversations=[]object{id=string,type=string,buyer_id=string,seller_id=string,subject=string,status=string,last_message_at=string,created_at=string},total=int,page=int,page_size=int}
// @Failure      401  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /conversations [get]
// @Security     BearerAuth
func (h *Handler) ListConversations(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	conversations, total, err := h.conversationUC.ListUserConversations(c.Request.Context(), userID, status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"conversations": conversations, "total": total, "page": page, "page_size": pageSize})
}

// GetConversation godoc
// @Summary      Get a conversation with recent messages
// @Tags         Chat
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Conversation ID"
// @Success      200  {object}  object{conversation=object{id=string,type=string,buyer_id=string,seller_id=string,subject=string,status=string,created_at=string},messages=[]object{id=string,conversation_id=string,sender_id=string,sender_role=string,content=string,message_type=string,is_read=bool,created_at=string}}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /conversations/{id} [get]
func (h *Handler) GetConversation(c *gin.Context) {
	id := c.Param("id")
	conversation, err := h.conversationUC.GetConversation(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
		return
	}

	// Fetch recent messages
	messages, _, err := h.messageUC.ListMessages(c.Request.Context(), id, 1, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"conversation": conversation, "messages": messages})
}

// ArchiveConversation godoc
// @Summary      Archive a conversation
// @Tags         Chat
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Conversation ID"
// @Success      200  {object}  object{status=string}
// @Failure      400  {object}  object{error=string}
// @Router       /conversations/{id}/archive [patch]
func (h *Handler) ArchiveConversation(c *gin.Context) {
	id := c.Param("id")
	if err := h.conversationUC.ArchiveConversation(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "archived"})
}

// --- Message Handlers ---

type sendMessageRequest struct {
	Content     string   `json:"content" binding:"required"`
	SenderRole  string   `json:"sender_role"`
	MessageType string   `json:"message_type"`
	Attachments []string `json:"attachments"`
}

// SendMessage godoc
// @Summary      Send a message in a conversation
// @Tags         Chat
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string              true  "User ID"
// @Param        id         path    string              true  "Conversation ID"
// @Param        body       body    sendMessageRequest   true  "Message payload"
// @Success      201  {object}  object{message=object{id=string,conversation_id=string,sender_id=string,sender_role=string,content=string,message_type=string,is_read=bool,created_at=string}}
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Router       /conversations/{id}/messages [post]
// @Security     BearerAuth
func (h *Handler) SendMessage(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	conversationID := c.Param("id")
	var req sendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := h.messageUC.SendMessage(c.Request.Context(), usecase.SendMessageRequest{
		ConversationID: conversationID,
		SenderID:       userID,
		SenderRole:     req.SenderRole,
		Content:        req.Content,
		MessageType:    req.MessageType,
		Attachments:    req.Attachments,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": message})
}

// ListMessages godoc
// @Summary      List messages in a conversation
// @Tags         Chat
// @Accept       json
// @Produce      json
// @Param        id         path   string  true   "Conversation ID"
// @Param        page       query  int     false  "Page number"      default(1)
// @Param        page_size  query  int     false  "Items per page"   default(50)
// @Success      200  {object}  object{messages=[]object{id=string,conversation_id=string,sender_id=string,sender_role=string,content=string,message_type=string,is_read=bool,created_at=string},total=int,page=int,page_size=int}
// @Failure      500  {object}  object{error=string}
// @Router       /conversations/{id}/messages [get]
func (h *Handler) ListMessages(c *gin.Context) {
	conversationID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	messages, total, err := h.messageUC.ListMessages(c.Request.Context(), conversationID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages, "total": total, "page": page, "page_size": pageSize})
}

// MarkAsRead godoc
// @Summary      Mark messages as read
// @Tags         Chat
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string  true  "User ID"
// @Param        id         path    string  true  "Conversation ID"
// @Success      200  {object}  object{status=string}
// @Failure      401  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /conversations/{id}/messages/read [patch]
// @Security     BearerAuth
func (h *Handler) MarkAsRead(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	conversationID := c.Param("id")
	if err := h.messageUC.MarkAsRead(c.Request.Context(), conversationID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "read"})
}

// GetUnreadCount godoc
// @Summary      Get unread message count
// @Tags         Chat
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string  true  "User ID"
// @Param        id         path    string  true  "Conversation ID"
// @Success      200  {object}  object{unread_count=int}
// @Failure      401  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /conversations/{id}/unread [get]
// @Security     BearerAuth
func (h *Handler) GetUnreadCount(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	conversationID := c.Param("id")
	count, err := h.messageUC.GetUnreadCount(c.Request.Context(), conversationID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}
