package grpc

import (
	"context"
	"fmt"

	"github.com/southern-martin/ecommerce/services/order/internal/domain"
	"github.com/southern-martin/ecommerce/services/order/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// OrderService defines the gRPC service interface for inter-service communication.
type OrderService interface {
	GetOrder(ctx context.Context, req *GetOrderRequest) (*GetOrderResponse, error)
	UpdateOrderStatus(ctx context.Context, req *UpdateOrderStatusRequest) (*UpdateOrderStatusResponse, error)
}

// --- Request/Response types ---

// GetOrderRequest is the gRPC request for GetOrder.
type GetOrderRequest struct {
	OrderID string
}

// GetOrderResponse is the gRPC response for GetOrder.
type GetOrderResponse struct {
	ID            string
	OrderNumber   string
	BuyerID       string
	Status        string
	SubtotalCents int64
	ShippingCents int64
	TaxCents      int64
	DiscountCents int64
	TotalCents    int64
	Currency      string
	Items         []*OrderItemProto
}

// OrderItemProto represents an order item in the gRPC response.
type OrderItemProto struct {
	ID             string
	ProductID      string
	VariantID      string
	ProductName    string
	Quantity       int32
	UnitPriceCents int64
	TotalCents     int64
	SellerID       string
}

// UpdateOrderStatusRequest is the gRPC request for UpdateOrderStatus.
type UpdateOrderStatusRequest struct {
	OrderID   string
	NewStatus string
}

// UpdateOrderStatusResponse is the gRPC response for UpdateOrderStatus.
type UpdateOrderStatusResponse struct {
	ID          string
	OrderNumber string
	Status      string
}

// Server implements the OrderService gRPC interface.
type Server struct {
	getOrderUC     *usecase.GetOrderUseCase
	updateStatusUC *usecase.UpdateOrderStatusUseCase
}

// NewServer creates a new gRPC Server.
func NewServer(
	getOrderUC *usecase.GetOrderUseCase,
	updateStatusUC *usecase.UpdateOrderStatusUseCase,
) *Server {
	return &Server{
		getOrderUC:     getOrderUC,
		updateStatusUC: updateStatusUC,
	}
}

// GetOrder retrieves an order by ID via gRPC.
func (s *Server) GetOrder(ctx context.Context, req *GetOrderRequest) (*GetOrderResponse, error) {
	if req.OrderID == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	order, err := s.getOrderUC.GetOrder(ctx, req.OrderID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	resp := &GetOrderResponse{
		ID:            order.ID,
		OrderNumber:   order.OrderNumber,
		BuyerID:       order.BuyerID,
		Status:        string(order.Status),
		SubtotalCents: order.SubtotalCents,
		ShippingCents: order.ShippingCents,
		TaxCents:      order.TaxCents,
		DiscountCents: order.DiscountCents,
		TotalCents:    order.TotalCents,
		Currency:      order.Currency,
	}

	for _, item := range order.Items {
		resp.Items = append(resp.Items, &OrderItemProto{
			ID:             item.ID,
			ProductID:      item.ProductID,
			VariantID:      item.VariantID,
			ProductName:    item.ProductName,
			Quantity:       int32(item.Quantity),
			UnitPriceCents: item.UnitPriceCents,
			TotalCents:     item.TotalCents,
			SellerID:       item.SellerID,
		})
	}

	return resp, nil
}

// UpdateOrderStatus updates an order's status via gRPC (used by payment service).
func (s *Server) UpdateOrderStatus(ctx context.Context, req *UpdateOrderStatusRequest) (*UpdateOrderStatusResponse, error) {
	if req.OrderID == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}
	if req.NewStatus == "" {
		return nil, status.Error(codes.InvalidArgument, "new_status is required")
	}

	order, err := s.updateStatusUC.UpdateOrderStatus(ctx, req.OrderID, domain.OrderStatus(req.NewStatus))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &UpdateOrderStatusResponse{
		ID:          order.ID,
		OrderNumber: order.OrderNumber,
		Status:      string(order.Status),
	}, nil
}

// --- gRPC ServiceDesc for manual registration ---

// handlerGetOrder is the gRPC handler wrapper for GetOrder.
func handlerGetOrder(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &GetOrderRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(OrderService).GetOrder(ctx, req)
}

// handlerUpdateOrderStatus is the gRPC handler wrapper for UpdateOrderStatus.
func handlerUpdateOrderStatus(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &UpdateOrderStatusRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(OrderService).UpdateOrderStatus(ctx, req)
}

// OrderServiceDesc is the gRPC service descriptor for manual registration.
var OrderServiceDesc = grpc.ServiceDesc{
	ServiceName: "order.OrderService",
	HandlerType: (*OrderService)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetOrder",
			Handler:    handlerGetOrder,
		},
		{
			MethodName: "UpdateOrderStatus",
			Handler:    handlerUpdateOrderStatus,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fmt.Sprintf("order_service.proto"),
}

// RegisterOrderServiceServer registers the OrderService with a gRPC server.
func RegisterOrderServiceServer(s *grpc.Server, srv OrderService) {
	s.RegisterService(&OrderServiceDesc, srv)
}
