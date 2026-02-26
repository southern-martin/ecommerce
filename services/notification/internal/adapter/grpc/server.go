package grpc

import (
	"context"
	"fmt"

	"github.com/southern-martin/ecommerce/services/notification/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NotificationService defines the gRPC service interface.
type NotificationService interface {
	SendNotification(ctx context.Context, req *SendNotificationGRPCRequest) (*NotificationGRPCResponse, error)
	GetUnreadCount(ctx context.Context, req *GetUnreadCountRequest) (*UnreadCountResponse, error)
}

// --- Request/Response types ---

// SendNotificationGRPCRequest holds the data for sending a notification via gRPC.
type SendNotificationGRPCRequest struct {
	UserID  string
	Type    string
	Channel string
	Subject string
	Body    string
	Data    string
}

// NotificationGRPCResponse holds the gRPC response for a notification.
type NotificationGRPCResponse struct {
	ID      string
	UserID  string
	Type    string
	Channel string
	Subject string
	Status  string
}

// GetUnreadCountRequest holds the data for getting unread count via gRPC.
type GetUnreadCountRequest struct {
	UserID string
}

// UnreadCountResponse holds the gRPC response for unread count.
type UnreadCountResponse struct {
	Count int64
}

// Server implements the NotificationService gRPC interface.
type Server struct {
	notificationUC *usecase.NotificationUseCase
}

// NewServer creates a new gRPC Server.
func NewServer(notificationUC *usecase.NotificationUseCase) *Server {
	return &Server{
		notificationUC: notificationUC,
	}
}

func (s *Server) SendNotification(ctx context.Context, req *SendNotificationGRPCRequest) (*NotificationGRPCResponse, error) {
	if req.UserID == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.Subject == "" {
		return nil, status.Error(codes.InvalidArgument, "subject is required")
	}

	notification, err := s.notificationUC.SendNotification(ctx, usecase.SendNotificationRequest{
		UserID:  req.UserID,
		Type:    req.Type,
		Channel: req.Channel,
		Subject: req.Subject,
		Body:    req.Body,
		Data:    req.Data,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &NotificationGRPCResponse{
		ID:      notification.ID,
		UserID:  notification.UserID,
		Type:    string(notification.Type),
		Channel: string(notification.Channel),
		Subject: notification.Subject,
		Status:  string(notification.Status),
	}, nil
}

func (s *Server) GetUnreadCount(ctx context.Context, req *GetUnreadCountRequest) (*UnreadCountResponse, error) {
	if req.UserID == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	count, err := s.notificationUC.GetUnreadCount(ctx, req.UserID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &UnreadCountResponse{Count: count}, nil
}

// --- gRPC ServiceDesc ---

func handlerSendNotification(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &SendNotificationGRPCRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(NotificationService).SendNotification(ctx, req)
}

func handlerGetUnreadCount(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &GetUnreadCountRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(NotificationService).GetUnreadCount(ctx, req)
}

// NotificationServiceDesc is the gRPC service descriptor.
var NotificationServiceDesc = grpc.ServiceDesc{
	ServiceName: "notification.NotificationService",
	HandlerType: (*NotificationService)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "SendNotification", Handler: handlerSendNotification},
		{MethodName: "GetUnreadCount", Handler: handlerGetUnreadCount},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fmt.Sprintf("notification_service.proto"),
}

// RegisterNotificationServiceServer registers the NotificationService with a gRPC server.
func RegisterNotificationServiceServer(s *grpc.Server, srv NotificationService) {
	s.RegisterService(&NotificationServiceDesc, srv)
}
