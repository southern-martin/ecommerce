package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/southern-martin/ecommerce/services/chat/internal/usecase"
)

// Handler holds all HTTP handlers for the chat service.
type Handler struct {
	conversationUC *usecase.ConversationUseCase
	messageUC      *usecase.MessageUseCase
}

// NewHandler creates a new Handler.
func NewHandler(conversationUC *usecase.ConversationUseCase, messageUC *usecase.MessageUseCase) *Handler {
	return &Handler{
		conversationUC: conversationUC,
		messageUC:      messageUC,
	}
}

// Health returns a health check response.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "chat"})
}

// --- Conversation Handlers ---

type createConversationRequest struct {
	Type     string `json:"type"`
	BuyerID  string `json:"buyer_id" binding:"required"`
	SellerID string `json:"seller_id" binding:"required"`
	OrderID  string `json:"order_id"`
	Subject  string `json:"subject"`
}

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
