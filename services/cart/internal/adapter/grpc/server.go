package grpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/southern-martin/ecommerce/services/cart/internal/domain"
	"github.com/southern-martin/ecommerce/services/cart/internal/usecase"
	"google.golang.org/grpc"
)

// --- Request / Response types for gRPC ---

// GetCartRequest is the request for GetCart RPC.
type GetCartRequest struct {
	UserID string `json:"user_id"`
}

// ClearCartRequest is the request for ClearCart RPC.
type ClearCartRequest struct {
	UserID string `json:"user_id"`
}

// CartItemResponse is a cart item in the gRPC response.
type CartItemResponse struct {
	ProductID   string `json:"product_id"`
	VariantID   string `json:"variant_id"`
	ProductName string `json:"product_name"`
	VariantName string `json:"variant_name"`
	SKU         string `json:"sku"`
	PriceCents  int64  `json:"price_cents"`
	Quantity    int32  `json:"quantity"`
	ImageURL    string `json:"image_url"`
	SellerID    string `json:"seller_id"`
}

// GetCartResponse is the response for GetCart RPC.
type GetCartResponse struct {
	UserID        string             `json:"user_id"`
	Items         []CartItemResponse `json:"items"`
	TotalItems    int32              `json:"total_items"`
	SubtotalCents int64              `json:"subtotal_cents"`
}

// ClearCartResponse is the response for ClearCart RPC.
type ClearCartResponse struct {
	Success bool `json:"success"`
}

// --- CartService interface ---

// CartService defines the gRPC service interface for inter-service cart operations.
type CartService interface {
	GetCart(ctx context.Context, req *GetCartRequest) (*GetCartResponse, error)
	ClearCart(ctx context.Context, req *ClearCartRequest) (*ClearCartResponse, error)
}

// --- Server implementation ---

// cartServiceServer implements CartService.
type cartServiceServer struct {
	cartUC *usecase.CartUseCase
	logger zerolog.Logger
}

// NewCartServiceServer creates a new gRPC cart service server.
func NewCartServiceServer(cartUC *usecase.CartUseCase, logger zerolog.Logger) CartService {
	return &cartServiceServer{
		cartUC: cartUC,
		logger: logger.With().Str("component", "grpc_cart_server").Logger(),
	}
}

func (s *cartServiceServer) GetCart(ctx context.Context, req *GetCartRequest) (*GetCartResponse, error) {
	cart, err := s.cartUC.GetCart(ctx, req.UserID)
	if err != nil {
		s.logger.Error().Err(err).Str("user_id", req.UserID).Msg("grpc: failed to get cart")
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	return toGetCartResponse(cart), nil
}

func (s *cartServiceServer) ClearCart(ctx context.Context, req *ClearCartRequest) (*ClearCartResponse, error) {
	if err := s.cartUC.ClearCart(ctx, req.UserID); err != nil {
		s.logger.Error().Err(err).Str("user_id", req.UserID).Msg("grpc: failed to clear cart")
		return nil, fmt.Errorf("failed to clear cart: %w", err)
	}

	return &ClearCartResponse{Success: true}, nil
}

func toGetCartResponse(cart *domain.Cart) *GetCartResponse {
	items := make([]CartItemResponse, len(cart.Items))
	for i, item := range cart.Items {
		items[i] = CartItemResponse{
			ProductID:   item.ProductID,
			VariantID:   item.VariantID,
			ProductName: item.ProductName,
			VariantName: item.VariantName,
			SKU:         item.SKU,
			PriceCents:  item.PriceCents,
			Quantity:    int32(item.Quantity),
			ImageURL:    item.ImageURL,
			SellerID:    item.SellerID,
		}
	}

	return &GetCartResponse{
		UserID:        cart.UserID,
		Items:         items,
		TotalItems:    int32(cart.TotalItems()),
		SubtotalCents: cart.SubtotalCents(),
	}
}

// --- Manual gRPC ServiceDesc ---

// jsonCodec is a simple JSON-based gRPC codec for manual service registration.
type jsonCodec struct{}

func (jsonCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (jsonCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (jsonCodec) Name() string {
	return "json"
}

// CartServiceDesc is the gRPC ServiceDesc for the CartService.
var CartServiceDesc = grpc.ServiceDesc{
	ServiceName: "cart.CartService",
	HandlerType: (*CartService)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCart",
			Handler:    getCartHandler,
		},
		{
			MethodName: "ClearCart",
			Handler:    clearCartHandler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cart.proto",
}

func getCartHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	req := new(GetCartRequest)
	if err := dec(req); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CartService).GetCart(ctx, req)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cart.CartService/GetCart",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CartService).GetCart(ctx, req.(*GetCartRequest))
	}
	return interceptor(ctx, req, info, handler)
}

func clearCartHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	req := new(ClearCartRequest)
	if err := dec(req); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CartService).ClearCart(ctx, req)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cart.CartService/ClearCart",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CartService).ClearCart(ctx, req.(*ClearCartRequest))
	}
	return interceptor(ctx, req, info, handler)
}

// RegisterCartServiceServer registers the CartService server with the gRPC server.
func RegisterCartServiceServer(s *grpc.Server, srv CartService) {
	s.RegisterService(&CartServiceDesc, srv)
}
