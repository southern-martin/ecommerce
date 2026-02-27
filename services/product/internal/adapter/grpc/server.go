package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/southern-martin/ecommerce/services/product/internal/usecase"
)

// ProductServiceServer defines the gRPC handler interface.
type ProductServiceServer interface {
	GetProduct(ctx context.Context, req *GetProductRequest) (*GetProductResponse, error)
	GetVariant(ctx context.Context, req *GetVariantRequest) (*GetVariantResponse, error)
	UpdateStock(ctx context.Context, req *UpdateStockRequest) (*UpdateStockResponse, error)
	ListVariantsByProduct(ctx context.Context, req *ListVariantsByProductRequest) (*ListVariantsByProductResponse, error)
}

// --- Request/Response Types ---

type GetProductRequest struct {
	ProductID string
}

type GetProductResponse struct {
	ID             string
	SellerID       string
	CategoryID     string
	Name           string
	Slug           string
	Description    string
	BasePriceCents int64
	Currency       string
	Status         string
	HasVariants    bool
}

type GetVariantRequest struct {
	VariantID string
}

type GetVariantResponse struct {
	ID         string
	ProductID  string
	SKU        string
	Name       string
	PriceCents int64
	Stock      int
	IsActive   bool
}

type UpdateStockRequest struct {
	VariantID string
	Delta     int32
}

type UpdateStockResponse struct {
	Success bool
	Message string
}

type ListVariantsByProductRequest struct {
	ProductID string
}

type VariantInfo struct {
	ID         string
	SKU        string
	Name       string
	PriceCents int64
	Stock      int
	IsActive   bool
}

type ListVariantsByProductResponse struct {
	Variants []VariantInfo
}

// --- Server Implementation ---

// Server implements the gRPC product service.
type Server struct {
	productUC *usecase.ProductUseCase
	variantUC *usecase.VariantUseCase
	server    *grpc.Server
}

// NewServer creates a new gRPC server.
func NewServer(productUC *usecase.ProductUseCase, variantUC *usecase.VariantUseCase) *Server {
	return &Server{
		productUC: productUC,
		variantUC: variantUC,
	}
}

func (s *Server) GetProduct(ctx context.Context, req *GetProductRequest) (*GetProductResponse, error) {
	if req.ProductID == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	product, err := s.productUC.GetProduct(ctx, req.ProductID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &GetProductResponse{
		ID:             product.ID,
		SellerID:       product.SellerID,
		CategoryID:     product.CategoryID,
		Name:           product.Name,
		Slug:           product.Slug,
		Description:    product.Description,
		BasePriceCents: product.BasePriceCents,
		Currency:       product.Currency,
		Status:         string(product.Status),
		HasVariants:    product.HasVariants,
	}, nil
}

func (s *Server) GetVariant(ctx context.Context, req *GetVariantRequest) (*GetVariantResponse, error) {
	if req.VariantID == "" {
		return nil, status.Error(codes.InvalidArgument, "variant_id is required")
	}

	variant, err := s.variantUC.GetVariant(ctx, req.VariantID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &GetVariantResponse{
		ID:         variant.ID,
		ProductID:  variant.ProductID,
		SKU:        variant.SKU,
		Name:       variant.Name,
		PriceCents: variant.PriceCents,
		Stock:      variant.Stock,
		IsActive:   variant.IsActive,
	}, nil
}

func (s *Server) UpdateStock(ctx context.Context, req *UpdateStockRequest) (*UpdateStockResponse, error) {
	if req.VariantID == "" {
		return nil, status.Error(codes.InvalidArgument, "variant_id is required")
	}

	if err := s.variantUC.UpdateStockDirect(ctx, req.VariantID, int(req.Delta)); err != nil {
		return &UpdateStockResponse{
			Success: false,
			Message: err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &UpdateStockResponse{
		Success: true,
		Message: "stock updated",
	}, nil
}

func (s *Server) ListVariantsByProduct(ctx context.Context, req *ListVariantsByProductRequest) (*ListVariantsByProductResponse, error) {
	if req.ProductID == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	variants, err := s.variantUC.ListVariantsByProduct(ctx, req.ProductID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var infos []VariantInfo
	for _, v := range variants {
		infos = append(infos, VariantInfo{
			ID:         v.ID,
			SKU:        v.SKU,
			Name:       v.Name,
			PriceCents: v.PriceCents,
			Stock:      v.Stock,
			IsActive:   v.IsActive,
		})
	}

	return &ListVariantsByProductResponse{Variants: infos}, nil
}

// --- gRPC ServiceDesc (manual registration, no proto codegen) ---

var _ProductService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "product.ProductService",
	HandlerType: (*ProductServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetProduct",
			Handler:    _ProductService_GetProduct_Handler,
		},
		{
			MethodName: "GetVariant",
			Handler:    _ProductService_GetVariant_Handler,
		},
		{
			MethodName: "UpdateStock",
			Handler:    _ProductService_UpdateStock_Handler,
		},
		{
			MethodName: "ListVariantsByProduct",
			Handler:    _ProductService_ListVariantsByProduct_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

func _ProductService_GetProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	req := new(GetProductRequest)
	if err := dec(req); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).GetProduct(ctx, req)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/product.ProductService/GetProduct",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).GetProduct(ctx, req.(*GetProductRequest))
	}
	return interceptor(ctx, req, info, handler)
}

func _ProductService_GetVariant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	req := new(GetVariantRequest)
	if err := dec(req); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).GetVariant(ctx, req)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/product.ProductService/GetVariant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).GetVariant(ctx, req.(*GetVariantRequest))
	}
	return interceptor(ctx, req, info, handler)
}

func _ProductService_UpdateStock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	req := new(UpdateStockRequest)
	if err := dec(req); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).UpdateStock(ctx, req)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/product.ProductService/UpdateStock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).UpdateStock(ctx, req.(*UpdateStockRequest))
	}
	return interceptor(ctx, req, info, handler)
}

func _ProductService_ListVariantsByProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	req := new(ListVariantsByProductRequest)
	if err := dec(req); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).ListVariantsByProduct(ctx, req)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/product.ProductService/ListVariantsByProduct",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).ListVariantsByProduct(ctx, req.(*ListVariantsByProductRequest))
	}
	return interceptor(ctx, req, info, handler)
}

// Start starts the gRPC server on the given port.
func (s *Server) Start(port string, opts ...grpc.ServerOption) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	s.server = grpc.NewServer(opts...)
	s.server.RegisterService(&_ProductService_serviceDesc, s)
	reflection.Register(s.server)

	log.Info().Str("port", port).Msg("gRPC server listening")
	return s.server.Serve(lis)
}

// Stop gracefully stops the gRPC server.
func (s *Server) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
	}
}
