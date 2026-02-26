package grpc

import (
	"context"
	"fmt"

	"github.com/southern-martin/ecommerce/services/review/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ReviewService defines the gRPC service interface.
type ReviewService interface {
	GetReview(ctx context.Context, req *GetReviewRequest) (*GetReviewResponse, error)
	GetProductSummary(ctx context.Context, req *GetProductSummaryRequest) (*GetProductSummaryResponse, error)
}

// --- Request/Response types ---

type GetReviewRequest struct {
	ReviewID string
}

type GetReviewResponse struct {
	ID        string
	ProductID string
	UserID    string
	UserName  string
	Rating    int32
	Title     string
	Content   string
	Status    string
}

type GetProductSummaryRequest struct {
	ProductID string
}

type GetProductSummaryResponse struct {
	ProductID          string
	AverageRating      float64
	TotalReviews       int32
	RatingDistribution map[int32]int32
}

// Server implements the ReviewService gRPC interface.
type Server struct {
	reviewUC *usecase.ReviewUseCase
}

// NewServer creates a new gRPC Server.
func NewServer(reviewUC *usecase.ReviewUseCase) *Server {
	return &Server{reviewUC: reviewUC}
}

func (s *Server) GetReview(ctx context.Context, req *GetReviewRequest) (*GetReviewResponse, error) {
	if req.ReviewID == "" {
		return nil, status.Error(codes.InvalidArgument, "review_id is required")
	}

	review, err := s.reviewUC.GetReview(ctx, req.ReviewID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &GetReviewResponse{
		ID:        review.ID,
		ProductID: review.ProductID,
		UserID:    review.UserID,
		UserName:  review.UserName,
		Rating:    int32(review.Rating),
		Title:     review.Title,
		Content:   review.Content,
		Status:    string(review.Status),
	}, nil
}

func (s *Server) GetProductSummary(ctx context.Context, req *GetProductSummaryRequest) (*GetProductSummaryResponse, error) {
	if req.ProductID == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	summary, err := s.reviewUC.GetProductSummary(ctx, req.ProductID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	dist := make(map[int32]int32)
	for k, v := range summary.RatingDistribution {
		dist[int32(k)] = int32(v)
	}

	return &GetProductSummaryResponse{
		ProductID:          summary.ProductID,
		AverageRating:      summary.AverageRating,
		TotalReviews:       int32(summary.TotalReviews),
		RatingDistribution: dist,
	}, nil
}

// --- gRPC ServiceDesc ---

func handlerGetReview(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &GetReviewRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(ReviewService).GetReview(ctx, req)
}

func handlerGetProductSummary(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &GetProductSummaryRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(ReviewService).GetProductSummary(ctx, req)
}

// ReviewServiceDesc is the gRPC service descriptor.
var ReviewServiceDesc = grpc.ServiceDesc{
	ServiceName: "review.ReviewService",
	HandlerType: (*ReviewService)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "GetReview", Handler: handlerGetReview},
		{MethodName: "GetProductSummary", Handler: handlerGetProductSummary},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fmt.Sprintf("review_service.proto"),
}

// RegisterReviewServiceServer registers the ReviewService with a gRPC server.
func RegisterReviewServiceServer(s *grpc.Server, srv ReviewService) {
	s.RegisterService(&ReviewServiceDesc, srv)
}

// Ensure Server implements ReviewService at compile time.
var _ ReviewService = (*Server)(nil)
