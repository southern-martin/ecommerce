package grpc

import (
	"context"
	"fmt"

	"github.com/southern-martin/ecommerce/services/affiliate/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AffiliateService defines the gRPC service interface.
type AffiliateService interface {
	TrackClick(ctx context.Context, req *TrackClickRequest) (*TrackClickResponse, error)
	TrackConversion(ctx context.Context, req *TrackConversionRequest) (*TrackConversionResponse, error)
	GetProgram(ctx context.Context, req *GetProgramRequest) (*GetProgramResponse, error)
}

// --- Request/Response types ---

type TrackClickRequest struct {
	Code string
}

type TrackClickResponse struct {
	LinkID    string
	UserID    string
	TargetURL string
}

type TrackConversionRequest struct {
	LinkID          string
	ReferredID      string
	OrderID         string
	OrderTotalCents int64
}

type TrackConversionResponse struct {
	ReferralID      string
	ReferrerID      string
	CommissionCents int64
	Status          string
}

type GetProgramRequest struct{}

type GetProgramResponse struct {
	ID                 string
	CommissionRate     float64
	MinPayoutCents     int64
	CookieDays         int32
	ReferrerBonusCents int64
	ReferredBonusCents int64
	IsActive           bool
}

// Server implements the AffiliateService gRPC interface.
type Server struct {
	linkUC     *usecase.LinkUseCase
	referralUC *usecase.ReferralUseCase
	programUC  *usecase.ProgramUseCase
}

// NewServer creates a new gRPC Server.
func NewServer(linkUC *usecase.LinkUseCase, referralUC *usecase.ReferralUseCase, programUC *usecase.ProgramUseCase) *Server {
	return &Server{
		linkUC:     linkUC,
		referralUC: referralUC,
		programUC:  programUC,
	}
}

func (s *Server) TrackClick(ctx context.Context, req *TrackClickRequest) (*TrackClickResponse, error) {
	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "code is required")
	}

	link, err := s.linkUC.TrackClick(ctx, req.Code)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &TrackClickResponse{
		LinkID:    link.ID,
		UserID:    link.UserID,
		TargetURL: link.TargetURL,
	}, nil
}

func (s *Server) TrackConversion(ctx context.Context, req *TrackConversionRequest) (*TrackConversionResponse, error) {
	if req.LinkID == "" {
		return nil, status.Error(codes.InvalidArgument, "link_id is required")
	}
	if req.OrderID == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	referral, err := s.referralUC.TrackConversion(ctx, usecase.TrackConversionRequest{
		LinkID:          req.LinkID,
		ReferredID:      req.ReferredID,
		OrderID:         req.OrderID,
		OrderTotalCents: req.OrderTotalCents,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &TrackConversionResponse{
		ReferralID:      referral.ID,
		ReferrerID:      referral.ReferrerID,
		CommissionCents: referral.CommissionCents,
		Status:          string(referral.Status),
	}, nil
}

func (s *Server) GetProgram(ctx context.Context, req *GetProgramRequest) (*GetProgramResponse, error) {
	program, err := s.programUC.GetProgram(ctx)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &GetProgramResponse{
		ID:                 program.ID,
		CommissionRate:     program.CommissionRate,
		MinPayoutCents:     program.MinPayoutCents,
		CookieDays:         int32(program.CookieDays),
		ReferrerBonusCents: program.ReferrerBonusCents,
		ReferredBonusCents: program.ReferredBonusCents,
		IsActive:           program.IsActive,
	}, nil
}

// --- gRPC ServiceDesc ---

func handlerTrackClick(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &TrackClickRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(AffiliateService).TrackClick(ctx, req)
}

func handlerTrackConversion(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &TrackConversionRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(AffiliateService).TrackConversion(ctx, req)
}

func handlerGetProgram(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &GetProgramRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(AffiliateService).GetProgram(ctx, req)
}

// AffiliateServiceDesc is the gRPC service descriptor.
var AffiliateServiceDesc = grpc.ServiceDesc{
	ServiceName: "affiliate.AffiliateService",
	HandlerType: (*AffiliateService)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "TrackClick", Handler: handlerTrackClick},
		{MethodName: "TrackConversion", Handler: handlerTrackConversion},
		{MethodName: "GetProgram", Handler: handlerGetProgram},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fmt.Sprintf("affiliate_service.proto"),
}

// RegisterAffiliateServiceServer registers the AffiliateService with a gRPC server.
func RegisterAffiliateServiceServer(s *grpc.Server, srv AffiliateService) {
	s.RegisterService(&AffiliateServiceDesc, srv)
}

// Verify Server implements AffiliateService at compile time.
var _ AffiliateService = (*Server)(nil)
