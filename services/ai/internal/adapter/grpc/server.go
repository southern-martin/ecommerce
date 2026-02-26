package grpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/southern-martin/ecommerce/services/ai/internal/domain"
	"github.com/southern-martin/ecommerce/services/ai/internal/usecase"
	"google.golang.org/grpc"
)

// AIService defines the gRPC interface for the AI service.
type AIService interface {
	GenerateEmbedding(ctx context.Context, req *GenerateEmbeddingRequest) (*EmbeddingResponse, error)
	GetRecommendations(ctx context.Context, req *GetRecommendationsRequest) (*RecommendationsResponse, error)
	GenerateDescription(ctx context.Context, req *GenerateDescriptionRequest) (*GeneratedContentResponse, error)
}

// Request/Response types

type GenerateEmbeddingRequest struct {
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	Text       string `json:"text"`
}

type EmbeddingResponse struct {
	ID           string `json:"id"`
	EntityType   string `json:"entity_type"`
	EntityID     string `json:"entity_id"`
	ModelVersion string `json:"model_version"`
	Dimensions   int    `json:"dimensions"`
}

type GetRecommendationsRequest struct {
	UserID   string `json:"user_id"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

type RecommendationsResponse struct {
	Recommendations []RecommendationItem `json:"recommendations"`
	Total           int64                `json:"total"`
}

type RecommendationItem struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	Score     float64 `json:"score"`
	Reason    string  `json:"reason"`
}

type GenerateDescriptionRequest struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Category    string `json:"category"`
}

type GeneratedContentResponse struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Model   string `json:"model"`
}

// Server implements the AIService interface.
type Server struct {
	embeddingUC      *usecase.EmbeddingUseCase
	recommendationUC *usecase.RecommendationUseCase
	contentUC        *usecase.ContentUseCase
}

var _ AIService = (*Server)(nil)

// NewServer creates a new gRPC server.
func NewServer(
	embeddingUC *usecase.EmbeddingUseCase,
	recommendationUC *usecase.RecommendationUseCase,
	contentUC *usecase.ContentUseCase,
) *Server {
	return &Server{
		embeddingUC:      embeddingUC,
		recommendationUC: recommendationUC,
		contentUC:        contentUC,
	}
}

func (s *Server) GenerateEmbedding(ctx context.Context, req *GenerateEmbeddingRequest) (*EmbeddingResponse, error) {
	embedding, err := s.embeddingUC.GenerateEmbedding(ctx, usecase.GenerateEmbeddingRequest{
		EntityType: domain.EntityType(req.EntityType),
		EntityID:   req.EntityID,
		Text:       req.Text,
	})
	if err != nil {
		return nil, err
	}
	return &EmbeddingResponse{
		ID:           embedding.ID,
		EntityType:   string(embedding.EntityType),
		EntityID:     embedding.EntityID,
		ModelVersion: embedding.ModelVersion,
		Dimensions:   embedding.Dimensions,
	}, nil
}

func (s *Server) GetRecommendations(ctx context.Context, req *GetRecommendationsRequest) (*RecommendationsResponse, error) {
	recs, total, err := s.recommendationUC.GetRecommendations(ctx, req.UserID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	items := make([]RecommendationItem, len(recs))
	for i, r := range recs {
		items[i] = RecommendationItem{
			ID:        r.ID,
			ProductID: r.ProductID,
			Score:     r.Score,
			Reason:    r.Reason,
		}
	}
	return &RecommendationsResponse{
		Recommendations: items,
		Total:           total,
	}, nil
}

func (s *Server) GenerateDescription(ctx context.Context, req *GenerateDescriptionRequest) (*GeneratedContentResponse, error) {
	content, err := s.contentUC.GenerateDescription(ctx, usecase.GenerateDescriptionRequest{
		ProductID:   req.ProductID,
		ProductName: req.ProductName,
		Category:    req.Category,
	})
	if err != nil {
		return nil, err
	}
	return &GeneratedContentResponse{
		ID:      content.ID,
		Content: content.Content,
		Model:   content.Model,
	}, nil
}

// RegisterAIServiceServer registers the AI gRPC service.
func RegisterAIServiceServer(s *grpc.Server, srv AIService) {
	s.RegisterService(&AIServiceDesc, srv)
}

func _AIService_GenerateEmbedding_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := new(GenerateEmbeddingRequest)
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(AIService).GenerateEmbedding(ctx, req)
}

func _AIService_GetRecommendations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := new(GetRecommendationsRequest)
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(AIService).GetRecommendations(ctx, req)
}

func _AIService_GenerateDescription_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := new(GenerateDescriptionRequest)
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(AIService).GenerateDescription(ctx, req)
}

// AIServiceDesc is the gRPC service descriptor.
var AIServiceDesc = grpc.ServiceDesc{
	ServiceName: "ai.AIService",
	HandlerType: (*AIService)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "GenerateEmbedding", Handler: _AIService_GenerateEmbedding_Handler},
		{MethodName: "GetRecommendations", Handler: _AIService_GetRecommendations_Handler},
		{MethodName: "GenerateDescription", Handler: _AIService_GenerateDescription_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ai/ai.proto",
}

func init() {
	// Register JSON codec for manual gRPC
	_ = json.Marshal
	_ = fmt.Sprintf
}
