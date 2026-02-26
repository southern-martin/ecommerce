package grpc

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/southern-martin/ecommerce/services/payment/internal/domain"
	"github.com/southern-martin/ecommerce/services/payment/internal/usecase"
)

// PaymentService defines the gRPC service interface for payment operations.
type PaymentService interface {
	GetPayment(ctx context.Context, req *GetPaymentRequest) (*GetPaymentResponse, error)
	ProcessRefund(ctx context.Context, req *ProcessRefundRequest) (*ProcessRefundResponse, error)
}

// GetPaymentRequest is the request for GetPayment.
type GetPaymentRequest struct {
	OrderID string
}

// GetPaymentResponse is the response for GetPayment.
type GetPaymentResponse struct {
	PaymentID   string
	OrderID     string
	BuyerID     string
	AmountCents int64
	Currency    string
	Status      string
	Method      string
}

// ProcessRefundRequest is the request for ProcessRefund.
type ProcessRefundRequest struct {
	OrderID     string
	AmountCents int64
	SellerID    string
}

// ProcessRefundResponse is the response for ProcessRefund.
type ProcessRefundResponse struct {
	Success bool
	Message string
}

// PaymentGRPCServer implements the PaymentService gRPC interface.
type PaymentGRPCServer struct {
	paymentRepo domain.PaymentRepository
	refundUC    *usecase.RefundUseCase
}

// NewPaymentGRPCServer creates a new PaymentGRPCServer.
func NewPaymentGRPCServer(paymentRepo domain.PaymentRepository, refundUC *usecase.RefundUseCase) *PaymentGRPCServer {
	return &PaymentGRPCServer{
		paymentRepo: paymentRepo,
		refundUC:    refundUC,
	}
}

// GetPayment retrieves payment details by order ID.
func (s *PaymentGRPCServer) GetPayment(ctx context.Context, req *GetPaymentRequest) (*GetPaymentResponse, error) {
	if req.OrderID == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	payment, err := s.paymentRepo.GetByOrderID(ctx, req.OrderID)
	if err != nil {
		log.Error().Err(err).Str("order_id", req.OrderID).Msg("Failed to get payment")
		return nil, status.Errorf(codes.NotFound, "payment not found for order %s", req.OrderID)
	}

	return &GetPaymentResponse{
		PaymentID:   payment.ID,
		OrderID:     payment.OrderID,
		BuyerID:     payment.BuyerID,
		AmountCents: payment.AmountCents,
		Currency:    payment.Currency,
		Status:      string(payment.Status),
		Method:      string(payment.Method),
	}, nil
}

// ProcessRefund processes a refund for a given order.
func (s *PaymentGRPCServer) ProcessRefund(ctx context.Context, req *ProcessRefundRequest) (*ProcessRefundResponse, error) {
	if req.OrderID == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	input := usecase.RefundInput{
		OrderID:     req.OrderID,
		AmountCents: req.AmountCents,
		SellerID:    req.SellerID,
	}

	if err := s.refundUC.ProcessRefund(ctx, input); err != nil {
		log.Error().Err(err).Str("order_id", req.OrderID).Msg("Failed to process refund")
		return nil, status.Errorf(codes.Internal, "failed to process refund: %v", err)
	}

	return &ProcessRefundResponse{
		Success: true,
		Message: fmt.Sprintf("Refund processed for order %s", req.OrderID),
	}, nil
}

// paymentServiceDesc is the gRPC ServiceDesc for PaymentService.
var paymentServiceDesc = grpc.ServiceDesc{
	ServiceName: "payment.PaymentService",
	HandlerType: (*PaymentService)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPayment",
			Handler:    getPaymentHandler,
		},
		{
			MethodName: "ProcessRefund",
			Handler:    processRefundHandler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

func getPaymentHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &GetPaymentRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentService).GetPayment(ctx, req)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/payment.PaymentService/GetPayment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentService).GetPayment(ctx, req.(*GetPaymentRequest))
	}
	return interceptor(ctx, req, info, handler)
}

func processRefundHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &ProcessRefundRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentService).ProcessRefund(ctx, req)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/payment.PaymentService/ProcessRefund",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentService).ProcessRefund(ctx, req.(*ProcessRefundRequest))
	}
	return interceptor(ctx, req, info, handler)
}

// RegisterPaymentService registers the PaymentService with a gRPC server.
func RegisterPaymentService(s *grpc.Server, srv PaymentService) {
	s.RegisterService(&paymentServiceDesc, srv)
}

// Ensure PaymentGRPCServer implements PaymentService.
var _ PaymentService = (*PaymentGRPCServer)(nil)
