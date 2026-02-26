package grpc

import (
	"context"
	"fmt"

	"github.com/southern-martin/ecommerce/services/loyalty/internal/domain"
	"github.com/southern-martin/ecommerce/services/loyalty/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoyaltyService defines the gRPC service interface.
type LoyaltyService interface {
	GetMembership(ctx context.Context, req *GetMembershipRequest) (*GetMembershipResponse, error)
	EarnPoints(ctx context.Context, req *EarnPointsRequest) (*EarnPointsResponse, error)
	RedeemPoints(ctx context.Context, req *RedeemPointsRequest) (*RedeemPointsResponse, error)
	GetPointsBalance(ctx context.Context, req *GetPointsBalanceRequest) (*GetPointsBalanceResponse, error)
}

// --- Request/Response types ---

type GetMembershipRequest struct {
	UserID string
}

type GetMembershipResponse struct {
	UserID         string
	Tier           string
	PointsBalance  int64
	LifetimePoints int64
}

type EarnPointsRequest struct {
	UserID      string
	Points      int64
	Source      string
	ReferenceID string
	Description string
}

type EarnPointsResponse struct {
	TransactionID string
	PointsEarned  int64
	NewBalance    int64
}

type RedeemPointsRequest struct {
	UserID  string
	Points  int64
	OrderID string
}

type RedeemPointsResponse struct {
	TransactionID string
	PointsRedeemed int64
	NewBalance     int64
}

type GetPointsBalanceRequest struct {
	UserID string
}

type GetPointsBalanceResponse struct {
	UserID        string
	PointsBalance int64
}

// Server implements the LoyaltyService gRPC interface.
type Server struct {
	membershipUC *usecase.MembershipUseCase
	pointsUC     *usecase.PointsUseCase
}

// NewServer creates a new gRPC Server.
func NewServer(membershipUC *usecase.MembershipUseCase, pointsUC *usecase.PointsUseCase) *Server {
	return &Server{
		membershipUC: membershipUC,
		pointsUC:     pointsUC,
	}
}

func (s *Server) GetMembership(ctx context.Context, req *GetMembershipRequest) (*GetMembershipResponse, error) {
	if req.UserID == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	membership, err := s.membershipUC.GetMembership(ctx, req.UserID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &GetMembershipResponse{
		UserID:         membership.UserID,
		Tier:           string(membership.Tier),
		PointsBalance:  membership.PointsBalance,
		LifetimePoints: membership.LifetimePoints,
	}, nil
}

func (s *Server) EarnPoints(ctx context.Context, req *EarnPointsRequest) (*EarnPointsResponse, error) {
	if req.UserID == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.Points <= 0 {
		return nil, status.Error(codes.InvalidArgument, "points must be positive")
	}

	tx, err := s.pointsUC.EarnPoints(ctx, usecase.EarnPointsRequest{
		UserID:      req.UserID,
		Points:      req.Points,
		Source:      domain.PointsSource(req.Source),
		ReferenceID: req.ReferenceID,
		Description: req.Description,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	balance, _ := s.pointsUC.GetBalance(ctx, req.UserID)

	return &EarnPointsResponse{
		TransactionID: tx.ID,
		PointsEarned:  tx.Points,
		NewBalance:    balance,
	}, nil
}

func (s *Server) RedeemPoints(ctx context.Context, req *RedeemPointsRequest) (*RedeemPointsResponse, error) {
	if req.UserID == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.Points <= 0 {
		return nil, status.Error(codes.InvalidArgument, "points must be positive")
	}

	tx, err := s.pointsUC.RedeemPoints(ctx, usecase.RedeemPointsRequest{
		UserID:  req.UserID,
		Points:  req.Points,
		OrderID: req.OrderID,
	})
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	balance, _ := s.pointsUC.GetBalance(ctx, req.UserID)

	return &RedeemPointsResponse{
		TransactionID:  tx.ID,
		PointsRedeemed: tx.Points,
		NewBalance:     balance,
	}, nil
}

func (s *Server) GetPointsBalance(ctx context.Context, req *GetPointsBalanceRequest) (*GetPointsBalanceResponse, error) {
	if req.UserID == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	balance, err := s.pointsUC.GetBalance(ctx, req.UserID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &GetPointsBalanceResponse{
		UserID:        req.UserID,
		PointsBalance: balance,
	}, nil
}

// --- gRPC ServiceDesc ---

func handlerGetMembership(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &GetMembershipRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(LoyaltyService).GetMembership(ctx, req)
}

func handlerEarnPoints(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &EarnPointsRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(LoyaltyService).EarnPoints(ctx, req)
}

func handlerRedeemPoints(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &RedeemPointsRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(LoyaltyService).RedeemPoints(ctx, req)
}

func handlerGetPointsBalance(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &GetPointsBalanceRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(LoyaltyService).GetPointsBalance(ctx, req)
}

// LoyaltyServiceDesc is the gRPC service descriptor.
var LoyaltyServiceDesc = grpc.ServiceDesc{
	ServiceName: "loyalty.LoyaltyService",
	HandlerType: (*LoyaltyService)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "GetMembership", Handler: handlerGetMembership},
		{MethodName: "EarnPoints", Handler: handlerEarnPoints},
		{MethodName: "RedeemPoints", Handler: handlerRedeemPoints},
		{MethodName: "GetPointsBalance", Handler: handlerGetPointsBalance},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fmt.Sprintf("loyalty_service.proto"),
}

// RegisterLoyaltyServiceServer registers the LoyaltyService with a gRPC server.
func RegisterLoyaltyServiceServer(s *grpc.Server, srv LoyaltyService) {
	s.RegisterService(&LoyaltyServiceDesc, srv)
}
