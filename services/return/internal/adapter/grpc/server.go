package grpc

import (
	"context"
	"fmt"

	"github.com/southern-martin/ecommerce/services/return/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ReturnService defines the gRPC service interface.
type ReturnService interface {
	CreateReturn(ctx context.Context, req *CreateReturnGRPCRequest) (*ReturnGRPCResponse, error)
	GetReturn(ctx context.Context, req *GetReturnRequest) (*ReturnGRPCResponse, error)
	ApproveReturn(ctx context.Context, req *ApproveReturnRequest) (*ReturnGRPCResponse, error)
	RejectReturn(ctx context.Context, req *RejectReturnRequest) (*ReturnGRPCResponse, error)
	CreateDispute(ctx context.Context, req *CreateDisputeGRPCRequest) (*DisputeGRPCResponse, error)
	ResolveDispute(ctx context.Context, req *ResolveDisputeGRPCRequest) (*DisputeGRPCResponse, error)
}

// --- Request/Response types ---

type CreateReturnGRPCRequest struct {
	OrderID  string
	BuyerID  string
	SellerID string
	Reason   string
}

type GetReturnRequest struct {
	ReturnID string
}

type ApproveReturnRequest struct {
	ReturnID          string
	SellerID          string
	RefundAmountCents int64
}

type RejectReturnRequest struct {
	ReturnID string
	SellerID string
}

type ReturnGRPCResponse struct {
	ID       string
	OrderID  string
	Status   string
	Reason   string
	BuyerID  string
	SellerID string
}

type CreateDisputeGRPCRequest struct {
	OrderID     string
	ReturnID    string
	BuyerID     string
	SellerID    string
	Type        string
	Description string
}

type ResolveDisputeGRPCRequest struct {
	DisputeID  string
	Resolution string
	ResolvedBy string
	Status     string
}

type DisputeGRPCResponse struct {
	ID         string
	OrderID    string
	Status     string
	Type       string
	Resolution string
}

// Server implements the ReturnService gRPC interface.
type Server struct {
	createReturnUC *usecase.CreateReturnUseCase
	manageReturnUC *usecase.ManageReturnUseCase
	disputeUC      *usecase.DisputeUseCase
}

// NewServer creates a new gRPC Server.
func NewServer(
	createReturnUC *usecase.CreateReturnUseCase,
	manageReturnUC *usecase.ManageReturnUseCase,
	disputeUC *usecase.DisputeUseCase,
) *Server {
	return &Server{
		createReturnUC: createReturnUC,
		manageReturnUC: manageReturnUC,
		disputeUC:      disputeUC,
	}
}

func (s *Server) CreateReturn(ctx context.Context, req *CreateReturnGRPCRequest) (*ReturnGRPCResponse, error) {
	if req.OrderID == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	ret, err := s.createReturnUC.Execute(ctx, usecase.CreateReturnRequest{
		OrderID:  req.OrderID,
		BuyerID:  req.BuyerID,
		SellerID: req.SellerID,
		Reason:   req.Reason,
		Items:    []usecase.CreateReturnItemRequest{},
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ReturnGRPCResponse{
		ID:       ret.ID,
		OrderID:  ret.OrderID,
		Status:   string(ret.Status),
		Reason:   string(ret.Reason),
		BuyerID:  ret.BuyerID,
		SellerID: ret.SellerID,
	}, nil
}

func (s *Server) GetReturn(ctx context.Context, req *GetReturnRequest) (*ReturnGRPCResponse, error) {
	if req.ReturnID == "" {
		return nil, status.Error(codes.InvalidArgument, "return_id is required")
	}

	ret, err := s.manageReturnUC.GetReturn(ctx, req.ReturnID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &ReturnGRPCResponse{
		ID:       ret.ID,
		OrderID:  ret.OrderID,
		Status:   string(ret.Status),
		Reason:   string(ret.Reason),
		BuyerID:  ret.BuyerID,
		SellerID: ret.SellerID,
	}, nil
}

func (s *Server) ApproveReturn(ctx context.Context, req *ApproveReturnRequest) (*ReturnGRPCResponse, error) {
	ret, err := s.manageReturnUC.ApproveReturn(ctx, req.ReturnID, req.SellerID, req.RefundAmountCents)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ReturnGRPCResponse{
		ID:       ret.ID,
		OrderID:  ret.OrderID,
		Status:   string(ret.Status),
		Reason:   string(ret.Reason),
		BuyerID:  ret.BuyerID,
		SellerID: ret.SellerID,
	}, nil
}

func (s *Server) RejectReturn(ctx context.Context, req *RejectReturnRequest) (*ReturnGRPCResponse, error) {
	ret, err := s.manageReturnUC.RejectReturn(ctx, req.ReturnID, req.SellerID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ReturnGRPCResponse{
		ID:       ret.ID,
		OrderID:  ret.OrderID,
		Status:   string(ret.Status),
		Reason:   string(ret.Reason),
		BuyerID:  ret.BuyerID,
		SellerID: ret.SellerID,
	}, nil
}

func (s *Server) CreateDispute(ctx context.Context, req *CreateDisputeGRPCRequest) (*DisputeGRPCResponse, error) {
	dispute, err := s.disputeUC.CreateDispute(ctx, usecase.CreateDisputeRequest{
		OrderID:     req.OrderID,
		ReturnID:    req.ReturnID,
		BuyerID:     req.BuyerID,
		SellerID:    req.SellerID,
		Type:        req.Type,
		Description: req.Description,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &DisputeGRPCResponse{
		ID:      dispute.ID,
		OrderID: dispute.OrderID,
		Status:  string(dispute.Status),
		Type:    string(dispute.Type),
	}, nil
}

func (s *Server) ResolveDispute(ctx context.Context, req *ResolveDisputeGRPCRequest) (*DisputeGRPCResponse, error) {
	dispute, err := s.disputeUC.ResolveDispute(ctx, usecase.ResolveDisputeRequest{
		DisputeID:  req.DisputeID,
		Resolution: req.Resolution,
		ResolvedBy: req.ResolvedBy,
		Status:     req.Status,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &DisputeGRPCResponse{
		ID:         dispute.ID,
		OrderID:    dispute.OrderID,
		Status:     string(dispute.Status),
		Type:       string(dispute.Type),
		Resolution: dispute.Resolution,
	}, nil
}

// --- gRPC ServiceDesc ---

func handlerCreateReturn(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &CreateReturnGRPCRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(ReturnService).CreateReturn(ctx, req)
}

func handlerGetReturn(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &GetReturnRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(ReturnService).GetReturn(ctx, req)
}

func handlerApproveReturn(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &ApproveReturnRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(ReturnService).ApproveReturn(ctx, req)
}

func handlerRejectReturn(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &RejectReturnRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(ReturnService).RejectReturn(ctx, req)
}

func handlerCreateDispute(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &CreateDisputeGRPCRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(ReturnService).CreateDispute(ctx, req)
}

func handlerResolveDispute(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &ResolveDisputeGRPCRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(ReturnService).ResolveDispute(ctx, req)
}

// ReturnServiceDesc is the gRPC service descriptor.
var ReturnServiceDesc = grpc.ServiceDesc{
	ServiceName: "return.ReturnService",
	HandlerType: (*ReturnService)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreateReturn", Handler: handlerCreateReturn},
		{MethodName: "GetReturn", Handler: handlerGetReturn},
		{MethodName: "ApproveReturn", Handler: handlerApproveReturn},
		{MethodName: "RejectReturn", Handler: handlerRejectReturn},
		{MethodName: "CreateDispute", Handler: handlerCreateDispute},
		{MethodName: "ResolveDispute", Handler: handlerResolveDispute},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fmt.Sprintf("return_service.proto"),
}

// RegisterReturnServiceServer registers the ReturnService with a gRPC server.
func RegisterReturnServiceServer(s *grpc.Server, srv ReturnService) {
	s.RegisterService(&ReturnServiceDesc, srv)
}
