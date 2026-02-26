package grpc

import (
	"context"

	apperrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/user/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserProfileResponse is the gRPC response for a user profile.
type UserProfileResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DisplayName string `json:"display_name"`
	Phone       string `json:"phone"`
	AvatarURL   string `json:"avatar_url"`
	Bio         string `json:"bio"`
	Role        string `json:"role"`
}

// SellerProfileResponse is the gRPC response for a seller profile.
type SellerProfileResponse struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	StoreName   string  `json:"store_name"`
	Description string  `json:"description"`
	LogoURL     string  `json:"logo_url"`
	Rating      float64 `json:"rating"`
	TotalSales  int     `json:"total_sales"`
	Status      string  `json:"status"`
}

// UserService is the interface required by the gRPC service descriptor.
type UserService interface {
	getProfile(ctx context.Context, req *GetProfileRequest) (*UserProfileResponse, error)
	getSellerProfile(ctx context.Context, req *GetSellerProfileRequest) (*SellerProfileResponse, error)
}

// UserServiceServer implements a simple gRPC service for inter-service calls.
type UserServiceServer struct {
	profileUC *usecase.ProfileUseCase
	sellerUC  *usecase.SellerUseCase
}

// NewUserServiceServer creates a new UserServiceServer.
func NewUserServiceServer(profileUC *usecase.ProfileUseCase, sellerUC *usecase.SellerUseCase) *UserServiceServer {
	return &UserServiceServer{
		profileUC: profileUC,
		sellerUC:  sellerUC,
	}
}

// GetProfileRequest is the request message for GetProfile.
type GetProfileRequest struct {
	UserID string `json:"user_id"`
}

// GetSellerProfileRequest is the request message for GetSellerProfile.
type GetSellerProfileRequest struct {
	SellerID string `json:"seller_id"`
}

// RegisterServer registers the user service gRPC handlers on the given grpc.Server.
func RegisterServer(s *grpc.Server, srv *UserServiceServer) {
	sd := grpc.ServiceDesc{
		ServiceName: "user.UserService",
		HandlerType: (*UserService)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "GetProfile",
				Handler:    handleGetProfile,
			},
			{
				MethodName: "GetSellerProfile",
				Handler:    handleGetSellerProfile,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "user/user.proto",
	}
	s.RegisterService(&sd, srv)
}

func handleGetProfile(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	s := srv.(*UserServiceServer)
	var req GetProfileRequest
	if err := dec(&req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to decode request: %v", err)
	}

	if interceptor == nil {
		return s.getProfile(ctx, &req)
	}

	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserService/GetProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.getProfile(ctx, req.(*GetProfileRequest))
	}
	return interceptor(ctx, &req, info, handler)
}

func handleGetSellerProfile(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	s := srv.(*UserServiceServer)
	var req GetSellerProfileRequest
	if err := dec(&req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to decode request: %v", err)
	}

	if interceptor == nil {
		return s.getSellerProfile(ctx, &req)
	}

	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserService/GetSellerProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.getSellerProfile(ctx, req.(*GetSellerProfileRequest))
	}
	return interceptor(ctx, &req, info, handler)
}

func (s *UserServiceServer) getProfile(ctx context.Context, req *GetProfileRequest) (*UserProfileResponse, error) {
	profile, err := s.profileUC.GetProfile(ctx, req.UserID)
	if err != nil {
		st := apperrors.ToGRPCStatus(err)
		return nil, st.Err()
	}

	return &UserProfileResponse{
		ID:          profile.ID,
		Email:       profile.Email,
		FirstName:   profile.FirstName,
		LastName:    profile.LastName,
		DisplayName: profile.DisplayName,
		Phone:       profile.Phone,
		AvatarURL:   profile.AvatarURL,
		Bio:         profile.Bio,
		Role:        profile.Role,
	}, nil
}

func (s *UserServiceServer) getSellerProfile(ctx context.Context, req *GetSellerProfileRequest) (*SellerProfileResponse, error) {
	seller, err := s.sellerUC.GetSeller(ctx, req.SellerID)
	if err != nil {
		st := apperrors.ToGRPCStatus(err)
		return nil, st.Err()
	}

	return &SellerProfileResponse{
		ID:          seller.ID,
		UserID:      seller.UserID,
		StoreName:   seller.StoreName,
		Description: seller.Description,
		LogoURL:     seller.LogoURL,
		Rating:      seller.Rating,
		TotalSales:  seller.TotalSales,
		Status:      seller.Status,
	}, nil
}
