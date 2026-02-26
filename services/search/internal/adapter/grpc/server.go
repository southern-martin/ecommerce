package grpc

import (
	"context"
	"fmt"

	"github.com/southern-martin/ecommerce/services/search/internal/domain"
	"github.com/southern-martin/ecommerce/services/search/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SearchService defines the gRPC service interface.
type SearchService interface {
	Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error)
	Suggest(ctx context.Context, req *SuggestRequest) (*SuggestResponse, error)
}

// --- Request/Response types ---

type SearchRequest struct {
	Query      string
	CategoryID string
	MinPrice   int64
	MaxPrice   int64
	InStock    *bool
	SellerID   string
	SortBy     string
	SortOrder  string
	Page       int32
	PageSize   int32
}

type SearchResponse struct {
	Results []*SearchResultProto
	Total   int64
}

type SearchResultProto struct {
	ID          string
	ProductID   string
	Name        string
	Slug        string
	Description string
	PriceCents  int64
	Currency    string
	ImageURL    string
	SellerID    string
	CategoryID  string
	Rating      float64
	ReviewCount int32
	InStock     bool
	Score       float64
}

type SuggestRequest struct {
	Query string
	Limit int32
}

type SuggestResponse struct {
	Suggestions []*SuggestionProto
}

type SuggestionProto struct {
	Text      string
	Type      string
	ProductID string
}

// Server implements the SearchService gRPC interface.
type Server struct {
	searchUC *usecase.SearchUseCase
	indexUC  *usecase.IndexUseCase
}

// NewServer creates a new gRPC Server.
func NewServer(searchUC *usecase.SearchUseCase, indexUC *usecase.IndexUseCase) *Server {
	return &Server{
		searchUC: searchUC,
		indexUC:  indexUC,
	}
}

func (s *Server) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	filter := domain.SearchFilter{
		Query:      req.Query,
		CategoryID: req.CategoryID,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
		InStock:    req.InStock,
		SellerID:   req.SellerID,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
		Page:       int(req.Page),
		PageSize:   int(req.PageSize),
	}

	results, total, err := s.searchUC.Search(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &SearchResponse{Total: total}
	for _, r := range results {
		resp.Results = append(resp.Results, &SearchResultProto{
			ID:          r.ID,
			ProductID:   r.ProductID,
			Name:        r.Name,
			Slug:        r.Slug,
			Description: r.Description,
			PriceCents:  r.PriceCents,
			Currency:    r.Currency,
			ImageURL:    r.ImageURL,
			SellerID:    r.SellerID,
			CategoryID:  r.CategoryID,
			Rating:      r.Rating,
			ReviewCount: int32(r.ReviewCount),
			InStock:     r.InStock,
			Score:       r.Score,
		})
	}

	return resp, nil
}

func (s *Server) Suggest(ctx context.Context, req *SuggestRequest) (*SuggestResponse, error) {
	if req.Query == "" {
		return nil, status.Error(codes.InvalidArgument, "query is required")
	}

	suggestions, err := s.searchUC.Suggest(ctx, req.Query, int(req.Limit))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &SuggestResponse{}
	for _, s := range suggestions {
		resp.Suggestions = append(resp.Suggestions, &SuggestionProto{
			Text:      s.Text,
			Type:      s.Type,
			ProductID: s.ProductID,
		})
	}

	return resp, nil
}

// --- gRPC ServiceDesc ---

func handlerSearch(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &SearchRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(SearchService).Search(ctx, req)
}

func handlerSuggest(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &SuggestRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(SearchService).Suggest(ctx, req)
}

// SearchServiceDesc is the gRPC service descriptor.
var SearchServiceDesc = grpc.ServiceDesc{
	ServiceName: "search.SearchService",
	HandlerType: (*SearchService)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Search", Handler: handlerSearch},
		{MethodName: "Suggest", Handler: handlerSuggest},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fmt.Sprintf("search_service.proto"),
}

// RegisterSearchServiceServer registers the SearchService with a gRPC server.
func RegisterSearchServiceServer(s *grpc.Server, srv SearchService) {
	s.RegisterService(&SearchServiceDesc, srv)
}
