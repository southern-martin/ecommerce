package aiclient

import (
	"fmt"
	"math/rand"
)

// ChatMessage represents a chat message for the AI client.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// MockAIClient provides mock AI functionality for development.
// In production, this would proxy requests to the Python FastAPI service.
type MockAIClient struct {
	baseURL string
}

// NewMockAIClient creates a new MockAIClient.
func NewMockAIClient(baseURL string) *MockAIClient {
	return &MockAIClient{baseURL: baseURL}
}

// GenerateEmbedding returns a mock 384-dimensional embedding vector.
func (c *MockAIClient) GenerateEmbedding(text string) ([]float64, error) {
	dims := 384
	vector := make([]float64, dims)
	for i := 0; i < dims; i++ {
		vector[i] = rand.Float64()*2 - 1 // range [-1, 1]
	}
	return vector, nil
}

// Chat returns a mock AI chat response.
func (c *MockAIClient) Chat(messages []ChatMessage) (string, error) {
	if len(messages) == 0 {
		return "Hello! How can I help you today?", nil
	}
	lastMsg := messages[len(messages)-1].Content
	return fmt.Sprintf("I understand you're asking about: %q. This is a mock AI response. In production, this would be processed by our AI model through the Python service at %s.", lastMsg, c.baseURL), nil
}

// GenerateDescription returns a mock AI-generated product description.
func (c *MockAIClient) GenerateDescription(productName, category string) (string, error) {
	return fmt.Sprintf(
		"Discover the %s - a premium offering in our %s collection. "+
			"Crafted with attention to detail, this product combines quality materials "+
			"with thoughtful design to deliver an exceptional experience. "+
			"Whether you're a seasoned enthusiast or new to the category, "+
			"the %s is designed to exceed your expectations. "+
			"[AI-generated mock description - production would use %s]",
		productName, category, productName, c.baseURL,
	), nil
}
