package grpc

import (
	"context"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pkgauth "github.com/southern-martin/ecommerce/pkg/auth"
	"github.com/southern-martin/ecommerce/services/auth/internal/usecase"
)

// ValidateTokenRequest mirrors the expected proto request for ValidateToken.
type ValidateTokenRequest struct {
	Token string
}

// ValidateTokenResponse mirrors the expected proto response for ValidateToken.
type ValidateTokenResponse struct {
	UserID string
	Email  string
	Role   string
}

// UpdateUserRoleRequest mirrors the expected proto request for UpdateUserRole.
type UpdateUserRoleRequest struct {
	UserID string
	Role   string
}

// UpdateUserRoleResponse mirrors the expected proto response for UpdateUserRole.
type UpdateUserRoleResponse struct {
	Success bool
}

// AuthService is the interface that the gRPC service descriptor requires.
type AuthService interface {
	ValidateToken(ctx context.Context, token string) (*ValidateTokenResponse, error)
	UpdateUserRole(ctx context.Context, userID, role string) (*UpdateUserRoleResponse, error)
}

// AuthServiceServer implements the gRPC auth service.
type AuthServiceServer struct {
	jwtSecret  string
	updateRole *usecase.UpdateRoleUseCase
	logger     zerolog.Logger
}

// NewAuthServiceServer creates a new AuthServiceServer.
func NewAuthServiceServer(
	jwtSecret string,
	updateRole *usecase.UpdateRoleUseCase,
	logger zerolog.Logger,
) *AuthServiceServer {
	return &AuthServiceServer{
		jwtSecret:  jwtSecret,
		updateRole: updateRole,
		logger:     logger,
	}
}

// ValidateToken validates a JWT token and returns user information.
func (s *AuthServiceServer) ValidateToken(_ context.Context, token string) (*ValidateTokenResponse, error) {
	claims, err := pkgauth.ValidateToken(token, s.jwtSecret)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	return &ValidateTokenResponse{
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
	}, nil
}

// UpdateUserRole updates a user's role.
func (s *AuthServiceServer) UpdateUserRole(ctx context.Context, userID, role string) (*UpdateUserRoleResponse, error) {
	input := usecase.UpdateRoleInput{
		UserID: userID,
		Role:   role,
	}

	if err := s.updateRole.Execute(ctx, input); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update role: %v", err)
	}

	return &UpdateUserRoleResponse{Success: true}, nil
}

// Register registers the AuthServiceServer on the given gRPC server.
func (s *AuthServiceServer) Register(srv *grpc.Server) {
	desc := grpc.ServiceDesc{
		ServiceName: "auth.AuthService",
		HandlerType: (*AuthService)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "ValidateToken",
				Handler:    validateTokenHandler,
			},
			{
				MethodName: "UpdateUserRole",
				Handler:    updateUserRoleHandler,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "auth.proto",
	}
	srv.RegisterService(&desc, s)
}

func validateTokenHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &ValidateTokenRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(*AuthServiceServer).ValidateToken(ctx, req.Token)
}

func updateUserRoleHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &UpdateUserRoleRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(*AuthServiceServer).UpdateUserRole(ctx, req.UserID, req.Role)
}
