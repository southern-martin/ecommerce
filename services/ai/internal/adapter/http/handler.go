package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/southern-martin/ecommerce/services/ai/internal/domain"
	"github.com/southern-martin/ecommerce/services/ai/internal/usecase"
	"gorm.io/gorm"
)

// Handler holds all HTTP handlers for the AI service.
type Handler struct {
	embeddingUC      *usecase.EmbeddingUseCase
	recommendationUC *usecase.RecommendationUseCase
	chatbotUC        *usecase.ChatbotUseCase
	contentUC        *usecase.ContentUseCase
	db               *gorm.DB
}

// NewHandler creates a new Handler.
func NewHandler(
	embeddingUC *usecase.EmbeddingUseCase,
	recommendationUC *usecase.RecommendationUseCase,
	chatbotUC *usecase.ChatbotUseCase,
	contentUC *usecase.ContentUseCase,
	db *gorm.DB,
) *Handler {
	return &Handler{
		embeddingUC:      embeddingUC,
		recommendationUC: recommendationUC,
		chatbotUC:        chatbotUC,
		contentUC:        contentUC,
		db:               db,
	}
}

// Health returns a health check response.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "ai"})
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

// --- Chat Handlers ---

type chatRequest struct {
	ConversationID string `json:"conversation_id"`
	Message        string `json:"message" binding:"required"`
}

// Chat godoc
// @Summary      Send a message to the AI assistant
// @Tags         AI
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string       true  "User ID"
// @Param        body       body    chatRequest  true  "Chat message"
// @Success      200  {object}  usecase.ChatResponse
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /ai/chat [post]
// @Security     BearerAuth
func (h *Handler) Chat(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	var req chatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.chatbotUC.Chat(c.Request.Context(), usecase.ChatRequest{
		UserID:         userID,
		ConversationID: req.ConversationID,
		Message:        req.Message,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListConversations godoc
// @Summary      List AI conversations for the current user
// @Tags         AI
// @Produce      json
// @Param        X-User-ID  header  string  true   "User ID"
// @Param        page       query   int     false  "Page number"
// @Param        page_size  query   int     false  "Page size"
// @Success      200  {object}  object{conversations=[]domain.AIConversation,total=int,page=int,page_size=int}
// @Failure      401  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /ai/chat [get]
// @Security     BearerAuth
func (h *Handler) ListConversations(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	conversations, total, err := h.chatbotUC.ListConversations(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"conversations": conversations, "total": total, "page": page, "page_size": pageSize})
}

// GetConversation godoc
// @Summary      Retrieve an AI conversation by ID
// @Tags         AI
// @Produce      json
// @Param        id  path  string  true  "Conversation ID"
// @Success      200  {object}  object{conversation=domain.AIConversation}
// @Failure      404  {object}  object{error=string}
// @Router       /ai/chat/{id} [get]
func (h *Handler) GetConversation(c *gin.Context) {
	id := c.Param("id")
	conversation, err := h.chatbotUC.GetConversation(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"conversation": conversation})
}

// --- Recommendation Handlers ---

// GetRecommendations godoc
// @Summary      Get product recommendations for the current user
// @Tags         AI
// @Produce      json
// @Param        X-User-ID  header  string  true   "User ID"
// @Param        page       query   int     false  "Page number"
// @Param        page_size  query   int     false  "Page size"
// @Success      200  {object}  object{recommendations=[]domain.Recommendation,total=int,page=int,page_size=int}
// @Failure      401  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /ai/recommendations [get]
// @Security     BearerAuth
func (h *Handler) GetRecommendations(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	recommendations, total, err := h.recommendationUC.GetRecommendations(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recommendations": recommendations, "total": total, "page": page, "page_size": pageSize})
}

// --- Content Generation Handlers ---

type generateDescriptionRequest struct {
	ProductID   string `json:"product_id" binding:"required"`
	ProductName string `json:"product_name" binding:"required"`
	Category    string `json:"category" binding:"required"`
}

// GenerateDescription godoc
// @Summary      Generate a product description using AI
// @Tags         AI
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string                      true  "User ID"
// @Param        body       body    generateDescriptionRequest  true  "Product details"
// @Success      200  {object}  object{generated_content=domain.GeneratedContent}
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /ai/generate-description [post]
// @Security     BearerAuth
func (h *Handler) GenerateDescription(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	var req generateDescriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	content, err := h.contentUC.GenerateDescription(c.Request.Context(), usecase.GenerateDescriptionRequest{
		ProductID:   req.ProductID,
		ProductName: req.ProductName,
		Category:    req.Category,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"generated_content": content})
}

// --- Image Search Handler ---

// ImageSearch godoc
// @Summary      Search products by image
// @Tags         AI
// @Produce      json
// @Success      200  {object}  object{results=[]interface{},message=string}
// @Router       /search/image [post]
func (h *Handler) ImageSearch(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"results": []interface{}{},
		"message": "Image search is a mock endpoint. In production, images would be processed by the Python AI service.",
	})
}

// --- Embedding Handlers ---

type generateEmbeddingRequest struct {
	EntityType string `json:"entity_type" binding:"required"`
	EntityID   string `json:"entity_id" binding:"required"`
	Text       string `json:"text" binding:"required"`
}

// GenerateEmbedding godoc
// @Summary      Generate an embedding for an entity
// @Tags         AI
// @Accept       json
// @Produce      json
// @Param        body  body  generateEmbeddingRequest  true  "Embedding request"
// @Success      201  {object}  object{embedding=object{id=string,entity_type=string,entity_id=string,model_version=string,dimensions=int}}
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /ai/embeddings [post]
func (h *Handler) GenerateEmbedding(c *gin.Context) {
	var req generateEmbeddingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	embedding, err := h.embeddingUC.GenerateEmbedding(c.Request.Context(), usecase.GenerateEmbeddingRequest{
		EntityType: domain.EntityType(req.EntityType),
		EntityID:   req.EntityID,
		Text:       req.Text,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"embedding": map[string]interface{}{
		"id":            embedding.ID,
		"entity_type":   embedding.EntityType,
		"entity_id":     embedding.EntityID,
		"model_version": embedding.ModelVersion,
		"dimensions":    embedding.Dimensions,
	}})
}
