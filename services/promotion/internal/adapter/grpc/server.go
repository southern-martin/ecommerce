package grpc

import (
	"context"
	"fmt"

	"github.com/southern-martin/ecommerce/services/promotion/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PromotionService defines the gRPC service interface for inter-service communication.
type PromotionService interface {
	ValidateCoupon(ctx context.Context, req *ValidateCouponRequest) (*ValidateCouponResponse, error)
	RedeemCoupon(ctx context.Context, req *RedeemCouponRequest) (*RedeemCouponResponse, error)
	GetFlashSalePrice(ctx context.Context, req *GetFlashSalePriceRequest) (*GetFlashSalePriceResponse, error)
}

// --- Request/Response types ---

// ValidateCouponRequest is the gRPC request for ValidateCoupon.
type ValidateCouponRequest struct {
	Code       string
	UserID     string
	OrderCents int64
}

// ValidateCouponResponse is the gRPC response for ValidateCoupon.
type ValidateCouponResponse struct {
	Valid         bool
	CouponID      string
	CouponCode    string
	CouponType    string
	DiscountCents int64
}

// RedeemCouponRequest is the gRPC request for RedeemCoupon.
type RedeemCouponRequest struct {
	Code       string
	UserID     string
	OrderID    string
	OrderCents int64
}

// RedeemCouponResponse is the gRPC response for RedeemCoupon.
type RedeemCouponResponse struct {
	UsageID       string
	CouponID      string
	DiscountCents int64
}

// GetFlashSalePriceRequest is the gRPC request for GetFlashSalePrice.
type GetFlashSalePriceRequest struct {
	ProductID string
	VariantID string
}

// GetFlashSalePriceResponse is the gRPC response for GetFlashSalePrice.
type GetFlashSalePriceResponse struct {
	HasFlashSale   bool
	SalePriceCents int64
	FlashSaleID    string
	FlashSaleName  string
}

// Server implements the PromotionService gRPC interface.
type Server struct {
	couponUC    *usecase.CouponUseCase
	flashSaleUC *usecase.FlashSaleUseCase
}

// NewServer creates a new gRPC Server.
func NewServer(
	couponUC *usecase.CouponUseCase,
	flashSaleUC *usecase.FlashSaleUseCase,
) *Server {
	return &Server{
		couponUC:    couponUC,
		flashSaleUC: flashSaleUC,
	}
}

// ValidateCoupon validates a coupon via gRPC.
func (s *Server) ValidateCoupon(ctx context.Context, req *ValidateCouponRequest) (*ValidateCouponResponse, error) {
	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "code is required")
	}
	if req.UserID == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	coupon, discountCents, err := s.couponUC.ValidateCoupon(ctx, usecase.ValidateCouponInput{
		Code:       req.Code,
		UserID:     req.UserID,
		OrderCents: req.OrderCents,
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &ValidateCouponResponse{
		Valid:         true,
		CouponID:      coupon.ID,
		CouponCode:    coupon.Code,
		CouponType:    string(coupon.Type),
		DiscountCents: discountCents,
	}, nil
}

// RedeemCoupon redeems a coupon via gRPC.
func (s *Server) RedeemCoupon(ctx context.Context, req *RedeemCouponRequest) (*RedeemCouponResponse, error) {
	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "code is required")
	}
	if req.UserID == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.OrderID == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	usage, err := s.couponUC.RedeemCoupon(ctx, req.Code, req.UserID, req.OrderID, req.OrderCents)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &RedeemCouponResponse{
		UsageID:       usage.ID,
		CouponID:      usage.CouponID,
		DiscountCents: usage.DiscountCents,
	}, nil
}

// GetFlashSalePrice checks if a product is in an active flash sale via gRPC.
func (s *Server) GetFlashSalePrice(ctx context.Context, req *GetFlashSalePriceRequest) (*GetFlashSalePriceResponse, error) {
	if req.ProductID == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	flashSales, err := s.flashSaleUC.ListActiveFlashSales(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	for _, fs := range flashSales {
		for _, item := range fs.Items {
			if item.ProductID == req.ProductID && (req.VariantID == "" || item.VariantID == req.VariantID) {
				if item.QuantityLimit > 0 && item.SoldCount >= item.QuantityLimit {
					continue // sold out
				}
				return &GetFlashSalePriceResponse{
					HasFlashSale:   true,
					SalePriceCents: item.SalePriceCents,
					FlashSaleID:    fs.ID,
					FlashSaleName:  fs.Name,
				}, nil
			}
		}
	}

	return &GetFlashSalePriceResponse{
		HasFlashSale: false,
	}, nil
}

// --- gRPC ServiceDesc for manual registration ---

// handlerValidateCoupon is the gRPC handler wrapper for ValidateCoupon.
func handlerValidateCoupon(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &ValidateCouponRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(PromotionService).ValidateCoupon(ctx, req)
}

// handlerRedeemCoupon is the gRPC handler wrapper for RedeemCoupon.
func handlerRedeemCoupon(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &RedeemCouponRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(PromotionService).RedeemCoupon(ctx, req)
}

// handlerGetFlashSalePrice is the gRPC handler wrapper for GetFlashSalePrice.
func handlerGetFlashSalePrice(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &GetFlashSalePriceRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(PromotionService).GetFlashSalePrice(ctx, req)
}

// PromotionServiceDesc is the gRPC service descriptor for manual registration.
var PromotionServiceDesc = grpc.ServiceDesc{
	ServiceName: "promotion.PromotionService",
	HandlerType: (*PromotionService)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ValidateCoupon",
			Handler:    handlerValidateCoupon,
		},
		{
			MethodName: "RedeemCoupon",
			Handler:    handlerRedeemCoupon,
		},
		{
			MethodName: "GetFlashSalePrice",
			Handler:    handlerGetFlashSalePrice,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fmt.Sprintf("promotion_service.proto"),
}

// RegisterPromotionServiceServer registers the PromotionService with a gRPC server.
func RegisterPromotionServiceServer(s *grpc.Server, srv PromotionService) {
	s.RegisterService(&PromotionServiceDesc, srv)
}
